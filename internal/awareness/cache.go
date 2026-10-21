package awareness

import (
	"context"
	"sync"
	"time"

	loc "new-robot-self/pkg/location"
	amap "new-robot-self/pkg/location/provider/amap"
)

// 简单内存缓存条目

type cacheItem[T any] struct {
	val   T
	expTs int64 // 过期时间(ms)
}

type regeoKey struct {
	lat float64
	lon float64
}

type aroundKey struct {
	lat    float64
	lon    float64
	kw     string
	types  string
	radius int
}

type weatherKey struct {
	city string
	ext  string
}

// cacheStore 管理三类缓存

type cacheStore struct {
	mu      sync.Mutex
	regeo   map[regeoKey]cacheItem[loc.ReverseGeocodeResult]
	around  map[aroundKey]cacheItem[loc.PlaceAroundResult]
	weather map[weatherKey]cacheItem[loc.WeatherResult]
}

func newCacheStore() *cacheStore {
	return &cacheStore{
		regeo:   make(map[regeoKey]cacheItem[loc.ReverseGeocodeResult]),
		around:  make(map[aroundKey]cacheItem[loc.PlaceAroundResult]),
		weather: make(map[weatherKey]cacheItem[loc.WeatherResult]),
	}
}

// 带缓存与重试的封装：挂在 Service 上以便使用其配置

func (s *Service) regeoWithCache(ctx context.Context, lat, lon float64, p amap.RegeoParams) (loc.ReverseGeocodeResult, error) {
	key := regeoKey{lat: round6(lat), lon: round6(lon)}
	// 命中
	s.cache.mu.Lock()
	if it, ok := s.cache.regeo[key]; ok && it.expTs > time.Now().Unix() {
		v := it.val
		s.cache.mu.Unlock()
		return v, nil
	}
	s.cache.mu.Unlock()
	// 未命中：重试调用
	var lastErr error
	var out loc.ReverseGeocodeResult
	for i := 0; i < s.retryTimes; i++ {
		v, err := s.amapClient.ReverseGeocode(ctx, lat, lon, p)
		if err == nil {
			out = v
			break
		}
		lastErr = err
		t := time.Duration(s.retryBaseDelayMs*(1<<i)) * time.Millisecond
		time.Sleep(t)
	}
	if out.TS == 0 && lastErr != nil {
		return out, lastErr
	}
	// 写入缓存
	s.cache.mu.Lock()
	s.cache.regeo[key] = cacheItem[loc.ReverseGeocodeResult]{val: out, expTs: time.Now().Unix() + int64(s.regeoTTLms/1000)}
	s.cache.mu.Unlock()
	return out, nil
}

func (s *Service) aroundWithCache(ctx context.Context, lat, lon float64, p amap.AroundParams) (loc.PlaceAroundResult, error) {
	key := aroundKey{lat: round6(lat), lon: round6(lon), kw: p.Keywords, types: p.Types, radius: p.Radius}
	s.cache.mu.Lock()
	if it, ok := s.cache.around[key]; ok && it.expTs > time.Now().Unix() {
		v := it.val
		s.cache.mu.Unlock()
		return v, nil
	}
	s.cache.mu.Unlock()

	var lastErr error
	var out loc.PlaceAroundResult
	for i := 0; i < s.retryTimes; i++ {
		v, err := s.amapClient.PlaceAround(ctx, lat, lon, p)
		if err == nil {
			out = v
			break
		}
		lastErr = err
		t := time.Duration(s.retryBaseDelayMs*(1<<i)) * time.Millisecond
		time.Sleep(t)
	}
	if out.Count == 0 && lastErr != nil {
		return out, lastErr
	}
	// cache
	s.cache.mu.Lock()
	s.cache.around[key] = cacheItem[loc.PlaceAroundResult]{val: out, expTs: time.Now().Unix() + int64(s.aroundTTLms/1000)}
	s.cache.mu.Unlock()
	return out, nil
}

func (s *Service) weatherWithCache(ctx context.Context, p amap.WeatherParams) (loc.WeatherResult, error) {
	key := weatherKey{city: p.City, ext: p.Extensions}
	s.cache.mu.Lock()
	if it, ok := s.cache.weather[key]; ok && it.expTs > time.Now().Unix() {
		v := it.val
		s.cache.mu.Unlock()
		return v, nil
	}
	s.cache.mu.Unlock()

	var lastErr error
	var out loc.WeatherResult
	for i := 0; i < s.retryTimes; i++ {
		v, err := s.amapClient.Weather(ctx, p)
		if err == nil {
			out = v
			break
		}
		lastErr = err
		t := time.Duration(s.retryBaseDelayMs*(1<<i)) * time.Millisecond
		time.Sleep(t)
	}
	if (len(out.Lives) == 0 && len(out.Forecasts) == 0) && lastErr != nil {
		return out, lastErr
	}
	// cache
	s.cache.mu.Lock()
	s.cache.weather[key] = cacheItem[loc.WeatherResult]{val: out, expTs: time.Now().Unix() + int64(s.weatherTTLms/1000)}
	// 限制缓存大小（粗暴淘汰）
	if len(s.cache.weather) > 1024 {
		for k := range s.cache.weather {
			delete(s.cache.weather, k)
			break
		}
	}
	s.cache.mu.Unlock()
	return out, nil
}

// round6 将经纬度保留 1e-6 精度用于缓存key，以避免抖动
func round6(v float64) float64 {
	const m = 1e6
	iv := int64(v * m)
	return float64(iv) / m
}
