package tools

import (
	"context"
	"fmt"
	"time"

	"github.com/bytedance/sonic"
	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/schema"
)

// ExerciseAdvisorTool 运动建议工具
// 基于健康数据提供个性化运动建议、运动计划制定、效果跟踪
type ExerciseAdvisorTool struct {
	BaseToolImpl
}

type exerciseParams struct {
	Action      string  `json:"action"`       // plan, advice, track
	Weight      float64 `json:"weight,omitempty"`      // 体重(kg)
	Height      float64 `json:"height,omitempty"`      // 身高(cm)
	BMI         float64 `json:"bmi,omitempty"`         // BMI
	ExerciseType string  `json:"exercise_type,omitempty"` // 运动类型
	Duration    int     `json:"duration,omitempty"`    // 运动时长(分钟)
	Goal        string  `json:"goal,omitempty"`        // 目标: weight_loss, fitness, health
}

type exerciseResult struct {
	Action       string           `json:"action"`
	Plan         *exercisePlan    `json:"plan,omitempty"`
	Advice       string           `json:"advice,omitempty"`
	Tracking     *exerciseTracking `json:"tracking,omitempty"`
	Recommendations []string      `json:"recommendations,omitempty"`
	CaloriesBurned int            `json:"calories_burned,omitempty"`
	GeneratedAt  time.Time        `json:"generated_at"`
	Message      string           `json:"message"`
}

type exercisePlan struct {
	WeeklySchedule []exerciseSession `json:"weekly_schedule"`
	TotalDuration  int              `json:"total_duration"` // 每周总时长(分钟)
	EstimatedCalories int           `json:"estimated_calories"` // 每周消耗卡路里
}

type exerciseSession struct {
	Day         string `json:"day"`         // 周几
	Type        string `json:"type"`        // 运动类型
	Duration    int    `json:"duration"`    // 时长(分钟)
	Intensity   string `json:"intensity"`   // 强度: low, medium, high
	Description string `json:"description"` // 描述
}

type exerciseTracking struct {
	CurrentWeight     float64 `json:"current_weight"`
	TargetWeight      float64 `json:"target_weight"`
	WeightProgress    float64 `json:"weight_progress"` // 减重进度(kg)
	ExerciseDays      int     `json:"exercise_days"`    // 运动天数
	TotalDuration     int     `json:"total_duration"`    // 总运动时长(分钟)
	ConsistencyScore  float64 `json:"consistency_score"` // 坚持度评分(0-100)
	NextMilestone     string  `json:"next_milestone"`    // 下一个里程碑
}

func NewExerciseAdvisorTool() *ExerciseAdvisorTool {
	return &ExerciseAdvisorTool{
		BaseToolImpl: BaseToolImpl{
			Category:       "health",
			EnableLearning: true,
			ToolName:       "exercise_advisor",
		},
	}
}

func (t *ExerciseAdvisorTool) Info(_ context.Context) (*schema.ToolInfo, error) {
	return &schema.ToolInfo{
		Name: "exercise_advisor",
		Desc: "基于健康数据提供个性化运动建议、运动计划制定、效果跟踪。特别针对糖尿病患者的运动指导",
		ParamsOneOf: schema.NewParamsOneOfByParams(map[string]*schema.ParameterInfo{
			"action":        {Desc: "操作类型: plan(制定计划), advice(获取建议), track(效果跟踪)", Type: schema.String, Required: true},
			"weight":        {Desc: "体重(kg)", Type: schema.Number, Required: false},
			"height":        {Desc: "身高(cm)", Type: schema.Number, Required: false},
			"bmi":           {Desc: "BMI", Type: schema.Number, Required: false},
			"exercise_type": {Desc: "运动类型: running, walking, swimming, cycling", Type: schema.String, Required: false},
			"duration":      {Desc: "运动时长(分钟)", Type: schema.Integer, Required: false},
			"goal":          {Desc: "目标: weight_loss(减肥), fitness(健身), health(健康)", Type: schema.String, Required: false},
		}),
	}, nil
}

func (t *ExerciseAdvisorTool) InvokableRun(ctx context.Context, argumentsInJSON string, _ ...tool.Option) (string, error) {
	p := exerciseParams{}
	if err := sonic.UnmarshalString(argumentsInJSON, &p); err != nil {
		return "", err
	}

	// 验证参数
	if p.Action == "" {
		return "", fmt.Errorf("action is required")
	}

	var result exerciseResult

	switch p.Action {
	case "plan":
		result = t.generateExercisePlan(p.Weight, p.Height, p.BMI, p.Goal)
	case "advice":
		result = t.getExerciseAdvice(p.ExerciseType, fmt.Sprintf("%d", p.Duration), p.Weight)
	case "track":
		result = t.trackExerciseProgress(p.Weight, p.Height, p.BMI)
	default:
		return "", fmt.Errorf("unsupported action: %s", p.Action)
	}

	result.GeneratedAt = time.Now()

	b, _ := sonic.Marshal(result)
	return string(b), nil
}

// generateExercisePlan 生成运动计划
func (t *ExerciseAdvisorTool) generateExercisePlan(weight, height, bmi float64, goal string) exerciseResult {
	// 默认使用威猛先生的数据
	if weight == 0 {
		weight = 76.6
	}
	if height == 0 {
		height = 165
	}
	if bmi == 0 {
		bmi = weight / ((height / 100) * (height / 100))
	}
	if goal == "" {
		goal = "weight_loss" // 减肥目标
	}

	// 根据目标制定不同的计划
	var weeklySchedule []exerciseSession
	var totalDuration int
	var estimatedCalories int

	switch goal {
	case "weight_loss":
		// 减肥计划：有氧运动为主
		weeklySchedule = []exerciseSession{
			{Day: "周一", Type: "室内慢跑", Duration: 60, Intensity: "medium", Description: "室内慢跑1小时，控制心率在120-140次/分钟"},
			{Day: "周二", Type: "室内慢走", Duration: 60, Intensity: "low", Description: "室内慢走1小时，放松身心"},
			{Day: "周三", Type: "有氧全身运动", Duration: 20, Intensity: "medium", Description: "有氧全身运动20分钟，燃脂塑形"},
			{Day: "周四", Type: "室内慢跑", Duration: 60, Intensity: "medium", Description: "室内慢跑1小时，控制心率"},
			{Day: "周五", Type: "室内慢走", Duration: 60, Intensity: "low", Description: "室内慢走1小时，恢复体力"},
			{Day: "周六", Type: "有氧全身运动", Duration: 20, Intensity: "medium", Description: "有氧全身运动20分钟"},
			{Day: "周日", Type: "休息", Duration: 0, Intensity: "low", Description: "充分休息，让身体恢复"},
		}
		totalDuration = 320 // 5天运动
		estimatedCalories = 2000 // 每周约消耗2000卡路里

	case "fitness":
		// 健身计划：力量+有氧结合
		weeklySchedule = []exerciseSession{
			{Day: "周一", Type: "力量训练", Duration: 45, Intensity: "high", Description: "上肢力量训练45分钟"},
			{Day: "周二", Type: "室内慢跑", Duration: 30, Intensity: "medium", Description: "室内慢跑30分钟"},
			{Day: "周三", Type: "力量训练", Duration: 45, Intensity: "high", Description: "下肢力量训练45分钟"},
			{Day: "周四", Type: "室内慢走", Duration: 30, Intensity: "low", Description: "室内慢走30分钟"},
			{Day: "周五", Type: "有氧全身运动", Duration: 30, Intensity: "medium", Description: "有氧全身运动30分钟"},
			{Day: "周六", Type: "力量训练", Duration: 60, Intensity: "high", Description: "全身力量训练60分钟"},
			{Day: "周日", Type: "休息", Duration: 0, Intensity: "low", Description: "充分休息"},
		}
		totalDuration = 240
		estimatedCalories = 1800

	default: // health
		// 健康计划：温和运动
		weeklySchedule = []exerciseSession{
			{Day: "周一", Type: "室内慢走", Duration: 60, Intensity: "low", Description: "室内慢走1小时，促进血液循环"},
			{Day: "周二", Type: "有氧全身运动", Duration: 20, Intensity: "low", Description: "有氧全身运动20分钟"},
			{Day: "周三", Type: "室内慢走", Duration: 60, Intensity: "low", Description: "室内慢走1小时"},
			{Day: "周四", Type: "有氧全身运动", Duration: 20, Intensity: "low", Description: "有氧全身运动20分钟"},
			{Day: "周五", Type: "室内慢走", Duration: 60, Intensity: "low", Description: "室内慢走1小时"},
			{Day: "周六", Type: "休息", Duration: 0, Intensity: "low", Description: "休息日"},
			{Day: "周日", Type: "休息", Duration: 0, Intensity: "low", Description: "休息日"},
		}
		totalDuration = 220
		estimatedCalories = 1200
	}

	plan := &exercisePlan{
		WeeklySchedule:    weeklySchedule,
		TotalDuration:     totalDuration,
		EstimatedCalories: estimatedCalories,
	}

	recommendations := t.generateRecommendations(goal, bmi)

	return exerciseResult{
		Action:         "plan",
		Plan:           plan,
		Recommendations: recommendations,
		Message:        fmt.Sprintf("已为您制定%s计划，每周运动%d分钟，预计消耗%d卡路里", goal, totalDuration, estimatedCalories),
	}
}

// getExerciseAdvice 获取运动建议
func (t *ExerciseAdvisorTool) getExerciseAdvice(exerciseType, duration string, weight float64) exerciseResult {
	if weight == 0 {
		weight = 76.6
	}

	advice := ""
	caloriesBurned := 0

	// 解析时长
	durationMinutes := 60 // 默认60分钟
	if duration != "" {
		if d, err := fmt.Sscanf(duration, "%d", &durationMinutes); err == nil && d > 0 {
			// 使用解析的时长
		}
	}

	switch exerciseType {
	case "running":
		advice = "室内慢跑是很好的有氧运动，建议：1) 控制心率在120-140次/分钟 2) 保持均匀呼吸 3) 运动前后做拉伸 4) 穿着合适的运动鞋"
		caloriesBurned = 8 * durationMinutes
	case "walking":
		advice = "室内慢走适合恢复和放松，建议：1) 保持轻松的步伐 2) 可以听音乐或播客 3) 注意姿势端正 4) 逐渐增加时长"
		caloriesBurned = 4 * durationMinutes
	case "swimming":
		advice = "游泳是全身运动，建议：1) 注意安全 2) 控制游泳强度 3) 运动后及时擦干 4) 补充水分"
		caloriesBurned = 10 * durationMinutes
	default:
		advice = "有氧全身运动建议：1) 根据自身体力调整强度 2) 运动前后充分热身和拉伸 3) 保持规律运动 4) 注意运动安全"
		caloriesBurned = 6 * durationMinutes
	}

	recommendations := []string{
		"运动前30分钟不要大量进食",
		"运动中保持适当的水分补充",
		"运动后进行拉伸放松",
		"根据血糖情况调整运动强度",
	}

	return exerciseResult{
		Action:          "advice",
		Advice:          advice,
		CaloriesBurned:  caloriesBurned,
		Recommendations: recommendations,
		Message:         fmt.Sprintf("运动建议已生成，预计消耗%d卡路里", caloriesBurned),
	}
}

// trackExerciseProgress 跟踪运动效果
func (t *ExerciseAdvisorTool) trackExerciseProgress(weight, height, bmi float64) exerciseResult {
	if weight == 0 {
		weight = 76.6
	}
	if height == 0 {
		height = 165
	}
	if bmi == 0 {
		bmi = weight / ((height / 100) * (height / 100))
	}

	// 模拟跟踪数据
	targetWeight := 70.0 // 目标体重
	weightProgress := weight - targetWeight
	exerciseDays := 30 // 假设已运动30天
	totalDuration := 9600 // 假设总运动时长(分钟)
	consistencyScore := 85.0 // 坚持度评分

	var nextMilestone string
	if weightProgress > 5 {
		nextMilestone = "距离目标体重还有5kg，继续保持！"
	} else if weightProgress > 2 {
		nextMilestone = "距离目标体重还有2kg，最后冲刺！"
	} else {
		nextMilestone = "即将达到目标体重，保持当前节奏！"
	}

	tracking := &exerciseTracking{
		CurrentWeight:    weight,
		TargetWeight:     targetWeight,
		WeightProgress:   weightProgress,
		ExerciseDays:     exerciseDays,
		TotalDuration:    totalDuration,
		ConsistencyScore: consistencyScore,
		NextMilestone:    nextMilestone,
	}

	return exerciseResult{
		Action:   "track",
		Tracking: tracking,
		Message:  fmt.Sprintf("运动进度跟踪：当前体重%.1fkg，已运动%d天，坚持度%.1f分", weight, exerciseDays, consistencyScore),
	}
}

// generateRecommendations 生成个性化建议
func (t *ExerciseAdvisorTool) generateRecommendations(goal string, bmi float64) []string {
	recommendations := []string{
		"保持规律的运动习惯，每周至少5天",
		"运动前后监测血糖，避免低血糖",
		"根据身体感觉调整运动强度",
		"运动后及时补充水分和营养",
	}

	if bmi > 28 {
		recommendations = append(recommendations, "BMI较高，建议以低强度有氧运动为主")
	}

	if goal == "weight_loss" {
		recommendations = append(recommendations, "减肥期间配合饮食控制效果更佳")
		recommendations = append(recommendations, "建议每周测量体重，记录变化")
	}

	return recommendations
}

// AfterRun 工具调用后的后处理
// 集成学习引擎，保存运动建议到知识库
func (t *ExerciseAdvisorTool) AfterRun(ctx context.Context, result string, err error) {
	// 使用默认的BaseToolImpl后处理
	t.BaseToolImpl.AfterRun(ctx, result, err)
}
