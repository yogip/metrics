package models

import (
	"fmt"
	"strconv"

	"github.com/yogip/metrics/internal/storage"
)

type Metric interface {
	Type() MetricType
	ParseString(string) error
	String() string
	StringValue() string
	GetName() string
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

func (g *Gauge) GetName() string {
	return g.Name
}

func (g *Gauge) String() string {
	return fmt.Sprintf("<Gauge %s: %f>", g.Name, g.Value)
}

func (g *Gauge) StringValue() string {
	return strconv.FormatFloat(g.Value, 'f', -1, 64)
}

func (g *Gauge) Set(value float64) {
	g.Value = value
}

// Set and convert value from sting, return error for wrong type
func (g *Gauge) ParseString(value string) error {
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

func (c *Counter) GetName() string {
	return c.Name
}

func (c *Counter) String() string {
	return fmt.Sprintf("<Countre %s: %d>", c.Name, c.Value)
}

func (c *Counter) StringValue() string {
	return strconv.FormatInt(c.Value, 10)
}

func (c *Counter) Incremet(value int64) {
	c.Value += value
}

// Set and convert value from sting, return error for wrong type
func (c *Counter) ParseString(value string) error {
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
	m, ok := storage.Storage.Get(pk)
	if !ok {
		return nil, false
	}
	metric, ok := m.(Metric)
	if !ok {
		return nil, false
	}
	return metric, ok
}

func SaveMetric(m Metric) error {
	return storage.Storage.Save(m.pk(), m)
}

func pkey(metricType MetricType, metricName string) string {
	return fmt.Sprintf("%s:%s", metricType, metricName)
}
