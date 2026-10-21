package intelligent

import (
	"context"
	"fmt"
	"time"
)

// IntelligentAssistant 智能助手
type IntelligentAssistant struct {
	IntentRecognizer *IntentRecognizer
	ToolCoordinator *ToolCoordinator
	ContextMemory   *ContextMemory
	TriggerManager  *TriggerManager
}

func NewIntelligentAssistant() *IntelligentAssistant {
	return &IntelligentAssistant{
		IntentRecognizer: NewIntentRecognizer(),
		ToolCoordinator: NewToolCoordinator(),
		ContextMemory:   NewContextMemory(),
		TriggerManager:  NewTriggerManager(),
	}
}

// Start 启动智能助手
func (a *IntelligentAssistant) Start() error {
	// 启动触发管理器
	a.TriggerManager.Start()
	
	// TODO: 初始化其他组件
	
	return nil
}

// Stop 停止智能助手
func (a *IntelligentAssistant) Stop() error {
	a.TriggerManager.Stop()
	return nil
}

// Process 处理用户输入
func (a *IntelligentAssistant) Process(ctx context.Context, input string) (*AssistantResponse, error) {
	// 1. 意图识别
	intent, err := a.IntentRecognizer.Recognize(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("intent recognition failed: %w", err)
	}

	// 2. 记录用户输入到上下文
	a.ContextMemory.AddMessage("user", input, intent)

	// 3. 工具协调规划
	plan, err := a.ToolCoordinator.Plan(ctx, intent)
	if err != nil {
		return nil, fmt.Errorf("tool planning failed: %w", err)
	}

	// 4. 执行工具调用
	toolResults, err := a.ToolCoordinator.Execute(ctx, plan)
	if err != nil {
		return nil, fmt.Errorf("tool execution failed: %w", err)
	}

	// 5. 生成回复
	response := a.generateResponse(input, intent, toolResults)

	// 6. 记录助手回复到上下文
	a.ContextMemory.AddMessage("assistant", response.Text, intent)

	// 7. 记录活动
	activity := map[string]interface{}{
		"input":  input,
		"intent": intent.Type,
		"results": toolResults,
	}
	a.ContextMemory.AddActivity(intent.Type, activity)

	// 8. 学习用户偏好（在锁外调用）
	// 用户消息学习
	userMessage := Message{
		Role:      "user",
		Content:   input,
		Intent:    intent,
		Timestamp: intent.Timestamp,
	}
	a.ContextMemory.LearnFromMessage(userMessage)
	
	// 活动学习
	activityData := Activity{
		Type:      intent.Type,
		Details:   activity,
		Timestamp: intent.Timestamp,
	}
	a.ContextMemory.LearnFromActivity(activityData)

	return response, nil
}

// generateResponse 生成回复
func (a *IntelligentAssistant) generateResponse(input string, intent *Intent, toolResults map[string]interface{}) *AssistantResponse {
	// 根据意图类型生成不同的回复
	switch intent.Type {
	case "schedule":
		return a.generateScheduleResponse(input, intent, toolResults)
	case "navigation":
		return a.generateNavigationResponse(input, intent, toolResults)
	case "health":
		return a.generateHealthResponse(input, intent, toolResults)
	case "work":
		return a.generateWorkResponse(input, intent, toolResults)
	case "info":
		return a.generateInfoResponse(input, intent, toolResults)
	case "life":
		return a.generateLifeResponse(input, intent, toolResults)
	default:
		return a.generateDefaultResponse(input, intent, toolResults)
	}
}

// generateScheduleResponse 生成日程相关回复
func (a *IntelligentAssistant) generateScheduleResponse(input string, intent *Intent, toolResults map[string]interface{}) *AssistantResponse {
	event, _ := intent.Entities["event"].(string)
	timeStr, _ := intent.Entities["time"].(string)

	text := fmt.Sprintf("好的，我已经为您安排了日程：\n")
	text += fmt.Sprintf("📅 事件：%s\n", event)
	if timeStr != "" {
		text += fmt.Sprintf("⏰ 时间：%s\n", timeStr)
	}
	text += fmt.Sprintf("✅ 我会在合适的时间提醒您")

	actions := []string{
		"查看详细日程",
		"修改日程",
		"设置提醒",
	}

	return &AssistantResponse{
		Text:        text,
		Intent:      intent,
		ToolResults: toolResults,
		Actions:     actions,
		Confidence:  intent.Confidence,
		Timestamp:   time.Now().Unix(),
	}
}

// generateNavigationResponse 生成导航相关回复
func (a *IntelligentAssistant) generateNavigationResponse(input string, intent *Intent, toolResults map[string]interface{}) *AssistantResponse {
	location, _ := intent.Entities["location"].(string)

	text := fmt.Sprintf("收到，我来帮您规划去%s的路线：\n", location)
	text += fmt.Sprintf("🗺️ 正在为您查询最佳路线...\n")
	text += fmt.Sprintf("🚗 同时查询周边情况...\n")
	text += fmt.Sprintf("🌤️ 查看当地天气情况...\n")
	text += fmt.Sprintf("✅ 综合信息已准备好，随时可以出发")

	actions := []string{
		"查看详细路线",
		"查询周边更多",
		"查看其他交通方式",
	}

	return &AssistantResponse{
		Text:        text,
		Intent:      intent,
		ToolResults: toolResults,
		Actions:     actions,
		Confidence:  intent.Confidence,
		Timestamp:   time.Now().Unix(),
	}
}

// generateHealthResponse 生成健康相关回复
func (a *IntelligentAssistant) generateHealthResponse(input string, intent *Intent, toolResults map[string]interface{}) *AssistantResponse {
	food, _ := intent.Entities["food"].(string)
	quantity, _ := intent.Entities["quantity"].(string)

	text := fmt.Sprintf("我来帮您查询%s的营养信息：\n", food)
	if quantity != "" {
		text += fmt.Sprintf("📊 数量：%s\n", quantity)
	}
	text += fmt.Sprintf("🔥 正在计算卡路里...\n")
	text += fmt.Sprintf("🥗 分析营养成分...\n")
	text += fmt.Sprintf("✅ 营养信息已准备好")

	actions := []string{
		"查看详细营养信息",
		"记录饮食",
		"获取饮食建议",
	}

	return &AssistantResponse{
		Text:        text,
		Intent:      intent,
		ToolResults: toolResults,
		Actions:     actions,
		Confidence:  intent.Confidence,
		Timestamp:   time.Now().Unix(),
	}
}

// generateWorkResponse 生成工作相关回复
func (a *IntelligentAssistant) generateWorkResponse(input string, intent *Intent, toolResults map[string]interface{}) *AssistantResponse {
	searchTerm, _ := intent.Entities["search_term"].(string)

	text := fmt.Sprintf("我来帮您查询工作相关信息：\n")
	if searchTerm != "" {
		text += fmt.Sprintf("🔍 搜索：%s\n", searchTerm)
	}
	text += fmt.Sprintf("📚 正在搜索技术文档...\n")
	text += fmt.Sprintf("💻 查询代码片段...\n")
	text += fmt.Sprintf("✅ 相关资料已准备好")

	actions := []string{
		"查看详细文档",
		"保存代码片段",
		"搜索更多资料",
	}

	return &AssistantResponse{
		Text:        text,
		Intent:      intent,
		ToolResults: toolResults,
		Actions:     actions,
		Confidence:  intent.Confidence,
		Timestamp:   time.Now().Unix(),
	}
}

// generateInfoResponse 生成信息查询回复
func (a *IntelligentAssistant) generateInfoResponse(input string, intent *Intent, toolResults map[string]interface{}) *AssistantResponse {
	text := fmt.Sprintf("我来帮您查询相关信息：\n")
	text += fmt.Sprintf("🔍 正在搜索...\n")
	text += fmt.Sprintf("✅ 查询结果已准备好")

	actions := []string{
		"查看详细结果",
		"搜索更多信息",
	}

	return &AssistantResponse{
		Text:        text,
		Intent:      intent,
		ToolResults: toolResults,
		Actions:     actions,
		Confidence:  intent.Confidence,
		Timestamp:   time.Now().Unix(),
	}
}

// generateLifeResponse 生成生活服务回复
func (a *IntelligentAssistant) generateLifeResponse(input string, intent *Intent, toolResults map[string]interface{}) *AssistantResponse {
	location, _ := intent.Entities["location"].(string)
	item, _ := intent.Entities["item"].(string)

	text := fmt.Sprintf("收到，我来帮您安排：\n")
	if item != "" {
		text += fmt.Sprintf("🛒 购买：%s\n", item)
	}
	if location != "" {
		text += fmt.Sprintf("📍 地点：%s\n", location)
	}
	text += fmt.Sprintf("🗺️ 正在规划路线...\n")
	text += fmt.Sprintf("🏪 查询周边购物信息...\n")
	text += fmt.Sprintf("✅ 购物信息已准备好")

	actions := []string{
		"查看详细路线",
		"查询周边更多",
		"开始导航",
	}

	return &AssistantResponse{
		Text:        text,
		Intent:      intent,
		ToolResults: toolResults,
		Actions:     actions,
		Confidence:  intent.Confidence,
		Timestamp:   time.Now().Unix(),
	}
}

// generateDefaultResponse 生成默认回复
func (a *IntelligentAssistant) generateDefaultResponse(input string, intent *Intent, toolResults map[string]interface{}) *AssistantResponse {
	text := fmt.Sprintf("我收到您的消息：%s\n", input)
	text += fmt.Sprintf("正在为您处理...\n")
	text += fmt.Sprintf("✅ 处理完成")

	return &AssistantResponse{
		Text:        text,
		Intent:      intent,
		ToolResults: toolResults,
		Actions:     []string{},
		Confidence:  intent.Confidence,
		Timestamp:   time.Now().Unix(),
	}
}

// GetContextSummary 获取上下文摘要
func (a *IntelligentAssistant) GetContextSummary() map[string]interface{} {
	return a.ContextMemory.GetContextSummary()
}

// GetRecentMessages 获取最近消息
func (a *IntelligentAssistant) GetRecentMessages(count int) []Message {
	return a.ContextMemory.GetRecentMessages(count)
}

// GetContextMemory 获取上下文记忆
func (a *IntelligentAssistant) GetContextMemory() *ContextMemory {
	return a.ContextMemory
}

// GetTriggerManager 获取触发管理器
func (a *IntelligentAssistant) GetTriggerManager() *TriggerManager {
	return a.TriggerManager
}