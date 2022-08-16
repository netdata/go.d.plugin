// SPDX-License-Identifier: GPL-3.0-or-later

package mysql

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"time"

	"github.com/blang/semver/v4"
)

func (m *MySQL) collect() (map[string]int64, error) {
	if m.db == nil {
		if err := m.openConnection(); err != nil {
			return nil, err
		}
	}
	if m.version == nil {
		ver, isMariaDB, err := m.collectVersion()
		if err != nil {
			return nil, err
		}

		m.version = ver
		m.isMariaDB = isMariaDB
		// https://mariadb.com/kb/en/user-statistics/
		minVer := semver.Version{Major: 10, Minor: 1, Patch: 1}
		m.doUserStatistics = m.isMariaDB && m.version.GTE(minVer)
	}

	collected := make(map[string]int64)

	// TODO: do we really need to collect global vars on every iteration?
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

	if hasMyISAMStorageEngine(collected) || m.isMariaDB {
		m.addMyISAMOnce.Do(m.addMyISAMCharts)
	}
	if hasBinlogEnabled(collected) {
		m.addBinlogOnce.Do(m.addBinlogCharts)
	}

	if m.doSlaveStatus {
		if err := m.collectSlaveStatus(collected); err != nil {
			m.Errorf("error on collecting slave status: %v", err)
			// TODO: shouldn't disable on any error
			m.doSlaveStatus = false
		}
	}

	if m.doUserStatistics {
		if err := m.collectUserStatistics(collected); err != nil {
			m.Errorf("error on collecting user statistics: %v", err)
			// TODO: shouldn't disable on any error
			m.doUserStatistics = false
		}
	}

	if err := m.collectProcessListStatistics(collected); err != nil {
		m.Errorf("error on collecting process list statistics: %v", err)
	}

	calcThreadCacheMisses(collected)
	return collected, nil
}

func (m *MySQL) openConnection() error {
	db, err := sql.Open("mysql", m.DSN)
	if err != nil {
		return fmt.Errorf("error on opening a connection with the mysql database [%s]: %v", m.DSN, err)
	}

	db.SetConnMaxLifetime(10 * time.Minute)

	ctx, cancel := context.WithTimeout(context.Background(), m.Timeout.Duration)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		_ = db.Close()
		return fmt.Errorf("error on pinging the mysql database [%s]: %v", m.DSN, err)
	}

	m.db = db
	return nil
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

func hasMyISAMStorageEngine(collected map[string]int64) bool {
	return collected["disabled_storage_engines"] == 0
}

func hasBinlogEnabled(collected map[string]int64) bool {
	return collected["log_bin"] == 1
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

// ----

func (m *MySQL) collectQuery(query string, assign func(column, value string, lineEnd bool)) (duration int64, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), m.Timeout.Duration)
	defer cancel()

	s := time.Now()
	rows, err := m.db.QueryContext(ctx, query)
	if err != nil {
		return 0, err
	}
	duration = time.Since(s).Milliseconds()
	defer func() { _ = rows.Close() }()

	columns, err := rows.Columns()
	if err != nil {
		return duration, err
	}

	vs := makeValues(len(columns))
	for rows.Next() {
		if err := rows.Scan(vs...); err != nil {
			return duration, err
		}
		for i, l := 0, len(vs); i < l; i++ {
			assign(columns[i], valueToString(vs[i]), i == l-1)
		}
	}
	return duration, rows.Err()
}

func parseInt(s string) int64 {
	v, _ := strconv.ParseInt(s, 10, 64)
	return v
}

func makeValues(size int) []any {
	vs := make([]any, size)
	for i := range vs {
		vs[i] = &sql.NullString{}
	}
	return vs
}

func valueToString(value any) string {
	v, ok := value.(*sql.NullString)
	if !ok || !v.Valid {
		return ""
	}
	return v.String
}
