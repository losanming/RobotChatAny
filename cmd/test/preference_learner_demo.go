package main

import (
	"context"
	"fmt"

	"new-robot-self/internal/brain"
	"new-robot-self/internal/tools"
)

func main() {
	fmt.Println("=== 偏好学习工具测试 ===")

	// 创建知识库
	knowledgeBase := brain.NewKnowledgeBase()

	// 注册所有工具
	fmt.Println("1. 注册所有工具...")
	if err := tools.RegisterAllTools(knowledgeBase); err != nil {
		fmt.Printf("注册工具失败: %v\n", err)
		return
	}
	fmt.Println("✓ 工具注册成功")

	// 获取偏好学习工具
	fmt.Println("\n2. 获取偏好学习工具...")
	preferenceTool, exists := tools.Get("preference_learner")
	if !exists {
		fmt.Println("✗ 偏好学习工具不存在")
		return
	}
	fmt.Println("✓ 偏好学习工具获取成功")

	// 测试分析偏好
	fmt.Println("\n3. 测试分析偏好...")
	ctx := context.Background()
	result, err := preferenceTool.InvokableRun(ctx, `{"action":"analyze"}`)
	if err != nil {
		fmt.Printf("✗ 分析偏好失败: %v\n", err)
	} else {
		fmt.Printf("✓ 分析偏好成功，结果: %s\n", result)
	}

	// 测试更新偏好
	fmt.Println("\n4. 测试更新偏好...")
	result, err = preferenceTool.InvokableRun(ctx, `{"action":"update","category":"tech","user_action":"like","duration":120}`)
	if err != nil {
		fmt.Printf("✗ 更新偏好失败: %v\n", err)
	} else {
		fmt.Printf("✓ 更新偏好成功，结果: %s\n", result)
	}

	// 测试获取偏好
	fmt.Println("\n5. 测试获取偏好...")
	result, err = preferenceTool.InvokableRun(ctx, `{"action":"get"}`)
	if err != nil {
		fmt.Printf("✗ 获取偏好失败: %v\n", err)
	} else {
		fmt.Printf("✓ 获取偏好成功，结果: %s\n", result)
	}

	// 测试分享行为更新
	fmt.Println("\n6. 测试分享行为更新...")
	result, err = preferenceTool.InvokableRun(ctx, `{"action":"update","category":"health","user_action":"share","duration":300}`)
	if err != nil {
		fmt.Printf("✗ 分享行为更新失败: %v\n", err)
	} else {
		fmt.Printf("✓ 分享行为更新成功，结果: %s\n", result)
	}

	// 再次获取偏好查看变化
	fmt.Println("\n7. 再次获取偏好...")
	result, err = preferenceTool.InvokableRun(ctx, `{"action":"get"}`)
	if err != nil {
		fmt.Printf("✗ 获取偏好失败: %v\n", err)
	} else {
		fmt.Printf("✓ 获取偏好成功，结果: %s\n", result)
	}

	// 列出所有工具
	fmt.Println("\n8. 列出所有工具...")
	toolNames := tools.List()
	fmt.Printf("工具总数: %d\n", len(toolNames))
	for _, name := range toolNames {
		fmt.Printf("  - %s\n", name)
	}

	fmt.Println("\n=== 测试完成 ===")
}
