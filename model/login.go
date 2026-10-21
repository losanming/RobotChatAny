package model

// LoginRequest 登录请求
type LoginRequest struct {
	// 登录类型: username(用户名密码), mobile_pwd(手机号密码), mobile_sms(手机号短信)
	LoginType    string `json:"login_type" binding:"required,oneof=username mobile_pwd mobile_sms"`
	Username     string `json:"username"`     // 用户名
	Mobile       string `json:"mobile"`       // 手机号
	Password     string `json:"password"`     // 密码
	SmsCode      string `json:"sms_code"`     // 短信验证码
	CaptchaToken string `json:"captcha_token" binding:"required"` // 滑块验证成功后的token
}

// LoginResponse 登录响应
type LoginResponse struct {
	Token     string `json:"token"`      // JWT token
	ExpiresAt int64  `json:"expires_at"` // 过期时间戳
	User      *User  `json:"user"`       // 用户信息
}

// TokenInfo JWT token信息
type TokenInfo struct {
	UserId    int64  `json:"user_id"`
	Username  string `json:"username"`
	ExpiresAt int64  `json:"expires_at"`
}
