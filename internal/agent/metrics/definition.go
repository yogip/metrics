package metrics

import (
	"github.com/yogip/metrics/internal/models"
)

var PollCountCounter = models.Counter{Name: "PollCount"}

var AllocGauge = models.Gauge{Name: "Alloc"}
var BuckHashSysGauge = models.Gauge{Name: "BuckHashSys"}
var FreesGauge = models.Gauge{Name: "Frees"}
var GCCPUFractionGauge = models.Gauge{Name: "GCCPUFraction"}
var GCSysGauge = models.Gauge{Name: "GCSys"}
var HeapAllocGauge = models.Gauge{Name: "HeapAlloc"}
var HeapIdleGauge = models.Gauge{Name: "HeapIdle"}
var HeapInuseGauge = models.Gauge{Name: "HeapInuse"}
var HeapObjectsGauge = models.Gauge{Name: "HeapObjects"}
var HeapReleasedGauge = models.Gauge{Name: "HeapReleased"}
var HeapSysGauge = models.Gauge{Name: "HeapSys"}
var LastGCGauge = models.Gauge{Name: "LastGC"}
var LookupsGauge = models.Gauge{Name: "Lookups"}
var MCacheInuseGauge = models.Gauge{Name: "MCacheInuse"}
var MCacheSysGauge = models.Gauge{Name: "MCacheSys"}
var MSpanInuseGauge = models.Gauge{Name: "MSpanInuse"}
var MSpanSysGauge = models.Gauge{Name: "MSpanSys"}
var MallocsGauge = models.Gauge{Name: "Mallocs"}
var NextGCGauge = models.Gauge{Name: "NextGC"}
var NumForcedGCGauge = models.Gauge{Name: "NumForcedGC"}
var NumGCGauge = models.Gauge{Name: "NumGC"}
var OtherSysGauge = models.Gauge{Name: "OtherSys"}
var PauseTotalNsGauge = models.Gauge{Name: "PauseTotalNs"}
var StackInuseGauge = models.Gauge{Name: "StackInuse"}
var StackSysGauge = models.Gauge{Name: "StackSys"}
var SysGauge = models.Gauge{Name: "Sys"}
var TotalAllocGauge = models.Gauge{Name: "TotalAlloc"}

var AllMetrics []models.Metric = []models.Metric{
	&PollCountCounter,
	&AllocGauge,
	&BuckHashSysGauge,
	&FreesGauge,
	&GCCPUFractionGauge,
	&GCSysGauge,
	&HeapAllocGauge,
	&HeapIdleGauge,
	&HeapInuseGauge,
	&HeapObjectsGauge,
	&HeapReleasedGauge,
	&HeapSysGauge,
	&LastGCGauge,
	&LookupsGauge,
	&MCacheInuseGauge,
	&MCacheSysGauge,
	&MSpanInuseGauge,
	&MSpanSysGauge,
	&MallocsGauge,
	&NextGCGauge,
	&NumForcedGCGauge,
	&NumGCGauge,
	&OtherSysGauge,
	&PauseTotalNsGauge,
	&StackInuseGauge,
	&StackSysGauge,
	&SysGauge,
	&TotalAllocGauge,
}
