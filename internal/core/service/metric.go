package service

import (
	"fmt"

	"metrics/internal/core/model"
)

type Store interface {
	GetGauge(req *model.MetricRequest) (*model.Gauge, error)
	SetGauge(req *model.MetricRequest, gauge *model.Gauge) error
	GetCounter(req *model.MetricRequest) (*model.Counter, error)
	SetCounter(req *model.MetricRequest, counter *model.Counter) error
}

type Metric interface {
	StringValue() string
	ParseString(value string) error
}

type MetricService struct {
	store Store
}

func NewMetricService(store Store) *MetricService {
	return &MetricService{
		store: store,
	}
}

func (m *MetricService) GetMetric(req *model.MetricRequest) (*model.MetricResponse, error) {
	var metric Metric
	var err error

	switch req.Type {
	case model.CounterType:
		metric, err = m.store.GetCounter(req)
	case model.GaugeType:
		metric, err = m.store.GetGauge(req)
	default:
		return nil, fmt.Errorf("unknown metric type: %s", req.Type)
	}

	if err != nil {
		return &model.MetricResponse{}, fmt.Errorf("failed to fetch %s from the store: %w", req.Type, err)
	}
	switch req.Type {
	case model.CounterType:
		if counter, ok := metric.(*model.Counter); ok && counter == nil {
			return &model.MetricResponse{}, fmt.Errorf("counter %s not found", req.Name)
		}
	case model.GaugeType:
		if gauge, ok := metric.(*model.Gauge); ok && gauge == nil {
			return &model.MetricResponse{}, fmt.Errorf("gauge %s not found", req.Name)
		}
	}

	return &model.MetricResponse{
		Name:  req.Name,
		Type:  model.CounterType,
		Value: metric.StringValue(),
	}, err
}

func (m *MetricService) SetMetricValue(req *model.MetricUpdateRequest) (*model.MetricResponse, error) {
	var metric Metric
	var err error

	getReq := &model.MetricRequest{Name: req.Name, Type: req.Type}

	switch req.Type {
	case model.CounterType:
		metric, err = m.store.GetCounter(getReq)
	case model.GaugeType:
		metric, err = m.store.GetGauge(getReq)
	default:
		return nil, fmt.Errorf("unknown metric type: %s", req.Type)
	}

	if err != nil {
		return &model.MetricResponse{}, fmt.Errorf("failed to fetch %s from the store: %w", req.Type, err)
	}
	switch req.Type {
	case model.CounterType:
		if counter, ok := metric.(*model.Counter); ok && counter == nil {
			metric = model.NewCounter(req.Name)
		}
	case model.GaugeType:
		if gauge, ok := metric.(*model.Gauge); ok && gauge == nil {
			metric = model.NewGauge(req.Name)
		}
	}

	err = metric.ParseString(req.Value)
	if err != nil {
		return &model.MetricResponse{}, fmt.Errorf("failed to parse %s value: %w", req.Type, err)
	}

	switch req.Type {
	case model.CounterType:
		err = m.store.SetCounter(getReq, metric.(*model.Counter))
	case model.GaugeType:
		err = m.store.SetGauge(getReq, metric.(*model.Gauge))
	}
	if err != nil {
		return &model.MetricResponse{}, fmt.Errorf("failed to save %s value: %w", req.Type, err)
	}

	return &model.MetricResponse{
		Name:  req.Name,
		Type:  model.CounterType,
		Value: metric.StringValue(),
	}, err
}
