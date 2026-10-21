package multimedia

import (
	"encoding/json"
	"fmt"
)

func toJSON(data interface{}) string {
	bytes, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Sprintf("{\"error\": \"JSON序列化失败: %v\"}", err)
	}
	return string(bytes)
}

func getStringParam(params map[string]interface{}, key, defaultValue string) string {
	if val, ok := params[key]; ok {
		return val.(string)
	}
	return defaultValue
}

func getIntParam(params map[string]interface{}, key string, defaultValue int) int {
	if val, ok := params[key]; ok {
		return int(val.(float64))
	}
	return defaultValue
}