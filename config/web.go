package config

// WebConfig Web服务配置
type WebConfig struct {
	// Port HTTP服务端口
	Port int `yaml:"port"`
	// ReadTimeout 读取超时时间（秒）
	ReadTimeout int `yaml:"read_timeout"`
	// WriteTimeout 写入超时时间（秒）
	WriteTimeout int `yaml:"write_timeout"`
	// MaxHeaderBytes HTTP头部最大字节数
	MaxHeaderBytes int `yaml:"max_header_bytes"`
}
