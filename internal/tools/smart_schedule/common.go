package smart_schedule

import (
	"github.com/google/uuid"
)

// SmartScheduleConfig 智能日程类工具配置
func SmartScheduleConfig() map[string]interface{} {
	return map[string]interface{}{
		"category":       "smart_schedule",
		"category_desc":  "智能日程类工具",
		"category_icon":  "📅",
		"category_color": "#DDA0DD",
		"enable_learning": true,
	}
}

// GenerateUUID 生成唯一ID
func GenerateUUID() string {
	return uuid.New().String()
}