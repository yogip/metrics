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

func (s *Store) GetGauge(ctx context.Context, req *model.MetricRequest) (*model.Gauge, error) {
	s.mux.RLock()
	defer s.mux.RUnlock()

	res, ok := s.gauge[req.Name]
	if !ok {
		return nil, nil
	}
	return res, nil
}

func (s *Store) SetGauge(ctx context.Context, gauge *model.Gauge) error {
	s.mux.Lock()
	defer s.mux.Unlock()

	s.gauge[gauge.Name] = gauge
	if s.config.StoreIntreval == 0 && s.config.FileStoragePath != "" {
		s.saveDump()
	}
	return nil
}

func (s *Store) GetCounter(ctx context.Context, req *model.MetricRequest) (*model.Counter, error) {
	s.mux.RLock()
	defer s.mux.RUnlock()

	res, ok := s.counter[req.Name]
	if !ok {
		return nil, nil
	}
	return res, nil
}

func (s *Store) SetCounter(ctx context.Context, counter *model.Counter) error {
	s.mux.Lock()
	defer s.mux.Unlock()

	s.counter[counter.Name] = counter
	if s.config.StoreIntreval == 0 && s.config.FileStoragePath != "" {
		s.saveDump()
	}
	return nil
}

func (s *Store) ListGauge(ctx context.Context) ([]*model.Gauge, error) {
	s.mux.RLock()
	defer s.mux.RUnlock()

	res := make([]*model.Gauge, 0, len(s.gauge))
	for _, v := range s.gauge {
		res = append(res, v)
	}
	return res, nil
}

func (s *Store) ListCounter(ctx context.Context) ([]*model.Counter, error) {
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

func (s *Store) Ping(ctx context.Context) error {
	return errors.New("memory store not supported ping")
}
