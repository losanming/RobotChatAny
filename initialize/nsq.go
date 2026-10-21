package initialize

import (
	"fmt"
	"new-robot-self/global"
	nsqManager "new-robot-self/internal/nsq"
	"new-robot-self/public"
	public_nsq "new-robot-self/public/nsq"
	"time"

	"new-robot-self/pkg/nsqmq"
)

// InitNSQConsumer 初始化 NSQ 消费者并开始监听图片消息
// 环境变量（可选）：
//
//	NSQ_NSQD        如 101.126.17.221:4150（直连 nsqd）
//	NSQ_LOOKUPD     如 101.126.17.221:4161（通过 lookupd 发现）
//	NSQ_TOPIC       默认 screenshot.ingest
//	NSQ_CHANNEL     默认 ai
//	NSQ_DLQ         默认 screenshot.dlq
//	NSQ_WORKERS     并发数，默认 8（可不设置）
//	NSQ_MAX_INFLIGHT 默认 100
//	NSQ_MAX_ATTEMPTS 默认 5
//	NSQ_MSG_TIMEOUT 秒，默认 30
func InitNSQConsumer() error {

	cfg := nsqmq.NSQConfig{
		NSQDAddress:      public.Config.Nsq.Nsqd,
		LookupdAddresses: public.Config.Nsq.Nsqlookupd,
		Topic:            public.Config.Nsq.Topic,
		Channel:          public.Config.Nsq.Channel,
		MaxInFlight:      public.Config.Nsq.MaxInFlight,
		MaxAttempts:      public.Config.Nsq.MaxAttempts,
		MsgTimeout:       time.Duration(public.Config.Nsq.MsgTimeout) * time.Second,
		DLQTopic:         public.Config.Nsq.Dlq,
		ClientID:         "nsq-consumer",
		Hostname:         global.HostName,
	}

	client, err := nsqmq.NewClient(fallbackNSQD(cfg))
	if err != nil {
		return fmt.Errorf("init nsq client failed: %w", err)
	}
	if client == nil {
		return fmt.Errorf("nsq client is nil")
	}
	public_nsq.NsqClient = client
	err = nsqManager.StartNsqConsumer()
	if err != nil {
		return fmt.Errorf("start nsq consumer failed: %w", err)
	}
	return nil
}

// InitNSQ 兼容项目现有初始化流水线：无上下文参数的初始化入口
func InitNSQ() error {
	if public.Config.Nsq.Off {
		return InitNSQConsumer()
	}
	fmt.Printf("nsq 消息队列 初始化配置成功 \n")
	return nil
}

// fallbackNSQD: 当仅配置了 lookupd（生产常用）但 nsqmq.NewClient 需要 NSQDAddress 时，给一个占位地址避免报错。
// 注意：Producer 才强依赖 NSQDAddress，这里我们用它只为复用 Client 结构（其中包含连接配置）。
func fallbackNSQD(cfg nsqmq.NSQConfig) nsqmq.NSQConfig {
	c := cfg
	if c.NSQDAddress == "" {
		c.NSQDAddress = "127.0.0.1:4150" // 不会被消费端连接使用，仅为通过构造校验
	}
	return c
}
