// V1 version of API handlers.
// The API provides methods for create, update, get and list metics.
// V2 is prefered to usage.
package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"metrics/internal/core/model"
	"metrics/internal/core/service"
	"metrics/internal/logger"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// @Title Metrics REST API V1
// @Version 1.0

type HandlerV1 struct {
	metricService *service.MetricService
}

func NewHandlerV1(metricService *service.MetricService) *HandlerV1 {
	return &HandlerV1{metricService: metricService}
}

// Update metrics API handler
// @Tags V1 API
// @Summary Update metrics
// @Description
// @ID UpdateHandler
// @Param name path string true "Metric name"
// @Param type path string true "Metric type"
// @Param value path string true "Metric value "
// @Success 200
// @Failure 400 {string} string "Bad request"
// @Failure 500 {string} string "Inernal Server Error"
// @Router /update/{type}/{name}/{value}/ [POST]
func (h *HandlerV1) UpdateHandler(ctx *gin.Context) {
	req := &model.MetricUpdateRequest{}
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.String(http.StatusBadRequest, err.Error())
		logger.Log.Error("Error binding uri", zap.Error(err))
		return
	}
	log := logger.Log.With(
		zap.String("name", req.Name),
		zap.String("type", req.Type.String()),
		zap.String("value", req.Value),
	)
	log.Debug("Getting update request")

	reqV2, err := h.metricService.BuildMetricRequest(req.Name, req.Type, req.Value, true)
	if err != nil {
		ctx.String(http.StatusBadRequest, err.Error())
		log.Error("Error parsing metric value", zap.Error(err))
		return
	}

	_, err = h.metricService.UpsertMetricValue(ctx, reqV2)
	if err != nil {
		ctx.String(http.StatusBadRequest, err.Error())
		log.Error("Metric value update error", zap.Error(err))
		return
	}
}

// Get metric API handler
// @Tags V1 API
// @Summary Get metric
// @Description Get metric from storage
// @ID GetHandler
// @Param name path string true "Metric name"
// @Param type path string true "Metric type"
// @Success 200 {string} string "Ok"
// @Failure 400 {string} string "Bad request"
// @Failure 404 {string} string "Not found"
// @Failure 500 {string} string "Inernal Server Error"
// @Router /update/{type}/{name}/ [GET]
func (h *HandlerV1) GetHandler(ctx *gin.Context) {
	req := &model.MetricRequest{}
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.String(http.StatusBadRequest, err.Error())
		logger.Log.Error("Error binding uri", zap.Error(err))
		return
	}
	log := logger.Log.With(
		zap.String("name", req.Name),
		zap.String("type", req.Type.String()),
	)

	log.Debug("Getting value for metric")

	reqV2, err := h.metricService.BuildMetricRequest(req.Name, req.Type, "", false)
	if err != nil {
		ctx.String(http.StatusBadRequest, err.Error())
		log.Error("Error converting metric request to V2", zap.Error(err))
		return
	}

	metric, err := h.metricService.GetMetric(ctx, reqV2)
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

	if metric.MType == model.CounterType {
		ctx.String(http.StatusOK, strconv.FormatInt(*metric.Delta, 10))
	} else {
		ctx.String(http.StatusOK, strconv.FormatFloat(*metric.Value, 'f', -1, 64))
	}
}

// List all metrics API handler
// @Tags V1 API
// @Summary List metrics
// @Description Get metric all from storage
// @ID ListHandler
// @Success 200 {string} string "Metrics list"
// @Failure 400 {string} string "Bad request"
// @Failure 500 {string} string "Inernal Server Error"
// @Router / [GET]
func (h *HandlerV1) ListHandler(ctx *gin.Context) {
	metrics, err := h.metricService.ListMetrics(ctx)
	if err != nil {
		ctx.String(http.StatusBadRequest, err.Error())
		return
	}

	content := "<html><body><ul>%s</ul></body></html>"
	var listItems string
	for _, m := range metrics.Metrics {
		listItems += "<li><strong>" + m.Name + "</strong>: " + m.Value + "</li>"
	}

	ctx.Header("Content-Type", "text/html")
	ctx.String(http.StatusOK, fmt.Sprintf(content, listItems))
}
