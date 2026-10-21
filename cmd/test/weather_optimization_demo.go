package main

import (
	"context"
	"fmt"

	"new-robot-self/internal/brain"
	"new-robot-self/internal/tools"
)

func main() {
	fmt.Println("=== 天气查询工具优化测试 ===")

	// 创建知识库
	knowledgeBase := brain.NewKnowledgeBase()

	// 注册所有工具
	fmt.Println("1. 注册所有工具...")
	if err := tools.RegisterAllTools(knowledgeBase); err != nil {
		fmt.Printf("注册工具失败: %v\n", err)
		return
	}
	fmt.Println("✓ 工具注册成功")

	// 获取天气查询工具
	fmt.Println("\n2. 获取天气查询工具...")
	weatherTool, exists := tools.Get("weather_query")
	if !exists {
		fmt.Println("✗ 天气查询工具不存在")
		return
	}
	fmt.Println("✓ 天气查询工具获取成功")

	// 测试实时天气
	fmt.Println("\n3. 测试实时天气...")
	ctx := context.Background()
	result, err := weatherTool.InvokableRun(ctx, `{"city":"北京","extensions":"base"}`)
	if err != nil {
		fmt.Printf("✗ 实时天气查询失败: %v\n", err)
	} else {
		fmt.Printf("✓ 实时天气查询成功，结果长度: %d字符\n", len(result))
		fmt.Printf("结果预览: %.200s...\n", result)
	}

	// 测试多日预报
	fmt.Println("\n4. 测试多日预报...")
	result, err = weatherTool.InvokableRun(ctx, `{"city":"上海","extensions":"forecast","days":3}`)
	if err != nil {
		fmt.Printf("✗ 多日预报失败: %v\n", err)
	} else {
		fmt.Printf("✓ 多日预报成功，结果长度: %d字符\n", len(result))
		fmt.Printf("结果预览: %.200s...\n", result)
	}

	// 测试空气质量
	fmt.Println("\n5. 测试空气质量...")
	result, err = weatherTool.InvokableRun(ctx, `{"city":"广州","extensions":"air"}`)
	if err != nil {
		fmt.Printf("✗ 空气质量查询失败: %v\n", err)
	} else {
		fmt.Printf("✓ 空气质量查询成功，结果长度: %d字符\n", len(result))
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
