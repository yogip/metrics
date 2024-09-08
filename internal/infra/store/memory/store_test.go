package memory

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"metrics/internal/core/config"
	"metrics/internal/core/model"
)

func TestMemStorageSetAndGetCounter(t *testing.T) {
	tests := []struct {
		value      *model.Counter
		name       string
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

	repo, err := NewStore(&config.StorageConfig{
		StoreIntreval:   1000,
		FileStoragePath: "/tmp/storage_dump.json",
		Restore:         false,
	})
	require.NoError(t, err)

	for _, test := range tests {
		if test.valueExist {
			err := repo.SetCounter(context.Background(), test.value)
			require.NoError(t, err)
		}
	}
	for _, test := range tests {

		t.Run(test.name, func(t *testing.T) {
			req := &model.MetricsV2{ID: test.name, MType: model.CounterType}

			counter, err := repo.GetCounter(context.Background(), req)
			require.NoError(t, err)
			assert.Equal(t, test.valueExist, counter != nil)

			if !test.valueExist {
				return
			}

			assert.Equal(t, test.name, counter.Name)
			assert.Equal(t, test.value.Value, counter.Value)
		})
	}
}

func TestMemStorageSetAndGetGauge(t *testing.T) {
	var value01 float64 = 1.01
	var value02 float64 = 2.01
	var delta01 int64 = 1
	var delta02 int64 = 10
	var delta03 int64 = 101
	batch := []*model.MetricsV2{
		{
			ID:    "gauge_01",
			MType: model.GaugeType,
			Value: &value01,
		},
		{
			ID:    "gauge_02",
			MType: model.GaugeType,
			Value: &value02,
		},
		{
			ID:    "counter_01",
			MType: model.CounterType,
			Delta: &delta01,
		},
		{
			ID:    "counter_01",
			MType: model.CounterType,
			Delta: &delta02,
		},
		{
			ID:    "counter_01",
			MType: model.CounterType,
			Delta: &delta03,
		},
	}

	var expDelta01 int64 = delta01
	var expDelta02 int64 = delta01 + delta02
	var expDelta03 int64 = delta01 + delta02 + delta03
	expected := []*model.MetricsV2{
		{
			ID:    "gauge_01",
			MType: model.GaugeType,
			Value: &value01,
		},
		{
			ID:    "gauge_02",
			MType: model.GaugeType,
			Value: &value02,
		},
		{
			ID:    "counter_01",
			MType: model.CounterType,
			Delta: &expDelta01,
		},
		{
			ID:    "counter_01",
			MType: model.CounterType,
			Delta: &expDelta02,
		},
		{
			ID:    "counter_01",
			MType: model.CounterType,
			Delta: &expDelta03,
		},
	}

	repo, err := NewStore(&config.StorageConfig{
		StoreIntreval:   1000,
		FileStoragePath: "/tmp/storage_dump.json",
		Restore:         false,
	})
	require.NoError(t, err)

	actual, err := repo.BatchUpsertMetrics(context.Background(), batch)
	require.NoError(t, err)
	assert.Equal(t, expected, actual)
}

func TestListGauge(t *testing.T) {
	store, err := NewStore(&config.StorageConfig{
		StoreIntreval:   1000,
		FileStoragePath: "/tmp/storage_dump.json",
		Restore:         false,
	})
	require.NoError(t, err)
	ctx := context.Background()

	expected := []*model.Gauge{
		{
			Name:  "gauge_01",
			Value: 12.0,
		},
		{
			Name:  "gauge_02",
			Value: 1.01,
		},
	}
	for _, test := range expected {
		err := store.SetGauge(context.Background(), test)
		require.NoError(t, err)
	}

	actual, err := store.ListGauge(ctx)
	require.NoError(t, err)
	assert.Equal(t, expected, actual)
}

func TestListCounter(t *testing.T) {
	store, err := NewStore(&config.StorageConfig{
		StoreIntreval:   1000,
		FileStoragePath: "/tmp/storage_dump.json",
		Restore:         false,
	})
	require.NoError(t, err)
	ctx := context.Background()

	expected := []*model.Counter{
		{
			Name:  "counter_01",
			Value: 1230,
		},
		{
			Name:  "counter_02",
			Value: 101,
		},
	}
	for _, test := range expected {
		err := store.SetCounter(context.Background(), test)
		require.NoError(t, err)
	}

	actual, err := store.ListCounter(ctx)
	require.NoError(t, err)
	assert.Equal(t, expected, actual)
}

func TestPing(t *testing.T) {
	store, err := NewStore(&config.StorageConfig{
		StoreIntreval:   1000,
		FileStoragePath: "/tmp/storage_dump.json",
		Restore:         false,
	})
	require.NoError(t, err)
	ctx := context.Background()

	err = store.Ping(ctx)
	require.Error(t, err)
}

func TestClose(t *testing.T) {
	ch := make(chan bool)
	store := &Store{
		quit:    ch,
		gauge:   make(map[string]*model.Gauge),
		counter: make(map[string]*model.Counter),
	}

	store.Close()

	var closed bool
	select {
	case <-ch:
		closed = true
	case <-time.After(1 * time.Second):
		closed = false
	}

	assert.True(t, closed)
}
