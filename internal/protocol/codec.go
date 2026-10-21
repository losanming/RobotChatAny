package protocol

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"sync/atomic"
	"time"
)

// MessageCodec 消息编解码器
type MessageCodec struct{}

// headerSize 返回消息头大小
func (c *MessageCodec) headerSize() int {
	return binary.Size(MessageHeader{})
}

// encodeHeader 编码消息头
func (c *MessageCodec) encodeHeader(header *MessageHeader) ([]byte, error) {
	buf := new(bytes.Buffer)
	if err := binary.Write(buf, binary.BigEndian, header); err != nil {
		return nil, fmt.Errorf("编码消息头失败: %v", err)
	}
	return buf.Bytes(), nil
}

// encodeBody 编码消息体
func (c *MessageCodec) encodeBody(payload interface{}) ([]byte, error) {
	if payload == nil {
		return nil, nil
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("序列化消息体失败: %v", err)
	}
	return body, nil
}

// Encode 编码消息
func (c *MessageCodec) Encode(msgType uint8, payload interface{}) ([]byte, error) {
	// 1. 序列化消息体
	var body []byte
	var err error
	if payload != nil {
		body, err = json.Marshal(payload)
		if err != nil {
			return nil, fmt.Errorf("序列化消息体失败: %v", err)
		}
	}

	// 2. 构建消息头
	header := MessageHeader{
		Version:   1,
		Type:      msgType,
		Sequence:  generateSequence(),
		Timestamp: time.Now().Unix(),
		Length:    uint32(len(body)),
	}

	// 3. 编码消息头
	buf := new(bytes.Buffer)
	if err := binary.Write(buf, binary.BigEndian, header); err != nil {
		return nil, fmt.Errorf("编码消息头失败: %v", err)
	}

	// 4. 写入消息体
	if body != nil {
		buf.Write(body)
	}

	return buf.Bytes(), nil
}

// Decode 解码消息
func (c *MessageCodec) Decode(data []byte) (*MessageHeader, interface{}, error) {
	if len(data) < binary.Size(MessageHeader{}) {
		return nil, nil, fmt.Errorf("消息太短")
	}

	// 1. 解码消息头
	header := &MessageHeader{}
	headerBuf := bytes.NewReader(data[:binary.Size(MessageHeader{})])
	if err := binary.Read(headerBuf, binary.BigEndian, header); err != nil {
		return nil, nil, fmt.Errorf("解码消息头失败: %v", err)
	}

	// 2. 检查消息体长度
	if uint32(len(data)) < uint32(binary.Size(MessageHeader{}))+header.Length {
		return nil, nil, fmt.Errorf("消息不完整")
	}

	// 3. 根据消息类型解析消息体
	var payload interface{}
	if header.Length > 0 {
		body := data[binary.Size(MessageHeader{}) : binary.Size(MessageHeader{})+int(header.Length)]

		switch header.Type {
		case MsgTypeHeartbeat:
			payload = &HeartbeatMessage{}
		case MsgTypeAuth:
			payload = &AuthMessage{}
		case MsgTypeAuthAck:
			payload = &AuthAckMessage{}
		case MsgTypeText:
			payload = &TextMessage{}
		case MsgTypeTextAck:
			payload = &TextAckMessage{}
		case MsgTypeNotify:
			payload = &NotifyMessage{}
		default:
			return nil, nil, fmt.Errorf("未知的消息类型: %d", header.Type)
		}

		if err := json.Unmarshal(body, payload); err != nil {
			return nil, nil, fmt.Errorf("解析消息体失败: %v", err)
		}
	}

	return header, payload, nil
}

// 生成序列号
var sequence uint32 = 0

func generateSequence() uint32 {
	return atomic.AddUint32(&sequence, 1)
}
