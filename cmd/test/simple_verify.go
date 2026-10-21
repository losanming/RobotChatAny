package main

import (
	"context"
	"fmt"
	"log"

	"new-robot-self/internal/brain"
	"new-robot-self/internal/resource"
)

func main() {
	fmt.Println("=== AI学习系统验证 ===")
	fmt.Println("问题: \"最近发生了什么大事\"")
	
	brainInstance := brain.NewBrain()
	
	learningConfig := brain.LearningConfig{
		MaxItemsPerLearn:   8,
		EnabledTypes: []resource.ResourceType{resource.Social, resource.Tech},
		EnableAutoLearn:    false,
		ResourceConfigPath: "asset/resources.json",
	}
	
	err := brainInstance.InitLearningEngine(learningConfig)
	if err != nil {
		log.Fatal(err)
	}
	
	fmt.Println("开始学习...")
	ctx := context.Background()
	result, err := brainInstance.GetLearningEngine().LearnOnce(ctx)
	if err != nil {
		log.Fatal(err)
	}
	
	fmt.Printf("学习完成: %d条信息\n", result.ProcessedCount)
	
	kb := brainInstance.GetLearningEngine().GetKnowledgeBase()
	hotItems := kb.GetHotItems(5)
	
	fmt.Println("\nAI回答:")
	for i, item := range hotItems {
		fmt.Printf("%d. %s\n", i+1, item.Title)
		fmt.Printf("   %s\n", item.Content[:min(80, len(item.Content))])
		fmt.Printf("   热度: %d\n", item.HotValue)
	}
	
	fmt.Println("\n验证成功!")
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}