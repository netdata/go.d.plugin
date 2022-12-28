// SPDX-License-Identifier: GPL-3.0-or-later

package mongo

import (
	"fmt"
)

// https://www.mongodb.com/docs/manual/reference/replica-states/#replica-set-member-states
var replicaSetMemberStates = map[string]int{
	"startup":    0,
	"primary":    1,
	"secondary":  2,
	"recovering": 3,
	"startup2":   5,
	"unknown":    6,
	"arbiter":    7,
	"down":       8,
	"rollback":   9,
	"removed":    10,
}

// TODO: deal with duplicates if we collect metrics from all cluster nodes
// should we only collect ReplSetStatus (at least by default) from primary nodes? (db.runCommand( { isMaster: 1 } ))
func (m *Mongo) collectReplSetStatus(mx map[string]int64) error {
	s, err := m.conn.replSetGetStatus()
	if err != nil {
		return fmt.Errorf("error get status of the replica set from mongo: %s", err)
	}

	// https://www.mongodb.com/docs/manual/reference/command/replSetGetStatus/

	seen := make(map[string]replSetMember)

	for _, member := range s.Members {
		seen[member.Name] = member

		px := fmt.Sprintf("repl_set_member_%s_", member.Name)

		mx[px+"replication_lag"] = s.Date.Sub(member.OptimeDate).Milliseconds()

		for k, v := range replicaSetMemberStates {
			mx[px+"state_"+k] = boolToInt(member.State == v)
		}

		mx[px+"health_status_up"] = boolToInt(member.Health == 1)
		mx[px+"health_status_down"] = boolToInt(member.Health == 0)

		if member.Self == nil {
			mx[px+"uptime"] = member.Uptime
			if v := member.LastHeartbeatRecv; v != nil && !v.IsZero() {
				mx[px+"heartbeat_latency"] = s.Date.Sub(*v).Milliseconds()
			}
			if v := member.PingMs; v != nil {
				mx[px+"ping_rtt"] = *v
			}
		}
	}

	for name, member := range seen {
		if !m.replSetMembers[name] {
			m.replSetMembers[name] = true
			m.Debugf("new replica set member '%s': adding charts", name)
			m.addReplSetMemberCharts(member)
		}
	}

	for name := range m.replSetMembers {
		if _, ok := seen[name]; !ok {
			delete(m.replSetMembers, name)
			m.Debugf("stale replica set member '%s': removing charts", name)
			m.removeReplSetMemberCharts(name)
		}
	}

	return nil
}
