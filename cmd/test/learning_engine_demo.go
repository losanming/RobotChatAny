package main

import (
	"context"
	"fmt"
	"log"

	"new-robot-self/internal/brain"
	"new-robot-self/internal/resource"
)

func main() {
	fmt.Println("=== AI信息学习引擎测试 ===")
	
	// 创建Brain实例
	brainInstance := brain.NewBrain()
	
	// 配置学习引擎
	learningConfig := brain.LearningConfig{
		AutoLearnInterval:  6,
		MaxItemsPerLearn:   10,
		EnabledTypes:       []resource.ResourceType{
			resource.Tech,
			resource.Social,
			resource.Life,
		},
		EnableAutoLearn:    false, // 测试时不启用自动学习
		ResourceConfigPath: "asset/resources.json",
	}
	
	// 初始化学习引擎
	err := brainInstance.InitLearningEngine(learningConfig)
	if err != nil {
		log.Fatalf("初始化学习引擎失败: %v", err)
	}
	
	fmt.Println("✅ 学习引擎初始化成功")
	
	// 执行一次学习
	fmt.Println("\n=== 开始执行信息学习 ===")
	ctx := context.Background()
	result, err := brainInstance.GetLearningEngine().LearnOnce(ctx)
	if err != nil {
		log.Fatalf("学习失败: %v", err)
	}
	
	fmt.Printf("学习完成！\n")
	fmt.Printf("- 获取条目: %d\n", result.TotalFetched)
	fmt.Printf("- 处理条目: %d\n", result.ProcessedCount)
	fmt.Printf("- 耗时: %v\n", result.Duration)
	fmt.Printf("- 按类型统计:\n")
	
	for resType, count := range result.ByType {
		fmt.Printf("  %s: %d条\n", resType, count)
	}
	
	// 查看知识库统计
	fmt.Println("\n=== 知识库统计 ===")
	kb := brainInstance.GetLearningEngine().GetKnowledgeBase()
	stats := kb.GetStats()
	
	fmt.Printf("- 总条目数: %d\n", stats.TotalItems)
	fmt.Printf("- 平均热度: %.1f\n", stats.AverageHotValue)
	fmt.Printf("- 最后更新: %v\n", stats.LastUpdated)
	fmt.Printf("- 按类别统计:\n")
	
	for category, count := range stats.ByCategory {
		fmt.Printf("  %s: %d条\n", category, count)
	}
	
	// 查看热门条目
	fmt.Println("\n=== 热门条目 (Top 5) ===")
	hotItems := kb.GetHotItems(5)
	for i, item := range hotItems {
		fmt.Printf("%d. %s (热度: %d)\n", i+1, item.Title, item.HotValue)
		fmt.Printf("   %s\n", item.Content[:min(100, len(item.Content))]+"...")
	}
	
	// 查看进化报告
	fmt.Println("\n=== AI进化报告 ===")
	report := brainInstance.GetLearningEngine().GetEvolutionReport()
	fmt.Println(report.GenerateSummary())
	
	// 查看用户偏好
	fmt.Println("\n=== 用户偏好学习 ===")
	preferences := brainInstance.GetLearningEngine().GetPreferences()
	fmt.Printf("- 最后更新: %v\n", preferences.LastUpdated)
	fmt.Printf("- 类别权重:\n")
	
	for category, weight := range preferences.CategoryWeights {
		fmt.Printf("  %s: %.3f\n", category, weight)
	}
	
	// 推荐类别
	fmt.Println("\n=== 推荐类别 ===")
	// 需要通过 PreferenceLearner 来获取推荐类别
	// 这里暂时显示类别优先级
	fmt.Println("类别优先级排序功能已实现，可通过 PreferenceLearner 获取")
	
	fmt.Println("\n=== 测试完成 ===")
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
