package location

import "time"

// Config 位置模块运行配置
// 说明：本模块所有时间戳均为毫秒（int64, ms），但配置项中的 TTL/超时可能以“秒/毫秒”区分，字段注释已标明。
type Config struct {
	Enabled bool   `json:"enabled"` // 是否启用位置模块
	Mode    string `json:"mode"`    // 模式：amap_primary(高德优先) | local_primary(本地算法优先)

	Sampling SamplingConfig `json:"sampling"` // 采样与去抖配置
	Quality  QualityConfig  `json:"quality"`  // 质量评估配置
	Cache    CacheConfig    `json:"cache"`    // 缓存相关配置
	AMap     AMapConfig     `json:"amap"`     // 高德相关配置
}

type SamplingConfig struct {
	MinIntervalSec int     `json:"min_interval_sec"` // 最小时间间隔（秒），低于该间隔的点会被丢弃
	MinDistanceM   float64 `json:"min_distance_m"`   // 最小距离（米），移动距离低于该值不触发
	MinHeadingDeg  float64 `json:"min_heading_deg"`  // 最小航向变化（度），在速度>阈值时生效
}

type QualityConfig struct {
	MinScore float64 `json:"min_score"` // 质量评分阈值 [0,1]，低于该值可丢弃或延迟
}
type CacheConfig struct {
	RecentTTLSec int `json:"recent_ttl_sec"` // 最近可信位置缓存 TTL（秒）
}

type AMapConfig struct {
    Key          string        `json:"key"`            // 高德 Web 服务 KEY
    Timeout      time.Duration `json:"-"`              // 内部使用（由 TimeoutMs 转换）
    TimeoutMs    int           `json:"timeout_ms"`     // HTTP 超时（毫秒）
    RateLimitRPS int           `json:"rate_limit_rps"` // 每秒请求上限（RPS）
    RetryMs      []int         `json:"retry_ms"`       // 重试间隔（毫秒），建议指数退避

    Convert ConvertConfig `json:"convert"` // 坐标转换配置
    Regeo   RegeoConfig   `json:"regeo"`   // 逆地理配置
    Around  AroundConfig  `json:"around"`  // 周边搜索配置
}

type ConvertConfig struct {
    BatchMax            int      `json:"batch_max"`              // 单次批量最大坐标对数（<=40）
    RequestSep          string   `json:"request_sep"`            // 请求 locations 分隔符（默认“|”）
    ResponseSepFallback []string `json:"response_sep_fallbacks"` // 响应 locations 可能使用的分隔符集合（兼容 ";"、"|")
    CacheTTLSec         int      `json:"cache_ttl_sec"`          // 转换结果缓存 TTL（秒）
}

// RegeoConfig 逆地理配置
type RegeoConfig struct {
    Extensions string `json:"extensions"` // 返回内容层级：base | all（all 才有 POI/道路/路口/AOI）
    Radius     int    `json:"radius"`     // 搜索半径：0~3000（米），默认 1000
    PoiType    string `json:"poitype"`    // POI 类型编码，多个使用“|”分隔（extensions=all 生效）
    RoadLevel  int    `json:"roadlevel"`  // 道路等级：0 全部；1 主干道（extensions=all 生效）
    HomeOrCorp int    `json:"homeorcorp"` // 0 不干预；1 居家相关优先；2 公司相关优先（extensions=all 生效）
}

func DefaultConfig() Config {
    return Config{
        Enabled: true,
        Mode:    "local_primary", // 初期默认走本地算法，待拿到 Key 可切到 amap_primary
        Sampling: SamplingConfig{MinIntervalSec: 60, MinDistanceM: 200, MinHeadingDeg: 35},
        Quality:  QualityConfig{MinScore: 0.3},
        Cache:    CacheConfig{RecentTTLSec: 600},
        AMap: AMapConfig{
            Key:           "",
            TimeoutMs:     800,
            RateLimitRPS:  5,
            RetryMs:       []int{200, 500, 1000},
            Convert:       ConvertConfig{BatchMax: 40, RequestSep: "|", ResponseSepFallback: []string{";", "|"}, CacheTTLSec: 300},
            Regeo:         RegeoConfig{Extensions: "base", Radius: 1000, PoiType: "", RoadLevel: 0, HomeOrCorp: 0},
            Around:        AroundConfig{Radius: 5000, SortRule: "distance", PageSize: 10, PageNum: 1, CityLimit: false},
        },
    }
}

// AroundConfig 周边搜索配置（对应 v5/place/around）
type AroundConfig struct {
    Radius    int    `json:"radius"`               // 搜索半径 0-50000（米）
    SortRule  string `json:"sortrule"`             // 距离排序 distance | 综合排序 weight
    Region    string `json:"region,omitempty"`     // 区域权重（搭配 city_limit）
    CityLimit bool   `json:"city_limit,omitempty"` // 仅召回 region 对应区域
    PageSize  int    `json:"page_size"`            // 1-25，默认10
    PageNum   int    `json:"page_num"`             // 默认1
    Types     string `json:"types,omitempty"`      // typecode，用“|”分隔
}
