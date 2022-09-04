package proxysql

import (
	"database/sql"
	"fmt"
	"time"
)

func (p *ProxySQL) collect() (map[string]int64, error) {
	if p.db == nil {
		if err := p.openConnection(); err != nil {
			return nil, err
		}
	}

	collected := make(map[string]int64)

	if err := p.collectGlobalVars(collected); err != nil {
		return nil, fmt.Errorf("error on collecting global vars: %v", err)
	}

	if err := p.collectMemoryMetrics(collected); err != nil {
		return nil, fmt.Errorf("error on collecting memory metrics: %v", err)
	}
	if err := p.collectMysqlCommandCounters(collected); err != nil {
		return nil, fmt.Errorf("error on collecting mysql command counters: %v", err)
	}
	if err := p.collectMysqlGlobalStatus(collected); err != nil {
		return nil, fmt.Errorf("error on collecting mysql global status: %v", err)
	}
	if err := p.collectMysqlUsers(collected); err != nil {
		return nil, fmt.Errorf("error on collecting mysql users: %v", err)
	}

	return collected, nil
}

func (p *ProxySQL) openConnection() error {
	db, err := sql.Open("mysql", p.DSN)
	if err != nil {
		return fmt.Errorf("error on opening a connection with the proxysql instance [%s]: %v", p.DSN, err)
	}

	db.SetConnMaxLifetime(10 * time.Minute)

	if err := db.Ping(); err != nil {
		_ = db.Close()
		return fmt.Errorf("error on pinging the proxysql instance [%s]: %v", p.DSN, err)
	}

	p.db = db
	return nil
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
