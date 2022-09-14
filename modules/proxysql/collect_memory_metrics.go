package proxysql

import (
	"strconv"
	"strings"
)

const queryMemoryMetrics = "SELECT Variable_Name, Variable_Value FROM stats_memory_metrics"

var memoryMetrics = []string{
	"Auth_memory",
	"SQLite3_memory_bytes",
	"jemalloc_active",
	"jemalloc_allocated",
	"jemalloc_mapped",
	"jemalloc_metadata",
	"jemalloc_resident",
	"jemalloc_retained",
	"mysql_firewall_rules_config",
	"mysql_firewall_rules_table",
	"mysql_firewall_users_config",
	"mysql_firewall_users_table",
	"mysql_query_rules_memory",
	"query_digest_memory",
	"stack_memory_admin_threads",
	"stack_memory_cluster_threads",
	"stack_memory_mysql_threads",
}

func (p *ProxySQL) collectMemoryMetrics(collected map[string]int64) error {
	// https://proxysql.com/documentation/stats-statistics/#stats_mysql_memory_metrics
	p.Debugf("executing query: '%s'", queryMemoryMetrics)

	rows, err := p.db.Query(queryMemoryMetrics)
	if err != nil {
		return err
	}
	defer rows.Close()

	set, err := rowsAsMap(rows)
	if err != nil {
		return err
	}

	for _, name := range memoryMetrics {
		strValue, ok := set[name]
		if !ok {
			continue
		}
		value, err := parseMemoryMetrics(strValue)
		if err != nil {
			continue
		}
		collected[strings.ToLower(name)] = value
	}
	return nil
}

func parseMemoryMetrics(value string) (int64, error) {
	return strconv.ParseInt(value, 10, 64)
}
