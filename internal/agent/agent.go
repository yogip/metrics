package agent

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"runtime"
	"sync"
	"syscall"
	"time"

	"metrics/internal/agent/metrics"
	"metrics/internal/agent/transport"
)

const pollInterval int = 2
const reportInterval int = 10

// Method gather all metrics from runtime
func pollFromRuntime() {
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

// Method gather all metrics from runtime and sleep for pollInterval
func pollLauncher() {
	for {
		log.Println("Polling all metrics")
		pollFromRuntime()
		time.Sleep(time.Duration(pollInterval) * time.Second)
	}
}

func reportMetrics(metric metrics.Metric, client metrics.Transporter) error {
	err := metric.Send(client)
	if err != nil {
		return fmt.Errorf("sending metric error: %s", err)
	}
	return nil
}

// Method send all metrics to server and sleep for reportInterval
func reportLauncher() {
	client := transport.NewClient("http://localhost:8080")
	for {
		time.Sleep(time.Duration(reportInterval) * time.Second)
		log.Println("Reporting all metrics")
		for _, metric := range metrics.AllMetrics {
			reportMetrics(metric, client)
		}
	}
}

// Method create 2 goroutines for polling and reporting metrics
func Run() {
	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		pollLauncher()
	}()

	go func() {
		defer wg.Done()
		reportLauncher()
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	select {
	case <-quit:
		log.Println("Received Ctrl+C, stopping...")
	case <-func() chan struct{} {
		done := make(chan struct{})
		go func() {
			wg.Wait()
			close(done)
		}()
		return done
	}():
		log.Println("All goroutines have finished, stopping...")
	}
}
