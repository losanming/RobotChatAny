package config

// HuaweiCloudOBSConfig 华为云OBS配置
type HuaweiCloudOBSConfig struct {
	// Endpoint OBS服务地址
	Endpoint string `yaml:"endpoint"`
	// AccessKey 访问密钥ID
	AccessKey string `yaml:"access_key"`
	// SecretKey 访问密钥密码
	SecretKey string `yaml:"secret_key"`
	// BucketName 存储桶名称
	BucketName string `yaml:"bucket_name"`
	Path       string `yaml:"path"`
}
