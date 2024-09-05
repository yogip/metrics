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

// MetricsV2 описывает схему ответа и запроса для метрик.
type MetricsV2 struct {
	Delta *int64     `json:"delta,omitempty"`
	Value *float64   `json:"value,omitempty"`
	ID    string     `json:"id"`
	MType MetricType `json:"type"`
}
