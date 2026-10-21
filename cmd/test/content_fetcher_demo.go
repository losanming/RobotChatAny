package main

import (
	"context"
	"fmt"

	"new-robot-self/internal/brain"
	"new-robot-self/internal/tools"
)

func main() {
	fmt.Println("=== 内容获取工具测试 ===")

	// 创建知识库
	knowledgeBase := brain.NewKnowledgeBase()

	// 注册所有工具
	fmt.Println("1. 注册所有工具...")
	if err := tools.RegisterAllTools(knowledgeBase); err != nil {
		fmt.Printf("注册工具失败: %v\n", err)
		return
	}
	fmt.Println("✓ 工具注册成功")

	// 获取内容获取工具
	fmt.Println("\n2. 获取内容获取工具...")
	fetcherTool, exists := tools.Get("content_fetcher")
	if !exists {
		fmt.Println("✗ 内容获取工具不存在")
		return
	}
	fmt.Println("✓ 内容获取工具获取成功")

	// 测试内容获取 (使用一个简单的网页)
	fmt.Println("\n3. 测试内容获取...")
	ctx := context.Background()
	result, err := fetcherTool.InvokableRun(ctx, `{"url":"https://example.com","extract_images":false,"max_length":1000}`)
	if err != nil {
		fmt.Printf("✗ 内容获取失败: %v\n", err)
	} else {
		fmt.Printf("✓ 内容获取成功，结果长度: %d字符\n", len(result))
		fmt.Printf("结果预览: %.200s...\n", result)
	}

	// 测试带有图片的内容获取
	fmt.Println("\n4. 测试带图片的内容获取...")
	result, err = fetcherTool.InvokableRun(ctx, `{"url":"https://example.com","extract_images":true,"max_length":500}`)
	if err != nil {
		fmt.Printf("✗ 带图片内容获取失败: %v\n", err)
	} else {
		fmt.Printf("✓ 带图片内容获取成功，结果长度: %d字符\n", len(result))
		fmt.Printf("结果预览: %.200s...\n", result)
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
