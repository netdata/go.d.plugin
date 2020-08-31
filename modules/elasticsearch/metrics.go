package elasticsearch

// https://www.elastic.co/guide/en/elasticsearch/reference/current/docker.html

type esMetrics struct {
	// https://www.elastic.co/guide/en/elasticsearch/reference/current/cluster-nodes-stats.html
	LocalNodeStats *esNodeStats
	// https://www.elastic.co/guide/en/elasticsearch/reference/current/cluster-health.html
	ClusterHealth *esClusterHealth
	// https://www.elastic.co/guide/en/elasticsearch/reference/current/cluster-stats.html
	ClusterStats *esClusterStats
	// https://www.elastic.co/guide/en/elasticsearch/reference/current/cat-indices.html
	LocalIndicesStats []esIndexStats
}

func (m esMetrics) empty() bool {
	switch {
	case m.hasLocalNodeStats(), m.hasClusterHealth(), m.hasClusterStats(), m.hasLocalIndicesStats():
		return false
	}
	return true
}

func (m esMetrics) hasLocalNodeStats() bool    { return m.LocalNodeStats != nil }
func (m esMetrics) hasClusterHealth() bool     { return m.ClusterHealth != nil }
func (m esMetrics) hasClusterStats() bool      { return m.ClusterStats != nil }
func (m esMetrics) hasLocalIndicesStats() bool { return len(m.LocalIndicesStats) > 0 }

type esNodeStats struct {
	Indices struct {
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
		} `stm:"flush"`
		FieldData struct {
			MemorySizeInBytes int `stm:"memory_size_in_bytes" json:"memory_size_in_bytes"`
			Evictions         int `stm:"evictions"`
		} `stm:"fielddata"`
		Segments struct {
			Count                     int `stm:"count" json:"count"`
			MemoryInBytes             int `stm:"memory_in_bytes" json:"memory_in_bytes"`
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
	} `stm:"indices"`
	Process struct {
		OpenFileDescriptors int `stm:"open_file_descriptors" json:"open_file_descriptors"`
		MaxFileDescriptors  int `stm:"max_file_descriptors" json:"max_file_descriptors"`
	} `stm:"process"`
	JVM struct {
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
		} `stm:"buffer_pool" json:"buffer_pool"`
	} `stm:"jvm"`
	ThreadPool struct {
		Search struct {
			Queue    int `stm:"queue"`
			Rejected int `stm:"rejected"`
		} `stm:"search"`
		Write struct {
			Queue    int `stm:"queue"`
			Rejected int `stm:"rejected"`
		} `stm:"write"`
	} `stm:"thread_pool" json:"thread_pool"`
	Transport struct {
		RxCount       int `stm:"rx_count" json:"rx_count"`
		RxSizeInBytes int `stm:"rx_size_in_bytes" json:"rx_size_in_bytes"`
		TxCount       int `stm:"tx_count" json:"tx_count"`
		TxSizeInBytes int `stm:"tx_size_in_bytes" json:"tx_size_in_bytes"`
	} `stm:"transport"`
	HTTP struct {
		CurrentOpen int `stm:"current_open" json:"current_open"`
	} `stm:"http"`
	Breakers struct {
		Requests struct {
			Tripped int `stm:"tripped"`
		} `stm:"requests"`
		FieldData struct {
			Tripped int `stm:"tripped"`
		} `stm:"fielddata"`
		InFlightRequests struct {
			Tripped int `stm:"tripped"`
		} `stm:"in_flight_requests" json:"in_flight_requests"`
		ModelInference struct {
			Tripped int `stm:"tripped"`
		} `stm:"model_inference" json:"model_inference"`
		Accounting struct {
			Tripped int `stm:"tripped"`
		} `stm:"accounting"`
		Parent struct {
			Tripped int `stm:"tripped"`
		} `stm:"parent"`
	} `stm:"breakers"`
}

type esClusterHealth struct {
	Status                      string
	NumOfNodes                  int `stm:"number_of_nodes" json:"number_of_nodes"`
	NumOfDataNodes              int `stm:"number_of_data_nodes" json:"number_of_data_nodes"`
	ActivePrimaryShards         int `stm:"active_primary_shards" json:"active_primary_shards"`
	ActiveShards                int `stm:"active_shards" json:"active_shards"`
	RelocatingShards            int `stm:"relocating_shards" json:"relocating_shards"`
	InitializingShards          int `stm:"initializing_shards" json:"initializing_shards"`
	UnassignedShards            int `stm:"unassigned_shards" json:"unassigned_shards"`
	DelayedUnassignedShards     int `stm:"delayed_unassigned_shards" json:"delayed_unassigned_shards"`
	NumOfPendingTasks           int `stm:"number_of_pending_tasks" json:"number_of_pending_tasks"`
	NumOfInFlightFetch          int `stm:"number_of_in_flight_fetch" json:"number_of_in_flight_fetch"`
	ActiveShardsPercentAsNumber int `stm:"active_shards_percent_as_number" json:"active_shards_percent_as_number"`
}

/*
   "total": 1,
   "coordinating_only": 0,
   "data": 1,
   "ingest": 1,
   "master": 1,
   "ml": 1,
   "remote_cluster_client": 1,
   "transform": 1,
   "voting_only": 0
*/

type esClusterStats struct {
	Nodes struct {
		Count struct {
			Total               int `stm:"total"`
			CoordinatingOnly    int `stm:"coordinating_only" json:"coordinating_only"`
			Data                int `stm:"data"`
			Ingest              int `stm:"ingest"`
			Master              int `stm:"master"`
			ML                  int `stm:"ml"`
			RemoteClusterClient int `stm:"remote_cluster_client" json:"remote_cluster_client"`
			Transform           int `stm:"transform"`
			VotingOnly          int `stm:"voting_only" json:"voting_only"`
		} `stm:"count"`
	} `stm:"nodes"`
	Indices struct {
		Count  int `stm:"count"`
		Shards struct {
			Total       int `stm:"total"`
			Primaries   int `stm:"primaries"`
			Replication int `stm:"replication"`
		} `stm:"shards"`
		Docs struct {
			Count int `stm:"count"`
		} `stm:"docs"`
		Store struct {
			SizeInBytes int `stm:"size_in_bytes" json:"size_in_bytes"`
		} `stm:"store"`
		QueryCache struct {
			HitCount  int `stm:"hit_count" json:"hit_count"`
			MissCount int `stm:"miss_count" json:"miss_count"`
		} `stm:"query_cache" json:"query_cache"`
	} `stm:"indices"`
}

type esIndexStats struct {
	Index     string
	Health    string
	Rep       string
	DocsCount string `json:"docs.count"`
	StoreSize string `json:"store.size"`
}
