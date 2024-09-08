package store

import (
	"metrics/internal/core/config"
	"metrics/internal/infra/store/db"
	"metrics/internal/infra/store/memory"
	"testing"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/stretchr/testify/require"
)

func TestNewDBStore(t *testing.T) {
	cfg := &config.StorageConfig{
		DatabaseDSN: "...",
	}

	actual, err := NewStore(cfg)
	require.NoError(t, err)
	require.IsType(t, &db.Store{}, actual)
}

func TestNewMemoryStore(t *testing.T) {
	cfg := &config.StorageConfig{
		Restore: false,
	}

	actual, err := NewStore(cfg)
	require.NoError(t, err)
	require.IsType(t, &memory.Store{}, actual)
}
