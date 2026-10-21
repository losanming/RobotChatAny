package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"new-robot-self/internal/brain"
	"new-robot-self/internal/tools"
	"time"
)

func main() {
	fmt.Println("========== 测试 time_manager 工具 ==========")

	// 初始化工具系统
	ctx := context.Background()
	knowledgeBase := brain.NewKnowledgeBase()

	err := tools.RegisterAllTools(knowledgeBase)
	if err != nil {
		log.Fatal("工具注册失败:", err)
	}

	fmt.Println("✅ 工具系统初始化成功")

	// 获取时间管理工具
	timeTool, exists := tools.Get("time_manager")
	if !exists {
		log.Fatal("时间管理工具不存在")
	}

	fmt.Println("✅ 时间管理工具获取成功\n")

	// 测试1: 开始番茄钟
	fmt.Println("⏱️ 测试1: 开始番茄钟")
	startArgs := map[string]interface{}{
		"action":    "start_pomodoro",
		"task":      "代码编写",
		"duration":  1500.0, // 25分钟
		"break_type": "short",
	}

	argsJSON, _ := json.Marshal(startArgs)
	result, err := timeTool.InvokableRun(ctx, string(argsJSON))
	if err != nil {
		log.Fatal("开始番茄钟失败:", err)
	}

	fmt.Printf("✅ 番茄钟启动结果: %s\n\n", result)

	// 测试2: 获取番茄钟状态
	fmt.Println("⏱️ 测试2: 获取番茄钟状态")
	statusArgs := map[string]interface{}{
		"action": "get_pomodoro_status",
	}

	argsJSON, _ = json.Marshal(statusArgs)
	result, err = timeTool.InvokableRun(ctx, string(argsJSON))
	if err != nil {
		log.Fatal("获取番茄钟状态失败:", err)
	}

	fmt.Printf("✅ 番茄钟状态: %s\n\n", result)

	// 测试3: 暂停番茄钟
	fmt.Println("⏱️ 测试3: 暂停番茄钟")
	pauseArgs := map[string]interface{}{
		"action": "pause_pomodoro",
	}

	argsJSON, _ = json.Marshal(pauseArgs)
	result, err = timeTool.InvokableRun(ctx, string(argsJSON))
	if err != nil {
		log.Fatal("暂停番茄钟失败:", err)
	}

	fmt.Printf("✅ 暂停结果: %s\n\n", result)

	// 测试4: 恢复番茄钟
	fmt.Println("⏱️ 测试4: 恢复番茄钟")
	resumeArgs := map[string]interface{}{
		"action": "resume_pomodoro",
	}

	argsJSON, _ = json.Marshal(resumeArgs)
	result, err = timeTool.InvokableRun(ctx, string(argsJSON))
	if err != nil {
		log.Fatal("恢复番茄钟失败:", err)
	}

	fmt.Printf("✅ 恢复结果: %s\n\n", result)

	// 测试5: 停止番茄钟
	fmt.Println("⏱️ 测试5: 停止番茄钟")
	stopArgs := map[string]interface{}{
		"action": "stop_pomodoro",
	}

	argsJSON, _ = json.Marshal(stopArgs)
	result, err = timeTool.InvokableRun(ctx, string(argsJSON))
	if err != nil {
		log.Fatal("停止番茄钟失败:", err)
	}

	fmt.Printf("✅ 停止结果: %s\n\n", result)

	// 测试6: 记录时间
	fmt.Println("⏱️ 测试6: 记录时间")
	recordArgs := map[string]interface{}{
		"action":     "record_time",
		"task":       "代码审查",
		"category":   "work",
		"start_time": float64(time.Now().Add(-2 * time.Hour).Unix()),
		"end_time":   float64(time.Now().Add(-1 * time.Hour).Unix()),
		"is_break":   false,
	}

	argsJSON, _ = json.Marshal(recordArgs)
	result, err = timeTool.InvokableRun(ctx, string(argsJSON))
	if err != nil {
		log.Fatal("记录时间失败:", err)
	}

	fmt.Printf("✅ 记录结果: %s\n\n", result)

	// 测试7: 记录更多时间数据
	fmt.Println("⏱️ 测试7: 记录更多时间数据")
	moreRecords := []map[string]interface{}{
		{
			"action":     "record_time",
			"task":       "文档编写",
			"category":   "work",
			"start_time": float64(time.Now().Add(-4 * time.Hour).Unix()),
			"end_time":   float64(time.Now().Add(-3 * time.Hour).Unix()),
			"is_break":   false,
		},
		{
			"action":     "record_time",
			"task":       "休息",
			"category":   "break",
			"start_time": float64(time.Now().Add(-3 * time.Hour).Unix()),
			"end_time":   float64(time.Now().Add(-2 * time.Hour).Unix()),
			"is_break":   true,
		},
		{
			"action":     "record_time",
			"task":       "学习新技术",
			"category":   "learning",
			"start_time": float64(time.Now().Add(-6 * time.Hour).Unix()),
			"end_time":   float64(time.Now().Add(-5 * time.Hour).Unix()),
			"is_break":   false,
		},
	}

	for _, recordData := range moreRecords {
		argsJSON, _ = json.Marshal(recordData)
		result, err = timeTool.InvokableRun(ctx, string(argsJSON))
		if err != nil {
			log.Printf("记录时间失败: %v\n", err)
		} else {
			fmt.Printf("✅ 记录成功: %s\n", recordData["task"])
		}
	}
	fmt.Println()

	// 测试8: 获取时间统计
	fmt.Println("⏱️ 测试8: 获取时间统计")
	today := time.Now().Format("2006-01-02")
	statsArgs := map[string]interface{}{
		"action": "get_time_stats",
		"date":   today,
		"days":   7,
	}

	argsJSON, _ = json.Marshal(statsArgs)
	result, err = timeTool.InvokableRun(ctx, string(argsJSON))
	if err != nil {
		log.Fatal("获取时间统计失败:", err)
	}

	fmt.Printf("✅ 时间统计结果: %s\n\n", result)

	// 测试9: 获取效率报告
	fmt.Println("⏱️ 测试9: 获取效率报告")
	reportArgs := map[string]interface{}{
		"action": "get_efficiency_report",
		"date":   today,
		"days":   7,
	}

	argsJSON, _ = json.Marshal(reportArgs)
	result, err = timeTool.InvokableRun(ctx, string(argsJSON))
	if err != nil {
		log.Fatal("获取效率报告失败:", err)
	}

	fmt.Printf("✅ 效率报告结果: %s\n\n", result)

	fmt.Println("========== 测试完成 ==========")
}