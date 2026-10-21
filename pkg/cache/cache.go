package cache

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/allegro/bigcache/v3"
	"github.com/redis/go-redis/v9"
	"github.com/spaolacci/murmur3"
)

// 错误定义
var (
	ErrKeyNotFound      = errors.New("key not found in cache")
	ErrRedisUnavailable = errors.New("redis is unavailable")
	ErrInvalidValue     = errors.New("invalid cache value")
)

// Cache 定义缓存接口
type Cache interface {
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error
	Get(ctx context.Context, key string, value interface{}) error
	Delete(ctx context.Context, key string) error
	Exists(ctx context.Context, key string) (bool, error)
	Clear(ctx context.Context) error
}

// RedisClusterConfig Redis集群配置
type RedisClusterConfig struct {
	Name        string `json:"name"`
	Host        string `json:"host"`
	Port        int    `json:"port"`
	Password    string `json:"password"`
	DB          int    `json:"db"`
	PoolSize    int    `json:"pool_size"`
	MinIdleConn int    `json:"min_idle_conns"`
}

// RedisConfig Redis配置
type RedisConfig struct {
	// 单Redis配置（向后兼容）
	Host        string `json:"host"`
	Port        int    `json:"port"`
	Password    string `json:"password"`
	DB          int    `json:"db"`
	PoolSize    int    `json:"pool_size"`
	MinIdleConn int    `json:"min_idle_conns"`

	// 多Redis集群配置
	Clusters []RedisClusterConfig `json:"clusters"`

	// 缓存降级配置
	Fallback Fall `json:"fallback"`
}

type Fall struct {
	Enabled        bool `json:"enabled"`
	LocalCacheSize int  `json:"local_cache_size"` // MB
	RetryInterval  int  `json:"retry_interval"`   // 秒
}

// RedisCache Redis缓存实现
type RedisCache struct {
	client *redis.Client
}

// NewRedisCache 创建Redis缓存
func NewRedisCache(config RedisClusterConfig) *RedisCache {
	client := redis.NewClient(&redis.Options{
		Addr:         fmt.Sprintf("%s:%d", config.Host, config.Port),
		Password:     config.Password,
		DB:           config.DB,
		PoolSize:     config.PoolSize,
		MinIdleConns: config.MinIdleConn,
	})

	return &RedisCache{client: client}
}

// Set 设置缓存
func (rc *RedisCache) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return rc.client.Set(ctx, key, data, expiration).Err()
}

// Get 获取缓存
func (rc *RedisCache) Get(ctx context.Context, key string, value interface{}) error {
	data, err := rc.client.Get(ctx, key).Bytes()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return ErrKeyNotFound
		}
		return err
	}

	return json.Unmarshal(data, value)
}

// Delete 删除缓存
func (rc *RedisCache) Delete(ctx context.Context, key string) error {
	return rc.client.Del(ctx, key).Err()
}

// Exists 检查缓存是否存在
func (rc *RedisCache) Exists(ctx context.Context, key string) (bool, error) {
	val, err := rc.client.Exists(ctx, key).Result()
	if err != nil {
		return false, err
	}
	return val > 0, nil
}

// Clear 清空缓存（危险操作，仅用于测试）
func (rc *RedisCache) Clear(ctx context.Context) error {
	return rc.client.FlushDB(ctx).Err()
}

// Ping 检查Redis连接
func (rc *RedisCache) Ping(ctx context.Context) error {
	return rc.client.Ping(ctx).Err()
}

// LocalCache 基于BigCache的本地缓存实现
type LocalCache struct {
	cache       *bigcache.BigCache
	expirations sync.Map // 存储过期时间
}

// NewLocalCache 创建本地缓存
func NewLocalCache(sizeInMB int) (*LocalCache, error) {
	config := bigcache.DefaultConfig(24 * time.Hour)
	config.HardMaxCacheSize = sizeInMB

	cache, err := bigcache.New(context.Background(), config)
	if err != nil {
		return nil, err
	}

	return &LocalCache{
		cache:       cache,
		expirations: sync.Map{},
	}, nil
}

// Set 设置本地缓存
func (lc *LocalCache) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}

	// 设置过期时间
	if expiration > 0 {
		lc.expirations.Store(key, time.Now().Add(expiration).Unix())
	}

	return lc.cache.Set(key, data)
}

// Get 获取本地缓存
func (lc *LocalCache) Get(ctx context.Context, key string, value interface{}) error {
	// 检查是否过期
	if expired, _ := lc.isExpired(key); expired {
		lc.cache.Delete(key)
		lc.expirations.Delete(key)
		return ErrKeyNotFound
	}

	data, err := lc.cache.Get(key)
	if err != nil {
		if err == bigcache.ErrEntryNotFound {
			return ErrKeyNotFound
		}
		return err
	}

	return json.Unmarshal(data, value)
}

// isExpired 检查键是否过期
func (lc *LocalCache) isExpired(key string) (bool, error) {
	expVal, ok := lc.expirations.Load(key)
	if !ok {
		return false, nil // 没有设置过期时间
	}

	expTime, ok := expVal.(int64)
	if !ok {
		return false, ErrInvalidValue
	}

	return time.Now().Unix() > expTime, nil
}

// Delete 删除本地缓存
func (lc *LocalCache) Delete(ctx context.Context, key string) error {
	lc.expirations.Delete(key)
	return lc.cache.Delete(key)
}

// Exists 检查本地缓存是否存在
func (lc *LocalCache) Exists(ctx context.Context, key string) (bool, error) {
	// 检查是否过期
	if expired, _ := lc.isExpired(key); expired {
		lc.cache.Delete(key)
		lc.expirations.Delete(key)
		return false, nil
	}

	_, err := lc.cache.Get(key)
	if err != nil {
		if errors.Is(err, bigcache.ErrEntryNotFound) {
			return false, nil
		}
		return false, err
	}

	return true, nil
}

// Clear 清空本地缓存
func (lc *LocalCache) Clear(ctx context.Context) error {
	return lc.cache.Reset()
}

// CacheManager 缓存管理器
type CacheManager struct {
	redisClients  []*RedisCache // Redis客户端列表
	localCache    *LocalCache   // 本地缓存
	useLocalOnly  bool          // 是否仅使用本地缓存
	retryInterval time.Duration // Redis故障重试间隔
	mutex         sync.RWMutex  // 读写锁
}

// NewCacheManager 创建缓存管理器
func NewCacheManager(config RedisConfig) (*CacheManager, error) {
	manager := &CacheManager{
		retryInterval: time.Duration(config.Fallback.RetryInterval) * time.Second,
	}

	// 创建本地缓存
	if config.Fallback.Enabled {
		localCache, err := NewLocalCache(config.Fallback.LocalCacheSize)
		if err != nil {
			return nil, fmt.Errorf("failed to create local cache: %w", err)
		}
		manager.localCache = localCache
	}

	// 尝试连接Redis
	if len(config.Clusters) > 0 {
		// 使用多Redis集群配置
		for _, clusterConfig := range config.Clusters {
			redisCache := NewRedisCache(clusterConfig)
			manager.redisClients = append(manager.redisClients, redisCache)
		}
	} else if config.Host != "" {
		// 使用单Redis配置
		clusterConfig := RedisClusterConfig{
			Name:        "default",
			Host:        config.Host,
			Port:        config.Port,
			Password:    config.Password,
			DB:          config.DB,
			PoolSize:    config.PoolSize,
			MinIdleConn: config.MinIdleConn,
		}
		redisCache := NewRedisCache(clusterConfig)
		manager.redisClients = append(manager.redisClients, redisCache)
	}

	// 检查Redis连接状态
	if len(manager.redisClients) > 0 {
		ctx := context.Background()
		for _, client := range manager.redisClients {
			err := client.Ping(ctx)
			if err != nil {
				log.Printf("Redis连接失败: %v, 将使用本地缓存", err)
				manager.useLocalOnly = true
				// 启动Redis重连协程
				if config.Fallback.Enabled {
					go manager.reconnectRedis()
				}
				break
			}
		}
	} else if manager.localCache != nil {
		// 没有配置Redis，使用本地缓存
		manager.useLocalOnly = true
	} else {
		return nil, errors.New("no cache available: neither Redis nor local cache is configured")
	}

	return manager, nil
}

// reconnectRedis 尝试重新连接Redis
func (cm *CacheManager) reconnectRedis() {
	for {
		time.Sleep(cm.retryInterval)

		ctx := context.Background()
		allConnected := true

		for _, client := range cm.redisClients {
			err := client.Ping(ctx)
			if err != nil {
				allConnected = false
				break
			}
		}

		if allConnected && len(cm.redisClients) > 0 {
			cm.mutex.Lock()
			cm.useLocalOnly = false
			cm.mutex.Unlock()
			log.Println("Redis连接恢复，切换回Redis缓存")
			return
		}
	}
}

// getRedisClient 根据key获取Redis客户端
func (cm *CacheManager) getRedisClient(key string) *RedisCache {
	if len(cm.redisClients) == 1 {
		return cm.redisClients[0]
	}

	// 使用一致性哈希选择Redis客户端
	hashValue := murmur3.Sum32([]byte(key))
	index := int(hashValue) % len(cm.redisClients)
	return cm.redisClients[index]
}

// getRedisClient 根据key获取Redis客户端
func (cm *CacheManager) GetRedisClient() *redis.Client {
	if len(cm.redisClients) == 0 {
		return nil
	}
	return cm.redisClients[0].client
}

// Set 设置缓存
func (cm *CacheManager) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	// 如果有本地缓存，先设置本地缓存
	if cm.localCache != nil {
		err := cm.localCache.Set(ctx, key, value, expiration)
		if err != nil {
			log.Printf("设置本地缓存失败: %v", err)
		}
	}

	// 检查是否只使用本地缓存
	cm.mutex.RLock()
	useLocalOnly := cm.useLocalOnly
	cm.mutex.RUnlock()

	if useLocalOnly || len(cm.redisClients) == 0 {
		return nil // 已经设置了本地缓存
	}

	// 设置Redis缓存
	redisClient := cm.getRedisClient(key)
	return redisClient.Set(ctx, key, value, expiration)
}

// Get 获取缓存
func (cm *CacheManager) Get(ctx context.Context, key string, value interface{}) error {
	// 先尝试从本地缓存获取
	if cm.localCache != nil {
		err := cm.localCache.Get(ctx, key, value)
		if err == nil {
			return nil // 本地缓存命中
		}

		if err != ErrKeyNotFound {
			log.Printf("获取本地缓存失败: %v", err)
		}
	}

	// 检查是否只使用本地缓存
	cm.mutex.RLock()
	useLocalOnly := cm.useLocalOnly
	cm.mutex.RUnlock()

	if useLocalOnly || len(cm.redisClients) == 0 {
		return ErrKeyNotFound
	}

	// 从Redis获取
	redisClient := cm.getRedisClient(key)
	err := redisClient.Get(ctx, key, value)
	if err != nil {
		return err
	}

	// 如果Redis命中，更新本地缓存
	if cm.localCache != nil {
		_ = cm.localCache.Set(ctx, key, value, 0) // 忽略错误
	}

	return nil
}

// Delete 删除缓存
func (cm *CacheManager) Delete(ctx context.Context, key string) error {
	// 删除本地缓存
	if cm.localCache != nil {
		_ = cm.localCache.Delete(ctx, key) // 忽略错误
	}

	// 检查是否只使用本地缓存
	cm.mutex.RLock()
	useLocalOnly := cm.useLocalOnly
	cm.mutex.RUnlock()

	if useLocalOnly || len(cm.redisClients) == 0 {
		return nil
	}

	// 删除Redis缓存
	redisClient := cm.getRedisClient(key)
	return redisClient.Delete(ctx, key)
}

// Exists 检查缓存是否存在
func (cm *CacheManager) Exists(ctx context.Context, key string) (bool, error) {
	// 先检查本地缓存
	if cm.localCache != nil {
		exists, err := cm.localCache.Exists(ctx, key)
		if err == nil && exists {
			return true, nil // 本地缓存存在
		}
	}

	// 检查是否只使用本地缓存
	cm.mutex.RLock()
	useLocalOnly := cm.useLocalOnly
	cm.mutex.RUnlock()

	if useLocalOnly || len(cm.redisClients) == 0 {
		return false, nil
	}

	// 检查Redis缓存
	redisClient := cm.getRedisClient(key)
	return redisClient.Exists(ctx, key)
}

// Clear 清空缓存
func (cm *CacheManager) Clear(ctx context.Context) error {
	// 清空本地缓存
	if cm.localCache != nil {
		_ = cm.localCache.Clear(ctx) // 忽略错误
	}

	// 检查是否只使用本地缓存
	cm.mutex.RLock()
	useLocalOnly := cm.useLocalOnly
	cm.mutex.RUnlock()

	if useLocalOnly || len(cm.redisClients) == 0 {
		return nil
	}

	// 清空所有Redis缓存
	for _, client := range cm.redisClients {
		err := client.Clear(ctx)
		if err != nil {
			return err
		}
	}

	return nil
}
