package model

// SliderCaptcha 滑块验证码模型
type SliderCaptcha struct {
	Token         string `json:"token"`          // 验证token
	OriginalImage string `json:"original_image"` // 原始图片base64
	SliderImage   string `json:"slider_image"`   // 滑块图片base64
	XPosition     int    `json:"-"`              // 正确的X轴位置(后端保存，不返回前端)
	ExpireTime    int64  `json:"-"`              // 过期时间戳
}

// SliderCaptchaVerify 滑块验证请求
type SliderCaptchaVerify struct {
	Token      string `json:"token" binding:"required"`      // 验证token
	XPosition  int    `json:"x_position" binding:"required"` // 用户滑动的X轴位置
}

// VerificationResult 验证结果
type VerificationResult struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Token   string `json:"token,omitempty"` // 验证成功后的token
}
