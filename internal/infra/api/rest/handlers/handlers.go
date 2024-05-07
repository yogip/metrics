package handlers

import (
	"log"
	"net/http"

	"metrics/internal/core/model"
	"metrics/internal/core/service"

	"github.com/gin-gonic/gin"
	// "metrics/internal/infra/"
)

type Handler struct {
	metricService *service.MetricService
}

func NewHandler(metricService *service.MetricService) *Handler {
	return &Handler{metricService: metricService}
}

func (h *Handler) UpdateHandler(ctx *gin.Context) {
	metricType := model.MetricType(ctx.Param("type"))
	metricName := ctx.Param("name")
	metricValue := ctx.Param("value")
	log.Printf("Got update input %s:%s set %s\n", metricType, metricName, metricValue)

	if metricType != model.GaugeType && metricType != model.CounterType {
		ctx.String(http.StatusBadRequest, "Incorrect metric type: %s", metricType)
		return
	}

	_, err := h.metricService.SetMetricValue(
		&model.MetricUpdateRequest{Name: metricName, Type: metricType, Value: metricValue},
	)
	if err != nil {
		ctx.String(http.StatusBadRequest, err.Error())
		return
	}
}

func (h *Handler) GetHandler(ctx *gin.Context) {
	metricType := model.MetricType(ctx.Param("type"))
	metricName := ctx.Param("name")
	log.Printf("Getting value for %s:%s\n", metricType, metricName)

	if metricType != model.GaugeType && metricType != model.CounterType {
		ctx.String(http.StatusBadRequest, "Incorrect metric type: %s", metricType)
		return
	}

	metric, err := h.metricService.GetMetric(
		&model.MetricRequest{Name: metricName, Type: metricType},
	)
	if err != nil {
		ctx.String(http.StatusBadRequest, err.Error())
		return
	}
	if metric == nil {
		ctx.String(http.StatusNotFound, "Not found")
		return
	}

	ctx.String(http.StatusOK, metric.Value)
}
