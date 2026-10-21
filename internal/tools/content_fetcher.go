package tools

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/bytedance/sonic"
	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/schema"
)

// ContentFetcherTool 内容深度获取工具
// 从网页URL提取正文内容，去除广告和无关信息
type ContentFetcherTool struct {
	BaseToolImpl
	client *http.Client
}

type fetchParams struct {
	URL           string `json:"url"`                      // 网页URL
	ExtractImages bool   `json:"extract_images,omitempty"` // 是否提取图片
	MaxLength     int    `json:"max_length,omitempty"`     // 最大内容长度，默认5000
}

type fetchResult struct {
	URL         string   `json:"url"`
	Title       string   `json:"title"`
	Content     string   `json:"content"`
	Images      []string `json:"images,omitempty"`
	WordCount   int      `json:"word_count"`
	FetchedAt   time.Time `json:"fetched_at"`
	ContentType string   `json:"content_type"`
}

func NewContentFetcherTool() *ContentFetcherTool {
	return &ContentFetcherTool{
		BaseToolImpl: BaseToolImpl{
			Category:       "info_fetcher",
			EnableLearning: true,
			ToolName:       "content_fetcher",
		},
		client: &http.Client{Timeout: 15 * time.Second},
	}
}

func (t *ContentFetcherTool) Info(_ context.Context) (*schema.ToolInfo, error) {
	return &schema.ToolInfo{
		Name: "content_fetcher",
		Desc: "从网页URL提取正文内容，去除广告和无关信息。返回标题、正文、图片等结构化内容",
		ParamsOneOf: schema.NewParamsOneOfByParams(map[string]*schema.ParameterInfo{
			"url":           {Desc: "网页URL", Type: schema.String, Required: true},
			"extract_images": {Desc: "是否提取图片", Type: schema.Boolean, Required: false},
			"max_length":    {Desc: "最大内容长度，默认5000", Type: schema.Integer, Required: false},
		}),
	}, nil
}

func (t *ContentFetcherTool) InvokableRun(ctx context.Context, argumentsInJSON string, _ ...tool.Option) (string, error) {
	p := fetchParams{}
	if err := sonic.UnmarshalString(argumentsInJSON, &p); err != nil {
		return "", err
	}

	// 验证参数
	if p.URL == "" {
		return "", fmt.Errorf("url is required")
	}
	if p.MaxLength == 0 || p.MaxLength > 20000 {
		p.MaxLength = 5000
	}

	// 获取网页内容
	req, err := http.NewRequestWithContext(ctx, "GET", p.URL, nil)
	if err != nil {
		return "", err
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36")

	resp, err := t.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to fetch URL: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("HTTP status: %d", resp.StatusCode)
	}

	// 解析HTML
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to parse HTML: %w", err)
	}

	// 提取标题
	title := t.extractTitle(doc)

	// 提取正文内容
	content := t.extractContent(doc)

	// 限制长度
	if len(content) > p.MaxLength {
		content = content[:p.MaxLength] + "..."
	}

	// 提取图片
	var images []string
	if p.ExtractImages {
		images = t.extractImages(doc)
	}

	// 统计字数
	wordCount := len(strings.Fields(content))

	result := fetchResult{
		URL:         p.URL,
		Title:       title,
		Content:     content,
		Images:      images,
		WordCount:   wordCount,
		FetchedAt:   time.Now(),
		ContentType: resp.Header.Get("Content-Type"),
	}

	b, _ := sonic.Marshal(result)
	return string(b), nil
}

// extractTitle 提取网页标题
func (t *ContentFetcherTool) extractTitle(doc *goquery.Document) string {
	// 尝试多种方式获取标题
	title := ""
	if title = doc.Find("title").Text(); title != "" {
		return strings.TrimSpace(title)
	}
	if title = doc.Find("h1").First().Text(); title != "" {
		return strings.TrimSpace(title)
	}
	if title = doc.Find("meta[property='og:title']").AttrOr("content", ""); title != "" {
		return strings.TrimSpace(title)
	}
	return "未知标题"
}

// extractContent 提取正文内容
func (t *ContentFetcherTool) extractContent(doc *goquery.Document) string {
	// 尝试常见的正文选择器
	selectors := []string{
		"article",
		"main",
		".content",
		".article-content",
		".post-content",
		".entry-content",
		"#content",
		"#main",
		".markdown-body",
	}

	var content string
	for _, selector := range selectors {
		if selection := doc.Find(selector); selection.Length() > 0 {
			content = selection.Text()
			if len(content) > 200 { // 确保内容足够长
				break
			}
		}
	}

	// 如果没有找到特定选择器，使用通用方法
	if content == "" || len(content) < 200 {
		content = t.extractGenericContent(doc)
	}

	// 清理内容
	content = t.cleanContent(content)

	return content
}

// extractGenericContent 通用内容提取方法
func (t *ContentFetcherTool) extractGenericContent(doc *goquery.Document) string {
	// 移除不需要的标签
	doc.Find("script, style, nav, footer, header, aside, iframe, noscript").Remove()

	// 获取body文本
	content := doc.Find("body").Text()

	return content
}

// cleanContent 清理内容
func (t *ContentFetcherTool) cleanContent(content string) string {
	// 移除多余空白
	content = strings.TrimSpace(content)

	// 替换多个空白为单个空格
	content = strings.Join(strings.Fields(content), " ")

	// 移除常见的广告词
	adKeywords := []string{"广告", "advertisement", "赞助", "sponsored", "点击查看", "了解更多"}
	for _, keyword := range adKeywords {
		content = strings.ReplaceAll(content, keyword, "")
	}

	return content
}

// extractImages 提取图片URL
func (t *ContentFetcherTool) extractImages(doc *goquery.Document) []string {
	images := make([]string, 0)
	doc.Find("img").Each(func(i int, s *goquery.Selection) {
		if i >= 10 { // 最多提取10张图片
			return
		}

		src, exists := s.Attr("src")
		if exists && src != "" {
			// 跳过小图片（可能是图标）
			if width, _ := s.Attr("width"); width != "" && width < "50" {
				return
			}
			images = append(images, src)
		}
	})

	return images
}

// AfterRun 工具调用后的后处理
// 集成学习引擎，保存获取的内容到知识库
func (t *ContentFetcherTool) AfterRun(ctx context.Context, result string, err error) {
	// 使用默认的BaseToolImpl后处理
	t.BaseToolImpl.AfterRun(ctx, result, err)
}
