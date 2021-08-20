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
	}
	m.Config.Databases.Includes = []string{"* *"}
	m.URI = "mongodb://localhost"
	m.Init()
	ms := m.Collect()
	assert.Len(t, ms, reflect.ValueOf(dbStats{}).NumField())
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
	assert.True(t, m.charts.Has("replication_lag"))
	assert.True(t, m.charts.Has("replication_heartbeat_latency"))
	assert.True(t, m.charts.Has("replication_node_ping"))
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
	return []string{"db"}, nil
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
	return &replSetStatus{
		Date: time.Now(),
		Members: []struct {
			Name                  string     `bson:"name"`
			State                 int        `bson:"state"`
			OptimeDate            time.Time  `bson:"optimeDate"`
			LastHeartbeat         *time.Time `bson:"lastHeartbeat"`
			LastHeartbeatReceived *time.Time `bson:"lastHeartbeatRecv"`
			PingMs                *int64     `bson:"pingMs"`
		}{
			{
				Name:                  "1",
				State:                 0,
				OptimeDate:            now,
				LastHeartbeat:         &now,
				LastHeartbeatReceived: &now,
				PingMs:                &ping,
			},
			{Name: "2"},
		},
	}, nil
}
