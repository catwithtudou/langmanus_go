package nodes

import (
	"context"
	"strings"

	"github.com/cloudwego/eino/components/model"
	"github.com/cloudwego/eino/compose"
	"go.uber.org/zap"

	"catwithtudou/langmanus_go/config"
	"catwithtudou/langmanus_go/internal/llm"
	"catwithtudou/langmanus_go/internal/prompts"
	"catwithtudou/langmanus_go/log"
)

const (
	handoffToPlanner = "handoff_to_planner"
	coordinatorGoto  = "coordinator_goto"
)

type CoordinatorNode struct {
	name      config.AgentType
	chatModel model.ChatModel
}

func CoordinatorCompose() (*compose.Lambda, *CoordinatorNode, error) {
	n := newCoordinatorNode()
	lambda, err := compose.AnyLambda(n.Invoke, nil, nil, nil)
	if err != nil {
		return nil, n, err
	}
	return lambda, n, nil
}

func newCoordinatorNode() *CoordinatorNode {
	return &CoordinatorNode{
		name:      config.CoordinatorAgent,
		chatModel: llm.GetLLMClient(config.AgentLLMap[config.CoordinatorAgent]),
	}
}

func (n *CoordinatorNode) Name() string {
	return string(n.name)
}

func (n *CoordinatorNode) Invoke(ctx context.Context, input map[string]any, opts ...model.Option) (output map[string]any, err error) {
	userQuery, ok := input["input"].(string)
	if !ok {
		log.GetLogger().Error("[CoordinatorNode]failed to get user query from input")
		return nil, nil
	}

	promptMsg := prompts.GetSystemPromptSchemaMsgWithInput(ctx, n.name, map[string]any{
		prompts.UserQueryKey: userQuery,
	})

	result, err := n.chatModel.Generate(ctx, promptMsg, opts...)
	if err != nil {
		log.GetLogger().Error("[CoordinatorNode]failed to invoke chat model", zap.Error(err))
		return nil, err
	}

	if strings.Contains(result.Content, handoffToPlanner) {
		input[coordinatorGoto] = config.PlannerAgent
	}

	return input, nil
}

// func (n *CoordinatorNode) Stream(ctx context.Context, input map[string]any, opts ...model.Option) (output *schema.StreamReader[*schema.Message], err error) {
// 	userQuery, ok := input["input"].(string)
// 	if !ok {
// 		log.GetLogger().Error("[CoordinatorNode]failed to get user query from input")
// 		return nil, nil
// 	}
//
// 	promptMsg := prompts.GetSystemPromptSchemaMsg(ctx, n.name, map[string]any{
// 		prompts.UserQueryKey: userQuery,
// 	})
//
// 	result, err := n.chatModel.Stream(ctx, promptMsg, opts...)
// 	if err != nil {
// 		log.GetLogger().Error("[CoordinatorNode]failed to stream chat model", zap.Error(err))
// 		return nil, err
// 	}
//
// 	return result, nil
// }

func (n *CoordinatorNode) Branch(ctx context.Context, in map[string]any) (endNode string, err error) {
	coordinatorGoto, ok := in[coordinatorGoto].(config.AgentType)
	if !ok {
		return compose.END, nil
	}

	if coordinatorGoto == config.PlannerAgent {
		return string(config.PlannerAgent), nil
	}

	return compose.END, nil
}

func (n *CoordinatorNode) BranchNodes() map[string]bool {
	return map[string]bool{string(config.PlannerAgent): true, compose.END: true}
}
