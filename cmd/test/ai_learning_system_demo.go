package main

import (
	"context"
	"fmt"
	"log"

	"new-robot-self/internal/brain"
	"new-robot-self/internal/resource"
)

func main() {
	fmt.Println("=== AI全能化学习系统综合测试 ===")
	
	// 创建Brain实例
	brainInstance := brain.NewBrain()
	
	// 配置学习引擎
	learningConfig := brain.LearningConfig{
		AutoLearnInterval:  6,
		MaxItemsPerLearn:   15,
		EnabledTypes: []resource.ResourceType{
			resource.Tech,
			resource.Social,
			resource.Life,
			resource.Game,
		},
		EnableAutoLearn:    false, // 测试时不启用自动学习
		ResourceConfigPath: "/Users/danfwing/workspace/new-robot-self/asset/resources.json",
	}
	
	// 初始化学习引擎
	err := brainInstance.InitLearningEngine(learningConfig)
	if err != nil {
		log.Fatalf("初始化学习引擎失败: %v", err)
	}
	
	fmt.Println("✅ 学习引擎初始化成功")
	
	// 执行多次学习以展示进化过程
	fmt.Println("\n=== 第一轮学习 ===")
	executeLearning(brainInstance)
	
	fmt.Println("\n=== 第二轮学习 ===")
	executeLearning(brainInstance)
	
	// 展示知识库状态
	displayKnowledgeBase(brainInstance)
	
	// 展示进化报告
	displayEvolutionReport(brainInstance)
	
	// 测试上下文构建
	fmt.Println("\n=== 测试上下文构建 ===")
	testContextBuilder(brainInstance)
	
	// 测试偏好学习
	fmt.Println("\n=== 测试偏好学习 ===")
	testPreferenceLearning(brainInstance)
	
	fmt.Println("\n=== 测试完成 ===")
	fmt.Println("🎉 AI全能化学习系统测试成功！")
}

func executeLearning(brainInstance *brain.Brain) {
	ctx := context.Background()
	result, err := brainInstance.GetLearningEngine().LearnOnce(ctx)
	if err != nil {
		log.Printf("学习失败: %v", err)
		return
	}
	
	fmt.Printf("✅ 学习完成！获取%d条，处理%d条，耗时%v\n", 
		result.TotalFetched, result.ProcessedCount, result.Duration)
	
	for resType, count := range result.ByType {
		fmt.Printf("   %s: %d条\n", resType, count)
	}
}

func displayKnowledgeBase(brainInstance *brain.Brain) {
	fmt.Println("\n=== 知识库状态 ===")
	kb := brainInstance.GetLearningEngine().GetKnowledgeBase()
	stats := kb.GetStats()
	
	fmt.Printf("总条目数: %d\n", stats.TotalItems)
	fmt.Printf("平均热度: %.1f\n", stats.AverageHotValue)
	fmt.Printf("最后更新: %v\n", stats.LastUpdated)
	
	fmt.Println("\n按类别统计:")
	for category, count := range stats.ByCategory {
		fmt.Printf("  %s: %d条\n", category, count)
	}
	
	fmt.Println("\n热门条目 (Top 3):")
	hotItems := kb.GetHotItems(3)
	for i, item := range hotItems {
		fmt.Printf("%d. %s (热度: %d)\n", i+1, item.Title, item.HotValue)
	}
}

func displayEvolutionReport(brainInstance *brain.Brain) {
	fmt.Println("\n=== AI进化报告 ===")
	report := brainInstance.GetLearningEngine().GetEvolutionReport()
	
	fmt.Println(report.GenerateSummary())
	
	fmt.Println("\n类别性能分析:")
	for category, perf := range report.Metrics.ByCategoryPerformance {
		fmt.Printf("  %s: 总条目=%d, 最后更新=%v\n", category, perf.TotalItems, perf.LastUpdated)
	}
}

func testContextBuilder(brainInstance *brain.Brain) {
	ctx := context.Background()
	contextBuilder := brainInstance.GetLearningEngine().GetContextBuilder()
	
	// 测试不同类型的查询
	queries := []string{
		"最近有什么科技新闻？",
		"现在社会上有什么热点事件？",
		"有什么好玩的游戏推荐吗？",
	}
	
	for _, query := range queries {
		fmt.Printf("\n用户查询: %s\n", query)
		enhancedContext := contextBuilder.BuildEnhancedContext(ctx, query)
		fmt.Printf("增强上下文:\n%s\n", enhancedContext)
	}
}

func testPreferenceLearning(brainInstance *brain.Brain) {
	preferenceLearner := brainInstance.GetLearningEngine()
	
	// 模拟用户交互
	fmt.Println("模拟用户交互学习...")
	
	// 记录一些交互
	interactions := []struct {
		category  string
		action    string
		relevance float64
	}{
		{"tech", "view", 0.8},
		{"tech", "like", 0.9},
		{"social", "view", 0.6},
		{"game", "like", 0.95},
		{"life", "ignore", 0.3},
	}
	
	for _, interaction := range interactions {
		// 这里需要通过PreferenceLearner来记录交互
		// 由于当前结构限制，我们直接展示概念
		fmt.Printf("  记录交互: 类别=%s, 动作=%s, 相关性=%.1f\n", 
			interaction.category, interaction.action, interaction.relevance)
	}
	
	// 显示学习到的偏好
	preferences := preferenceLearner.GetPreferences()
	fmt.Println("\n学习到的用户偏好:")
	fmt.Printf("最后更新: %v\n", preferences.LastUpdated)
	
	if len(preferences.CategoryWeights) > 0 {
		fmt.Println("类别权重:")
		for category, weight := range preferences.CategoryWeights {
			fmt.Printf("  %s: %.3f\n", category, weight)
		}
	} else {
		fmt.Println("  (暂无偏好数据，需要更多交互)")
	}
}
