package initialize

import (
	"errors"
	"fmt"
	"new-robot-self/internal/context"
	"new-robot-self/public"
	public_contextManager "new-robot-self/public/contextManager"
)

func InitContext() error {
	// 这里的数据只存在redis 所以redis必须要通
	if public.Cache == nil {
		return errors.New("redis未初始化")
	}
	manager := context.NewContextManager(public.DB, public.Cache.GetRedisClient())
	public_contextManager.ContextManager = manager
	fmt.Printf("context 初始化配置成功 \n")
	return nil
}
