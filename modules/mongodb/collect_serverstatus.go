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
func (m *Mongo) collectServerStatus(mx map[string]int64) error {
	s, err := m.conn.serverStatus()
	if err != nil {
		return fmt.Errorf("serverStatus command failed: %s", err)
	}

	m.addOptionalCharts(s)

	for k, v := range stm.ToMap(s) {
		mx[k] = v
	}

	if s.Transactions != nil && s.Transactions.CommitTypes != nil {
		px := "transactions_commit_types_"
		v := s.Transactions.CommitTypes
		mx[px+"no_shards_unsuccessful"] = v.NoShards.Initiated - v.NoShards.Successful
		mx[px+"single_shard_unsuccessful"] = v.SingleShard.Initiated - v.SingleShard.Successful
		mx[px+"single_write_shard_unsuccessful"] = v.SingleWriteShard.Initiated - v.SingleWriteShard.Successful
		mx[px+"read_only_unsuccessful"] = v.ReadOnly.Initiated - v.ReadOnly.Successful
		mx[px+"two_phase_commit_unsuccessful"] = v.TwoPhaseCommit.Initiated - v.TwoPhaseCommit.Successful
		mx[px+"recover_with_token_unsuccessful"] = v.RecoverWithToken.Initiated - v.RecoverWithToken.Successful
	}

	return nil
}

func (m *Mongo) addOptionalCharts(s *documentServerStatus) {
	m.addOptionalChart(s.OpLatencies,
		&chartOperations,
		&chartOperationsLatency,
	)
	m.addOptionalChart(s.WiredTiger,
		&chartWiredTigerConcurrentReadTransactions,
		&chartWiredTigerConcurrentWriteTransactions,
		&chartWiredTigerCacheUsage,
		&chartWiredTigerCacheDirtySpaceSize,
		&chartWiredTigerCacheIORate,
		&chartWiredTigerCacheEvictionsRate,
	)
	m.addOptionalChart(s.Tcmalloc,
		&chartMemoryTCMallocStatsChart,
	)
	m.addOptionalChart(s.GlobalLock,
		&chartGlobalLockActiveClients,
		&chartGlobalLockCurrentQueue,
	)
	m.addOptionalChart(s.Network.NumSlowDNSOperations,
		&chartNetworkSlowDNSResolutions,
	)
	m.addOptionalChart(s.Network.NumSlowSSLOperations,
		&chartNetworkSlowSSLHandshakes,
	)
	if s.Transactions != nil {
		m.addOptionalChart(s.Transactions,
			&chartTransactionsCurrent,
			&chartTransactionsRate,
		)
		m.addOptionalChart(s.Transactions.CommitTypes,
			&chartTransactionsNoShardsCommitsRate,
			&chartTransactionsNoShardsCommitsDuration,
			&chartTransactionsSingleShardCommitsRate,
			&chartTransactionsSingleShardCommitsDuration,
			&chartTransactionsSingleWriteShardCommitsRate,
			&chartTransactionsSingleWriteShardCommitsDuration,
			&chartTransactionsReadOnlyCommitsRate,
			&chartTransactionsReadOnlyCommitsDuration,
			&chartTransactionsTwoPhaseCommitCommitsRate,
			&chartTransactionsTwoPhaseCommitCommitsDuration,
			&chartTransactionsRecoverWithTokenCommitsRate,
			&chartTransactionsRecoverWithTokenCommitsDuration,
		)
	}
	if s.Locks != nil {
		m.addOptionalChart(s.Locks.Global, &chartGlobalLockAcquisitions)
		m.addOptionalChart(s.Locks.Database, &chartDatabaseLockAcquisitions)
		m.addOptionalChart(s.Locks.Collection, &chartCollectionLockAcquisitions)
		m.addOptionalChart(s.Locks.Mutex, &chartMutexLockAcquisitions)
		m.addOptionalChart(s.Locks.Metadata, &chartMetadataLockAcquisitions)
		m.addOptionalChart(s.Locks.Oplog, &chartOpLogLockAcquisitions)
	}
}

func (m *Mongo) addOptionalChart(iface any, charts ...*module.Chart) {
	if reflect.ValueOf(iface).IsNil() {
		return
	}
	for _, chart := range charts {
		if m.optionalCharts[chart.ID] {
			continue
		}
		m.optionalCharts[chart.ID] = true

		if err := m.charts.Add(chart.Copy()); err != nil {
			m.Warning(err)
		}
	}
}
