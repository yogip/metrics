package transport

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"metrics/internal/core/model"
	"metrics/internal/logger"

	"go.uber.org/zap"
)

type HTTPClient struct {
	ServerHost     string
	MetricEndpoint string
}

func NewClient(serverHost string) *HTTPClient {
	return &HTTPClient{
		ServerHost:     serverHost,
		MetricEndpoint: "/update",
	}
}

// HTTTP Client to sent metrics to MetricEndpoint
func (client *HTTPClient) SendMetric(req *model.MetricsV2) error {
	body, err := json.Marshal(req)
	if err != nil {
		return fmt.Errorf("error marshalling request body: %w", err)
	}

	url := fmt.Sprintf(client.ServerHost + client.MetricEndpoint)
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(body))
	if err != nil {
		return fmt.Errorf("request error: %w", err)
	}
	defer resp.Body.Close()
	logger.Log.Debug("Metric was send", zap.String("url", url), zap.Int("status", resp.StatusCode))

	if resp.StatusCode != http.StatusOK {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("reading body error: %w, code: %d", err, resp.StatusCode)
		}

		return fmt.Errorf("request error: %s, code: %d", string(body), resp.StatusCode)
	}
	return nil
}
