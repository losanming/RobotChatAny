package authContract

// 滑块验证码信息
type GetSliderCaptchaResponse struct {
	Token         string `json:"token"`          // 验证token
	OriginalImage string `json:"original_image"` // 原始图片base64
	SliderImage   string `json:"slider_image"`   // 滑块图片base64
}

// 滑块验证码验证成功返回数据
type SliderCaptchaVerificationResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Token   string `json:"token,omitempty"` // 验证成功后的token
}

// LoginResponse 登录响应
type LoginResponse struct {
	Token     string   `json:"token"`      // JWT token
	ExpiresAt int64    `json:"expires_at"` // 过期时间戳
	User      UserData `json:"user"`       // 用户信息
}

type UserData struct {
	Id       int64  `json:"id"`
	Code     string `json:"code"`
	Username string `json:"username"`
	Nickname string `json:"nickname"`
	Role     []struct {
		Id   int64  `json:"id"`
		Code string `json:"code"`
		Name string `json:"name"`
	}
}
