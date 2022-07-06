// SPDX-License-Identifier: GPL-3.0-or-later

package elasticsearch

import (
	"github.com/netdata/go.d.plugin/agent/module"
)

type (
	Charts = module.Charts
	Dims   = module.Dims
	Vars   = module.Vars
)

// TODO: indices operations charts: query_latency, fetch_latency, index_latency, flush_latency
// if they are needed at all..

var (
	nodeCharts = Charts{
		// Indices Indexing
		{
			ID:    "node_indices_indexing_operations",
			Title: "Indexing Operations",
			Units: "operations/s",
			Fam:   "indices indexing",
			Ctx:   "elasticsearch.node_indices_indexing",
			Dims: Dims{
				{ID: "node_indices_indexing_index_total", Name: "index", Algo: module.Incremental},
			},
		},
		{
			ID:    "node_indices_indexing_operations_current",
			Title: "Indexing Operations Current",
			Units: "operations",
			Fam:   "indices indexing",
			Ctx:   "elasticsearch.node_indices_indexing_current",
			Dims: Dims{
				{ID: "node_indices_indexing_index_current", Name: "index"},
			},
		},
		{
			ID:    "node_indices_indexing_operations_time",
			Title: "Time Spent On Indexing Operations",
			Units: "milliseconds",
			Fam:   "indices indexing",
			Ctx:   "elasticsearch.node_indices_indexing_time",
			Dims: Dims{
				{ID: "node_indices_indexing_index_time_in_millis", Name: "index", Algo: module.Incremental},
			},
		},
		// Indices Search
		{
			ID:    "node_indices_search_operations",
			Title: "Search Operations",
			Units: "operations/s",
			Fam:   "indices search",
			Ctx:   "elasticsearch.node_indices_search",
			Type:  module.Stacked,
			Dims: Dims{
				{ID: "node_indices_search_query_total", Name: "queries", Algo: module.Incremental},
				{ID: "node_indices_search_fetch_total", Name: "fetches", Algo: module.Incremental},
			},
		},
		{
			ID:    "node_indices_search_operations_current",
			Title: "Search Operations Current",
			Units: "operations",
			Fam:   "indices search",
			Ctx:   "elasticsearch.node_indices_search_current",
			Type:  module.Stacked,
			Dims: Dims{
				{ID: "node_indices_search_query_current", Name: "queries"},
				{ID: "node_indices_search_fetch_current", Name: "fetches"},
			},
		},
		{
			ID:    "node_indices_search_operations_time",
			Title: "Time Spent On Search Operations",
			Units: "milliseconds",
			Fam:   "indices search",
			Ctx:   "elasticsearch.node_indices_search_time",
			Type:  module.Stacked,
			Dims: Dims{
				{ID: "node_indices_search_query_time_in_millis", Name: "query", Algo: module.Incremental},
				{ID: "node_indices_search_fetch_time_in_millis", Name: "fetch", Algo: module.Incremental},
			},
		},
		// Indices Refresh
		{
			ID:    "node_indices_refresh_operations",
			Title: "Refresh Operations",
			Units: "operations/s",
			Fam:   "indices refresh",
			Ctx:   "elasticsearch.node_indices_refresh",
			Dims: Dims{
				{ID: "node_indices_refresh_total", Name: "refresh", Algo: module.Incremental},
			},
		},
		{
			ID:    "node_indices_refresh_operations_time",
			Title: "Time Spent On Refresh Operations",
			Units: "milliseconds",
			Fam:   "indices refresh",
			Ctx:   "elasticsearch.node_indices_refresh_time",
			Dims: Dims{
				{ID: "node_indices_refresh_total_time_in_millis", Name: "refresh", Algo: module.Incremental},
			},
		},
		// Indices Flush
		{
			ID:    "node_indices_flush_operations",
			Title: "Flush Operations",
			Units: "operations/s",
			Fam:   "indices flush",
			Ctx:   "elasticsearch.node_indices_flush",
			Dims: Dims{
				{ID: "node_indices_flush_total", Name: "flush", Algo: module.Incremental},
			},
		},
		{
			ID:    "node_indices_flush_operations_time",
			Title: "Time Spent On Flush Operations",
			Units: "milliseconds",
			Fam:   "indices flush",
			Ctx:   "elasticsearch.node_indices_flush_time",
			Dims: Dims{
				{ID: "node_indices_flush_total_time_in_millis", Name: "flush", Algo: module.Incremental},
			},
		},
		// Indices Fielddata
		{
			ID:    "node_indices_fielddata_memory_usage",
			Title: "Fielddata Cache Memory Usage",
			Units: "bytes",
			Fam:   "indices fielddata",
			Ctx:   "elasticsearch.node_indices_fielddata_memory_usage",
			Type:  module.Area,
			Dims: Dims{
				{ID: "node_indices_fielddata_memory_size_in_bytes", Name: "used"},
			},
		},
		{
			ID:    "node_indices_fielddata_evictions",
			Title: "Fielddata Evictions",
			Units: "operations/s",
			Fam:   "indices fielddata",
			Ctx:   "elasticsearch.node_indices_fielddata_evictions",
			Dims: Dims{
				{ID: "node_indices_fielddata_evictions", Name: "evictions", Algo: module.Incremental},
			},
		},
		// Indices Segments
		{
			ID:    "node_indices_segments_count",
			Title: "Segments Count",
			Units: "segments",
			Fam:   "indices segments",
			Ctx:   "elasticsearch.node_indices_segments_count",
			Dims: Dims{
				{ID: "node_indices_segments_count", Name: "segments"},
			},
		},
		{
			ID:    "node_indices_segments_memory_usage_total",
			Title: "Segments Memory Usage Total",
			Units: "bytes",
			Fam:   "indices segments",
			Ctx:   "elasticsearch.node_indices_segments_memory_usage_total",
			Dims: Dims{
				{ID: "node_indices_segments_memory_in_bytes", Name: "used"},
			},
		},
		{
			ID:    "node_indices_segments_memory_usage",
			Title: "Segments Memory Usage",
			Units: "bytes",
			Fam:   "indices segments",
			Ctx:   "elasticsearch.node_indices_segments_memory_usage",
			Type:  module.Stacked,
			Dims: Dims{
				{ID: "node_indices_segments_terms_memory_in_bytes", Name: "terms"},
				{ID: "node_indices_segments_stored_fields_memory_in_bytes", Name: "stored_fields"},
				{ID: "node_indices_segments_term_vectors_memory_in_bytes", Name: "term_vectors"},
				{ID: "node_indices_segments_norms_memory_in_bytes", Name: "norms"},
				{ID: "node_indices_segments_points_memory_in_bytes", Name: "points"},
				{ID: "node_indices_segments_doc_values_memory_in_bytes", Name: "doc_values"},
				{ID: "node_indices_segments_index_writer_memory_in_bytes", Name: "index_writer"},
				{ID: "node_indices_segments_version_map_memory_in_bytes", Name: "version_map"},
				{ID: "node_indices_segments_fixed_bit_set_memory_in_bytes", Name: "fixed_bit_set"},
			},
		},
		// Indices Translog
		{
			ID:    "node_indices_translog_operations",
			Title: "Translog Operations",
			Units: "operations",
			Fam:   "indices translog",
			Ctx:   "elasticsearch.node_indices_translog_operations",
			Type:  module.Area,
			Dims: Dims{
				{ID: "node_indices_translog_operations", Name: "total"},
				{ID: "node_indices_translog_uncommitted_operations", Name: "uncommitted"},
			},
		},
		{
			ID:    "node_index_translog_size",
			Title: "Translog Size",
			Units: "bytes",
			Fam:   "indices translog",
			Ctx:   "elasticsearch.node_indices_translog_size",
			Type:  module.Area,
			Dims: Dims{
				{ID: "node_indices_translog_size_in_bytes", Name: "total"},
				{ID: "node_indices_translog_uncommitted_size_in_bytes", Name: "uncommitted"},
			},
		},
		// Process
		{
			ID:    "node_file_descriptors",
			Title: "Process File Descriptors",
			Units: "fd",
			Fam:   "process",
			Ctx:   "elasticsearch.node_file_descriptors",
			Dims: Dims{
				{ID: "node_process_open_file_descriptors", Name: "open"},
			},
			Vars: Vars{
				{ID: "node_process_max_file_descriptors"},
			},
		},
		// JVM
		{
			ID:    "node_jvm_mem_heap",
			Title: "JVM Heap Percentage Currently in Use",
			Units: "percentage",
			Fam:   "jvm",
			Ctx:   "elasticsearch.node_jvm_heap",
			Type:  module.Area,
			Dims: Dims{
				{ID: "node_jvm_mem_heap_used_percent", Name: "inuse"},
			},
		},
		{
			ID:    "node_jvm_mem_heap_bytes",
			Title: "JVM Heap Commit And Usage",
			Units: "bytes",
			Fam:   "jvm",
			Ctx:   "elasticsearch.node_jvm_heap_bytes",
			Type:  module.Area,
			Dims: Dims{
				{ID: "node_jvm_mem_heap_committed_in_bytes", Name: "committed"},
				{ID: "node_jvm_mem_heap_used_in_bytes", Name: "used"},
			},
		},
		{
			ID:    "node_jvm_buffer_pools_count",
			Title: "JVM Buffer Pools Count",
			Units: "pools",
			Fam:   "jvm",
			Ctx:   "elasticsearch.node_jvm_buffer_pools_count",
			Dims: Dims{
				{ID: "node_jvm_buffer_pools_direct_count", Name: "direct"},
				{ID: "node_jvm_buffer_pools_mapped_count", Name: "mapped"},
			},
		},
		{
			ID:    "node_jvm_buffer_pool_direct_memory",
			Title: "JVM Buffer Pool Direct Memory",
			Units: "bytes",
			Fam:   "jvm",
			Ctx:   "elasticsearch.node_jvm_buffer_pool_direct_memory",
			Type:  module.Area,
			Dims: Dims{
				{ID: "node_jvm_buffer_pools_direct_total_capacity_in_bytes", Name: "total"},
				{ID: "node_jvm_buffer_pools_direct_used_in_bytes", Name: "used"},
			},
		},
		{
			ID:    "node_jvm_buffer_pool_mapped_memory",
			Title: "JVM Buffer Pool Mapped Memory",
			Units: "bytes",
			Fam:   "jvm",
			Ctx:   "elasticsearch.node_jvm_buffer_pool_mapped_memory",
			Type:  module.Area,
			Dims: Dims{
				{ID: "node_jvm_buffer_pools_mapped_total_capacity_in_bytes", Name: "total"},
				{ID: "node_jvm_buffer_pools_mapped_used_in_bytes", Name: "used"},
			},
		},
		{
			ID:    "node_jvm_gc_count",
			Title: "JVM Garbage Collections",
			Units: "gc/s",
			Fam:   "jvm",
			Ctx:   "elasticsearch.node_jvm_gc_count",
			Type:  module.Stacked,
			Dims: Dims{
				{ID: "node_jvm_gc_collectors_young_collection_count", Name: "young", Algo: module.Incremental},
				{ID: "node_jvm_gc_collectors_old_collection_count", Name: "old", Algo: module.Incremental},
			},
		},
		{
			ID:    "node_jvm_gc_time",
			Title: "JVM Time Spent On Garbage Collections",
			Units: "milliseconds",
			Fam:   "jvm",
			Ctx:   "elasticsearch.node_jvm_gc_time",
			Type:  module.Stacked,
			Dims: Dims{
				{ID: "node_jvm_gc_collectors_young_collection_time_in_millis", Name: "young", Algo: module.Incremental},
				{ID: "node_jvm_gc_collectors_old_collection_time_in_millis", Name: "old", Algo: module.Incremental},
			},
		},
		// Thread Pool
		{
			ID:    "node_thread_pool_queued",
			Title: "Thread Pool Queued Threads Count",
			Units: "threads",
			Fam:   "thread pool",
			Ctx:   "elasticsearch.node_thread_pool_queued",
			Type:  module.Stacked,
			Dims: Dims{
				{ID: "node_thread_pool_generic_queue", Name: "generic"},
				{ID: "node_thread_pool_search_queue", Name: "search"},
				{ID: "node_thread_pool_search_throttled_queue", Name: "search_throttled"},
				{ID: "node_thread_pool_get_queue", Name: "get"},
				{ID: "node_thread_pool_analyze_queue", Name: "analyze"},
				{ID: "node_thread_pool_write_queue", Name: "write"},
				{ID: "node_thread_pool_snapshot_queue", Name: "snapshot"},
				{ID: "node_thread_pool_warmer_queue", Name: "warmer"},
				{ID: "node_thread_pool_refresh_queue", Name: "refresh"},
				{ID: "node_thread_pool_listener_queue", Name: "listener"},
				{ID: "node_thread_pool_fetch_shard_started_queue", Name: "fetch_shard_started"},
				{ID: "node_thread_pool_fetch_shard_store_queue", Name: "fetch_shard_store"},
				{ID: "node_thread_pool_flush_queue", Name: "flush"},
				{ID: "node_thread_pool_force_merge_queue", Name: "force_merge"},
				{ID: "node_thread_pool_management_queue", Name: "management"},
			},
		},
		{
			ID:    "node_thread_pool_rejected",
			Title: "Thread Pool Rejected Threads Count",
			Units: "threads",
			Fam:   "thread pool",
			Ctx:   "elasticsearch.node_thread_pool_rejected",
			Type:  module.Stacked,
			Dims: Dims{
				{ID: "node_thread_pool_generic_rejected", Name: "generic"},
				{ID: "node_thread_pool_search_rejected", Name: "search"},
				{ID: "node_thread_pool_search_throttled_rejected", Name: "search_throttled"},
				{ID: "node_thread_pool_get_rejected", Name: "get"},
				{ID: "node_thread_pool_analyze_rejected", Name: "analyze"},
				{ID: "node_thread_pool_write_rejected", Name: "write"},
				{ID: "node_thread_pool_snapshot_rejected", Name: "snapshot"},
				{ID: "node_thread_pool_warmer_rejected", Name: "warmer"},
				{ID: "node_thread_pool_refresh_rejected", Name: "refresh"},
				{ID: "node_thread_pool_listener_rejected", Name: "listener"},
				{ID: "node_thread_pool_fetch_shard_started_rejected", Name: "fetch_shard_started"},
				{ID: "node_thread_pool_fetch_shard_store_rejected", Name: "fetch_shard_store"},
				{ID: "node_thread_pool_flush_rejected", Name: "flush"},
				{ID: "node_thread_pool_force_merge_rejected", Name: "force_merge"},
				{ID: "node_thread_pool_management_rejected", Name: "management"},
			},
		},
		// Transport
		{
			ID:    "cluster_communication_packets",
			Title: "Cluster Communication",
			Units: "pps",
			Fam:   "transport",
			Ctx:   "elasticsearch.cluster_communication_packets",
			Dims: Dims{
				{ID: "node_transport_rx_count", Name: "received", Algo: module.Incremental},
				{ID: "node_transport_tx_count", Name: "sent", Mul: -1, Algo: module.Incremental},
			},
		},
		{
			ID:    "cluster_communication",
			Title: "Cluster Communication Bandwidth",
			Units: "bytes/s",
			Fam:   "transport",
			Ctx:   "elasticsearch.cluster_communication",
			Dims: Dims{
				{ID: "node_transport_rx_size_in_bytes", Name: "received", Algo: module.Incremental},
				{ID: "node_transport_tx_size_in_bytes", Name: "sent", Mul: -1, Algo: module.Incremental},
			},
		},
		// HTTP
		{
			ID:    "http_connections",
			Title: "HTTP Connections",
			Units: "connections",
			Fam:   "http",
			Ctx:   "elasticsearch.http_connections",
			Dims: Dims{
				{ID: "node_http_current_open", Name: "open"},
			},
		},
		// Breakers
		{
			ID:    "breakers_trips",
			Title: "Circuit Breaker Trips Count",
			Units: "trips/s",
			Fam:   "circuit breakers",
			Ctx:   "elasticsearch.breakers_trips",
			Type:  module.Stacked,
			Dims: Dims{
				{ID: "node_breakers_request_tripped", Name: "requests", Algo: module.Incremental},
				{ID: "node_breakers_fielddata_tripped", Name: "fielddata", Algo: module.Incremental},
				{ID: "node_breakers_in_flight_requests_tripped", Name: "in_flight_requests", Algo: module.Incremental},
				{ID: "node_breakers_model_inference_tripped", Name: "model_inference", Algo: module.Incremental},
				{ID: "node_breakers_accounting_tripped", Name: "accounting", Algo: module.Incremental},
				{ID: "node_breakers_parent_tripped", Name: "parent", Algo: module.Incremental},
			},
		},
	}
)

var nodeIndicesStatsCharts = Charts{
	{
		ID:    "node_index_health",
		Title: "Index Health (0: green, 1: yellow, 2: red)",
		Units: "status",
		Fam:   "indices stats",
		Ctx:   "elasticsearch.node_index_health",
	},
	{
		ID:    "node_index_shards_count",
		Title: "Index Shards Count",
		Units: "shards",
		Fam:   "indices stats",
		Ctx:   "elasticsearch.node_index_shards_count",
		Type:  module.Stacked,
	},
	{
		ID:    "node_index_docs_count",
		Title: "Index Docs Count",
		Units: "docs",
		Fam:   "indices stats",
		Ctx:   "elasticsearch.node_index_docs_count",
		Type:  module.Stacked,
	},
	{
		ID:    "node_index_store_size",
		Title: "Index Store Size",
		Units: "bytes",
		Fam:   "indices stats",
		Ctx:   "elasticsearch.node_index_store_size",
		Type:  module.Stacked,
	},
}

var clusterHealthCharts = Charts{
	{
		ID:    "cluster_status",
		Title: "Cluster Status (0: green, 1: yellow, 2: red)",
		Units: "status",
		Fam:   "cluster health",
		Ctx:   "elasticsearch.cluster_health_status",
		Dims: Dims{
			{ID: "cluster_status", Name: "status"},
		},
	},
	{
		ID:    "cluster_number_of_nodes",
		Title: "Cluster Nodes Count",
		Units: "nodes",
		Fam:   "cluster health",
		Ctx:   "elasticsearch.cluster_number_of_nodes",
		Dims: Dims{
			{ID: "cluster_number_of_nodes", Name: "nodes"},
			{ID: "cluster_number_of_data_nodes", Name: "data_nodes"},
		},
	},
	{
		ID:    "cluster_shards_count",
		Title: "Cluster Shards Count",
		Units: "shards",
		Fam:   "cluster health",
		Ctx:   "elasticsearch.cluster_shards_count",
		Dims: Dims{
			{ID: "cluster_active_primary_shards", Name: "active_primary"},
			{ID: "cluster_active_shards", Name: "active"},
			{ID: "cluster_relocating_shards", Name: "relocating"},
			{ID: "cluster_initializing_shards", Name: "initializing"},
			{ID: "cluster_unassigned_shards", Name: "unassigned"},
			{ID: "cluster_delayed_unassigned_shards", Name: "delayed_unassigned"},
		},
	},
	{
		ID:    "cluster_pending_tasks",
		Title: "Cluster Pending Tasks",
		Units: "tasks",
		Fam:   "cluster health",
		Ctx:   "elasticsearch.cluster_pending_tasks",
		Dims: Dims{
			{ID: "cluster_number_of_pending_tasks", Name: "pending"},
		},
	},
	{
		ID:    "cluster_number_of_in_flight_fetch",
		Title: "Cluster Unfinished Fetches",
		Units: "fetches",
		Fam:   "cluster health",
		Ctx:   "elasticsearch.cluster_number_of_in_flight_fetch",
		Dims: Dims{
			{ID: "cluster_number_of_in_flight_fetch", Name: "in_flight_fetch"},
		},
	},
}

var clusterStatsCharts = Charts{
	{
		ID:    "cluster_indices_count",
		Title: "Cluster Indices Count",
		Units: "indices",
		Fam:   "cluster stats",
		Ctx:   "elasticsearch.cluster_indices_count",
		Dims: Dims{
			{ID: "cluster_indices_count", Name: "indices"},
		},
	},
	{
		ID:    "cluster_indices_shards_count",
		Title: "Cluster Indices Shards Count",
		Units: "shards",
		Fam:   "cluster stats",
		Ctx:   "elasticsearch.cluster_indices_shards_count",
		Dims: Dims{
			{ID: "cluster_indices_shards_total", Name: "total"},
			{ID: "cluster_indices_shards_primaries", Name: "primaries"},
			{ID: "cluster_indices_shards_replication", Name: "replication"},
		},
	},
	{
		ID:    "cluster_indices_docs_count",
		Title: "Cluster Indices Docs Count",
		Units: "docs",
		Fam:   "cluster stats",
		Ctx:   "elasticsearch.cluster_indices_docs_count",
		Dims: Dims{
			{ID: "cluster_indices_docs_count", Name: "docs"},
		},
	},
	{
		ID:    "cluster_indices_store_size",
		Title: "Cluster Indices Store Size",
		Units: "bytes",
		Fam:   "cluster stats",
		Ctx:   "elasticsearch.cluster_indices_store_size",
		Dims: Dims{
			{ID: "cluster_indices_store_size_in_bytes", Name: "size"},
		},
	},
	{
		ID:    "cluster_indices_query_cache",
		Title: "Cluster Indices Query Cache",
		Units: "events/s",
		Fam:   "cluster stats",
		Ctx:   "elasticsearch.cluster_indices_query_cache",
		Type:  module.Stacked,
		Dims: Dims{
			{ID: "cluster_indices_query_cache_hit_count", Name: "hit", Algo: module.Incremental},
			{ID: "cluster_indices_query_cache_miss_count", Name: "miss", Algo: module.Incremental},
		},
	},
	{
		ID:    "cluster_nodes_by_role_count",
		Title: "Cluster Nodes By Role Count",
		Units: "nodes",
		Fam:   "cluster stats",
		Ctx:   "elasticsearch.cluster_nodes_by_role_count",
		Dims: Dims{
			{ID: "cluster_nodes_count_coordinating_only", Name: "coordinating_only"},
			{ID: "cluster_nodes_count_data", Name: "data"},
			{ID: "cluster_nodes_count_ingest", Name: "ingest"},
			{ID: "cluster_nodes_count_master", Name: "master"},
			{ID: "cluster_nodes_count_ml", Name: "ml"},
			{ID: "cluster_nodes_count_remote_cluster_client", Name: "remote_cluster_client"},
			{ID: "cluster_nodes_count_voting_only", Name: "voting_only"},
		},
	},
}
