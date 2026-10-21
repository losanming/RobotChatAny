package intelligent

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"new-robot-self/internal/tools"
)

// TriggerManager 触发管理器
type TriggerManager struct {
	rules       map[string]*TriggerRule
	activeRules map[string]bool
	mu          sync.RWMutex
	ticker      *time.Ticker
	ctx         context.Context
	cancel      context.CancelFunc
	assistant   *IntelligentAssistant
}

func NewTriggerManager() *TriggerManager {
	ctx, cancel := context.WithCancel(context.Background())
	
	manager := &TriggerManager{
		rules:       make(map[string]*TriggerRule),
		activeRules: make(map[string]bool),
		ctx:         ctx,
		cancel:      cancel,
	}
	
	manager.initDefaultRules()
	return manager
}

// SetAssistant 设置智能助手引用
func (m *TriggerManager) SetAssistant(assistant *IntelligentAssistant) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.assistant = assistant
}

// initDefaultRules 初始化默认触发规则
func (m *TriggerManager) initDefaultRules() {
	// 时间触发规则：每小时检查一次日程提醒
	m.rules["hourly_schedule_check"] = &TriggerRule{
		ID:         "hourly_schedule_check",
		Type:       "time",
		Condition:  "every_hour",
		Action:     "check_schedule_reminders",
		Parameters: map[string]interface{}{"interval": 3600},
		Enabled:    true,
		LastTrigger: 0,
	}

	// 状态触发规则：健康指标异常
	m.rules["health_alert"] = &TriggerRule{
		ID:         "health_alert",
		Type:       "state",
		Condition:  "health_metric_abnormal",
		Action:     "send_health_alert",
		Parameters: map[string]interface{}{
			"blood_sugar_high": 7.0,
			"blood_sugar_low": 3.9,
		},
		Enabled:    true,
		LastTrigger: 0,
	}

	// 事件触发规则：任务完成提醒
	m.rules["task_completion"] = &TriggerRule{
		ID:         "task_completion",
		Type:       "event",
		Condition:  "task_completed",
		Action:     "send_completion_notification",
		Parameters: map[string]interface{}{},
		Enabled:    true,
		LastTrigger: 0,
	}
}

// Start 启动触发管理器
func (m *TriggerManager) Start() error {
	m.ticker = time.NewTicker(1 * time.Minute)
	
	go m.run()
	return nil
}

// Stop 停止触发管理器
func (m *TriggerManager) Stop() error {
	m.cancel()
	if m.ticker != nil {
		m.ticker.Stop()
	}
	return nil
}

// run 运行触发检查循环
func (m *TriggerManager) run() {
	for {
		select {
		case <-m.ctx.Done():
			return
		case <-m.ticker.C:
			m.checkTriggers()
		}
	}
}

// checkTriggers 检查所有触发规则
func (m *TriggerManager) checkTriggers() {
	m.mu.RLock()
	defer m.mu.RUnlock()

	for id, rule := range m.rules {
		if rule.Enabled && m.activeRules[id] {
			if m.shouldTrigger(rule) {
				m.executeRule(rule)
			}
		}
	}
}

// shouldTrigger 判断是否应该触发
func (m *TriggerManager) shouldTrigger(rule *TriggerRule) bool {
	now := time.Now().Unix()
	
	switch rule.Type {
	case "time":
		// 时间触发：检查时间间隔
		if interval, ok := rule.Parameters["interval"].(int); ok {
			return now-rule.LastTrigger >= int64(interval)
		}
	case "state":
		// 状态触发：检查状态条件
		return m.checkStateCondition(rule)
	case "event":
		// 事件触发：检查事件条件
		// TODO: 实现事件检查逻辑
		return false
	case "smart":
		// 智能触发：基于预测的触发
		// TODO: 实现智能触发逻辑
		return false
	default:
		return false
	}
	
	return false
}

// checkStateCondition 检查状态条件
func (m *TriggerManager) checkStateCondition(rule *TriggerRule) bool {
	if rule.Condition == "health_metric_abnormal" {
		// 检查健康指标是否异常
		// 调用health_tracker工具检查最新健康数据
		tool, exists := tools.Get("health_tracker")
		if !exists {
			return false
		}
		
		// 调用工具查询健康数据
		ctx := context.Background()
		params := map[string]interface{}{
			"action":     "alert",
			"start_date": time.Now().Add(-24 * time.Hour).Unix(),
			"end_date":   time.Now().Unix(),
		}
		
		paramsJSON, _ := json.Marshal(params)
		result, err := tool.InvokableRun(ctx, string(paramsJSON))
		if err != nil {
			return false
		}
		
		// 解析结果，检查是否有异常
		var alertData map[string]interface{}
		json.Unmarshal([]byte(result), &alertData)
		
		if alerts, ok := alertData["alerts"].([]interface{}); ok && len(alerts) > 0 {
			return true // 有异常，需要触发
		}
	}
	
	return false
}

// executeRule 执行触发规则
func (m *TriggerManager) executeRule(rule *TriggerRule) {
	m.mu.Lock()
	defer m.mu.Unlock()

	// 根据规则类型执行不同的动作
	switch rule.Action {
	case "check_schedule_reminders":
		m.checkScheduleReminders(rule)
	case "send_health_alert":
		m.sendHealthAlert(rule)
	case "send_completion_notification":
		m.sendCompletionNotification(rule)
	default:
		fmt.Printf("执行触发规则: %s, 动作: %s\n", rule.ID, rule.Action)
	}
	
	// 更新最后触发时间
	rule.LastTrigger = time.Now().Unix()
}

// checkScheduleReminders 检查日程提醒
func (m *TriggerManager) checkScheduleReminders(rule *TriggerRule) {
	// 调用smart_scheduler工具检查即将到期的任务
	tool, exists := tools.Get("smart_scheduler")
	if !exists {
		return
	}
	
	ctx := context.Background()
	params := map[string]interface{}{
		"action":     "query",
		"start_date": time.Now().Unix(),
		"end_date":   time.Now().Add(24 * time.Hour).Unix(),
	}
	
	paramsJSON, _ := json.Marshal(params)
	result, err := tool.InvokableRun(ctx, string(paramsJSON))
	if err != nil {
		fmt.Printf("检查日程提醒失败: %v\n", err)
		return
	}
	
	fmt.Printf("日程提醒检查完成: %s\n", result)
}

// sendHealthAlert 发送健康提醒
func (m *TriggerManager) sendHealthAlert(rule *TriggerRule) {
	// 调用health_tracker工具获取健康数据并检查异常
	tool, exists := tools.Get("health_tracker")
	if !exists {
		return
	}
	
	ctx := context.Background()
	params := map[string]interface{}{
		"action":     "alert",
		"start_date": time.Now().Add(-24 * time.Hour).Unix(),
		"end_date":   time.Now().Unix(),
	}
	
	paramsJSON, _ := json.Marshal(params)
	result, err := tool.InvokableRun(ctx, string(paramsJSON))
	if err != nil {
		fmt.Printf("健康提醒失败: %v\n", err)
		return
	}
	
	fmt.Printf("健康提醒已发送: %s\n", result)
}

// sendCompletionNotification 发送完成通知
func (m *TriggerManager) sendCompletionNotification(rule *TriggerRule) {
	// 调用notification_manager工具发送通知
	tool, exists := tools.Get("notification_manager")
	if !exists {
		return
	}
	
	ctx := context.Background()
	params := map[string]interface{}{
		"action":    "send",
		"title":     "任务完成通知",
		"message":   "您的任务已完成",
		"recipient": "weimeng_xiansheng",
	}
	
	paramsJSON, _ := json.Marshal(params)
	result, err := tool.InvokableRun(ctx, string(paramsJSON))
	if err != nil {
		fmt.Printf("完成通知发送失败: %v\n", err)
		return
	}
	
	fmt.Printf("完成通知已发送: %s\n", result)
}

// AddRule 添加触发规则
func (m *TriggerManager) AddRule(rule *TriggerRule) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.rules[rule.ID]; exists {
		return fmt.Errorf("rule %s already exists", rule.ID)
	}

	m.rules[rule.ID] = rule
	m.activeRules[rule.ID] = true
	
	return nil
}

// RemoveRule 移除触发规则
func (m *TriggerManager) RemoveRule(ruleID string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.rules[ruleID]; !exists {
		return fmt.Errorf("rule %s not found", ruleID)
	}

	delete(m.rules, ruleID)
	delete(m.activeRules, ruleID)
	
	return nil
}

// EnableRule 启用触发规则
func (m *TriggerManager) EnableRule(ruleID string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.rules[ruleID]; !exists {
		return fmt.Errorf("rule %s not found", ruleID)
	}

	m.activeRules[ruleID] = true
	return nil
}

// DisableRule 禁用触发规则
func (m *TriggerManager) DisableRule(ruleID string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.rules[ruleID]; !exists {
		return fmt.Errorf("rule %s not found", ruleID)
	}

	m.activeRules[ruleID] = false
	return nil
}

// GetRule 获取触发规则
func (m *TriggerManager) GetRule(ruleID string) (*TriggerRule, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	rule, exists := m.rules[ruleID]
	if !exists {
		return nil, fmt.Errorf("rule %s not found", ruleID)
	}

	return rule, nil
}

// ListRules 列出所有触发规则
func (m *TriggerManager) ListRules() map[string]*TriggerRule {
	m.mu.RLock()
	defer m.mu.RUnlock()

	// 返回副本
	rules := make(map[string]*TriggerRule)
	for id, rule := range m.rules {
		ruleCopy := *rule
		rules[id] = &ruleCopy
	}

	return rules
}

// IsRunning 检查触发管理器是否正在运行
func (m *TriggerManager) IsRunning() bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.ticker != nil
}