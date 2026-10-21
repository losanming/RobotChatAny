package main

import (
	"context"
	"fmt"

	"new-robot-self/internal/brain"
	"new-robot-self/internal/tools"
)

func main() {
	fmt.Println("=== 运动建议工具测试 ===")

	// 创建知识库
	knowledgeBase := brain.NewKnowledgeBase()

	// 注册所有工具
	fmt.Println("1. 注册所有工具...")
	if err := tools.RegisterAllTools(knowledgeBase); err != nil {
		fmt.Printf("注册工具失败: %v\n", err)
		return
	}
	fmt.Println("✓ 工具注册成功")

	// 获取运动建议工具
	fmt.Println("\n2. 获取运动建议工具...")
	exerciseTool, exists := tools.Get("exercise_advisor")
	if !exists {
		fmt.Println("✗ 运动建议工具不存在")
		return
	}
	fmt.Println("✓ 运动建议工具获取成功")

	// 测试制定减肥计划
	fmt.Println("\n3. 测试制定减肥计划...")
	ctx := context.Background()
	result, err := exerciseTool.InvokableRun(ctx, `{"action":"plan","weight":76.6,"height":165,"goal":"weight_loss"}`)
	if err != nil {
		fmt.Printf("✗ 制定减肥计划失败: %v\n", err)
	} else {
		fmt.Printf("✓ 制定减肥计划成功，结果长度: %d字符\n", len(result))
		fmt.Printf("结果预览: %.300s...\n", result)
	}

	// 测试获取跑步建议
	fmt.Println("\n4. 测试获取跑步建议...")
	result, err = exerciseTool.InvokableRun(ctx, `{"action":"advice","exercise_type":"running","duration":"60"}`)
	if err != nil {
		fmt.Printf("✗ 获取跑步建议失败: %v\n", err)
	} else {
		fmt.Printf("✓ 获取跑步建议成功，结果长度: %d字符\n", len(result))
		fmt.Printf("结果预览: %.200s...\n", result)
	}

	// 测试效果跟踪
	fmt.Println("\n5. 测试效果跟踪...")
	result, err = exerciseTool.InvokableRun(ctx, `{"action":"track","weight":76.6,"height":165}`)
	if err != nil {
		fmt.Printf("✗ 效果跟踪失败: %v\n", err)
	} else {
		fmt.Printf("✓ 效果跟踪成功，结果长度: %d字符\n", len(result))
		fmt.Printf("结果预览: %.200s...\n", result)
	}

	// 测试健身计划
	fmt.Println("\n6. 测试制定健身计划...")
	result, err = exerciseTool.InvokableRun(ctx, `{"action":"plan","weight":76.6,"height":165,"goal":"fitness"}`)
	if err != nil {
		fmt.Printf("✗ 制定健身计划失败: %v\n", err)
	} else {
		fmt.Printf("✓ 制定健身计划成功，结果长度: %d字符\n", len(result))
		fmt.Printf("结果预览: %.300s...\n", result)
	}

	// 列出所有工具
	fmt.Println("\n7. 列出所有工具...")
	toolNames := tools.List()
	fmt.Printf("工具总数: %d\n", len(toolNames))
	for _, name := range toolNames {
		fmt.Printf("  - %s\n", name)
	}

	fmt.Println("\n=== 测试完成 ===")
}
