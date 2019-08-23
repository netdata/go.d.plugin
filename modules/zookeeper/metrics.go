package zookeeper

// zk_version	3.5.5-390fe37ea45dee01bf87dc1c042b5e3dcce88653, built on 05/03/2019 12:07 GMT
//zk_avg_latency	0
//zk_max_latency	0
//zk_min_latency	0
//zk_packets_received	140
//zk_packets_sent	139
//zk_num_alive_connections	1
//zk_outstanding_requests	0
//zk_server_state	standalone
//zk_znode_count	5
//zk_watch_count	0
//zk_ephemerals_count	0
//zk_approximate_data_size	44
//zk_open_file_descriptor_count	46
//zk_max_file_descriptor_count	1048576
type mntr struct {
	AvgLatency          int64 `stm:"avg_latency"`
	MaxLatency          int64 `stm:"max_latency"`
	MinLatency          int64 `stm:"min_latency"`
	PacketsReceived     int64 `stm:"packets_received"`
	PacketsSent         int64 `stm:"packets_sent"`
	NumAliveConnections int64 `stm:"num_alive_connections"`
	OutstandingRequests int64 `stm:"outstanding_requests"`
}
