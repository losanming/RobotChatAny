package main

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"io/ioutil"
	"new-robot-self/pkg/screenshot"
	"new-robot-self/pkg/screenshot/quality"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	file1, err := ioutil.ReadFile("pkg/screenshot/example/screenshot_1755162152.jpg")
	file2, err := ioutil.ReadFile("pkg/screenshot/example/pic2.jpg")
	file3, err := ioutil.ReadFile("pkg/screenshot/example/zpic3.jpg")
	file4, err := ioutil.ReadFile("pkg/screenshot/example/zpic4.jpg")
	if err != nil {
		panic(err)
	}
	analyzer := quality.NewAnalyzer()
	result1, err := analyzer.AnalyzeBytes(file1)
	result2, err := analyzer.AnalyzeBytes(file2)
	result3, err := analyzer.AnalyzeBytes(file3)
	result4, err := analyzer.AnalyzeBytes(file4)
	fmt.Println(result1)
	fmt.Println(result2)
	fmt.Println(result3)
	fmt.Println(result4)
}

func screenBegin() {
	// 初始化日志
	logger, _ := zap.NewDevelopment()
	defer logger.Sync()

	// 创建截图配置
	config := screenshot.DefaultConfig()
	config.TargetWidth = 1280
	config.TargetHeight = 720
	config.JpegQuality = 80

	// 创建截图服务
	service := screenshot.NewService(
		screenshot.WithLogger(logger),
		screenshot.WithInterval(time.Minute), // 每分钟截图一次
		screenshot.WithConfig(config),
	)

	// 注册回调函数
	service.RegisterCallback(func(result *screenshot.CaptureResult) {
		logger.Info("收到新的截图",
			zap.Int("原始宽度", result.OriginalWidth),
			zap.Int("原始高度", result.OriginalHeight),
			zap.Int("压缩宽度", result.CompressedWidth),
			zap.Int("压缩高度", result.CompressedHeight),
			zap.Int("压缩大小(KB)", result.CompressedSize/1024),
			zap.Int64("时间戳", result.Timestamp),
		)

		// 这里可以添加将图像发送给AI处理的代码
		// 例如：
		// aiClient.AnalyzeImage(result.ImageData)

		// 保存图像到文件（仅用于测试）
		filename := fmt.Sprintf("screenshot_%d.jpg", result.Timestamp)
		err := os.WriteFile(filename, result.ImageData, 0644)
		if err != nil {
			logger.Error("保存截图失败", zap.Error(err))
		} else {
			logger.Info("截图已保存", zap.String("文件名", filename))
		}
	})

	// 创建上下文，用于控制服务生命周期
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// 启动服务
	if err := service.Start(ctx); err != nil {
		logger.Fatal("启动截图服务失败", zap.Error(err))
	}

	logger.Info("截图服务已启动，按Ctrl+C退出")

	// 等待中断信号
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	<-sigCh

	// 停止服务
	service.Stop()
	logger.Info("程序已退出")
}
