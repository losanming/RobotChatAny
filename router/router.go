package router

import (
	"github.com/gin-gonic/gin"
	"new-robot-self/pkg/response"
	authRouter "new-robot-self/router/auth"
	systemRouter "new-robot-self/router/system"
	userRouter "new-robot-self/router/user"
)

// SetupRouter 设置路由
func SetupRouter(r *gin.Engine) {
	// 健康检查
	r.GET("/ping", func(c *gin.Context) {
		response.Resp(c, response.StatusOK)
	})
	group := r.Group("api/v1.0")
	// 初始化用户路由
	userRouterGroup := &userRouter.UserRouter{}
	authRouterGroup := &authRouter.AuthRouter{}
	systemRouterGroup := &systemRouter.SystemRouter{}

	//用户相关
	userRouterGroup.InitUserRouter(group)

	//认证相关
	authRouterGroup.InitAuthRouter(group)

	//系统相关
	systemRouterGroup.InitSystemRouter(group)
}
