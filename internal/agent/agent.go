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
				tools.GetCrawHtmlTool(),
			},
		},
		MessageModifier: getMessageModifier(config.ResearcherAgent),
	})
}

func CreateCoderReactAgent(ctx context.Context) (*react.Agent, error) {
	return react.NewAgent(ctx, &react.AgentConfig{
		Model: llm.GetLLMClient(config.AgentLLMap[config.CoderAgent]),
		ToolsConfig: compose.ToolsNodeConfig{
			Tools: []tool.BaseTool{
				tools.GetBashTool(),
				// TODO: add tools
			},
		},
		MessageModifier: getMessageModifier(config.CoderAgent),
	})
}

func CreateBrowserReactAgent(ctx context.Context) (*react.Agent, error) {
	return react.NewAgent(ctx, &react.AgentConfig{
		Model: llm.GetLLMClient(config.AgentLLMap[config.BrowserAgent]),
		ToolsConfig: compose.ToolsNodeConfig{
			Tools: []tool.BaseTool{
				// TODO: add tools
			},
		},
		MessageModifier: getMessageModifier(config.BrowserAgent),
	})
}

func getMessageModifier(agentType config.AgentType) func(ctx context.Context, input []*schema.Message) []*schema.Message {
	return func(ctx context.Context, input []*schema.Message) []*schema.Message {
		res := make([]*schema.Message, 0, len(input)+1)
		res = append(res, schema.SystemMessage(prompts.GetSystemPrompt(agentType)))
		res = append(res, input...)
		return res
	}
}
