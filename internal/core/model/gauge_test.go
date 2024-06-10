package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGaugeSetSuccess(t *testing.T) {
	g := &Gauge{Name: "test", Value: 0}

	testCases := []float64{10, 9992, 10.0, 10.05, 0.0000001, -10.0, -1000.000001}

	for _, tc := range testCases {
		g.Set(tc)
		assert.Equal(t, tc, g.Value)
	}
}
