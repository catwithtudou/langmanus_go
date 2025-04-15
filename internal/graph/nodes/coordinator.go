package nodes

import (
	"context"

	"github.com/cloudwego/eino/components/model"
	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino/schema"
	"go.uber.org/zap"

	"catwithtudou/langmanus_go/config"
	"catwithtudou/langmanus_go/internal/llm"
	"catwithtudou/langmanus_go/internal/prompts"
	"catwithtudou/langmanus_go/log"
)

type CoordinatorNode struct {
	name      config.AgentType
	chatModel model.ChatModel
}

func CoordinatorCompose() (*compose.Lambda, error) {
	n := newCoordinatorNode()
	return compose.AnyLambda(n.Invoke, n.Stream, nil, nil)
}

func newCoordinatorNode() *CoordinatorNode {
	return &CoordinatorNode{
		name:      config.CoordinatorAgent,
		chatModel: llm.GetLLMClient(config.AgentLLMap[config.CoordinatorAgent]),
	}
}

func (n *CoordinatorNode) Invoke(ctx context.Context, input string, opts ...model.Option) (output *schema.Message, err error) {
	promptMsg := prompts.GetSystemPromptSchemaMsg(ctx, n.name, map[string]any{
		prompts.UserQueryKey: input,
	})

	result, err := n.chatModel.Generate(ctx, promptMsg, opts...)
	if err != nil {
		log.GetLogger().Error("[CoordinatorNode]failed to invoke chat model", zap.Error(err))
		return nil, err
	}

	return result, nil
}

func (n *CoordinatorNode) Stream(ctx context.Context, input string, opts ...model.Option) (output *schema.StreamReader[*schema.Message], err error) {
	promptMsg := prompts.GetSystemPromptSchemaMsg(ctx, n.name, map[string]any{
		prompts.UserQueryKey: input,
	})

	result, err := n.chatModel.Stream(ctx, promptMsg, opts...)
	if err != nil {
		log.GetLogger().Error("[CoordinatorNode]failed to stream chat model", zap.Error(err))
		return nil, err
	}

	return result, nil
}
