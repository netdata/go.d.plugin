package mysql

import (
	"database/sql"
	"regexp"
	"strconv"
	"time"

	"github.com/netdata/go.d.plugin/modules"
)

func init() {
	creator := modules.Creator{
		Create: func() modules.Module { return New() },
	}

	modules.Register("mysql", creator)
}

// MySQL is the mysql database module.
type MySQL struct {
	modules.Base
	db *sql.DB

	DSN string `yaml:"dsn"` // i.e user:password@/dbname
	// // or
	// User string `yaml:"user"`
	// Pass string `yaml:"pass"`

	// Host string `yaml:"host"`
	// Port int    `yaml:"port"`
	// // or
	// Socket string `yaml:"socket"`
}

// New creates and returns a new empty MySQL module.
func New() *MySQL {
	return &MySQL{}
}

// CompatibleMinimumVersion is the minimum required version of the mysql server.
const CompatibleMinimumVersion = 5.1

func (m *MySQL) getMySQLVersion() float64 {
	var versionStr string
	var versionNum float64
	if err := m.db.QueryRow("SELECT @@version").Scan(&versionStr); err == nil {
		versionNum, _ = strconv.ParseFloat(regexp.MustCompile(`^\d+\.\d+`).FindString(versionStr), 64)
	}

	return versionNum
}

// Cleanup performs cleanup.
func (m *MySQL) Cleanup() {
	err := m.db.Close()
	if err != nil {
		m.Errorf("cleanup: error on closing the mysql database [%s]: %v", m.DSN, err)
	}
}

// Init makes initialization of the MySQL mod.
func (m *MySQL) Init() bool {
	if m.DSN == "" {
		m.Errorf("dsn is missing")
		return false
	}

	// test the connectivity here.
	if err := m.openConnection(); err != nil {
		return false
	}

	if min, got := CompatibleMinimumVersion, m.getMySQLVersion(); min > 0 && got < min {
		m.Warningf("running with uncompatible mysql version [%v<%v]", got, min)
	}

	// post Init debug info.
	m.Debugf("using DSN [%s]", m.DSN)
	return true
}

func (m *MySQL) openConnection() error {
	if m.db != nil {
		if err := m.db.Ping(); err != nil {
			m.db.Close()
			m.db = nil

			return m.openConnection()
		}

		return nil
	}

	db, err := sql.Open("mysql", m.DSN)
	if err != nil {
		m.Errorf("error on opening a connection with the mysql database [%s]: %v", m.DSN, err)
		return err
	}
	db.SetConnMaxLifetime(1 * time.Minute)

	if err = db.Ping(); err != nil {
		db.Close()
		m.Errorf("error on pinging the mysql database [%s]: %v", m.DSN, err)
		return err
	}

	m.db = db
	return nil
}

// Check makes check.
func (m *MySQL) Check() bool {
	return len(m.Collect()) > 0
}

// Charts creates Charts.
func (m *MySQL) Charts() *Charts {
	return charts.Copy()
}

// Collect collects health checks and metrics for MySQL.
func (m *MySQL) Collect() map[string]int64 {
	return nil
}
