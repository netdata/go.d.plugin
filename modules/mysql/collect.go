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
	return collected, nil
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
