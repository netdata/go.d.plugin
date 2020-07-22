package mysql

import (
	"strconv"
)

const querySlaveStatus = "SHOW SLAVE STATUS"

var slaveStatusMetrics = []string{
	"Seconds_Behind_Master",
	"Slave_SQL_Running",
	"Slave_IO_Running",
}

func (m *MySQL) collectSlaveStatus(collected map[string]int64) error {
	// https://dev.mysql.com/doc/refman/18.0/en/show-slave-status.html
	// https://dev.mysql.com/doc/refman/8.0/en/replication-channels.html
	// https://github.com/gdaws/mysql-slave-status
	rows, err := m.db.Query(querySlaveStatus)
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
		// NOTE: 'Connection_name' in MariaDB 10.0+, exposed only in 'SHOW ALL SLAVES STATUS'
		// https://mariadb.com/kb/en/show-slave-status/#multi-source
		channel := set["Channel_Name"]
		suffix := slaveMetricSuffix(channel)

		if !m.collectedChannels[channel] {
			m.collectedChannels[channel] = true
			m.addSlaveReplicationChannelCharts(channel)
		}

		for _, name := range slaveStatusMetrics {
			strValue, ok := set[name]
			if !ok {
				continue
			}
			value, err := parseSlaveStatusValue(name, strValue)
			if err != nil {
				continue
			}
			collected[name+suffix] = value
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
	switch value {
	case "Yes":
		return "1"
	case "Connecting":
		return "2"
	default:
		return "0"
	}
}

func slaveMetricSuffix(channel string) string {
	if channel == "" {
		return ""
	}
	return "_" + channel
}
