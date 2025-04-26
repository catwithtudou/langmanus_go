package tools

import (
	"context"
	"encoding/json"
	neturl "net/url"
	"regexp"

	"catwithtudou/langmanus_go/log"

	urlloader "github.com/cloudwego/eino-ext/components/document/loader/url"
	"github.com/cloudwego/eino/components/document"
	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/schema"
	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
	"go.uber.org/zap"
)

type CrawHtmlTool struct {
}

type crawHtmlToolInput struct {
	URL string `json:"url"`
}

type markdownContent struct {
	Title   string `json:"title"`
	Content string `json:"content"`
	URL     string `json:"url"`
}

var urlLoader *urlloader.Loader

func initCrawHtmlTool(ctx context.Context) {
	loader, err := urlloader.NewLoader(ctx, nil)
	if err != nil {
		log.GetLogger().Error("[CrawHtmlTool]初始化失败", zap.Error(err))
		panic(err)
	}
	urlLoader = loader
}

func GetCrawHtmlTool() tool.InvokableTool {
	return &CrawHtmlTool{}
}

func (t *CrawHtmlTool) Info(_ context.Context) (*schema.ToolInfo, error) {
	return &schema.ToolInfo{
		Name: "crawl_tool",
		Desc: "用于从指定URL爬取网页内容并转换为Markdown格式的工具。该工具会：\n" +
			"1. 从URL加载HTML内容\n" +
			"2. 提取页面标题和正文内容\n" +
			"3. 将HTML转换为Markdown格式，确保内容具有良好的可读性\n" +
			"4. 处理图片链接，确保使用绝对URL\n" +
			"5. 返回结构化的JSON结果，包含标题和Markdown格式的内容\n\n" +
			"设计目的：\n" +
			"- 帮助LLM更好地理解网页内容\n" +
			"- 提取干净的、结构化的文章内容\n" +
			"- 将内容转换为易于LLM处理的格式\n" +
			"- 保持图片和文本的完整性和可读性\n\n" +
			"使用限制：\n" +
			"- 仅用于爬取内容，不支持页面交互\n" +
			"- 不支持执行数学计算\n" +
			"- 不支持文件操作\n" +
			"- 仅能使用搜索结果或用户提供的URL",
		ParamsOneOf: schema.NewParamsOneOfByParams(map[string]*schema.ParameterInfo{
			"url": {
				Type:     schema.String,
				Desc:     "要爬取的网页URL，必须是有效的URL地址",
				Required: true,
			},
		}),
	}, nil
}

func (t *CrawHtmlTool) InvokableRun(ctx context.Context, argumentsInJSON string, _ ...tool.Option) (string, error) {
	input := &crawHtmlToolInput{}
	err := json.Unmarshal([]byte(argumentsInJSON), input)
	if err != nil {
		log.GetLogger().Error("[CrawHtmlTool]反序列化失败", zap.Error(err))
		return "", err
	}

	docs, err := urlLoader.Load(ctx, document.Source{
		URI: input.URL,
	})
	if err != nil {
		log.GetLogger().Error("[CrawHtmlTool]加载url失败", zap.Error(err))
		return "", err
	}

	result := &markdownContent{
		URL: input.URL,
	}
	for _, doc := range docs {
		// 从 MetaData 中获取标题
		if title, ok := doc.MetaData["title"].(string); ok {
			result.Title = title
		}

		// 转换 HTML 到 Markdown
		markdownContent := htmlToMarkdown(doc.Content)

		// 处理图片链接
		result.Content = processImages(markdownContent, input.URL)
	}

	// 将结果转换为 JSON
	jsonResult, err := json.Marshal(result)
	if err != nil {
		log.GetLogger().Error("[CrawHtmlTool]序列化结果失败", zap.Error(err))
		return "", err
	}

	log.GetLogger().Info("[CrawHtmlTool]结果", zap.String("result", string(jsonResult)))

	return string(jsonResult), nil
}

// htmlToMarkdown 将 HTML 内容转换为 Markdown 格式
func htmlToMarkdown(htmlContent string) string {
	// 创建解析器
	extensions := parser.CommonExtensions | parser.AutoHeadingIDs
	p := parser.NewWithExtensions(extensions)

	// 创建 HTML 渲染器
	htmlFlags := html.CommonFlags | html.HrefTargetBlank
	opts := html.RendererOptions{Flags: htmlFlags}
	renderer := html.NewRenderer(opts)

	// 将 HTML 转换为 Markdown
	md := markdown.ToHTML([]byte(htmlContent), p, renderer)
	return string(md)
}

// processImages 处理 Markdown 中的图片链接
func processImages(markdownContent string, baseURL string) string {
	// 匹配 Markdown 图片语法
	re := regexp.MustCompile(`!\[(.*?)\]\((.*?)\)`)

	// 替换图片链接
	return re.ReplaceAllStringFunc(markdownContent, func(match string) string {
		parts := re.FindStringSubmatch(match)
		if len(parts) != 3 {
			return match
		}

		// 解析相对 URL
		imgURL, err := neturl.Parse(parts[2])
		if err != nil {
			return match
		}

		base, err := neturl.Parse(baseURL)
		if err != nil {
			return match
		}

		// 合并 URL
		absoluteURL := base.ResolveReference(imgURL).String()
		return "![" + parts[1] + "](" + absoluteURL + ")"
	})
}
