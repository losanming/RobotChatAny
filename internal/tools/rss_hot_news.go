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

// RSSHotNewsTool RSS热点聚合工具
// 从多个RSS源获取热点新闻，支持去重、排序、内容增强
type RSSHotNewsTool struct {
	BaseToolImpl
	client *http.Client
}

type rssParams struct {
	Category string `json:"category,omitempty"` // 分类：tech, life, finance, all
	Limit    int    `json:"limit,omitempty"`    // 返回数量，默认20，最大50
}

type rssItem struct {
	Title       string    `json:"title"`
	Link        string    `json:"link"`
	Description string    `json:"description"`
	PubDate     time.Time `json:"pub_date"`
	Source      string    `json:"source"`
}

type rssResult struct {
	Items      []rssItem `json:"items"`
	Total      int       `json:"total"`
	Category   string    `json:"category"`
	FetchedAt  time.Time `json:"fetched_at"`
}

// RSS源配置
var rssSources = map[string]string{
	"36kr":   "https://36kr.com/feed",
	"sspai":  "https://sspai.com/feed",
	"ithome": "https://www.ithome.com/rss/",
}

func NewRSSHotNewsTool() *RSSHotNewsTool {
	return &RSSHotNewsTool{
		BaseToolImpl: BaseToolImpl{
			Category:       "info_fetcher",
			EnableLearning: true,
			ToolName:       "rss_hot_news",
		},
		client: &http.Client{Timeout: 15 * time.Second},
	}
}

func (t *RSSHotNewsTool) Info(_ context.Context) (*schema.ToolInfo, error) {
	return &schema.ToolInfo{
		Name: "rss_hot_news",
		Desc: "从多个RSS源获取热点新闻，支持分类筛选和数量限制。返回20-50条最新热点，自动去重和排序",
		ParamsOneOf: schema.NewParamsOneOfByParams(map[string]*schema.ParameterInfo{
			"category": {Desc: "分类筛选：tech(科技), life(生活), finance(财经), all(全部)", Type: schema.String, Required: false},
			"limit":    {Desc: "返回数量，默认20，最大50", Type: schema.Integer, Required: false},
		}),
	}, nil
}

func (t *RSSHotNewsTool) InvokableRun(ctx context.Context, argumentsInJSON string, _ ...tool.Option) (string, error) {
	p := rssParams{}
	if err := sonic.UnmarshalString(argumentsInJSON, &p); err != nil {
		return "", err
	}

	// 设置默认值
	if p.Category == "" {
		p.Category = "all"
	}
	if p.Limit == 0 || p.Limit > 50 {
		p.Limit = 20
	}

	// 从所有RSS源获取数据
	allItems := make([]rssItem, 0)
	for source, url := range rssSources {
		items, err := t.fetchRSS(ctx, source, url)
		if err != nil {
			fmt.Printf("Failed to fetch RSS from %s: %v\n", source, err)
			continue
		}
		allItems = append(allItems, items...)
	}

	// 去重
	uniqueItems := t.deduplicateItems(allItems)

	// 按发布时间排序
	t.sortItemsByDate(uniqueItems)

	// 应用分类筛选
	filteredItems := t.filterByCategory(uniqueItems, p.Category)

	// 限制数量
	if len(filteredItems) > p.Limit {
		filteredItems = filteredItems[:p.Limit]
	}

	result := rssResult{
		Items:     filteredItems,
		Total:     len(filteredItems),
		Category:  p.Category,
		FetchedAt: time.Now(),
	}

	b, _ := sonic.Marshal(result)
	return string(b), nil
}

// fetchRSS 从指定RSS源获取数据
func (t *RSSHotNewsTool) fetchRSS(ctx context.Context, source, url string) ([]rssItem, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
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

	items := make([]rssItem, 0)
	doc.Find("item").Each(func(i int, s *goquery.Selection) {
		title := s.Find("title").Text()
		link := s.Find("link").Text()
		description := s.Find("description").Text()
		pubDateStr := s.Find("pubDate").Text()

		// 解析发布时间
		pubDate, err := time.Parse(time.RFC1123Z, pubDateStr)
		if err != nil {
			pubDate = time.Now()
		}

		// 清理HTML标签
		description = t.cleanHTML(description)

		items = append(items, rssItem{
			Title:       strings.TrimSpace(title),
			Link:        strings.TrimSpace(link),
			Description: strings.TrimSpace(description),
			PubDate:     pubDate,
			Source:      source,
		})
	})

	return items, nil
}

// deduplicateItems 去重（基于标题）
func (t *RSSHotNewsTool) deduplicateItems(items []rssItem) []rssItem {
	seen := make(map[string]bool)
	unique := make([]rssItem, 0)

	for _, item := range items {
		key := strings.ToLower(item.Title)
		if !seen[key] {
			seen[key] = true
			unique = append(unique, item)
		}
	}

	return unique
}

// sortItemsByDate 按发布时间排序（最新的在前）
func (t *RSSHotNewsTool) sortItemsByDate(items []rssItem) {
	for i := 0; i < len(items)-1; i++ {
		for j := i + 1; j < len(items); j++ {
			if items[i].PubDate.Before(items[j].PubDate) {
				items[i], items[j] = items[j], items[i]
			}
		}
	}
}

// filterByCategory 按分类筛选
func (t *RSSHotNewsTool) filterByCategory(items []rssItem, category string) []rssItem {
	if category == "all" {
		return items
	}

	// 简单的关键词匹配分类
	keywords := map[string][]string{
		"tech":    {"科技", "AI", "人工智能", "互联网", "数码", "手机", "电脑", "芯片"},
		"life":   {"生活", "健康", "美食", "旅游", "娱乐", "电影", "音乐"},
		"finance": {"财经", "股票", "投资", "经济", "金融", "基金", "银行"},
	}

	filtered := make([]rssItem, 0)
	if kw, ok := keywords[category]; ok {
		for _, item := range items {
			for _, keyword := range kw {
				if strings.Contains(item.Title, keyword) || strings.Contains(item.Description, keyword) {
					filtered = append(filtered, item)
					break
				}
			}
		}
	}

	return filtered
}

// cleanHTML 清理HTML标签
func (t *RSSHotNewsTool) cleanHTML(html string) string {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return html
	}
	return doc.Text()
}

// AfterRun 工具调用后的后处理
// 集成学习引擎，保存热点新闻到知识库
func (t *RSSHotNewsTool) AfterRun(ctx context.Context, result string, err error) {
	// 使用默认的BaseToolImpl后处理
	t.BaseToolImpl.AfterRun(ctx, result, err)
}
