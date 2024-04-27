package server

import (
	"log"
	"net/http"
	"strings"

	"github.com/yogip/metrics/internal/shared"
)

func mainPage(res http.ResponseWriter, req *http.Request) {
	log.Printf("handler method: [%s] %s\n", req.Method, req.URL.Path)
	if req.Method != http.MethodPost {
		http.Error(res, "Method Not Allowed", http.StatusMethodNotAllowed)
	}

	// http://<АДРЕС_СЕРВЕРА>/update/<ТИП_МЕТРИКИ>/<ИМЯ_МЕТРИКИ>/<ЗНАЧЕНИЕ_МЕТРИКИ>
	pathParts := strings.Split(req.URL.Path, "/")
	if len(pathParts) < 5 || pathParts[1] != "update" {
		http.NotFound(res, req)
		return
	}

	metricType := shared.MetricType(pathParts[2])
	metricName := pathParts[3]
	metricValue := pathParts[4]

	if metricType != shared.GaugeType && metricType != shared.CounterType {
		http.Error(res, "Bad Request", http.StatusBadRequest)
		return
	}

	log.Printf("handler method: %s\n", req.Method)
	log.Printf("metric type: %s\n", metricType)
	log.Printf("metric name: %s\n", metricName)
	log.Printf("metric value: %s\n", metricValue)

	// body := fmt.Sprintf("Method: %s\r\n", req.Method)
	// body += "Header ===============\r\n"
	// for k, v := range req.Header {
	// 	body += fmt.Sprintf("%s: %v\r\n", k, v)
	// }

	// if err := req.ParseForm(); err != nil {
	// 	body += fmt.Sprintf("ParseForm error: %v\r\n", err)
	// 	res.Write([]byte(body))
	// 	return
	// }

	// body += "Query parameters ===============\r\n"
	// for k, v := range req.Form {
	// 	body += fmt.Sprintf("%s: %v\r\n", k, v)
	// }
	// res.Write([]byte(body))
}

func Run() {
	log.Printf("Start server\n")
	mux := http.NewServeMux()
	mux.HandleFunc("/update/", mainPage)

	err := http.ListenAndServe(`localhost:8080`, mux)
	if err != nil {
		panic(err)
	}
}
