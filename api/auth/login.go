package authApi

import (
	authContract "new-robot-self/api/contract/auth"
	"new-robot-self/pkg/response"
	"new-robot-self/service"

	"github.com/gin-gonic/gin"
)

// LoginApi 登录API
type LoginApi struct {
	LoginService service.LoginService
}

// Login 登录
func (api *LoginApi) Login(c *gin.Context) {
	var req authContract.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.RespContent(c, response.ErrBadRequest, nil)
		return
	}
	// 获取客户端IP
	clientIP := c.ClientIP()

	// 调用登录服务
	result, err := api.LoginService.Login(&req, clientIP)
	if err != nil {
		response.RespContent(c, response.ErrLoginFailed, nil)
		return
	}
	resp := authContract.LoginResponse{
		ExpiresAt: result.ExpiresAt,
		Token:     result.Token,
		User: authContract.UserData{
			Code:     result.User.Code,
			Id:       result.User.Id,
			Nickname: result.User.Nickname,
			Role:     nil,
			Username: result.User.Username,
		},
	}
	response.RespContent(c, response.StatusOK, resp)
}
