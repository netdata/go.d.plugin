package mongo

import (
	"time"

	"github.com/netdata/go.d.plugin/agent/module"
	"github.com/netdata/go.d.plugin/pkg/matcher"
)

type Config struct {
	Uri       string             `yaml:"uri"`
	Timeout   time.Duration      `yaml:"timeout"`
	Databases matcher.SimpleExpr `yaml:"databases"`
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
			Databases: matcher.SimpleExpr{
				Includes: []string{},
				Excludes: []string{},
			},
		},
		charts:                &module.Charts{},
		optionalChartsEnabled: make(map[string]bool),
		dimsEnabled:           make(map[string]bool),
		databases:             make([]string, 0),
		mongoCollector:        &mongoCollector{},
	}
}

type Mongo struct {
	module.Base
	Config                `yaml:",inline"`
	mongoCollector        connector
	charts                *module.Charts
	databasesMatcher      matcher.Matcher
	optionalChartsEnabled map[string]bool
	databases             []string
	dimsEnabled           map[string]bool
}

func (m *Mongo) Init() bool {
	m.Infof("initializing mongodb")
	if m.Uri == "" {
		m.Errorf("connection URI is empty")
		return false
	}

	if !m.Databases.Empty() {
		mMatcher, err := m.Databases.Parse()
		if err != nil {
			m.Errorf("error on creating 'databases' matcher : %v", err)
			return false
		}
		m.databasesMatcher = mMatcher
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

func (m Mongo) initCharts() (*module.Charts, error) {
	var charts module.Charts
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
