package communication

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/schema"
)

// NotificationManager 通知管理工具
type NotificationManager struct {
	enableLearning bool
	category       string
	toolName       string
}

// NewNotificationManager 创建通知管理工具
func NewNotificationManager() *NotificationManager {
	return &NotificationManager{
		enableLearning: true,
		category:       "communication",
		toolName:       "notification_manager",
	}
}

// Info 返回工具信息
func (n *NotificationManager) Info(ctx context.Context) (*schema.ToolInfo, error) {
	return &schema.ToolInfo{
		Name: "notification_manager",
		Desc: "通知管理工具，支持通知发送、定时提醒、通知历史",
		ParamsOneOf: schema.NewParamsOneOfByParams(map[string]*schema.ParameterInfo{
			"action": {
				Desc:     "操作类型",
				Type:     schema.String,
				Required: true,
				Enum:     []string{"send", "schedule", "list", "delete"},
			},
			"title": {
				Desc:     "通知标题",
				Type:     schema.String,
				Required: false,
			},
			"message": {
				Desc:     "通知内容",
				Type:     schema.String,
				Required: false,
			},
			"type": {
				Desc:     "通知类型",
				Type:     schema.String,
				Required: false,
				Enum:     []string{"info", "warning", "error", "success"},
			},
			"scheduled_time": {
				Desc:     "预定时间",
				Type:     schema.String,
				Required: false,
			},
		}),
	}, nil
}

// InvokableRun 执行工具调用
func (n *NotificationManager) InvokableRun(ctx context.Context, args string, opts ...tool.Option) (string, error) {
	var params map[string]interface{}
	if err := json.Unmarshal([]byte(args), &params); err != nil {
		return "", fmt.Errorf("参数解析失败: %w", err)
	}

	action := params["action"].(string)
	switch action {
	case "send":
		return n.sendNotification(params)
	case "schedule":
		return n.scheduleNotification(params)
	case "list":
		return n.listNotifications(params)
	case "delete":
		return n.deleteNotification(params)
	default:
		return "", fmt.Errorf("不支持的操作类型: %s", action)
	}
}

// sendNotification 发送通知
func (n *NotificationManager) sendNotification(params map[string]interface{}) (string, error) {
	title := getStringParam(params, "title", "")
	message := params["message"].(string)
	notifyType := getStringParam(params, "type", "info")

	notificationID := fmt.Sprintf("notif_%d", time.Now().UnixNano())

	result := map[string]interface{}{
		"success":        true,
		"notification_id": notificationID,
		"title":          title,
		"message":        message,
		"type":           notifyType,
		"sent_at":        time.Now().Format(time.RFC3339),
		"status":         "通知发送成功（演示版本）",
	}

	return toJSON(result), nil
}

// scheduleNotification 预定通知
func (n *NotificationManager) scheduleNotification(params map[string]interface{}) (string, error) {
	title := getStringParam(params, "title", "")
	message := params["message"].(string)
	scheduledTime := params["scheduled_time"].(string)

	result := map[string]interface{}{
		"success":        true,
		"notification_id": fmt.Sprintf("notif_%d", time.Now().UnixNano()),
		"title":          title,
		"message":        message,
		"scheduled_time": scheduledTime,
		"status":         "通知预定成功（演示版本）",
	}

	return toJSON(result), nil
}

// listNotifications 列出通知
func (n *NotificationManager) listNotifications(params map[string]interface{}) (string, error) {
	notifications := []map[string]interface{}{
		{
			"id":           "notif_001",
			"title":         "运动提醒",
			"message":       "该进行室内慢跑了",
			"type":          "info",
			"sent_at":       time.Now().Add(-1 * time.Hour).Format(time.RFC3339),
			"read":          false,
		},
	}

	result := map[string]interface{}{
		"success":       true,
		"notifications": notifications,
		"total":         len(notifications),
	}

	return toJSON(result), nil
}

// deleteNotification 删除通知
func (n *NotificationManager) deleteNotification(params map[string]interface{}) (string, error) {
	notificationID := params["id"].(string)

	result := map[string]interface{}{
		"success":         true,
		"notification_id": notificationID,
		"message":         "通知删除成功（演示版本）",
	}

	return toJSON(result), nil
}

// AfterRun 工具调用后的后处理
func (n *NotificationManager) AfterRun(ctx context.Context, result string, err error) {
	// 简化实现
}