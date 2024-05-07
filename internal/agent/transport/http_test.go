package transport

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"metrics/internal/core/model"

	"github.com/stretchr/testify/assert"
)

func TestSendMetric(t *testing.T) {
	tests := []struct {
		mType model.MetricType
		name  string
		value string
	}{
		{
			mType: model.CounterType,
			name:  "counter",
			value: "3",
		},
		{
			mType: model.CounterType,
			name:  "counter_zero",
			value: "0",
		},
		{
			mType: model.CounterType,
			name:  "counter_negative_zero",
			value: "-0",
		},
		{
			mType: model.CounterType,
			name:  "counter_negative",
			value: "-1",
		},
		{
			mType: model.CounterType,
			name:  "counter_big_negative",
			value: "-10000000000000000",
		},
		{
			mType: model.CounterType,
			name:  "counter_big",
			value: "10000000000000000",
		},
		{
			mType: model.GaugeType,
			name:  "gauge",
			value: "3.0",
		},
		{
			mType: model.GaugeType,
			name:  "gauge_with_int",
			value: "3",
		},
		{
			mType: model.GaugeType,
			name:  "gauge_zero",
			value: "0",
		},
		{
			mType: model.GaugeType,
			name:  "gauge_zero_2",
			value: "0.0",
		},
		{
			mType: model.GaugeType,
			name:  "gauge_zero_3",
			value: ".0",
		},
		{
			mType: model.GaugeType,
			name:  "gauge_near_zero",
			value: "0.0000000000001",
		},
		{
			mType: model.GaugeType,
			name:  "gauge_negative_zero",
			value: "-0",
		},
		{
			mType: model.GaugeType,
			name:  "gauge_negative",
			value: "-1.01",
		},
		{
			mType: model.GaugeType,
			name:  "gauge_big_negative",
			value: "-100000000.1",
		},
		{
			mType: model.GaugeType,
			name:  "gauge_big",
			value: "1000000000000000.2",
		},
		{
			mType: model.GaugeType,
			name:  "GaugeCapitalizedName",
			value: "12.12",
		},
		{
			mType: model.CounterType,
			name:  "CaunterCapitalizedName",
			value: "12",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Create a test server
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				// Verify the request URL
				expectedURL := fmt.Sprintf("/update/%s/%s/%s", test.mType, test.name, test.value)
				assert.Equal(t, expectedURL, r.URL.String())

				// Verify the request method
				assert.Equal(t, http.MethodPost, r.Method)
			}))
			defer server.Close()

			client := NewClient(server.URL)

			// Call the function being tested
			err := client.SendMetric(&model.MetricResponse{
				Type:  test.mType,
				Name:  test.name,
				Value: test.value,
			})

			// Verify the result
			assert.NoError(t, err)
		})
	}
}
