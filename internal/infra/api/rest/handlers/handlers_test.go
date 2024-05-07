package handlers

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"metrics/internal/core/model"
	"metrics/internal/core/service"
	"metrics/internal/infra/store/memory"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUpdateHandler(t *testing.T) {
	type want struct {
		code     int
		response string
	}
	tests := []struct {
		name       string
		metricName string
		metricType string
		value      string
		method     string
		want       want
	}{
		{
			name:       "counter - positive #1",
			metricName: "name",
			metricType: string(model.CounterType),
			value:      "200",
			method:     http.MethodPost,
			want: want{
				code:     200,
				response: "",
			},
		},
		{
			name:       "counter - positive #2",
			metricName: "name",
			metricType: string(model.CounterType),
			value:      "2",
			method:     http.MethodPost,
			want: want{
				code:     200,
				response: "",
			},
		},
		{
			name:       "counter - zero 0",
			metricName: "name",
			metricType: string(model.CounterType),
			value:      "0",
			method:     http.MethodPost,
			want: want{
				code:     200,
				response: "",
			},
		},
		{
			name:       "counter - zero 0.0",
			metricName: "name",
			metricType: string(model.CounterType),
			value:      "0.0",
			method:     http.MethodPost,
			want: want{
				code:     400,
				response: "failed to parse counter value: could not set value (0.0) to Counter: strconv.ParseInt: parsing \"0.0\": invalid syntax\n",
			},
		},
		{
			name:       "counter - negative zero",
			metricName: "name",
			metricType: string(model.CounterType),
			value:      "-0",
			method:     http.MethodPost,
			want: want{
				code:     200,
				response: "",
			},
		},
		{
			name:       "counter - negative value",
			metricName: "name",
			metricType: string(model.CounterType),
			value:      "-1",
			method:     http.MethodPost,
			want: want{
				code:     400,
				response: "failed to parse counter value: could not set negative value (-1) to Counter\n",
			},
		},
		{
			name:       "counter - negative big value",
			metricName: "name",
			metricType: string(model.CounterType),
			value:      "-1000000000000",
			method:     http.MethodPost,
			want: want{
				code:     400,
				response: "failed to parse counter value: could not set negative value (-1000000000000) to Counter\n",
			},
		},
		{
			name:       "counter - negative float value",
			metricName: "name",
			metricType: string(model.CounterType),
			value:      "-10.0",
			method:     http.MethodPost,
			want: want{
				code:     400,
				response: "failed to parse counter value: could not set value (-10.0) to Counter: strconv.ParseInt: parsing \"-10.0\": invalid syntax\n",
			},
		},
		{
			name:       "counter - not number value",
			metricName: "name",
			metricType: string(model.CounterType),
			value:      "incorrect",
			method:     http.MethodPost,
			want: want{
				code:     400,
				response: "failed to parse counter value: could not set value (incorrect) to Counter: strconv.ParseInt: parsing \"incorrect\": invalid syntax\n",
			},
		},
		{
			name:       "counter - value 00",
			metricName: "name",
			metricType: string(model.CounterType),
			value:      "00",
			method:     http.MethodPost,
			want: want{
				code:     200,
				response: "",
			},
		},
		{
			name:       "counter - get method",
			metricName: "name",
			metricType: string(model.CounterType),
			value:      "00",
			method:     http.MethodGet,
			want: want{
				code:     405,
				response: "Method Not Allowed\n",
			},
		},

		// Gauge test cases
		{
			name:       "gauge - positive #1",
			metricName: "name",
			metricType: string(model.GaugeType),
			value:      "20.1230",
			method:     http.MethodPost,
			want: want{
				code:     200,
				response: "",
			},
		},
		{
			name:       "gauge - positive #2",
			metricName: "name",
			metricType: string(model.GaugeType),
			value:      "0.1230111",
			method:     http.MethodPost,
			want: want{
				code:     200,
				response: "",
			},
		},
		{
			name:       "gauge - zero 0",
			metricName: "name",
			metricType: string(model.GaugeType),
			value:      "0",
			method:     http.MethodPost,
			want: want{
				code:     200,
				response: "",
			},
		},
		{
			name:       "gauge - zero 0.0",
			metricName: "name",
			metricType: string(model.GaugeType),
			value:      "0.0",
			method:     http.MethodPost,
			want: want{
				code:     200,
				response: "",
			},
		},
		{
			name:       "gauge - negative zero",
			metricName: "name",
			metricType: string(model.GaugeType),
			value:      "-0",
			method:     http.MethodPost,
			want: want{
				code:     200,
				response: "",
			},
		},
		{
			name:       "gauge - negative value",
			metricName: "name",
			metricType: string(model.GaugeType),
			value:      "-1",
			method:     http.MethodPost,
			want: want{
				code:     200,
				response: "",
			},
		},
		{
			name:       "gauge - big value",
			metricName: "name",
			metricType: string(model.GaugeType),
			value:      "1000000000000.123",
			method:     http.MethodPost,
			want: want{
				code:     200,
				response: "",
			},
		},
		{
			name:       "gauge - negative value 2",
			metricName: "name",
			metricType: string(model.GaugeType),
			value:      "-10.0",
			method:     http.MethodPost,
			want: want{
				code:     200,
				response: "",
			},
		},
		{
			name:       "gauge - int value",
			metricName: "name",
			metricType: string(model.GaugeType),
			value:      "100",
			method:     http.MethodPost,
			want: want{
				code:     200,
				response: "",
			},
		},
		{
			name:       "gauge - not number value",
			metricName: "name",
			metricType: string(model.GaugeType),
			value:      "incorrect",
			method:     http.MethodPost,
			want: want{
				code:     400,
				response: "failed to parse gauge value: could not set value (incorrect) to Gauge: strconv.ParseFloat: parsing \"incorrect\": invalid syntax\n",
			},
		},
		{
			name:       "gauge - value 00",
			metricName: "name",
			metricType: string(model.GaugeType),
			value:      "00",
			method:     http.MethodPost,
			want: want{
				code:     200,
				response: "",
			},
		},
	}
	store := memory.NewStore()
	service := service.NewMetricService(store)
	handler := NewHandler(service)

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			url := fmt.Sprintf("/update/%s/%s/%s", test.metricType, test.metricName, test.value)
			request := httptest.NewRequest(test.method, url, nil)

			w := httptest.NewRecorder()
			handler.UpdateHandler(w, request)

			res := w.Result()
			assert.Equal(t, test.want.code, res.StatusCode)

			defer res.Body.Close()
			resBody, err := io.ReadAll(res.Body)

			require.NoError(t, err)
			assert.Equal(t, test.want.response, string(resBody))
		})
	}
}
