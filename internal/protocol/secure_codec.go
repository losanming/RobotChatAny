package protocol

import (
	"fmt"
	"time"
)

// SecureMessageCodec 安全消息编解码器
type SecureMessageCodec struct {
	codec    *MessageCodec    // 基础编解码器
	security *SecurityManager // 安全管理器
}

// NewSecureMessageCodec 创建安全消息编解码器
func NewSecureMessageCodec(secretKey []byte) *SecureMessageCodec {
	return &SecureMessageCodec{
		codec:    &MessageCodec{},
		security: NewSecurityManager(secretKey, 5*time.Minute),
	}
}

// EncodeWithSecurity 安全编码消息
func (c *SecureMessageCodec) EncodeWithSecurity(sessionID string, msgType uint8, payload interface{}) ([]byte, error) {
	// 1. 生成nonce
	nonce := c.security.GenerateNonce()

	// 2. 获取当前时间戳
	timestamp := time.Now().Unix()

	// 3. 序列化消息体
	body, err := c.codec.encodeBody(payload)
	if err != nil {
		return nil, fmt.Errorf("序列化消息体失败: %v", err)
	}

	// 4. 构建消息头
	header := &MessageHeader{
		Version:   1,
		Type:      msgType,
		Sequence:  generateSequence(),
		Timestamp: timestamp,
		Length:    uint32(len(body)),
		Nonce:     nonce,
	}

	// 5. 设置会话标识
	copy(header.SessionID[:], StringToFixedBytes(sessionID, SessionIDSize))

	// 6. 计算签名
	header.Signature = c.security.GenerateSignature(header, body)

	// 7. 编码消息头
	data, err := c.codec.encodeHeader(header)
	if err != nil {
		return nil, fmt.Errorf("编码消息头失败: %v", err)
	}

	// 8. 组合消息头和消息体
	if body != nil {
		data = append(data, body...)
	}

	return data, nil
}

// DecodeWithSecurity 安全解码消息
func (c *SecureMessageCodec) DecodeWithSecurity(data []byte) (*MessageHeader, interface{}, error) {
	// 1. 解码消息
	header, payload, err := c.codec.Decode(data)
	if err != nil {
		return nil, nil, err
	}

	// 2. 提取消息体
	bodyStart := c.codec.headerSize()
	var body []byte
	if len(data) > bodyStart {
		body = data[bodyStart:]
	}

	// 3. 验证消息安全性
	if err := c.security.ValidateMessage(header, body); err != nil {
		return nil, nil, fmt.Errorf("安全验证失败: %v", err)
	}

	return header, payload, nil
}
