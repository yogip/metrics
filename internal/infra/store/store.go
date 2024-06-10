package store

import (
	"context"
	"metrics/internal/core/config"
	"metrics/internal/core/model"
	"metrics/internal/infra/store/db"
	"metrics/internal/infra/store/memory"
)

type Store interface {
	BatchUpsertMetrics(ctx context.Context, metrics []*model.MetricsV2) ([]*model.MetricsV2, error)
	GetGauge(ctx context.Context, req *model.MetricsV2) (*model.Gauge, error)
	SetGauge(ctx context.Context, gauge *model.Gauge) error
	ListGauge(ctx context.Context) ([]*model.Gauge, error)
	GetCounter(ctx context.Context, req *model.MetricsV2) (*model.Counter, error)
	SetCounter(ctx context.Context, counter *model.Counter) error
	ListCounter(ctx context.Context) ([]*model.Counter, error)
	Ping(ctx context.Context) error
	Close()
}

func NewStore(cfg *config.StorageConfig) (Store, error) {
	if cfg.DatabaseDSN != "" {
		return db.NewStore(cfg)
	}
	return memory.NewStore(cfg)
}
