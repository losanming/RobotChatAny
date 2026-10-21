package main

import (
	"context"
	"fmt"

	"new-robot-self/internal/intelligent"
)

func main() {
	fmt.Println("=== 时间优先级测试 ===")
	fmt.Println()

	ctx := context.Background()

	// 创建意图识别器
	recognizer := intelligent.NewIntentRecognizer()

	// 测试不同时间优先级的场景
	testCases := []struct {
		name  string
		input string
	}{
		{"当前紧急需求", "我现在需要去徐州市的苏宁购物"},
		{"未来预约需求", "明天晚上我想要去铜山万达看电影"},
		{"复合场景", "明天晚上我想要去铜山万达看电影,而且我现在需要去徐州市的苏宁购物"},
		{"多个时间点", "今天下午3点开会，明天上午10点健身"},
	}

	for _, tc := range testCases {
		fmt.Printf("【%s】\n", tc.name)
		fmt.Printf("输入: %s\n", tc.input)
		
		intent, err := recognizer.Recognize(ctx, tc.input)
		if err != nil {
			fmt.Printf("❌ 意图识别失败: %v\n", err)
		} else {
			fmt.Printf("✅ 类型: %s, 置信度: %.2f, 优先级: %s\n", 
				intent.Type, intent.Confidence, intent.Priority)
			fmt.Printf("   时间: %v\n", intent.Entities["time"])
			fmt.Printf("   地点: %v\n", intent.Entities["location"])
			fmt.Printf("   活动: %v\n", intent.Entities["event"])
		}
		fmt.Println()
	}

	fmt.Println("=== 测试完成 ===")
}