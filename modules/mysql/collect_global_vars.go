// SPDX-License-Identifier: GPL-3.0-or-later

package mysql

import (
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

func (m *MySQL) collectGlobalVariables(mx map[string]int64) error {
	// MariaDB: https://mariadb.com/kb/en/server-system-variables/
	// MySQL: https://dev.mysql.com/doc/refman/8.0/en/server-system-variable-reference.html
	m.Debugf("executing query: '%s'", queryGlobalVariables)

	var name string
	return m.collectQuery(queryGlobalVariables, func(column, value string) {
		switch column {
		case "Variable_name":
			name = value
		case "Value":
			switch name {
			case "disabled_storage_engines":
				mx[name] = parseInt(convertStorageEngineValue(value))
			case "log_bin":
				mx[name] = parseInt(convertBinlogValue(value))
			case "max_connections", "table_open_cache":
				mx[name] = parseInt(value)
			}
		}
	})
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
