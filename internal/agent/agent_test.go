package agent

import (
	"testing"

	"metrics/internal/agent/metrics"
	"metrics/internal/core/model"

	"github.com/stretchr/testify/assert"
)

func TestPollFromRuntime(t *testing.T) {
	pollFromRuntime()

	for _, metric := range metrics.AllMetrics {
		if metric.Type() == model.GaugeType {
			gauge, _ := metric.(*metrics.Gauge)
			if gauge.Name == "GCCPUFraction" {
				continue
			}
			assert.Greater(t, gauge.Value, 0.0)
		} else {
			counter, _ := metric.(*metrics.Counter)
			assert.Greater(t, counter.Value, int64(0))
		}
	}
}
