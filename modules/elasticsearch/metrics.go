package elasticsearch

type esMetrics struct {
	ClusterHealth *esClusterHealth `stm:"cluster_health"`
	ClusterStats  *esClusterStats  `stm:"cluster_stats"`
}

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

type esClusterStats struct {
	Nodes struct {
		Count struct {
			Data             int `stm:"data"`
			Master           int `stm:"master"`
			Total            int `stm:"total"`
			CoordinatingOnly int `stm:"coordinating_only" json:"coordinating_only"`
			Ingest           int `stm:"ingest"`
		} `stm:"count"`
	} `stm:"nodes"`
	Indices struct {
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
	} `stm:"indices"`
}
