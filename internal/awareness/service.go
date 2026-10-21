package awareness

import (
	"context"
	"errors"
	"log/slog"
	"strings"
	"time"

	loc "new-robot-self/pkg/location"
	amap "new-robot-self/pkg/location/provider/amap"
	transform "new-robot-self/pkg/location/transform"
)

// Service 负责汇聚定位/周边/天气（可选）与视觉（占位），形成可喂给大模型的上下文与 Prompt。
// 依赖通过构造函数注入，便于在外层统一治理（限流/重试/缓存/埋点）。

type Service struct {
	amapClient          *amap.Client
	defaultAroundRadius int
	weatherTimeoutMs    int
	cache               *cacheStore
	regeoTTLms          int64
	aroundTTLms         int64
	weatherTTLms        int64
	retryTimes          int
	retryBaseDelayMs    int
}

// Option 函数式配置
type Option func(*Service)

// WithDefaultAroundRadius 配置周边默认半径（米）
func WithDefaultAroundRadius(r int) Option {
	return func(s *Service) {
		if r > 0 {
			s.defaultAroundRadius = r
		}
	}
}

// WithWeatherTimeoutMs 配置天气查询超时时间（毫秒）
func WithWeatherTimeoutMs(ms int) Option {
	return func(s *Service) {
		if ms > 0 {
			s.weatherTimeoutMs = ms
		}
	}
}

// WithRetry 配置重试次数与基线回退（毫秒）
func WithRetry(times, baseDelayMs int) Option {
	return func(s *Service) {
		if times >= 0 {
			s.retryTimes = times
		}
		if baseDelayMs > 0 {
			s.retryBaseDelayMs = baseDelayMs
		}
	}
}

// WithTTL 配置各类缓存 TTL（毫秒）
func WithTTL(regeo, around, weather int64) Option {
	return func(s *Service) {
		if regeo > 0 {
			s.regeoTTLms = regeo
		}
		if around > 0 {
			s.aroundTTLms = around
		}
		if weather > 0 {
			s.weatherTTLms = weather
		}
	}
}

// NewService 构造
func NewService(amapClient *amap.Client, opts ...Option) (*Service, error) {
	if amapClient == nil {
		return nil, errors.New("amap client is nil")
	}
	s := &Service{
		amapClient:          amapClient,
		defaultAroundRadius: 3000,
		weatherTimeoutMs:    2500,
		cache:               newCacheStore(),
		regeoTTLms:          60_000,      // 1 min
		aroundTTLms:         5 * 60_000,  // 5 min
		weatherTTLms:        15 * 60_000, // 15 min
		retryTimes:          2,
		retryBaseDelayMs:    200,
	}
	for _, o := range opts {
		o(s)
	}
	return s, nil
}

// Gather 汇聚数据，生成上下文与 Prompt
func (s *Service) Gather(ctx context.Context, in Input) (Result, error) {
	if in.RawLoc.TS == 0 || strings.TrimSpace(in.RawLoc.DeviceID) == "" {
		return Result{}, errors.New("invalid input: missing RawLoc.TS or DeviceID")
	}

	// 1) 坐标转换：WGS-84 -> GCJ-02（若不在中国大陆范围则原样返回）
	latW, lonW := in.RawLoc.LatWGS84, in.RawLoc.LonWGS84
	latGCJ, lonGCJ := transform.WGS84ToGCJ02(latW, lonW)
	locTrans := loc.TransformedLocation{
		TS:         time.Now().Unix(),
		LatGCJ02:   latGCJ,
		LonGCJ02:   lonGCJ,
		LatWGS84:   latW,
		LonWGS84:   lonW,
		InChina:    transform.InChina(latW, lonW),
		Quality:    in.RawLoc.AccuracyMeter,
		Adcode:     "",
		ReasonCode: "",
	}

	// 2) 逆地理（基础地址 + 可选POI）
	regeo, err := s.regeoWithCache(ctx, latGCJ, lonGCJ, amap.RegeoParams{Extensions: "all", Radius: 1000})
	if err != nil {
		slog.Error("ReverseGeocode error:", err)
		// 降级：返回空的 regeo
		r := amapEmptyRegeo(latGCJ, lonGCJ)
		// 将空结果赋值回来
		regeo = r
	}

	// 3) 周边（可选）
	var aroundResPtr *loc.PlaceAroundResult
	if in.NeedAround {
		radius := in.AroundRadius
		if radius <= 0 {
			radius = s.defaultAroundRadius
		}
		ar, err := s.aroundWithCache(ctx, latGCJ, lonGCJ, amap.AroundParams{
			Keywords:   in.AroundKeywords,
			Types:      in.AroundTypes,
			Radius:     radius,
			SortRule:   "distance",
			PageSize:   10,
			PageNum:    1,
			ShowFields: "",
		})
		if err == nil {
			aroundResPtr = &ar
		}
	}

	// 4) 天气（以 adcode 为主，缺失时尝试 citycode）
	var weatherPtr *loc.WeatherResult
	{
		adcode := regeo.Adcode
		if adcode != "" {
			wCtx, cancel := context.WithTimeout(ctx, time.Duration(s.weatherTimeoutMs)*time.Millisecond)
			defer cancel()
			ext := in.WeatherExtensions
			if strings.TrimSpace(ext) == "" {
				ext = "base"
			}
			w, err := s.weatherWithCache(wCtx, amap.WeatherParams{City: adcode, Extensions: ext})
			if err == nil {
				weatherPtr = &w
			}
		}
	}

	// 5) 视觉（占位，当前不做同步识别；预留外部注入结果通道）
	var vision *ImageInsight
	if len(in.ImageJPEG) > 0 {
		vision = &ImageInsight{Summary: "(占位) 待接入图片理解模块", Labels: nil, Objects: nil}
	}

	// 6) 组装上下文
	ctxOut := Context{
		TS:       time.Now().Unix(),
		DeviceID: in.RawLoc.DeviceID,
		Source:   in.RawLoc.Source,
		LocTrans: locTrans,
		Regeo:    &regeo,
		Around:   aroundResPtr,
		Weather:  weatherPtr,
		Vision:   vision,
	}

	// 7) 生成 Prompt
	prompt := BuildPrompt(ctxOut)
	return Result{Ctx: ctxOut, Prompt: prompt}, nil
}

func amapEmptyRegeo(latGCJ, lonGCJ float64) loc.ReverseGeocodeResult {
	return loc.ReverseGeocodeResult{
		TS:       time.Now().Unix(),
		LatGCJ02: latGCJ,
		LonGCJ02: lonGCJ,
		Provider: "amap",
	}
}
