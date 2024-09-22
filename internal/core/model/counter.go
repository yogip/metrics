package model

import (
	"fmt"
	"strconv"
)

type Counter struct {
	Name  string
	Value int64
}

func NewCounter(name string) *Counter {
	return &Counter{Name: name}
}

func (c *Counter) Type() MetricType {
	return CounterType
}

func (c *Counter) StringValue() string {
	return strconv.FormatInt(c.Value, 10)
}

func (c *Counter) Increment(value int64) error {
	if value < 0 {
		return fmt.Errorf("could not increment Counter to negative value (%d)", value)
	}
	c.Value += value
	return nil
}
