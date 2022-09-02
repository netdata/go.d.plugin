// SPDX-License-Identifier: GPL-3.0-or-later

package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"time"
)

const (
	pgVersion94 = 9_04_00
	pgVersion10 = 10_00_00
	pgVersion11 = 11_00_00
)

func (p *Postgres) collect() (map[string]int64, error) {
	if p.db == nil {
		if err := p.openConnection(); err != nil {
			return nil, err
		}
	}

	if p.pgVersion == 0 {
		ver, err := p.queryServerVersion()
		if err != nil {
			return nil, fmt.Errorf("querying server version error: %v", err)
		}
		p.pgVersion = ver
		p.Debugf("connected to PostgreSQL v%s", p.pgVersion)
	}

	if p.superUser == nil {
		v, err := p.queryIsSuperUser()
		if err != nil {
			return nil, fmt.Errorf("querying is super user error: %v", err)
		}
		p.superUser = &v
		p.Debugf("connected as super user: %s", *p.superUser)
	}

	now := time.Now()

	if now.Sub(p.recheckSettingsTime) > p.recheckSettingsEvery {
		p.recheckSettingsTime = now
		maxConn, err := p.querySettingsMaxConnections()
		if err != nil {
			return nil, fmt.Errorf("querying settings max connections error: %v", err)
		}
		p.metrics.maxConnections = maxConn
	}

	if err := p.queryGlobalMetrics(); err != nil {
		return nil, err
	}
	if err := p.doQueryReplicationMetrics(); err != nil {
		return nil, err
	}
	if err := p.doQueryDatabasesMetrics(); err != nil {
		return nil, err
	}

	mx := make(map[string]int64)

	return mx, nil
}

func (p *Postgres) openConnection() error {
	db, err := sql.Open("pgx", p.DSN)
	if err != nil {
		return fmt.Errorf("error on opening a connection with the Postgres database [%s]: %v", p.DSN, err)
	}

	db.SetMaxOpenConns(1)
	db.SetMaxIdleConns(1)
	db.SetConnMaxLifetime(10 * time.Minute)

	ctx, cancel := context.WithTimeout(context.Background(), p.Timeout.Duration)
	defer cancel()
	if err := db.PingContext(ctx); err != nil {
		_ = db.Close()
		return fmt.Errorf("error on pinging the Postgres database [%s]: %v", p.DSN, err)
	}
	p.db = db

	return nil
}

func (p *Postgres) querySettingsMaxConnections() (int64, error) {
	q := querySettingsMaxConnections()

	var s string
	if err := p.doQueryRow(q).Scan(&s); err != nil {
		return 0, err
	}

	return strconv.ParseInt(s, 10, 64)
}

func (p *Postgres) queryServerVersion() (int, error) {
	q := queryServerVersion()

	var s string
	if err := p.doQueryRow(q).Scan(&s); err != nil {
		return 0, err
	}

	return strconv.Atoi(s)
}

func (p *Postgres) queryIsSuperUser() (bool, error) {
	q := queryIsSuperUser()

	var v bool
	if err := p.doQueryRow(q).Scan(&v); err != nil {
		return false, err
	}

	return v, nil
}

func (p *Postgres) isSuperUser() bool { return p.superUser != nil && *p.superUser }

func (p *Postgres) doQueryRow(query string) *sql.Row {
	ctx, cancel := context.WithTimeout(context.Background(), p.Timeout.Duration)
	defer cancel()

	return p.db.QueryRowContext(ctx, query)
}

func (p *Postgres) doQueryRows(query string, assign func(column, value string, rowEnd bool)) error {
	ctx, cancel := context.WithTimeout(context.Background(), p.Timeout.Duration)
	defer cancel()

	rows, err := p.db.QueryContext(ctx, query)
	if err != nil {
		return err
	}
	defer func() { _ = rows.Close() }()

	return readRows(rows, assign)
}

func (p *Postgres) getDBMetrics(name string) *dbMetrics {
	db, ok := p.metrics.dbs[name]
	if !ok {
		db = &dbMetrics{name: name}
		p.metrics.dbs[name] = db
	}
	return db
}

func (p *Postgres) getReplAppMetrics(name string) *replStandbyAppMetrics {
	app, ok := p.metrics.replApps[name]
	if !ok {
		app = &replStandbyAppMetrics{name: name}
		p.metrics.replApps[name] = app
	}
	return app
}

func (p *Postgres) getReplSlotMetrics(name string) *replSlotMetrics {
	slot, ok := p.metrics.replSlots[name]
	if !ok {
		slot = &replSlotMetrics{name: name}
		p.metrics.replSlots[name] = slot
	}
	return slot
}

func readRows(rows *sql.Rows, assign func(column, value string, rowEnd bool)) error {
	if assign == nil {
		return nil
	}

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
	return nil
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

func parseInt(s string) int64 {
	v, _ := strconv.ParseInt(s, 10, 64)
	return v
}

func calcPercentage(value, total int64) int64 {
	if total == 0 {
		return 0
	}
	return value * 100 / total
}
