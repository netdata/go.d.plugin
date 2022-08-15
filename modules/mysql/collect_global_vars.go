// SPDX-License-Identifier: GPL-3.0-or-later

package mysql

import (
	"context"
	"strconv"
	"strings"
)

const (
	queryGlobalVariables = "SHOW GLOBAL VARIABLES WHERE " +
		"Variable_name LIKE 'max_connections' " +
		"OR " +
		"Variable_name LIKE 'table_open_cache' " +
		"OR " +
		"Variable_name LIKE 'disabled_storage_engines' " +
		"OR " +
		"Variable_name LIKE 'log_bin'"
)

var globalVariablesMetrics = []string{
	"max_connections",
	"table_open_cache",
	"disabled_storage_engines",
	"log_bin",
}

func (m *MySQL) collectGlobalVariables(collected map[string]int64) error {
	// MariaDB: https://mariadb.com/kb/en/server-system-variables/
	// MySQL: https://dev.mysql.com/doc/refman/8.0/en/server-system-variable-reference.html
	m.Debugf("executing query: '%s'", queryGlobalVariables)

	ctx, cancel := context.WithTimeout(context.Background(), m.Timeout.Duration)
	defer cancel()

	rows, err := m.db.QueryContext(ctx, queryGlobalVariables)
	if err != nil {
		return err
	}
	defer func() { _ = rows.Close() }()

	set, err := rowsAsMap(rows)
	if err != nil {
		return err
	}

	for _, name := range globalVariablesMetrics {
		v, ok := set[name]
		if !ok {
			continue
		}
		value, err := parseGlobalVariable(name, v)
		if err != nil {
			continue
		}
		collected[strings.ToLower(name)] = value
	}
	return nil
}

func parseGlobalVariable(name, value string) (int64, error) {
	switch name {
	case "disabled_storage_engines":
		value = convertStorageEngineValue(value)
	case "log_bin":
		value = convertBinlogValue(value)
	}
	return strconv.ParseInt(value, 10, 64)
}

func convertStorageEngineValue(val string) string {
	if strings.Contains(val, "MyISAM") {
		return "1"
	}
	return "0"
}

func convertBinlogValue(val string) string {
	if val == "OFF" {
		return "0"
	}
	return "1"
}
