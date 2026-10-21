package config

// TCPConfig TCP服务配置
type TCPConfig struct {
	// Port TCP服务端口
	Port int `yaml:"port"`
	// ReadTimeout 读取超时（秒）
	ReadTimeout int `yaml:"read_timeout"`
	// WriteTimeout 写入超时（秒）
	WriteTimeout int `yaml:"write_timeout"`
	// MaxConnections 最大连接数
	MaxConnections int `yaml:"max_connections"`
	// HeartbeatInterval 心跳间隔（秒）
	HeartbeatInterval int `yaml:"heartbeat_interval"`
	// HeartbeatTimeout 心跳超时（秒）
	HeartbeatTimeout int `yaml:"heartbeat_timeout"`
	// SecretKey 安全密钥,用于消息加密
	SecretKey string `yaml:"secret_key"`
	// TLS TLS配置
	TLS struct {
		// Enabled 是否启用TLS
		Enabled bool `yaml:"enabled"`
		// CertFile 证书文件路径
		CertFile string `yaml:"cert_file"`
		// KeyFile 密钥文件路径
		KeyFile string `yaml:"key_file"`
	} `yaml:"tls"`
}
