package main

import (
	"context"
	"fmt"

	"new-robot-self/internal/intelligent"
)

func main() {
	fmt.Println("=== 复杂场景测试 ===")
	fmt.Println()

	ctx := context.Background()

	// 创建意图识别器
	recognizer := intelligent.NewIntentRecognizer()

	// 测试复杂场景
	input := "明天晚上我想要去铜山万达看电影,而且我现在需要去徐州市的苏宁购物"
	
	fmt.Printf("用户输入: %s\n", input)
	fmt.Println()

	intent, err := recognizer.Recognize(ctx, input)
	if err != nil {
		fmt.Printf("❌ 意图识别失败: %v\n", err)
		return
	}

	fmt.Printf("✅ 意图识别成功\n")
	fmt.Printf("🎯 类型: %s\n", intent.Type)
	fmt.Printf("📊 置信度: %.2f\n", intent.Confidence)
	fmt.Printf("🔍 实体信息: %v\n", intent.Entities)
	fmt.Printf("⚡ 优先级: %s\n", intent.Priority)
	fmt.Printf("📝 上下文: %s\n", intent.Context)

	// 分析识别结果
	fmt.Println("\n=== 识别结果分析 ===")
	
	// 检查是否识别出多个意图
	if len(intent.Entities) > 0 {
		fmt.Println("识别到的关键信息:")
		for key, value := range intent.Entities {
			fmt.Printf("  • %s: %v\n", key, value)
		}
	}

	// 根据识别结果给出建议
	fmt.Println("\n=== 智能建议 ===")
	
	// 检查是否有时间信息
	if time, ok := intent.Entities["time"].(string); ok {
		fmt.Printf("⏰ 时间识别: %s\n", time)
		if time == "现在" || time == "今天" {
			fmt.Println("   → 这是紧急需求，应该优先处理")
		}
	}

	// 检查是否有地点信息
	if location, ok := intent.Entities["location"].(string); ok {
		fmt.Printf("📍 地点识别: %s\n", location)
	}

	// 检查是否有活动信息
	if event, ok := intent.Entities["event"].(string); ok {
		fmt.Printf("🎬 活动识别: %s\n", event)
	}

	fmt.Println("\n=== 建议的处理顺序 ===")
	fmt.Println("1. 优先处理：现在去徐州市苏宁购物（紧急需求）")
	fmt.Println("2. 后续安排：明天晚上去铜山万达看电影（预约提醒）")

	fmt.Println("\n=== 测试完成 ===")
}