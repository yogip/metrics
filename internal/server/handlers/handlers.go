package handlers

import (
	"log"
	"net/http"
	"strings"

	"github.com/yogip/metrics/internal/storage"
)

func UpdateHandler(res http.ResponseWriter, req *http.Request) {
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

	metricType := storage.MetricType(pathParts[2])
	metricName := pathParts[3]
	metricValue := pathParts[4]
	log.Printf("Got update input %s:%s set %s\n", metricType, metricName, metricValue)

	if metricType != storage.GaugeType && metricType != storage.CounterType {
		http.Error(res, "Incorrect metric type", http.StatusBadRequest)
		return
	}

	metric, ok := storage.Get(metricType, metricName)
	if !ok {
		metric, _ = storage.NewMetric(metricType, metricName)
	}

	if err := metric.Set(metricValue); err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	if err := storage.SaveMetric(metric); err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}
}

func GetHandler(res http.ResponseWriter, req *http.Request) {
	log.Printf("handler method: [%s] %s\n", req.Method, req.URL.Path)
	if req.Method != http.MethodGet {
		http.Error(res, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
	pathParts := strings.Split(req.URL.Path, "/")
	if len(pathParts) < 4 {
		http.NotFound(res, req)
		return
	}

	metricType := storage.MetricType(pathParts[2])
	metricName := pathParts[3]
	log.Printf("Get value for %s:%s\n", metricType, metricName)

	if metricType != storage.GaugeType && metricType != storage.CounterType {
		http.Error(res, "Bad Request", http.StatusBadRequest)
		return
	}

	metric, ok := storage.Get(metricType, metricName)
	if !ok {
		http.NotFound(res, req)
		return
	}

	res.Write([]byte(metric.String()))
}
