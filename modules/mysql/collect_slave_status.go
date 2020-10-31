package mysql

import (
	"strconv"
	"strings"

	"github.com/blang/semver/v4"
)

const (
	querySlaveStatus     = "SHOW SLAVE STATUS"
	queryAllSlavesStatus = "SHOW ALL SLAVES STATUS"
)

var slaveStatusMetrics = []string{
	"Seconds_Behind_Master",
	"Slave_SQL_Running",
	"Slave_IO_Running",
}

func (m *MySQL) collectSlaveStatus(collected map[string]int64) error {
	// https://mariadb.com/docs/reference/es/sql-statements/SHOW_ALL_SLAVES_STATUS/
	mariaDBMinVer := semver.Version{Major: 10, Minor: 2, Patch: 0}
	var query string
	if m.isMariaDB && m.version.GTE(mariaDBMinVer) {
		query = queryAllSlavesStatus
	} else {
		query = querySlaveStatus
	}
	m.Debugf("executing query: '%s'", query)

	rows, err := m.db.Query(query)
	if err != nil {
		return err
	}
	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		return err
	}

	values := nullStringsFromColumns(columns)

	for rows.Next() {
		if err := rows.Scan(values...); err != nil {
			return err
		}

		set := rowAsMap(columns, values)

		var conn string
		if m.isMariaDB {
			conn = set["Connection_name"]
		} else {
			conn = set["Channel_Name"]
		}

		suffix := slaveMetricSuffix(conn)

		if !m.collectedReplConns[conn] {
			m.collectedReplConns[conn] = true
			m.addSlaveReplicationConnCharts(conn)
		}

		for _, name := range slaveStatusMetrics {
			v, ok := set[name]
			if !ok {
				continue
			}
			value, err := parseSlaveStatusValue(name, v)
			if err != nil {
				continue
			}
			collected[strings.ToLower(name+suffix)] = value
		}
	}
	return rows.Err()
}

func parseSlaveStatusValue(name, value string) (int64, error) {
	switch name {
	case "Slave_SQL_Running":
		value = convertSlaveSQLRunning(value)
	case "Slave_IO_Running":
		value = convertSlaveIORunning(value)
	}
	return strconv.ParseInt(value, 10, 64)
}

func convertSlaveSQLRunning(value string) string {
	switch value {
	case "Yes":
		return "1"
	default:
		return "0"
	}
}

func convertSlaveIORunning(value string) string {
	// NOTE: There is 'Connecting' state and probably others
	switch value {
	case "Yes":
		return "1"
	default:
		return "0"
	}
}

func slaveMetricSuffix(conn string) string {
	if conn == "" {
		return ""
	}
	return "_" + conn
}
