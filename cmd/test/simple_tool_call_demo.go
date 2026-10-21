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
	fmt.Println("=== 简单AI工具调用测试 ===")
	fmt.Println()

	ctx := context.Background()

	// 1. 初始化工具
	fmt.Println("【步骤1】初始化工具...")
	toolRegistry := tools.GetGlobalRegistry()
	toolNames := toolRegistry.List()
	if len(toolNames) == 0 {
		knowledgeBase := brain.NewKnowledgeBase()
		tools.RegisterAllTools(knowledgeBase)
		toolNames = toolRegistry.List()
	}
	fmt.Printf("✅ 已注册 %d 个工具\n", len(toolNames))

	// 2. 初始化AI服务
	fmt.Println("\n【步骤2】初始化AI服务...")
	aiService := ai.NewService(global.DefaultConfig)
	err := aiService.InitModel()
	if err != nil {
		fmt.Printf("❌ AI服务初始化失败: %v\n", err)
		return
	}
	fmt.Println("✅ AI服务初始化成功")

	// 3. 创建AI Agent
	fmt.Println("\n【步骤3】创建AI Agent...")
	aiAgent := agent.NewAIAgent(aiService, toolRegistry)
	fmt.Println("✅ AI Agent创建成功")

	// 4. 测试明确的工具调用场景
	fmt.Println("\n【步骤4】测试工具调用...")
	
	testCases := []string{
		"帮我查一下北京今天的天气",
		"我想去天安门广场，帮我规划路线",
		"查询一下苹果的营养成分",
	}

	for i, input := range testCases {
		fmt.Printf("\n【测试%d】%s\n", i+1, input)
		
		response, err := aiAgent.Process(ctx, input)
		if err != nil {
			fmt.Printf("❌ 处理失败: %v\n", err)
		} else {
			fmt.Printf("✅ 处理成功\n")
			fmt.Printf("AI回复: %s\n", response.Content)
			
			if len(response.ToolCalls) > 0 {
				fmt.Printf("调用了 %d 个工具:\n", len(response.ToolCalls))
				for j, toolCall := range response.ToolCalls {
					fmt.Printf("  %d. 工具: %s\n", j+1, toolCall.Function.Name)
					fmt.Printf("     参数: %s\n", toolCall.Function.Arguments)
				}
			} else {
				fmt.Println("  AI未调用工具（直接回复）")
			}
		}
	}

	fmt.Println("\n=== 测试完成 ===")
	fmt.Println("注意：如果AI未调用工具，可能是因为：")
	fmt.Println("1. AI模型不支持工具调用功能")
	fmt.Println("2. 提示词不够明确")
	fmt.Println("3. 需要使用支持工具调用的模型")
}