package config

// RedisClusterConfig Redis集群配置
type RedisClusterConfig struct {
	// Name 集群名称
	Name string `yaml:"name"`
	// Host Redis主机地址
	Host string `yaml:"host"`
	// Port Redis端口
	Port int `yaml:"port"`
	// Password Redis密码
	Password string `yaml:"password"`
	// DB Redis数据库索引
	DB int `yaml:"db"`
	// PoolSize 连接池大小
	PoolSize int `yaml:"pool_size"`
	// MinIdleConns 最小空闲连接数
	MinIdleConns int `yaml:"min_idle_conns"`
}

// RedisFallbackConfig Redis降级配置
type RedisFallbackConfig struct {
	// Enabled 是否启用降级
	Enabled bool `yaml:"enabled"`
	// LocalCacheSize 本地缓存大小(MB)
	LocalCacheSize int `yaml:"local_cache_size"`
	// RetryInterval Redis故障后重试间隔(秒)
	RetryInterval int `yaml:"retry_interval"`
}

// RedisConfig Redis配置
type RedisConfig struct {
	// 单Redis配置（向后兼容）
	// Host Redis主机地址
	Host string `yaml:"host"`
	// Port Redis端口
	Port int `yaml:"port"`
	// Password Redis密码
	Password string `yaml:"password"`
	// DB Redis数据库索引
	DB int `yaml:"db"`
	// PoolSize 连接池大小
	PoolSize int `yaml:"pool_size"`
	// MinIdleConns 最小空闲连接数
	MinIdleConns int `yaml:"min_idle_conns"`

	// Clusters 多Redis集群配置
	Clusters []RedisClusterConfig `yaml:"clusters"`

	// Fallback 缓存降级配置
	Fallback RedisFallbackConfig `yaml:"fallback"`
}
