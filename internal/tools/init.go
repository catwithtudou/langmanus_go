package tools

import "context"

func InitTools(ctx context.Context) {
	initDuckSearchTool(ctx)
	initCrawHtmlTool(ctx)
	initGoReplTool()
	initBrowserTool(ctx)
}
