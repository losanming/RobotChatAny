package brain

import (
	"fmt"
	"sync"
	"time"
)

// EvolutionEngine 进化引擎
type EvolutionEngine struct {
	learningHistory []*LearningResult
	metrics        *EvolutionMetrics
	mutex          sync.RWMutex
}

// EvolutionMetrics 进化指标
type EvolutionMetrics struct {
	TotalLearningSessions  int                    `json:"total_learning_sessions"`
	TotalItemsLearned      int                    `json:"total_items_learned"`
	AverageItemsPerSession float64                `json:"average_items_per_session"`
	LearningRate           float64                `json:"learning_rate"` // 学习效率
	KnowledgeGrowthRate    float64                `json:"knowledge_growth_rate"`
	CoverageScore          float64                `json:"coverage_score"`
	QualityScore           float64                `json:"quality_score"`
	FirstLearningTime      time.Time              `json:"first_learning_time"`
	LastLearningTime       time.Time              `json:"last_learning_time"`
	ByCategoryPerformance  map[string]*CategoryPerformance `json:"by_category_performance"`
}

// CategoryPerformance 类别性能
type CategoryPerformance struct {
	TotalItems      int     `json:"total_items"`
	AverageQuality  float64 `json:"average_quality"`
	LearningRate    float64 `json:"learning_rate"`
	LastUpdated     time.Time `json:"last_updated"`
}

// NewEvolutionEngine 创建进化引擎
func NewEvolutionEngine() *EvolutionEngine {
	return &EvolutionEngine{
		learningHistory: make([]*LearningResult, 0),
		metrics: &EvolutionMetrics{
			ByCategoryPerformance: make(map[string]*CategoryPerformance),
		},
	}
}

// RecordLearning 记录学习结果
func (ee *EvolutionEngine) RecordLearning(result *LearningResult) {
	ee.mutex.Lock()
	defer ee.mutex.Unlock()
	
	// 添加到历史记录
	ee.learningHistory = append(ee.learningHistory, result)
	
	// 更新基础指标
	ee.metrics.TotalLearningSessions++
	ee.metrics.TotalItemsLearned += result.ProcessedCount
	ee.metrics.AverageItemsPerSession = float64(ee.metrics.TotalItemsLearned) / float64(ee.metrics.TotalLearningSessions)
	
	// 更新时间
	if ee.metrics.FirstLearningTime.IsZero() {
		ee.metrics.FirstLearningTime = result.StartTime
	}
	ee.metrics.LastLearningTime = result.EndTime
	
	// 更新类别性能
	for category, count := range result.ByType {
		if perf, exists := ee.metrics.ByCategoryPerformance[string(category)]; exists {
			perf.TotalItems += count
			perf.LastUpdated = time.Now()
		} else {
			ee.metrics.ByCategoryPerformance[string(category)] = &CategoryPerformance{
				TotalItems:  count,
				LastUpdated: time.Now(),
			}
		}
	}
	
	// 计算进化指标
	ee.calculateEvolutionMetrics()
}

// calculateEvolutionMetrics 计算进化指标
func (ee *EvolutionEngine) calculateEvolutionMetrics() {
	if len(ee.learningHistory) < 2 {
		return
	}
	
	// 计算学习率（最近几次学习的平均效率）
	recentSessions := ee.learningHistory
	if len(recentSessions) > 5 {
		recentSessions = ee.learningHistory[len(ee.learningHistory)-5:]
	}
	
	var totalEfficiency float64
	for _, session := range recentSessions {
		efficiency := float64(session.ProcessedCount) / session.Duration.Seconds()
		totalEfficiency += efficiency
	}
	ee.metrics.LearningRate = totalEfficiency / float64(len(recentSessions))
	
	// 计算知识增长率
	if len(ee.learningHistory) >= 2 {
		firstSession := ee.learningHistory[0]
		lastSession := ee.learningHistory[len(ee.learningHistory)-1]
		
		timeDiff := lastSession.StartTime.Sub(firstSession.StartTime).Hours()
		itemsDiff := lastSession.ProcessedCount - firstSession.ProcessedCount
		
		if timeDiff > 0 {
			ee.metrics.KnowledgeGrowthRate = float64(itemsDiff) / timeDiff
		}
	}
	
	// 计算覆盖度分数
	totalCategories := len(ee.metrics.ByCategoryPerformance)
	ee.metrics.CoverageScore = float64(totalCategories) / 9.0 // 假设有9个主要类别
	
	// 计算质量分数（基于处理成功率）
	totalFetched := 0
	totalProcessed := 0
	for _, session := range ee.learningHistory {
		totalFetched += session.TotalFetched
		totalProcessed += session.ProcessedCount
	}
	
	if totalFetched > 0 {
		ee.metrics.QualityScore = float64(totalProcessed) / float64(totalFetched)
	}
}

// GenerateReport 生成进化报告
func (ee *EvolutionEngine) GenerateReport() *EvolutionReport {
	ee.mutex.RLock()
	defer ee.mutex.RUnlock()
	
	return &EvolutionReport{
		Metrics:        *ee.metrics,
		LearningHistory: ee.learningHistory,
		GeneratedAt:    time.Now(),
	}
}

// GetMetrics 获取进化指标
func (ee *EvolutionEngine) GetMetrics() *EvolutionMetrics {
	ee.mutex.RLock()
	defer ee.mutex.RUnlock()
	
	metricsCopy := *ee.metrics
	metricsCopy.ByCategoryPerformance = make(map[string]*CategoryPerformance)
	
	for k, v := range ee.metrics.ByCategoryPerformance {
		perfCopy := *v
		metricsCopy.ByCategoryPerformance[k] = &perfCopy
	}
	
	return &metricsCopy
}

// EvolutionReport 进化报告
type EvolutionReport struct {
	Metrics         EvolutionMetrics   `json:"metrics"`
	LearningHistory []*LearningResult `json:"learning_history"`
	GeneratedAt     time.Time          `json:"generated_at"`
	Summary         string             `json:"summary"`
}

// GenerateSummary 生成摘要
func (report *EvolutionReport) GenerateSummary() string {
	summary := fmt.Sprintf("AI进化报告：\n")
	summary += fmt.Sprintf("- 总学习次数: %d\n", report.Metrics.TotalLearningSessions)
	summary += fmt.Sprintf("- 总学习条目: %d\n", report.Metrics.TotalItemsLearned)
	summary += fmt.Sprintf("- 平均每次学习: %.1f条\n", report.Metrics.AverageItemsPerSession)
	summary += fmt.Sprintf("- 学习效率: %.2f条/秒\n", report.Metrics.LearningRate)
	summary += fmt.Sprintf("- 知识增长率: %.2f条/小时\n", report.Metrics.KnowledgeGrowthRate)
	summary += fmt.Sprintf("- 覆盖度分数: %.1f%%\n", report.Metrics.CoverageScore*100)
	summary += fmt.Sprintf("- 质量分数: %.1f%%\n", report.Metrics.QualityScore*100)
	
	return summary
}
