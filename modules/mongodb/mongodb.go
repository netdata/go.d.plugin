package mongo

import (
	"context"
	"fmt"
	"time"

	"github.com/netdata/go.d.plugin/agent/module"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func init() {
	module.Register("mongodb", module.Creator{
		Defaults: module.Defaults{
			UpdateEvery:        module.UpdateEvery,
			AutoDetectionRetry: module.AutoDetectionRetry,
			Priority:           module.Priority,
			Disabled:           false,
		},
		Create: func() module.Module { return New() },
	})
}

func New() *Mongo {
	return &Mongo{
		Config: Config{
			Local: Local{
				Host: defaultHost,
				Port: defaultPort,
			},
			Auth: Auth{
				Host:   defaultHost,
				Port:   defaultPort,
				Authdb: defaultAuthDb,
				User:   defaultUser,
				Pass:   defaultPass,
			},
			Timeout:       defaultTimeout,
			ConnectionStr: defaultConnectionStr,
		},
		charts: &module.Charts{},
	}
}

type Mongo struct {
	module.Base
	Config `yaml:",inline"`
	client *mongo.Client
	charts *module.Charts
}

func (m *Mongo) Init() bool {
	m.Infof("initializing mongodb")
	validLocalConfig := m.Local.valid()
	validAuthConfig := m.Auth.valid()
	if !validLocalConfig && !validAuthConfig && m.ConnectionStr == "" {
		m.Errorf("config validation: all local and auth "+
			"and connection string config values are empty."+
			"Attempting to connect to %s:%d", defaultHost, defaultPort)
		return false
	}

	var err error
	m.charts, err = m.initCharts()
	if err != nil {
		m.Errorf("init charts: %v", err)
		return false
	}
	return true
}

func (m *Mongo) Check() bool {
	return len(m.Collect()) > 0
}

func (m *Mongo) Charts() *module.Charts {
	return m.charts
}

func (m *Mongo) Collect() map[string]int64 {
	if m.client == nil {
		var err error
		err = m.initMongoClient()
		if err != nil {
			m.Errorf("init mongo client: %v", err)
			return nil
		}
	}

	ms := make(map[string]int64)
	m.serverStatusCollect(ms)
	if len(ms) == 0 {
		m.Warning("zero collected values")
		return nil
	}
	return ms
}

func (m *Mongo) Cleanup() {
	if m.client == nil {
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), m.Timeout*time.Second)
	defer cancel()
	err := m.client.Disconnect(ctx)
	if err != nil {
		m.Warningf("cleanup: error on closing mongo client: %v", err)
	}
	m.client = nil
}

func (m *Mongo) initMongoClient() error {
	var connectionString string
	switch {
	case m.ConnectionStr != "":
		connectionString = m.ConnectionStr
	case m.Auth.valid():
		connectionString = m.Auth.connectionString()
	case m.Local.valid():
		connectionString = m.Local.connectionString()
	default:
		connectionString = fmt.Sprintf("mongodb://%s:%d", defaultHost, defaultPort)
	}
	client, err := mongo.NewClient(options.Client().ApplyURI(connectionString))
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), m.Timeout*time.Second)
	defer cancel()
	err = client.Connect(ctx)
	if err != nil {
		return err
	}
	m.client = client
	return nil
}

func (m *Mongo) initCharts() (*module.Charts, error) {
	return &serverStatusCharts, nil
}
