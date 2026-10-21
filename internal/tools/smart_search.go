package tools

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/bytedance/sonic"
	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/schema"
)

// SmartSearchTool 智能搜索工具
// 提供真实的联网搜索能力，支持多个搜索引擎
type SmartSearchTool struct {
	BaseToolImpl
	client *http.Client
}

type searchParams struct {
	Query      string `json:"query"`                // 搜索关键词
	Engine     string `json:"engine,omitempty"`     // 搜索引擎: bing, google, duckduckgo
	MaxResults int    `json:"max_results,omitempty"` // 最大结果数，默认10，最大20
	Language   string `json:"language,omitempty"`   // 语言: zh-CN, en-US
}

type searchResult struct {
	Title   string `json:"title"`
	URL     string `json:"url"`
	Snippet string `json:"snippet"`
	Source  string `json:"source"`
}

type searchResponse struct {
	Query      string         `json:"query"`
	Engine     string         `json:"engine"`
	Results    []searchResult `json:"results"`
	Total      int            `json:"total"`
	SearchedAt time.Time      `json:"searched_at"`
}

func NewSmartSearchTool() *SmartSearchTool {
	return &SmartSearchTool{
		BaseToolImpl: BaseToolImpl{
			Category:       "info_fetcher",
			EnableLearning: true,
			ToolName:       "smart_search",
		},
		client: &http.Client{Timeout: 10 * time.Second},
	}
}

func (t *SmartSearchTool) Info(_ context.Context) (*schema.ToolInfo, error) {
	return &schema.ToolInfo{
		Name: "smart_search",
		Desc: "智能联网搜索工具，支持多个搜索引擎。返回10-20条搜索结果，包含标题、URL和摘要",
		ParamsOneOf: schema.NewParamsOneOfByParams(map[string]*schema.ParameterInfo{
			"query":       {Desc: "搜索关键词", Type: schema.String, Required: true},
			"engine":      {Desc: "搜索引擎: bing, google, duckduckgo (默认bing)", Type: schema.String, Required: false},
			"max_results": {Desc: "最大结果数，默认10，最大20", Type: schema.Integer, Required: false},
			"language":    {Desc: "语言: zh-CN, en-US (默认zh-CN)", Type: schema.String, Required: false},
		}),
	}, nil
}

func (t *SmartSearchTool) InvokableRun(ctx context.Context, argumentsInJSON string, _ ...tool.Option) (string, error) {
	p := searchParams{}
	if err := sonic.UnmarshalString(argumentsInJSON, &p); err != nil {
		return "", err
	}

	// 设置默认值
	if p.Query == "" {
		return "", fmt.Errorf("query is required")
	}
	if p.Engine == "" {
		p.Engine = "bing"
	}
	if p.MaxResults == 0 || p.MaxResults > 20 {
		p.MaxResults = 10
	}
	if p.Language == "" {
		p.Language = "zh-CN"
	}

	// 根据搜索引擎执行搜索
	var results []searchResult
	var err error

	switch p.Engine {
	case "bing":
		results, err = t.searchBing(ctx, p.Query, p.MaxResults, p.Language)
	case "google":
		results, err = t.searchGoogle(ctx, p.Query, p.MaxResults, p.Language)
	case "duckduckgo":
		results, err = t.searchDuckDuckGo(ctx, p.Query, p.MaxResults)
	default:
		return "", fmt.Errorf("unsupported search engine: %s", p.Engine)
	}

	if err != nil {
		return "", fmt.Errorf("search failed: %w", err)
	}

	response := searchResponse{
		Query:      p.Query,
		Engine:     p.Engine,
		Results:    results,
		Total:      len(results),
		SearchedAt: time.Now(),
	}

	b, _ := sonic.Marshal(response)
	return string(b), nil
}

// searchBing 使用Bing搜索
func (t *SmartSearchTool) searchBing(ctx context.Context, query string, maxResults int, language string) ([]searchResult, error) {
	// 构建Bing搜索URL
	searchURL := fmt.Sprintf("https://www.bing.com/search?q=%s&setlang=%s&count=%d",
		url.QueryEscape(query), language, maxResults)

	req, err := http.NewRequestWithContext(ctx, "GET", searchURL, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36")

	resp, err := t.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP status: %d", resp.StatusCode)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}

	results := make([]searchResult, 0)
	doc.Find(".b_algo").Each(func(i int, s *goquery.Selection) {
		if i >= maxResults {
			return
		}

		title := s.Find("h2 a").Text()
		href, _ := s.Find("h2 a").Attr("href")
		snippet := s.Find(".b_caption p").Text()

		if title != "" && href != "" {
			results = append(results, searchResult{
				Title:   strings.TrimSpace(title),
				URL:     strings.TrimSpace(href),
				Snippet: strings.TrimSpace(snippet),
				Source:  "bing",
			})
		}
	})

	return results, nil
}

// searchGoogle 使用Google搜索
func (t *SmartSearchTool) searchGoogle(ctx context.Context, query string, maxResults int, language string) ([]searchResult, error) {
	// 构建Google搜索URL
	searchURL := fmt.Sprintf("https://www.google.com/search?q=%s&hl=%s&num=%d",
		url.QueryEscape(query), language, maxResults)

	req, err := http.NewRequestWithContext(ctx, "GET", searchURL, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36")

	resp, err := t.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP status: %d", resp.StatusCode)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}

	results := make([]searchResult, 0)
	doc.Find("div.g").Each(func(i int, s *goquery.Selection) {
		if i >= maxResults {
			return
		}

		title := s.Find("h3").Text()
		href, _ := s.Find("a").First().Attr("href")
		snippet := s.Find(".VwiC3b").Text()

		if title != "" && href != "" {
			// 清理Google的URL重定向
			if strings.HasPrefix(href, "/url?q=") {
				parsedURL, err := url.Parse(href)
				if err == nil {
					queryParams := parsedURL.Query()
					if realURL := queryParams.Get("q"); realURL != "" {
						href = realURL
					}
				}
			}

			results = append(results, searchResult{
				Title:   strings.TrimSpace(title),
				URL:     strings.TrimSpace(href),
				Snippet: strings.TrimSpace(snippet),
				Source:  "google",
			})
		}
	})

	return results, nil
}

// searchDuckDuckGo 使用DuckDuckGo搜索
func (t *SmartSearchTool) searchDuckDuckGo(ctx context.Context, query string, maxResults int) ([]searchResult, error) {
	// 构建DuckDuckGo搜索URL
	searchURL := fmt.Sprintf("https://html.duckduckgo.com/html/?q=%s", url.QueryEscape(query))

	req, err := http.NewRequestWithContext(ctx, "GET", searchURL, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36")

	resp, err := t.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP status: %d", resp.StatusCode)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}

	results := make([]searchResult, 0)
	doc.Find(".result").Each(func(i int, s *goquery.Selection) {
		if i >= maxResults {
			return
		}

		title := s.Find(".result__a").Text()
		href, _ := s.Find(".result__a").Attr("href")
		snippet := s.Find(".result__snippet").Text()

		if title != "" && href != "" {
			results = append(results, searchResult{
				Title:   strings.TrimSpace(title),
				URL:     strings.TrimSpace(href),
				Snippet: strings.TrimSpace(snippet),
				Source:  "duckduckgo",
			})
		}
	})

	return results, nil
}

// AfterRun 工具调用后的后处理
// 集成学习引擎，保存搜索结果到知识库
func (t *SmartSearchTool) AfterRun(ctx context.Context, result string, err error) {
	// 使用默认的BaseToolImpl后处理
	t.BaseToolImpl.AfterRun(ctx, result, err)
}
