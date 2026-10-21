package health

import (
	"github.com/google/uuid"
)

// HealthConfig 健康管理类工具配置
// 注意：为了避免循环导入，这个函数在 register_all.go 中不直接使用
// 而是通过 tools.HealthConfig() 获取配置
func HealthConfig() map[string]interface{} {
	return map[string]interface{}{
		"category":       "health",
		"category_desc":  "健康管理类工具",
		"category_icon":  "❤️",
		"category_color": "#FF6B6B",
		"enable_learning": true,
	}
}

// GenerateUUID 生成唯一ID
func GenerateUUID() string {
	return uuid.New().String()
}