package service

import (
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"new-robot-self/model"
	"new-robot-self/public"
	"time"
)

// UserService 用户服务
type UserService struct{}

// GetUserByUsername 通过用户名获取用户
func (s *UserService) GetUserByUsername(username string) (*model.User, error) {
	var user model.User
	has, err := public.DB.Where("username = ?", username).Get(&user)
	if err != nil {
		return nil, err
	}
	if !has {
		return nil, errors.New("用户不存在")
	}
	return &user, nil
}

// GetUserByMobile 通过手机号获取用户
func (s *UserService) GetUserByMobile(mobile string) (*model.User, error) {
	var user model.User
	has, err := public.DB.Where("mobile = ?", mobile).Get(&user)
	if err != nil {
		return nil, err
	}
	if !has {
		return nil, errors.New("用户不存在")
	}
	return &user, nil
}

// GetUserByMobile 通过手机号获取用户
func (s *UserService) GetUserByCode(code string) (*model.User, error) {
	var user model.User
	has, err := public.DB.Where("code = ?", code).Get(&user)
	if err != nil {
		return nil, err
	}
	if !has {
		return nil, errors.New("用户不存在")
	}
	return &user, nil
}

// GetUserByMobile 通过手机号获取用户
func (s *UserService) Create(user *model.User) error {
	_, err := public.DB.Insert(user)
	return err
}

// VerifyPassword 验证密码
func (s *UserService) VerifyPassword(user *model.User, password string) bool {
	// 在实际应用中，应该使用bcrypt等加密算法比较密码
	// 这里为了简单，直接比较明文密码
	return user.Password == password
}

// UpdateLoginInfo 更新登录信息
func (s *UserService) UpdateLoginInfo(userId int64, ip string) error {
	now := time.Now().Unix()
	_, err := public.DB.Exec("UPDATE users SET last_login_time = ?, last_login_ip = ? WHERE id = ?", now, ip, userId)
	return err
}

// GenerateToken 生成JWT token
func (s *UserService) GenerateToken(user *model.User) (string, int64, error) {
	// 设置token有效期为24小时
	expiresAt := time.Now().Add(24 * time.Hour).Unix()

	claims := jwt.MapClaims{
		"user_id":  user.Id,
		"userCode": user.Code,
		"username": user.Username,
		"exp":      expiresAt,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// 使用密钥签名token（在实际应用中，应该从配置文件中读取密钥）
	tokenString, err := token.SignedString([]byte("your-secret-key"))
	if err != nil {
		return "", 0, err
	}

	return tokenString, expiresAt, nil
}
