package brain

import (
	"context"
	"fmt"
	"strings"
)

// ContextBuilder 上下文构建器
type ContextBuilder struct {
	knowledgeBase *KnowledgeBase
	preferences  *UserPreferences
}

// NewContextBuilder 创建上下文构建器
func NewContextBuilder(knowledgeBase *KnowledgeBase, preferences *UserPreferences) *ContextBuilder {
	return &ContextBuilder{
		knowledgeBase: knowledgeBase,
		preferences:  preferences,
	}
}

// BuildEnhancedContext 构建增强的上下文
func (cb *ContextBuilder) BuildEnhancedContext(ctx context.Context, userQuery string) string {
	var contextBuilder strings.Builder
	
	// 1. 分析用户查询，确定相关领域
	relevantCategories := cb.analyzeQueryCategories(userQuery)
	
	// 2. 从知识库中检索相关信息
	relevantKnowledge := cb.retrieveRelevantKnowledge(relevantCategories, userQuery)
	
	// 3. 构建上下文
	if len(relevantKnowledge) > 0 {
		contextBuilder.WriteString("【相关知识】\n")
		for i, knowledge := range relevantKnowledge {
			contextBuilder.WriteString(fmt.Sprintf("%d. %s\n", i+1, knowledge.Title))
			contextBuilder.WriteString(fmt.Sprintf("   %s\n\n", knowledge.Content[:min(200, len(knowledge.Content))]))
		}
	}
	
	// 4. 添加时效性信息
	timeSensitiveInfo := cb.getTimeSensitiveInfo(relevantCategories)
	if timeSensitiveInfo != "" {
		contextBuilder.WriteString("【最新资讯】\n")
		contextBuilder.WriteString(timeSensitiveInfo)
		contextBuilder.WriteString("\n")
	}
	
	// 5. 添加用户偏好上下文
	if cb.preferences != nil && len(cb.preferences.CategoryWeights) > 0 {
		contextBuilder.WriteString("【用户兴趣偏好】\n")
		topCategories := cb.getTopCategories(3)
		for i, category := range topCategories {
			weight := cb.preferences.CategoryWeights[category]
			contextBuilder.WriteString(fmt.Sprintf("%d. %s (兴趣度: %.1f%%)\n", i+1, category, weight*100))
		}
		contextBuilder.WriteString("\n")
	}
	
	return contextBuilder.String()
}

// analyzeQueryCategories 分析查询相关的类别
func (cb *ContextBuilder) analyzeQueryCategories(query string) []string {
	// 简化的关键词匹配
	categoryKeywords := map[string][]string{
		"tech":     {"科技", "技术", "AI", "人工智能", "编程", "开发"},
		"social":   {"社会", "新闻", "时事", "政治", "经济"},
		"life":     {"生活", "美食", "旅游", "健康", "运动"},
		"game":     {"游戏", "手游", "端游", "steam"},
		"work":     {"工作", "职场", "求职", "面试"},
	}
	
	var relevantCategories []string
	queryLower := strings.ToLower(query)
	
	for category, keywords := range categoryKeywords {
		for _, keyword := range keywords {
			if strings.Contains(queryLower, strings.ToLower(keyword)) {
				relevantCategories = append(relevantCategories, category)
				break
			}
		}
	}
	
	// 如果没有匹配到，返回所有类别
	if len(relevantCategories) == 0 {
		relevantCategories = []string{"tech", "social", "life"}
	}
	
	return relevantCategories
}

// retrieveRelevantKnowledge 检索相关知识
func (cb *ContextBuilder) retrieveRelevantKnowledge(categories []string, query string) []*KnowledgeItem {
	var relevantKnowledge []*KnowledgeItem
	
	for _, category := range categories {
		items := cb.knowledgeBase.SearchByCategory(category, 3)
		relevantKnowledge = append(relevantKnowledge, items...)
	}
	
	// 按热度排序
	for i := 0; i < len(relevantKnowledge); i++ {
		for j := i + 1; j < len(relevantKnowledge); j++ {
			if relevantKnowledge[i].HotValue < relevantKnowledge[j].HotValue {
				relevantKnowledge[i], relevantKnowledge[j] = relevantKnowledge[j], relevantKnowledge[i]
			}
		}
	}
	
	// 限制数量
	if len(relevantKnowledge) > 5 {
		relevantKnowledge = relevantKnowledge[:5]
	}
	
	return relevantKnowledge
}

// getTimeSensitiveInfo 获取时效性信息
func (cb *ContextBuilder) getTimeSensitiveInfo(categories []string) string {
	// 获取最近的热门信息
	hotItems := cb.knowledgeBase.GetHotItems(3)
	
	if len(hotItems) == 0 {
		return ""
	}
	
	var infoBuilder strings.Builder
	infoBuilder.WriteString("以下是当前热门资讯：\n")
	
	for i, item := range hotItems {
		infoBuilder.WriteString(fmt.Sprintf("%d. %s (热度: %d)\n", i+1, item.Title, item.HotValue))
	}
	
	return infoBuilder.String()
}

// getTopCategories 获取权重最高的类别
func (cb *ContextBuilder) getTopCategories(limit int) []string {
	if cb.preferences == nil || len(cb.preferences.CategoryWeights) == 0 {
		return []string{}
	}
	
	// 按权重排序
	type weightPair struct {
		category string
		weight   float64
	}
	
	pairs := make([]weightPair, 0, len(cb.preferences.CategoryWeights))
	for category, weight := range cb.preferences.CategoryWeights {
		pairs = append(pairs, weightPair{category, weight})
	}
	
	// 简单排序
	for i := 0; i < len(pairs); i++ {
		for j := i + 1; j < len(pairs); j++ {
			if pairs[i].weight < pairs[j].weight {
				pairs[i], pairs[j] = pairs[j], pairs[i]
			}
		}
	}
	
	// 提取类别名称
	categories := make([]string, min(limit, len(pairs)))
	for i := 0; i < len(categories); i++ {
		categories[i] = pairs[i].category
	}
	
	return categories
}

// BuildConversationContext 构建对话上下文
func (cb *ContextBuilder) BuildConversationContext(ctx context.Context, userQuery, conversationHistory string) string {
	var contextBuilder strings.Builder
	
	// 基础上下文
	contextBuilder.WriteString(cb.BuildEnhancedContext(ctx, userQuery))
	
	// 对话历史上下文
	if conversationHistory != "" {
		contextBuilder.WriteString("【对话历史】\n")
		contextBuilder.WriteString(conversationHistory)
		contextBuilder.WriteString("\n")
	}
	
	// 用户查询
	contextBuilder.WriteString("【当前问题】\n")
	contextBuilder.WriteString(userQuery)
	
	return contextBuilder.String()
}

// min 辅助函数
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
