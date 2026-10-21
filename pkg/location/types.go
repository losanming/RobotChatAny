package location

// 位置模块通用类型定义
// 约定：所有时间字段一律使用毫秒时间戳（int64, ms），与项目全局规范一致。

type RawLocation struct {
	TS             int64   `json:"ts"`                        // 定位数据时间戳（ms）
	LatWGS84       float64 `json:"lat_wgs84"`                 // 原始纬度（WGS-84）
	LonWGS84       float64 `json:"lon_wgs84"`                 // 原始经度（WGS-84）
	AccuracyMeter  float64 `json:"accuracy_m,omitempty"`      // 精度（米）
	SpeedMps       float64 `json:"speed_mps,omitempty"`       // 速度（米/秒）
	HeadingDeg     float64 `json:"heading_deg,omitempty"`     // 航向（度，0~360）
	SatelliteCount int     `json:"satellite_count,omitempty"` // 卫星数量（可选）
	Source         string  `json:"source"`                    // 来源：gps/wifi/cell/ip
	DeviceID       string  `json:"device_id"`                 // 设备ID
}

type TransformedLocation struct {
	TS         int64   `json:"ts"`                    // 转换结果时间戳（ms）
	LatGCJ02   float64 `json:"lat_gcj02"`             // 转换后纬度（GCJ-02）
	LonGCJ02   float64 `json:"lon_gcj02"`             // 转换后经度（GCJ-02）
	LatWGS84   float64 `json:"lat_wgs84"`             // 原始纬度（回显）
	LonWGS84   float64 `json:"lon_wgs84"`             // 原始经度（回显）
	InChina    bool    `json:"is_in_china"`           // 是否位于中国境内（用于决定是否转换）
	Quality    float64 `json:"quality"`               // 质量评分 [0,1]
	Adcode     string  `json:"adcode,omitempty"`      // 行政区编码（逆地理得到，预留）
	ReasonCode string  `json:"reason_code,omitempty"` // 原因码（如：amap_failed/out_of_china）
}

type LocationEvent struct {
	TS       int64   `json:"ts"`        // 事件时间戳（ms）
	DeviceID string  `json:"device_id"` // 设备ID
	Source   string  `json:"source"`    // 来源：gps/wifi/cell/ip
	Quality  float64 `json:"quality"`   // 质量评分 [0,1]

	LatWGS84 float64 `json:"lat_wgs84"`   // 原始纬度（WGS-84）
	LonWGS84 float64 `json:"lon_wgs84"`   // 原始经度（WGS-84）
	LatGCJ02 float64 `json:"lat_gcj02"`   // 转换后纬度（GCJ-02）
	LonGCJ02 float64 `json:"lon_gcj02"`   // 转换后经度（GCJ-02）
	InChina  bool    `json:"is_in_china"` // 是否在中国境内

	Geohash    string `json:"geohash,omitempty"`     // 可选：地理哈希，用于缓存与模糊化
	Adcode     string `json:"adcode,omitempty"`      // 行政区编码（逆地理）
	ReasonCode string `json:"reason_code,omitempty"` // 原因码（过滤/降级/回退）
}

// 简单点结构体用于批量转换
type WGS84Point struct {
	Lat float64 // 纬度
	Lon float64 // 经度
}

type GCJ02Point struct {
	Lat float64 // 纬度
	Lon float64 // 经度
}

// ReverseGeocodeResult 逆地理最小输出（与 Step2 设计一致，时间戳为毫秒）
type ReverseGeocodeResult struct {
	TS       int64    `json:"ts"`                 // 逆地理时间戳（ms）
	LatGCJ02 float64  `json:"lat_gcj02"`          // 输入纬度（GCJ-02）
	LonGCJ02 float64  `json:"lon_gcj02"`          // 输入经度（GCJ-02）
	Adcode   string   `json:"adcode"`             // 行政区编码
	Address  string   `json:"addr"`               // 格式化地址
	Country  string   `json:"country,omitempty"`  // 国家
	Province string   `json:"province,omitempty"` // 省/直辖市
	City     []string `json:"city,omitempty"`     // 市
	CityCode string   `json:"citycode,omitempty"` // 城市编码
	District string   `json:"district,omitempty"` // 区/县
	Township string   `json:"township,omitempty"` // 街道/乡镇
	TownCode string   `json:"towncode,omitempty"` // 街道编码
	POIs     []POI    `json:"poi_topk,omitempty"` // 附近POI（当extensions=all）
	Provider string   `json:"provider"`           // 数据提供方，如：amap
}

type POI struct {
	ID        string  `json:"id,omitempty"`           // POI唯一标识
	Name      string  `json:"name,omitempty"`         // 名称
	Type      string  `json:"type,omitempty"`         // 类型（TypeCode）
	TypeCode  string  `json:"typecode,omitempty"`     // 分类编码
	Address   string  `json:"address,omitempty"`      // 详细地址
	DistanceM float64 `json:"distance,omitempty"`     // 与输入点距离（米）
	Location  string  `json:"location,omitempty"`     // 位置字符串："lon,lat"
	BizArea   string  `json:"businessarea,omitempty"` // 所在商圈
}

// PlaceAroundResult 周边搜索最小输出（与全局约定一致，时间戳为毫秒）
// 说明：仅保留通用字段，详细字段可通过 show_fields 在 Provider 内部扩展返回
type PlaceAroundResult struct {
	TS        int64   `json:"ts"`               // 查询时间戳（ms）
	LatGCJ02  float64 `json:"lat_gcj02"`        // 中心点纬度（GCJ-02）
	LonGCJ02  float64 `json:"lon_gcj02"`        // 中心点经度（GCJ-02）
	Radius    int     `json:"radius,omitempty"` // 搜索半径（米）
	Keywords  string  `json:"keywords,omitempty"`
	Types     string  `json:"types,omitempty"`
	SortRule  string  `json:"sortrule,omitempty"` // distance | weight
	Region    string  `json:"region,omitempty"`
	CityLimit bool    `json:"city_limit,omitempty"`
	PageSize  int     `json:"page_size,omitempty"`
	PageNum   int     `json:"page_num,omitempty"`
	Count     int     `json:"count"`
	POIs      []POI   `json:"pois"`
	Provider  string  `json:"provider"`
}

// WeatherLive 实况天气（extensions=base）
type WeatherLive struct {
	Province      string `json:"province"`
	City          string `json:"city"`
	Adcode        string `json:"adcode"`
	Weather       string `json:"weather"`
	TemperatureC  string `json:"temperature"` // 摄氏度，字符串
	WindDirection string `json:"winddirection"`
	WindPower     string `json:"windpower"`
	Humidity      string `json:"humidity"`
	ReportTimeMS  int64  `json:"reporttime_ms"` // AMap 字符串时间解析为毫秒
}

// WeatherCast 预报天气条目（extensions=all）
type WeatherCast struct {
	Date         string `json:"date"`
	Week         string `json:"week"`
	DayWeather   string `json:"dayweather"`
	NightWeather string `json:"nightweather"`
	DayTempC     string `json:"daytemp"`
	NightTempC   string `json:"nighttemp"`
	DayWind      string `json:"daywind"`
	NightWind    string `json:"nightwind"`
	DayPower     string `json:"daypower"`
	NightPower   string `json:"nightpower"`
}

// WeatherForecast 预报集合（最多四天）
type WeatherForecast struct {
	City         string        `json:"city"`
	Province     string        `json:"province"`
	Adcode       string        `json:"adcode"`
	ReportTimeMS int64         `json:"reporttime_ms"`
	Casts        []WeatherCast `json:"casts"`
}

type WeatherResult struct {
	TS         int64             `json:"ts"`          // 查询时间戳（ms）
	Provider   string            `json:"provider"`    // amap
	Extensions string            `json:"extensions"`  // base | all
	CityAdcode string            `json:"city_adcode"` // 查询的城市 adcode
	Count      int               `json:"count"`
	Lives      []WeatherLive     `json:"lives,omitempty"`
	Forecasts  []WeatherForecast `json:"forecasts,omitempty"`
}

// GeocodeItem 地理编码条目（address -> location）
type GeocodeItem struct {
	Country  string   `json:"country,omitempty"`
	Province string   `json:"province,omitempty"`
	City     string   `json:"city,omitempty"`
	CityCode string   `json:"citycode,omitempty"`
	District string   `json:"district,omitempty"`
	Street   []string `json:"street,omitempty"`
	Number   []string `json:"number,omitempty"`
	Adcode   string   `json:"adcode,omitempty"`
	Location string   `json:"location,omitempty"` // "lon,lat"
	Level    string   `json:"level,omitempty"`
}

// GeocodeResult 地理编码统一输出
type GeocodeResult struct {
	TS       int64         `json:"ts"`             // 查询时间戳（ms）
	Provider string        `json:"provider"`       // amap
	Address  string        `json:"address"`        // 输入的结构化地址
	City     string        `json:"city,omitempty"` // 可选：限定城市（中文/citycode/adcode）
	Count    int           `json:"count"`
	Items    []GeocodeItem `json:"items"`
}
