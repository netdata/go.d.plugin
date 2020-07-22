package mysql

import (
	"database/sql"
	"fmt"
)

func (m *MySQL) collect() (map[string]int64, error) {
	collected := make(map[string]int64)

	if err := m.collectGlobalStatus(collected); err != nil {
		return nil, fmt.Errorf("error on collecting global status: %v", err)
	}
	m.checkGalera.Do(func() {
		if hasGaleraMetrics(collected) {
			m.addGaleraCharts()
		}
	})

	if err := m.collectGlobalVariables(collected); err != nil {
		return nil, fmt.Errorf("error on collecting global variables: %v", err)
	}

	if m.doSlaveStats {
		if err := m.collectSlaveStatus(collected); err != nil {
			m.Errorf("error on collecting slave status: %v", err)
			m.doSlaveStats = false
		}
	}

	if m.doUserStatistics {
		if err := m.collectUserStatistics(collected); err != nil {
			m.Errorf("error on collecting user statistics: %v", err)
			m.doUserStatistics = false
		}
	}

	calcThreadCacheMisses(collected)
	return collected, nil
}

func calcThreadCacheMisses(collected map[string]int64) {
	threads, ok1 := collected["threads_created"]
	cons, ok2 := collected["connections"]
	if !ok1 || !ok2 {
		return
	}
	if threads == 0 || cons == 0 {
		collected["thread_cache_misses"] = 0
	} else {
		collected["thread_cache_misses"] = int64(float64(threads) / float64(cons) * 10000)
	}
}

func hasGaleraMetrics(collected map[string]int64) bool {
	_, ok := collected["wsrep_received"]
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
