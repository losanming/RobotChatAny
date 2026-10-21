package initialize

import (
	"new-robot-self/config"
	"new-robot-self/pkg/web"
	"new-robot-self/public"
	"new-robot-self/webHandler"
)

func InitWeb() error {
	var c web.WebConfig
	conf := config.GetConfig().Web
	c.Port = conf.Port
	c.MaxHeaderBytes = conf.MaxHeaderBytes
	c.ReadTimeout = conf.ReadTimeout
	c.WriteTimeout = conf.WriteTimeout
	server := web.NewServer(c)
	public.Web = server
	err := RunWeb()
	return err
}

func RunWeb() error {
	webHandler.InitRouter(public.Web)
	return public.Web.Run()
}
