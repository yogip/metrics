package metrics

import (
	"metrics/internal/core/model"
)

type Transporter interface {
	SendMetric(req *model.MetricResponse) error
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
		&model.MetricResponse{Name: g.Name, Type: g.Type(), Value: g.StringValue()},
	)
}

type Counter struct {
	model.Counter
}

func (c *Counter) Send(transport Transporter) error {
	err := transport.SendMetric(
		&model.MetricResponse{Name: c.Name, Type: c.Type(), Value: c.StringValue()},
	)
	if err == nil {
		c.Value = 0
	}
	return err
}
