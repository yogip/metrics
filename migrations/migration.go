package migrations

import (
	"context"
	"database/sql"
	"embed"
	"fmt"
	"syscall"
	"time"

	"metrics/internal/core/config"
	"metrics/internal/logger"
	"metrics/internal/retrier"

	"github.com/pressly/goose/v3"
	"go.uber.org/zap"
)

//go:embed *.sql
var embedMigrations embed.FS

func RunMigration(ctx context.Context, cfg *config.Config) error {
	logger.Log.Debug("Run migrations")

	ret := &retrier.Retrier{
		Strategy: retrier.Backoff(
			3,             // max attempts
			1*time.Second, // initial delay
			3,             // multiplier
			5*time.Second, // max delay
		),
		OnRetry: func(ctx context.Context, n int, err error) {
			logger.Log.Debug(fmt.Sprintf("RunMigration retry #%d: %v", n, err))
		},
	}

	db, err := sql.Open("pgx", cfg.Storage.DatabaseDSN)
	if err != nil {
		return fmt.Errorf("failed to initialize Database: %w", err)
	}

	goose.SetBaseFS(embedMigrations)

	if err := goose.SetDialect("postgres"); err != nil {
		return fmt.Errorf("failed to set dialect: %w", err)
	}

	fun := func() error {
		if err := goose.Up(db, "."); err != nil {
			return fmt.Errorf("failed to migrate: %w", err)
		}
		return nil
	}

	if err := ret.Do(ctx, fun, syscall.ECONNREFUSED); err != nil {
		logger.Log.Error("sending metric error", zap.String("error", err.Error()))
		return err
	}

	logger.Log.Debug("Migrations done")
	return nil
}
