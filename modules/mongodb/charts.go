package mongo

import (
	"strings"

	"github.com/netdata/go.d.plugin/agent/module"
)

var serverStatusCharts = module.Charts{
	chartOpcounter.Copy(),
	chartGlobalLockCurrentQueue.Copy(),
	chartGlobalLockActiveClients.Copy(),
	chartConnections.Copy(),
	chartMemory.Copy(),
	chartAsserts.Copy(),
	chartPageFaults.Copy(),
	chartNetwork.Copy(),
	chartNetworkRequests.Copy(),
	chartTransactionsCurrent.Copy(),
	chartCollections.Copy(),
	chartFlowControl.Copy(),
	chartTcmallocGeneric.Copy(),
	chartTcmalloc.Copy(),
	chartWiredTigerBlockManager.Copy(),
	chartWiredTigerCache.Copy(),
	chartWiredTigerCapacity.Copy(),
	chartWiredTigerConnection.Copy(),
	chartWiredTigerCursor.Copy(),
	chartWiredTigerLock.Copy(),
	chartWiredTigerLockDuration.Copy(),
	chartWiredTigerLogOps.Copy(),
	chartWiredTigerLogBytes.Copy(),
	chartWiredTigerTransactions.Copy(),
	chartGlobalLocks.Copy(),
	chartMetricsCommands.Copy(),
}

var (
	chartOpcounter = module.Chart{
		ID:    "opcounters",
		Title: "Commands total rate",
		Units: "commands/s",
		Fam:   "commands",
		Ctx:   "mongodb.command_total_rate",
		Type:  module.Stacked,
		Dims: module.Dims{
			{ID: "opcounters.insert", Name: "insert", Algo: module.Incremental},
			{ID: "opcounters.query", Name: "query", Algo: module.Incremental},
			{ID: "opcounters.update", Name: "update", Algo: module.Incremental},
			{ID: "opcounters.delete", Name: "delete", Algo: module.Incremental},
			{ID: "opcounters.getmore", Name: "getmore", Algo: module.Incremental},
			{ID: "opcounters.command", Name: "command", Algo: module.Incremental},
		},
	}
	chartGlobalLockCurrentQueue = module.Chart{
		ID:    "globalLockCurrentQueue",
		Title: "Current Queue Clients",
		Units: "clients",
		Fam:   "clients",
		Ctx:   "mongodb.currentQueue",
		Type:  module.Stacked,
		Dims: module.Dims{
			{ID: "globalLock.currentQueue.readers", Name: "readers", Algo: module.Absolute},
			{ID: "globalLock.currentQueue.writers", Name: "writers", Algo: module.Absolute},
		},
	}
	chartGlobalLockActiveClients = module.Chart{
		ID:    "globalLockActiveClients",
		Title: "Active Clients",
		Units: "clients",
		Fam:   "clients",
		Ctx:   "mongodb.currentQueue",
		Type:  module.Stacked,
		Dims: module.Dims{
			{ID: "globalLock.activeClients.readers", Name: "readers", Algo: module.Absolute},
			{ID: "globalLock.activeClients.writers", Name: "writers", Algo: module.Absolute},
		},
	}
	chartConnections = module.Chart{
		ID:    "connections",
		Title: "Connections",
		Units: "connections/s",
		Fam:   "connections",
		Ctx:   "mongodb.connections",
		Type:  module.Line,
		Dims: module.Dims{
			{ID: "connections.available", Name: "available", Algo: module.Absolute},
			{ID: "connections.current", Name: "current", Algo: module.Absolute},
			{ID: "connections.active", Name: "active", Algo: module.Absolute},
			{ID: "connections.threaded", Name: "threaded", Algo: module.Absolute},
			{ID: "connections.exhaustIsMaster", Name: "exhaustIsMaster", Algo: module.Absolute},
			{ID: "connections.exhaustHello", Name: "exhaustHello", Algo: module.Absolute},
			{ID: "connections.awaitingTopologyChanges", Name: "awaiting topology changes", Algo: module.Absolute},
		},
	}
	chartMemory = module.Chart{
		ID:    "mem",
		Title: "Memory",
		Units: "MiB",
		Fam:   "memory",
		Ctx:   "mongodb.memory",
		Type:  module.Line,
		Dims: module.Dims{
			{ID: "mem.resident", Name: "resident", Algo: module.Absolute},
			{ID: "mem.virtual", Name: "virtual", Algo: module.Absolute},
			{ID: "mem.mapped", Name: "mapped", Algo: module.Absolute},
			{ID: "mem.mappedWithJournal", Name: "mapped with journal", Algo: module.Absolute},
		},
	}
	chartAsserts = module.Chart{
		ID:    "asserts",
		Title: "Asserts",
		Units: "number",
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
	chartPageFaults = module.Chart{
		ID:    "page_faults",
		Title: "Page faults",
		Units: "page faults/s",
		Fam:   "page_faults",
		Ctx:   "mongodb.page_faults",
		Type:  module.Line,
		Dims: module.Dims{
			{ID: "extra_info.page_faults", Name: "Page Faults", Algo: module.Incremental},
		},
	}
	chartNetwork = module.Chart{
		ID:    "network",
		Title: "Network IO",
		Units: "bytes/s",
		Fam:   "network",
		Ctx:   "mongodb.network",
		Type:  module.Line,
		Dims: module.Dims{
			{ID: "network.bytesIn", Name: "Bytes In", Algo: module.Incremental},
			{ID: "network.bytesOut", Name: "Bytes Out", Algo: module.Incremental},
		},
	}
	chartNetworkRequests = module.Chart{
		ID:    "networkRequests",
		Title: "Network Requests",
		Units: "requests/s",
		Fam:   "network requests",
		Ctx:   "mongodb.networkRequests",
		Type:  module.Line,
		Dims: module.Dims{
			{ID: "network.numRequests", Name: "Requests", Algo: module.Incremental},
		},
	}
	chartTransactionsCurrent = module.Chart{
		ID:    "transactionsCurrent",
		Title: "Current Transactions",
		Units: "transactions",
		Fam:   "transactionsCurrent",
		Ctx:   "mongodb.transactionsCurrent",
		Type:  module.Line,
		Dims: module.Dims{
			{ID: "network.currentActive", Name: "current active", Algo: module.Absolute},
			{ID: "network.currentInactive", Name: "current inactive", Algo: module.Absolute},
			{ID: "network.currentOpen", Name: "current open", Algo: module.Absolute},
			{ID: "network.currentPrepared", Name: "current prepared", Algo: module.Absolute},
		},
	}
	chartCollections = module.Chart{
		ID:    "collections",
		Title: "Collections",
		Units: "collections",
		Fam:   "collections",
		Ctx:   "mongodb.collections",
		Type:  module.Line,
		Dims: module.Dims{
			{ID: "catalogStats.collections", Name: "collections", Algo: module.Incremental},
			{ID: "catalogStats.capped", Name: "capped", Algo: module.Incremental},
			{ID: "catalogStats.timeseries", Name: "timeseries", Algo: module.Incremental},
			{ID: "catalogStats.views", Name: "views", Algo: module.Incremental},
			{ID: "catalogStats.internalCollections", Name: "internalCollections", Algo: module.Incremental},
			{ID: "catalogStats.internalViews", Name: "internalViews", Algo: module.Incremental},
		},
	}
	chartFlowControl = module.Chart{
		ID:    "flowControl",
		Title: "Flow Control Stats",
		Units: "number",
		Fam:   "flowControl",
		Ctx:   "mongodb.flowControl",
		Type:  module.Line,
		Dims: module.Dims{
			{ID: "flowControl.timeAcquiringMicros", Name: "timeAcquiringMicros", Algo: module.Absolute},
			{ID: "flowControl.isLaggedTimeMicros", Name: "isLaggedTimeMicros", Algo: module.Absolute},
		},
	}
	chartGlobalLocks = module.Chart{
		ID:    "locks",
		Title: "Global Locks",
		Units: "locks",
		Fam:   "locks",
		Ctx:   "mongodb.locks",
		Type:  module.Line,
		Dims: module.Dims{
			{ID: "locks.Global.acquireCount.r", Name: "Global Read Locks", Algo: module.Incremental},
			{ID: "locks.Global.acquireCount.w", Name: "Global Write Locks", Algo: module.Incremental},
			{ID: "locks.Database.acquireCount.r", Name: "Database Read Locks", Algo: module.Incremental},
			{ID: "locks.Database.acquireCount.w", Name: "Database Write Locks", Algo: module.Incremental},
			{ID: "locks.Collection.acquireCount.r", Name: "Collection Read Locks", Algo: module.Incremental},
			{ID: "locks.Collection.acquireCount.w", Name: "Collection Write Locks", Algo: module.Incremental},
		},
	}
	chartMetricsCommands = module.Chart{
		ID:    "metricsCommand",
		Title: "Command Metrics",
		Units: "commands",
		Fam:   "commands",
		Ctx:   "mongodb.metricsCommand",
		Type:  module.Line,
		Dims: module.Dims{
			{ID: "metrics.commands.eval.total", Name: "", Algo: module.Incremental},
			{ID: "metrics.commands.delete.total", Name: "", Algo: module.Incremental},
			{ID: "metrics.commands.count.failed", Name: "", Algo: module.Incremental},
			{ID: "metrics.commands.createIndexes.", Name: "", Algo: module.Incremental},
			{ID: "metrics.commands.delete.failed", Name: "", Algo: module.Incremental},
			{ID: "metrics.commands.eval.failed", Name: "", Algo: module.Incremental},
			{ID: "metrics.commands.findAndModify", Name: "", Algo: module.Incremental},
			{ID: "metrics.commands.insert.failed", Name: "", Algo: module.Incremental},
		},
	}
	chartTcmallocGeneric = module.Chart{
		ID:    "tcmallocGeneric",
		Title: "Tcmalloc generic metrics",
		Units: "MiB",
		Fam:   "tcmalloc",
		Ctx:   "mongodb.tcmallocGeneric",
		Type:  module.Stacked,
		Dims: module.Dims{
			{ID: "tcmalloc.generic.current_allocated_bytes", Name: "current_allocated_bytes", Algo: module.Absolute, Div: 1 << 20},
			{ID: "tcmalloc.generic.heap_size", Name: "heap_size", Algo: module.Absolute, Div: 1 << 20},
		},
	}
	chartTcmalloc = module.Chart{
		ID:    "tcmalloc",
		Title: "Tcmalloc",
		Units: "KiB",
		Fam:   "tcmalloc",
		Ctx:   "mongodb.tcmalloc",
		Type:  module.Line,
		Dims: module.Dims{
			{ID: "tcmalloc.tcmalloc.pageheap_free_bytes", Name: "", Algo: module.Absolute, Div: 1024},
			{ID: "tcmalloc.tcmalloc.pageheap_unmapped_bytes", Name: "", Algo: module.Absolute, Div: 1024},
			{ID: "tcmalloc.tcmalloc.max_total_thread_cache_bytes", Name: "", Algo: module.Absolute, Div: 1024},
			{ID: "tcmalloc.tcmalloc.current_total_thread_cache_bytes", Name: "", Algo: module.Absolute, Div: 1024},
			{ID: "tcmalloc.tcmalloc.total_free_bytes", Name: "", Algo: module.Absolute, Div: 1024},
			{ID: "tcmalloc.tcmalloc.central_cache_free_bytes", Name: "", Algo: module.Absolute, Div: 1024},
			{ID: "tcmalloc.tcmalloc.transfer_cache_free_bytes", Name: "", Algo: module.Absolute, Div: 1024},
			{ID: "tcmalloc.tcmalloc.thread_cache_free_bytes", Name: "", Algo: module.Absolute, Div: 1024},
			{ID: "tcmalloc.tcmalloc.aggressive_memory_decommit", Name: "", Algo: module.Absolute, Div: 1024},
			{ID: "tcmalloc.tcmalloc.pageheap_committed_bytes", Name: "", Algo: module.Absolute, Div: 1024},
			{ID: "tcmalloc.tcmalloc.pageheap_total_commit_bytes", Name: "", Algo: module.Absolute, Div: 1024},
			{ID: "tcmalloc.tcmalloc.pageheap_total_decommit_bytes", Name: "", Algo: module.Absolute, Div: 1024},
			{ID: "tcmalloc.tcmalloc.pageheap_total_reserve_bytes", Name: "", Algo: module.Absolute, Div: 1024},
		},
	}
	chartWiredTigerBlockManager = module.Chart{
		ID:    "wiredtigerBlockManager",
		Title: "Wired Tiger Block Manager",
		Units: "KiB",
		Fam:   "wiredtiger",
		Ctx:   "mongodb.wiredTigerBlocks",
		Type:  module.Stacked,
		Dims: module.Dims{
			{ID: toID("wiredtiger.block-manager.bytes read"), Name: "bytes read", Algo: module.Absolute, Div: 1024},
			{ID: toID("wiredtiger.block-manager.bytes read via memory map API"), Name: "bytes read via memory map API", Algo: module.Absolute, Div: 1024},
			{ID: toID("wiredtiger.block-manager.bytes read via system call API"), Name: "bytes read via system call API", Algo: module.Absolute, Div: 1024},
			{ID: toID("wiredtiger.block-manager.bytes written"), Name: "bytes written", Algo: module.Absolute, Div: 1024},
			{ID: toID("wiredtiger.block-manager.bytes written for checkpoint"), Name: "bytes written for checkpoint", Algo: module.Absolute, Div: 1024},
			{ID: toID("wiredtiger.block-manager.bytes written via memory map API"), Name: "bytes written via memory map API", Algo: module.Absolute, Div: 1024},
		},
	}
	chartWiredTigerCache = module.Chart{
		ID:    "wiredtigerCache",
		Title: "Wired Tiger Cache",
		Units: "KiB",
		Fam:   "wiredtiger",
		Ctx:   "mongodb.wiredTigerCache",
		Type:  module.Line,
		Dims: module.Dims{
			{ID: toID("wiredtiger.cache.bytes allocated for updates"), Name: "bytes allocated for updates", Algo: module.Absolute, Div: 1024},
			{ID: toID("wiredtiger.cache.bytes read into cache"), Name: "bytes read into cache", Algo: module.Absolute, Div: 1024},
			{ID: toID("wiredtiger.cache.bytes written from cache"), Name: "bytes written from cache", Algo: module.Absolute, Div: 1024},
		},
	}
	chartWiredTigerCapacity = module.Chart{
		ID:    "wiredtigerCapacity",
		Title: "Wired Tiger Capacity",
		Units: "usec",
		Fam:   "wiredtiger",
		Ctx:   "mongodb.wiredTigerCapacity",
		Type:  module.Line,
		Dims: module.Dims{
			{ID: toID("wiredtiger.capacity.time waiting due to total capacity (usecs)"), Name: "time waiting due to total capacity (usecs)", Algo: module.Absolute},
			{ID: toID("wiredtiger.capacity.time waiting during checkpoint (usecs)"), Name: "time waiting during checkpoint (usecs)", Algo: module.Absolute},
			{ID: toID("wiredtiger.capacity.time waiting during eviction (usecs)"), Name: "time waiting during eviction (usecs)", Algo: module.Absolute},
			{ID: toID("wiredtiger.capacity.time waiting during logging (usecs)"), Name: "time waiting during logging (usecs)", Algo: module.Absolute},
			{ID: toID("wiredtiger.capacity.time waiting during read (usecs)"), Name: "time waiting during read (usecs)", Algo: module.Absolute},
		},
	}
	chartWiredTigerConnection = module.Chart{
		ID:    "wiredtigerConnection",
		Title: "Wired Tiger Connections",
		Units: "ops",
		Fam:   "wiredtiger",
		Ctx:   "mongodb.wiredTigerConnection",
		Type:  module.Line,
		Dims: module.Dims{
			{ID: toID("wiredtiger.connection.memory allocations"), Name: "memory allocations", Algo: module.Incremental},
			{ID: toID("wiredtiger.connection.memory frees"), Name: "memory frees", Algo: module.Incremental},
			{ID: toID("wiredtiger.connection.memory re-allocations"), Name: "memory re-allocations", Algo: module.Incremental},
		},
	}
	chartWiredTigerCursor = module.Chart{
		ID:    "wiredtigerCursor",
		Title: "Wired Tiger Cursor",
		Units: "calls",
		Fam:   "wiredtiger",
		Ctx:   "mongodb.wiredTigerCursor",
		Type:  module.Line,
		Dims: module.Dims{
			{ID: toID("wiredtiger.cursor.open cursor count"), Name: "open cursor count", Algo: module.Incremental},
			{ID: toID("wiredtiger.cursor.cached cursor count"), Name: "cached cursor count", Algo: module.Incremental},
			{ID: toID("wiredtiger.cursor.cursor bulk loaded cursor insert calls"), Name: "cursor bulk loaded cursor insert calls", Algo: module.Incremental},
			{ID: toID("wiredtiger.cursor.cursor close calls that result in cache"), Name: "cursor close calls that result in cache", Algo: module.Incremental},
			{ID: toID("wiredtiger.cursor.cursor create calls"), Name: "cursor create calls", Algo: module.Incremental},
			{ID: toID("wiredtiger.cursor.cursor insert calls"), Name: "cursor insert calls", Algo: module.Incremental},
			{ID: toID("wiredtiger.cursor.cursor modify calls"), Name: "cursor modify calls", Algo: module.Incremental},
			{ID: toID("wiredtiger.cursor.cursor next calls"), Name: "cursor next calls", Algo: module.Incremental},
			{ID: toID("wiredtiger.cursor.cursor operation restarted"), Name: "cursor operation restarted", Algo: module.Incremental},
			{ID: toID("wiredtiger.cursor.cursor prev calls"), Name: "cursor prev calls", Algo: module.Incremental},
			{ID: toID("wiredtiger.cursor.cursor remove calls"), Name: "cursor remove calls", Algo: module.Incremental},
			{ID: toID("wiredtiger.cursor.cursor remove key bytes removed"), Name: "cursor remove key bytes removed", Algo: module.Incremental},
			{ID: toID("wiredtiger.cursor.cursor reserve calls"), Name: "cursor reserve calls", Algo: module.Incremental},
			{ID: toID("wiredtiger.cursor.cursor reset calls"), Name: "cursor reset calls", Algo: module.Incremental},
			{ID: toID("wiredtiger.cursor.cursor search calls"), Name: "cursor search calls", Algo: module.Incremental},
			{ID: toID("wiredtiger.cursor.cursor search history store calls"), Name: "cursor search history store calls", Algo: module.Incremental},
			{ID: toID("wiredtiger.cursor.cursor search near calls"), Name: "cursor search near calls", Algo: module.Incremental},
			{ID: toID("wiredtiger.cursor.cursor sweep buckets"), Name: "cursor sweep buckets", Algo: module.Incremental},
			{ID: toID("wiredtiger.cursor.cursor sweep cursors closed"), Name: "cursor sweep cursors closed", Algo: module.Incremental},
			{ID: toID("wiredtiger.cursor.cursor sweep cursors examined"), Name: "cursor sweep cursors examined", Algo: module.Incremental},
			{ID: toID("wiredtiger.cursor.cursor sweeps"), Name: "cursor sweeps", Algo: module.Incremental},
			{ID: toID("wiredtiger.cursor.cursor truncate calls"), Name: "cursor truncate calls", Algo: module.Incremental},
			{ID: toID("wiredtiger.cursor.cursor update calls"), Name: "cursor update calls", Algo: module.Incremental},
			{ID: toID("wiredtiger.cursor.cursor update value size change"), Name: "cursor update value size change", Algo: module.Incremental},
		},
	}
	chartWiredTigerLock = module.Chart{
		ID:    "wiredtigerLock",
		Title: "Wired Tiger Lock",
		Units: "ops",
		Fam:   "wiredtiger",
		Ctx:   "mongodb.wiredTigerLock",
		Type:  module.Line,
		Dims: module.Dims{
			{ID: toID("wiredtiger.lock.checkpoint lock acquisitions"), Name: "checkpoint lock acquisitions", Algo: module.Incremental},
			{ID: toID("wiredtiger.lock.dhandle read lock acquisitions"), Name: "dhandle read lock acquisitions", Algo: module.Incremental},
			{ID: toID("wiredtiger.lock.dhandle write lock acquisitions"), Name: "dhandle write lock acquisitions", Algo: module.Incremental},
			{ID: toID("wiredtiger.lock.durable timestamp queue read lock acquisitions"), Name: "durable timestamp queue read lock acquisitions", Algo: module.Incremental},
			{ID: toID("wiredtiger.lock.durable timestamp queue write lock acquisitions"), Name: "durable timestamp queue write lock acquisitions", Algo: module.Incremental},
			{ID: toID("wiredtiger.lock.metadata lock acquisitions"), Name: "metadata lock acquisitions", Algo: module.Incremental},
			{ID: toID("wiredtiger.lock.read timestamp queue read lock acquisitions"), Name: "read timestamp queue read lock acquisitions", Algo: module.Incremental},
			{ID: toID("wiredtiger.lock.read timestamp queue write lock acquisitions"), Name: "read timestamp queue write lock acquisitions", Algo: module.Incremental},
			{ID: toID("wiredtiger.lock.schema lock acquisitions"), Name: "schema lock acquisitions", Algo: module.Incremental},
			{ID: toID("wiredtiger.lock.table read lock acquisitions"), Name: "table read lock acquisitions", Algo: module.Incremental},
			{ID: toID("wiredtiger.lock.table write lock acquisitions"), Name: "table write lock acquisitions", Algo: module.Incremental},
			{ID: toID("wiredtiger.lock.txn global read lock acquisitions"), Name: "txn global read lock acquisitions", Algo: module.Incremental},
		},
	}
	chartWiredTigerLockDuration = module.Chart{
		ID:    "wiredtigerLockDuration",
		Title: "Wired Tiger Lock Duration",
		Units: "usec",
		Fam:   "wiredtiger",
		Ctx:   "mongodb.wiredTigerLockDuration",
		Type:  module.Line,
		Dims: module.Dims{
			{ID: toID("wiredtiger.lock.checkpoint lock application thread wait time (usecs)"), Name: "checkpoint lock application thread wait time (usecs)", Algo: module.Absolute},
			{ID: toID("wiredtiger.lock.checkpoint lock internal thread wait time (usecs)"), Name: "checkpoint lock internal thread wait time (usecs)", Algo: module.Absolute},
			{ID: toID("wiredtiger.lock.dhandle lock application thread time waiting (usecs)"), Name: "dhandle lock application thread time waiting (usecs)", Algo: module.Absolute},
			{ID: toID("wiredtiger.lock.dhandle lock internal thread time waiting (usecs)"), Name: "dhandle lock internal thread time waiting (usecs)", Algo: module.Absolute},
			{ID: toID("wiredtiger.lock.durable timestamp queue lock application thread time waiting (usecs)"), Name: "durable timestamp queue lock application thread time waiting (usecs)", Algo: module.Absolute},
			{ID: toID("wiredtiger.lock.durable timestamp queue lock internal thread time waiting (usecs)"), Name: "durable timestamp queue lock internal thread time waiting (usecs)", Algo: module.Absolute},
			{ID: toID("wiredtiger.lock.metadata lock application thread wait time (usecs)"), Name: "metadata lock application thread wait time (usecs)", Algo: module.Absolute},
			{ID: toID("wiredtiger.lock.metadata lock internal thread wait time (usecs)"), Name: "metadata lock internal thread wait time (usecs)", Algo: module.Absolute},
			{ID: toID("wiredtiger.lock.read timestamp queue lock application thread time waiting (usecs)"), Name: "read timestamp queue lock application thread time waiting (usecs)", Algo: module.Absolute},
			{ID: toID("wiredtiger.lock.read timestamp queue lock internal thread time waiting (usecs)"), Name: "read timestamp queue lock internal thread time waiting (usecs)", Algo: module.Absolute},
			{ID: toID("wiredtiger.lock.schema lock application thread wait time (usecs)"), Name: "schema lock application thread wait time (usecs)", Algo: module.Absolute},
			{ID: toID("wiredtiger.lock.schema lock internal thread wait time (usecs)"), Name: "schema lock internal thread wait time (usecs)", Algo: module.Absolute},
			{ID: toID("wiredtiger.lock.table lock application thread time waiting for the table lock (usecs)"), Name: "table lock application thread time waiting for the table lock (usecs)", Algo: module.Absolute},
			{ID: toID("wiredtiger.lock.table lock internal thread time waiting for the table lock (usecs)"), Name: "table lock internal thread time waiting for the table lock (usecs)", Algo: module.Absolute},
			{ID: toID("wiredtiger.lock.txn global lock application thread time waiting (usecs)"), Name: "txn global lock application thread time waiting (usecs)", Algo: module.Absolute},
			{ID: toID("wiredtiger.lock.txn global lock internal thread time waiting (usecs)"), Name: "txn global lock internal thread time waiting (usecs)", Algo: module.Absolute},
		},
	}
	chartWiredTigerLogOps = module.Chart{
		ID:    "wiredtigerLogOps",
		Title: "Wired Tiger Log Operations",
		Units: "ops",
		Fam:   "wiredtiger",
		Ctx:   "mongodb.wiredTigerLogOps",
		Type:  module.Line,
		Dims: module.Dims{
			{ID: toID("wiredtiger.log.log flush operations"), Name: "log flush operations", Algo: module.Incremental},
			{ID: toID("wiredtiger.log.log force write operations"), Name: "log force write operations", Algo: module.Incremental},
			{ID: toID("wiredtiger.log.log force write operations skipped"), Name: "log force write operations skipped", Algo: module.Incremental},
			{ID: toID("wiredtiger.log.log scan operations"), Name: "log scan operations", Algo: module.Incremental},
			{ID: toID("wiredtiger.log.log sync operations"), Name: "log sync operations", Algo: module.Incremental},
			{ID: toID("wiredtiger.log.log sync_dir operations"), Name: "log sync_dir operations", Algo: module.Incremental},
			{ID: toID("wiredtiger.log.log write operations"), Name: "log write operations", Algo: module.Incremental},
		},
	}
	chartWiredTigerLogBytes = module.Chart{
		ID:    "wiredtigerLogOps",
		Title: "Wired Tiger Log Operations",
		Units: "bytes",
		Fam:   "wiredtiger",
		Ctx:   "mongodb.wiredTigerLogOps",
		Type:  module.Line,
		Dims: module.Dims{
			{ID: toID("wiredtiger.log.log bytes of payload data"), Name: "log bytes of payload data", Algo: module.Incremental},
			{ID: toID("wiredtiger.log.log bytes written"), Name: "log bytes written", Algo: module.Incremental},
			{ID: toID("wiredtiger.log.logging bytes consolidated"), Name: "logging bytes consolidated", Algo: module.Incremental},
			{ID: toID("wiredtiger.log.total log buffer size"), Name: "total log buffer size", Algo: module.Incremental},
		},
	}
	chartWiredTigerTransactions = module.Chart{
		ID:    "wiredtigerTransactions",
		Title: "Wired Tiger Log Transactions",
		Units: "transactions",
		Fam:   "wiredtiger",
		Ctx:   "mongodb.wiredTigerTransactions",
		Type:  module.Line,
		Dims: module.Dims{
			{ID: toID("wiredtiger.transaction.prepared transactions"), Name: "prepared transactions", Algo: module.Incremental},
			{ID: toID("wiredtiger.transaction.query timestamp calls"), Name: "query timestamp calls", Algo: module.Incremental},
			{ID: toID("wiredtiger.transaction.rollback to stable calls"), Name: "rollback to stable calls", Algo: module.Incremental},
			{ID: toID("wiredtiger.transaction.set timestamp calls"), Name: "set timestamp calls", Algo: module.Incremental},
			{ID: toID("wiredtiger.transaction.transaction begins"), Name: "transaction begins", Algo: module.Incremental},
			{ID: toID("wiredtiger.transaction.transaction sync calls"), Name: "transaction sync calls", Algo: module.Incremental},
			{ID: toID("wiredtiger.transaction.transactions committed"), Name: "transactions committed", Algo: module.Incremental},
			{ID: toID("wiredtiger.transaction.transactions rolled back"), Name: "transactions rolled back", Algo: module.Incremental},
		},
	}
)

func toID(in string) string {
	id := strings.ReplaceAll(in, " ", "_")
	id = strings.ReplaceAll(id, "-", "_")
	return id
}
