package llm

import (
	"context"

	"github.com/cloudwego/eino-ext/components/model/ark"
	"github.com/cloudwego/eino/components/model"

	"catwithtudou/langmanus_go/config"
)

type llmClient struct {
	BasicModel     model.ChatModel
	ReasoningModel model.ChatModel
	VisionModel    model.ChatModel
}

var llmClientInstance *llmClient

func InitLLMClient() error {
	ctx := context.Background()

	basic, err := ark.NewChatModel(ctx, &ark.ChatModelConfig{
		APIKey:  config.GetConfig().BasicLLM.APIKey,
		Model:   config.GetConfig().BasicLLM.Model,
		BaseURL: config.GetConfig().BasicLLM.BaseURL,
	})
	if err != nil {
		return err
	}

	reasoning, err := ark.NewChatModel(ctx, &ark.ChatModelConfig{
		APIKey:  config.GetConfig().ReasoningLLM.APIKey,
		Model:   config.GetConfig().ReasoningLLM.Model,
		BaseURL: config.GetConfig().ReasoningLLM.BaseURL,
	})
	if err != nil {
		return err
	}

	vision, err := ark.NewChatModel(ctx, &ark.ChatModelConfig{
		APIKey:  config.GetConfig().VisionLLM.APIKey,
		Model:   config.GetConfig().VisionLLM.Model,
		BaseURL: config.GetConfig().VisionLLM.BaseURL,
	})
	if err != nil {
		return err
	}

	llmClientInstance = &llmClient{
		BasicModel:     basic,
		ReasoningModel: reasoning,
		VisionModel:    vision,
	}

	return nil
}

func (l *llmClient) GetBasicModel() model.ChatModel {
	return l.BasicModel
}

func (l *llmClient) GetReasoningModel() model.ChatModel {
	return l.ReasoningModel
}

func (l *llmClient) GetVisionModel() model.ChatModel {
	return l.VisionModel
}

func GetLLMClient(llmType config.LLMType) model.ChatModel {
	if llmClientInstance == nil {
		panic("llm client not initialized")
	}

	switch llmType {
	case config.BasicLLM:
		return llmClientInstance.GetBasicModel()
	case config.ReasoningLLM:
		return llmClientInstance.GetReasoningModel()
	case config.VisionLLM:
		return llmClientInstance.GetVisionModel()
	default:
		panic("unknown llm type: " + llmType)
	}
}
