package metrics

import (
	"metrics/internal/core/model"
)

type Transporter interface {
	SendMetric(req []*model.MetricsV2) error
}

type Metric interface {
	Payload() *model.MetricsV2
	Type() model.MetricType
	WasSend()
}

type Gauge struct {
	model.Gauge
}

func (g *Gauge) WasSend() {}

func (g *Gauge) Payload() *model.MetricsV2 {
	return &model.MetricsV2{ID: g.Name, MType: g.Type(), Value: &g.Value}
	// return transport.SendMetric(
	// 	&model.MetricsV2{ID: g.Name, MType: g.Type(), Value: &g.Value},
	// )
}

type Counter struct {
	model.Counter
}

func (c *Counter) Payload() *model.MetricsV2 {
	return &model.MetricsV2{ID: c.Name, MType: c.Type(), Delta: &c.Value}
}

func (c *Counter) WasSend() {
	c.Value = 0
}
