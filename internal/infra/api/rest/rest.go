package rest

import (
	"metrics/internal/core/service"
	"metrics/internal/infra/api/rest/handlers"

	"github.com/gin-gonic/gin"
)

type API struct {
	srv *gin.Engine
}

func NewAPI(metricService *service.MetricService) *API {
	handler := handlers.NewHandler(metricService)

	srv := gin.Default()
	srv.Use(gin.Logger())
	srv.Use(gin.Recovery())

	srv.GET("/", handler.ListHandler)
	srv.GET("/value/:type/:name", handler.GetHandler)
	srv.POST("/update/:type/:name/:value", handler.UpdateHandler)

	return &API{
		srv: srv,
	}
}

func (app *API) Run() error {
	return app.srv.Run("localhost:8080")
}
