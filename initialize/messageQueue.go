package initialize

import (
	"errors"
	"fmt"
	"new-robot-self/internal/queue"
	"new-robot-self/public"
)

func InitQueueAndDispatcher() error {
	messageQueue := queue.NewMessageQueue(public.Cache.GetRedisClient())
	if messageQueue == nil {
		return errors.New("messageQueue is nil")
	}
	dispatcher := queue.NewDispatcher(messageQueue, 10, &queue.DefaultRetryConfig)
	if dispatcher == nil {
		return errors.New("dispatcher is nil")
	}
	dispatcher.Start()
	public.Dispatcher = dispatcher
	public.MessageQueue = messageQueue
	fmt.Printf("queue 初始化配置成功 \n")
	return nil
}
