package storage

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type Stub struct {
	Name  string
	Value int64
}

func TestMemStorageSetAndGet(t *testing.T) {
	tests := []struct {
		name       string
		value      int64
		valueExist bool
	}{
		{
			name:       "success",
			value:      123,
			valueExist: true,
		},
		{
			name:       "failed",
			value:      123,
			valueExist: false,
		},
	}

	repo := NewMemRepo()

	for _, test := range tests {
		stub := &Stub{Name: test.name, Value: test.value}

		if test.valueExist {
			err := repo.Save(test.name, stub)
			assert.NoError(t, err)
		}
	}
	for _, test := range tests {

		t.Run(test.name, func(t *testing.T) {
			stored_stub, ok := repo.Get(test.name)
			assert.Equal(t, test.valueExist, ok)

			if !test.valueExist {
				return
			}

			stub, ok := stored_stub.(*Stub)
			assert.True(t, ok)

			assert.Equal(t, test.name, stub.Name)
			assert.Equal(t, test.value, stub.Value)
		})
	}
}
