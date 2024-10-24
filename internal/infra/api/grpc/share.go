package grpc

import (
	"metrics/internal/core/model"
	pb "metrics/internal/proto"
)

func MetricsV2ToPbMetric(metric *model.MetricsV2) *pb.Metric {
	if metric.MType == model.GaugeType {
		return &pb.Metric{
			Name:  metric.ID,
			Type:  pb.MetricType_GAUGE,
			Value: *metric.Value,
		}
	}

	return &pb.Metric{
		Name:      metric.ID,
		Type:      pb.MetricType_COUNTER,
		Increment: *metric.Delta,
	}
}

func MetricsV2FromPbMetric(metric *pb.Metric) *model.MetricsV2 {
	if metric.Type == pb.MetricType_GAUGE {
		return &model.MetricsV2{
			ID:    metric.Name,
			MType: model.GaugeType,
			Value: &metric.Value,
		}
	}

	return &model.MetricsV2{
		ID:    metric.Name,
		MType: model.CounterType,
		Delta: &metric.Increment,
	}
}
