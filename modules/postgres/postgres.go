// SPDX-License-Identifier: GPL-3.0-or-later

package postgres

import (
	"database/sql"
	"time"

	"github.com/netdata/go.d.plugin/agent/module"
	"github.com/netdata/go.d.plugin/pkg/matcher"
	"github.com/netdata/go.d.plugin/pkg/metrics"
	"github.com/netdata/go.d.plugin/pkg/web"

	_ "github.com/jackc/pgx/v4/stdlib"
)

func init() {
	module.Register("postgres", module.Creator{
		Create: func() module.Module { return New() },
	})
}

func New() *Postgres {
	return &Postgres{
		Config: Config{
			Timeout:            web.Duration{Duration: time.Second * 2},
			DSN:                "postgres://postgres:postgres@127.0.0.1:5432/postgres",
			XactTimeHistogram:  []float64{.1, .5, 1, 2.5, 5, 10},
			QueryTimeHistogram: []float64{.1, .5, 1, 2.5, 5, 10},
		},
		charts:  baseCharts.Copy(),
		dbConns: make(map[string]*dbConn),
		mx: &pgMetrics{
			dbs:       make(map[string]*dbMetrics),
			tables:    make(map[string]*tableMetrics),
			replApps:  make(map[string]*replStandbyAppMetrics),
			replSlots: make(map[string]*replSlotMetrics),
		},
		recheckSettingsEvery: time.Minute * 30,
		doBloatEvery:         time.Minute * 5,
	}
}

type Config struct {
	DSN                string       `yaml:"dsn"`
	Timeout            web.Duration `yaml:"timeout"`
	DBSelector         string       `yaml:"collect_databases_matching"`
	XactTimeHistogram  []float64    `yaml:"transaction_time_histogram"`
	QueryTimeHistogram []float64    `yaml:"query_time_histogram"`
}

type (
	Postgres struct {
		module.Base
		Config `yaml:",inline"`

		charts *module.Charts

		db      *sql.DB
		dbConns map[string]*dbConn

		superUser *bool
		pgVersion int

		currentDB string
		dbSr      matcher.Matcher

		mx *pgMetrics

		recheckSettingsTime  time.Time
		recheckSettingsEvery time.Duration

		doBloatTime  time.Time
		doBloatEvery time.Duration
	}
	dbConn struct {
		db         *sql.DB
		connErrors int
	}
)

func (p *Postgres) Init() bool {
	err := p.validateConfig()
	if err != nil {
		p.Errorf("config validation: %v", err)
		return false
	}

	sr, err := p.initDBSelector()
	if err != nil {
		p.Errorf("config validation: %v", err)
		return false
	}
	p.dbSr = sr

	p.mx.xactTimeHist = metrics.NewHistogramWithRangeBuckets(p.XactTimeHistogram)
	p.addTransactionsRunTimeHistogramChart()

	p.mx.queryTimeHist = metrics.NewHistogramWithRangeBuckets(p.QueryTimeHistogram)
	p.addQueriesRunTimeHistogramChart()

	return true
}

func (p *Postgres) Check() bool {
	return len(p.Collect()) > 0
}

func (p *Postgres) Charts() *module.Charts {
	return p.charts
}

func (p *Postgres) Collect() map[string]int64 {
	mx, err := p.collect()
	if err != nil {
		p.Error(err)
	}

	if len(mx) == 0 {
		return nil
	}
	return mx
}

func (p *Postgres) Cleanup() {
	if p.db == nil {
		return
	}
	if err := p.db.Close(); err != nil {
		p.Warningf("cleanup: error on closing the Postgres database [%s]: %v", p.DSN, err)
	}
	p.db = nil

	for dbname, conn := range p.dbConns {
		delete(p.dbConns, dbname)
		if conn.db != nil {
			_ = conn.db.Close()
		}
	}
}
