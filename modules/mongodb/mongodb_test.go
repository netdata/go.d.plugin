package mongo

import (
	"encoding/json"
	"reflect"
	"testing"
	"time"

	v5_0_0 "github.com/netdata/go.d.plugin/modules/mongodb/testdata/v5.0.0"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			m := New()
			m.Config = test.config
			assert.Equal(t, test.success, m.Init())
		})
	}
}

func TestMongo_Charts(t *testing.T) {
	m := New()
	require.True(t, m.Init())
	assert.NotNil(t, m.Charts())
}

func TestMongo_ChartsOptional(t *testing.T) {
	m := New()
	require.True(t, m.Init())
	assert.NotNil(t, m.Charts())
}

func TestMongo_initMongoClient_uri(t *testing.T) {
	m := New()
	m.Config.URI = "mongodb://user:pass@localhost:27017"
	assert.True(t, m.Init())
}

func TestMongo_CheckFail(t *testing.T) {
	m := New()
	m.Config.Timeout = 0
	assert.False(t, m.Check())
}

func TestMongo_Success(t *testing.T) {
	m := New()
	m.Config.Timeout = 1
	m.Config.URI = ""
	obj := &mockMongo{serverStatusResponse: v5_0_0.ServerStatus}
	m.mongoCollector = obj
	assert.True(t, m.Check())
}

func TestMongo_Collect_DbStats(t *testing.T) {
	m := New()
	m.mongoCollector = &mockMongo{
		serverStatusResponse: "{}",
		dbStatsResponse:      "{}",
		databaseNames:        []string{"db"},
	}
	m.Config.Databases.Includes = []string{"* *"}
	m.URI = "mongodb://localhost"
	m.Init()
	ms := m.Collect()
	assert.Len(t, ms, reflect.ValueOf(dbStats{}).NumField())
}

func TestMongo_Collect_DbStatsRemoveDropped(t *testing.T) {
	m := New()
	m.mongoCollector = &mockMongo{
		serverStatusResponse: "{}",
		dbStatsResponse:      "{}",
		databaseNames:        []string{"db1", "db2"},
	}
	m.Config.Databases.Includes = []string{"* *"}
	m.URI = "mongodb://localhost"
	m.Init()
	ms := m.Collect()
	assert.Len(t, ms, 14)

	// remove a database
	m.mongoCollector = &mockMongo{
		serverStatusResponse: "{}",
		dbStatsResponse:      "{}",
		databaseNames:        []string{"db1"},
	}
	ms = m.Collect()
	assert.True(t, m.charts.Get("database_collections").Dims[1].Obsolete)
	assert.Len(t, ms, 7)

	// add two databases
	m.mongoCollector = &mockMongo{
		serverStatusResponse: "{}",
		dbStatsResponse:      "{}",
		databaseNames:        []string{"db1", "db2", "db3"},
	}
	ms = m.Collect()
	assert.Len(t, m.charts.Get("database_collections").Dims, 3)
	assert.Len(t, ms, 21)
}

func TestMongo_Collect_DbStats_Fail(t *testing.T) {
	m := New()
	m.mongoCollector = &mockMongo{
		serverStatusResponse: "{}",
		dbStatsResponse:      "",
	}
	m.Config.Databases.Includes = []string{"* *"}
	m.URI = "mongodb://localhost"
	m.Init()
	ms := m.Collect()
	assert.Len(t, ms, 0)
}

func TestMongo_Collect_DbStats_EmptyMatcher(t *testing.T) {
	m := New()
	m.mongoCollector = &mockMongo{
		serverStatusResponse: "{}",
		dbStatsResponse:      "{}",
	}
	m.Config.Databases.Includes = []string{"* not_matching"}
	m.URI = "mongodb://localhost"
	m.Init()
	ms := m.Collect()
	assert.Len(t, ms, 0)
}

func TestMongo_Collect_ReplSetStatus(t *testing.T) {
	m := New()
	m.mongoCollector = &mockMongo{
		serverStatusResponse: "{}",
		dbStatsResponse:      "{}",
		replicaSet:           true,
	}
	m.Config.Databases.Includes = []string{"* *"}
	m.URI = "mongodb://localhost"
	m.Init()
	_ = m.Collect()
	assert.True(t, m.charts.Has(replicationLag))
	assert.True(t, m.charts.Has(replicationHeartbeatLatency))
	assert.True(t, m.charts.Has(replicationNodePing))
}

func TestMongo_Collect_ReplSetStatusAddRemove(t *testing.T) {
	m := New()
	m.mongoCollector = &mockMongo{
		serverStatusResponse: "{}",
		dbStatsResponse:      "{}",
		replicaSet:           true,
		replicaNodes:         []string{"node1"},
	}
	m.Config.Databases.Includes = []string{"* *"}
	m.URI = "mongodb://localhost"
	m.Init()
	_ = m.Collect()
	assert.True(t, m.charts.Get(replicationLag).HasDim(replicationLagDimPrefix+"node1"))
	assert.True(t, m.charts.Get(replicationHeartbeatLatency).HasDim(replicationHeartbeatLatencyDimPrefix+"node1"))
	assert.True(t, m.charts.Get(replicationNodePing).HasDim(replicationNodePingDimPrefix+"node1"))

	m.mongoCollector = &mockMongo{
		serverStatusResponse: "{}",
		dbStatsResponse:      "{}",
		replicaSet:           true,
		replicaNodes:         []string{"node2"},
	}
	_ = m.Collect()
	// node2 dimensions added
	assert.True(t, m.charts.Get(replicationLag).HasDim(replicationLagDimPrefix+"node2"))
	assert.True(t, m.charts.Get(replicationHeartbeatLatency).HasDim(replicationHeartbeatLatencyDimPrefix+"node2"))
	assert.True(t, m.charts.Get(replicationNodePing).HasDim(replicationNodePingDimPrefix+"node2"))

	// node1 dimensions removed
	assert.True(t, m.charts.Get(replicationLag).GetDim(replicationLagDimPrefix+"node1").Obsolete)
	assert.True(t, m.charts.Get(replicationHeartbeatLatency).GetDim(replicationHeartbeatLatencyDimPrefix+"node1").Obsolete)
	assert.True(t, m.charts.Get(replicationNodePing).GetDim(replicationNodePingDimPrefix+"node1").Obsolete)
}

func TestMongo_Incomplete(t *testing.T) {
	m := New()
	m.mongoCollector = &mockMongo{}
	ms := m.Collect()
	assert.Len(t, ms, 0)
}

func TestMongo_Cleanup(t *testing.T) {
	m := New()
	connector := &mockMongo{}
	m.mongoCollector = connector
	m.Cleanup()
	assert.True(t, connector.closeCalled)
}

// mock for unit testing the mongo database as the driver
// doesn't use interfaces.
type mockMongo struct {
	connector
	serverStatusResponse string
	dbStatsResponse      string
	closeCalled          bool
	replicaSet           bool
	databaseNames        []string
	replicaNodes         []string
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
	return m.databaseNames, nil
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
	var ping int64 = 10
	var now = time.Now()
	status := replSetStatus{
		Date:    now,
		Members: nil,
	}
	for _, node := range m.replicaNodes {
		status.Members = append(status.Members,
			struct {
				Name                  string     `bson:"name"`
				State                 int        `bson:"state"`
				OptimeDate            time.Time  `bson:"optimeDate"`
				LastHeartbeat         *time.Time `bson:"lastHeartbeat"`
				LastHeartbeatReceived *time.Time `bson:"lastHeartbeatRecv"`
				PingMs                *int64     `bson:"pingMs"`
			}{
				Name:                  node,
				State:                 0,
				OptimeDate:            now,
				LastHeartbeat:         &now,
				LastHeartbeatReceived: &now,
				PingMs:                &ping,
			},
		)
	}
	return &status, nil
}
