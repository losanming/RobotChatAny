package main

import (
	"context"
	"fmt"

	"new-robot-self/internal/intelligent"
)

func main() {
	fmt.Println("=== 智能全能助手真实功能对接测试 ===")
	fmt.Println()

	// 创建智能助手
	assistant := intelligent.NewIntelligentAssistant()
	
	// 启动助手
	err := assistant.Start()
	if err != nil {
		fmt.Printf("启动助手失败: %v\n", err)
		return
	}
	defer assistant.Stop()

	ctx := context.Background()

	// 测试真实功能对接
	testScenario1(ctx, assistant)
	testScenario2(ctx, assistant)
	testScenario3(ctx, assistant)
	testScenario4(ctx, assistant)
	testScenario5(ctx, assistant) // 测试用户偏好学习
	testScenario6(ctx, assistant) // 测试主动触发

	// 显示上下文摘要
	fmt.Println("\n=== 上下文摘要 ===")
	summary := assistant.GetContextSummary()
	fmt.Printf("最近消息数: %v\n", summary["recent_messages_count"])
	fmt.Printf("最近活动数: %v\n", summary["recent_activities_count"])
	fmt.Printf("偏好设置数: %v\n", summary["preferences_count"])
	fmt.Printf("习惯记录数: %v\n", summary["habits_count"])

	// 显示最近消息
	fmt.Println("\n=== 最近消息 ===")
	messages := assistant.GetRecentMessages(5)
	for i, msg := range messages {
		fmt.Printf("%d. [%s] %s (置信度: %.2f)\n", i+1, msg.Role, msg.Content, msg.Intent.Confidence)
	}

	// 显示学习到的偏好
	fmt.Println("\n=== 学习到的用户偏好 ===")
	if preferredTime, ok := assistant.GetContextMemory().GetPreference("preferred_time"); ok {
		fmt.Printf("✓ 时间偏好: %v\n", preferredTime)
	}
	if preferredActivity, ok := assistant.GetContextMemory().GetPreference("preferred_activity"); ok {
		fmt.Printf("✓ 活动偏好: %v\n", preferredActivity)
	}

	fmt.Println("\n=== 测试完成 ===")
}

// 场景1: 日程安排 + 主动提醒
func testScenario1(ctx context.Context, assistant *intelligent.IntelligentAssistant) {
	fmt.Println("【场景1】日程安排 + 主动提醒")
	fmt.Println("用户输入: 我今日下午三点打算去公园玩")

	input := "我今日下午三点打算去公园玩"
	response, err := assistant.Process(ctx, input)
	if err != nil {
		fmt.Printf("❌ 处理失败: %v\n", err)
	} else {
		fmt.Printf("✅ 助手回复: %s\n", response.Text)
		fmt.Printf("🎯 识别意图: %s (置信度: %.2f)\n", response.Intent.Type, response.Confidence)
		fmt.Printf("🔧 建议行动: %v\n", response.Actions)
	}
	fmt.Println()
}

// 场景2: 购物 + 导航 + 周边查询
func testScenario2(ctx context.Context, assistant *intelligent.IntelligentAssistant) {
	fmt.Println("【场景2】购物 + 导航 + 周边查询")
	fmt.Println("用户输入: 我打算去市场买菜")

	input := "我打算去市场买菜"
	response, err := assistant.Process(ctx, input)
	if err != nil {
		fmt.Printf("❌ 处理失败: %v\n", err)
	} else {
		fmt.Printf("✅ 助手回复: %s\n", response.Text)
		fmt.Printf("🎯 识别意图: %s (置信度: %.2f)\n", response.Intent.Type, response.Confidence)
		fmt.Printf("🔧 建议行动: %v\n", response.Actions)
	}
	fmt.Println()
}

// 场景3: 健康咨询 + 饮食建议
func testScenario3(ctx context.Context, assistant *intelligent.IntelligentAssistant) {
	fmt.Println("【场景3】健康咨询 + 饮食建议")
	fmt.Println("用户输入: 我中午吃了米饭，想知道多少卡路里")

	input := "我中午吃了米饭，想知道多少卡路里"
	response, err := assistant.Process(ctx, input)
	if err != nil {
		fmt.Printf("❌ 处理失败: %v\n", err)
	} else {
		fmt.Printf("✅ 助手回复: %s\n", response.Text)
		fmt.Printf("🎯 识别意图: %s (置信度: %.2f)\n", response.Intent.Type, response.Confidence)
		fmt.Printf("🔧 建议行动: %v\n", response.Actions)
	}
	fmt.Println()
}

// 场景4: 工作需求 + 代码查询
func testScenario4(ctx context.Context, assistant *intelligent.IntelligentAssistant) {
	fmt.Println("【场景4】工作需求 + 代码查询")
	fmt.Println("用户输入: Go语言怎么连接数据库")

	input := "Go语言怎么连接数据库"
	response, err := assistant.Process(ctx, input)
	if err != nil {
		fmt.Printf("❌ 处理失败: %v\n", err)
	} else {
		fmt.Printf("✅ 助手回复: %s\n", response.Text)
		fmt.Printf("🎯 识别意图: %s (置信度: %.2f)\n", response.Intent.Type, response.Confidence)
		fmt.Printf("🔧 建议行动: %v\n", response.Actions)
		fmt.Printf("🔍 使用真实工具: %v\n", response.ToolResults)
	}
	fmt.Println()
}

// 场景5: 用户偏好学习测试
func testScenario5(ctx context.Context, assistant *intelligent.IntelligentAssistant) {
	fmt.Println("【场景5】用户偏好学习测试")
	fmt.Println("模拟多次相似交互以学习用户偏好")

	inputs := []string{
		"帮我安排明天上午9点的晨跑",
		"帮我安排明天上午9点的会议",
		"帮我安排明天上午9点的健身",
	}

	for _, input := range inputs {
		response, err := assistant.Process(ctx, input)
		if err != nil {
			fmt.Printf("❌ 处理失败: %v\n", err)
		} else {
			fmt.Printf("✅ 处理成功: %s (意图: %s)\n", input, response.Intent.Type)
		}
	}
	
	fmt.Println("✓ 用户偏好学习完成")
	fmt.Println()
}

// 场景6: 主动触发测试
func testScenario6(ctx context.Context, assistant *intelligent.IntelligentAssistant) {
	fmt.Println("【场景6】主动触发测试")
	fmt.Println("测试触发管理器的主动提醒功能")

	// 获取触发管理器
	triggerManager := assistant.GetTriggerManager()
	
	// 检查触发器状态
	rules := triggerManager.ListRules()
	fmt.Printf("✓ 当前有 %d 个触发规则\n", len(rules))
	
	// 检查触发器是否正在运行
	if triggerManager.IsRunning() {
		fmt.Println("✓ 触发管理器正在运行")
	} else {
		fmt.Println("❌ 触发管理器未运行")
	}
	
	fmt.Println("✓ 主动触发测试完成")
	fmt.Println()
}