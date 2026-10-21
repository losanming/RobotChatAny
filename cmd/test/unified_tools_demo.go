package main

import (
	"context"
	"fmt"

	"new-robot-self/internal/brain"
	"new-robot-self/internal/tools"
)

func main() {
	fmt.Println("=== 统一Tools系统测试 ===")

	// 创建知识库
	knowledgeBase := brain.NewKnowledgeBase()

	// 注册所有工具
	fmt.Println("1. 注册所有工具...")
	if err := tools.RegisterAllTools(knowledgeBase); err != nil {
		fmt.Printf("注册工具失败: %v\n", err)
		return
	}
	fmt.Println("✓ 工具注册成功")

	// 列出所有工具
	fmt.Println("\n2. 列出所有工具...")
	toolNames := tools.List()
	fmt.Printf("工具列表 (%d个):\n", len(toolNames))
	for _, name := range toolNames {
		fmt.Printf("  - %s\n", name)
	}

	// 按分类列出工具
	fmt.Println("\n3. 按分类列出工具...")
	categories := []string{"location", "health", "info_fetcher", "learning"}
	for _, category := range categories {
		toolsInCategory := tools.ListByCategory(category)
		fmt.Printf("  %s: %d个工具\n", category, len(toolsInCategory))
		for _, name := range toolsInCategory {
			fmt.Printf("    - %s\n", name)
		}
	}

	// 获取工具信息
	fmt.Println("\n4. 获取工具信息...")
	ctx := context.Background()
	for _, name := range toolNames {
		info, err := tools.GetToolInfo(ctx, name)
		if err != nil {
			fmt.Printf("  ✗ %s: %v\n", name, err)
			continue
		}
		fmt.Printf("  ✓ %s: %s\n", info.Name, info.Desc)
	}

	// 测试RSS热点工具
	fmt.Println("\n5. 测试RSS热点工具...")
	rssTool, exists := tools.Get("rss_hot_news")
	if !exists {
		fmt.Println("  ✗ RSS热点工具不存在")
	} else {
		result, err := rssTool.InvokableRun(ctx, `{"category":"tech","limit":10}`)
		if err != nil {
			fmt.Printf("  ✗ RSS热点工具调用失败: %v\n", err)
		} else {
			fmt.Printf("  ✓ RSS热点工具调用成功，结果长度: %d字符\n", len(result))
		}
	}

	// 测试知识库
	fmt.Println("\n6. 检查知识库...")
	stats := knowledgeBase.GetStats()
	fmt.Printf("  知识库统计:\n")
	fmt.Printf("    总条目数: %d\n", stats.TotalItems)
	fmt.Printf("    按分类统计:\n")
	for category, count := range stats.ByCategory {
		fmt.Printf("      %s: %d\n", category, count)
	}

	// 转换为Eino工具
	fmt.Println("\n7. 转换为Eino工具...")
	einoTools := tools.ToEinoTools()
	fmt.Printf("  Eino工具数量: %d\n", len(einoTools))

	fmt.Println("\n=== 测试完成 ===")
}
