package grpc

import (
	"context"
	"crypto/rsa"
	"fmt"
	"net"
	"strings"

	"metrics/internal/core/config"
	"metrics/internal/core/model"
	"metrics/internal/core/service"
	"metrics/internal/logger"
	pb "metrics/internal/proto"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"
)

// MetricsServer реализует gRPC сервер метрик
type MetricsServer struct {
	pb.UnimplementedMetricsServer

	cfg           *config.Config
	metricService *service.MetricService
	systemService *service.SystemService

	srv *grpc.Server
}

// subnetInterceptor create Interceptor to checks client's subnet.
func subnetInterceptor(subnet *net.IPNet) func(context.Context, interface{}, *grpc.UnaryServerInfo, grpc.UnaryHandler) (interface{}, error) {
	return func(
		ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler,
	) (interface{}, error) {
		p, ok := peer.FromContext(ctx)
		if !ok {
			return nil, status.Error(codes.PermissionDenied, "Failed to get peer from context")
		}

		ipPort := strings.Split(p.Addr.String(), ":")
		ipStr := ipPort[0]
		clientIP := net.ParseIP(ipStr)
		logger.Log.Debug(fmt.Sprintf("Got request from Client with IP %s", clientIP))

		if clientIP == nil || !subnet.Contains(clientIP) {
			logger.Log.Warn(fmt.Sprintf("Client IP %s is not allowed", clientIP))
			return nil, status.Error(codes.PermissionDenied, "Access denied")
		}

		return handler(ctx, req)
	}
}

func NewMetricsServer(
	cfg *config.Config,
	metricService *service.MetricService,
	systemService *service.SystemService,
	privateKey *rsa.PrivateKey,
) *MetricsServer {
	s := grpc.NewServer()
	if cfg.TrustedSubnet != nil {
		s = grpc.NewServer(grpc.UnaryInterceptor(subnetInterceptor(cfg.TrustedSubnet)))
	}

	m := MetricsServer{
		cfg:           cfg,
		metricService: metricService,
		systemService: systemService,
		srv:           s,
	}
	pb.RegisterMetricsServer(s, &m)

	return &m
}

// Run MetricsServer server. It blocks until the server is stopped.
func (s *MetricsServer) Run(address string) error {
	logger.Log.Info("Run gRPC server", zap.String("Addres", address))
	listen, err := net.Listen("tcp", address)
	if err != nil {
		return fmt.Errorf("failed to listen gRPC adress at %s. error: %w", address, err)
	}

	return s.srv.Serve(listen)
}

// Shutdown MetricsServer server. It blocks until the server is stopped. Under the hood calls http.Server.Shutdown.
func (s *MetricsServer) Shutdown(ctx context.Context) error {
	logger.Log.Info("Starting gracefull shutdown of gRPC server")
	s.srv.GracefulStop()
	logger.Log.Info("gRPC server is down")
	return nil
}

// GetMetric реализует интерфейс чтения метрик.
func (s *MetricsServer) GetMetric(ctx context.Context, in *pb.GetRequest) (*pb.GetResponse, error) {
	logger.Log.Info("GetMetric request", zap.Any("request", in))
	metric, err := s.metricService.GetMetric(ctx, &model.MetricsV2{ID: in.Name, MType: model.MetricType(in.Type)})
	if err != nil {
		return nil, status.Errorf(codes.Unknown, "Reading metric error: %s", err)
	}
	if metric != nil {
		return nil, status.Errorf(codes.NotFound, "Not found")
	}

	response := pb.GetResponse{
		Metric: MetricsV2ToPbMetric(metric),
	}
	return &response, nil
}

// BatchUpdateMetric реализует интерфейс батч обновление метрик.
func (s *MetricsServer) BatchUpdateMetric(ctx context.Context, in *pb.BatchUpdateRequest) (*pb.BatchUpdateResponse, error) {
	logger.Log.Info("BatchUpdateMetric request")

	batch := []*model.MetricsV2{}
	for _, metric := range in.Metric {
		batch = append(batch, MetricsV2FromPbMetric(metric))
	}
	updated, err := s.metricService.BatchUpsertMetricValue(ctx, batch)
	if err != nil {
		return nil, status.Errorf(codes.Unknown, "Batch update metric error: %s", err)
	}

	var response pb.BatchUpdateResponse
	for _, metric := range updated {
		response.Metric = append(response.Metric, MetricsV2ToPbMetric(metric))
	}

	return &response, nil
}
