package amap

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	loc "new-robot-self/pkg/location"
)

// WeatherParams 天气查询参数
// - City: 必填，城市 adcode
// - Extensions: base(实况) | all(预报)
// - Output: 仅支持 json，保持兼容
// 注：报告时间 reporttime 为字符串，需要解析为毫秒时间戳。
type WeatherParams struct {
	City       string
	Extensions string // base | all
	Output     string // json (default)
}

// Weather 查询高德天气 v3/weather/weatherInfo
func (c *Client) Weather(ctx context.Context, p WeatherParams) (loc.WeatherResult, error) {
	var empty loc.WeatherResult
	if c == nil || c.Key == "" {
		return empty, errors.New("amap client or key missing")
	}
	if strings.TrimSpace(p.City) == "" {
		return empty, errors.New("amap weather city(adcode) is empty")
	}

	u := fmt.Sprintf("%s/v3/weather/weatherInfo", strings.TrimRight(c.BaseURL, "/"))
	q := url.Values{}
	q.Set("key", c.Key)
	q.Set("city", p.City)
	if p.Extensions != "" {
		q.Set("extensions", p.Extensions)
	}
	if p.Output != "" {
		q.Set("output", p.Output)
	}
	n_ctx := context.Background()
	req, _ := http.NewRequestWithContext(n_ctx, http.MethodGet, u+"?"+q.Encode(), nil)
	resp, err := c.HTTP.Do(req)
	if err != nil {
		return empty, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return empty, fmt.Errorf("amap weather http %d", resp.StatusCode)
	}

	var aw weatherResponse
	if err := json.NewDecoder(resp.Body).Decode(&aw); err != nil {
		return empty, err
	}
	if aw.Status != "1" {
		return empty, fmt.Errorf("amap weather failed: status=%s info=%s", aw.Status, aw.Info)
	}

	res := loc.WeatherResult{
		TS:         time.Now().Unix(),
		Provider:   "amap",
		Extensions: p.Extensions,
		CityAdcode: p.City,
	}
	if aw.Count != "" {
		if n, err := strconv.Atoi(aw.Count); err == nil {
			res.Count = n
		}
	}

	for _, lv := range aw.Lives {
		res.Lives = append(res.Lives, loc.WeatherLive{
			Province:      lv.Province,
			City:          lv.City,
			Adcode:        lv.Adcode,
			Weather:       lv.Weather,
			TemperatureC:  lv.Temperature,
			WindDirection: lv.WindDirection,
			WindPower:     lv.WindPower,
			Humidity:      lv.Humidity,
			ReportTimeMS:  parseAMapTimeToMs(lv.ReportTime),
		})
	}
	for _, fc := range aw.Forecasts {
		f := loc.WeatherForecast{
			City:         fc.City,
			Province:     fc.Province,
			Adcode:       fc.Adcode,
			ReportTimeMS: parseAMapTimeToMs(fc.ReportTime),
		}
		for _, cst := range fc.Casts {
			f.Casts = append(f.Casts, WeatherCastToLoc(cst))
		}
		res.Forecasts = append(res.Forecasts, f)
	}
	return res, nil
}

// WeatherCastToLoc 映射 AMap 预报条目到本地结构
func WeatherCastToLoc(c castItem) loc.WeatherCast {
	return loc.WeatherCast{
		Date:         c.Date,
		Week:         c.Week,
		DayWeather:   c.DayWeather,
		NightWeather: c.NightWeather,
		DayTempC:     c.DayTemp,
		NightTempC:   c.NightTemp,
		DayWind:      c.DayWind,
		NightWind:    c.NightWind,
		DayPower:     c.DayPower,
		NightPower:   c.NightPower,
	}
}

func parseAMapTimeToMs(s string) int64 {
	if strings.TrimSpace(s) == "" {
		return 0
	}
	// AMap 示例："2020-06-24 11:00:00"
	const layout = "2006-01-02 15:04:05"
	t, err := time.ParseInLocation(layout, s, time.Local)
	if err != nil {
		return 0
	}
	return t.Unix()
}

// --- AMap Weather API JSON subset ---

type weatherResponse struct {
	Status    string         `json:"status"`
	Info      string         `json:"info"`
	InfoCode  string         `json:"infocode"`
	Count     string         `json:"count"`
	Lives     []liveItem     `json:"lives"`
	Forecasts []forecastItem `json:"forecasts"`
}

type liveItem struct {
	Province      string `json:"province"`
	City          string `json:"city"`
	Adcode        string `json:"adcode"`
	Weather       string `json:"weather"`
	Temperature   string `json:"temperature"`
	WindDirection string `json:"winddirection"`
	WindPower     string `json:"windpower"`
	Humidity      string `json:"humidity"`
	ReportTime    string `json:"reporttime"`
}

type forecastItem struct {
	City       string     `json:"city"`
	Adcode     string     `json:"adcode"`
	Province   string     `json:"province"`
	ReportTime string     `json:"reporttime"`
	Casts      []castItem `json:"casts"`
}

type castItem struct {
	Date         string `json:"date"`
	Week         string `json:"week"`
	DayWeather   string `json:"dayweather"`
	NightWeather string `json:"nightweather"`
	DayTemp      string `json:"daytemp"`
	NightTemp    string `json:"nighttemp"`
	DayWind      string `json:"daywind"`
	NightWind    string `json:"nightwind"`
	DayPower     string `json:"daypower"`
	NightPower   string `json:"nightpower"`
}
