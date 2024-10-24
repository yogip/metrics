package transport

import (
	"context"
	"crypto/rsa"
	"fmt"
	"log"

	"metrics/internal/core/model"
	pb "metrics/internal/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type GRPCClient struct {
	conn        *grpc.ClientConn
	client      pb.MetricsClient
	signHashKey string
	pubKey      *rsa.PublicKey
}

func NewGRPCClient(serverHost string, signHashKey string, pubKey *rsa.PublicKey) *GRPCClient {
	conn, err := grpc.NewClient(serverHost, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}

	client := pb.NewMetricsClient(conn)

	return &GRPCClient{
		conn:        conn,
		client:      client,
		signHashKey: signHashKey,
		pubKey:      pubKey,
	}
}

func (c *GRPCClient) Close() {
	c.conn.Close()
}

func (c *GRPCClient) SendMetrics(ctx context.Context, data []model.MetricsV2) error {
	metrics := make([]*pb.Metric, len(data))
	for _, m := range data {
		var pbM pb.Metric
		switch m.MType {
		case model.CounterType:
			pbM = pb.Metric{
				Name:      m.ID,
				Type:      pb.MetricType_COUNTER,
				Increment: *m.Delta,
			}
		case model.GaugeType:
			pbM = pb.Metric{
				Name:  m.ID,
				Type:  pb.MetricType_GAUGE,
				Value: *m.Value,
			}
		default:
			return fmt.Errorf("unsupported metric type: %s", m.MType)
		}
		metrics = append(metrics, &pbM)
	}
	payload := pb.BatchUpdateRequest{Metric: metrics}
	_, err := c.client.BatchUpdateMetric(ctx, &payload)
	if err != nil {
		return fmt.Errorf("grc batch update error: %w", err)
	}

	return nil
}
