// SPDX-License-Identifier: GPL-3.0-or-later

package postgres

import (
	"context"
)

func (p *Postgres) collectDatabaseStats(mx map[string]int64) error {
	q := queryDatabaseStats(p.databases)

	ctx, cancel := context.WithTimeout(context.Background(), p.Timeout.Duration)
	defer cancel()
	rows, err := p.db.QueryContext(ctx, q)
	if err != nil {
		return err
	}
	defer func() { _ = rows.Close() }()

	var db string
	var currConn int64
	return collectRows(rows, func(column, value string) {
		switch column {
		case "datname":
			db = value
		case "numbackends":
			currConn = safeParseInt(value)
			mx["db_"+db+"_"+column] = currConn
		case "datconnlimit":
			limit := safeParseInt(value)
			if limit <= 0 {
				limit = p.maxConnections
			}
			mx["db_"+db+"_numbackends_utilization"] = calcPercentage(currConn, limit)
		default:
			mx["db_"+db+"_"+column] = safeParseInt(value)
		}
	})
}

func (p *Postgres) collectDatabaseConflicts(mx map[string]int64) error {
	q := queryDatabaseConflicts(p.databases)

	ctx, cancel := context.WithTimeout(context.Background(), p.Timeout.Duration)
	defer cancel()
	rows, err := p.db.QueryContext(ctx, q)
	if err != nil {
		return err
	}
	defer func() { _ = rows.Close() }()

	var db string
	return collectRows(rows, func(column, value string) {
		switch column {
		case "datname":
			db = value
		default:
			mx["db_"+db+"_"+column] = safeParseInt(value)
		}
	})
}

func (p *Postgres) collectDatabaseLocks(mx map[string]int64) error {
	q := queryDatabaseLocks(p.databases)

	ctx, cancel := context.WithTimeout(context.Background(), p.Timeout.Duration)
	defer cancel()
	rows, err := p.db.QueryContext(ctx, q)
	if err != nil {
		return err
	}
	defer func() { _ = rows.Close() }()

	// https://github.com/postgres/postgres/blob/7c34555f8c39eeefcc45b3c3f027d7a063d738fc/src/include/storage/lockdefs.h#L36-L45
	// https://www.postgresql.org/docs/7.2/locking-tables.html
	var lockModes = []string{
		"AccessShareLock",
		"RowShareLock",
		"RowExclusiveLock",
		"ShareUpdateExclusiveLock",
		"ShareLock",
		"ShareRowExclusiveLock",
		"ExclusiveLock",
		"AccessExclusiveLock",
	}

	for _, db := range p.databases {
		for _, mode := range lockModes {
			mx["db_"+db+"_lock_mode_"+mode+"_held"] = 0
			mx["db_"+db+"_lock_mode_"+mode+"_awaited"] = 0
		}
	}

	var db, mode string
	var granted bool
	return collectRows(rows, func(column, value string) {
		switch column {
		case "datname":
			db = value
		case "mode":
			mode = value
		case "granted":
			granted = value == "true" || value == "t"
		case "locks_count":
			if granted {
				mx["db_"+db+"_lock_mode_"+mode+"_held"] = safeParseInt(value)
			} else {
				mx["db_"+db+"_lock_mode_"+mode+"_awaited"] = safeParseInt(value)
			}
		}
	})
}

func (p *Postgres) queryDatabaseList() ([]string, error) {
	q := queryDatabaseList()

	ctx, cancel := context.WithTimeout(context.Background(), p.Timeout.Duration)
	defer cancel()
	rows, err := p.db.QueryContext(ctx, q)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	var dbs []string
	var db string
	for rows.Next() {
		if err := rows.Scan(&db); err != nil {
			return nil, err
		}
		dbs = append(dbs, db)
	}

	return dbs, nil
}

func (p *Postgres) collectDatabaseList(dbs []string) {
	if len(dbs) == 0 {
		return
	}

	collected := make(map[string]bool)
	for _, db := range p.databases {
		collected[db] = true
	}
	p.databases = dbs

	seen := make(map[string]bool)
	for _, db := range dbs {
		seen[db] = true
		if !collected[db] {
			collected[db] = true
			p.addNewDatabaseCharts(db)
		}
	}
	for db := range collected {
		if !seen[db] {
			p.removeDatabaseCharts(db)
		}
	}
}
