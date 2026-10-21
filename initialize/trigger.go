package initialize

import (
	"errors"
	"fmt"
	"new-robot-self/global"
	"new-robot-self/internal/trigger"
	"new-robot-self/public"
)

func InitTrigger() error {
	manager := trigger.NewManager()
	if manager == nil {
		return errors.New("trigger manager is nil")
	}
	config, err := trigger.LoadFromFile(global.TriggerPath)
	if err != nil {
		return err
	}
	manager.LoadTriggers(
		config.TimePoints,
		config.RealTime,
		config.UserBehavior,
		config.RandomTopics,
	)
	public.Trigger = manager
	fmt.Printf("trigger 初始化配置成功 \n")
	return nil
}
