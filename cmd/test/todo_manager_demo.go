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
	fmt.Println("========== 测试 todo_manager 工具 ==========")

	// 初始化工具系统
	ctx := context.Background()
	knowledgeBase := brain.NewKnowledgeBase()

	err := tools.RegisterAllTools(knowledgeBase)
	if err != nil {
		log.Fatal("工具注册失败:", err)
	}

	fmt.Println("✅ 工具系统初始化成功")

	// 获取待办事项管理工具
	todoTool, exists := tools.Get("todo_manager")
	if !exists {
		log.Fatal("待办事项管理工具不存在")
	}

	fmt.Println("✅ 待办事项管理工具获取成功\n")

	// 测试1: 创建待办事项
	fmt.Println("✅ 测试1: 创建待办事项")
	createArgs := map[string]interface{}{
		"action":     "create",
		"title":      "完成代码重构",
		"description": "重构用户模块代码",
		"priority":   "high",
		"tags":       []string{"工作", "开发"},
		"due_date":   float64(time.Now().Add(7 * 24 * time.Hour).Unix()),
	}

	argsJSON, _ := json.Marshal(createArgs)
	result, err := todoTool.InvokableRun(ctx, string(argsJSON))
	if err != nil {
		log.Fatal("创建待办事项失败:", err)
	}

	var createResult map[string]interface{}
	json.Unmarshal([]byte(result), &createResult)
	todo := createResult["todo"].(map[string]interface{})
	todoID := todo["id"].(string)

	fmt.Printf("✅ 待办事项创建成功，ID: %s\n\n", todoID)

	// 测试2: 查询待办事项
	fmt.Println("✅ 测试2: 查询待办事项")
	queryArgs := map[string]interface{}{
		"action": "query",
		"id":    todoID,
	}

	argsJSON, _ = json.Marshal(queryArgs)
	result, err = todoTool.InvokableRun(ctx, string(argsJSON))
	if err != nil {
		log.Fatal("查询待办事项失败:", err)
	}

	fmt.Printf("✅ 查询结果: %s\n\n", result)

	// 测试3: 创建更多待办事项
	fmt.Println("✅ 测试3: 创建更多待办事项")
	moreTodos := []map[string]interface{}{
		{
			"action":     "create",
			"title":      "编写测试用例",
			"description": "为新功能编写单元测试",
			"priority":   "medium",
			"tags":       []string{"工作", "测试"},
		},
		{
			"action":     "create",
			"title":      "健身运动",
			"description": "晚上室内慢跑1小时",
			"priority":   "medium",
			"tags":       []string{"健康", "运动"},
		},
		{
			"action":     "create",
			"title":      "学习新技术",
			"description": "学习Go语言高级特性",
			"priority":   "low",
			"tags":       []string{"学习", "技术"},
		},
	}

	for _, todoData := range moreTodos {
		argsJSON, _ = json.Marshal(todoData)
		result, err = todoTool.InvokableRun(ctx, string(argsJSON))
		if err != nil {
			log.Printf("创建待办事项失败: %v\n", err)
		} else {
			fmt.Printf("✅ 创建成功: %s\n", todoData["title"])
		}
	}
	fmt.Println()

	// 测试4: 列出所有待办事项
	fmt.Println("✅ 测试4: 列出所有待办事项")
	listArgs := map[string]interface{}{
		"action": "list",
	}

	argsJSON, _ = json.Marshal(listArgs)
	result, err = todoTool.InvokableRun(ctx, string(argsJSON))
	if err != nil {
		log.Fatal("列出待办事项失败:", err)
	}

	fmt.Printf("✅ 列出结果: %s\n\n", result)

	// 测试5: 按状态筛选
	fmt.Println("✅ 测试5: 按状态筛选")
	filterArgs := map[string]interface{}{
		"action":        "list",
		"filter_status": "pending",
	}

	argsJSON, _ = json.Marshal(filterArgs)
	result, err = todoTool.InvokableRun(ctx, string(argsJSON))
	if err != nil {
		log.Fatal("筛选待办事项失败:", err)
	}

	fmt.Printf("✅ 筛选结果: %s\n\n", result)

	// 测试6: 搜索待办事项
	fmt.Println("✅ 测试6: 搜索待办事项")
	searchArgs := map[string]interface{}{
		"action":         "search",
		"search_keyword": "代码",
	}

	argsJSON, _ = json.Marshal(searchArgs)
	result, err = todoTool.InvokableRun(ctx, string(argsJSON))
	if err != nil {
		log.Fatal("搜索待办事项失败:", err)
	}

	fmt.Printf("✅ 搜索结果: %s\n\n", result)

	// 测试7: 完成待办事项
	fmt.Println("✅ 测试7: 完成待办事项")
	completeArgs := map[string]interface{}{
		"action": "complete",
		"id":     todoID,
	}

	argsJSON, _ = json.Marshal(completeArgs)
	result, err = todoTool.InvokableRun(ctx, string(argsJSON))
	if err != nil {
		log.Fatal("完成待办事项失败:", err)
	}

	fmt.Printf("✅ 完成结果: %s\n\n", result)

	// 测试8: 删除待办事项
	fmt.Println("✅ 测试8: 删除待办事项")
	deleteArgs := map[string]interface{}{
		"action": "delete",
		"id":     todoID,
	}

	argsJSON, _ = json.Marshal(deleteArgs)
	result, err = todoTool.InvokableRun(ctx, string(argsJSON))
	if err != nil {
		log.Fatal("删除待办事项失败:", err)
	}

	fmt.Printf("✅ 删除结果: %s\n\n", result)

	fmt.Println("========== 测试完成 ==========")
}