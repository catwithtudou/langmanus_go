package agent

import (
	"context"

	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino/flow/agent/react"

	"catwithtudou/langmanus_go/config"
	"catwithtudou/langmanus_go/internal/llm"
)

func CreateResearcherReactAgent(ctx context.Context) (*react.Agent, error) {

	return react.NewAgent(ctx, &react.AgentConfig{
		Model:       llm.GetLLMClient(config.AgentLLMap[config.ResearcherAgent]),
		ToolsConfig: compose.ToolsNodeConfig{},
	})
}
