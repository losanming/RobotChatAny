package response

import (
	"github.com/gin-gonic/gin"
)

type CB func(c *gin.Context)

func API(cb CB) gin.HandlerFunc {
	return func(c *gin.Context) {
		cb(c)
	}
}
