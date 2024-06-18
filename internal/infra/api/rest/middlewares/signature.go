package middlewares

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"metrics/internal/logger"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func VerifySignature(hashKey string) gin.HandlerFunc {
	return func(c *gin.Context) {
		log := logger.Log.With(
			zap.String("hashKey", hashKey),
			zap.String("method", c.Request.Method),
			zap.String("method", c.Request.URL.String()),
			zap.String("HashSHA256", c.GetHeader("HashSHA256")),
			zap.String("Hash", c.GetHeader("Hash")),
		)
		log.Info("--!! Start VerifySignature")

		if c.Request.Method != http.MethodPut && c.Request.Method != http.MethodPost {
			c.Next()
			return
		}

		signature := c.GetHeader("HashSHA256")
		// if signature = c.GetHeader("Hashsha256"); signature == "" {
		// 	c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"status": false, "message": "There is no signature header"})
		// 	return
		// }

		body, err := io.ReadAll(c.Request.Body)
		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		// Восстановление тела запроса для последующего использования
		c.Request.Body = io.NopCloser(bytes.NewBuffer(body))

		if len(body) == 0 {
			log.Debug("Skip Signature verification beacuse of empty body")
			c.Next()
			return
		}

		h := hmac.New(sha256.New, []byte(hashKey))
		h.Write(body)
		validSignature := hex.EncodeToString(h.Sum(nil))

		log.Info(fmt.Sprintf("%s == %s", validSignature, signature))
		log.Info("==VerifySignature =====")

		if validSignature != signature {
			log.Warn("Signature error", zap.String("validSignature", validSignature), zap.String("signature", signature))
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}

		c.Next()
	}
}

type responseWriter struct {
	gin.ResponseWriter
	hashKey string
}

func (rw *responseWriter) Write(b []byte) (int, error) {
	h := hmac.New(sha256.New, []byte(rw.hashKey))
	h.Write(b)
	signature := hex.EncodeToString(h.Sum(nil))

	logger.Log.Info(fmt.Sprintf("Set HashSHA256 == %s", signature))

	rw.Header().Set("HashSHA256", signature)
	return rw.ResponseWriter.Write(b)
}

func SignBody(hashKey string) gin.HandlerFunc {
	return func(c *gin.Context) {
		rw := &responseWriter{
			ResponseWriter: c.Writer,
			hashKey:        hashKey,
		}
		c.Writer = rw

		c.Next()
	}
}
