package authContract

// LoginRequest 登录请求
type LoginRequest struct {
	// 登录类型: username(用户名密码), mobile_pwd(手机号密码), mobile_sms(手机号短信), gitee(Gitee登录)
	LoginType    string `json:"login_type" binding:"required,oneof=username mobile_pwd mobile_sms gitee"`
	Username     string `json:"username"`      // 用户名
	Mobile       string `json:"mobile"`        // 手机号
	Password     string `json:"password"`      // 密码
	SmsCode      string `json:"sms_code"`      // 短信验证码
	CaptchaToken string `json:"captcha_token"` // 滑块验证成功后的token
	Code         string `json:"code"`          // OAuth授权码
}

// SliderCaptchaVerify 滑块验证请求
type SliderCaptchaVerifyRequest struct {
	Token     string `json:"token" binding:"required"`      // 验证token
	XPosition int    `json:"x_position" binding:"required"` // 用户滑动的X轴位置
}
