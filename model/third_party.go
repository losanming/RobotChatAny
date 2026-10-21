package model

// UserThirdParty 第三方账号关联
type UserThirdParty struct {
	Id               int64  `json:"id" xorm:"pk autoincr"`
	Code             string `json:"code" xorm:"notnull varchar(32) unique"`  // 唯一标识
	UserCode         string `json:"user_code" xorm:"notnull varchar(32) index"`  // 用户唯一标识
	Platform         string `json:"platform" xorm:"notnull varchar(20)"`
	PlatformUserId   string `json:"platform_user_id" xorm:"notnull varchar(100)"`
	PlatformUsername string `json:"platform_username" xorm:"varchar(100)"`
	PlatformAvatar   string `json:"platform_avatar" xorm:"varchar(255)"`
	AccessToken      string `json:"-" xorm:"varchar(255)"`
	RefreshToken     string `json:"-" xorm:"varchar(255)"`
	ExpiresAt        int64  `json:"expires_at" xorm:"bigint"`
	CreatedAt        int64  `json:"created_at" xorm:"notnull bigint created"`
	UpdatedAt        int64  `json:"updated_at" xorm:"notnull bigint updated"`
}
