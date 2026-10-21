package initialize

import (
	"fmt"
	"new-robot-self/global"
	"new-robot-self/internal/ai"
	public_ai "new-robot-self/public/ai"
)

func InitAiService() error {
	service := ai.NewService(global.DefaultConfig)
	err := service.InitModel()
	if err != nil {
		return err
	}
	public_ai.AiService = service
	fmt.Printf("ai 初始化配置成功 \n")
	return nil
}
