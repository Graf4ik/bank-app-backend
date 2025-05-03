package middleware

import (
	"bank-app-backend/internal/lib/logger"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"time"
)

func ZapLoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		c.Next()

		lib.Log.Info("request",
			zap.String("method", c.Request.Method),
			zap.String("path", path),
			zap.Int("status", c.Writer.Status()),
			zap.Duration("duration", time.Since(start)),
		)
	}
}
