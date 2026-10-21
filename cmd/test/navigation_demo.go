package main

import (
	"context"
	"fmt"

	"new-robot-self/internal/brain"
	"new-robot-self/internal/tools"
)

func main() {
	fmt.Println("=== 导航路径工具测试 ===")

	// 创建知识库
	knowledgeBase := brain.NewKnowledgeBase()

	// 注册所有工具
	fmt.Println("1. 注册所有工具...")
	if err := tools.RegisterAllTools(knowledgeBase); err != nil {
		fmt.Printf("注册工具失败: %v\n", err)
		return
	}
	fmt.Println("✓ 工具注册成功")

	// 获取导航工具
	fmt.Println("\n2. 获取导航工具...")
	navTool, exists := tools.Get("navigation")
	if !exists {
		fmt.Println("✗ 导航工具不存在")
		return
	}
	fmt.Println("✓ 导航工具获取成功")

	// 测试驾车导航
	fmt.Println("\n3. 测试驾车导航...")
	ctx := context.Background()
	result, err := navTool.InvokableRun(ctx, `{"origin":"116.481028,39.989643","destination":"116.327421,39.903738","strategy":"fast","mode":"driving"}`)
	if err != nil {
		fmt.Printf("✗ 驾车导航失败: %v\n", err)
	} else {
		fmt.Printf("✓ 驾车导航成功，结果长度: %d字符\n", len(result))
		fmt.Printf("结果预览: %.200s...\n", result)
	}

	// 测试步行导航
	fmt.Println("\n4. 测试步行导航...")
	result, err = navTool.InvokableRun(ctx, `{"origin":"116.481028,39.989643","destination":"116.327421,39.903738","strategy":"short","mode":"walking"}`)
	if err != nil {
		fmt.Printf("✗ 步行导航失败: %v\n", err)
	} else {
		fmt.Printf("✓ 步行导航成功，结果长度: %d字符\n", len(result))
		fmt.Printf("结果预览: %.200s...\n", result)
	}

	// 测试经济策略
	fmt.Println("\n5. 测试经济策略...")
	result, err = navTool.InvokableRun(ctx, `{"origin":"116.481028,39.989643","destination":"116.327421,39.903738","strategy":"economy","mode":"driving"}`)
	if err != nil {
		fmt.Printf("✗ 经济策略失败: %v\n", err)
	} else {
		fmt.Printf("✓ 经济策略成功，结果长度: %d字符\n", len(result))
		fmt.Printf("结果预览: %.200s...\n", result)
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
