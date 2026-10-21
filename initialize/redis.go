package initialize

import (
	"fmt"
	"log/slog"
	"new-robot-self/config"
	"new-robot-self/pkg/cache"
	"new-robot-self/public"
)

func InitCache() error {
	conf := config.GetConfig()
	var c cache.RedisConfig
	c.Host = conf.Redis.Host
	c.Port = conf.Redis.Port
	c.Password = conf.Redis.Password
	c.DB = conf.Redis.DB
	c.PoolSize = conf.Redis.PoolSize
	c.MinIdleConn = conf.Redis.MinIdleConns

	var clusters []cache.RedisClusterConfig
	for _, v := range conf.Redis.Clusters {
		clusters = append(clusters, cache.RedisClusterConfig{
			Host:        v.Host,
			Port:        v.Port,
			Password:    v.Password,
			DB:          v.DB,
			PoolSize:    v.PoolSize,
			MinIdleConn: v.MinIdleConns,
			Name:        v.Name,
		})
	}
	c.Clusters = clusters
	c.Fallback = cache.Fall{
		Enabled:        c.Fallback.Enabled,
		LocalCacheSize: c.Fallback.LocalCacheSize,
		RetryInterval:  c.Fallback.RetryInterval,
	}

	manager, err := cache.NewCacheManager(c)
	if err != nil {
		slog.Error("failed to create cache manager", "error", err)
		return err
	}
	public.Cache = manager
	fmt.Printf("cache 初始化配置成功 \n")
	return nil
}
