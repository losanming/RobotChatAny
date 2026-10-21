package utility

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/schema"
)

// UnitConverter 单位转换工具
type UnitConverter struct {
	enableLearning bool
	category       string
	toolName       string
}

// NewUnitConverter 创建单位转换工具
func NewUnitConverter() *UnitConverter {
	return &UnitConverter{
		enableLearning: true,
		category:       "utility",
		toolName:       "unit_converter",
	}
}

// Info 返回工具信息
func (u *UnitConverter) Info(ctx context.Context) (*schema.ToolInfo, error) {
	return &schema.ToolInfo{
		Name: "unit_converter",
		Desc: "单位转换工具，支持长度、重量、温度、时间、货币等单位转换",
		ParamsOneOf: schema.NewParamsOneOfByParams(map[string]*schema.ParameterInfo{
			"action": {
				Desc:     "操作类型",
				Type:     schema.String,
				Required: true,
				Enum:     []string{"convert", "get_supported_units"},
			},
			"value": {
				Desc:     "要转换的数值",
				Type:     schema.Number,
				Required: false,
			},
			"from_unit": {
				Desc:     "源单位",
				Type:     schema.String,
				Required: false,
			},
			"to_unit": {
				Desc:     "目标单位",
				Type:     schema.String,
				Required: false,
			},
			"category": {
				Desc:     "单位类别",
				Type:     schema.String,
				Required: false,
				Enum:     []string{"length", "weight", "temperature", "time", "data"},
			},
		}),
	}, nil
}

// InvokableRun 执行工具调用
func (u *UnitConverter) InvokableRun(ctx context.Context, args string, opts ...tool.Option) (string, error) {
	var params map[string]interface{}
	if err := json.Unmarshal([]byte(args), &params); err != nil {
		return "", fmt.Errorf("参数解析失败: %w", err)
	}

	action := params["action"].(string)
	switch action {
	case "convert":
		return u.convert(params)
	case "get_supported_units":
		return u.getSupportedUnits(params)
	default:
		return "", fmt.Errorf("不支持的操作类型: %s", action)
	}
}

// convert 转换单位
func (u *UnitConverter) convert(params map[string]interface{}) (string, error) {
	value := getFloatParam(params, "value", 0)
	fromUnit := getStringParam(params, "from_unit", "")
	toUnit := getStringParam(params, "to_unit", "")

	if fromUnit == "" || toUnit == "" {
		return "", fmt.Errorf("源单位和目标单位不能为空")
	}

	// 执行转换
	convertedValue, err := u.performConversion(value, fromUnit, toUnit)
	if err != nil {
		return "", fmt.Errorf("转换失败: %w", err)
	}

	result := map[string]interface{}{
		"success":        true,
		"original_value": value,
		"from_unit":      fromUnit,
		"to_unit":        toUnit,
		"converted_value": convertedValue,
	}

	return toJSON(result), nil
}

// getSupportedUnits 获取支持的单位列表
func (u *UnitConverter) getSupportedUnits(params map[string]interface{}) (string, error) {
	category := getStringParam(params, "category", "all")

	units := map[string][]map[string]string{
		"length": {
			{"name": "米", "code": "m"},
			{"name": "千米", "code": "km"},
			{"name": "厘米", "code": "cm"},
			{"name": "毫米", "code": "mm"},
			{"name": "英里", "code": "mi"},
			{"name": "英尺", "code": "ft"},
		},
		"weight": {
			{"name": "千克", "code": "kg"},
			{"name": "克", "code": "g"},
			{"name": "毫克", "code": "mg"},
			{"name": "磅", "code": "lb"},
			{"name": "盎司", "code": "oz"},
		},
		"temperature": {
			{"name": "摄氏度", "code": "c"},
			{"name": "华氏度", "code": "f"},
			{"name": "开尔文", "code": "k"},
		},
		"time": {
			{"name": "秒", "code": "s"},
			{"name": "分钟", "code": "min"},
			{"name": "小时", "code": "h"},
			{"name": "天", "code": "d"},
		},
		"data": {
			{"name": "字节", "code": "b"},
			{"name": "千字节", "code": "kb"},
			{"name": "兆字节", "code": "mb"},
			{"name": "吉字节", "code": "gb"},
		},
	}

	var resultUnits []map[string]string
	if category == "all" {
		for cat := range units {
			resultUnits = append(resultUnits, units[cat]...)
		}
	} else {
		resultUnits = units[category]
	}

	result := map[string]interface{}{
		"success": true,
		"units":   resultUnits,
		"total":   len(resultUnits),
	}

	return toJSON(result), nil
}

// performConversion 执行转换
func (u *UnitConverter) performConversion(value float64, fromUnit, toUnit string) (float64, error) {
	// 温度转换特殊处理
	if isTemperatureUnit(fromUnit) && isTemperatureUnit(toUnit) {
		return u.convertTemperature(value, fromUnit, toUnit)
	}

	// 转换为基本单位（米、千克等）
	baseValue, err := u.toBaseUnit(value, fromUnit)
	if err != nil {
		return 0, err
	}

	// 从基本单位转换为目标单位
	result, err := u.fromBaseUnit(baseValue, toUnit)
	if err != nil {
		return 0, err
	}

	return result, nil
}

// toBaseUnit 转换为基本单位
func (u *UnitConverter) toBaseUnit(value float64, unit string) (float64, error) {
	conversions := map[string]float64{
		// 长度（基本单位：米）
		"km": 1000, "cm": 0.01, "mm": 0.001, "mi": 1609.344, "ft": 0.3048,
		// 重量（基本单位：千克）
		"g": 0.001, "mg": 0.000001, "lb": 0.453592, "oz": 0.0283495,
		// 时间（基本单位：秒）
		"min": 60, "h": 3600, "d": 86400,
		// 数据（基本单位：字节）
		"kb": 1024, "mb": 1048576, "gb": 1073741824,
	}

	if factor, ok := conversions[unit]; ok {
		return value * factor, nil
	}

	// 基本单位直接返回
	if unit == "m" || unit == "kg" || unit == "s" || unit == "b" {
		return value, nil
	}

	return 0, fmt.Errorf("不支持的源单位: %s", unit)
}

// fromBaseUnit 从基本单位转换
func (u *UnitConverter) fromBaseUnit(value float64, unit string) (float64, error) {
	conversions := map[string]float64{
		// 长度（基本单位：米）
		"km": 1 / 1000, "cm": 1 / 0.01, "mm": 1 / 0.001, "mi": 1 / 1609.344, "ft": 1 / 0.3048,
		// 重量（基本单位：千克）
		"g": 1 / 0.001, "mg": 1 / 0.000001, "lb": 1 / 0.453592, "oz": 1 / 0.0283495,
		// 时间（基本单位：秒）
		"min": 1 / 60, "h": 1 / 3600, "d": 1 / 86400,
		// 数据（基本单位：字节）
		"kb": 1 / 1024, "mb": 1 / 1048576, "gb": 1 / 1073741824,
	}

	if factor, ok := conversions[unit]; ok {
		return value * factor, nil
	}

	// 基本单位直接返回
	if unit == "m" || unit == "kg" || unit == "s" || unit == "b" {
		return value, nil
	}

	return 0, fmt.Errorf("不支持的目标单位: %s", unit)
}

// convertTemperature 温度转换
func (u *UnitConverter) convertTemperature(value float64, fromUnit, toUnit string) (float64, error) {
	// 先转换为摄氏度
	var celsius float64
	switch fromUnit {
	case "c":
		celsius = value
	case "f":
		celsius = (value - 32) * 5 / 9
	case "k":
		celsius = value - 273.15
	default:
		return 0, fmt.Errorf("不支持的温度单位: %s", fromUnit)
	}

	// 从摄氏度转换为目标单位
	switch toUnit {
	case "c":
		return celsius, nil
	case "f":
		return celsius*9/5 + 32, nil
	case "k":
		return celsius + 273.15, nil
	default:
		return 0, fmt.Errorf("不支持的温度单位: %s", toUnit)
	}
}

// isTemperatureUnit 检查是否为温度单位
func isTemperatureUnit(unit string) bool {
	return unit == "c" || unit == "f" || unit == "k"
}

// AfterRun 工具调用后的后处理
func (u *UnitConverter) AfterRun(ctx context.Context, result string, err error) {
	// 简化实现，暂时不做学习集成
}