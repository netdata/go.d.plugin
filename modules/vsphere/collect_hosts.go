package vsphere

import (
	"errors"
	"fmt"

	rs "github.com/netdata/go.d.plugin/modules/vsphere/resources"

	"github.com/vmware/govmomi/performance"
)

func (vs *VSphere) collectHosts(mx map[string]int64) error {
	// NOTE: returns unsorted if at least one types.PerfMetricId Instance is not ""
	metrics := vs.CollectHostsMetrics(vs.resources.Hosts)
	if len(metrics) == 0 {
		return errors.New("failed to collect hosts metrics")
	}

	vs.processHostsMetrics(mx, metrics)
	return nil
}

func (vs *VSphere) processHostsMetrics(mx map[string]int64, metrics []performance.EntityMetric) {
	updated := make(map[string]bool)
	for _, m := range metrics {
		host := vs.resources.Hosts.Get(m.Entity.Value)
		if host == nil {
			continue
		}
		writeHostMetrics(mx, host, m.Value)
		updated[host.ID] = true
		vs.collectedHosts[host.ID] = 0
	}

	for k := range vs.collectedHosts {
		if updated[k] {
			continue
		}
		vs.collectedHosts[k] += 1
	}
}

func writeHostMetrics(dst map[string]int64, host *rs.Host, metrics []performance.MetricSeries) {
	for _, m := range metrics {
		if len(m.Value) == 0 || m.Value[0] == -1 {
			continue
		}
		key := hostMetricKey(host, m.Instance, m.Name)
		dst[key] = m.Value[0]
	}
}

func hostMetricKey(host *rs.Host, instance, metricName string) string {
	if instance == "" {
		return fmt.Sprintf("%s_%s_%s", host.ID, host.Name, metricName)
	}
	return fmt.Sprintf("%s_%s_%s_%s", host.ID, host.Name, metricName, instance)
}
