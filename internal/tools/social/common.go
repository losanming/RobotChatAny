package social

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