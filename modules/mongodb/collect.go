// SPDX-License-Identifier: GPL-3.0-or-later

package mongo

import "fmt"

func (m *Mongo) collect() (map[string]int64, error) {
	if err := m.mongoCollector.initClient(m.URI, m.Timeout); err != nil {
		return nil, fmt.Errorf("init mongo client: %v", err)
	}

	mx := map[string]int64{}
	if err := m.collectServerStatus(mx); err != nil {
		return nil, fmt.Errorf("couldn't collecting server status metrics: %v", err)
	}

	if err := m.collectDbStats(mx); err != nil {
		return mx, fmt.Errorf("couldn't collecting dbstats metrics: %v", err)
	}

	if m.mongoCollector.isReplicaSet() {
		// if we have replica set based on the serverStatus response
		// we add once the charts during runtime
		m.addReplChartsOnce.Do(func() {
			if err := m.charts.Add(*replCharts.Copy()...); err != nil {
				m.Errorf("failed to add replica set chart: %v", err)
			}
		})

		if err := m.collectReplSetStatus(mx); err != nil {
			return mx, fmt.Errorf("couldn't collecting replSetStatus metrics: %v", err)
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
			return mx, fmt.Errorf("couldn't collecting shard metrics: %v", err)
		}
	}

	return mx, nil
}

// sliceDiff calculates the diff between to slices
func sliceDiff(slice1, slice2 []string) []string {
	mb := make(map[string]struct{}, len(slice2))
	for _, x := range slice2 {
		mb[x] = struct{}{}
	}
	var diff []string
	for _, x := range slice1 {
		if _, found := mb[x]; !found {
			diff = append(diff, x)
		}
	}
	return diff
}
