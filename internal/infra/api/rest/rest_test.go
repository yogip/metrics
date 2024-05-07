package rest

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"metrics/internal/core/model"
	"metrics/internal/core/service"
	"metrics/internal/infra/store/memory"

	"github.com/stretchr/testify/assert"
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
			name:       "incorrect_type",
			metricName: "name",
			metricType: "incorrectType",
			value:      "200",
			method:     http.MethodPost,
			want: want{
				code:     400,
				response: "unknown metric type: incorrectType",
			},
		},
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
				response: "failed to parse counter value: could not set value (0.0) to Counter: strconv.ParseInt: parsing \"0.0\": invalid syntax",
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
				response: "failed to parse counter value: could not set negative value (-1) to Counter",
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
				response: "failed to parse counter value: could not set negative value (-1000000000000) to Counter",
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
				response: "failed to parse counter value: could not set value (-10.0) to Counter: strconv.ParseInt: parsing \"-10.0\": invalid syntax",
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
				response: "failed to parse counter value: could not set value (incorrect) to Counter: strconv.ParseInt: parsing \"incorrect\": invalid syntax",
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
				code:     404,
				response: "404 page not found",
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
				response: "failed to parse gauge value: could not set value (incorrect) to Gauge: strconv.ParseFloat: parsing \"incorrect\": invalid syntax",
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

	api := NewAPI(service)

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			url := fmt.Sprintf("/update/%s/%s/%s", test.metricType, test.metricName, test.value)
			request := httptest.NewRequest(test.method, url, nil)

			w := httptest.NewRecorder()
			api.srv.ServeHTTP(w, request)

			assert.Equal(t, test.want.code, w.Code)
			assert.Equal(t, test.want.response, w.Body.String())
		})
	}
}
