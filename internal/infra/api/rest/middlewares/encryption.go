package middlewares

import (
	"bytes"
	"crypto/rsa"
	"fmt"
	"io"
	"net/http"

	"metrics/internal/core/service"

	"github.com/gin-gonic/gin"
)

func DecryptReqBody(privateKey *rsa.PrivateKey) gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Method != http.MethodPut && c.Request.Method != http.MethodPost {
			c.Next()
			return
		}

		encBody, err := io.ReadAll(c.Request.Body)
		if err != nil {
			c.AbortWithStatusJSON(
				http.StatusBadRequest,
				gin.H{"status": false, "message": fmt.Sprintf("read body error: %s", err)},
			)
			return
		}

		if len(encBody) == 0 {
			c.Next()
			return
		}
		body, err := service.Decrypt(privateKey, encBody)
		if err != nil {
			c.AbortWithStatusJSON(
				http.StatusBadRequest,
				gin.H{"status": false, "message": fmt.Sprintf("decrypt body error: %s", err)},
			)
		}

		// Восстановление тела запроса для последующего использования
		c.Request.Body = io.NopCloser(bytes.NewBuffer(body))
		c.Next()
	}
}
