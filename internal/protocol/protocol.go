package protocol

// 消息类型定义
const (
	// 系统消息 (0-99)
	MsgTypeHeartbeat    uint8 = 1  // 心跳
	MsgTypeHeartbeatAck uint8 = 2  // 心跳响应
	MsgTypeAuth         uint8 = 10 // 认证请求
	MsgTypeAuthAck      uint8 = 11 // 认证响应
	MsgTypeKick         uint8 = 20 // 踢出通知

	// 聊天消息 (100-199)
	MsgTypeText     uint8 = 100 // 文本消息
	MsgTypeTextAck  uint8 = 101 // 文本消息确认
	MsgTypeImage    uint8 = 110 // 图片消息
	MsgTypeImageAck uint8 = 111 // 图片消息确认

	// 通知消息 (200-299)
	MsgTypeNotify    uint8 = 200 // 系统通知
	MsgTypeNotifyAck uint8 = 201 // 系统通知确认

	// 设备消息
	DeviceAuto uint8 = 220
)

// 消息状态定义
const (
	MsgStatusSuccess uint8 = 0 // 成功
	MsgStatusFailed  uint8 = 1 // 失败
)

const (
	SessionIDSize = 32 // 会话标识长度
	SignatureSize = 64 // 签名长度
)

// 基础消息头
type MessageHeader struct {
	// 基本信息
	Version   uint8  // 协议版本号
	Type      uint8  // 消息类型
	Sequence  uint32 // 消息序列号
	Timestamp int64  // 消息时间戳
	Length    uint32 // 消息体长度

	// 安全相关
	SessionID [SessionIDSize]byte // 会话标识，用于验证消息来源
	Nonce     uint64              // 随机数，防止重放攻击
	Signature [SignatureSize]byte // 消息签名，验证消息完整性
}

// 错误码定义
const (
	// 系统错误 (0-99)
	ErrSuccess        uint8 = 0 // 成功
	ErrSystemError    uint8 = 1 // 系统错误
	ErrInvalidParam   uint8 = 2 // 参数错误
	ErrInvalidMessage uint8 = 3 // 消息格式错误

	// 认证错误 (100-199)
	ErrAuthFailed   uint8 = 100 // 认证失败
	ErrTokenExpired uint8 = 101 // Token过期
	ErrTokenInvalid uint8 = 102 // Token无效

	// 业务错误 (200-299)
	ErrUserNotFound   uint8 = 200 // 用户不存在
	ErrUserOffline    uint8 = 201 // 用户离线
	ErrMessageTooLong uint8 = 202 // 消息过长
)

// 错误信息映射
var ErrorMessages = map[uint8]string{
	ErrSuccess:        "成功",
	ErrSystemError:    "系统错误",
	ErrInvalidParam:   "参数错误",
	ErrInvalidMessage: "消息格式错误",
	ErrAuthFailed:     "认证失败",
	ErrTokenExpired:   "Token过期",
	ErrTokenInvalid:   "Token无效",
	ErrUserNotFound:   "用户不存在",
	ErrUserOffline:    "用户离线",
	ErrMessageTooLong: "消息过长",
}

// GetErrorMessage 获取错误信息
func GetErrorMessage(code uint8) string {
	if msg, ok := ErrorMessages[code]; ok {
		return msg
	}
	return "未知错误"
}
