package personal

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

// getStringArrayParam 获取字符串数组参数
func getStringArrayParam(params map[string]interface{}, key string) []string {
	if val, ok := params[key]; ok {
		if arr, ok := val.([]interface{}); ok {
			result := make([]string, len(arr))
			for i, v := range arr {
				result[i] = v.(string)
			}
			return result
		}
	}
	return []string{}
}

// contains 检查字符串是否包含关键词
func contains(text, keyword string) bool {
	if len(text) == 0 || len(keyword) == 0 {
		return false
	}
	
	if text == keyword {
		return true
	}
	
	if len(text) < len(keyword) {
		return false
	}
	
	// 检查前缀
	if text[:len(keyword)] == keyword {
		return true
	}
	
	// 检查后缀
	if text[len(text)-len(keyword):] == keyword {
		return true
	}
	
	// 检查中间
	return indexOf(text, keyword) >= 0
}

// indexOf 查找子字符串位置
func indexOf(text, keyword string) int {
	for i := 0; i <= len(text)-len(keyword); i++ {
		if text[i:i+len(keyword)] == keyword {
			return i
		}
	}
	return -1
}

// toJSON 转换为JSON字符串
func toJSON(data interface{}) string {
	bytes, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Sprintf("{\"error\": \"JSON序列化失败: %v\"}", err)
	}
	return string(bytes)
}