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

// GeoParams 地理编码（结构化地址 -> 坐标）参数
// - Address: 必填，结构化地址（国家 省 市 区 街道 门牌 ...）
// - City: 可选，限定城市（中文名/拼音/citycode/adcode），不支持县级市
// - Output: 可选，默认 json
// 返回：GeocodeResult，包含 Items（可能多个候选）
type GeoParams struct {
	Address string
	City    string
	Output  string // json
}

// Geocode 请求高德 /v3/geocode/geo
func (c *Client) Geocode(ctx context.Context, p GeoParams) (loc.GeocodeResult, error) {
	var empty loc.GeocodeResult
	if c == nil || c.Key == "" {
		return empty, errors.New("amap client or key missing")
	}
	if strings.TrimSpace(p.Address) == "" {
		return empty, errors.New("amap geocode address is empty")
	}

	u := fmt.Sprintf("%s/v3/geocode/geo", strings.TrimRight(c.BaseURL, "/"))
	q := url.Values{}
	q.Set("key", c.Key)
	q.Set("address", p.Address)
	if p.City != "" {
		q.Set("city", p.City)
	}
	if p.Output != "" {
		q.Set("output", p.Output)
	}

	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, u+"?"+q.Encode(), nil)
	resp, err := c.HTTP.Do(req)
	if err != nil {
		return empty, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return empty, fmt.Errorf("amap geocode http %d", resp.StatusCode)
	}

	var ag geoResponse
	if err := json.NewDecoder(resp.Body).Decode(&ag); err != nil {
		return empty, err
	}
	if ag.Status != "1" {
		return empty, fmt.Errorf("amap geocode failed: status=%s info=%s", ag.Status, ag.Info)
	}

	res := loc.GeocodeResult{
		TS:       time.Now().Unix(),
		Provider: "amap",
		Address:  p.Address,
		City:     p.City,
	}
	if ag.Count != "" {
		if n, err := strconv.Atoi(ag.Count); err == nil {
			res.Count = n
		}
	}
	for _, g := range ag.Geocodes {
		res.Items = append(res.Items, loc.GeocodeItem{
			Country:  g.Country,
			Province: g.Province,
			City:     g.City,
			CityCode: g.CityCode,
			District: g.District,
			Street:   g.Street,
			Number:   g.Number,
			Adcode:   g.Adcode,
			Location: g.Location,
			Level:    g.Level,
		})
	}
	return res, nil
}

// --- AMap geocode/geo 响应结构（子集）---

type geoResponse struct {
	Status   string    `json:"status"`
	Info     string    `json:"info"`
	InfoCode string    `json:"infocode"`
	Count    string    `json:"count"`
	Geocodes []geoItem `json:"geocodes"`
}

type geoItem struct {
	Country  string   `json:"country"`
	Province string   `json:"province"`
	City     string   `json:"city"`
	CityCode string   `json:"citycode"`
	District string   `json:"district"`
	Street   []string `json:"street"`
	Number   []string `json:"number"`
	Adcode   string   `json:"adcode"`
	Location string   `json:"location"` // lon,lat
	Level    string   `json:"level"`
}
