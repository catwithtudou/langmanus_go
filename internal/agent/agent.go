package agent

import (
	"context"

	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino/flow/agent/react"
	"github.com/cloudwego/eino/schema"

	"catwithtudou/langmanus_go/config"
	"catwithtudou/langmanus_go/internal/llm"
	"catwithtudou/langmanus_go/internal/prompts"
	"catwithtudou/langmanus_go/internal/tools"
)

func CreateResearcherReactAgent(ctx context.Context) (*react.Agent, error) {
	return react.NewAgent(ctx, &react.AgentConfig{
		Model: llm.GetLLMClient(config.AgentLLMap[config.ResearcherAgent]),
		ToolsConfig: compose.ToolsNodeConfig{
			Tools: []tool.BaseTool{
				tools.GetDuckSearchTool(),
				// TODO: 实现 crawer tools
			},
		},
		MessageModifier: func(ctx context.Context, input []*schema.Message) []*schema.Message {
			res := make([]*schema.Message, 0, len(input)+1)
			res = append(res, schema.SystemMessage(prompts.GetSystemPrompt(config.ResearcherAgent)))
			res = append(res, input...)
			return res
		},
	})
}
