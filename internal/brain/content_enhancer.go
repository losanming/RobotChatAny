package brain

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

// ContentEnhancer 内容增强器
type ContentEnhancer struct {
	httpClient *http.Client
}

// NewContentEnhancer 创建内容增强器
func NewContentEnhancer() *ContentEnhancer {
	return &ContentEnhancer{
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// EnhanceContent 增强内容，获取更详细的信息
func (ce *ContentEnhancer) EnhanceContent(title, url, currentContent string) (string, error) {
	// 如果当前内容已经足够详细，直接返回
	if len(currentContent) > 200 {
		return currentContent, nil
	}
	
	// 如果没有URL，无法获取详细内容
	if url == "" {
		return ce.generateFallbackContent(title), nil
	}
	
	// 尝试从URL获取详细内容
	detailedContent, err := ce.fetchContentFromURL(url)
	if err != nil {
		fmt.Printf("获取详细内容失败: %v\n", err)
		return ce.generateFallbackContent(title), nil
	}
	
	if detailedContent != "" {
		return detailedContent, nil
	}
	
	return ce.generateFallbackContent(title), nil
}

// fetchContentFromURL 从URL获取内容
func (ce *ContentEnhancer) fetchContentFromURL(url string) (string, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}
	
	// 添加常见的请求头
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8")
	
	resp, err := ce.httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("HTTP status: %d", resp.StatusCode)
	}
	
	// 解析HTML
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return "", err
	}
	
	// 根据URL特征选择合适的选择器
	content := ce.extractContentByURL(doc, url)
	
	return content, nil
}

// extractContentByURL 根据URL提取内容
func (ce *ContentEnhancer) extractContentByURL(doc *goquery.Document, url string) string {
	var content string
	
	// 尝试多种选择器
	selectors := []string{
		".article-content",
		".RichContent-inner",
		".RichContent",
		".news_info",
		".content",
		".c-content",
		".detail_txt",
		".WB_text",
		".video-info",
		".desc-info",
		"article",
		".post-content",
		".entry-content",
	}
	
	for _, selector := range selectors {
		if text := doc.Find(selector).Text(); text != "" {
			content = text
			break
		}
	}
	
	// 清理内容
	content = strings.TrimSpace(content)
	content = strings.ReplaceAll(content, "\n", " ")
	content = strings.ReplaceAll(content, "\t", " ")
	
	// 限制长度
	if len(content) > 1000 {
		content = content[:1000] + "..."
	}
	
	return content
}

// generateFallbackContent 生成备用内容
func (ce *ContentEnhancer) generateFallbackContent(title string) string {
	// 基于标题生成一些通用描述
	return fmt.Sprintf("关于「%s」的详细信息正在获取中。这是一个热门话题，引起了广泛关注。", title)
}

// BatchEnhanceItems 批量增强知识条目
func (ce *ContentEnhancer) BatchEnhanceItems(items []*KnowledgeItem) []*KnowledgeItem {
	enhancedItems := make([]*KnowledgeItem, 0, len(items))
	
	for _, item := range items {
		enhancedContent, err := ce.EnhanceContent(item.Title, item.URL, item.Content)
		if err != nil {
			fmt.Printf("增强内容失败 [%s]: %v\n", item.Title, err)
			enhancedItems = append(enhancedItems, item)
		} else {
			// 创建增强后的副本
			enhancedItem := *item
			enhancedItem.Content = enhancedContent
			enhancedItem.UpdatedAt = time.Now()
			enhancedItems = append(enhancedItems, &enhancedItem)
		}
	}
	
	return enhancedItems
}
