package handlers

import (
	"net/http"

	"metrics/internal/core/model"
	"metrics/internal/core/service"
	"metrics/internal/logger"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type HandlerV2 struct {
	metricService *service.MetricService
}

func NewHandlerV2(metricService *service.MetricService) *HandlerV2 {
	return &HandlerV2{metricService: metricService}
}

func (h *HandlerV2) UpdateHandler(ctx *gin.Context) {
	req := &model.MetricsV2{}
	if err := ctx.ShouldBindBodyWithJSON(&req); err != nil {
		ctx.String(http.StatusBadRequest, err.Error())
		logger.Log.Error("Error binding uri", zap.Error(err))
		return
	}
	log := logger.Log.With(
		zap.String("name", req.ID),
		zap.String("type", req.MType.String()),
		zap.Int64p("delta", req.Delta),
		zap.Float64p("value", req.Value),
	)
	log.Debug("Getting update request")

	metric, err := h.metricService.UpsertMetricValue(ctx, req)
	if err != nil {
		ctx.String(http.StatusBadRequest, err.Error())
		log.Error("Error setting metric value", zap.Error(err))
		return
	}
	ctx.JSON(http.StatusOK, metric)
}

func (h *HandlerV2) GetHandler(ctx *gin.Context) {
	req := &model.MetricsV2{}
	if err := ctx.ShouldBindBodyWithJSON(&req); err != nil {
		ctx.String(http.StatusBadRequest, err.Error())
		logger.Log.Error("Error binding uri", zap.Error(err))
		return
	}
	log := logger.Log.With(
		zap.String("name", req.ID),
		zap.String("type", req.MType.String()),
	)

	log.Debug("Getting value for metric")

	metric, err := h.metricService.GetMetric(ctx, req)
	if err != nil {
		ctx.String(http.StatusBadRequest, err.Error())
		logger.Log.Error("Error getting metric", zap.Error(err))
		return
	}
	if metric == nil {
		ctx.String(http.StatusNotFound, "Not found")
		logger.Log.Error("Metric not found", zap.Error(err))
		return
	}

	ctx.JSON(http.StatusOK, metric)
}
