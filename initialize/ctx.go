package initialize

import (
	"context"
	"fmt"
	"new-robot-self/public"
)

func InitCtx() error {
	public.Ctx = context.Background()
	fmt.Printf("ctx 初始化配置成功 \n")
	return nil
}
func NewCtx() context.Context {
	return context.Background()
}
