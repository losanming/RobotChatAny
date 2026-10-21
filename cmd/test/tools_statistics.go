package main

import (
	"context"
	"fmt"
	"new-robot-self/internal/brain"
	"new-robot-self/internal/tools"
)

func main() {
	fmt.Println("========== 工具系统统计 ==========")

	// 初始化工具系统
	ctx := context.Background()
	knowledgeBase := brain.NewKnowledgeBase()

	err := tools.RegisterAllTools(knowledgeBase)
	if err != nil {
		fmt.Printf("工具注册失败: %v\n", err)
		return
	}

	fmt.Println("✅ 工具系统初始化成功")

	// 列出所有工具
	toolNames := tools.List()
	fmt.Printf("📊 工具总数: %d\n\n", len(toolNames))

	// 按分类统计工具
	categories := make(map[string][]string)
	for _, name := range toolNames {
		tool, exists := tools.Get(name)
		if exists {
			info, _ := tool.Info(ctx)
			if info != nil {
				category := getCategoryFromDesc(info.Desc)
				categories[category] = append(categories[category], name)
			}
		}
	}

	fmt.Println("📋 按分类统计:")
	for category, tools := range categories {
		fmt.Printf("  %s: %d个工具\n", category, len(tools))
		for _, tool := range tools {
			fmt.Printf("    - %s\n", tool)
		}
		fmt.Println()
	}

	fmt.Println("========== 统计完成 ==========")
}

func getCategoryFromDesc(desc string) string {
	// 简单的分类推断
	if contains(desc, "日程") || contains(desc, "待办") || contains(desc, "时间") {
		return "个人管理"
	}
	if contains(desc, "文件") || contains(desc, "数据") || contains(desc, "报告") {
		return "信息分析"
	}
	if contains(desc, "代码") || contains(desc, "翻译") || contains(desc, "计算") || contains(desc, "时区") {
		return "通用能力"
	}
	if contains(desc, "位置") || contains(desc, "地理") || contains(desc, "天气") || contains(desc, "导航") {
		return "位置服务"
	}
	if contains(desc, "健康") || contains(desc, "运动") {
		return "健康管理"
	}
	if contains(desc, "学习") || contains(desc, "知识") || contains(desc, "偏好") {
		return "学习进化"
	}
	if contains(desc, "搜索") || contains(desc, "新闻") || contains(desc, "内容") {
		return "信息获取"
	}
	return "其他"
}

func contains(text, substr string) bool {
	return len(text) > 0 && len(substr) > 0 && 
		   findSubstring(text, substr) >= 0
}

func findSubstring(text, substr string) int {
	for i := 0; i <= len(text)-len(substr); i++ {
		if text[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}