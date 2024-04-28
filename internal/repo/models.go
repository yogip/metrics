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

func (s *MemStorage) pkey(metricType MetricType, metricName string) string {
	return fmt.Sprintf("%s:%s", metricType, metricName)
}

func (s *MemStorage) SetValue(metricType MetricType, metricName string, value string) error {
	metric, error := s.GetOrCreate(metricType, metricName)
	if error != nil {
		return fmt.Errorf("could not get or create metric: %s", error)
	}

	if err := metric.Set(value); err != nil {
		return fmt.Errorf("could not set value: %s", error)
	}

	// use pointer due to remove extra copy
	pk := s.pkey(metricType, metricName)
	s.storage[pk] = metric
	return nil
}

func (s *MemStorage) Get(metricType MetricType, metricName string) (Metric, bool) {
	pk := s.pkey(metricType, metricName)
	m, ok := s.storage[pk]
	return m, ok
}

func (s *MemStorage) GetOrCreate(metricType MetricType, metricName string) (Metric, error) {
	metric, ok := s.Get(metricType, metricName)
	if ok {
		return metric, nil
	}

	switch metricType {
	case GaugeType:
		return &Gauge{Name: metricName, Value: 0}, nil
	case CounterType:
		return &Counter{Name: metricName, Value: 0}, nil
	default:
		return nil, fmt.Errorf("unknown metric type: %s", metricType)
	}
}
