package middleware

import (
	"blog-management/logger"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"time"
)

func LogMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		latency := time.Since(start)
		statusCode := c.Writer.Status()
		traceID := c.GetString("TraceID")

		logger.Logger.Info("request log",
			zap.String("trace_id", traceID),
			zap.String("method", c.Request.Method),
			zap.String("path", c.FullPath()),
			zap.Int("status_code", statusCode),
			zap.Duration("latency", latency),
		)
	}
}
