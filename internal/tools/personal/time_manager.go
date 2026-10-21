package personal

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/schema"
)

// TimeRecord 时间记录结构
type TimeRecord struct {
	ID          string    `json:"id"`
	Task        string    `json:"task"`        // 任务名称
	StartTime   int64     `json:"start_time"`   // 开始时间戳
	EndTime     int64     `json:"end_time"`     // 结束时间戳
	Duration    int       `json:"duration"`     // 持续时间（秒）
	Category    string    `json:"category"`     // 分类
	IsBreak     bool      `json:"is_break"`     // 是否为休息时间
	CreatedAt   time.Time `json:"created_at"`
}

// PomodoroSession 番茄钟会话
type PomodoroSession struct {
	ID           string    `json:"id"`
	Task         string    `json:"task"`
	StartTime    int64     `json:"start_time"`
	Duration     int       `json:"duration"`      // 番茄钟时长（秒）
	BreakType    string    `json:"break_type"`    // 休息类型 short | long
	Status       string    `json:"status"`        // running | paused | completed | cancelled
	CompletedAt  time.Time `json:"completed_at"`
	CreatedAt    time.Time `json:"created_at"`
}

// TimeManager 时间管理工具
type TimeManager struct {
	enableLearning  bool
	category        string
	toolName        string
	timeRecords     map[string]*TimeRecord     // 时间记录
	pomodoroSession *PomodoroSession            // 当前番茄钟会话
}

// NewTimeManager 创建时间管理工具
func NewTimeManager() *TimeManager {
	return &TimeManager{
		enableLearning: true,
		category:       "personal",
		toolName:       "time_manager",
		timeRecords:    make(map[string]*TimeRecord),
	}
}

// Info 返回工具信息
func (t *TimeManager) Info(ctx context.Context) (*schema.ToolInfo, error) {
	return &schema.ToolInfo{
		Name: "time_manager",
		Desc: "时间管理工具，支持番茄钟计时、时间记录和统计、专注时间分析、效率报告生成",
		ParamsOneOf: schema.NewParamsOneOfByParams(map[string]*schema.ParameterInfo{
			"action": {
				Desc:     "操作类型",
				Type:     schema.String,
				Required: true,
				Enum:     []string{"start_pomodoro", "stop_pomodoro", "pause_pomodoro", "resume_pomodoro", "get_pomodoro_status", "record_time", "get_time_stats", "get_efficiency_report"},
			},
			"task": {
				Desc:     "任务名称",
				Type:     schema.String,
				Required: false,
			},
			"duration": {
				Desc:     "持续时间（秒）",
				Type:     schema.Integer,
				Required: false,
			},
			"break_type": {
				Desc:     "休息类型",
				Type:     schema.String,
				Required: false,
				Enum:     []string{"short", "long"},
			},
			"category": {
				Desc:     "分类",
				Type:     schema.String,
				Required: false,
			},
			"start_time": {
				Desc:     "开始时间戳",
				Type:     schema.Integer,
				Required: false,
			},
			"end_time": {
				Desc:     "结束时间戳",
				Type:     schema.Integer,
				Required: false,
			},
			"is_break": {
				Desc:     "是否为休息时间",
				Type:     schema.Boolean,
				Required: false,
			},
			"date": {
				Desc:     "日期（YYYY-MM-DD）",
				Type:     schema.String,
				Required: false,
			},
			"days": {
				Desc:     "统计天数",
				Type:     schema.Integer,
				Required: false,
			},
		}),
	}, nil
}

// InvokableRun 执行工具调用
func (t *TimeManager) InvokableRun(ctx context.Context, args string, opts ...tool.Option) (string, error) {
	var params map[string]interface{}
	if err := json.Unmarshal([]byte(args), &params); err != nil {
		return "", fmt.Errorf("参数解析失败: %w", err)
	}

	action := params["action"].(string)
	switch action {
	case "start_pomodoro":
		return t.startPomodoro(params)
	case "stop_pomodoro":
		return t.stopPomodoro()
	case "pause_pomodoro":
		return t.pausePomodoro()
	case "resume_pomodoro":
		return t.resumePomodoro()
	case "get_pomodoro_status":
		return t.getPomodoroStatus()
	case "record_time":
		return t.recordTime(params)
	case "get_time_stats":
		return t.getTimeStats(params)
	case "get_efficiency_report":
		return t.getEfficiencyReport(params)
	default:
		return "", fmt.Errorf("不支持的操作类型: %s", action)
	}
}

// startPomodoro 开始番茄钟
func (t *TimeManager) startPomodoro(params map[string]interface{}) (string, error) {
	// 检查是否已有运行中的番茄钟
	if t.pomodoroSession != nil && t.pomodoroSession.Status == "running" {
		return "", fmt.Errorf("已有运行中的番茄钟，请先停止或暂停")
	}

	task := getStringParam(params, "task", "未命名任务")
	duration := getIntParam(params, "duration", 1500) // 默认25分钟
	breakType := getStringParam(params, "break_type", "short")

	session := &PomodoroSession{
		ID:        GenerateUUID(),
		Task:      task,
		StartTime: time.Now().Unix(),
		Duration:  duration,
		BreakType: breakType,
		Status:    "running",
		CreatedAt: time.Now(),
	}

	t.pomodoroSession = session

	result := map[string]interface{}{
		"success": true,
		"session": session,
		"message": "番茄钟已启动",
	}

	return toJSON(result), nil
}

// stopPomodoro 停止番茄钟
func (t *TimeManager) stopPomodoro() (string, error) {
	if t.pomodoroSession == nil {
		return "", fmt.Errorf("没有运行中的番茄钟")
	}

	t.pomodoroSession.Status = "completed"
	t.pomodoroSession.CompletedAt = time.Now()

	// 记录时间
	record := &TimeRecord{
		ID:        GenerateUUID(),
		Task:      t.pomodoroSession.Task,
		StartTime: t.pomodoroSession.StartTime,
		EndTime:   time.Now().Unix(),
		Duration:  int(time.Now().Unix() - t.pomodoroSession.StartTime),
		Category:  "pomodoro",
		IsBreak:   false,
		CreatedAt: time.Now(),
	}

	t.timeRecords[record.ID] = record

	session := t.pomodoroSession
	t.pomodoroSession = nil

	result := map[string]interface{}{
		"success": true,
		"session": session,
		"record":  record,
		"message": "番茄钟已完成并记录",
	}

	return toJSON(result), nil
}

// pausePomodoro 暂停番茄钟
func (t *TimeManager) pausePomodoro() (string, error) {
	if t.pomodoroSession == nil || t.pomodoroSession.Status != "running" {
		return "", fmt.Errorf("没有运行中的番茄钟")
	}

	t.pomodoroSession.Status = "paused"

	result := map[string]interface{}{
		"success": true,
		"session": t.pomodoroSession,
		"message": "番茄钟已暂停",
	}

	return toJSON(result), nil
}

// resumePomodoro 恢复番茄钟
func (t *TimeManager) resumePomodoro() (string, error) {
	if t.pomodoroSession == nil || t.pomodoroSession.Status != "paused" {
		return "", fmt.Errorf("没有暂停中的番茄钟")
	}

	t.pomodoroSession.Status = "running"

	result := map[string]interface{}{
		"success": true,
		"session": t.pomodoroSession,
		"message": "番茄钟已恢复",
	}

	return toJSON(result), nil
}

// getPomodoroStatus 获取番茄钟状态
func (t *TimeManager) getPomodoroStatus() (string, error) {
	if t.pomodoroSession == nil {
		result := map[string]interface{}{
			"success": true,
			"status":  "none",
			"message": "没有运行中的番茄钟",
		}
		return toJSON(result), nil
	}

	// 计算剩余时间
	elapsed := time.Now().Unix() - t.pomodoroSession.StartTime
	remaining := t.pomodoroSession.Duration - int(elapsed)
	if remaining < 0 {
		remaining = 0
	}

	result := map[string]interface{}{
		"success":   true,
		"status":    t.pomodoroSession.Status,
		"session":   t.pomodoroSession,
		"elapsed":   elapsed,
		"remaining": remaining,
		"progress":  float64(elapsed) / float64(t.pomodoroSession.Duration) * 100,
	}

	return toJSON(result), nil
}

// recordTime 记录时间
func (t *TimeManager) recordTime(params map[string]interface{}) (string, error) {
	task := getStringParam(params, "task", "未命名任务")
	category := getStringParam(params, "category", "general")
	startTime := getInt64Param(params, "start_time", time.Now().Unix())
	endTime := getInt64Param(params, "end_time", time.Now().Unix())
	isBreak := getBoolParam(params, "is_break", false)

	if endTime < startTime {
		return "", fmt.Errorf("结束时间不能早于开始时间")
	}

	record := &TimeRecord{
		ID:        GenerateUUID(),
		Task:      task,
		StartTime: startTime,
		EndTime:   endTime,
		Duration:  int(endTime - startTime),
		Category:  category,
		IsBreak:   isBreak,
		CreatedAt: time.Now(),
	}

	t.timeRecords[record.ID] = record

	result := map[string]interface{}{
		"success": true,
		"record":  record,
		"message": "时间记录成功",
	}

	return toJSON(result), nil
}

// getTimeStats 获取时间统计
func (t *TimeManager) getTimeStats(params map[string]interface{}) (string, error) {
	date := getStringParam(params, "date", time.Now().Format("2006-01-02"))
	days := getIntParam(params, "days", 7)

	targetDate, err := time.Parse("2006-01-02", date)
	if err != nil {
		return "", fmt.Errorf("日期格式错误: %w", err)
	}

	startDate := targetDate.AddDate(0, 0, -days+1).Truncate(24 * time.Hour)
	endDate := startDate.AddDate(0, 0, days)

	var totalDuration int
	var workDuration int
	var breakDuration int
	categoryStats := make(map[string]int)
	dailyStats := make(map[string]int)

	for _, record := range t.timeRecords {
		recordTime := time.Unix(record.StartTime, 0)
		if recordTime.Before(startDate) || recordTime.After(endDate) {
			continue
		}

		totalDuration += record.Duration

		if record.IsBreak {
			breakDuration += record.Duration
		} else {
			workDuration += record.Duration
		}

		// 分类统计
		categoryStats[record.Category] += record.Duration

		// 每日统计
		dayKey := recordTime.Format("2006-01-02")
		dailyStats[dayKey] += record.Duration
	}

	// 计算平均值
	avgDailyDuration := 0
	if len(dailyStats) > 0 {
		avgDailyDuration = totalDuration / len(dailyStats)
	}

	result := map[string]interface{}{
		"success":           true,
		"period":            map[string]interface{}{"start": startDate.Format("2006-01-02"), "end": endDate.Format("2006-01-02")},
		"total_duration":    totalDuration,
		"work_duration":     workDuration,
		"break_duration":    breakDuration,
		"avg_daily_duration": avgDailyDuration,
		"category_stats":    categoryStats,
		"daily_stats":       dailyStats,
		"record_count":      len(t.timeRecords),
	}

	return toJSON(result), nil
}

// getEfficiencyReport 获取效率报告
func (t *TimeManager) getEfficiencyReport(params map[string]interface{}) (string, error) {
	date := getStringParam(params, "date", time.Now().Format("2006-01-02"))
	days := getIntParam(params, "days", 7)

	targetDate, err := time.Parse("2006-01-02", date)
	if err != nil {
		return "", fmt.Errorf("日期格式错误: %w", err)
	}

	startDate := targetDate.AddDate(0, 0, -days+1).Truncate(24 * time.Hour)
	endDate := startDate.AddDate(0, 0, days)

	var workDuration int
	var breakDuration int
	taskCount := make(map[string]int)
	productiveHours := 0

	for _, record := range t.timeRecords {
		recordTime := time.Unix(record.StartTime, 0)
		if recordTime.Before(startDate) || recordTime.After(endDate) {
			continue
		}

		if !record.IsBreak {
			workDuration += record.Duration
			taskCount[record.Task]++
			
			// 统计高效工作时间（假设每次工作超过1小时为高效）
			if record.Duration >= 3600 {
				productiveHours++
			}
		} else {
			breakDuration += record.Duration
		}
	}

	// 计算效率指标
	totalDuration := workDuration + breakDuration
	efficiencyRatio := 0.0
	if totalDuration > 0 {
		efficiencyRatio = float64(workDuration) / float64(totalDuration) * 100
	}

	// 找出最常做的任务
	topTasks := make([]string, 0, 5)
	for task, count := range taskCount {
		topTasks = append(topTasks, fmt.Sprintf("%s(%d次)", task, count))
	}

	result := map[string]interface{}{
		"success":           true,
		"period":            map[string]interface{}{"start": startDate.Format("2006-01-02"), "end": endDate.Format("2006-01-02")},
		"efficiency_ratio":  efficiencyRatio,
		"work_duration":     workDuration,
		"break_duration":    breakDuration,
		"total_duration":    totalDuration,
		"productive_hours":  productiveHours,
		"task_count":        len(taskCount),
		"top_tasks":         topTasks,
		"recommendations":   t.generateRecommendations(efficiencyRatio, productiveHours, len(taskCount)),
	}

	return toJSON(result), nil
}

// generateRecommendations 生成效率建议
func (t *TimeManager) generateRecommendations(efficiencyRatio float64, productiveHours, taskCount int) []string {
	var recommendations []string

	if efficiencyRatio < 60 {
		recommendations = append(recommendations, "工作效率较低，建议减少休息时间或增加专注时间")
	} else if efficiencyRatio > 85 {
		recommendations = append(recommendations, "工作效率很高，注意适当休息以避免疲劳")
	} else {
		recommendations = append(recommendations, "工作效率良好，保持当前节奏")
	}

	if productiveHours < 3 {
		recommendations = append(recommendations, "深度工作时间较少，建议安排更多连续专注时段")
	} else if productiveHours > 8 {
		recommendations = append(recommendations, "深度工作时间充足，注意劳逸结合")
	}

	if taskCount < 5 {
		recommendations = append(recommendations, "任务种类较少，可以尝试多样化工作内容")
	} else if taskCount > 15 {
		recommendations = append(recommendations, "任务切换频繁，建议集中精力处理相似任务")
	}

	if len(recommendations) == 0 {
		recommendations = append(recommendations, "时间管理良好，继续保持")
	}

	return recommendations
}

// AfterRun 工具调用后的后处理
func (t *TimeManager) AfterRun(ctx context.Context, result string, err error) {
	// 简化实现，暂时不做学习集成
}