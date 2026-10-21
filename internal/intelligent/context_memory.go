package intelligent

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"time"
)

// ContextMemory 上下文记忆模块
type ContextMemory struct {
	shortTerm   *ConversationHistory
	mediumTerm  *RecentActivity
	longTerm    *UserPreferences
	mu          sync.RWMutex
	storagePath string // 持久化存储路径
}

func NewContextMemory() *ContextMemory {
	memory := &ContextMemory{
		shortTerm: &ConversationHistory{
			Messages: make([]Message, 0, 50),
			MaxSize:  50,
		},
		mediumTerm: &RecentActivity{
			Activities: make([]Activity, 0, 100),
			TimeWindow: 7 * 24 * 3600, // 7天
		},
		longTerm: &UserPreferences{
			Profile:     make(map[string]interface{}),
			Preferences: make(map[string]interface{}),
			Habits:      make(map[string]interface{}),
			LastUpdated: time.Now().Unix(),
		},
		storagePath: "data/context_memory.json",
	}

	// 尝试从持久化存储加载
	memory.loadFromStorage()

	return memory
}

// AddMessage 添加消息到短期记忆
func (m *ContextMemory) AddMessage(role, content string, intent *Intent) {
	m.mu.Lock()
	defer m.mu.Unlock()

	message := Message{
		Role:      role,
		Content:   content,
		Intent:    intent,
		Timestamp: time.Now().Unix(),
	}

	m.shortTerm.Messages = append(m.shortTerm.Messages, message)

	// 如果超过最大容量，删除最旧的消息
	if len(m.shortTerm.Messages) > m.shortTerm.MaxSize {
		m.shortTerm.Messages = m.shortTerm.Messages[1:]
	}
}

// LearnFromMessage 学习用户偏好的包装函数（不持有锁）
func (m *ContextMemory) LearnFromMessage(message Message) {
	// 学习用户偏好
	m.learnFromMessage(message)
	
	// 持久化到存储
	m.saveToStorage()
}

// GetRecentMessages 获取最近的消息
func (m *ContextMemory) GetRecentMessages(count int) []Message {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if count <= 0 || count > len(m.shortTerm.Messages) {
		count = len(m.shortTerm.Messages)
	}

	return m.shortTerm.Messages[len(m.shortTerm.Messages)-count:]
}

// AddActivity 添加活动到中期记忆
func (m *ContextMemory) AddActivity(activityType string, details map[string]interface{}) {
	m.mu.Lock()
	defer m.mu.Unlock()

	activity := Activity{
		Type:      activityType,
		Details:   details,
		Timestamp: time.Now().Unix(),
	}

	m.mediumTerm.Activities = append(m.mediumTerm.Activities, activity)

	// 清理过期活动
	m.cleanupOldActivities()
}

// LearnFromActivity 学习用户习惯的包装函数（不持有锁）
func (m *ContextMemory) LearnFromActivity(activity Activity) {
	m.mu.Lock()
	defer m.mu.Unlock()
	
	// 学习用户习惯
	m.learnFromActivity(activity)
	
	// 持久化到存储
	m.saveToStorage()
}

// GetRecentActivities 获取最近的活动
func (m *ContextMemory) GetRecentActivities(count int) []Activity {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if count <= 0 || count > len(m.mediumTerm.Activities) {
		count = len(m.mediumTerm.Activities)
	}

	return m.mediumTerm.Activities[len(m.mediumTerm.Activities)-count:]
}

// cleanupOldActivities 清理过期活动
func (m *ContextMemory) cleanupOldActivities() {
	now := time.Now().Unix()
	cutoff := now - m.mediumTerm.TimeWindow

	var validActivities []Activity
	for _, activity := range m.mediumTerm.Activities {
		if activity.Timestamp >= cutoff {
			validActivities = append(validActivities, activity)
		}
	}

	m.mediumTerm.Activities = validActivities
}

// learnFromMessage 从消息中学习用户偏好
func (m *ContextMemory) learnFromMessage(message Message) {
	if message.Intent == nil {
		return
	}

	// 根据意图类型学习偏好
	switch message.Intent.Type {
	case "schedule":
		m.learnSchedulePreference(message)
	case "navigation":
		m.learnNavigationPreference(message)
	case "health":
		m.learnHealthPreference(message)
	case "work":
		m.learnWorkPreference(message)
	case "life":
		m.learnLifePreference(message)
	}
}

// learnSchedulePreference 学习日程偏好
func (m *ContextMemory) learnSchedulePreference(message Message) {
	// 提取时间偏好
	if timeStr, ok := message.Intent.Entities["time"].(string); ok {
		m.UpdatePreference("preferred_time", timeStr)
	}

	// 提取事件类型偏好
	if event, ok := message.Intent.Entities["event"].(string); ok {
		m.UpdatePreference("preferred_activity", event)
	}
}

// learnNavigationPreference 学习导航偏好
func (m *ContextMemory) learnNavigationPreference(message Message) {
	// 提取交通方式偏好
	if location, ok := message.Intent.Entities["location"].(string); ok {
		// 记录常去的地点
		locations := make([]string, 0)
		if existingLocs, ok := m.GetPreference("frequent_locations"); ok {
			if locs, ok := existingLocs.([]string); ok {
				locations = locs
			}
		}
		locations = append(locations, location)
		m.UpdatePreference("frequent_locations", locations)
	}
}

// learnHealthPreference 学习健康偏好
func (m *ContextMemory) learnHealthPreference(message Message) {
	// 提取食物偏好
	if food, ok := message.Intent.Entities["food"].(string); ok {
		// 记录常吃的食物
		foods := make([]string, 0)
		if existingFoods, ok := m.GetPreference("frequent_foods"); ok {
			if fs, ok := existingFoods.([]string); ok {
				foods = fs
			}
		}
		foods = append(foods, food)
		m.UpdatePreference("frequent_foods", foods)
	}
}

// learnWorkPreference 学习工作偏好
func (m *ContextMemory) learnWorkPreference(message Message) {
	// 提取技术栈偏好
	// TODO: 实现工作偏好学习
}

// learnLifePreference 学习生活偏好
func (m *ContextMemory) learnLifePreference(message Message) {
	// 提取购物地点偏好
	// TODO: 实现生活偏好学习
}

// learnFromActivity 从活动中学习用户习惯
func (m *ContextMemory) learnFromActivity(activity Activity) {
	switch activity.Type {
	case "schedule":
		m.learnScheduleHabit(activity)
	case "navigation":
		m.learnNavigationHabit(activity)
	case "health":
		m.learnHealthHabit(activity)
	}
}

// learnScheduleHabit 学习日程习惯
func (m *ContextMemory) learnScheduleHabit(activity Activity) {
	// 分析活动时间模式
	timestamp := activity.Timestamp
	hour := time.Unix(timestamp, 0).Hour()

	// 统计每个时间段的活动频率
	timeSlot := fmt.Sprintf("hour_%d", hour)
	currentCount := 0
	if count, ok := m.longTerm.Habits[timeSlot].(int); ok {
		currentCount = count
	}
	m.UpdateHabit(timeSlot, currentCount+1)
}

// learnNavigationHabit 学习导航习惯
func (m *ContextMemory) learnNavigationHabit(activity Activity) {
	if location, ok := activity.Details["location"].(string); ok {
		// 统计访问地点的频率
		locationKey := fmt.Sprintf("location_%s", location)
		currentCount := 0
		if count, ok := m.longTerm.Habits[locationKey].(int); ok {
			currentCount = count
		}
		m.UpdateHabit(locationKey, currentCount+1)
	}
}

// learnHealthHabit 学习健康习惯
func (m *ContextMemory) learnHealthHabit(activity Activity) {
	// 分析健康数据模式
	// TODO: 实现健康习惯学习
}

// UpdatePreference 更新用户偏好
func (m *ContextMemory) UpdatePreference(key string, value interface{}) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.longTerm.Preferences[key] = value
	m.longTerm.LastUpdated = time.Now().Unix()
}

// GetPreference 获取用户偏好
func (m *ContextMemory) GetPreference(key string) (interface{}, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	value, ok := m.longTerm.Preferences[key]
	return value, ok
}

// UpdateProfile 更新用户档案
func (m *ContextMemory) UpdateProfile(key string, value interface{}) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.longTerm.Profile[key] = value
	m.longTerm.LastUpdated = time.Now().Unix()
}

// GetProfile 获取用户档案
func (m *ContextMemory) GetProfile(key string) (interface{}, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	value, ok := m.longTerm.Profile[key]
	return value, ok
}

// UpdateHabit 更新用户习惯
func (m *ContextMemory) UpdateHabit(key string, value interface{}) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.longTerm.Habits[key] = value
	m.longTerm.LastUpdated = time.Now().Unix()
}

// GetHabit 获取用户习惯
func (m *ContextMemory) GetHabit(key string) (interface{}, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	value, ok := m.longTerm.Habits[key]
	return value, ok
}

// GetContextSummary 获取上下文摘要
func (m *ContextMemory) GetContextSummary() map[string]interface{} {
	m.mu.RLock()
	defer m.mu.RUnlock()

	return map[string]interface{}{
		"recent_messages_count":   len(m.shortTerm.Messages),
		"recent_activities_count": len(m.mediumTerm.Activities),
		"preferences_count":       len(m.longTerm.Preferences),
		"habits_count":            len(m.longTerm.Habits),
		"last_updated":            m.longTerm.LastUpdated,
	}
}


// saveToStorage 持久化到存储
func (m *ContextMemory) saveToStorage() {
	data := struct {
		ShortTerm  *ConversationHistory `json:"short_term"`
		MediumTerm *RecentActivity      `json:"medium_term"`
		LongTerm   *UserPreferences     `json:"long_term"`
	}{
		ShortTerm:  m.shortTerm,
		MediumTerm: m.mediumTerm,
		LongTerm:   m.longTerm,
	}

	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		fmt.Printf("持久化上下文记忆失败: %v\n", err)
		return
	}

	// 确保目录存在
	dir := "data"
	if err := os.MkdirAll(dir, 0755); err != nil {
		fmt.Printf("创建数据目录失败: %v\n", err)
		return
	}

	if err := os.WriteFile(m.storagePath, jsonData, 0644); err != nil {
		fmt.Printf("写入上下文记忆文件失败: %v\n", err)
	}
}

// loadFromStorage 从存储加载
func (m *ContextMemory) loadFromStorage() {
	data, err := os.ReadFile(m.storagePath)
	if err != nil {
		// 文件不存在是正常的，首次运行
		return
	}

	var loadedData struct {
		ShortTerm  *ConversationHistory `json:"short_term"`
		MediumTerm *RecentActivity      `json:"medium_term"`
		LongTerm   *UserPreferences     `json:"long_term"`
	}

	if err := json.Unmarshal(data, &loadedData); err != nil {
		fmt.Printf("加载上下文记忆失败: %v\n", err)
		return
	}

	if loadedData.ShortTerm != nil {
		m.shortTerm = loadedData.ShortTerm
	}
	if loadedData.MediumTerm != nil {
		m.mediumTerm = loadedData.MediumTerm
	}
	if loadedData.LongTerm != nil {
		m.longTerm = loadedData.LongTerm
	}
}
