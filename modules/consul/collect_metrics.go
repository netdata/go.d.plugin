// SPDX-License-Identifier: GPL-3.0-or-later

package consul

type agentMetrics struct {
	Gauges []struct {
		Name   string
		Value  int64
		Labels map[string]string
	}
	Counters []struct {
		Name   string
		Count  int64
		Labels map[string]string
	}
	Samples []struct {
		Name   string
		Count  int64
		Labels map[string]string
	}
}

// https://www.consul.io/api-docs/agent#view-metrics
const urlPathAgentMetrics = "/v1/agent/metrics"

func (c *Consul) collectAgentMetrics(mx map[string]int64) error {
	var metrics agentMetrics

	if err := c.doOKDecode(urlPathAgentMetrics, &metrics); err != nil {
		return err
	}

	for _, m := range metrics.Gauges {
		switch m.Name {
		case "consul.server.isLeader":
			mx[m.Name+".yes"] = boolToInt(m.Value == 1)
			mx[m.Name+".no"] = boolToInt(m.Value != 1)
		case "consul.autopilot.healthy":
			mx[m.Name+".yes"] = boolToInt(m.Value == 1)
			mx[m.Name+".no"] = boolToInt(m.Value != 1)
		case
			"consul.autopilot.failure_tolerance",
			"consul.runtime.alloc_bytes",
			"consul.runtime.sys_bytes",
			"consul.runtime.total_gc_pause_ns":
			mx[m.Name] = m.Value
		}
	}

	for _, m := range metrics.Counters {
		switch m.Name {
		case "consul.client.rpc":
			mx[m.Name] = m.Count
		}
	}

	for _, m := range metrics.Samples {
		switch m.Name {
		case "consul.client.rpc":
			mx[m.Name] = m.Count
		}
	}

	return nil
}
