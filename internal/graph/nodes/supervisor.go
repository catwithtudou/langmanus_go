package nodes

import (
	"context"

	"github.com/bytedance/sonic"
	"github.com/cloudwego/eino/components/model"
	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino/schema"
	"go.uber.org/zap"

	"catwithtudou/langmanus_go/config"
	"catwithtudou/langmanus_go/internal/llm"
	"catwithtudou/langmanus_go/internal/prompts"
	"catwithtudou/langmanus_go/log"
)

// SupervisorNode is responsible for coordinating the entire workflow and deciding which agent should execute next
type SupervisorNode struct {
	name      config.AgentType
	chatModel model.ChatModel
}

func SupervisorCompose() (*compose.Lambda, *SupervisorNode, error) {
	n := newSupervisorNode()
	lambda, err := compose.AnyLambda(n.Invoke, nil, nil, nil)
	if err != nil {
		return nil, n, err
	}
	return lambda, n, nil
}

func newSupervisorNode() *SupervisorNode {
	return &SupervisorNode{
		name:      config.SupervisorAgent,
		chatModel: llm.GetLLMClient(config.AgentLLMap[config.SupervisorAgent]),
	}
}

func (n *SupervisorNode) Name() string {
	return string(n.name)
}

func (n *SupervisorNode) Invoke(ctx context.Context, input *schema.Message, opts ...model.Option) (output *schema.Message, err error) {
	log.GetLogger().Info("[SupervisorNode] Evaluating next action")

	messages := prompts.GetSystemPromptSchemaMsgWithMsg(ctx, n.name, input)

	result, err := n.chatModel.Generate(ctx, messages, opts...)
	if err != nil {
		log.GetLogger().Error("[SupervisorNode] Failed to call chat model", zap.Error(err))
		return nil, err
	}

	fullResponse := cleanJSONResponse(result.Content)

	log.GetLogger().Info("[SupervisorNode] full response", zap.String("fullResponse", fullResponse))

	var router Router
	if err := sonic.UnmarshalString(fullResponse, &router); err != nil {
		log.GetLogger().Error("[SupervisorNode] Failed to parse routing decision", zap.Error(err))
		return nil, err
	}

	log.GetLogger().Info("[SupervisorNode] Delegating task to", zap.String("agent", router.Next))

	if input.Extra == nil {
		input.Extra = make(map[string]any)
	}
	input.Extra["next"] = router.Next

	return input, nil
}

func (n *SupervisorNode) Branch(ctx context.Context, in *schema.Message) (endNode string, err error) {
	next, ok := in.Extra["next"].(string)
	if !ok {
		log.GetLogger().Warn("[SupervisorNode] Planning response is empty")
		return compose.END, nil
	}

	if next == string(Finish) {
		log.GetLogger().Info("[SupervisorNode] Workflow completed")
		return compose.END, nil
	}

	if err = compose.ProcessState(ctx, func(ctx context.Context, state *State) error {
		state.next = next
		return nil
	}); err != nil {
		log.GetLogger().Error("[SupervisorNode] Failed to update state", zap.Error(err))
		return compose.END, err
	}

	return next, nil
}

func (n *SupervisorNode) BranchNodes() map[string]bool {
	return map[string]bool{
		string(config.ResearcherAgent): true,
		string(config.CoderAgent):      true,
		string(config.BrowserAgent):    true,
		string(config.ReporterAgent):   true,
		compose.END:                    true,
	}
}
