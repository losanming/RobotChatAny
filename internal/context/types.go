package context

import "time"

// Message 表示一条对话消息
type Message struct {
	Content   string    `json:"content" xorm:"text 'content'"`
	Time      int64     `json:"time" xorm:"'time'"`      // 时间戳
	Priority  int       `json:"priority" xorm:"'priority'"`  // 优先级
	Type      string    `json:"type" xorm:"'type'"`      // 消息类型
	Keywords  []string  `json:"keywords" xorm:"json 'keywords'"`  // 关键词
}

// Summary 表示对话摘要
type Summary struct {
	Content    string    `json:"content" xorm:"text 'content'"`
	StartTime  int64     `json:"start_time" xorm:"'start_time'"`
	EndTime    int64     `json:"end_time" xorm:"'end_time'"`
	Topics     []string  `json:"topics" xorm:"json 'topics'"`
	KeyPoints  []string  `json:"key_points" xorm:"json 'key_points'"`
}

// KeyInfo 表示长期记忆中的关键信息
type KeyInfo struct {
	Content    string    `json:"content" xorm:"text 'content'"`
	Time       int64     `json:"time" xorm:"'time'"`
	Category   string    `json:"category" xorm:"'category'"`  // 分类
	Weight     float64   `json:"weight" xorm:"'weight'"`      // 重要性权重
}

// ChatContext 表示完整的对话上下文
type ChatContext struct {
	UserID           string    `json:"user_id" xorm:"pk 'user_id'"`
	ShortTermMemory  []Message `json:"short_term" xorm:"json 'short_term'"`    // 短期记忆
	MediumTermMemory []Summary `json:"medium_term" xorm:"json 'medium_term'"`  // 中期记忆
	LongTermMemory   []KeyInfo `json:"long_term" xorm:"json 'long_term'"`      // 长期记忆
	UpdatedAt        int64     `json:"updated_at" xorm:"'updated_at'"`
}

// TokenConfig 定义不同场景的token限制
type TokenConfig struct {
	MaxTokensForRecent   int `json:"max_tokens_for_recent"`   // 最近对话的token限制
	MaxTokensForSummary  int `json:"max_tokens_for_summary"`  // 摘要信息的token限制
	MaxTokensForUserInfo int `json:"max_tokens_for_user"`     // 用户信息的token限制
	MaxTokensForPrompt   int `json:"max_tokens_for_prompt"`   // 基础prompt的token限制
}

// CompressedMessage 表示压缩后的消息
type CompressedMessage struct {
	KeyPoints []string `json:"key_points" xorm:"json 'key_points'"`  // 关键信息点
	Sentiment float64  `json:"sentiment" xorm:"'sentiment'"`         // 情感值
	Time      int64    `json:"time" xorm:"'time'"`                  // 时间戳
}

// NewChatContext 创建新的对话上下文
func NewChatContext(userID string) *ChatContext {
	return &ChatContext{
		UserID:           userID,
		ShortTermMemory:  make([]Message, 0),
		MediumTermMemory: make([]Summary, 0),
		LongTermMemory:   make([]KeyInfo, 0),
		UpdatedAt:        time.Now().Unix(),
	}
}

// DefaultTokenConfig 返回默认的token配置
func DefaultTokenConfig() TokenConfig {
	return TokenConfig{
		MaxTokensForRecent:   1000,
		MaxTokensForSummary:  500,
		MaxTokensForUserInfo: 300,
		MaxTokensForPrompt:   200,
	}
}

// SceneTokenLimits 定义不同场景的token限制
var SceneTokenLimits = map[string]int{
	"日常闲聊": 800,   // 较少token
	"情感倾诉": 1500,  // 需要更多情感上下文
	"深度话题": 2000,  // 需要更多知识背景
	"默认":    1000,
}
