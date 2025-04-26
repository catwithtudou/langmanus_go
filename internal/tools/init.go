package tools

import "context"

func InitTools(ctx context.Context) {
	initDuckSearchTool(ctx)
	initCrawHtmlTool(ctx)
}
