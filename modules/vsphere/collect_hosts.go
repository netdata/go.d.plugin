package vsphere

import (
	"errors"
	"fmt"

	rs "github.com/netdata/go.d.plugin/modules/vsphere/resources"

	"github.com/vmware/govmomi/performance"
)

func (vs *VSphere) collectHosts(mx map[string]int64) error {
	// NOTE: returns unsorted if at least one types.PerfMetricId Instance is not ""
	metrics := vs.ScrapeHostsMetrics(vs.resources.Hosts)
	if len(metrics) == 0 {
		return errors.New("failed to scrape hosts metrics")
	}

	collected := vs.collectHostsMetrics(mx, metrics)
	vs.updateHostsCharts(collected)
	return nil
}

func (vs *VSphere) collectHostsMetrics(mx map[string]int64, metrics []performance.EntityMetric) map[string]bool {
	collected := make(map[string]bool)
	for _, m := range metrics {
		host := vs.resources.Hosts.Get(m.Entity.Value)
		if host == nil {
			continue
		}
		writeHostMetrics(mx, host, m.Value)
		collected[host.ID] = true
	}

	for k := range vs.discoveredHosts {
		if collected[k] {
			vs.discoveredHosts[k] = 0
		} else {
			vs.discoveredHosts[k] += 1
		}
	}
	return collected
}

func writeHostMetrics(dst map[string]int64, host *rs.Host, metrics []performance.MetricSeries) {
	for _, m := range metrics {
		if len(m.Value) == 0 || m.Value[0] == -1 {
			continue
		}
		key := fmt.Sprintf("%s_%s", host.ID, m.Name)
		dst[key] = m.Value[0]
	}
	key := fmt.Sprintf("%s_overall.status", host.ID)
	dst[key] = overallStatusToInt(host.OverallStatus)
}

func overallStatusToInt(status string) int64 {
	// ManagedEntityStatus
	switch status {
	default:
		return 0
	case "grey":
		return 1
	case "green":
		return 2
	case "yellow":
		return 3
	case "red":
		return 4
	}
}
