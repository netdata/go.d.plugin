package cockroachdb

import "github.com/netdata/go-orchestrator/module"

type (
	Charts = module.Charts
	Dims   = module.Dims
	Vars   = module.Vars
)

// TODO:
// following metrics are Gauges, but they are Counters?
// - rocksdb_block_cache_hits
// - rocksdb_block_cache_misses
// - rocksdb_compactions
// - rocksdb_flushes

// TODO:
// better grouping

var charts = Charts{
	{
		ID:    "total_storage_capacity",
		Title: "Total Storage Capacity",
		Units: "KiB",
		Fam:   "capacity",
		Ctx:   "cockroachdb.total_storage_capacity",
		Dims: Dims{
			{ID: "storage_capacity_total", Name: "total", Div: 1024},
		},
	},
	{
		ID:    "storage_capacity_usability",
		Title: "Storage Capacity Usability",
		Units: "KiB",
		Fam:   "capacity",
		Ctx:   "cockroachdb.storage_capacity_usability",
		Type:  module.Stacked,
		Dims: Dims{
			{ID: "storage_capacity_usable", Name: "usable", Div: 1024},
			{ID: "storage_capacity_unusable", Name: "unusable", Div: 1024},
		},
	},
	{
		ID:    "storage_usable_capacity",
		Title: "Storage Usable Capacity",
		Units: "KiB",
		Fam:   "capacity",
		Ctx:   "cockroachdb.storage_usable_capacity",
		Type:  module.Stacked,
		Dims: Dims{
			{ID: "storage_capacity_available", Name: "available", Div: 1024},
			{ID: "storage_capacity_used", Name: "used", Div: 1024},
		},
	},
	{
		ID:    "storage_used_capacity_percentage",
		Title: "Storage Used Capacity",
		Units: "percentage",
		Fam:   "capacity",
		Ctx:   "cockroachdb.storage_used_capacity_percentage",
		Dims: Dims{
			{ID: "storage_capacity_total_used_percentage", Name: "total", Div: percentagePrecision},
			{ID: "storage_capacity_usable_used_percentage", Name: "usable", Div: percentagePrecision},
		},
	},
	{
		ID:    "live_bytes",
		Title: "The Amount of Used Live Data",
		Units: "KiB",
		Fam:   "storage",
		Ctx:   "cockroachdb.live_bytes",
		Dims: Dims{
			{ID: "storage_live_bytes", Name: "applications", Div: 1024},
			{ID: "storage_sys_bytes", Name: "system", Div: 1024},
		},
	},
	{
		ID:    "rocksdb_read_amplification",
		Title: "RocksDB  Read Amplification",
		Units: "reads/query",
		Fam:   "rocksdb",
		Ctx:   "cockroachdb.rocksdb_read_amplification",
		Dims: Dims{
			{ID: "storage_rocksdb_read_amplification", Name: "reads"},
		},
	},
	{
		ID:    "rocksdb_table_operations",
		Title: "RocksDB Table Operations",
		Units: "operations",
		Fam:   "rocksdb",
		Ctx:   "cockroachdb.rocksdb_table_operations",
		Dims: Dims{
			{ID: "storage_rocksdb_compactions", Name: "compactions", Algo: module.Incremental},
			{ID: "storage_rocksdb_flushes", Name: "flushes", Algo: module.Incremental},
		},
	},
	{
		ID:    "rocksdb_cache_usage",
		Title: "RocksDB Block Cache Usage",
		Units: "KiB",
		Fam:   "rocksdb",
		Ctx:   "cockroachdb.rocksdb_cache_usage",
		Type:  module.Area,
		Dims: Dims{
			{ID: "storage_rocksdb_block_cache_bytes", Name: "used", Div: 1024},
		},
	},
	{
		ID:    "rocksdb_cache_operations",
		Title: "RocksDB Block Cache Operations",
		Units: "operations/s",
		Fam:   "rocksdb",
		Ctx:   "cockroachdb.rocksdb_cache_operations",
		Type:  module.Area,
		Dims: Dims{
			{ID: "storage_rocksdb_block_cache_hits", Name: "hits", Algo: module.Incremental},
			{ID: "storage_rocksdb_block_cache_misses", Name: "misses", Algo: module.Incremental},
		},
	},
	{
		ID:    "rocksdb_cache_hit_rate",
		Title: "RocksDB Block Cache Hit Rate",
		Units: "percentage",
		Fam:   "rocksdb",
		Ctx:   "cockroachdb.rocksdb_cache_hit_rate",
		Type:  module.Area,
		Dims: Dims{
			{ID: "storage_rocksdb_block_cache_hit_rate", Name: "hit rate", Div: percentagePrecision},
		},
	},
	{
		ID:    "rocksdb_sstables",
		Title: "RocksDB SSTables",
		Units: "sstables",
		Fam:   "rocksdb",
		Ctx:   "cockroachdb.rocksdb_sstables",
		Dims: Dims{
			{ID: "storage_rocksdb_num_sstables", Name: "sstables"},
		},
	},
	{
		ID:    "file_descriptors",
		Title: "File Descriptors Statistics",
		Units: "file descriptors",
		Fam:   "storage",
		Ctx:   "cockroachdb.file_descriptors",
		Dims: Dims{
			{ID: "storage_file_descriptors_open", Name: "open"},
		},
		Vars: Vars{
			{ID: "storage_file_descriptors_soft_limit"},
		},
	},
	{
		ID:    "timeseries_samples",
		Title: "Time Series Written Samples",
		Units: "samples/s",
		Fam:   "timeseries",
		Ctx:   "cockroachdb.timeseries_samples",
		Dims: Dims{
			{ID: "storage_timeseries_write_samples", Name: "written", Algo: module.Incremental},
		},
	},
	{
		ID:    "timeseries_write_errors",
		Title: "Time Series Write Errors",
		Units: "errors/s",
		Fam:   "timeseries",
		Ctx:   "cockroachdb.timeseries_write_errors",
		Dims: Dims{
			{ID: "storage_timeseries_write_errors", Name: "write", Algo: module.Incremental},
		},
	},
	{
		ID:    "timeseries_write_bytes",
		Title: "Time Series Bytes Written",
		Units: "KiB/s",
		Fam:   "timeseries",
		Ctx:   "cockroachdb.timeseries_write_bytes",
		Dims: Dims{
			{ID: "storage_timeseries_write_bytes", Name: "written", Algo: module.Incremental},
		},
	},

	{
		ID:    "nodes",
		Title: "Nodes",
		Units: "nodes",
		Fam:   "runtime",
		Ctx:   "cockroachdb.nodes",
		Dims: Dims{
			{ID: "runtime_live_nodes", Name: "live"},
		},
	},

	{
		ID:    "system_uptime",
		Title: "Nodes",
		Units: "seconds",
		Fam:   "runtime",
		Ctx:   "cockroachdb.uptime",
		Dims: Dims{
			{ID: "runtime_uptime", Name: "uptime"},
		},
	},

	{
		ID:    "rss_memory_usage",
		Title: "RSS Memory Usage",
		Units: "KiB",
		Fam:   "runtime",
		Ctx:   "cockroachdb.rss_memory_usage",
		Dims: Dims{
			{ID: "runtime_memory_rss", Name: "rss", Div: 1024},
		},
	},
	{
		ID:    "code_memory_usage",
		Title: "GO/CGO Memory Usage",
		Units: "KiB",
		Fam:   "runtime",
		Ctx:   "cockroachdb.code_memory_usage",
		Type:  module.Stacked,
		Dims: Dims{
			{ID: "runtime_memory_go_alloc_bytes", Name: "go", Div: 1024},
			{ID: "runtime_memory_cgo_alloc_bytes", Name: "cgo", Div: 1024},
		},
	},
	{
		ID:    "code_memory_allocations",
		Title: "GO/CGO Memory Allocations",
		Units: "KiB/s",
		Fam:   "runtime",
		Ctx:   "cockroachdb.code_memory_allocations",
		Type:  module.Stacked,
		Dims: Dims{
			{ID: "runtime_memory_go_total_bytes", Name: "go", Div: 1024, Algo: module.Incremental},
			{ID: "runtime_memory_cgo_total_bytes", Name: "cgo", Div: 1024, Algo: module.Incremental},
		},
	},
}
