package trigger

// TriggerType 定义触发类型
type TriggerType int

const (
	TimePoint TriggerType = iota // 时间点触发
	RealTime                     // 实时事件触发
	UserBehavior                 // 用户行为触发
	RandomTopic                  // 随机话题触发
)

// String 返回TriggerType的字符串表示
func (t TriggerType) String() string {
	switch t {
	case TimePoint:
		return "TimePoint"
	case RealTime:
		return "RealTime"
	case UserBehavior:
		return "UserBehavior"
	case RandomTopic:
		return "RandomTopic"
	default:
		return "Unknown"
	}
}

// TriggerResult 触发结果
type TriggerResult struct {
	Type    TriggerType  // 触发器类型
	Content interface{} // 触发器内容
}

// TimePointTrigger 时间点触发配置
type TimePointTrigger struct {
	Hour   int    `json:"hour"`
	Minute int    `json:"minute"`
	Topic  string `json:"topic"`
}

// RealTimeTrigger 实时事件触发配置
type RealTimeTrigger struct {
	Type    string `json:"type"` // game_update, hot_news, weather, holiday
	Content string `json:"content"`
}

// UserBehaviorTrigger 用户行为触发配置
type UserBehaviorTrigger struct {
	Type    string `json:"type"` // online, offline_return, late_night
	Content string `json:"content"`
}

// RandomTopicTrigger 随机话题触发配置
type RandomTopicTrigger struct {
	Type    string `json:"type"` // game, work, life, anime
	Content string `json:"content"`
}
