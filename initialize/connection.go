package initialize

import (
	"errors"
	"fmt"
	"new-robot-self/internal/connection"
	"new-robot-self/public"
	public_connectionManager "new-robot-self/public/connectionManager"
	public_heartbeatManager "new-robot-self/public/heartbeatManager"
)

func InitConnectionAndHeartbeat() error {
	manager := connection.NewConnectionManager(public.Dispatcher, []byte(public.Config.TCP.SecretKey), 4, 5)
	if manager == nil {
		return errors.New("connection manager is nil")
	}
	heartbeatManager := connection.NewHeartbeatManager(manager, &connection.DefaultHeartbeatConfig)
	if heartbeatManager == nil {
		return errors.New("heartbeat manager is nil")
	}
	heartbeatManager.Start()

	connection.RegisterHeartbeatUpdater(heartbeatManager.UpdateHeartbeat)
	connection.RegisterAuthHandler(manager.AuthenticateConnection)
	connection.RegisterEncodeWithSecurity(manager.GetCodec().EncodeWithSecurity)
	connection.RegisterGeneratePrompt(public.ContextPromptGenerator.GeneratePrompt)

	public_connectionManager.Connection = manager
	public_heartbeatManager.HeartbeatManager = heartbeatManager
	fmt.Printf("connection 初始化配置成功 \n")
	return nil
}
