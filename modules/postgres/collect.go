// SPDX-License-Identifier: GPL-3.0-or-later

package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/stdlib"
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
		db, err := p.openPrimaryConnection()
		if err != nil {
			return nil, err
		}
		p.db = db
	}

	if p.pgVersion == 0 {
		ver, err := p.doQueryServerVersion()
		if err != nil {
			return nil, fmt.Errorf("querying server version error: %v", err)
		}
		p.pgVersion = ver
		p.Debugf("connected to PostgreSQL v%d", p.pgVersion)
	}

	if p.superUser == nil {
		v, err := p.doQueryIsSuperUser()
		if err != nil {
			return nil, fmt.Errorf("querying is super user error: %v", err)
		}
		p.superUser = &v
		p.Debugf("connected as super user: %v", *p.superUser)
	}

	now := time.Now()

	if now.Sub(p.recheckSettingsTime) > p.recheckSettingsEvery {
		p.recheckSettingsTime = now
		maxConn, err := p.doQuerySettingsMaxConnections()
		if err != nil {
			return nil, fmt.Errorf("querying settings max connections error: %v", err)
		}
		p.mx.maxConnections = maxConn
	}

	p.resetMetrics()

	if err := p.doQueryGlobalMetrics(); err != nil {
		return nil, err
	}
	if err := p.doQueryReplicationMetrics(); err != nil {
		return nil, err
	}
	if err := p.doQueryDatabasesMetrics(); err != nil {
		return nil, err
	}
	if err := p.doQueryTablesMetrics(); err != nil {
		return nil, err
	}

	mx := make(map[string]int64)
	p.collectMetrics(mx)

	return mx, nil
}

func (p *Postgres) collectMetrics(mx map[string]int64) {
	mx["server_connections_used"] = p.mx.connUsed
	if p.mx.maxConnections > 0 {
		mx["server_connections_available"] = p.mx.maxConnections - p.mx.connUsed
		mx["server_connections_utilization"] = calcPercentage(p.mx.connUsed, p.mx.maxConnections)
	}
	mx["server_uptime"] = p.mx.uptime
	mx["server_connections_state_active"] = p.mx.connStateActive
	mx["server_connections_state_idle"] = p.mx.connStateIdle
	mx["server_connections_state_idle_in_transaction"] = p.mx.connStateIdleInTrans
	mx["server_connections_state_idle_in_transaction_aborted"] = p.mx.connStateIdleInTransAborted
	mx["server_connections_state_fastpath_function_call"] = p.mx.connStateFastpathFunctionCall
	mx["server_connections_state_disabled"] = p.mx.connStateDisabled
	mx["checkpoints_timed"] = p.mx.checkpointsTimed
	mx["checkpoints_req"] = p.mx.checkpointsReq
	mx["checkpoint_write_time"] = p.mx.checkpointWriteTime
	mx["checkpoint_sync_time"] = p.mx.checkpointSyncTime
	mx["buffers_checkpoint"] = p.mx.buffersCheckpoint
	mx["buffers_clean"] = p.mx.buffersClean
	mx["maxwritten_clean"] = p.mx.maxwrittenClean
	mx["buffers_backend"] = p.mx.buffersBackend
	mx["buffers_backend_fsync"] = p.mx.buffersBackendFsync
	mx["buffers_alloc"] = p.mx.buffersAlloc
	mx["oldest_current_xid"] = p.mx.oldestXID
	mx["percent_towards_wraparound"] = p.mx.percentTowardsWraparound
	mx["percent_towards_emergency_autovacuum"] = p.mx.percentTowardsEmergencyAutovacuum
	mx["wal_writes"] = p.mx.walWrites
	mx["wal_recycled_files"] = p.mx.walRecycledFiles
	mx["wal_written_files"] = p.mx.walWrittenFiles
	mx["wal_archive_files_ready_count"] = p.mx.walArchiveFilesReady
	mx["wal_archive_files_done_count"] = p.mx.walArchiveFilesDone
	mx["catalog_relkind_r_count"] = p.mx.relkindOrdinaryTable
	mx["catalog_relkind_i_count"] = p.mx.relkindIndex
	mx["catalog_relkind_S_count"] = p.mx.relkindSequence
	mx["catalog_relkind_t_count"] = p.mx.relkindTOASTTable
	mx["catalog_relkind_v_count"] = p.mx.relkindView
	mx["catalog_relkind_m_count"] = p.mx.relkindMatView
	mx["catalog_relkind_c_count"] = p.mx.relkindCompositeType
	mx["catalog_relkind_f_count"] = p.mx.relkindForeignTable
	mx["catalog_relkind_p_count"] = p.mx.relkindPartitionedTable
	mx["catalog_relkind_I_count"] = p.mx.relkindPartitionedIndex
	mx["catalog_relkind_r_size"] = p.mx.relkindOrdinaryTableSize
	mx["catalog_relkind_i_size"] = p.mx.relkindIndexSize
	mx["catalog_relkind_S_size"] = p.mx.relkindSequenceSize
	mx["catalog_relkind_t_size"] = p.mx.relkindTOASTTableSize
	mx["catalog_relkind_v_size"] = p.mx.relkindViewSize
	mx["catalog_relkind_m_size"] = p.mx.relkindMatViewSize
	mx["catalog_relkind_c_size"] = p.mx.relkindCompositeTypeSize
	mx["catalog_relkind_f_size"] = p.mx.relkindForeignTableSize
	mx["catalog_relkind_p_size"] = p.mx.relkindPartitionedTableSize
	mx["catalog_relkind_I_size"] = p.mx.relkindPartitionedIndexSize
	mx["autovacuum_analyze"] = p.mx.autovacuumWorkersAnalyze
	mx["autovacuum_vacuum_analyze"] = p.mx.autovacuumWorkersVacuumAnalyze
	mx["autovacuum_vacuum"] = p.mx.autovacuumWorkersVacuum
	mx["autovacuum_vacuum_freeze"] = p.mx.autovacuumWorkersVacuumFreeze
	mx["autovacuum_brin_summarize"] = p.mx.autovacuumWorkersBrinSummarize

	for name, m := range p.mx.dbs {
		if !m.updated {
			delete(p.mx.dbs, name)
			p.removeDatabaseCharts(name)
			continue
		}
		if !m.hasCharts {
			m.hasCharts = true
			p.addNewDatabaseCharts(name)
		}
		px := "db_" + m.name + "_"
		mx[px+"numbackends"] = m.numBackends
		if m.datConnLimit <= 0 {
			mx[px+"numbackends_utilization"] = calcPercentage(m.numBackends, p.mx.maxConnections)
		} else {
			mx[px+"numbackends_utilization"] = calcPercentage(m.numBackends, m.datConnLimit)
		}
		mx[px+"xact_commit"] = m.xactCommit
		mx[px+"xact_rollback"] = m.xactRollback
		mx[px+"blks_read"] = m.blksRead
		mx[px+"blks_hit"] = m.blksHit
		mx[px+"tup_returned"] = m.tupReturned
		mx[px+"tup_fetched"] = m.tupFetched
		mx[px+"tup_inserted"] = m.tupInserted
		mx[px+"tup_updated"] = m.tupUpdated
		mx[px+"tup_deleted"] = m.tupDeleted
		mx[px+"conflicts"] = m.conflicts
		mx[px+"size"] = m.size
		mx[px+"temp_files"] = m.tempFiles
		mx[px+"temp_bytes"] = m.tempBytes
		mx[px+"deadlocks"] = m.deadlocks
		mx[px+"confl_tablespace"] = m.conflTablespace
		mx[px+"confl_lock"] = m.conflLock
		mx[px+"confl_snapshot"] = m.conflSnapshot
		mx[px+"confl_bufferpin"] = m.conflBufferpin
		mx[px+"confl_deadlock"] = m.conflDeadlock
		mx[px+"lock_mode_AccessShareLock_held"] = m.accessShareLockHeld
		mx[px+"lock_mode_RowShareLock_held"] = m.rowShareLockHeld
		mx[px+"lock_mode_RowExclusiveLock_held"] = m.rowExclusiveLockHeld
		mx[px+"lock_mode_ShareUpdateExclusiveLock_held"] = m.shareUpdateExclusiveLockHeld
		mx[px+"lock_mode_ShareLock_held"] = m.shareLockHeld
		mx[px+"lock_mode_ShareRowExclusiveLock_held"] = m.shareRowExclusiveLockHeld
		mx[px+"lock_mode_ExclusiveLock_held"] = m.exclusiveLockHeld
		mx[px+"lock_mode_AccessExclusiveLock_held"] = m.accessExclusiveLockHeld
		mx[px+"lock_mode_AccessShareLock_awaited"] = m.accessShareLockAwaited
		mx[px+"lock_mode_RowShareLock_awaited"] = m.rowShareLockAwaited
		mx[px+"lock_mode_RowExclusiveLock_awaited"] = m.rowExclusiveLockAwaited
		mx[px+"lock_mode_ShareUpdateExclusiveLock_awaited"] = m.shareUpdateExclusiveLockAwaited
		mx[px+"lock_mode_ShareLock_awaited"] = m.shareLockAwaited
		mx[px+"lock_mode_ShareRowExclusiveLock_awaited"] = m.shareRowExclusiveLockAwaited
		mx[px+"lock_mode_ExclusiveLock_awaited"] = m.exclusiveLockAwaited
		mx[px+"lock_mode_AccessExclusiveLock_awaited"] = m.accessExclusiveLockAwaited
	}

	for name, m := range p.mx.tables {
		if !m.updated {
			delete(p.mx.tables, name)
			p.removeTableCharts(m.db, m.schema, m.name)
			continue
		}
		if !m.hasCharts {
			m.hasCharts = true
			p.addNewTableCharts(m.db, m.schema, m.name)
		}
		if !m.hasLastAutoVacuumChart && m.lastAutoVacuumAgo != -1 {
			m.hasLastAutoVacuumChart = true
			p.addTableLastAutoVacuumAgoChart(m.db, m.schema, m.name)
		}
		if !m.hasLastVacuumChart && m.lastVacuumAgo != -1 {
			m.hasLastVacuumChart = true
			p.addTableLastVacuumAgoChart(m.db, m.schema, m.name)
		}
		if !m.hasLastAutoAnalyzeChart && m.lastAutoAnalyzeAgo != -1 {
			m.hasLastAutoAnalyzeChart = true
			p.addTableLastAutoAnalyzeAgoChart(m.db, m.schema, m.name)
		}
		if !m.hasLastAnalyzeChart && m.lastAnalyzeAgo != -1 {
			m.hasLastAnalyzeChart = true
			p.addTableLastAnalyzeAgoChart(m.db, m.schema, m.name)
		}

		px := fmt.Sprintf("db_%s_schema_%s_table_%s_", m.db, m.schema, m.name)

		mx[px+"seq_scan"] = m.seqScan
		mx[px+"seq_tup_read"] = m.seqTupRead
		mx[px+"idx_scan"] = m.idxScan
		mx[px+"idx_tup_fetch"] = m.idxTupFetch
		mx[px+"n_live_tup"] = m.nLiveTup
		mx[px+"n_dead_tup"] = m.nDeadTup
		mx[px+"n_tup_ins"] = m.nTupIns
		mx[px+"n_tup_upd"] = m.nTupUpd
		mx[px+"n_tup_del"] = m.nTupDel
		mx[px+"n_tup_hot_upd"] = m.nTupHotUpd
		if m.lastAutoVacuumAgo != -1 {
			mx[px+"last_autovacuum_ago"] = m.lastAutoVacuumAgo
		}
		if m.lastVacuumAgo != -1 {
			mx[px+"last_vacuum_ago"] = m.lastVacuumAgo
		}
		if m.lastAutoAnalyzeAgo != -1 {
			mx[px+"last_autoanalyze_ago"] = m.lastAutoAnalyzeAgo
		}
		if m.lastAnalyzeAgo != -1 {
			mx[px+"last_analyze_ago"] = m.lastAnalyzeAgo
		}
		mx[px+"total_size"] = m.totalSize
	}

	for name, m := range p.mx.replApps {
		if !m.updated {
			delete(p.mx.replApps, name)
			p.removeReplicationStandbyAppCharts(name)
			continue
		}
		if !m.hasCharts {
			m.hasCharts = true
			p.addNewReplicationStandbyAppCharts(name)
		}
		px := "repl_standby_app_" + m.name + "_wal_"
		mx[px+"sent_delta"] = m.walSentDelta
		mx[px+"write_delta"] = m.walWriteDelta
		mx[px+"flush_delta"] = m.walFlushDelta
		mx[px+"replay_delta"] = m.walReplayDelta
		mx[px+"write_lag"] = m.walWriteLag
		mx[px+"flush_lag"] = m.walFlushLag
		mx[px+"replay_lag"] = m.walReplayLag
	}

	for name, m := range p.mx.replSlots {
		if !m.updated {
			delete(p.mx.replSlots, name)
			p.removeReplicationSlotCharts(name)
			continue
		}
		if !m.hasCharts {
			m.hasCharts = true
			p.addNewReplicationSlotCharts(name)
		}
		px := "repl_slot_" + m.name + "_"
		mx[px+"replslot_wal_keep"] = m.walKeep
		mx[px+"replslot_files"] = m.files
	}
}

func (p *Postgres) resetMetrics() {
	p.mx.srvMetrics = srvMetrics{
		maxConnections: p.mx.maxConnections,
	}
	for name, m := range p.mx.dbs {
		p.mx.dbs[name] = &dbMetrics{
			name:      m.name,
			hasCharts: m.hasCharts,
		}
	}
	for name, m := range p.mx.tables {
		p.mx.tables[name] = &tableMetrics{
			db:        m.db,
			schema:    m.schema,
			name:      m.name,
			hasCharts: m.hasCharts,
		}
	}
	for name, m := range p.mx.replApps {
		p.mx.replApps[name] = &replStandbyAppMetrics{
			name:      m.name,
			hasCharts: m.hasCharts,
		}
	}
	for name, m := range p.mx.replSlots {
		p.mx.replSlots[name] = &replSlotMetrics{
			name:      m.name,
			hasCharts: m.hasCharts,
		}
	}
}

func (p *Postgres) openPrimaryConnection() (*sql.DB, error) {
	db, err := sql.Open("pgx", p.DSN)
	if err != nil {
		return nil, fmt.Errorf("error on opening a connection with the Postgres database [%s]: %v", p.DSN, err)
	}

	db.SetMaxOpenConns(1)
	db.SetMaxIdleConns(1)
	db.SetConnMaxLifetime(10 * time.Minute)

	ctx, cancel := context.WithTimeout(context.Background(), p.Timeout.Duration)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("error on pinging the Postgres database [%s]: %v", p.DSN, err)
	}

	return db, nil
}

func (p *Postgres) openSecondaryConnection(dbname string) (*sql.DB, string, error) {
	cfg, err := pgx.ParseConfig(p.DSN)
	if err != nil {
		return nil, "", fmt.Errorf("error on parsing DSN [%s]: %v", p.DSN, err)
	}

	cfg.Database = dbname
	connString := stdlib.RegisterConnConfig(cfg)

	db, err := sql.Open("pgx", connString)
	if err != nil {
		return nil, "", fmt.Errorf("error on opening a secondary connection with the Postgres database [%s]: %v", dbname, err)
	}

	db.SetMaxOpenConns(1)
	db.SetMaxIdleConns(1)
	db.SetConnMaxLifetime(10 * time.Minute)

	ctx, cancel := context.WithTimeout(context.Background(), p.Timeout.Duration)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		_ = db.Close()
		return nil, "", fmt.Errorf("error on pinging the secondary Postgres database [%s]: %v", dbname, err)
	}

	return db, connString, nil
}

func (p *Postgres) isSuperUser() bool { return p.superUser != nil && *p.superUser }

func (p *Postgres) getDBMetrics(name string) *dbMetrics {
	db, ok := p.mx.dbs[name]
	if !ok {
		db = &dbMetrics{name: name}
		p.mx.dbs[name] = db
	}
	return db
}

func (p *Postgres) getTableMetrics(db, schema, name string) *tableMetrics {
	key := db + "_" + schema + "_" + name
	m, ok := p.mx.tables[key]
	if !ok {
		m = &tableMetrics{db: db, schema: schema, name: name}
		p.mx.tables[key] = m
	}
	return m
}

func (p *Postgres) getReplAppMetrics(name string) *replStandbyAppMetrics {
	app, ok := p.mx.replApps[name]
	if !ok {
		app = &replStandbyAppMetrics{name: name}
		p.mx.replApps[name] = app
	}
	return app
}

func (p *Postgres) getReplSlotMetrics(name string) *replSlotMetrics {
	slot, ok := p.mx.replSlots[name]
	if !ok {
		slot = &replSlotMetrics{name: name}
		p.mx.replSlots[name] = slot
	}
	return slot
}

func parseInt(s string) int64 {
	v, _ := strconv.ParseInt(s, 10, 64)
	return v
}

func parseFloat(s string) int64 {
	v, _ := strconv.ParseFloat(s, 64)
	return int64(v)
}

func calcPercentage(value, total int64) int64 {
	if total == 0 {
		return 0
	}
	return value * 100 / total
}
