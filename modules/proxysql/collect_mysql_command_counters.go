package proxysql

import (
	"strconv"
	"strings"
)

const queryMysqlCommandCounters = "SELECT * FROM stats_mysql_commands_counters"

var mysqlCommandsCounters = []string{
	"Total_Time_us",
	"Total_cnt",
	"cnt_100us",
	"cnt_500us",
	"cnt_1ms",
	"cnt_5ms",
	"cnt_10ms",
	"cnt_10ms",
	"cnt_100ms",
	"cnt_500ms",
	"cnt_1s",
	"cnt_5s",
	"cnt_10s",
	"cnt_INFs",
}

func (p *ProxySQL) collectMysqlCommandCounters(collected map[string]int64) error {
	// https://proxysql.com/documentation/stats-statistics/#stats_mysql_commands_counters
	p.Debugf("executing query: '%s'", queryMysqlCommandCounters)

	rows, err := p.db.Query(queryMysqlCommandCounters)
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
		command := set["Command"]
		prefix := "mysql_command_" + command + "_"

		p.addMysqlCommandCountersCharts(command)

		for _, name := range mysqlCommandsCounters {
			v, ok := set[name]
			if !ok {
				continue
			}
			value, err := parseMysqlCommandCountersValue(v)
			if err != nil {
				continue
			}
			collected[strings.ToLower(prefix+name)] = value
		}
	}
	return rows.Err()
}

func parseMysqlCommandCountersValue(value string) (int64, error) {
	return strconv.ParseInt(value, 10, 64)
}
