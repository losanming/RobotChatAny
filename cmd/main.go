package main

import (
	"fmt"
	"new-robot-self/initialize"
	"new-robot-self/pkg/banner"
	"new-robot-self/pkg/logger"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	banner.PrintCustomBanner("原神启动", "v0.0.1")
	var f []func() error
	var e []error
	f = append(f,
		initialize.InitConfig,
		initialize.InitCache,
		initialize.InitPgsql,
		initialize.InitCtx,
		initialize.InitAdcode, // 高德城市编码
		initialize.InitTools, // 统一工具系统
		//以下需要在数据库、缓存等硬性资源完成后的初始化工具
		initialize.InitAiService,
		initialize.InitFetcher,         // 资源加载器
		initialize.InitContext,         // 用户上下文加载器
		initialize.InitPromptGenerator, // prompt加载器
		initialize.InitTrigger,         // 触发器加载器

		initialize.InitQueueAndDispatcher,     // 消息和分发器
		initialize.InitConnectionAndHeartbeat, // 连接和心跳
		initialize.InitTcpServer,              // tcp服务
		initialize.InitPushService,            // 推送服务

		initialize.InitAnalyzer, // 截图质量分析管理器
		initialize.InitNSQ,      // NSQ 消费者（云端只写消费者）
		initialize.InitRAG,      // RAG服务初始化

	)
	for _, v := range f {
		err := v()
		if err != nil {
			e = append(e, err)
		}
	}
	if len(e) > 0 {
		for _, v := range e {
			fmt.Printf("初始化错误err: %s\n", v)
		}
		fmt.Printf("初始化错误,程序退出")
		os.Exit(22)
	}
	logger.InitLogger()
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	fmt.Println("程序执行完毕!")
}
