// SPDX-License-Identifier: GPL-3.0-or-later

package mongo

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type connector interface {
	serverStatus() (*serverStatus, error)
	listDatabaseNames() ([]string, error)
	dbStats(name string) (*dbStats, error)
	isReplicaSet() bool
	isMongos() bool
	replSetGetStatus() (*replSetStatus, error)
	shardNodes() (*shardNodesResult, error)
	shardDatabasesPartitioning() (*partitionedResult, error)
	shardCollectionsPartitioning() (*partitionedResult, error)
	shardChunks() (map[string]int64, error)
	initClient(uri string, timeout time.Duration) error
	close() error
}

// mongoCollector interface that helps to abstract and mock the database layer.
type mongoCollector struct {
	client         *mongo.Client
	timeout        time.Duration
	replicaSetFlag *bool
	mongosFlag     *bool
	aggrFunc       func(ctx context.Context, client *mongo.Client, collection string, aggr []bson.D) ([]aggrResults, error)
}

// serverStatus connects to the database and return the output of the
// `serverStatus` command.
func (m *mongoCollector) serverStatus() (*serverStatus, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*m.timeout)
	defer cancel()

	cmd := bson.D{{Key: "serverStatus", Value: 1}, {Key: "metrics", Value: 0}, {Key: "repl", Value: 1}}
	var status *serverStatus

	err := m.client.Database("admin").RunCommand(ctx, cmd).Decode(&status)
	if err != nil {
		return nil, err
	}

	isReplSet := status.Repl != nil
	m.replicaSetFlag = &isReplSet

	isMongos := status.Process == mongos
	m.mongosFlag = &isMongos

	return status, err
}

// listDatabaseNames returns a string slice with the available databases on the server.
func (m *mongoCollector) listDatabaseNames() ([]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*m.timeout)
	defer cancel()

	return m.client.ListDatabaseNames(ctx, bson.M{})
}

// dbStats gets the `dbstats` metrics for a specific database.
func (m *mongoCollector) dbStats(name string) (*dbStats, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*m.timeout)
	defer cancel()

	cmd := bson.M{"dbStats": 1}
	var stats dbStats

	if err := m.client.Database(name).RunCommand(ctx, cmd).Decode(&stats); err != nil {
		return nil, err
	}

	return &stats, nil
}

func (m *mongoCollector) isReplicaSet() bool {
	if m.replicaSetFlag != nil {
		return *m.replicaSetFlag
	}

	status, err := m.serverStatus()
	if err != nil {
		return false
	}

	return status.Repl != nil
}

// isMongos checks if the queried node is a mongos or mongod process
func (m *mongoCollector) isMongos() bool {
	if m.mongosFlag != nil {
		return *m.mongosFlag
	}

	status, err := m.serverStatus()
	if err != nil {
		return false
	}

	return status.Process == mongos
}

// replSetGetStatus gets the `replSetGetStatus` from the server
func (m *mongoCollector) replSetGetStatus() (*replSetStatus, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*m.timeout)
	defer cancel()

	var status *replSetStatus
	cmd := bson.M{"replSetGetStatus": 1}

	err := m.client.Database("admin").RunCommand(ctx, cmd).Decode(&status)
	if err != nil {
		return nil, err
	}

	return status, err
}

func (m *mongoCollector) shardNodes() (*shardNodesResult, error) {
	collection := "shards"
	groupStage := bson.D{{Key: "$sortByCount", Value: "$state"}}

	nodesByState, err := m.shardCollectAggregation(collection, []bson.D{groupStage})
	if err != nil {
		return nil, err
	}

	return &shardNodesResult{nodesByState.True, nodesByState.False}, nil
}

func (m *mongoCollector) shardDatabasesPartitioning() (*partitionedResult, error) {
	collection := "databases"
	groupStage := bson.D{{Key: "$sortByCount", Value: "$partitioned"}}

	partitioning, err := m.shardCollectAggregation(collection, []bson.D{groupStage})
	if err != nil {
		return nil, err
	}

	return &partitionedResult{partitioning.True, partitioning.False}, nil
}

func (m *mongoCollector) shardCollectionsPartitioning() (*partitionedResult, error) {
	collection := "collections"
	matchStage := bson.D{{Key: "$match", Value: bson.D{{Key: "dropped", Value: false}}}}
	countStage := bson.D{{Key: "$sortByCount", Value: bson.D{{Key: "$eq", Value: bson.A{"$distributionMode", "sharded"}}}}}

	partitioning, err := m.shardCollectAggregation(collection, []bson.D{matchStage, countStage})
	if err != nil {
		return nil, err
	}

	return &partitionedResult{partitioning.True, partitioning.False}, nil
}

func (m *mongoCollector) shardCollectAggregation(collection string, aggregation []bson.D) (*aggrResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*m.timeout)
	defer cancel()

	rows, err := m.aggrFunc(ctx, m.client, collection, aggregation)
	if err != nil {
		return nil, err
	}

	result := &aggrResult{}
	for _, row := range rows {
		if row.Bool {
			result.True = row.Count
		} else {
			result.False = row.Count
		}
	}
	return result, err
}

// dbAggregate is not a method in order to mock it out in the tests
func dbAggregate(ctx context.Context, client *mongo.Client, collection string, aggregation []bson.D) ([]aggrResults, error) {
	col := client.Database("config").Collection(collection)

	cursor, err := col.Aggregate(ctx, aggregation)
	if err != nil {
		return nil, err
	}

	defer func() { _ = cursor.Close(ctx) }()

	var rows []aggrResults
	if err = cursor.All(ctx, &rows); err != nil {
		return nil, err
	}

	return rows, nil
}

func (m *mongoCollector) shardChunks() (map[string]int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*m.timeout)
	defer cancel()

	col := m.client.Database("config").Collection("chunks")

	cursor, err := col.Aggregate(ctx, mongo.Pipeline{bson.D{{Key: "$sortByCount", Value: "$shard"}}})
	if err != nil {
		return nil, err
	}

	var rows []bson.M
	if err = cursor.All(ctx, &rows); err != nil {
		return nil, err
	}

	defer func() { _ = cursor.Close(ctx) }()

	result := map[string]int64{}
	for _, row := range rows {
		k, ok := row["_id"].(string)
		if !ok {
			return nil, fmt.Errorf("shard name is not a string: %v", row["_id"])
		}
		v, ok := row["count"].(int32)
		if !ok {
			return nil, fmt.Errorf("shard chunk count is not a int32: %v", row["count"])
		}
		result[k] = int64(v)
	}

	return result, err
}

// initClient initialises the database client if is not initialised.
func (m *mongoCollector) initClient(uri string, timeout time.Duration) error {
	if m.client != nil {
		return nil
	}

	m.timeout = timeout

	client, err := mongo.NewClient(options.Client().ApplyURI(uri))
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), m.timeout*time.Second)
	defer cancel()

	err = client.Connect(ctx)
	if err != nil {
		return err
	}

	m.client = client
	m.aggrFunc = dbAggregate

	return nil
}

// close the database client and all its background goroutines.
func (m *mongoCollector) close() error {
	if m.client == nil {
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), m.timeout*time.Second)
	defer cancel()

	if err := m.client.Disconnect(ctx); err != nil {
		return err
	}

	m.client = nil

	return nil
}
