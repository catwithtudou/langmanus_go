package nodes

import (
	"context"
	"encoding/json"
	"strings"

	"github.com/cloudwego/eino-ext/components/tool/duckduckgo"
	"github.com/cloudwego/eino/components/model"
	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino/schema"
	"go.uber.org/zap"

	"catwithtudou/langmanus_go/config"
	"catwithtudou/langmanus_go/internal/llm"
	"catwithtudou/langmanus_go/internal/prompts"
	"catwithtudou/langmanus_go/internal/tools"
	"catwithtudou/langmanus_go/log"
)

// PlannerNode is responsible for generating the overall execution plan
type PlannerNode struct {
	name           config.AgentType
	chatModel      model.ChatModel
	reasoningModel model.ChatModel
}

func PlannerCompose() (*compose.Lambda, *PlannerNode, error) {
	n := newPlannerNode()
	lambda, err := compose.AnyLambda(n.Invoke, nil, nil, nil)
	if err != nil {
		return nil, n, err
	}
	return lambda, n, nil
}

func newPlannerNode() *PlannerNode {
	return &PlannerNode{
		name:           config.PlannerAgent,
		chatModel:      llm.GetLLMClient(config.BasicLLM),
		reasoningModel: llm.GetLLMClient(config.ReasoningLLM),
	}
}

func (n *PlannerNode) Name() string {
	return string(n.name)
}

func (n *PlannerNode) Invoke(ctx context.Context, input map[string]any, opts ...model.Option) (output *schema.Message, err error) {
	log.GetLogger().Info("[PlannerNode] Generating complete execution plan")

	// Receive input from coordinator node
	userQuery, ok := input["input"].(string)
	if !ok {
		log.GetLogger().Error("[PlannerNode] Failed to get user query from input")
		return nil, nil
	}

	// Check if search is needed before planning
	searchResults := n.checkSearchBeforePlanning(ctx, input, userQuery)
	if len(searchResults) > 0 {
		searchResultsJson, err := json.Marshal(searchResults)
		if err == nil {
			userQuery += "\n\n# Relative Search Results\n\n" + string(searchResultsJson)
		}
	}

	// Build prompt message
	promptMsg := prompts.GetSystemPromptSchemaMsgWithInput(ctx, n.name, userQuery)

	// Call model to generate plan
	streamReader, err := n.getTargetModel(input).Stream(ctx, promptMsg, opts...)
	if err != nil {
		log.GetLogger().Error("[PlannerNode] Failed to call chat model", zap.Error(err))
		return nil, err
	}

	fullResponse := ""

	for {
		chunk, err := streamReader.Recv()
		if err != nil {
			log.GetLogger().Error("[PlannerNode]stream.Recv() failed, err=%v", zap.Error(err))
			break
		}

		fullResponse += chunk.Content
	}
	log.GetLogger().Info("[PlannerNode] full response", zap.String("fullResponse", fullResponse))

	fullResponse = cleanJSONResponse(fullResponse)

	log.GetLogger().Debug("[PlannerNode] full response clean json response", zap.String("fullResponse", fullResponse))

	return &schema.Message{
		Content: fullResponse,
		Role:    schema.User,
		Name:    n.Name(),
	}, nil
}

func (n *PlannerNode) Branch(ctx context.Context, in *schema.Message) (endNode string, err error) {
	fullResponse := in.Content
	if fullResponse == "" {
		log.GetLogger().Warn("[PlannerNode] Planning response is empty")
		return compose.END, nil
	}

	if !isValidJSON(fullResponse) {
		log.GetLogger().Warn("[PlannerNode] Planning response is not valid JSON")
		return compose.END, nil
	}

	if err = compose.ProcessState(ctx, func(ctx context.Context, state *State) error {
		state.fullPlan = fullResponse
		return nil
	}); err != nil {
		log.GetLogger().Error("[PlannerNode] Failed to update state", zap.Error(err))
		return compose.END, err
	}

	return string(config.SupervisorAgent), nil
}

func (n *PlannerNode) BranchNodes() map[string]bool {
	return map[string]bool{string(config.SupervisorAgent): true, compose.END: true}
}

// Check if deep thinking mode is enabled and return the corresponding model
func (n *PlannerNode) getTargetModel(input map[string]any) model.ChatModel {
	if isDeepThink, _ := input["deep_thinking_mode"].(bool); isDeepThink {
		log.GetLogger().Info("[PlannerNode] Using deep thinking mode")
		return n.reasoningModel
	}
	return n.chatModel
}

// Check if search is needed before planning
func (n *PlannerNode) checkSearchBeforePlanning(ctx context.Context, input map[string]any, userQuery string) []*duckduckgo.SearchResult {
	searchBeforePlanning, _ := input["search_before_planning"].(bool)

	if !searchBeforePlanning {
		return nil
	}

	searchReq := &duckduckgo.SearchRequest{
		Query: userQuery,
		Page:  1,
	}

	searchReqJson, err := json.Marshal(searchReq)
	if err != nil {
		log.GetLogger().Error("[PlannerNode] Failed to serialize search request", zap.Error(err))
		return nil
	}

	searchResults, err := tools.GetDuckSearchTool().InvokableRun(ctx, string(searchReqJson))
	if err != nil {
		log.GetLogger().Error("[PlannerNode] Search failed", zap.Error(err))
		return nil
	}

	var searchResp duckduckgo.SearchResponse
	if err := json.Unmarshal([]byte(searchResults), &searchResp); err != nil {
		log.GetLogger().Error("[PlannerNode] Failed to deserialize search response", zap.Error(err))
		return nil
	}

	if len(searchResp.Results) == 0 {
		return nil
	}

	return searchResp.Results
}

// cleanJSONResponse cleans the JSON format in the response
func cleanJSONResponse(response string) string {
	// Remove JSON code block markers
	response = strings.TrimPrefix(response, "```json")
	response = strings.TrimPrefix(response, "```")
	response = strings.TrimSuffix(response, "```")
	return strings.TrimSpace(response)
}

// isValidJSON verifies if the string is valid JSON
func isValidJSON(str string) bool {
	var js json.RawMessage
	return json.Unmarshal([]byte(str), &js) == nil
}
