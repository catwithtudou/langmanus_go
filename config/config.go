package config

import (
	"os"

	"github.com/joho/godotenv"
)

// Config 定义了应用的所有配置项
type Config struct {
	// LLM配置
	ReasoningLLM LLMConfig // 用于复杂推理任务的LLM
	BasicLLM     LLMConfig // 用于简单直接任务的LLM
	VisionLLM    LLMConfig // 用于需要视觉理解的任务的LLM

	// 其他配置
	TavilyMaxResults   int
	ChromeInstancePath string
}

// LLMConfig 定义了LLM模型的基本配置
type LLMConfig struct {
	Model   string
	BaseURL string
	APIKey  string
}

type LLMType string

const (
	BasicLLM     LLMType = "basic"
	ReasoningLLM LLMType = "reasoning"
	VisionLLM    LLMType = "vision"
)

const (
	ReasoningModel   = "REASONING_MODEL"
	ReasoningBaseURL = "REASONING_BASE_URL"
	ReasoningAPIKey  = "REASONING_API_KEY"

	BasicModel   = "BASIC_MODEL"
	BasicBaseURL = "BASIC_BASE_URL"
	BasicAPIKey  = "BASIC_API_KEY"

	VLLMModel   = "VL_MODEL"
	VLLMBaseURL = "VL_BASE_URL"
	VLLMAPIKey  = "VL_API_KEY"

	ChromeInstancePath = "CHROME_INSTANCE_PATH"
)

const (
	TavilyMaxResults = 5
)

var config *Config

// LoadConfig 加载配置
func LoadConfig() *Config {
	// 加载.env文件
	_ = godotenv.Load()

	config = &Config{
		ReasoningLLM: LLMConfig{
			Model:   os.Getenv(ReasoningModel),
			BaseURL: os.Getenv(ReasoningBaseURL),
			APIKey:  os.Getenv(ReasoningAPIKey),
		},
		BasicLLM: LLMConfig{
			Model:   os.Getenv(BasicModel),
			BaseURL: os.Getenv(BasicBaseURL),
			APIKey:  os.Getenv(BasicAPIKey),
		},
		VisionLLM: LLMConfig{
			Model:   os.Getenv(VLLMModel),
			BaseURL: os.Getenv(VLLMBaseURL),
			APIKey:  os.Getenv(VLLMAPIKey),
		},
		TavilyMaxResults:   TavilyMaxResults,
		ChromeInstancePath: os.Getenv(ChromeInstancePath),
	}

	return config
}

// GetConfig 获取当前配置
func GetConfig() *Config {
	return config
}

func GetAgentLLMConfig(agent AgentType) *LLMConfig {
	llmType := AgentLLMap[agent]
	switch llmType {
	case ReasoningLLM:
		return &config.ReasoningLLM
	case BasicLLM:
		return &config.BasicLLM
	case VisionLLM:
		return &config.VisionLLM
	default:
		panic("unknown llm type: " + llmType)
	}
}
