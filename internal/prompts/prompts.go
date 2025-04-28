package prompts

import (
	"context"
	"os"
	"path/filepath"
	"regexp"
	"time"

	"github.com/cloudwego/eino/components/prompt"
	"github.com/cloudwego/eino/schema"
	"go.uber.org/zap"

	"catwithtudou/langmanus_go/config"
	"catwithtudou/langmanus_go/log"
)

const (
	CurrentTimeKey = "CURRENT_TIME"
	TeamMembersKey = "TEAM_MEMBERS"
)

func GetSystemPromptSchemaMsgWithInput(ctx context.Context, agent config.AgentType, userQuery string) []*schema.Message {
	params := make(map[string]any)
	params[CurrentTimeKey] = time.Now().Format("2006-01-02 15:04:05")
	params[TeamMembersKey] = config.TeamMembers

	template := prompt.FromMessages(schema.FString,
		schema.SystemMessage(GetSystemPrompt(agent)),
		schema.UserMessage(userQuery))

	messages, err := template.Format(ctx, params)
	if err != nil {
		log.GetLogger().Error("[GetSystemPromptSchemaMsg]failed to format template",
			zap.Error(err))
	}

	return messages
}

func GetSystemPromptSchemaMsgWithMsg(ctx context.Context, agent config.AgentType, schemaMsg *schema.Message) []*schema.Message {
	params := make(map[string]any)
	params[CurrentTimeKey] = time.Now().Format("2006-01-02 15:04:05")
	params[TeamMembersKey] = config.TeamMembers

	template := prompt.FromMessages(schema.FString,
		schema.SystemMessage(GetSystemPrompt(agent)),
		schemaMsg)

	messages, err := template.Format(ctx, params)
	if err != nil {
		log.GetLogger().Error("[GetSystemPromptSchemaMsgWithMsg]failed to format template",
			zap.Error(err))
	}

	return messages
}

func GetSystemPrompt(agent config.AgentType) string {
	switch agent {
	case config.CoordinatorAgent:
		return loadPromptTemplate(string(config.CoordinatorAgent))
	case config.PlannerAgent:
		return loadPromptTemplate(string(config.PlannerAgent))
	case config.SupervisorAgent:
		return loadPromptTemplate(string(config.SupervisorAgent))
	case config.ResearcherAgent:
		return loadPromptTemplate(string(config.ResearcherAgent))
	case config.CoderAgent:
		return loadPromptTemplate(string(config.CoderAgent))
	case config.BrowserAgent:
		return loadPromptTemplate(string(config.BrowserAgent))
	case config.ReporterAgent:
		return loadPromptTemplate(string(config.ReporterAgent))
	default:
		log.GetLogger().Warn("[getSystemPrompt]unknown agent type",
			zap.String("agent", string(agent)))
		return ""
	}
}

// loadPromptTemplate reads and processes prompt templates from the file system
func loadPromptTemplate(promptName string) string {
	// Get current working directory to support running from different paths
	workDir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.GetLogger().Error("[loadPromptTemplate]failed to get working directory",
			zap.Error(err))
		return ""
	}

	// Try different possible paths to find the template file
	possiblePaths := []string{
		filepath.Join(workDir, "internal", "prompts", promptName+".md"),
		filepath.Join("internal", "prompts", promptName+".md"),
		filepath.Join("..", "internal", "prompts", promptName+".md"),
	}

	var content []byte
	var templatePath string
	var readErr error

	// Try each possible path
	for _, path := range possiblePaths {
		content, readErr = os.ReadFile(path)
		if readErr == nil {
			templatePath = path
			break
		}
	}

	// If all paths fail, log error and return empty string
	if readErr != nil {
		log.GetLogger().Error("[loadPromptTemplate]failed to read template file",
			zap.Strings("attempted_paths", possiblePaths),
			zap.Error(readErr))
		return ""
	}

	log.GetLogger().Debug("[loadPromptTemplate]successfully loaded template",
		zap.String("path", templatePath))

	// Convert template content to string
	template := string(content)

	// Convert special format placeholders <<VAR>> to standard {VAR} format
	re := regexp.MustCompile(`<<([^>>]+)>>`)
	template = re.ReplaceAllString(template, "{$1}")

	return template
}
