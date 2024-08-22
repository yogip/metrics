package service

import (
	"context"
	"testing"

	"metrics/internal/core/model"
	"metrics/internal/mocks"

	"go.uber.org/mock/gomock"
)

func BenchmarkListMetrics(b *testing.B) {
	ctrl := gomock.NewController(b)

	ctx := context.Background()

	gauges := []*model.Gauge{
		{Name: "gauge_1", Value: 1.01},
		{Name: "gauge_2", Value: 2.01},
	}

	counters := []*model.Counter{
		{Name: "counter_1", Value: 1},
		{Name: "counter_2", Value: 2},
	}

	mock := mocks.NewMockStore(ctrl)
	mock.EXPECT().ListGauge(ctx).Return(gauges, nil).AnyTimes()
	mock.EXPECT().ListCounter(ctx).Return(counters, nil).AnyTimes()

	metricService := NewMetricService(mock)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = metricService.ListMetrics(ctx)
	}
}

func BenchmarkGetCounter(b *testing.B) {
	ctrl := gomock.NewController(b)

	ctx := context.Background()
	counter := model.Counter{Name: "counter_1", Value: 1}
	req := model.MetricsV2{ID: counter.Name, MType: model.CounterType, Delta: &counter.Value}

	mock := mocks.NewMockStore(ctrl)
	mock.EXPECT().GetCounter(ctx, &req).Return(&counter, nil).AnyTimes()

	metricService := NewMetricService(mock)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = metricService.GetCounter(ctx, &req)
	}
}

func BenchmarkGetGauge(b *testing.B) {
	ctrl := gomock.NewController(b)

	ctx := context.Background()
	gauge := model.Gauge{Name: "gauge_1", Value: 1}
	req := model.MetricsV2{ID: gauge.Name, MType: model.CounterType, Value: &gauge.Value}

	mock := mocks.NewMockStore(ctrl)
	mock.EXPECT().GetGauge(ctx, &req).Return(&gauge, nil).AnyTimes()

	metricService := NewMetricService(mock)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = metricService.GetGauge(ctx, &req)
	}
}
