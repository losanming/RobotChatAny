package trigger

import (
	"fmt"
	"reflect"
)

// getTriggerID 获取触发器的唯一标识
func getTriggerID(trigger interface{}) string {
	switch t := trigger.(type) {
	case TimePointTrigger:
		return fmt.Sprintf("time_%d_%d_%s", t.Hour, t.Minute, t.Topic)
	case RealTimeTrigger:
		return fmt.Sprintf("real_%s_%s", t.Type, t.Content)
	case UserBehaviorTrigger:
		return fmt.Sprintf("user_%s_%s", t.Type, t.Content)
	case RandomTopicTrigger:
		return fmt.Sprintf("random_%s_%s", t.Type, t.Content)
	default:
		// 如果是未知类型，使用类型名称和内存地址作为ID
		return fmt.Sprintf("unknown_%s_%p", reflect.TypeOf(trigger).String(), trigger)
	}
}

// isTriggered 检查指定的触发器是否已经在今天触发过
func (m *Manager) isTriggered(timestamp int64, triggerID string) bool {
	if history, ok := m.triggerHistory[timestamp]; ok {
		for _, id := range history {
			if id == triggerID {
				return true
			}
		}
	}
	return false
}

// markTriggered 标记触发器已被触发
func (m *Manager) markTriggered(timestamp int64, triggerID string) {
	if _, ok := m.triggerHistory[timestamp]; !ok {
		m.triggerHistory[timestamp] = make([]string, 0)
	}
	m.triggerHistory[timestamp] = append(m.triggerHistory[timestamp], triggerID)
}
