package model

// User 用户模型
type User struct {
	Id            int64  `xorm:"pk autoincr" json:"id"`
	Code          string `xorm:"varchar(50) notnull unique" json:"code"`
	Username      string `xorm:"varchar(50) notnull unique" json:"username"`
	Password      string `xorm:"varchar(100) notnull" json:"-"`
	Mobile        string `xorm:"varchar(20) unique" json:"mobile"`
	Email         string `xorm:"varchar(100) unique" json:"email"`
	Nickname      string `xorm:"varchar(50)" json:"nickname"`
	Avatar        string `xorm:"varchar(255)" json:"avatar"`
	Status        int    `xorm:"int default 1" json:"status"` // 1:正常 0:禁用
	LastLoginTime int64  `xorm:"bigint" json:"last_login_time"`
	LastLoginIp   string `xorm:"varchar(50)" json:"last_login_ip"`
	CreatedAt     int64  `xorm:"bigint created" json:"created_at"`
	UpdatedAt     int64  `xorm:"bigint updated" json:"updated_at"`
	DeletedAt     int64  `xorm:"bigint deleted" json:"-"`
}

// TableName 表名
func (u *User) TableName() string {
	return "users"
}
