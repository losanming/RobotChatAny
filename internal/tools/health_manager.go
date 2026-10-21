package tools

import (
	"context"
	"errors"

	"github.com/bytedance/sonic"
	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/schema"

	"new-robot-self/internal/brain"
)

// HealthManagerTool 让 AI 设置或补丁更新健康信息（秒级时间戳）
type HealthManagerTool struct {
	BaseToolImpl
	brainService *brain.Brain
}

func NewHealthManagerTool() *HealthManagerTool {
	return &HealthManagerTool{
		BaseToolImpl: BaseToolImpl{
			Category:       "health",
			EnableLearning: true,
			ToolName:       "health_manager",
		},
		brainService: nil, // 将在InitTools中设置
	}
}

// SetBrainService 设置Brain服务实例
func (t *HealthManagerTool) SetBrainService(brainService *brain.Brain) {
	t.brainService = brainService
}

type healthParams struct {
	Action       string   `json:"action"` // set | patch
	WeightKg     *float64 `json:"weight_kg,omitempty"`
	BMI          *float64 `json:"bmi,omitempty"`
	ExecDesc     *string  `json:"exec_desc,omitempty"` // 运动计划备注
	BloodAvg     *float64 `json:"blood_recent_avg,omitempty"`
	BloodRisk    *bool    `json:"blood_risk_flag,omitempty"`
}

func (t *HealthManagerTool) Info(_ context.Context) (*schema.ToolInfo, error) {
	return &schema.ToolInfo{
		Name: "health_manager",
		Desc: "设置或补丁更新健康信息，时间戳单位为秒。action: set|patch",
		ParamsOneOf: schema.NewParamsOneOfByParams(map[string]*schema.ParameterInfo{
			"action":           {Desc: "set 或 patch", Type: schema.String, Required: true},
			"weight_kg":        {Desc: "体重(kg)", Type: schema.Number, Required: false},
			"bmi":              {Desc: "BMI", Type: schema.Number, Required: false},
			"exec_desc":        {Desc: "运动计划备注", Type: schema.String, Required: false},
			"blood_recent_avg": {Desc: "血糖近期均值", Type: schema.Number, Required: false},
			"blood_risk_flag":  {Desc: "血糖风险标记", Type: schema.Boolean, Required: false},
		}),
	}, nil
}

func (t *HealthManagerTool) InvokableRun(ctx context.Context, args string, _ ...tool.Option) (string, error) {
	if t.brainService == nil {
		return "", errors.New("health_manager: brain service not initialized")
	}
	var p healthParams
	if err := sonic.UnmarshalString(args, &p); err != nil {
		return "", err
	}
	switch p.Action {
	case "set":
		h := brain.HealthState{}
		if p.WeightKg != nil {
			h.WeightKg = *p.WeightKg
		}
		if p.BMI != nil {
			h.BMI = *p.BMI
		}
		if p.ExecDesc != nil {
			h.ExecDesc = *p.ExecDesc
		}
		if p.BloodAvg != nil || p.BloodRisk != nil {
			bs := &brain.BloodSugarBrief{}
			if p.BloodAvg != nil {
				bs.RecentAvg = *p.BloodAvg
			}
			if p.BloodRisk != nil {
				bs.RiskFlag = *p.BloodRisk
			}
			h.BloodSugar = bs
		}
		t.brainService.SetHealth(h)
		return `{"ok":true,"action":"set"}`, nil
	case "patch":
		t.brainService.PatchHealth(func(h *brain.HealthState) {
			if p.WeightKg != nil {
				h.WeightKg = *p.WeightKg
			}
			if p.BMI != nil {
				h.BMI = *p.BMI
			}
			if p.ExecDesc != nil {
				h.ExecDesc = *p.ExecDesc
			}
			if p.BloodAvg != nil || p.BloodRisk != nil {
				if h.BloodSugar == nil {
					h.BloodSugar = &brain.BloodSugarBrief{}
				}
				if p.BloodAvg != nil {
					h.BloodSugar.RecentAvg = *p.BloodAvg
				}
				if p.BloodRisk != nil {
					h.BloodSugar.RiskFlag = *p.BloodRisk
				}
			}
		})
		return `{"ok":true,"action":"patch"}`, nil
	default:
		return "", errors.New("health_manager: unsupported action")
	}
}

// AfterRun 工具调用后的后处理
// 集成学习引擎，保存健康信息变更到知识库
func (t *HealthManagerTool) AfterRun(ctx context.Context, result string, err error) {
	// 使用默认的BaseToolImpl后处理
	t.BaseToolImpl.AfterRun(ctx, result, err)
}
