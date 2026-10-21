package communication

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/schema"
)

// EmailManager 邮件管理工具
type EmailManager struct {
	enableLearning bool
	category       string
	toolName       string
}

// NewEmailManager 创建邮件管理工具
func NewEmailManager() *EmailManager {
	return &EmailManager{
		enableLearning: true,
		category:       "communication",
		toolName:       "email_manager",
	}
}

// Info 返回工具信息
func (e *EmailManager) Info(ctx context.Context) (*schema.ToolInfo, error) {
	return &schema.ToolInfo{
		Name: "email_manager",
		Desc: "邮件管理工具，支持邮件发送、接收、分类、附件处理（演示版本）",
		ParamsOneOf: schema.NewParamsOneOfByParams(map[string]*schema.ParameterInfo{
			"action": {
				Desc:     "操作类型",
				Type:     schema.String,
				Required: true,
				Enum:     []string{"send", "list", "read", "delete", "classify"},
			},
			"to": {
				Desc:     "收件人邮箱",
				Type:     schema.String,
				Required: false,
			},
			"subject": {
				Desc:     "邮件主题",
				Type:     schema.String,
				Required: false,
			},
			"body": {
				Desc:     "邮件正文",
				Type:     schema.String,
				Required: false,
			},
			"from": {
				Desc:     "发件人邮箱",
				Type:     schema.String,
				Required: false,
			},
			"smtp_host": {
				Desc:     "SMTP服务器地址",
				Type:     schema.String,
				Required: false,
			},
			"smtp_port": {
				Desc:     "SMTP端口",
				Type:     schema.Integer,
				Required: false,
			},
			"username": {
				Desc:     "SMTP用户名",
				Type:     schema.String,
				Required: false,
			},
			"password": {
				Desc:     "SMTP密码",
				Type:     schema.String,
				Required: false,
			},
		}),
	}, nil
}

// InvokableRun 执行工具调用
func (e *EmailManager) InvokableRun(ctx context.Context, args string, opts ...tool.Option) (string, error) {
	var params map[string]interface{}
	if err := json.Unmarshal([]byte(args), &params); err != nil {
		return "", fmt.Errorf("参数解析失败: %w", err)
	}

	action := params["action"].(string)
	switch action {
	case "send":
		return e.sendEmail(params)
	case "list":
		return e.listEmails(params)
	case "read":
		return e.readEmail(params)
	case "delete":
		return e.deleteEmail(params)
	case "classify":
		return e.classifyEmail(params)
	default:
		return "", fmt.Errorf("不支持的操作类型: %s", action)
	}
}

// sendEmail 发送邮件
func (e *EmailManager) sendEmail(params map[string]interface{}) (string, error) {
	to := getStringParam(params, "to", "")
	subject := getStringParam(params, "subject", "")
	from := getStringParam(params, "from", "")
	smtpHost := getStringParam(params, "smtp_host", "smtp.gmail.com")
	smtpPort := getIntParam(params, "smtp_port", 587)

	if to == "" {
		return "", fmt.Errorf("收件人邮箱不能为空")
	}

	// 演示版本：模拟发送邮件
	emailID := generateEmailID()
	
	// 实际发送逻辑（注释掉，需要真实的SMTP配置）
	/*
	auth := smtp.PlainAuth("", username, password, smtpHost)
	msg := fmt.Sprintf("From: %s\r\nTo: %s\r\nSubject: %s\r\n\r\n%s", from, to, subject, body)
	err := smtp.SendMail(fmt.Sprintf("%s:%d", smtpHost, smtpPort), auth, from, []string{to}, []byte(msg))
	if err != nil {
		return "", fmt.Errorf("邮件发送失败: %w", err)
	}
	*/

	result := map[string]interface{}{
		"success":     true,
		"email_id":    emailID,
		"to":          to,
		"from":        from,
		"subject":     subject,
		"sent_at":     time.Now().Format(time.RFC3339),
		"message":     "邮件发送成功（演示版本）",
		"smtp_config": map[string]interface{}{
			"host": smtpHost,
			"port": smtpPort,
		},
	}

	return toJSON(result), nil
}

// listEmails 列出邮件
func (e *EmailManager) listEmails(params map[string]interface{}) (string, error) {
	// 演示版本：返回模拟邮件列表
	emails := []map[string]interface{}{
		{
			"id":       "email_001",
			"from":     "sender@example.com",
			"to":       "user@example.com",
			"subject":  "项目进度更新",
			"body":     "本周项目进度良好...",
			"sent_at":  time.Now().Add(-2 * time.Hour).Format(time.RFC3339),
			"read":     false,
			"category": "工作",
		},
		{
			"id":       "email_002",
			"from":     "newsletter@tech.com",
			"to":       "user@example.com",
			"subject":  "技术周刊",
			"body":     "本周技术热点...",
			"sent_at":  time.Now().Add(-24 * time.Hour).Format(time.RFC3339),
			"read":     true,
			"category": "订阅",
		},
	}

	result := map[string]interface{}{
		"success": true,
		"emails":  emails,
		"total":   len(emails),
	}

	return toJSON(result), nil
}

// readEmail 读取邮件
func (e *EmailManager) readEmail(params map[string]interface{}) (string, error) {
	emailID := params["id"].(string)
	if emailID == "" {
		return "", fmt.Errorf("邮件ID不能为空")
	}

	// 演示版本：返回模拟邮件详情
	email := map[string]interface{}{
		"id":       emailID,
		"from":     "sender@example.com",
		"to":       "user@example.com",
		"subject":  "项目进度更新",
		"body":     "本周项目进度良好，主要完成了以下工作：\n1. 完成工具系统重构\n2. 实现P0阶段工具\n3. 开始P1阶段开发",
		"sent_at":  time.Now().Add(-2 * time.Hour).Format(time.RFC3339),
		"read_at":  time.Now().Format(time.RFC3339),
		"attachments": []string{"report.pdf"},
		"category": "工作",
	}

	result := map[string]interface{}{
		"success": true,
		"email":   email,
	}

	return toJSON(result), nil
}

// deleteEmail 删除邮件
func (e *EmailManager) deleteEmail(params map[string]interface{}) (string, error) {
	emailID := params["id"].(string)
	if emailID == "" {
		return "", fmt.Errorf("邮件ID不能为空")
	}

	result := map[string]interface{}{
		"success": true,
		"email_id": emailID,
		"message":  "邮件删除成功（演示版本）",
	}

	return toJSON(result), nil
}

// classifyEmail 分类邮件
func (e *EmailManager) classifyEmail(params map[string]interface{}) (string, error) {
	subject := params["subject"].(string)
	body := params["body"].(string)

	// 简单的分类逻辑
	content := subject + " " + body
	content = strings.ToLower(content)

	category := "其他"
	if strings.Contains(content, "项目") || strings.Contains(content, "工作") || strings.Contains(content, "会议") {
		category = "工作"
	} else if strings.Contains(content, "订阅") || strings.Contains(content, "新闻") || strings.Contains(content, "资讯") {
		category = "订阅"
	} else if strings.Contains(content, "账单") || strings.Contains(content, "发票") || strings.Contains(content, "财务") {
		category = "财务"
	} else if strings.Contains(content, "健康") || strings.Contains(content, "运动") || strings.Contains(content, "医疗") {
		category = "健康"
	}

	result := map[string]interface{}{
		"success":  true,
		"category": category,
		"confidence": 0.85,
	}

	return toJSON(result), nil
}

// AfterRun 工具调用后的后处理
func (e *EmailManager) AfterRun(ctx context.Context, result string, err error) {
	// 简化实现，暂时不做学习集成
}

// 辅助函数
func generateEmailID() string {
	return fmt.Sprintf("email_%d", time.Now().UnixNano())
}