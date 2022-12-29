// SPDX-License-Identifier: GPL-3.0-or-later

package mongo

import (
	"github.com/netdata/go.d.plugin/agent/module"
)

const (
	prioOperations = module.Priority + iota
	prioOperationsLatency

	prioConnectionsUsage
	prioConnectionsRate
	prioConnectionsState

	prioNetworkIO
	prioNetworkRequests

	prioPageFaults
	prioTCMallocGeneric
	prioTCMalloc

	prioAsserts

	prioTransactionsCurrent
	prioTransactionsCommitTypes

	prioActiveClients
	prioQueuedOperations

	prioLocks

	prioFlowControlTimings

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
	prioDatabaseStorageSize

	prioReplSetMemberState
	prioReplSetMemberHealthStatus
	prioReplSetMemberReplicationLag
	prioReplSetMemberHeartbeatLatency
	prioReplSetMemberPingRTT
	prioReplSetMemberUptime

	prioShardNodesCount
	prioShardDatabasesStatus
	prioShardCollectionsStatus
	prioShardChunks
)

// these charts are expected to be available in many versions
var chartsServerStatus = module.Charts{
	chartOperations.Copy(),
	chartOperationsLatency.Copy(),
	chartConnectionsUsage.Copy(),
	chartConnectionsRate.Copy(),
	chartConnectionsByState.Copy(),
	chartNetworkIO.Copy(),
	chartNetworkRequests.Copy(),
	chartPageFaults.Copy(),
	chartAsserts.Copy(),
}

var chartsTmplDatabase = module.Charts{
	chartTmplDatabaseCollections.Copy(),
	chartTmplDatabaseIndexes.Copy(),
	chartTmplDatabaseViews.Copy(),
	chartTmplDatabaseDocuments.Copy(),
	chartTmplDatabaseStorageSize.Copy(),
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
	chartShardNodes.Copy(),
	chartShardDatabases.Copy(),
	chartShardCollections.Copy(),
	chartShardChunksTmpl.Copy(),
}

var (
	chartOperations = module.Chart{
		ID:       "operations",
		Title:    "Operations by type",
		Units:    "ops/s",
		Fam:      "operations",
		Ctx:      "mongodb.operations",
		Priority: prioOperations,
		Dims: module.Dims{
			{ID: "operations_insert", Name: "insert", Algo: module.Incremental},
			{ID: "operations_query", Name: "query", Algo: module.Incremental},
			{ID: "operations_update", Name: "update", Algo: module.Incremental},
			{ID: "operations_delete", Name: "delete", Algo: module.Incremental},
			{ID: "operations_getmore", Name: "getmore", Algo: module.Incremental},
			{ID: "operations_command", Name: "command", Algo: module.Incremental},
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
			{ID: "operations_latency_read", Name: "reads", Algo: module.Incremental, Div: 1000},
			{ID: "operations_latency_write", Name: "writes", Algo: module.Incremental, Div: 1000},
			{ID: "operations_latency_command", Name: "commands", Algo: module.Incremental, Div: 1000},
		},
	}

	chartConnectionsUsage = module.Chart{
		ID:       "connections",
		Title:    "Connections",
		Units:    "connections",
		Fam:      "connections",
		Ctx:      "mongodb.connections",
		Type:     module.Stacked,
		Priority: prioConnectionsUsage,
		Dims: module.Dims{
			{ID: "connections_current", Name: "current"},
			{ID: "connections_available", Name: "available"},
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
	chartConnectionsByState = module.Chart{
		ID:       "connections_state",
		Title:    "Connections By State",
		Units:    "connections",
		Fam:      "connections",
		Ctx:      "mongodb.connections_state",
		Priority: prioConnectionsState,
		Dims: module.Dims{
			{ID: "connections_active", Name: "active"},
			{ID: "connections_threaded", Name: "threaded"},
			{ID: "connections_exhaustIsMaster", Name: "exhaustIsMaster"},
			{ID: "connections_exhaustHello", Name: "exhaustHello"},
			{ID: "connections_awaitingTopologyChanges", Name: "awaiting topology changes"},
		},
	}

	chartNetworkIO = module.Chart{
		ID:       "network",
		Title:    "Network IO",
		Units:    "bytes/s",
		Fam:      "network",
		Ctx:      "mongodb.network_io",
		Priority: prioNetworkIO,
		Type:     module.Area,
		Dims: module.Dims{
			{ID: "network_bytes_in", Name: "in", Algo: module.Incremental, Mul: -1},
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

	chartPageFaults = module.Chart{
		ID:       "page_faults",
		Title:    "Page faults",
		Units:    "page faults/s",
		Fam:      "memory",
		Ctx:      "mongodb.page_faults",
		Priority: prioPageFaults,
		Dims: module.Dims{
			{ID: "extra_info_page_faults", Name: "page Faults", Algo: module.Incremental},
		},
	}
	tcMallocGenericChart = module.Chart{
		ID:       "tcmalloc_generic",
		Title:    "Tcmalloc generic metrics",
		Units:    "bytes",
		Fam:      "memory",
		Ctx:      "mongodb.tcmalloc_generic",
		Priority: prioTCMallocGeneric,
		Dims: module.Dims{
			{ID: "tcmalloc_generic_current_allocated", Name: "current_allocated"},
			{ID: "tcmalloc_generic_heap_size", Name: "heap_size"},
		},
	}
	tcMallocChart = module.Chart{
		ID:       "tcmalloc",
		Title:    "Tcmalloc",
		Units:    "bytes",
		Fam:      "memory",
		Ctx:      "mongodb.tcmalloc",
		Priority: prioTCMalloc,
		Dims: module.Dims{
			{ID: "tcmalloc_tcmalloc_pageheap_free", Name: "pageheap_free"},
			{ID: "tcmalloc_tcmalloc_pageheap_unmapped", Name: "pageheap_unmapped"},
			{ID: "tcmalloc_tcmalloc_max_total_thread_cache", Name: "max_total_thread_cache"},
			{ID: "tcmalloc_tcmalloc_total_free", Name: "total_free"},
			{ID: "tcmalloc_tcmalloc_pageheap_committed", Name: "pageheap_committed"},
			{ID: "tcmalloc_tcmalloc_pageheap_total_commit", Name: "pageheap_total_commit"},
			{ID: "tcmalloc_tcmalloc_pageheap_total_decommit", Name: "pageheap_total_decommit"},
			{ID: "tcmalloc_tcmalloc_pageheap_total_reserve", Name: "pageheap_total_reserve"},
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
		Title:    "Current Transactions",
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
	chartTransactionsCommitTypes = module.Chart{
		ID:       "transactions_commit_types",
		Title:    "Transactions Commit Types",
		Units:    "commits/s",
		Fam:      "transactions",
		Ctx:      "mongodb.transactions_commit_types",
		Priority: prioTransactionsCommitTypes,
		Dims: module.Dims{
			{ID: "transactions_commit_types_no_shards_initiated", Name: "no_shards_initiated", Algo: module.Incremental},
			{ID: "transactions_commit_types_no_shards_successful", Name: "no_shards_successful", Algo: module.Incremental},
			{ID: "transactions_commit_types_single_shard_initiated", Name: "single_shard_initiated", Algo: module.Incremental},
			{ID: "transactions_commit_types_single_shard_successful", Name: "single_shard_successful", Algo: module.Incremental},
			{ID: "transactions_commit_types_single_write_shard_initiated", Name: "single_write_shard_initiated", Algo: module.Incremental},
			{ID: "transactions_commit_types_single_write_shard_successful", Name: "single_write_shard_successful", Algo: module.Incremental},
			{ID: "transactions_commit_types_two_phase_initiated", Name: "two_phase_initiated", Algo: module.Incremental},
			{ID: "transactions_commit_types_two_phase_successful", Name: "two_phase_successful", Algo: module.Incremental},
		},
	}

	chartGlobalLockActiveClients = module.Chart{
		ID:       "active_clients",
		Title:    "Active Clients",
		Units:    "clients",
		Fam:      "clients",
		Ctx:      "mongodb.active_clients",
		Priority: prioActiveClients,
		Dims: module.Dims{
			{ID: "glock_active_clients_readers", Name: "readers"},
			{ID: "glock_active_clients_writers", Name: "writers"},
		},
	}
	chartGlobalLockCurrentQueue = module.Chart{
		ID:       "queued_operations",
		Title:    "Queued operations because of a lock",
		Units:    "operations",
		Fam:      "clients",
		Ctx:      "mongodb.queued_operations",
		Priority: prioQueuedOperations,
		Type:     module.Stacked,
		Dims: module.Dims{
			{ID: "glock_current_queue_readers", Name: "readers"},
			{ID: "glock_current_queue_writers", Name: "writers"},
		},
	}

	chartLocks = module.Chart{
		ID:       "locks",
		Title:    "Acquired locks",
		Units:    "locks/s",
		Fam:      "locks",
		Ctx:      "mongodb.locks",
		Priority: prioLocks,
		Dims: module.Dims{
			{ID: "locks_global_read", Name: "global read", Algo: module.Incremental},
			{ID: "locks_global_write", Name: "global write", Algo: module.Incremental},
			{ID: "locks_database_read", Name: "database read", Algo: module.Incremental},
			{ID: "locks_database_write", Name: "database write", Algo: module.Incremental},
			{ID: "locks_collection_read", Name: "collection read", Algo: module.Incremental},
			{ID: "locks_collection_write", Name: "collection write", Algo: module.Incremental},
		},
	}

	chartFlowControlTimings = module.Chart{
		ID:       "flow_control_timings",
		Title:    "Flow Control Stats",
		Units:    "milliseconds",
		Fam:      "flow_control",
		Ctx:      "mongodb.flow_control_timings",
		Priority: prioFlowControlTimings,
		Dims: module.Dims{
			{ID: "flow_target_rate_limit", Name: "acquiring", Algo: module.Incremental, Div: 1000},
			{ID: "flow_time_acquiring_micros", Name: "lagged", Algo: module.Incremental, Div: 1000},
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
			{ID: "wiredtiger_cache_dirty_space_size", Name: "dirty"},
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
	chartWiredTigerCacheEvictionRate = module.Chart{
		ID:       "wiredtiger_cache_eviction_rate",
		Title:    "Wired Tiger cache evictions",
		Units:    "pages/s",
		Fam:      "wiredtiger",
		Ctx:      "mongodb.wiredtiger_cache_dirty_space_size",
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
		Title:    "Collections",
		Units:    "collections",
		Fam:      "database stats",
		Ctx:      "mongodb.database_collections",
		Priority: prioDatabaseCollections,
		Dims: module.Dims{
			{ID: "database_%s_collections", Name: "collections"},
		},
	}
	chartTmplDatabaseIndexes = &module.Chart{
		ID:       "database_%s_indexes",
		Title:    "Indexes",
		Units:    "indexes",
		Fam:      "database stats",
		Ctx:      "mongodb.database_indexes",
		Priority: prioDatabaseIndexes,
		Dims: module.Dims{
			{ID: "database_%s_indexes", Name: "indexes"},
		},
	}
	chartTmplDatabaseViews = &module.Chart{
		ID:       "database_%s_views",
		Title:    "Views",
		Units:    "views",
		Fam:      "database stats",
		Ctx:      "mongodb.database_views",
		Priority: prioDatabaseViews,
		Dims: module.Dims{
			{ID: "database_%s_views", Name: "views"},
		},
	}
	chartTmplDatabaseDocuments = &module.Chart{
		ID:       "database_%s_documents",
		Title:    "Documents",
		Units:    "documents",
		Fam:      "database stats",
		Ctx:      "mongodb.database_documents",
		Priority: prioDatabaseDocuments,
		Dims: module.Dims{
			{ID: "database_%s_documents", Name: "documents"},
		},
	}
	chartTmplDatabaseStorageSize = &module.Chart{
		ID:       "database_%s_storage_size",
		Title:    "Disk Size",
		Units:    "bytes",
		Fam:      "database stats",
		Ctx:      "mongodb.database_storage_size",
		Priority: prioDatabaseStorageSize,
		Dims: module.Dims{
			{ID: "database_%s_storage_size", Name: "storage_size"},
		},
	}
)

var (
	chartTmplReplSetMemberState = &module.Chart{
		ID:       "replica_set_member_%s_state",
		Title:    "Replica Set member state",
		Units:    "state",
		Fam:      "replica set stats",
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
		Fam:      "replica set stats",
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
		Fam:      "replica set stats",
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
		Fam:      "replica set stats",
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
		Fam:      "replica set stats",
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
		Fam:      "replica set stats",
		Ctx:      "mongodb.repl_set_member_uptime",
		Priority: prioReplSetMemberUptime,
		Dims: module.Dims{
			{ID: "repl_set_member_%s_uptime", Name: "uptime"},
		},
	}
)

var (
	chartShardNodes = &module.Chart{
		ID:       "shard_nodes_count",
		Title:    "Shard Nodes",
		Units:    "nodes",
		Fam:      "sharding",
		Ctx:      "mongodb.shard_nodes_count",
		Type:     module.Stacked,
		Priority: prioShardNodesCount,
		Dims: module.Dims{
			{ID: "shard_nodes_count_aware", Name: "shard aware"},
			{ID: "shard_nodes_count_unaware", Name: "shard unaware"},
		},
	}

	chartShardDatabases = &module.Chart{
		ID:       "shard_databases_status",
		Title:    "Databases Sharding Status",
		Units:    "databases",
		Fam:      "sharding",
		Ctx:      "mongodb.shard_databases_status",
		Type:     module.Stacked,
		Priority: prioShardDatabasesStatus,
		Dims: module.Dims{
			{ID: "shard_databases_partitioned", Name: "partitioned"},
			{ID: "shard_databases_unpartitioned", Name: "un-partitioned"},
		},
	}

	chartShardCollections = &module.Chart{
		ID:       "shard_collections_status",
		Title:    "Collections Sharding Status",
		Units:    "collections",
		Fam:      "sharding",
		Ctx:      "mongodb.shard_collections_status",
		Type:     module.Stacked,
		Priority: prioShardCollectionsStatus,
		Dims: module.Dims{
			{ID: "shard_collections_partitioned", Name: "partitioned"},
			{ID: "shard_collections_unpartitioned", Name: "un-partitioned"},
		},
	}

	shardChartsTmpl = module.Charts{
		chartShardChunksTmpl.Copy(),
	}

	chartShardChunksTmpl = &module.Chart{
		ID:       "shard_id_%s_chunks",
		Title:    "Shard chunks",
		Units:    "chunks",
		Fam:      "sharding",
		Ctx:      "mongodb.shard_chunks",
		Priority: prioShardChunks,
		Dims: module.Dims{
			{ID: "shard_id_%s_chunks", Name: "chunks"},
		},
	}
)
