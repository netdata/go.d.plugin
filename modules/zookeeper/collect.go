package zookeeper

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

func (z *Zookeeper) collect() (map[string]int64, error) {
	mx, err := z.collectMntr()
	if err != nil {
		return nil, fmt.Errorf("error on collecting 'mntr' : %v", err)
	}
	return mx, nil
}

func (z *Zookeeper) collectMntr() (map[string]int64, error) {
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
	rows, err := z.fetch("mntr")
	if err != nil {
		return nil, err
	}
	if len(rows) == 0 {
		return nil, errors.New("empty response")
	}
	if len(rows) == 1 {
		// mntr is not executed because it is not in the whitelist.
		return nil, fmt.Errorf("bad response : %s", rows[0])
	}
	mx := make(map[string]int64)

	for _, row := range rows {
		parseMntrRow(row, mx)
	}

	if len(mx) == 0 {
		return nil, fmt.Errorf("failed to parse reponse")
	}

	return mx, nil
}

func parseMntrRow(row string, mx map[string]int64) {
	parts := strings.Fields(row)
	if len(parts) != 2 {
		return
	}

	key := strings.TrimPrefix(parts[0], "zk_")
	value := parts[1]
	switch key {
	case "version":
	case "server_state":
		mx[key] = parseServerState(value)
	default:
		v, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return
		}
		mx[key] = v
	}
}

func parseServerState(state string) int64 {
	switch state {
	default:
		return 0
	case "leader":
		return 1
	case "follower":
		return 2
	case "observer":
		return 3
	case "standalone":
		return 4
	}
}
