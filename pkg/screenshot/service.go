package screenshot

import (
	"context"
	"fmt"
	"sync"
	"time"

	"go.uber.org/zap"
)

// CaptureCallback 截图回调函数类型
type CaptureCallback func(*CaptureResult)

// Service 截图服务
type Service struct {
	config         *ScreenshotConfig
	interval       time.Duration
	callbacks      []CaptureCallback
	logger         *zap.Logger
	stopChan       chan struct{}
	wg             sync.WaitGroup
	latestCapture  *CaptureResult
	latestCaptureMu sync.RWMutex
	isRunning      bool
	runningMu      sync.Mutex
}

// ServiceOption 服务配置选项
type ServiceOption func(*Service)

// WithLogger 设置日志记录器
func WithLogger(logger *zap.Logger) ServiceOption {
	return func(s *Service) {
		s.logger = logger
	}
}

// WithInterval 设置截图间隔
func WithInterval(interval time.Duration) ServiceOption {
	return func(s *Service) {
		s.interval = interval
	}
}

// WithConfig 设置截图配置
func WithConfig(config *ScreenshotConfig) ServiceOption {
	return func(s *Service) {
		s.config = config
	}
}

// NewService 创建新的截图服务
func NewService(options ...ServiceOption) *Service {
	s := &Service{
		config:    DefaultConfig(),
		interval:  time.Minute, // 默认1分钟
		callbacks: make([]CaptureCallback, 0),
		stopChan:  make(chan struct{}),
		logger:    zap.NewNop(), // 默认不记录日志
	}

	// 应用选项
	for _, option := range options {
		option(s)
	}

	return s
}

// RegisterCallback 注册截图回调
func (s *Service) RegisterCallback(callback CaptureCallback) {
	s.callbacks = append(s.callbacks, callback)
}

// Start 启动截图服务
func (s *Service) Start(ctx context.Context) error {
	s.runningMu.Lock()
	defer s.runningMu.Unlock()

	if s.isRunning {
		return fmt.Errorf("截图服务已经在运行")
	}

	s.isRunning = true
	s.wg.Add(1)

	go func() {
		defer s.wg.Done()
		ticker := time.NewTicker(s.interval)
		defer ticker.Stop()

		// 立即执行一次截图
		s.captureAndNotify()

		for {
			select {
			case <-ticker.C:
				s.captureAndNotify()
			case <-s.stopChan:
				return
			case <-ctx.Done():
				return
			}
		}
	}()

	s.logger.Info("截图服务已启动", zap.Duration("interval", s.interval))
	return nil
}

// Stop 停止截图服务
func (s *Service) Stop() {
	s.runningMu.Lock()
	defer s.runningMu.Unlock()

	if !s.isRunning {
		return
	}

	close(s.stopChan)
	s.wg.Wait()
	s.isRunning = false
	s.stopChan = make(chan struct{})
	s.logger.Info("截图服务已停止")
}

// GetLatestCapture 获取最新的截图结果
func (s *Service) GetLatestCapture() *CaptureResult {
	s.latestCaptureMu.RLock()
	defer s.latestCaptureMu.RUnlock()
	return s.latestCapture
}

// captureAndNotify 执行截图并通知回调
func (s *Service) captureAndNotify() {
	result, err := Capture(s.config)
	if err != nil {
		s.logger.Error("截图失败", zap.Error(err))
		return
	}

	// 保存最新截图
	s.latestCaptureMu.Lock()
	s.latestCapture = result
	s.latestCaptureMu.Unlock()

	// 记录日志
	s.logger.Debug("截图成功",
		zap.Int("original_width", result.OriginalWidth),
		zap.Int("original_height", result.OriginalHeight),
		zap.Int("compressed_width", result.CompressedWidth),
		zap.Int("compressed_height", result.CompressedHeight),
		zap.Int("compressed_size_kb", result.CompressedSize/1024),
	)

	// 通知所有回调
	for _, callback := range s.callbacks {
		go func(cb CaptureCallback) {
			cb(result)
		}(callback)
	}
}
