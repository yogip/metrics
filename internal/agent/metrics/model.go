package metrics

import (
	"metrics/internal/core/model"
)

type Transporter interface {
	SendMetric(req []model.MetricsV2) error
}

type Metric interface {
	Payload() model.MetricsV2
	Type() model.MetricType
	WasSent()
}

type Gauge struct {
	model.Gauge
}

func (g *Gauge) WasSent() {}

func (g *Gauge) Payload() model.MetricsV2 {
	return model.MetricsV2{ID: g.Name, MType: g.Type(), Value: &g.Value}
}

type Counter struct {
	model.Counter
}

func (c *Counter) Payload() model.MetricsV2 {
	value := c.Value // Create new variable due to prevent of erasing by WasSent
	return model.MetricsV2{ID: c.Name, MType: c.Type(), Delta: &value}
}

func (c *Counter) WasSent() {
	c.Value = 0
}
