package db

import (
	"context"
	"database/sql"
	"fmt"

	"metrics/internal/core/model"

	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func TestPing(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.MonitorPingsOption(true))
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	mock.ExpectPing()

	store := newStore(db)

	ctx := context.Background()

	err = store.Ping(ctx)
	require.NoError(t, err)
}

func TestBatchUpsertMetrics(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	var value01 float64 = 1.01
	var value02 float64 = 2.01
	var delta01 int64 = 1
	var delta02 int64 = 10
	var delta03 int64 = 101
	batch := []*model.MetricsV2{
		{
			ID:    "gauge_01",
			MType: model.GaugeType,
			Value: &value01,
		},
		{
			ID:    "gauge_02",
			MType: model.GaugeType,
			Value: &value02,
		},
		{
			ID:    "counter_01",
			MType: model.CounterType,
			Delta: &delta01,
		},
		{
			ID:    "counter_01",
			MType: model.CounterType,
			Delta: &delta02,
		},
		{
			ID:    "counter_01",
			MType: model.CounterType,
			Delta: &delta03,
		},
	}

	var expDelta01 int64 = delta01
	var expDelta02 int64 = delta01 + delta02
	var expDelta03 int64 = delta01 + delta02 + delta03
	expected := []*model.MetricsV2{
		{
			ID:    "gauge_01",
			MType: model.GaugeType,
			Value: &value01,
		},
		{
			ID:    "gauge_02",
			MType: model.GaugeType,
			Value: &value02,
		},
		{
			ID:    "counter_01",
			MType: model.CounterType,
			Delta: &expDelta01,
		},
		{
			ID:    "counter_01",
			MType: model.CounterType,
			Delta: &expDelta02,
		},
		{
			ID:    "counter_01",
			MType: model.CounterType,
			Delta: &expDelta03,
		},
	}

	store := newStore(db)
	mock.ExpectBegin()

	var total int64
	for _, m := range batch {
		switch m.MType {
		case model.GaugeType:
			mock.ExpectQuery(
				`INSERT INTO gauge\(id, value\) values\(\$1, \$2\) ON conflict\(id\) 
					DO UPDATE SET value \= excluded.value
					RETURNING value`,
			).
				WithArgs(m.ID, m.Value).
				WillReturnRows(
					sqlmock.NewRows([]string{"current"}).AddRow(m.Value),
				)
		case model.CounterType:

			mock.ExpectQuery(
				`INSERT INTO counter\(id, value\) values\(\$1, \$2\) ON conflict\(id\) 
				DO UPDATE SET value \= counter.value \+ excluded.value 
				RETURNING value`,
			).
				WithArgs(m.ID, m.Delta).
				WillReturnRows(
					sqlmock.NewRows([]string{"current"}).AddRow(total + *m.Delta),
				)
			total += *m.Delta
		}
	}
	mock.ExpectCommit()

	ctx := context.Background()

	actual, err := store.BatchUpsertMetrics(ctx, batch)
	require.NoError(t, err)

	assert.Equal(t, expected, actual)
}

func TestRollbackBatchUpsertMetrics(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	var value01 float64 = 1.01
	var delta01 int64 = 1

	tests := [][]*model.MetricsV2{
		{
			{
				ID:    "gauge_01",
				MType: model.GaugeType,
				Value: &value01,
			},
			{
				ID:    "gauge_01_nil",
				MType: model.GaugeType,
				Value: nil,
			},
		},
		{
			{
				ID:    "counter_01",
				MType: model.CounterType,
				Delta: &delta01,
			},
			{
				ID:    "counter_01",
				MType: model.CounterType,
				Delta: nil,
			},
		},
	}

	for n, batch := range tests {
		t.Run(fmt.Sprintf("batch rollback - %d", n), func(t *testing.T) {
			mock.ExpectBegin()

			for _, m := range batch {
				switch m.MType {
				case model.GaugeType:
					if m.Value == nil {
						mock.ExpectRollback()
					} else {
						mock.ExpectQuery(
							`INSERT INTO gauge\(id, value\) values\(\$1, \$2\) ON conflict\(id\) 
								DO UPDATE SET value \= excluded.value
								RETURNING value`,
						).
							WithArgs(m.ID, m.Value).
							WillReturnRows(
								sqlmock.NewRows([]string{"current"}).AddRow(m.Value),
							)
					}
				case model.CounterType:
					if m.Delta == nil {
						mock.ExpectRollback()
					} else {
						mock.ExpectQuery(
							`INSERT INTO counter\(id, value\) values\(\$1, \$2\) ON conflict\(id\) 
							DO UPDATE SET value \= counter.value \+ excluded.value 
							RETURNING value`,
						).
							WithArgs(m.ID, m.Delta).
							WillReturnRows(
								sqlmock.NewRows([]string{"current"}).AddRow(*m.Delta),
							)
					}
				}
			}

			store := newStore(db)
			ctx := context.Background()

			actual, err := store.BatchUpsertMetrics(ctx, batch)
			require.Error(t, err)
			assert.Nil(t, actual)
		})
	}
}

func TestGetGauge(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	store := newStore(db)
	ctx := context.Background()

	tests := []struct {
		metric   *model.MetricsV2
		want     *model.Gauge
		testName string
	}{
		{
			testName: "get gauge success",
			metric: &model.MetricsV2{
				ID:    "gauge_01",
				MType: model.GaugeType,
			},
			want: &model.Gauge{
				Name:  "gauge_01",
				Value: 123.001,
			},
		},
		{
			testName: "get gauge not found",
			metric: &model.MetricsV2{
				ID:    "gauge_02",
				MType: model.GaugeType,
			},
			want: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			if tt.want != nil {
				mock.ExpectQuery(
					`SELECT id, value FROM gauge WHERE id=\$1`,
				).
					WithArgs(tt.metric.ID).
					WillReturnRows(
						sqlmock.NewRows([]string{"id", "value"}).AddRow(tt.want.Name, tt.want.Value),
					)
			} else {
				mock.ExpectQuery(
					`SELECT id, value FROM gauge WHERE id=\$1`,
				).
					WithArgs(tt.metric.ID).
					WillReturnError(sql.ErrNoRows)
			}

			actual, err := store.GetGauge(ctx, tt.metric)
			require.NoError(t, err)
			assert.Equal(t, tt.want, actual)

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestGetCounter(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	store := newStore(db)
	ctx := context.Background()

	tests := []struct {
		metric   *model.MetricsV2
		want     *model.Counter
		testName string
	}{
		{
			testName: "get counter success",
			metric: &model.MetricsV2{
				ID:    "counter_01",
				MType: model.CounterType,
			},
			want: &model.Counter{
				Name:  "counter_01",
				Value: 123,
			},
		},
		{
			testName: "get counter not found",
			metric: &model.MetricsV2{
				ID:    "counter_02",
				MType: model.CounterType,
			},
			want: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			if tt.want != nil {
				mock.ExpectQuery(
					`SELECT id, value FROM counter WHERE id=\$1`,
				).
					WithArgs(tt.metric.ID).
					WillReturnRows(
						sqlmock.NewRows([]string{"id", "value"}).AddRow(tt.want.Name, tt.want.Value),
					)
			} else {
				mock.ExpectQuery(
					`SELECT id, value FROM counter WHERE id=\$1`,
				).
					WithArgs(tt.metric.ID).
					WillReturnError(sql.ErrNoRows)
			}

			actual, err := store.GetCounter(ctx, tt.metric)
			require.NoError(t, err)
			assert.Equal(t, tt.want, actual)

			require.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestSetGauge(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	store := newStore(db)
	ctx := context.Background()

	tests := []struct {
		metric   *model.Gauge
		testName string
	}{
		{
			testName: "set gauge success",
			metric: &model.Gauge{
				Name:  "gauge_01",
				Value: 123.0,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			mock.ExpectExec(
				`INSERT INTO gauge\(id, value\) values\(\$1, \$2\) ON conflict\(id\) DO UPDATE SET value \= excluded\.value`,
			).
				WithArgs(tt.metric.Name, tt.metric.Value).
				WillReturnResult(sqlmock.NewResult(1, 1))

			err := store.SetGauge(ctx, tt.metric)
			require.NoError(t, err)
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestSetCounter(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	store := newStore(db)
	ctx := context.Background()

	tests := []struct {
		metric   *model.Counter
		testName string
	}{
		{
			testName: "set counter success",
			metric: &model.Counter{
				Name:  "counter_01",
				Value: 123,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			mock.ExpectExec(
				`INSERT INTO counter\(id, value\) values\(\$1, \$2\) ON conflict\(id\) DO UPDATE SET value \= excluded\.value`,
			).
				WithArgs(tt.metric.Name, tt.metric.Value).
				WillReturnResult(sqlmock.NewResult(1, 1))

			err := store.SetCounter(ctx, tt.metric)
			require.NoError(t, err)
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestListGauge(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	store := newStore(db)
	ctx := context.Background()

	metrics := []*model.Gauge{
		{
			Name:  "gauge_01",
			Value: 123.0,
		},
		{
			Name:  "gauge_02",
			Value: 1.01,
		},
	}

	rows := sqlmock.NewRows([]string{"id", "value"}).
		AddRow("gauge_01", 123.0).
		AddRow("gauge_02", 1.01)

	mock.ExpectQuery(`SELECT id, value FROM gauge`).WillReturnRows(rows)

	actual, err := store.ListGauge(ctx)
	require.NoError(t, err)
	require.NoError(t, mock.ExpectationsWereMet())
	assert.Equal(t, metrics, actual)
}

func TestListCounter(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	store := newStore(db)
	ctx := context.Background()

	metrics := []*model.Counter{
		{
			Name:  "counter_01",
			Value: 123,
		},
		{
			Name:  "counter_02",
			Value: 1,
		},
	}

	rows := sqlmock.NewRows([]string{"id", "value"}).
		AddRow("counter_01", 123).
		AddRow("counter_02", 1)

	mock.ExpectQuery(`SELECT id, value FROM counter`).WillReturnRows(rows)

	actual, err := store.ListCounter(ctx)
	require.NoError(t, err)
	require.NoError(t, mock.ExpectationsWereMet())
	assert.Equal(t, metrics, actual)
}
