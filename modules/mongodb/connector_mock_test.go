// SPDX-License-Identifier: GPL-3.0-or-later

package mongo

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/netdata/go.d.plugin/modules/mongodb/testdata/v5.0.0"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// mock for unit testing the mongo database as the driver
// doesn't use interfaces.
type mockMongo struct {
	connector
	serverStatusResponse      string
	listDatabaseNamesResponse []string
	dbStatsResponse           string
	closeCalled               bool
	replicaSet                bool
	replicaSetResponse        string
	shardNodesResponse        string
	shardDbPartitionResponse  string
	shardColPartitionResponse string
	chunksShardNum            int
	mongos                    bool
}

func (m *mockMongo) initClient(_ string, _ time.Duration) error {
	return nil
}

func (m *mockMongo) serverStatus() (*serverStatus, error) {
	var status serverStatus
	err := json.Unmarshal([]byte(m.serverStatusResponse), &status)
	if err != nil {
		return nil, err
	}
	return &status, nil
}

func (m *mockMongo) close() error {
	m.closeCalled = true
	return nil
}

func (m *mockMongo) listDatabaseNames() ([]string, error) {
	if m.listDatabaseNamesResponse == nil {
		return nil, errors.New("mocked error")
	}
	return m.listDatabaseNamesResponse, nil
}

func (m *mockMongo) dbStats(_ string) (*dbStats, error) {
	stats := &dbStats{}
	err := json.Unmarshal([]byte(m.dbStatsResponse), stats)
	if err != nil {
		return nil, err
	}
	return stats, nil
}

func (m *mockMongo) isReplicaSet() bool {
	return m.replicaSet
}

func (m *mockMongo) replSetGetStatus() (*replSetStatus, error) {
	status := &replSetStatus{}
	err := json.Unmarshal([]byte(m.replicaSetResponse), status)
	if err != nil {
		return nil, err
	}
	return status, nil
}

func (m *mockMongo) isMongos() bool {
	return m.mongos
}

func (m *mockMongo) shardNodes() (*shardNodesResult, error) {
	return m.connector.shardNodes()
}

func (m *mockMongo) shardDatabasesPartitioning() (*partitionedResult, error) {
	return m.connector.shardDatabasesPartitioning()
}

func (m *mockMongo) shardCollectionsPartitioning() (*partitionedResult, error) {
	return m.connector.shardCollectionsPartitioning()
}

func (m *mockMongo) shardChunks() (map[string]int64, error) {
	res := map[string]int64{}
	for i := 1; i <= m.chunksShardNum; i++ {
		res[fmt.Sprintf("shard%d", i)] = int64(i)
	}
	return res, nil
}

func (m *mockMongo) dbAggregate(_ context.Context, _ *mongo.Client, collection string, _ []bson.D) ([]aggrResults, error) {
	var res []aggrResults
	var response string
	switch collection {
	case "shards":
		response = v5_0_0.ShardNodes
	case "databases":
		response = v5_0_0.ShardDatabases
	case "collections":
		response = v5_0_0.ShardCollections
	}
	err := json.Unmarshal([]byte(response), &res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

type mockMongoServerStatusOnly struct {
	mongoCollector
	serverStatusResponse  string
	mockListDatabaseNames bool
	mockDbStats           bool
}

func (m *mockMongoServerStatusOnly) listDatabaseNames() ([]string, error) {
	if !m.mockListDatabaseNames {
		return m.mongoCollector.listDatabaseNames()
	}
	return []string{"db1", "db2"}, nil
}

func (m *mockMongoServerStatusOnly) serverStatus() (*serverStatus, error) {
	var status serverStatus
	err := json.Unmarshal([]byte(m.serverStatusResponse), &status)
	if err != nil {
		return nil, err
	}
	return &status, nil
}

func (m *mockMongoServerStatusOnly) dbStats(_ string) (*dbStats, error) {
	if !m.mockDbStats {
		return m.mongoCollector.dbStats("db")
	}
	stats := &dbStats{}
	err := json.Unmarshal([]byte(v5_0_0.DbStats), stats)
	if err != nil {
		return nil, err
	}
	return stats, nil
}

func (m *mockMongoServerStatusOnly) isReplicaSet() bool {
	return true
}

type mockMongoErrors struct {
	mockMongo
	shardNodesError                   bool
	shardDatabasesPartitioningError   bool
	shardCollectionsPartitioningError bool
	shardChunksError                  bool
}

func (m *mockMongoErrors) shardNodes() (*shardNodesResult, error) {
	if m.shardNodesError {
		return nil, errors.New("test error")
	}
	return m.mockMongo.shardNodes()
}

func (m *mockMongoErrors) shardDatabasesPartitioning() (*partitionedResult, error) {
	if m.shardDatabasesPartitioningError {
		return nil, errors.New("test error")
	}
	return m.mockMongo.shardDatabasesPartitioning()
}

func (m *mockMongoErrors) shardCollectionsPartitioning() (*partitionedResult, error) {
	if m.shardCollectionsPartitioningError {
		return nil, errors.New("test error")
	}
	return m.mockMongo.shardCollectionsPartitioning()
}

func (m *mockMongoErrors) shardChunks() (map[string]int64, error) {
	if m.shardChunksError {
		return m.mockMongo.connector.shardChunks()
	}
	return m.mockMongo.shardChunks()
}
