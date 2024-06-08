package migrations

import (
	"database/sql"
	"embed"
	"fmt"
	"metrics/internal/core/config"
	"metrics/internal/logger"

	"github.com/pressly/goose/v3"
)

//go:embed *.sql
var embedMigrations embed.FS

func RunMigration(cfg *config.Config) error {
	logger.Log.Debug("Run migrations")
	db, err := sql.Open("pgx", cfg.Storage.DatabaseDSN)
	if err != nil {
		return fmt.Errorf("failed to initialize Database: %w", err)
	}

	goose.SetBaseFS(embedMigrations)

	if err := goose.SetDialect("postgres"); err != nil {
		return fmt.Errorf("failed to set dialect: %w", err)
	}

	if err := goose.Up(db, "."); err != nil {
		return fmt.Errorf("failed to migrate: %w", err)
	}

	logger.Log.Debug("Migrations done")
	return nil
}
