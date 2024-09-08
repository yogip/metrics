package middlewares

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCompression(t *testing.T) {
	body := []byte("Test compression")
	compressed, err := compress(&body)
	require.NoError(t, err)

	gin.SetMode(gin.TestMode)
	router := setupRouter()
	req, err := http.NewRequest(http.MethodPost, "/", compressed)
	require.NoError(t, err)
	req.Header.Set("Content-Encoding", "gzip")
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Result().StatusCode)
	assert.Equal(t, "Test compression", w.Body.String())
	defer w.Result().Body.Close()

}

func compress(data *[]byte) (*bytes.Buffer, error) {
	var buf bytes.Buffer

	gz, err := gzip.NewWriterLevel(&buf, gzip.BestCompression)
	if err != nil {
		return nil, fmt.Errorf("error creating gzip writer: %w", err)
	}
	if _, err := gz.Write(*data); err != nil {
		return nil, fmt.Errorf("error writing to gzip writer: %w", err)
	}
	if err := gz.Close(); err != nil {
		return nil, fmt.Errorf("error closing gzip writer: %w", err)
	}
	return &buf, nil
}

func setupRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.Use(GzipDecompressMiddleware())
	router.Use(GzipCompressMiddleware())
	router.POST("/", func(c *gin.Context) {
		body, err := io.ReadAll(c.Request.Body)
		if err != nil || body == nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		c.String(http.StatusOK, string(body))
	})
	return router
}
