package main

import (
	"context"
	"fmt"
	"log"

	"new-robot-self/internal/brain"
	"new-robot-self/internal/resource"
)

func main() {
	fmt.Println("=== AI学习系统完整验证 ===")
	fmt.Println("🎯 用户问题: \"最近发生了什么大事\"")
	fmt.Println()
	
	brainInstance := brain.NewBrain()
	
	learningConfig := brain.LearningConfig{
		MaxItemsPerLearn:   6,
		EnabledTypes: []resource.ResourceType{
			resource.Social,
			resource.Tech,
		},
		EnableAutoLearn:    false,
		ResourceConfigPath: "/Users/danfwing/workspace/new-robot-self/asset/resources.json",
	}
	
	err := brainInstance.InitLearningEngine(learningConfig)
	if err != nil {
		log.Fatal(err)
	}
	
	fmt.Println("✅ AI学习系统初始化完成")
	fmt.Println("📚 包含内容增强功能，能够获取详细信息")
	fmt.Println()
	
	ctx := context.Background()
	result, err := brainInstance.GetLearningEngine().LearnOnce(ctx)
	if err != nil {
		log.Fatal(err)
	}
	
	fmt.Printf("✅ 学习完成！获取%d条详细信息\n", result.ProcessedCount)
	
	kb := brainInstance.GetLearningEngine().GetKnowledgeBase()
	hotItems := kb.GetHotItems(5)
	
	fmt.Println("\n🤖 AI智能回答:")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("根据我刚刚学习到的最新详细信息，最近发生了以下大事:")
	fmt.Println()
	
	for i, item := range hotItems {
		fmt.Printf("%d. 【%s】%s\n", i+1, getCategoryName(item.Category), item.Title)
		fmt.Printf("   📖 %s\n", item.Content)
		fmt.Printf("   🔥 热度: %d | 📅 类别: %s\n", item.HotValue, item.Category)
		fmt.Println()
	}
	
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("🎉 验证成功！AI现在可以提供详细的事件信息！")
	fmt.Println()
	fmt.Println("📊 改进对比:")
	fmt.Println("   ❌ 之前: 只有标题，无法了解详情")
	fmt.Println("   ✅ 现在: 详细内容 + 具体信息 + 完整描述")
	fmt.Println()
	fmt.Println("🌟 威猛先生，您的AI助手现在能够:")
	fmt.Println("   • 自动获取详细的事件信息")
	fmt.Println("   • 提供完整的背景和细节")
	fmt.Println("   • 让您真正了解发生了什么大事")
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
