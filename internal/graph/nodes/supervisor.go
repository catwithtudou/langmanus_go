package nodes

import (
	"context"
	"encoding/json"
	"errors"
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

// Router 用于定义下一步路由的结构体
type Router struct {
	Next string `json:"next"`
}

// SupervisorNode 监督节点，负责协调整个工作流并决定下一步应该由哪个代理执行
type SupervisorNode struct {
	name      config.AgentType
	chatModel model.ChatModel
}

func SupervisorCompose() (*compose.Lambda, *SupervisorNode, error) {
	n := newSupervisorNode()
	lambda, err := compose.AnyLambda(n.Invoke, n.Stream, nil, nil)
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

func (n *SupervisorNode) Invoke(ctx context.Context, input map[string]any, opts ...model.Option) (output *schema.Message, err error) {
	log.GetLogger().Info("[SupervisorNode]评估下一步行动")

	// 从planner节点获取输入
	planResponse, ok := input["input"].(string)
	if !ok {
		log.GetLogger().Error("[SupervisorNode]无法从输入获取计划响应")
		return nil, errors.New("无法从输入获取计划响应")
	}

	// 构建提示消息
	promptMsg := prompts.GetSystemPromptSchemaMsg(ctx, n.name, map[string]any{
		"plan_response": planResponse,
	})

	// 调用模型生成路由决策
	result, err := n.chatModel.Generate(ctx, promptMsg, opts...)
	if err != nil {
		log.GetLogger().Error("[SupervisorNode]调用聊天模型失败", zap.Error(err))
		return nil, err
	}

	// 处理响应结果
	fullResponse := result.Content

	// 清理JSON格式
	fullResponse = cleanJSONResponse(fullResponse)

	// 解析路由决策
	var router Router
	if err := json.Unmarshal([]byte(fullResponse), &router); err != nil {
		log.GetLogger().Error("[SupervisorNode]解析路由决策失败", zap.Error(err))
		return nil, err
	}

	// 检查是否结束工作流
	if strings.ToUpper(router.Next) == "FINISH" {
		log.GetLogger().Info("[SupervisorNode]工作流程已完成")
		router.Next = "__end__"
	} else {
		log.GetLogger().Info("[SupervisorNode]将任务委派给", zap.String("agent", router.Next))
	}

	// 创建包含路由信息的响应
	routerJSON, _ := json.Marshal(map[string]string{
		"content": fullResponse,
		"next":    router.Next,
	})

	// 更新结果内容，包含路由信息
	result.Content = string(routerJSON)

	return result, nil
}

func (n *SupervisorNode) Stream(ctx context.Context, input map[string]any, opts ...model.Option) (output *schema.StreamReader[*schema.Message], err error) {
	log.GetLogger().Info("[SupervisorNode]流式评估下一步行动")

	// 从planner节点获取输入
	planResponse, ok := input["input"].(string)
	if !ok {
		log.GetLogger().Error("[SupervisorNode]无法从输入获取计划响应")
		return nil, errors.New("无法从输入获取计划响应")
	}

	// 构建提示消息
	promptMsg := prompts.GetSystemPromptSchemaMsg(ctx, n.name, map[string]any{
		"plan_response": planResponse,
	})

	// 调用模型流式生成路由决策
	streamReader, err := n.chatModel.Stream(ctx, promptMsg, opts...)
	if err != nil {
		log.GetLogger().Error("[SupervisorNode]流式调用聊天模型失败", zap.Error(err))
		return nil, err
	}

	// 返回流式响应
	return streamReader, nil
}
