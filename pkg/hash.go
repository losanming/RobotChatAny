package pkg

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"time"
)

// BcryptHash 使用 bcrypt 对密码进行加密
func BcryptHash(password string) string {
	bytes, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes)
}

// BcryptCheck 对比明文密码和数据库的哈希值
func BcryptCheck(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
func MD5V(str []byte, b ...byte) string {
	h := md5.New()
	h.Write(str)
	return hex.EncodeToString(h.Sum(b))
}

// generateConnID 生成连接ID
func GenerateConnID() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}
