package transport

import (
	"bytes"
	"compress/gzip"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"metrics/internal/core/model"
	"metrics/internal/logger"

	"go.uber.org/zap"
)

type HTTPClient struct {
	serverHost     string
	metricEndpoint string
	signHashKey    string
	client         *http.Client
}

func NewClient(serverHost string, signHashKey string) *HTTPClient {
	return &HTTPClient{
		serverHost:     serverHost,
		signHashKey:    signHashKey,
		metricEndpoint: "/updates",
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

func (c *HTTPClient) sign(data *[]byte) string {
	h := hmac.New(sha256.New, []byte(c.signHashKey))
	h.Write(*data)
	return hex.EncodeToString(h.Sum(nil))
}

// HTTTP Client to sent metrics to MetricEndpoint
func (c *HTTPClient) SendMetric(data []model.MetricsV2) error {
	body, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("error marshalling request body: %w", err)
	}

	buf, err := c.compress(&body)
	if err != nil {
		return fmt.Errorf("error compressiong request body: %w", err)
	}

	url := fmt.Sprintf(c.serverHost + c.metricEndpoint)
	req, err := http.NewRequest("POST", url, buf)
	if err != nil {
		return fmt.Errorf("request creation error: %w", err)
	}
	req.Header.Set("Content-Encoding", "gzip")
	req.Header.Set("Content-Type", "application/json")

	if c.signHashKey != "" {
		signature := c.sign(&body)
		logger.Log.Debug(fmt.Sprintf("signature for body - %s", string(signature)))
		req.Header.Set("HashSHA256", string(signature))
	}

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
