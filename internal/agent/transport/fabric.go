package transport

import (
	"context"
	"crypto/rsa"
	"fmt"

	"metrics/internal/agent/config"
	"metrics/internal/core/model"
	"metrics/internal/logger"
)

type Transporter interface {
	SendMetrics(context.Context, []model.MetricsV2) error
	Close()
}

func NewClient(tType config.TransportType, serverHost string, signHashKey string, pubKey *rsa.PublicKey) (Transporter, error) {
	logger.Log.Info(fmt.Sprintf("Create %s transport for %s", tType, serverHost))
	switch tType {
	case config.HTTPTransportType:
		return NewHTTPClient(serverHost, signHashKey, pubKey), nil
	case config.GRPCTransportType:
		return NewGRPCClient(serverHost, signHashKey, pubKey), nil
	}
	return nil, fmt.Errorf("unknown transport type: %s", tType)
}
