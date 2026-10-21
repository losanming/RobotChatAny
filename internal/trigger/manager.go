package trigger

import (
	"math/rand"
	"sync"
	"time"
)

// Manager 触发管理器
type Manager struct {
	timePoints     []TimePointTrigger
	realTime       []RealTimeTrigger
	userBehavior   []UserBehaviorTrigger
	randomTopics   []RandomTopicTrigger
	triggerHistory map[int64][]string // key为日期时间戳(精确到天), value为已触发的消息ID
	mutex          sync.RWMutex
}

// NewManager 创建触发管理器
func NewManager() *Manager {
	return &Manager{
		triggerHistory: make(map[int64][]string),
	}
}

// LoadTriggers 加载触发器配置
func (m *Manager) LoadTriggers(
	timePoints []TimePointTrigger,
	realTime []RealTimeTrigger,
	userBehavior []UserBehaviorTrigger,
	randomTopics []RandomTopicTrigger,
) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	m.timePoints = timePoints
	m.realTime = realTime
	m.userBehavior = userBehavior
	m.randomTopics = randomTopics
}

// GetRandomTriggers 获取当天的随机触发消息
// maxCount: 最大消息数
// types: 指定要获取的触发类型，如果为空则获取所有类型
func (m *Manager) GetRandomTriggers(maxCount int, types []TriggerType) []TriggerResult {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	// 获取今天的时间戳(精确到天)
	now := time.Now()
	todayTimestamp := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location()).Unix()

	// 如果没有指定类型，则使用所有类型
	if len(types) == 0 {
		types = []TriggerType{TimePoint, RealTime, UserBehavior, RandomTopic}
	}

	// 准备所有可用的触发器
	var availableTriggers []TriggerResult
	for _, t := range types {
		switch t {
		case TimePoint:
			for _, trigger := range m.timePoints {
				// 检查是否已经触发过
				if !m.isTriggered(todayTimestamp, getTriggerID(trigger)) {
					availableTriggers = append(availableTriggers, TriggerResult{Type: TimePoint, Content: trigger})
				}
			}
		case RealTime:
			for _, trigger := range m.realTime {
				if !m.isTriggered(todayTimestamp, getTriggerID(trigger)) {
					availableTriggers = append(availableTriggers, TriggerResult{Type: RealTime, Content: trigger})
				}
			}
		case UserBehavior:
			for _, trigger := range m.userBehavior {
				if !m.isTriggered(todayTimestamp, getTriggerID(trigger)) {
					availableTriggers = append(availableTriggers, TriggerResult{Type: UserBehavior, Content: trigger})
				}
			}
		case RandomTopic:
			for _, trigger := range m.randomTopics {
				if !m.isTriggered(todayTimestamp, getTriggerID(trigger)) {
					availableTriggers = append(availableTriggers, TriggerResult{Type: RandomTopic, Content: trigger})
				}
			}
		}
	}

	// 如果没有可用的触发器，返回空
	if len(availableTriggers) == 0 {
		return nil
	}

	// 随机打乱触发器顺序
	rand.Shuffle(len(availableTriggers), func(i, j int) {
		availableTriggers[i], availableTriggers[j] = availableTriggers[j], availableTriggers[i]
	})

	// 选择指定数量的触发器
	resultCount := min(maxCount, len(availableTriggers))
	selected := availableTriggers[:resultCount]

	// 记录已触发的消息
	for _, trigger := range selected {
		m.markTriggered(todayTimestamp, getTriggerID(trigger.Content))
	}

	return selected
}

// min 返回两个整数中的较小值
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// ClearHistory 清除历史记录
func (m *Manager) ClearHistory() {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.triggerHistory = make(map[int64][]string)
}
