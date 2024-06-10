package agent

import (
	"math/rand"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"metrics/internal/agent/config"
	"metrics/internal/agent/metrics"
	"metrics/internal/agent/transport"
	"metrics/internal/core/model"
	"metrics/internal/logger"

	"go.uber.org/zap"
)

func pollFromRuntime() {
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

func reportMetrics(client metrics.Transporter) {
	logger.Log.Debug("Reporting all metrics")
	data := []*model.MetricsV2{}
	for _, metric := range metrics.AllMetrics {
		data = append(data, metric.Payload())
	}
	err := client.SendMetric(data)
	if err != nil {
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

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	for {
		select {
		case <-quit:
			logger.Log.Info("Received Ctrl+C, stopping...")
			return
		case <-pollTicker.C:
			pollFromRuntime()
		case <-reportTicker.C:
			reportMetrics(client)
		}
	}
}
