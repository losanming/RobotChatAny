package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"new-robot-self/internal/brain"
	"new-robot-self/internal/tools"
)

func main() {
	fmt.Println("========== 测试 report_generator 工具 ==========")

	// 初始化工具系统
	ctx := context.Background()
	knowledgeBase := brain.NewKnowledgeBase()

	err := tools.RegisterAllTools(knowledgeBase)
	if err != nil {
		log.Fatal("工具注册失败:", err)
	}

	fmt.Println("✅ 工具系统初始化成功")

	// 获取报告生成工具
	reportTool, exists := tools.Get("report_generator")
	if !exists {
		log.Fatal("报告生成工具不存在")
	}

	fmt.Println("✅ 报告生成工具获取成功\n")

	// 测试1: 列出可用模板
	fmt.Println("📄 测试1: 列出可用模板")
	listArgs := map[string]interface{}{
		"action": "list_templates",
	}

	argsJSON, _ := json.Marshal(listArgs)
	result, err := reportTool.InvokableRun(ctx, string(argsJSON))
	if err != nil {
		log.Fatal("列出模板失败:", err)
	}

	fmt.Printf("✅ 模板列表: %s\n\n", result)

	// 测试2: 生成周报
	fmt.Println("📄 测试2: 生成周报")
	weeklyData := map[string]interface{}{
		"week_start":     "2026-05-11",
		"week_end":       "2026-05-18",
		"tasks_completed": "完成统一Tools系统重构",
		"tasks_in_progress": "实现信息分析工具",
		"issues":         "无重大问题",
		"next_week_plan": "继续实现P2阶段工具",
	}

	weeklyArgs := map[string]interface{}{
		"action":        "generate",
		"template":      "weekly_report",
		"title":         "周报 - 全能AI助手开发",
		"author":        "威猛先生",
		"output_format": "markdown",
		"data":          weeklyData,
		"include_charts": false,
	}

	argsJSON, _ = json.Marshal(weeklyArgs)
	result, err = reportTool.InvokableRun(ctx, string(argsJSON))
	if err != nil {
		log.Fatal("生成周报失败:", err)
	}

	fmt.Printf("✅ 周报生成结果: %s\n\n", result)

	// 测试3: 生成月报
	fmt.Println("📄 测试3: 生成月报")
	monthlyData := map[string]interface{}{
		"month":           "2026年05月",
		"key_achievements": "完成12个核心工具",
		"metrics":         "工具覆盖率100%",
		"challenges":      "时间管理优化",
		"next_month_plan": "实现高级功能工具",
	}

	monthlyArgs := map[string]interface{}{
		"action":        "generate",
		"template":      "monthly_report",
		"title":         "月报 - 工具系统开发",
		"author":        "威猛先生",
		"output_format": "html",
		"data":          monthlyData,
		"include_charts": true,
	}

	argsJSON, _ = json.Marshal(monthlyArgs)
	result, err = reportTool.InvokableRun(ctx, string(argsJSON))
	if err != nil {
		log.Fatal("生成月报失败:", err)
	}

	fmt.Printf("✅ 月报生成结果（前500字符）: %s...\n\n", result[:min(500, len(result))])

	// 测试4: 生成项目总结
	fmt.Println("📄 测试4: 生成项目总结")
	projectData := map[string]interface{}{
		"project_name": "全能AI助手工具系统",
		"status":       "进行中",
		"progress":     "60%",
		"milestones":   "完成P0和P1阶段",
		"risks":        "技术风险可控",
		"next_steps":   "开始P2阶段开发",
	}

	projectArgs := map[string]interface{}{
		"action":        "generate",
		"template":      "project_summary",
		"title":         "项目总结 - 工具系统扩展",
		"author":        "威猛先生",
		"output_format": "json",
		"data":          projectData,
		"include_charts": false,
	}

	argsJSON, _ = json.Marshal(projectArgs)
	result, err = reportTool.InvokableRun(ctx, string(argsJSON))
	if err != nil {
		log.Fatal("生成项目总结失败:", err)
	}

	fmt.Printf("✅ 项目总结结果: %s\n\n", result)

	// 测试5: 生成数据分析报告
	fmt.Println("📄 测试5: 生成数据分析报告")
	analysisData := map[string]interface{}{
		"data_source": "工具使用数据",
		"time_range":  "最近7天",
		"statistics": "总调用次数: 156次",
		"trends":     "使用量呈上升趋势",
		"insights":   "个人管理工具使用最频繁",
	}

	analysisArgs := map[string]interface{}{
		"action":        "generate",
		"template":      "data_analysis",
		"title":         "数据分析报告 - 工具使用情况",
		"author":        "威猛先生",
		"output_format": "markdown",
		"data":          analysisData,
		"include_charts": true,
	}

	argsJSON, _ = json.Marshal(analysisArgs)
	result, err = reportTool.InvokableRun(ctx, string(argsJSON))
	if err != nil {
		log.Fatal("生成数据分析报告失败:", err)
	}

	fmt.Printf("✅ 数据分析报告结果（前300字符）: %s...\n\n", result[:min(300, len(result))])

	// 测试6: 生成自定义报告
	fmt.Println("📄 测试6: 生成自定义报告")
	customData := map[string]interface{}{
		"section1": "这是自定义内容的第一部分",
		"section2": "这是自定义内容的第二部分",
		"summary":  "自定义报告生成成功",
	}

	customArgs := map[string]interface{}{
		"action":        "generate",
		"template":      "custom",
		"title":         "自定义报告",
		"author":        "威猛先生",
		"output_format": "text",
		"data":          customData,
		"include_charts": false,
	}

	argsJSON, _ = json.Marshal(customArgs)
	result, err = reportTool.InvokableRun(ctx, string(argsJSON))
	if err != nil {
		log.Fatal("生成自定义报告失败:", err)
	}

	fmt.Printf("✅ 自定义报告结果: %s\n\n", result)

	fmt.Println("========== 测试完成 ==========")
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}