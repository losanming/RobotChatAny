package main

import (
	"context"
	"fmt"
	"log"
	"strings"

	"new-robot-self/internal/brain"
	"new-robot-self/internal/resource"
)

func main() {
	fmt.Println("=== AI学习系统验证测试 ===")
	fmt.Println("用户问题: \"最近发生了什么大事\"")
	fmt.Println()
	
	// 创建Brain实例
	brainInstance := brain.NewBrain()
	
	// 配置学习引擎
	learningConfig := brain.LearningConfig{
		AutoLearnInterval:  6,
		MaxItemsPerLearn:   20, // 获取更多条目用于验证
		EnabledTypes: []resource.ResourceType{
			resource.Tech,
			resource.Social,
			resource.Life,
			resource.Game,
			resource.Work,
		},
		EnableAutoLearn:    false,
		ResourceConfigPath: "asset/resources.json",
	}
	
	// 初始化学习引擎
	err := brainInstance.InitLearningEngine(learningConfig)
	if err != nil {
		log.Fatalf("初始化学习引擎失败: %v", err)
	}
	
	fmt.Println("✅ 学习引擎初始化成功")
	
	// 步骤1: 执行信息学习
	fmt.Println("\n=== 步骤1: AI正在学习最新信息 ===")
	ctx := context.Background()
	result, err := brainInstance.GetLearningEngine().LearnOnce(ctx)
	if err != nil {
		log.Fatalf("学习失败: %v", err)
	}
	
	fmt.Printf("✅ 学习完成！获取%d条，处理%d条，耗时%v\n", 
		result.TotalFetched, result.ProcessedCount, result.Duration)
	
	// 步骤2: 查看学习到的信息
	fmt.Println("\n=== 步骤2: AI学习到的最新信息 ===")
	displayLearnedInformation(brainInstance)
	
	// 步骤3: 构建增强的上下文
	fmt.Println("\n=== 步骤3: 为用户问题构建增强上下文 ===")
	userQuery := "最近发生了什么大事"
	contextBuilder := brainInstance.GetLearningEngine().GetContextBuilder()
	enhancedContext := contextBuilder.BuildEnhancedContext(ctx, userQuery)
	
	fmt.Println("增强上下文内容:")
	fmt.Println(strings.Repeat("=", 60))
	fmt.Println(enhancedContext)
	fmt.Println(strings.Repeat("=", 60))
	
	// 步骤4: 展示AI的回答能力
	fmt.Println("\n=== 步骤4: AI基于学习信息的回答 ===")
	generateAIResponse(userQuery, enhancedContext, brainInstance)
	
	// 步骤5: 展示进化状态
	fmt.Println("\n=== 步骤5: AI进化状态报告 ===")
	displayEvolutionStatus(brainInstance)
	
	fmt.Println("\n=== 验证测试完成 ===")
	fmt.Println("🎉 AI学习系统验证成功！")
}

func displayLearnedInformation(brainInstance *brain.Brain) {
	kb := brainInstance.GetLearningEngine().GetKnowledgeBase()
	hotItems := kb.GetHotItems(10) // 获取前10条热门信息
	
	if len(hotItems) == 0 {
		fmt.Println("暂无学习到的信息")
		return
	}
	
	fmt.Printf("AI学习到了 %d 条热门信息:\n", len(hotItems))
	for i, item := range hotItems {
		fmt.Printf("\n%d. 【%s】%s\n", i+1, item.Category, item.Title)
		fmt.Printf("   热度: %d\n", item.HotValue)
		fmt.Printf("   内容: %s\n", truncateString(item.Content, 150))
		if item.URL != "" {
			fmt.Printf("   来源: %s\n", item.URL)
		}
	}
}

func generateAIResponse(userQuery, enhancedContext string, brainInstance *brain.Brain) {
	fmt.Println("AI分析:")
	fmt.Println("- 用户询问:", userQuery)
	fmt.Println("- 利用学习到的信息构建上下文")
	fmt.Println("- 结合知识库中的热门事件")
	
	kb := brainInstance.GetLearningEngine().GetKnowledgeBase()
	hotItems := kb.GetHotItems(5)
	
	if len(hotItems) == 0 {
		fmt.Println("\nAI回答: 抱歉，我暂时没有学习到最新的信息。")
		return
	}
	
	fmt.Println("\nAI回答:")
	fmt.Println("根据我刚刚学习到的最新信息，最近发生了以下大事:")
	
	for i, item := range hotItems {
		fmt.Printf("\n%d. %s\n", i+1, item.Title)
		fmt.Printf("   %s\n", truncateString(item.Content, 200))
		fmt.Printf("   (热度: %d, 类别: %s)\n", item.HotValue, item.Category)
	}
	
	fmt.Println("\n这些信息来源于我的自动学习系统，我会持续获取和更新最新资讯。")
}

func displayEvolutionStatus(brainInstance *brain.Brain) {
	report := brainInstance.GetLearningEngine().GetEvolutionReport()
	
	fmt.Println("AI进化状态:")
	fmt.Printf("- 学习次数: %d次\n", report.Metrics.TotalLearningSessions)
	fmt.Printf("- 知识条目: %d条\n", report.Metrics.TotalItemsLearned)
	fmt.Printf("- 覆盖领域: %.1f%%\n", report.Metrics.CoverageScore*100)
	fmt.Printf("- 质量分数: %.1f%%\n", report.Metrics.QualityScore*100)
	
	kb := brainInstance.GetLearningEngine().GetKnowledgeBase()
	stats := kb.GetStats()
	
	fmt.Printf("\n知识库详情:")
	fmt.Printf("- 总条目: %d\n", stats.TotalItems)
	fmt.Printf("- 平均热度: %.1f\n", stats.AverageHotValue)
	fmt.Printf("- 最后更新: %v\n", stats.LastUpdated)
	
	fmt.Println("\n按类别分布:")
	for category, count := range stats.ByCategory {
		fmt.Printf("- %s: %d条\n", category, count)
	}
}

func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}
