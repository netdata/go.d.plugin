package wmi

import "github.com/netdata/go.d.plugin/pkg/prometheus"

const (
	metricCSLogicalProcessors   = "wmi_cs_logical_processors"
	metricCSPhysicalMemoryBytes = "wmi_cs_physical_memory_bytes"
)

func (w *WMI) collectCS(mx *metrics, pms prometheus.Metrics) bool {
	enabled, success := checkCollector(pms, collectorCS)
	if !(enabled && success) {
		return false
	}
	mx.CS = &cs{}

	mx.CS.LogicalProcessors = pms.FindByName(metricCSLogicalProcessors).Max()
	mx.CS.PhysicalMemoryBytes = pms.FindByName(metricCSPhysicalMemoryBytes).Max()

	return true
}
