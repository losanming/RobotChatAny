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

// TodoItemInternal 待办事项内部结构（用于内存操作）
type TodoItemInternal struct {
	ID          string   `json:"id"`
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Priority    string   `json:"priority"` // high | medium | low
	Status      string   `json:"status"`   // pending | in_progress | completed | cancelled
	Tags        []string `json:"tags"`
	DueDate     int64    `json:"due_date"` // 截止日期时间戳
	CompletedAt int64    `json:"completed_at"`
	CreatedAt   int64    `json:"created_at"`
	UpdatedAt   int64    `json:"updated_at"`
}

// TodoManager 待办事项管理工具
type TodoManager struct {
	enableLearning bool
	category       string
	toolName       string
	db             *xorm.Engine // 数据库实例
}

// NewTodoManager 创建待办事项管理工具
func NewTodoManager() *TodoManager {
	return &TodoManager{
		enableLearning: true,
		category:       "personal",
		toolName:       "todo_manager",
	}
}

// SetDB 设置数据库实例
func (t *TodoManager) SetDB(db *xorm.Engine) {
	t.db = db
}

// getDB 获取数据库实例
func (t *TodoManager) getDB() *xorm.Engine {
	if t.db == nil {
		return public.DB
	}
	return t.db
}

// Info 返回工具信息
func (t *TodoManager) Info(ctx context.Context) (*schema.ToolInfo, error) {
	return &schema.ToolInfo{
		Name: "todo_manager",
		Desc: "待办事项管理工具，支持创建、查询、更新、删除待办事项，优先级设置，标签分类，完成度跟踪",
		ParamsOneOf: schema.NewParamsOneOfByParams(map[string]*schema.ParameterInfo{
			"action": {
				Desc:     "操作类型",
				Type:     schema.String,
				Required: true,
				Enum:     []string{"create", "query", "update", "delete", "list", "complete", "search"},
			},
			"id": {
				Desc:     "待办事项ID（update/delete/query/complete时需要）",
				Type:     schema.String,
				Required: false,
			},
			"title": {
				Desc:     "待办事项标题（create时需要）",
				Type:     schema.String,
				Required: false,
			},
			"description": {
				Desc:     "待办事项描述",
				Type:     schema.String,
				Required: false,
			},
			"priority": {
				Desc:     "优先级",
				Type:     schema.String,
				Required: false,
				Enum:     []string{"high", "medium", "low"},
			},
			"status": {
				Desc:     "状态",
				Type:     schema.String,
				Required: false,
				Enum:     []string{"pending", "in_progress", "completed", "cancelled"},
			},
			"tags": {
				Desc:     "标签列表",
				Type:     schema.Array,
				Required: false,
			},
			"due_date": {
				Desc:     "截止日期时间戳",
				Type:     schema.Integer,
				Required: false,
			},
			"filter_status": {
				Desc:     "筛选状态（list时使用）",
				Type:     schema.String,
				Required: false,
				Enum:     []string{"pending", "in_progress", "completed", "cancelled", "all"},
			},
			"filter_priority": {
				Desc:     "筛选优先级（list时使用）",
				Type:     schema.String,
				Required: false,
				Enum:     []string{"high", "medium", "low", "all"},
			},
			"search_keyword": {
				Desc:     "搜索关键词（search时使用）",
				Type:     schema.String,
				Required: false,
			},
		}),
	}, nil
}

// InvokableRun 执行工具调用
func (t *TodoManager) InvokableRun(ctx context.Context, args string, opts ...tool.Option) (string, error) {
	var params map[string]interface{}
	if err := json.Unmarshal([]byte(args), &params); err != nil {
		return "", fmt.Errorf("参数解析失败: %w", err)
	}

	action := params["action"].(string)
	switch action {
	case "create":
		return t.createTodo(params)
	case "query":
		return t.queryTodo(params)
	case "update":
		return t.updateTodo(params)
	case "delete":
		return t.deleteTodo(params)
	case "list":
		return t.listTodos(params)
	case "complete":
		return t.completeTodo(params)
	case "search":
		return t.searchTodos(params)
	default:
		return "", fmt.Errorf("不支持的操作类型: %s", action)
	}
}

// createTodo 创建待办事项
func (t *TodoManager) createTodo(params map[string]interface{}) (string, error) {
	title := getStringParam(params, "title", "")
	if title == "" {
		return "", fmt.Errorf("标题不能为空")
	}

	db := t.getDB()
	if db == nil {
		return "", fmt.Errorf("database not initialized")
	}

	priority := getStringParam(params, "priority", "medium")
	status := getStringParam(params, "status", "pending")
	
	// 验证优先级
	if priority != "high" && priority != "medium" && priority != "low" {
		priority = "medium"
	}
	
	// 验证状态
	if status != "pending" && status != "in_progress" && status != "completed" && status != "cancelled" {
		status = "pending"
	}

	tags := getStringArrayParam(params, "tags")
	tagsJSON, _ := json.Marshal(tags)
	
	dueDate := getInt64Param(params, "due_date", 0)

	todo := &model.TodoItem{
		UserID:      "weimeng_xiansheng",
		Title:       title,
		Description: getStringParam(params, "description", ""),
		Priority:    priority,
		Status:      status,
		Tags:        string(tagsJSON),
		DueDate:     dueDate,
	}

	_, err := db.Insert(todo)
	if err != nil {
		return "", fmt.Errorf("failed to create todo: %w", err)
	}

	result := map[string]interface{}{
		"success": true,
		"todo": map[string]interface{}{
			"id":          todo.Id,
			"title":       todo.Title,
			"description": todo.Description,
			"priority":    todo.Priority,
			"status":      todo.Status,
			"tags":        tags,
			"due_date":    todo.DueDate,
			"created_at":  todo.CreatedAt,
		},
		"message": "待办事项创建成功",
	}

	return toJSON(result), nil
}

// queryTodo 查询待办事项
func (t *TodoManager) queryTodo(params map[string]interface{}) (string, error) {
	id := getInt64Param(params, "id", 0)
	if id == 0 {
		return "", fmt.Errorf("待办事项ID不能为空")
	}

	db := t.getDB()
	if db == nil {
		return "", fmt.Errorf("database not initialized")
	}

	var todo model.TodoItem
	has, err := db.ID(id).Where("user_id = ?", "weimeng_xiansheng").Get(&todo)
	if err != nil {
		return "", fmt.Errorf("failed to query todo: %w", err)
	}
	if !has {
		return "", fmt.Errorf("待办事项不存在: %d", id)
	}

	var tags []string
	json.Unmarshal([]byte(todo.Tags), &tags)

	result := map[string]interface{}{
		"success": true,
		"todo": map[string]interface{}{
			"id":          todo.Id,
			"title":       todo.Title,
			"description": todo.Description,
			"priority":    todo.Priority,
			"status":      todo.Status,
			"tags":        tags,
			"due_date":    todo.DueDate,
			"completed_at": todo.CompletedAt,
			"created_at":  todo.CreatedAt,
			"updated_at":  todo.UpdatedAt,
		},
	}

	return toJSON(result), nil
}

// updateTodo 更新待办事项
func (t *TodoManager) updateTodo(params map[string]interface{}) (string, error) {
	id := getInt64Param(params, "id", 0)
	if id == 0 {
		return "", fmt.Errorf("待办事项ID不能为空")
	}

	db := t.getDB()
	if db == nil {
		return "", fmt.Errorf("database not initialized")
	}

	var todo model.TodoItem
	has, err := db.ID(id).Where("user_id = ?", "weimeng_xiansheng").Get(&todo)
	if err != nil {
		return "", fmt.Errorf("failed to query todo: %w", err)
	}
	if !has {
		return "", fmt.Errorf("待办事项不存在: %d", id)
	}

	// 构建更新字段
	updates := make(map[string]interface{})
	
	if title := getStringParam(params, "title", ""); title != "" {
		updates["title"] = title
	}
	if desc := getStringParam(params, "description", ""); desc != "" {
		updates["description"] = desc
	}
	
	priority := getStringParam(params, "priority", "")
	if priority == "high" || priority == "medium" || priority == "low" {
		updates["priority"] = priority
	}
	
	status := getStringParam(params, "status", "")
	if status == "pending" || status == "in_progress" || status == "completed" || status == "cancelled" {
		updates["status"] = status
		if status == "completed" && todo.Status != "completed" {
			updates["completed_at"] = time.Now().Unix()
		}
	}
	
	tags := getStringArrayParam(params, "tags")
	if len(tags) > 0 {
		tagsJSON, _ := json.Marshal(tags)
		updates["tags"] = string(tagsJSON)
	}
	
	dueDate := getInt64Param(params, "due_date", 0)
	if dueDate > 0 {
		updates["due_date"] = dueDate
	}

	if len(updates) == 0 {
		return "", fmt.Errorf("没有需要更新的字段")
	}

	_, err = db.ID(id).Update(&todo, updates)
	if err != nil {
		return "", fmt.Errorf("failed to update todo: %w", err)
	}

	// 重新查询获取最新数据
	db.ID(id).Get(&todo)
	
	var tagList []string
	json.Unmarshal([]byte(todo.Tags), &tagList)

	result := map[string]interface{}{
		"success": true,
		"todo": map[string]interface{}{
			"id":          todo.Id,
			"title":       todo.Title,
			"description": todo.Description,
			"priority":    todo.Priority,
			"status":      todo.Status,
			"tags":        tagList,
			"due_date":    todo.DueDate,
			"completed_at": todo.CompletedAt,
			"updated_at":  todo.UpdatedAt,
		},
		"message": "待办事项更新成功",
	}

	return toJSON(result), nil
}

// deleteTodo 删除待办事项
func (t *TodoManager) deleteTodo(params map[string]interface{}) (string, error) {
	id := getInt64Param(params, "id", 0)
	if id == 0 {
		return "", fmt.Errorf("待办事项ID不能为空")
	}

	db := t.getDB()
	if db == nil {
		return "", fmt.Errorf("database not initialized")
	}

	_, err := db.Where("id = ? AND user_id = ?", id, "weimeng_xiansheng").Delete(&model.TodoItem{})
	if err != nil {
		return "", fmt.Errorf("failed to delete todo: %w", err)
	}

	result := map[string]interface{}{
		"success": true,
		"message": "待办事项删除成功",
	}

	return toJSON(result), nil
}

// listTodos 列出待办事项
func (t *TodoManager) listTodos(params map[string]interface{}) (string, error) {
	db := t.getDB()
	if db == nil {
		return "", fmt.Errorf("database not initialized")
	}

	filterStatus := getStringParam(params, "filter_status", "all")
	filterPriority := getStringParam(params, "filter_priority", "all")

	query := db.Where("user_id = ?", "weimeng_xiansheng")
	
	if filterStatus != "all" {
		query = query.And("status = ?", filterStatus)
	}
	if filterPriority != "all" {
		query = query.And("priority = ?", filterPriority)
	}

	var todos []model.TodoItem
	err := query.OrderBy("CASE priority WHEN 'high' THEN 1 WHEN 'medium' THEN 2 WHEN 'low' THEN 3 END, due_date ASC").Find(&todos)
	if err != nil {
		return "", fmt.Errorf("failed to list todos: %w", err)
	}

	resultTodos := make([]map[string]interface{}, 0, len(todos))
	for _, todo := range todos {
		var tags []string
		json.Unmarshal([]byte(todo.Tags), &tags)
		
		resultTodos = append(resultTodos, map[string]interface{}{
			"id":          todo.Id,
			"title":       todo.Title,
			"description": todo.Description,
			"priority":    todo.Priority,
			"status":      todo.Status,
			"tags":        tags,
			"due_date":    todo.DueDate,
			"created_at":  todo.CreatedAt,
		})
	}

	result := map[string]interface{}{
		"success": true,
		"todos":   resultTodos,
		"total":   len(todos),
	}

	return toJSON(result), nil
}

// completeTodo 完成待办事项
func (t *TodoManager) completeTodo(params map[string]interface{}) (string, error) {
	id := getInt64Param(params, "id", 0)
	if id == 0 {
		return "", fmt.Errorf("待办事项ID不能为空")
	}

	db := t.getDB()
	if db == nil {
		return "", fmt.Errorf("database not initialized")
	}

	var todo model.TodoItem
	has, err := db.ID(id).Where("user_id = ?", "weimeng_xiansheng").Get(&todo)
	if err != nil {
		return "", fmt.Errorf("failed to query todo: %w", err)
	}
	if !has {
		return "", fmt.Errorf("待办事项不存在: %d", id)
	}

	if todo.Status == "completed" {
		return "", fmt.Errorf("待办事项已完成")
	}

	_, err = db.ID(id).Cols("status", "completed_at").Update(&model.TodoItem{
		Status:      "completed",
		CompletedAt: time.Now().Unix(),
	})
	if err != nil {
		return "", fmt.Errorf("failed to complete todo: %w", err)
	}

	result := map[string]interface{}{
		"success": true,
		"todo": map[string]interface{}{
			"id":          todo.Id,
			"title":       todo.Title,
			"status":      "completed",
			"completed_at": time.Now().Unix(),
		},
		"message": "待办事项已完成",
	}

	return toJSON(result), nil
}

// searchTodos 搜索待办事项
func (t *TodoManager) searchTodos(params map[string]interface{}) (string, error) {
	keyword := getStringParam(params, "search_keyword", "")
	if keyword == "" {
		return "", fmt.Errorf("搜索关键词不能为空")
	}

	db := t.getDB()
	if db == nil {
		return "", fmt.Errorf("database not initialized")
	}

	var todos []model.TodoItem
	err := db.Where("user_id = ? AND (title LIKE ? OR description LIKE ?)",
		"weimeng_xiansheng", "%"+keyword+"%", "%"+keyword+"%").
		Find(&todos)
	if err != nil {
		return "", fmt.Errorf("failed to search todos: %w", err)
	}

	resultTodos := make([]map[string]interface{}, 0, len(todos))
	for _, todo := range todos {
		var tags []string
		json.Unmarshal([]byte(todo.Tags), &tags)
		
		resultTodos = append(resultTodos, map[string]interface{}{
			"id":          todo.Id,
			"title":       todo.Title,
			"description": todo.Description,
			"priority":    todo.Priority,
			"status":      todo.Status,
			"tags":        tags,
			"due_date":    todo.DueDate,
		})
	}

	result := map[string]interface{}{
		"success": true,
		"todos":   resultTodos,
		"total":   len(todos),
		"keyword": keyword,
	}

	return toJSON(result), nil
}

// AfterRun 工具调用后的后处理
func (t *TodoManager) AfterRun(ctx context.Context, result string, err error) {
	// 简化实现，暂时不做学习集成
}

// GenerateUUID 生成UUID
func GenerateUUID() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}
