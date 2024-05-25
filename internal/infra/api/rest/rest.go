package rest

import (
	"fmt"
	"metrics/internal/core/service"
	"metrics/internal/infra/api/rest/handlers"
	"metrics/internal/infra/api/rest/middlewares"
	"metrics/internal/logger"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type API struct {
	srv *gin.Engine
}

func ZapLogger(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		duration := time.Since(start)
		msg := fmt.Sprintf(
			"%s %s %d %d %s %s",
			c.Request.Method, c.Request.URL.Path, c.Writer.Status(), c.Writer.Size(), duration, c.Request.UserAgent(),
		)
		logger.Info(msg,
			zap.String("method", c.Request.Method),
			zap.String("path", c.Request.URL.Path),
			zap.Int("status", c.Writer.Status()),
			zap.Duration("duration", duration),
			zap.String("user-agent", c.Request.UserAgent()),
		)
	}
}

func NewAPI(metricService *service.MetricService) *API {
	handlerV1 := handlers.NewHandlerV1(metricService)
	handlerV2 := handlers.NewHandlerV2(metricService)

	srv := gin.Default()
	srv.Use(ZapLogger(logger.Log))
	srv.Use(gin.Recovery())
	srv.Use(middlewares.GzipMiddleware())

	srv.GET("/", handlerV1.ListHandler)
	srv.GET("/value/:type/:name", handlerV1.GetHandler)
	srv.POST("/update/:type/:name/:value", handlerV1.UpdateHandler)

	srv.POST("/value", handlerV2.GetHandler)
	srv.POST("/update", handlerV2.UpdateHandler)

	return &API{
		srv: srv,
	}
}

func (app *API) Run(runAddr string) error {
	return app.srv.Run(runAddr)
}
