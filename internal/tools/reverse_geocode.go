package tools

import (
	"context"

	"github.com/bytedance/sonic"
	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/schema"

	amap "new-robot-self/pkg/location/provider/amap"
)

// ReverseGeocodeTool 封装 AMap 逆地理为统一工具
// 输入 GCJ-02 坐标，输出项目统一的 ReverseGeocodeResult（毫秒时间戳）
type ReverseGeocodeTool struct {
	BaseToolImpl
}

type reverseParams struct {
	Lat        float64 `json:"lat"`                  // GCJ-02 纬度，必填
	Lng        float64 `json:"lng"`                  // GCJ-02 经度，必填
	Extensions string  `json:"extensions,omitempty"` // base | all
	Radius     int     `json:"radius,omitempty"`     // 0~3000
	PoiType    string  `json:"poitype,omitempty"`    // 仅在 extensions=all 生效
	RoadLevel  int     `json:"roadlevel,omitempty"`  // 0|1，extensions=all
	HomeOrCorp int     `json:"homeorcorp,omitempty"` // 0|1|2，extensions=all
}

func NewReverseGeocodeTool() *ReverseGeocodeTool {
	return &ReverseGeocodeTool{
		BaseToolImpl: BaseToolImpl{
			Category:       "location",
			EnableLearning: true,
			ToolName:       "reverse_geocode",
		},
	}
}

func (t *ReverseGeocodeTool) Info(_ context.Context) (*schema.ToolInfo, error) {
	return &schema.ToolInfo{
		Name: "reverse_geocode",
		Desc: "基于 GCJ-02 坐标进行逆地理解析，返回格式化地址与行政区划，可选返回附近POI（extensions=all）",
		ParamsOneOf: schema.NewParamsOneOfByParams(map[string]*schema.ParameterInfo{
			"lat":        {Desc: "GCJ-02 纬度", Type: schema.Number, Required: true},
			"lng":        {Desc: "GCJ-02 经度", Type: schema.Number, Required: true},
			"extensions": {Desc: "返回层级 base|all，all 返回POI/道路等扩展信息", Type: schema.String, Required: false},
			"radius":     {Desc: "搜索半径 0~3000 米（extensions=all 时有效）", Type: schema.Integer, Required: false},
			"poitype":    {Desc: "POI 类型编码，多个用 '|' 分隔（extensions=all 有效）", Type: schema.String, Required: false},
			"roadlevel":  {Desc: "道路等级 0|1（extensions=all 有效）", Type: schema.Integer, Required: false},
			"homeorcorp": {Desc: "排序优化 0不干预|1居家|2公司（extensions=all 有效）", Type: schema.Integer, Required: false},
		}),
	}, nil
}

func (t *ReverseGeocodeTool) InvokableRun(ctx context.Context, argumentsInJSON string, _ ...tool.Option) (string, error) {
	p := reverseParams{}
	if err := sonic.UnmarshalString(argumentsInJSON, &p); err != nil {
		return "", err
	}
	key, err := GetAMapKey()
	if err != nil {
		return "", err
	}
	cli, err := amap.NewClient(key, 5000)
	if err != nil {
		return "", err
	}
	res, err := cli.ReverseGeocode(ctx, p.Lat, p.Lng, amap.RegeoParams{
		Extensions: p.Extensions,
		Radius:     p.Radius,
		PoiType:    p.PoiType,
		RoadLevel:  p.RoadLevel,
		HomeOrCorp: p.HomeOrCorp,
	})
	if err != nil {
		return "", err
	}
	b, _ := sonic.Marshal(res)
	return string(b), nil
}

// AfterRun 工具调用后的后处理
// 集成学习引擎，保存位置信息到知识库
func (t *ReverseGeocodeTool) AfterRun(ctx context.Context, result string, err error) {
	// 使用默认的BaseToolImpl后处理
	t.BaseToolImpl.AfterRun(ctx, result, err)
}
