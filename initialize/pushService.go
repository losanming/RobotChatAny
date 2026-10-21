package initialize

import (
	"errors"
	"fmt"
	"new-robot-self/internal/push"
	"new-robot-self/public"
	public_ai "new-robot-self/public/ai"
	public_connectionManager "new-robot-self/public/connectionManager"
)

func InitPushService() error {
	service := push.NewService(public_connectionManager.Connection, public.Trigger, public.Fetcher, public_ai.AiService, public.ContextPromptGenerator)
	if service == nil {
		return errors.New("推送服务初始化失败")
	}
	service.Start()
	public.PushService = service
	fmt.Printf("push 初始化配置成功 \n")
	return nil
}
