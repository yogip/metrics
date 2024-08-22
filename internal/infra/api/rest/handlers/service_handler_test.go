package handlers

import (
	"errors"
	"net/http/httptest"
	"testing"

	"metrics/internal/core/service"
	"metrics/internal/mocks"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestSystemHandler_Ping(t *testing.T) {
	type want struct {
		code     int
		response string
	}

	ctrl := gomock.NewController(t)

	tests := []struct {
		name        string
		method      string
		returnValue error
		want        want
	}{
		{
			name:        "success",
			returnValue: nil,
			want: want{
				code:     200,
				response: "OK",
			},
		},
		{
			name:        "error",
			returnValue: errors.New("Some error"),
			want: want{
				code:     500,
				response: "Internal Server Error",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Создаем новый контекст gin с тестовым запросом
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			mock := mocks.NewMockPinger(ctrl)
			mock.EXPECT().Ping(c).Return(tt.returnValue)
			handler := NewSystemHandler(service.NewSystemService(mock))

			// Вызываем обработчик
			handler.Ping(c)

			// Проверяем результаты
			assert.Equal(t, tt.want.code, w.Code)
			assert.Equal(t, tt.want.response, w.Body.String())
		})
	}
}
