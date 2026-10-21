package amap

import (
	"errors"
	"net/http"
	"time"
)

// Client 高德 Web 服务最小客户端
// 作用：
// - 维护调用所需的 Key、BaseURL 与 http.Client（带超时）。
// - 由更高层的 Provider 方法（如 ReverseGeocode）复用。
// 字段：
// - Key：高德 Web 服务 API Key（敏感信息，建议从配置/环境变量注入）
// - BaseURL：服务基础地址，默认 https://restapi.amap.com
// - HTTP：底层 HTTP 客户端；仅设置超时，其他治理（重试/限流/Tracing）建议在上层实现
type Client struct {
	Key     string
	BaseURL string
	HTTP    *http.Client
}

// NewClient 创建高德客户端
// 参数：
// - key：高德 Web 服务 Key，必填
// - timeoutMs：HTTP 超时时间（毫秒）
// 返回：
// - *Client：客户端实例
// - error：当 key 为空时返回错误
// 说明：
// - 仅做最小封装，复杂治理（重试/限流/签名/日志）建议在调用层处理或在 Client 外层包装。
func NewClient(key string, timeoutMs int) (*Client, error) {
	if key == "" {
		return nil, errors.New("amap key is empty")
	}
	cli := &http.Client{Timeout: time.Duration(timeoutMs) * time.Millisecond}
	return &Client{
		Key:     key,
		BaseURL: "https://restapi.amap.com",
		HTTP:    cli,
	}, nil
}
