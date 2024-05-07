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
	req := &model.MetricUpdateRequest{}
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.String(http.StatusBadRequest, err.Error())
	}
	log.Printf("Getting update request %s", req)

	_, err := h.metricService.SetMetricValue(req)
	if err != nil {
		ctx.String(http.StatusBadRequest, err.Error())
		return
	}
}

func (h *Handler) GetHandler(ctx *gin.Context) {
	req := &model.MetricRequest{}
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.String(http.StatusBadRequest, err.Error())
	}

	log.Printf("Getting value for %s:%s", req.Name, req.Type)

	metric, err := h.metricService.GetMetric(req)
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
