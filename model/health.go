package model

// HealthRecord 健康记录模型
type HealthRecord struct {
	Id         int64   `xorm:"pk autoincr" json:"id"`
	UserID     string  `xorm:"varchar(50) notnull index" json:"user_id"`
	RecordType string  `xorm:"varchar(20) notnull index" json:"record_type"` // blood_sugar, weight, exercise, steps, sleep, heart_rate, water_intake
	Value      float64 `xorm:"double notnull" json:"value"`
	Unit       string  `xorm:"varchar(10)" json:"unit"` // mmol/L, kg, min, kcal, steps, hours, bpm, L
	Notes      string  `xorm:"text" json:"notes"`
	Metadata   string  `xorm:"json" json:"metadata"` // JSON 格式的额外数据
	RecordedAt int64   `xorm:"bigint notnull index" json:"recorded_at"` // 记录时间戳（秒）
	CreatedAt  int64   `xorm:"bigint created" json:"created_at"`
	UpdatedAt  int64   `xorm:"bigint updated" json:"updated_at"`
}

// TableName 表名
func (h *HealthRecord) TableName() string {
	return "health_records"
}

// TodoItem 待办事项模型
type TodoItem struct {
	Id          int64  `xorm:"pk autoincr" json:"id"`
	UserID      string `xorm:"varchar(50) notnull index" json:"user_id"`
	Title       string `xorm:"varchar(200) notnull" json:"title"`
	Description string `xorm:"text" json:"description"`
	Priority    string `xorm:"varchar(10) default 'medium'" json:"priority"` // high, medium, low
	Status      string `xorm:"varchar(20) default 'pending'" json:"status"`  // pending, in_progress, completed, cancelled
	Tags        string `xorm:"json" json:"tags"`                             // JSON 数组
	DueDate     int64  `xorm:"bigint index" json:"due_date"`                 // 截止时间戳
	CompletedAt int64  `xorm:"bigint" json:"completed_at"`                   // 完成时间戳
	CreatedAt   int64  `xorm:"bigint created" json:"created_at"`
	UpdatedAt   int64  `xorm:"bigint updated" json:"updated_at"`
}

// TableName 表名
func (t *TodoItem) TableName() string {
	return "todos"
}

// CalendarEvent 日历事件模型
type CalendarEvent struct {
	Id          int64  `xorm:"pk autoincr" json:"id"`
	UserID      string `xorm:"varchar(50) notnull index" json:"user_id"`
	Title       string `xorm:"varchar(200) notnull" json:"title"`
	Description string `xorm:"text" json:"description"`
	Location    string `xorm:"varchar(200)" json:"location"`
	StartTime   int64  `xorm:"bigint notnull index" json:"start_time"` // 开始时间戳
	EndTime     int64  `xorm:"bigint notnull" json:"end_time"`         // 结束时间戳
	IsAllDay    bool   `xorm:"bool default false" json:"is_all_day"`   // 是否全天
	Reminder    int64  `xorm:"bigint default 0" json:"reminder"`       // 提醒提前量（分钟）
	Status      string `xorm:"varchar(20) default 'confirmed'" json:"status"` // confirmed, cancelled, tentative
	CreatedAt   int64  `xorm:"bigint created" json:"created_at"`
	UpdatedAt   int64  `xorm:"bigint updated" json:"updated_at"`
}

// TableName 表名
func (c *CalendarEvent) TableName() string {
	return "calendar_events"
}

// DietRecord 饮食记录模型
type DietRecord struct {
	Id         int64   `xorm:"pk autoincr" json:"id"`
	UserID     string  `xorm:"varchar(50) notnull index" json:"user_id"`
	MealType   string  `xorm:"varchar(10) notnull index" json:"meal_type"` // breakfast, lunch, dinner, snack
	FoodName   string  `xorm:"varchar(200) notnull" json:"food_name"`
	Calories   float64 `xorm:"double" json:"calories"`            // 卡路里
	Protein    float64 `xorm:"double" json:"protein"`             // 蛋白质(g)
	Carbs      float64 `xorm:"double" json:"carbs"`               // 碳水化合物(g)
	Fat        float64 `xorm:"double" json:"fat"`                 // 脂肪(g)
	Serving    string  `xorm:"varchar(50)" json:"serving"`        // 份量描述
	Notes      string  `xorm:"text" json:"notes"`
	RecordedAt int64   `xorm:"bigint notnull index" json:"recorded_at"` // 记录时间戳
	CreatedAt  int64   `xorm:"bigint created" json:"created_at"`
}

// TableName 表名
func (d *DietRecord) TableName() string {
	return "diet_records"
}

// UserPreference 用户偏好模型
type UserPreference struct {
	Id        int64  `xorm:"pk autoincr" json:"id"`
	UserID    string `xorm:"varchar(50) notnull unique" json:"user_id"`
	Prefs     string `xorm:"json" json:"prefs"` // JSON 格式的偏好数据
	UpdatedAt int64  `xorm:"bigint updated" json:"updated_at"`
}

// TableName 表名
func (u *UserPreference) TableName() string {
	return "user_preferences"
}

// Reminder 提醒模型
type Reminder struct {
	Id          int64  `xorm:"pk autoincr" json:"id"`
	UserID      string `xorm:"varchar(50) notnull index" json:"user_id"`
	Type        string `xorm:"varchar(20) notnull index" json:"type"` // blood_sugar, exercise, diet, medication, work_rest
	Title       string `xorm:"varchar(200) notnull" json:"title"`
	Description string `xorm:"text" json:"description"`
	ScheduledAt int64  `xorm:"bigint notnull index" json:"scheduled_at"` // 计划提醒时间戳
	RepeatRule  string `xorm:"varchar(20) default ''" json:"repeat_rule"` // daily, weekly, custom, empty for once
	Status      string `xorm:"varchar(20) default 'pending'" json:"status"` // pending, completed, skipped, cancelled
	Priority    string `xorm:"varchar(10) default 'medium'" json:"priority"` // low, medium, high
	Metadata    string `xorm:"json" json:"metadata"` // JSON 格式的额外数据
	CreatedAt   int64  `xorm:"bigint created" json:"created_at"`
	UpdatedAt   int64  `xorm:"bigint updated" json:"updated_at"`
	CompletedAt int64  `xorm:"bigint default 0" json:"completed_at"` // 完成时间戳
}

// TableName 表名
func (r *Reminder) TableName() string {
	return "reminders"
}
