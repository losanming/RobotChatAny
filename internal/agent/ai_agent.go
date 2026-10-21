package agent

import (
	"context"
	"fmt"

	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/schema"

	"new-robot-self/internal/ai"
	"new-robot-self/internal/tools"
)

// AIAgent AI智能代理
// 负责协调AI服务和工具调用，实现智能业务处理
type AIAgent struct {
	aiService    *ai.Service
	toolRegistry *tools.ToolRegistry
}

// NewAIAgent 创建新的AI代理
func NewAIAgent(aiService *ai.Service, toolRegistry *tools.ToolRegistry) *AIAgent {
	return &AIAgent{
		aiService:    aiService,
		toolRegistry: toolRegistry,
	}
}

// Process 处理用户输入，让AI自主选择和调用工具
func (a *AIAgent) Process(ctx context.Context, userInput string) (*AgentResponse, error) {
	// 1. 将工具注册表转换为AI服务需要的格式
	toolMap, err := a.convertToolsToInvokable()
	if err != nil {
		return nil, fmt.Errorf("convert tools failed: %w", err)
	}

	// 2. 构建消息
	messages := []*schema.Message{
		{
			Role:    schema.System,
			Content: `你是一个智能助手，能够使用各种工具来帮助用户完成他们的需求。

重要规则：
1. 当用户提到导航、路线、去某地时，必须调用navigation工具
2. 当用户提到食物、卡路里、营养时，必须调用diet_manager工具  
3. 当用户提到天气、气温、下雨时，必须调用weather_query工具
4. 当用户提到健康、血糖、体重时，必须调用health_manager工具
5. 只有当工具无法满足用户需求时，才直接回复用户

请根据用户的问题选择合适的工具来解决问题。`,
		},
		{
			Role:    schema.User,
			Content: userInput,
		},
	}

	// 3. 调用AI服务，让AI自主选择和调用工具
	responseMessages, err := a.aiService.InvokeWithTools(ctx, toolMap, messages)
	if err != nil {
		return nil, fmt.Errorf("AI invoke with tools failed: %w", err)
	}

	// 4. 解析AI的响应
	return a.parseResponse(responseMessages), nil
}

// convertToolsToInvokable 将工具注册表转换为InvokableTool格式
func (a *AIAgent) convertToolsToInvokable() (map[string]tool.InvokableTool, error) {
	toolNames := a.toolRegistry.List()
	
	toolMap := make(map[string]tool.InvokableTool)
	for _, name := range toolNames {
		baseTool, exists := a.toolRegistry.Get(name)
		if !exists {
			continue
		}
		
		// BaseTool已经实现了tool.InvokableTool接口
		invokableTool, ok := baseTool.(tool.InvokableTool)
		if !ok {
			return nil, fmt.Errorf("tool %s does not implement InvokableTool interface", name)
		}
		toolMap[name] = invokableTool
	}
	
	return toolMap, nil
}

// parseResponse 解析AI的响应
func (a *AIAgent) parseResponse(messages []*schema.Message) *AgentResponse {
	if len(messages) == 0 {
		return &AgentResponse{
			Success: false,
			Error:   "No response from AI",
		}
	}

	lastMessage := messages[len(messages)-1]
	
	response := &AgentResponse{
		Success: true,
		Content: lastMessage.Content,
		Messages: messages,
	}

	// 检查是否有工具调用
	for _, msg := range messages {
		if len(msg.ToolCalls) > 0 {
			for _, toolCall := range msg.ToolCalls {
				response.ToolCalls = append(response.ToolCalls, &toolCall)
			}
		}
	}

	return response
}

// AgentResponse AI代理响应
type AgentResponse struct {
	Success    bool
	Content    string
	Error      string
	Messages   []*schema.Message
	ToolCalls  []*schema.ToolCall
	ToolResults map[string]interface{}
}
