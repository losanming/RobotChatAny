package amap

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	loc "new-robot-self/pkg/location"
)

// RegeoParams 逆地理可选参数
// 说明：与高德官方文档一致。其中部分参数仅在 extensions=all 时生效。
// - Extensions：返回内容层级，base 仅地址；all 返回地址+POI+道路+路口+AOI
// - Radius：搜索半径（米），0~3000，默认 1000
// - PoiType：POI 类型编码，多个用 "|" 分隔，仅在 extensions=all 时生效
// - roadLevel：道路等级，0 全部；1 主干道，仅在 extensions=all 时生效
// - HomeOrCorp：排序优化，0 不干预；1 居家相关优先；2 公司相关优先，仅在 extensions=all 时生效
type RegeoParams struct {
	Extensions string // base | all
	Radius     int    // 0~3000 meters
	PoiType    string // POI typecode list separated by '|', effective when Extensions=all
	RoadLevel  int    // 0 or 1, effective when Extensions=all
	HomeOrCorp int    // 0/1/2, effective when Extensions=all
}

// strOrArr 兼容高德返回的字符串或字符串数组
type strOrArr string

func (s *strOrArr) UnmarshalJSON(b []byte) error {
	// 如果是数组，取第一个元素；如果是字符串，直接解析
	if len(b) > 0 && b[0] == '[' {
		var arr []string
		if err := json.Unmarshal(b, &arr); err != nil {
			return err
		}
		if len(arr) > 0 {
			*s = strOrArr(arr[0])
		} else {
			*s = ""
		}
		return nil
	}
	var v string
	if err := json.Unmarshal(b, &v); err != nil {
		return err
	}
	*s = strOrArr(v)
	return nil
}

// ReverseGeocode 基于 GCJ-02 坐标调用高德逆地理，并映射为项目统一结构。
// 用法与注意：
// - 传入坐标须为 GCJ-02；如源为 WGS-84，请先转换（见 pkg/location/transform）。
// - location 参数格式严格为 "经度,纬度"（lon,lat），且小数位不超过 6 位。
// - 当 p.Extensions = "all" 时，才会返回 POI/道路/路口/AOI 等扩展信息。
// - 返回的时间戳 TS 均为毫秒（int64）。
// - 错误处理：status != "1" 或 HTTP 非 200 视为失败，error 中包含 info 便于定位。
func (c *Client) ReverseGeocode(ctx context.Context, latGCJ02, lonGCJ02 float64, p RegeoParams) (loc.ReverseGeocodeResult, error) {
	var empty loc.ReverseGeocodeResult
	if c == nil || c.Key == "" {
		return empty, errors.New("amap client or key missing")
	}

	u := fmt.Sprintf("%s/v3/geocode/regeo", strings.TrimRight(c.BaseURL, "/"))
	q := url.Values{}
	q.Set("key", c.Key)
	q.Set("location", fmt.Sprintf("%f,%f", lonGCJ02, latGCJ02)) // lon,lat
	if p.Extensions != "" {
		q.Set("extensions", p.Extensions)
	}
	if p.Radius > 0 {
		q.Set("radius", fmt.Sprintf("%d", p.Radius))
	}
	if p.PoiType != "" {
		q.Set("poitype", p.PoiType)
	}
	if p.RoadLevel == 0 || p.RoadLevel == 1 {
		q.Set("roadlevel", fmt.Sprintf("%d", p.RoadLevel))
	}
	if p.HomeOrCorp == 0 || p.HomeOrCorp == 1 || p.HomeOrCorp == 2 {
		q.Set("homeorcorp", fmt.Sprintf("%d", p.HomeOrCorp))
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
		return empty, fmt.Errorf("amap regeo http %d", resp.StatusCode)
	}
	var ar regeoResponse
	all, err := io.ReadAll(resp.Body)
	if err != nil {
		return empty, err
	}
	err = json.Unmarshal(all, &ar)
	if err != nil {
		return empty, err
	}
	//if err := json.NewDecoder(resp.Body).Decode(&ar); err != nil {
	//	return empty, err
	//}
	if ar.Status != "1" {
		return empty, fmt.Errorf("amap regeo failed: status=%s info=%s", ar.Status, ar.Info)
	}
	res := loc.ReverseGeocodeResult{
		TS:       time.Now().Unix(),
		LatGCJ02: latGCJ02,
		LonGCJ02: lonGCJ02,
		Provider: "amap",
	}
	if ar.Regeo.AddressComponent != nil {
		ac := ar.Regeo.AddressComponent
		res.Country = ac.Country
		res.Province = ac.Province
		res.CityCode = ac.CityCode
		res.District = ac.District
		res.Township = ac.Township
		res.TownCode = ac.TownCode
		res.Adcode = ac.Adcode
	}
	res.Address = ar.Regeo.FormattedAddress
	// POIs（尽力而为，只有在 extensions=all 时才会有数据）
	for _, p := range ar.Regeo.POIs {
		res.POIs = append(res.POIs, loc.POI{
			ID:        p.ID,
			Name:      p.Name,
			Type:      p.Type,
			DistanceM: p.Distance,
			Location:  p.Location,
			BizArea:   p.BusinessArea,
		})
	}
	_ = start // 预留耗时指标采集位置
	return res, nil
}

// --- AMap JSON 响应结构（子集，已满足当前需求） ---

// regeoResponse 为高德逆地理接口的顶层响应结构（已截取所需字段）
type regeoResponse struct {
	Status string     `json:"status"`
	Info   string     `json:"info"`
	Regeo  *regeoBody `json:"regeocode"`
}

type regeoBody struct {
	FormattedAddress string            `json:"formatted_address"`
	AddressComponent *addressComponent `json:"addressComponent"`
	POIs             []poiItem         `json:"pois"`
}

type addressComponent struct {
	Country  string   `json:"country"`
	Province string   `json:"province"`
	City     []string `json:"-"`
	CityCode string   `json:"citycode"`
	District string   `json:"district"`
	Adcode   string   `json:"adcode"`
	Township string   `json:"township"`
	TownCode string   `json:"towncode"`
}

type poiItem struct {
	ID           string   `json:"id"`
	Name         string   `json:"name"`
	Type         string   `json:"type"`
	Tel          []string `json:"-"`
	Distance     float64  `json:"distance,string"`
	Direction    string   `json:"direction"`
	Address      string   `json:"-"`
	Location     string   `json:"location"` // lon,lat
	BusinessArea string   `json:"-"`
}
