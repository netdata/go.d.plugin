// SPDX-License-Identifier: GPL-3.0-or-later

package mongo

import (
	"fmt"

	"github.com/netdata/go.d.plugin/agent/module"
)

// replSetCollect creates the map[string]int64 for the available dims.
// nil values will be ignored and not added to the map and thus metrics
// should not appear on the dashboard.
// if the querying node does not belong to a replica set
func (m *Mongo) collectReplSetStatus(ms map[string]int64) error {
	status, err := m.mongoCollector.replSetGetStatus()
	if err != nil {
		return fmt.Errorf("error get status of the replica set from mongo: %s", err)
	}

	var currentMembers []string
	for _, member := range status.Members {
		currentMembers = append(currentMembers, member.Name)
	}

	// replica nodes may be removed
	// we should collect metrics for these anymore
	m.removeReplicaSetMembers(currentMembers)
	m.replSetMembers = currentMembers

	for _, member := range status.Members {
		if member.LastHeartbeatRecv != nil {
			id := replicationHeartbeatLatencyDimPrefix + member.Name
			ms[id] = status.Date.Sub(*member.LastHeartbeatRecv).Milliseconds()

			if !m.replSetDimsEnabled[id] {
				m.replSetDimsEnabled[id] = true

				if chart := m.charts.Get(replicationHeartbeatLatency); chart != nil {
					if err := chart.AddDim(&module.Dim{ID: id, Name: member.Name}); err != nil {
						m.Warningf("failed to add dim: %v", err)
					} else {
						chart.MarkNotCreated()
					}
				}
			}
		}

		id := replicationLagDimPrefix + member.Name
		// Replica set time diff between current time and time when last entry from the oplog was applied
		ms[id] = status.Date.Sub(member.OptimeDate).Milliseconds()

		if !m.replSetDimsEnabled[id] {
			m.replSetDimsEnabled[id] = true

			if chart := m.charts.Get(replicationLag); chart != nil {
				if err := chart.AddDim(&module.Dim{ID: id, Name: member.Name}); err != nil {
					m.Warningf("failed to add dim: %v", err)
				} else {
					chart.MarkNotCreated()
				}
			}
		}

		if member.PingMs != nil {
			id := replicationNodePingDimPrefix + member.Name
			ms[id] = *member.PingMs

			if !m.replSetDimsEnabled[id] {
				m.replSetDimsEnabled[id] = true

				if chart := m.charts.Get(replicationNodePing); chart != nil {
					if err := chart.AddDim(&module.Dim{ID: id, Name: member.Name}); err != nil {
						m.Warningf("failed to add dim: %v", err)
					} else {
						chart.MarkNotCreated()
					}
				}
			}
		}
	}

	return nil
}

// removeReplicaSetMember removes dimensions for not existing
// replica set members
func (m *Mongo) removeReplicaSetMembers(newMembers []string) {
	diff := sliceDiff(m.replSetMembers, newMembers)
	for _, name := range diff {
		for _, v := range []struct{ chartID, dimPrefix string }{
			{replicationLag, replicationLagDimPrefix},
			{replicationHeartbeatLatency, replicationHeartbeatLatencyDimPrefix},
			{replicationNodePing, replicationNodePingDimPrefix},
		} {
			id := v.dimPrefix + name
			if !m.replSetDimsEnabled[id] {
				continue
			}
			delete(m.replSetDimsEnabled, id)

			chart := m.charts.Get(v.chartID)
			if chart == nil {
				m.Warningf("failed to remove dimension: %s. job doesn't have chart: %s", id, v.chartID)
				continue
			}

			err := chart.MarkDimRemove(id, true)
			if err != nil {
				m.Warningf("failed to remove dimension: %v", err)
				continue
			}
			chart.MarkNotCreated()
		}
	}
}
