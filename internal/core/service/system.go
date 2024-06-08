package service

import "context"

type Pinger interface {
	Ping(ctx context.Context) error
}

type SystemService struct {
	store Pinger
}

func NewSystemService(store Pinger) *SystemService {
	return &SystemService{store: store}
}

func (s *SystemService) Ping(ctx context.Context) error {
	return s.store.Ping(ctx)
}
