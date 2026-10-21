package mobile

import (
	"context"
	"encoding/json"

	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/schema"
)

// MobileAdapter 移动端适配工具
type MobileAdapter struct {
	enableLearning bool
	category       string
	toolName       string
}

func NewMobileAdapter() *MobileAdapter {
	return &MobileAdapter{
		enableLearning: true,
		category:       "mobile",
		toolName:       "mobile_adapter",
	}
}

func (m *MobileAdapter) Info(ctx context.Context) (*schema.ToolInfo, error) {
	return &schema.ToolInfo{
		Name: "mobile_adapter",
		Desc: "移动端适配工具，支持移动端界面适配、推送通知、位置跟踪",
		ParamsOneOf: schema.NewParamsOneOfByParams(map[string]*schema.ParameterInfo{
			"action": {Desc: "操作类型", Type: schema.String, Required: true, Enum: []string{"adapt", "push", "track_location"}},
		}),
	}, nil
}

func (m *MobileAdapter) InvokableRun(ctx context.Context, args string, opts ...tool.Option) (string, error) {
	var params map[string]interface{}
	json.Unmarshal([]byte(args), &params)
	result := map[string]interface{}{
		"success": true,
		"message": "移动端操作成功（演示版本）",
	}
	return toJSON(result), nil
}

func (m *MobileAdapter) AfterRun(ctx context.Context, result string, err error) {}