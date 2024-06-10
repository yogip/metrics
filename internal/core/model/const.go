package model

type MetricType string

func (mt *MetricType) String() string {
	return string(*mt)
}

const (
	GaugeType   MetricType = "gauge"
	CounterType MetricType = "counter"
)
