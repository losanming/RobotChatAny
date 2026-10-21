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
	fmt.Println("========== 搜索 danfwing ==========")

	// 初始化工具系统
	ctx := context.Background()

	// 创建知识库
	knowledgeBase := brain.NewKnowledgeBase()

	// 注册所有工具
	err := tools.RegisterAllTools(knowledgeBase)
	if err != nil {
		log.Fatal("工具注册失败:", err)
	}

	fmt.Println("✅ 工具系统初始化成功")

	// 列出所有可用工具
	toolNames := tools.List()
	fmt.Printf("📋 可用工具数量: %d\n", len(toolNames))
	fmt.Printf("🔧 工具列表: %v\n\n", toolNames)

	// 获取智能搜索工具
	searchTool, exists := tools.Get("smart_search")
	if !exists {
		log.Fatal("智能搜索工具不存在")
	}

	fmt.Println("✅ 智能搜索工具获取成功")

	// 构造搜索参数
	searchArgs := map[string]interface{}{
		"query":       "danfwing",
		"engine":      "bing",
		"max_results": 10,
		"language":    "zh-CN",
	}

	argsJSON, err := json.Marshal(searchArgs)
	if err != nil {
		log.Fatal("参数序列化失败:", err)
	}

	fmt.Printf("🔍 搜索参数: %s\n\n", string(argsJSON))

	// 调用搜索工具
	result, err := searchTool.InvokableRun(ctx, string(argsJSON))
	if err != nil {
		log.Fatal("搜索失败:", err)
	}

	fmt.Println("📊 搜索结果:")
	fmt.Println(result)
	fmt.Println("\n========== 搜索完成 ==========")
}