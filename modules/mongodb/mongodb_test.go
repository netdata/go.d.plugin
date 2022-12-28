// SPDX-License-Identifier: GPL-3.0-or-later

package mongo

import (
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/netdata/go.d.plugin/agent/module"
	"github.com/netdata/go.d.plugin/modules/mongodb/testdata/v5.0.0"
	"github.com/netdata/go.d.plugin/pkg/matcher"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/mongo"
)

func TestMongo_Init(t *testing.T) {
	tests := map[string]struct {
		config  Config
		success bool
	}{
		"success on default config": {
			success: true,
			config:  New().Config,
		},
		"fails on unset 'address'": {
			success: true,
			config: Config{
				URI:     "mongodb://localhost:27017",
				Timeout: 10,
			},
		},
		"fails on invalid port": {
			success: false,
			config: Config{
				URI:     "",
				Timeout: 0,
			},
		},
	}

	msg := "Init() result does not match Init()"

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			m := New()
			m.Config = test.config
			assert.Equal(t, test.success, m.Init(), msg)
		})
	}
}

func TestMongo_Init_AddServerChartsTwiceFails(t *testing.T) {
	msg := "adding duplicate server status charts is expected to fail Init()"
	m := New()
	// duplicate charts
	temp := serverStatusCharts
	defer func() { serverStatusCharts = temp }()
	serverStatusCharts = append(serverStatusCharts, serverStatusCharts...)
	assert.Equal(t, false, m.Init(), msg)
}

func TestMongo_Init_AddDbChartsTwiceFails(t *testing.T) {
	msg := "adding duplicate db stats charts is expected to fail Init()"
	m := New()
	// duplicate charts
	temp := dbStatsChartsTmpl
	defer func() { dbStatsChartsTmpl = temp }()
	dbStatsChartsTmpl = append(dbStatsChartsTmpl, dbStatsChartsTmpl...)
	assert.Equal(t, false, m.Init(), msg)
}

func TestMongo_Init_BadMatcher(t *testing.T) {
	msg := "bad database matcher value is expected to fail Init()"
	m := New()
	m.Databases = matcher.SimpleExpr{
		Includes: []string{"bad value"},
		Excludes: nil,
	}
	assert.Equal(t, false, m.Init(), msg)
}

func TestMongo_Charts(t *testing.T) {
	msg := "after Init() we expect to have server status and db stats charts"
	m := New()
	require.True(t, m.Init())
	assert.Len(t, *m.Charts(), 14, msg)
}

func TestMongo_ChartsOptional(t *testing.T) {
	// optional charts are the serer status charts
	// depending on the supported metrics by
	// the database version and configuration
	msg := "optional charts should be added after the first Collect()"
	m := New()
	require.True(t, m.Init())
	charts := *m.Charts()
	var IDs []string
	for _, chart := range charts {
		IDs = append(IDs, chart.ID)
	}
	require.NotNil(t, m.Charts())
	for _, id := range []string{
		"current_transactions",
		"flow_control_timings",
		"active_clients",
		"queued_operations",
		"tcmalloc",
		"tcmalloc_generic",
		"wiredtiger_cache",
		"wiredtiger_capacity",
		"wiredtiger_connection",
		"wiredtiger_cursor",
		"wiredtiger_lock",
		"wiredtiger_lock_duration",
		"wiredtiger_log_ops",
		"wiredtiger_log_ops_size",
		"wiredtiger_transactions",
	} {
		assert.NotContainsf(t, IDs, id, msg)
	}
}

func TestMongo_initMongoClient_uri(t *testing.T) {
	m := New()
	m.Config.URI = "mongodb://user:pass@localhost:27017"
	assert.True(t, m.Init())
}

func TestMongo_CheckFail(t *testing.T) {
	m := New()
	m.Config.Timeout = 0
	assert.False(t, m.Check(), "Check() should fail with context deadline exceeded")
}

func TestMongo_Success(t *testing.T) {
	m := New()
	m.Config.Timeout = 1
	m.Config.URI = ""
	obj := &mockMongo{serverStatusResponse: v5_0_0.ServerStatus}
	m.mongoCollector = obj
	assert.True(t, m.Check(), "check should success with the mocker serverStatus response")
}

func TestMongo_Collect_DbStats(t *testing.T) {
	m := New()
	m.mongoCollector = &mockMongo{
		serverStatusResponse:      "{}",
		listDatabaseNamesResponse: v5_0_0.ListDatabaseNames,
		dbStatsResponse:           v5_0_0.DbStats,
		replicaSetResponse:        "{}",
		closeCalled:               false,
		replicaSet:                false,
	}
	m.Config.Databases.Includes = []string{"* *"} // matcher
	m.URI = "mongodb://localhost"
	require.True(t, m.Init())
	ms := m.Collect()
	msg := "collected values should be equal to the number of the dbStats field * number of databases"
	assert.Len(t, ms, reflect.ValueOf(dbStats{}).NumField()*len(v5_0_0.ListDatabaseNames), msg)

	charts := []*module.Chart{
		chartDBStatsCollectionsTmpl.Copy(),
		chartDBStatsIndexesTmpl.Copy(),
		chartDBStatsViewsTmpl.Copy(),
		chartDBStatsDocumentsTmpl.Copy(),
		chartDBStatsSizeTmpl.Copy(),
	}
	for _, chart := range charts {
		require.True(t, m.charts.Has(chart.ID))
		for _, dbName := range v5_0_0.ListDatabaseNames {
			dimID := fmt.Sprintf("%s_%s", chart.ID, dbName)
			assert.True(t, m.charts.Get(chart.ID).HasDim(dimID), "dimension is expected")
			assert.EqualValues(t, 1, ms[dimID], "all values are hardcode to 1 in the test data")
		}
	}
}

func TestMongo_Collect_DbStatsRemoveDropped(t *testing.T) {
	m := New()
	m.mongoCollector = &mockMongo{
		serverStatusResponse:      "{}",
		dbStatsResponse:           "{}",
		listDatabaseNamesResponse: []string{"db1", "db2"},
	}
	m.Config.Databases.Includes = []string{"* *"}
	m.URI = "mongodb://localhost"
	require.True(t, m.Init())
	ms := m.Collect()
	assert.Len(t, ms, 14)

	// remove a database
	m.mongoCollector = &mockMongo{
		serverStatusResponse:      "{}",
		dbStatsResponse:           "{}",
		listDatabaseNamesResponse: []string{"db1"},
	}
	ms = m.Collect()
	msg := "dimension was removed but is still active"
	assert.True(t, m.charts.Get("database_collections").Dims[1].Obsolete, msg)
	assert.Len(t, ms, 7, "we should have collected exactly 7 metrics")

	// add two databases
	m.mongoCollector = &mockMongo{
		serverStatusResponse:      "{}",
		dbStatsResponse:           "{}",
		listDatabaseNamesResponse: []string{"db1", "db2", "db3"},
	}
	ms = m.Collect()
	msg = "after adding two databases we should still have 3 charts"
	assert.Len(t, m.charts.Get("database_collections").Dims, 3, msg)
	msg = "after adding two databases we should still have 3 charts with 7 dimensions each"
	assert.Len(t, ms, 21, msg)
}

func TestMongo_Collect_DbStats_Fail(t *testing.T) {
	m := New()
	m.mongoCollector = &mockMongo{
		serverStatusResponse:      "{}",
		dbStatsResponse:           "",
		listDatabaseNamesResponse: []string{},
		replicaSetResponse:        "{}",
	}
	m.Config.Databases.Includes = []string{"* *"}
	m.URI = "mongodb://localhost"
	require.True(t, m.Init())
	ms := m.Collect()
	assert.Len(t, ms, 0)
}

func TestMongo_Collect_DbStats_EmptyMatcher(t *testing.T) {
	m := New()
	m.mongoCollector = &mockMongo{
		serverStatusResponse:      "{}",
		dbStatsResponse:           "{}",
		listDatabaseNamesResponse: []string{},
		replicaSetResponse:        "{}",
	}
	m.Config.Databases.Includes = []string{"* not_matching"}
	m.URI = "mongodb://localhost"
	require.True(t, m.Init())
	ms := m.Collect()
	msg := "we shouldn't have any metrics with a bad matcher"
	assert.Len(t, ms, 0, msg)
}

func TestMongo_Collect_ReplSetStatus(t *testing.T) {
	m := New()
	m.mongoCollector = &mockMongo{
		serverStatusResponse:      "{}",
		dbStatsResponse:           "{}",
		listDatabaseNamesResponse: []string{},
		replicaSet:                true,
		replicaSetResponse:        v5_0_0.ReplSetGetStatus,
	}
	m.Config.Databases.Includes = []string{"* *"}
	m.URI = "mongodb://localhost"
	require.True(t, m.Init())
	_ = m.Collect()
	msg := "%s chart should have been added"
	assert.True(t, m.charts.Has(replicationLag), msg, replicationLag)
	assert.True(t, m.charts.Has(replicationHeartbeatLatency), msg, replicationHeartbeatLatency)
	assert.True(t, m.charts.Has(replicationNodePing), msg, replicationNodePing)
}

func TestMongo_Collect_ReplSetStatusAddRemove(t *testing.T) {
	m := New()
	m.mongoCollector = &mockMongo{
		serverStatusResponse:      "{}",
		dbStatsResponse:           "{}",
		listDatabaseNamesResponse: []string{},
		replicaSet:                true,
		replicaSetResponse:        v5_0_0.ReplSetGetStatusNode1,
	}
	m.Config.Databases.Includes = []string{"* *"}
	m.URI = "mongodb://localhost"
	require.True(t, m.Init())
	_ = m.Collect()
	msg := "node1 dimension is missing"
	assert.True(t, m.charts.Get(replicationLag).HasDim(replicationLagDimPrefix+"node1"), msg)
	assert.True(t, m.charts.Get(replicationHeartbeatLatency).HasDim(replicationHeartbeatLatencyDimPrefix+"node1"), msg)
	assert.True(t, m.charts.Get(replicationNodePing).HasDim(replicationNodePingDimPrefix+"node1"), msg)

	m.mongoCollector = &mockMongo{
		serverStatusResponse:      "{}",
		dbStatsResponse:           "{}",
		listDatabaseNamesResponse: []string{},
		replicaSet:                true,
		replicaSetResponse:        v5_0_0.ReplSetGetStatusNode2,
	}
	_ = m.Collect()
	// node2 dimensions added
	msg = "node2 dimension is missing"
	assert.True(t, m.charts.Get(replicationLag).HasDim(replicationLagDimPrefix+"node2"), msg)
	assert.True(t, m.charts.Get(replicationHeartbeatLatency).HasDim(replicationHeartbeatLatencyDimPrefix+"node2"), msg)
	assert.True(t, m.charts.Get(replicationNodePing).HasDim(replicationNodePingDimPrefix+"node2"), msg)

	// node1 dimensions removed
	msg = "node1 was remove but dimension is still active"
	assert.True(t, m.charts.Get(replicationLag).GetDim(replicationLagDimPrefix+"node1").Obsolete, msg)
	assert.True(t, m.charts.Get(replicationHeartbeatLatency).GetDim(replicationHeartbeatLatencyDimPrefix+"node1").Obsolete, msg)
	assert.True(t, m.charts.Get(replicationNodePing).GetDim(replicationNodePingDimPrefix+"node1").Obsolete, msg)
}

func TestMongo_Collect_Shard(t *testing.T) {
	m := New()
	mockClient := &mockMongo{
		serverStatusResponse:      "{}",
		listDatabaseNamesResponse: []string{},
		dbStatsResponse:           "{}",
		replicaSetResponse:        "{}",
		mongos:                    true,
		shardNodesResponse:        v5_0_0.ShardNodes,
		shardDbPartitionResponse:  v5_0_0.ShardDatabases,
		shardColPartitionResponse: v5_0_0.ShardCollections,
		chunksShardNum:            2,
	}
	mockClient.connector = &mongoCollector{aggrFunc: mockClient.dbAggregate}
	m.mongoCollector = mockClient
	m.Config.Databases.Includes = []string{"* *"}
	m.URI = "mongodb://localhost"
	require.True(t, m.Init())
	ms := m.Collect()
	msg := "%s chart should have been added"
	for _, chart := range shardCharts {
		assert.True(t, m.charts.Has(chart.ID), msg, chart.ID)
		assert.Len(t, m.charts.Get(chart.ID).Dims, 2)
	}
	assert.Len(t, ms, 8)
}

func TestMongo_Collect_Shard_Fail(t *testing.T) {
	m := New()
	mockClient := &mockMongoErrors{
		mockMongo: mockMongo{
			serverStatusResponse:      "{}",
			listDatabaseNamesResponse: []string{},
			dbStatsResponse:           "{}",
			replicaSetResponse:        "{}",
			mongos:                    true,
			shardNodesResponse:        v5_0_0.ShardNodes,
			shardDbPartitionResponse:  v5_0_0.ShardDatabases,
			shardColPartitionResponse: v5_0_0.ShardCollections,
			chunksShardNum:            2,
		},
	}
	mockClient.connector = &mongoCollector{
		client:   &mongo.Client{},
		aggrFunc: mockClient.dbAggregate,
	}
	m.mongoCollector = mockClient
	m.Config.Databases.Includes = []string{"* *"}
	m.URI = "mongodb://localhost"
	require.True(t, m.Init())

	ms := m.Collect()
	assert.Len(t, ms, 8)

	mockClient.shardChunksError = true
	ms = m.Collect()
	assert.Len(t, ms, 6)

	mockClient.shardCollectionsPartitioningError = true
	ms = m.Collect()
	assert.Len(t, ms, 4)

	mockClient.shardDatabasesPartitioningError = true
	ms = m.Collect()
	assert.Len(t, ms, 2)

	mockClient.shardNodesError = true
	ms = m.Collect()
	assert.Len(t, ms, 0)

}

func TestMongo_ShardUpdateNodeChart(t *testing.T) {
	m := New()
	mockClient := &mockMongo{
		serverStatusResponse:      "{}",
		listDatabaseNamesResponse: []string{},
		dbStatsResponse:           "{}",
		replicaSetResponse:        "{}",
		mongos:                    true,
		shardNodesResponse:        v5_0_0.ShardNodes,
		shardDbPartitionResponse:  v5_0_0.ShardDatabases,
		shardColPartitionResponse: v5_0_0.ShardCollections,
		chunksShardNum:            2,
	}
	mockClient.connector = &mongoCollector{aggrFunc: mockClient.dbAggregate}
	m.mongoCollector = mockClient
	m.Config.Databases.Includes = []string{"* *"}
	m.URI = "mongodb://localhost"
	require.True(t, m.Init())
	_ = m.Collect()
	assert.Len(t, m.charts.Get("shard_chucks_per_node").Dims, 2)

	mockClient.chunksShardNum = 1
	_ = m.Collect()
	chart := m.charts.Get("shard_chucks_per_node")
	assert.True(t, chart.GetDim("shard_chucks_per_node_shard2").Obsolete)
}

func TestMongo_Incomplete(t *testing.T) {
	m := New()
	m.mongoCollector = &mockMongo{}
	ms := m.Collect()
	msg := "uninitialized client should collect any data"
	assert.Len(t, ms, 0, msg)
}

func TestMongo_Cleanup(t *testing.T) {
	m := New()
	connector := &mockMongo{}
	m.mongoCollector = connector
	m.Cleanup()
	msg := "Cleanup() should have closed the mongo client"
	assert.True(t, connector.closeCalled, msg)
}

func TestCollectUpToServerStatus(t *testing.T) {
	m := New()
	m.Timeout = 0
	m.mongoCollector = &mockMongoServerStatusOnly{
		serverStatusResponse: v5_0_0.ServerStatus,
	}

	m.Config.Databases.Includes = []string{"* *"}
	m.URI = "mongodb://localhost"
	require.True(t, m.Init())
	ms := m.Collect()
	msg := "dim should not have been added: %s"
	for dim := range ms {
		assert.False(t, strings.HasPrefix(dim, "database"), msg, dim)
		assert.False(t, strings.HasPrefix(dim, "replication"), msg, dim)
	}
}

func TestCollectUpToServerStatusListDbNamesFails(t *testing.T) {
	m := New()
	m.Timeout = 0
	m.mongoCollector = &mockMongoServerStatusOnly{
		serverStatusResponse:  v5_0_0.ServerStatus,
		mockListDatabaseNames: true,
	}
	m.Config.Databases.Includes = []string{"* *"}
	m.URI = "mongodb://localhost"
	require.True(t, m.Init())
	ms := m.Collect()
	msg := "dim should not have been added: %s"
	for dim := range ms {
		assert.False(t, strings.HasPrefix(dim, "database"), msg, dim)
		assert.False(t, strings.HasPrefix(dim, "replication"), msg, dim)
	}
}

func TestCollectUpToDbStats(t *testing.T) {
	m := New()
	m.Timeout = 0
	m.mongoCollector = &mockMongoServerStatusOnly{
		serverStatusResponse:  v5_0_0.ServerStatus,
		mockListDatabaseNames: true,
		mockDbStats:           true,
	}
	m.Config.Databases.Includes = []string{"* *"}
	m.URI = "mongodb://localhost"
	require.True(t, m.Init())
	ms := m.Collect()

	foundDbStatsDims := false
	for dim := range ms {
		if strings.HasPrefix(dim, "database") {
			foundDbStatsDims = true
			break
		}
	}
	assert.True(t, foundDbStatsDims, "some dims 'database_*' were not found")
}
