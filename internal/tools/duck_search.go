package tools

import (
	"catwithtudou/langmanus_go/log"
	"context"

	"github.com/cloudwego/eino-ext/components/tool/duckduckgo"
	"github.com/cloudwego/eino-ext/components/tool/duckduckgo/ddgsearch"
	"github.com/cloudwego/eino/components/tool"
	"go.uber.org/zap"
)

const (
	maxResults = 5
	cache      = true
	maxRetries = 5
)

var duckSearchTool tool.InvokableTool

func GetDuckSearchTool() tool.InvokableTool {
	return duckSearchTool
}

func initDuckSearchTool(ctx context.Context) {
	tool, err := duckduckgo.NewTool(ctx, &duckduckgo.Config{
		MaxResults: maxResults,
		Region:     ddgsearch.RegionCN,
		DDGConfig:  &ddgsearch.Config{Cache: cache},
	})
	if err != nil {
		log.GetLogger().Error("[InitDuckSearchTool]failed to init duckduckgo tool", zap.Error(err))
		panic(err)
	}

	duckSearchTool = tool
}
