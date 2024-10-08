package model

import (
	"strconv"
)

type Gauge struct {
	Name  string
	Value float64
}

func NewGauge(name string) *Gauge {
	return &Gauge{Name: name}
}

func (g *Gauge) Type() MetricType {
	return GaugeType
}

func (g *Gauge) StringValue() string {
	return strconv.FormatFloat(g.Value, 'f', -1, 64)
}

func (g *Gauge) Set(value float64) {
	g.Value = value
}
