package model

import (
	"fmt"
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

func (g *Gauge) String() string {
	return fmt.Sprintf("<Gauge %s: %f>", g.Name, g.Value)
}

func (g *Gauge) GetName() string {
	return g.Name
}

func (g *Gauge) StringValue() string {
	return strconv.FormatFloat(g.Value, 'f', -1, 64)
}

func (g *Gauge) Set(value float64) {
	g.Value = value
}

// func (g *Gauge) ParseString(value string) error {
// 	if v, err := strconv.ParseFloat(value, 64); err != nil {
// 		return fmt.Errorf("could not set value (%s) to Gauge: %s", value, err)
// 	} else {
// 		g.Value = v
// 	}
// 	return nil
// }
