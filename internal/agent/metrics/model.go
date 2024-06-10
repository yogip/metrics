package metrics

import (
	"metrics/internal/core/model"
)

type Transporter interface {
	SendMetric(req *model.MetricsV2) error
}

type Metric interface {
	Send(transport Transporter) error
	Type() model.MetricType
}

type Gauge struct {
	model.Gauge
}

func (g *Gauge) Send(transport Transporter) error {
	return transport.SendMetric(
		&model.MetricsV2{ID: g.Name, MType: g.Type(), Value: &g.Value},
	)
}

type Counter struct {
	model.Counter
}

func (c *Counter) Send(transport Transporter) error {
	err := transport.SendMetric(
		&model.MetricsV2{ID: c.Name, MType: c.Type(), Delta: &c.Value},
	)
	if err != nil {
		return err
	}
	c.Value = 0
	return nil
}
