package model

import "fmt"

type MetricRequest struct {
	Name string     `uri:"name" binding:"required"`
	Type MetricType `uri:"type" binding:"required" oneof=gauge counter`
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
	Name  string     `uri:"name" binding:"required"`
	Type  MetricType `uri:"type" binding:"required" oneof=gauge counter`
	Value string     `uri:"value" binding:"required"`
}
