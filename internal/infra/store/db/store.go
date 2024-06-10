package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"metrics/internal/core/config"
	"metrics/internal/core/model"
	"metrics/internal/logger"

	"go.uber.org/zap"
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

func (s *Store) BatchUpsertMetrics(ctx context.Context, metrics []*model.MetricsV2) ([]*model.MetricsV2, error) {
	results := make([]*model.MetricsV2, 0, len(metrics))
	tx, err := s.db.Begin()
	if err != nil {
		return nil, fmt.Errorf("error starting transaction: %w", err)
	}
	for _, m := range metrics {
		logger.Log.Debug("BatchUpsertMetrics input",
			zap.String("metric", m.ID),
			zap.String("type", m.MType.String()),
			zap.Float64p("value", m.Value),
			zap.Int64p("delta", m.Delta),
		)
		switch m.MType {
		case model.GaugeType:
			if m.Value == nil {
				tx.Rollback()
				return nil, fmt.Errorf("gauge Value clould not be nil: %v", m)
			}
			res := model.MetricsV2{ID: m.ID, MType: m.MType}
			row := tx.QueryRowContext(
				ctx,
				`INSERT INTO gauge(id, value) values($1, $2) ON conflict(id) 
				 DO UPDATE SET value = excluded.value
				 RETURNING value`,
				m.ID, m.Value,
			)
			err = row.Scan(&res.Value)
			logger.Log.Debug("BatchUpsertMetrics returning",
				zap.String("metric", res.ID),
				zap.String("type", res.MType.String()),
				zap.Float64p("value", res.Value),
			)
			if err != nil {
				tx.Rollback()
				return nil, fmt.Errorf("error upserting gauge: %w", err)
			}
			results = append(results, &res)
		case model.CounterType:
			if m.Delta == nil {
				tx.Rollback()
				return nil, fmt.Errorf("counter Delta clould not be nil: %v", m)
			}
			res := model.MetricsV2{ID: m.ID, MType: m.MType}
			row := tx.QueryRowContext(
				ctx,
				`INSERT INTO counter(id, value) values($1, $2) ON conflict(id) 
				 DO UPDATE SET value = counter.value + excluded.value
				 RETURNING value`,
				m.ID, m.Delta,
			)
			err = row.Scan(&res.Delta)
			if err != nil {
				tx.Rollback()
				return nil, fmt.Errorf("error upserting counter: %w", err)
			}
			logger.Log.Debug("BatchUpsertMetrics returning",
				zap.String("metric", res.ID),
				zap.String("type", res.MType.String()),
				zap.Int64p("delta", res.Delta),
			)
			results = append(results, &res)
		default:
			tx.Rollback()
			return nil, fmt.Errorf("unknown metric type: %s", m.MType.String())
		}
	}
	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("batch transaction commit error: %w", err)
	}
	return results, nil
}

func (s *Store) GetGauge(ctx context.Context, req *model.MetricsV2) (*model.Gauge, error) {
	gauge := model.Gauge{}

	row := s.db.QueryRowContext(ctx, "SELECT id, value FROM gauge WHERE id=$1", req.ID)
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
		"INSERT INTO gauge(id, value) values($1, $2) ON conflict(id) DO UPDATE SET value = excluded.value",
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
	rows.Err()

	for rows.Next() {
		var g model.Gauge
		err := rows.Scan(&g.Name, &g.Value)
		if err != nil {
			return nil, fmt.Errorf("error reading gauge row: %w", err)
		}

		gauges = append(gauges, &g)
	}

	err = rows.Err()
	if err != nil {
		return nil, fmt.Errorf("error scaning gauges: %w", err)
	}

	return gauges, nil
}

func (s *Store) GetCounter(ctx context.Context, req *model.MetricsV2) (*model.Counter, error) {
	counter := &model.Counter{}

	row := s.db.QueryRowContext(ctx, "SELECT id, value FROM counter WHERE id=$1", req.ID)
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
		"INSERT INTO counter(id, value) values($1, $2) ON conflict(id) DO UPDATE SET value = excluded.value",
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
			return nil, fmt.Errorf("error reading counter row: %w", err)
		}

		counters = append(counters, &c)
	}

	err = rows.Err()
	if err != nil {
		return nil, fmt.Errorf("error scaning counters: %w", err)
	}

	return counters, nil
}
