package tools

import (
	"context"

	"github.com/cloudwego/eino-ext/components/tool/browseruse"
	"github.com/cloudwego/eino/components/tool"
	"go.uber.org/zap"

	"catwithtudou/langmanus_go/log"
)

var browserTool tool.InvokableTool

// GetBrowserTool returns the initialized browser tool instance
func GetBrowserTool() tool.InvokableTool {
	return browserTool
}

// initBrowserTool initializes the browser tool with default configuration
func initBrowserTool(ctx context.Context) {
	tool, err := browseruse.NewBrowserUseTool(ctx, &browseruse.Config{})
	if err != nil {
		log.GetLogger().Error("[InitBrowserTool]failed to init browser tool", zap.Error(err))
		panic(err)
	}

	browserTool = tool
}
