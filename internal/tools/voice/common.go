package voice

import (
	"github.com/google/uuid"
)

// VoiceConfig 语音类工具配置
func VoiceConfig() map[string]interface{} {
	return map[string]interface{}{
		"category":       "voice",
		"category_desc":  "语音类工具",
		"category_icon":  "🎤",
		"category_color": "#FFB6C1",
		"enable_learning": true,
	}
}

// GenerateUUID 生成唯一ID
func GenerateUUID() string {
	return uuid.New().String()
}