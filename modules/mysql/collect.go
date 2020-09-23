package mysql

import (
	"database/sql"
	"fmt"
	"strings"
)

func (m *MySQL) collect() (map[string]int64, error) {
	collected := make(map[string]int64)

	if err := m.collectGlobalStatus(collected); err != nil {
		return nil, fmt.Errorf("error on collecting global status: %v", err)
	}
	if hasInnodbDeadlocks(collected) {
		m.addInnodbDeadlocksOnce.Do(m.addInnodbDeadlocksChart)
	}
	if hasQCacheMetrics(collected) {
		m.addQCacheOnce.Do(m.addQCacheCharts)
	}
	if hasGaleraMetrics(collected) {
		m.addGaleraOnce.Do(m.addGaleraCharts)
	}

	if err := m.collectGlobalVariables(collected); err != nil {
		return nil, fmt.Errorf("error on collecting global variables: %v", err)
	}

	if m.doSlaveStatus {
		// TODO: shouldn't disable on any error
		if err := m.collectSlaveStatus(collected); err != nil {
			m.Errorf("error on collecting slave status: %v", err)
			m.doSlaveStatus = false
		}
	}

	if m.doUserStatistics {
		// TODO: shouldn't disable on any error
		if err := m.collectUserStatistics(collected); err != nil {
			m.Errorf("error on collecting user statistics: %v", err)
			m.doUserStatistics = false
		}
	}

	calcThreadCacheMisses(collected)
	return collected, nil
}

func (m MySQL) isMariaDB() bool {
	return strings.Contains(strings.ToLower(m.version), "mariadb")
}

func calcThreadCacheMisses(collected map[string]int64) {
	threads, cons := collected["threads_created"], collected["connections"]
	if threads == 0 || cons == 0 {
		collected["thread_cache_misses"] = 0
	} else {
		collected["thread_cache_misses"] = int64(float64(threads) / float64(cons) * 10000)
	}
}

func hasInnodbDeadlocks(collected map[string]int64) bool {
	_, ok := collected["innodb_deadlocks"]
	return ok
}

func hasGaleraMetrics(collected map[string]int64) bool {
	_, ok := collected["wsrep_received"]
	return ok
}

func hasQCacheMetrics(collected map[string]int64) bool {
	_, ok := collected["qcache_hits"]
	return ok
}

func rowsAsMap(rows *sql.Rows) (map[string]string, error) {
	set := make(map[string]string)
	for rows.Next() {
		var name, value string
		if err := rows.Scan(&name, &value); err != nil {
			return nil, err
		}
		set[name] = value
	}
	return set, rows.Err()
}

func rowAsMap(columns []string, values []interface{}) map[string]string {
	set := make(map[string]string, len(columns))
	for i, name := range columns {
		if v, ok := values[i].(*sql.NullString); ok && v.Valid {
			set[name] = v.String
		}
	}
	return set
}

func nullStringsFromColumns(columns []string) []interface{} {
	values := make([]interface{}, len(columns))
	for i := range values {
		values[i] = &sql.NullString{}
	}
	return values
}
