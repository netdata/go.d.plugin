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

func (m *Mongo) addOptionalCharts(s *documentServerStatus) {
	m.addOptionalChart(s.FlowControl, &chartFlowControlTimings)

	if s.Transactions != nil {
		m.addOptionalChart(s.Transactions, &chartTransactionsCurrent)
		if m.conn.isMongos() {
			m.addOptionalChart(s.Transactions.CommitTypes, &chartTransactionsCommitTypes)
		}
	}

	if s.GlobalLock != nil {
		m.addOptionalChart(s.GlobalLock.ActiveClients, &chartGlobalLockActiveClients)
		m.addOptionalChart(s.GlobalLock.CurrentQueue, &chartGlobalLockCurrentQueue)
	}
	if s.Tcmalloc != nil {
		m.addOptionalChart(s.Tcmalloc.Tcmalloc, &chartMemoryTCMallocStatsChart)
	}
	if s.Locks != nil {
		m.addOptionalChart(s.Locks.Global, &chartLocks)
		m.addOptionalChart(s.Locks.Database, &chartLocks)
		m.addOptionalChart(s.Locks.Collection, &chartLocks)
	}
	if s.WiredTiger != nil {
		m.addOptionalChart(s.WiredTiger, &chartWiredTigerConcurrentReadTransactions)
		m.addOptionalChart(s.WiredTiger, &chartWiredTigerConcurrentWriteTransactions)
		m.addOptionalChart(s.WiredTiger, &chartWiredTigerCacheUsage)
		m.addOptionalChart(s.WiredTiger, &chartWiredTigerCacheDirtySpaceSize)
		m.addOptionalChart(s.WiredTiger, &chartWiredTigerCacheIORate)
		m.addOptionalChart(s.WiredTiger, &chartWiredTigerCacheEvictionRate)
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
