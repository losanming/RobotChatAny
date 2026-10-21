package initialize

import (
	"fmt"
	"new-robot-self/config"
	"new-robot-self/public"
)

func InitConfig() error {
	loadConfig, err := config.LoadConfig()
	if err != nil {
		return err
	}
	public.Config = loadConfig
	fmt.Printf("config 初始化配置成功 \n")
	return nil
}
