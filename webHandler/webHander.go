package webHandler

import (
	"new-robot-self/middleware"
	"new-robot-self/pkg/web"
	"new-robot-self/router"
)

var basePath = "/api/v1.0"

func InitRouter(w *web.Server) {
	// server init
	whiteList := middleware.NewWhiteListConfig()

	// 添加该服务特定的白名单路径
	whiteList.Add(
		basePath+"/auth/captcha",
		basePath+"/system/version",
		"ping",
	)
	w.Engine.Use(middleware.Cors(), middleware.JWTAuth(whiteList))
	router.SetupRouter(w.Engine)
}
