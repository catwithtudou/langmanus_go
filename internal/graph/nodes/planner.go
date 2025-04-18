package nodes

import (
	"context"
	"encoding/json"
	"strings"

	"github.com/cloudwego/eino/components/model"
	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino/schema"
	"go.uber.org/zap"

	"catwithtudou/langmanus_go/config"
	"catwithtudou/langmanus_go/internal/llm"
	"catwithtudou/langmanus_go/internal/prompts"
	"catwithtudou/langmanus_go/log"
)

// PlannerNode 规划节点，负责生成整体执行计划
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
	log.GetLogger().Info("[PlannerNode]生成完整执行计划")

	// 从coordinator节点接收输入
	userQuery, ok := input["input"].(string)
	if !ok {
		log.GetLogger().Error("[PlannerNode]无法从输入获取用户查询")
		return nil, nil
	}

	// 构建提示消息
	promptMsg := prompts.GetSystemPromptSchemaMsg(ctx, n.name, map[string]any{
		prompts.UserQueryKey: userQuery,
	})

	// 检查是否需要在规划前进行搜索
	n.checkSearchBeforePlanning(input, userQuery)

	// 调用模型生成计划
	streamReader, err := n.getTargetModel(input).Stream(ctx, promptMsg, opts...)
	if err != nil {
		log.GetLogger().Error("[PlannerNode]调用聊天模型失败", zap.Error(err))
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

		log.GetLogger().Debug("[PlannerNode]stream.Recv()", zap.Any("chunk", chunk))
	}

	fullResponse = cleanJSONResponse(fullResponse)

	return &schema.Message{
		Content: fullResponse,
		Role:    schema.User,
		Name:    n.Name(),
		Extra: map[string]any{
			"full_response": fullResponse,
		},
	}, nil
}

// func (n *PlannerNode) Stream(ctx context.Context, input map[string]any, opts ...model.Option) (output *schema.StreamReader[*schema.Message], err error) {
// 	log.GetLogger().Info("[PlannerNode]流式生成执行计划")
//
// 	// 从coordinator节点接收输入
// 	userQuery, ok := input["input"].(string)
// 	if !ok {
// 		log.GetLogger().Error("[PlannerNode]无法从输入获取用户查询")
// 		return nil, nil
// 	}
//
// 	// 构建提示消息
// 	promptMsg := prompts.GetSystemPromptSchemaMsg(ctx, n.name, map[string]any{
// 		prompts.UserQueryKey: userQuery,
// 	})
//
// 	// 检查是否需要在规划前进行搜索
// 	n.checkSearchBeforePlanning(input, userQuery)
//
// 	// 调用模型流式生成计划
// 	streamReader, err := n.getTargetModel(input).Stream(ctx, promptMsg, opts...)
// 	if err != nil {
// 		log.GetLogger().Error("[PlannerNode]流式调用聊天模型失败", zap.Error(err))
// 		return nil, err
// 	}
//
// 	// 返回流式响应
// 	return streamReader, nil
// }

func (n *PlannerNode) Branch(ctx context.Context, in *schema.Message) (endNode string, err error) {
	fullResponse := in.Content
	if fullResponse == "" {
		log.GetLogger().Warn("[PlannerNode]规划响应为空")
		return compose.END, nil
	}

	if !isValidJSON(fullResponse) {
		log.GetLogger().Warn("[PlannerNode]规划响应不是有效的JSON")
		return compose.END, nil
	}

	return string(config.SupervisorAgent), nil
}

func (n *PlannerNode) BranchNodes() map[string]bool {
	return map[string]bool{string(config.SupervisorAgent): true, compose.END: true}
}

// 检测是否启用深度思考模式，并返回对应的模型
func (n *PlannerNode) getTargetModel(input map[string]any) model.ChatModel {
	if isDeepThink, _ := input["deep_thinking_mode"].(bool); isDeepThink {
		log.GetLogger().Info("[PlannerNode]使用深度思考模式")
		return n.reasoningModel
	}
	return n.chatModel
}

// 检测是否需要在规划前进行搜索
func (n *PlannerNode) checkSearchBeforePlanning(input map[string]any, userQuery string) {
	searchBeforePlanning, _ := input["search_before_planning"].(bool)

	if !searchBeforePlanning {
		return
	}

	// 这里应实现搜索逻辑，类似Python中的tavily_tool
	// searchResults := search.Invoke(userQuery)
	// 将搜索结果添加到提示中

	// TODO: 实现搜索功能

	log.GetLogger().Info("[PlannerNode]搜索后规划功能待实现")
}

// cleanJSONResponse 清理响应中的JSON格式
func cleanJSONResponse(response string) string {
	// 移除JSON代码块标记
	response = strings.TrimPrefix(response, "```json")
	response = strings.TrimPrefix(response, "```")
	response = strings.TrimSuffix(response, "```")
	return strings.TrimSpace(response)
}

// isValidJSON 验证字符串是否为有效的JSON
func isValidJSON(str string) bool {
	var js json.RawMessage
	return json.Unmarshal([]byte(str), &js) == nil
}
