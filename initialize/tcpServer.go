package initialize

import (
	"errors"
	"fmt"
	"new-robot-self/internal/tcp"
	"new-robot-self/public"
	public_connectionManager "new-robot-self/public/connectionManager"
)

func InitTcpServer() error {
	var tcp_config = tcp.Config{
		Port:              public.Config.TCP.Port,
		ReadTimeout:       int64(public.Config.TCP.ReadTimeout),
		WriteTimeout:      int64(public.Config.TCP.WriteTimeout),
		MaxConnections:    public.Config.TCP.MaxConnections,
		HeartbeatInterval: public.Config.TCP.HeartbeatInterval,
		HeartbeatTimeout:  public.Config.TCP.HeartbeatTimeout,
		TLS: struct {
			Enabled  bool   `yaml:"enabled"`
			CertFile string `yaml:"cert_file"`
			KeyFile  string `yaml:"key_file"`
		}{
			Enabled:  public.Config.TCP.TLS.Enabled,
			CertFile: public.Config.TCP.TLS.CertFile,
			KeyFile:  public.Config.TCP.TLS.KeyFile,
		},
	}
	server := tcp.NewServer(&tcp_config, public_connectionManager.Connection)
	if server == nil {
		return errors.New("tcp server is nil")
	}
	err := server.Start()
	if err != nil {
		return err
	}
	public.TcpServer = server
	fmt.Printf("tcp 初始化配置成功 \n")
	return nil
}
