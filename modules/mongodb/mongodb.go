// SPDX-License-Identifier: GPL-3.0-or-later

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
			Timeout: 1,
			URI:     "mongodb://localhost:27017",
			Databases: matcher.SimpleExpr{
				Includes: []string{"* *"},
				Excludes: []string{},
			},
		},
		charts:             &module.Charts{},
		optionalCharts:     make(map[string]bool),
		discoveredDBs:      make([]string, 0),
		shardNodesDims:     make(map[string]bool),
		mongoCollector:     &mongoCollector{},
		addReplChartsOnce:  sync.Once{},
		addShardChartsOnce: sync.Once{},
		replSetMembers:     make([]string, 0),
		replSetDimsEnabled: make(map[string]bool),

		dbs: make(map[string]bool),
	}
}

type Mongo struct {
	module.Base
	Config             `yaml:",inline"`
	mongoCollector     connector
	charts             *module.Charts
	dbMatcher          matcher.Matcher
	optionalCharts     map[string]bool
	discoveredDBs      []string
	shardNodesDims     map[string]bool
	chartsDbStats      *module.Charts
	replSetMembers     []string
	replSetDimsEnabled map[string]bool
	addReplChartsOnce  sync.Once
	addShardChartsOnce sync.Once

	dbs map[string]bool
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
		m.dbMatcher = mMatcher
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
	ms, err := m.collect()
	if err != nil {
		m.Error(err)
	}
	if len(ms) == 0 {
		m.Warning("no values collected")
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
	err := charts.Add(*serverStatusCharts.Copy()...)
	if err != nil {
		return nil, err
	}

	return &charts, nil
}
