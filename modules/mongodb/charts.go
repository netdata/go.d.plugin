package mongo

import (
	"github.com/netdata/go.d.plugin/agent/module"
)

// these charts are expected to be available in many versions
// and build in mongoDB and we are always creating them
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

var dbStatsCharts = []*module.Chart{
	dbStatsCollectionsChart,
	dbStatsIndexesChart,
	dbStatsViewsChart,
	dbStatsDocumentsChart,
	dbStatsSizeChart,
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
		Ctx:   "mongb.connections_rate",
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
			{ID: "tcmalloc_tcmalloc_pageheap_unmapped", Name: "pageheap unmapped "},
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
			{ID: "wiredtiger_cursor_search", Name: "cursor eset calls", Algo: module.Incremental},
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
		Title: "Wired Tiger Lock",
		Units: "ops/s",
		Fam:   "wiredtiger",
		Ctx:   "mongodb.wiredtiger_lock",

		Dims: module.Dims{
			{ID: "wiredtiger_lock_checkpoint_acquisitions", Name: "checkpoint lock acquisitions", Algo: module.Incremental},
			{ID: "wiredtiger_lock_read_acquisitions", Name: "dhandle read lock acquisitions", Algo: module.Incremental},
			{ID: "wiredtiger_lock_write_acquisitions", Name: "dhandle write lock acquisitions", Algo: module.Incremental},
			{ID: "wiredtiger_lock_durable_timestamp_queue_read_acquisitions", Name: "durable timestamp queue read lock acquisitions", Algo: module.Incremental},
			{ID: "wiredtiger_lock_durable", Name: "durable timestamp queue write lock acquisitions", Algo: module.Incremental},
			{ID: "wiredtiger_lock_metadata_acquisitions", Name: "metadata lock acquisitions", Algo: module.Incremental},
			{ID: "wiredtiger_lock_read_timestamp_queue_read_acquisitions", Name: "read timestamp queue read lock acquisitions", Algo: module.Incremental},
			{ID: "wiredtiger_lock_read", Name: "read timestamp queue write lock acquisitions", Algo: module.Incremental},
			{ID: "wiredtiger_lock_schema_acquisitions", Name: "schema lock acquisitions", Algo: module.Incremental},
			{ID: "wiredtiger_lock_table_read_acquisitions", Name: "table read lock acquisitions", Algo: module.Incremental},
			{ID: "wiredtiger_lock_table_write_acquisitions", Name: "table write lock acquisitions", Algo: module.Incremental},
			{ID: "wiredtiger_lock_txn_global_read_acquisitions", Name: "txn global read lock acquisitions", Algo: module.Incremental},
		},
	}

	chartWiredTigerLockDuration = module.Chart{
		ID:    "wiredtiger_lock_duration",
		Title: "Wired Tiger Lock Duration",
		Units: "usec",
		Fam:   "wiredtiger",
		Ctx:   "mongodb.wiredtiger_lock_duration",

		Dims: module.Dims{
			{ID: "wiredtiger_lock_checkpoint_wait_time", Name: "checkpoint lock application thread wait time"},
			{ID: "wiredtiger_lock_checkpoint_internal_thread_wait_time", Name: "checkpoint lock internal thread wait time"},
			{ID: "wiredtiger_lock_application_thread_time_waiting", Name: "dhandle lock application thread time waiting"},
			{ID: "wiredtiger_lock_internal_thread_time_waiting", Name: "dhandle lock internal thread time waiting"},
			{ID: "wiredtiger_lock_durable_timestamp_queue_application_thread_time_waiting", Name: "durable timestamp queue lock application thread time waiting"},
			{ID: "wiredtiger_lock_durable_timestamp_queue_internal_thread_time_waiting", Name: "durable timestamp queue lock internal thread time waiting"},
			{ID: "wiredtiger_lock_metadata_application_thread_wait_time", Name: "metadata lock application thread wait time"},
			{ID: "wiredtiger_lock_metadata_internal_thread_wait_time", Name: "metadata lock internal thread wait time"},
			{ID: "wiredtiger_lock_read_timestamp_queue_application_thread_time_waiting", Name: "read timestamp queue lock application thread time waiting"},
			{ID: "wiredtiger_lock_read_timestamp_queue_internal_thread_time_waiting", Name: "read timestamp queue lock internal thread time waiting"},
			{ID: "wiredtiger_lock_schema_application_thread_wait_time", Name: "schema lock application thread wait time"},
			{ID: "wiredtiger_lock_schema_internal_thread_wait_time", Name: "schema lock internal thread wait time"},
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
		Units: "transactions/s",
		Fam:   "wiredtiger",
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
	dbStatsCollectionsChart = &module.Chart{
		ID:    "collections",
		Title: "Collections",
		Units: "collections",
		Fam:   "database_statistics",
		Ctx:   "mongodb.collections",
		Type:  module.Stacked,
	}
	dbStatsIndexesChart = &module.Chart{
		ID:    "indexes",
		Title: "Indexes",
		Units: "indexes",
		Fam:   "database_statistics",
		Ctx:   "mongodb.indexes",
		Type:  module.Stacked,
	}
	dbStatsViewsChart = &module.Chart{
		ID:    "views",
		Title: "Views",
		Units: "views",
		Fam:   "database_statistics",
		Ctx:   "mongodb.views",
		Type:  module.Stacked,
	}

	dbStatsDocumentsChart = &module.Chart{
		ID:    "documents",
		Title: "Documents",
		Units: "documents",
		Fam:   "database_statistics",
		Ctx:   "mongodb.documents",
		Type:  module.Stacked,
	}

	dbStatsSizeChart = &module.Chart{
		ID:    "storage_size",
		Title: "Disk Size",
		Units: "bytes",
		Fam:   "database_statistics",
		Ctx:   "mongodb.storage_size",
		Type:  module.Stacked,
	}
)

func (m *Mongo) dimsForDbStats(newDatabases []string) {
	if len(newDatabases) == 0 {
		return
	}
	if m.databasesMatcher == nil {
		return
	}

	// remove dims for not existing databases
	diff := sliceDiff(m.databases, newDatabases)
	for _, name := range diff {
		for _, chart := range dbStatsCharts {
			id := chart.ID + "_" + name
			err := chart.MarkDimRemove(id, true)
			if err != nil {
				m.Warningf("failed to remove dimension %s with error: %s", id, err.Error())
				continue
			}
			chart.MarkNotCreated()
		}
	}

	// add dimensions for new databases
	for _, chart := range dbStatsCharts {
		for _, name := range newDatabases {
			if !m.databasesMatcher.MatchString(name) {
				continue
			}
			if !m.dimsEnabled[chart.ID+"_"+name] {
				id := chart.ID + "_" + name
				err := chart.AddDim(&module.Dim{ID: id, Name: name, Algo: module.Absolute})
				if err != nil {
					m.Errorf("failed to add dim: %s, %s", id, err)
					continue
				}
				chart.MarkNotCreated()
				m.dimsEnabled[chart.ID+"_"+name] = true
			}
		}
	}
}

func sliceDiff(slice1, slice2 []string) []string {
	mb := make(map[string]struct{}, len(slice2))
	for _, x := range slice2 {
		mb[x] = struct{}{}
	}
	var diff []string
	for _, x := range slice1 {
		if _, found := mb[x]; !found {
			diff = append(diff, x)
		}
	}
	return diff
}
