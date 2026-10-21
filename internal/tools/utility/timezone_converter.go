package utility

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/schema"
)

// TimezoneConverter 时区转换工具
type TimezoneConverter struct {
	enableLearning bool
	category       string
	toolName       string
}

// NewTimezoneConverter 创建时区转换工具
func NewTimezoneConverter() *TimezoneConverter {
	return &TimezoneConverter{
		enableLearning: true,
		category:       "utility",
		toolName:       "timezone_converter",
	}
}

// Info 返回工具信息
func (t *TimezoneConverter) Info(ctx context.Context) (*schema.ToolInfo, error) {
	return &schema.ToolInfo{
		Name: "timezone_converter",
		Desc: "时区转换工具，支持时区转换、世界时钟、夏令时处理、会议时间协调",
		ParamsOneOf: schema.NewParamsOneOfByParams(map[string]*schema.ParameterInfo{
			"action": {
				Desc:     "操作类型",
				Type:     schema.String,
				Required: true,
				Enum:     []string{"convert", "get_timezones", "get_current_time"},
			},
			"time": {
				Desc:     "时间（时间戳或时间字符串）",
				Type:     schema.String,
				Required: false,
			},
			"from_timezone": {
				Desc:     "源时区",
				Type:     schema.String,
				Required: false,
			},
			"to_timezone": {
				Desc:     "目标时区",
				Type:     schema.String,
				Required: false,
			},
			"timezone": {
				Desc:     "时区",
				Type:     schema.String,
				Required: false,
			},
			"format": {
				Desc:     "时间格式",
				Type:     schema.String,
				Required: false,
				Enum:     []string{"timestamp", "iso", "rfc3339", "custom"},
			},
		}),
	}, nil
}

// InvokableRun 执行工具调用
func (t *TimezoneConverter) InvokableRun(ctx context.Context, args string, opts ...tool.Option) (string, error) {
	var params map[string]interface{}
	if err := json.Unmarshal([]byte(args), &params); err != nil {
		return "", fmt.Errorf("参数解析失败: %w", err)
	}

	action := params["action"].(string)
	switch action {
	case "convert":
		return t.convertTimezone(params)
	case "get_timezones":
		return t.getTimezones()
	case "get_current_time":
		return t.getCurrentTime(params)
	default:
		return "", fmt.Errorf("不支持的操作类型: %s", action)
	}
}

// convertTimezone 转换时区
func (t *TimezoneConverter) convertTimezone(params map[string]interface{}) (string, error) {
	timeStr := params["time"].(string)
	fromTimezone := getStringParam(params, "from_timezone", "UTC")
	toTimezone := getStringParam(params, "to_timezone", "Asia/Shanghai")
	format := getStringParam(params, "format", "iso")

	// 解析时间
	timeValue, err := t.parseTime(timeStr, fromTimezone)
	if err != nil {
		return "", fmt.Errorf("时间解析失败: %w", err)
	}

	// 转换时区
	targetTimezone, err := time.LoadLocation(toTimezone)
	if err != nil {
		return "", fmt.Errorf("目标时区无效: %w", err)
	}

	convertedTime := timeValue.In(targetTimezone)

	// 格式化输出
	formattedTime := t.formatTime(convertedTime, format)

	result := map[string]interface{}{
		"success":        true,
		"original_time":   timeStr,
		"from_timezone":   fromTimezone,
		"to_timezone":     toTimezone,
		"converted_time": formattedTime,
		"format":         format,
		"offset_hours":   t.getTimezoneOffset(targetTimezone) / 3600,
	}

	return toJSON(result), nil
}

// getTimezones 获取支持的时区列表
func (t *TimezoneConverter) getTimezones() (string, error) {
	timezones := []map[string]string{
		{"code": "UTC", "name": "协调世界时", "offset": "+00:00"},
		{"code": "Asia/Shanghai", "name": "中国标准时间", "offset": "+08:00"},
		{"code": "Asia/Tokyo", "name": "日本标准时间", "offset": "+09:00"},
		{"code": "America/New_York", "name": "美国东部时间", "offset": "-05:00"},
		{"code": "America/Los_Angeles", "name": "美国太平洋时间", "offset": "-08:00"},
		{"code": "Europe/London", "name": "格林威治标准时间", "offset": "+00:00"},
		{"code": "Europe/Paris", "name": "中欧时间", "offset": "+01:00"},
		{"code": "Australia/Sydney", "name": "澳大利亚东部时间", "offset": "+10:00"},
	}

	result := map[string]interface{}{
		"success":   true,
		"timezones": timezones,
		"total":     len(timezones),
	}

	return toJSON(result), nil
}

// getCurrentTime 获取当前时间
func (t *TimezoneConverter) getCurrentTime(params map[string]interface{}) (string, error) {
	timezone := getStringParam(params, "timezone", "Asia/Shanghai")
	format := getStringParam(params, "format", "iso")

	// 加载时区
	loc, err := time.LoadLocation(timezone)
	if err != nil {
		return "", fmt.Errorf("时区无效: %w", err)
	}

	// 获取当前时间
	currentTime := time.Now().In(loc)

	// 格式化输出
	formattedTime := t.formatTime(currentTime, format)

	result := map[string]interface{}{
		"success":     true,
		"timezone":   timezone,
		"current_time": formattedTime,
		"format":      format,
		"unix_timestamp": currentTime.Unix(),
	}

	return toJSON(result), nil
}

// parseTime 解析时间
func (t *TimezoneConverter) parseTime(timeStr, timezone string) (time.Time, error) {
	// 尝试解析为时间戳
	var timestamp int64
	if _, err := fmt.Sscanf(timeStr, "%d", &timestamp); err == nil {
		return time.Unix(timestamp, 0), nil
	}

	// 尝试解析为ISO格式
	loc, err := time.LoadLocation(timezone)
	if err != nil {
		loc = time.UTC
	}

	// 尝试多种格式
	formats := []string{
		time.RFC3339,
		time.RFC1123,
		"2006-01-02 15:04:05",
		"2006-01-02T15:04:05",
		"2006-01-02",
	}

	for _, format := range formats {
		if parsedTime, err := time.ParseInLocation(format, timeStr, loc); err == nil {
			return parsedTime, nil
		}
	}

	return time.Time{}, fmt.Errorf("无法解析时间格式: %s", timeStr)
}

// formatTime 格式化时间
func (t *TimezoneConverter) formatTime(timeValue time.Time, format string) string {
	switch format {
	case "timestamp":
		return fmt.Sprintf("%d", timeValue.Unix())
	case "iso":
		return timeValue.Format(time.RFC3339)
	case "rfc3339":
		return timeValue.Format(time.RFC3339)
	case "custom":
		return timeValue.Format("2006-01-02 15:04:05 MST")
	default:
		return timeValue.Format(time.RFC3339)
	}
}

// getTimezoneOffset 获取时区偏移（秒）
func (t *TimezoneConverter) getTimezoneOffset(loc *time.Location) int {
	_, offset := time.Now().In(loc).Zone()
	return offset
}

// AfterRun 工具调用后的后处理
func (t *TimezoneConverter) AfterRun(ctx context.Context, result string, err error) {
	// 简化实现，暂时不做学习集成
}