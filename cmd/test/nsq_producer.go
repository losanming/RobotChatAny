package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/google/uuid"
	"new-robot-self/pkg/nsqmq"
)

type ImgMsTO struct {
	ImgID        string  `json:"img_id"`
	TS           int64   `json:"ts"`
	DeviceID     string  `json:"device_id"`
	URLOrKey     string  `json:"url_or_key"`
	QualityScore float64 `json:"quality_score"`
	Hash         string  `json:"hash"`
}

func main() {
	var (
		nsqd     = "101.126.17.221:24150"
		topic    = "screenshot.ingest"
		count    = flag.Int("count", 10, "number of messages to publish")
		deviceID = flag.String("device", "dev-001", "device id in message")
		baseURL  = flag.String("base", "obs://bucket/path/", "base url_or_key prefix (e.g., obs:// or redis key prefix)")
		qps      = flag.Int("qps", 20, "limit publish QPS (approximate)")
	)
	flag.Parse()

	cfg := nsqmq.NSQConfig{
		NSQDAddress: nsqd,
		Topic:       topic,
		MaxAttempts: 5,
		MaxInFlight: 100,
		MsgTimeout:  30 * time.Second,
		ClientID:    "nsq-producer",
		Hostname:    "nsq-producer",
	}

	client, err := nsqmq.NewClient(cfg)
	if err != nil {
		log.Fatalf("init producer failed: %v", err)
	}
	defer client.Close()

	interval := time.Second
	if *qps > 0 {
		interval = time.Second / time.Duration(*qps)
	}
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for i := 0; i < *count; i++ {
		<-ticker.C
		id := uuid.New().String()
		msg := ImgMsTO{
			ImgID:        id,
			TS:           time.Now().Unix(),
			DeviceID:     *deviceID,
			URLOrKey:     strings.TrimRight(*baseURL, "/") + "/" + id + ".jpg",
			QualityScore: 0.8,
			Hash:         "aHash:" + id[:8],
		}
		b, _ := json.Marshal(msg)
		if err := client.PublishMessage(topic, b); err != nil {
			log.Printf("publish failed: %v", err)
		} else {
			fmt.Printf("published: %s\n", msg.ImgID)
		}
	}

	fmt.Println("producer done")
}
