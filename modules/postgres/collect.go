// SPDX-License-Identifier: GPL-3.0-or-later

package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"time"
)

func (p *Postgres) collect() (map[string]int64, error) {
	if p.db == nil {
		if err := p.openConnection(); err != nil {
			return nil, err
		}
	}

	if p.serverVersion == 0 {
		ver, err := p.queryServerVersion()
		if err != nil {
			return nil, fmt.Errorf("querying server version error: %v", err)
		}
		p.serverVersion = ver
	}

	if now := time.Now(); now.Sub(p.relistDatabaseTime) > p.relistDatabasesEvery {
		p.relistDatabaseTime = now
		dbs, err := p.queryDatabaseList()
		if err != nil {
			return nil, fmt.Errorf("querying database list error: %v", err)
		}
		p.collectDatabaseList(dbs)
	}

	mx := make(map[string]int64)

	if err := p.collectDatabaseStats(mx); err != nil {
		return mx, fmt.Errorf("querying database stats error: %v", err)
	}

	// TODO: This view will only contain information on standby servers, since conflicts do not occur on primary servers.
	// see if possible to identify primary/standby and disable on primary if yes.
	if err := p.collectDatabaseConflicts(mx); err != nil {
		return mx, fmt.Errorf("querying database conflicts error: %v", err)
	}

	if err := p.collectCheckpoints(mx); err != nil {
		return mx, fmt.Errorf("querying database conflicts error: %v", err)
	}

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

func (p *Postgres) queryServerVersion() (int, error) {
	q := queryServerVersion()

	var v string
	ctx, cancel := context.WithTimeout(context.Background(), p.Timeout.Duration)
	defer cancel()
	if err := p.db.QueryRowContext(ctx, q).Scan(&v); err != nil {
		return 0, err
	}
	return strconv.Atoi(v)
}

//func (p *Postgres) queryIsSuperUser() (bool, error) {
//	q := queryIsSuperUser()
//
//	var v bool
//	ctx, cancel := context.WithTimeout(context.Background(), p.Timeout.Duration)
//	defer cancel()
//	if err := p.db.QueryRowContext(ctx, q).Scan(&v); err != nil {
//		return false, err
//	}
//	return v, nil
//}

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
