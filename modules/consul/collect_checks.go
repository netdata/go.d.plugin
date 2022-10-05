// SPDX-License-Identifier: GPL-3.0-or-later

package consul

const urlPathChecks = "/v1/agent/checks"

type healthCheck struct {
	Node        string
	CheckID     string
	Name        string
	Status      string
	ServiceID   string
	ServiceName string
	ServiceTags []string
}

func (c *Consul) collectLocalChecks(mx map[string]int64) error {
	var checks map[string]*healthCheck

	if err := c.doOKDecode(urlPathChecks, &checks); err != nil {
		return err
	}

	for id, check := range checks {
		if !c.checksSr.MatchString(id) {
			c.Debugf("check with id '%s' does not match the selector ('%s'), skipping it", id, c.ChecksSelector)
			continue
		}

		if !c.checks[id] {
			c.checks[id] = true
			c.addHealthCheckCharts(check)
		}

		mx["health_check_"+id+"_passing_status"] = boolToInt(check.Status == "passing")
		mx["health_check_"+id+"_warning_status"] = boolToInt(check.Status == "warning")
		mx["health_check_"+id+"_critical_status"] = boolToInt(check.Status == "critical")
		mx["health_check_"+id+"_maintenance_status"] = boolToInt(check.Status == "maintenance")
	}

	for id := range c.checks {
		if _, ok := checks[id]; !ok {
			delete(c.checks, id)
			c.removeHealthCheckCharts(id)
		}
	}

	return nil
}

func boolToInt(v bool) int64 {
	if v {
		return 1
	}
	return 0
}
