package proxysql

import (
	"strconv"
	"strings"
)

const queryGlobalVars = "SELECT Variable_Name, Variable_Value FROM global_variables"

var globalVarsMetrics = []string{
	"mysql-max_allowed_packet",
}

func (p *ProxySQL) collectGlobalVars(collected map[string]int64) error {
	// https://proxysql.com/documentation/stats-statistics/#stats_mysql_global
	p.Debugf("executing query: '%s'", queryGlobalVars)

	rows, err := p.db.Query(queryGlobalVars)
	if err != nil {
		return err
	}
	defer rows.Close()

	set, err := rowsAsMap(rows)
	if err != nil {
		return err
	}

	for _, name := range globalVarsMetrics {
		strValue, ok := set[name]
		if !ok {
			continue
		}
		value, err := parseGlobalVarsValue(strValue)
		if err != nil {
			continue
		}
		collected[parseGlobalVarsKey(name)] = value
	}
	return nil
}

func parseGlobalVarsValue(value string) (int64, error) {
	return strconv.ParseInt(value, 10, 64)
}

func parseGlobalVarsKey(key string) string {
	key = strings.ToLower(key)
	key = strings.Replace(key, "-", "_", -1)
	return key
}
