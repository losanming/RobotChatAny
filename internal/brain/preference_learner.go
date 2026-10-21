package brain

import (
	"sync"
	"time"

	"new-robot-self/internal/resource"
)

// PreferenceLearner 偏好学习器
type PreferenceLearner struct {
	preferences      *UserPreferences
	interactionHistory []*InteractionRecord
	mutex           sync.RWMutex
}

// UserPreferences 用户偏好
type UserPreferences struct {
	CategoryWeights map[string]float64 `json:"category_weights"` // 类别权重
	TimePreferences map[string]string `json:"time_preferences"` // 时间偏好
	TopicInterests  []string          `json:"topic_interests"`  // 主题兴趣
	LastUpdated     time.Time         `json:"last_updated"`
}

// InteractionRecord 交互记录
type InteractionRecord struct {
	Timestamp    time.Time
	Category     string
	Action       string // view, like, share, ignore
	Duration     time.Duration
	Relevance    float64 // 相关性评分
}

// NewPreferenceLearner 创建偏好学习器
func NewPreferenceLearner() *PreferenceLearner {
	return &PreferenceLearner{
		preferences: &UserPreferences{
			CategoryWeights: make(map[string]float64),
			TimePreferences: make(map[string]string),
			TopicInterests:  make([]string, 0),
			LastUpdated:     time.Now(),
		},
		interactionHistory: make([]*InteractionRecord, 0),
	}
}

// LearnFromResource 从资源学习偏好
func (pl *PreferenceLearner) LearnFromResource(res resource.Resource) {
	pl.mutex.Lock()
	defer pl.mutex.Unlock()
	
	// 基于资源类型更新权重
	for _, resType := range res.Types {
		category := string(resType)
		currentWeight := pl.preferences.CategoryWeights[category]
		
		// 热度越高，权重增加越多
		weightIncrease := float64(res.HotValue) / 10000.0
		if weightIncrease > 0.1 {
			weightIncrease = 0.1 // 限制最大增加量
		}
		
		pl.preferences.CategoryWeights[category] = currentWeight + weightIncrease
	}
	
	pl.preferences.LastUpdated = time.Now()
}

// RecordInteraction 记录用户交互
func (pl *PreferenceLearner) RecordInteraction(category, action string, duration time.Duration, relevance float64) {
	pl.mutex.Lock()
	defer pl.mutex.Unlock()
	
	record := &InteractionRecord{
		Timestamp: time.Now(),
		Category:  category,
		Action:    action,
		Duration:  duration,
		Relevance: relevance,
	}
	
	pl.interactionHistory = append(pl.interactionHistory, record)
	
	// 根据交互调整权重
	pl.adjustWeightFromInteraction(record)
	
	pl.preferences.LastUpdated = time.Now()
}

// adjustWeightFromInteraction 根据交互调整权重
func (pl *PreferenceLearner) adjustWeightFromInteraction(record *InteractionRecord) {
	category := record.Category
	currentWeight := pl.preferences.CategoryWeights[category]
	
	var weightChange float64
	
	switch record.Action {
	case "like":
		weightChange = 0.2
	case "share":
		weightChange = 0.15
	case "view":
		// 基于观看时长和相关性
		if record.Duration > 30*time.Second && record.Relevance > 0.7 {
			weightChange = 0.1
		} else if record.Duration > 10*time.Second {
			weightChange = 0.05
		}
	case "ignore":
		weightChange = -0.1
	}
	
	// 应用权重变化
	newWeight := currentWeight + weightChange
	if newWeight < 0 {
		newWeight = 0
	}
	if newWeight > 2.0 {
		newWeight = 2.0 // 限制最大权重
	}
	
	pl.preferences.CategoryWeights[category] = newWeight
}

// GetPreferences 获取用户偏好
func (pl *PreferenceLearner) GetPreferences() *UserPreferences {
	pl.mutex.RLock()
	defer pl.mutex.RUnlock()
	
	// 归一化权重
	normalizedPrefs := *pl.preferences
	normalizedPrefs.CategoryWeights = pl.normalizeWeights(pl.preferences.CategoryWeights)
	
	return &normalizedPrefs
}

// normalizeWeights 归一化权重
func (pl *PreferenceLearner) normalizeWeights(weights map[string]float64) map[string]float64 {
	if len(weights) == 0 {
		return weights
	}
	
	// 找到最大权重
	maxWeight := 0.0
	for _, weight := range weights {
		if weight > maxWeight {
			maxWeight = weight
		}
	}
	
	if maxWeight == 0 {
		return weights
	}
	
	// 归一化到0-1范围
	normalized := make(map[string]float64)
	for category, weight := range weights {
		normalized[category] = weight / maxWeight
	}
	
	return normalized
}

// GetCategoryPriority 获取类别优先级排序
func (pl *PreferenceLearner) GetCategoryPriority() []string {
	pl.mutex.RLock()
	defer pl.mutex.RUnlock()
	
	preferences := pl.GetPreferences()
	
	// 按权重排序
	type weightPair struct {
		category string
		weight   float64
	}
	
	pairs := make([]weightPair, 0, len(preferences.CategoryWeights))
	for category, weight := range preferences.CategoryWeights {
		pairs = append(pairs, weightPair{category, weight})
	}
	
	// 排序
	for i := 0; i < len(pairs); i++ {
		for j := i + 1; j < len(pairs); j++ {
			if pairs[i].weight < pairs[j].weight {
				pairs[i], pairs[j] = pairs[j], pairs[i]
			}
		}
	}
	
	// 提取类别名称
	categories := make([]string, len(pairs))
	for i, pair := range pairs {
		categories[i] = pair.category
	}
	
	return categories
}

// RecommendCategories 推荐类别
func (pl *PreferenceLearner) RecommendCategories(limit int) []string {
	priority := pl.GetCategoryPriority()
	
	if limit > 0 && len(priority) > limit {
		priority = priority[:limit]
	}
	
	return priority
}

// GetInteractionStats 获取交互统计
func (pl *PreferenceLearner) GetInteractionStats() map[string]interface{} {
	pl.mutex.RLock()
	defer pl.mutex.RUnlock()
	
	stats := make(map[string]interface{})
	
	// 按类别统计交互
	categoryStats := make(map[string]int)
	actionStats := make(map[string]int)
	totalDuration := time.Duration(0)
	
	for _, record := range pl.interactionHistory {
		categoryStats[record.Category]++
		actionStats[record.Action]++
		totalDuration += record.Duration
	}
	
	stats["total_interactions"] = len(pl.interactionHistory)
	stats["by_category"] = categoryStats
	stats["by_action"] = actionStats
	stats["total_duration"] = totalDuration.String()
	
	return stats
}
