package agent

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"runtime"
	"sync"
	"syscall"
	"time"

	"metrics/internal/agent/config"
	"metrics/internal/agent/metrics"
	"metrics/internal/agent/transport"
	"metrics/internal/core/model"
	"metrics/internal/logger"
	"metrics/internal/retrier"

	"go.uber.org/zap"
)

func pollFromRuntime(lock *sync.Mutex) {
	if !lock.TryLock() {
		logger.Log.Warn("Skip gathering metrics")
		return
	}
	defer lock.Unlock()
	logger.Log.Debug("Gathering metrics")

	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	var rtm runtime.MemStats
	runtime.ReadMemStats(&rtm)

	// Counters metric
	metrics.PollCountCounter.Increment(1)

	// Gauge metrics
	metrics.RandomValue.Set(r.Float64())
	metrics.AllocGauge.Set(float64(rtm.Alloc))
	metrics.BuckHashSysGauge.Set(float64(rtm.BuckHashSys))
	metrics.FreesGauge.Set(float64(rtm.Frees))
	metrics.GCCPUFractionGauge.Set(float64(rtm.GCCPUFraction))
	metrics.GCSysGauge.Set(float64(rtm.GCSys))
	metrics.HeapAllocGauge.Set(float64(rtm.HeapAlloc))
	metrics.HeapIdleGauge.Set(float64(rtm.HeapIdle))
	metrics.HeapInuseGauge.Set(float64(rtm.HeapInuse))
	metrics.HeapObjectsGauge.Set(float64(rtm.HeapObjects))
	metrics.HeapReleasedGauge.Set(float64(rtm.HeapReleased))
	metrics.HeapSysGauge.Set(float64(rtm.HeapSys))
	metrics.LastGCGauge.Set(float64(rtm.LastGC))
	metrics.LookupsGauge.Set(float64(rtm.Lookups))
	metrics.MCacheInuseGauge.Set(float64(rtm.MCacheInuse))
	metrics.MCacheSysGauge.Set(float64(rtm.MCacheSys))
	metrics.MSpanInuseGauge.Set(float64(rtm.MSpanInuse))
	metrics.MSpanSysGauge.Set(float64(rtm.MSpanSys))
	metrics.MallocsGauge.Set(float64(rtm.Mallocs))
	metrics.NextGCGauge.Set(float64(rtm.NextGC))
	metrics.NumForcedGCGauge.Set(float64(rtm.NumForcedGC))
	metrics.NumGCGauge.Set(float64(rtm.NumGC))
	metrics.OtherSysGauge.Set(float64(rtm.OtherSys))
	metrics.PauseTotalNsGauge.Set(float64(rtm.PauseTotalNs))
	metrics.StackInuseGauge.Set(float64(rtm.StackInuse))
	metrics.StackSysGauge.Set(float64(rtm.StackSys))
	metrics.SysGauge.Set(float64(rtm.Sys))
	metrics.TotalAllocGauge.Set(float64(rtm.TotalAlloc))
}

func reportMetrics(ctx context.Context, client metrics.Transporter, lock *sync.Mutex) {
	if !lock.TryLock() {
		logger.Log.Warn("Skip reporting all metrics")
		return
	}
	defer lock.Unlock()
	logger.Log.Debug("Reporting all metrics")

	ret := &retrier.Retrier{
		Strategy: retrier.Backoff(
			3,             // max attempts
			1*time.Second, // initial delay
			3,             // multiplier
			5*time.Second, // max delay
		),
		OnRetry: func(ctx context.Context, n int, err error) {
			logger.Log.Debug(fmt.Sprintf("reportMetrics retry #%d: %v", n, err))
		},
	}

	data := []*model.MetricsV2{}
	for _, metric := range metrics.AllMetrics {
		data = append(data, metric.Payload())
	}

	fun := func() error {
		return client.SendMetric(data)
	}

	if err := ret.Do(ctx, fun, syscall.ECONNREFUSED); err != nil {
		logger.Log.Error("sending metric error", zap.String("error", err.Error()))
		return
	}
	for _, metric := range metrics.AllMetrics {
		metric.WasSend()
	}
}

func Run(config *config.AgentConfig) {
	pollTicker := time.NewTicker(time.Duration(config.PollInterval) * time.Second)
	defer pollTicker.Stop()

	reportTicker := time.NewTicker(time.Duration(config.ReportInterval) * time.Second)
	defer reportTicker.Stop()

	client := transport.NewClient(config.ServerAddresPort)
	ctx, cancel := context.WithCancel(context.Background())

	lock := &sync.Mutex{}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	for {
		select {
		case <-quit:
			logger.Log.Info("Received Ctrl+C, stopping...")
			cancel()
			return
		case <-pollTicker.C:
			go pollFromRuntime(lock)
		case <-reportTicker.C:
			go reportMetrics(ctx, client, lock)
		}
	}
}
