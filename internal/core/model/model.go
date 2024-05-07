package model

import "fmt"

type MetricRequest struct {
	Name string
	Type MetricType
}

func (m *MetricRequest) ID() string {
	return fmt.Sprintf("%s:%s", m.Type, m.Name)
}

type MetricResponse struct {
	Name  string
	Type  MetricType
	Value string
}

type MetricUpdateRequest struct {
	Name  string
	Type  MetricType
	Value string
}
