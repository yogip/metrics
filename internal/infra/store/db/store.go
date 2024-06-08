package db

import (
	"database/sql"
	"fmt"
	"metrics/internal/core/config"
	"metrics/internal/logger"

	"github.com/gin-gonic/gin"
)

type Store struct {
	db *sql.DB
}

func NewStore(cfg *config.StorageConfig) (*Store, error) {
	db, err := sql.Open("pgx", cfg.DatabaseDSN)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize Database: %w", err)
	}

	store := &Store{db: db}

	logger.Log.Info("DB Store initialized")
	return store, nil
}

func (s *Store) Close() {
	s.db.Close()
}

func (s *Store) Ping(ctx *gin.Context) error {
	return s.db.PingContext(ctx)
}
