// SPDX-License-Identifier: GPL-3.0-or-later

package mongo

import (
	"strings"

	"github.com/netdata/go.d.plugin/agent/module"
)

// collectShard adds the sharding stats for the available dims.
func (m *Mongo) collectShard(ms map[string]int64) error {
	// nodes in the shard cluster
	nodes, err := m.mongoCollector.shardNodes()
	if err != nil {
		return err
	}
	ms["shard_nodes_count_aware"] = nodes.ShardAware
	ms["shard_nodes_count_unaware"] = nodes.ShardUnaware

	// databases partitioning
	databasesPartitioning, err := m.mongoCollector.shardDatabasesPartitioning()
	if err != nil {
		return err
	}
	ms["shard_databases_partitioned"] = databasesPartitioning.Partitioned
	ms["shard_databases_unpartitioned"] = databasesPartitioning.UnPartitioned

	// collections partitioning
	collectionsPartitioning, err := m.mongoCollector.shardCollectionsPartitioning()
	if err != nil {
		return err
	}
	ms["shard_collections_partitioned"] = collectionsPartitioning.Partitioned
	ms["shard_collections_unpartitioned"] = collectionsPartitioning.UnPartitioned

	// chunks per shard node
	chunksPerShard, err := m.mongoCollector.shardChunks()
	if err != nil {
		return err
	}
	m.updateShardChunkChartDims(chunksPerShard)
	for shard, count := range chunksPerShard {
		ms["shard_chucks_per_node_"+shard] = count
	}

	return nil
}

func (m *Mongo) updateShardChunkChartDims(chunksPerShard map[string]int64) {
	chart := m.charts.Get("shard_chucks_per_node")
	if chart == nil {
		return
	}
	for _, dim := range chart.Dims {
		if _, ok := chunksPerShard[strings.TrimPrefix(dim.ID, chart.ID+"_")]; !ok {
			delete(m.shardNodesDims, dim.ID)
			if err := chart.MarkDimRemove(dim.ID, true); err != nil {
				m.Warningf("updateShardChunkChartDims failed to remove dim %v", err)
			} else {
				chart.MarkNotCreated()
			}
		}
	}
	for shard := range chunksPerShard {
		id := chart.ID + "_" + shard
		if !m.shardNodesDims[id] {
			m.shardNodesDims[id] = true
			err := chart.AddDim(&module.Dim{ID: id, Name: shard, Algo: module.Absolute})
			if err != nil {
				m.Warningf("updateShardChunkChartDims failed to add dim %v", err)
				continue
			}
			chart.MarkNotCreated()
		}
	}
}
