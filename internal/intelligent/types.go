package intelligent

// Intent 用户意图
type Intent struct {
	Type        string                 // 意图类型: schedule, navigation, health, work, info, life
	Confidence  float64                // 置信度 0-1
	Entities    map[string]interface{} // 提取的实体信息
	Priority    string                 // 优先级: high, medium, low
	Context     string                 // 上下文信息
	Timestamp   int64                  // 意图识别时间
}

// Entity 实体信息
type Entity struct {
	Type   string      // 实体类型: time, location, person, event, quantity
	Value  interface{} // 实体值
	Confidence float64  // 置信度
}

// ToolCallStep 工具调用步骤
type ToolCallStep struct {
	ToolName    string                 // 工具名称
	Parameters  map[string]interface{} // 调用参数
	DependsOn   []string               // 依赖的工具步骤
	Condition   string                 // 调用条件
	Priority    int                    // 调用优先级
}

// ToolCallPlan 工具调用计划
type ToolCallPlan struct {
	Intent         Intent          // 用户意图
	Steps          []ToolCallStep // 工具调用步骤
	Strategy       string          // 协调策略: sequential, parallel, conditional, iterative
	ExpectedOutput string          // 期望输出格式
	CreatedAt      int64           // 创建时间
}

// ConversationHistory 对话历史
type ConversationHistory struct {
	Messages []Message // 消息列表
	MaxSize  int       // 最大保存条数
}

// Message 消息
type Message struct {
	Role      string    // 角色: user, assistant
	Content   string    // 消息内容
	Intent    *Intent   // 识别的意图
	Timestamp int64     // 时间戳
}

// RecentActivity 最近活动
type RecentActivity struct {
	Activities []Activity // 活动列表
	TimeWindow int64      // 时间窗口（秒）
}

// Activity 活动
type Activity struct {
	Type      string                 // 活动类型
	Details   map[string]interface{} // 活动详情
	Timestamp int64                  // 时间戳
}

// UserPreferences 用户偏好
type UserPreferences struct {
	Profile      map[string]interface{} // 用户档案
	Preferences  map[string]interface{} // 用户偏好
	Habits       map[string]interface{} // 用户习惯
	LastUpdated  int64                  // 最后更新时间
}

// TriggerRule 触发规则
type TriggerRule struct {
	ID         string                 // 规则ID
	Type       string                 // 触发类型: time, state, event, smart
	Condition  string                 // 触发条件
	Action     string                 // 触发动作
	Parameters map[string]interface{} // 触发参数
	Enabled    bool                   // 是否启用
	LastTrigger int64                 // 最后触发时间
}

// AssistantResponse 助手回复
type AssistantResponse struct {
	Text        string                 // 回复文本
	Intent      *Intent               // 识别的意图
	ToolResults map[string]interface{} // 工具调用结果
	Actions     []string              // 建议的行动
	Confidence  float64               // 整体置信度
	Timestamp   int64                 // 回复时间
}