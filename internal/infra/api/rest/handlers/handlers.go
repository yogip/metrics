package handlers

import (
	"fmt"
	"net/http"

	"metrics/internal/core/model"
	"metrics/internal/core/service"
	"metrics/internal/logger"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
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
		logger.Log.Error("Error binding uri", zap.Error(err))
		return
	}
	log := logger.Log.With(
		zap.String("name", req.Name),
		zap.String("type", req.Type.String()),
		zap.String("value", req.Value),
	)
	log.Debug("Getting update request")

	_, err := h.metricService.SetMetricValue(req)
	if err != nil {
		ctx.String(http.StatusBadRequest, err.Error())
		log.Error("Error setting metric value", zap.Error(err))
		return
	}
}

func (h *Handler) GetHandler(ctx *gin.Context) {
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

	metric, err := h.metricService.GetMetric(req)
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

	ctx.String(http.StatusOK, metric.Value)
}

func (h *Handler) ListHandler(ctx *gin.Context) {
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
