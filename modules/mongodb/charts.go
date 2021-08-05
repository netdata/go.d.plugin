package mongo

import (
	"strings"

	"github.com/netdata/go.d.plugin/agent/module"
)

var serverStatusCharts = module.Charts{
	// default charts
	chartOpcounter.Copy(),
	chartOpLatencies.Copy(),
	chartConnections.Copy(),
	chartNetwork.Copy(),
	chartNetworkRequests.Copy(),
	chartMemory.Copy(),
	chartPageFaults.Copy(),
	chartAsserts.Copy(),

	//Optional charts
	//chartTransactionsCurrent.Copy(),
	//chartGlobalLockActiveClients.Copy(),
	//chartCollections.Copy(),
	//chartTcmallocGeneric.Copy(),
	//chartTcmalloc.Copy(),
	//chartGlobalLockCurrentQueue.Copy(),
	//chartMetricsCommands.Copy(),
	//chartGlobalLocks.Copy(),
	//chartFlowControl.Copy(),
	//chartWiredTigerBlockManager.Copy(),
	//chartWiredTigerCache.Copy(),
	//chartWiredTigerCapacity.Copy(),
	//chartWiredTigerConnection.Copy(),
	//chartWiredTigerCursor.Copy(),
	//chartWiredTigerLock.Copy(),
	//chartWiredTigerLockDuration.Copy(),
	//chartWiredTigerLogOps.Copy(),
	//chartWiredTigerLogBytes.Copy(),
	//chartWiredTigerTransactions.Copy(),
}

// default charts
var (
	chartOpcounter = module.Chart{
		ID:    "opcounters",
		Title: "Operations by type",
		Units: "ops/s",
		Fam:   "operations",
		Ctx:   "mongodb.operations",
		Dims: module.Dims{
			{ID: "opcounters.insert", Name: "insert", Algo: module.Incremental},
			{ID: "opcounters.query", Name: "query", Algo: module.Incremental},
			{ID: "opcounters.update", Name: "update", Algo: module.Incremental},
			{ID: "opcounters.delete", Name: "delete", Algo: module.Incremental},
			{ID: "opcounters.getmore", Name: "getmore", Algo: module.Incremental},
			{ID: "opcounters.command", Name: "command", Algo: module.Incremental},
		},
	}
	chartOpLatencies = module.Chart{
		ID:    "opLatencies",
		Title: "Operations Latency",
		Units: "milliseconds",
		Fam:   "operations",
		Ctx:   "mongodb.operations_latency",
		Dims: module.Dims{
			{ID: "opLatencies.reads.latency", Name: "Reads", Algo: module.Incremental, Div: 1000},
			{ID: "opLatencies.writes.latency", Name: "Writes", Algo: module.Incremental, Div: 1000},
			{ID: "opLatencies.commands.latency", Name: "Commands", Algo: module.Incremental, Div: 1000},
		},
	}
)

var chartConnections = module.Chart{
	ID:    "connections",
	Title: "Connections",
	Units: "connections",
	Fam:   "connections",
	Ctx:   "mongodb.connections",
	Dims: module.Dims{
		{ID: "connections.current", Name: "current"},
		{ID: "connections.active", Name: "active"},
		{ID: "connections.threaded", Name: "threaded"},
		{ID: "connections.exhaustIsMaster", Name: "exhaustIsMaster"},
		{ID: "connections.exhaustHello", Name: "exhaustHello"},
		{ID: "connections.awaitingTopologyChanges", Name: "awaiting topology changes"},
		{ID: "connections.available", Name: "available"},
	},
}

var (
	chartNetwork = module.Chart{
		ID:    "network",
		Title: "Network IO",
		Units: "bytes/s",
		Fam:   "network",
		Ctx:   "mongodb.network_io",
		Type:  module.Area,
		Dims: module.Dims{
			{ID: "network.bytesIn", Name: "in", Algo: module.Incremental},
			{ID: "network.bytesOut", Name: "out", Algo: module.Incremental, Mul: -1},
		},
	}
	chartNetworkRequests = module.Chart{
		ID:    "networkRequests",
		Title: "Network Requests",
		Units: "requests/s",
		Fam:   "network",
		Ctx:   "mongodb.network_requests",
		Dims: module.Dims{
			{ID: "network.numRequests", Name: "Requests", Algo: module.Incremental},
		},
	}
)

var chartMemory = module.Chart{
	ID:    "mem",
	Title: "Memory usage",
	Units: "MiB",
	Fam:   "memory",
	Ctx:   "mongodb.memory_usage",
	Dims: module.Dims{
		{ID: "mem.resident", Name: "resident"},
		{ID: "mem.virtual", Name: "virtual"},
		{ID: "mem.mapped", Name: "mapped"},
		{ID: "mem.mappedWithJournal", Name: "mapped with journal"},
	},
}

var chartPageFaults = module.Chart{
	ID:    "page_faults",
	Title: "Page faults",
	Units: "page faults/s",
	Fam:   "memory",
	Ctx:   "mongodb.page_faults",
	Dims: module.Dims{
		{ID: "extra_info.page_faults", Name: "Page Faults", Algo: module.Incremental},
	},
}

var chartAsserts = module.Chart{
	ID:    "asserts",
	Title: "Raised assertions",
	Units: "asserts/s",
	Fam:   "asserts",
	Ctx:   "mongodb.asserts",
	Type:  module.Stacked,
	Dims: module.Dims{
		{ID: "asserts.regular", Name: "regular", Algo: module.Incremental},
		{ID: "asserts.warning", Name: "warning", Algo: module.Incremental},
		{ID: "asserts.msg", Name: "msg", Algo: module.Incremental},
		{ID: "asserts.user", Name: "user", Algo: module.Incremental},
		{ID: "asserts.tripwire", Name: "tripwire", Algo: module.Incremental},
		{ID: "asserts.rollovers", Name: "rollovers", Algo: module.Incremental},
	},
}

// option charts
var chartTransactionsCurrent = module.Chart{
	ID:    "transactionsCurrent",
	Title: "Current Transactions",
	Units: "transactions",
	Fam:   "transactionsCurrent",
	Ctx:   "mongodb.current_transactions",
	Dims: module.Dims{
		{ID: "transactions.currentActive", Name: "active"},
		{ID: "transactions.currentInactive", Name: "inactive"},
		{ID: "transactions.currentOpen", Name: "open"},
		{ID: "transactions.currentPrepared", Name: "prepared"},
	},
}

var chartGlobalLockActiveClients = module.Chart{
	ID:    "globalLockActiveClients",
	Title: "Active Clients",
	Units: "clients",
	Fam:   "clients",
	Ctx:   "mongodb.currentQueue",
	Dims: module.Dims{
		{ID: "globalLock.activeClients.readers", Name: "readers"},
		{ID: "globalLock.activeClients.writers", Name: "writers"},
	},
}

var chartCollections = module.Chart{
	ID:    "catalogStats",
	Title: "Catalog Stats",
	Units: "objects",
	Fam:   "catalogStats",
	Ctx:   "mongodb.catalogStats",
	Type:  module.Stacked,
	Dims: module.Dims{
		{ID: "catalogStats.collections", Name: "collections"},
		{ID: "catalogStats.capped", Name: "capped"},
		{ID: "catalogStats.timeseries", Name: "timeseries"},
		{ID: "catalogStats.views", Name: "views"},
		{ID: "catalogStats.internalCollections", Name: "internalCollections"},
		{ID: "catalogStats.internalViews", Name: "internalViews"},
	},
}

var (
	chartTcmallocGeneric = module.Chart{
		ID:    "tcmallocGeneric",
		Title: "Tcmalloc generic metrics",
		Units: "MiB",
		Fam:   "tcmalloc",
		Ctx:   "mongodb.tcmallocGeneric",
		Dims: module.Dims{
			{ID: "tcmalloc.generic.current_allocated_bytes", Name: "current_allocated_bytes", Div: 1 << 20},
			{ID: "tcmalloc.generic.heap_size", Name: "heap_size", Div: 1 << 20},
		},
	}
	chartTcmalloc = module.Chart{
		ID:    "tcmalloc",
		Title: "Tcmalloc",
		Units: "KiB",
		Fam:   "tcmalloc",
		Ctx:   "mongodb.tcmalloc",
		Dims: module.Dims{
			{ID: "tcmalloc.tcmalloc.pageheap_free_bytes", Name: "Pageheap free", Div: 1024},
			{ID: "tcmalloc.tcmalloc.pageheap_unmapped_bytes", Name: "Pageheap unmapped ", Div: 1024},
			{ID: "tcmalloc.tcmalloc.current_total_thread_cache_bytes", Name: "Total threaded cache", Div: 1024},
			{ID: "tcmalloc.tcmalloc.total_free_bytes", Name: "Free", Div: 1024},
			{ID: "tcmalloc.tcmalloc.pageheap_committed_bytes", Name: "Pageheap committed", Div: 1024},
			{ID: "tcmalloc.tcmalloc.pageheap_total_commit_bytes", Name: "Pageheap total commit", Div: 1024},
			{ID: "tcmalloc.tcmalloc.pageheap_total_decommit_bytes", Name: "Pageheap decommit", Div: 1024},
			{ID: "tcmalloc.tcmalloc.pageheap_total_reserve_bytes", Name: "Pageheap reserve", Div: 1024},
		},
	}
)

var chartGlobalLockCurrentQueue = module.Chart{
	ID:    "globalLockCurrentQueue",
	Title: "Current Queue Clients",
	Units: "clients",
	Fam:   "clients",
	Ctx:   "mongodb.currentQueue",
	Type:  module.Stacked,
	Dims: module.Dims{
		{ID: "globalLock.currentQueue.readers", Name: "readers"},
		{ID: "globalLock.currentQueue.writers", Name: "writers"},
	},
}

var chartMetricsCommands = module.Chart{
	ID:    "metricsCommand",
	Title: "Command Metrics",
	Units: "commands/s",
	Fam:   "commands",
	Ctx:   "mongodb.metricsCommand",
	Dims: module.Dims{
		{ID: "metrics.commands.eval.total", Name: "Eval", Algo: module.Incremental},
		{ID: "metrics.commands.eval.failed", Name: "Eval Failed", Algo: module.Incremental},
		{ID: "metrics.commands.delete.total", Name: "Delete", Algo: module.Incremental},
		{ID: "metrics.commands.delete.failed", Name: "Delete Failed", Algo: module.Incremental},
		{ID: "metrics.commands.count.failed", Name: "Count Failed", Algo: module.Incremental},
		{ID: "metrics.commands.createIndexes", Name: "Create Indexes", Algo: module.Incremental},
		{ID: "metrics.commands.findAndModify", Name: "Find And Modify", Algo: module.Incremental},
		{ID: "metrics.commands.insert.failed", Name: "Insert Fail", Algo: module.Incremental},
	},
}

var chartGlobalLocks = module.Chart{
	ID:    "locks",
	Title: "Acquired locks",
	Units: "locks/s",
	Fam:   "locks",
	Ctx:   "mongodb.locks",
	Dims: module.Dims{
		{ID: "locks.Global.acquireCount.r", Name: "global read", Algo: module.Incremental},
		{ID: "locks.Global.acquireCount.w", Name: "global write", Algo: module.Incremental},
		{ID: "locks.Database.acquireCount.r", Name: "database read", Algo: module.Incremental},
		{ID: "locks.Database.acquireCount.w", Name: "database write", Algo: module.Incremental},
		{ID: "locks.Collection.acquireCount.r", Name: "collection read", Algo: module.Incremental},
		{ID: "locks.Collection.acquireCount.w", Name: "collection write", Algo: module.Incremental},
	},
}

var chartFlowControl = module.Chart{
	ID:    "flowControl",
	Title: "Flow Control Stats",
	Units: "milliseconds",
	Fam:   "flowControl",
	Ctx:   "mongodb.flow_control_timings",
	Dims: module.Dims{
		{ID: "flowControl.timeAcquiringMicros", Name: "timeAcquiring", Algo: module.Incremental, Div: 1000},
		{ID: "flowControl.isLaggedTimeMicros", Name: "isLaggedTime", Algo: module.Incremental, Div: 1000},
	},
}

var (
	// WiredTiger (optional)
	chartWiredTigerBlockManager = module.Chart{
		ID:    "wiredtigerBlockManager",
		Title: "Wired Tiger Block Manager",
		Units: "KiB",
		Fam:   "wiredtiger",
		Ctx:   "mongodb.wiredTigerBlocks",
		Type:  module.Stacked,
		Dims: module.Dims{
			{ID: toID("wiredTiger.block-manager.bytes read"), Name: "bytes read", Div: 1024},
			{ID: toID("wiredTiger.block-manager.bytes read via memory map API"), Name: "bytes read via memory map API", Div: 1024},
			{ID: toID("wiredTiger.block-manager.bytes read via system call API"), Name: "bytes read via system call API", Div: 1024},
			{ID: toID("wiredTiger.block-manager.bytes written"), Name: "bytes written", Div: 1024},
			{ID: toID("wiredTiger.block-manager.bytes written for checkpoint"), Name: "bytes written for checkpoint", Div: 1024},
			{ID: toID("wiredTiger.block-manager.bytes written via memory map API"), Name: "bytes written via memory map API", Div: 1024},
		},
	}
	chartWiredTigerCache = module.Chart{
		ID:    "wiredtigerCache",
		Title: "Wired Tiger Cache",
		Units: "KiB",
		Fam:   "wiredtiger",
		Ctx:   "mongodb.wiredTigerCache",
		Dims: module.Dims{
			{ID: toID("wiredTiger.cache.bytes allocated for updates"), Name: "bytes allocated for updates", Div: 1024},
			{ID: toID("wiredTiger.cache.bytes read into cache"), Name: "bytes read into cache", Div: 1024},
			{ID: toID("wiredTiger.cache.bytes written from cache"), Name: "bytes written from cache", Div: 1024},
		},
	}
	chartWiredTigerCapacity = module.Chart{
		ID:    "wiredtigerCapacity",
		Title: "Wired Tiger Capacity",
		Units: "usec",
		Fam:   "wiredtiger",
		Ctx:   "mongodb.wiredTigerCapacity",

		Dims: module.Dims{
			{ID: toID("wiredTiger.capacity.time waiting due to total capacity (usecs)"), Name: "time waiting due to total capacity (usecs)"},
			{ID: toID("wiredTiger.capacity.time waiting during checkpoint (usecs)"), Name: "time waiting during checkpoint (usecs)"},
			{ID: toID("wiredTiger.capacity.time waiting during eviction (usecs)"), Name: "time waiting during eviction (usecs)"},
			{ID: toID("wiredTiger.capacity.time waiting during logging (usecs)"), Name: "time waiting during logging (usecs)"},
			{ID: toID("wiredTiger.capacity.time waiting during read (usecs)"), Name: "time waiting during read (usecs)"},
		},
	}
	chartWiredTigerConnection = module.Chart{
		ID:    "wiredtigerConnection",
		Title: "Wired Tiger Connections",
		Units: "ops/s",
		Fam:   "wiredtiger",
		Ctx:   "mongodb.wiredTigerConnection",

		Dims: module.Dims{
			{ID: toID("wiredTiger.connection.memory allocations"), Name: "memory allocations", Algo: module.Incremental},
			{ID: toID("wiredTiger.connection.memory frees"), Name: "memory frees", Algo: module.Incremental},
			{ID: toID("wiredTiger.connection.memory re-allocations"), Name: "memory re-allocations", Algo: module.Incremental},
		},
	}
	chartWiredTigerCursor = module.Chart{
		ID:    "wiredtigerCursor",
		Title: "Wired Tiger Cursor",
		Units: "calls/s",
		Fam:   "wiredtiger",
		Ctx:   "mongodb.wiredTigerCursor",

		Dims: module.Dims{
			{ID: toID("wiredTiger.cursor.open cursor count"), Name: "open cursor count", Algo: module.Incremental},
			{ID: toID("wiredTiger.cursor.cached cursor count"), Name: "cached cursor count", Algo: module.Incremental},
			{ID: toID("wiredTiger.cursor.cursor bulk loaded cursor insert calls"), Name: "cursor bulk loaded cursor insert calls", Algo: module.Incremental},
			{ID: toID("wiredTiger.cursor.cursor close calls that result in cache"), Name: "cursor close calls that result in cache", Algo: module.Incremental},
			{ID: toID("wiredTiger.cursor.cursor create calls"), Name: "cursor create calls", Algo: module.Incremental},
			{ID: toID("wiredTiger.cursor.cursor insert calls"), Name: "cursor insert calls", Algo: module.Incremental},
			{ID: toID("wiredTiger.cursor.cursor modify calls"), Name: "cursor modify calls", Algo: module.Incremental},
			{ID: toID("wiredTiger.cursor.cursor next calls"), Name: "cursor next calls", Algo: module.Incremental},
			{ID: toID("wiredTiger.cursor.cursor operation restarted"), Name: "cursor operation restarted", Algo: module.Incremental},
			{ID: toID("wiredTiger.cursor.cursor prev calls"), Name: "cursor prev calls", Algo: module.Incremental},
			{ID: toID("wiredTiger.cursor.cursor remove calls"), Name: "cursor remove calls", Algo: module.Incremental},
			{ID: toID("wiredTiger.cursor.cursor remove key bytes removed"), Name: "cursor remove key bytes removed", Algo: module.Incremental},
			{ID: toID("wiredTiger.cursor.cursor reserve calls"), Name: "cursor reserve calls", Algo: module.Incremental},
			{ID: toID("wiredTiger.cursor.cursor reset calls"), Name: "cursor reset calls", Algo: module.Incremental},
			{ID: toID("wiredTiger.cursor.cursor search calls"), Name: "cursor search calls", Algo: module.Incremental},
			{ID: toID("wiredTiger.cursor.cursor search history store calls"), Name: "cursor search history store calls", Algo: module.Incremental},
			{ID: toID("wiredTiger.cursor.cursor search near calls"), Name: "cursor search near calls", Algo: module.Incremental},
			{ID: toID("wiredTiger.cursor.cursor sweep buckets"), Name: "cursor sweep buckets", Algo: module.Incremental},
			{ID: toID("wiredTiger.cursor.cursor sweep cursors closed"), Name: "cursor sweep cursors closed", Algo: module.Incremental},
			{ID: toID("wiredTiger.cursor.cursor sweep cursors examined"), Name: "cursor sweep cursors examined", Algo: module.Incremental},
			{ID: toID("wiredTiger.cursor.cursor sweeps"), Name: "cursor sweeps", Algo: module.Incremental},
			{ID: toID("wiredTiger.cursor.cursor truncate calls"), Name: "cursor truncate calls", Algo: module.Incremental},
			{ID: toID("wiredTiger.cursor.cursor update calls"), Name: "cursor update calls", Algo: module.Incremental},
			{ID: toID("wiredTiger.cursor.cursor update value size change"), Name: "cursor update value size change", Algo: module.Incremental},
		},
	}
	chartWiredTigerLock = module.Chart{
		ID:    "wiredtigerLock",
		Title: "Wired Tiger Lock",
		Units: "ops/s",
		Fam:   "wiredtiger",
		Ctx:   "mongodb.wiredTigerLock",

		Dims: module.Dims{
			{ID: toID("wiredTiger.lock.checkpoint lock acquisitions"), Name: "checkpoint lock acquisitions", Algo: module.Incremental},
			{ID: toID("wiredTiger.lock.dhandle read lock acquisitions"), Name: "dhandle read lock acquisitions", Algo: module.Incremental},
			{ID: toID("wiredTiger.lock.dhandle write lock acquisitions"), Name: "dhandle write lock acquisitions", Algo: module.Incremental},
			{ID: toID("wiredTiger.lock.durable timestamp queue read lock acquisitions"), Name: "durable timestamp queue read lock acquisitions", Algo: module.Incremental},
			{ID: toID("wiredTiger.lock.durable timestamp queue write lock acquisitions"), Name: "durable timestamp queue write lock acquisitions", Algo: module.Incremental},
			{ID: toID("wiredTiger.lock.metadata lock acquisitions"), Name: "metadata lock acquisitions", Algo: module.Incremental},
			{ID: toID("wiredTiger.lock.read timestamp queue read lock acquisitions"), Name: "read timestamp queue read lock acquisitions", Algo: module.Incremental},
			{ID: toID("wiredTiger.lock.read timestamp queue write lock acquisitions"), Name: "read timestamp queue write lock acquisitions", Algo: module.Incremental},
			{ID: toID("wiredTiger.lock.schema lock acquisitions"), Name: "schema lock acquisitions", Algo: module.Incremental},
			{ID: toID("wiredTiger.lock.table read lock acquisitions"), Name: "table read lock acquisitions", Algo: module.Incremental},
			{ID: toID("wiredTiger.lock.table write lock acquisitions"), Name: "table write lock acquisitions", Algo: module.Incremental},
			{ID: toID("wiredTiger.lock.txn global read lock acquisitions"), Name: "txn global read lock acquisitions", Algo: module.Incremental},
		},
	}
	chartWiredTigerLockDuration = module.Chart{
		ID:    "wiredtigerLockDuration",
		Title: "Wired Tiger Lock Duration",
		Units: "usec",
		Fam:   "wiredtiger",
		Ctx:   "mongodb.wiredTigerLockDuration",

		Dims: module.Dims{
			{ID: toID("wiredTiger.lock.checkpoint lock application thread wait time (usecs)"), Name: "checkpoint lock application thread wait time (usecs)"},
			{ID: toID("wiredTiger.lock.checkpoint lock internal thread wait time (usecs)"), Name: "checkpoint lock internal thread wait time (usecs)"},
			{ID: toID("wiredTiger.lock.dhandle lock application thread time waiting (usecs)"), Name: "dhandle lock application thread time waiting (usecs)"},
			{ID: toID("wiredTiger.lock.dhandle lock internal thread time waiting (usecs)"), Name: "dhandle lock internal thread time waiting (usecs)"},
			{ID: toID("wiredTiger.lock.durable timestamp queue lock application thread time waiting (usecs)"), Name: "durable timestamp queue lock application thread time waiting (usecs)"},
			{ID: toID("wiredTiger.lock.durable timestamp queue lock internal thread time waiting (usecs)"), Name: "durable timestamp queue lock internal thread time waiting (usecs)"},
			{ID: toID("wiredTiger.lock.metadata lock application thread wait time (usecs)"), Name: "metadata lock application thread wait time (usecs)"},
			{ID: toID("wiredTiger.lock.metadata lock internal thread wait time (usecs)"), Name: "metadata lock internal thread wait time (usecs)"},
			{ID: toID("wiredTiger.lock.read timestamp queue lock application thread time waiting (usecs)"), Name: "read timestamp queue lock application thread time waiting (usecs)"},
			{ID: toID("wiredTiger.lock.read timestamp queue lock internal thread time waiting (usecs)"), Name: "read timestamp queue lock internal thread time waiting (usecs)"},
			{ID: toID("wiredTiger.lock.schema lock application thread wait time (usecs)"), Name: "schema lock application thread wait time (usecs)"},
			{ID: toID("wiredTiger.lock.schema lock internal thread wait time (usecs)"), Name: "schema lock internal thread wait time (usecs)"},
			{ID: toID("wiredTiger.lock.table lock application thread time waiting for the table lock (usecs)"), Name: "table lock application thread time waiting for the table lock (usecs)"},
			{ID: toID("wiredTiger.lock.table lock internal thread time waiting for the table lock (usecs)"), Name: "table lock internal thread time waiting for the table lock (usecs)"},
			{ID: toID("wiredTiger.lock.txn global lock application thread time waiting (usecs)"), Name: "txn global lock application thread time waiting (usecs)"},
			{ID: toID("wiredTiger.lock.txn global lock internal thread time waiting (usecs)"), Name: "txn global lock internal thread time waiting (usecs)"},
		},
	}
	chartWiredTigerLogOps = module.Chart{
		ID:    "wiredtigerLogOps",
		Title: "Wired Tiger Log Operations",
		Units: "ops/s",
		Fam:   "wiredtiger",
		Ctx:   "mongodb.wiredTigerLogOps",

		Dims: module.Dims{
			{ID: toID("wiredTiger.log.log flush operations"), Name: "log flush operations", Algo: module.Incremental},
			{ID: toID("wiredTiger.log.log force write operations"), Name: "log force write operations", Algo: module.Incremental},
			{ID: toID("wiredTiger.log.log force write operations skipped"), Name: "log force write operations skipped", Algo: module.Incremental},
			{ID: toID("wiredTiger.log.log scan operations"), Name: "log scan operations", Algo: module.Incremental},
			{ID: toID("wiredTiger.log.log sync operations"), Name: "log sync operations", Algo: module.Incremental},
			{ID: toID("wiredTiger.log.log sync_dir operations"), Name: "log sync_dir operations", Algo: module.Incremental},
			{ID: toID("wiredTiger.log.log write operations"), Name: "log write operations", Algo: module.Incremental},
		},
	}
	chartWiredTigerLogBytes = module.Chart{
		ID:    "wiredtigerLogOps",
		Title: "Wired Tiger Log Operations",
		Units: "bytes/s",
		Fam:   "wiredtiger",
		Ctx:   "mongodb.wiredTigerLogOps",

		Dims: module.Dims{
			{ID: toID("wiredTiger.log.log bytes of payload data"), Name: "log bytes of payload data", Algo: module.Incremental},
			{ID: toID("wiredTiger.log.log bytes written"), Name: "log bytes written", Algo: module.Incremental},
			{ID: toID("wiredTiger.log.logging bytes consolidated"), Name: "logging bytes consolidated", Algo: module.Incremental},
			{ID: toID("wiredTiger.log.total log buffer size"), Name: "total log buffer size", Algo: module.Incremental},
		},
	}
	chartWiredTigerTransactions = module.Chart{
		ID:    "wiredtigerTransactions",
		Title: "Wired Tiger Log Transactions",
		Units: "transactions/s",
		Fam:   "wiredtiger",
		Ctx:   "mongodb.wiredTigerTransactions",

		Dims: module.Dims{
			{ID: toID("wiredTiger.transaction.prepared transactions"), Name: "prepared transactions", Algo: module.Incremental},
			{ID: toID("wiredTiger.transaction.query timestamp calls"), Name: "query timestamp calls", Algo: module.Incremental},
			{ID: toID("wiredTiger.transaction.rollback to stable calls"), Name: "rollback to stable calls", Algo: module.Incremental},
			{ID: toID("wiredTiger.transaction.set timestamp calls"), Name: "set timestamp calls", Algo: module.Incremental},
			{ID: toID("wiredTiger.transaction.transaction begins"), Name: "transaction begins", Algo: module.Incremental},
			{ID: toID("wiredTiger.transaction.transaction sync calls"), Name: "transaction sync calls", Algo: module.Incremental},
			{ID: toID("wiredTiger.transaction.transactions committed"), Name: "transactions committed", Algo: module.Incremental},
			{ID: toID("wiredTiger.transaction.transactions rolled back"), Name: "transactions rolled back", Algo: module.Incremental},
		},
	}
)

func toID(in string) string {
	id := strings.ReplaceAll(in, " ", "%20")
	return id
}

func fromID(in string) string {
	id := strings.ReplaceAll(in, "%20", " ")
	return id
}
