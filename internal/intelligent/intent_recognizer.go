package intelligent

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"new-robot-self/internal/ai"
	"new-robot-self/global"
)

// IntentRecognizer 意图识别器
type IntentRecognizer struct {
	aiService *ai.Service
}

func NewIntentRecognizer() *IntentRecognizer {
	// 使用全局配置的AI服务
	aiService := ai.NewService(global.DefaultConfig)
	err := aiService.InitModel()
	if err != nil {
		fmt.Printf("初始化AI服务失败: %v\n", err)
	}
	
	return &IntentRecognizer{
		aiService: aiService,
	}
}

// Recognize 使用AI服务识别用户意图
func (r *IntentRecognizer) Recognize(ctx context.Context, input string) (*Intent, error) {
	prompt := fmt.Sprintf(`你是一个意图识别专家。请分析以下用户输入，识别其意图类型和提取关键实体。

用户输入：%s

请返回JSON格式的结果，包含以下字段：
- type: 意图类型，必须是以下之一：schedule(日程安排), navigation(导航), health(健康), work(工作), info(信息查询), life(生活服务)
- confidence: 置信度(0-1之间的浮点数)
- entities: 提取的实体信息，如time(时间), location(地点), event(事件), food(食物), quantity(数量)等
- priority: 优先级(high/medium/low)
- context: 上下文信息(原始输入)

只返回JSON，不要有其他内容。`, input)

	message := ai.Message{
		Role:    "intent_recognizer",
		Content: prompt,
	}

	response, err := r.aiService.Chat(ctx, message)
	if err != nil {
		return nil, fmt.Errorf("AI service intent recognition failed: %w", err)
	}

	// 解析AI返回的JSON
	var intent Intent
	err = json.Unmarshal([]byte(response), &intent)
	if err != nil {
		// 如果解析失败，使用基于规则的回退方案
		return r.fallbackRecognize(ctx, input), nil
	}

	// 优先级重写：检查时间紧迫性
	if r.isUrgent(input) {
		intent.Priority = "urgent"
	}

	intent.Timestamp = time.Now().Unix()
	return &intent, nil
}

// fallbackRecognize 回退的基于规则的意图识别
func (r *IntentRecognizer) fallbackRecognize(ctx context.Context, input string) *Intent {
	return &Intent{
		Type:       r.detectIntentType(input),
		Confidence: r.calculateConfidence(input),
		Entities:   r.extractEntities(input),
		Priority:   r.determinePriority(input),
		Context:    input,
		Timestamp:  time.Now().Unix(),
	}
}

// detectIntentType 检测意图类型（回退方案）
func (r *IntentRecognizer) detectIntentType(input string) string {
	keywords := map[string][]string{
		"schedule":   {"打算", "计划", "准备", "提醒", "记得", "别忘了", "设置", "创建", "新建", "三点", "下午", "上午", "晚上", "今天", "明天"},
		"navigation": {"去", "到", "去往", "前往", "怎么走", "导航", "路线", "菜市场", "超市", "商场", "公园", "医院", "公司", "学校"},
		"health":     {"吃了", "食用", "卡路里", "热量", "营养", "GI", "运动", "锻炼", "健身", "血糖", "体重", "健康", "记录", "测量", "检查"},
		"work":       {"代码", "怎么写", "如何实现", "Go", "Python", "Java", "JavaScript", "查询", "搜索", "查找", "文档", "API", "函数", "库", "保存", "记录", "片段", "命令"},
		"info":       {"天气", "气温", "下雨", "晴天", "新闻", "热点", "头条", "搜索", "查一下", "找一下", "信息", "资料"},
		"life":       {"买", "购物", "选购", "菜", "水果", "零食", "日用品", "吃", "喝", "用餐", "哪里", "什么地方", "推荐", "娱乐", "玩", "看电影", "听音乐"},
	}
	
	maxMatches := 0
	var bestType string
	
	for intentType, words := range keywords {
		matches := 0
		for _, word := range words {
			if contains(input, word) {
				matches++
			}
		}
		if matches > maxMatches {
			maxMatches = matches
			bestType = intentType
		}
	}
	
	if bestType == "" {
		return "info"
	}
	
	return bestType
}

// calculateConfidence 计算置信度（回退方案）
func (r *IntentRecognizer) calculateConfidence(input string) float64 {
	keywords := map[string]int{
		"schedule":   3,
		"navigation": 2,
		"health":     2,
		"work":       2,
		"info":       1,
		"life":       2,
	}
	
	intentType := r.detectIntentType(input)
	if weight, ok := keywords[intentType]; ok {
		baseConfidence := 0.6
		return baseConfidence + float64(weight)*0.1
	}
	
	return 0.5
}

// extractEntities 提取实体信息（回退方案）
func (r *IntentRecognizer) extractEntities(input string) map[string]interface{} {
	entities := make(map[string]interface{})
	
	if time := extractTime(input); time != "" {
		entities["time"] = time
	}
	
	if location := extractLocation(input); location != "" {
		entities["location"] = location
	}
	
	if event := extractEvent(input); event != "" {
		entities["event"] = event
	}
	
	if food := extractFood(input); food != "" {
		entities["food"] = food
	}
	
	if quantity := extractQuantity(input); quantity != "" {
		entities["quantity"] = quantity
	}
	
	return entities
}

// determinePriority 确定优先级（回退方案）
func (r *IntentRecognizer) determinePriority(input string) string {
	intentType := r.detectIntentType(input)
	
	// 检查时间紧迫性
	if r.isUrgent(input) {
		return "urgent" // 最高优先级
	}
	
	priorities := map[string]string{
		"schedule":   "high",
		"navigation": "medium",
		"health":     "high",
		"work":       "medium",
		"info":       "low",
		"life":       "medium",
	}
	
	if priority, ok := priorities[intentType]; ok {
		return priority
	}
	return "low"
}

// isUrgent 检查是否为紧急需求
func (r *IntentRecognizer) isUrgent(input string) bool {
	urgentKeywords := []string{"现在", "马上", "立即", "急需", "紧急", "赶紧", "快点", "等一下"}
	for _, keyword := range urgentKeywords {
		if contains(input, keyword) {
			return true
		}
	}
	return false
}

// DetectCompositeIntent 检测并拆分复合意图
func (r *IntentRecognizer) DetectCompositeIntent(ctx context.Context, input string) ([]*Intent, error) {
	// 复合意图的连接词
	connectors := []string{"而且", "并且", "同时", "还有", "另外", "再者", "此外", "而且"}
	
	// 检查是否包含连接词
	for _, connector := range connectors {
		if contains(input, connector) {
			// 拆分输入
			parts := r.splitByConnector(input, connector)
			if len(parts) > 1 {
				var intents []*Intent
				for _, part := range parts {
					intent, err := r.Recognize(ctx, part)
					if err != nil {
						// 如果单个部分识别失败，继续处理其他部分
						continue
					}
					intents = append(intents, intent)
				}
				if len(intents) > 0 {
					return intents, nil
				}
			}
		}
	}
	
	// 如果不是复合意图，返回单个意图
	intent, err := r.Recognize(ctx, input)
	if err != nil {
		return nil, err
	}
	return []*Intent{intent}, nil
}

// splitByConnector 按连接词拆分输入
func (r *IntentRecognizer) splitByConnector(input, connector string) []string {
	var parts []string
	current := []rune{}
	inputRunes := []rune(input)
	connectorRunes := []rune(connector)
	
	for i := 0; i < len(inputRunes); i++ {
		// 检查是否匹配连接词
		if i+len(connectorRunes) <= len(inputRunes) {
			match := true
			for j := 0; j < len(connectorRunes); j++ {
				if inputRunes[i+j] != connectorRunes[j] {
					match = false
					break
				}
			}
			if match {
				if len(current) > 0 {
					parts = append(parts, string(current))
				}
				current = []rune{}
				i += len(connectorRunes) - 1 // 跳过连接词
				continue
			}
		}
		current = append(current, inputRunes[i])
	}
	
	if len(current) > 0 {
		parts = append(parts, string(current))
	}
	
	return parts
}

// SortIntentsByPriority 按优先级排序意图
func (r *IntentRecognizer) SortIntentsByPriority(intents []*Intent) []*Intent {
	priorityOrder := map[string]int{
		"urgent": 4,
		"high":   3,
		"medium": 2,
		"low":    1,
	}
	
	// 冒泡排序（简单实现）
	for i := 0; i < len(intents)-1; i++ {
		for j := 0; j < len(intents)-i-1; j++ {
			if priorityOrder[intents[j].Priority] < priorityOrder[intents[j+1].Priority] {
				intents[j], intents[j+1] = intents[j+1], intents[j]
			}
		}
	}
	
	return intents
}

// 辅助函数
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) && containsHelper(s, substr))
}

func containsHelper(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

func extractTime(input string) string {
	timeKeywords := []string{"三点", "下午", "上午", "晚上", "今天", "明天", "点", "时"}
	for _, keyword := range timeKeywords {
		if contains(input, keyword) {
			return keyword
		}
	}
	return ""
}

func extractLocation(input string) string {
	locationKeywords := []string{"菜市场", "超市", "商场", "公园", "医院", "公司", "学校"}
	for _, keyword := range locationKeywords {
		if contains(input, keyword) {
			return keyword
		}
	}
	return ""
}

func extractEvent(input string) string {
	words := []string{"去", "玩", "买", "吃", "做", "参加"}
	for _, word := range words {
		if contains(input, word) && len(input) > 0 {
			return word + "活动"
		}
	}
	return ""
}

func extractFood(input string) string {
	foodKeywords := []string{"米饭", "面条", "面包", "牛奶", "鸡蛋", "蔬菜", "水果", "肉"}
	for _, keyword := range foodKeywords {
		if contains(input, keyword) {
			return keyword
		}
	}
	return ""
}

func extractQuantity(input string) string {
	for i := 0; i < len(input); i++ {
		if input[i] >= '0' && input[i] <= '9' {
			j := i
			for j < len(input) && (input[j] >= '0' && input[j] <= '9' || input[j] == '.') {
				j++
			}
			if j > i {
				quantity := input[i:j]
				if j < len(input) {
					unit := input[j:j+2]
					return quantity + unit
				}
				return quantity
			}
		}
	}
	return ""
}
