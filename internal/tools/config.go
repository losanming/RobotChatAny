package tools

import "time"

// ToolConfig 工具配置
// 用于统一管理所有工具的配置参数
type ToolConfig struct {
	// Name 工具名称
	Name string

	// Description 工具描述
	Description string

	// Category 工具分类
	// 可选值: info_fetcher, location, health, learning, lifecycle
	Category string

	// EnableCache 是否启用缓存
	EnableCache bool

	// CacheTTL 缓存过期时间
	CacheTTL time.Duration

	// EnableLearn 是否启用学习集成
	EnableLearn bool

	// MaxResults 最大结果数量
	// 用于限制搜索结果、热点新闻等返回的数量
	MaxResults int

	// QualityThreshold 质量阈值
	// 用于过滤低质量结果，范围 0.0-1.0
	QualityThreshold float64

	// Timeout 工具执行超时时间
	Timeout time.Duration

	// RetryCount 重试次数
	RetryCount int

	// RateLimit 速率限制 (每秒请求数)
	RateLimit float64
}

// DefaultToolConfig 返回默认的工具配置
func DefaultToolConfig() ToolConfig {
	return ToolConfig{
		EnableCache:       true,
		CacheTTL:          5 * time.Minute,
		EnableLearn:       true,
		MaxResults:        20,
		QualityThreshold:  0.6,
		Timeout:           30 * time.Second,
		RetryCount:        2,
		RateLimit:         10.0,
	}
}

// InfoFetcherConfig 信息获取类工具配置
func InfoFetcherConfig() ToolConfig {
	config := DefaultToolConfig()
	config.Category = "info_fetcher"
	config.MaxResults = 50
	config.CacheTTL = 10 * time.Minute
	return config
}

// LocationConfig 位置服务类工具配置
func LocationConfig() ToolConfig {
	config := DefaultToolConfig()
	config.Category = "location"
	config.MaxResults = 20
	config.CacheTTL = 30 * time.Minute
	return config
}

// HealthConfig 健康管理类工具配置
func HealthConfig() ToolConfig {
	config := DefaultToolConfig()
	config.Category = "health"
	config.MaxResults = 10
	config.CacheTTL = 1 * time.Hour
	return config
}

// LearningConfig 学习进化类工具配置
func LearningConfig() ToolConfig {
	config := DefaultToolConfig()
	config.Category = "learning"
	config.MaxResults = 100
	config.CacheTTL = 24 * time.Hour
	return config
}

// LifecycleConfig 生活服务类工具配置
func LifecycleConfig() ToolConfig {
	config := DefaultToolConfig()
	config.Category = "lifecycle"
	config.MaxResults = 15
	config.CacheTTL = 1 * time.Hour
	return config
}

// PersonalConfig 个人管理类工具配置
func PersonalConfig() ToolConfig {
	config := DefaultToolConfig()
	config.Category = "personal"
	config.MaxResults = 50
	config.CacheTTL = 24 * time.Hour
	config.EnableLearn = true
	return config
}

// AnalysisConfig 信息分析类工具配置
func AnalysisConfig() ToolConfig {
	config := DefaultToolConfig()
	config.Category = "analysis"
	config.MaxResults = 100
	config.CacheTTL = 30 * time.Minute
	config.EnableLearn = true
	return config
}

// UtilityConfig 通用能力类工具配置
func UtilityConfig() ToolConfig {
	config := DefaultToolConfig()
	config.Category = "utility"
	config.MaxResults = 20
	config.CacheTTL = 5 * time.Minute
	config.EnableLearn = false
	return config
}

// CommunicationConfig 通信类工具配置
func CommunicationConfig() ToolConfig {
	config := DefaultToolConfig()
	config.Category = "communication"
	config.MaxResults = 50
	config.CacheTTL = 10 * time.Minute
	config.EnableLearn = true
	return config
}

// MultimediaConfig 多媒体类工具配置
func MultimediaConfig() ToolConfig {
	config := DefaultToolConfig()
	config.Category = "multimedia"
	config.MaxResults = 30
	config.CacheTTL = 15 * time.Minute
	config.EnableLearn = false
	return config
}

// SecurityConfig 安全类工具配置
func SecurityConfig() ToolConfig {
	config := DefaultToolConfig()
	config.Category = "security"
	config.MaxResults = 20
	config.CacheTTL = 1 * time.Hour
	config.EnableLearn = false
	return config
}

// FinanceConfig 金融管理类工具配置
func FinanceConfig() ToolConfig {
	config := DefaultToolConfig()
	config.Category = "finance"
	config.MaxResults = 50
	config.CacheTTL = 24 * time.Hour
	config.EnableLearn = true
	return config
}

// AutomationConfig 自动化类工具配置
func AutomationConfig() ToolConfig {
	config := DefaultToolConfig()
	config.Category = "automation"
	config.MaxResults = 30
	config.CacheTTL = 1 * time.Hour
	config.EnableLearn = true
	return config
}

// CloudConfig 云服务类工具配置
func CloudConfig() ToolConfig {
	config := DefaultToolConfig()
	config.Category = "cloud"
	config.MaxResults = 20
	config.CacheTTL = 30 * time.Minute
	config.EnableLearn = false
	return config
}

// SocialConfig 社交媒体类工具配置
func SocialConfig() ToolConfig {
	config := DefaultToolConfig()
	config.Category = "social"
	config.MaxResults = 40
	config.CacheTTL = 1 * time.Hour
	config.EnableLearn = true
	return config
}

// IntelligenceConfig 智能化类工具配置
func IntelligenceConfig() ToolConfig {
	config := DefaultToolConfig()
	config.Category = "intelligence"
	config.MaxResults = 100
	config.CacheTTL = 24 * time.Hour
	config.EnableLearn = true
	return config
}

// MobileConfig 移动端类工具配置
func MobileConfig() ToolConfig {
	config := DefaultToolConfig()
	config.Category = "mobile"
	config.MaxResults = 20
	config.CacheTTL = 10 * time.Minute
	config.EnableLearn = false
	return config
}

// WorkConfig 工作助手类工具配置
func WorkConfig() ToolConfig {
	config := DefaultToolConfig()
	config.Category = "work"
	config.MaxResults = 50
	config.CacheTTL = 24 * time.Hour
	config.EnableLearn = true
	return config
}

// VisualizationConfig 数据可视化类工具配置
func VisualizationConfig() ToolConfig {
	config := DefaultToolConfig()
	config.Category = "visualization"
	config.MaxResults = 20
	config.CacheTTL = 30 * time.Minute
	config.EnableLearn = false
	return config
}

// SmartScheduleConfig 智能日程类工具配置
func SmartScheduleConfig() ToolConfig {
	config := DefaultToolConfig()
	config.Category = "smart_schedule"
	config.MaxResults = 30
	config.CacheTTL = 1 * time.Hour
	config.EnableLearn = true
	return config
}

// VoiceConfig 语音类工具配置
func VoiceConfig() ToolConfig {
	config := DefaultToolConfig()
	config.Category = "voice"
	config.MaxResults = 10
	config.CacheTTL = 5 * time.Minute
	config.EnableLearn = false
	return config
}
