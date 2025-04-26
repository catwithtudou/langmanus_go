package nodes

import (
	"catwithtudou/langmanus_go/config"
	"catwithtudou/langmanus_go/internal/agent"
	"catwithtudou/langmanus_go/internal/llm"
	"catwithtudou/langmanus_go/log"
	"context"

	"github.com/cloudwego/eino/components/model"
	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino/schema"
	"go.uber.org/zap"
)

// ResearcherNode     研究代理节点,负责执行研究任务
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
	log.GetLogger().Info("[ResearcherNode]开始执行研究任务")

	researcherAgent, err := agent.CreateResearcherReactAgent(ctx)
	if err != nil {
		log.GetLogger().Error("[ResearcherNode]创建研究代理失败", zap.Error(err))
		return nil, err
	}

	output, err = researcherAgent.Generate(ctx, []*schema.Message{input})
	if err != nil {
		log.GetLogger().Error("[ResearcherNode]执行研究代理失败", zap.Error(err))
		return nil, err
	}

	log.GetLogger().Info("[ResearcherNode]研究代理执行完成", zap.Any("output", output))

	return &schema.Message{
		Content: output.Content,
		Role:    schema.User,
		Name:    n.Name(),
	}, nil
}
