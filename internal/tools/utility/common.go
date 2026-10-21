package utility

import (
	"encoding/json"
	"fmt"
)

// getStringParam 获取字符串参数
func getStringParam(params map[string]interface{}, key, defaultValue string) string {
	if val, ok := params[key]; ok {
		return val.(string)
	}
	return defaultValue
}

// getIntParam 获取整数参数
func getIntParam(params map[string]interface{}, key string, defaultValue int) int {
	if val, ok := params[key]; ok {
		return int(val.(float64))
	}
	return defaultValue
}

// getInt64Param 获取64位整数参数
func getInt64Param(params map[string]interface{}, key string, defaultValue int64) int64 {
	if val, ok := params[key]; ok {
		return int64(val.(float64))
	}
	return defaultValue
}

// getBoolParam 获取布尔参数
func getBoolParam(params map[string]interface{}, key string, defaultValue bool) bool {
	if val, ok := params[key]; ok {
		return val.(bool)
	}
	return defaultValue
}

// getFloatParam 获取浮点数参数
func getFloatParam(params map[string]interface{}, key string, defaultValue float64) float64 {
	if val, ok := params[key]; ok {
		return val.(float64)
	}
	return defaultValue
}

// toJSON 转换为JSON字符串
func toJSON(data interface{}) string {
	bytes, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Sprintf("{\"error\": \"JSON序列化失败: %v\"}", err)
	}
	return string(bytes)
}