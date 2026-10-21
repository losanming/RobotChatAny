package brain

import (
	"context"
	"fmt"
	"math"
	"new-robot-self/internal/awareness"
	loc "new-robot-self/pkg/location"
	"new-robot-self/pkg/location/transform"
	"sync"
	"time"
)

// GPSService 负责接收硬件上报的原始GPS数据，汇聚为 awareness.Context 并注入 Brain
// 设计：
// - 通过内存通道缓冲上报；
type GPSService struct {
	brain   *Brain
	aware   *awareness.Service
	inCh    chan gpsReport
	wg      sync.WaitGroup
	ctx     context.Context
	cancel  context.CancelFunc
	started bool

	// 配置项
	queueSize int
	gatherTO  time.Duration // awareness.Gather 超时

	// 资源优化配置
	cfg GPSPolicy

	// 每用户最近一次上下文与时间
	mu      sync.Mutex
	perUser map[string]*userTrace
}

type gpsReport struct {
	UserID string
	Raw    loc.RawLocation
}

type GPSPolicy struct {
	MinRegeoIntervalMs int64   // 相邻逆地理最小间隔
	MinRegeoMoveMeters float64 // 位移超过阈值才触发逆地理
	AroundTTLms        int64   // 周边结果最短复用TTL
	WeatherTTLms       int64   // 天气结果最短复用TTL
	AroundKeywords     string
	AroundTypes        string
	AroundRadius       int
	WeatherExtensions  string // base|all
}

func DefaultGPSPolicy() GPSPolicy {
	return GPSPolicy{
		MinRegeoIntervalMs: 60_000,      // 60s
		MinRegeoMoveMeters: 50.0,        // 50米
		AroundTTLms:        5 * 60_000,  // 5分钟
		WeatherTTLms:       15 * 60_000, // 15分钟
		AroundKeywords:     "",
		AroundTypes:        "050000|060000",
		AroundRadius:       2000,
		WeatherExtensions:  "base",
	}
}

type userTrace struct {
	lastCtx       *awareness.Context
	lastRegeoTs   int64
	lastAroundTs  int64
	lastWeatherTs int64
	lastLat       float64
	lastLon       float64
}

// 可选配置项
type GPSOption func(*GPSService)

// WithGPSQueueSize 配置上报缓冲队列大小
func WithGPSQueueSize(n int) GPSOption {
	return func(s *GPSService) {
		if n > 0 {
			s.queueSize = n
		}
	}
}

// WithGPSGatherTimeout 配置感知汇聚超时时间
func WithGPSGatherTimeout(d time.Duration) GPSOption {
	return func(s *GPSService) {
		if d > 0 {
			s.gatherTO = d
		}
	}
}

// WithGPSPolicy 配置运行策略
func WithGPSPolicy(p GPSPolicy) GPSOption {
	return func(s *GPSService) {
		s.cfg = p
	}
}

// NewGPSService 构造函数
// parent 用于派生可取消上下文；brain/aware 为外部注入依赖
func NewGPSService(parent context.Context, brain *Brain, aware *awareness.Service, opts ...GPSOption) (*GPSService, error) {
	if brain == nil {
		return nil, fmt.Errorf("brain is nil")
	}
	if aware == nil {
		return nil, fmt.Errorf("awareness service is nil")
	}
	if parent == nil {
		parent = context.Background()
	}

	ctx, cancel := context.WithCancel(parent)

	s := &GPSService{
		brain:     brain,
		aware:     aware,
		inCh:      nil, // 延后按 queueSize 设置
		ctx:       ctx,
		cancel:    cancel,
		started:   false,
		queueSize: 256,
		gatherTO:  3 * time.Second,
		cfg:       DefaultGPSPolicy(),
		perUser:   make(map[string]*userTrace),
	}

	// 应用可选项
	for _, o := range opts {
		if o != nil {
			o(s)
		}
	}

	// 按最终队列大小创建通道
	s.inCh = make(chan gpsReport, s.queueSize)

	return s, nil
}

// 粗略 Haversine 距离（米）
func haversineMeters(lat1, lon1, lat2, lon2 float64) float64 {
	const R = 6371000.0
	toRad := func(d float64) float64 { return d * math.Pi / 180 }
	dLat := toRad(lat2 - lat1)
	dLon := toRad(lon2 - lon1)
	a := math.Sin(dLat/2)*math.Sin(dLat/2) + math.Cos(toRad(lat1))*math.Cos(toRad(lat2))*math.Sin(dLon/2)*math.Sin(dLon/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	return R * c
}

// Start 启动后台消费
func (s *GPSService) Start() {
	if s.started {
		return
	}
	s.started = true
	s.wg.Add(1)
	go s.loop()
}

// Stop 优雅停止
func (s *GPSService) Stop() {
	if !s.started {
		return
	}
	s.cancel()
	s.wg.Wait()
	s.started = false
}

// ReportRaw 设备上报原始定位（WGS-84），时间戳为毫秒
func (s *GPSService) ReportRaw(userID string, raw loc.RawLocation) {
	if !s.started {
		return
	}
	select {
	case s.inCh <- gpsReport{UserID: userID, Raw: raw}:
	default:
		// 丢弃过载数据，避免阻塞（也可改为阻塞或最新覆盖策略）
	}
}

func (s *GPSService) loop() {
	defer s.wg.Done()
	for {
		select {
		case <-s.ctx.Done():
			return
		case r := <-s.inCh:
			s.handleReport(r)
		}
	}
}

// SetPolicy 运行时更新资源策略（线程安全）
func (s *GPSService) SetPolicy(p GPSPolicy) {
	s.mu.Lock()
	s.cfg = p
	s.mu.Unlock()
}

// GetPolicy 获取当前策略快照
func (s *GPSService) GetPolicy() GPSPolicy {
	s.mu.Lock()
	p := s.cfg
	s.mu.Unlock()
	return p
}

// 内部：处理单条上报，采用条件触发策略
func (s *GPSService) handleReport(r gpsReport) {
	now := time.Now().Unix()

	s.mu.Lock()
	ut, ok := s.perUser[r.UserID]
	if !ok {
		ut = &userTrace{}
		s.perUser[r.UserID] = ut
	}
	// 快速复制上一次上下文引用，避免长时间持锁
	lastCtx := ut.lastCtx
	lastRegeoTs := ut.lastRegeoTs
	lastAroundTs := ut.lastAroundTs
	lastWeatherTs := ut.lastWeatherTs
	lastLat := ut.lastLat
	lastLon := ut.lastLon
	s.mu.Unlock()

	// 计算位移距离
	dist := 0.0
	if lastLat != 0 || lastLon != 0 {
		dist = haversineMeters(lastLat, lastLon, r.Raw.LatWGS84, r.Raw.LonWGS84)
	}

	// 始终快速生成基础位置上下文（无需外部API）：WGS84->GCJ02
	latGCJ, lonGCJ := transform.WGS84ToGCJ02(r.Raw.LatWGS84, r.Raw.LonWGS84)
	inChina := transform.InChina(r.Raw.LatWGS84, r.Raw.LonWGS84)
	baseCtx := awareness.Context{
		TS:       now,
		DeviceID: r.Raw.DeviceID,
		Source:   r.Raw.Source,
		LocTrans: loc.TransformedLocation{
			TS:       now,
			LatGCJ02: latGCJ,
			LonGCJ02: lonGCJ,
			LatWGS84: r.Raw.LatWGS84,
			LonWGS84: r.Raw.LonWGS84,
			InChina:  inChina,
			Quality:  r.Raw.AccuracyMeter,
		},
	}

	// 判断是否需要刷新重服务：逆地理/周边/天气
	needRegeo := (now-lastRegeoTs >= int64(s.cfg.MinRegeoIntervalMs/1000)) || (dist >= s.cfg.MinRegeoMoveMeters) || lastCtx == nil
	needAround := (now-lastAroundTs >= int64(s.cfg.AroundTTLms/1000)) || needRegeo // 周边通常伴随地址变化
	needWeather := (now-lastWeatherTs >= int64(s.cfg.WeatherTTLms/1000))

	if needRegeo || needAround || needWeather {
		ctxTO, cancel := context.WithTimeout(s.ctx, s.gatherTO)
		awIn := awareness.Input{
			RawLoc:            r.Raw,
			NeedAround:        needAround,
			AroundKeywords:    s.cfg.AroundKeywords,
			AroundTypes:       s.cfg.AroundTypes,
			AroundRadius:      s.cfg.AroundRadius,
			WeatherExtensions: s.cfg.WeatherExtensions,
		}
		res, err := s.aware.Gather(ctxTO, awIn)
		cancel()
		if err == nil {
			s.brain.UpdateFromAwareness(s.ctx, res.Ctx)
			// 更新 per-user 轨迹
			s.mu.Lock()
			ut.lastCtx = &res.Ctx
			ut.lastRegeoTs = now
			if needAround {
				ut.lastAroundTs = now
			}
			if needWeather {
				ut.lastWeatherTs = now
			}
			ut.lastLat = r.Raw.LatWGS84
			ut.lastLon = r.Raw.LonWGS84
			s.mu.Unlock()
			return
		}
		// 如果外部失败，仍走快速路径
	}

	// 快速路径：复用上一次的 Regeo/Around/Weather（若存在）
	if lastCtx != nil {
		baseCtx.Regeo = lastCtx.Regeo
		baseCtx.Around = lastCtx.Around
		baseCtx.Weather = lastCtx.Weather
	}
	s.brain.UpdateFromAwareness(s.ctx, baseCtx)
	// 更新最近基础位置
	s.mu.Lock()
	ut.lastCtx = &baseCtx
	ut.lastLat = r.Raw.LatWGS84
	ut.lastLon = r.Raw.LonWGS84
	s.mu.Unlock()
}
