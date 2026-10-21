package tcp

import (
	"crypto/tls"
	"fmt"
	"log"
	"net"
	"new-robot-self/internal/connection"
	"os"
	"sync"
	"sync/atomic"
)

// Server TCP服务器结构体
type Server struct {
	// 配置相关
	config *Config

	// 服务器状态
	running atomic.Bool

	// 连接管理
	listener    net.Listener
	connCount   int32
	connManager *connection.ConnectionManager

	// 日志
	logger *log.Logger

	// 关闭控制
	closeOnce sync.Once
	closeChan chan struct{}
}

// Config TCP服务器配置
type Config struct {
	Port              int   `yaml:"port"`
	ReadTimeout       int64 `yaml:"read_timeout"`
	WriteTimeout      int64 `yaml:"write_timeout"`
	MaxConnections    int   `yaml:"max_connections"`
	HeartbeatInterval int   `yaml:"heartbeat_interval"`
	HeartbeatTimeout  int   `yaml:"heartbeat_timeout"`

	TLS struct {
		Enabled  bool   `yaml:"enabled"`
		CertFile string `yaml:"cert_file"`
		KeyFile  string `yaml:"key_file"`
	} `yaml:"tls"`
}

// NewServer 创建新的TCP服务器
func NewServer(config *Config, connManager *connection.ConnectionManager) *Server {
	return &Server{
		config:      config,
		closeChan:   make(chan struct{}),
		connManager: connManager,
		logger:      log.New(os.Stderr, "[TCPServer] ", log.LstdFlags),
	}
}

// Start 启动TCP服务器
func (s *Server) Start() error {
	// 已经在运行
	if !s.running.CompareAndSwap(false, true) {
		return fmt.Errorf("server is already running")
	}

	var listener net.Listener
	var err error

	// 根据配置决定是否使用TLS
	if s.config.TLS.Enabled {
		// 加载TLS证书
		cert, err := tls.LoadX509KeyPair(s.config.TLS.CertFile, s.config.TLS.KeyFile)
		if err != nil {
			return fmt.Errorf("failed to load TLS cert: %v", err)
		}

		// TLS配置
		tlsConfig := &tls.Config{
			Certificates: []tls.Certificate{cert},
			MinVersion:   tls.VersionTLS12,
		}

		// 创建TLS监听器
		listener, err = tls.Listen("tcp", fmt.Sprintf(":%d", s.config.Port), tlsConfig)
	} else {
		// 创建普通TCP监听器
		listener, err = net.Listen("tcp", fmt.Sprintf(":%d", s.config.Port))
	}

	if err != nil {
		s.running.Store(false)
		return fmt.Errorf("failed to start listener: %v", err)
	}

	s.listener = listener
	s.logger.Printf("TCP server listening on port %d, TLS enabled: %v", s.config.Port, s.config.TLS.Enabled)

	// 开始接受连接
	go s.acceptLoop()

	return nil
}

// acceptLoop 接受新的连接
func (s *Server) acceptLoop() {
	for {
		select {
		case <-s.closeChan:
			return
		default:
			conn, err := s.listener.Accept()
			if err != nil {
				if s.running.Load() {
					s.logger.Printf("Failed to accept connection: %v", err)
				}
				continue
			}

			// 检查连接数限制
			if atomic.LoadInt32(&s.connCount) >= int32(s.config.MaxConnections) {
				conn.Close()
				s.logger.Printf("Rejected connection: max connections reached")
				continue
			}

			// 交给连接管理器处理
			s.connManager.HandleConnection(conn)
			atomic.AddInt32(&s.connCount, 1)
		}
	}
}

// Stop 停止TCP服务器
func (s *Server) Stop() error {
	s.closeOnce.Do(func() {
		s.running.Store(false)
		close(s.closeChan)

		if s.listener != nil {
			s.listener.Close()
		}

		// 关闭所有连接
		s.connManager.CloseAllConnections()
	})

	return nil
}
