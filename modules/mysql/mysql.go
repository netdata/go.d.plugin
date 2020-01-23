package mysql

import (
	"database/sql"
	"fmt"
	"strconv"
	"strings"
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

const (
	querySlaveStatus    = "SHOW SLAVE STATUS"
	queryGlobalStatus   = "SHOW GLOBAL STATUS"
	queryMaxConnections = "SHOW GLOBAL VARIABLES LIKE 'max_connections'"
)

// MySQL is the mysql database module.
type MySQL struct {
	module.Base
	db *sql.DB
	// i.e user:password@/dbname
	DSN     string `yaml:"dsn"`
	doSlave bool

	charts *Charts
}

// New creates and returns a new empty MySQL module.
func New() *MySQL { return &MySQL{charts: charts.Copy(), doSlave: true} }

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

	if err = db.Ping(); err != nil {
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

	if _, ok := metrics["seconds_behind_master"]; ok {
		_ = m.charts.Add(*slaveCharts.Copy()...)
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
	metrics := make(map[string]int64)

	if err := m.collectGlobalStats(metrics); err != nil {
		m.Errorf("error on collecting global stats: %v", err)
		return nil
	}

	// TODO: do better
	if m.doSlave {
		if err := m.collectSlaveStatus(metrics); err != nil {
			m.Errorf("error on collecting slave status: %v", err)
			m.doSlave = false
		}
	}

	if err := m.collectMaxConnections(metrics); err != nil {
		m.Errorf("error on determining max connections: %v", err)
		return nil
	}

	return metrics
}

func (m *MySQL) collectGlobalStats(metrics map[string]int64) error {
	rows, err := m.db.Query(queryGlobalStatus)

	if err != nil {
		return err
	}

	defer rows.Close()

	for rows.Next() {
		var (
			varName string
			// it does not always returns int64,
			// i.e public key []uint8(bytes) we can't catch it before Scan, so interface{}.
			value interface{}
		)

		if err = rows.Scan(&varName, &value); err != nil {
			return err
		}

		val, err := strconv.ParseInt(fmt.Sprintf("%s", value), 10, 64)

		if err != nil {
			continue
		}

		metrics[strings.ToLower(varName)] = val
	}

	v1, ok1 := metrics["threads_created"]
	v2, ok2 := metrics["connections"]

	// NOTE: not sure this check is needed
	if ok1 && ok2 {
		metrics["thread_cache_misses"] = int64(float64(v1) / float64(v2) * 10000)
	}

	return nil
}

// https://dev.mysql.com/doc/refman/8.0/en/show-variables.html
func (m *MySQL) collectMaxConnections(metrics map[string]int64) error {
	// only one result, i.e "max_conections" = 151
	rows, err := m.db.Query(queryMaxConnections)
	if err != nil {
		return err
	}

	defer rows.Close()

	if !rows.Next() {
		return nil
	}

	var (
		varName string
		value   int64
	)

	err = rows.Scan(&varName, &value)
	if err != nil {
		return err
	}

	metrics["max_connections"] = value

	return nil
}

// https://dev.mysql.com/doc/refman/8.0/en/show-slave-status.html
func (m *MySQL) collectSlaveStatus(metrics map[string]int64) error {
	// // https://github.com/gdaws/mysql-slave-status
	rows, err := m.db.Query(querySlaveStatus)
	if err != nil {
		return err
	}

	defer rows.Close()

	columns, err := rows.Columns()

	if err != nil {
		return err
	}

	if !rows.Next() {
		if err = rows.Err(); err != nil {
			return err
		}
		return nil
	}

	values := make([]interface{}, len(columns))

	for index := range values {
		values[index] = new(sql.NullString)
	}

	err = rows.Scan(values...)

	if err != nil {
		return err
	}

	for index, columnName := range columns {
		switch columnName {
		case "Slave_SQL_Running", "Slave_IO_Running":
			val := *values[index].(*sql.NullString)
			if !val.Valid {
				continue
			}
			if val.String == "Yes" {
				metrics[strings.ToLower(columnName)] = 1
			} else {
				metrics[strings.ToLower(columnName)] = 0
			}
		case "Seconds_Behind_Master":
			val := *values[index].(*sql.NullString)
			if !val.Valid {
				continue
			}
			v, err := strconv.ParseInt(val.String, 10, 64)
			if err != nil {
				continue
			}
			metrics[strings.ToLower(columnName)] = v

		}
	}

	return nil
}

// // CompatibleMinimumVersion is the minimum required version of the mysql server.
// const CompatibleMinimumVersion = 5.1
//
// func (m *MySQL) getMySQLVersion() float64 {
// 	var versionStr string
// 	var versionNum float64
// 	if err := m.db.QueryRow("SELECT @@version").Scan(&versionStr); err == nil {
// 		versionNum, _ = strconv.ParseFloat(regexp.MustCompile(`^\d+\.\d+`).FindString(versionStr), 64)
// 	}
//
// 	return versionNum
// }
