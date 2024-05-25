package middlewares

import (
	"compress/gzip"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type gzipResponseWriter struct {
	gin.ResponseWriter
	Writer *gzip.Writer
}

func (g *gzipResponseWriter) Write(data []byte) (int, error) {
	return g.Writer.Write(data)
}

func GzipCompressMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if !strings.Contains(
			c.GetHeader("Accept-Encoding"),
			"gzip",
		) {
			c.Next()
			return
		}
		isValidContentType := false
		validContentTypes := []string{"application/json", "text/html", "html/text"}
		for _, validContentType := range validContentTypes {
			if c.GetHeader("Content-Type") == validContentType {
				isValidContentType = true
				break
			}
		}

		if !isValidContentType {
			c.Next()
			return
		}

		gz, err := gzip.NewWriterLevel(c.Writer, gzip.BestSpeed)
		if err != nil {
			c.String(
				http.StatusInternalServerError,
				fmt.Errorf("error creating gzip writer: %w", err).Error(),
			)
			return
		}
		defer gz.Close()

		c.Header("Content-Encoding", "gzip")
		c.Writer = &gzipResponseWriter{Writer: gz, ResponseWriter: c.Writer}
		c.Next()
	}
}

func GzipDecompressMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if !strings.Contains(
			c.GetHeader("Content-Encoding"),
			"gzip",
		) {
			c.Next()
			return
		}

		gz, err := gzip.NewReader(c.Request.Body)
		if err != nil {
			c.String(
				http.StatusInternalServerError,
				fmt.Errorf("error creating gzip reader: %w", err).Error(),
			)
			return
		}
		defer gz.Close()

		c.Request.Body = gz
		c.Next()
	}
}
