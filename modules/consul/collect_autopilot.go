// SPDX-License-Identifier: GPL-3.0-or-later

package consul

import "time"

const (
	// https://developer.hashicorp.com/consul/api-docs/operator/autopilot#read-health
	urlPathOperationAutopilotHealth = "/v1/operator/autopilot/health"
)

type autopilotHealth struct {
	Servers []struct {
		ID          string
		SerfStatus  string
		Leader      bool
		LastContact string
		Healthy     bool
		Voter       bool
		StableSince time.Time
	}
}

func (c *Consul) collectAutopilotHealth(mx map[string]int64) error {
	var health autopilotHealth

	if err := c.doOKDecode(urlPathOperationAutopilotHealth, &health); err != nil {
		return err
	}

	for _, srv := range health.Servers {
		c.Infof("my id: '%s', compare with: '%s'", c.cfg.Config.NodeID, srv.ID)
		if srv.ID == c.cfg.Config.NodeID {
			// SerfStatus: alive, left, failed or none:
			// https://github.com/hashicorp/consul/blob/c7ef04c5979dbc311ff3c67b7bf3028a93e8b0f1/agent/consul/operator_autopilot_endpoint.go#L124-L133
			mx["autopilot_server_sefStatus_alive"] = boolToInt(srv.SerfStatus == "alive")
			mx["autopilot_server_sefStatus_left"] = boolToInt(srv.SerfStatus == "left")
			mx["autopilot_server_sefStatus_failed"] = boolToInt(srv.SerfStatus == "failed")
			mx["autopilot_server_sefStatus_none"] = boolToInt(srv.SerfStatus == "none")
			mx["autopilot_server_healthy_yes"] = boolToInt(srv.Healthy)
			mx["autopilot_server_healthy_no"] = boolToInt(!srv.Healthy)
			mx["autopilot_server_voter_yes"] = boolToInt(srv.Voter)
			mx["autopilot_server_voter_no"] = boolToInt(!srv.Voter)
			if srv.Healthy {
				mx["autopilot_server_stable_time"] = int64(time.Now().Sub(srv.StableSince).Seconds())
			}
			if !srv.Leader {
				if v, err := time.ParseDuration(srv.LastContact); err == nil {
					mx["autopilot_server_lastContact_leader"] = v.Milliseconds()
				}
			}

			break
		}
	}

	return nil
}
