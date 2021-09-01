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
		// Heartbeat lag calculation
		if member.LastHeartbeatRecv != nil {
			id := replicationHeartbeatLatencyDimPrefix + member.Name
			// add dimension if not exists yet
			if !m.replSetDimsEnabled[id] {
				m.replSetDimsEnabled[id] = true
				chart := m.charts.Get(replicationHeartbeatLatency)
				if chart != nil {
					if err := chart.AddDim(&module.Dim{ID: id, Name: member.Name}); err != nil {
						m.Warningf("failed to add dim: %v", err)
					}
				}
			}
			ms[id] = status.Date.Sub(*member.LastHeartbeatRecv).Milliseconds()
		}

		// Replica set time diff between current time and time when last entry from the oplog was applied
		id := replicationLagDimPrefix + member.Name
		// add dimension if not exists yet
		if !m.replSetDimsEnabled[id] {
			m.replSetDimsEnabled[id] = true
			chart := m.charts.Get(replicationLag)
			if chart != nil {
				if err := chart.AddDim(&module.Dim{ID: id, Name: member.Name}); err != nil {
					m.Warningf("failed to add dim: %v", err)
				}
			}
		}
		ms[id] = status.Date.Sub(member.OptimeDate).Milliseconds()

		// Ping time
		if member.PingMs != nil {
			id := replicationNodePingDimPrefix + member.Name
			// add dimension if not exists yet
			if !m.replSetDimsEnabled[id] {
				m.replSetDimsEnabled[id] = true
				chart := m.charts.Get(replicationNodePing)
				if chart != nil {
					if err := chart.AddDim(&module.Dim{ID: id, Name: member.Name}); err != nil {
						m.Warningf("failed to add dim: %v", err)
					}
				}
			}
			ms[id] = *member.PingMs
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
			if chart != nil {
				if err := chart.MarkDimRemove(id, true); err != nil {
					m.Warningf("failed to remove dimension: %v", err)
				}
			} else {
				m.Warningf("failed to remove dimension:%s. job doesn't have chart: %s", id, v.chartID)
			}
		}
	}
}
