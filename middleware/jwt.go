package middleware

import (
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"new-robot-self/model"
	"strings"
	"sync"

	"github.com/gin-gonic/gin"
)

type WhiteListConfig struct {
	paths map[string]bool
	mu    sync.RWMutex
}

// NewWhiteListConfig 创建新的白名单配置
func NewWhiteListConfig() *WhiteListConfig {
	return &WhiteListConfig{
		paths: make(map[string]bool),
	}
}

// Add 添加白名单路径
func (w *WhiteListConfig) Add(paths ...string) {
	w.mu.Lock()
	defer w.mu.Unlock()
	for _, path := range paths {
		w.paths[path] = true
	}
}

// Remove 移除白名单路径
func (w *WhiteListConfig) Remove(paths ...string) {
	w.mu.Lock()
	defer w.mu.Unlock()
	for _, path := range paths {
		delete(w.paths, path)
	}
}

// Clear 清空白名单
func (w *WhiteListConfig) Clear() {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.paths = make(map[string]bool)
}

// JWTAuth JWT认证中间件
func JWTAuth(whiteList *WhiteListConfig) gin.HandlerFunc {
	return func(c *gin.Context) {

		// 检查是否是白名单路径
		if whiteList != nil && whiteList.Contains(c.Request.URL.Path) {
			c.Next()
			return
		}

		// 从请求头中获取token
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code": 401,
				"msg":  "请求未携带token，无权限访问",
			})
			c.Abort()
			return
		}

		// 检查token格式
		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code": 401,
				"msg":  "token格式不正确",
			})
			c.Abort()
			return
		}

		// 解析token
		tokenString := parts[1]
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// 验证签名方法
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			// 返回签名密钥
			return []byte("your-secret-key"), nil
		})

		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code": 401,
				"msg":  "token无效: " + err.Error(),
			})
			c.Abort()
			return
		}

		// 验证token有效性
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			// 将用户信息保存到上下文
			userId := int64(claims["user_id"].(float64))
			username := claims["username"].(string)
			expiresAt := int64(claims["exp"].(float64))

			c.Set("token_info", &model.TokenInfo{
				UserId:    userId,
				Username:  username,
				ExpiresAt: expiresAt,
			})
			c.Next()
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code": 401,
				"msg":  "token无效",
			})
			c.Abort()
			return
		}
	}
}

// Contains 检查路径是否在白名单中
func (w *WhiteListConfig) Contains(path string) bool {
	w.mu.RLock()
	defer w.mu.RUnlock()

	// 直接匹配
	if w.paths[path] {
		return true
	}

	// 通配符匹配
	for whitePath := range w.paths {
		if strings.HasSuffix(whitePath, "/*") {
			prefix := strings.TrimSuffix(whitePath, "/*")
			if strings.HasPrefix(path, prefix) {
				return true
			}
		}
	}

	return false
}
