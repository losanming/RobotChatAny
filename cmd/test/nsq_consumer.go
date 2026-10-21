package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"new-robot-self/pkg/nsqmq"
)

type ImgMsg struct {
	ImgID        string  `json:"img_id"`
	TS           int64   `json:"ts"`
	DeviceID     string  `json:"device_id"`
	URLOrKey     string  `json:"url_or_key"`
	QualityScore float64 `json:"quality_score"`
	Hash         string  `json:"hash"`
}

func main() {
	var (
		lookupd   = "127.0.0.1:24161"
		nsqd      = "127.0.0.1:24150"
		topic     = "screenshot.ingest"
		channel   = "ai"
		workers   = 8
		maxFlight = 100
		maxAtt    = 5
		dlq       = "screenshot.dlq"
	)

	cfg := nsqmq.NSQConfig{
		NSQDAddress:      nsqd, // 若留空则使用 lookupd 连接
		LookupdAddresses: []string{},
		Topic:            topic,
		Channel:          channel,
		MaxInFlight:      maxFlight,
		MaxAttempts:      uint16(maxAtt),
		MsgTimeout:       30 * time.Second,
		DLQTopic:         dlq,
		ClientID:         "nsq-consumer",
		Hostname:         "nsq-consumer",
	}
	if lookupd != "" {
		cfg.LookupdAddresses = []string{lookupd}
	}

	client, err := nsqmq.NewClient(cfg)
	if err != nil {
		log.Fatalf("init nsq client failed: %v", err)
	}
	defer client.Close()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// 优雅退出
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sig
		fmt.Println("\nshutting down...")
		cancel()
	}()

	handler := func(b []byte) error {
		var m ImgMsg
		if err := json.Unmarshal(b, &m); err != nil {
			log.Printf("invalid message: %v", err)
			return err
		}
		// 模拟处理（取图 -> 调用 AI -> 落库），这里只打印
		log.Printf("consume img_id=%s device=%s ts=%d url_or_key=%s score=%.2f attempts<=%d",
			m.ImgID, m.DeviceID, m.TS, m.URLOrKey, m.QualityScore, maxAtt)
		// 示例：轻微延迟，体现并发
		time.Sleep(50 * time.Millisecond)
		return nil
	}

	err = client.ConsumeMessagesCtx(ctx, nsqmq.ConsumerConfig{
		Topic:       topic,
		Channel:     channel,
		WorkerCount: workers,
		MaxInFlight: maxFlight,
		MaxAttempts: uint16(maxAtt),
		DLQTopic:    dlq,
		MsgTimeout:  30 * time.Second,
	}, handler)
	if err != nil {
		log.Fatalf("start consumer failed: %v", err)
	}

	// 阻塞直到退出
	<-ctx.Done()
}
