package handlers

import (
	"context"
	"fmt"
	"net/http"
	"time"

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

	tOutCtx, cancel := context.WithTimeout(ctx, time.Second*15)
	defer cancel()

	metric, err := h.metricService.UpsertMetricValue(tOutCtx, req)
	if err != nil {
		ctx.String(http.StatusBadRequest, err.Error())
		log.Error("Error setting metric value", zap.Error(err))
		return
	}
	ctx.JSON(http.StatusOK, metric)
}

func (h *HandlerV2) BatchUpdateHandler(ctx *gin.Context) {
	req := []*model.MetricsV2{}
	if err := ctx.ShouldBindBodyWithJSON(&req); err != nil {
		ctx.String(http.StatusBadRequest, err.Error())
		logger.Log.Error("Error binding body", zap.Error(err))
		return
	}
	logger.Log.Debug("Getting update request")

	metrics, err := h.metricService.BatchUpsertMetricValue(ctx, req)
	if err != nil {
		ctx.String(http.StatusBadRequest, err.Error())
		logger.Log.Error("Batch update error", zap.Error(err))
		return
	}
	ctx.JSON(http.StatusOK, metrics)
}

func (h *HandlerV2) GetHandler(ctx *gin.Context) {
	req := &model.MetricsV2{}
	if err := ctx.ShouldBindBodyWithJSON(&req); err != nil {
		logger.Log.Error("Error binding uri", zap.Error(err))
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"status": false, "message": fmt.Sprintf("Error binding body: %s", err)})
		return
	}
	log := logger.Log.With(
		zap.String("name", req.ID),
		zap.String("type", req.MType.String()),
	)

	log.Debug("Getting value for metric")

	metric, err := h.metricService.GetMetric(ctx, req)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"status": false, "message": fmt.Sprintf("Error getting metric: %s", err)})
		logger.Log.Error("Error getting metric", zap.Error(err))
		return
	}
	if metric == nil {
		ctx.String(http.StatusNotFound, "Not found")
		return
	}

	ctx.JSON(http.StatusOK, metric)
}
