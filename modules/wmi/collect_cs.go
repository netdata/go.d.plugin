package wmi

import (
	"github.com/netdata/go.d.plugin/pkg/prometheus"
)

const (
	collectorCS = "cs"

	metricCSLogicalProcessors   = "windows_cs_logical_processors"
	metricCSPhysicalMemoryBytes = "windows_cs_physical_memory_bytes"
)

var csMetricsNames = []string{
	metricCSLogicalProcessors,
	metricCSPhysicalMemoryBytes,
}

func doCollectCS(pms prometheus.Metrics) bool {
	enabled, success := checkCollector(pms, collectorCS)
	return enabled && success
}

func collectCS(pms prometheus.Metrics) *csMetrics {
	if !doCollectCS(pms) {
		return nil
	}

	csm := &csMetrics{}
	for _, name := range csMetricsNames {
		collectCSMetric(csm, pms, name)
	}
	return csm
}

func collectCSMetric(csm *csMetrics, pms prometheus.Metrics, name string) {
	value := pms.FindByName(name).Max()
	assignCSMetric(csm, name, value)
}

func assignCSMetric(csm *csMetrics, name string, value float64) {
	switch name {
	case metricCSLogicalProcessors:
		csm.LogicalProcessors = value
	case metricCSPhysicalMemoryBytes:
		csm.PhysicalMemoryBytes = value
	}
}
