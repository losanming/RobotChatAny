package logger

import (
	"fmt"
	"io"
	"log/slog"
	"new-robot-self/config"
	"os"
	"path/filepath"
	"time"

	"gopkg.in/natefinch/lumberjack.v2"
)

var (
	// 全局日志记录器
	defaultLogger *slog.Logger
)

// 初始化日志记录器
func InitLogger() {
	cfg := config.GetConfig()
	logCfg := cfg.Log

	// 确保日志目录存在
	logDir := filepath.Dir(logCfg.Filename)
	if err := os.MkdirAll(logDir, 0755); err != nil {
		fmt.Printf("创建日志目录失败: %v\n", err)
		os.Exit(1)
	}

	// 配置日志轮转
	rotator := &lumberjack.Logger{
		Filename:   logCfg.Filename,
		MaxSize:    logCfg.MaxSize,    // 单位：MB
		MaxAge:     logCfg.MaxAge,     // 单位：天
		MaxBackups: logCfg.MaxBackups, // 保留文件个数
		Compress:   true,              // 压缩旧日志文件
		LocalTime:  true,              // 使用本地时间
	}

	// 创建多输出写入器（同时写入文件和控制台）
	multiWriter := io.MultiWriter(rotator, os.Stdout)

	// 设置日志级别
	var level slog.Level
	switch logCfg.Level {
	case "debug":
		level = slog.LevelDebug
	case "info":
		level = slog.LevelInfo
	case "warn":
		level = slog.LevelWarn
	case "error":
		level = slog.LevelError
	default:
		level = slog.LevelInfo
	}

	// 创建JSON处理器
	opts := &slog.HandlerOptions{
		Level:     level,
		AddSource: true,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			// 自定义时间格式
			if a.Key == slog.TimeKey {
				if t, ok := a.Value.Any().(time.Time); ok {
					a.Value = slog.StringValue(t.Format("2006-01-02 15:04:05"))
				}
			}
			return a
		},
	}

	// 创建JSON格式的处理器
	handler := slog.NewJSONHandler(multiWriter, opts)

	// 创建日志记录器
	defaultLogger = slog.New(handler)

	// 设置为默认日志记录器
	slog.SetDefault(defaultLogger)
}
