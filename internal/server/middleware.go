package server

import (
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func LoggerMiddleware(log *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		now := time.Now()
		c.Next()
		log.Sugar().Infof("| %d | %s | %s | %s |",
			c.Writer.Status(),
			time.Since(now),
			c.Request.Method,
			c.Request.URL.String(),
		)
	}
}
