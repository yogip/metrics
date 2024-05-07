package agent

import (
	"log"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"metrics/internal/agent/config"
	"metrics/internal/agent/metrics"
	"metrics/internal/agent/transport"
)

func pollFromRuntime() {
	log.Println("Polling metrics")
	var rtm runtime.MemStats
	runtime.ReadMemStats(&rtm)

	// Counters metric
	metrics.PollCountCounter.Incremet(1)

	// Gauge metrics
	metrics.AllocGauge.Set(float64(rtm.Alloc))
	metrics.BuckHashSysGauge.Set(float64(rtm.BuckHashSys))
	metrics.FreesGauge.Set(float64(rtm.Frees))
	metrics.GCCPUFractionGauge.Set(float64(rtm.GCCPUFraction))
	metrics.GCSysGauge.Set(float64(rtm.GCSys))
	metrics.HeapAllocGauge.Set(float64(rtm.HeapAlloc))
	metrics.HeapIdleGauge.Set(float64(rtm.HeapAlloc))
	metrics.HeapInuseGauge.Set(float64(rtm.HeapAlloc))
	metrics.HeapObjectsGauge.Set(float64(rtm.HeapAlloc))
	metrics.HeapReleasedGauge.Set(float64(rtm.HeapAlloc))
	metrics.HeapSysGauge.Set(float64(rtm.HeapAlloc))
	metrics.LastGCGauge.Set(float64(rtm.HeapAlloc))
	metrics.LookupsGauge.Set(float64(rtm.HeapAlloc))
	metrics.MCacheInuseGauge.Set(float64(rtm.HeapAlloc))
	metrics.MCacheSysGauge.Set(float64(rtm.HeapAlloc))
	metrics.MSpanInuseGauge.Set(float64(rtm.HeapAlloc))
	metrics.MSpanSysGauge.Set(float64(rtm.HeapAlloc))
	metrics.MallocsGauge.Set(float64(rtm.HeapAlloc))
	metrics.NextGCGauge.Set(float64(rtm.HeapAlloc))
	metrics.NumForcedGCGauge.Set(float64(rtm.HeapAlloc))
	metrics.NumGCGauge.Set(float64(rtm.HeapAlloc))
	metrics.OtherSysGauge.Set(float64(rtm.HeapAlloc))
	metrics.PauseTotalNsGauge.Set(float64(rtm.HeapAlloc))
	metrics.StackInuseGauge.Set(float64(rtm.HeapAlloc))
	metrics.StackSysGauge.Set(float64(rtm.HeapAlloc))
	metrics.SysGauge.Set(float64(rtm.HeapAlloc))
	metrics.TotalAllocGauge.Set(float64(rtm.HeapAlloc))
}

func reportMetrics(client metrics.Transporter) {
	log.Println("Reporting all metrics")
	for _, metric := range metrics.AllMetrics {
		err := metric.Send(client)
		if err != nil {
			log.Printf("sending metric error: %s", err)
		}
	}
}

func Run(config *config.AgentConfig) {
	pollTicker := time.NewTicker(config.PollInterval)
	reportTicker := time.NewTicker(config.ReportInterval)
	client := transport.NewClient(config.ServerAddresPort)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	for {
		select {
		case <-quit:
			log.Println("Received Ctrl+C, stopping...")
			return
		case <-pollTicker.C:
			pollFromRuntime()
		case <-reportTicker.C:
			reportMetrics(client)
		}
	}
}
