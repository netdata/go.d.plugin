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
			name: defaultName,
			Local: Local{
				host: defaultHost,
				port: defaultPort,
			},
			Auth: Auth{
				host:   defaultHost,
				port:   defaultPort,
				authdb: defaultAuthDb,
				user:   defaultUser,
				pass:   defaultPass,
			},
			Timeout:       defaultTimeout,
			ConnectionStr: defaultConnectionStr,
		},
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
		m.Warningf("config validation: all local and auth "+
			"and connection string config values are empty."+
			"Attempting to connect to %s:%d", defaultHost, defaultPort)
	}

	var err error
	// init client but do not attempt any IO
	m.client, err = m.initMongoClient()
	if err != nil {
		m.Errorf("init mongo client: %v", err)
		return false
	}

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

func (m *Mongo) initMongoClient() (*mongo.Client, error) {
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
		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), m.Timeout*time.Second)
	defer cancel()
	err = client.Connect(ctx)
	if err != nil {
		return nil, err
	}
	return client, nil
}

func (m *Mongo) initCharts() (*module.Charts, error) {
	//charts := module.Charts{}
	//serverStatusCharts := serverStatusCharts()
	//m.Info("get serverStatusCharts")
	//err := charts.Add(serverStatusCharts...)
	//if err != nil {
	//	m.Error(err)
	//	return nil, err
	//}
	return &serverStatusCharts, nil
}
