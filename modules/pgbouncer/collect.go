package pgbouncer

import (
	"context"
	"database/sql"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/blang/semver/v4"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/stdlib"
)

func (p *PgBouncer) collect() (map[string]int64, error) {
	if p.db == nil {
		if err := p.openConnection(); err != nil {
			return nil, err
		}
	}
	if p.version == nil {
		ver, err := p.queryVersion()
		if err != nil {
			return nil, err
		}
		p.version = ver
		p.Debugf("connected to PgBouncer v%s", p.version)
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
	p.Debugf("executing query: %v", q)

	var name string
	return p.collectQuery(q, func(column, value string) {
		switch column {
		case "list":
			name = value
		case "items":
			if v, err := strconv.ParseInt(value, 10, 64); err == nil {
				mx[name] = v
			}
		}
	})
}

func (p *PgBouncer) collectDatabases(mx map[string]int64) error {
	q := "SHOW DATABASES;"
	p.Debugf("executing query: %v", q)

	var db string
	return p.collectQuery(q, func(column, value string) {
		switch column {
		case "database":
			db = value
		default:
			if v, err := strconv.ParseInt(value, 10, 64); err == nil {
				mx["db_"+db+"_"+column] = v
			}
		}
	})
}

func (p *PgBouncer) collectStats(mx map[string]int64) error {
	q := "SHOW STATS;"
	p.Debugf("executing query: %v", q)

	var db string
	return p.collectQuery(q, func(column, value string) {
		switch column {
		case "database":
			db = value
		default:
			if v, err := strconv.ParseInt(value, 10, 64); err == nil {
				mx["db_"+db+"_"+column] = v
			}
		}
	})
}

func (p *PgBouncer) collectPools(mx map[string]int64) error {
	q := "SHOW POOLS;"
	p.Debugf("executing query: %v", q)

	var db, user string
	return p.collectQuery(q, func(column, value string) {
		switch column {
		case "database":
			db = value
		case "user":
			user = value
		default:
			if v, err := strconv.ParseInt(value, 10, 64); err == nil {
				mx["db_"+db+"_user_"+user+"_"+column] = v
			}
		}
	})
}

var reVersion = regexp.MustCompile(`\d+\.\d+\.\d+`)

func (p *PgBouncer) queryVersion() (*semver.Version, error) {
	q := "SHOW VERSION;"
	p.Debugf("executing query: %v", q)

	var resp string
	ctx, cancel := context.WithTimeout(context.Background(), p.Timeout.Duration)
	defer cancel()
	if err := p.db.QueryRowContext(ctx, q).Scan(&resp); err != nil {
		return nil, err
	}

	if !strings.Contains(resp, "PgBouncer") {
		return nil, fmt.Errorf("not PgBouncer instance: version response: %s", resp)
	}

	ver := reVersion.FindString(resp)
	if ver == "" {
		return nil, fmt.Errorf("couldn't parse version string '%s' (expected pattern '%s')", resp, reVersion)
	}

	v, err := semver.New(ver)
	if err != nil {
		return nil, fmt.Errorf("couldn't parse version string '%s': %v", ver, err)
	}

	return v, nil
}

func (p *PgBouncer) openConnection() error {
	cfg, err := pgx.ParseConfig(p.DSN)
	if err != nil {
		return err
	}
	cfg.PreferSimpleProtocol = true

	db, err := sql.Open("pgx", stdlib.RegisterConnConfig(cfg))
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
