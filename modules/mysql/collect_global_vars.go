package mysql

import (
	"strconv"
	"strings"
)

const (
	queryGlobalVariables = "SHOW GLOBAL VARIABLES"
)

var globalVariablesMetrics = []string{
	"max_connections",
	"table_open_cache",
}

func (m *MySQL) collectGlobalVariables(collected map[string]int64) error {
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
		strValue, ok := set[name]
		if !ok {
			continue
		}
		value, err := parseGlobalVariable(strValue)
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
