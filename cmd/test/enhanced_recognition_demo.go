package main

import (
	"context"
	"fmt"
	"strings"

	"new-robot-self/internal/intelligent"
)

func main() {
	fmt.Println("=== 增强意图识别测试 ===")
	fmt.Println()

	ctx := context.Background()

	// 创建意图识别器
	recognizer := intelligent.NewIntentRecognizer()

	// 测试1: 原始复合场景
	fmt.Println("【测试1】原始复合场景")
	input1 := "明天晚上我想要去铜山万达看电影,而且我现在需要去徐州市的苏宁购物"
	fmt.Printf("输入: %s\n", input1)
	
	intents1, err := recognizer.DetectCompositeIntent(ctx, input1)
	if err != nil {
		fmt.Printf("❌ 复合意图检测失败: %v\n", err)
	} else {
		fmt.Printf("✅ 检测到 %d 个意图\n", len(intents1))
		
		// 按优先级排序
		sortedIntents := recognizer.SortIntentsByPriority(intents1)
		
		for i, intent := range sortedIntents {
			fmt.Printf("\n意图 %d:\n", i+1)
			fmt.Printf("  类型: %s\n", intent.Type)
			fmt.Printf("  置信度: %.2f\n", intent.Confidence)
			fmt.Printf("  优先级: %s\n", intent.Priority)
			fmt.Printf("  实体: %v\n", intent.Entities)
		}
	}

	fmt.Println("\n" + strings.Repeat("=", 50) + "\n")

	// 测试2: 时间紧迫性
	fmt.Println("【测试2】时间紧迫性测试")
	testCases := []struct {
		name  string
		input string
	}{
		{"紧急需求", "我现在需要去徐州市的苏宁购物"},
		{"普通需求", "明天晚上我想要去铜山万达看电影"},
		{"立即行动", "马上帮我查一下天气"},
		{"预约需求", "下周三下午2点开会"},
	}

	for _, tc := range testCases {
		fmt.Printf("\n%s: %s\n", tc.name, tc.input)
		intent, err := recognizer.Recognize(ctx, tc.input)
		if err != nil {
			fmt.Printf("❌ 识别失败: %v\n", err)
		} else {
			urgency := "普通"
			if intent.Priority == "urgent" {
				urgency = "🚨 紧急"
			}
			fmt.Printf("✅ 类型: %s, 优先级: %s %s\n", intent.Type, intent.Priority, urgency)
		}
	}

	fmt.Println("\n" + strings.Repeat("=", 50) + "\n")

	// 测试3: 复合意图拆分
	fmt.Println("【测试3】复合意图拆分测试")
	compositeCases := []struct {
		name  string
		input string
	}{
		{"而且连接", "我要去超市买菜，而且我想去银行取钱"},
		{"同时连接", "今天下午3点开会，同时明天上午10点健身"},
		{"另外连接", "帮我查一下天气，另外告诉我明天的日程"},
	}

	for _, tc := range compositeCases {
		fmt.Printf("\n%s: %s\n", tc.name, tc.input)
		intents, err := recognizer.DetectCompositeIntent(ctx, tc.input)
		if err != nil {
			fmt.Printf("❌ 拆分失败: %v\n", err)
		} else {
			fmt.Printf("✅ 拆分为 %d 个意图\n", len(intents))
			for i, intent := range intents {
				fmt.Printf("  %d. %s (优先级: %s)\n", i+1, intent.Type, intent.Priority)
			}
		}
	}

	fmt.Println("\n=== 测试完成 ===")
}