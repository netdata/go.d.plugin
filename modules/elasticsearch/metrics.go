package elasticsearch

// https://www.elastic.co/guide/en/elasticsearch/reference/current/docker.html

type esMetrics struct {
	LocalNodeStats *esNodeStats     `stm:"node_stats"`
	ClusterHealth  *esClusterHealth `stm:"cluster_health"`
	ClusterStats   *esClusterStats  `stm:"cluster_stats"`
}

// https://www.elastic.co/guide/en/elasticsearch/reference/current/cluster-nodes-stats.html
type (
	esNodeStats struct {
		Indices    esNodeIndicesStats    `stm:"indices"`
		JVM        esNodeJVMStats        `stm:"jvm"`
		ThreadPool esNodeThreadPoolStats `stm:"thread_pool" json:"thread_pool"`
	}

	esNodeIndicesStats struct {
		Indexing struct {
			IndexTotal        int `stm:"index_total" json:"index_total"`
			IndexCurrent      int `stm:"index_current" json:"index_current"`
			IndexTimeInMillis int `stm:"index_time_in_millis" json:"index_time_in_millis"`
		} `stm:"indexing"`
		Search struct {
			FetchTotal        int `stm:"fetch_total" json:"fetch_total"`
			FetchCurrent      int `stm:"fetch_current" json:"fetch_current"`
			FetchTimeInMillis int `stm:"fetch_time_in_millis" json:"fetch_time_in_millis"`
			QueryTotal        int `stm:"query_total" json:"query_total"`
			QueryCurrent      int `stm:"query_current" json:"query_current"`
			QueryTimeInMillis int `stm:"query_time_in_millis" json:"query_time_in_millis"`
		} `stm:"search"`
		Refresh struct {
			Total        int `stm:"total"`
			TimeInMillis int `stm:"total_time_in_millis" json:"total_time_in_millis"`
		} `stm:"refresh"`
		Flush struct {
			Total        int `stm:"total"`
			TimeInMillis int `stm:"total_time_in_millis" json:"total_time_in_millis"`
		} `stm:"refresh"`
		Segments struct {
			Count                     int `json:"count"`
			TermsMemoryInBytes        int `stm:"terms_memory_in_bytes" json:"terms_memory_in_bytes"`
			StoredFieldsMemoryInBytes int `stm:"stored_fields_memory_in_bytes" json:"stored_fields_memory_in_bytes"`
			TermVectorsMemoryInBytes  int `stm:"term_vectors_memory_in_bytes" json:"term_vectors_memory_in_bytes"`
			NormsMemoryInBytes        int `stm:"norms_memory_in_bytes" json:"norms_memory_in_bytes"`
			PointsMemoryInBytes       int `stm:"points_memory_in_bytes" json:"points_memory_in_bytes"`
			DocValuesMemoryInBytes    int `stm:"doc_values_memory_in_bytes" json:"doc_values_memory_in_bytes"`
			IndexWriterMemoryInBytes  int `stm:"index_writer_memory_in_bytes" json:"index_writer_memory_in_bytes"`
			VersionMapMemoryInBytes   int `stm:"version_map_memory_in_bytes" json:"version_map_memory_in_bytes"`
			FixedBitSetMemoryInBytes  int `stm:"fixed_bit_set_memory_in_bytes" json:"fixed_bit_set_memory_in_bytes"`
		} `stm:"segments"`
		Translog struct {
			Operations             int `stm:"operations"`
			SizeInBytes            int `stm:"size_in_bytes" json:"size_in_bytes"`
			UncommittedOperations  int `stm:"uncommitted_operations" json:"uncommitted_operations"`
			UncommittedSizeInBytes int `stm:"uncommitted_size_in_bytes" json:"uncommitted_size_in_bytes"`
		} `stm:"translog"`
	}

	esNodeJVMStats struct {
		Mem struct {
			HeapUsedPercent      int `stm:"heap_used_percent" json:"heap_used_percent"`
			HeapUsedInBytes      int `stm:"heap_used_in_bytes" json:"heap_used_in_bytes"`
			HeapCommittedInBytes int `stm:"heap_committed_in_bytes" json:"heap_committed_in_bytes"`
		} `stm:"mem"`
		GC struct {
			Collectors struct {
				Young struct {
					CollectionCount        int `stm:"collection_count" json:"collection_count"`
					CollectionTimeInMillis int `stm:"collection_time_in_millis" json:"collection_time_in_millis"`
				} `stm:"young"`
				Old struct {
					CollectionCount        int `stm:"collection_count" json:"collection_count"`
					CollectionTimeInMillis int `stm:"collection_time_in_millis" json:"collection_time_in_millis"`
				} `stm:"old"`
			} `stm:"collectors"`
		} `stm:"gc"`
		BufferPool struct {
			Mapped struct {
				Count                int `stm:"count"`
				UsedInBytes          int `stm:"used_in_bytes" json:"used_in_bytes"`
				TotalCapacityInBytes int `stm:"total_capacity_in_bytes" json:"total_capacity_in_bytes"`
			} `stm:"mapped"`
			Direct struct {
				Count                int `stm:"count"`
				UsedInBytes          int `stm:"used_in_bytes" json:"used_in_bytes"`
				TotalCapacityInBytes int `stm:"total_capacity_in_bytes" json:"total_capacity_in_bytes"`
			} `stm:"direct"`
		} `stm:"buffer_pools" json:"buffer_pools"`
	}

	esNodeThreadPoolStats struct {
		Search struct {
			Queue    int `stm:"queue"`
			Rejected int `stm:"rejected"`
		} `stm:"search"`
		Write struct {
			Queue    int `stm:"queue"`
			Rejected int `stm:"rejected"`
		} `stm:"write"`
	}
)

// https://www.elastic.co/guide/en/elasticsearch/reference/current/cluster-health.html
type esClusterHealth struct {
	//Status                      string `stm:"status"`
	NumOfNodes                  int `stm:"number_of_nodes" json:"number_of_nodes"`
	NumOfDataNodes              int `stm:"number_of_data_nodes" json:"number_of_data_nodes"`
	ActiveShards                int `stm:"active_shards" json:"active_shards"`
	RelocatingShards            int `stm:"relocating_shards" json:"relocating_shards"`
	InitializingShards          int `stm:"initializing_shards" json:"initializing_shards"`
	UnassignedShards            int `stm:"unassigned_shards" json:"unassigned_shards"`
	DelayedUnassignedShards     int `stm:"delayed_unassigned_shards" json:"delayed_unassigned_shards"`
	NumOfPendingTasks           int `stm:"number_of_pending_tasks" json:"number_of_pending_tasks"`
	NumOfInFlightFetch          int `stm:"number_of_in_flight_fetch" json:"number_of_in_flight_fetch"`
	ActiveShardsPercentAsNumber int `stm:"active_shards_percent_as_number" json:"active_shards_percent_as_number"`
}

// https://www.elastic.co/guide/en/elasticsearch/reference/current/cluster-stats.html
type (
	esClusterStats struct {
		Nodes   esClusterNodesStats   `stm:"nodes"`
		Indices esClusterIndicesStats `stm:"indices"`
	}

	esClusterNodesStats struct {
		Count struct {
			Data             int `stm:"data"`
			Master           int `stm:"master"`
			Total            int `stm:"total"`
			CoordinatingOnly int `stm:"coordinating_only" json:"coordinating_only"`
			Ingest           int `stm:"ingest"`
		} `stm:"count"`
	}

	esClusterIndicesStats struct {
		Count int `stm:"count"`
		Docs  struct {
			Count int `stm:"count"`
		} `stm:"docs"`
		QueryCache struct {
			HitCount  int `stm:"hit_count" json:"hit_count"`
			MissCount int `stm:"miss_count" json:"miss_count"`
		} `stm:"query_cache" json:"query_cache"`
		Store struct {
			SizeInBytes int `stm:"size_in_bytes" json:"size_in_bytes"`
		} `stm:"store"`
		Shards struct {
			Total int `stm:"total"`
		} `stm:"shards"`
	}
)
