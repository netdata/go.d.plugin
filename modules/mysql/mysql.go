package mysql

import (
	"database/sql"
	"fmt"
	"sync"
	"time"

	"github.com/netdata/go.d.plugin/agent/module"

	"github.com/blang/semver/v4"
	_ "github.com/go-sql-driver/mysql"
)

func init() {
	module.Register("mysql", module.Creator{
		Create: func() module.Module { return New() },
	})
}

type (
	Config struct {
		DSN         string `yaml:"dsn"`
		MyCNF       string `yaml:"my.cnf"`
		UpdateEvery int    `yaml:"update_every"`
	}
	MySQL struct {
		module.Base
		Config `yaml:",inline"`

		db        *sql.DB
		isMariaDB bool
		version   *semver.Version

		addInnodbDeadlocksOnce *sync.Once
		addGaleraOnce          *sync.Once
		addQCacheOnce          *sync.Once
		addUserStatsCPUOnce    *sync.Once

		doSlaveStatus      bool
		collectedReplConns map[string]bool
		doUserStatistics   bool
		collectedUsers     map[string]bool

		charts *Charts
	}
)

func New() *MySQL {
	return &MySQL{
		Config: Config{
			DSN: "root@tcp(localhost:3306)/",
		},

		charts:                 charts.Copy(),
		addInnodbDeadlocksOnce: &sync.Once{},
		addGaleraOnce:          &sync.Once{},
		addQCacheOnce:          &sync.Once{},
		addUserStatsCPUOnce:    &sync.Once{},
		doSlaveStatus:          true,
		doUserStatistics:       true,
		collectedReplConns:     make(map[string]bool),
		collectedUsers:         make(map[string]bool),
	}
}

func (m *MySQL) Cleanup() {
	if m.db == nil {
		return
	}
	if err := m.db.Close(); err != nil {
		m.Errorf("cleanup: error on closing the mysql database [%s]: %v", m.DSN, err)
	}
	m.db = nil
}

func (m *MySQL) Init() bool {
	if m.MyCNF != "" {
		dsn, err := dsnFromFile(m.MyCNF)
		if err != nil {
			m.Error(err)
			return false
		}
		m.DSN = dsn
	}

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
		return fmt.Errorf("error on opening a connection with the mysql database [%s]: %v", m.DSN, err)
	}

	db.SetConnMaxLifetime(1 * time.Minute)

	if err := db.Ping(); err != nil {
		_ = db.Close()
		return fmt.Errorf("error on pinging the mysql database [%s]: %v", m.DSN, err)
	}

	m.db = db
	return nil
}

func (m *MySQL) Check() bool {
	return len(m.Collect()) > 0
}

func (m *MySQL) Charts() *Charts {
	return m.charts
}

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
