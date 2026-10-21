package nsqmq

import (
	"context"
	"errors"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/nsqio/go-nsq"
)

// NSQConfig NSQ 基础配置
// - NSQDAddress: 直连 nsqd 发布（如 "127.0.0.1:4150"）
// - LookupdAddresses: 通过 lookupd 发现 nsqd（消费方常用，"127.0.0.1:4161"）
// - Topic/Channel: 默认主题/通道（消费时可覆盖）
// - MaxInFlight: 每个消费者并行处理上限
// - MaxAttempts: 单条消息最大尝试次数（超过后认为失败）
// - DLQTopic: 失败消息投递的主题（可选）
// - ClientID/Hostname: 便于定位
// - MsgTimeout: 处理超时（超过后视为失败并重试）
// - RequeueBackoff: 失败后本地重试退避时间（仅用于我们手动 Requeue 的场景）
//
// 注意：NSQ 默认是“至少一次”投递，需在业务侧幂等。
// 没有内建的持久化回放语义，适合即产即消费。
//
type NSQConfig struct {
	NSQDAddress       string
	LookupdAddresses  []string
	Topic             string
	Channel           string
	MaxInFlight       int
	MaxAttempts       uint16
	DLQTopic          string
	ClientID          string
	Hostname          string
	MsgTimeout        time.Duration
	RequeueBackoff    time.Duration
}

// ConsumerConfig 消费者配置（可覆盖部分默认）
type ConsumerConfig struct {
	Topic        string
	Channel      string
	WorkerCount  int           // 并发处理数（AddConcurrentHandlers）
	MaxInFlight  int           // 覆盖并行上限
	MaxAttempts  uint16        // 覆盖最大尝试次数
	DLQTopic     string        // 覆盖 DLQ 主题
	MsgTimeout   time.Duration // 覆盖消息超时
}

// Client NSQ 客户端（包含 Producer；Consumer 按需创建）
type Client struct {
	producer *nsq.Producer
	cfg      NSQConfig
	mu       sync.Mutex
}

// NewClient 创建 NSQ 客户端（Producer）
func NewClient(cfg NSQConfig) (*Client, error) {
	if cfg.NSQDAddress == "" {
		return nil, fmt.Errorf("NSQDAddress 不能为空，例如 127.0.0.1:4150")
	}

	ncfg := nsq.NewConfig()
	if cfg.ClientID != "" {
		ncfg.ClientID = cfg.ClientID
	}
	if cfg.Hostname != "" {
		ncfg.Hostname = cfg.Hostname
	}
	if cfg.MaxInFlight > 0 {
		ncfg.MaxInFlight = cfg.MaxInFlight
	}
	if cfg.MaxAttempts > 0 {
		ncfg.MaxAttempts = cfg.MaxAttempts
	}
	if cfg.MsgTimeout > 0 {
		ncfg.MsgTimeout = cfg.MsgTimeout
	}

	p, err := nsq.NewProducer(cfg.NSQDAddress, ncfg)
	if err != nil {
		return nil, fmt.Errorf("创建 NSQ Producer 失败: %w", err)
	}
	// 可选：简单健康检查
	if err := p.Ping(); err != nil {
		p.Stop()
		return nil, fmt.Errorf("NSQ Producer Ping 失败: %w", err)
	}

	return &Client{producer: p, cfg: cfg}, nil
}

// Close 关闭 Producer
func (c *Client) Close() {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.producer != nil {
		c.producer.Stop()
	}
}

// PublishMessage 发布消息（指定 topic；留空则使用默认 cfg.Topic）
func (c *Client) PublishMessage(topic string, body []byte) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.producer == nil {
		return errors.New("producer 未初始化")
	}
	if topic == "" {
		topic = c.cfg.Topic
	}
	if topic == "" {
		return errors.New("未指定 topic")
	}
	return c.producer.Publish(topic, body)
}

// PublishBatch 批量发布
func (c *Client) PublishBatch(topic string, bodies [][]byte) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.producer == nil {
		return errors.New("producer 未初始化")
	}
	if topic == "" {
		topic = c.cfg.Topic
	}
	if topic == "" {
		return errors.New("未指定 topic")
	}
	return c.producer.MultiPublish(topic, bodies)
}

// ConsumeMessagesCtx 开始消费（并发处理、自动重试、可投递 DLQ）
// handleMessage: 返回 nil 表示成功；返回 error 表示失败，NSQ 会按 attempts 策略重试。
// 注意：NSQ 的重试由服务端 attempts 控制，超过 MaxAttempts 后我们将消息转发至 DLQ（若配置）。
func (c *Client) ConsumeMessagesCtx(parent context.Context, cc ConsumerConfig, handleMessage func([]byte) error) error {
	if handleMessage == nil {
		return errors.New("handleMessage 不能为空")
	}

	topic := cc.Topic
	if topic == "" {
		topic = c.cfg.Topic
	}
	channel := cc.Channel
	if channel == "" {
		channel = c.cfg.Channel
	}
	if topic == "" || channel == "" {
		return errors.New("topic/channel 不能为空")
	}

	ncfg := nsq.NewConfig()
	// 继承全局默认
	if c.cfg.MaxInFlight > 0 {
		ncfg.MaxInFlight = c.cfg.MaxInFlight
	}
	if c.cfg.MaxAttempts > 0 {
		ncfg.MaxAttempts = c.cfg.MaxAttempts
	}
	if c.cfg.MsgTimeout > 0 {
		ncfg.MsgTimeout = c.cfg.MsgTimeout
	}
	// 局部覆盖
	if cc.MaxInFlight > 0 {
		ncfg.MaxInFlight = cc.MaxInFlight
	}
	if cc.MaxAttempts > 0 {
		ncfg.MaxAttempts = cc.MaxAttempts
	}
	if cc.MsgTimeout > 0 {
		ncfg.MsgTimeout = cc.MsgTimeout
	}

	consumer, err := nsq.NewConsumer(topic, channel, ncfg)
	if err != nil {
		return fmt.Errorf("创建 NSQ Consumer 失败: %w", err)
	}

	// 处理逻辑：
	// - 调用 handleMessage
	// - 返回 nil 则 Finish
	// - 返回 error：如果还没超过 MaxAttempts，让 NSQ 重试；超过后（or 达到 attempts 上限）转发 DLQ 再 Finish
	consumer.SetLogger(log.Default(), nsq.LogLevelWarning)

	// WorkerCount 默认 1
	wc := cc.WorkerCount
	if wc <= 0 {
		wc = 1
	}

	var dlqTopic string
	if cc.DLQTopic != "" {
		dlqTopic = cc.DLQTopic
	} else {
		dlqTopic = c.cfg.DLQTopic
	}

	h := &nsqHandler{
		handle:   handleMessage,
		client:   c,
		dlqTopic: dlqTopic,
		maxAttempts: func() uint16 {
			if cc.MaxAttempts > 0 {
				return cc.MaxAttempts
			}
			if c.cfg.MaxAttempts > 0 {
				return c.cfg.MaxAttempts
			}
			return 5 // 兜底：5 次视为最终失败
		}(),
	}
	consumer.AddConcurrentHandlers(h, wc)

	// 连接方式：优先 lookupd，否则直连 nsqd
	var connectErr error
	if len(c.cfg.LookupdAddresses) > 0 {
		connectErr = consumer.ConnectToNSQLookupds(c.cfg.LookupdAddresses)
	} else if c.cfg.NSQDAddress != "" {
		connectErr = consumer.ConnectToNSQD(c.cfg.NSQDAddress)
	} else {
		connectErr = fmt.Errorf("必须配置 LookupdAddresses 或 NSQDAddress 其一用于消费端连接")
	}
	if connectErr != nil {
		consumer.Stop()
		return fmt.Errorf("消费者连接失败: %w", connectErr)
	}

	// 生命周期：监听 ctx 取消
	go func() {
		<-parent.Done()
		consumer.Stop()
	}()

	return nil
}

// nsqHandler 封装业务处理 + 失败转 DLQ
// 注意：在 NSQ 中，每条消息自带 Attempts 计数，由服务端维护。
// 当业务返回 error 时，不调用 Finish，NSQ 会自动安排重试；
// 当达到 attempts 上限时（NSQ 将再次投递），我们将其视为最终失败，转发 DLQ 并 Finish。
type nsqHandler struct {
	handle      func([]byte) error
	client      *Client
	dlqTopic    string
	maxAttempts uint16
}

func (h *nsqHandler) HandleMessage(m *nsq.Message) error {
	// 业务处理
	if err := h.handle(m.Body); err != nil {
		// 达到最大重试次数后转发 DLQ
		if h.dlqTopic != "" && h.maxAttempts > 0 && m.Attempts >= h.maxAttempts {
			// 最终失败，转 DLQ（忽略错误，以免阻塞）
			_ = h.client.PublishMessage(h.dlqTopic, m.Body)
			m.Finish()
			return nil
		}
		// 让 NSQ 按默认策略重试（不调用 Finish 即可）。
		return err
	}
	// 成功
	m.Finish()
	return nil
}
