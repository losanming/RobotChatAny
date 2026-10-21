package queue

import (
	"time"
)

// RetryConfig 重试配置
type RetryConfig struct {
	MaxRetries  int           // 最大重试次数
	InitialWait time.Duration // 初始等待时间
	MaxWait     time.Duration // 最大等待时间
}

// DefaultRetryConfig 默认重试配置
var DefaultRetryConfig = RetryConfig{
	MaxRetries:  3,
	InitialWait: time.Second,
	MaxWait:     time.Second * 10,
}

// RetryableError 可重试的错误
type RetryableError struct {
	Err error
}

func (e *RetryableError) Error() string {
	return e.Err.Error()
}

// NewRetryableError 创建可重试错误
func NewRetryableError(err error) *RetryableError {
	return &RetryableError{Err: err}
}

// IsRetryable 判断错误是否可重试
func IsRetryable(err error) bool {
	_, ok := err.(*RetryableError)
	return ok
}
