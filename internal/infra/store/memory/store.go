package memory

import (
	"metrics/internal/core/model"
	"sync"
)

type Store struct {
	mux_gauge   *sync.RWMutex
	mux_counter *sync.RWMutex
	gauge       map[string]*model.Gauge
	counter     map[string]*model.Counter
}

func NewStore() *Store {
	return &Store{
		mux_gauge:   &sync.RWMutex{},
		mux_counter: &sync.RWMutex{},
		gauge:       make(map[string]*model.Gauge),
		counter:     make(map[string]*model.Counter),
	}
}

func (s *Store) GetGauge(req *model.MetricRequest) (*model.Gauge, error) {
	s.mux_gauge.RLock()
	defer s.mux_gauge.RUnlock()

	res, ok := s.gauge[req.ID()]
	if !ok {
		return nil, nil
	}
	return res, nil
}

func (s *Store) SetGauge(req *model.MetricRequest, gauge *model.Gauge) error {
	s.mux_gauge.Lock()
	defer s.mux_gauge.Unlock()

	s.gauge[req.ID()] = gauge
	return nil
}

func (s *Store) GetCounter(req *model.MetricRequest) (*model.Counter, error) {
	s.mux_counter.RLock()
	defer s.mux_counter.RUnlock()

	res, ok := s.counter[req.ID()]
	if !ok {
		return nil, nil
	}
	return res, nil
}

func (s *Store) SetCounter(req *model.MetricRequest, counter *model.Counter) error {
	s.mux_counter.Lock()
	defer s.mux_counter.Unlock()

	s.counter[req.ID()] = counter
	return nil
}
