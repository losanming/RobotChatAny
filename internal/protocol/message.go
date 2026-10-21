package protocol

// MessageBody 消息体接口
type MessageBody interface {
	Type() uint8
}

// HeartbeatMessage 心跳消息
type HeartbeatMessage struct {
	Timestamp int64 `json:"timestamp"` // 发送时间戳
}

func (m *HeartbeatMessage) Type() uint8 { return MsgTypeHeartbeat }

// AuthMessage 认证消息
type AuthMessage struct {
	Token     string `json:"token"`     // 认证令牌
	DeviceID  string `json:"deviceId"`  // 设备ID
	Platform  string `json:"platform"`  // 平台类型(ios/android/web)
	Version   string `json:"version"`   // 客户端版本
}

func (m *AuthMessage) Type() uint8 { return MsgTypeAuth }

// AuthAckMessage 认证响应消息
type AuthAckMessage struct {
	Status  uint8  `json:"status"`  // 状态码
	Message string `json:"message"` // 状态信息
	UserID  int64  `json:"userId"`  // 用户ID
}

func (m *AuthAckMessage) Type() uint8 { return MsgTypeAuthAck }

// TextMessage 文本消息
type TextMessage struct {
	From    int64  `json:"from"`    // 发送者ID
	To      int64  `json:"to"`      // 接收者ID
	Content string `json:"content"` // 消息内容
}

func (m *TextMessage) Type() uint8 { return MsgTypeText }

// TextAckMessage 文本消息确认
type TextAckMessage struct {
	MessageID uint64 `json:"messageId"` // 消息ID
	Status    uint8  `json:"status"`    // 状态码
	Message   string `json:"message"`   // 状态信息
}

func (m *TextAckMessage) Type() uint8 { return MsgTypeTextAck }

// ImageMessage 图片消息
type ImageMessage struct {
	From     int64  `json:"from"`     // 发送者ID
	To       int64  `json:"to"`       // 接收者ID
	ImageURL string `json:"imageUrl"` // 图片URL
	Width    int    `json:"width"`    // 图片宽度
	Height   int    `json:"height"`   // 图片高度
}

func (m *ImageMessage) Type() uint8 { return MsgTypeImage }

// ImageAckMessage 图片消息确认
type ImageAckMessage struct {
	MessageID uint64 `json:"messageId"` // 消息ID
	Status    uint8  `json:"status"`    // 状态码
	Message   string `json:"message"`   // 状态信息
}

func (m *ImageAckMessage) Type() uint8 { return MsgTypeImageAck }

// NotifyMessage 系统通知消息
type NotifyMessage struct {
	NotifyType uint8  `json:"type"`    // 通知类型
	Content    string `json:"content"` // 通知内容
}

func (m *NotifyMessage) Type() uint8 { return MsgTypeNotify }

// NotifyAckMessage 系统通知确认
type NotifyAckMessage struct {
	MessageID uint64 `json:"messageId"` // 消息ID
	Status    uint8  `json:"status"`    // 状态码
	Message   string `json:"message"`   // 状态信息
}

func (m *NotifyAckMessage) Type() uint8 { return MsgTypeNotifyAck }
