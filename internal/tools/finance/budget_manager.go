package finance

import (
	"context"
	"encoding/json"

	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/schema"
)

// BudgetManager 预算管理工具
type BudgetManager struct {
	enableLearning bool
	category       string
	toolName       string
}

func NewBudgetManager() *BudgetManager {
	return &BudgetManager{
		enableLearning: true,
		category:       "finance",
		toolName:       "budget_manager",
	}
}

func (b *BudgetManager) Info(ctx context.Context) (*schema.ToolInfo, error) {
	return &schema.ToolInfo{
		Name: "budget_manager",
		Desc: "预算管理工具，支持预算设置、支出跟踪、预算分析",
		ParamsOneOf: schema.NewParamsOneOfByParams(map[string]*schema.ParameterInfo{
			"action": {Desc: "操作类型", Type: schema.String, Required: true, Enum: []string{"set_budget", "track_spending", "analyze"}},
		}),
	}, nil
}

func (b *BudgetManager) InvokableRun(ctx context.Context, args string, opts ...tool.Option) (string, error) {
	var params map[string]interface{}
	json.Unmarshal([]byte(args), &params)
	result := map[string]interface{}{
		"success": true,
		"message": "预算操作成功（演示版本）",
	}
	return toJSON(result), nil
}

func (b *BudgetManager) AfterRun(ctx context.Context, result string, err error) {}