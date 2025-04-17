package nodes

import (
	"context"
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

const (
	handoffToPlanner = "handoff_to_planner"
)

type CoordinatorNode struct {
	name      config.AgentType
	chatModel model.ChatModel
}

func CoordinatorCompose() (*compose.Lambda, *CoordinatorNode, error) {
	n := newCoordinatorNode()
	lambda, err := compose.AnyLambda(n.Invoke, n.Stream, nil, nil)
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

func (n *CoordinatorNode) Invoke(ctx context.Context, input map[string]any, opts ...model.Option) (output *schema.Message, err error) {
	userQuery, ok := input["input"].(string)
	if !ok {
		log.GetLogger().Error("[CoordinatorNode]failed to get user query from input")
		return nil, nil
	}

	promptMsg := prompts.GetSystemPromptSchemaMsg(ctx, n.name, map[string]any{
		prompts.UserQueryKey: userQuery,
	})

	result, err := n.chatModel.Generate(ctx, promptMsg, opts...)
	if err != nil {
		log.GetLogger().Error("[CoordinatorNode]failed to invoke chat model", zap.Error(err))
		return nil, err
	}

	return result, nil
}

func (n *CoordinatorNode) Stream(ctx context.Context, input map[string]any, opts ...model.Option) (output *schema.StreamReader[*schema.Message], err error) {
	userQuery, ok := input["input"].(string)
	if !ok {
		log.GetLogger().Error("[CoordinatorNode]failed to get user query from input")
		return nil, nil
	}

	promptMsg := prompts.GetSystemPromptSchemaMsg(ctx, n.name, map[string]any{
		prompts.UserQueryKey: userQuery,
	})

	result, err := n.chatModel.Stream(ctx, promptMsg, opts...)
	if err != nil {
		log.GetLogger().Error("[CoordinatorNode]failed to stream chat model", zap.Error(err))
		return nil, err
	}

	return result, nil
}

func (n *CoordinatorNode) Branch(ctx context.Context, in *schema.Message) (endNode string, err error) {
	if strings.Contains(in.Content, handoffToPlanner) {
		return string(config.PlannerAgent), nil
	}

	return compose.END, nil
}

func (n *CoordinatorNode) BranchNodes() map[string]bool {
	return map[string]bool{string(config.PlannerAgent): true, compose.END: true}
}
