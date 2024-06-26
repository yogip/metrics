package metrics

import (
	"metrics/internal/core/model"
)

var PollCountCounter = Counter{model.Counter{Name: "PollCount"}}

var RandomValue = Gauge{model.Gauge{Name: "RandomValue"}}
var AllocGauge = Gauge{model.Gauge{Name: "Alloc"}}
var BuckHashSysGauge = Gauge{model.Gauge{Name: "BuckHashSys"}}
var FreesGauge = Gauge{model.Gauge{Name: "Frees"}}
var GCCPUFractionGauge = Gauge{model.Gauge{Name: "GCCPUFraction"}}
var GCSysGauge = Gauge{model.Gauge{Name: "GCSys"}}
var HeapAllocGauge = Gauge{model.Gauge{Name: "HeapAlloc"}}
var HeapIdleGauge = Gauge{model.Gauge{Name: "HeapIdle"}}
var HeapInuseGauge = Gauge{model.Gauge{Name: "HeapInuse"}}
var HeapObjectsGauge = Gauge{model.Gauge{Name: "HeapObjects"}}
var HeapReleasedGauge = Gauge{model.Gauge{Name: "HeapReleased"}}
var HeapSysGauge = Gauge{model.Gauge{Name: "HeapSys"}}
var LastGCGauge = Gauge{model.Gauge{Name: "LastGC"}}
var LookupsGauge = Gauge{model.Gauge{Name: "Lookups"}}
var MCacheInuseGauge = Gauge{model.Gauge{Name: "MCacheInuse"}}
var MCacheSysGauge = Gauge{model.Gauge{Name: "MCacheSys"}}
var MSpanInuseGauge = Gauge{model.Gauge{Name: "MSpanInuse"}}
var MSpanSysGauge = Gauge{model.Gauge{Name: "MSpanSys"}}
var MallocsGauge = Gauge{model.Gauge{Name: "Mallocs"}}
var NextGCGauge = Gauge{model.Gauge{Name: "NextGC"}}
var NumForcedGCGauge = Gauge{model.Gauge{Name: "NumForcedGC"}}
var NumGCGauge = Gauge{model.Gauge{Name: "NumGC"}}
var OtherSysGauge = Gauge{model.Gauge{Name: "OtherSys"}}
var PauseTotalNsGauge = Gauge{model.Gauge{Name: "PauseTotalNs"}}
var StackInuseGauge = Gauge{model.Gauge{Name: "StackInuse"}}
var StackSysGauge = Gauge{model.Gauge{Name: "StackSys"}}
var SysGauge = Gauge{model.Gauge{Name: "Sys"}}
var TotalAllocGauge = Gauge{model.Gauge{Name: "TotalAlloc"}}

// todo
// var TotalMemory = Gauge{model.Gauge{Name: "TotalMemory"}}
// var FreeMemory = Gauge{model.Gauge{Name: "FreeMemory"}}
// var CPUutilization1 = Gauge{model.Gauge{Name: "CPUutilization1"}}

var AllMetrics []Metric = []Metric{
	&PollCountCounter,
	&RandomValue,
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
