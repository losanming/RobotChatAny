package tools

import (
	"context"
	"errors"

	"github.com/bytedance/sonic"
	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/schema"

	public_city "new-robot-self/public/cityCode"
)

// CityAdcodeTool 提供根据地名查询城市 adcode 的功能
// 输入城市名称（可为常见中文城市名，如"北京/上海/广州市/成都市"等），返回 {name, adcode}
// 注意：只查命中词典内的标准行政区名称集合（来源于 asset/cities_adcode.json）
type CityAdcodeTool struct {
	BaseToolImpl
}

type cityQueryParams struct {
	Name string `json:"name"` // 城市名称，如 "北京"
}

type cityQueryResult struct {
	Name   string `json:"name"`
	Adcode string `json:"adcode"`
}

func NewCityAdcodeTool() *CityAdcodeTool {
	return &CityAdcodeTool{
		BaseToolImpl: BaseToolImpl{
			Category:       "location",
			EnableLearning: true,
			ToolName:       "city_adcode",
		},
	}
}

func (t *CityAdcodeTool) Info(_ context.Context) (*schema.ToolInfo, error) {
	return &schema.ToolInfo{
		Name: "city_adcode",
		Desc: "根据地名返回城市 adcode（城市编码）。若未命中将返回错误",
		ParamsOneOf: schema.NewParamsOneOfByParams(map[string]*schema.ParameterInfo{
			"name": {Desc: "城市中文名称，如 北京/上海/深圳/杭州市 等", Type: schema.String, Required: true},
		}),
	}, nil
}

func (t *CityAdcodeTool) InvokableRun(_ context.Context, argumentsInJSON string, _ ...tool.Option) (string, error) {
	if public_city.CityStore == nil {
		return "", errors.New("city store not initialized")
	}
	p := cityQueryParams{}
	if err := sonic.UnmarshalString(argumentsInJSON, &p); err != nil {
		return "", err
	}
	if p.Name == "" {
		return "", errors.New("name is required")
	}
	// 优先使用模糊匹配，兼容"北京/北京市"等差异
	entry, ok := public_city.CityStore.GetByFuzzy(p.Name)
	if !ok {
		return "", errors.New("city not found")
	}
	res := cityQueryResult{Name: entry.Name, Adcode: entry.Adcode}
	b, _ := sonic.Marshal(res)
	return string(b), nil
}

// AfterRun 工具调用后的后处理
// 集成学习引擎，保存城市查询结果到知识库
func (t *CityAdcodeTool) AfterRun(ctx context.Context, result string, err error) {
	// 使用默认的BaseToolImpl后处理
	t.BaseToolImpl.AfterRun(ctx, result, err)
}
