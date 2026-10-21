package response

import (
	"github.com/gin-gonic/gin"
)

func resp(c *gin.Context, code StatusCode, tips string, count int64, content any) {
	var respData = apiResponse{
		Code:  code,
		Msg:   code.StatusText(),
		Tips:  tips,
		Count: count,
	}
	respData.Data = content
	c.JSON(200, respData)
	return
}

func Resp(c *gin.Context, code StatusCode) {
	resp(c, code, "", 0, nil)
}
func RespTips(c *gin.Context, code StatusCode, tips string) {
	resp(c, code, tips, 0, nil)
}

func RespContent(c *gin.Context, code StatusCode, content any) {
	resp(c, code, "", 0, content)
}

func RespContentTips(c *gin.Context, code StatusCode, tips string, content any) {
	resp(c, code, tips, 0, content)
}
func RespContentWithNum(c *gin.Context, code StatusCode, num int64, content any) {
	resp(c, code, "", num, content)
}

func RespContentWithNumTips(c *gin.Context, code StatusCode, num int64, tips string, content any) {
	resp(c, code, tips, num, content)
}

type apiResponse struct {
	// 错误码
	Code StatusCode `json:"code,omitempty"`
	// 消息短语
	Msg string `json:"msg"`
	// Tips 提示信息
	Tips string `json:"tips,omitempty"`
	// 数据数量
	Count int64 `json:"count,omitempty"`
	// 内容
	Data any `json:"data,omitempty"`
}
