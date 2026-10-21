package main

import (
	"context"
	"fmt"

	"new-robot-self/internal/intelligent"
)

func main() {
	fmt.Println("=== 简单智能助手测试 ===")
	fmt.Println()

	ctx := context.Background()

	// 测试意图识别
	fmt.Println("【测试1】意图识别")
	recognizer := intelligent.NewIntentRecognizer()
	input := "帮我安排明天上午9点的会议"
	intent, err := recognizer.Recognize(ctx, input)
	if err != nil {
		fmt.Printf("❌ 意图识别失败: %v\n", err)
	} else {
		fmt.Printf("✅ 意图识别成功\n")
		fmt.Printf("   类型: %s\n", intent.Type)
		fmt.Printf("   置信度: %.2f\n", intent.Confidence)
		fmt.Printf("   实体: %v\n", intent.Entities)
	}

	fmt.Println("\n=== 测试完成 ===")
}