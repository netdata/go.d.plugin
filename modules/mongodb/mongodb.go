package mongo

import (
	"time"

	"github.com/netdata/go.d.plugin/agent/module"
)

type Config struct {
	Uri     string        `yaml:"uri"`
	Timeout time.Duration `yaml:"timeout"`
}

func init() {
	module.Register("mongodb", module.Creator{
		Create: func() module.Module { return New() },
	})
}

func New() *Mongo {
	return &Mongo{
		Config: Config{
			Timeout: 20,
			Uri:     "mongodb://localhost:27017",
		},
		charts:                &module.Charts{},
		optionalChartsEnabled: make(map[string]bool),
		dimsEnabled:           make(map[string]bool),
		mongoCollector:        &mongoCollector{},
	}
}

type Mongo struct {
	module.Base
	Config                `yaml:",inline"`
	mongoCollector        connector
	charts                *module.Charts
	optionalChartsEnabled map[string]bool
	dimsEnabled           map[string]bool
}

func (m *Mongo) Init() bool {
	m.Infof("initializing mongodb")
	if m.Uri == "" {
		m.Errorf("connection URI is empty")
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
	if err := m.mongoCollector.initClient(m.Uri, m.Timeout); err != nil {
		m.Errorf("init mongo client: %v", err)
		return nil
	}

	ms := map[string]int64{}
	for _, result := range []map[string]int64{
		m.serverStatusCollect(),
		m.dbStatsCollect(),
	} {
		for k, v := range result {
			ms[k] = v
		}
	}

	if len(ms) == 0 {
		m.Warning("zero collected values")
		return nil
	}
	return ms
}

func (m *Mongo) Cleanup() {
	err := m.mongoCollector.close()
	if err != nil {
		m.Warningf("cleanup: error on closing mongo client: %v", err)
	}
}

func (m *Mongo) initCharts() (*module.Charts, error) {
	charts := module.Charts{}
	err := charts.Add(serverStatusCharts...)
	if err != nil {
		return nil, err
	}
	err = charts.Add(dbStatsCharts...)
	if err != nil {
		return &charts, err
	}
	return &charts, nil
}
