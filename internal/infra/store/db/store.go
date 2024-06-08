package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"metrics/internal/core/config"
	"metrics/internal/core/model"
	"metrics/internal/logger"
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

func (s *Store) Ping(ctx context.Context) error {
	return s.db.PingContext(ctx)
}

func (s *Store) GetGauge(ctx context.Context, req *model.MetricRequest) (*model.Gauge, error) {
	gauge := model.Gauge{}

	row := s.db.QueryRowContext(ctx, "SELECT id, value FROM gauge WHERE id=$1", req.Name)
	err := row.Scan(&gauge.Name, &gauge.Value)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("error reading gauge: %w", err)
	}

	return &gauge, nil
}

func (s *Store) SetGauge(ctx context.Context, gauge *model.Gauge) error {
	result, err := s.db.ExecContext(
		ctx,
		"INSERT INTO gauge(id, value) values($1, $2) ON conflict(id) DO UPDATE SET value = $2",
		gauge.Name, gauge.Value,
	)
	if err != nil {
		return fmt.Errorf("error writing gauge: %w", err)
	}

	count, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error getting RowsAffected for gauge: %w", err)
	}
	if count != 1 {
		return fmt.Errorf("incorrect rows affected: %d", count)
	}

	return nil
}

func (s *Store) ListGauge(ctx context.Context) ([]*model.Gauge, error) {
	gauges := make([]*model.Gauge, 0, 10)

	rows, err := s.db.QueryContext(ctx, "SELECT id, value FROM gauge")
	if err != nil {
		return nil, fmt.Errorf("error reading gauge: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var g model.Gauge
		err := rows.Scan(&g.Name, &g.Value)
		if err != nil {
			return nil, err
		}

		gauges = append(gauges, &g)
	}

	return gauges, nil
}

func (s *Store) GetCounter(ctx context.Context, req *model.MetricRequest) (*model.Counter, error) {
	counter := &model.Counter{}

	row := s.db.QueryRowContext(ctx, "SELECT id, value FROM counter WHERE id=$1", req.Name)
	err := row.Scan(&counter.Name, &counter.Value)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("error reading gauge: %w", err)
	}

	return counter, nil
}

func (s *Store) SetCounter(ctx context.Context, counter *model.Counter) error {
	result, err := s.db.ExecContext(
		ctx,
		"INSERT INTO counter(id, value) values($1, $2) ON conflict(id) DO UPDATE SET value = $2",
		counter.Name, counter.Value,
	)
	if err != nil {
		return fmt.Errorf("error writing counter: %w", err)
	}

	count, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error getting RowsAffected for counter: %w", err)
	}
	if count != 1 {
		return fmt.Errorf("incorrect rows affected: %d", count)
	}

	return nil
}

func (s *Store) ListCounter(ctx context.Context) ([]*model.Counter, error) {
	counters := make([]*model.Counter, 0, 10)

	rows, err := s.db.QueryContext(ctx, "SELECT id, value FROM counter")
	if err != nil {
		return nil, fmt.Errorf("error reading counters: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var c model.Counter
		err := rows.Scan(&c.Name, &c.Value)
		if err != nil {
			return nil, err
		}

		counters = append(counters, &c)
	}

	return counters, nil
}
