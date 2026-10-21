package main

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"new-robot-self/internal/tools/health"
	"new-robot-self/internal/tools/smart_schedule"
	"new-robot-self/internal/tools/visualization"
	"new-robot-self/internal/tools/work"
)

func main() {
	fmt.Println("=== 威猛先生个人健康增强功能测试 ===")
	fmt.Println()

	// 测试健康数据追踪工具
	testHealthTracker()
	
	// 测试饮食管理工具
	testDietManager()
	
	// 测试智能提醒工具
	testSmartReminder()
	
	// 测试工作助手工具
	testWorkAssistant()
	
	// 测试数据可视化工具
	testDataVisualizer()
	
	// 测试智能日程工具
	testSmartScheduler()
	
	fmt.Println("=== 所有测试完成 ===")
}

func testHealthTracker() {
	fmt.Println("【1】测试健康数据追踪工具 (health_tracker)")
	
	tool := health.NewHealthTrackerTool()
	ctx := context.Background()
	
	// 测试记录血糖数据
	recordParams := map[string]interface{}{
		"action":      "record",
		"record_type": "blood_sugar",
		"value":       6.5,
		"unit":        "mmol/L",
		"notes":       "空腹血糖",
		"recorded_at": time.Now().Unix(),
	}
	
	paramsJSON, _ := json.Marshal(recordParams)
	result, err := tool.InvokableRun(ctx, string(paramsJSON))
	if err != nil {
		fmt.Printf("❌ 记录血糖数据失败: %v\n", err)
	} else {
		fmt.Printf("✅ 记录血糖数据成功: %s\n", result)
	}
	
	// 测试查询健康数据
	queryParams := map[string]interface{}{
		"action":     "query",
		"start_date": time.Now().AddDate(0, 0, -7).Unix(),
		"end_date":   time.Now().Unix(),
	}
	
	queryJSON, _ := json.Marshal(queryParams)
	_, err = tool.InvokableRun(ctx, string(queryJSON))
	if err != nil {
		fmt.Printf("❌ 查询健康数据失败: %v\n", err)
	} else {
		fmt.Printf("✅ 查询健康数据成功\n")
	}
	
	// 测试生成健康报告
	reportParams := map[string]interface{}{
		"action":     "report",
		"start_date": time.Now().AddDate(0, 0, -30).Unix(),
		"end_date":   time.Now().Unix(),
	}
	
	reportJSON, _ := json.Marshal(reportParams)
	_, err = tool.InvokableRun(ctx, string(reportJSON))
	if err != nil {
		fmt.Printf("❌ 生成健康报告失败: %v\n", err)
	} else {
		fmt.Printf("✅ 生成健康报告成功\n")
	}
	
	fmt.Println()
}

func testDietManager() {
	fmt.Println("【2】测试饮食管理工具 (diet_manager)")
	
	tool := health.NewDietManagerTool()
	ctx := context.Background()
	
	// 测试计算卡路里
	calcParams := map[string]interface{}{
		"action":       "calculate",
		"food_name":    "米饭",
		"amount":       100,
		"amount_unit":  "g",
	}
	
	paramsJSON, _ := json.Marshal(calcParams)
	result, err := tool.InvokableRun(ctx, string(paramsJSON))
	if err != nil {
		fmt.Printf("❌ 计算卡路里失败: %v\n", err)
	} else {
		fmt.Printf("✅ 计算卡路里成功: %s\n", result)
	}
	
	// 测试获取饮食建议
	suggestParams := map[string]interface{}{
		"action": "suggest",
	}
	
	suggestJSON, _ := json.Marshal(suggestParams)
	_, err = tool.InvokableRun(ctx, string(suggestJSON))
	if err != nil {
		fmt.Printf("❌ 获取饮食建议失败: %v\n", err)
	} else {
		fmt.Printf("✅ 获取饮食建议成功\n")
	}
	
	fmt.Println()
}

func testSmartReminder() {
	fmt.Println("【3】测试智能提醒工具 (smart_reminder)")
	
	tool := health.NewSmartReminderTool()
	ctx := context.Background()
	
	// 测试创建提醒
	createParams := map[string]interface{}{
		"action":        "create",
		"reminder_type": "exercise",
		"title":         "晚间运动",
		"description":   "根据运动循环进行相应运动",
		"repeat_rule":   "daily",
	}
	
	paramsJSON, _ := json.Marshal(createParams)
	result, err := tool.InvokableRun(ctx, string(paramsJSON))
	if err != nil {
		fmt.Printf("❌ 创建提醒失败: %v\n", err)
	} else {
		fmt.Printf("✅ 创建提醒成功: %s\n", result)
	}
	
	// 测试触发提醒检查
	triggerParams := map[string]interface{}{
		"action": "trigger",
	}
	
	triggerJSON, _ := json.Marshal(triggerParams)
	_, err = tool.InvokableRun(ctx, string(triggerJSON))
	if err != nil {
		fmt.Printf("❌ 触发提醒检查失败: %v\n", err)
	} else {
		fmt.Printf("✅ 触发提醒检查成功\n")
	}
	
	fmt.Println()
}

func testWorkAssistant() {
	fmt.Println("【4】测试工作助手工具 (work_assistant)")
	
	tool := work.NewWorkAssistantTool()
	ctx := context.Background()
	
	// 测试保存代码片段
	snippetParams := map[string]interface{}{
		"action":      "save_snippet",
		"language":    "go",
		"code":        "func main() {\n    fmt.Println(\"Hello, World!\")\n}",
		"description": "简单的Go程序",
		"tags":        []string{"go", "basic", "hello"},
	}
	
	paramsJSON, _ := json.Marshal(snippetParams)
	result, err := tool.InvokableRun(ctx, string(paramsJSON))
	if err != nil {
		fmt.Printf("❌ 保存代码片段失败: %v\n", err)
	} else {
		fmt.Printf("✅ 保存代码片段成功: %s\n", result)
	}
	
	// 测试搜索技术文档
	docParams := map[string]interface{}{
		"action":       "search_doc",
		"search_query": "Go语言教程",
		"doc_type":     "tutorial",
	}
	
	docJSON, _ := json.Marshal(docParams)
	_, err = tool.InvokableRun(ctx, string(docJSON))
	if err != nil {
		fmt.Printf("❌ 搜索技术文档失败: %v\n", err)
	} else {
		fmt.Printf("✅ 搜索技术文档成功\n")
	}
	
	fmt.Println()
}

func testDataVisualizer() {
	fmt.Println("【5】测试数据可视化工具 (data_visualizer)")
	
	tool := visualization.NewDataVisualizerTool()
	ctx := context.Background()
	
	// 测试生成血糖趋势图
	chartParams := map[string]interface{}{
		"action":     "generate_chart",
		"chart_type": "line",
		"data_type":  "blood_sugar",
		"start_date": time.Now().AddDate(0, 0, -30).Unix(),
		"end_date":   time.Now().Unix(),
		"resolution": "daily",
	}
	
	paramsJSON, _ := json.Marshal(chartParams)
	_, err := tool.InvokableRun(ctx, string(paramsJSON))
	if err != nil {
		fmt.Printf("❌ 生成血糖趋势图失败: %v\n", err)
	} else {
		fmt.Printf("✅ 生成血糖趋势图成功\n")
	}
	
	// 测试生成综合报告
	reportParams := map[string]interface{}{
		"action":     "generate_report",
		"start_date": time.Now().AddDate(0, 0, -30).Unix(),
		"end_date":   time.Now().Unix(),
	}
	
	reportJSON, _ := json.Marshal(reportParams)
	_, err = tool.InvokableRun(ctx, string(reportJSON))
	if err != nil {
		fmt.Printf("❌ 生成综合报告失败: %v\n", err)
	} else {
		fmt.Printf("✅ 生成综合报告成功\n")
	}
	
	fmt.Println()
}

func testSmartScheduler() {
	fmt.Println("【6】测试智能日程工具 (smart_scheduler)")
	
	tool := smart_schedule.NewSmartSchedulerTool()
	ctx := context.Background()
	
	// 测试创建任务
	createParams := map[string]interface{}{
		"action":         "create",
		"task_title":     "项目开发",
		"task_type":      "work",
		"duration":       240,
		"priority":       "high",
		"preferred_time": "上午",
	}
	
	paramsJSON, _ := json.Marshal(createParams)
	result, err := tool.InvokableRun(ctx, string(paramsJSON))
	if err != nil {
		fmt.Printf("❌ 创建任务失败: %v\n", err)
	} else {
		fmt.Printf("✅ 创建任务成功: %s\n", result)
	}
	
	// 测试优化日程
	optimizeParams := map[string]interface{}{
		"action":     "optimize",
		"start_date": time.Now().Unix(),
		"end_date":   time.Now().AddDate(0, 0, 7).Unix(),
	}
	
	optimizeJSON, _ := json.Marshal(optimizeParams)
	_, err = tool.InvokableRun(ctx, string(optimizeJSON))
	if err != nil {
		fmt.Printf("❌ 优化日程失败: %v\n", err)
	} else {
		fmt.Printf("✅ 优化日程成功\n")
	}
	
	fmt.Println()
}