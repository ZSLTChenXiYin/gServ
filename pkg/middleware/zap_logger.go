package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func ZapLogger(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery
		request_id := c.GetHeader("X-Request-Id")

		// 处理请求
		c.Next()

		end := time.Now()
		latency := end.Sub(start)

		if len(c.Errors) > 0 {
			// 如果有错误，记录错误信息
			logger.Error("HTTP Request",
				zap.Int("status", c.Writer.Status()),
				zap.String("method", c.Request.Method),
				zap.String("path", path),
				zap.String("query", query),
				zap.String("ip", c.ClientIP()),
				zap.String("user-agent", c.Request.UserAgent()),
				zap.Duration("latency", latency),
				zap.String("time", end.Format(time.RFC3339)),
				zap.String("request_id", request_id),
				zap.String("errors", c.Errors.String()),
			)
		} else {
			// 记录正常的请求日志
			logger.Info("HTTP Request",
				zap.Int("status", c.Writer.Status()),
				zap.String("method", c.Request.Method),
				zap.String("path", path),
				zap.String("query", query),
				zap.String("ip", c.ClientIP()),
				zap.String("user-agent", c.Request.UserAgent()),
				zap.Duration("latency", latency),
				zap.String("time", end.Format(time.RFC3339)),
				zap.String("request_id", request_id),
			)
		}
	}
}
