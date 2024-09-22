package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"metrics/internal/core/config"
	"metrics/internal/core/model"
	"metrics/internal/core/service"
	"metrics/internal/infra/store/memory"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUpdateHandler(t *testing.T) {
	var ten int64 = 10

	type want struct {
		response string
		code     int
	}
	tests := []struct {
		metric model.MetricsV2
		method string
		want   want
	}{
		{
			metric: model.MetricsV2{
				ID:    "counter",
				MType: model.CounterType,
				Delta: &ten,
			},
			method: http.MethodPost,
			want: want{
				code:     200,
				response: `{"delta":10,"id":"counter","type":"counter"}`,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.metric.ID, func(t *testing.T) {
			// Создаем новый обработчик с поддельным сервисом
			store, err := memory.NewStore(&config.StorageConfig{
				StoreIntreval:   1000,
				FileStoragePath: "/tmp/storage_dump.json",
				Restore:         false,
			})
			require.NoError(t, err)
			service := service.NewMetricService(store)
			handler := NewHandlerV2(service)

			// Создаем новый контекст gin с тестовым запросом
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			// Convert tt.metric to JSON format
			body, err := json.Marshal(tt.metric)
			require.NoError(t, err)

			c.Request = httptest.NewRequest(tt.method, "/update", strings.NewReader(string(body)))

			// Вызываем обработчик
			handler.UpdateHandler(c)

			// Проверяем результаты
			assert.Equal(t, tt.want.code, w.Code)
			assert.JSONEq(t, tt.want.response, w.Body.String())
		})
	}
}

func TestUpdateRequestsWithTheSameStore(t *testing.T) {
	var ten int64 = 10

	type want struct {
		response string
		code     int
	}
	tests := []struct {
		metric model.MetricsV2
		method string
		want   want
	}{
		{
			metric: model.MetricsV2{
				ID:    "counter01",
				MType: model.CounterType,
				Delta: &ten,
			},
			method: http.MethodPost,
			want: want{
				code:     200,
				response: `{"delta":10, "id":"counter01", "type":"counter"}`,
			},
		},
		{
			metric: model.MetricsV2{
				ID:    "counter01",
				MType: model.CounterType,
				Delta: &ten,
			},
			method: http.MethodPost,
			want: want{
				code:     200,
				response: `{"delta":20, "id":"counter01", "type":"counter"}`,
			},
		},
	}

	// один стор для всех запросов, результат будет накопительный
	store, err := memory.NewStore(&config.StorageConfig{
		StoreIntreval:   1000,
		FileStoragePath: "/tmp/storage_dump.json",
		Restore:         false,
	})
	require.NoError(t, err)
	for _, tt := range tests {
		// Создаем новый обработчик с поддельным сервисом
		service := service.NewMetricService(store)
		handler := NewHandlerV2(service)

		// Создаем новый контекст gin с тестовым запросом
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		// Convert tt.metric to JSON format
		body, err := json.Marshal(tt.metric)
		require.NoError(t, err)

		c.Request = httptest.NewRequest(tt.method, "/update", strings.NewReader(string(body)))

		// Вызываем обработчик
		handler.UpdateHandler(c)

		// Проверяем результаты
		assert.Equal(t, tt.want.code, w.Code)
		assert.JSONEq(t, tt.want.response, w.Body.String())
	}
}

func TestGetHandler(t *testing.T) {
	var ten int64 = 10

	type want struct {
		response string
		code     int
	}
	tests := []struct {
		metric model.MetricsV2
		method string
		want   want
	}{
		{
			metric: model.MetricsV2{
				ID:    "counter",
				MType: model.CounterType,
				Delta: &ten,
			},
			method: http.MethodPost,
			want: want{
				code:     200,
				response: `{"delta":10, "id":"counter", "type":"counter"}`,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.metric.ID, func(t *testing.T) {
			// Создаем новый обработчик с поддельным сервисом
			store, err := memory.NewStore(&config.StorageConfig{
				StoreIntreval:   1000,
				FileStoragePath: "/tmp/storage_dump.json",
				Restore:         false,
			})
			require.NoError(t, err)
			store.SetCounter(context.Background(), &model.Counter{Name: tt.metric.ID, Value: *tt.metric.Delta})
			service := service.NewMetricService(store)
			handler := NewHandlerV2(service)

			// Создаем новый контекст gin с тестовым запросом
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			body, err := json.Marshal(model.MetricsV2{ID: tt.metric.ID, MType: tt.metric.MType})
			require.NoError(t, err)

			c.Request = httptest.NewRequest(tt.method, "/value", strings.NewReader(string(body)))

			// Вызываем обработчик
			handler.GetHandler(c)

			// Проверяем результаты
			assert.Equal(t, tt.want.code, w.Code)
			assert.JSONEq(t, tt.want.response, w.Body.String())
		})
	}
}

func ExampleHandlerV2_UpdateHandler() {
	var ten int64 = 10

	metric := model.MetricsV2{
		ID:    "counter",
		MType: model.CounterType,
		Delta: &ten,
	}

	store, _ := memory.NewStore(&config.StorageConfig{Restore: false})

	service := service.NewMetricService(store)
	handler := NewHandlerV2(service)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	body, _ := json.Marshal(metric)

	c.Request = httptest.NewRequest(http.MethodPost, "/update", strings.NewReader(string(body)))

	handler.UpdateHandler(c)

	fmt.Println(w.Code)
	fmt.Println(w.Body.String())

	// Output:
	// 200
	// {"delta":10,"id":"counter","type":"counter"}
}
