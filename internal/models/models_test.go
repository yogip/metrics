package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
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
		err := g.ParseString(tc.Value)
		assert.NoError(t, err)
		assert.Equal(t, tc.Expected, g.Value)
	}
}

func TestGaugeSetFailed(t *testing.T) {
	g := &Gauge{Name: "test", Value: 0}
	err := g.ParseString("invalid")

	assert.Error(t, err)
}

func TestCounterSet(t *testing.T) {
	c := &Counter{Name: "test", Value: 0}

	testCases := []CounterTestCase{
		{"10", 10},
		{"10", 20},
		{"1", 21},
		{"0", 21},
		{"11122", 11143},
	}

	for _, tc := range testCases {
		err := c.ParseString(tc.Value)
		assert.NoError(t, err)
		assert.Equal(t, tc.Expected, c.Value)
	}
}

func TestCounterSetFailed(t *testing.T) {
	c := &Counter{Name: "test", Value: 0}

	for _, v := range []string{"invalid", "0.1", "0.00001", "2.000001", "-0.000000003", "-1"} {
		err := c.ParseString(v)
		if err == nil {
			t.Errorf("Tried to set %v and expecting error, got nil", v)
		}
	}
}
