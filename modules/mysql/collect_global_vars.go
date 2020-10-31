package mysql

import (
	"strconv"
	"strings"
)

const (
	queryGlobalVariables = "SHOW GLOBAL VARIABLES WHERE " +
		"Variable_name LIKE 'max_connections' " +
		"OR " +
		"Variable_name LIKE 'table_open_cache'"
)

var globalVariablesMetrics = []string{
	"max_connections",
	"table_open_cache",
}

func (m *MySQL) collectGlobalVariables(collected map[string]int64) error {
	// MariaDB: https://mariadb.com/kb/en/server-system-variables/
	// MySQL: https://dev.mysql.com/doc/refman/8.0/en/server-system-variable-reference.html
	m.Debugf("executing query: '%s'", queryGlobalVariables)

	rows, err := m.db.Query(queryGlobalVariables)
	if err != nil {
		return err
	}
	defer rows.Close()

	set, err := rowsAsMap(rows)
	if err != nil {
		return err
	}

	for _, name := range globalVariablesMetrics {
		v, ok := set[name]
		if !ok {
			continue
		}
		value, err := parseGlobalVariable(v)
		if err != nil {
			continue
		}
		collected[strings.ToLower(name)] = value
	}
	return nil
}

func parseGlobalVariable(value string) (int64, error) {
	return strconv.ParseInt(value, 10, 64)
}
