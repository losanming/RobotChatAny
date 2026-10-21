package connection

import (
	"context"
	"fmt"
	"new-robot-self/internal/protocol"
	"sync"
	"time"
)

// HeartbeatConfig 心跳配置
type HeartbeatConfig struct {
	Interval       time.Duration // 心跳间隔
	Timeout        time.Duration // 超时时间
	MaxMissedBeats int           // 最大允许丢失心跳次数
}

// DefaultHeartbeatConfig 默认心跳配置
var DefaultHeartbeatConfig = HeartbeatConfig{
	Interval:       time.Second * 30,
	Timeout:        time.Second * 90,
	MaxMissedBeats: 3,
}

// HeartbeatManager 心跳管理器
type HeartbeatManager struct {
	config    HeartbeatConfig
	connMgr   *ConnectionManager
	beatTimes sync.Map // 记录最后心跳时间 map[string]int64
	ctx       context.Context
	cancel    context.CancelFunc
}

// NewHeartbeatManager 创建心跳管理器
func NewHeartbeatManager(connMgr *ConnectionManager, config *HeartbeatConfig) *HeartbeatManager {
	if config == nil {
		config = &DefaultHeartbeatConfig
	}

	ctx, cancel := context.WithCancel(context.Background())
	return &HeartbeatManager{
		config:  *config,
		connMgr: connMgr,
		ctx:     ctx,
		cancel:  cancel,
	}
}

// NewHeartbeatManagerWithParams 使用直接参数创建心跳管理器
func NewHeartbeatManagerWithParams(connMgr *ConnectionManager, intervalSeconds int, timeoutSeconds int) *HeartbeatManager {
	config := HeartbeatConfig{
		Interval:       time.Duration(intervalSeconds) * time.Second,
		Timeout:        time.Duration(timeoutSeconds) * time.Second,
		MaxMissedBeats: 3,
	}

	ctx, cancel := context.WithCancel(context.Background())
	return &HeartbeatManager{
		config:  config,
		connMgr: connMgr,
		ctx:     ctx,
		cancel:  cancel,
	}
}

// Start 启动心跳管理
func (h *HeartbeatManager) Start() {
	// 启动心跳检查协程
	go h.checkHeartbeats()
}

// Stop 停止心跳管理
func (h *HeartbeatManager) Stop() {
	h.cancel()
}

// UpdateHeartbeat 更新心跳时间
func (h *HeartbeatManager) UpdateHeartbeat(connID string) {
	h.beatTimes.Store(connID, time.Now().Unix())
}

// checkHeartbeats 检查心跳状态
func (h *HeartbeatManager) checkHeartbeats() {
	ticker := time.NewTicker(h.config.Interval)
	defer ticker.Stop()

	for {
		select {
		case <-h.ctx.Done():
			return
		case <-ticker.C:
			now := time.Now().Unix()
			h.beatTimes.Range(func(key, value interface{}) bool {
				connID := key.(string)
				lastBeat := value.(int64)

				// 检查是否超时
				if now-lastBeat > int64(h.config.Timeout.Seconds()) {
					// 移除连接
					h.connMgr.RemoveConnection(connID)
					h.beatTimes.Delete(connID)
				}
				return true
			})
		}
	}
}

// SendHeartbeat 发送心跳包
func (h *HeartbeatManager) SendHeartbeat(conn *Connection) error {
	// 创建心跳消息
	heartbeat := &protocol.HeartbeatMessage{
		Timestamp: time.Now().Unix(),
	}

	// 编码消息
	data, err := conn.Codec.EncodeWithSecurity(conn.SessionID, protocol.MsgTypeHeartbeat, heartbeat)
	if err != nil {
		return fmt.Errorf("编码心跳消息失败: %v", err)
	}

	// 发送消息
	_, err = conn.Conn.Write(data)
	if err != nil {
		return fmt.Errorf("发送心跳消息失败: %v", err)
	}

	return nil
}

// StartHeartbeat 启动心跳发送
func (h *HeartbeatManager) StartHeartbeat(conn *Connection) {
	go func() {
		ticker := time.NewTicker(h.config.Interval)
		defer ticker.Stop()

		for {
			select {
			case <-h.ctx.Done():
				return
			case <-ticker.C:
				// 尝试发送心跳,最多重试3次
				var err error
				for retries := 0; retries < 3; retries++ {
					if err = h.SendHeartbeat(conn); err == nil {
						break
					}
					h.connMgr.logger.Printf("发送心跳失败 [ConnID=%s] 重试(%d/3): %v", conn.ID, retries+1, err)
					// 等待1秒后重试
					time.Sleep(time.Second)
				}

				// 如果3次重试都失败了,才移除连接
				if err != nil {
					h.connMgr.logger.Printf("发送心跳最终失败 [ConnID=%s]: %v", conn.ID, err)
					h.connMgr.RemoveConnection(conn.ID)
					return
				}
			}
		}
	}()
}
