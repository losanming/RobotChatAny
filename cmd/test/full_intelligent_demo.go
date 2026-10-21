package main

import (
	"context"
	"fmt"

	"new-robot-self/internal/intelligent"
)

func main() {
	fmt.Println("=== 智能意图识别与学习测试 ===")
	fmt.Println()

	ctx := context.Background()

	// 创建意图识别器
	recognizer := intelligent.NewIntentRecognizer()

	// 测试场景
	testScenarios := []struct {
		name  string
		input string
	}{
		{"日程安排", "帮我安排明天上午9点的会议"},
		{"导航查询", "导航到北京西站"},
		{"健康咨询", "我中午吃了100g米饭，想知道多少卡路里"},
		{"工作查询", "Go语言怎么连接数据库"},
		{"信息查询", "今天天气怎么样"},
		{"生活服务", "我想买一些水果"},
	}

	for _, scenario := range testScenarios {
		fmt.Printf("【%s】\n", scenario.name)
		fmt.Printf("用户输入: %s\n", scenario.input)
		
		intent, err := recognizer.Recognize(ctx, scenario.input)
		if err != nil {
			fmt.Printf("❌ 意图识别失败: %v\n", err)
		} else {
			fmt.Printf("✅ 意图识别成功\n")
			fmt.Printf("   类型: %s\n", intent.Type)
			fmt.Printf("   置信度: %.2f\n", intent.Confidence)
			fmt.Printf("   实体: %v\n", intent.Entities)
			fmt.Printf("   优先级: %s\n", intent.Priority)
		}
		fmt.Println()
	}

	fmt.Println("=== 测试完成 ===")
}