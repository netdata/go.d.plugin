package wmi

import (
	"github.com/netdata/go.d.plugin/pkg/prometheus"
)

const (
	collectorOS = "os"

	metricOSPhysicalMemoryFreeBytes = "windows_os_physical_memory_free_bytes"
	metricOSPagingFreeBytes         = "windows_os_paging_free_bytes"
	metricOSVirtualMemoryFreeBytes  = "windows_os_virtual_memory_free_bytes"
	metricOSProcessesLimit          = "windows_os_processes_limit"
	metricOSProcessMemoryLimitBytes = "windows_os_process_memory_limit_bytes"
	metricOSProcesses               = "windows_os_processes"
	metricOSUsers                   = "windows_os_users"
	metricOSPagingLimitBytes        = "windows_os_paging_limit_bytes"
	metricOSVirtualMemoryBytes      = "windows_os_virtual_memory_bytes"
	metricOSVisibleMemoryBytes      = "windows_os_visible_memory_bytes"
	metricOSTime                    = "windows_os_time"
)

var osMetricsNames = []string{
	metricOSPhysicalMemoryFreeBytes,
	metricOSPagingFreeBytes,
	metricOSVirtualMemoryFreeBytes,
	metricOSProcessesLimit,
	metricOSProcessMemoryLimitBytes,
	metricOSProcesses,
	metricOSUsers,
	metricOSPagingLimitBytes,
	metricOSVirtualMemoryBytes,
	metricOSVisibleMemoryBytes,
	metricOSTime,
}

func doCollectOS(pms prometheus.Metrics) bool {
	enabled, success := checkCollector(pms, collectorOS)
	return enabled && success
}

func collectOS(pms prometheus.Metrics) *osMetrics {
	if !doCollectOS(pms) {
		return nil
	}

	osm := &osMetrics{}
	for _, name := range osMetricsNames {
		collectOSMetric(osm, pms, name)
	}
	return osm
}

func collectOSMetric(osm *osMetrics, pms prometheus.Metrics, name string) {
	value := pms.FindByName(name).Max()
	assignOSMetric(osm, name, value)
}

func assignOSMetric(mx *osMetrics, name string, value float64) {
	switch name {
	case metricOSPhysicalMemoryFreeBytes:
		mx.PhysicalMemoryFreeBytes = value
	case metricOSPagingFreeBytes:
		mx.PagingFreeBytes = value
	case metricOSVirtualMemoryFreeBytes:
		mx.VirtualMemoryFreeBytes = value
	case metricOSProcessesLimit:
		mx.ProcessesLimit = value
	case metricOSProcessMemoryLimitBytes:
		mx.ProcessMemoryLimitBytes = value
	case metricOSProcesses:
		mx.Processes = value
	case metricOSUsers:
		mx.Users = value
	case metricOSPagingLimitBytes:
		mx.PagingLimitBytes = value
	case metricOSVirtualMemoryBytes:
		mx.VirtualMemoryBytes = value
	case metricOSVisibleMemoryBytes:
		mx.VisibleMemoryBytes = value
	case metricOSTime:
		mx.Time = value
	}
}
