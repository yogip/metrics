package handlers

import (
	"log"
	"net/http"
	"strings"

	"metrics/internal/core/model"
	"metrics/internal/core/service"
	// "metrics/internal/infra/"
)

type Handler struct {
	metricService *service.MetricService
}

func NewHandler(metricService *service.MetricService) *Handler {
	return &Handler{metricService: metricService}
}

func (h *Handler) UpdateHandler(res http.ResponseWriter, req *http.Request) {
	log.Printf("handler method: [%s] %s\n", req.Method, req.URL.Path)
	if req.Method != http.MethodPost {
		http.Error(res, "Method Not Allowed", http.StatusMethodNotAllowed)
	}

	// http://<АДРЕС_СЕРВЕРА>/update/<ТИП_МЕТРИКИ>/<ИМЯ_МЕТРИКИ>/<ЗНАЧЕНИЕ_МЕТРИКИ>
	pathParts := strings.Split(req.URL.Path, "/")
	if len(pathParts) < 5 {
		http.NotFound(res, req)
		return
	}

	metricType := model.MetricType(pathParts[2])
	metricName := pathParts[3]
	metricValue := pathParts[4]
	log.Printf("Got update input %s:%s set %s\n", metricType, metricName, metricValue)

	if metricType != model.GaugeType && metricType != model.CounterType {
		http.Error(res, "Incorrect metric type", http.StatusBadRequest)
		return
	}

	_, err := h.metricService.SetMetricValue(
		&model.MetricUpdateRequest{Name: metricName, Type: metricType, Value: metricValue},
	)
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}
}

func (h *Handler) GetHandler(res http.ResponseWriter, req *http.Request) {
	log.Printf("handler method: [%s] %s\n", req.Method, req.URL.Path)
	if req.Method != http.MethodGet {
		http.Error(res, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
	pathParts := strings.Split(req.URL.Path, "/")
	if len(pathParts) < 4 {
		http.NotFound(res, req)
		return
	}

	metricType := model.MetricType(pathParts[2])
	metricName := pathParts[3]
	log.Printf("Get value for %s:%s\n", metricType, metricName)

	if metricType != model.GaugeType && metricType != model.CounterType {
		http.Error(res, "Bad Request", http.StatusBadRequest)
		return
	}

	metric, err := h.metricService.GetMetric(
		&model.MetricRequest{Name: metricName, Type: metricType},
	)
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}
	if metric == nil {
		http.NotFound(res, req)
		return
	}

	res.Write([]byte(metric.Value))
}
