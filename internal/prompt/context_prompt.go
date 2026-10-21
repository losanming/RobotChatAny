package prompt

import (
	"fmt"
	"new-robot-self/internal/context"
	"strings"
	"time"
)

// ContextPromptGenerator 上下文相关的prompt生成器
type ContextPromptGenerator struct {
	contextManager *context.ContextManager
	tokenConfig    context.TokenConfig
	lastResponse   string // 记录最后一次AI的回答
}

// NewContextPromptGenerator 创建一个新的上下文prompt生成器
func NewContextPromptGenerator(cm *context.ContextManager) *ContextPromptGenerator {
	return &ContextPromptGenerator{
		contextManager: cm,
		tokenConfig:    context.DefaultTokenConfig(),
		lastResponse:   "",
	}
}

// GeneratePrompt 生成完整的prompt
func (g *ContextPromptGenerator) GeneratePrompt(userID string, currentInput string) (string, error) {
	var sb strings.Builder

	// 1. 添加基础人设
	sb.WriteString(`你是一个性格活泼、说话幽默的老哥。你喜欢用口语化的方式交谈，经常使用"老哥"、"兄弟"等称呼。
`)
	sb.WriteString(`你应该:
`)
	sb.WriteString(`1. 用轻松自然的语气交谈
`)
	sb.WriteString(`2. 适当使用表情符号
`)
	sb.WriteString(`3. 分享自己的经验和见解
`)
	sb.WriteString(`4. 给出实用的建议和解决方案

`)

	// 2. 获取用户上下文
	ctx, err := g.contextManager.GetContext(userID)
	if err != nil {
		return "", fmt.Errorf("获取上下文失败: %v", err)
	}

	// 3. 添加用户特征和偏好
	if len(ctx.LongTermMemory) > 0 {
		sb.WriteString(`用户的特征和偏好：
`)
		for _, info := range ctx.LongTermMemory {
			if info.Category == "user_preference" && info.Weight > 0.7 {
				sb.WriteString(fmt.Sprintf("- %s\n", info.Content))
			}
		}
		sb.WriteString("\n")
	}

	// 4. 添加最近的对话历史
	if len(ctx.ShortTermMemory) > 0 {
		sb.WriteString(`最近的对话历史：
`)
		// 只取最近的5条消息
		start := 0
		if len(ctx.ShortTermMemory) > 5 {
			start = len(ctx.ShortTermMemory) - 5
		}
		for _, msg := range ctx.ShortTermMemory[start:] {
			timeStr := time.Unix(msg.Time, 0).Format("15:04:05")
			// 区分用户消息和AI回答
			prefix := "用户"
			if msg.Type == "ai_response" {
				prefix = "AI"
			}
			sb.WriteString(fmt.Sprintf("[%s] %s: %s\n", timeStr, prefix, msg.Content))
		}
		sb.WriteString("\n")
	}

	// 4.1 特别提醒上一次的回答
	if g.lastResponse != "" {
		sb.WriteString(`你上一次的回答是：
`)
		sb.WriteString(g.lastResponse + "\n\n")
	}

	// 5. 添加相关的历史摘要
	if len(ctx.MediumTermMemory) > 0 {
		sb.WriteString(`相关的历史话题：
`)
		for _, summary := range ctx.MediumTermMemory {
			if g.isRelevantToCurrentInput(summary, currentInput) {
				sb.WriteString(fmt.Sprintf("- %s\n", summary.Content))
				for _, point := range summary.KeyPoints {
					sb.WriteString(fmt.Sprintf("  * %s\n", point))
				}
			}
		}
		sb.WriteString("\n")
	}

	// 6. 添加当前输入的处理指南
	sb.WriteString(`用户当前的输入是：
`)
	sb.WriteString(currentInput + "\n\n")
	sb.WriteString(`请根据以上上下文，以老哥的身份回复用户。回复应该：
`)
	sb.WriteString(`1. 考虑用户的历史对话和偏好
`)
	sb.WriteString(`2. 保持一致的对话风格
`)
	sb.WriteString(`3. 给出有针对性的回应
`)

	return sb.String(), nil
}

// isRelevantToCurrentInput 判断历史摘要是否与当前输入相关
func (g *ContextPromptGenerator) isRelevantToCurrentInput(summary context.Summary, input string) bool {
	// 1. 检查时间相关性（一天内的对话）
	if time.Now().Unix()-summary.EndTime < 24*60*60 {
		return true
	}

	// 2. 检查主题相关性
	for _, topic := range summary.Topics {
		if strings.Contains(input, topic) {
			return true
		}
	}

	// 3. 检查关键点相关性
	for _, point := range summary.KeyPoints {
		if strings.Contains(input, point) {
			return true
		}
	}

	return false
}

// GetUserPreference 获取用户偏好设置
func (g *ContextPromptGenerator) GetUserPreference(userID string) (map[string]string, error) {
	ctx, err := g.contextManager.GetContext(userID)
	if err != nil {
		return nil, err
	}

	prefs := make(map[string]string)
	for _, info := range ctx.LongTermMemory {
		if info.Category == "user_preference" {
			parts := strings.SplitN(info.Content, ":", 2)
			if len(parts) == 2 {
				prefs[strings.TrimSpace(parts[0])] = strings.TrimSpace(parts[1])
			}
		}
	}

	return prefs, nil
}

// RecordAIResponse 记录AI的回答
func (g *ContextPromptGenerator) RecordAIResponse(userID string, response string) error {
	// 保存最后一次回答
	g.lastResponse = response

	// 添加到上下文
	return g.contextManager.AddMessage(userID, context.Message{
		Content:  response,
		Time:     time.Now().Unix(),
		Priority: 2, // AI回答的优先级设为中等
		Type:     "ai_response",
		Keywords: g.ExtractKeywords(response), // 从回答中提取关键词
	})
}

// extractKeywords 从文本中提取关键词
func (g *ContextPromptGenerator) ExtractKeywords(text string) []string {
	// 这里可以实现更复杂的关键词提取算法
	// 目前简单实现：按空格分词，取长度大于2的词
	words := strings.Fields(text)
	keywords := make([]string, 0)
	for _, word := range words {
		if len(word) > 2 {
			keywords = append(keywords, word)
		}
	}
	return keywords
}

// UpdateUserPreference 更新用户偏好设置
func (g *ContextPromptGenerator) UpdateUserPreference(userID string, key, value string) error {
	info := context.KeyInfo{
		Content:  fmt.Sprintf("%s: %s", key, value),
		Time:     time.Now().Unix(),
		Category: "user_preference",
		Weight:   0.8, // 默认较高权重
	}

	ctx, err := g.contextManager.GetContext(userID)
	if err != nil {
		return err
	}

	// 更新或添加偏好
	found := false
	for i, existing := range ctx.LongTermMemory {
		if existing.Category == "user_preference" {
			parts := strings.SplitN(existing.Content, ":", 2)
			if len(parts) == 2 && strings.TrimSpace(parts[0]) == key {
				ctx.LongTermMemory[i] = info
				found = true
				break
			}
		}
	}

	if !found {
		ctx.LongTermMemory = append(ctx.LongTermMemory, info)
	}

	// 保存更新后的上下文
	return g.contextManager.AddMessage(userID, context.Message{
		Content:  fmt.Sprintf("更新用户偏好：%s = %s", key, value),
		Time:     time.Now().Unix(),
		Priority: 1,
		Type:     "system",
		Keywords: []string{"user_preference", key},
	})
}
