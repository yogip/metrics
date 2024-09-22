package agent

import (
	"context"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	"metrics/internal/agent/config"
	"metrics/internal/agent/metrics"
	"metrics/internal/core/model"
	"metrics/internal/infra/api/rest/middlewares"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPollFromRuntime(t *testing.T) {
	lock := &sync.Mutex{}
	cfg := config.AgentConfig{PollInterval: 1}

	ctx := context.Background()
	ctxT, cancel := context.WithTimeout(ctx, time.Duration(cfg.PollInterval)*time.Second+time.Second)
	defer cancel()

	metricRuntimePoller(ctxT, &cfg, lock)

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
		assert.Positive(t, metrics.PollCountCounter.Value)
	})
}

func TestPollFromPsutils(t *testing.T) {
	lock := &sync.Mutex{}
	cfg := config.AgentConfig{PollInterval: 1}

	ctx := context.Background()
	ctxT, cancel := context.WithTimeout(ctx, time.Duration(cfg.PollInterval)*time.Second+time.Second)
	defer cancel()

	metricPollerPsutils(ctxT, &cfg, lock)

	gauges := []metrics.Gauge{
		metrics.TotalMemory,
		metrics.FreeMemory,
	}

	for _, gauge := range gauges {
		t.Run(gauge.Name, func(t *testing.T) {
			assert.Greater(t, gauge.Value, 0.0)
		})
	}
	for _, gauge := range metrics.CPUutilizations {
		t.Run(gauge.Name, func(t *testing.T) {
			assert.GreaterOrEqual(t, gauge.Value, 0.0)
		})
	}
}

func TestRun(t *testing.T) {
	var called bool

	// Create a test server
	srv := gin.New()
	srv.Use(middlewares.GzipDecompressMiddleware())
	srv.POST("/updates", func(c *gin.Context) {
		var actualMetrics []model.MetricsV2
		err := c.BindJSON(&actualMetrics)
		require.NoError(t, err)

		assert.NotEmpty(t, actualMetrics)
		called = true
	})
	testSrv := httptest.NewServer(srv)

	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, time.Duration(4)*time.Second)
	defer cancel()

	Run(ctx, &config.AgentConfig{
		ServerAddresPort: testSrv.URL,
		ReportInterval:   2,
		PollInterval:     1,
		RateLimit:        1,
	})

	<-ctx.Done()

	// Verify the result
	assert.True(t, called)
}
