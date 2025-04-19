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

// SupervisorNode 监督节点，负责协调整个工作流并决定下一步应该由哪个代理执行
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
	log.GetLogger().Info("[SupervisorNode]评估下一步行动")

	messages := prompts.GetSystemPromptSchemaMsgWithMsg(ctx, n.name, input)

	result, err := n.chatModel.Generate(ctx, messages, opts...)
	if err != nil {
		log.GetLogger().Error("[SupervisorNode]调用聊天模型失败", zap.Error(err))
		return nil, err
	}

	fullResponse := cleanJSONResponse(result.Content)

	log.GetLogger().Info("[SupervisorNode] full response", zap.String("fullResponse", fullResponse))

	var router Router
	if err := sonic.UnmarshalString(fullResponse, &router); err != nil {
		log.GetLogger().Error("[SupervisorNode]解析路由决策失败", zap.Error(err))
		return nil, err
	}

	log.GetLogger().Info("[SupervisorNode]将任务委派给", zap.String("agent", router.Next))

	if input.Extra == nil {
		input.Extra = make(map[string]any)
	}
	input.Extra["next"] = router.Next

	return input, nil
}

func (n *SupervisorNode) Branch(ctx context.Context, in *schema.Message) (endNode string, err error) {
	next, ok := in.Extra["next"].(string)
	if !ok {
		log.GetLogger().Warn("[SupervisorNode]规划响应为空")
		return compose.END, nil
	}

	if next == string(Finish) {
		log.GetLogger().Info("[SupervisorNode]工作流程已完成")
		return compose.END, nil
	}

	if err = compose.ProcessState[*State](ctx, func(ctx context.Context, state *State) error {
		state.next = next
		return nil
	}); err != nil {
		log.GetLogger().Error("[SupervisorNode]更新状态失败", zap.Error(err))
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
