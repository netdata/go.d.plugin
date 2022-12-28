// SPDX-License-Identifier: GPL-3.0-or-later

package mongo

import "fmt"

func (m *Mongo) collect() (map[string]int64, error) {
	if err := m.mongoCollector.initClient(m.URI, m.Timeout); err != nil {
		return nil, fmt.Errorf("init mongo client: %v", err)
	}

	mx := map[string]int64{}

	if err := m.collectServerStatus(mx); err != nil {
		return nil, fmt.Errorf("couldn't collect server status metrics: %v", err)
	}

	if err := m.collectDbStats(mx); err != nil {
		return mx, fmt.Errorf("couldn't collect dbstats metrics: %v", err)
	}

	if m.mongoCollector.isReplicaSet() {
		if err := m.collectReplSetStatus(mx); err != nil {
			return mx, fmt.Errorf("couldn't collect replSetStatus metrics: %v", err)
		}
	}

	if m.mongoCollector.isMongos() {
		// if we are on a shard based on the serverStatus response
		// we add once the charts during runtime
		m.addShardChartsOnce.Do(func() {
			if err := m.charts.Add(*shardCharts.Copy()...); err != nil {
				m.Errorf("failed to add shard chart: %v", err)
			}
		})

		if err := m.collectShard(mx); err != nil {
			return mx, fmt.Errorf("couldn't collect shard metrics: %v", err)
		}
	}

	return mx, nil
}

func boolToInt(v bool) int64 {
	if v {
		return 1
	}
	return 0
}
