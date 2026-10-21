package queue

import (
	"context"
	"fmt"
	"log"
	"new-robot-self/internal/protocol"
	"os"
	"sync"
	"time"
)

// 音频数据消息
type AudioData struct {
	DeviceID   string `json:"device_id"`
	Timestamp  int64  `json:"timestamp"`
	Format     string `json:"format"`      // 音频格式 PCM/WAV
	SampleRate int    `json:"sample_rate"` // 采样率
	BitDepth   int    `json:"bit_depth"`   // 位深度
	Data       string `json:"data"`        // base64编码的音频数据
}

// MessageHandler 消息处理器接口
type MessageHandler interface {
	HandleMessage(ctx context.Context, msg *Message) error
}

// Dispatcher 消息分发器
type Dispatcher struct {
	mq          *MessageQueue
	handlers    map[uint8]MessageHandler
	workerCount int
	wg          sync.WaitGroup
	ctx         context.Context
	cancel      context.CancelFunc
	userQueues  *sync.Map   // 用户消息队列映射
	retryConfig RetryConfig // 重试配置
	errorLog    *log.Logger // 错误日志
}

// NewDispatcher 创建消息分发器
func NewDispatcher(mq *MessageQueue, workerCount int, retryConfig *RetryConfig) *Dispatcher {
	if workerCount <= 0 {
		workerCount = 10 // 默认10个工作协程
	}

	// 使用默认重试配置如果没有提供
	if retryConfig == nil {
		retryConfig = &DefaultRetryConfig
	}

	ctx, cancel := context.WithCancel(context.Background())
	return &Dispatcher{
		mq:          mq,
		handlers:    make(map[uint8]MessageHandler),
		workerCount: workerCount,
		ctx:         ctx,
		cancel:      cancel,
		userQueues:  &sync.Map{},
		retryConfig: *retryConfig,
		errorLog:    log.New(os.Stderr, "[Dispatcher] ", log.LstdFlags),
	}
}

// RegisterHandler 注册消息处理器
func (d *Dispatcher) RegisterHandler(msgType uint8, handler MessageHandler) {
	d.handlers[msgType] = handler
}

func (d *Dispatcher) GetMq() *MessageQueue {
	return d.mq
}

func (d *Dispatcher) GetCtx() context.Context {
	return d.ctx
}

// Start 启动消息分发系统
func (d *Dispatcher) Start() {
	// 启动多个工作协程
	for i := 0; i < d.workerCount; i++ {
		d.wg.Add(1)
		go func(workerID int) {
			defer d.wg.Done()
			d.dispatchLoop(workerID)
		}(i)
	}
}

// Stop 停止消息分发
func (d *Dispatcher) Stop() {
	d.cancel()
	d.wg.Wait()
}

// dispatchLoop 消息分发循环
func (d *Dispatcher) dispatchLoop(workerID int) {
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-d.ctx.Done():
			return
		case <-ticker.C:
			// 获取所有活跃用户
			activeUsers := d.getActiveUsers()

			// 遍历处理每个用户的消息
			for _, userID := range activeUsers {
				// 获取待处理消息数量
				count, err := d.mq.GetPendingMessageCount(d.ctx, userID)
				if err != nil {
					log.Printf("获取待处理消息数量失败 [UserID=%d]: %v", userID, err)
					continue
				}

				if count == 0 {
					continue
				}

				// 批量获取消息
				messages, err := d.mq.BatchConsumeMessages(d.ctx, userID, count)
				if err != nil {
					log.Printf("批量获取消息失败 [UserID=%d]: %v", userID, err)
					continue
				}

				// 处理消息
				for _, msg := range messages {
					if handler, ok := d.handlers[msg.Type]; ok {
						// 使用重试机制处理消息
						if err := d.retryMessageHandler(msg, handler); err != nil {
							d.errorLog.Printf("消息处理最终失败 [ID=%d, UserID=%d]: %v", msg.ID, userID, err)
							// TODO: 可以将失败的消息存储到死信队列
							//d.mq.PublishToDeadLetterQueue(d.ctx, msg)
						}
					} else {
						d.errorLog.Printf("未找到消息处理器 [Type=%d]", msg.Type)
					}
				}
			}
		}
	}
}

// getActiveUsers 获取活跃用户列表
func (d *Dispatcher) getActiveUsers() []int64 {
	var users []int64
	d.userQueues.Range(func(key, value interface{}) bool {
		userID := key.(int64)
		users = append(users, userID)
		return true
	})
	return users
}

// AddUser 添加用户到分发系统
func (d *Dispatcher) AddUser(userID int64) {
	d.userQueues.Store(userID, true)
}

// RemoveUser 从分发系统移除用户
func (d *Dispatcher) RemoveUser(userID int64) {
	d.userQueues.Delete(userID)
}

// retryMessageHandler 重试处理消息
func (d *Dispatcher) retryMessageHandler(msg *Message, handler MessageHandler) error {
	var lastErr error
	backoff := d.retryConfig.InitialWait

	// 重试循环
	for attempt := 0; attempt <= d.retryConfig.MaxRetries; attempt++ {
		// 处理消息
		err := handler.HandleMessage(d.ctx, msg)
		if err == nil {
			// 处理成功
			return nil
		}

		// 记录错误
		lastErr = err
		d.errorLog.Printf("消息处理失败 [MsgID=%d, Type=%d, Attempt=%d]: %v",
			msg.ID, msg.Type, attempt+1, err)

		// 如果错误不可重试，直接返回
		if !IsRetryable(err) {
			return fmt.Errorf("消息处理失败，不可重试: %v", err)
		}

		// 检查是否需要继续重试
		if attempt == d.retryConfig.MaxRetries {
			break
		}

		// 等待后重试
		select {
		case <-d.ctx.Done():
			return d.ctx.Err()
		case <-time.After(backoff):
			// 增加退避时间，但不超过最大等待时间
			backoff *= 2
			if backoff > d.retryConfig.MaxWait {
				backoff = d.retryConfig.MaxWait
			}
		}
	}

	// 所有重试都失败了
	return fmt.Errorf("消息处理失败，已重试%d次: %v",
		d.retryConfig.MaxRetries, lastErr)
}

// DefaultTextHandler 默认文本消息处理器
type DefaultTextHandler struct {
	codec *protocol.SecureMessageCodec
	// 这里可以添加连接管理器等其他依赖
}
