package protocol

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"sync"
	"time"
)

// SecurityManager 安全管理器
type SecurityManager struct {
	nonceMap     sync.Map    // 用于存储已使用的nonce
	nonceTimeout time.Duration // nonce的有效期
	secretKey    []byte      // 用于签名的密钥
}

// NewSecurityManager 创建安全管理器
func NewSecurityManager(secretKey []byte, nonceTimeout time.Duration) *SecurityManager {
	return &SecurityManager{
		secretKey:    secretKey,
		nonceTimeout: nonceTimeout,
	}
}

// StringToFixedBytes 将字符串转换为固定大小的字节数组
func StringToFixedBytes(s string, size int) []byte {
	bytes := make([]byte, size)
	copy(bytes, []byte(s))
	return bytes
}

// BytesToString 将字节数组转换为字符串，去除末尾的空字节
func BytesToString(bytes []byte) string {
	n := 0
	for i := len(bytes) - 1; i >= 0; i-- {
		if bytes[i] != 0 {
			n = i + 1
			break
		}
	}
	return string(bytes[:n])
}

// GenerateSignature 生成消息签名
func (s *SecurityManager) GenerateSignature(header *MessageHeader, body []byte) [SignatureSize]byte {
	// 1. 组合待签名数据
	data := fmt.Sprintf("%d:%d:%d:%d:%s:%d:%s",
		header.Version,
		header.Type,
		header.Sequence,
		header.Timestamp,
		BytesToString(header.SessionID[:]),
		header.Nonce,
		body,
	)

	// 2. 使用HMAC-SHA256计算签名
	h := hmac.New(sha256.New, s.secretKey)
	h.Write([]byte(data))
	signature := h.Sum(nil)

	// 3. 转换为固定大小的字节数组
	var result [SignatureSize]byte
	hex := hex.EncodeToString(signature)
	copy(result[:], []byte(hex))
	return result
}

// ValidateMessage 验证消息
func (s *SecurityManager) ValidateMessage(header *MessageHeader, body []byte) error {
	// 1. 验证时间戳
	now := time.Now().Unix()
	if now-header.Timestamp > 300 { // 5分钟过期
		return fmt.Errorf("消息已过期")
	}

	// 2. 验证nonce（防重放）
	sessionID := BytesToString(header.SessionID[:])
	nonceKey := fmt.Sprintf("%s:%d", sessionID, header.Nonce)
	if _, exists := s.nonceMap.LoadOrStore(nonceKey, now); exists {
		return fmt.Errorf("重复的消息")
	}

	// 3. 定时清理过期的nonce
	go func() {
		time.Sleep(s.nonceTimeout)
		s.nonceMap.Delete(nonceKey)
	}()

	// 4. 验证签名
	expectedSignature := s.GenerateSignature(header, body)
	if header.Signature != expectedSignature {
		return fmt.Errorf("无效的消息签名")
	}

	return nil
}

// GenerateNonce 生成nonce
func (s *SecurityManager) GenerateNonce() uint64 {
	// 使用时间戳和随机数组合生成nonce
	return uint64(time.Now().UnixNano())
}

// ClearExpiredNonces 清理过期的nonce
func (s *SecurityManager) ClearExpiredNonces() {
	now := time.Now().Unix()
	s.nonceMap.Range(func(key, value interface{}) bool {
		if timestamp, ok := value.(int64); ok {
			if now-timestamp > int64(s.nonceTimeout.Seconds()) {
				s.nonceMap.Delete(key)
			}
		}
		return true
	})
}
