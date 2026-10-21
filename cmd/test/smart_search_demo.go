package main

import (
	"context"
	"fmt"

	"new-robot-self/internal/brain"
	"new-robot-self/internal/tools"
)

func main() {
	fmt.Println("=== 智能搜索工具测试 ===")

	// 创建知识库
	knowledgeBase := brain.NewKnowledgeBase()

	// 注册所有工具
	fmt.Println("1. 注册所有工具...")
	if err := tools.RegisterAllTools(knowledgeBase); err != nil {
		fmt.Printf("注册工具失败: %v\n", err)
		return
	}
	fmt.Println("✓ 工具注册成功")

	// 获取智能搜索工具
	fmt.Println("\n2. 获取智能搜索工具...")
	searchTool, exists := tools.Get("smart_search")
	if !exists {
		fmt.Println("✗ 智能搜索工具不存在")
		return
	}
	fmt.Println("✓ 智能搜索工具获取成功")

	// 测试Bing搜索
	fmt.Println("\n3. 测试Bing搜索...")
	ctx := context.Background()
	result, err := searchTool.InvokableRun(ctx, `{"query":"人工智能最新发展","engine":"bing","max_results":5}`)
	if err != nil {
		fmt.Printf("✗ Bing搜索失败: %v\n", err)
	} else {
		fmt.Printf("✓ Bing搜索成功，结果长度: %d字符\n", len(result))
		fmt.Printf("结果预览: %.100s...\n", result)
	}

	// 测试DuckDuckGo搜索
	fmt.Println("\n4. 测试DuckDuckGo搜索...")
	result, err = searchTool.InvokableRun(ctx, `{"query":"Go语言编程","engine":"duckduckgo","max_results":3}`)
	if err != nil {
		fmt.Printf("✗ DuckDuckGo搜索失败: %v\n", err)
	} else {
		fmt.Printf("✓ DuckDuckGo搜索成功，结果长度: %d字符\n", len(result))
		fmt.Printf("结果预览: %.100s...\n", result)
	}

	// 检查知识库
	fmt.Println("\n5. 检查知识库...")
	stats := knowledgeBase.GetStats()
	fmt.Printf("  知识库统计:\n")
	fmt.Printf("    总条目数: %d\n", stats.TotalItems)
	fmt.Printf("    按分类统计:\n")
	for category, count := range stats.ByCategory {
		fmt.Printf("      %s: %d\n", category, count)
	}

	// 列出所有工具
	fmt.Println("\n6. 列出所有工具...")
	toolNames := tools.List()
	fmt.Printf("工具总数: %d\n", len(toolNames))
	for _, name := range toolNames {
		fmt.Printf("  - %s\n", name)
	}

	fmt.Println("\n=== 测试完成 ===")
}
