package transport

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"metrics/internal/core/model"

	"github.com/stretchr/testify/assert"
)

func TestSendMetric(t *testing.T) {
	var ten int64 = 10
	var zero int64 = 0
	var big int64 = 10000000000000000
	var minusOne int64 = -1
	var minusBig int64 = -10000000000000000

	tenFloat := 10.
	zeroFloat := 0.
	bigFloat := 10000000000000000.0
	nearZeroFloat := 0.000000000000001
	minusTenFloat := -10.
	minusBigFloat := -10000000000000000.

	tests := []model.MetricsV2{
		{
			ID:    "counter",
			MType: model.CounterType,
			Delta: &ten,
		},
		{
			ID:    "counter_zero",
			MType: model.CounterType,
			Delta: &zero,
		},
		{
			ID:    "counter_negative",
			MType: model.CounterType,
			Delta: &minusOne,
		},
		{
			MType: model.CounterType,
			ID:    "counter_big_negative",
			Delta: &minusBig,
		},
		{
			ID:    "counter_big",
			MType: model.CounterType,
			Delta: &big,
		},
		{
			MType: model.CounterType,
			ID:    "CaunterCapitalizedName",
			Delta: &ten,
		},

		// Gauges
		{
			MType: model.GaugeType,
			ID:    "gauge",
			Value: &tenFloat,
		},
		{
			MType: model.GaugeType,
			ID:    "gauge_zero",
			Value: &zeroFloat,
		},
		{
			MType: model.GaugeType,
			ID:    "gauge_near_zero",
			Value: &nearZeroFloat,
		},
		{
			MType: model.GaugeType,
			ID:    "gauge_negative",
			Value: &minusTenFloat,
		},
		{
			MType: model.GaugeType,
			ID:    "gauge_big_negative",
			Value: &minusBigFloat,
		},
		{
			MType: model.GaugeType,
			ID:    "gauge_big",
			Value: &bigFloat,
		},
		{
			MType: model.GaugeType,
			ID:    "GaugeCapitalizedName",
			Value: &tenFloat,
		},
	}

	for _, expectedMetric := range tests {
		t.Run(expectedMetric.ID, func(t *testing.T) {
			// Create a test server
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				// Verify the request URL
				assert.Equal(t, "/update", r.URL.String())

				// Verify the request method
				assert.Equal(t, http.MethodPost, r.Method)

				// Verify the request body
				body, err := io.ReadAll(r.Body)
				assert.NoError(t, err)

				var acutalMetric model.MetricsV2
				err = json.Unmarshal(body, &acutalMetric)

				assert.NoError(t, err)
				assert.Equal(t, expectedMetric, acutalMetric)
			}))
			defer server.Close()

			client := NewClient(server.URL)

			// Call the function being tested
			err := client.SendMetric(&expectedMetric)

			// Verify the result
			assert.NoError(t, err)
		})
	}
}
