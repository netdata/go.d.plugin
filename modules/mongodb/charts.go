// SPDX-License-Identifier: GPL-3.0-or-later

package mongo

import (
	"github.com/netdata/go.d.plugin/agent/module"
)

const (
	prioOperations = module.Priority + iota
	prioOperationsLatency
	prioOperationsByType
	prioDocumentOperations

	prioActiveClients
	prioQueuedOperations

	prioGlobalLockAcquisitions
	prioDatabaseLockAcquisitions
	prioCollectionLockAcquisitions
	prioMutexLockAcquisitions
	prioMetadataLockAcquisitions
	prioOpLogLockAcquisitions

	prioTransactionsCurrent
	prioTransactionsRate
	prioTransactionsNoShardsCommitsRate
	prioTransactionsNoShardsCommitsDuration
	prioTransactionsSingleShardCommitsRate
	prioTransactionsSingleShardCommitsDuration
	prioTransactionsSingleWriteShardCommitsRate
	prioTransactionsSingleWriteShardCommitsDuration
	prioTransactionsReadOnlyCommitsRate
	prioTransactionsReadOnlyCommitsDuration
	prioTransactionsTwoPhaseCommitCommitsRate
	prioTransactionsTwoPhaseCommitCommitsDuration
	prioTransactionsRecoverWithTokenCommitsRate
	prioTransactionsRecoverWithTokenCommitsDuration

	prioConnectionsUsage
	prioConnectionsByState
	prioConnectionsRate

	prioAsserts

	prioNetworkTraffic
	prioNetworkRequests
	prioNetworkSlowDNSResolutions
	prioNetworkSlowSSLHandshakes

	prioMemoryResidentSize
	prioMemoryVirtualSize
	prioMemoryPageFaults
	prioMemoryTCMallocStats

	prioWiredTigerConcurrentReadTransactionsUsage
	prioWiredTigerConcurrentWriteTransactionsUsage
	prioWiredTigerCacheUsage
	prioWiredTigerCacheDirtySpaceSize
	prioWiredTigerCacheIORate
	prioWiredTigerCacheEvictionsRate

	prioDatabaseCollections
	prioDatabaseIndexes
	prioDatabaseViews
	prioDatabaseDocuments
	prioDatabaseDataSize
	prioDatabaseStorageSize
	prioDatabaseIndexSize

	prioReplSetMemberState
	prioReplSetMemberHealthStatus
	prioReplSetMemberReplicationLag
	prioReplSetMemberHeartbeatLatency
	prioReplSetMemberPingRTT
	prioReplSetMemberUptime

	prioShardingNodesCount
	prioShardingShardedDatabasesCount
	prioShardingShardedCollectionsCount
	prioShardChunks
)

// these charts are expected to be available in many versions
var chartsServerStatus = module.Charts{
	chartOperationsByType.Copy(),
	chartDocumentOperations.Copy(),

	chartConnectionsUsage.Copy(),
	chartConnectionsByState.Copy(),
	chartConnectionsRate.Copy(),

	chartNetworkTraffic.Copy(),
	chartNetworkRequests.Copy(),

	chartMemoryResident.Copy(),
	chartMemoryVirtual.Copy(),
	chartMemoryPageFaults.Copy(),

	chartAsserts.Copy(),
}

var chartsTmplDatabase = module.Charts{
	chartTmplDatabaseCollections.Copy(),
	chartTmplDatabaseIndexes.Copy(),
	chartTmplDatabaseViews.Copy(),
	chartTmplDatabaseDocuments.Copy(),
	chartTmplDatabaseDataSize.Copy(),
	chartTmplDatabaseStorageSize.Copy(),
	chartTmplDatabaseIndexSize.Copy(),
}

var chartsTmplReplSetMember = module.Charts{
	chartTmplReplSetMemberState.Copy(),
	chartTmplReplSetMemberHealthStatus.Copy(),
	chartTmplReplSetMemberReplicationLag.Copy(),
	chartTmplReplSetMemberHeartbeatLatency.Copy(),
	chartTmplReplSetMemberPingRTT.Copy(),
	chartTmplReplSetMemberUptime.Copy(),
}

var chartsSharding = module.Charts{
	chartShardingNodesCount.Copy(),
	chartShardingShardedDatabases.Copy(),
	chartShardingShardedCollectionsCount.Copy(),
}

var chartsTmplShardingShard = module.Charts{
	chartTmplShardChunks.Copy(),
}

var (
	chartOperations = module.Chart{
		ID:       "operations",
		Title:    "Operations",
		Units:    "operations/s",
		Fam:      "operations",
		Ctx:      "mongodb.operations",
		Priority: prioOperations,
		Dims: module.Dims{
			{ID: "operations_latencies_reads_ops", Name: "reads", Algo: module.Incremental},
			{ID: "operations_latencies_writes_ops", Name: "writes", Algo: module.Incremental},
			{ID: "operations_latencies_commands_ops", Name: "commands", Algo: module.Incremental},
		},
	}
	chartOperationsLatency = module.Chart{
		ID:       "operations_latency",
		Title:    "Operations Latency",
		Units:    "milliseconds",
		Fam:      "operations",
		Ctx:      "mongodb.operations_latency",
		Priority: prioOperationsLatency,
		Dims: module.Dims{
			{ID: "operations_latencies_reads_latency", Name: "reads", Algo: module.Incremental, Div: 1000},
			{ID: "operations_latencies_writes_latency", Name: "writes", Algo: module.Incremental, Div: 1000},
			{ID: "operations_latencies_commands_latency", Name: "commands", Algo: module.Incremental, Div: 1000},
		},
	}
	chartOperationsByType = module.Chart{
		ID:       "operations_by_type",
		Title:    "Operations by type",
		Units:    "operations/s",
		Fam:      "operations",
		Ctx:      "mongodb.operations_by_type",
		Priority: prioOperationsByType,
		Dims: module.Dims{
			{ID: "operations_insert", Name: "insert", Algo: module.Incremental},
			{ID: "operations_query", Name: "query", Algo: module.Incremental},
			{ID: "operations_update", Name: "update", Algo: module.Incremental},
			{ID: "operations_delete", Name: "delete", Algo: module.Incremental},
			{ID: "operations_getmore", Name: "getmore", Algo: module.Incremental},
			{ID: "operations_command", Name: "command", Algo: module.Incremental},
		},
	}
	chartDocumentOperations = module.Chart{
		ID:       "document_operations",
		Title:    "Document operations",
		Units:    "operations/s",
		Fam:      "operations",
		Ctx:      "mongodb.document_operations",
		Type:     module.Stacked,
		Priority: prioDocumentOperations,
		Dims: module.Dims{
			{ID: "metrics_document_inserted", Name: "inserted", Algo: module.Incremental},
			{ID: "metrics_document_deleted", Name: "deleted", Algo: module.Incremental},
			{ID: "metrics_document_returned", Name: "returned", Algo: module.Incremental},
			{ID: "metrics_document_updated", Name: "updated", Algo: module.Incremental},
		},
	}

	chartGlobalLockActiveClients = module.Chart{
		ID:       "active_clients",
		Title:    "Connection clients",
		Units:    "clients",
		Fam:      "clients",
		Ctx:      "mongodb.active_clients",
		Priority: prioActiveClients,
		Dims: module.Dims{
			{ID: "global_lock_active_clients_readers", Name: "readers"},
			{ID: "global_lock_active_clients_writers", Name: "writers"},
		},
	}
	chartGlobalLockCurrentQueue = module.Chart{
		ID:       "queued_operations",
		Title:    "Queued operations because of a lock",
		Units:    "operations",
		Fam:      "clients",
		Ctx:      "mongodb.queued_operations",
		Priority: prioQueuedOperations,
		Dims: module.Dims{
			{ID: "global_lock_current_queue_readers", Name: "readers"},
			{ID: "global_lock_current_queue_writers", Name: "writers"},
		},
	}

	chartConnectionsUsage = module.Chart{
		ID:       "connections_usage",
		Title:    "Connections usage",
		Units:    "connections",
		Fam:      "connections",
		Ctx:      "mongodb.connections_usage",
		Type:     module.Stacked,
		Priority: prioConnectionsUsage,
		Dims: module.Dims{
			{ID: "connections_available", Name: "available"},
			{ID: "connections_current", Name: "used"},
		},
	}
	chartConnectionsByState = module.Chart{
		ID:       "connections_by_state",
		Title:    "Connections By State",
		Units:    "connections",
		Fam:      "connections",
		Ctx:      "mongodb.connections_by_state",
		Priority: prioConnectionsByState,
		Dims: module.Dims{
			{ID: "connections_active", Name: "active"},
			{ID: "connections_threaded", Name: "threaded"},
			{ID: "connections_exhaust_is_master", Name: "exhaust_is_master"},
			{ID: "connections_exhaust_hello", Name: "exhaust_hello"},
			{ID: "connections_awaiting_topology_changes", Name: "awaiting_topology_changes"},
		},
	}
	chartConnectionsRate = module.Chart{
		ID:       "connections_rate",
		Title:    "Connections Rate",
		Units:    "connections/s",
		Fam:      "connections",
		Ctx:      "mongodb.connections_rate",
		Priority: prioConnectionsRate,
		Dims: module.Dims{
			{ID: "connections_total_created", Name: "created", Algo: module.Incremental},
		},
	}

	chartNetworkTraffic = module.Chart{
		ID:       "network_traffic",
		Title:    "Network traffic",
		Units:    "bytes/s",
		Fam:      "network",
		Ctx:      "mongodb.network_traffic",
		Priority: prioNetworkTraffic,
		Type:     module.Area,
		Dims: module.Dims{
			{ID: "network_bytes_in", Name: "in", Algo: module.Incremental},
			{ID: "network_bytes_out", Name: "out", Algo: module.Incremental},
		},
	}
	chartNetworkRequests = module.Chart{
		ID:       "network_requests",
		Title:    "Network Requests",
		Units:    "requests/s",
		Fam:      "network",
		Ctx:      "mongodb.network_requests",
		Priority: prioNetworkRequests,
		Dims: module.Dims{
			{ID: "network_requests", Name: "requests", Algo: module.Incremental},
		},
	}
	chartNetworkSlowDNSResolutions = module.Chart{
		ID:       "network_slow_dns_resolutions",
		Title:    "Slow DNS resolution operations",
		Units:    "resolutions/s",
		Fam:      "network",
		Ctx:      "mongodb.network_slow_dns_resolutions",
		Priority: prioNetworkSlowDNSResolutions,
		Dims: module.Dims{
			{ID: "network_slow_dns_operations", Name: "slow_dns", Algo: module.Incremental},
		},
	}
	chartNetworkSlowSSLHandshakes = module.Chart{
		ID:       "network_slow_ssl_handshakes",
		Title:    "Slow SSL handshake operations",
		Units:    "handshakes/s",
		Fam:      "network",
		Ctx:      "mongodb.network_slow_ssl_handshakes",
		Priority: prioNetworkSlowSSLHandshakes,
		Dims: module.Dims{
			{ID: "network_slow_ssl_operations", Name: "slow_ssl", Algo: module.Incremental},
		},
	}

	chartMemoryResident = module.Chart{
		ID:       "memory_resident_size",
		Title:    "Used resident memory",
		Units:    "bytes",
		Fam:      "memory",
		Ctx:      "mongodb.memory_resident_size",
		Priority: prioMemoryResidentSize,
		Dims: module.Dims{
			{ID: "memory_resident", Name: "used"},
		},
	}
	chartMemoryVirtual = module.Chart{
		ID:       "memory_virtual_size",
		Title:    "Used virtual memory",
		Units:    "bytes",
		Fam:      "memory",
		Ctx:      "mongodb.memory_virtual_size",
		Priority: prioMemoryVirtualSize,
		Dims: module.Dims{
			{ID: "memory_virtual", Name: "used"},
		},
	}
	chartMemoryPageFaults = module.Chart{
		ID:       "memory_page_faults",
		Title:    "Memory page faults",
		Units:    "pgfaults/s",
		Fam:      "memory",
		Ctx:      "mongodb.memory_page_faults",
		Priority: prioMemoryPageFaults,
		Dims: module.Dims{
			{ID: "extra_info_page_faults", Name: "pgfaults", Algo: module.Incremental},
		},
	}
	chartMemoryTCMallocStatsChart = module.Chart{
		ID:       "memory_tcmalloc_stats",
		Title:    "TCMalloc statistics",
		Units:    "bytes",
		Fam:      "memory",
		Ctx:      "mongodb.memory_tcmalloc_stats",
		Priority: prioMemoryTCMallocStats,
		Dims: module.Dims{
			{ID: "tcmalloc_generic_current_allocated_bytes", Name: "allocated"},
			{ID: "tcmalloc_central_cache_free_bytes", Name: "central_cache_freelist"},
			{ID: "tcmalloc_transfer_cache_free_bytes", Name: "transfer_cache_freelist"},
			{ID: "tcmalloc_thread_cache_free_bytes", Name: "thread_cache_freelists"},
			{ID: "tcmalloc_pageheap_free_bytes", Name: "pageheap_freelist"},
			{ID: "tcmalloc_pageheap_unmapped_bytes", Name: "pageheap_unmapped"},
		},
	}

	chartAsserts = module.Chart{
		ID:       "asserts",
		Title:    "Raised assertions",
		Units:    "asserts/s",
		Fam:      "asserts",
		Ctx:      "mongodb.asserts",
		Type:     module.Stacked,
		Priority: prioAsserts,
		Dims: module.Dims{
			{ID: "asserts_regular", Name: "regular", Algo: module.Incremental},
			{ID: "asserts_warning", Name: "warning", Algo: module.Incremental},
			{ID: "asserts_msg", Name: "msg", Algo: module.Incremental},
			{ID: "asserts_user", Name: "user", Algo: module.Incremental},
			{ID: "asserts_tripwire", Name: "tripwire", Algo: module.Incremental},
			{ID: "asserts_rollovers", Name: "rollovers", Algo: module.Incremental},
		},
	}

	chartTransactionsCurrent = module.Chart{
		ID:       "current_transactions",
		Title:    "Current transactions",
		Units:    "transactions",
		Fam:      "transactions",
		Ctx:      "mongodb.current_transactions",
		Priority: prioTransactionsCurrent,
		Dims: module.Dims{
			{ID: "transactions_active", Name: "active"},
			{ID: "transactions_inactive", Name: "inactive"},
			{ID: "transactions_open", Name: "open"},
			{ID: "transactions_prepared", Name: "prepared"},
		},
	}
	chartTransactionsRate = module.Chart{
		ID:       "transactions_rate",
		Title:    "Transactions rate",
		Units:    "transactions/s",
		Fam:      "transactions",
		Ctx:      "mongodb.transactions_rate",
		Priority: prioTransactionsRate,
		Dims: module.Dims{
			{ID: "transactions_total_started", Name: "started", Algo: module.Incremental},
			{ID: "transactions_total_aborted", Name: "aborted", Algo: module.Incremental},
			{ID: "transactions_total_committed", Name: "committed", Algo: module.Incremental},
			{ID: "transactions_total_prepared", Name: "prepared", Algo: module.Incremental},
		},
	}
	chartTransactionsNoShardsCommitsRate = module.Chart{
		ID:       "transactions_no_shards_commits_rate",
		Title:    "Transactions commits",
		Units:    "commits/s",
		Fam:      "transactions",
		Ctx:      "mongodb.transactions_commits_rate",
		Priority: prioTransactionsNoShardsCommitsRate,
		Type:     module.Stacked,
		Labels:   []module.Label{{Key: "commit_type", Value: "noShards"}},
		Dims: module.Dims{
			{ID: "transactions_commit_types_no_shards_successful", Name: "success", Algo: module.Incremental},
			{ID: "transactions_commit_types_no_shards_unsuccessful", Name: "fail", Algo: module.Incremental},
		},
	}
	chartTransactionsNoShardsCommitsDuration = module.Chart{
		ID:       "transactions_no_shards_commits_duration",
		Title:    "Transactions successful commits duration",
		Units:    "milliseconds",
		Fam:      "transactions",
		Ctx:      "mongodb.transactions_commits_duration",
		Priority: prioTransactionsNoShardsCommitsDuration,
		Labels:   []module.Label{{Key: "commit_type", Value: "noShards"}},
		Dims: module.Dims{
			{ID: "transactions_commit_types_no_shards_successful_duration_micros", Name: "commits", Algo: module.Incremental, Div: 1000},
		},
	}
	chartTransactionsSingleShardCommitsRate = module.Chart{
		ID:       "transactions_single_shard_commits_rate",
		Title:    "Transactions commits",
		Units:    "commits/s",
		Fam:      "transactions",
		Ctx:      "mongodb.transactions_commits_rate",
		Priority: prioTransactionsSingleShardCommitsRate,
		Type:     module.Stacked,
		Labels:   []module.Label{{Key: "commit_type", Value: "singleShard"}},
		Dims: module.Dims{
			{ID: "transactions_commit_types_single_shard_successful", Name: "success", Algo: module.Incremental},
			{ID: "transactions_commit_types_single_shard_unsuccessful", Name: "fail", Algo: module.Incremental},
		},
	}
	chartTransactionsSingleShardCommitsDuration = module.Chart{
		ID:       "transactions_single_shard_commits_duration",
		Title:    "Transactions successful commits duration",
		Units:    "milliseconds",
		Fam:      "transactions",
		Ctx:      "mongodb.transactions_commits_duration",
		Priority: prioTransactionsSingleShardCommitsDuration,
		Labels:   []module.Label{{Key: "commit_type", Value: "singleShard"}},
		Dims: module.Dims{
			{ID: "transactions_commit_types_single_shard_successful_duration_micros", Name: "commits", Algo: module.Incremental, Div: 1000},
		},
	}
	chartTransactionsSingleWriteShardCommitsRate = module.Chart{
		ID:       "transactions_single_write_shard_commits_rate",
		Title:    "Transactions commits",
		Units:    "commits/s",
		Fam:      "transactions",
		Ctx:      "mongodb.transactions_commits_rate",
		Priority: prioTransactionsSingleWriteShardCommitsRate,
		Type:     module.Stacked,
		Labels:   []module.Label{{Key: "commit_type", Value: "singleWriteShard"}},
		Dims: module.Dims{
			{ID: "transactions_commit_types_single_write_shard_successful", Name: "success", Algo: module.Incremental},
			{ID: "transactions_commit_types_single_write_shard_unsuccessful", Name: "fail", Algo: module.Incremental},
		},
	}
	chartTransactionsSingleWriteShardCommitsDuration = module.Chart{
		ID:       "transactions_single_write_shard_commits_duration",
		Title:    "Transactions successful commits duration",
		Units:    "milliseconds",
		Fam:      "transactions",
		Ctx:      "mongodb.transactions_commits_duration",
		Priority: prioTransactionsSingleWriteShardCommitsDuration,
		Labels:   []module.Label{{Key: "commit_type", Value: "singleWriteShard"}},
		Dims: module.Dims{
			{ID: "transactions_commit_types_single_write_shard_successful_duration_micros", Name: "commits", Algo: module.Incremental, Div: 1000},
		},
	}
	chartTransactionsReadOnlyCommitsRate = module.Chart{
		ID:       "transactions_read_only_commits_rate",
		Title:    "Transactions commits",
		Units:    "commits/s",
		Fam:      "transactions",
		Ctx:      "mongodb.transactions_commits_rate",
		Priority: prioTransactionsReadOnlyCommitsRate,
		Type:     module.Stacked,
		Labels:   []module.Label{{Key: "commit_type", Value: "readOnly"}},
		Dims: module.Dims{
			{ID: "transactions_commit_types_read_only_successful", Name: "success", Algo: module.Incremental},
			{ID: "transactions_commit_types_read_only_unsuccessful", Name: "fail", Algo: module.Incremental},
		},
	}
	chartTransactionsReadOnlyCommitsDuration = module.Chart{
		ID:       "transactions_read_only_commits_duration",
		Title:    "Transactions successful commits duration",
		Units:    "milliseconds",
		Fam:      "transactions",
		Ctx:      "mongodb.transactions_commits_duration",
		Priority: prioTransactionsReadOnlyCommitsDuration,
		Labels:   []module.Label{{Key: "commit_type", Value: "readOnly"}},
		Dims: module.Dims{
			{ID: "transactions_commit_types_read_only_successful_duration_micros", Name: "commits", Algo: module.Incremental, Div: 1000},
		},
	}
	chartTransactionsTwoPhaseCommitCommitsRate = module.Chart{
		ID:       "transactions_two_phase_commit_commits_rate",
		Title:    "Transactions commits",
		Units:    "commits/s",
		Fam:      "transactions",
		Ctx:      "mongodb.transactions_commits_rate",
		Priority: prioTransactionsTwoPhaseCommitCommitsRate,
		Type:     module.Stacked,
		Labels:   []module.Label{{Key: "commit_type", Value: "twoPhaseCommit"}},
		Dims: module.Dims{
			{ID: "transactions_commit_types_two_phase_commit_successful", Name: "success", Algo: module.Incremental},
			{ID: "transactions_commit_types_two_phase_commit_unsuccessful", Name: "fail", Algo: module.Incremental},
		},
	}
	chartTransactionsTwoPhaseCommitCommitsDuration = module.Chart{
		ID:       "transactions_two_phase_commit_commits_duration",
		Title:    "Transactions successful commits duration",
		Units:    "milliseconds",
		Fam:      "transactions",
		Ctx:      "mongodb.transactions_commits_duration",
		Priority: prioTransactionsTwoPhaseCommitCommitsDuration,
		Labels:   []module.Label{{Key: "commit_type", Value: "twoPhaseCommit"}},
		Dims: module.Dims{
			{ID: "transactions_commit_types_two_phase_commit_successful_duration_micros", Name: "commits", Algo: module.Incremental, Div: 1000},
		},
	}
	chartTransactionsRecoverWithTokenCommitsRate = module.Chart{
		ID:       "transactions_recover_with_token_commits_rate",
		Title:    "Transactions commits",
		Units:    "commits/s",
		Fam:      "transactions",
		Ctx:      "mongodb.transactions_commits_rate",
		Priority: prioTransactionsRecoverWithTokenCommitsRate,
		Type:     module.Stacked,
		Labels:   []module.Label{{Key: "commit_type", Value: "recoverWithToken"}},
		Dims: module.Dims{
			{ID: "transactions_commit_types_recover_with_token_successful", Name: "success", Algo: module.Incremental},
			{ID: "transactions_commit_types_recover_with_token_unsuccessful", Name: "fail", Algo: module.Incremental},
		},
	}
	chartTransactionsRecoverWithTokenCommitsDuration = module.Chart{
		ID:       "transactions_recover_with_token_commits_duration",
		Title:    "Transactions successful commits duration",
		Units:    "milliseconds",
		Fam:      "transactions",
		Ctx:      "mongodb.transactions_commits_duration",
		Priority: prioTransactionsRecoverWithTokenCommitsDuration,
		Labels:   []module.Label{{Key: "commit_type", Value: "recoverWithToken"}},
		Dims: module.Dims{
			{ID: "transactions_commit_types_recover_with_token_successful_duration_micros", Name: "commits", Algo: module.Incremental, Div: 1000},
		},
	}

	chartGlobalLockAcquisitions = module.Chart{
		ID:       "global_lock_acquisitions",
		Title:    "Global lock acquisitions",
		Units:    "acquisitions/s",
		Fam:      "locks",
		Ctx:      "mongodb.lock_acquisitions",
		Priority: prioGlobalLockAcquisitions,
		Labels:   []module.Label{{Key: "lock_type", Value: "global"}},
		Dims: module.Dims{
			{ID: "locks_global_acquire_shared", Name: "shared", Algo: module.Incremental},
			{ID: "locks_global_acquire_exclusive", Name: "exclusive", Algo: module.Incremental},
			{ID: "locks_global_acquire_intent_shared", Name: "intent_shared", Algo: module.Incremental},
			{ID: "locks_global_acquire_intent_exclusive", Name: "intent_exclusive", Algo: module.Incremental},
		},
	}
	chartDatabaseLockAcquisitions = module.Chart{
		ID:       "database_lock_acquisitions",
		Title:    "Database lock acquisitions",
		Units:    "acquisitions/s",
		Fam:      "locks",
		Ctx:      "mongodb.lock_acquisitions",
		Priority: prioDatabaseLockAcquisitions,
		Labels:   []module.Label{{Key: "lock_type", Value: "database"}},
		Dims: module.Dims{
			{ID: "locks_database_acquire_shared", Name: "shared", Algo: module.Incremental},
			{ID: "locks_database_acquire_exclusive", Name: "exclusive", Algo: module.Incremental},
			{ID: "locks_database_acquire_intent_shared", Name: "intent_shared", Algo: module.Incremental},
			{ID: "locks_database_acquire_intent_exclusive", Name: "intent_exclusive", Algo: module.Incremental},
		},
	}
	chartCollectionLockAcquisitions = module.Chart{
		ID:       "collection_lock_acquisitions",
		Title:    "Collection lock acquisitions",
		Units:    "acquisitions/s",
		Fam:      "locks",
		Ctx:      "mongodb.lock_acquisitions",
		Priority: prioCollectionLockAcquisitions,
		Labels:   []module.Label{{Key: "lock_type", Value: "collection"}},
		Dims: module.Dims{
			{ID: "locks_collection_acquire_shared", Name: "shared", Algo: module.Incremental},
			{ID: "locks_collection_acquire_exclusive", Name: "exclusive", Algo: module.Incremental},
			{ID: "locks_collection_acquire_intent_shared", Name: "intent_shared", Algo: module.Incremental},
			{ID: "locks_collection_acquire_intent_exclusive", Name: "intent_exclusive", Algo: module.Incremental},
		},
	}
	chartMutexLockAcquisitions = module.Chart{
		ID:       "mutex_lock_acquisitions",
		Title:    "Mutex lock acquisitions",
		Units:    "acquisitions/s",
		Fam:      "locks",
		Ctx:      "mongodb.lock_acquisitions",
		Priority: prioMutexLockAcquisitions,
		Labels:   []module.Label{{Key: "lock_type", Value: "mutex"}},
		Dims: module.Dims{
			{ID: "locks_mutex_acquire_shared", Name: "shared", Algo: module.Incremental},
			{ID: "locks_mutex_acquire_exclusive", Name: "exclusive", Algo: module.Incremental},
			{ID: "locks_mutex_acquire_intent_shared", Name: "intent_shared", Algo: module.Incremental},
			{ID: "locks_mutex_acquire_intent_exclusive", Name: "intent_exclusive", Algo: module.Incremental},
		},
	}
	chartMetadataLockAcquisitions = module.Chart{
		ID:       "metadata_lock_acquisitions",
		Title:    "Metadata lock acquisitions",
		Units:    "acquisitions/s",
		Fam:      "locks",
		Ctx:      "mongodb.lock_acquisitions",
		Priority: prioMetadataLockAcquisitions,
		Labels:   []module.Label{{Key: "lock_type", Value: "metadata"}},
		Dims: module.Dims{
			{ID: "locks_metadata_acquire_shared", Name: "shared", Algo: module.Incremental},
			{ID: "locks_metadata_acquire_exclusive", Name: "exclusive", Algo: module.Incremental},
			{ID: "locks_metadata_acquire_intent_shared", Name: "intent_shared", Algo: module.Incremental},
			{ID: "locks_metadata_acquire_intent_exclusive", Name: "intent_exclusive", Algo: module.Incremental},
		},
	}
	chartOpLogLockAcquisitions = module.Chart{
		ID:       "oplog_lock_acquisitions",
		Title:    "Operations log lock acquisitions",
		Units:    "acquisitions/s",
		Fam:      "locks",
		Ctx:      "mongodb.lock_acquisitions",
		Priority: prioOpLogLockAcquisitions,
		Labels:   []module.Label{{Key: "lock_type", Value: "oplog"}},
		Dims: module.Dims{
			{ID: "locks_oplog_acquire_shared", Name: "shared", Algo: module.Incremental},
			{ID: "locks_oplog_acquire_exclusive", Name: "exclusive", Algo: module.Incremental},
			{ID: "locks_oplog_acquire_intent_shared", Name: "intent_shared", Algo: module.Incremental},
			{ID: "locks_oplog_acquire_intent_exclusive", Name: "intent_exclusive", Algo: module.Incremental},
		},
	}

	chartWiredTigerConcurrentReadTransactions = module.Chart{
		ID:       "wiredtiger_concurrent_read_transactions_usage",
		Title:    "Wired Tiger concurrent read transactions usage",
		Units:    "transactions",
		Fam:      "wiredtiger",
		Ctx:      "mongodb.wiredtiger_concurrent_read_transactions_usage",
		Priority: prioWiredTigerConcurrentReadTransactionsUsage,
		Type:     module.Stacked,
		Dims: module.Dims{
			{ID: "wiredtiger_concurrent_transactions_read_available", Name: "available"},
			{ID: "wiredtiger_concurrent_transactions_read_out", Name: "used"},
		},
	}
	chartWiredTigerConcurrentWriteTransactions = module.Chart{
		ID:       "wiredtiger_concurrent_write_transactions_usage",
		Title:    "Wired Tiger concurrent write transactions usage",
		Units:    "transactions",
		Fam:      "wiredtiger",
		Ctx:      "mongodb.wiredtiger_concurrent_write_transactions_usage",
		Priority: prioWiredTigerConcurrentWriteTransactionsUsage,
		Type:     module.Stacked,
		Dims: module.Dims{
			{ID: "wiredtiger_concurrent_transactions_write_available", Name: "available"},
			{ID: "wiredtiger_concurrent_transactions_write_out", Name: "used"},
		},
	}
	chartWiredTigerCacheUsage = module.Chart{
		ID:       "wiredtiger_cache_usage",
		Title:    "Wired Tiger cache usage",
		Units:    "bytes",
		Fam:      "wiredtiger",
		Ctx:      "mongodb.wiredtiger_cache_usage",
		Priority: prioWiredTigerCacheUsage,
		Type:     module.Stacked,
		Dims: module.Dims{
			{ID: "wiredtiger_cache_currently_in_cache_bytes", Name: "used"},
		},
	}
	chartWiredTigerCacheDirtySpaceSize = module.Chart{
		ID:       "wiredtiger_cache_dirty_space_size",
		Title:    "Wired Tiger cache dirty space size",
		Units:    "bytes",
		Fam:      "wiredtiger",
		Ctx:      "mongodb.wiredtiger_cache_dirty_space_size",
		Priority: prioWiredTigerCacheDirtySpaceSize,
		Dims: module.Dims{
			{ID: "wiredtiger_cache_tracked_dirty_in_the_cache_bytes", Name: "dirty"},
		},
	}
	chartWiredTigerCacheIORate = module.Chart{
		ID:       "wiredtiger_cache_io_rate",
		Title:    "Wired Tiger IO activity",
		Units:    "pages/s",
		Fam:      "wiredtiger",
		Ctx:      "mongodb.wiredtiger_cache_io_rate",
		Priority: prioWiredTigerCacheIORate,
		Dims: module.Dims{
			{ID: "wiredtiger_cache_read_into_cache_pages", Name: "read", Algo: module.Incremental},
			{ID: "wiredtiger_cache_written_from_cache_pages", Name: "written", Algo: module.Incremental},
		},
	}
	chartWiredTigerCacheEvictionsRate = module.Chart{
		ID:       "wiredtiger_cache_eviction_rate",
		Title:    "Wired Tiger cache evictions",
		Units:    "pages/s",
		Fam:      "wiredtiger",
		Ctx:      "mongodb.wiredtiger_cache_evictions_rate",
		Type:     module.Stacked,
		Priority: prioWiredTigerCacheEvictionsRate,
		Dims: module.Dims{
			{ID: "wiredtiger_cache_unmodified_evicted_pages", Name: "unmodified", Algo: module.Incremental},
			{ID: "wiredtiger_cache_modified_evicted_pages", Name: "modified", Algo: module.Incremental},
		},
	}
)

var (
	chartTmplDatabaseCollections = &module.Chart{
		ID:       "database_%s_collections",
		Title:    "Database collections",
		Units:    "collections",
		Fam:      "databases",
		Ctx:      "mongodb.database_collections",
		Priority: prioDatabaseCollections,
		Dims: module.Dims{
			{ID: "database_%s_collections", Name: "collections"},
		},
	}
	chartTmplDatabaseIndexes = &module.Chart{
		ID:       "database_%s_indexes",
		Title:    "Database indexes",
		Units:    "indexes",
		Fam:      "databases",
		Ctx:      "mongodb.database_indexes",
		Priority: prioDatabaseIndexes,
		Dims: module.Dims{
			{ID: "database_%s_indexes", Name: "indexes"},
		},
	}
	chartTmplDatabaseViews = &module.Chart{
		ID:       "database_%s_views",
		Title:    "Database views",
		Units:    "views",
		Fam:      "databases",
		Ctx:      "mongodb.database_views",
		Priority: prioDatabaseViews,
		Dims: module.Dims{
			{ID: "database_%s_views", Name: "views"},
		},
	}
	chartTmplDatabaseDocuments = &module.Chart{
		ID:       "database_%s_documents",
		Title:    "Database documents",
		Units:    "documents",
		Fam:      "databases",
		Ctx:      "mongodb.database_documents",
		Priority: prioDatabaseDocuments,
		Dims: module.Dims{
			{ID: "database_%s_documents", Name: "documents"},
		},
	}
	chartTmplDatabaseDataSize = &module.Chart{
		ID:       "database_%s_data_size",
		Title:    "Database data size",
		Units:    "bytes",
		Fam:      "databases",
		Ctx:      "mongodb.database_data_size",
		Priority: prioDatabaseDataSize,
		Dims: module.Dims{
			{ID: "database_%s_data_size", Name: "data_size"},
		},
	}
	chartTmplDatabaseStorageSize = &module.Chart{
		ID:       "database_%s_storage_size",
		Title:    "Database storage size",
		Units:    "bytes",
		Fam:      "databases",
		Ctx:      "mongodb.database_storage_size",
		Priority: prioDatabaseStorageSize,
		Dims: module.Dims{
			{ID: "database_%s_storage_size", Name: "storage_size"},
		},
	}
	chartTmplDatabaseIndexSize = &module.Chart{
		ID:       "database_%s_index_size",
		Title:    "Database index size",
		Units:    "bytes",
		Fam:      "databases",
		Ctx:      "mongodb.database_index_size",
		Priority: prioDatabaseIndexSize,
		Dims: module.Dims{
			{ID: "database_%s_index_size", Name: "index_size"},
		},
	}
)

var (
	chartTmplReplSetMemberState = &module.Chart{
		ID:       "replica_set_member_%s_state",
		Title:    "Replica Set member state",
		Units:    "state",
		Fam:      "replica sets",
		Ctx:      "mongodb.repl_set_member_state",
		Priority: prioReplSetMemberState,
		Dims: module.Dims{
			{ID: "repl_set_member_%s_state_primary", Name: "primary"},
			{ID: "repl_set_member_%s_state_startup", Name: "startup"},
			{ID: "repl_set_member_%s_state_secondary", Name: "secondary"},
			{ID: "repl_set_member_%s_state_recovering", Name: "recovering"},
			{ID: "repl_set_member_%s_state_startup2", Name: "startup2"},
			{ID: "repl_set_member_%s_state_unknown", Name: "unknown"},
			{ID: "repl_set_member_%s_state_arbiter", Name: "arbiter"},
			{ID: "repl_set_member_%s_state_down", Name: "down"},
			{ID: "repl_set_member_%s_state_rollback", Name: "rollback"},
			{ID: "repl_set_member_%s_state_removed", Name: "removed"},
		},
	}
	chartTmplReplSetMemberHealthStatus = &module.Chart{
		ID:       "replica_set_member_%s_health_status",
		Title:    "Replica Set member health status",
		Units:    "status",
		Fam:      "replica sets",
		Ctx:      "mongodb.repl_set_member_health_status",
		Priority: prioReplSetMemberHealthStatus,
		Dims: module.Dims{
			{ID: "repl_set_member_%s_health_status_up", Name: "up"},
			{ID: "repl_set_member_%s_health_status_down", Name: "down"},
		},
	}
	chartTmplReplSetMemberReplicationLag = &module.Chart{
		ID:       "replica_set_member_%s_replication_lag",
		Title:    "Replica Set member replication lag",
		Units:    "milliseconds",
		Fam:      "replica sets",
		Ctx:      "mongodb.repl_set_member_replication_lag",
		Priority: prioReplSetMemberReplicationLag,
		Dims: module.Dims{
			{ID: "repl_set_member_%s_replication_lag", Name: "replication_lag"},
		},
	}
	chartTmplReplSetMemberHeartbeatLatency = &module.Chart{
		ID:       "replica_set_member_%s_heartbeat_latency",
		Title:    "Replica Set member heartbeat latency",
		Units:    "milliseconds",
		Fam:      "replica sets",
		Ctx:      "mongodb.repl_set_member_heartbeat_latency",
		Priority: prioReplSetMemberHeartbeatLatency,
		Dims: module.Dims{
			{ID: "repl_set_member_%s_heartbeat_latency", Name: "heartbeat_latency"},
		},
	}
	chartTmplReplSetMemberPingRTT = &module.Chart{
		ID:       "replica_set_member_%s_ping_rtt",
		Title:    "Replica Set member ping RTT",
		Units:    "milliseconds",
		Fam:      "replica sets",
		Ctx:      "mongodb.repl_set_member_ping_rtt",
		Priority: prioReplSetMemberPingRTT,
		Dims: module.Dims{
			{ID: "repl_set_member_%s_ping_rtt", Name: "ping_rtt"},
		},
	}
	chartTmplReplSetMemberUptime = &module.Chart{
		ID:       "replica_set_member_%s_uptime",
		Title:    "Replica Set member uptime",
		Units:    "seconds",
		Fam:      "replica sets",
		Ctx:      "mongodb.repl_set_member_uptime",
		Priority: prioReplSetMemberUptime,
		Dims: module.Dims{
			{ID: "repl_set_member_%s_uptime", Name: "uptime"},
		},
	}
)

var (
	chartShardingNodesCount = &module.Chart{
		ID:       "sharding_nodes_count",
		Title:    "Sharding Nodes",
		Units:    "nodes",
		Fam:      "sharding",
		Ctx:      "mongodb.sharding_nodes_count",
		Type:     module.Stacked,
		Priority: prioShardingNodesCount,
		Dims: module.Dims{
			{ID: "shard_nodes_aware", Name: "shard_aware"},
			{ID: "shard_nodes_unaware", Name: "shard_unaware"},
		},
	}
	chartShardingShardedDatabases = &module.Chart{
		ID:       "sharding_sharded_databases_count",
		Title:    "Sharded databases",
		Units:    "databases",
		Fam:      "sharding",
		Ctx:      "mongodb.sharding_sharded_databases_count",
		Type:     module.Stacked,
		Priority: prioShardingShardedDatabasesCount,
		Dims: module.Dims{
			{ID: "shard_databases_partitioned", Name: "partitioned"},
			{ID: "shard_databases_unpartitioned", Name: "unpartitioned"},
		},
	}

	chartShardingShardedCollectionsCount = &module.Chart{
		ID:       "sharding_sharded_collections_count",
		Title:    "Sharded collections",
		Units:    "collections",
		Fam:      "sharding",
		Ctx:      "mongodb.sharding_sharded_collections_count",
		Type:     module.Stacked,
		Priority: prioShardingShardedCollectionsCount,
		Dims: module.Dims{
			{ID: "shard_collections_partitioned", Name: "partitioned"},
			{ID: "shard_collections_unpartitioned", Name: "unpartitioned"},
		},
	}

	chartTmplShardChunks = &module.Chart{
		ID:       "sharding_shard_id_%s_chunks",
		Title:    "Shard chunks",
		Units:    "chunks",
		Fam:      "sharding",
		Ctx:      "mongodb.sharding_shard_chunks_count",
		Priority: prioShardChunks,
		Dims: module.Dims{
			{ID: "shard_id_%s_chunks", Name: "chunks"},
		},
	}
)
