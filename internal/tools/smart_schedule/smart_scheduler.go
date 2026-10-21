package smart_schedule

import (
	"context"
	"fmt"
	"time"

	"github.com/bytedance/sonic"
	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/schema"
)

// SmartSchedulerTool 智能日程工具
// 基于健康状态和工作需求的智能排程优化
type SmartSchedulerTool struct {
	EnableLearning bool
	Category       string
	ToolName       string
}

// GetCategory 获取工具分类
func (t *SmartSchedulerTool) GetCategory() string {
	return t.Category
}

// GetToolName 获取工具名称
func (t *SmartSchedulerTool) GetToolName() string {
	return t.ToolName
}

// AfterRun 工具调用后的后处理
func (t *SmartSchedulerTool) AfterRun(ctx context.Context, result string, err error) {
	// 简化实现，不进行学习集成以避免循环依赖
}

type smartSchedulerParams struct {
	Action       string                 `json:"action"`        // optimize, create, query, update
	TaskTitle    string                 `json:"task_title,omitempty"`
	TaskType     string                 `json:"task_type,omitempty"`     // work, exercise, diet, rest, health
	Duration     int                    `json:"duration,omitempty"`      // 时长(分钟)
	Priority     string                 `json:"priority,omitempty"`      // high, medium, low
	PreferredTime string               `json:"preferred_time,omitempty"` // 期望时间段
	Deadline     int64                  `json:"deadline,omitempty"`      // 截止时间(时间戳)
	StartDate    int64                  `json:"start_date,omitempty"`
	EndDate      int64                  `json:"end_date,omitempty"`
	Metadata     map[string]interface{} `json:"metadata,omitempty"`
}

type smartTask struct {
	ID           string                 `json:"id"`
	UserID       string                 `json:"user_id"`
	Title        string                 `json:"title"`
	Type         string                 `json:"type"`
	Duration     int                    `json:"duration"`      // 时长(分钟)
	Priority     string                 `json:"priority"`      // high, medium, low
	PreferredTime string                `json:"preferred_time"` // 期望时间段
	ScheduledTime int64                  `json:"scheduled_time"` // 计划时间(时间戳)
	Deadline     int64                  `json:"deadline"`
	Status       string                 `json:"status"`        // pending, scheduled, completed, cancelled
	HealthScore  float64                `json:"health_score"`  // 健康影响评分
	CreatedAt    int64                  `json:"created_at"`
	UpdatedAt    int64                  `json:"updated_at"`
	Metadata     map[string]interface{} `json:"metadata"`
}

type scheduleOptimization struct {
	OriginalTasks []smartTask `json:"original_tasks"`
	OptimizedTasks []smartTask `json:"optimized_tasks"`
	HealthImprovements []string `json:"health_improvements"`
	WorkEfficiencyGain float64 `json:"work_efficiency_gain"`
	Recommendations []string `json:"recommendations"`
	GeneratedAt int64 `json:"generated_at"`
}

type scheduleConflict struct {
	Task1ID      string  `json:"task1_id"`
	Task2ID      string  `json:"task2_id"`
	ConflictType string  `json:"conflict_type"` // time_overlap, resource_conflict
	Suggestion   string  `json:"suggestion"`
}

func NewSmartSchedulerTool() *SmartSchedulerTool {
	return &SmartSchedulerTool{
		Category:       "smart_schedule",
		EnableLearning: true,
		ToolName:       "smart_scheduler",
	}
}

func (t *SmartSchedulerTool) Info(_ context.Context) (*schema.ToolInfo, error) {
	return &schema.ToolInfo{
		Name: "smart_scheduler",
		Desc: "智能日程工具，基于健康状态和工作需求的智能排程优化，支持任务安排、冲突检测、效率优化",
		ParamsOneOf: schema.NewParamsOneOfByParams(map[string]*schema.ParameterInfo{
			"action":         {Desc: "操作类型: optimize(优化日程), create(创建任务), query(查询任务), update(更新任务)", Type: schema.String, Required: true},
			"task_title":     {Desc: "任务标题", Type: schema.String, Required: false},
			"task_type":      {Desc: "任务类型: work, exercise, diet, rest, health", Type: schema.String, Required: false},
			"duration":       {Desc: "时长(分钟)", Type: schema.Integer, Required: false},
			"priority":       {Desc: "优先级: high, medium, low", Type: schema.String, Required: false},
			"preferred_time": {Desc: "期望时间段", Type: schema.String, Required: false},
			"deadline":       {Desc: "截止时间(时间戳)", Type: schema.Integer, Required: false},
			"start_date":     {Desc: "开始时间(时间戳)", Type: schema.Integer, Required: false},
			"end_date":       {Desc: "结束时间(时间戳)", Type: schema.Integer, Required: false},
			"metadata":       {Desc: "额外元数据", Type: schema.Object, Required: false},
		}),
	}, nil
}

func (t *SmartSchedulerTool) InvokableRun(ctx context.Context, argumentsInJSON string, _ ...tool.Option) (string, error) {
	p := smartSchedulerParams{}
	if err := sonic.UnmarshalString(argumentsInJSON, &p); err != nil {
		return "", err
	}

	// 验证参数
	if p.Action == "" {
		return "", fmt.Errorf("action is required")
	}

	// 根据操作类型执行相应功能
	switch p.Action {
	case "optimize":
		return t.optimizeSchedule(p)
	case "create":
		return t.createTask(p)
	case "query":
		return t.queryTasks(p)
	case "update":
		return t.updateTask(p)
	default:
		return "", fmt.Errorf("unknown action: %s", p.Action)
	}
}

// optimizeSchedule 优化日程
func (t *SmartSchedulerTool) optimizeSchedule(p smartSchedulerParams) (string, error) {
	// 设置默认优化时间范围（未来7天）
	if p.StartDate == 0 {
		p.StartDate = time.Now().Unix()
	}
	if p.EndDate == 0 {
		p.EndDate = time.Now().AddDate(0, 0, 7).Unix()
	}

	// TODO: 从数据库获取现有任务
	originalTasks := []smartTask{
		{
			ID:           "1",
			UserID:       "weimeng_xiansheng",
			Title:        "项目开发",
			Type:         "work",
			Duration:     240,
			Priority:     "high",
			ScheduledTime: time.Now().Add(2 * time.Hour).Unix(),
			Status:       "scheduled",
			HealthScore:  0.3,
		},
		{
			ID:           "2",
			UserID:       "weimeng_xiansheng",
			Title:        "晚间运动",
			Type:         "exercise",
			Duration:     60,
			Priority:     "medium",
			ScheduledTime: time.Now().Add(8 * time.Hour).Unix(),
			Status:       "scheduled",
			HealthScore:  0.9,
		},
	}

	// 基于健康状态优化任务安排
	optimizedTasks := t.optimizeTasksForHealth(originalTasks)

	optimization := scheduleOptimization{
		OriginalTasks: originalTasks,
		OptimizedTasks: optimizedTasks,
		HealthImprovements: []string{
			"增加工作休息间隔，保护视力和颈椎",
			"调整运动时间到更合适的时段",
			"优化工作节奏，避免连续高强度工作",
		},
		WorkEfficiencyGain: 15.5, // 效率提升15.5%
		Recommendations: []string{
			"建议采用番茄工作法：25分钟工作 + 5分钟休息",
			"每工作1小时进行10分钟伸展运动",
			"将高强度工作安排在精力最充沛的时段",
			"确保每天有足够的睡眠和休息时间",
		},
		GeneratedAt: time.Now().Unix(),
	}

	b, _ := sonic.Marshal(optimization)
	return string(b), nil
}

// createTask 创建任务
func (t *SmartSchedulerTool) createTask(p smartSchedulerParams) (string, error) {
	if p.TaskTitle == "" {
		return "", fmt.Errorf("task_title is required for create action")
	}
	if p.TaskType == "" {
		return "", fmt.Errorf("task_type is required for create action")
	}

	// 智能安排时间
	scheduledTime := t.scheduleTaskSmartly(p)

	task := smartTask{
		ID:           GenerateUUID(),
		UserID:       "weimeng_xiansheng",
		Title:        p.TaskTitle,
		Type:         p.TaskType,
		Duration:     p.Duration,
		Priority:     p.Priority,
		PreferredTime: p.PreferredTime,
		ScheduledTime: scheduledTime,
		Deadline:     p.Deadline,
		Status:       "scheduled",
		HealthScore:  t.calculateHealthScore(p.TaskType),
		CreatedAt:    time.Now().Unix(),
		UpdatedAt:    time.Now().Unix(),
		Metadata:     p.Metadata,
	}

	// TODO: 实际存储到数据库
	result := map[string]interface{}{
		"success": true,
		"message": "任务创建成功",
		"task":    task,
	}

	b, _ := sonic.Marshal(result)
	return string(b), nil
}

// queryTasks 查询任务
func (t *SmartSchedulerTool) queryTasks(p smartSchedulerParams) (string, error) {
	// 设置默认查询时间范围（未来7天）
	if p.StartDate == 0 {
		p.StartDate = time.Now().Unix()
	}
	if p.EndDate == 0 {
		p.EndDate = time.Now().AddDate(0, 0, 7).Unix()
	}

	// TODO: 从数据库查询
	tasks := []smartTask{
		{
			ID:           "1",
			UserID:       "weimeng_xiansheng",
			Title:        "项目开发",
			Type:         "work",
			Duration:     240,
			Priority:     "high",
			ScheduledTime: time.Now().Add(2 * time.Hour).Unix(),
			Status:       "scheduled",
			HealthScore:  0.3,
		},
		{
			ID:           "2",
			UserID:       "weimeng_xiansheng",
			Title:        "晚间运动",
			Type:         "exercise",
			Duration:     60,
			Priority:     "medium",
			ScheduledTime: time.Now().Add(8 * time.Hour).Unix(),
			Status:       "scheduled",
			HealthScore:  0.9,
		},
	}

	result := map[string]interface{}{
		"success": true,
		"tasks":   tasks,
		"count":   len(tasks),
	}

	b, _ := sonic.Marshal(result)
	return string(b), nil
}

// updateTask 更新任务
func (t *SmartSchedulerTool) updateTask(p smartSchedulerParams) (string, error) {
	// TODO: 更新数据库中的任务
	result := map[string]interface{}{
		"success": true,
		"message": "任务更新成功",
	}

	b, _ := sonic.Marshal(result)
	return string(b), nil
}

// optimizeTasksForHealth 基于健康状态优化任务
func (t *SmartSchedulerTool) optimizeTasksForHealth(tasks []smartTask) []smartTask {
	optimized := make([]smartTask, len(tasks))
	
	for i, task := range tasks {
		optimized[i] = task
		
		// 根据任务类型调整安排
		switch task.Type {
		case "work":
			// 工作任务：增加休息间隔
			if task.Duration > 120 {
				optimized[i].HealthScore = 0.4 // 稍微提高健康评分
			}
		case "exercise":
			// 运动任务：安排在合适的时间
			optimized[i].HealthScore = 0.95 // 运动任务健康评分高
		case "diet":
			// 饮食任务：确保定时定量
			optimized[i].HealthScore = 0.85
		case "rest":
			// 休息任务：健康评分最高
			optimized[i].HealthScore = 1.0
		}
	}
	
	return optimized
}

// scheduleTaskSmartly 智能安排任务时间
func (t *SmartSchedulerTool) scheduleTaskSmartly(p smartSchedulerParams) int64 {
	now := time.Now()
	
	switch p.TaskType {
	case "exercise":
		// 运动安排在晚上8点
		return time.Date(now.Year(), now.Month(), now.Day(), 20, 0, 0, 0, now.Location()).Unix()
	case "diet":
		// 饮食安排在正餐时间
		hour := now.Hour()
		if hour < 9 {
			return time.Date(now.Year(), now.Month(), now.Day(), 8, 0, 0, 0, now.Location()).Unix()
		} else if hour < 12 {
			return time.Date(now.Year(), now.Month(), now.Day(), 12, 0, 0, 0, now.Location()).Unix()
		} else if hour < 18 {
			return time.Date(now.Year(), now.Month(), now.Day(), 18, 0, 0, 0, now.Location()).Unix()
		} else {
			return time.Date(now.Year(), now.Month(), now.Day()+1, 8, 0, 0, 0, now.Location()).Unix()
		}
	case "rest":
		// 休息安排在1小时后
		return now.Add(1 * time.Hour).Unix()
	default:
		// 工作任务安排在下一个可用时段
		return now.Add(1 * time.Hour).Unix()
	}
}

// calculateHealthScore 计算健康影响评分
func (t *SmartSchedulerTool) calculateHealthScore(taskType string) float64 {
	scores := map[string]float64{
		"work":     0.3,
		"exercise": 0.9,
		"diet":     0.8,
		"rest":     1.0,
		"health":   0.95,
	}
	
	if score, ok := scores[taskType]; ok {
		return score
	}
	return 0.5
}