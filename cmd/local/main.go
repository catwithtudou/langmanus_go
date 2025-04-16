package main

import (
	"context"

	"go.uber.org/zap"

	"catwithtudou/langmanus_go/config"
	"catwithtudou/langmanus_go/internal/llm"
	internal "catwithtudou/langmanus_go/internal/service"
	"catwithtudou/langmanus_go/log"
)

func main() {
	ctx := context.Background()

	log.InitLogger()
	defer log.GetLogger().Sync()
	config.LoadConfig()
	llm.InitLLMClient(ctx)

	output, err := internal.RunAgentWorkflow(ctx, "你好")
	if err != nil {
		log.GetLogger().Error("[main]failed to run agent workflow", zap.Error(err))
		return
	}

	// 记录完整输出（包含元数据）
	log.GetLogger().Debug("[main]agent workflow complete output", zap.Any("output", output))
	// 仅记录模型回复内容
	log.GetLogger().Info("[main]agent workflow response", zap.String("content", output.Content))
}
