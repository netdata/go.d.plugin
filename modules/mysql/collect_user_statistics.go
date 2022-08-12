// SPDX-License-Identifier: GPL-3.0-or-later

package mysql

import (
	"context"
	"strconv"
	"strings"
)

const queryUserStatistics = "SHOW USER_STATISTICS"

var userStatisticsMetrics = []string{
	"Cpu_time",
	"Rows_read",
	"Rows_sent",
	"Rows_deleted",
	"Rows_inserted",
	"Rows_updated",
	"Select_commands",
	"Update_commands",
	"Other_commands",
}

func (m *MySQL) collectUserStatistics(collected map[string]int64) error {
	// https://mariadb.com/kb/en/user-statistics/
	// https://mariadb.com/kb/en/information-schema-user_statistics-table/
	m.Debugf("executing query: '%s'", queryUserStatistics)

	ctx, cancel := context.WithTimeout(context.Background(), m.Timeout.Duration)
	defer cancel()

	rows, err := m.db.QueryContext(ctx, queryUserStatistics)
	if err != nil {
		return err
	}
	defer func() { _ = rows.Close() }()

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
		user := set["User"]
		prefix := "userstats_" + user + "_"

		if !m.collectedUsers[user] {
			m.collectedUsers[user] = true
			m.addUserStatisticsCharts(user)
		}

		for _, name := range userStatisticsMetrics {
			v, ok := set[name]
			if !ok {
				continue
			}
			value, err := parseUserStatisticsValue(name, v)
			if err != nil {
				continue
			}
			collected[strings.ToLower(prefix+name)] = value
		}
	}
	return rows.Err()
}

func parseUserStatisticsValue(name, value string) (int64, error) {
	if name == "Cpu_time" {
		v, err := strconv.ParseFloat(value, 64)
		return int64(v * 1000), err
	}
	return strconv.ParseInt(value, 10, 64)
}
