package service

import (
	"errors"
	"fmt"
	"strconv"

	"metrics/internal/core/model"
)

type Store interface {
	GetGauge(req *model.MetricRequest) (*model.Gauge, error)
	SetGauge(req *model.MetricRequest, gauge *model.Gauge) error
	ListGauge() ([]*model.Gauge, error)
	GetCounter(req *model.MetricRequest) (*model.Counter, error)
	SetCounter(req *model.MetricRequest, counter *model.Counter) error
	ListCounter() ([]*model.Counter, error)
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

func (m *MetricService) ListMetrics() (*model.ListMetricResponse, error) {
	result := model.ListMetricResponse{Metrics: []*model.MetricResponse{}}
	gagues, err := m.store.ListGauge()
	if err != nil {
		return nil, fmt.Errorf("failed to list gauges: %w", err)
	}
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

	counters, err := m.store.ListCounter()
	if err != nil {
		return nil, fmt.Errorf("failed to list gauges: %w", err)
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

func (m *MetricService) GetCounter(req *model.MetricsV2) (*model.MetricsV2, error) {
	storeReq := &model.MetricRequest{Name: req.ID, Type: req.MType}

	counter, err := m.store.GetCounter(storeReq)
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

func (m *MetricService) GetGauge(req *model.MetricsV2) (*model.MetricsV2, error) {
	storeReq := &model.MetricRequest{Name: req.ID, Type: req.MType}

	gauge, err := m.store.GetGauge(storeReq)
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

func (m *MetricService) GetMetric(req *model.MetricsV2) (*model.MetricsV2, error) {
	switch req.MType {
	case model.CounterType:
		return m.GetCounter(req)
	case model.GaugeType:
		return m.GetGauge(req)
	default:
		return nil, fmt.Errorf("unknown metric type: %s", req.MType)
	}
}

func (m *MetricService) UpsertGaugeValue(req *model.MetricsV2) (*model.MetricsV2, error) {
	storeReq := &model.MetricRequest{Name: req.ID, Type: req.MType}

	gauge, err := m.store.GetGauge(storeReq)
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
	err = m.store.SetGauge(storeReq, gauge)
	if err != nil {
		return &model.MetricsV2{}, fmt.Errorf("failed to save gauge to store: %w", err)
	}

	return &model.MetricsV2{
		ID:    req.ID,
		MType: model.GaugeType,
		Value: &gauge.Value,
	}, nil
}

func (m *MetricService) UpsertCounterValue(req *model.MetricsV2) (*model.MetricsV2, error) {
	storeReq := &model.MetricRequest{Name: req.ID, Type: req.MType}

	counter, err := m.store.GetCounter(storeReq)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch counter from the store: %w", err)
	}
	if counter == nil {
		counter = model.NewCounter(req.ID)
	}

	if req.Delta == nil {
		return nil, errors.New("incorrect value")
	}
	err = counter.Incremet(*req.Delta)
	if err != nil {
		return nil, fmt.Errorf("failed to increment metric (%s): %w", req.ID, err)
	}

	err = m.store.SetCounter(storeReq, counter)
	if err != nil {
		return nil, fmt.Errorf("failed to save %s value: %w", req.MType.String(), err)
	}

	return &model.MetricsV2{
		ID:    req.ID,
		MType: model.CounterType,
		Delta: &counter.Value,
	}, nil
}

func (m *MetricService) UpsertMetricValue(req *model.MetricsV2) (*model.MetricsV2, error) {
	switch req.MType {
	case model.CounterType:
		return m.UpsertCounterValue(req)
	case model.GaugeType:
		return m.UpsertGaugeValue(req)
	default:
		return nil, fmt.Errorf("unknown metric type: %s", req.MType.String())
	}
}

func (m *MetricService) BuildMetricRequest(req *model.MetricUpdateRequest, mustParseValue bool) (*model.MetricsV2, error) {
	reqV2 := &model.MetricsV2{
		ID:    req.Name,
		MType: req.Type,
	}

	switch req.Type {
	case model.CounterType:
		v, err := strconv.ParseInt(req.Value, 10, 64)
		if err != nil && mustParseValue {
			return nil, fmt.Errorf("failed to parse counter value: %w", err)
		}
		if err != nil {
			break
		}
		reqV2.Delta = &v
	case model.GaugeType:
		v, err := strconv.ParseFloat(req.Value, 64)
		if err != nil && mustParseValue {
			return nil, fmt.Errorf("failed to parse gauge value: %w", err)
		}
		if err != nil {
			break
		}
		reqV2.Value = &v
	default:
		return nil, fmt.Errorf("unknown metric type: %s", req.Type.String())
	}

	return reqV2, nil
}
