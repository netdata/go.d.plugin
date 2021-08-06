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
	ms := stm.ToMap(status)
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
