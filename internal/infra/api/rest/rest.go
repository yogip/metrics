package rest

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"metrics/internal/core/config"
	"metrics/internal/core/service"
	"metrics/internal/infra/api/rest/handlers"
	"metrics/internal/infra/api/rest/middlewares"
	"metrics/internal/logger"

	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type API struct {
	srv *http.Server
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

func NewAPI(cfg *config.Config, metricService *service.MetricService, systemService *service.SystemService) *API {
	serviceHandler := handlers.NewSystemHandler(systemService)
	handlerV1 := handlers.NewHandlerV1(metricService)
	handlerV2 := handlers.NewHandlerV2(metricService)

	router := gin.Default()
	router.Use(ZapLogger(logger.Log))
	router.Use(gin.Recovery())
	router.Use(middlewares.GzipDecompressMiddleware())
	router.Use(middlewares.GzipCompressMiddleware())

	if cfg.HashKey != "" {
		router.Use(middlewares.VerifySignature(cfg.HashKey))
		router.Use(middlewares.SignBody(cfg.HashKey))
	}

	router.GET("/ping", serviceHandler.Ping)

	router.GET("/", handlerV1.ListHandler)
	router.GET("/value/:type/:name/", handlerV1.GetHandler)
	router.POST("/update/:type/:name/:value/", handlerV1.UpdateHandler)

	router.POST("/value/", handlerV2.GetHandler)
	router.POST("/update/", handlerV2.UpdateHandler)
	router.POST("/updates/", handlerV2.BatchUpdateHandler)

	pprof.Register(router)
	srv := &http.Server{Handler: router}
	return &API{
		srv: srv,
	}
}

func (api *API) Run(runAddr string) error {
	logger.Log.Info("Run API server", zap.String("Addres", runAddr))
	api.srv.Addr = runAddr
	return api.srv.ListenAndServe()
}

func (api *API) Shutdown(ctx context.Context) error {
	return api.srv.Shutdown(ctx)
}
