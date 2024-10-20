package middlewares

import (
	"fmt"
	"net"
	"net/http"

	"metrics/internal/logger"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func CheckSubnet(subnet *net.IPNet) gin.HandlerFunc {
	return func(c *gin.Context) {
		addr := c.ClientIP()
		log := logger.Log.With(
			zap.String("method", c.Request.Method),
			zap.String("url", c.Request.URL.String()),
			zap.String("RemouteAddr", addr),
		)

		clientIP := net.ParseIP(addr)
		if clientIP == nil || !subnet.Contains(clientIP) {
			log.Warn(fmt.Sprintf("Client IP %s is not allowed", clientIP))
			c.AbortWithStatus(http.StatusForbidden)
			return
		}

		c.Next()
	}
}
