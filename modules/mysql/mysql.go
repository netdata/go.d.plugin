package mysql

import (
	"database/sql"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/netdata/go-orchestrator/module"
)

func init() {
	creator := module.Creator{
		Defaults: module.Defaults{
			Disabled: true,
		},
		Create: func() module.Module { return New() },
	}

	module.Register("mysql", creator)
}

// MySQL is the mysql database module.
type MySQL struct {
	module.Base
	db  *sql.DB
	DSN string `yaml:"dsn"`

	doSlaveStats      bool
	doUserStatistics  bool
	collectedChannels map[string]bool
	collectedUsers    map[string]bool

	charts *Charts
}

// New creates and returns a new empty MySQL module.
func New() *MySQL {
	return &MySQL{
		DSN:               "netdata:password@tcp(127.0.0.1:3306)/",
		charts:            &module.Charts{},
		collectedChannels: make(map[string]bool),
		collectedUsers:    make(map[string]bool),
		doSlaveStats:      true,
		doUserStatistics:  true,
	}
}

// Cleanup performs cleanup.
func (m *MySQL) Cleanup() {
	if m.db == nil {
		return
	}
	if err := m.db.Close(); err != nil {
		m.Errorf("cleanup: error on closing the mysql database [%s]: %v", m.DSN, err)
	}
}

// Init makes initialization of the MySQL mod.
func (m *MySQL) Init() bool {
	if m.DSN == "" {
		m.Error("DSN not set")
		return false
	}

	if err := m.openConnection(); err != nil {
		m.Error(err)
		return false
	}

	m.Debugf("connected using DSN [%s]", m.DSN)
	return true
}

func (m *MySQL) openConnection() error {
	db, err := sql.Open("mysql", m.DSN)
	if err != nil {
		m.Errorf("error on opening a connection with the mysql database [%s]: %v", m.DSN, err)
		return err
	}

	db.SetConnMaxLifetime(1 * time.Minute)

	if err := db.Ping(); err != nil {
		_ = db.Close()
		m.Errorf("error on pinging the mysql database [%s]: %v", m.DSN, err)
		return err
	}

	m.db = db
	return nil
}

// Check makes check.
func (m *MySQL) Check() bool {
	metrics := m.Collect()

	if len(metrics) == 0 {
		return false
	}

	// FIXME: not sure this check is valid
	if _, ok := metrics["wsrep_local_recv_queue"]; ok {
		_ = m.charts.Add(*galeraCharts.Copy()...)
	}

	return true
}

// Charts creates Charts.
func (m *MySQL) Charts() *Charts {
	return m.charts
}

// Collect collects health checks and metrics for MySQL.
func (m *MySQL) Collect() map[string]int64 {
	mx, err := m.collect()
	if err != nil {
		m.Error(err)
	}

	if len(mx) == 0 {
		return nil
	}
	return mx
}
