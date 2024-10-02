package db

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"
	"fmt"
	"io"
	"sync"
	"syscall"
	"time"

	"metrics/internal/core/config"
	"metrics/internal/core/model"
	"metrics/internal/logger"
	"metrics/internal/retrier"

	"github.com/jackc/pgx"
	"go.uber.org/zap"
)

type Store struct {
	db      *sql.DB
	retrier *retrier.Retrier
}

var recoverableErrors = []error{
	syscall.ECONNREFUSED,
	pgx.ErrDeadConn,
	sql.ErrConnDone,
	pgx.ErrConnBusy,
	driver.ErrBadConn,
	io.EOF,
}

func newStore(db *sql.DB) *Store {
	ret := &retrier.Retrier{
		Strategy: retrier.Backoff(
			3,             // max attempts
			1*time.Second, // initial delay
			3,             // multiplier
			5*time.Second, // max delay
		),
		OnRetry: func(ctx context.Context, n int, err error) {
			logger.Log.Debug(fmt.Sprintf("Retrying DB. retry #%d: %v", n, err))
		},
	}

	store := &Store{db: db, retrier: ret}

	logger.Log.Info("DB Store initialized")
	return store
}

func NewStore(ctx context.Context, wg *sync.WaitGroup, cfg *config.StorageConfig) (*Store, error) {
	db, err := sql.Open("pgx", cfg.DatabaseDSN)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize Database: %w", err)
	}

	return newStore(db), nil
}

func (s *Store) Close() {
	s.db.Close()
}

func (s *Store) Ping(ctx context.Context) error {
	fun := func() error {
		return s.db.PingContext(ctx)
	}
	return s.retrier.Do(ctx, fun, recoverableErrors...)
}

func (s *Store) BatchUpsertMetrics(ctx context.Context, metrics []*model.MetricsV2) ([]*model.MetricsV2, error) {
	results := []*model.MetricsV2{}
	fun := func() error {
		var err error
		results, err = s.doBatchUpsertMetrics(ctx, metrics)
		return err
	}

	err := s.retrier.Do(ctx, fun, recoverableErrors...)
	return results, err
}

func (s *Store) doBatchUpsertMetrics(ctx context.Context, metrics []*model.MetricsV2) ([]*model.MetricsV2, error) {
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
	gauge := &model.Gauge{}

	fun := func() error {
		var err error
		row := s.db.QueryRowContext(ctx, "SELECT id, value FROM gauge WHERE id=$1", req.ID)
		err = row.Scan(&gauge.Name, &gauge.Value)
		if errors.Is(err, sql.ErrNoRows) {
			gauge = nil
			return nil
		}
		return err
	}
	err := s.retrier.Do(ctx, fun, recoverableErrors...)
	if err != nil {
		return nil, fmt.Errorf("error reading gauge: %w", err)
	}
	return gauge, nil
}

func (s *Store) SetGauge(ctx context.Context, gauge *model.Gauge) error {
	fun := func() error {
		result, err := s.db.ExecContext(
			ctx,
			"INSERT INTO gauge(id, value) values($1, $2) ON conflict(id) DO UPDATE SET value = excluded.value",
			gauge.Name, gauge.Value,
		)
		if err != nil {
			return fmt.Errorf("error ExecContext for gauge: %w", err)
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
	err := s.retrier.Do(ctx, fun, recoverableErrors...)
	if err != nil {
		return fmt.Errorf("error writing gauge: %w", err)
	}
	return nil
}

func (s *Store) ListGauge(ctx context.Context) ([]*model.Gauge, error) {
	results := []*model.Gauge{}
	fun := func() error {
		var err error
		results, err = s.doListGauge(ctx)
		return err
	}

	err := s.retrier.Do(ctx, fun, recoverableErrors...)
	return results, err
}

func (s *Store) doListGauge(ctx context.Context) ([]*model.Gauge, error) {
	gauges := make([]*model.Gauge, 0, 10)

	rows, err := s.db.QueryContext(ctx, "SELECT id, value FROM gauge")
	if err != nil {
		return nil, fmt.Errorf("error reading gauge: %w", err)
	}
	defer rows.Close()
	rows.Err()

	for rows.Next() {
		var g model.Gauge
		err = rows.Scan(&g.Name, &g.Value)
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

	fun := func() error {
		var err error
		row := s.db.QueryRowContext(ctx, "SELECT id, value FROM counter WHERE id=$1", req.ID)
		err = row.Scan(&counter.Name, &counter.Value)
		if errors.Is(err, sql.ErrNoRows) {
			counter = nil
			return nil
		}
		return err
	}
	err := s.retrier.Do(ctx, fun, recoverableErrors...)
	if err != nil {
		return nil, fmt.Errorf("error reading gauge: %w", err)
	}
	return counter, nil
}

func (s *Store) SetCounter(ctx context.Context, counter *model.Counter) error {
	fun := func() error {
		result, err := s.db.ExecContext(
			ctx,
			"INSERT INTO counter(id, value) values($1, $2) ON conflict(id) DO UPDATE SET value = excluded.value",
			counter.Name, counter.Value,
		)
		if err != nil {
			return fmt.Errorf("error ExecContext for counter: %w", err)
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
	err := s.retrier.Do(ctx, fun, recoverableErrors...)
	if err != nil {
		return fmt.Errorf("error writing counter: %w", err)
	}
	return nil
}

func (s *Store) ListCounter(ctx context.Context) ([]*model.Counter, error) {
	results := []*model.Counter{}
	fun := func() error {
		var err error
		results, err = s.doListCounter(ctx)
		return err
	}

	err := s.retrier.Do(ctx, fun, recoverableErrors...)
	return results, err
}

func (s *Store) doListCounter(ctx context.Context) ([]*model.Counter, error) {
	counters := make([]*model.Counter, 0, 10)

	rows, err := s.db.QueryContext(ctx, "SELECT id, value FROM counter")
	if err != nil {
		return nil, fmt.Errorf("error reading counters: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var c model.Counter
		err = rows.Scan(&c.Name, &c.Value)
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
