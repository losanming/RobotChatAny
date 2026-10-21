package personal

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

// CalendarManager 日程管理工具
type CalendarManager struct {
	enableLearning bool
	category       string
	toolName       string
	db             *xorm.Engine // 数据库实例
}

// NewCalendarManager 创建日程管理工具
func NewCalendarManager() *CalendarManager {
	return &CalendarManager{
		enableLearning: true,
		category:       "personal",
		toolName:       "calendar_manager",
	}
}

// SetDB 设置数据库实例
func (c *CalendarManager) SetDB(db *xorm.Engine) {
	c.db = db
}

// getDB 获取数据库实例
func (c *CalendarManager) getDB() *xorm.Engine {
	if c.db == nil {
		return public.DB
	}
	return c.db
}

// Info 返回工具信息
func (c *CalendarManager) Info(ctx context.Context) (*schema.ToolInfo, error) {
	return &schema.ToolInfo{
		Name: "calendar_manager",
		Desc: "日程管理工具，支持创建、查询、更新、删除日程事件，日历视图生成，冲突检测",
		ParamsOneOf: schema.NewParamsOneOfByParams(map[string]*schema.ParameterInfo{
			"action": {
				Desc:     "操作类型",
				Type:     schema.String,
				Required: true,
				Enum:     []string{"create", "query", "update", "delete", "list", "check_conflict", "calendar_view"},
			},
			"id": {
				Desc:     "事件ID（update/delete/query时需要）",
				Type:     schema.String,
				Required: false,
			},
			"title": {
				Desc:     "事件标题（create时需要）",
				Type:     schema.String,
				Required: false,
			},
			"start_time": {
				Desc:     "开始时间戳（create时需要）",
				Type:     schema.Integer,
				Required: false,
			},
			"duration": {
				Desc:     "持续时间（秒）（create时需要）",
				Type:     schema.Integer,
				Required: false,
			},
			"location": {
				Desc:     "地点",
				Type:     schema.String,
				Required: false,
			},
			"description": {
				Desc:     "描述",
				Type:     schema.String,
				Required: false,
			},
			"reminder": {
				Desc:     "提前提醒时间（秒）",
				Type:     schema.Integer,
				Required: false,
			},
			"start_date": {
				Desc:     "查询开始日期（YYYY-MM-DD）",
				Type:     schema.String,
				Required: false,
			},
			"end_date": {
				Desc:     "查询结束日期（YYYY-MM-DD）",
				Type:     schema.String,
				Required: false,
			},
			"view_type": {
				Desc:     "日历视图类型（day/week/month）",
				Type:     schema.String,
				Required: false,
				Enum:     []string{"day", "week", "month"},
			},
		}),
	}, nil
}

// InvokableRun 执行工具调用
func (c *CalendarManager) InvokableRun(ctx context.Context, args string, opts ...tool.Option) (string, error) {
	var params map[string]interface{}
	if err := json.Unmarshal([]byte(args), &params); err != nil {
		return "", fmt.Errorf("参数解析失败: %w", err)
	}

	action := params["action"].(string)
	switch action {
	case "create":
		return c.createEvent(params)
	case "query":
		return c.queryEvent(params)
	case "update":
		return c.updateEvent(params)
	case "delete":
		return c.deleteEvent(params)
	case "list":
		return c.listEvents(params)
	case "check_conflict":
		return c.checkConflict(params)
	case "calendar_view":
		return c.calendarView(params)
	default:
		return "", fmt.Errorf("不支持的操作类型: %s", action)
	}
}

// createEvent 创建事件
func (c *CalendarManager) createEvent(params map[string]interface{}) (string, error) {
	title := getStringParam(params, "title", "")
	if title == "" {
		return "", fmt.Errorf("标题不能为空")
	}

	startTime := getInt64Param(params, "start_time", 0)
	if startTime == 0 {
		return "", fmt.Errorf("开始时间不能为空")
	}

	duration := getInt64Param(params, "duration", 0)
	if duration == 0 {
		return "", fmt.Errorf("持续时间不能为空")
	}

	db := c.getDB()
	if db == nil {
		return "", fmt.Errorf("database not initialized")
	}

	// 检查冲突
	conflicts := c.checkTimeConflict(startTime, duration, 0)
	if len(conflicts) > 0 {
		return "", fmt.Errorf("时间冲突，与以下事件重叠: %v", conflicts)
	}

	event := &model.CalendarEvent{
		UserID:      "weimeng_xiansheng",
		Title:       title,
		Description: getStringParam(params, "description", ""),
		Location:    getStringParam(params, "location", ""),
		StartTime:   startTime,
		EndTime:     startTime + duration,
		IsAllDay:    getBoolParam(params, "is_all_day", false),
		Reminder:    getInt64Param(params, "reminder", 0),
	}

	_, err := db.Insert(event)
	if err != nil {
		return "", fmt.Errorf("failed to create event: %w", err)
	}

	result := map[string]interface{}{
		"success": true,
		"event": map[string]interface{}{
			"id":          event.Id,
			"title":       event.Title,
			"description": event.Description,
			"location":    event.Location,
			"start_time":  event.StartTime,
			"end_time":    event.EndTime,
			"is_all_day":  event.IsAllDay,
			"reminder":    event.Reminder,
			"created_at":  event.CreatedAt,
		},
		"message": "日程创建成功",
	}

	return toJSON(result), nil
}

// queryEvent 查询事件
func (c *CalendarManager) queryEvent(params map[string]interface{}) (string, error) {
	id := getInt64Param(params, "id", 0)
	if id == 0 {
		return "", fmt.Errorf("事件ID不能为空")
	}

	db := c.getDB()
	if db == nil {
		return "", fmt.Errorf("database not initialized")
	}

	var event model.CalendarEvent
	has, err := db.ID(id).Where("user_id = ?", "weimeng_xiansheng").Get(&event)
	if err != nil {
		return "", fmt.Errorf("failed to query event: %w", err)
	}
	if !has {
		return "", fmt.Errorf("事件不存在: %d", id)
	}

	result := map[string]interface{}{
		"success": true,
		"event": map[string]interface{}{
			"id":          event.Id,
			"title":       event.Title,
			"description": event.Description,
			"location":    event.Location,
			"start_time":  event.StartTime,
			"end_time":    event.EndTime,
			"is_all_day":  event.IsAllDay,
			"reminder":    event.Reminder,
			"status":      event.Status,
			"created_at":  event.CreatedAt,
			"updated_at":  event.UpdatedAt,
		},
	}

	return toJSON(result), nil
}

// updateEvent 更新事件
func (c *CalendarManager) updateEvent(params map[string]interface{}) (string, error) {
	id := getInt64Param(params, "id", 0)
	if id == 0 {
		return "", fmt.Errorf("事件ID不能为空")
	}

	db := c.getDB()
	if db == nil {
		return "", fmt.Errorf("database not initialized")
	}

	var event model.CalendarEvent
	has, err := db.ID(id).Where("user_id = ?", "weimeng_xiansheng").Get(&event)
	if err != nil {
		return "", fmt.Errorf("failed to query event: %w", err)
	}
	if !has {
		return "", fmt.Errorf("事件不存在: %d", id)
	}

	// 构建更新字段
	updates := make(map[string]interface{})

	if title := getStringParam(params, "title", ""); title != "" {
		updates["title"] = title
	}
	if desc := getStringParam(params, "description", ""); desc != "" {
		updates["description"] = desc
	}
	if location := getStringParam(params, "location", ""); location != "" {
		updates["location"] = location
	}

	startTime := getInt64Param(params, "start_time", 0)
	endTime := getInt64Param(params, "end_time", 0)
	
	if startTime > 0 {
		updates["start_time"] = startTime
		if endTime == 0 {
			// 如果只提供了开始时间，保持原有持续时间
			duration := event.EndTime - event.StartTime
			updates["end_time"] = startTime + duration
		} else {
			updates["end_time"] = endTime
		}
		
		// 检查时间冲突
		checkEndTime := endTime
		if checkEndTime == 0 {
			checkEndTime = startTime + (event.EndTime - event.StartTime)
		}
		conflicts := c.checkTimeConflict(startTime, checkEndTime-startTime, id)
		if len(conflicts) > 0 {
			return "", fmt.Errorf("时间冲突，与以下事件重叠: %v", conflicts)
		}
	} else if endTime > 0 {
		updates["end_time"] = endTime
	}

	if reminder := getInt64Param(params, "reminder", -1); reminder >= 0 {
		updates["reminder"] = reminder
	}

	status := getStringParam(params, "status", "")
	if status == "confirmed" || status == "cancelled" || status == "tentative" {
		updates["status"] = status
	}

	if len(updates) == 0 {
		return "", fmt.Errorf("没有需要更新的字段")
	}

	_, err = db.ID(id).Update(&event, updates)
	if err != nil {
		return "", fmt.Errorf("failed to update event: %w", err)
	}

	// 重新查询获取最新数据
	db.ID(id).Get(&event)

	result := map[string]interface{}{
		"success": true,
		"event": map[string]interface{}{
			"id":          event.Id,
			"title":       event.Title,
			"description": event.Description,
			"location":    event.Location,
			"start_time":  event.StartTime,
			"end_time":    event.EndTime,
			"is_all_day":  event.IsAllDay,
			"reminder":    event.Reminder,
			"status":      event.Status,
			"updated_at":  event.UpdatedAt,
		},
		"message": "日程更新成功",
	}

	return toJSON(result), nil
}

// deleteEvent 删除事件
func (c *CalendarManager) deleteEvent(params map[string]interface{}) (string, error) {
	id := getInt64Param(params, "id", 0)
	if id == 0 {
		return "", fmt.Errorf("事件ID不能为空")
	}

	db := c.getDB()
	if db == nil {
		return "", fmt.Errorf("database not initialized")
	}

	_, err := db.Where("id = ? AND user_id = ?", id, "weimeng_xiansheng").Delete(&model.CalendarEvent{})
	if err != nil {
		return "", fmt.Errorf("failed to delete event: %w", err)
	}

	result := map[string]interface{}{
		"success": true,
		"message": "日程删除成功",
	}

	return toJSON(result), nil
}

// listEvents 列出事件
func (c *CalendarManager) listEvents(params map[string]interface{}) (string, error) {
	db := c.getDB()
	if db == nil {
		return "", fmt.Errorf("database not initialized")
	}

	query := db.Where("user_id = ?", "weimeng_xiansheng")

	// 日期范围过滤
	startDate := getStringParam(params, "start_date", "")
	endDate := getStringParam(params, "end_date", "")

	if startDate != "" {
		startTS, _ := parseDateToTimestamp(startDate)
		query = query.And("start_time >= ?", startTS)
	}
	if endDate != "" {
		endTS, _ := parseDateToTimestamp(endDate)
		// 包含结束日期的全天
		endTS += 86400
		query = query.And("start_time < ?", endTS)
	}

	// 状态过滤
	status := getStringParam(params, "status", "")
	if status != "" && status != "all" {
		query = query.And("status = ?", status)
	}

	var events []model.CalendarEvent
	err := query.OrderBy("start_time ASC").Find(&events)
	if err != nil {
		return "", fmt.Errorf("failed to list events: %w", err)
	}

	resultEvents := make([]map[string]interface{}, 0, len(events))
	for _, event := range events {
		resultEvents = append(resultEvents, map[string]interface{}{
			"id":          event.Id,
			"title":       event.Title,
			"description": event.Description,
			"location":    event.Location,
			"start_time":  event.StartTime,
			"end_time":    event.EndTime,
			"is_all_day":  event.IsAllDay,
			"reminder":    event.Reminder,
			"status":      event.Status,
		})
	}

	result := map[string]interface{}{
		"success": true,
		"events":  resultEvents,
		"total":   len(events),
	}

	return toJSON(result), nil
}

// checkConflict 检查时间冲突
func (c *CalendarManager) checkConflict(params map[string]interface{}) (string, error) {
	startTime := getInt64Param(params, "start_time", 0)
	duration := getInt64Param(params, "duration", 0)
	excludeID := getInt64Param(params, "exclude_id", 0)

	if startTime == 0 || duration == 0 {
		return "", fmt.Errorf("开始时间和持续时间不能为空")
	}

	conflicts := c.checkTimeConflict(startTime, duration, excludeID)

	result := map[string]interface{}{
		"success":   true,
		"conflicts": conflicts,
		"has_conflict": len(conflicts) > 0,
	}

	return toJSON(result), nil
}

// calendarView 日历视图
func (c *CalendarManager) calendarView(params map[string]interface{}) (string, error) {
	viewType := getStringParam(params, "view_type", "week") // day, week, month
	date := getStringParam(params, "date", time.Now().Format("2006-01-02"))

	startTS, _ := parseDateToTimestamp(date)
	
	var endTS int64
	switch viewType {
	case "day":
		endTS = startTS + 86400
	case "week":
		endTS = startTS + 7*86400
	case "month":
		endTS = startTS + 30*86400
	default:
		endTS = startTS + 7*86400
	}

	db := c.getDB()
	if db == nil {
		return "", fmt.Errorf("database not initialized")
	}

	var events []model.CalendarEvent
	err := db.Where("user_id = ? AND start_time >= ? AND start_time < ?",
		"weimeng_xiansheng", startTS, endTS).
		OrderBy("start_time ASC").
		Find(&events)
	if err != nil {
		return "", fmt.Errorf("failed to query events: %w", err)
	}

	resultEvents := make([]map[string]interface{}, 0, len(events))
	for _, event := range events {
		resultEvents = append(resultEvents, map[string]interface{}{
			"id":          event.Id,
			"title":       event.Title,
			"start_time":  event.StartTime,
			"end_time":    event.EndTime,
			"is_all_day":  event.IsAllDay,
			"location":    event.Location,
		})
	}

	result := map[string]interface{}{
		"success": true,
		"view_type": viewType,
		"date":      date,
		"events":    resultEvents,
		"total":     len(events),
	}

	return toJSON(result), nil
}

// checkTimeConflict 检查时间冲突（内部方法）
func (c *CalendarManager) checkTimeConflict(startTime int64, duration int64, excludeID int64) []map[string]interface{} {
	db := c.getDB()
	if db == nil {
		return nil
	}

	endTime := startTime + duration

	var events []model.CalendarEvent
	query := db.Where("user_id = ? AND start_time < ? AND end_time > ?",
		"weimeng_xiansheng", endTime, startTime)

	if excludeID > 0 {
		query = query.And("id != ?", excludeID)
	}

	err := query.Find(&events)
	if err != nil {
		return nil
	}

	conflicts := make([]map[string]interface{}, 0, len(events))
	for _, event := range events {
		conflicts = append(conflicts, map[string]interface{}{
			"id":         event.Id,
			"title":      event.Title,
			"start_time": event.StartTime,
			"end_time":   event.EndTime,
		})
	}

	return conflicts
}

// parseDateToTimestamp 将日期字符串转换为时间戳
func parseDateToTimestamp(dateStr string) (int64, error) {
	t, err := time.ParseInLocation("2006-01-02", dateStr, time.Local)
	if err != nil {
		return 0, err
	}
	return t.Unix(), nil
}

// getBoolParam 获取布尔参数
func getBoolParam(params map[string]interface{}, key string, defaultVal bool) bool {
	if val, ok := params[key]; ok {
		if b, ok := val.(bool); ok {
			return b
		}
		if f, ok := val.(float64); ok {
			return f != 0
		}
	}
	return defaultVal
}

// AfterRun 工具调用后的后处理
func (c *CalendarManager) AfterRun(ctx context.Context, result string, err error) {
	// 简化实现，暂时不做学习集成
	// 后续可以通过接口注入学习功能
}