package amap

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	loc "new-robot-self/pkg/location"
)

// AroundParams 周边搜索可选参数，映射 AMap v5/place/around 文档
// 注意：中心点坐标需为 GCJ-02，location=lon,lat（经度在前，最多6位小数）
// - Keywords: 仅支持一个关键词，长度<=80
// - Types: 多个 typecode 用 "|" 分隔
// - Radius: 0-50000（米），默认 5000
// - SortRule: distance | weight
// - Region/CityLimit: 参见官方文档
// - ShowFields: 可选字段控制，逗号分隔
// - PageSize: 1-25，默认10
// - PageNum: 从1开始
type AroundParams struct {
	Keywords   string
	Types      string
	Radius     int
	SortRule   string
	Region     string
	CityLimit  bool
	ShowFields string
	PageSize   int
	PageNum    int
}

// PlaceAround 调用高德 v5/place/around 周边搜索
// 输入 GCJ-02 坐标，返回标准化的 PlaceAroundResult。
func (c *Client) PlaceAround(ctx context.Context, latGCJ02, lonGCJ02 float64, p AroundParams) (loc.PlaceAroundResult, error) {
	var empty loc.PlaceAroundResult
	if c == nil || c.Key == "" {
		return empty, errors.New("amap client or key missing")
	}

	u := fmt.Sprintf("%s/v5/place/around", strings.TrimRight(c.BaseURL, "/"))
	q := url.Values{}
	q.Set("key", c.Key)
	q.Set("location", fmt.Sprintf("%f,%f", lonGCJ02, latGCJ02)) // lon,lat
	if p.Keywords != "" {
		q.Set("keywords", p.Keywords)
	}
	if p.Types != "" {
		q.Set("types", p.Types)
	}
	if p.Radius > 0 {
		q.Set("radius", fmt.Sprintf("%d", p.Radius))
	}
	if p.SortRule != "" {
		q.Set("sortrule", p.SortRule)
	}
	if p.Region != "" {
		q.Set("region", p.Region)
	}
	if p.CityLimit {
		q.Set("city_limit", "true")
	}
	if p.ShowFields != "" {
		q.Set("show_fields", p.ShowFields)
	}
	if p.PageSize > 0 {
		q.Set("page_size", fmt.Sprintf("%d", p.PageSize))
	}
	if p.PageNum > 0 {
		q.Set("page_num", fmt.Sprintf("%d", p.PageNum))
	}
	n_ctx := context.Background()
	req, _ := http.NewRequestWithContext(n_ctx, http.MethodGet, u+"?"+q.Encode(), nil)
	start := time.Now()
	resp, err := c.HTTP.Do(req)
	if err != nil {
		return empty, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return empty, fmt.Errorf("amap around http %d", resp.StatusCode)
	}

	var ar aroundResponse
	if err := json.NewDecoder(resp.Body).Decode(&ar); err != nil {
		return empty, err
	}
	if ar.Status != "1" {
		return empty, fmt.Errorf("amap around failed: status=%s info=%s", ar.Status, ar.Info)
	}

	res := loc.PlaceAroundResult{
		TS:        time.Now().Unix(),
		LatGCJ02:  latGCJ02,
		LonGCJ02:  lonGCJ02,
		Radius:    p.Radius,
		Keywords:  p.Keywords,
		Types:     p.Types,
		SortRule:  p.SortRule,
		Region:    p.Region,
		CityLimit: p.CityLimit,
		PageSize:  p.PageSize,
		PageNum:   p.PageNum,
		Provider:  "amap",
	}

	// count 字段为字符串，转 int
	if ar.Count != "" {
		var cInt int
		if _, err := fmt.Sscanf(ar.Count, "%d", &cInt); err == nil {
			res.Count = cInt
		}
	}

	for _, p := range ar.POIs {
		res.POIs = append(res.POIs, loc.POI{
			ID:       p.ID,
			Name:     p.Name,
			Type:     p.Type,
			TypeCode: p.TypeCode,
			Address:  p.Address,
			Location: p.Location,
		})
	}
	_ = start // 预留耗时指标
	return res, nil
}

// --- AMap v5/place/around 响应结构（子集）---

type aroundResponse struct {
	Status   string      `json:"status"`
	Info     string      `json:"info"`
	InfoCode string      `json:"infocode"`
	Count    string      `json:"count"`
	POIs     []aroundPOI `json:"pois"`
}

type aroundPOI struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Location string `json:"location"` // lon,lat
	Type     string `json:"type"`
	TypeCode string `json:"typecode"`
	Address  string `json:"address"`
	// 其余可选字段可按需补充（如 tel、photos 等）
}
