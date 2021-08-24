package mongo

import (
	"sync"
	"time"

	"github.com/netdata/go.d.plugin/agent/module"
	"github.com/netdata/go.d.plugin/pkg/matcher"
)

type Config struct {
	URI       string             `yaml:"uri"`
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
			URI:     "mongodb://localhost:27017",
			Databases: matcher.SimpleExpr{
				Includes: []string{},
				Excludes: []string{},
			},
		},
		charts:                &module.Charts{},
		optionalChartsEnabled: make(map[string]bool),
		discoveredDBs:         make([]string, 0),
		mongoCollector:        &mongoCollector{},
		addReplChartsOnce:     sync.Once{},
		replSetMembers:        make([]string, 0),
		replSetDimsEnabled:    make(map[string]bool),
	}
}

type Mongo struct {
	module.Base
	Config                `yaml:",inline"`
	mongoCollector        connector
	charts                *module.Charts
	databasesMatcher      matcher.Matcher
	optionalChartsEnabled map[string]bool
	discoveredDBs         []string
	chartsDbStats         *module.Charts
	replSetMembers        []string
	replSetDimsEnabled    map[string]bool
	addReplChartsOnce     sync.Once
}

func (m *Mongo) Init() bool {
	m.Infof("initializing mongodb")
	if m.URI == "" {
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
	if err := m.mongoCollector.initClient(m.URI, m.Timeout); err != nil {
		m.Errorf("init mongo client: %v", err)
		return nil
	}

	ms := map[string]int64{}
	if err := m.collectServerStatus(ms); err != nil {
		m.Errorf("couldn't collecting server status metrics: %s", err)
		return nil
	}

	if err := m.collectDbStats(ms); err != nil {
		m.Errorf("couldn't collecting dbstats metrics: %s", err)
	}

	if m.mongoCollector.isReplicaSet() {
		// if we have replica set based on the serverStatus response
		// we add once the charts during runtime
		m.addReplChartsOnce.Do(func() {
			if err := m.charts.Add(*replCharts.Copy()...); err != nil {
				m.Errorf("failed to add replica set chart: %v", err)
			}
		})

		if err := m.collectReplSetStatus(ms); err != nil {
			m.Errorf("couldn't collecting replSetStatus metrics: %s", err)
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
	var charts module.Charts
	err := charts.Add(serverStatusCharts...)
	if err != nil {
		return nil, err
	}

	m.chartsDbStats = dbStatsCharts.Copy()
	for _, chart := range *m.chartsDbStats {
		err = charts.Add(chart)
		if err != nil {
			return &charts, err
		}
	}
	return &charts, nil
}
