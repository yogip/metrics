package middlewares

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestVerifySignature(t *testing.T) {
	tests := []struct {
		name      string
		signature string
		code      int
	}{
		{
			name:      "Test signature 1",
			signature: "88a0f14742ff0dfd031600abc677726a7fe9f81271ad6b1fed0bfdec30c3f284",
			code:      http.StatusOK,
		},
		{
			name:      "Test signature 2",
			signature: "88a0f14742ff0dfd031600abc677726a7fe9f81271ad6b1fed0bfdec30c3f284",
			code:      http.StatusOK,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gin.SetMode(gin.TestMode)
			router := setupApp()
			req, _ := http.NewRequest(http.MethodPost, "/", bytes.NewBuffer([]byte(tt.name)))
			req.Header.Set("HashSHA256", tt.signature)
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)
			result := w.Result()
			assert.Equal(t, tt.code, result.StatusCode)
			result.Body.Close()
		})
	}
}

func setupApp() *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.Use(VerifySignature("test_hash_key"))
	router.POST("/", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})
	return router
}
