package mongo

import (
	"fmt"
	"reflect"

	"github.com/netdata/go.d.plugin/agent/module"
	"github.com/netdata/go.d.plugin/pkg/stm"
)

// collectServerStatus creates the map[string]int64 for the available dims.
// nil values will be ignored and not added to the map and thus metrics
// should not appear on the dashboard.
// Because mongo reports a metric only after it first appears, some dims might
// take a while to appear. For example, in order to report number of create
// commands, a document must be created first.
func (m *Mongo) collectServerStatus(ms map[string]int64) error {
	status, err := m.mongoCollector.serverStatus()
	if err != nil {
		return fmt.Errorf("serverStatus command failed: %s", err)
	}

	m.addOptionalCharts(status)
	for k, v := range stm.ToMap(status) {
		ms[k] = v
	}
	return nil
}

// addOptionalCharts tries to add charts based on the availability of
// metrics coming back from the `serverStatus` command.
// Nil pointer structs will be skipped and we won't produce metrics for
// unavailable metrics.
func (m *Mongo) addOptionalCharts(status *serverStatus) {
	m.metricExists(status.Transactions, &chartTransactionsCurrent)
	m.metricExists(status.FlowControl, &chartFlowControl)

	if status.GlobalLock != nil {
		m.metricExists(status.GlobalLock.ActiveClients, &chartGlobalLockActiveClients)
		m.metricExists(status.GlobalLock.CurrentQueue, &chartGlobalLockCurrentQueue)
	}
	if status.Tcmalloc != nil {
		m.metricExists(status.Tcmalloc.Generic, &chartTcmallocGeneric)
		m.metricExists(status.Tcmalloc.Tcmalloc, &chartTcmalloc)
	}
	if status.Locks != nil {
		m.metricExists(status.Locks.Global, &chartLocks)
		m.metricExists(status.Locks.Database, &chartLocks)
		m.metricExists(status.Locks.Collection, &chartLocks)
	}
	if status.WiredTiger != nil {
		m.metricExists(status.WiredTiger.BlockManager, &chartWiredTigerBlockManager)
		m.metricExists(status.WiredTiger.Cache, &chartWiredTigerCache)
		m.metricExists(status.WiredTiger.Capacity, &chartWiredTigerCapacity)
		m.metricExists(status.WiredTiger.Connection, &chartWiredTigerConnection)
		m.metricExists(status.WiredTiger.Cursor, &chartWiredTigerCursor)
		m.metricExists(status.WiredTiger.Lock, &chartWiredTigerLock)
		m.metricExists(status.WiredTiger.Lock, &chartWiredTigerLockDuration)
		m.metricExists(status.WiredTiger.Log, &chartWiredTigerLogOps)
		m.metricExists(status.WiredTiger.Log, &chartWiredTigerLogBytes)
		m.metricExists(status.WiredTiger.Transaction, &chartWiredTigerTransactions)
	}
}

// metricExists checks if the paces interface(iface) is not nil
// and if so add the passed chart to the index.
func (m *Mongo) metricExists(iface interface{}, chart *module.Chart) {
	if reflect.ValueOf(iface).IsNil() {
		return
	}
	if !m.optionalChartsEnabled[chart.ID] {
		err := m.charts.Add(chart.Copy())
		if err != nil {
			m.Warning(err)
		}
		m.optionalChartsEnabled[chart.ID] = true
		return
	}
}

// collectDbStats creates the map[string]int64 for the available dims
// it calls listDatabase and then dbstats for each database internally
func (m *Mongo) collectDbStats(ms map[string]int64) error {
	if m.databasesMatcher == nil {
		return nil
	}
	allDatabases, err := m.mongoCollector.listDatabaseNames()
	if err != nil {
		return fmt.Errorf("cannot get database names: %s", err)
	}

	// filter matching databases and exclude non-matching
	var databases []string
	for _, database := range allDatabases {
		if m.databasesMatcher.MatchString(database) {
			databases = append(databases, database)
		}
	}

	// add dims for each database
	m.updateDBStatsCharts(databases)

	for _, database := range databases {
		stats, err := m.mongoCollector.dbStats(database)
		if err != nil {
			return fmt.Errorf("dbStats command failed: %s", err)
		}
		stats.toMap(database, ms)
	}
	return nil
}

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
		if member.LastHeartbeatReceived != nil {
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
			ms[id] = status.Date.Sub(*member.LastHeartbeatReceived).Milliseconds()
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
