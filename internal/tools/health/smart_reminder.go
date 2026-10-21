package health

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/schema"
	"new-robot-self/model"
	"new-robot-self/public"
	"xorm.io/xorm"
)

// SmartReminderTool 智能提醒工具
// 基于健康数据和日程的智能提醒系统
type SmartReminderTool struct {
	EnableLearning bool
	Category       string
	ToolName       string
	db             *xorm.Engine // 数据库实例
}

// SetDB 设置数据库实例
func (t *SmartReminderTool) SetDB(db *xorm.Engine) {
	t.db = db
}

// getDB 获取数据库实例
func (t *SmartReminderTool) getDB() *xorm.Engine {
	if t.db == nil {
		return public.DB
	}
	return t.db
}

// GetCategory 获取工具分类
func (t *SmartReminderTool) GetCategory() string {
	return t.Category
}

// GetToolName 获取工具名称
func (t *SmartReminderTool) GetToolName() string {
	return t.ToolName
}

// AfterRun 工具调用后的后处理
func (t *SmartReminderTool) AfterRun(ctx context.Context, result string, err error) {
	// 简化实现，不进行学习集成以避免循环依赖
}

type smartReminderParams struct {
	Action       string                 `json:"action"`       // create, query, update, delete, trigger
	ReminderType string                 `json:"reminder_type,omitempty"` // blood_sugar, exercise, diet, medication, work_rest
	Title        string                 `json:"title,omitempty"`
	Description  string                 `json:"description,omitempty"`
	ScheduledAt  int64                  `json:"scheduled_at,omitempty"`  // 计划提醒时间(时间戳)
	RepeatRule   string                 `json:"repeat_rule,omitempty"`   // 重复规则: daily, weekly, custom
	Status       string                 `json:"status,omitempty"`        // pending, completed, skipped, cancelled
	ID           int64                  `json:"id,omitempty"`            // 提醒ID
	Metadata     map[string]interface{} `json:"metadata,omitempty"`
}

type reminderTrigger struct {
	Reminders     []map[string]interface{} `json:"reminders"`
	NextActions   []string                 `json:"next_actions"`
	HealthContext healthContext            `json:"health_context"`
}

type healthContext struct {
	CurrentBloodSugar float64 `json:"current_blood_sugar"`
	CurrentWeight     float64 `json:"current_weight"`
	LastExerciseTime  int64   `json:"last_exercise_time"`
	ExerciseCycleDay  int     `json:"exercise_cycle_day"` // 1: 跑步, 2: 走路, 3: 有氧
}

func NewSmartReminderTool() *SmartReminderTool {
	return &SmartReminderTool{
		Category:       "health",
		EnableLearning: true,
		ToolName:       "smart_reminder",
	}
}

func (t *SmartReminderTool) Info(_ context.Context) (*schema.ToolInfo, error) {
	return &schema.ToolInfo{
		Name: "smart_reminder",
		Desc: "智能提醒工具，基于健康数据和日程的智能提醒系统，支持测血糖、运动、饮食、用药、工作休息提醒",
		ParamsOneOf: schema.NewParamsOneOfByParams(map[string]*schema.ParameterInfo{
			"action":        {Desc: "操作类型: create(创建), query(查询), update(更新), delete(删除), trigger(触发)", Type: schema.String, Required: true},
			"reminder_type": {Desc: "提醒类型: blood_sugar(测血糖), exercise(运动), diet(饮食), medication(用药), work_rest(工作休息)", Type: schema.String, Required: false},
			"title":         {Desc: "提醒标题", Type: schema.String, Required: false},
			"description":   {Desc: "提醒描述", Type: schema.String, Required: false},
			"scheduled_at":  {Desc: "计划提醒时间(时间戳)", Type: schema.Integer, Required: false},
			"repeat_rule":   {Desc: "重复规则: daily, weekly, custom", Type: schema.String, Required: false},
			"status":        {Desc: "状态: pending, completed, skipped, cancelled", Type: schema.String, Required: false},
			"id":            {Desc: "提醒ID（update/delete时需要）", Type: schema.Integer, Required: false},
			"metadata":      {Desc: "额外元数据", Type: schema.Object, Required: false},
		}),
	}, nil
}

func (t *SmartReminderTool) InvokableRun(ctx context.Context, argumentsInJSON string, _ ...tool.Option) (string, error) {
	var p smartReminderParams
	if err := json.Unmarshal([]byte(argumentsInJSON), &p); err != nil {
		return "", fmt.Errorf("参数解析失败: %w", err)
	}

	// 验证参数
	if p.Action == "" {
		return "", fmt.Errorf("action is required")
	}

	// 根据操作类型执行相应功能
	switch p.Action {
	case "create":
		return t.createReminder(p)
	case "query":
		return t.queryReminders(p)
	case "update":
		return t.updateReminder(p)
	case "delete":
		return t.deleteReminder(p)
	case "trigger":
		return t.triggerReminders(p)
	default:
		return "", fmt.Errorf("unknown action: %s", p.Action)
	}
}

// createReminder 创建提醒
func (t *SmartReminderTool) createReminder(p smartReminderParams) (string, error) {
	if p.ReminderType == "" {
		return "", fmt.Errorf("reminder_type is required for create action")
	}
	if p.Title == "" {
		return "", fmt.Errorf("title is required for create action")
	}
	if p.ScheduledAt == 0 {
		// 如果没有指定时间，根据类型设置默认时间
		p.ScheduledAt = t.getDefaultReminderTime(p.ReminderType)
	}

	db := t.getDB()
	if db == nil {
		return "", fmt.Errorf("database not initialized")
	}

	metadataJSON := "{}"
	if len(p.Metadata) > 0 {
		data, _ := json.Marshal(p.Metadata)
		metadataJSON = string(data)
	}

	reminder := &model.Reminder{
		UserID:      "weimeng_xiansheng",
		Type:        p.ReminderType,
		Title:       p.Title,
		Description: p.Description,
		ScheduledAt: p.ScheduledAt,
		RepeatRule:  p.RepeatRule,
		Status:      "pending",
		Priority:    t.getReminderPriority(p.ReminderType),
		Metadata:    metadataJSON,
	}

	_, err := db.Insert(reminder)
	if err != nil {
		return "", fmt.Errorf("failed to create reminder: %w", err)
	}

	result := map[string]interface{}{
		"success": true,
		"message": "提醒创建成功",
		"reminder": map[string]interface{}{
			"id":           reminder.Id,
			"type":         reminder.Type,
			"title":        reminder.Title,
			"description":  reminder.Description,
			"scheduled_at": reminder.ScheduledAt,
			"repeat_rule":  reminder.RepeatRule,
			"status":       reminder.Status,
			"priority":     reminder.Priority,
		},
	}

	return toJSON(result), nil
}

// queryReminders 查询提醒
func (t *SmartReminderTool) queryReminders(p smartReminderParams) (string, error) {
	db := t.getDB()
	if db == nil {
		return "", fmt.Errorf("database not initialized")
	}

	query := db.Where("user_id = ?", "weimeng_xiansheng")

	// 可选的类型过滤
	if p.ReminderType != "" && p.ReminderType != "all" {
		query = query.And("type = ?", p.ReminderType)
	}

	// 可选的状态过滤
	if p.Status != "" && p.Status != "all" {
		query = query.And("status = ?", p.Status)
	}

	var reminders []model.Reminder
	err := query.OrderBy("scheduled_at ASC").Find(&reminders)
	if err != nil {
		return "", fmt.Errorf("failed to query reminders: %w", err)
	}

	resultReminders := make([]map[string]interface{}, 0, len(reminders))
	for _, reminder := range reminders {
		resultReminders = append(resultReminders, map[string]interface{}{
			"id":           reminder.Id,
			"type":         reminder.Type,
			"title":        reminder.Title,
			"description":  reminder.Description,
			"scheduled_at": reminder.ScheduledAt,
			"repeat_rule":  reminder.RepeatRule,
			"status":       reminder.Status,
			"priority":     reminder.Priority,
		})
	}

	result := map[string]interface{}{
		"success":   true,
		"reminders": resultReminders,
		"count":     len(reminders),
	}

	return toJSON(result), nil
}

// updateReminder 更新提醒
func (t *SmartReminderTool) updateReminder(p smartReminderParams) (string, error) {
	if p.ID == 0 {
		return "", fmt.Errorf("id is required for update action")
	}

	db := t.getDB()
	if db == nil {
		return "", fmt.Errorf("database not initialized")
	}

	var reminder model.Reminder
	has, err := db.ID(p.ID).Where("user_id = ?", "weimeng_xiansheng").Get(&reminder)
	if err != nil {
		return "", fmt.Errorf("failed to query reminder: %w", err)
	}
	if !has {
		return "", fmt.Errorf("reminder not found: %d", p.ID)
	}

	// 构建更新字段
	updates := make(map[string]interface{})

	if p.Title != "" {
		updates["title"] = p.Title
	}
	if p.Description != "" {
		updates["description"] = p.Description
	}
	if p.ScheduledAt != 0 {
		updates["scheduled_at"] = p.ScheduledAt
	}
	if p.RepeatRule != "" {
		updates["repeat_rule"] = p.RepeatRule
	}
	if p.Status != "" {
		updates["status"] = p.Status
		if p.Status == "completed" {
			updates["completed_at"] = time.Now().Unix()
		}
	}

	if len(updates) == 0 {
		return "", fmt.Errorf("没有需要更新的字段")
	}

	_, err = db.ID(p.ID).Update(&reminder, updates)
	if err != nil {
		return "", fmt.Errorf("failed to update reminder: %w", err)
	}

	// 重新查询获取最新数据
	db.ID(p.ID).Get(&reminder)

	result := map[string]interface{}{
		"success": true,
		"message": "提醒更新成功",
		"reminder": map[string]interface{}{
			"id":           reminder.Id,
			"type":         reminder.Type,
			"title":        reminder.Title,
			"description":  reminder.Description,
			"scheduled_at": reminder.ScheduledAt,
			"repeat_rule":  reminder.RepeatRule,
			"status":       reminder.Status,
			"priority":     reminder.Priority,
		},
	}

	return toJSON(result), nil
}

// deleteReminder 删除提醒
func (t *SmartReminderTool) deleteReminder(p smartReminderParams) (string, error) {
	if p.ID == 0 {
		return "", fmt.Errorf("id is required for delete action")
	}

	db := t.getDB()
	if db == nil {
		return "", fmt.Errorf("database not initialized")
	}

	_, err := db.Where("id = ? AND user_id = ?", p.ID, "weimeng_xiansheng").Delete(&model.Reminder{})
	if err != nil {
		return "", fmt.Errorf("failed to delete reminder: %w", err)
	}

	result := map[string]interface{}{
		"success": true,
		"message": "提醒删除成功",
	}

	return toJSON(result), nil
}

// triggerReminders 触发提醒检查
func (t *SmartReminderTool) triggerReminders(p smartReminderParams) (string, error) {
	// 获取当前健康上下文
	healthCtx := t.getHealthContext()

	// 根据健康上下文生成智能提醒
	reminders := t.generateSmartReminders(healthCtx)

	trigger := reminderTrigger{
		Reminders:   reminders,
		NextActions: t.generateNextActions(healthCtx),
		HealthContext: healthCtx,
	}

	return toJSON(trigger), nil
}

// getDefaultReminderTime 获取默认提醒时间
func (t *SmartReminderTool) getDefaultReminderTime(reminderType string) int64 {
	now := time.Now()

	switch reminderType {
	case "blood_sugar":
		// 默认明天早上7点
		return time.Date(now.Year(), now.Month(), now.Day()+1, 7, 0, 0, 0, now.Location()).Unix()
	case "exercise":
		// 默认今天晚上8点
		return time.Date(now.Year(), now.Month(), now.Day(), 20, 0, 0, 0, now.Location()).Unix()
	case "diet":
		// 默认下一餐时间
		hour := now.Hour()
		if hour < 9 {
			return time.Date(now.Year(), now.Month(), now.Day(), 9, 0, 0, 0, now.Location()).Unix()
		} else if hour < 12 {
			return time.Date(now.Year(), now.Month(), now.Day(), 12, 0, 0, 0, now.Location()).Unix()
		} else if hour < 18 {
			return time.Date(now.Year(), now.Month(), now.Day(), 18, 0, 0, 0, now.Location()).Unix()
		} else {
			return time.Date(now.Year(), now.Month(), now.Day()+1, 7, 0, 0, 0, now.Location()).Unix()
		}
	case "work_rest":
		// 默认1小时后
		return now.Add(1 * time.Hour).Unix()
	default:
		return now.Add(1 * time.Hour).Unix()
	}
}

// getReminderPriority 获取提醒优先级
func (t *SmartReminderTool) getReminderPriority(reminderType string) string {
	priorities := map[string]string{
		"blood_sugar": "high",
		"medication":  "high",
		"exercise":    "medium",
		"diet":        "medium",
		"work_rest":   "low",
	}

	if priority, ok := priorities[reminderType]; ok {
		return priority
	}
	return "medium"
}

// getHealthContext 获取健康上下文
func (t *SmartReminderTool) getHealthContext() healthContext {
	db := t.getDB()

	// 获取最新血糖记录
	var currentBloodSugar float64 = 6.5 // 默认值
	if db != nil {
		var record model.HealthRecord
		has, _ := db.Where("user_id = ? AND record_type = ?", "weimeng_xiansheng", "blood_sugar").
			OrderBy("recorded_at DESC").Get(&record)
		if has {
			currentBloodSugar = record.Value
		}
	}

	// 获取最新体重记录
	var currentWeight float64 = 76.6 // 默认值
	if db != nil {
		var record model.HealthRecord
		has, _ := db.Where("user_id = ? AND record_type = ?", "weimeng_xiansheng", "weight").
			OrderBy("recorded_at DESC").Get(&record)
		if has {
			currentWeight = record.Value
		}
	}

	// 获取最后运动时间
	var lastExerciseTime int64 = 0
	if db != nil {
		var record model.HealthRecord
		has, _ := db.Where("user_id = ? AND record_type = ?", "weimeng_xiansheng", "exercise").
			OrderBy("recorded_at DESC").Get(&record)
		if has {
			lastExerciseTime = record.RecordedAt
		}
	}

	// 模拟威猛先生的运动循环（3天循环）
	now := time.Now()
	dayOffset := int(now.Unix() / 86400) % 3
	cycleDay := dayOffset + 1 // 1, 2, 3

	return healthContext{
		CurrentBloodSugar: currentBloodSugar,
		CurrentWeight:     currentWeight,
		LastExerciseTime:  lastExerciseTime,
		ExerciseCycleDay:  cycleDay,
	}
}

// generateSmartReminders 生成智能提醒
func (t *SmartReminderTool) generateSmartReminders(ctx healthContext) []map[string]interface{} {
	now := time.Now()
	reminders := []map[string]interface{}{}

	// 根据运动循环生成运动提醒
	var exerciseTitle, exerciseDesc string
	switch ctx.ExerciseCycleDay {
	case 1:
		exerciseTitle = "室内慢跑1小时"
		exerciseDesc = "运动循环第1天：室内慢跑，保持心率在120-140次/分钟"
	case 2:
		exerciseTitle = "室内慢走1小时"
		exerciseDesc = "运动循环第2天：室内慢走，轻松放松"
	case 3:
		exerciseTitle = "有氧全身运动20分钟"
		exerciseDesc = "运动循环第3天：有氧全身运动，高强度间歇训练"
	}

	reminders = append(reminders, map[string]interface{}{
		"type":         "exercise",
		"title":        exerciseTitle,
		"description":  exerciseDesc,
		"scheduled_at": time.Date(now.Year(), now.Month(), now.Day(), 20, 0, 0, 0, now.Location()).Unix(),
		"repeat_rule":  "daily",
		"priority":     "medium",
	})

	// 血糖提醒（如果血糖偏高）
	if ctx.CurrentBloodSugar > 7.0 {
		reminders = append(reminders, map[string]interface{}{
			"type":         "blood_sugar",
			"title":        "血糖偏高提醒",
			"description":  fmt.Sprintf("当前血糖%.1fmmol/L偏高，建议检查饮食和增加运动", ctx.CurrentBloodSugar),
			"scheduled_at": now.Unix(),
			"repeat_rule":  "",
			"priority":     "high",
		})
	}

	// 工作休息提醒
	reminders = append(reminders, map[string]interface{}{
		"type":         "work_rest",
		"title":        "工作休息提醒",
		"description":  "程序员需要定时休息，保护视力和颈椎",
		"scheduled_at": now.Add(1 * time.Hour).Unix(),
		"repeat_rule":  "daily",
		"priority":     "low",
	})

	return reminders
}

// generateNextActions 生成下一步行动建议
func (t *SmartReminderTool) generateNextActions(ctx healthContext) []string {
	actions := []string{}

	// 基于运动循环的建议
	switch ctx.ExerciseCycleDay {
	case 1:
		actions = append(actions, "今天晚上安排室内慢跑1小时")
		actions = append(actions, "准备运动装备：跑鞋、运动服")
		actions = append(actions, "运动前不要进食过多")
	case 2:
		actions = append(actions, "今天晚上安排室内慢走1小时")
		actions = append(actions, "可以边看电视边走路")
		actions = append(actions, "保持轻松的节奏")
	case 3:
		actions = append(actions, "今天晚上安排有氧全身运动20分钟")
		actions = append(actions, "准备高强度间歇训练")
		actions = append(actions, "运动后适当补充蛋白质")
	}

	// 基于血糖的建议
	if ctx.CurrentBloodSugar > 7.0 {
		actions = append(actions, "避免高GI食物")
		actions = append(actions, "增加蔬菜摄入")
		actions = append(actions, "餐后30分钟进行轻度运动")
	}

	// 基于体重的建议
	if ctx.CurrentWeight > 70.0 {
		actions = append(actions, "继续控制饮食热量")
		actions = append(actions, "保持规律运动")
		actions = append(actions, "定期记录体重变化")
	}

	return actions
}
