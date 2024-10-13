package store

import (
	"context"
	"sync"
	"testing"

	"metrics/internal/core/config"
	"metrics/internal/infra/store/db"
	"metrics/internal/infra/store/memory"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/stretchr/testify/require"
)

func TestNewDBStore(t *testing.T) {
	cfg := &config.StorageConfig{
		DatabaseDSN: "...",
	}
	ctx := context.Background()

	var wg sync.WaitGroup
	actual, err := NewStore(ctx, &wg, cfg)
	require.NoError(t, err)
	require.IsType(t, &db.Store{}, actual)
}

func TestNewMemoryStore(t *testing.T) {
	ctx := context.Background()
	cfg := &config.StorageConfig{
		Restore: false,
	}

	var wg sync.WaitGroup
	actual, err := NewStore(ctx, &wg, cfg)
	require.NoError(t, err)
	require.IsType(t, &memory.Store{}, actual)
}
