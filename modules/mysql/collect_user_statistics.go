// SPDX-License-Identifier: GPL-3.0-or-later

package mysql

import (
	"strings"
)

const queryShowUserStatistics = "SHOW USER_STATISTICS;"

func (m *MySQL) collectUserStatistics(mx map[string]int64) error {
	// https://mariadb.com/kb/en/user-statistics/
	// https://mariadb.com/kb/en/information-schema-user_statistics-table/
	q := queryShowUserStatistics
	m.Debugf("executing query: '%s'", q)

	var user, prefix string
	_, err := m.collectQuery(q, func(column, value string, _ bool) {
		switch column {
		case "User":
			user = value
			prefix = "userstats_" + user + "_"
			if !m.collectedUsers[user] {
				m.collectedUsers[user] = true
				m.addUserStatisticsCharts(user)
			}
		case "Cpu_time":
			mx[strings.ToLower(prefix+column)] = int64(parseFloat(value) * 1000)
		case "Rows_read",
			"Rows_sent",
			"Rows_deleted",
			"Rows_inserted",
			"Rows_updated",
			"Select_commands",
			"Update_commands",
			"Other_commands":
			mx[strings.ToLower(prefix+column)] = parseInt(value)
		}
	})
	return err
}
