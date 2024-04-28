package repo

import (
	"fmt"
	"strconv"
)

type Metric interface {
	Set(string) error
	Type() MetricType
	ValidateType(MetricType) bool
}

type Gauge struct {
	Name  string
	Value float64
}

func (g *Gauge) Type() MetricType {
	return GaugeType
}

func (g *Gauge) ValidateType(t MetricType) bool {
	return t == g.Type()
}

// Set and convert value from sting, return error for wrong type
func (c *Gauge) Set(value string) error {
	if v, err := strconv.ParseFloat(value, 64); err != nil {
		return fmt.Errorf("could not set value (%s) to Gauge: %s", value, err)
	} else {
		c.Value = v
	}
	return nil
}

type Counter struct {
	Name  string
	Value int64
}

func (c *Counter) Type() MetricType {
	return CounterType
}

func (c *Counter) ValidateType(t MetricType) bool {
	return t == c.Type()
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

type MemStorage struct {
	storage map[string]Metric
}

func NewMemStorage() *MemStorage {
	return &MemStorage{storage: make(map[string]Metric)}
}

func (s *MemStorage) Set(metricType MetricType, metricName string, metric Metric) {
	s.storage[fmt.Sprintf("%s:%s", metricType, metricName)] = metric
}

func (s *MemStorage) Get(metricType MetricType, metricName string) (Metric, bool) {
	m, ok := s.storage[fmt.Sprintf("%s:%s", metricType, metricName)]
	return m, ok
}
