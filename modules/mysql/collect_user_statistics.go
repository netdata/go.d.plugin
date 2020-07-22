package mysql

import (
	"strconv"
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
	rows, err := m.db.Query(queryUserStatistics)
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
		user := set["User"]
		prefix := "userstats_" + user + "_"

		if !m.collectedUsers[user] {
			m.collectedUsers[user] = true
			m.addUserStatisticsCharts(user)
		}

		for _, name := range userStatisticsMetrics {
			strValue, ok := set[name]
			if !ok {
				continue
			}
			value, err := parseUserStatisticsValue(strValue)
			if err != nil {
				continue
			}
			collected[prefix+name] = value
		}
	}
	return rows.Err()
}

func parseUserStatisticsValue(value string) (int64, error) {
	return strconv.ParseInt(value, 10, 64)
}
