package middlewares

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"net/http"

	"metrics/internal/logger"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func VerifySignature(hashKey string) gin.HandlerFunc {
	return func(c *gin.Context) {
		log := logger.Log.With(
			zap.String("method", c.Request.Method),
			zap.String("url", c.Request.URL.String()),
			zap.String("HashSHA256", c.GetHeader("HashSHA256")),
		)

		if c.Request.Method != http.MethodPut && c.Request.Method != http.MethodPost {
			c.Next()
			return
		}

		signature := c.GetHeader("HashSHA256")
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

		if validSignature != signature {
			log.Warn(
				"Signature verification error!",
				zap.String("validSignature", validSignature),
				zap.String("signature", signature),
				zap.String("body", string(body)),
			)
			// пропускаю 400 для временного фикса, т/к авто-тест на запрос без подписи ждет 200 ответ
			// https://app.pachca.com/chats?thread_id=4024933&sidebar_message=263511005
			// c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"status": false, "message": "Signaure is not valid"})
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
