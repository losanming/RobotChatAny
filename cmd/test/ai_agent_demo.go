package main

import (
	"context"
	"fmt"

	"new-robot-self/internal/ai"
	"new-robot-self/internal/agent"
	"new-robot-self/internal/brain"
	"new-robot-self/internal/tools"
	"new-robot-self/global"
)

func main() {
	fmt.Println("=== AI Agent 工具调用测试 ===")
	fmt.Println()

	// 1. 初始化工具注册表
	fmt.Println("【步骤1】初始化工具注册表...")
	toolRegistry := tools.GetGlobalRegistry()
	
	// 如果全局注册表为空，则注册工具
	toolNames := toolRegistry.List()
	if len(toolNames) == 0 {
		fmt.Println("全局工具注册表为空，开始注册工具...")
		
		// 创建知识库
		knowledgeBase := brain.NewKnowledgeBase()
		
		// 注册所有工具
		regErr := tools.RegisterAllTools(knowledgeBase)
		if regErr != nil {
			fmt.Printf("❌ 注册工具失败: %v\n", regErr)
			return
		}
		
		toolNames = toolRegistry.List()
	}
	
	fmt.Printf("✅ 成功注册 %d 个工具\n", len(toolNames))
	
	// 只显示前10个工具，避免输出过多
	maxDisplay := 10
	if len(toolNames) > maxDisplay {
		fmt.Printf("   (显示前%d个工具，共%d个)\n", maxDisplay, len(toolNames))
		for i := 0; i < maxDisplay && i < len(toolNames); i++ {
			fmt.Printf("   • %s\n", toolNames[i])
		}
	} else {
		for _, name := range toolNames {
			fmt.Printf("   • %s\n", name)
		}
	}

	// 2. 初始化AI服务
	fmt.Println("\n【步骤2】初始化AI服务...")
	aiService := ai.NewService(global.DefaultConfig)
	initErr := aiService.InitModel()
	if initErr != nil {
		fmt.Printf("❌ 初始化AI服务失败: %v\n", initErr)
		return
	}
	fmt.Println("✅ AI服务初始化成功")

	// 3. 创建AI Agent
	fmt.Println("\n【步骤3】创建AI Agent...")
	aiAgent := agent.NewAIAgent(aiService, toolRegistry)
	fmt.Println("✅ AI Agent创建成功")

	// 4. 测试工具调用
	fmt.Println("\n【步骤4】测试AI Agent工具调用...")
	testCases := []struct {
		name  string
		input string
	}{
		{"导航测试", "我想去北京西站"},
		{"健康查询", "查询100g大米的卡路里"},
		{"天气查询", "今天北京天气怎么样"},
		{"复合场景", "明天晚上我想要去铜山万达看电影,而且我现在需要去徐州市的苏宁购物"},
	}

	ctx := context.Background()
	
	for _, tc := range testCases {
		fmt.Printf("\n【%s】\n", tc.name)
		fmt.Printf("用户输入: %s\n", tc.input)
		
		response, err := aiAgent.Process(ctx, tc.input)
		if err != nil {
			fmt.Printf("❌ 处理失败: %v\n", err)
		} else {
			fmt.Printf("✅ 处理成功\n")
			fmt.Printf("AI回复: %s\n", response.Content)
			
			if len(response.ToolCalls) > 0 {
				fmt.Printf("调用的工具: %d 个\n", len(response.ToolCalls))
				for i, toolCall := range response.ToolCalls {
					fmt.Printf("  %d. 工具: %s, 参数: %s\n", i+1, toolCall.Function.Name, toolCall.Function.Arguments)
				}
			} else {
				fmt.Println("  未调用工具（AI直接回复）")
			}
		}
	}

	fmt.Println("\n=== 测试完成 ===")
}