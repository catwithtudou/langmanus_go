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
		ReasoningLLM:       LLMConfig{},
		BasicLLM:           LLMConfig{},
		VisionLLM:          LLMConfig{},
		TavilyMaxResults:   TavilyMaxResults,
		ChromeInstancePath: "",
	}

	// 只覆盖环境变量中存在的配置项
	if model := os.Getenv(ReasoningModel); model != "" {
		config.ReasoningLLM.Model = model
	}
	if baseURL := os.Getenv(ReasoningBaseURL); baseURL != "" {
		config.ReasoningLLM.BaseURL = baseURL
	}
	if apiKey := os.Getenv(ReasoningAPIKey); apiKey != "" {
		config.ReasoningLLM.APIKey = apiKey
	}

	if model := os.Getenv(BasicModel); model != "" {
		config.BasicLLM.Model = model
	}
	if baseURL := os.Getenv(BasicBaseURL); baseURL != "" {
		config.BasicLLM.BaseURL = baseURL
	}
	if apiKey := os.Getenv(BasicAPIKey); apiKey != "" {
		config.BasicLLM.APIKey = apiKey
	}

	if model := os.Getenv(VLLMModel); model != "" {
		config.VisionLLM.Model = model
	}
	if baseURL := os.Getenv(VLLMBaseURL); baseURL != "" {
		config.VisionLLM.BaseURL = baseURL
	}
	if apiKey := os.Getenv(VLLMAPIKey); apiKey != "" {
		config.VisionLLM.APIKey = apiKey
	}

	if chromeInstancePath := os.Getenv(ChromeInstancePath); chromeInstancePath != "" {
		config.ChromeInstancePath = chromeInstancePath
	}

	return config
}

// GetConfig 获取当前配置
func GetConfig() *Config {
	return config
}
