package queue

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
	"time"
)

const (
	// 消息队列的key前缀
	MessageQueuePrefix = "mq:msg:"
	// 消息ID生成器的key
	MessageIDKey = "mq:msg:id"
	// 消息过期时间（7天）
	MessageExpiration = 7 * 24 * time.Hour
)

// Message 消息队列中的消息结构
type Message struct {
	ID        int64     `json:"id"`         // 消息ID
	From      int64     `json:"from"`       // 发送者ID
	To        int64     `json:"to"`         // 接收者ID
	Type      uint8     `json:"type"`       // 消息类型
	Content   []byte    `json:"content"`    // 消息内容
	Timestamp int64     `json:"timestamp"`  // 消息时间戳
	Status    int       `json:"status"`     // 消息状态
	CreatedAt time.Time `json:"created_at"` // 创建时间
}

// MessageQueue 消息队列管理器
type MessageQueue struct {
	rdb *redis.Client
}

// NewMessageQueue 创建消息队列管理器
func NewMessageQueue(rdb *redis.Client) *MessageQueue {
	return &MessageQueue{
		rdb: rdb,
	}
}

// PublishMessage 发布消息到队列
func (mq *MessageQueue) PublishMessage(ctx context.Context, msg *Message) error {
	// 1. 生成消息ID
	id, err := mq.rdb.Incr(ctx, MessageIDKey).Result()
	if err != nil {
		return fmt.Errorf("生成消息ID失败: %v", err)
	}
	msg.ID = id
	msg.CreatedAt = time.Now()

	// 2. 序列化消息
	data, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("序列化消息失败: %v", err)
	}

	// 3. 发布到用户的消息队列
	key := fmt.Sprintf("%s%d", MessageQueuePrefix, msg.To)
	err = mq.rdb.LPush(ctx, key, data).Err()
	if err != nil {
		return fmt.Errorf("发布消息失败: %v", err)
	}

	// 4. 设置过期时间
	mq.rdb.Expire(ctx, key, MessageExpiration)

	return nil
}

// ConsumeMessage 消费消息
func (mq *MessageQueue) ConsumeMessage(ctx context.Context, userID int64) (*Message, error) {
	key := fmt.Sprintf("%s%d", MessageQueuePrefix, userID)

	// 从队列右端弹出消息
	data, err := mq.rdb.RPop(ctx, key).Bytes()
	if errors.Is(err, redis.Nil) {
		return nil, nil // 队列为空
	}
	if err != nil {
		return nil, fmt.Errorf("获取消息失败: %v", err)
	}

	// 反序列化消息
	var msg Message
	if err := json.Unmarshal(data, &msg); err != nil {
		return nil, fmt.Errorf("反序列化消息失败: %v", err)
	}

	return &msg, nil
}

// GetPendingMessageCount 获取待处理消息数量
func (mq *MessageQueue) GetPendingMessageCount(ctx context.Context, userID int64) (int64, error) {
	key := fmt.Sprintf("%s%d", MessageQueuePrefix, userID)
	return mq.rdb.LLen(ctx, key).Result()
}

// BatchConsumeMessages 批量获取消息
func (mq *MessageQueue) BatchConsumeMessages(ctx context.Context, userID int64, count int64) ([]*Message, error) {
	key := fmt.Sprintf("%s%d", MessageQueuePrefix, userID)

	// 获取多条消息
	pipeline := mq.rdb.Pipeline()
	for i := int64(0); i < count; i++ {
		pipeline.RPop(ctx, key)
	}

	cmds, err := pipeline.Exec(ctx)
	if err != nil {
		return nil, fmt.Errorf("批量获取消息失败: %v", err)
	}

	messages := make([]*Message, 0, len(cmds))
	for _, cmd := range cmds {
		if cmd == nil {
			continue
		}

		popCmd := cmd.(*redis.StringCmd)
		data, err_pop := popCmd.Bytes()
		if errors.Is(err_pop, redis.Nil) {
			continue
		}
		if err_pop != nil {
			return nil, fmt.Errorf("获取消息数据失败: %v", err)
		}

		var msg Message
		if err_marshal := json.Unmarshal(data, &msg); err_marshal != nil {
			return nil, fmt.Errorf("反序列化消息失败: %v", err_marshal)
		}
		messages = append(messages, &msg)
	}

	return messages, nil
}
