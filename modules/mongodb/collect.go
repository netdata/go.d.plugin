package mongo

import (
	"reflect"

	"github.com/netdata/go.d.plugin/agent/module"
	"github.com/netdata/go.d.plugin/pkg/stm"
)

// serverStatusCollect creates the map[string]int64 for the available dims.
// nil values will be ignored and not added to the map and thus metrics
// should not appear on the dashboard.
// Because mongo reports a metric only after it first appears, some dims might
// take a while to appear. For example, in order to report number of create
// commands, a document must be created first.
func (m *Mongo) serverStatusCollect() map[string]int64 {
	status, err := m.mongoCollector.serverStatus()
	if err != nil {
		m.Errorf("error get server status from mongo: %s", err)
		return nil
	}

	m.addOptionalCharts(status)

	// values for the field below are expected to be present
	// in wide range of mongo versions and builds
	// please refer to https://docs.mongodb.com/manual/reference/command/serverStatus/
	// for version specific availability and changes
	var args = []interface{}{
		status.Opcounters,
		status.Connections,
		status.Network,
		status.ExtraInfo,
		status.Asserts,
	}

	// Available on mongod in 3.6.3+
	if status.Transactions != nil {
		args = append(args, status.Transactions)
	}
	// New in mongo version 4.2
	if status.FlowControl != nil {
		args = append(args, status.FlowControl)
	}
	// Only for `mongod` instances
	if status.OpLatencies != nil {
		args = append(args, status.OpLatencies.Reads)
		args = append(args, status.OpLatencies.Writes)
		args = append(args, status.OpLatencies.Commands)
	}
	if status.GlobalLock != nil {
		args = append(args, status.GlobalLock.ActiveClients)
		args = append(args, status.GlobalLock.CurrentQueue)
	}
	if status.Tcmalloc != nil {
		args = append(args, status.Tcmalloc.Generic)
		args = append(args, status.Tcmalloc.Tcmalloc)
	}
	if status.Locks != nil {
		args = append(args, status.Locks.Global)
		args = append(args, status.Locks.Database)
		args = append(args, status.Locks.Collection)
	}
	// available when WiredTiger is used as the storage engine
	// https://docs.mongodb.com/manual/reference/command/serverStatus/#wiredtiger
	if status.WiredTiger != nil {
		args = append(args, status.WiredTiger.BlockManager)
		args = append(args, status.WiredTiger.Cache)
		args = append(args, status.WiredTiger.Capacity)
		args = append(args, status.WiredTiger.Connection)
		args = append(args, status.WiredTiger.Cursor)
		args = append(args, status.WiredTiger.Lock)
		args = append(args, status.WiredTiger.Log)
		args = append(args, status.WiredTiger.Transaction)
	}

	ms := stm.ToMap(args...)
	return ms
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
