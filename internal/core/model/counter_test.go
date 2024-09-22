package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCounterIncrement(t *testing.T) {
	c := &Counter{Name: "test", Value: 0}

	testCases := []struct {
		Increment int64
		Expected  int64
	}{
		{10, 10},
		{10, 20},
		{1, 21},
		{0, 21},
		{11122, 11143},
	}

	for _, tc := range testCases {
		err := c.Increment(tc.Increment)
		require.NoError(t, err)
		assert.Equal(t, tc.Expected, c.Value)
	}
}

func TestCounterSetFailed(t *testing.T) {
	c := &Counter{Name: "test", Value: 0}

	for _, v := range []int64{-10, -15, -5} {
		err := c.Increment(v)
		assert.Error(t, err)
	}
}
