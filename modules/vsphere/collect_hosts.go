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
		return errors.New("failed to gather hosts metrics")
	}

	vs.collectHostsMetrics(mx, metrics)
	return nil
}

func (vs *VSphere) collectHostsMetrics(mx map[string]int64, metrics []performance.EntityMetric) {
	vs.nilHostsMetrics()

	for _, m := range metrics {
		host := vs.resources.Hosts.Get(m.Entity.Value)
		if host == nil {
			continue
		}

		host.Metrics = m.Value
		writeHostMetricsTo(mx, host)
	}

}

func (vs *VSphere) nilHostsMetrics() {
	for _, v := range vs.resources.Hosts {
		v.Metrics = nil
	}
}

func writeHostMetricsTo(to map[string]int64, host *rs.Host) {
	for _, m := range host.Metrics {
		if len(m.Value) == 0 || m.Value[0] == -1 {
			continue
		}
		key := buildHostKey(host, m.Instance, m.Name)
		to[key] = m.Value[0]
	}
}

func buildHostKey(h *rs.Host, instance string, metricName string) string {
	// NOTE: name is not unique
	if instance == "" {
		return fmt.Sprintf("%s_%s_%s", h.ID, h.Name, metricName)
	}
	return fmt.Sprintf("%s_%s_%s_%s", h.ID, h.Name, metricName, instance)
}
