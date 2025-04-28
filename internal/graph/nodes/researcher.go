package nodes

import (
	"context"

	"catwithtudou/langmanus_go/config"
	"catwithtudou/langmanus_go/internal/agent"
	"catwithtudou/langmanus_go/internal/llm"
	"catwithtudou/langmanus_go/log"

	"github.com/cloudwego/eino/components/model"
	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino/schema"
	"go.uber.org/zap"
)

// ResearcherNode is responsible for executing research tasks
type ResearcherNode struct {
	name      config.AgentType
	chatModel model.ChatModel
}

func ResearcherCompose() (*compose.Lambda, *ResearcherNode, error) {
	n := newResearcherNode()
	lambda, err := compose.AnyLambda(n.Invoke, nil, nil, nil)
	if err != nil {
		return nil, n, err
	}
	return lambda, n, nil
}

func newResearcherNode() *ResearcherNode {
	return &ResearcherNode{
		name:      config.ResearcherAgent,
		chatModel: llm.GetLLMClient(config.AgentLLMap[config.ResearcherAgent]),
	}
}

func (n *ResearcherNode) Name() string {
	return string(n.name)
}

func (n *ResearcherNode) Invoke(ctx context.Context, input *schema.Message, opts ...model.Option) (output *schema.Message, err error) {
	log.GetLogger().Info("[ResearcherNode] Starting research task execution")

	researcherAgent, err := agent.CreateResearcherReactAgent(ctx)
	if err != nil {
		log.GetLogger().Error("[ResearcherNode] Failed to create research agent", zap.Error(err))
		return nil, err
	}

	output, err = researcherAgent.Generate(ctx, []*schema.Message{input})
	if err != nil {
		log.GetLogger().Error("[ResearcherNode] Failed to execute research agent", zap.Error(err))
		return nil, err
	}

	log.GetLogger().Info("[ResearcherNode] Research agent execution completed", zap.Any("output", output))

	return &schema.Message{
		Content: output.Content,
		Role:    schema.User,
		Name:    n.Name(),
	}, nil
}
