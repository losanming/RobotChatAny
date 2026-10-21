package tools

import (
	"net/http"
	"new-robot-self/global"
	"time"
)

var httpClient = &http.Client{Timeout: 6 * time.Second}

// GetAMapKey 从全局配置加载高德地图API密钥
func GetAMapKey() (string, error) {
	return global.AmapKey, nil
}
