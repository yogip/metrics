package repo

import (
	"testing"
)

type GaugeTestCase struct {
	Value    string
	Expected float64
}

type GaugeFailedTestCase struct {
	Value string
}

type CounterTestCase struct {
	Value    string
	Expected int64
}

func TestGaugeSetSuccess(t *testing.T) {
	g := &Gauge{Name: "test", Value: 0}

	testCases := []GaugeTestCase{
		{"10", 10.0},
		{"10.5", 10.5},
		{"0.5", 0.5},
		{"0.0000001", 0.0000001},
		{"-10", -10.},
		{"-0.1", -0.1},
		{".01", 0.01},
	}

	for _, tc := range testCases {
		err := g.Set(tc.Value)
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}

		if g.Value != tc.Expected {
			t.Errorf("Expected value to be %v, got %v", tc.Expected, g.Value)
		}
	}
}

func TestGaugeSetFailed(t *testing.T) {
	g := &Gauge{Name: "test", Value: 0}
	err := g.Set("invalid")
	if err == nil {
		t.Error("Expected error, got nil")
	}
}

func TestCounterSet(t *testing.T) {
	c := &Counter{Name: "test", Value: 0}

	testCases := []CounterTestCase{
		{"10", 10},
		{"10", 20},
		{"1", 21},
		{"0", 21},
		{"11122", 11122},
	}

	for _, tc := range testCases {
		err := c.Set(tc.Value)
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}

		if c.Value != tc.Expected {
			t.Errorf("Expected value to be %v, got %v", tc.Expected, c.Value)
		}
	}
}

func TestCounterSetFailed(t *testing.T) {
	c := &Counter{Name: "test", Value: 0}
	for _, v := range []string{"invalid", "0.1", "0.00001", "-1"} {
		err := c.Set(v)
		if err == nil {
			t.Errorf("Tried to set %v and expecting error, got nil", v)
		}
	}
}

func TestMemStorageSetAndGet(t *testing.T) {
	s := NewMemStorage()

	g := &Gauge{Name: "test", Value: 0}
	err := g.Set("10")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	s.Set(GaugeType, "gauge_1", g)

	metric, ok := s.Get(GaugeType, "gauge_1")
	if !ok {
		t.Error("Expected metric to exist, got false")
	}

	if metric != g {
		t.Errorf("Expected metric to be %v, got %v", g, metric)
	}
}
