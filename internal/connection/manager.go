package connection

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"new-robot-self/config"
	"new-robot-self/internal/ai"
	context_user "new-robot-self/internal/context"
	"new-robot-self/internal/protocol"
	"new-robot-self/internal/queue"
	"new-robot-self/pkg"
	"new-robot-self/pkg/utils"
	public_ai "new-robot-self/public/ai"
	public_contextManager "new-robot-self/public/contextManager"
	public_context_prompt "new-robot-self/public/contextPromptGenerator"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

// Connection 连接信息
type Connection struct {
	ID        string                       // 连接ID
	UserID    int64                        // 用户ID
	Conn      net.Conn                     // TCP连接
	SessionID string                       // 会话ID
	Codec     *protocol.SecureMessageCodec // 安全消息编解码器
	LastSeen  int64                        // 最后活跃时间

	// 超时设置
	ReadTimeout  time.Duration // 读取超时
	WriteTimeout time.Duration // 写入超时
}

// ConnectionManager 连接管理器
type ConnectionManager struct {
	connections sync.Map                     // 所有连接 map[string]*Connection
	userConns   sync.Map                     // 用户连接映射 map[int64]*Connection
	dispatcher  *queue.Dispatcher            // 消息分发器
	codec       *protocol.SecureMessageCodec // 安全消息编解码器
	logger      *log.Logger                  // 日志记录器

	// 配置
	readTimeout  time.Duration // 读取超时
	writeTimeout time.Duration // 写入超时
}

/*
	因为出现了导包问题,我这边不得不采用接口的方式解决问题,好的处理方法还是没有想出来
*/

// 定义回调函数类型
type HeartbeatUpdateFunc func(connID string)
type AuthHandlerFunc func(connID string, userID int64) error
type EncodeWithSecurityFunc func(sessionID string, msgType uint8, payload interface{}) ([]byte, error)
type GeneratePromptFunc func(userID string, currentInput string) (string, error)

// 全局回调函数变量
var (
	GlobalHeartbeatUpdater   HeartbeatUpdateFunc
	GlobalAuthHandler        AuthHandlerFunc
	GlobalEncodeWithSecurity EncodeWithSecurityFunc
	GlobalGeneratePrompt     GeneratePromptFunc
)

// 注册函数
func RegisterHeartbeatUpdater(updater HeartbeatUpdateFunc) {
	GlobalHeartbeatUpdater = updater
}

func RegisterAuthHandler(handler AuthHandlerFunc) {
	GlobalAuthHandler = handler
}
func RegisterEncodeWithSecurity(updater EncodeWithSecurityFunc) {
	GlobalEncodeWithSecurity = updater
}

func RegisterGeneratePrompt(updater GeneratePromptFunc) {
	GlobalGeneratePrompt = updater
}

// NewConnectionManager 创建连接管理器
func NewConnectionManager(dispatcher *queue.Dispatcher, secretKey []byte, readTimeout, writeTimeout int64) *ConnectionManager {
	return &ConnectionManager{
		dispatcher:   dispatcher,
		codec:        protocol.NewSecureMessageCodec(secretKey),
		logger:       log.New(os.Stderr, "[ConnectionManager] ", log.LstdFlags),
		readTimeout:  time.Duration(readTimeout) * time.Second,
		writeTimeout: time.Duration(writeTimeout) * time.Second,
	}
}

// AddConnection 添加新连接
func (m *ConnectionManager) AddConnection(conn net.Conn) *Connection {
	connID := pkg.GenerateConnID()

	// 设置连接超时
	if m.readTimeout > 0 {
		conn.SetReadDeadline(time.Now().Add(m.readTimeout))
	}
	if m.writeTimeout > 0 {
		conn.SetWriteDeadline(time.Now().Add(m.writeTimeout))
	}

	connection := &Connection{
		ID:           connID,
		Conn:         conn,
		SessionID:    fmt.Sprintf("session_%s", connID),
		Codec:        m.codec,
		LastSeen:     time.Now().Unix(),
		ReadTimeout:  m.readTimeout,
		WriteTimeout: m.writeTimeout,
	}
	m.connections.Store(connID, connection)
	return connection
}

func (m *ConnectionManager) GetLogger() *log.Logger {
	return m.logger
}

func (m *ConnectionManager) GetCodec() *protocol.SecureMessageCodec {
	return m.codec
}

// AuthenticateConnection 认证连接
func (m *ConnectionManager) AuthenticateConnection(connID string, userID int64) error {
	// 1. 获取连接
	conn, ok := m.getConnection(connID)
	if !ok {
		return fmt.Errorf("连接不存在: %s", connID)
	}

	// 2. 检查用户是否已经有其他连接
	if oldConn, ok := m.getUserConnection(userID); ok {
		// 关闭旧连接
		m.RemoveConnection(oldConn.ID)
	}

	// 3. 更新连接信息
	conn.UserID = userID
	m.userConns.Store(userID, conn)

	// 4. 添加到消息分发系统
	//m.dispatcher.AddUser(userID)

	return nil
}

// RemoveConnection 移除连接
func (m *ConnectionManager) RemoveConnection(connID string) {
	if conn, ok := m.getConnection(connID); ok {
		// 1. 关闭连接
		conn.Conn.Close()

		// 2. 如果是认证用户，从用户连接映射中移除
		if conn.UserID > 0 {
			m.userConns.Delete(conn.UserID)
			// 从消息分发系统移除
			m.dispatcher.RemoveUser(conn.UserID)
		}

		// 3. 从连接列表移除
		m.connections.Delete(connID)
	}
}

// SendMessage 发送消息给指定用户
func (m *ConnectionManager) SendMessage(ctx context.Context, msg *queue.Message) error {
	// 1. 获取用户连接
	conn, ok := m.getUserConnection(msg.To)
	if !ok {
		return fmt.Errorf("用户不在线: %d", msg.To)
	}

	// 2. 将消息转换为协议消息
	textMsg := &protocol.TextMessage{
		From:    msg.From,
		To:      msg.To,
		Content: string(msg.Content),
	}

	// 3. 编码消息
	data, err := conn.Codec.EncodeWithSecurity(conn.SessionID, msg.Type, textMsg)
	if err != nil {
		return fmt.Errorf("编码消息失败: %v", err)
	}

	// 4. 设置写入超时
	if conn.WriteTimeout > 0 {
		conn.Conn.SetWriteDeadline(time.Now().Add(conn.WriteTimeout))
	}

	// 5. 发送消息
	_, err = conn.Conn.Write(data)
	if err != nil {
		m.logger.Printf("发送消息失败: %v", err)
		// 发送失败,可能连接已断开,移除连接
		m.RemoveConnection(conn.ID)
		return fmt.Errorf("发送消息失败: %v", err)
	}

	return nil
}

// UpdateLastSeen 更新最后活跃时间
func (m *ConnectionManager) UpdateLastSeen(connID string) {
	if conn, ok := m.getConnection(connID); ok {
		conn.LastSeen = time.Now().Unix()
	}
}

// CleanInactiveConnections 清理不活跃连接
func (m *ConnectionManager) CleanInactiveConnections(timeout int64) {
	now := time.Now().Unix()
	m.connections.Range(func(key, value interface{}) bool {
		conn := value.(*Connection)
		if now-conn.LastSeen > timeout {
			m.RemoveConnection(conn.ID)
		}
		return true
	})
}

// getConnection 获取连接
func (m *ConnectionManager) getConnection(connID string) (*Connection, bool) {
	if value, ok := m.connections.Load(connID); ok {
		return value.(*Connection), true
	}
	return nil, false
}

// getUserConnection 获取用户连接
func (m *ConnectionManager) getUserConnection(userID int64) (*Connection, bool) {
	if value, ok := m.userConns.Load(userID); ok {
		return value.(*Connection), true
	}
	return nil, false
}
func (m *ConnectionManager) GetDispatcher() (*queue.Dispatcher, bool) {
	if m.dispatcher != nil {
		return m.dispatcher, true
	}
	return nil, false
}

// 获取所有连接上的用户id
func (m *ConnectionManager) GetUserIDs() []int64 {
	var userIDs []int64
	m.userConns.Range(func(key, value interface{}) bool {
		userIDs = append(userIDs, key.(int64))
		return true
	})
	return userIDs
}

// HandleConnection 处理新的TCP连接
func (m *ConnectionManager) HandleConnection(conn net.Conn) {
	// 添加连接
	connection := m.AddConnection(conn)

	m.logger.Printf("新连接已建立: %s", connection.ID)

	// 启动读取循环
	go m.readLoop(connection)
}

// readLoop 读取连接数据
func (m *ConnectionManager) readLoop(conn *Connection) {
	defer func() {
		m.RemoveConnection(conn.ID)
		m.logger.Printf("连接已关闭: %s", conn.ID)
	}()

	buffer := make([]byte, 4096)
	for {
		// 设置读取超时
		if conn.ReadTimeout > 0 {
			conn.Conn.SetReadDeadline(time.Now().Add(conn.ReadTimeout))
		}

		// 读取数据
		n, err := conn.Conn.Read(buffer)
		if err != nil {
			m.logger.Printf("连接 %s 读取错误: %v", conn.ID, err)
			return
		}

		if n > 0 {
			// 处理接收到的数据
			data := buffer[:n]

			// 解码消息
			header, msgBody, err := conn.Codec.DecodeWithSecurity(data)
			if err != nil {
				m.logger.Printf("解码消息失败: %v", err)
				continue
			}
			marshal, _ := json.Marshal(msgBody)
			if marshal == nil {
				continue
			}
			// 处理消息
			m.HandleMessage(conn, header.Type, marshal)

			// 更新最后活跃时间
			m.UpdateLastSeen(conn.ID)
		}
	}
}

// CloseAllConnections 关闭所有连接
func (m *ConnectionManager) CloseAllConnections() {
	m.connections.Range(func(key, value interface{}) bool {
		if conn, ok := value.(*Connection); ok {
			conn.Conn.Close()
			m.RemoveConnection(conn.ID)
		}
		return true
	})
}

// handleMessage 处理接收到的消息
func (m *ConnectionManager) HandleMessage(conn *Connection, msgType uint8, msgBody []byte) {
	// 根据消息类型处理
	switch msgType {
	case protocol.MsgTypeAuth:
		if GlobalAuthHandler == nil {
			m.logger.Printf("未注册认证处理器")
			return
		}
		// 处理认证消息
		var authMsg protocol.AuthMessage
		if err := json.Unmarshal(msgBody, &authMsg); err != nil {
			m.logger.Printf("解析认证消息失败: %v", err)
			return
		}

		// TODO: 实现认证逻辑 token的校验,token为当前分钟时间戳 的 上下100范围 + dangwing + 设备编号12位
		code, token := CheckToken(authMsg.Token)
		if !token {
			m.logger.Printf("编码消息失败")
			return
		}
		var user_id int64
		conn.SessionID = pkg.GenerateConnID()
		user_id = code
		err := GlobalAuthHandler(conn.ID, user_id)
		if err != nil {
			m.logger.Printf("认证失败: %v", err)
			return
		}
		m.logger.Printf("收到认证消息: %+v", authMsg)
		var sendAuthMsg protocol.AuthAckMessage
		security, err := GlobalEncodeWithSecurity(conn.SessionID, protocol.MsgTypeAuth, sendAuthMsg)
		if err != nil {
			m.logger.Printf("编码消息失败: %v", err)
			return
		}
		conn.Conn.Write(security)
	case protocol.MsgTypeHeartbeat:
		// 处理心跳消息
		var heartbeatMsg protocol.HeartbeatMessage
		if err := json.Unmarshal(msgBody, &heartbeatMsg); err != nil {
			m.logger.Printf("解析心跳消息失败: %v", err)
			return
		}

		// 更新心跳时间
		m.logger.Printf("收到心跳消息: %+v", heartbeatMsg)
		if GlobalHeartbeatUpdater != nil {
			GlobalHeartbeatUpdater(conn.ID)
		} else {
			m.logger.Printf("GlobalHeartbeatUpdater is nil")
			return
		}
	case protocol.MsgTypeText:
		// 处理文本消息
		var textMsg protocol.TextMessage
		if err := json.Unmarshal(msgBody, &textMsg); err != nil {
			m.logger.Printf("解析文本消息失败: %v", err)
			return
		}
		userIDStr := fmt.Sprintf("%d", conn.UserID)

		// 生成个性化promp
		prompt, err := GlobalGeneratePrompt(userIDStr, textMsg.Content)
		if err != nil {
			log.Printf("生成prompt失败: %v", err)
			return
		}

		//ai处理
		chat, err := public_ai.AiService.Chat(context.Background(), ai.Message{
			Role:    "user",
			Content: prompt,
		})
		if err != nil {
			m.logger.Printf("ai处理失败: %s", err)
			return
		}
		fmt.Sprintf("ai处理成功: %s", chat)
		//添加用户消息和ai消息到消息缓存
		err = public_contextManager.ContextManager.AddMessage(userIDStr, context_user.Message{
			Content:  textMsg.Content,
			Time:     time.Now().Unix(),
			Priority: 2, // AI回答的优先级设为中等
			Type:     "user",
			Keywords: public_context_prompt.ContextPrompt.ExtractKeywords(textMsg.Content), // 从回答中提取关键词
		})
		if err != nil {
			log.Printf("添加用户消息失败: %v", err)
		}
		err = public_context_prompt.ContextPrompt.RecordAIResponse(userIDStr, chat)
		if err != nil {
			log.Printf("记录AI回答失败: %v", err)
		}
		var SendTextMsg protocol.TextMessage
		SendTextMsg.Content = chat
		SendTextMsg.From = textMsg.From
		SendTextMsg.To = textMsg.To
		marshal, _ := json.Marshal(&SendTextMsg)
		if marshal == nil {
			return
		}
		security, err := m.codec.EncodeWithSecurity(conn.SessionID, protocol.MsgTypeText, SendTextMsg)
		if err != nil {
			m.logger.Printf("编码消息失败: %v", err)
			return
		}
		conn.Conn.Write(security)
	default:
		m.logger.Printf("收到未知类型消息: %d", msgType)
	}
}

func CheckToken(token string) (int64, bool) {
	key := []byte(config.GetConfig().TCP.SecretKey)
	decrypt, err := utils.AesDecrypt(token, key)
	if err != nil {
		return 0, false
	}
	split := strings.Split(string(decrypt), "+")
	if len(split) != 3 {
		return 0, false
	}
	if split[1] != "dangwing" {
		return 0, false
	}
	if len(split[2]) != 12 {
		return 0, false
	}
	atoi, err := strconv.Atoi(split[0])
	if atoi == 0 || err != nil {
		return 0, false
	}
	//设备编号

	code, err := strconv.Atoi(split[2])
	if code == 0 || err != nil {
		return 0, false
	}

	// 获取当前分钟时间戳
	now := time.Now()
	loc := time.FixedZone("CST", 8*3600)
	currentMinute := time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), 0, 0, loc)
	// 获取时间戳（单位：秒）
	timestamp := currentMinute.Unix()

	// 判断是否在时间戳的上下100范围内
	if (atoi >= int(timestamp)-100) && (atoi <= int(timestamp)+100) {
		//@@TODO 校验code是否有效
		return int64(code), true
	}
	return 0, false
}
