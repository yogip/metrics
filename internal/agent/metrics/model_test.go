package metrics

import (
	"testing"

	"metrics/internal/core/model"

	"github.com/stretchr/testify/assert"
)

func TestPayload(t *testing.T) {
	var delta10 int64 = 10
	var delta101 int64 = 101

	value10 := 10.0
	value101 := 101.0

	tests := []struct {
		want   model.MetricsV2
		metric Metric
		name   string
	}{
		// gauges
		{
			metric: &Gauge{
				model.Gauge{
					Name:  "gauge_test_01",
					Value: value10,
				},
			},
			want: model.MetricsV2{
				ID:    "gauge_test_01",
				MType: model.GaugeType,
				Value: &value10,
			},
		},
		{
			metric: &Gauge{
				model.Gauge{
					Name:  "gauge_test_02",
					Value: value101,
				},
			},
			want: model.MetricsV2{
				ID:    "gauge_test_02",
				MType: model.GaugeType,
				Value: &value101,
			},
		},

		// counters
		{
			metric: &Counter{
				model.Counter{
					Name:  "counter_test_01",
					Value: delta10,
				},
			},
			want: model.MetricsV2{
				ID:    "counter_test_01",
				MType: model.CounterType,
				Delta: &delta10,
			},
		},
		{
			metric: &Counter{
				model.Counter{
					Name:  "counter_test_02",
					Value: delta101,
				},
			},
			want: model.MetricsV2{
				ID:    "counter_test_02",
				MType: model.CounterType,
				Delta: &delta101,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := tt.metric.Payload()
			assert.Equal(t, tt.want, actual)
		})
	}
}

func TestWasSent(t *testing.T) {
	var delta0 int64 = 0
	var delta10 int64 = 10
	var delta101 int64 = 101

	value10 := 10.0
	value101 := 101.0

	tests := []struct {
		want   model.MetricsV2
		metric Metric
		name   string
	}{
		// gauges
		{
			metric: &Gauge{
				model.Gauge{
					Name:  "gauge_test_01",
					Value: value10,
				},
			},
			want: model.MetricsV2{
				ID:    "gauge_test_01",
				MType: model.GaugeType,
				Value: &value10,
			},
		},
		{
			metric: &Gauge{
				model.Gauge{
					Name:  "gauge_test_02",
					Value: value101,
				},
			},
			want: model.MetricsV2{
				ID:    "gauge_test_02",
				MType: model.GaugeType,
				Value: &value101,
			},
		},

		// counters
		{
			metric: &Counter{
				model.Counter{
					Name:  "counter_test_01",
					Value: delta10,
				},
			},
			want: model.MetricsV2{
				ID:    "counter_test_01",
				MType: model.CounterType,
				Delta: &delta0,
			},
		},
		{
			metric: &Counter{
				model.Counter{
					Name:  "counter_test_02",
					Value: delta101,
				},
			},
			want: model.MetricsV2{
				ID:    "counter_test_02",
				MType: model.CounterType,
				Delta: &delta0,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			tt.metric.WasSent()
			actual := tt.metric.Payload()
			assert.Equal(t, tt.want, actual)
		})
	}
}
