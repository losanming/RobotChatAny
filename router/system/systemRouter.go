package systemRouter

import (
	"github.com/gin-gonic/gin"
	"new-robot-self/pkg/response"
	"new-robot-self/public"
)

type SystemRouter struct{}

func (s *SystemRouter) InitSystemRouter(r *gin.RouterGroup) {
	systemGroup := r.Group("/system")
	{
		systemGroup.GET("/version", func(c *gin.Context) {
			response.RespContent(c, response.StatusOK, public.Config.Sys.Version)
		})
	}
}
