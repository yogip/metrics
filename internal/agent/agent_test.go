package agent

import (
	"testing"

	// "time"

	"github.com/stretchr/testify/assert"
	"github.com/yogip/metrics/internal/agent/metrics"
	"github.com/yogip/metrics/internal/models"
)

func TestPollFromRuntime(t *testing.T) {
	pollFromRuntime()

	for _, metric := range metrics.AllMetrics {
		if metric.GetName() == "GCCPUFraction" {
			continue
		}
		t.Run(metric.GetName(), func(t *testing.T) {
			if metric.Type() == models.GaugeType {
				gauge, _ := metric.(*models.Gauge)
				assert.Greater(t, gauge.Value, 0.0)
			} else {
				counter, _ := metric.(*models.Counter)
				assert.Greater(t, counter.Value, int64(0))
			}
		})
	}
}
