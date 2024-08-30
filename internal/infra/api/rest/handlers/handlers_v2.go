// V2 version of API handlers.
// The API provides methods for create, update, batch update, get and list metics.
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

// @Title Metrics REST API V2
// @Version 2.0

type HandlerV2 struct {
	metricService *service.MetricService
}

func NewHandlerV2(metricService *service.MetricService) *HandlerV2 {
	return &HandlerV2{metricService: metricService}
}

// Update metrics API handler
// @Tags V2 API
// @Summary Update metrics
// @Description
// @ID UpdateHandlerV2
// @Accept  json
// @Produce json
// @Param req body model.MetricsV2 true "Metric Name"
// @Success 200 {object} model.MetricsV2
// @Failure 400 {string} string "Bad request"
// @Failure 500 {string} string "Inernal Server Error"
// @Router /update/ [POST]
func (h *HandlerV2) UpdateHandler(ctx *gin.Context) {
	req := &model.MetricsV2{}
	if err := ctx.ShouldBindBodyWithJSON(&req); err != nil {
		logger.Log.Error("Error binding uri", zap.Error(err))
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"status": false, "message": fmt.Sprintf("Error binding body: %s", err)})
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
		ctx.AbortWithStatusJSON(
			http.StatusBadRequest,
			gin.H{"status": false, "message": fmt.Sprintf("Error setting metric value: %s", err)},
		)
		return
	}
	ctx.JSON(http.StatusOK, metric)
}

// Batch metrics update API handler
// @Tags V2 API
// @Summary Batch update
// @Description
// @ID BatchUpdateHandler
// @Accept  json
// @Produce json
// @Param req body []model.MetricsV2 true "Metrics request"
// @Success 200 {object} []model.MetricsV2
// @Failure 400 {string} string "Bad request"
// @Failure 500 {string} string "Inernal Server Error"
// @Router /updates/ [POST]
func (h *HandlerV2) BatchUpdateHandler(ctx *gin.Context) {
	req := []*model.MetricsV2{}
	if err := ctx.ShouldBindBodyWithJSON(&req); err != nil {
		logger.Log.Error("Error binding body", zap.Error(err))
		ctx.AbortWithStatusJSON(
			http.StatusBadRequest,
			gin.H{"status": false, "message": fmt.Sprintf("Error binding body: %s", err)},
		)
		return
	}
	logger.Log.Debug("Getting update request")

	metrics, err := h.metricService.BatchUpsertMetricValue(ctx, req)
	if err != nil {
		logger.Log.Error("Batch update error", zap.Error(err))
		ctx.AbortWithStatusJSON(
			http.StatusBadRequest,
			gin.H{"status": false, "message": fmt.Sprintf("Batch upsert error: %s", err)},
		)
		return
	}
	ctx.JSON(http.StatusOK, metrics)
}

// Get metrics API handler
// @Tags V2 API
// @Summary Get metrics
// @Description
// @ID GetHandlerV2
// @Accept  json
// @Produce json
// @Param req body model.MetricsV2 true "Metric request"
// @Success 200 {object} model.MetricsV2
// @Failure 400 {string} string "Bad request"
// @Failure 500 {string} string "Inernal Server Error"
// @Router /value/ [POST]
func (h *HandlerV2) GetHandler(ctx *gin.Context) {
	req := &model.MetricsV2{}
	if err := ctx.ShouldBindBodyWithJSON(&req); err != nil {
		logger.Log.Error("Error binding uri", zap.Error(err))
		ctx.AbortWithStatusJSON(
			http.StatusNotFound,
			gin.H{"status": false, "message": fmt.Sprintf("Error binding body: %s", err)},
		)
		return
	}
	log := logger.Log.With(
		zap.String("name", req.ID),
		zap.String("type", req.MType.String()),
	)

	log.Debug("Getting value for metric")

	metric, err := h.metricService.GetMetric(ctx, req)
	if err != nil {
		logger.Log.Error("Error getting metric", zap.Error(err))
		ctx.AbortWithStatusJSON(
			http.StatusBadRequest,
			gin.H{"status": false, "message": fmt.Sprintf("Error getting metric: %s", err)},
		)
		return
	}
	if metric == nil {
		ctx.AbortWithStatusJSON(
			http.StatusNotFound,
			gin.H{"status": false, "message": "Not found"},
		)
		return
	}

	ctx.JSON(http.StatusOK, metric)
}
