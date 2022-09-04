package proxysql

import (
	"strconv"
	"strings"
)

const queryMysqlusers = "SELECT * FROM stats_mysql_users"

var mysqlusers = []string{
	"frontend_connections",
	"frontend_max_connections",
}

func (p *ProxySQL) collectMysqlUsers(collected map[string]int64) error {
	// https://proxysql.com/documentation/stats-statistics/#stats_mysql_users
	p.Debugf("executing query: '%s'", queryMysqlusers)

	rows, err := p.db.Query(queryMysqlusers)
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
		username := set["username"]
		prefix := "mysql_user_" + username + "_"

		p.addMysqlUsersCharts(username)

		for _, name := range mysqlusers {
			v, ok := set[name]
			if !ok {
				continue
			}
			value, err := parseMysqlUsersValue(v)
			if err != nil {
				continue
			}
			collected[strings.ToLower(prefix+name)] = value
		}
	}
	return rows.Err()
}

func parseMysqlUsersValue(value string) (int64, error) {
	return strconv.ParseInt(value, 10, 64)
}
