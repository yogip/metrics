package model

type MetricRequest struct {
	Name string     `uri:"name" binding:"required"`
	Type MetricType `uri:"type" binding:"required" oneof:"gauge counter"`
}

type MetricResponse struct {
	Name  string
	Type  MetricType
	Value string
}

type ListMetricResponse struct {
	Metrics []*MetricResponse
}

type MetricUpdateRequest struct {
	Name  string     `uri:"name" binding:"required"`
	Type  MetricType `uri:"type" binding:"required" oneof:"gauge counter"`
	Value string     `uri:"value" binding:"required"`
}

// ID задае имя метрики
// Delta и Value задают значение метрики в случае передачи counter и gauge соответственно
type MetricsV2 struct {
	ID    string     `json:"id"`
	MType MetricType `json:"type"`
	Delta *int64     `json:"delta,omitempty"`
	Value *float64   `json:"value,omitempty"`
}
