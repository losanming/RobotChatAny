package visualization

import (
	"github.com/google/uuid"
)

// VisualizationConfig 数据可视化类工具配置
func VisualizationConfig() map[string]interface{} {
	return map[string]interface{}{
		"category":       "visualization",
		"category_desc":  "数据可视化类工具",
		"category_icon":  "📊",
		"category_color": "#95E1D3",
		"enable_learning": true,
	}
}

// GenerateUUID 生成唯一ID
func GenerateUUID() string {
	return uuid.New().String()
}