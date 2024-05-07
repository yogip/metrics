package memory

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"metrics/internal/core/model"
)

func TestMemStorageSetAndGetCounter(t *testing.T) {
	tests := []struct {
		name       string
		value      *model.Counter
		valueExist bool
	}{
		{
			name:       "success",
			value:      &model.Counter{Name: "success", Value: 123},
			valueExist: true,
		},
		{
			name:       "failed",
			value:      &model.Counter{Name: "failed", Value: 123},
			valueExist: false,
		},
	}

	repo := NewStore()

	for _, test := range tests {
		req := &model.MetricRequest{Name: test.name, Type: model.CounterType}

		if test.valueExist {
			err := repo.SetCounter(req, test.value)
			assert.NoError(t, err)
		}
	}
	for _, test := range tests {

		t.Run(test.name, func(t *testing.T) {
			req := &model.MetricRequest{Name: test.name, Type: model.CounterType}

			counter, err := repo.GetCounter(req)
			assert.NoError(t, err)
			assert.Equal(t, test.valueExist, counter != nil)

			if !test.valueExist {
				return
			}

			assert.Equal(t, test.name, counter.Name)
			assert.Equal(t, test.value.Value, counter.Value)
		})
	}
}
