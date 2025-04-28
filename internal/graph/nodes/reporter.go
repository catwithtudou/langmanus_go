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

// ReporterNode is responsible for reporting the results of the workflow
type ReporterNode struct {
	name      config.AgentType
	chatModel model.ChatModel
}

func ReporterCompose() (*compose.Lambda, *ReporterNode, error) {
	n := newReporterNode()
	lambda, err := compose.AnyLambda(n.Invoke, nil, nil, nil)
	if err != nil {
		return nil, n, err
	}
	return lambda, n, nil
}

func newReporterNode() *ReporterNode {
	return &ReporterNode{
		name:      config.ReporterAgent,
		chatModel: llm.GetLLMClient(config.AgentLLMap[config.ReporterAgent]),
	}
}

func (n *ReporterNode) Name() string {
	return string(n.name)
}

func (n *ReporterNode) Invoke(ctx context.Context, input *schema.Message, opts ...model.Option) (output *schema.Message, err error) {
	log.GetLogger().Info("[ReporterNode] Reporting results")

	messages := prompts.GetSystemPromptSchemaMsgWithMsg(ctx, n.name, input)

	result, err := n.chatModel.Generate(ctx, messages, opts...)
	if err != nil {
		log.GetLogger().Error("[ReporterNode] Failed to call chat model", zap.Error(err))
		return nil, err
	}

	log.GetLogger().Info("[ReporterNode] Reporting results", zap.String("result", result.Content))

	return &schema.Message{
		Content: result.Content,
		Role:    schema.User,
		Name:    n.Name(),
	}, nil
}
