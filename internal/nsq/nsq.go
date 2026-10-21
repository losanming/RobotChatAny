package nsqManager

import (
	"encoding/json"
	"fmt"
	"log"
	"new-robot-self/internal/queue"
	"new-robot-self/pkg/nsqmq"
	"new-robot-self/pkg/upload"
	"new-robot-self/public"
	public_connectionManager "new-robot-self/public/connectionManager"
	public_nsq "new-robot-self/public/nsq"
	"time"
)

// ImgMsg 与测试用消息体一致：指针 + 元数据（时间戳）
type ImgMsg struct {
	ImgID        string  `json:"img_id"`
	TS           int64   `json:"ts"`
	DeviceID     string  `json:"device_id"`
	From         int64   `json:"from"`
	URLOrKey     string  `json:"url_or_key"`
	QualityScore float64 `json:"quality_score"`
	Hash         string  `json:"hash"`
}

// nsq msgHandler
func NsqMsgHandler() func(b []byte) error {
	handler := func(b []byte) error {
		var m ImgMsg
		if err := json.Unmarshal(b, &m); err != nil {
			log.Printf("invalid message: %v", err)
			return err
		}
		// 这里先实现到：拿到图片 URL（或 Key）
		// 先排除掉垃圾图片 差值为1分钟的全部不处理
		//@@TODO 如果from不在线也不处理

		if time.Now().Unix()-m.TS < 60 {
			return nil
		}
		if m.ImgID == "" || m.URLOrKey == "" {
			return nil
		}
		// 从华为云获取
		oss := upload.NewOss(public.Config.Sys.Oss)
		if oss == nil {
			log.Println("OSS is nil")
			return nil
		}

		byByte, err := oss.DownloadFileByByte(m.URLOrKey)
		if err != nil {
			log.Println("DownloadFileByByte err", "err", err)
			return err
		}

		if len(byByte) > 0 {
			//@@TODO 图片数据移交ai图片理解

			// ai理解结果
			var msg queue.Message
			msg.Type = 2
			msg.ID = time.Now().Unix()
			msg.From = 121
			msg.To = m.From
			msg.Content = nil
			msg.Status = 1
			msg.Timestamp = time.Now().Unix()

			err_send := public_connectionManager.Connection.SendMessage(public.Ctx, &msg)
			if err_send != nil {
				log.Println("SendMessage err", "err", err)
				return err_send
			}
		}

		return nil
	}
	return handler
}

func StartNsqConsumer() error {
	if public.Config.Nsq.Off {
		log.Println("NSQ is off")
		return nil
	}
	log.Println("Starting NSQ consumer...")
	// 开始消费（并发）
	consCfg := nsqmq.ConsumerConfig{
		Topic:       public.Config.Nsq.Topic,
		Channel:     public.Config.Nsq.Channel,
		WorkerCount: public.Config.Nsq.Workers,
		MaxInFlight: public.Config.Nsq.MaxInFlight,
		MaxAttempts: public.Config.Nsq.MaxAttempts,
		MsgTimeout:  time.Duration(public.Config.Nsq.MsgTimeout) * time.Second,
		DLQTopic:    public.Config.Nsq.Dlq,
	}
	if err := public_nsq.NsqClient.ConsumeMessagesCtx(public.Ctx, consCfg, NsqMsgHandler()); err != nil {
		return fmt.Errorf("start consumer failed: %w", err)
	}

	log.Printf("NSQ consumer started: topic=%s channel=%s lookupd=%s nsqd=%s", public.Config.Nsq.Topic, public.Config.Nsq.Channel, public.Config.Nsq.Nsqlookupd, public.Config.Nsq.Nsqd)
	return nil
}
