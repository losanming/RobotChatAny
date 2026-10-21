package brain

import (
	"fmt"
	"sync"
	"time"
)

// KnowledgeBase 知识库
type KnowledgeBase struct {
	items    map[string]*KnowledgeItem
	mutex    sync.RWMutex
	stats    *KnowledgeBaseStats
}

// KnowledgeItem 知识条目
type KnowledgeItem struct {
	ID        string                 `json:"id"`
	Title     string                 `json:"title"`
	Content   string                 `json:"content"`
	URL       string                 `json:"url"`
	Category  string                 `json:"category"`
	Source    string                 `json:"source"`
	HotValue  int64                  `json:"hot_value"`
	CreatedAt time.Time              `json:"created_at"`
	UpdatedAt time.Time              `json:"updated_at"`
	Tags      []string               `json:"tags"`
	Metadata  map[string]interface{} `json:"metadata"`
}

// KnowledgeBaseStats 知识库统计
type KnowledgeBaseStats struct {
	TotalItems      int                    `json:"total_items"`
	ByCategory      map[string]int         `json:"by_category"`
	LastUpdated     time.Time              `json:"last_updated"`
	AverageHotValue float64                `json:"average_hot_value"`
}

// NewKnowledgeBase 创建知识库
func NewKnowledgeBase() *KnowledgeBase {
	return &KnowledgeBase{
		items: make(map[string]*KnowledgeItem),
		stats: &KnowledgeBaseStats{
			ByCategory: make(map[string]int),
		},
	}
}

// AddItem 添加知识条目
func (kb *KnowledgeBase) AddItem(item *KnowledgeItem) error {
	kb.mutex.Lock()
	defer kb.mutex.Unlock()
	
	// 检查是否已存在
	if _, exists := kb.items[item.ID]; exists {
		return fmt.Errorf("知识条目已存在: %s", item.ID)
	}
	
	// 添加条目
	kb.items[item.ID] = item
	
	// 更新统计
	kb.stats.TotalItems++
	kb.stats.ByCategory[item.Category]++
	kb.stats.LastUpdated = time.Now()
	
	// 更新平均热度
	kb.updateAverageHotValue()
	
	return nil
}

// GetItem 获取知识条目
func (kb *KnowledgeBase) GetItem(id string) (*KnowledgeItem, bool) {
	kb.mutex.RLock()
	defer kb.mutex.RUnlock()
	
	item, exists := kb.items[id]
	return item, exists
}

// SearchByCategory 按类别搜索
func (kb *KnowledgeBase) SearchByCategory(category string, limit int) []*KnowledgeItem {
	kb.mutex.RLock()
	defer kb.mutex.RUnlock()
	
	var results []*KnowledgeItem
	for _, item := range kb.items {
		if item.Category == category {
			results = append(results, item)
			if limit > 0 && len(results) >= limit {
				break
			}
		}
	}
	
	return results
}

// SearchByTags 按标签搜索
func (kb *KnowledgeBase) SearchByTags(tags []string, limit int) []*KnowledgeItem {
	kb.mutex.RLock()
	defer kb.mutex.RUnlock()
	
	var results []*KnowledgeItem
	for _, item := range kb.items {
		if kb.hasMatchingTags(item.Tags, tags) {
			results = append(results, item)
			if limit > 0 && len(results) >= limit {
				break
			}
		}
	}
	
	return results
}

// GetHotItems 获取热门条目
func (kb *KnowledgeBase) GetHotItems(limit int) []*KnowledgeItem {
	kb.mutex.RLock()
	defer kb.mutex.RUnlock()
	
	// 按热度排序
	items := make([]*KnowledgeItem, 0, len(kb.items))
	for _, item := range kb.items {
		items = append(items, item)
	}
	
	// 简单排序（实际应该使用更高效的排序算法）
	for i := 0; i < len(items); i++ {
		for j := i + 1; j < len(items); j++ {
			if items[i].HotValue < items[j].HotValue {
				items[i], items[j] = items[j], items[i]
			}
		}
	}
	
	if limit > 0 && len(items) > limit {
		items = items[:limit]
	}
	
	return items
}

// GetStats 获取统计信息
func (kb *KnowledgeBase) GetStats() *KnowledgeBaseStats {
	kb.mutex.RLock()
	defer kb.mutex.RUnlock()
	
	// 返回统计信息的副本
	statsCopy := *kb.stats
	statsCopy.ByCategory = make(map[string]int)
	for k, v := range kb.stats.ByCategory {
		statsCopy.ByCategory[k] = v
	}
	
	return &statsCopy
}

// hasMatchingTags 检查是否有匹配的标签
func (kb *KnowledgeBase) hasMatchingTags(itemTags []string, searchTags []string) bool {
	tagMap := make(map[string]bool)
	for _, tag := range itemTags {
		tagMap[tag] = true
	}
	
	for _, searchTag := range searchTags {
		if tagMap[searchTag] {
			return true
		}
	}
	
	return false
}

// updateAverageHotValue 更新平均热度
func (kb *KnowledgeBase) updateAverageHotValue() {
	if kb.stats.TotalItems == 0 {
		kb.stats.AverageHotValue = 0
		return
	}
	
	var totalHotValue int64
	for _, item := range kb.items {
		totalHotValue += item.HotValue
	}
	
	kb.stats.AverageHotValue = float64(totalHotValue) / float64(kb.stats.TotalItems)
}

// generateKnowledgeID 生成知识ID
func generateKnowledgeID(title, url string) string {
	// 简单的ID生成策略
	return fmt.Sprintf("%s_%d", title, time.Now().UnixNano())
}

// extractTags 提取标签
func extractTags(title, content string) []string {
	// 简化的标签提取逻辑
	// 实际应该使用更复杂的NLP技术
	tags := make([]string, 0)
	
	// 这里可以添加关键词提取逻辑
	// 暂时返回空数组
	
	return tags
}
