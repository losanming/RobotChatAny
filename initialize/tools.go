package initialize

import (
	"log"
	public_brain "new-robot-self/public/brain"
)

func InitTools() error {
	log.Println("开始初始化工具系统...")
	public_brain.InitTools()
	log.Println("工具系统初始化完成")
	return nil
}
