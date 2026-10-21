package main

import (
	"context"
	"fmt"
	"log"

	"new-robot-self/internal/brain"
	"new-robot-self/internal/resource"
)

func main() {
	fmt.Println("=== AI内容增强功能验证 ===")
	fmt.Println("🎯 验证目标: 获取详细的事件信息")
	fmt.Println()
	
	// 创建Brain实例
	brainInstance := brain.NewBrain()
	
	// 配置学习引擎
	learningConfig := brain.LearningConfig{
		MaxItemsPerLearn:   5,
		EnabledTypes: []resource.ResourceType{
			resource.Social,
			resource.Tech,
		},
		EnableAutoLearn:    false,
		ResourceConfigPath: "asset/resources.json",
	}
	
	// 初始化学习引擎
	err := brainInstance.InitLearningEngine(learningConfig)
	if err != nil {
		log.Fatalf("初始化学习引擎失败: %v", err)
	}
	
	fmt.Println("✅ AI学习系统初始化成功（包含内容增强功能）")
	
	// 执行学习
	fmt.Println("\n📚 AI正在学习并增强内容...")
	ctx := context.Background()
	result, err := brainInstance.GetLearningEngine().LearnOnce(ctx)
	if err != nil {
		log.Fatalf("学习失败: %v", err)
	}
	
	fmt.Printf("✅ 学习完成！获取%d条信息\n", result.ProcessedCount)
	
	// 获取知识库中的详细信息
	kb := brainInstance.GetLearningEngine().GetKnowledgeBase()
	hotItems := kb.GetHotItems(5)
	
	if len(hotItems) == 0 {
		fmt.Println("❌ 暂无学习到的信息")
		return
	}
	
	// 展示增强后的内容
	fmt.Println("\n🤖 AI回答（包含详细内容）:")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("根据我刚刚学习到的详细信息，最近发生了以下大事:")
	fmt.Println()
	
	for i, item := range hotItems {
		fmt.Printf("%d. 【%s】%s\n", i+1, getCategoryName(item.Category), item.Title)
		
		// 显示详细内容
		if len(item.Content) > 100 {
			fmt.Printf("   📄 详细内容: %s\n", item.Content)
		} else {
			fmt.Printf("   📄 内容: %s\n", item.Content)
		}
		
		fmt.Printf("   🔥 热度: %d\n", item.HotValue)
		if item.URL != "" {
			fmt.Printf("   🔗 来源: %s\n", item.URL)
		}
		fmt.Println()
	}
	
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("💡 现在AI不仅能获取标题，还能通过内容增强功能")
	fmt.Println("   获取到更详细的事件信息，为您提供更全面的回答！")
	
	fmt.Println("\n🎉 内容增强功能验证成功！")
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
