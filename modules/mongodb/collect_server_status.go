// SPDX-License-Identifier: GPL-3.0-or-later

package mongo

import (
	"fmt"
	"reflect"

	"github.com/netdata/go.d.plugin/agent/module"
	"github.com/netdata/go.d.plugin/pkg/stm"
)

// collectServerStatus creates the map[string]int64 for the available dims.
// nil values will be ignored and not added to the map and thus metrics should not appear on the dashboard.
// Because mongo reports a metric only after it first appears,some dims might take a while to appear.
// For example, in order to report number of create commands, a document must be created first.
func (m *Mongo) collectServerStatus(ms map[string]int64) error {
	status, err := m.conn.serverStatus()
	if err != nil {
		return fmt.Errorf("serverStatus command failed: %s", err)
	}

	m.addOptionalCharts(status)

	for k, v := range stm.ToMap(status) {
		ms[k] = v
	}

	return nil
}

func (m *Mongo) addOptionalCharts(s *serverStatus) {
	m.addOptionalChart(s.FlowControl, &flowControlTimingsChart)

	if s.Transactions != nil {
		m.addOptionalChart(s.Transactions, &transactionsCurrentChart)
		if m.conn.isMongos() {
			m.addOptionalChart(s.Transactions.CommitTypes, &transactionsCommitTypesChart)
		}
	}

	if s.GlobalLock != nil {
		m.addOptionalChart(s.GlobalLock.ActiveClients, &globalLockActiveClientsChart)
		m.addOptionalChart(s.GlobalLock.CurrentQueue, &globalLockCurrentQueueChart)
	}
	if s.Tcmalloc != nil {
		m.addOptionalChart(s.Tcmalloc.Generic, &tcMallocGenericChart)
		m.addOptionalChart(s.Tcmalloc.Tcmalloc, &tcMallocChart)
	}
	if s.Locks != nil {
		m.addOptionalChart(s.Locks.Global, &locksChart)
		m.addOptionalChart(s.Locks.Database, &locksChart)
		m.addOptionalChart(s.Locks.Collection, &locksChart)
	}
	if s.WiredTiger != nil {
		m.addOptionalChart(s.WiredTiger.BlockManager, &wiredTigerBlocksChart)
		m.addOptionalChart(s.WiredTiger.Cache, &wiredTigerCacheChart)
		m.addOptionalChart(s.WiredTiger.Capacity, &chartWiredTigerCapacity)
		m.addOptionalChart(s.WiredTiger.Connection, &chartWiredTigerConnection)
		m.addOptionalChart(s.WiredTiger.Cursor, &chartWiredTigerCursor)
		m.addOptionalChart(s.WiredTiger.Lock, &chartWiredTigerLock)
		m.addOptionalChart(s.WiredTiger.Lock, &chartWiredTigerLockDuration)
		m.addOptionalChart(s.WiredTiger.Log, &chartWiredTigerLogOps)
		m.addOptionalChart(s.WiredTiger.Log, &chartWiredTigerLogBytes)
		m.addOptionalChart(s.WiredTiger.Transaction, &chartWiredTigerTransactions)
	}
}

func (m *Mongo) addOptionalChart(iface interface{}, chart *module.Chart) {
	if reflect.ValueOf(iface).IsNil() || m.optionalCharts[chart.ID] {
		return
	}

	m.optionalCharts[chart.ID] = true

	if err := m.charts.Add(chart.Copy()); err != nil {
		m.Warning(err)
	}
}
