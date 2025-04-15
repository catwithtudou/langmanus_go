package prompts

import (
	"context"
	"time"

	"github.com/cloudwego/eino/components/prompt"
	"github.com/cloudwego/eino/schema"

	"catwithtudou/langmanus_go/config"
)

const (
	UserQueryKey   = "user_query"
	CurrentTimeKey = "CURRENT_TIME"
)

func GetSystemPromptSchemaMsg(ctx context.Context, agent config.AgentType, input map[string]any) []*schema.Message {

	template := prompt.FromMessages(schema.FString,
		schema.SystemMessage(getSystemPrompt(agent)),
		schema.UserMessage(UserQueryKey))

	messages, _ := template.Format(ctx, map[string]any{
		CurrentTimeKey: time.Now().Format("2006-01-02 15:04:05"),
		UserQueryKey:   input[UserQueryKey],
	})

	return messages
}

func getSystemPrompt(agent config.AgentType) string {
	switch agent {
	case config.CoordinatorAgent:
		return ""
	}
	return ""
}
