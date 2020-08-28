package elasticsearch

type esMetrics struct {
	ClusterHealth *esClusterHealth `stm:"cluster_health"`
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
