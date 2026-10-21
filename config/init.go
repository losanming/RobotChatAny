package config

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
	"sync"
)

// Config 应用程序总配置结构体
type Config struct {
	// Postgres 数据库配置
	Postgres PostgresConfig `yaml:"postgres"`
	// Redis 缓存配置
	Redis RedisConfig `yaml:"redis"`
	// HuaweiCloudOBS 华为云对象存储配置
	HuaweiCloudOBS HuaweiCloudOBSConfig `yaml:"huaweicloud_obs"`
	MinIo          Minio                `yaml:"minio"`
	OssLocal       Local                `yaml:"local"`
	// Web HTTP服务配置
	Web WebConfig `yaml:"web"`
	// TCP 长连接服务配置
	TCP TCPConfig `yaml:"tcp"`
	// Log 日志配置
	Log   LogConfig   `yaml:"log"`
	Sys   System      `yaml:"system"`
	Gitee GiteeConfig `yaml:"gitee"`
	Nsq   NsqConfig   `yaml:"nsq"`
	// RAG RAG相关配置
	RAG   RAGConfig   `yaml:"rag"`
}

var (
	SelfConfig *Config
	once       sync.Once
)

func LoadConfig() (*Config, error) {
	data, err := os.ReadFile("config.yaml")
	if err != nil {
		return nil, err
	}

	config := &Config{}
	if err := yaml.Unmarshal(data, config); err != nil {
		return nil, err
	}
	SelfConfig = config
	return config, nil
}

func GetConfig() *Config {
	if SelfConfig == nil {
		if _, err := LoadConfig(); err != nil {
			panic(fmt.Sprintf("初始化配置失败: %v", err))
		}
	}
	return SelfConfig
}
