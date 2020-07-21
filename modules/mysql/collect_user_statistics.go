package mysql

import (
	"strconv"
)

/*
MariaDB [(none)]> show user_statistics\G
*************************** 1. row ***************************
                       User: netdata
          Total_connections: 4
     Concurrent_connections: 0
             Connected_time: 1270
                  Busy_time: 0.020273
                   Cpu_time: 0.0192592
             Bytes_received: 646
                 Bytes_sent: 25426
       Binlog_bytes_written: 0
                  Rows_read: 0
                  Rows_sent: 22
               Rows_deleted: 0
              Rows_inserted: 0
               Rows_updated: 0
            Select_commands: 2
            Update_commands: 0
             Other_commands: 0
        Commit_transactions: 0
      Rollback_transactions: 0
         Denied_connections: 10
           Lost_connections: 0
              Access_denied: 2
              Empty_queries: 0
      Total_ssl_connections: 0
Max_statement_time_exceeded: 0
*/

const queryUserStatistics = "SHOW USER_STATISTICS"

var userStatisticsMetrics = []string{
	"Select_commands",
	"Update_commands",
	"Other_commands",
	"Cpu_time",
	"Rows_read",
	"Rows_sent",
	"Rows_deleted",
	"Rows_inserted",
	"Rows_updated",
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
