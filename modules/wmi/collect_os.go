package wmi

import (
	"fmt"

	"github.com/netdata/go.d.plugin/pkg/prometheus"
)

const (
	collectorOS = "os"

	metricOSPhysicalMemoryFreeBytes = "wmi_os_physical_memory_free_bytes"
	metricOSPagingFreeBytes         = "wmi_os_paging_free_bytes"
	metricOSVirtualMemoryFreeBytes  = "wmi_os_virtual_memory_free_bytes"
	metricOSProcessesLimit          = "wmi_os_processes_limit"
	metricOSProcessMemoryLimitBytes = "wmi_os_process_memory_limit_bytes"
	metricOSProcesses               = "wmi_os_processes"
	metricOSUsers                   = "wmi_os_users"
	metricOSPagingLimitBytes        = "wmi_os_paging_limit_bytes"
	metricOSVirtualMemoryBytes      = "wmi_os_virtual_memory_bytes"
	metricOSVisibleMemoryBytes      = "wmi_os_visible_memory_bytes"
	metricOSTime                    = "wmi_os_time"
)

func collectOS(mx *metrics, pms prometheus.Metrics) bool {
	enabled, success := checkCollector(pms, collectorOS)
	if !(enabled && success) {
		return false
	}
	mx.OS = &os{}

	names := []string{
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

	for _, name := range names {
		collectOSAny(mx, pms, name)
	}

	return true
}

func collectOSAny(mx *metrics, pms prometheus.Metrics, name string) {
	value := pms.FindByName(name).Max()

	switch name {
	default:
		panic(fmt.Sprintf("unknown metric name during OS collection : %s", name))
	case metricOSPhysicalMemoryFreeBytes:
		mx.OS.PhysicalMemoryFreeBytes = value
	case metricOSPagingFreeBytes:
		mx.OS.PagingFreeBytes = value
	case metricOSVirtualMemoryFreeBytes:
		mx.OS.VirtualMemoryFreeBytes = value
	case metricOSProcessesLimit:
		mx.OS.ProcessesLimit = value
	case metricOSProcessMemoryLimitBytes:
		mx.OS.ProcessMemoryLimitBytes = value
	case metricOSProcesses:
		mx.OS.Processes = value
	case metricOSUsers:
		mx.OS.Users = value
	case metricOSPagingLimitBytes:
		mx.OS.PagingLimitBytes = value
	case metricOSVirtualMemoryBytes:
		mx.OS.VirtualMemoryBytes = value
	case metricOSVisibleMemoryBytes:
		mx.OS.VisibleMemoryBytes = value
	case metricOSTime:
		mx.OS.Time = value
	}
}
