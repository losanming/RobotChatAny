package protocol

import (
	"fmt"
	"testing"
)

func BenchmarkCodec(b *testing.B) {
	// 准备测试数据
	msg := &TextMessage{
		From:    1001,
		To:      1002,
		Content: "这是一条测试消息，用于性能测试。这是一条测试消息，用于性能测试。",
	}
	
	// 普通编解码器
	normalCodec := &MessageCodec{}
	
	// 运行普通编解码基准测试
	b.Run("Normal", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			data, _ := normalCodec.Encode(MsgTypeText, msg)
			_, _, _ = normalCodec.Decode(data)
		}
	})
	
	// 安全编解码器
	secureCodec := NewSecureMessageCodec([]byte("test-secret-key-123"))
	sessionID := "test-session-001"
	
	// 运行安全编解码基准测试
	b.Run("Secure", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			data, _ := secureCodec.EncodeWithSecurity(sessionID, MsgTypeText, msg)
			_, _, _ = secureCodec.DecodeWithSecurity(data)
		}
	})
}

func BenchmarkHMAC(b *testing.B) {
	// 准备测试数据
	msg := &TextMessage{
		From:    1001,
		To:      1002,
		Content: "这是一条测试消息，用于性能测试。这是一条测试消息，用于性能测试。",
	}
	
	// 安全管理器
	security := NewSecurityManager([]byte("test-secret-key-123"), 0)
	
	// 准备消息头
	header := &MessageHeader{
		Version:   1,
		Type:      MsgTypeText,
		Sequence:  1,
		Timestamp: 1234567890,
		Length:    100,
		Nonce:     12345,
	}
	copy(header.SessionID[:], StringToFixedBytes("test-session-001", SessionIDSize))
	
	// 序列化消息体
	body, _ := (&MessageCodec{}).encodeBody(msg)
	
	// 运行HMAC签名基准测试
	b.Run("HMAC", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = security.GenerateSignature(header, body)
		}
	})
}

func BenchmarkNonceVerification(b *testing.B) {
	// 安全管理器
	security := NewSecurityManager([]byte("test-secret-key-123"), 0)
	
	// 准备测试数据
	sessionID := "test-session-001"
	
	// 运行Nonce验证基准测试
	b.Run("NonceVerify", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			nonceKey := fmt.Sprintf("%s:%d", sessionID, i)
			security.nonceMap.LoadOrStore(nonceKey, int64(i))
		}
	})
}
