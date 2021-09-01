package mongo

import (
	"encoding/json"
	"errors"
	"time"

	v5_0_0 "github.com/netdata/go.d.plugin/modules/mongodb/testdata/v5.0.0"
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
