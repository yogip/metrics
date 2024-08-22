package service

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"metrics/internal/core/model"
)

type Store interface {
	BatchUpsertMetrics(ctx context.Context, metrics []*model.MetricsV2) ([]*model.MetricsV2, error)
	GetGauge(ctx context.Context, req *model.MetricsV2) (*model.Gauge, error)
	SetGauge(ctx context.Context, gauge *model.Gauge) error
	ListGauge(ctx context.Context) ([]*model.Gauge, error)
	GetCounter(ctx context.Context, req *model.MetricsV2) (*model.Counter, error)
	SetCounter(ctx context.Context, counter *model.Counter) error
	ListCounter(ctx context.Context) ([]*model.Counter, error)
}

type Metric interface {
	StringValue() string
	Type() model.MetricType
}

type MetricService struct {
	store Store
}

func NewMetricService(store Store) *MetricService {
	return &MetricService{
		store: store,
	}
}

func (m *MetricService) ListMetrics(ctx context.Context) (*model.ListMetricResponse, error) {
	gagues, err := m.store.ListGauge(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list gauges: %w", err)
	}

	counters, err := m.store.ListCounter(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list gauges: %w", err)
	}

	result := model.ListMetricResponse{Metrics: make([]*model.MetricResponse, 0, len(gagues)+len(counters))}

	for _, gauge := range gagues {
		result.Metrics = append(
			result.Metrics,
			&model.MetricResponse{
				Name:  gauge.Name,
				Type:  gauge.Type(),
				Value: gauge.StringValue(),
			},
		)
	}

	for _, counter := range counters {
		result.Metrics = append(
			result.Metrics,
			&model.MetricResponse{
				Name:  counter.Name,
				Type:  counter.Type(),
				Value: counter.StringValue(),
			},
		)
	}

	return &result, nil
}

func (m *MetricService) GetCounter(ctx context.Context, req *model.MetricsV2) (*model.MetricsV2, error) {
	counter, err := m.store.GetCounter(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch counter from the store: %w", err)
	}
	if counter == nil {
		return nil, nil
	}
	return &model.MetricsV2{
		ID:    req.ID,
		MType: model.CounterType,
		Delta: &counter.Value,
	}, nil
}

func (m *MetricService) GetGauge(ctx context.Context, req *model.MetricsV2) (*model.MetricsV2, error) {
	gauge, err := m.store.GetGauge(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch cougauge from the store: %w", err)
	}
	if gauge == nil {
		return nil, nil
	}
	return &model.MetricsV2{
		ID:    req.ID,
		MType: model.GaugeType,
		Value: &gauge.Value,
	}, nil
}

func (m *MetricService) GetMetric(ctx context.Context, req *model.MetricsV2) (*model.MetricsV2, error) {
	switch req.MType {
	case model.CounterType:
		return m.GetCounter(ctx, req)
	case model.GaugeType:
		return m.GetGauge(ctx, req)
	default:
		return nil, fmt.Errorf("unknown metric type: %s", req.MType)
	}
}

func (m *MetricService) upsertGaugeValue(ctx context.Context, req *model.MetricsV2) (*model.MetricsV2, error) {
	gauge, err := m.store.GetGauge(ctx, req)
	if err != nil {
		return &model.MetricsV2{}, fmt.Errorf("failed to fetch %s from the store: %w", req.MType.String(), err)
	}
	if gauge == nil {
		gauge = model.NewGauge(req.ID)
	}

	if req.Value == nil {
		return nil, errors.New("incorrect value")
	}
	gauge.Set(*req.Value)
	err = m.store.SetGauge(ctx, gauge)
	if err != nil {
		return &model.MetricsV2{}, fmt.Errorf("failed to save gauge to store: %w", err)
	}

	return &model.MetricsV2{
		ID:    req.ID,
		MType: model.GaugeType,
		Value: &gauge.Value,
	}, nil
}

func (m *MetricService) upsertCounterValue(ctx context.Context, req *model.MetricsV2) (*model.MetricsV2, error) {
	counter, err := m.store.GetCounter(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch counter from the store: %w", err)
	}
	if counter == nil {
		counter = model.NewCounter(req.ID)
	}

	if req.Delta == nil {
		return nil, errors.New("incorrect value")
	}
	err = counter.Increment(*req.Delta)
	if err != nil {
		return nil, fmt.Errorf("failed to increment metric (%s): %w", req.ID, err)
	}

	err = m.store.SetCounter(ctx, counter)
	if err != nil {
		return nil, fmt.Errorf("failed to save counter value: %w", err)
	}

	return &model.MetricsV2{
		ID:    req.ID,
		MType: model.CounterType,
		Delta: &counter.Value,
	}, nil
}

func (m *MetricService) UpsertMetricValue(ctx context.Context, req *model.MetricsV2) (*model.MetricsV2, error) {
	switch req.MType {
	case model.CounterType:
		return m.upsertCounterValue(ctx, req)
	case model.GaugeType:
		return m.upsertGaugeValue(ctx, req)
	default:
		return nil, fmt.Errorf("unknown metric type: %s", req.MType.String())
	}
}

func (m *MetricService) BatchUpsertMetricValue(ctx context.Context, batch []*model.MetricsV2) ([]*model.MetricsV2, error) {
	return m.store.BatchUpsertMetrics(ctx, batch)
}

func (m *MetricService) BuildMetricRequest(
	name string, mType model.MetricType, value string, mustParseValue bool,
) (*model.MetricsV2, error) {
	reqV2 := &model.MetricsV2{
		ID:    name,
		MType: mType,
	}
	if value == "" && !mustParseValue {
		return reqV2, nil
	}

	switch mType {
	case model.CounterType:
		v, err := strconv.ParseInt(value, 10, 64)
		if err != nil && mustParseValue {
			return nil, fmt.Errorf("failed to parse counter value: %w", err)
		}
		if err != nil {
			break
		}
		reqV2.Delta = &v
	case model.GaugeType:
		v, err := strconv.ParseFloat(value, 64)
		if err != nil && mustParseValue {
			return nil, fmt.Errorf("failed to parse gauge value: %w", err)
		}
		if err != nil {
			break
		}
		reqV2.Value = &v
	default:
		return nil, fmt.Errorf("unknown metric type: %s", mType.String())
	}

	return reqV2, nil
}
