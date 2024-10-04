// The package provides persistent metrics storage.
// The Store has two backend implementations in memry with periodic dump to a file and in a database (PostgreSQL).
package store

import (
	"context"
	"sync"

	"metrics/internal/core/config"
	"metrics/internal/core/model"
	"metrics/internal/infra/store/db"
	"metrics/internal/infra/store/memory"
)

// Store interfave for all public methods.
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

// NewStore create new Store object.
// The Store will use database backend If the environment variable DATABASE_DSN or -d command arg is specified,
// otherwise memory storage.
func NewStore(ctx context.Context, wg *sync.WaitGroup, cfg *config.StorageConfig) (Store, error) {
	if cfg.DatabaseDSN != "" {
		return db.NewStore(ctx, wg, cfg)
	}
	return memory.NewStore(ctx, wg, cfg)
}
