package web

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// WebConfig HTTP服务配置
type WebConfig struct {
	// Port HTTP服务端口
	Port int `yaml:"port"`
	// ReadTimeout 读取超时时间（秒）
	ReadTimeout int `yaml:"read_timeout"`
	// WriteTimeout 写入超时时间（秒）
	WriteTimeout int `yaml:"write_timeout"`
	// MaxHeaderBytes HTTP头部最大字节数
	MaxHeaderBytes int `yaml:"max_header_bytes"`
	// Mode gin模式 (debug, release, test)
	Mode string `yaml:"mode"`
}

// DefaultWebConfig 默认配置
var DefaultWebConfig = WebConfig{
	Port:           8080,
	ReadTimeout:    60,
	WriteTimeout:   60,
	MaxHeaderBytes: 1 << 20, // 1MB
	Mode:           "debug",
}

// Server Web服务器
type Server struct {
	Engine     *gin.Engine
	HttpServer *http.Server
	Config     WebConfig
}

// NewServer 创建新的Web服务器
func NewServer(config WebConfig) *Server {
	// 设置gin模式
	gin.SetMode(config.Mode)

	// 创建gin引擎
	engine := gin.Default()

	// 创建HTTP服务器
	addr := fmt.Sprintf(":%d", config.Port)
	httpServer := &http.Server{
		Addr:           addr,
		Handler:        engine,
		ReadTimeout:    time.Duration(config.ReadTimeout) * time.Second,
		WriteTimeout:   time.Duration(config.WriteTimeout) * time.Second,
		MaxHeaderBytes: config.MaxHeaderBytes,
	}

	return &Server{
		Engine:     engine,
		HttpServer: httpServer,
		Config:     config,
	}
}

// Run 启动Web服务器
func (s *Server) Run() error {
	return s.HttpServer.ListenAndServe()
}

// RunTLS 启动HTTPS服务器
func (s *Server) RunTLS(certFile, keyFile string) error {
	return s.HttpServer.ListenAndServeTLS(certFile, keyFile)
}

// Shutdown 优雅关闭服务器
func (s *Server) Shutdown(timeout time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	return s.HttpServer.Shutdown(ctx)
}

// RunWeb 快速启动Web服务 (向后兼容旧接口)
func RunWeb(host string, port int) (*gin.Engine, *http.Server, error) {
	if port <= 0 || port > 65535 {
		return nil, nil, errors.New("端口号无效")
	}

	addr := host
	if host == "" {
		addr = fmt.Sprintf(":%d", port)
	}

	engine := gin.Default()
	httpServer := &http.Server{
		Addr:    addr,
		Handler: engine,
	}

	// 非阻塞方式启动服务器
	go func() {
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Printf("HTTP服务器启动失败: %v\n", err)
		}
	}()

	return engine, httpServer, nil
}
