package transport

import (
	"bytes"
	"compress/gzip"
	"crypto/hmac"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"metrics/internal/core/model"
	"metrics/internal/core/service"
	"metrics/internal/logger"

	"go.uber.org/zap"
)

type HTTPClient struct {
	client         *http.Client
	serverHost     string
	metricEndpoint string
	signHashKey    string
	pubKey         *rsa.PublicKey
}

func NewClient(serverHost string, signHashKey string, pubKey *rsa.PublicKey) *HTTPClient {
	return &HTTPClient{
		serverHost:     serverHost,
		signHashKey:    signHashKey,
		pubKey:         pubKey,
		metricEndpoint: "/updates",
		client:         &http.Client{},
	}
}

func (c *HTTPClient) compress(data []byte) ([]byte, error) {
	var buf bytes.Buffer
	gz, err := gzip.NewWriterLevel(&buf, gzip.BestCompression)
	if err != nil {
		return nil, fmt.Errorf("error creating gzip writer: %w", err)
	}
	if _, err := gz.Write(data); err != nil {
		return nil, fmt.Errorf("error writing to gzip writer: %w", err)
	}
	if err := gz.Close(); err != nil {
		return nil, fmt.Errorf("error closing gzip writer: %w", err)
	}
	return buf.Bytes(), nil
}

func (c *HTTPClient) encrypt(data []byte) ([]byte, error) {
	if c.pubKey == nil {
		return data, nil
	}
	return service.Encrypt(c.pubKey, data)
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

	body, err = c.compress(body)
	if err != nil {
		return fmt.Errorf("error compressiong request body: %w", err)
	}

	body, err = c.encrypt(body)
	if err != nil {
		return fmt.Errorf("error encrypting request body: %w", err)
	}

	url := c.serverHost + c.metricEndpoint
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
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
