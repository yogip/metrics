package transport

import (
	"fmt"
	"io"
	"log"
	"net/http"

	"metrics/internal/core/model"
)

type HTTPClient struct {
	ServerHost     string
	MetricEndpoint string
}

func NewClient(serverHost string) *HTTPClient {
	return &HTTPClient{
		ServerHost:     serverHost,
		MetricEndpoint: "/update/%s/%s/%s", // POST http://<АДРЕС_СЕРВЕРА>/update/<ТИП_МЕТРИКИ>/<ИМЯ_МЕТРИКИ>/<ЗНАЧЕНИЕ_МЕТРИКИ>
	}
}

// HTTTP Client to sent metrics to MetricEndpoint
func (client *HTTPClient) SendMetric(req *model.MetricResponse) error {
	url := "http://" + fmt.Sprintf(client.ServerHost+client.MetricEndpoint, req.Type, req.Name, req.Value)
	resp, err := http.Post(url, "text/plain", nil)
	if err != nil {
		log.Println(err)
		return fmt.Errorf("request error: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		log.Printf("Error sending metric: %s", resp.Status)
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Println("Error reading body", err)
			return fmt.Errorf("reading body error: %w, code: %d", err, resp.StatusCode)
		}
		log.Printf("Response body: %s", string(body))
		return fmt.Errorf("request error: %s, code: %d", string(body), resp.StatusCode)
	}
	return nil
}
