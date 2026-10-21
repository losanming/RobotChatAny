package intelligent

import (
	"context"
	"encoding/json"
	"fmt"

	"new-robot-self/internal/tools"
)

// ToolCoordinator 工具协调器
type ToolCoordinator struct {
	toolRegistry *tools.ToolRegistry
	planners      map[string]IntentPlanner
}

// ToolHandler 工具处理器接口
type ToolHandler interface {
	CanHandle(toolName string) bool
	Execute(ctx context.Context, toolName string, params map[string]interface{}) (map[string]interface{}, error)
}

// IntentPlanner 意图规划器接口
type IntentPlanner interface {
	Plan(intent *Intent) (*ToolCallPlan, error)
}

func NewToolCoordinator() *ToolCoordinator {
	coordinator := &ToolCoordinator{
		toolRegistry: tools.GetGlobalRegistry(),
		planners:      make(map[string]IntentPlanner),
	}
	coordinator.initPlanners()
	return coordinator
}

// initPlanners 初始化意图规划器
func (c *ToolCoordinator) initPlanners() {
	// 注册各种意图的规划器
	c.planners["schedule"] = &ScheduleIntentPlanner{}
	c.planners["navigation"] = &NavigationIntentPlanner{}
	c.planners["health"] = &HealthIntentPlanner{}
	c.planners["work"] = &WorkIntentPlanner{}
	c.planners["info"] = &InfoIntentPlanner{}
	c.planners["life"] = &LifeIntentPlanner{}
}

// Plan 根据意图规划工具调用
func (c *ToolCoordinator) Plan(ctx context.Context, intent *Intent) (*ToolCallPlan, error) {
	if planner, ok := c.planners[intent.Type]; ok {
		return planner.Plan(intent)
	}

	// 默认规划器
	return c.defaultPlan(intent)
}

// Execute 执行工具调用计划
func (c *ToolCoordinator) Execute(ctx context.Context, plan *ToolCallPlan) (map[string]interface{}, error) {
	results := make(map[string]interface{})
	var lastError error

	switch plan.Strategy {
	case "sequential":
		results, lastError = c.executeSequential(ctx, plan.Steps)
	case "parallel":
		results, lastError = c.executeParallel(ctx, plan.Steps)
	default:
		results, lastError = c.executeSequential(ctx, plan.Steps)
	}

	if lastError != nil {
		return results, fmt.Errorf("tool execution failed: %w", lastError)
	}

	return results, nil
}

// executeSequential 串行执行工具调用
func (c *ToolCoordinator) executeSequential(ctx context.Context, steps []ToolCallStep) (map[string]interface{}, error) {
	results := make(map[string]interface{})

	for i, step := range steps {
		// 检查依赖条件
		if !c.checkDependencies(step, results) {
			continue
		}

		// 检查调用条件
		if !c.checkCondition(step, results) {
			continue
		}

		// 执行工具调用
		result, err := c.executeTool(ctx, step)
		if err != nil {
			// 记录错误但继续执行后续步骤
			results[fmt.Sprintf("step_%d_error", i)] = err.Error()
			continue
		}

		results[fmt.Sprintf("step_%d", i)] = result
	}

	return results, nil
}

// executeParallel 并行执行工具调用
func (c *ToolCoordinator) executeParallel(ctx context.Context, steps []ToolCallStep) (map[string]interface{}, error) {
	results := make(map[string]interface{})
	resultChan := make(chan struct {
		index int
		result map[string]interface{}
		err    error
	}, len(steps))

	// 并行执行
	for i, step := range steps {
		go func(idx int, s ToolCallStep) {
			result, err := c.executeTool(ctx, s)
			resultChan <- struct {
				index int
				result map[string]interface{}
				err    error
			}{idx, result, err}
		}(i, step)
	}

	// 收集结果
	for i := 0; i < len(steps); i++ {
		res := <-resultChan
		if res.err != nil {
			results[fmt.Sprintf("step_%d_error", res.index)] = res.err.Error()
		} else {
			results[fmt.Sprintf("step_%d", res.index)] = res.result
		}
	}

	return results, nil
}

// executeTool 执行单个工具调用
func (c *ToolCoordinator) executeTool(ctx context.Context, step ToolCallStep) (map[string]interface{}, error) {
	// 从工具注册中心获取工具
	tool, exists := c.toolRegistry.Get(step.ToolName)
	if !exists {
		return nil, fmt.Errorf("tool %s not found", step.ToolName)
	}

	// 将参数转换为JSON字符串
	paramsJSON, err := json.Marshal(step.Parameters)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal parameters: %w", err)
	}

	// 调用真实的工具
	result, err := tool.InvokableRun(ctx, string(paramsJSON))
	if err != nil {
		return nil, fmt.Errorf("tool execution failed: %w", err)
	}

	// 解析工具返回结果
	var resultData map[string]interface{}
	if err := json.Unmarshal([]byte(result), &resultData); err != nil {
		// 如果解析失败，返回原始字符串
		return map[string]interface{}{
			"raw_result": result,
			"status":     "success",
		}, nil
	}

	return resultData, nil
}

// checkDependencies 检查依赖条件
func (c *ToolCoordinator) checkDependencies(step ToolCallStep, results map[string]interface{}) bool {
	if len(step.DependsOn) == 0 {
		return true
	}

	for _, dep := range step.DependsOn {
		if _, ok := results[dep]; !ok {
			return false
		}
	}

	return true
}

// checkCondition 检查调用条件
func (c *ToolCoordinator) checkCondition(step ToolCallStep, results map[string]interface{}) bool {
	if step.Condition == "" {
		return true
	}

	// TODO: 实现条件检查逻辑
	return true
}

// defaultPlan 默认规划器
func (c *ToolCoordinator) defaultPlan(intent *Intent) (*ToolCallPlan, error) {
	return &ToolCallPlan{
		Intent:         *intent,
		Steps:          []ToolCallStep{},
		Strategy:       "sequential",
		ExpectedOutput: "text",
		CreatedAt:      0,
	}, nil
}

// ScheduleIntentPlanner 日程意图规划器
type ScheduleIntentPlanner struct{}

func (p *ScheduleIntentPlanner) Plan(intent *Intent) (*ToolCallPlan, error) {
	steps := []ToolCallStep{
		{
			ToolName:   "smart_scheduler",
			Parameters: map[string]interface{}{
				"action":     "create",
				"task_title": intent.Entities["event"],
				"task_type":  "life",
			},
			Priority: 1,
		},
		{
			ToolName:   "smart_reminder",
			Parameters: map[string]interface{}{
				"action":        "create",
				"reminder_type": "schedule",
				"title":         "日程提醒",
				"scheduled_at":  intent.Entities["time"],
			},
			DependsOn: []string{"step_0"},
			Priority: 2,
		},
	}

	return &ToolCallPlan{
		Intent:         *intent,
		Steps:          steps,
		Strategy:       "sequential",
		ExpectedOutput: "schedule_created",
	}, nil
}

// NavigationIntentPlanner 导航意图规划器
type NavigationIntentPlanner struct{}

func (p *NavigationIntentPlanner) Plan(intent *Intent) (*ToolCallPlan, error) {
	location, _ := intent.Entities["location"].(string)

	steps := []ToolCallStep{
		{
			ToolName:   "navigation",
			Parameters: map[string]interface{}{
				"action":      "plan_route",
				"destination": location,
				"mode":        "driving",
			},
			Priority: 1,
		},
		{
			ToolName:   "around_search",
			Parameters: map[string]interface{}{
				"location": location,
				"radius":   1000,
				"types":    []string{"food", "shopping"},
			},
			Priority: 2,
		},
		{
			ToolName:   "weather_query",
			Parameters: map[string]interface{}{
				"location": location,
			},
			Priority: 3,
		},
	}

	return &ToolCallPlan{
		Intent:         *intent,
		Steps:          steps,
		Strategy:       "parallel",
		ExpectedOutput: "navigation_info",
	}, nil
}

// HealthIntentPlanner 健康意图规划器
type HealthIntentPlanner struct{}

func (p *HealthIntentPlanner) Plan(intent *Intent) (*ToolCallPlan, error) {
	food, _ := intent.Entities["food"].(string)
	quantity, _ := intent.Entities["quantity"].(string)

	steps := []ToolCallStep{
		{
			ToolName:   "diet_manager",
			Parameters: map[string]interface{}{
				"action":       "calculate",
				"food_name":    food,
				"amount":       quantity,
				"amount_unit":  "g",
			},
			Priority: 1,
		},
		{
			ToolName:   "health_tracker",
			Parameters: map[string]interface{}{
				"action":      "record",
				"record_type": "diet",
				"food_name":   food,
				"amount":      quantity,
			},
			DependsOn: []string{"step_0"},
			Priority: 2,
		},
	}

	return &ToolCallPlan{
		Intent:         *intent,
		Steps:          steps,
		Strategy:       "sequential",
		ExpectedOutput: "nutrition_info",
	}, nil
}

// WorkIntentPlanner 工作意图规划器
type WorkIntentPlanner struct{}

func (p *WorkIntentPlanner) Plan(intent *Intent) (*ToolCallPlan, error) {
	searchTerm, _ := intent.Entities["search_term"].(string)
	task, _ := intent.Entities["task"].(string)

	steps := []ToolCallStep{
		{
			ToolName:   "work_assistant",
			Parameters: map[string]interface{}{
				"action":       "search_doc",
				"search_query": searchTerm,
				"doc_type":     "tutorial",
			},
			Priority: 1,
		},
		{
			ToolName:   "smart_search",
			Parameters: map[string]interface{}{
				"query": task,
			},
			Priority: 2,
		},
	}

	return &ToolCallPlan{
		Intent:         *intent,
		Steps:          steps,
		Strategy:       "parallel",
		ExpectedOutput: "work_assistance",
	}, nil
}

// InfoIntentPlanner 信息意图规划器
type InfoIntentPlanner struct{}

func (p *InfoIntentPlanner) Plan(intent *Intent) (*ToolCallPlan, error) {
	location, _ := intent.Entities["location"].(string)
	query, _ := intent.Entities["query"].(string)

	steps := []ToolCallStep{
		{
			ToolName:   "weather_query",
			Parameters: map[string]interface{}{
				"location": location,
			},
			Priority: 1,
		},
		{
			ToolName:   "smart_search",
			Parameters: map[string]interface{}{
				"query": query,
			},
			Priority: 2,
		},
	}

	return &ToolCallPlan{
		Intent:         *intent,
		Steps:          steps,
		Strategy:       "parallel",
		ExpectedOutput: "info_result",
	}, nil
}

// LifeIntentPlanner 生活意图规划器
type LifeIntentPlanner struct{}

func (p *LifeIntentPlanner) Plan(intent *Intent) (*ToolCallPlan, error) {
	location, _ := intent.Entities["location"].(string)

	steps := []ToolCallStep{
		{
			ToolName:   "navigation",
			Parameters: map[string]interface{}{
				"action":      "plan_route",
				"destination": location,
			},
			Priority: 1,
		},
		{
			ToolName:   "around_search",
			Parameters: map[string]interface{}{
				"location": location,
				"radius":   1000,
				"types":    []string{"shopping"},
			},
			Priority: 2,
		},
	}

	return &ToolCallPlan{
		Intent:         *intent,
		Steps:          steps,
		Strategy:       "sequential",
		ExpectedOutput: "life_assistance",
	}, nil
}