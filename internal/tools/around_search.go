package tools

import (
	"context"

	"github.com/bytedance/sonic"
	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/schema"

	amap "new-robot-self/pkg/location/provider/amap"
)

// AroundSearchTool 封装 AMap v5/place/around 为统一工具
// 输入 GCJ-02 坐标，支持关键词/类型/半径等参数，返回标准化周边结果
type AroundSearchTool struct {
	BaseToolImpl
}

type aroundParams struct {
	Lat       float64 `json:"lat"`                 // GCJ-02 纬度，必填
	Lng       float64 `json:"lng"`                 // GCJ-02 经度，必填
	Keywords  string  `json:"keywords,omitempty"`
	Types     string  `json:"types,omitempty"`     // 多个 typecode 用 "|" 分隔
	Radius    int     `json:"radius,omitempty"`    // 0-50000，默认 5000
	SortRule  string  `json:"sortrule,omitempty"`  // distance | weight
	Region    string  `json:"region,omitempty"`
	CityLimit bool    `json:"city_limit,omitempty"`
	ShowFields string `json:"show_fields,omitempty"`
	PageSize  int     `json:"page_size,omitempty"` // 1-25
	PageNum   int     `json:"page_num,omitempty"`  // 从1开始
}

func NewAroundSearchTool() *AroundSearchTool {
	return &AroundSearchTool{
		BaseToolImpl: BaseToolImpl{
			Category:       "location",
			EnableLearning: true,
			ToolName:       "around_search",
		},
	}
}

func (t *AroundSearchTool) Info(_ context.Context) (*schema.ToolInfo, error) {
	return &schema.ToolInfo{
		Name: "around_search",
		Desc: "基于 GCJ-02 坐标执行高德周边搜索，支持关键词、类型与半径等参数",
		ParamsOneOf: schema.NewParamsOneOfByParams(map[string]*schema.ParameterInfo{
			"lat":        {Desc: "GCJ-02 纬度", Type: schema.Number, Required: true},
			"lng":        {Desc: "GCJ-02 经度", Type: schema.Number, Required: true},
			"keywords":   {Desc: "搜索关键词（可选，<=80）", Type: schema.String, Required: false},
			"types":      {Desc: "POI 类型编码，多个用 '|' 分隔", Type: schema.String, Required: false},
			"radius":     {Desc: "搜索半径 0~50000 米，默认 5000", Type: schema.Integer, Required: false},
			"sortrule":   {Desc: "排序规则 distance|weight", Type: schema.String, Required: false},
			"region":     {Desc: "限定区域/城市", Type: schema.String, Required: false},
			"city_limit": {Desc: "是否限制在城市内", Type: schema.Boolean, Required: false},
			"show_fields": {Desc: "返回字段控制（逗号分隔）", Type: schema.String, Required: false},
			"page_size":  {Desc: "分页大小 1-25", Type: schema.Integer, Required: false},
			"page_num":   {Desc: "页码，从1开始", Type: schema.Integer, Required: false},
		}),
	}, nil
}

func (t *AroundSearchTool) InvokableRun(ctx context.Context, argumentsInJSON string, _ ...tool.Option) (string, error) {
	p := aroundParams{}
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
	res, err := cli.PlaceAround(ctx, p.Lat, p.Lng, amap.AroundParams{
		Keywords:   p.Keywords,
		Types:      p.Types,
		Radius:     p.Radius,
		SortRule:   p.SortRule,
		Region:     p.Region,
		CityLimit:  p.CityLimit,
		ShowFields: p.ShowFields,
		PageSize:   p.PageSize,
		PageNum:    p.PageNum,
	})
	if err != nil {
		return "", err
	}
	b, _ := sonic.Marshal(res)
	return string(b), nil
}

// AfterRun 工具调用后的后处理
// 集成学习引擎，保存周边搜索结果到知识库
func (t *AroundSearchTool) AfterRun(ctx context.Context, result string, err error) {
	// 使用默认的BaseToolImpl后处理
	t.BaseToolImpl.AfterRun(ctx, result, err)
}
