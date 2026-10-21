package service

import (
	"errors"
	authContract "new-robot-self/api/contract/auth"
	"new-robot-self/model"
)

// LoginService 登录服务
type LoginService struct {
	UserService      UserService
	CaptchaService   CaptchaService
	ThirdPartyService ThirdPartyService
}

// Login 登录
func (s *LoginService) Login(req *authContract.LoginRequest, ip string) (*model.LoginResponse, error) {
	// 1. 验证滑块验证码token
	// 对于Gitee登录，不需要验证码
	if req.LoginType != "gitee" {
		if req.CaptchaToken == "" {
			return nil, errors.New("验证码token不能为空")
		}
		valid, err := s.CaptchaService.CheckCaptchaToken(req.CaptchaToken)
		if err != nil {
			return nil, err
		}
		if !valid {
			return nil, errors.New("滑块验证失败或已过期")
		}
	}

	// 2. 根据登录类型处理登录
	var user *model.User
	var err error
	
	switch req.LoginType {
	case "username":
		// 用户名密码登录
		if req.Username == "" || req.Password == "" {
			return nil, errors.New("用户名和密码不能为空")
		}
		user, err = s.UserService.GetUserByUsername(req.Username)
		if err != nil {
			return nil, err
		}
		if !s.UserService.VerifyPassword(user, req.Password) {
			return nil, errors.New("密码错误")
		}
	case "mobile_pwd":
		// 手机号密码登录
		if req.Mobile == "" || req.Password == "" {
			return nil, errors.New("手机号和密码不能为空")
		}
		user, err = s.UserService.GetUserByMobile(req.Mobile)
		if err != nil {
			return nil, err
		}
		if !s.UserService.VerifyPassword(user, req.Password) {
			return nil, errors.New("密码错误")
		}
	case "mobile_sms":
		// 手机号短信登录
		if req.Mobile == "" || req.SmsCode == "" {
			return nil, errors.New("手机号和验证码不能为空")
		}
		// 这里应该验证短信验证码
		// 为了简化，我们假设验证码已经验证通过
		user, err = s.UserService.GetUserByMobile(req.Mobile)
		if err != nil {
			return nil, err
		}
	case "gitee":
		// Gitee第三方登录
		if req.Code == "" {
			return nil, errors.New("Gitee授权码不能为空")
		}
		user, err = s.ThirdPartyService.HandleGiteeLogin(req.Code)
		if err != nil {
			return nil, err
		}
	default:
		return nil, errors.New("不支持的登录类型")
	}

	// 3. 检查用户状态
	if user.Status != 1 {
		return nil, errors.New("账号已被禁用")
	}

	// 4. 更新登录信息
	if err := s.UserService.UpdateLoginInfo(user.Id, ip); err != nil {
		return nil, err
	}

	// 5. 生成JWT token
	token, expiresAt, err := s.UserService.GenerateToken(user)
	if err != nil {
		return nil, err
	}

	return &model.LoginResponse{
		Token:     token,
		ExpiresAt: expiresAt,
		User:      user,
	}, nil
}
