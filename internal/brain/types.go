package brain

import (
	"time"

	"new-robot-self/internal/awareness"
)

// WorldState 全能环境状态（可拓展）。
// 设计目标：
// - 模块化：按领域拆分，便于独立演化；
// - 低耦合：每个子域可为空；
// - 可拓展：保留 Extras/Reserved 以及版本号；
// - 全局时间戳一律毫秒。

type WorldState struct {
	// 版本与元信息
	Version   string         `json:"version"`          // 状态结构版本，便于灰度升级
	UpdatedTS int64          `json:"updated_ts_ms"`    // 最后更新时间（ms）
	Extras    map[string]any `json:"extras,omitempty"` // 额外信息（自由扩展）

	// 用户侧
	Profile *UserProfile `json:"profile,omitempty"`
	Health  *HealthState `json:"health,omitempty"` // 预留健康拓展（体重、血糖、作息等）

	// 环境与感知（来自 internal/awareness）
	Env *Environment `json:"environment,omitempty"`

	// 对话与记忆
	Dialog *DialogContext `json:"dialog,omitempty"`
	Todos  *TodosState    `json:"todos,omitempty"`
	Events *EventsState   `json:"events,omitempty"`

	// 目的/任务
	CurrentTask    *Task            `json:"current_task,omitempty"`    // 显式目标（用户直接表达）
	GoalHypotheses []GoalHypothesis `json:"goal_hypotheses,omitempty"` // 目的假设列表（隐式推断）
}

// UserProfile 用户画像（可与 xorm 模型对齐）

type UserProfile struct {
	UserID     string         `json:"user_id"`
	Nick       string         `json:"nick,omitempty"`
	Age        int            `json:"age,omitempty"`
	Gender     int            `json:"gender,omitempty"` // 1男 2女
	HomeAdcode string         `json:"home_adcode,omitempty"`
	Prefs      map[string]any `json:"prefs,omitempty"` // 习惯/偏好，如：技术栈、运动频率、饮食等
}

// HealthState 健康数据（可逐步扩展）。
// 注意：仅做轻量聚合，详细时序数据应在独立健康模块存储/分析。

/*
@@TODO
healthData := HealthState{
    LastUpdatedTS: 1703129400000, // 2023-12-21 00:00:00 UTC时间戳(毫秒)
    WeightKg:      70.5,          // 体重70.5公斤
    BMI:           22.8,          // BMI指数22.8 (正常范围)
    BloodSugar: &BloodSugarBrief{
        RecentAvg: 5.6,  // 最近平均血糖5.6 mmol/L (正常范围)
        RiskFlag:  false, // 无血糖风险
    },
    ExercisePlan: "三天循环", // 运动计划
    Extras: map[string]any{
        "steps":        8500,    // 当日步数
        "sleep_hours":  7.5,     // 睡眠时长
        "water_intake": 2.0,     // 饮水量(升)
        "heart_rate":   72,      // 静息心率
    },
}
	// extras可以继续拓展身体的各项数据。但是是否需要抽离出来？？？
*/

type HealthState struct {
	LastUpdatedTS int64            `json:"last_updated_ts_ms"`
	WeightKg      float64          `json:"weight_kg,omitempty"` // 体重
	HeightM       float64          `json:"height_m,omitempty"`  // 身高
	BMI           float64          `json:"bmi,omitempty"`
	BloodSugar    *BloodSugarBrief `json:"blood_sugar,omitempty"`
	BloodPressure *BloodPressure   `json:"blood_pressure,omitempty"`
	ExercisePlan  []ExecPlan       `json:"exercise_plan,omitempty"` // 运动计划
	ExecDesc      string           `json:"exec_desc,omitempty"`     // 运动计划备注
	Steps         int64            `json:"steps,omitempty"`         // 步数
	WaterIntake   float64          `json:"water_intake,omitempty"`  // 水摄入量 L
	HeartRate     int64            `json:"heart_rate,omitempty"`    // 心率(静心)
	SleepHours    float64          `json:"sleep_hours,omitempty"`   // 睡眠时长avg hours
}

// 运动执行计划 (简易)
type ExecPlan struct {
	Id          int64  `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

// 血压
type BloodPressure struct {
	Systolic  float64 `json:"systolic,omitempty"`  // 舒张压，单位mmHg
	Diastolic float64 `json:"diastolic,omitempty"` // 收缩压，单位mmHg
}

type BloodSugarBrief struct {
	RecentAvg float64 `json:"recent_avg,omitempty"` // 最近 N 分钟的均值
	RiskFlag  bool    `json:"risk_flag,omitempty"`  // 是否有血糖风险
}

// Environment 汇聚 awareness 的结果（以引用方式存放，避免重复拷贝大字段）。

type Environment struct {
	Awareness awareness.Context `json:"awareness"`
}

// DialogContext 对话与ASR的窗口

type DialogContext struct {
	History        []DialogTurn   `json:"history,omitempty"` // 近 N 轮
	LastUserIntent string         `json:"last_user_intent,omitempty"`
	Extras         map[string]any `json:"extras,omitempty"`
}

type DialogTurn struct {
	TS      int64  `json:"ts_ms"`
	Role    string `json:"role"` // user|assistant|system
	Content string `json:"content"`
}

// Todos 与事件流

type TodosState struct {
	Items []TodoItem `json:"items,omitempty"`
}

type TodoItem struct {
	ID         string `json:"id"`
	Title      string `json:"title"`
	CreateTime int64  `json:"create_time"` // 这个计划的创建时间
	DueTS      int64  `json:"due_ts_ms,omitempty"`
	Priority   int    `json:"priority,omitempty"` // 1高 2中 3低
	Status     string `json:"status,omitempty"`   // pending|in_progress|done
}

type EventsState struct {
	Recent []Event `json:"recent,omitempty"` // 近一段时间的关键事件
}

type Event struct {
	TS   int64          `json:"ts_ms"`
	Type string         `json:"type"`
	Data map[string]any `json:"data,omitempty"`
}

// Task 显式当前任务

type Task struct {
	TaskID    string         `json:"task_id"`
	Title     string         `json:"title"`
	StartedTS int64          `json:"started_ts_ms"`
	Deadline  int64          `json:"deadline_ts_ms,omitempty"`
	Meta      map[string]any `json:"meta,omitempty"`
}

// GoalHypothesis 隐式目的假设

type GoalHypothesis struct {
	GoalID      string         `json:"goal_id"`
	Title       string         `json:"title"`
	Confidence  float64        `json:"confidence"`         // 0..1
	Evidence    []EvidenceItem `json:"evidence,omitempty"` // 证据摘要
	LastUpdated int64          `json:"last_updated_ts_ms"`
	HalfLifeMS  int64          `json:"half_life_ms"`     // 置信度衰减半衰期（ms）
	Status      string         `json:"status,omitempty"` // active|decayed|rejected
}

type EvidenceItem struct {
	TS      int64   `json:"ts_ms"`
	Source  string  `json:"source"` // awareness|asr|dialog|todo|manual
	Summary string  `json:"summary"`
	Weight  float64 `json:"weight"`
}

func NewWorldState() *WorldState {
	return &WorldState{
		Version:   "v1",
		UpdatedTS: time.Now().Unix(),
	}
}
