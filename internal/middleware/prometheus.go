package middleware

import (
	"strconv"
	"time"

	prometheus2 "github.com/Citadelas/api-gateway/internal/app/prometheus"
	"github.com/gin-gonic/gin"
)

func PrometheusMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.FullPath()
		if path == "" {
			path = c.Request.URL.Path
		}
		method := c.Request.Method
		c.Next()
		duration := time.Since(start).Seconds()
		status := strconv.Itoa(c.Writer.Status())

		prometheus2.RequestsTotal.WithLabelValues(method, path, status).Inc()
		prometheus2.RequestDuration.WithLabelValues(method, path).Observe(duration)
	}
}
