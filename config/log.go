package config

// LogConfig 日志配置
type LogConfig struct {
	// Level 日志级别(debug, info, warn, error)
	Level string `yaml:"level"`
	// Filename 日志文件路径
	Filename string `yaml:"filename"`
	// MaxSize 单个日志文件最大尺寸，单位MB
	MaxSize int `yaml:"max_size"`
	// MaxAge 日志文件保留天数
	MaxAge int `yaml:"max_age"`
	// MaxBackups 保留的旧日志文件个数
	MaxBackups int `yaml:"max_backups"`
}
