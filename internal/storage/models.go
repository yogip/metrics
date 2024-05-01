package storage

import (
	"fmt"
	"strconv"
)

type Metric interface {
	Type() MetricType
	Set(string) error
	String() string
	pk() string
}

type Gauge struct {
	Name  string
	Value float64
}

func (g *Gauge) pk() string {
	return pkey(g.Type(), g.Name)
}

func (g *Gauge) Type() MetricType {
	return GaugeType
}

func (g *Gauge) String() string {
	return fmt.Sprintf("<Gauge %s: %f>", g.Name, g.Value)
}

// Set and convert value from sting, return error for wrong type
func (g *Gauge) Set(value string) error {
	if v, err := strconv.ParseFloat(value, 64); err != nil {
		return fmt.Errorf("could not set value (%s) to Gauge: %s", value, err)
	} else {
		g.Value = v
	}
	return nil
}

type Counter struct {
	Name  string
	Value int64
}

func (c *Counter) pk() string {
	return pkey(c.Type(), c.Name)
}

func (c *Counter) Type() MetricType {
	return CounterType
}

func (c *Counter) String() string {
	return fmt.Sprintf("<Countre %s: %d>", c.Name, c.Value)
}

// Set and convert value from sting, return error for wrong type
func (c *Counter) Set(value string) error {
	v, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return fmt.Errorf("could not set value (%s) to Counter: %s", value, err)
	}
	if v < 0 {
		return fmt.Errorf("could not set negative value (%s) to Counter", value)
	}
	c.Value += v
	return nil
}

// Init New Metric object. Without saving at storage.
func NewMetric(mType MetricType, name string) (Metric, error) {
	switch mType {
	case GaugeType:
		gauge := &Gauge{Name: name}
		// if err := gauge.Set(value); err != nil {
		// return nil, fmt.Errorf("could not set value: %s for %s", err, gauge)
		// }
		return gauge, nil
	case CounterType:
		counter := &Counter{Name: name}
		// if err := counter.Set(value); err != nil {
		// 	return nil, fmt.Errorf("could not set value: %s for %s", err, counter)
		// }
		return counter, nil
	default:
		return nil, fmt.Errorf("unknown metric type: %s", mType)
	}
}

func GetMetric(mType MetricType, name string) (Metric, bool) {
	pk := pkey(mType, name)
	m, ok := storage.Get(pk)
	return m, ok
}

func SaveMetric(m Metric) error {
	return storage.Save(m.pk(), m)
}

func pkey(metricType MetricType, metricName string) string {
	return fmt.Sprintf("%s:%s", metricType, metricName)
}
