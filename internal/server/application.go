package server

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/yogip/metrics/internal/repo"
)

var storage *repo.MemStorage

func updateHandler(res http.ResponseWriter, req *http.Request) {
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

	metricType := repo.MetricType(pathParts[2])
	metricName := pathParts[3]
	metricValue := pathParts[4]
	log.Printf("Got update input %s:%s set %s\n", metricType, metricName, metricValue)

	if metricType != repo.GaugeType && metricType != repo.CounterType {
		http.Error(res, "Incorrect metric type", http.StatusBadRequest)
		return
	}

	if err := storage.SetValue(metricType, metricName, metricValue); err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}
}

func getHandler(res http.ResponseWriter, req *http.Request) {
	log.Printf("handler method: [%s] %s\n", req.Method, req.URL.Path)
	if req.Method != http.MethodGet {
		http.Error(res, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
	pathParts := strings.Split(req.URL.Path, "/")
	if len(pathParts) < 4 {
		http.NotFound(res, req)
		return
	}

	metricType := repo.MetricType(pathParts[2])
	metricName := pathParts[3]
	log.Printf("Get value for %s:%s\n", metricType, metricName)

	if metricType != repo.GaugeType && metricType != repo.CounterType {
		http.Error(res, "Bad Request", http.StatusBadRequest)
		return
	}

	metric, ok := storage.Get(metricType, metricName)
	if !ok {
		log.Printf("storage: %v \n", storage)
		http.NotFound(res, req)
		return
	}

	text := `
	<html>
		<head>
			<title>Metric value</title>
		</head>
		<body>
			<h1>Metric %s</h1>
			<p><b>Value:</b> %v</p>
	</html>`

	if metricType == repo.GaugeType {
		gaugeMetric := metric.(*repo.Gauge)
		text = fmt.Sprintf(text, gaugeMetric.Name, gaugeMetric.Value)
	} else {
		counterMetric := metric.(*repo.Counter)
		text = fmt.Sprintf(text, counterMetric.Name, counterMetric.Value)
	}

	res.Write([]byte(text))
}

func Run() {
	log.Println("Start server")
	mux := http.NewServeMux()
	mux.HandleFunc("/update/", updateHandler)
	mux.HandleFunc("/value/", getHandler)

	log.Println("Init storage")
	storage = repo.NewMemStorage()

	err := http.ListenAndServe(`localhost:8080`, mux)
	if err != nil {
		panic(err)
	}
}
