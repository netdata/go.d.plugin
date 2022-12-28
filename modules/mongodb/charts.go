// SPDX-License-Identifier: GPL-3.0-or-later

package mongo

import (
	"fmt"
	"github.com/netdata/go.d.plugin/agent/module"
	"strings"
)

const (
	_ = module.Priority + iota

	replSetMemberState
	replSetMemberHealthStatus
	replSetMemberReplicationLag
	replSetMemberHeartbeatLatency
	replSetMemberPingRTT
	replSetMemberUptime
)

// these charts are expected to be available in many versions
// and build in mongoDB, and we are always creating them
var serverStatusCharts = module.Charts{
	chartOpcounter.Copy(),
	chartOpLatencies.Copy(),
	chartConnectionsUsage.Copy(),
	chartConnectionsRate.Copy(),
	chartConnectionsByState.Copy(),
	chartNetwork.Copy(),
	chartNetworkRequests.Copy(),
	chartPageFaults.Copy(),
	chartAsserts.Copy(),
}

// dbStatsChartsTmpl are used to collect per database metrics
var dbStatsChartsTmpl = module.Charts{
	chartDBStatsCollectionsTmpl.Copy(),
	chartDBStatsIndexesTmpl,
	chartDBStatsViewsTmpl,
	chartDBStatsDocumentsTmpl,
	chartDBStatsSizeTmpl,
}

// replSetMemberChartsTmpl on used on replica sets
var replSetMemberChartsTmpl = module.Charts{
	replSetMemberStateChartTmpl.Copy(),
	replSetMemberHealthStatusChartTmpl.Copy(),
	replSetMemberReplicationLagChartTmpl,
	replSetMemberHeartbeatLatencyChartTmpl,
	replSetMemberPingRTTChartTmpl,
	replSetMemberUptimeChartTmpl.Copy(),
}

var shardingCharts = module.Charts{
	chartShardNodes,
	chartShardDatabases,
	chartShardCollections,
	chartShardChunksTmpl,
}

var (
	chartOpcounter = module.Chart{
		ID:    "operations",
		Title: "Operations by type",
		Units: "ops/s",
		Fam:   "operations",
		Ctx:   "mongodb.operations",
		Dims: module.Dims{
			{ID: "operations_insert", Name: "insert", Algo: module.Incremental},
			{ID: "operations_query", Name: "query", Algo: module.Incremental},
			{ID: "operations_update", Name: "update", Algo: module.Incremental},
			{ID: "operations_delete", Name: "delete", Algo: module.Incremental},
			{ID: "operations_getmore", Name: "getmore", Algo: module.Incremental},
			{ID: "operations_command", Name: "command", Algo: module.Incremental},
		},
	}
	chartOpLatencies = module.Chart{
		ID:    "operations_latency",
		Title: "Operations Latency",
		Units: "milliseconds",
		Fam:   "operations",
		Ctx:   "mongodb.operations_latency",
		Dims: module.Dims{
			{ID: "operations_latency_read", Name: "reads", Algo: module.Incremental, Div: 1000},
			{ID: "operations_latency_write", Name: "writes", Algo: module.Incremental, Div: 1000},
			{ID: "operations_latency_command", Name: "commands", Algo: module.Incremental, Div: 1000},
		},
	}
)

var (
	chartConnectionsUsage = module.Chart{
		ID:    "connections",
		Title: "Connections",
		Units: "connections",
		Fam:   "connections",
		Ctx:   "mongodb.connections",
		Type:  module.Stacked,
		Dims: module.Dims{
			{ID: "connections_current", Name: "current"},
			{ID: "connections_available", Name: "available"},
		},
	}
	chartConnectionsRate = module.Chart{
		ID:    "connections_rate",
		Title: "Connections Rate",
		Units: "connections/s",
		Fam:   "connections",
		Ctx:   "mongodb.connections_rate",
		Dims: module.Dims{
			{ID: "connections_total_created", Name: "created", Algo: module.Incremental},
		},
	}
	chartConnectionsByState = module.Chart{
		ID:    "connections_state",
		Title: "Connections By State",
		Units: "connections",
		Fam:   "connections",
		Ctx:   "mongodb.connections_state",
		Dims: module.Dims{
			{ID: "connections_active", Name: "active"},
			{ID: "connections_threaded", Name: "threaded"},
			{ID: "connections_exhaustIsMaster", Name: "exhaustIsMaster"},
			{ID: "connections_exhaustHello", Name: "exhaustHello"},
			{ID: "connections_awaitingTopologyChanges", Name: "awaiting topology changes"},
		},
	}
)

var (
	chartNetwork = module.Chart{
		ID:    "network",
		Title: "Network IO",
		Units: "bytes/s",
		Fam:   "network",
		Ctx:   "mongodb.network_io",
		Type:  module.Area,
		Dims: module.Dims{
			{ID: "network_bytes_in", Name: "in", Algo: module.Incremental, Mul: -1},
			{ID: "network_bytes_out", Name: "out", Algo: module.Incremental},
		},
	}
	chartNetworkRequests = module.Chart{
		ID:    "network_requests",
		Title: "Network Requests",
		Units: "requests/s",
		Fam:   "network",
		Ctx:   "mongodb.network_requests",
		Dims: module.Dims{
			{ID: "network_requests", Name: "requests", Algo: module.Incremental},
		},
	}
)

var (
	chartPageFaults = module.Chart{
		ID:    "page_faults",
		Title: "Page faults",
		Units: "page faults/s",
		Fam:   "memory",
		Ctx:   "mongodb.page_faults",
		Dims: module.Dims{
			{ID: "extra_info_page_faults", Name: "page Faults", Algo: module.Incremental},
		},
	}

	chartTcmallocGeneric = module.Chart{
		ID:    "tcmalloc_generic",
		Title: "Tcmalloc generic metrics",
		Units: "bytes",
		Fam:   "memory",
		Ctx:   "mongodb.tcmalloc_generic",
		Dims: module.Dims{
			{ID: "tcmalloc_generic_current_allocated", Name: "current_allocated"},
			{ID: "tcmalloc_generic_heap_size", Name: "heap_size"},
		},
	}

	chartTcmalloc = module.Chart{
		ID:    "tcmalloc",
		Title: "Tcmalloc",
		Units: "bytes",
		Fam:   "memory",
		Ctx:   "mongodb.tcmalloc",
		Dims: module.Dims{
			{ID: "tcmalloc_tcmalloc_pageheap_free", Name: "pageheap free"},
			{ID: "tcmalloc_tcmalloc_pageheap_unmapped", Name: "pageheap unmapped"},
			{ID: "tcmalloc_tcmalloc_max_total_thread_cache", Name: "total threaded cache"},
			{ID: "tcmalloc_tcmalloc_total_free", Name: "free"},
			{ID: "tcmalloc_tcmalloc_pageheap_committed", Name: "pageheap committed"},
			{ID: "tcmalloc_tcmalloc_pageheap_total_commit", Name: "pageheap total commit"},
			{ID: "tcmalloc_tcmalloc_pageheap_total_decommit", Name: "pageheap decommit"},
			{ID: "tcmalloc_tcmalloc_pageheap_total_reserve", Name: "pageheap reserve"},
		},
	}
)

var chartAsserts = module.Chart{
	ID:    "asserts",
	Title: "Raised assertions",
	Units: "asserts/s",
	Fam:   "asserts",
	Ctx:   "mongodb.asserts",
	Type:  module.Stacked,
	Dims: module.Dims{
		{ID: "asserts_regular", Name: "regular", Algo: module.Incremental},
		{ID: "asserts_warning", Name: "warning", Algo: module.Incremental},
		{ID: "asserts_msg", Name: "msg", Algo: module.Incremental},
		{ID: "asserts_user", Name: "user", Algo: module.Incremental},
		{ID: "asserts_tripwire", Name: "tripwire", Algo: module.Incremental},
		{ID: "asserts_rollovers", Name: "rollovers", Algo: module.Incremental},
	},
}

var chartTransactionsCurrent = module.Chart{
	ID:    "current_transactions",
	Title: "Current Transactions",
	Units: "transactions",
	Fam:   "transactions",
	Ctx:   "mongodb.current_transactions",
	Dims: module.Dims{
		{ID: "transactions_active", Name: "active"},
		{ID: "transactions_inactive", Name: "inactive"},
		{ID: "transactions_open", Name: "open"},
		{ID: "transactions_prepared", Name: "prepared"},
	},
}

var chartTransactionsCommitTypes = module.Chart{
	ID:    "shard_commit_types",
	Title: "Shard Commit Types",
	Units: "commits",
	Fam:   "shard stats",
	Ctx:   "mongodb.shard_commit_types",
	Dims: module.Dims{
		{ID: "transactions_commit_types_no_shards_initiated", Name: "no shard (init)"},
		{ID: "transactions_commit_types_no_shards_successful", Name: "no shard (successful)"},
		{ID: "transactions_commit_types_no_single_shard_initiated", Name: "single shard (init)"},
		{ID: "transactions_commit_types_no_single_shard_successful", Name: "single shard (successful)"},
		{ID: "transactions_commit_types_single_write_shard_initiated", Name: "shard write (init)"},
		{ID: "transactions_commit_types_single_write_shard_successful", Name: "shard write (successful)"},
		{ID: "transactions_commit_types_single_two_phase_initiated", Name: "two phase (init)"},
		{ID: "transactions_commit_types_single_two_phase_successful", Name: "two phase (successful)"},
	},
}

var (
	chartGlobalLockActiveClients = module.Chart{
		ID:    "active_clients",
		Title: "Active Clients",
		Units: "clients",
		Fam:   "clients",
		Ctx:   "mongodb.active_clients",
		Dims: module.Dims{
			{ID: "glock_active_clients_readers", Name: "readers"},
			{ID: "glock_active_clients_writers", Name: "writers"},
		},
	}
	chartGlobalLockCurrentQueue = module.Chart{
		ID:    "queued_operations",
		Title: "Queued operations because of a lock",
		Units: "operations",
		Fam:   "clients",
		Ctx:   "mongodb.queued_operations",
		Type:  module.Stacked,
		Dims: module.Dims{
			{ID: "glock_current_queue_readers", Name: "readers"},
			{ID: "glock_current_queue_writers", Name: "writers"},
		},
	}
)

var chartLocks = module.Chart{
	ID:    "locks",
	Title: "Acquired locks",
	Units: "locks/s",
	Fam:   "locks",
	Ctx:   "mongodb.locks",
	Dims: module.Dims{
		{ID: "locks_global_read", Name: "global read", Algo: module.Incremental},
		{ID: "locks_global_write", Name: "global write", Algo: module.Incremental},
		{ID: "locks_database_read", Name: "database read", Algo: module.Incremental},
		{ID: "locks_database_write", Name: "database write", Algo: module.Incremental},
		{ID: "locks_collection_read", Name: "collection read", Algo: module.Incremental},
		{ID: "locks_collection_write", Name: "collection write", Algo: module.Incremental},
	},
}

var chartFlowControl = module.Chart{
	ID:    "flow_control_timings",
	Title: "Flow Control Stats",
	Units: "milliseconds",
	Fam:   "flow_control",
	Ctx:   "mongodb.flow_control_timings",
	Dims: module.Dims{
		{ID: "flow_target_rate_limit", Name: "acquiring", Algo: module.Incremental, Div: 1000},
		{ID: "flow_time_acquiring_micros", Name: "lagged", Algo: module.Incremental, Div: 1000},
	},
}

var (
	chartWiredTigerBlockManager = module.Chart{
		ID:    "wiredtiger_blocks",
		Title: "Wired Tiger Block Manager",
		Units: "bytes",
		Fam:   "wiredtiger",
		Ctx:   "mongodb.wiredtiger_blocks",
		Type:  module.Stacked,
		Dims: module.Dims{
			{ID: "wiredtiger_block_manager_read", Name: "read"},
			{ID: "wiredtiger_block_manager_read_via_memory", Name: "read via memory map API"},
			{ID: "wiredtiger_block_manager_read_via_system_api", Name: "read via system call API"},
			{ID: "wiredtiger_block_manager_written", Name: "written"},
			{ID: "wiredtiger_block_manager_written_for_checkpoint", Name: "written for checkpoint"},
			{ID: "wiredtiger_block_manager_written_via_memory", Name: "written via memory map API"},
		},
	}

	chartWiredTigerCache = module.Chart{
		ID:    "wiredtiger_cache",
		Title: "Wired Tiger Cache",
		Units: "bytes",
		Fam:   "wiredtiger",
		Ctx:   "mongodb.wiredtiger_cache",
		Dims: module.Dims{
			{ID: "wiredtiger_cache_alloccated", Name: "allocated for updates"},
			{ID: "wiredtiger_cache_read", Name: "read into cache"},
			{ID: "wiredtiger_cache_write", Name: "written from cache"},
		},
	}

	chartWiredTigerCapacity = module.Chart{
		ID:    "wiredtiger_capacity",
		Title: "Wired Tiger Capacity Waiting",
		Units: "usec",
		Fam:   "wiredtiger",
		Ctx:   "mongodb.wiredtiger_capacity",

		Dims: module.Dims{
			{ID: "wiredtiger_capacity_wait_capacity", Name: "due to total capacity"},
			{ID: "wiredtiger_capacity_wait_checkpoint", Name: "during checkpoint"},
			{ID: "wiredtiger_capacity_wait_eviction", Name: "during eviction"},
			{ID: "wiredtiger_capacity_wait_logging", Name: "during logging"},
			{ID: "wiredtiger_capacity_wait_read", Name: "during read"},
		},
	}

	chartWiredTigerConnection = module.Chart{
		ID:    "wiredtiger_connection",
		Title: "Wired Tiger Connection",
		Units: "ops/s",
		Fam:   "wiredtiger",
		Ctx:   "mongodb.wiredtiger_connection",

		Dims: module.Dims{
			{ID: "wiredtiger_connection_allocations", Name: "memory allocations", Algo: module.Incremental},
			{ID: "wiredtiger_connection_frees", Name: "memory frees", Algo: module.Incremental},
			{ID: "wiredtiger_connection_reallocations", Name: "memory re-allocations", Algo: module.Incremental},
		},
	}

	chartWiredTigerCursor = module.Chart{
		ID:    "wiredtiger_cursor",
		Title: "Wired Tiger Cursor",
		Units: "calls/s",
		Fam:   "wiredtiger",
		Ctx:   "mongodb.wiredtiger_cursor",

		Dims: module.Dims{
			{ID: "wiredtiger_cursor_count", Name: "open count", Algo: module.Incremental},
			{ID: "wiredtiger_cursor_bulk", Name: "cached count", Algo: module.Incremental},
			{ID: "wiredtiger_cursor_close", Name: "bulk loaded insert calls", Algo: module.Incremental},
			{ID: "wiredtiger_cursor_create", Name: "close calls that result in cache", Algo: module.Incremental},
			{ID: "wiredtiger_cursor_insert", Name: "create calls", Algo: module.Incremental},
			{ID: "wiredtiger_cursor_modify", Name: "insert calls", Algo: module.Incremental},
			{ID: "wiredtiger_cursor_next", Name: "modify calls", Algo: module.Incremental},
			{ID: "wiredtiger_cursor_restarted", Name: "next calls", Algo: module.Incremental},
			{ID: "wiredtiger_cursor_prev", Name: "operation restarted", Algo: module.Incremental},
			{ID: "wiredtiger_cursor_remove", Name: "prev calls", Algo: module.Incremental},
			{ID: "wiredtiger_cursor_reserve", Name: "remove calls", Algo: module.Incremental},
			{ID: "wiredtiger_cursor_reset", Name: "reserve calls", Algo: module.Incremental},
			{ID: "wiredtiger_cursor_search", Name: "cursor reset calls", Algo: module.Incremental},
			{ID: "wiredtiger_cursor_search_history", Name: "search calls", Algo: module.Incremental},
			{ID: "wiredtiger_cursor_search_near", Name: "search history store calls", Algo: module.Incremental},
			{ID: "wiredtiger_cursor_sweep_buckets", Name: "search near calls", Algo: module.Incremental},
			{ID: "wiredtiger_cursor_sweep_cursors", Name: "sweep buckets", Algo: module.Incremental},
			{ID: "wiredtiger_cursor_sweep_examined", Name: "sweep cursors closed", Algo: module.Incremental},
			{ID: "wiredtiger_cursor_sweeps", Name: "sweep cursors examined", Algo: module.Incremental},
			{ID: "wiredtiger_cursor_truncate", Name: "sweeps", Algo: module.Incremental},
			{ID: "wiredtiger_cursor_update", Name: "truncate calls", Algo: module.Incremental},
			{ID: "wiredtiger_cursor_update_value", Name: "update calls", Algo: module.Incremental},
		},
	}

	chartWiredTigerLock = module.Chart{
		ID:    "wiredtiger_lock",
		Title: "Wired Tiger Lock Acquisitions",
		Units: "ops/s",
		Fam:   "wiredtiger",
		Ctx:   "mongodb.wiredtiger_lock",

		Dims: module.Dims{
			{ID: "wiredtiger_lock_checkpoint_acquisitions", Name: "checkpoint", Algo: module.Incremental},
			{ID: "wiredtiger_lock_read_acquisitions", Name: "dhandle read", Algo: module.Incremental},
			{ID: "wiredtiger_lock_write_acquisitions", Name: "dhandle write", Algo: module.Incremental},
			{ID: "wiredtiger_lock_durable_timestamp_queue_read_acquisitions", Name: "durable timestamp queue read", Algo: module.Incremental},
			{ID: "wiredtiger_lock_durable_timestamp_queue_write_acquisitions", Name: "durable timestamp queue write", Algo: module.Incremental},
			{ID: "wiredtiger_lock_metadata_acquisitions", Name: "metadata", Algo: module.Incremental},
			{ID: "wiredtiger_lock_read_timestamp_queue_read_acquisitions", Name: "read timestamp queue read", Algo: module.Incremental},
			{ID: "wiredtiger_lock_read_timestamp_queue_write_acquisitions", Name: "read timestamp queue write", Algo: module.Incremental},
			{ID: "wiredtiger_lock_schema_acquisitions", Name: "schema", Algo: module.Incremental},
			{ID: "wiredtiger_lock_table_read_acquisitions", Name: "table read", Algo: module.Incremental},
			{ID: "wiredtiger_lock_table_write_acquisitions", Name: "table write", Algo: module.Incremental},
			{ID: "wiredtiger_lock_txn_global_read_acquisitions", Name: "txn global read", Algo: module.Incremental},
		},
	}

	chartWiredTigerLockDuration = module.Chart{
		ID:    "wiredtiger_lock_duration",
		Title: "Wired Tiger Lock Duration",
		Units: "usec",
		Fam:   "wiredtiger",
		Ctx:   "mongodb.wiredtiger_lock_duration",

		Dims: module.Dims{
			{ID: "wiredtiger_lock_checkpoint_wait_time", Name: "checkpoint"},
			{ID: "wiredtiger_lock_checkpoint_internal_thread_wait_time", Name: "checkpoint internal thread"},
			{ID: "wiredtiger_lock_application_thread_time_waiting", Name: "dhandle application thread"},
			{ID: "wiredtiger_lock_internal_thread_time_waiting", Name: "dhandle internal thread"},
			{ID: "wiredtiger_lock_durable_timestamp_queue_application_thread_time_waiting", Name: "durable timestamp queue application thread"},
			{ID: "wiredtiger_lock_durable_timestamp_queue_internal_thread_time_waiting", Name: "durable timestamp queue internal thread"},
			{ID: "wiredtiger_lock_metadata_application_thread_wait_time", Name: "metadata application thread"},
			{ID: "wiredtiger_lock_metadata_internal_thread_wait_time", Name: "metadata internal thread"},
			{ID: "wiredtiger_lock_read_timestamp_queue_application_thread_time_waiting", Name: "read timestamp queue application thread"},
			{ID: "wiredtiger_lock_read_timestamp_queue_internal_thread_time_waiting", Name: "read timestamp queue internal thread"},
			{ID: "wiredtiger_lock_schema_application_thread_wait_time", Name: "schema application thread"},
			{ID: "wiredtiger_lock_schema_internal_thread_wait_time", Name: "schema internal thread"},
		},
	}

	chartWiredTigerLogOps = module.Chart{
		ID:    "wiredtiger_log_ops",
		Title: "Wired Tiger Log Operations",
		Units: "ops/s",
		Fam:   "wiredtiger",
		Ctx:   "mongodb.wiredtiger_log_ops",

		Dims: module.Dims{
			{ID: "wiredtiger_log_flush", Name: "flush", Algo: module.Incremental},
			{ID: "wiredtiger_log_force_write", Name: "force write", Algo: module.Incremental},
			{ID: "wiredtiger_log_write_skip", Name: "force write skipped", Algo: module.Incremental},
			{ID: "wiredtiger_log_scan", Name: "scan", Algo: module.Incremental},
			{ID: "wiredtiger_log_sync", Name: "sync", Algo: module.Incremental},
			{ID: "wiredtiger_log_sync_dir", Name: "sync_dir", Algo: module.Incremental},
			{ID: "wiredtiger_log_write", Name: "write", Algo: module.Incremental},
		},
	}

	chartWiredTigerLogBytes = module.Chart{
		ID:    "wiredtiger_log_ops_size",
		Title: "Wired Tiger Log Operations",
		Units: "bytes/s",
		Fam:   "wiredtiger",
		Ctx:   "mongodb.wiredtiger_log_ops_size",

		Dims: module.Dims{
			{ID: "wiredtiger_log_payload", Name: "payload data", Algo: module.Incremental},
			{ID: "wiredtiger_log_written", Name: "written", Algo: module.Incremental},
			{ID: "wiredtiger_log_consolidated", Name: "consolidated", Algo: module.Incremental},
			{ID: "wiredtiger_log_buffer_size", Name: "total buffer size", Algo: module.Incremental},
		},
	}

	chartWiredTigerTransactions = module.Chart{
		ID:    "wiredtiger_transactions",
		Title: "Wired Tiger Transactions",
		Fam:   "wiredtiger",
		Units: "transactions/s",
		Ctx:   "mongodb.wiredtiger_transactions",

		Dims: module.Dims{
			{ID: "wiredtiger_transaction_prepare", Name: "prepared", Algo: module.Incremental},
			{ID: "wiredtiger_transaction_query", Name: "query timestamp", Algo: module.Incremental},
			{ID: "wiredtiger_transaction_rollback", Name: "rollback to stable", Algo: module.Incremental},
			{ID: "wiredtiger_transaction_set_timestamp", Name: "set timestamp", Algo: module.Incremental},
			{ID: "wiredtiger_transaction_begin", Name: "begins", Algo: module.Incremental},
			{ID: "wiredtiger_transaction_sync", Name: "sync", Algo: module.Incremental},
			{ID: "wiredtiger_transaction_committed", Name: "committed", Algo: module.Incremental},
			{ID: "wiredtiger_transaction_rolled_back", Name: "rolled back", Algo: module.Incremental},
		},
	}
)

var (
	chartDBStatsCollectionsTmpl = &module.Chart{
		ID:    "database_%s_collections",
		Title: "Collections",
		Units: "collections",
		Fam:   "database_statistics",
		Ctx:   "mongodb.database_collections",
		Dims: module.Dims{
			{ID: "database_%s_collections", Name: "collections"},
		},
	}

	chartDBStatsIndexesTmpl = &module.Chart{
		ID:    "database_%s_indexes",
		Title: "Indexes",
		Units: "indexes",
		Fam:   "database_statistics",
		Ctx:   "mongodb.database_indexes",
		Dims: module.Dims{
			{ID: "database_%s_indexes", Name: "indexes"},
		},
	}

	chartDBStatsViewsTmpl = &module.Chart{
		ID:    "database_%s_views",
		Title: "Views",
		Units: "views",
		Fam:   "database_statistics",
		Ctx:   "mongodb.database_views",
		Dims: module.Dims{
			{ID: "database_%s_views", Name: "views"},
		},
	}

	chartDBStatsDocumentsTmpl = &module.Chart{
		ID:    "database_%s_documents",
		Title: "Documents",
		Units: "documents",
		Fam:   "database_statistics",
		Ctx:   "mongodb.database_documents",
		Dims: module.Dims{
			{ID: "database_%s_documents", Name: "documents"},
		},
	}

	chartDBStatsSizeTmpl = &module.Chart{
		ID:    "database_%s_storage_size",
		Title: "Disk Size",
		Units: "bytes",
		Fam:   "database_statistics",
		Ctx:   "mongodb.database_storage_size",
		Dims: module.Dims{
			{ID: "database_%s_storage_size", Name: "storage_size"},
		},
	}
)

var (
	replSetMemberStateChartTmpl = &module.Chart{
		ID:       "replica_set_member_%s_state",
		Title:    "Replica Set member state",
		Units:    "state",
		Fam:      "replica set",
		Ctx:      "mongodb.repl_set_member_state",
		Priority: replSetMemberState,
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
	replSetMemberHealthStatusChartTmpl = &module.Chart{
		ID:       "replica_set_member_%s_health_status",
		Title:    "Replica Set member health status",
		Units:    "status",
		Fam:      "replica set",
		Ctx:      "mongodb.repl_set_member_health_status",
		Priority: replSetMemberHealthStatus,
		Dims: module.Dims{
			{ID: "repl_set_member_%s_health_status_up", Name: "up"},
			{ID: "repl_set_member_%s_health_status_down", Name: "down"},
		},
	}
	replSetMemberReplicationLagChartTmpl = &module.Chart{
		ID:       "replica_set_member_%s_replication_lag",
		Title:    "Replica Set member replication lag",
		Units:    "milliseconds",
		Fam:      "replica set",
		Ctx:      "mongodb.repl_set_member_replication_lag",
		Priority: replSetMemberReplicationLag,
		Dims: module.Dims{
			{ID: "repl_set_member_%s_replication_lag", Name: "replication_lag"},
		},
	}
	replSetMemberHeartbeatLatencyChartTmpl = &module.Chart{
		ID:       "replica_set_member_%s_heartbeat_latency",
		Title:    "Replica Set member heartbeat latency",
		Units:    "milliseconds",
		Fam:      "replica set",
		Ctx:      "mongodb.repl_set_member_heartbeat_latency",
		Priority: replSetMemberHeartbeatLatency,
		Dims: module.Dims{
			{ID: "repl_set_member_%s_heartbeat_latency", Name: "heartbeat_latency"},
		},
	}
	replSetMemberPingRTTChartTmpl = &module.Chart{
		ID:       "replica_set_member_%s_ping_rtt",
		Title:    "Replica Set member ping RTT",
		Units:    "milliseconds",
		Fam:      "replica set",
		Ctx:      "mongodb.repl_set_member_ping_rtt",
		Priority: replSetMemberPingRTT,
		Dims: module.Dims{
			{ID: "repl_set_member_%s_ping_rtt", Name: "ping_rtt"},
		},
	}
	replSetMemberUptimeChartTmpl = &module.Chart{
		ID:       "replica_set_member_%s_uptime",
		Title:    "Replica Set member uptime",
		Units:    "seconds",
		Fam:      "replica set",
		Ctx:      "mongodb.repl_set_member_uptime",
		Priority: replSetMemberUptime,
		Dims: module.Dims{
			{ID: "repl_set_member_%s_uptime", Name: "uptime"},
		},
	}
)

var (
	chartShardNodes = &module.Chart{
		ID:    "shard_nodes_count",
		Title: "Shard Nodes",
		Units: "nodes",
		Fam:   "shard stats",
		Ctx:   "mongodb.shard_nodes_count",
		Type:  module.Stacked,
		Dims: module.Dims{
			{ID: "shard_nodes_count_aware", Name: "shard aware"},
			{ID: "shard_nodes_count_unaware", Name: "shard unaware"},
		},
	}

	chartShardDatabases = &module.Chart{
		ID:    "shard_databases_status",
		Title: "Databases Sharding Status",
		Units: "databases",
		Fam:   "shard stats",
		Ctx:   "mongodb.shard_databases_status",
		Type:  module.Stacked,
		Dims: module.Dims{
			{ID: "shard_databases_partitioned", Name: "partitioned"},
			{ID: "shard_databases_unpartitioned", Name: "un-partitioned"},
		},
	}

	chartShardCollections = &module.Chart{
		ID:    "shard_collections_status",
		Title: "Collections Sharding Status",
		Units: "collections",
		Fam:   "shard stats",
		Ctx:   "mongodb.shard_collections_status",
		Type:  module.Stacked,
		Dims: module.Dims{
			{ID: "shard_collections_partitioned", Name: "partitioned"},
			{ID: "shard_collections_unpartitioned", Name: "un-partitioned"},
		},
	}

	shardChartsTmpl = module.Charts{
		chartShardChunksTmpl.Copy(),
	}

	chartShardChunksTmpl = &module.Chart{
		ID:    "shard_id_%s_chunks",
		Title: "Shard chunks",
		Units: "chunks",
		Fam:   "shard stats",
		Ctx:   "mongodb.shard_chunks",
		Dims: module.Dims{
			{ID: "shard_id_%s_chunks", Name: "chunks"},
		},
	}
)

func (m *Mongo) addDatabaseCharts(name string) {
	charts := dbStatsChartsTmpl.Copy()

	for _, chart := range *charts {
		chart.ID = fmt.Sprintf(chart.ID, name)
		chart.Labels = []module.Label{
			{Key: "database", Value: name},
		}
		for _, dim := range chart.Dims {
			dim.ID = fmt.Sprintf(dim.ID, name)
		}
	}

	if err := m.Charts().Add(*charts...); err != nil {
		m.Warning(err)
	}
}

func (m *Mongo) removeDatabaseCharts(name string) {
	px := fmt.Sprintf("database_%s_", name)

	for _, chart := range *m.Charts() {
		if strings.HasPrefix(chart.ID, px) {
			chart.MarkRemove()
			chart.MarkNotCreated()
		}
	}
}

func (m *Mongo) addReplSetMemberCharts(v replSetMember) {
	charts := replSetMemberChartsTmpl.Copy()

	if v.Self != nil {
		_ = charts.Remove(replSetMemberHeartbeatLatencyChartTmpl.ID)
		_ = charts.Remove(replSetMemberPingRTTChartTmpl.ID)
		_ = charts.Remove(replSetMemberUptimeChartTmpl.ID)
	}

	for _, chart := range *charts {
		chart.ID = fmt.Sprintf(chart.ID, v.Name)
		chart.Labels = []module.Label{
			{Key: "repl_set_member", Value: v.Name},
		}
		for _, dim := range chart.Dims {
			dim.ID = fmt.Sprintf(dim.ID, v.Name)
		}
	}

	if err := m.Charts().Add(*charts...); err != nil {
		m.Warning(err)
	}
}

func (m *Mongo) removeReplSetMemberCharts(name string) {
	px := fmt.Sprintf("repl_set_member_%s_", name)

	for _, chart := range *m.Charts() {
		if strings.HasPrefix(chart.ID, px) {
			chart.MarkRemove()
			chart.MarkNotCreated()
		}
	}
}

func (m *Mongo) addShardingCharts() {
	charts := shardingCharts.Copy()

	if err := m.Charts().Add(*charts...); err != nil {
		m.Warning(err)
	}
}
