package memory

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"metrics/internal/core/config"
	"metrics/internal/core/model"
	"metrics/internal/logger"
	"os"
	"sync"
	"time"

	"go.uber.org/zap"
)

type Store struct {
	mux     *sync.RWMutex
	quit    chan bool
	config  *config.StorageConfig
	gauge   map[string]*model.Gauge
	counter map[string]*model.Counter
}

func NewStore(cfg *config.StorageConfig) (*Store, error) {
	store := &Store{
		mux:     &sync.RWMutex{},
		quit:    make(chan bool),
		config:  cfg,
		gauge:   make(map[string]*model.Gauge),
		counter: make(map[string]*model.Counter),
	}

	if cfg.Restore && cfg.FileStoragePath != "" {
		if err := store.loadDump(); err != nil {
			return nil, err
		}
	}
	if cfg.StoreIntreval > 0 && cfg.FileStoragePath != "" {
		go store.dumpPeriodicly()
	}
	logger.Log.Info(
		"Memory Store initialized",
		zap.String("FileStoragePath", cfg.FileStoragePath),
		zap.Int64("StoreIntreval", cfg.StoreIntreval),
		zap.Bool("Restore", cfg.Restore),
	)
	return store, nil
}

func (s *Store) Close() {
	logger.Log.Debug("Send close event to chanel")
	close(s.quit)
}

func (s *Store) GetGauge(_ context.Context, req *model.MetricsV2) (*model.Gauge, error) {
	s.mux.RLock()
	defer s.mux.RUnlock()

	res, ok := s.gauge[req.ID]
	if !ok {
		return nil, nil
	}
	return res, nil
}

func (s *Store) SetGauge(_ context.Context, gauge *model.Gauge) error {
	s.mux.Lock()
	defer s.mux.Unlock()

	s.gauge[gauge.Name] = gauge
	if s.config.StoreIntreval == 0 && s.config.FileStoragePath != "" {
		s.saveDump()
	}
	return nil
}

func (s *Store) GetCounter(_ context.Context, req *model.MetricsV2) (*model.Counter, error) {
	s.mux.RLock()
	defer s.mux.RUnlock()

	res, ok := s.counter[req.ID]
	if !ok {
		return nil, nil
	}
	return res, nil
}

func (s *Store) SetCounter(_ context.Context, counter *model.Counter) error {
	s.mux.Lock()
	defer s.mux.Unlock()

	s.counter[counter.Name] = counter
	if s.config.StoreIntreval == 0 && s.config.FileStoragePath != "" {
		s.saveDump()
	}
	return nil
}

func (s *Store) BatchUpsertMetrics(ctx context.Context, metrics []*model.MetricsV2) ([]*model.MetricsV2, error) {
	results := make([]*model.MetricsV2, 0, len(metrics))
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
				return nil, errors.New("incorrect value")
			}
			gauge, err := s.GetGauge(ctx, m)
			if err != nil {
				return nil, fmt.Errorf("failed to fetch %s from the store: %w", m.MType.String(), err)
			}
			if gauge == nil {
				gauge = model.NewGauge(m.ID)
			}
			gauge.Set(*m.Value)
			if err = s.SetGauge(ctx, gauge); err != nil {
				return nil, fmt.Errorf("failed to save gauge to store: %w", err)
			}
			results = append(results, &model.MetricsV2{ID: m.ID, MType: m.MType, Value: &gauge.Value})
		case model.CounterType:
			if m.Delta == nil {
				return nil, errors.New("incorrect value")
			}
			counter, err := s.GetCounter(ctx, m)
			if err != nil {
				return nil, fmt.Errorf("failed to fetch %s from the store: %w", m.MType.String(), err)
			}
			if counter == nil {
				counter = model.NewCounter(m.ID)
			}
			counter.Increment(*m.Delta)
			if err = s.SetCounter(ctx, counter); err != nil {
				return nil, fmt.Errorf("failed to save gauge to store: %w", err)
			}
			results = append(results, &model.MetricsV2{ID: m.ID, MType: m.MType, Delta: &counter.Value})
		default:
			return nil, fmt.Errorf("unknown metric type: %s", m.MType.String())
		}
	}
	return results, nil
}

func (s *Store) ListGauge(_ context.Context) ([]*model.Gauge, error) {
	s.mux.RLock()
	defer s.mux.RUnlock()

	res := make([]*model.Gauge, 0, len(s.gauge))
	for _, v := range s.gauge {
		res = append(res, v)
	}
	return res, nil
}

func (s *Store) ListCounter(_ context.Context) ([]*model.Counter, error) {
	s.mux.RLock()
	defer s.mux.RUnlock()

	res := make([]*model.Counter, 0, len(s.counter))
	for _, v := range s.counter {
		res = append(res, v)
	}
	return res, nil
}

func (s *Store) saveDump() error {
	logger.Log.Debug("Dump DB to file", zap.String("path", s.config.FileStoragePath))
	dump := struct {
		Gauge   map[string]*model.Gauge   `json:"gauge"`
		Counter map[string]*model.Counter `json:"counter"`
	}{
		Gauge:   s.gauge,
		Counter: s.counter,
	}

	data, err := json.MarshalIndent(dump, "", " ")
	if err != nil {
		logger.Log.Error("Dump DB to json error", zap.Error(err))
		return err
	}
	err = os.WriteFile(s.config.FileStoragePath, data, 0666)
	if err != nil {
		logger.Log.Error("Dump DB to file error", zap.Error(err))
		return err
	}
	return nil
}

func (s *Store) loadDump() error {
	s.mux.Lock()
	defer s.mux.Unlock()
	logger.Log.Info("Load DB dump", zap.String("path", s.config.FileStoragePath))
	dump := struct {
		Gauge   map[string]*model.Gauge   `json:"gauge"`
		Counter map[string]*model.Counter `json:"counter"`
	}{
		Gauge:   s.gauge,
		Counter: s.counter,
	}

	file, err := os.OpenFile(s.config.FileStoragePath, os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		logger.Log.Error("Load Dump DB from file error", zap.Error(err))
		return err
	}
	logger.Log.Debug("File len", zap.Int("len", len(data)))
	if len(data) == 0 {
		return nil
	}

	err = json.Unmarshal(data, &dump)
	if err != nil {
		logger.Log.Error("Load Dump DB from json error", zap.Error(err))
		return err
	}
	logger.Log.Debug(fmt.Sprintf("Data %v", dump))

	s.gauge = dump.Gauge
	s.counter = dump.Counter
	return nil
}

func (s *Store) dumpPeriodicly() {
	ticker := time.NewTicker(time.Duration(s.config.StoreIntreval) * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-s.quit:
			logger.Log.Info("Close dump DB cilcle")
			return
		case <-ticker.C:
			s.mux.Lock()
			if err := s.saveDump(); err != nil {
				logger.Log.Error("Store dumping Error", zap.Error(err))
			}
			s.mux.Unlock()
		}
	}
}

func (s *Store) Ping(_ context.Context) error {
	return errors.New("memory store not supported ping")
}
