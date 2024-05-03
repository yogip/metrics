package models

type MetricType string

const (
	GaugeType   MetricType = "gauge"
	CounterType MetricType = "counter"
)
