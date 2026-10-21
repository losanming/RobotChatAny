package context

import (
	"context"
	"encoding/json"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
	"xorm.io/xorm"
)

// ContextManager 管理对话上下文
type ContextManager struct {
	db          *xorm.Engine
	redis       *redis.Client
	ctx         context.Context
	tokenConfig TokenConfig
	mu          sync.RWMutex
}

// NewContextManager 创建上下文管理器
func NewContextManager(db *xorm.Engine, redis *redis.Client) *ContextManager {
	return &ContextManager{
		db:          db,
		redis:       redis,
		ctx:         context.Background(),
		tokenConfig: DefaultTokenConfig(),
	}
}

// AddMessage 添加新消息到短期记忆
func (m *ContextManager) AddMessage(userID string, msg Message) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	// 从Redis获取当前上下文
	ctx, err := m.getContextFromCache(userID)
	if err != nil {
		// 如果缓存中没有，从数据库加载
		ctx, err = m.getContextFromDB(userID)
		if err != nil {
			// 如果数据库中也没有，创建新的
			ctx = NewChatContext(userID)
		}
	}

	// 添加新消息到短期记忆
	ctx.ShortTermMemory = append(ctx.ShortTermMemory, msg)
	ctx.UpdatedAt = time.Now().Unix()

	// 如果短期记忆超过限制，触发压缩
	if m.shouldCompress(ctx.ShortTermMemory) {
		if err := m.compressShortTermMemory(ctx); err != nil {
			return err
		}
	}

	// 保存到缓存和数据库
	if err := m.saveContext(ctx); err != nil {
		return err
	}

	return nil
}

// GetContext 获取用户的完整上下文
func (m *ContextManager) GetContext(userID string) (*ChatContext, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	// 先从缓存获取
	ctx, err := m.getContextFromCache(userID)
	if err == nil {
		return ctx, nil
	}

	// 缓存没有，从数据库获取
	ctx, err = m.getContextFromDB(userID)
	if err != nil {
		return nil, err
	}

	// 放入缓存
	if err := m.setContextCache(ctx); err != nil {
		return nil, err
	}

	return ctx, nil
}

// ExtractRelevantContext 提取相关上下文
func (m *ContextManager) ExtractRelevantContext(userID string, currentTopic string, scene string) (string, error) {
	ctx, err := m.GetContext(userID)
	if err != nil {
		return "", err
	}

	// 获取场景的token限制
	tokenLimit := SceneTokenLimits[scene]
	if tokenLimit == 0 {
		tokenLimit = SceneTokenLimits["默认"]
	}

	// 构建上下文
	var relevantContext string

	// 1. 添加最近的短期记忆
	for i := len(ctx.ShortTermMemory) - 1; i >= 0; i-- {
		msg := ctx.ShortTermMemory[i]
		if isRelevant(msg, currentTopic) {
			relevantContext += msg.Content + "\n"
		}
	}

	// 2. 添加相关的中期记忆摘要
	for _, summary := range ctx.MediumTermMemory {
		if containsTopic(summary.Topics, currentTopic) {
			relevantContext += "摘要: " + summary.Content + "\n"
		}
	}

	// 3. 添加重要的长期记忆
	for _, info := range ctx.LongTermMemory {
		if info.Weight > 0.8 && isRelevant(Message{Content: info.Content}, currentTopic) {
			relevantContext += "重要信息: " + info.Content + "\n"
		}
	}

	return relevantContext, nil
}

// 内部辅助方法

func (m *ContextManager) getContextFromCache(userID string) (*ChatContext, error) {
	data, err := m.redis.Get(m.ctx, "chat_context:"+userID).Bytes()
	if err != nil {
		return nil, err
	}

	var ctx ChatContext
	if err := json.Unmarshal(data, &ctx); err != nil {
		return nil, err
	}

	return &ctx, nil
}

func (m *ContextManager) getContextFromDB(userID string) (*ChatContext, error) {
	ctx := &ChatContext{UserID: userID}
	has, err := m.db.Get(ctx)
	if err != nil {
		return nil, err
	}
	if !has {
		// 如果数据库中不存在，创建新的上下文
		return NewChatContext(userID), nil
	}
	return ctx, nil
}

func (m *ContextManager) setContextCache(ctx *ChatContext) error {
	data, err := json.Marshal(ctx)
	if err != nil {
		return err
	}

	// 设置缓存，有效期1小时
	return m.redis.Set(m.ctx, "chat_context:"+ctx.UserID, data, time.Hour).Err()
}

func (m *ContextManager) saveContext(ctx *ChatContext) error {
	// 保存到缓存
	if err := m.setContextCache(ctx); err != nil {
		return err
	}

	// 检查是否存在
	exists, err := m.db.Exist(&ChatContext{UserID: ctx.UserID})
	if err != nil {
		return err
	}

	// 如果存在就更新，不存在就插入
	if exists {
		_, err = m.db.ID(ctx.UserID).Update(ctx)
	} else {
		_, err = m.db.Insert(ctx)
	}
	return err
}

func (m *ContextManager) shouldCompress(messages []Message) bool {
	// 当消息数量超过20条或总token数超过限制时，触发压缩
	return len(messages) > 20 || estimateTokens(messages) > m.tokenConfig.MaxTokensForRecent
}

func (m *ContextManager) compressShortTermMemory(ctx *ChatContext) error {
	if len(ctx.ShortTermMemory) == 0 {
		return nil
	}

	// 创建新的摘要
	summary := Summary{
		StartTime: ctx.ShortTermMemory[0].Time,
		EndTime:   ctx.ShortTermMemory[len(ctx.ShortTermMemory)-1].Time,
		Topics:    extractTopics(ctx.ShortTermMemory),
		KeyPoints: extractKeyPoints(ctx.ShortTermMemory),
	}

	// 添加到中期记忆
	ctx.MediumTermMemory = append(ctx.MediumTermMemory, summary)

	// 只保留最近的5条消息
	if len(ctx.ShortTermMemory) > 5 {
		ctx.ShortTermMemory = ctx.ShortTermMemory[len(ctx.ShortTermMemory)-5:]
	}

	return nil
}

// 辅助函数

func isRelevant(msg Message, topic string) bool {
	// TODO: 实现相关性判断逻辑
	return true
}

func containsTopic(topics []string, topic string) bool {
	for _, t := range topics {
		if t == topic {
			return true
		}
	}
	return false
}

func estimateTokens(messages []Message) int {
	// TODO: 实现更准确的token计算
	tokens := 0
	for _, msg := range messages {
		tokens += len(msg.Content) / 4 // 粗略估计
	}
	return tokens
}

func extractTopics(messages []Message) []string {
	// TODO: 实现主题提取逻辑
	return []string{}
}

func extractKeyPoints(messages []Message) []string {
	// TODO: 实现关键点提取逻辑
	return []string{}
}
