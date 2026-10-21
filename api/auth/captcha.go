package authApi

import (
	"log/slog"
	authContract "new-robot-self/api/contract/auth"
	"new-robot-self/pkg/response"
	"new-robot-self/service"

	"github.com/gin-gonic/gin"
)

// CaptchaApi 验证码API
type CaptchaApi struct {
	CaptchaService service.CaptchaService
}

// GenerateSliderCaptcha 生成滑块验证码
func (api *CaptchaApi) GenerateSliderCaptcha(c *gin.Context) {
	captcha, err := api.CaptchaService.GenerateSliderCaptcha()
	if err != nil {
		response.Resp(c, response.ErrCaptchaGenerateFailed)
		return
	}
	var resp authContract.GetSliderCaptchaResponse
	resp.SliderImage = captcha.SliderImage
	resp.OriginalImage = captcha.OriginalImage
	resp.Token = captcha.Token
	slog.Info("生成滑块验证码成功", "滑块位置", captcha.XPosition)
	response.RespContent(c, response.StatusOK, resp)
}

// VerifySliderCaptcha 验证滑块验证码
func (api *CaptchaApi) VerifySliderCaptcha(c *gin.Context) {
	var req authContract.SliderCaptchaVerifyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Resp(c, response.ErrBadRequest)
		return
	}

	result, err := api.CaptchaService.VerifySliderCaptcha(req.Token, req.XPosition)
	if err != nil {
		response.Resp(c, response.ErrCaptchaVerifyFailed)
		return
	}
	var resp authContract.SliderCaptchaVerificationResponse
	resp.Message = result.Message
	resp.Success = result.Success
	resp.Token = result.Token
	response.RespContent(c, response.StatusOK, resp)
}
