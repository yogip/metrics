package metrics

import (
	"github.com/yogip/metrics/internal/storage"
)

var PollCountCounter = storage.Counter{Name: "PollCount"}

var AllocGauge = storage.Gauge{Name: "Alloc"}
var BuckHashSysGauge = storage.Gauge{Name: "BuckHashSys"}
var FreesGauge = storage.Gauge{Name: "Frees"}
var GCCPUFractionGauge = storage.Gauge{Name: "GCCPUFraction"}
var GCSysGauge = storage.Gauge{Name: "GCSys"}
var HeapAllocGauge = storage.Gauge{Name: "HeapAlloc"}
var HeapIdleGauge = storage.Gauge{Name: "HeapIdle"}
var HeapInuseGauge = storage.Gauge{Name: "HeapInuse"}
var HeapObjectsGauge = storage.Gauge{Name: "HeapObjects"}
var HeapReleasedGauge = storage.Gauge{Name: "HeapReleased"}
var HeapSysGauge = storage.Gauge{Name: "HeapSys"}
var LastGCGauge = storage.Gauge{Name: "LastGC"}
var LookupsGauge = storage.Gauge{Name: "Lookups"}
var MCacheInuseGauge = storage.Gauge{Name: "MCacheInuse"}
var MCacheSysGauge = storage.Gauge{Name: "MCacheSys"}
var MSpanInuseGauge = storage.Gauge{Name: "MSpanInuse"}
var MSpanSysGauge = storage.Gauge{Name: "MSpanSys"}
var MallocsGauge = storage.Gauge{Name: "Mallocs"}
var NextGCGauge = storage.Gauge{Name: "NextGC"}
var NumForcedGCGauge = storage.Gauge{Name: "NumForcedGC"}
var NumGCGauge = storage.Gauge{Name: "NumGC"}
var OtherSysGauge = storage.Gauge{Name: "OtherSys"}
var PauseTotalNsGauge = storage.Gauge{Name: "PauseTotalNs"}
var StackInuseGauge = storage.Gauge{Name: "StackInuse"}
var StackSysGauge = storage.Gauge{Name: "StackSys"}
var SysGauge = storage.Gauge{Name: "Sys"}
var TotalAllocGauge = storage.Gauge{Name: "TotalAlloc"}

var AllMetrics []storage.Metric = []storage.Metric{
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
