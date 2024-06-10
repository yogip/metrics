package transport

import (
	"bytes"
	"compress/gzip"
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
	client         *http.Client
}

func NewClient(serverHost string) *HTTPClient {
	return &HTTPClient{
		ServerHost:     serverHost,
		MetricEndpoint: "/updates",
		client:         &http.Client{},
	}
}

func (c *HTTPClient) compress(data *[]byte) (*bytes.Buffer, error) {
	var buf bytes.Buffer
	gz, err := gzip.NewWriterLevel(&buf, gzip.BestCompression)
	if err != nil {
		return nil, fmt.Errorf("error creating gzip writer: %w", err)
	}
	if _, err := gz.Write(*data); err != nil {
		return nil, fmt.Errorf("error writing to gzip writer: %w", err)
	}
	if err := gz.Close(); err != nil {
		return nil, fmt.Errorf("error closing gzip writer: %w", err)
	}
	return &buf, nil
}

// HTTTP Client to sent metrics to MetricEndpoint
func (c *HTTPClient) SendMetric(data []*model.MetricsV2) error {
	body, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("error marshalling request body: %w", err)
	}

	buf, err := c.compress(&body)
	if err != nil {
		return fmt.Errorf("error compressiong request body: %w", err)
	}

	url := fmt.Sprintf(c.ServerHost + c.MetricEndpoint)
	req, err := http.NewRequest("POST", url, buf)
	if err != nil {
		return fmt.Errorf("request creation error: %w", err)
	}
	req.Header.Set("Content-Encoding", "gzip")
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(req)
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
