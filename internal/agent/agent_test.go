package agent

import (
	"sync"
	"testing"

	"metrics/internal/agent/metrics"

	"github.com/stretchr/testify/assert"
)

func TestPollFromRuntime(t *testing.T) {
	lock := &sync.Mutex{}
	metricPollFromRuntime(lock)

	gauges := []metrics.Gauge{
		metrics.RandomValue,
		metrics.AllocGauge,
		metrics.BuckHashSysGauge,
		metrics.GCSysGauge,
		metrics.HeapAllocGauge,
		metrics.HeapIdleGauge,
		metrics.HeapInuseGauge,
		metrics.HeapObjectsGauge,
		metrics.HeapReleasedGauge,
		metrics.HeapSysGauge,
		metrics.MCacheInuseGauge,
		metrics.MCacheSysGauge,
		metrics.MSpanInuseGauge,
		metrics.MSpanSysGauge,
		metrics.MallocsGauge,
		metrics.OtherSysGauge,
		metrics.StackInuseGauge,
		metrics.StackSysGauge,
		metrics.SysGauge,
		metrics.TotalAllocGauge,
	}

	for _, gauge := range gauges {
		t.Run(gauge.Name, func(t *testing.T) {
			assert.Greater(t, gauge.Value, 0.0)
		})
	}

	t.Run(metrics.PollCountCounter.Name, func(t *testing.T) {
		assert.Greater(t, metrics.PollCountCounter.Value, int64(0))
	})
}
