package nodes

import (
	"context"

	"github.com/cloudwego/eino/components/model"
	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino/schema"
	"go.uber.org/zap"

	"catwithtudou/langmanus_go/config"
	"catwithtudou/langmanus_go/internal/agent"
	"catwithtudou/langmanus_go/internal/llm"
	"catwithtudou/langmanus_go/log"
)

// CoderNode is responsible for executing coding tasks
type CoderNode struct {
	name      config.AgentType
	chatModel model.ChatModel
}

func CoderCompose() (*compose.Lambda, *CoderNode, error) {
	n := newCoderNode()
	lambda, err := compose.AnyLambda(n.Invoke, nil, nil, nil)
	if err != nil {
		return nil, n, err
	}
	return lambda, n, nil
}

func newCoderNode() *CoderNode {
	return &CoderNode{
		name:      config.CoderAgent,
		chatModel: llm.GetLLMClient(config.AgentLLMap[config.CoderAgent]),
	}
}

func (n *CoderNode) Name() string {
	return string(n.name)
}

func (n *CoderNode) Invoke(ctx context.Context, input *schema.Message, opts ...model.Option) (output *schema.Message, err error) {
	log.GetLogger().Info("[CoderNode] Starting coding task execution")

	coderAgent, err := agent.CreateCoderReactAgent(ctx)
	if err != nil {
		log.GetLogger().Error("[CoderNode] Failed to create coder agent", zap.Error(err))
		return nil, err
	}

	output, err = coderAgent.Generate(ctx, []*schema.Message{input})
	if err != nil {
		log.GetLogger().Error("[CoderNode] Failed to execute coder agent", zap.Error(err))
		return nil, err
	}

	log.GetLogger().Info("[CoderNode] Coder agent execution completed", zap.Any("output", output))

	return &schema.Message{
		Content: output.Content,
		Role:    schema.User,
		Name:    n.Name(),
	}, nil
}
