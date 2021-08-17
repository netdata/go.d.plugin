package mongo

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type connector interface {
	serverStatus() (*serverStatus, error)
	listDatabaseNames() ([]string, error)
	dbStats(databaseName string) (*dbStats, error)
	initClient(uri string, timeout time.Duration) error
	close() error
}

// mongoCollector interface that helps abstracting and mocking the database layer.
type mongoCollector struct {
	Client  *mongo.Client
	Timeout time.Duration
}

// serverStatus connects to the database and return the output of the
// `serverStatus` command.
func (m *mongoCollector) serverStatus() (*serverStatus, error) {
	var status *serverStatus
	command := bson.D{{Key: "serverStatus", Value: 1}, {Key: "metrics", Value: 0}, {Key: "repl", Value: 0}}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*m.Timeout)
	defer cancel()
	err := m.Client.Database("admin").RunCommand(ctx, command).Decode(&status)
	if err != nil {
		return nil, err
	}
	return status, err
}

// listDatabaseNames returns a string slice with the available databases on the server.
func (m *mongoCollector) listDatabaseNames() ([]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*m.Timeout)
	defer cancel()
	return m.Client.ListDatabaseNames(ctx, bson.M{})
}

// dbStats gets the `dbstats` metrics for a specific database.
func (m *mongoCollector) dbStats(databaseName string) (*dbStats, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*m.Timeout)
	defer cancel()
	var dbStats dbStats
	db := m.Client.Database(databaseName)
	if err := db.RunCommand(ctx, bson.M{"dbStats": 1}).Decode(&dbStats); err != nil {
		return nil, err
	}
	return &dbStats, nil
}

// initClient initialises the database client if is not initialised.
func (m *mongoCollector) initClient(uri string, timeout time.Duration) error {
	if m.Client != nil {
		return nil
	}
	m.Timeout = timeout
	client, err := mongo.NewClient(options.Client().ApplyURI(uri))
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), m.Timeout*time.Second)
	defer cancel()
	err = client.Connect(ctx)
	if err != nil {
		return err
	}
	m.Client = client
	return nil
}

// close the database client and all its background goroutines.
func (m *mongoCollector) close() error {
	if m.Client == nil {
		return nil
	}
	ctx, cancel := context.WithTimeout(context.Background(), m.Timeout*time.Second)
	defer cancel()
	err := m.Client.Disconnect(ctx)
	if err != nil {
		return err
	}
	m.Client = nil
	return nil
}
