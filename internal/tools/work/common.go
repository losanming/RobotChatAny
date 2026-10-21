package work

import (
	"github.com/google/uuid"
)

// WorkConfig 工作助手类工具配置
func WorkConfig() map[string]interface{} {
	return map[string]interface{}{
		"category":       "work",
		"category_desc":  "工作助手类工具",
		"category_icon":  "💼",
		"category_color": "#4ECDC4",
		"enable_learning": true,
	}
}

// GenerateUUID 生成唯一ID
func GenerateUUID() string {
	return uuid.New().String()
}