package brain

import (
	"context"
	"log"
	"math"
	"sort"
	"strings"
	"sync"
	"time"

	"new-robot-self/internal/awareness"
)

// Brain 负责维护 WorldState，并基于多模态线索生成目的假设与 Prompt
type Brain struct {
	state               *WorldState
	mu                  sync.RWMutex
	predictionEngine    *PredictionEngine
	actionExecutor      *ActionExecutor
	learningEngine      *InformationLearningEngine
}

func NewBrain() *Brain {
	brain := &Brain{state: NewWorldState()}
	brain.predictionEngine = NewPredictionEngine(brain)
	brain.actionExecutor = NewActionExecutor(brain)
	
	// 初始化学习引擎（可选）
	// brain.learningEngine = ... 需要配置
	
	return brain
}

// InitLearningEngine 初始化学习引擎
func (b *Brain) InitLearningEngine(config LearningConfig) error {
	learningEngine, err := NewInformationLearningEngine(b, config)
	if err != nil {
		return err
	}
	
	b.learningEngine = learningEngine
	return nil
}

// GetLearningEngine 获取学习引擎
func (b *Brain) GetLearningEngine() *InformationLearningEngine {
	return b.learningEngine
}

// PredictNextActions 预测下一步可能的行动
func (b *Brain) PredictNextActions(ctx context.Context) ([]*Prediction, error) {
	return b.PredictNextActionsWithOptions(ctx, false, 3)
}

// PredictNextActionsWithOptions 预测下一步可能的行动，并执行高优先级建议
func (b *Brain) PredictNextActionsWithOptions(ctx context.Context, autoExecute bool, minPriority int) ([]*Prediction, error) {
	// 使用 PredictionEngine 进行预测
	predictions, err := b.predictionEngine.PredictNextActions(ctx)
	if err != nil {
		return nil, err
	}

	// 如果启用自动执行，执行高优先级的建议
	if autoExecute && len(predictions) > 0 {
		_, err := b.actionExecutor.ExecuteHighPriorityPredictions(ctx, predictions, minPriority)
		if err != nil {
			log.Printf("自动执行建议失败: %v", err)
		}
	}

	return predictions, nil
}

// UpdateFromAwareness 增量更新环境感知，并做一次最小目的规则推断
func (b *Brain) UpdateFromAwareness(ctx context.Context, a awareness.Context) {
	b.mu.Lock()
	defer b.mu.Unlock()
	
	// 更新 WorldState
	b.state.UpdatedTS = time.Now().Unix()
	if b.state.Env == nil {
		b.state.Env = &Environment{}
	}
	b.state.Env.Awareness = a

	b.ruleInferGoalsFromAwareness(a)
}

// UpdateFromDialog 记录对话并更新目的（极简）
func (b *Brain) UpdateFromDialog(ctx context.Context, role, text string, ts int64) {
	b.mu.Lock()
	defer b.mu.Unlock()
	if b.state.Dialog == nil {
		b.state.Dialog = &DialogContext{}
	}
	b.state.Dialog.History = append(b.state.Dialog.History, DialogTurn{TS: ts, Role: role, Content: text})
	b.state.UpdatedTS = time.Now().Unix()
	// 简单关键词：如果用户明确表达“我要…/我想…”，提升显式任务 or 假设
	lt := strings.TrimSpace(text)
	if role == "user" && (strings.Contains(lt, "我要") || strings.Contains(lt, "我想")) {
		b.pushOrBoostHypothesis("explicit_intent", lt, 0.85, "dialog")
		b.state.Dialog.LastUserIntent = lt
	}
}

// GetPrompt 生成给大模型的决策 Prompt（整合 TopK 目的与环境要点）
func (b *Brain) GetPrompt(ctx context.Context, topK int) string {
	b.mu.RLock()
	defer b.mu.RUnlock()
	if topK <= 0 {
	}
	return BuildBrainPrompt(b.state, topK)
}

// 内部：基于 awareness 的简易规则推断（全局统一秒级时间戳）
func (b *Brain) ruleInferGoalsFromAwareness(a awareness.Context) {
	// 时间段
	hour := time.Unix(a.TS, 0).In(time.Local).Hour()
	// 位置/周边
	var poiTypes []string
	if a.Around != nil {
		for _, p := range a.Around.POIs {
			if p.Type != "" {
				poiTypes = append(poiTypes, p.Type)
			}
		}
	}
	// 天气
	weather := ""
	if a.Weather != nil && len(a.Weather.Lives) > 0 {
		weather = strings.TrimSpace(a.Weather.Lives[0].Weather)
	}

	// 规则1：咖啡/餐饮附近 + 工作时段 => 可能办公/写代码
	if containsTypePrefix(poiTypes, "050") || strings.Contains(strings.Join(poiTypes, ","), "咖啡") {
		if hour >= 9 && hour <= 18 {
			b.pushOrBoostHypothesis("work_focus", "附近有餐饮/咖啡, 且处于工作时段, 可能在办公/编程", 0.6, "awareness")
		}
	}
	// 规则2：健身相关 => 可能在健身或计划运动
	if strings.Contains(strings.Join(poiTypes, ","), "健身") || containsTypePrefix(poiTypes, "080") {
		b.pushOrBoostHypothesis("fitness", "周边健身相关POI, 可能健身/运动计划", 0.55, "awareness")
	}
	// 规则3：下雨/降水 => 户外活动意愿下降，偏室内
	if weather != "" && (strings.Contains(weather, "雨") || strings.Contains(strings.ToLower(weather), "rain")) {
		b.pushOrBoostHypothesis("indoor_preference", "天气有降水, 倾向室内活动/远程办公", 0.5, "weather")
	}
	// 规则4：清晨/傍晚通勤时段
	if hour >= 7 && hour <= 9 {
		b.pushOrBoostHypothesis("commute_morning", "晨间通勤时段", 0.4, "time")
	}
	if hour >= 17 && hour <= 19 {
		b.pushOrBoostHypothesis("commute_evening", "晚间通勤时段", 0.4, "time")
	}

	b.decayHypotheses()
	b.sortHypotheses()
}

// AddOrUpdateTodo 新增或更新一条待办（按 ID upsert）
func (b *Brain) AddOrUpdateTodo(item TodoItem) {
	b.mu.Lock()
	defer b.mu.Unlock()
	if b.state == nil {
		b.state = NewWorldState()
	}
	if b.state.Todos == nil {
		b.state.Todos = &TodosState{}
	}
	for i := range b.state.Todos.Items {
		if b.state.Todos.Items[i].ID == item.ID {
			b.state.Todos.Items[i] = item
			b.state.UpdatedTS = time.Now().Unix()
			return
		}
	}
	b.state.Todos.Items = append(b.state.Todos.Items, item)
	b.state.UpdatedTS = time.Now().Unix()
}

func containsTypePrefix(types []string, prefix string) bool {
	for _, t := range types {
		if strings.HasPrefix(t, prefix) {
			return true
		}
	}
	return false
}

// UpdateTodoStatus 更新待办状态，返回是否成功
func (b *Brain) UpdateTodoStatus(id, status string) bool {
	if b.state == nil || b.state.Todos == nil {
		return false
	}
	for i := range b.state.Todos.Items {
		if b.state.Todos.Items[i].ID == id {
			b.state.Todos.Items[i].Status = status
			b.state.UpdatedTS = time.Now().Unix()
			return true
		}
	}
	return false
}

// RemoveTodo 按 ID 删除待办，返回是否删除
func (b *Brain) RemoveTodo(id string) bool {
	if b.state == nil || b.state.Todos == nil {
		return false
	}
	items := b.state.Todos.Items
	for i := range items {
		if items[i].ID == id {
			b.state.Todos.Items = append(items[:i], items[i+1:]...)
			b.state.UpdatedTS = time.Now().Unix()
			return true
		}
	}
	return false
}

// SetProfile 覆盖设置用户画像
func (b *Brain) SetProfile(p UserProfile) {
	if b.state == nil {
		b.state = NewWorldState()
	}
	b.state.Profile = &p
	b.state.UpdatedTS = time.Now().Unix()
}

// PatchProfile 以函数方式局部更新画像
func (b *Brain) PatchProfile(patch func(p *UserProfile)) {
	if patch == nil {
		return
	}
	if b.state == nil {
		b.state = NewWorldState()
	}
	if b.state.Profile == nil {
		b.state.Profile = &UserProfile{}
	}
	patch(b.state.Profile)
	b.state.UpdatedTS = time.Now().Unix()

}

// 推断目的
func (b *Brain) pushOrBoostHypothesis(id, title string, baseConf float64, source string) {
	if baseConf < 0 {
		baseConf = 0
	}
	if baseConf > 1 {
		baseConf = 1
	}
	now := time.Now().Unix()
	for i := range b.state.GoalHypotheses {
		if b.state.GoalHypotheses[i].GoalID == id {
			// 轻微提升并更新时间
			b.state.GoalHypotheses[i].Confidence = clamp01(b.state.GoalHypotheses[i].Confidence*0.7 + baseConf*0.3)
			b.state.GoalHypotheses[i].LastUpdated = now
			b.state.GoalHypotheses[i].Evidence = append(b.state.GoalHypotheses[i].Evidence, EvidenceItem{TS: now, Source: source, Summary: title, Weight: baseConf})
			return
		}
	}
	b.state.GoalHypotheses = append(b.state.GoalHypotheses, GoalHypothesis{
		GoalID:      id,
		Title:       title,
		Confidence:  baseConf,
		Evidence:    []EvidenceItem{{TS: now, Source: source, Summary: title, Weight: baseConf}},
		LastUpdated: now,
		HalfLifeMS:  int64(time.Hour * 2 / time.Millisecond), // 2h 半衰期
		Status:      "active",
	})
}

func (b *Brain) decayHypotheses() {
	now := time.Now().Unix()
	for i := range b.state.GoalHypotheses {
		h := &b.state.GoalHypotheses[i]
		if h.HalfLifeMS <= 0 {
			h.HalfLifeMS = int64(time.Hour / time.Millisecond)
		}
		dt := now - h.LastUpdated
		if dt <= 0 {
			continue
		}
		// 指数衰减：confidence *= 0.5^(dt/half_life)
		decay := math.Pow(0.5, float64(dt)/float64(h.HalfLifeMS))
		h.Confidence = clamp01(h.Confidence * decay)
		if h.Confidence < 0.05 {
			h.Status = "decayed"
		}
	}
}

func (b *Brain) sortHypotheses() {
	sort.SliceStable(b.state.GoalHypotheses, func(i, j int) bool {
		return b.state.GoalHypotheses[i].Confidence > b.state.GoalHypotheses[j].Confidence
	})
}

func clamp01(v float64) float64 {
	if v < 0 {
		return 0
	}
	if v > 1 {
		return 1
	}
	return v
}

// SetHealth 覆盖设置健康状态（同时刷新时间戳，秒）
func (b *Brain) SetHealth(h HealthState) {
	b.mu.Lock()
	defer b.mu.Unlock()
	now := time.Now().Unix()
	h.LastUpdatedTS = now
	if b.state == nil {
		b.state = NewWorldState()
	}
	b.state.Health = &h
	b.state.UpdatedTS = now
}

// PatchHealth 以函数方式局部修改健康状态
func (b *Brain) PatchHealth(patch func(h *HealthState)) {
	b.PatchHealthWithTimestamp(patch, time.Now().Unix())
}

// PatchHealthWithTimestamp 以函数方式局部修改健康状态，并指定时间戳
func (b *Brain) PatchHealthWithTimestamp(patch func(h *HealthState), timestamp int64) {
	if patch == nil {
		return
	}
	if b.state == nil {
		b.state = NewWorldState()
	}
	if b.state.Health == nil {
		b.state.Health = &HealthState{}
	}
	patch(b.state.Health)
	b.state.Health.LastUpdatedTS = timestamp
	b.state.UpdatedTS = timestamp
}

// ListTodos 返回当前待办列表（拷贝）
func (b *Brain) ListTodos() []TodoItem {
	if b.state == nil || b.state.Todos == nil {
		return nil
	}
	out := make([]TodoItem, len(b.state.Todos.Items))
	copy(out, b.state.Todos.Items)
	return out
}

// SetCurrentTask 设置/替换当前任务（StartedTS 为0则使用当前秒）
func (b *Brain) SetCurrentTask(t Task) {
	if b.state == nil {
		b.state = NewWorldState()
	}
	if t.StartedTS == 0 {
		t.StartedTS = time.Now().Unix()
	}
	b.state.CurrentTask = &t
	b.state.UpdatedTS = time.Now().Unix()
}

// ClearCurrentTask 清空当前任务
func (b *Brain) ClearCurrentTask() {
	if b.state == nil || b.state.CurrentTask == nil {
		return
	}
	b.state.CurrentTask = nil
	b.state.UpdatedTS = time.Now().Unix()
}

// GetCurrentTask 获取当前任务（拷贝返回）
func (b *Brain) GetCurrentTask() *Task {
	if b.state == nil || b.state.CurrentTask == nil {
		return nil
	}
	t := *b.state.CurrentTask
	return &t
}

// AddEvent 追加事件（ts 为0则用当前秒级时间戳）
func (b *Brain) AddEvent(e Event) {
	if b.state == nil {
		b.state = NewWorldState()
	}
	if b.state.Events == nil {
		b.state.Events = &EventsState{}
	}
	if e.TS == 0 {
		e.TS = time.Now().Unix()
	}
	b.state.Events.Recent = append(b.state.Events.Recent, e)
	b.state.UpdatedTS = time.Now().Unix()
}

// ListEvents 返回最近事件列表（拷贝）
func (b *Brain) ListEvents() []Event {
	if b.state == nil || b.state.Events == nil {
		return nil
	}
	out := make([]Event, len(b.state.Events.Recent))
	copy(out, b.state.Events.Recent)
	return out
}
