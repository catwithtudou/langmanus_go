package llm

import (
	"context"

	"github.com/cloudwego/eino-ext/components/model/ark"
	"github.com/cloudwego/eino/components/model"
	"go.uber.org/zap"

	"catwithtudou/langmanus_go/config"
	"catwithtudou/langmanus_go/log"
)

type llmClient struct {
	BasicModel     model.ChatModel
	ReasoningModel model.ChatModel
	VisionModel    model.ChatModel
}

var llmClientInstance *llmClient

func InitLLMClient(ctx context.Context) error {

	basic, err := createArkChatModel(ctx, &config.GetConfig().BasicLLM)
	if err != nil {
		return err
	}

	reasoning, err := createArkChatModel(ctx, &config.GetConfig().ReasoningLLM)
	if err != nil {
		return err
	}

	vision, err := createArkChatModel(ctx, &config.GetConfig().VisionLLM)
	if err != nil {
		return err
	}

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

func createArkChatModel(ctx context.Context, config *config.LLMConfig) (model.ChatModel, error) {
	if config == nil || config.APIKey == "" || config.Model == "" || config.BaseURL == "" {
		log.GetLogger().Info("[createArkChatModel]invalid config", zap.Any("config", config))
		return nil, nil
	}

	return ark.NewChatModel(ctx, &ark.ChatModelConfig{
		APIKey:  config.APIKey,
		Model:   config.Model,
		BaseURL: config.BaseURL,
	})
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
