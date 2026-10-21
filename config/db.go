package config

// PostgresConfig PostgreSQL数据库配置
type PostgresConfig struct {
	// Host 数据库主机地址
	Host string `yaml:"host"`
	// Port 数据库端口
	Port int `yaml:"port"`
	// User 数据库用户名
	User string `yaml:"user"`
	// Password 数据库密码
	Password string `yaml:"password"`
	// DBName 数据库名称
	DBName string `yaml:"dbname"`
	// MaxIdleConns 最大空闲连接数
	MaxIdleConns int `yaml:"max_idle_conns"`
	// MaxOpenConns 最大打开连接数
	MaxOpenConns int `yaml:"max_open_conns"`
	// ConnMaxLifetime 连接最大生命周期（秒）
	ConnMaxLifetime int `yaml:"conn_max_lifetime"`
}
