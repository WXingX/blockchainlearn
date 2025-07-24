package response

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
	TraceID string      `json:"trace_id"`
}

func Success(c *gin.Context, data interface{}) {
	traceID := c.GetString("TraceID")
	c.JSON(http.StatusOK, Response{
		Code:    0,
		Message: "success",
		Data:    data,
		TraceID: traceID,
	})
}

func Fail(c *gin.Context, httpCode int, code int, msgKey string) {
	traceID := c.GetString("TraceID")
	c.JSON(httpCode, Response{
		Code:    code,
		Message: msgKey,
		Data:    nil,
		TraceID: traceID,
	})
}
