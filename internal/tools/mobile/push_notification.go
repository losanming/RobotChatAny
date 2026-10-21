package mobile

import (
	"context"
	"encoding/json"

	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/schema"
)

// PushNotification 推送通知工具
type PushNotification struct {
	enableLearning bool
	category       string
	toolName       string
}

func NewPushNotification() *PushNotification {
	return &PushNotification{
		enableLearning: true,
		category:       "mobile",
		toolName:       "push_notification",
	}
}

func (p *PushNotification) Info(ctx context.Context) (*schema.ToolInfo, error) {
	return &schema.ToolInfo{
		Name: "push_notification",
		Desc: "推送通知工具，支持移动端推送、跨平台通知、通知管理",
		ParamsOneOf: schema.NewParamsOneOfByParams(map[string]*schema.ParameterInfo{
			"action": {Desc: "操作类型", Type: schema.String, Required: true, Enum: []string{"push", "list", "cancel"}},
		}),
	}, nil
}

func (p *PushNotification) InvokableRun(ctx context.Context, args string, opts ...tool.Option) (string, error) {
	var params map[string]interface{}
	json.Unmarshal([]byte(args), &params)
	result := map[string]interface{}{
		"success": true,
		"message": "推送通知成功（演示版本）",
	}
	return toJSON(result), nil
}

func (p *PushNotification) AfterRun(ctx context.Context, result string, err error) {}