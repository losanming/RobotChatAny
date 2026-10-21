package authRouter

import (
	"github.com/gin-gonic/gin"
	authApi "new-robot-self/api/auth"
)

type AuthRouter struct{}

func (s *AuthRouter) InitAuthRouter(r *gin.RouterGroup) {
	userGroup := r.Group("/auth")
	var captcha = authApi.CaptchaApi{}
	{
		userGroup.GET("/captcha", captcha.GenerateSliderCaptcha)
		userGroup.POST("/captcha", captcha.VerifySliderCaptcha)
	}
}
