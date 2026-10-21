package finance

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/schema"
)

// FinanceManager 财务管理工具
type FinanceManager struct {
	enableLearning bool
	category       string
	toolName       string
}

func NewFinanceManager() *FinanceManager {
	return &FinanceManager{
		enableLearning: true,
		category:       "finance",
		toolName:       "finance_manager",
	}
}

func (f *FinanceManager) Info(ctx context.Context) (*schema.ToolInfo, error) {
	return &schema.ToolInfo{
		Name: "finance_manager",
		Desc: "财务管理工具，支持收支记录、财务分析、账单管理",
		ParamsOneOf: schema.NewParamsOneOfByParams(map[string]*schema.ParameterInfo{
			"action": {Desc: "操作类型", Type: schema.String, Required: true, Enum: []string{"record", "analyze", "get_summary"}},
			"amount": {Desc: "金额", Type: schema.Number, Required: false},
			"category": {Desc: "分类", Type: schema.String, Required: false},
		}),
	}, nil
}

func (f *FinanceManager) InvokableRun(ctx context.Context, args string, opts ...tool.Option) (string, error) {
	var params map[string]interface{}
	json.Unmarshal([]byte(args), &params)
	action := params["action"].(string)
	
	result := map[string]interface{}{
		"success": true,
		"action": action,
		"message": fmt.Sprintf("财务操作成功（演示版本）"),
		"processed_at": time.Now().Format(time.RFC3339),
	}
	return toJSON(result), nil
}

func (f *FinanceManager) AfterRun(ctx context.Context, result string, err error) {}