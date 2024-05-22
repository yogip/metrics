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

type HandlerV1 struct {
	metricService *service.MetricService
}

func NewHandlerV1(metricService *service.MetricService) *HandlerV1 {
	return &HandlerV1{metricService: metricService}
}

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

	_, err = h.metricService.UpsertMetricValue(reqV2)
	if err != nil {
		ctx.String(http.StatusBadRequest, err.Error())
		log.Error("Metric value update error", zap.Error(err))
		return
	}
}

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

	metric, err := h.metricService.GetMetric(reqV2)
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

func (h *HandlerV1) ListHandler(ctx *gin.Context) {
	metrics, err := h.metricService.ListMetrics()
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
