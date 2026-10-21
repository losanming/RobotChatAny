package kafka

import (
	"context"
	"fmt"
	kafka "github.com/segmentio/kafka-go"
	"log"
	"sync"
	"time"
)

// KafkaConfig Kafka配置
type KafkaConfig struct {
	Brokers       []string
	Topic         string
	ConsumerGroup string
	// 主题配置（可选，EnsureTopicAndConsumerGroup 使用）
	Partitions        int
	ReplicationFactor int
	// Writer 批量与超时（可选）
	BatchSize    int
	BatchTimeout time.Duration
	WriteTimeout time.Duration
	// 失败消息投递的DLQ主题（可选）
	DLQTopic string
}

// Message Kafka消息结构体
type Message struct {
	Key       []byte `json:"key"`       // 消息Key
	Value     []byte `json:"value"`     // 消息内容
	Timestamp int64  `json:"timestamp"` // 时间戳
}

// ConsumerConfig 消费者配置
type ConsumerConfig struct {
	QueueName    string        // 队列名称（主题）
	WorkerCount  int           // 并发工作协程数
	GroupID      string        // 消费者组ID
	MaxRetries   int           // 处理失败最大重试次数（默认3）
	RetryBackoff time.Duration // 初始重试退避（默认200ms，指数退避）
	DLQTopic     string        // 失败消息投递的DLQ主题（可选）
}

// KafkaClient Kafka客户端
type KafkaClient struct {
	writer *kafka.Writer
	config KafkaConfig
	mu     sync.Mutex
}

// SendOptions 发送选项
// - TextChannel: 为true时，按文本通道策略发送（Base64+Smart压缩），并通过Header标注编码信息
// - MinSize/MinGainPct: Smart压缩阈值，留空则使用默认(1024, 5%)
type SendOptions struct {
	TextChannel bool
	MinSize     int
	MinGainPct  float64
}

// NewKafkaClient 创建新的Kafka客户端
func NewKafkaClient(config KafkaConfig) (*KafkaClient, error) {
	if len(config.Brokers) == 0 {
		return nil, fmt.Errorf("未配置 Kafka Brokers")
	}

	// 默认批量与超时
	if config.BatchSize == 0 {
		config.BatchSize = 100
	}
	if config.BatchTimeout == 0 {
		config.BatchTimeout = 200 * time.Millisecond
	}
	if config.WriteTimeout == 0 {
		config.WriteTimeout = 5 * time.Second
	}

	writer := &kafka.Writer{
		Addr:                   kafka.TCP(config.Brokers...),
		Topic:                  config.Topic,
		Balancer:               &kafka.Hash{},
		RequiredAcks:           kafka.RequireAll,
		Async:                  false,
		Compression:            kafka.Zstd,
		BatchSize:              config.BatchSize,
		BatchTimeout:           config.BatchTimeout,
		WriteTimeout:           config.WriteTimeout,
		AllowAutoTopicCreation: false,
	}

	// 创建客户端实例
	client := &KafkaClient{
		writer: writer,
		config: config,
	}

	// 确保 Topic 和消费组存在
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := client.EnsureTopicAndConsumerGroup(ctx); err != nil {
		// 如果初始化失败，关闭 writer 并返回错误
		writer.Close()
		return nil, fmt.Errorf("初始化Topic和消费组失败: %v", err)
	}

	// 返回客户端实例
	return client, nil
}

// Close 关闭Kafka客户端
func (k *KafkaClient) Close() {
	k.mu.Lock()
	defer k.mu.Unlock()

	if k.writer != nil {
		k.writer.Close()
	}
}

// PublishMessage 发布消息
func (k *KafkaClient) PublishMessage(topic string, message []byte) error {
	k.mu.Lock()
	defer k.mu.Unlock()

	ctx := context.Background()

	// 如果指定了不同的topic，创建一个临时writer
	var writer *kafka.Writer
	if topic != "" && topic != k.config.Topic {
		writer = &kafka.Writer{
			Addr:                   kafka.TCP(k.config.Brokers...),
			Topic:                  topic,
			Balancer:               &kafka.Hash{},
			RequiredAcks:           kafka.RequireAll,
			Async:                  false,
			Compression:            kafka.Zstd,
			BatchSize:              k.config.BatchSize,
			BatchTimeout:           k.config.BatchTimeout,
			WriteTimeout:           k.config.WriteTimeout,
			AllowAutoTopicCreation: false,
		}
		defer writer.Close()
	} else {
		writer = k.writer
		topic = k.config.Topic
	}

	// 发送消息
	err := writer.WriteMessages(ctx, kafka.Message{
		Value: message,
		Time:  time.Now(),
	})

	if err != nil {
		return fmt.Errorf("发送消息失败: %v", err)
	}

	return nil
}

// PublishMessageWithKey 发布带Key的消息
func (k *KafkaClient) PublishMessageWithKey(topic string, key, message []byte) error {
	k.mu.Lock()
	defer k.mu.Unlock()

	ctx := context.Background()

	// 如果指定了不同的topic，创建一个临时writer
	var writer *kafka.Writer
	if topic != "" && topic != k.config.Topic {
		writer = &kafka.Writer{
			Addr:                   kafka.TCP(k.config.Brokers...),
			Topic:                  topic,
			Balancer:               &kafka.Hash{},
			RequiredAcks:           kafka.RequireAll,
			Async:                  false,
			Compression:            kafka.Zstd,
			BatchSize:              k.config.BatchSize,
			BatchTimeout:           k.config.BatchTimeout,
			WriteTimeout:           k.config.WriteTimeout,
			AllowAutoTopicCreation: false,
		}
		defer writer.Close()
	} else {
		writer = k.writer
		topic = k.config.Topic
	}

	// 发送消息
	err := writer.WriteMessages(ctx, kafka.Message{
		Key:   key,
		Value: message,
		Time:  time.Now(),
	})

	if err != nil {
		return fmt.Errorf("发送消息失败: %v", err)
	}

	return nil
}

// publishDLQ 投递失败消息到DLQ
func (k *KafkaClient) publishDLQ(ctx context.Context, cfg ConsumerConfig, msg *kafka.Message, err error) error {
	if cfg.DLQTopic == "" {
		return nil
	}
	// 构造一个包含错误信息的消息（保留原Key/Value/部分Header）
	headers := append(msg.Headers, kafka.Header{Key: "error", Value: []byte(err.Error())})

	w := &kafka.Writer{
		Addr:                   kafka.TCP(k.config.Brokers...),
		Topic:                  cfg.DLQTopic,
		Balancer:               &kafka.Hash{},
		RequiredAcks:           kafka.RequireAll,
		Async:                  false,
		Compression:            kafka.Zstd,
		BatchSize:              k.config.BatchSize,
		BatchTimeout:           k.config.BatchTimeout,
		WriteTimeout:           k.config.WriteTimeout,
		AllowAutoTopicCreation: false,
	}
	defer w.Close()
	return w.WriteMessages(ctx, kafka.Message{
		Key:     msg.Key,
		Value:   msg.Value,
		Time:    time.Now(),
		Headers: headers,
	})
}

// ConsumeMessagesCtx 消费消息（并发处理，支持上下文与手动提交、重试与DLQ）
func (k *KafkaClient) ConsumeMessagesCtx(ctx context.Context, config ConsumerConfig, handleMessage func([]byte) error) error {
	if config.WorkerCount <= 0 {
		config.WorkerCount = 1
	}

	if config.QueueName == "" {
		config.QueueName = k.config.Topic
	}

	if config.GroupID == "" {
		config.GroupID = k.config.ConsumerGroup
	}

	// 创建消费者
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:        k.config.Brokers,
		GroupID:        config.GroupID,
		Topic:          config.QueueName,
		MinBytes:       10e3,             // 10KB
		MaxBytes:       10e6,             // 10MB
		CommitInterval: 0,                // 手动提交
		StartOffset:    kafka.LastOffset, // 首次消费从最新开始
	})

	// 创建工作协程池
	var wg sync.WaitGroup
	msgChan := make(chan kafka.Message, config.WorkerCount*2)
	stopChan := make(chan struct{})

	// 启动工作协程
	for i := 0; i < config.WorkerCount; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			log.Printf("工作协程[%d]已启动", workerID)

			for {
				select {
				case msg := <-msgChan:
					// 处理消息（重试+提交）
					var err error
					retries := config.MaxRetries
					if retries <= 0 {
						retries = 3
					}
					backoff := config.RetryBackoff
					if backoff <= 0 {
						backoff = 200 * time.Millisecond
					}
					for attempt := 0; attempt <= retries; attempt++ {
						err = handleMessage(msg.Value)
						if err == nil {
							// 成功后提交偏移
							if commitErr := reader.CommitMessages(ctx, msg); commitErr != nil {
								log.Printf("提交偏移失败: %v", commitErr)
							}
							break
						}
						if attempt < retries {
							time.Sleep(backoff)
							backoff *= 2
						}
					}
					if err != nil {
						log.Printf("消息处理失败，投递DLQ: %v", err)
						_ = k.publishDLQ(ctx, config, &msg, err)
						// 失败也提交，避免卡住分区（根据业务需求可改）
						if commitErr := reader.CommitMessages(ctx, msg); commitErr != nil {
							log.Printf("提交偏移失败(失败场景): %v", commitErr)
						}
					}

				case <-stopChan:
					log.Printf("工作协程[%d]收到停止信号", workerID)
					return
				case <-ctx.Done():
					return
				}
			}
		}(i)
	}

	// 分发消息到工作协程
	go func() {
		defer reader.Close()
		for {
			select {
			case <-stopChan:
				log.Println("消息分发协程收到停止信号")
				close(msgChan)
				return
			case <-ctx.Done():
				close(msgChan)
				return
			default:
				msg, err := reader.ReadMessage(ctx)
				if err != nil {
					log.Printf("读取消息失败: %v", err)
					time.Sleep(100 * time.Millisecond)
					continue
				}

				// 将消息发送到通道
				msgChan <- msg
			}
		}
	}()

	// 由调用方控制 ctx 取消；这里返回一个阻塞等待的函数由调用方选择。
	// 为保持兼容，也提供不阻塞的返回，调用方可选择自行管理生命周期。
	return nil
}

// ConsumeMessages 兼容旧签名（使用 context.Background()），建议改用 ConsumeMessagesCtx
func (k *KafkaClient) ConsumeMessages(config ConsumerConfig, handleMessage func([]byte) error) error {
	return k.ConsumeMessagesCtx(context.Background(), config, handleMessage)
}

// ConsumeMessagesWithCallback 消费消息（带完整回调）
func (k *KafkaClient) ConsumeMessagesWithCallback(config ConsumerConfig, callback func(*kafka.Message, error)) error {
	if config.WorkerCount <= 0 {
		config.WorkerCount = 1
	}

	if config.QueueName == "" {
		config.QueueName = k.config.Topic
	}

	if config.GroupID == "" {
		config.GroupID = k.config.ConsumerGroup
	}

	// 创建消费者
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:        k.config.Brokers,
		GroupID:        config.GroupID,
		Topic:          config.QueueName,
		MinBytes:       10e3, // 10KB
		MaxBytes:       10e6, // 10MB
		CommitInterval: 1 * time.Second,
		StartOffset:    kafka.FirstOffset, // 从最早的消息开始
	})

	// 创建工作协程池
	var wg sync.WaitGroup
	msgChan := make(chan kafka.Message, config.WorkerCount*2)
	errChan := make(chan error, config.WorkerCount*2)
	stopChan := make(chan struct{})

	// 启动工作协程
	for i := 0; i < config.WorkerCount; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()

			for {
				select {
				case msg := <-msgChan:
					if callback != nil {
						// 将segmentio/kafka-go的Message转换为我们自己的Message格式
						kafkaMsg := &kafka.Message{
							Topic:   msg.Topic,
							Key:     msg.Key,
							Value:   msg.Value,
							Time:    msg.Time,
							Headers: msg.Headers,
						}
						callback(kafkaMsg, nil)
					}

				case err := <-errChan:
					if callback != nil {
						callback(nil, err)
					}

				case <-stopChan:
					return
				}
			}
		}(i)
	}

	// 分发消息到工作协程
	go func() {
		defer reader.Close()

		ctx := context.Background()
		for {
			select {
			case <-stopChan:
				close(msgChan)
				close(errChan)
				return
			default:
				msg, err := reader.ReadMessage(ctx)
				if err != nil {
					errChan <- err
					time.Sleep(100 * time.Millisecond)
					continue
				}

				// 将消息发送到通道
				msgChan <- msg
			}
		}
	}()
	return nil
}

// EnsureTopicAndConsumerGroup 确保 Topic 和消费组存在，如果不存在则创建
func (k *KafkaClient) EnsureTopicAndConsumerGroup(ctx context.Context) error {
	// 检查并创建 Topic
	if k.config.Topic != "" {
		conn, err := kafka.DialContext(ctx, "tcp", k.config.Brokers[0])
		if err != nil {
			return fmt.Errorf("连接Kafka失败: %v", err)
		}
		defer conn.Close()

		// 读取现有分区
		partitions, readErr := conn.ReadPartitions(k.config.Topic)
		if readErr == nil && len(partitions) > 0 {
			log.Printf("Topic %s 已存在，包含 %d 个分区", k.config.Topic, len(partitions))
		} else {
			// 通过 controller 创建 topic
			controller, cerr := conn.Controller()
			if cerr != nil {
				return fmt.Errorf("获取Controller失败: %v", cerr)
			}
			c, derr := kafka.DialContext(ctx, "tcp", fmt.Sprintf("%s:%d", controller.Host, controller.Port))
			if derr != nil {
				return fmt.Errorf("连接Controller失败: %v", derr)
			}
			defer c.Close()

			np := k.config.Partitions
			if np <= 0 {
				np = 3
			}
			rf := k.config.ReplicationFactor
			if rf <= 0 {
				rf = 1
			}
			err = c.CreateTopics(kafka.TopicConfig{
				Topic:             k.config.Topic,
				NumPartitions:     np,
				ReplicationFactor: rf,
			})
			if err != nil {
				return fmt.Errorf("创建Topic失败: %v", err)
			}
			log.Printf("Topic %s 创建成功", k.config.Topic)
		}
	}

	// 消费组不需要预先创建，它会在首次消费时自动创建
	// 但我们可以记录一下
	if k.config.ConsumerGroup != "" {
		log.Printf("消费组 %s 将在首次消费时自动创建", k.config.ConsumerGroup)
	}

	return nil
}
