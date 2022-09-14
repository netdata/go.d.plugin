package proxysql

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"time"
)

const (
	queryVersion                    = "select version();"
	queryStatsMySQLGlobal           = "SELECT * FROM stats_mysql_global;"
	queryStatsMySQLMemoryMetrics    = "SELECT * FROM stats_memory_metrics;"
	queryStatsMySQLCommandsCounters = "SELECT * FROM stats_mysql_commands_counters;"
	queryStatsMySQLUsers            = "SELECT * FROM stats_mysql_users;"
)

func (p *ProxySQL) collect() (map[string]int64, error) {
	if p.db == nil {
		if err := p.openConnection(); err != nil {
			return nil, err
		}
	}

	mx := make(map[string]int64)

	if err := p.collectStatsMySQLGlobal(mx); err != nil {
		return nil, fmt.Errorf("error on collecting mysql global status: %v", err)
	}
	if err := p.collectStatsMySQLMemoryMetrics(mx); err != nil {
		return nil, fmt.Errorf("error on collecting memory metrics: %v", err)
	}
	if err := p.collectMySQLCommandsCounters(mx); err != nil {
		return nil, fmt.Errorf("error on collecting mysql command counters: %v", err)
	}
	if err := p.collectStatsMySQLUsers(mx); err != nil {
		return nil, fmt.Errorf("error on collecting mysql users: %v", err)
	}

	return mx, nil
}

func (p *ProxySQL) collectStatsMySQLGlobal(mx map[string]int64) error {
	// https://proxysql.com/documentation/stats-statistics/#stats_mysql_global
	q := queryStatsMySQLGlobal
	p.Debugf("executing query: '%s'", q)

	var name string
	return p.doQuery(q, func(column, value string, rowEnd bool) {
		switch column {
		case "Variable_Name":
			name = value
		case "Variable_Value":
			mx[name] = parseInt(value)
		}
	})
}

func (p *ProxySQL) collectStatsMySQLMemoryMetrics(mx map[string]int64) error {
	// https://proxysql.com/documentation/stats-statistics/#stats_mysql_memory_metrics
	q := queryStatsMySQLMemoryMetrics
	p.Debugf("executing query: '%s'", q)

	var name string
	return p.doQuery(q, func(column, value string, rowEnd bool) {
		switch column {
		case "Variable_Name":
			name = value
		case "Variable_Value":
			mx[name] = parseInt(value)
		}
	})
}

func (p *ProxySQL) collectMySQLCommandsCounters(mx map[string]int64) error {
	// https://proxysql.com/documentation/stats-statistics/#stats_mysql_commands_counters
	q := queryStatsMySQLCommandsCounters
	p.Debugf("executing query: '%s'", q)

	var command string
	return p.doQuery(q, func(column, value string, rowEnd bool) {
		switch column {
		case "Command":
			command = value
			if !p.commands[command] {
				p.commands[command] = true
				p.addMySQLCommandCountersCharts(command)
			}
		default:
			mx["mysql_command_"+command+"_"+column] = parseInt(value)
		}
	})
}

func (p *ProxySQL) collectStatsMySQLUsers(mx map[string]int64) error {
	// https://proxysql.com/documentation/stats-statistics/#stats_mysql_users
	q := queryStatsMySQLUsers
	p.Debugf("executing query: '%s'", q)

	var user string
	return p.doQuery(q, func(column, value string, rowEnd bool) {
		switch column {
		case "username":
			user = value
			if !p.users[user] {
				p.users[user] = true
				p.addMysqlUsersCharts(user)
			}
		default:
			mx["mysql_user_"+user+"_"+column] = parseInt(value)
		}

	})
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

func (p *ProxySQL) doQueryRow(query string, v any) error {
	ctx, cancel := context.WithTimeout(context.Background(), p.Timeout.Duration)
	defer cancel()

	return p.db.QueryRowContext(ctx, query).Scan(v)
}

func (p *ProxySQL) doQuery(query string, assign func(column, value string, rowEnd bool)) error {
	ctx, cancel := context.WithTimeout(context.Background(), p.Timeout.Duration)
	defer cancel()

	rows, err := p.db.QueryContext(ctx, query)
	if err != nil {
		return err
	}
	defer func() { _ = rows.Close() }()

	return readRows(rows, assign)
}

func readRows(rows *sql.Rows, assign func(column, value string, rowEnd bool)) error {
	columns, err := rows.Columns()
	if err != nil {
		return err
	}

	values := makeValues(len(columns))

	for rows.Next() {
		if err := rows.Scan(values...); err != nil {
			return err
		}
		for i, l := 0, len(values); i < l; i++ {
			assign(columns[i], valueToString(values[i]), i == l-1)
		}
	}
	return rows.Err()
}

func valueToString(value any) string {
	v, ok := value.(*sql.NullString)
	if !ok || !v.Valid {
		return ""
	}
	return v.String
}

func makeValues(size int) []any {
	vs := make([]any, size)
	for i := range vs {
		vs[i] = &sql.NullString{}
	}
	return vs
}

func parseInt(value string) int64 {
	v, _ := strconv.ParseInt(value, 10, 64)
	return v
}
