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
	initClient(uri string, timeout time.Duration) error
	close() error
}

type mongoCollector struct {
	Client  *mongo.Client
	Timeout time.Duration
}

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
