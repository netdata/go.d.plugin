package mongo

import (
	"reflect"

	"github.com/netdata/go.d.plugin/agent/module"
	"github.com/netdata/go.d.plugin/pkg/stm"
)

func (m *Mongo) serverStatusCollect() map[string]int64 {
	status, err := m.mongoCollector.serverStatus()
	if err != nil {
		m.Errorf("error get server status from mongo: %s", err)
		return nil
	}

	m.addOptionalCharts(status)

	var args = []interface{}{
		status.Opcounters,
		status.Connections,
		status.Network,
		status.ExtraInfo,
		status.Asserts,
	}

	if status.Transactions != nil {
		args = append(args, status.Transactions)
	}
	if status.FlowControl != nil {
		args = append(args, status.FlowControl)
	}
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
