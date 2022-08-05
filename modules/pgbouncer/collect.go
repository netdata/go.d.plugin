package pgbouncer

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"time"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/stdlib"
)

func (p *PgBouncer) collect() (map[string]int64, error) {
	if p.db == nil {
		if err := p.openConnection(); err != nil {
			return nil, err
		}
	}

	mx := make(map[string]int64)

	// http://www.pgbouncer.org/usage.html

	if err := p.collectLists(mx); err != nil {
		return mx, err
	}
	if err := p.collectDatabases(mx); err != nil {
		return mx, err
	}
	if err := p.collectStats(mx); err != nil {
		return mx, err
	}
	if err := p.collectPools(mx); err != nil {
		return mx, err
	}

	return mx, nil
}

func (p *PgBouncer) collectLists(mx map[string]int64) error {
	q := "SHOW LISTS;"

	var name string
	return p.collectQuery(q, func(column, value string) {
		switch column {
		case "list":
			name = value
		case "items":
			mx[name] = safeParseInt(value)
		}
	})
}

func (p *PgBouncer) collectDatabases(mx map[string]int64) error {
	//q := "SHOW DATABASES;"
	return nil
}

func (p *PgBouncer) collectStats(mx map[string]int64) error {
	//q := "SHOW STATS;"
	return nil
}

func (p *PgBouncer) collectPools(mx map[string]int64) error {
	//q := "SHOW POOLS;"
	return nil
}

func (p *PgBouncer) openConnection() error {
	if p.connString == "" {
		cfg, err := pgx.ParseConfig(p.DSN)
		if err != nil {
			return err
		}
		cfg.PreferSimpleProtocol = true
		p.connString = stdlib.RegisterConnConfig(cfg)
	}

	db, err := sql.Open("pgx", p.connString)
	if err != nil {
		return fmt.Errorf("error on opening a connection with the PgBouncer database [%s]: %v", p.DSN, err)
	}

	db.SetMaxOpenConns(1)
	db.SetMaxIdleConns(1)
	db.SetConnMaxLifetime(10 * time.Minute)

	p.db = db

	return nil
}

func (p *PgBouncer) collectQuery(query string, assign func(column, value string)) error {
	ctx, cancel := context.WithTimeout(context.Background(), p.Timeout.Duration)
	defer cancel()
	rows, err := p.db.QueryContext(ctx, query)
	if err != nil {
		return err
	}
	defer func() { _ = rows.Close() }()

	columns, err := rows.Columns()
	if err != nil {
		return err
	}

	values := makeNullStrings(len(columns))
	for rows.Next() {
		if err := rows.Scan(values...); err != nil {
			return err
		}
		for i, v := range values {
			assign(columns[i], valueToString(v))
		}
	}
	return nil
}

func valueToString(value interface{}) string {
	v, ok := value.(*sql.NullString)
	if !ok || !v.Valid {
		return ""
	}
	return v.String
}

func makeNullStrings(size int) []interface{} {
	vs := make([]interface{}, size)
	for i := range vs {
		vs[i] = &sql.NullString{}
	}
	return vs
}

func safeParseInt(s string) int64 {
	v, _ := strconv.ParseInt(s, 10, 64)
	return v
}
