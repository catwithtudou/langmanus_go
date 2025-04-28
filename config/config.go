package config

import (
	"os"

	"github.com/joho/godotenv"
)

// Config defines all configuration items for the application
type Config struct {
	// LLM Configuration
	ReasoningLLM LLMConfig // LLM for complex reasoning tasks
	BasicLLM     LLMConfig // LLM for simple direct tasks
	VisionLLM    LLMConfig // LLM for tasks requiring visual understanding

	// Other Configuration
	TavilyMaxResults   int
	ChromeInstancePath string
}

// LLMConfig defines the basic configuration for LLM models
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

// LoadConfig loads the configuration
func LoadConfig() *Config {
	// Load .env file
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

// GetConfig retrieves the current configuration
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
