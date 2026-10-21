package brain

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"
)

// ActionExecutor 建议执行器
type ActionExecutor struct {
	brain         *Brain
	executionLog  []*ExecutionRecord
	mu            sync.RWMutex
}

// ExecutionRecord 执行记录
type ExecutionRecord struct {
	PredictionID  string    `json:"prediction_id"`
	Action        string    `json:"action"`
	ExecutedAt    time.Time `json:"executed_at"`
	Status        string    `json:"status"` // success, failed, pending
	Result        string    `json:"result"`
	Error         string    `json:"error,omitempty"`
}

// NewActionExecutor 创建建议执行器
func NewActionExecutor(brain *Brain) *ActionExecutor {
	return &ActionExecutor{
		brain:        brain,
		executionLog: make([]*ExecutionRecord, 0),
	}
}

// ExecutePrediction 执行预测建议
func (e *ActionExecutor) ExecutePrediction(ctx context.Context, prediction *Prediction) (*ExecutionRecord, error) {
	record := &ExecutionRecord{
		Action:     prediction.Action,
		ExecutedAt: time.Now(),
		Status:     "pending",
	}

	// 根据不同的动作类型执行不同的逻辑
	switch prediction.Action {
	case "exercise_reminder":
		record.Status, record.Result = e.executeExerciseReminder(ctx, prediction)
	case "blood_sugar_reminder":
		record.Status, record.Result = e.executeBloodSugarReminder(ctx, prediction)
	case "diet_advice":
		record.Status, record.Result = e.executeDietAdvice(ctx, prediction)
	case "weight_management":
		record.Status, record.Result = e.executeWeightManagement(ctx, prediction)
	case "focus_time_reminder":
		record.Status, record.Result = e.executeFocusTimeReminder(ctx, prediction)
	case "break_reminder":
		record.Status, record.Result = e.executeBreakReminder(ctx, prediction)
	case "task_planning":
		record.Status, record.Result = e.executeTaskPlanning(ctx, prediction)
	case "exercise_suggestion":
		record.Status, record.Result = e.executeExerciseSuggestion(ctx, prediction)
	default:
		record.Status = "skipped"
		record.Result = fmt.Sprintf("未知的动作类型: %s", prediction.Action)
	}

	// 记录执行结果
	e.mu.Lock()
	e.executionLog = append(e.executionLog, record)
	e.mu.Unlock()

	log.Printf("执行建议 [%s]: 状态=%s, 结果=%s", prediction.Action, record.Status, record.Result)
	
	return record, nil
}

// executeExerciseReminder 执行运动提醒
func (e *ActionExecutor) executeExerciseReminder(ctx context.Context, prediction *Prediction) (status, result string) {
	// 这里可以集成通知系统，比如发送到飞书、微信等
	// 目前先记录日志
	return "success", fmt.Sprintf("已发送运动提醒: %s", prediction.Content)
}

// executeBloodSugarReminder 执行血糖监测提醒
func (e *ActionExecutor) executeBloodSugarReminder(ctx context.Context, prediction *Prediction) (status, result string) {
	// 可以提醒用户测量血糖，并记录测量结果
	return "success", fmt.Sprintf("已发送血糖监测提醒: %s", prediction.Content)
}

// executeDietAdvice 执行饮食建议
func (e *ActionExecutor) executeDietAdvice(ctx context.Context, prediction *Prediction) (status, result string) {
	return "success", fmt.Sprintf("已提供饮食建议: %s", prediction.Content)
}

// executeWeightManagement 执行体重管理
func (e *ActionExecutor) executeWeightManagement(ctx context.Context, prediction *Prediction) (status, result string) {
	return "success", fmt.Sprintf("已提供体重管理建议: %s", prediction.Content)
}

// executeFocusTimeReminder 执行专注时间提醒
func (e *ActionExecutor) executeFocusTimeReminder(ctx context.Context, prediction *Prediction) (status, result string) {
	return "success", fmt.Sprintf("已发送专注时间提醒: %s", prediction.Content)
}

// executeBreakReminder 执行休息提醒
func (e *ActionExecutor) executeBreakReminder(ctx context.Context, prediction *Prediction) (status, result string) {
	return "success", fmt.Sprintf("已发送休息提醒: %s", prediction.Content)
}

// executeTaskPlanning 执行任务规划
func (e *ActionExecutor) executeTaskPlanning(ctx context.Context, prediction *Prediction) (status, result string) {
	return "success", fmt.Sprintf("已提供任务规划建议: %s", prediction.Content)
}

// executeExerciseSuggestion 执行运动建议
func (e *ActionExecutor) executeExerciseSuggestion(ctx context.Context, prediction *Prediction) (status, result string) {
	return "success", fmt.Sprintf("已提供运动建议: %s", prediction.Content)
}

// GetExecutionLog 获取执行日志
func (e *ActionExecutor) GetExecutionLog() []*ExecutionRecord {
	e.mu.RLock()
	defer e.mu.RUnlock()
	
	// 返回副本
	logCopy := make([]*ExecutionRecord, len(e.executionLog))
	copy(logCopy, e.executionLog)
	return logCopy
}

// ExecuteHighPriorityPredictions 执行高优先级的预测建议
func (e *ActionExecutor) ExecuteHighPriorityPredictions(ctx context.Context, predictions []*Prediction, minPriority int) ([]*ExecutionRecord, error) {
	results := make([]*ExecutionRecord, 0)
	
	for _, prediction := range predictions {
		if prediction.Priority >= minPriority {
			record, err := e.ExecutePrediction(ctx, prediction)
			if err != nil {
				log.Printf("执行建议失败: %v", err)
				continue
			}
			results = append(results, record)
		}
	}
	
	return results, nil
}
