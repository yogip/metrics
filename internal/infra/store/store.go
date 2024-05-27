package store

import (
	"metrics/internal/core/config"
	"metrics/internal/core/model"
	"metrics/internal/infra/store/memory"
)

type Store interface {
	GetGauge(req *model.MetricRequest) (*model.Gauge, error)
	SetGauge(req *model.MetricRequest, gauge *model.Gauge) error
	ListGauge() ([]*model.Gauge, error)
	GetCounter(req *model.MetricRequest) (*model.Counter, error)
	SetCounter(req *model.MetricRequest, counter *model.Counter) error
	ListCounter() ([]*model.Counter, error)
	Close()
}

func NewStore(cfg *config.StorageConfig) (Store, error) {
	return memory.NewStore(cfg)
}
