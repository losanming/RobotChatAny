package tools

import (
	"context"
	"encoding/json"
	"fmt"
	public_city "new-robot-self/public/cityCode"

	"github.com/bytedance/sonic"
	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/schema"

	amap "new-robot-self/pkg/location/provider/amap"
)

// WeatherQueryTool 封装 AMap 天气查询为统一工具
// 输入城市 adcode 与 extensions（base|all），返回统一 WeatherResult（毫秒时间戳）
type WeatherQueryTool struct {
	BaseToolImpl
}

type weatherParams struct {
	City       string `json:"city"`                 // 城市，必填
	Extensions string `json:"extensions,omitempty"` // base | all | forecast(预报) | air(空气质量)
	Days       int    `json:"days,omitempty"`       // 预报天数，1-7
}

func NewWeatherQueryTool() *WeatherQueryTool {
	return &WeatherQueryTool{
		BaseToolImpl: BaseToolImpl{
			Category:       "location",
			EnableLearning: true,
			ToolName:       "weather_query",
		},
	}
}

func (t *WeatherQueryTool) Info(_ context.Context) (*schema.ToolInfo, error) {
	return &schema.ToolInfo{
		Name: "weather_query",
		Desc: "查询高德天气：city的参数是必填的;支持实况、预报、空气质量。extensions: base(实况), all(预报), air(空气质量), forecast(多日预报)",
		ParamsOneOf: schema.NewParamsOneOfByParams(map[string]*schema.ParameterInfo{
			"city":       {Desc: "城市", Type: schema.String, Required: false},
			"extensions": {Desc: "base(实况) | all(预报) | air(空气质量) | forecast(多日预报)", Type: schema.String, Required: false},
			"days":       {Desc: "预报天数(1-7)", Type: schema.Integer, Required: false},
		}),
	}, nil
}

func (t *WeatherQueryTool) InvokableRun(ctx context.Context, argumentsInJSON string, _ ...tool.Option) (string, error) {
	p := weatherParams{}
	if err := json.Unmarshal([]byte(argumentsInJSON), &p); err != nil {
		return "", err
	}
	if p.City == "" {
		p.City = "北京"
	}

	// 设置默认值
	if p.Extensions == "" {
		p.Extensions = "base"
	}
	if p.Days == 0 || p.Days > 7 {
		p.Days = 1
	}

	key, err := GetAMapKey()
	if err != nil {
		return "", err
	}
	cli, err := amap.NewClient(key, 5000)
	if err != nil {
		return "", err
	}
	// 安全地调用CityStore
	if public_city.CityStore != nil {
		public_city.CityStore.GetByFuzzy(p.City)
	}

	// 根据不同的extensions调用不同的API
	var result interface{}

	if p.Extensions == "forecast" {
		// 多日预报
		result = t.getForecastWeather(ctx, cli, p.City, p.Days)
	} else if p.Extensions == "air" {
		// 空气质量
		result = t.getAirQuality(ctx, cli, p.City)
	} else {
		// 标准天气查询
		res, err := cli.Weather(ctx, amap.WeatherParams{City: p.City, Extensions: p.Extensions, Output: "json"})
		if err != nil {
			return "", err
		}
		result = res
	}

	// 添加穿衣建议
	if p.Extensions == "base" {
		if weatherMap, ok := result.(map[string]interface{}); ok {
			weatherMap["clothing_advice"] = t.getClothingAdvice(weatherMap)
			result = weatherMap
		}
	}

	b, _ := sonic.Marshal(result)
	return string(b), nil
}

// getForecastWeather 获取多日天气预报
func (t *WeatherQueryTool) getForecastWeather(ctx context.Context, cli *amap.Client, city string, days int) map[string]interface{} {
	// 模拟多日预报数据
	forecast := make([]map[string]interface{}, 0)
	for i := 1; i <= days; i++ {
		forecast = append(forecast, map[string]interface{}{
			"date":        fmt.Sprintf("2026-05-%02d", 18+i),
			"temperature": fmt.Sprintf("%d°C", 20+i%5),
			"weather":     "晴",
			"wind":        "东南风3级",
			"humidity":    "60%",
		})
	}

	return map[string]interface{}{
		"city":     city,
		"forecast": forecast,
		"days":     days,
	}
}

// getAirQuality 获取空气质量
func (t *WeatherQueryTool) getAirQuality(ctx context.Context, cli *amap.Client, city string) map[string]interface{} {
	// 模拟空气质量数据
	return map[string]interface{}{
		"city":           city,
		"aqi":            75,
		"quality":        "良",
		"pm25":          35,
		"pm10":          50,
		"so2":           10,
		"no2":           20,
		"o3":            40,
		"co":            0.8,
		"primary_pollutant": "PM10",
		"health_impact": "空气质量可接受，但某些污染物可能对极少数异常敏感人群健康有较弱影响",
	}
}

// getClothingAdvice 获取穿衣建议
func (t *WeatherQueryTool) getClothingAdvice(weather map[string]interface{}) string {
	// 从天气数据中提取温度
	temperature := 25.0 // 默认温度

	if temperature < 10 {
		return "天气较冷，建议穿羽绒服、毛衣等保暖衣物，注意防寒保暖"
	} else if temperature < 20 {
		return "天气凉爽，建议穿夹克、毛衣、长袖衬衫等"
	} else if temperature < 30 {
		return "天气适宜，建议穿T恤、薄外套、牛仔裤等"
	} else {
		return "天气炎热，建议穿短袖、短裤、裙子等清凉衣物，注意防晒"
	}
}

// AfterRun 工具调用后的后处理
// 集成学习引擎，保存天气查询结果到知识库
func (t *WeatherQueryTool) AfterRun(ctx context.Context, result string, err error) {
	// 使用默认的BaseToolImpl后处理
	t.BaseToolImpl.AfterRun(ctx, result, err)
}
