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
		log.GetLogger().Error("[CrawHtmlTool] Initialization failed", zap.Error(err))
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
		Desc: "A tool for crawling web content from a specified URL and converting it to Markdown format. This tool will:\n" +
			"1. Load HTML content from the URL\n" +
			"2. Extract page title and main content\n" +
			"3. Convert HTML to Markdown format, ensuring good readability\n" +
			"4. Process image links to ensure absolute URLs\n" +
			"5. Return structured JSON result containing title and Markdown formatted content\n\n" +
			"Design Purpose:\n" +
			"- Help LLM better understand web content\n" +
			"- Extract clean, structured article content\n" +
			"- Convert content to a format easily processed by LLM\n" +
			"- Maintain integrity and readability of images and text\n\n" +
			"Usage Limitations:\n" +
			"- Only for content crawling, no page interaction support\n" +
			"- No mathematical computation support\n" +
			"- No file operations support\n" +
			"- Can only use search results or user-provided URLs",
		ParamsOneOf: schema.NewParamsOneOfByParams(map[string]*schema.ParameterInfo{
			"url": {
				Type:     schema.String,
				Desc:     "The URL of the webpage to crawl, must be a valid URL address",
				Required: true,
			},
		}),
	}, nil
}

func (t *CrawHtmlTool) InvokableRun(ctx context.Context, argumentsInJSON string, _ ...tool.Option) (string, error) {
	input := &crawHtmlToolInput{}
	err := json.Unmarshal([]byte(argumentsInJSON), input)
	if err != nil {
		log.GetLogger().Error("[CrawHtmlTool] Deserialization failed", zap.Error(err))
		return "", err
	}

	docs, err := urlLoader.Load(ctx, document.Source{
		URI: input.URL,
	})
	if err != nil {
		log.GetLogger().Error("[CrawHtmlTool] Failed to load URL", zap.Error(err))
		return "", err
	}

	result := &markdownContent{
		URL: input.URL,
	}
	for _, doc := range docs {
		// Get title from MetaData
		if title, ok := doc.MetaData["title"].(string); ok {
			result.Title = title
		}

		// Convert HTML to Markdown
		markdownContent := htmlToMarkdown(doc.Content)

		// Process image links
		result.Content = processImages(markdownContent, input.URL)
	}

	// Convert result to JSON
	jsonResult, err := json.Marshal(result)
	if err != nil {
		log.GetLogger().Error("[CrawHtmlTool] Failed to serialize result", zap.Error(err))
		return "", err
	}

	log.GetLogger().Info("[CrawHtmlTool] Result", zap.String("result", string(jsonResult)))

	return string(jsonResult), nil
}

// htmlToMarkdown converts HTML content to Markdown format
func htmlToMarkdown(htmlContent string) string {
	// Create parser
	extensions := parser.CommonExtensions | parser.AutoHeadingIDs
	p := parser.NewWithExtensions(extensions)

	// Create HTML renderer
	htmlFlags := html.CommonFlags | html.HrefTargetBlank
	opts := html.RendererOptions{Flags: htmlFlags}
	renderer := html.NewRenderer(opts)

	// Convert HTML to Markdown
	md := markdown.ToHTML([]byte(htmlContent), p, renderer)
	return string(md)
}

// processImages processes image links in Markdown content
func processImages(markdownContent string, baseURL string) string {
	// Match Markdown image syntax
	re := regexp.MustCompile(`!\[(.*?)\]\((.*?)\)`)

	// Replace image links
	return re.ReplaceAllStringFunc(markdownContent, func(match string) string {
		parts := re.FindStringSubmatch(match)
		if len(parts) != 3 {
			return match
		}

		// Parse relative URL
		imgURL, err := neturl.Parse(parts[2])
		if err != nil {
			return match
		}

		base, err := neturl.Parse(baseURL)
		if err != nil {
			return match
		}

		// Merge URLs
		absoluteURL := base.ResolveReference(imgURL).String()
		return "![" + parts[1] + "](" + absoluteURL + ")"
	})
}
