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

func Run(config *config.AgentConfig) {
	ctx, cancel := context.WithCancel(context.Background())

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	lock := &sync.Mutex{}

	go metricPoller(ctx, config, lock)
	go metricReporter(ctx, config, lock)

	<-quit
	logger.Log.Info("Received Ctrl+C, stopping...")
	cancel()
}

func metricReporter(ctx context.Context, cfg *config.AgentConfig, lock *sync.Mutex) {
	reportTicker := time.NewTicker(time.Duration(cfg.ReportInterval) * time.Second)
	defer reportTicker.Stop()

	metricsCh := make(chan []model.MetricsV2, cfg.RateLimit)
	for i := 0; i < cfg.RateLimit; i++ {
		go metricReporterWorker(ctx, cfg, metricsCh, i)
	}

	for {
		select {
		case <-ctx.Done():
			logger.Log.Info("Stop posting metrics")
			return
		case <-reportTicker.C:
			lock.Lock()

			data := []model.MetricsV2{}
			for _, metric := range metrics.AllMetrics {
				data = append(data, metric.Payload())
				metric.WasSent()
			}
			lock.Unlock()

			metricsCh <- data
		}
	}

}

func metricReporterWorker(
	ctx context.Context,
	cfg *config.AgentConfig,
	metricsCh chan []model.MetricsV2,
	workerID int,
) {
	logger.Log.Info(fmt.Sprintf("Start worker N: %d", workerID))
	client := transport.NewClient(cfg.ServerAddresPort, cfg.HashKey)

	for {
		select {
		case <-ctx.Done():
			logger.Log.Info("Exit from metricReporterWorker")
			return
		case data := <-metricsCh:
			postMetrics(ctx, client, data, workerID)
		}
	}
}

func postMetrics(ctx context.Context, client metrics.Transporter, data []model.MetricsV2, workerID int) {
	logger.Log.Debug(fmt.Sprintf("Reporting metrics. Worker ID: %d", workerID))

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

	fun := func() error {
		return client.SendMetric(data)
	}

	if err := ret.Do(ctx, fun, syscall.ECONNREFUSED); err != nil {
		logger.Log.Error("sending metric error", zap.String("error", err.Error()))
		return
	}
}

func metricPoller(ctx context.Context, cfg *config.AgentConfig, lock *sync.Mutex) {
	pollTicker := time.NewTicker(time.Duration(cfg.PollInterval) * time.Second)
	defer pollTicker.Stop()

	for {
		select {
		case <-ctx.Done():
			logger.Log.Info("Stop to poll metrics...")
			return
		case <-pollTicker.C:
			metricPollFromRuntime(lock)
		}
	}
}

func metricPollFromRuntime(lock *sync.Mutex) {
	logger.Log.Debug("Gathering metrics")
	lock.Lock()
	defer lock.Unlock()

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
