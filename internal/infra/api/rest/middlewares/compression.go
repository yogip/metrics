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

func GzipMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if !strings.Contains(
			c.GetHeader("Accept-Encoding"),
			"gzip",
		) {
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
