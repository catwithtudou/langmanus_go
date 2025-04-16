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
	UserQueryKey   = "user_query"
	CurrentTimeKey = "CURRENT_TIME"
)

func GetSystemPromptSchemaMsg(ctx context.Context, agent config.AgentType, input map[string]any) []*schema.Message {
	// 合并输入参数与系统变量
	params := make(map[string]any)
	for k, v := range input {
		params[k] = v
	}

	// 添加系统默认变量
	params[CurrentTimeKey] = time.Now().Format("2006-01-02 15:04:05")

	template := prompt.FromMessages(schema.FString,
		schema.SystemMessage(getSystemPrompt(agent)),
		schema.UserMessage(UserQueryKey))

	messages, err := template.Format(ctx, params)
	if err != nil {
		log.GetLogger().Error("[GetSystemPromptSchemaMsg]failed to format template",
			zap.Error(err))
	}

	return messages
}

func getSystemPrompt(agent config.AgentType) string {
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

// loadPromptTemplate 从文件系统中读取提示模板并处理
func loadPromptTemplate(promptName string) string {
	// 获取当前工作目录，使得能够支持从不同路径运行
	workDir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.GetLogger().Error("[loadPromptTemplate]failed to get working directory",
			zap.Error(err))
		return ""
	}

	// 尝试不同的可能路径来查找模板文件
	possiblePaths := []string{
		filepath.Join(workDir, "internal", "prompts", promptName+".md"),
		filepath.Join("internal", "prompts", promptName+".md"),
		filepath.Join("..", "internal", "prompts", promptName+".md"),
	}

	var content []byte
	var templatePath string
	var readErr error

	// 尝试每个可能的路径
	for _, path := range possiblePaths {
		content, readErr = os.ReadFile(path)
		if readErr == nil {
			templatePath = path
			break
		}
	}

	// 如果所有路径都失败，则记录错误并返回空字符串
	if readErr != nil {
		log.GetLogger().Error("[loadPromptTemplate]failed to read template file",
			zap.Strings("attempted_paths", possiblePaths),
			zap.Error(readErr))
		return ""
	}

	log.GetLogger().Debug("[loadPromptTemplate]successfully loaded template",
		zap.String("path", templatePath))

	// 将模板内容转换为字符串
	template := string(content)

	// 将特殊格式的占位符 <<VAR>> 转换为标准的 {VAR} 格式
	re := regexp.MustCompile(`<<([^>>]+)>>`)
	template = re.ReplaceAllString(template, "{$1}")

	return template
}
