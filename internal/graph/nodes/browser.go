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

// BrowserNode is responsible for executing browser tasks
type BrowserNode struct {
	name      config.AgentType
	chatModel model.ChatModel
}

func BrowserCompose() (*compose.Lambda, *BrowserNode, error) {
	n := newBrowserNode()
	lambda, err := compose.AnyLambda(n.Invoke, nil, nil, nil)
	if err != nil {
		return nil, n, err
	}
	return lambda, n, nil
}

func newBrowserNode() *BrowserNode {
	return &BrowserNode{
		name:      config.BrowserAgent,
		chatModel: llm.GetLLMClient(config.AgentLLMap[config.BrowserAgent]),
	}
}

func (n *BrowserNode) Name() string {
	return string(n.name)
}

func (n *BrowserNode) Invoke(ctx context.Context, input *schema.Message, opts ...model.Option) (output *schema.Message, err error) {
	log.GetLogger().Info("[BrowserNode] Starting browser task execution")

	browserAgent, err := agent.CreateBrowserReactAgent(ctx)
	if err != nil {
		log.GetLogger().Error("[BrowserNode] Failed to create browser agent", zap.Error(err))
		return nil, err
	}

	output, err = browserAgent.Generate(ctx, []*schema.Message{input})
	if err != nil {
		log.GetLogger().Error("[BrowserNode] Failed to execute browser agent", zap.Error(err))
		return nil, err
	}

	log.GetLogger().Info("[BrowserNode] Browser agent execution completed", zap.Any("output", output))

	return &schema.Message{
		Content: output.Content,
		Role:    schema.User,
		Name:    n.Name(),
	}, nil
}
