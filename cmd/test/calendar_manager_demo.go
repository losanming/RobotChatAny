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
	fmt.Println("========== 测试 calendar_manager 工具 ==========")

	// 初始化工具系统
	ctx := context.Background()
	knowledgeBase := brain.NewKnowledgeBase()

	err := tools.RegisterAllTools(knowledgeBase)
	if err != nil {
		log.Fatal("工具注册失败:", err)
	}

	fmt.Println("✅ 工具系统初始化成功")

	// 获取日程管理工具
	calendarTool, exists := tools.Get("calendar_manager")
	if !exists {
		log.Fatal("日程管理工具不存在")
	}

	fmt.Println("✅ 日程管理工具获取成功\n")

	// 测试1: 创建日程事件
	fmt.Println("📅 测试1: 创建日程事件")
	createArgs := map[string]interface{}{
		"action":     "create",
		"title":      "项目会议",
		"start_time": float64(time.Now().Add(1 * time.Hour).Unix()),
		"duration":   3600.0,
		"location":   "会议室A",
		"description": "讨论项目进度",
		"reminder":   300.0,
	}

	argsJSON, _ := json.Marshal(createArgs)
	result, err := calendarTool.InvokableRun(ctx, string(argsJSON))
	if err != nil {
		log.Fatal("创建日程失败:", err)
	}

	var createResult map[string]interface{}
	json.Unmarshal([]byte(result), &createResult)
	event := createResult["event"].(map[string]interface{})
	eventID := event["id"].(string)

	fmt.Printf("✅ 日程创建成功，ID: %s\n\n", eventID)

	// 测试2: 查询日程事件
	fmt.Println("📅 测试2: 查询日程事件")
	queryArgs := map[string]interface{}{
		"action": "query",
		"id":    eventID,
	}

	argsJSON, _ = json.Marshal(queryArgs)
	result, err = calendarTool.InvokableRun(ctx, string(argsJSON))
	if err != nil {
		log.Fatal("查询日程失败:", err)
	}

	fmt.Printf("✅ 查询结果: %s\n\n", result)

	// 测试3: 列出所有日程
	fmt.Println("📅 测试3: 列出所有日程")
	listArgs := map[string]interface{}{
		"action": "list",
	}

	argsJSON, _ = json.Marshal(listArgs)
	result, err = calendarTool.InvokableRun(ctx, string(argsJSON))
	if err != nil {
		log.Fatal("列出日程失败:", err)
	}

	fmt.Printf("✅ 列出结果: %s\n\n", result)

	// 测试4: 检查冲突
	fmt.Println("📅 测试4: 检查冲突")
	conflictArgs := map[string]interface{}{
		"action":     "check_conflict",
		"start_time": float64(time.Now().Add(30 * time.Minute).Unix()),
		"duration":   1800.0,
	}

	argsJSON, _ = json.Marshal(conflictArgs)
	result, err = calendarTool.InvokableRun(ctx, string(argsJSON))
	if err != nil {
		log.Fatal("检查冲突失败:", err)
	}

	fmt.Printf("✅ 冲突检查结果: %s\n\n", result)

	// 测试5: 日历视图
	fmt.Println("📅 测试5: 日历视图")
	today := time.Now().Format("2006-01-02")
	viewArgs := map[string]interface{}{
		"action":    "calendar_view",
		"view_type": "day",
		"date":      today,
	}

	argsJSON, _ = json.Marshal(viewArgs)
	result, err = calendarTool.InvokableRun(ctx, string(argsJSON))
	if err != nil {
		log.Fatal("生成日历视图失败:", err)
	}

	fmt.Printf("✅ 日历视图结果: %s\n\n", result)

	// 测试6: 删除日程
	fmt.Println("📅 测试6: 删除日程")
	deleteArgs := map[string]interface{}{
		"action": "delete",
		"id":    eventID,
	}

	argsJSON, _ = json.Marshal(deleteArgs)
	result, err = calendarTool.InvokableRun(ctx, string(argsJSON))
	if err != nil {
		log.Fatal("删除日程失败:", err)
	}

	fmt.Printf("✅ 删除结果: %s\n\n", result)

	fmt.Println("========== 测试完成 ==========")
}