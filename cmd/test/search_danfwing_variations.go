package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"new-robot-self/internal/brain"
	"new-robot-self/internal/tools"
)

func main() {
	fmt.Println("========== 多种方式搜索 danfwing ==========")

	// 初始化工具系统
	ctx := context.Background()
	knowledgeBase := brain.NewKnowledgeBase()

	err := tools.RegisterAllTools(knowledgeBase)
	if err != nil {
		log.Fatal("工具注册失败:", err)
	}

	fmt.Println("✅ 工具系统初始化成功\n")

	// 尝试不同的搜索策略
	searchStrategies := []map[string]interface{}{
		{
			"query":       "danfwing github",
			"engine":      "bing",
			"max_results": 5,
		},
		{
			"query":       "danfwing developer",
			"engine":      "bing",
			"max_results": 5,
		},
		{
			"query":       "danfwing 程序员",
			"engine":      "bing",
			"max_results": 5,
		},
	}

	searchTool, exists := tools.Get("smart_search")
	if !exists {
		log.Fatal("智能搜索工具不存在")
	}

	for i, strategy := range searchStrategies {
		fmt.Printf("🔍 搜索策略 %d: %s\n\n", i+1, strategy["query"])

		argsJSON, err := json.Marshal(strategy)
		if err != nil {
			log.Fatal("参数序列化失败:", err)
		}

		result, err := searchTool.InvokableRun(ctx, string(argsJSON))
		if err != nil {
			fmt.Printf("❌ 搜索失败: %v\n\n", err)
			continue
		}

		fmt.Printf("📊 搜索结果:\n%s\n\n", result)
		fmt.Println("---")
	}

	fmt.Println("========== 搜索完成 ==========")
}