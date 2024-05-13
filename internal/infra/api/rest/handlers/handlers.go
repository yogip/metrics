package handlers

import (
	"fmt"
	"log"
	"net/http"

	"metrics/internal/core/model"
	"metrics/internal/core/service"

	"github.com/gin-gonic/gin"
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
