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

	now := time.Now()

	if now.Sub(p.recheckSettingsTime) > p.recheckSettingsEvery {
		p.recheckSettingsTime = now
		maxConn, err := p.querySettingsMaxConnections()
		if err != nil {
			return nil, fmt.Errorf("querying settings max connections error: %v", err)
		}
		p.maxConnections = maxConn
	}

	if now.Sub(p.relistDatabaseTime) > p.relistDatabaseEvery {
		p.relistDatabaseTime = now
		dbs, err := p.queryDatabaseList()
		if err != nil {
			return nil, fmt.Errorf("querying database list error: %v", err)
		}
		p.collectDatabaseList(dbs)
	}

	if now.Sub(p.relistReplStandbyTime) > p.relistReplStandbyEvery {
		p.relistReplStandbyTime = now
		apps, err := p.queryReplicationStandbyAppList()
		if err != nil {
			return nil, fmt.Errorf("querying replication standby app list error: %v", err)
		}
		p.collectReplicationStandbyAppList(apps)
	}

	if now.Sub(p.relistReplSlotTime) > p.relistReplSlotEvery {
		p.relistReplSlotTime = now
		slots, err := p.queryReplicationSlotList()
		if err != nil {
			return nil, fmt.Errorf("querying replication slot list error: %v", err)
		}
		p.collectReplicationSlotList(slots)
	}

	mx := make(map[string]int64)

	if err := p.collectConnection(mx); err != nil {
		return mx, fmt.Errorf("querying server connections error: %v", err)
	}

	if err := p.collectCheckpoints(mx); err != nil {
		return mx, fmt.Errorf("querying database conflicts error: %v", err)
	}

	if err := p.collectUptime(mx); err != nil {
		return mx, fmt.Errorf("querying server uptime error: %v", err)
	}

	if err := p.collectTXIDWraparound(mx); err != nil {
		return mx, fmt.Errorf("querying txid wraparound error: %v", err)
	}

	if err := p.collectWALWrites(mx); err != nil {
		return mx, fmt.Errorf("querying wal writes error: %v", err)
	}

	// TODO: superuser only
	if err := p.collectWALFiles(mx); err != nil {
		return mx, fmt.Errorf("querying wal files error: %v", err)
	}

	// TODO: superuser only
	if err := p.collectWALArchiveFiles(mx); err != nil {
		return mx, fmt.Errorf("querying wal archive files error: %v", err)
	}

	if err := p.collectCatalog(mx); err != nil {
		return mx, fmt.Errorf("querying catalog relations error: %v", err)
	}

	if err := p.collectAutovacuumWorkers(mx); err != nil {
		return mx, fmt.Errorf("querying autovacuum workers error: %v", err)
	}

	if len(p.replStandbyApps) > 0 {
		if err := p.collectReplicationStandbyAppWALDelta(mx); err != nil {
			return mx, fmt.Errorf("querying replication standby app wal delta error: %v", err)
		}
		if p.serverVersion >= 100000 {
			if err := p.collectReplicationStandbyAppWALLag(mx); err != nil {
				return mx, fmt.Errorf("querying replication standby app wal lag error: %v", err)
			}
		}
	}

	// TODO: superuser only
	if len(p.replSlots) > 0 {
		if p.serverVersion >= 100000 {
			if err := p.collectReplicationSlotFiles(mx); err != nil {
				return mx, fmt.Errorf("querying replication slot files error: %v", err)
			}
		}
	}

	if len(p.databases) > 0 {
		if err := p.collectDatabaseStats(mx); err != nil {
			return mx, fmt.Errorf("querying database stats error: %v", err)
		}

		if err := p.collectDatabaseConflicts(mx); err != nil {
			return mx, fmt.Errorf("querying database conflicts error: %v", err)
		}

		if err := p.collectDatabaseLocks(mx); err != nil {
			return mx, fmt.Errorf("querying database locks error: %v", err)
		}
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

func (p *Postgres) querySettingsMaxConnections() (int64, error) {
	q := querySettingsMaxConnections()

	var s string
	ctx, cancel := context.WithTimeout(context.Background(), p.Timeout.Duration)
	defer cancel()
	if err := p.db.QueryRowContext(ctx, q).Scan(&s); err != nil {
		return 0, err
	}
	return strconv.ParseInt(s, 10, 64)
}

func (p *Postgres) queryServerVersion() (int, error) {
	q := queryServerVersion()

	var s string
	ctx, cancel := context.WithTimeout(context.Background(), p.Timeout.Duration)
	defer cancel()
	if err := p.db.QueryRowContext(ctx, q).Scan(&s); err != nil {
		return 0, err
	}
	return strconv.Atoi(s)
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

func collectRows(rows *sql.Rows, assign func(column, value string)) error {
	if assign == nil {
		return nil
	}
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

func calcPercentage(value, total int64) int64 {
	if total == 0 {
		return 0
	}
	return value * 100 / total
}
