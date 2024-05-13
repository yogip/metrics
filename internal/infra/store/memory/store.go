package memory

import (
	"metrics/internal/core/model"
	"sync"
)

type Store struct {
	muxGauge   *sync.RWMutex
	muxCounter *sync.RWMutex
	gauge      map[string]*model.Gauge
	counter    map[string]*model.Counter
}

func NewStore() *Store {
	return &Store{
		muxGauge:   &sync.RWMutex{},
		muxCounter: &sync.RWMutex{},
		gauge:      make(map[string]*model.Gauge),
		counter:    make(map[string]*model.Counter),
	}
}

func (s *Store) GetGauge(req *model.MetricRequest) (*model.Gauge, error) {
	s.muxGauge.RLock()
	defer s.muxGauge.RUnlock()

	res, ok := s.gauge[req.ID()]
	if !ok {
		return nil, nil
	}
	return res, nil
}

func (s *Store) SetGauge(req *model.MetricRequest, gauge *model.Gauge) error {
	s.muxGauge.Lock()
	defer s.muxGauge.Unlock()

	s.gauge[req.ID()] = gauge
	return nil
}

func (s *Store) GetCounter(req *model.MetricRequest) (*model.Counter, error) {
	s.muxCounter.RLock()
	defer s.muxCounter.RUnlock()

	res, ok := s.counter[req.ID()]
	if !ok {
		return nil, nil
	}
	return res, nil
}

func (s *Store) SetCounter(req *model.MetricRequest, counter *model.Counter) error {
	s.muxCounter.Lock()
	defer s.muxCounter.Unlock()

	s.counter[req.ID()] = counter
	return nil
}

func (s *Store) ListGauge() ([]*model.Gauge, error) {
	s.muxGauge.RLock()
	defer s.muxGauge.RUnlock()

	var res []*model.Gauge
	for _, v := range s.gauge {
		res = append(res, v)
	}
	return res, nil
}

func (s *Store) ListCounter() ([]*model.Counter, error) {
	s.muxCounter.RLock()
	defer s.muxCounter.RUnlock()

	var res []*model.Counter
	for _, v := range s.counter {
		res = append(res, v)
	}
	return res, nil
}
