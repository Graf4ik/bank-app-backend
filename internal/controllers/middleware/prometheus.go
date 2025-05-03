package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
)

func PrometheusMiddleware(requestCount *prometheus.CounterVec) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		statusCode := c.Writer.Status()
		method := c.Request.Method
		requestCount.WithLabelValues(method, fmt.Sprintf("%d", statusCode)).Inc()
	}
}
