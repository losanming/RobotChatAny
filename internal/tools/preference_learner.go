package tools

import (
	"context"
	"fmt"
	"time"

	"github.com/bytedance/sonic"
	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/schema"
)

// PreferenceLearnerTool 偏好学习工具
// 分析用户行为，提取兴趣偏好，个性化推荐
type PreferenceLearnerTool struct {
	BaseToolImpl
	// 这里可以添加对brain.PreferenceLearner的引用
	// 但为了避免循环导入，我们使用简单的本地实现
}

type preferenceParams struct {
	Action     string                 `json:"action"`     // analyze, update, get
	Category   string                 `json:"category,omitempty"`   // 分类
	UserAction string                 `json:"user_action,omitempty"` // 用户行为: view, like, share, ignore
	Duration   int                    `json:"duration,omitempty"`   // 持续时间(秒)
	Metadata   map[string]interface{} `json:"metadata,omitempty"`   // 扩展元数据
}

type preferenceResult struct {
	Action       string                 `json:"action"`
	Category     string                 `json:"category,omitempty"`
	Preferences  map[string]float64     `json:"preferences,omitempty"`
	Interests    []string               `json:"interests,omitempty"`
	LastUpdated  time.Time              `json:"last_updated"`
	Message      string                 `json:"message"`
	Timestamp    time.Time              `json:"timestamp"`
}

// 简单的本地偏好存储
var localPreferences = map[string]float64{
	"tech":     0.8,
	"health":   0.9,
	"life":     0.6,
	"finance":  0.4,
	"location": 0.7,
}

var localInterests = []string{"人工智能", "健康管理", "编程技术", "运动健身"}

func NewPreferenceLearnerTool() *PreferenceLearnerTool {
	return &PreferenceLearnerTool{
		BaseToolImpl: BaseToolImpl{
			Category:       "learning",
			EnableLearning: true,
			ToolName:       "preference_learner",
		},
	}
}

func (t *PreferenceLearnerTool) Info(_ context.Context) (*schema.ToolInfo, error) {
	return &schema.ToolInfo{
		Name: "preference_learner",
		Desc: "分析用户行为，提取兴趣偏好，个性化推荐。支持分析偏好、更新偏好、获取偏好",
		ParamsOneOf: schema.NewParamsOneOfByParams(map[string]*schema.ParameterInfo{
			"action":      {Desc: "操作类型: analyze(分析), update(更新), get(获取)", Type: schema.String, Required: true},
			"category":    {Desc: "分类 (update时需要)", Type: schema.String, Required: false},
			"user_action": {Desc: "用户行为: view, like, share, ignore (update时需要)", Type: schema.String, Required: false},
			"duration":    {Desc: "持续时间(秒) (update时需要)", Type: schema.Integer, Required: false},
			"metadata":    {Desc: "扩展元数据", Type: schema.Object, Required: false},
		}),
	}, nil
}

func (t *PreferenceLearnerTool) InvokableRun(ctx context.Context, argumentsInJSON string, _ ...tool.Option) (string, error) {
	p := preferenceParams{}
	if err := sonic.UnmarshalString(argumentsInJSON, &p); err != nil {
		return "", err
	}

	// 验证参数
	if p.Action == "" {
		return "", fmt.Errorf("action is required")
	}

	var result preferenceResult

	switch p.Action {
	case "analyze":
		result = t.analyzePreferences()
	case "update":
		if p.Category == "" {
			return "", fmt.Errorf("category is required for update action")
		}
		result = t.updatePreferences(p.Category, p.UserAction, p.Duration)
	case "get":
		result = t.getPreferences()
	default:
		return "", fmt.Errorf("unsupported action: %s", p.Action)
	}

	result.Timestamp = time.Now()

	b, _ := sonic.Marshal(result)
	return string(b), nil
}

// analyzePreferences 分析用户偏好
func (t *PreferenceLearnerTool) analyzePreferences() preferenceResult {
	// 分析本地偏好数据
	totalWeight := 0.0
	for _, weight := range localPreferences {
		totalWeight += weight
	}

	// 归一化
	normalizedPreferences := make(map[string]float64)
	for category, weight := range localPreferences {
		if totalWeight > 0 {
			normalizedPreferences[category] = weight / totalWeight
		}
	}

	return preferenceResult{
		Action:      "analyze",
		Preferences: normalizedPreferences,
		Interests:   localInterests,
		LastUpdated: time.Now(),
		Message:     "偏好分析完成",
	}
}

// updatePreferences 更新用户偏好
func (t *PreferenceLearnerTool) updatePreferences(category, userAction string, duration int) preferenceResult {
	// 根据用户行为调整权重
	weightChange := 0.0

	switch userAction {
	case "like":
		weightChange = 0.1
	case "share":
		weightChange = 0.15
	case "view":
		weightChange = 0.05
	case "ignore":
		weightChange = -0.1
	default:
		weightChange = 0.02 // 默认轻微增加
	}

	// 考虑持续时间的影响
	if duration > 60 { // 超过1分钟
		weightChange *= 1.5
	}

	// 更新权重
	if currentWeight, exists := localPreferences[category]; exists {
		localPreferences[category] = currentWeight + weightChange
		// 确保权重在合理范围内
		if localPreferences[category] < 0 {
			localPreferences[category] = 0
		}
		if localPreferences[category] > 2.0 {
			localPreferences[category] = 2.0
		}
	} else {
		localPreferences[category] = 0.5 + weightChange
	}

	return preferenceResult{
		Action:      "update",
		Category:    category,
		Preferences: localPreferences,
		Interests:   localInterests,
		LastUpdated: time.Now(),
		Message:     fmt.Sprintf("已更新 %s 分类偏好，权重变化: %.2f", category, weightChange),
	}
}

// getPreferences 获取用户偏好
func (t *PreferenceLearnerTool) getPreferences() preferenceResult {
	// 找出权重最高的分类
	maxWeight := 0.0
	topCategory := ""
	for category, weight := range localPreferences {
		if weight > maxWeight {
			maxWeight = weight
			topCategory = category
		}
	}

	// 生成推荐
	recommendations := t.generateRecommendations(topCategory)

	return preferenceResult{
		Action:      "get",
		Category:    topCategory,
		Preferences: localPreferences,
		Interests:   localInterests,
		LastUpdated: time.Now(),
		Message:     fmt.Sprintf("当前最感兴趣的分类: %s (权重: %.2f)。推荐: %s", topCategory, maxWeight, recommendations),
	}
}

// generateRecommendations 生成个性化推荐
func (t *PreferenceLearnerTool) generateRecommendations(topCategory string) string {
	recommendations := map[string]string{
		"tech":    "人工智能最新进展、编程技术教程、科技产品评测",
		"health":  "健康管理建议、运动健身计划、营养饮食指导",
		"life":    "生活技巧、美食推荐、旅游攻略",
		"finance": "投资理财建议、财经新闻分析、保险规划",
		"location": "周边美食推荐、旅游景点、本地活动",
	}

	if rec, exists := recommendations[topCategory]; exists {
		return rec
	}
	return "综合内容推荐"
}

// AfterRun 工具调用后的后处理
// 偏好学习工具本身不需要额外的学习集成
func (t *PreferenceLearnerTool) AfterRun(ctx context.Context, result string, err error) {
	// 偏好学习工具本身就是学习工具，不需要额外的学习集成
}
