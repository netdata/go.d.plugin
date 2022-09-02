// SPDX-License-Identifier: GPL-3.0-or-later

package postgres

import (
	"fmt"
)

func (p *Postgres) doQueryDatabasesMetrics() error {
	if err := p.doQueryDatabaseStats(); err != nil {
		return fmt.Errorf("querying database stats error: %v", err)
	}
	if err := p.doQueryDatabaseConflicts(); err != nil {
		return fmt.Errorf("querying database conflicts error: %v", err)
	}
	if err := p.doQueryDatabaseLocks(); err != nil {
		return fmt.Errorf("querying database locks error: %v", err)
	}
	return nil
}

func (p *Postgres) doQueryDatabaseStats() error {
	q := queryDatabaseStats(p.databases)

	var db string
	return p.doQueryRows(q, func(column, value string, _ bool) {
		switch column {
		case "datname":
			db = value
			p.getDBMetrics(db).updated = true
		case "numbackends":
			p.getDBMetrics(db).numBackends = parseInt(value)
		case "datconnlimit":
			p.getDBMetrics(db).datConnLimit = parseInt(value)
		case "xact_commit":
			p.getDBMetrics(db).xactCommit = parseInt(value)
		case "xact_rollback":
			p.getDBMetrics(db).xactRollback = parseInt(value)
		case "blks_read":
			p.getDBMetrics(db).blksRead = parseInt(value)
		case "blks_hit":
			p.getDBMetrics(db).blksHit = parseInt(value)
		case "tup_returned":
			p.getDBMetrics(db).tupReturned = parseInt(value)
		case "tup_fetched":
			p.getDBMetrics(db).tupFetched = parseInt(value)
		case "tup_inserted":
			p.getDBMetrics(db).tupInserted = parseInt(value)
		case "tup_updated":
			p.getDBMetrics(db).tupUpdated = parseInt(value)
		case "tup_deleted":
			p.getDBMetrics(db).tupDeleted = parseInt(value)
		case "conflicts":
			p.getDBMetrics(db).conflicts = parseInt(value)
		case "size":
			p.getDBMetrics(db).size = parseInt(value)
		case "temp_files":
			p.getDBMetrics(db).tempFiles = parseInt(value)
		case "temp_bytes":
			p.getDBMetrics(db).tempBytes = parseInt(value)
		case "deadlocks":
			p.getDBMetrics(db).deadlocks = parseInt(value)
		}
	})
}

func (p *Postgres) doQueryDatabaseConflicts() error {
	q := queryDatabaseConflicts(p.databases)

	var db string
	return p.doQueryRows(q, func(column, value string, _ bool) {
		switch column {
		case "datname":
			db = value
			p.getDBMetrics(db).updated = true
		case "confl_tablespace":
			p.getDBMetrics(db).conflTablespace = parseInt(value)
		case "confl_lock":
			p.getDBMetrics(db).conflLock = parseInt(value)
		case "confl_snapshot":
			p.getDBMetrics(db).conflSnapshot = parseInt(value)
		case "confl_bufferpin":
			p.getDBMetrics(db).conflBufferpin = parseInt(value)
		case "confl_deadlock":
			p.getDBMetrics(db).conflDeadlock = parseInt(value)
		}
	})
}

func (p *Postgres) doQueryDatabaseLocks() error {
	q := queryDatabaseLocks(p.databases)

	var db, mode, granted string
	var locks struct{ held, awaited int64 }
	return p.doQueryRows(q, func(column, value string, rowEnd bool) {
		switch column {
		case "datname":
			db = value
			p.getDBMetrics(db).updated = true
		case "mode":
			mode = value
		case "granted":
			granted = value
		case "locks_count":
			locks.held, locks.awaited = 0, 0
			if granted == "true" || granted == "t" {
				locks.held = parseInt(value)
			} else {
				locks.awaited = parseInt(value)
			}
		}
		if !rowEnd {
			return
		}
		// https://github.com/postgres/postgres/blob/7c34555f8c39eeefcc45b3c3f027d7a063d738fc/src/include/storage/lockdefs.h#L36-L45
		// https://www.postgresql.org/docs/7.2/locking-tables.html
		switch mode {
		case "AccessShareLock":
			p.getDBMetrics(db).accessShareLockHeld = locks.held
			p.getDBMetrics(db).accessShareLockAwaited = locks.awaited
		case "RowShareLock":
			p.getDBMetrics(db).rowShareLockHeld = locks.held
			p.getDBMetrics(db).rowShareLockAwaited = locks.awaited
		case "RowExclusiveLock":
			p.getDBMetrics(db).rowExclusiveLockHeld = locks.held
			p.getDBMetrics(db).rowExclusiveLockAwaited = locks.awaited
		case "ShareUpdateExclusiveLock":
			p.getDBMetrics(db).shareUpdateExclusiveLockHeld = locks.held
			p.getDBMetrics(db).shareUpdateExclusiveLockAwaited = locks.awaited
		case "ShareLock":
			p.getDBMetrics(db).shareLockHeld = locks.held
			p.getDBMetrics(db).shareLockAwaited = locks.awaited
		case "ShareRowExclusiveLock":
			p.getDBMetrics(db).shareRowExclusiveLockHeld = locks.held
			p.getDBMetrics(db).shareRowExclusiveLockAwaited = locks.awaited
		case "ExclusiveLock":
			p.getDBMetrics(db).exclusiveLockHeld = locks.held
			p.getDBMetrics(db).exclusiveLockAwaited = locks.awaited
		case "AccessExclusiveLock":
			p.getDBMetrics(db).accessExclusiveLockHeld = locks.held
			p.getDBMetrics(db).accessExclusiveLockAwaited = locks.awaited
		}
	})
}
