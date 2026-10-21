package main

import (
	"context"
	"fmt"
	"log"

	"new-robot-self/internal/brain"
	"new-robot-self/internal/resource"
)

func main() {
	fmt.Println("=== AI学习系统快速验证 ===")
	fmt.Println("🎯 验证问题: \"最近发生了什么大事\"")
	fmt.Println()
	
	// 创建Brain实例
	brainInstance := brain.NewBrain()
	
	// 配置学习引擎
	learningConfig := brain.LearningConfig{
		AutoLearnInterval:  6,
		MaxItemsPerLearn:   10,
		EnabledTypes: []resource.ResourceType{
			resource.Social,   // 重点关注社会热点
			resource.Tech,     // 科技资讯
		},
		EnableAutoLearn:    false,
		ResourceConfigPath: "/Users/danfwing/workspace/new-robot-self/asset/resources.json",
	}
	
	// 初始化学习引擎
	err := brainInstance.InitLearningEngine(learningConfig)
	if err != nil {
		log.Fatalf("初始化学习引擎失败: %v", err)
	}
	
	fmt.Println("✅ AI学习系统初始化成功")
	
	// 执行学习
	fmt.Println("\n📚 AI正在学习最新信息...")
	ctx := context.Background()
	result, err := brainInstance.GetLearningEngine().LearnOnce(ctx)
	if err != nil {
		log.Fatalf("学习失败: %v", err)
	}
	
	fmt.Printf("✅ 学习完成！获取%d条信息\n", result.ProcessedCount)
	
	// 获取热门信息
	kb := brainInstance.GetLearningEngine().GetKnowledgeBase()
	hotItems := kb.GetHotItems(5)
	
	if len(hotItems) == 0 {
		fmt.Println("❌ 暂无学习到的信息")
		return
	}
	
	// AI回答
	fmt.Println("\n🤖 AI回答:")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("根据我刚刚学习到的最新信息，最近发生了以下大事:")
	fmt.Println()
	
	for i, item := range hotItems {
		fmt.Printf("%d. 【%s】%s\n", i+1, getCategoryName(item.Category), item.Title)
		fmt.Printf("   %s\n", truncateContent(item.Content, 100))
		fmt.Printf("   🔥 热度: %d\n", item.HotValue)
		fmt.Println()
	}
	
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("💡 这些信息来源于我的自动学习系统，")
	fmt.Println("   我会持续获取和更新最新资讯来回答您的问题！")
	
	// 显示学习状态
	fmt.Println("\n📊 AI学习状态:")
	fmt.Printf("   知识条目: %d条\n", kb.GetStats().TotalItems)
	fmt.Printf("   学习效率: %.1f条/秒\n", float64(result.ProcessedCount)/result.Duration.Seconds())
	
	fmt.Println("\n🎉 验证成功！AI学习系统正常工作！")
}

func getCategoryName(category string) string {
	categoryNames := map[string]string{
		"social": "社会热点",
		"tech":   "科技资讯", 
		"life":   "生活相关",
		"game":   "游戏娱乐",
		"work":   "职场信息",
	}
	
	if name, ok := categoryNames[category]; ok {
		return name
	}
	return category
}

func truncateContent(content string, maxLen int) string {
	if len(content) <= maxLen {
		return content
	}
	return content[:maxLen] + "..."
}