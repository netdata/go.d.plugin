// SPDX-License-Identifier: GPL-3.0-or-later

package postgres

func queryServerVersion() string {
	return "SHOW server_version_num;"
}

func queryIsSuperUser() string {
	return "SELECT current_setting('is_superuser') = 'on' AS is_superuser;"
}

func querySettingsMaxConnections() string {
	return "SELECT current_setting('max_connections')::INT - current_setting('superuser_reserved_connections')::INT;"
}

// TODO: this is not correct and we should use pg_stat_activity.
// But we need to check what connections (backend_type) count towards 'max_connections'.
// I think python version query doesn't count it correctly.
// https://github.com/netdata/netdata/blob/1782e2d002bc5203128e5a5d2b801010e2822d2d/collectors/python.d.plugin/postgres/postgres.chart.py#L266
func queryServerCurrentConnectionsUsed() string {
	return "SELECT sum(numbackends) FROM pg_stat_database;"
}

func queryServerConnectionsState() string {
	return `
SELECT state,
       COUNT(*)
FROM pg_stat_activity
WHERE state IN
      (
       'active',
       'idle',
       'idle in transaction',
       'idle in transaction (aborted)',
       'fastpath function call',
       'disabled'
          )
GROUP BY state;
`
}

func queryCheckpoints() string {
	// definition by version: https://pgpedia.info/p/pg_stat_bgwriter.html
	// docs: https://www.postgresql.org/docs/current/monitoring-stats.html#MONITORING-PG-STAT-BGWRITER-VIEW
	// code: https://github.com/postgres/postgres/blob/366283961ac0ed6d89014444c6090f3fd02fce0a/src/backend/catalog/system_views.sql#L1104

	return `
SELECT checkpoints_timed,
       checkpoints_req,
       checkpoint_write_time,
       checkpoint_sync_time,
       buffers_checkpoint * current_setting('block_size')::numeric buffers_checkpoint,
       buffers_clean * current_setting('block_size')::numeric      buffers_clean,
       maxwritten_clean,
       buffers_backend * current_setting('block_size')::numeric    buffers_backend,
       buffers_backend_fsync,
       buffers_alloc * current_setting('block_size')::numeric      buffers_alloc
FROM pg_stat_bgwriter;
`
}

func queryServerUptime() string {
	return `SELECT EXTRACT(epoch FROM CURRENT_TIMESTAMP - pg_postmaster_start_time());`
}

func queryTXIDWraparound() string {
	// https://www.crunchydata.com/blog/managing-transaction-id-wraparound-in-postgresql
	return `
    WITH max_age AS ( SELECT
        2000000000 as max_old_xid,
        setting AS autovacuum_freeze_max_age 
    FROM
        pg_catalog.pg_settings 
    WHERE
        name = 'autovacuum_freeze_max_age'), per_database_stats AS ( SELECT
        datname ,
        m.max_old_xid::int ,
        m.autovacuum_freeze_max_age::int ,
        age(d.datfrozenxid) AS oldest_current_xid 
    FROM
        pg_catalog.pg_database d 
    JOIN
        max_age m 
            ON (true) 
    WHERE
        d.datallowconn) SELECT
        max(oldest_current_xid) AS oldest_current_xid ,
        max(ROUND(100*(oldest_current_xid/max_old_xid::float))) AS percent_towards_wraparound ,
        max(ROUND(100*(oldest_current_xid/autovacuum_freeze_max_age::float))) AS percent_towards_emergency_autovacuum 
    FROM
        per_database_stats;
`
}

func queryWALWrites(version int) string {
	if version < pgVersion10 {
		return `
SELECT
    pg_xlog_location_diff( 
    CASE
        pg_is_in_recovery() 
        WHEN
            TRUE 
        THEN
            pg_last_xlog_receive_location() 
        ELSE
            pg_current_xlog_location() 
    END
, '0/0') AS wal_writes ;
`
	}
	return `
SELECT
    pg_wal_lsn_diff( 
    CASE
        pg_is_in_recovery() 
        WHEN
            TRUE 
        THEN
            pg_last_wal_receive_lsn() 
        ELSE
            pg_current_wal_lsn() 
    END
, '0/0') AS wal_writes ;
`
}

func queryWALFiles(version int) string {
	if version < pgVersion10 {
		return `
SELECT count(*) FILTER (WHERE type = 'recycled') AS wal_recycled_files,
       count(*) FILTER (WHERE type = 'written')  AS wal_written_files
FROM (SELECT wal.name,
             pg_xlogfile_name(
                     CASE pg_is_in_recovery()
                         WHEN true THEN NULL
                         ELSE pg_current_xlog_location()
                         END),
             CASE
                 WHEN wal.name > pg_xlogfile_name(
                         CASE pg_is_in_recovery()
                             WHEN true THEN NULL
                             ELSE pg_current_xlog_location()
                             END) THEN 'recycled'
                 ELSE 'written'
                 END AS type
      FROM pg_catalog.pg_ls_dir('pg_xlog') AS wal(name)
      WHERE name ~ '^[0-9A-F]{24}$'
      ORDER BY (pg_stat_file('pg_xlog/' || name, true)).modification,
               wal.name DESC) sub;
`
	}
	return `
SELECT count(*) FILTER (WHERE type = 'recycled') AS wal_recycled_files,
       count(*) FILTER (WHERE type = 'written')  AS wal_written_files
FROM (SELECT wal.name,
             pg_walfile_name(
                     CASE pg_is_in_recovery()
                         WHEN true THEN NULL
                         ELSE pg_current_wal_lsn()
                         END),
             CASE
                 WHEN wal.name > pg_walfile_name(
                         CASE pg_is_in_recovery()
                             WHEN true THEN NULL
                             ELSE pg_current_wal_lsn()
                             END) THEN 'recycled'
                 ELSE 'written'
                 END AS type
      FROM pg_catalog.pg_ls_dir('pg_wal') AS wal(name)
      WHERE name ~ '^[0-9A-F]{24}$'
      ORDER BY (pg_stat_file('pg_wal/' || name, true)).modification,
               wal.name DESC) sub;
`
}

func queryWALArchiveFiles(version int) string {
	if version < pgVersion10 {
		return `
    SELECT
        CAST(COALESCE(SUM(CAST(archive_file ~ $r$\.ready$$r$ as INT)),
        0) AS INT) AS wal_archive_files_ready_count,
        CAST(COALESCE(SUM(CAST(archive_file ~ $r$\.done$$r$ AS INT)),
        0) AS INT)  AS wal_archive_files_done_count 
    FROM
        pg_catalog.pg_ls_dir('pg_xlog/archive_status') AS archive_files (archive_file);
`
	}
	return `
    SELECT
        CAST(COALESCE(SUM(CAST(archive_file ~ $r$\.ready$$r$ as INT)),
        0) AS INT) AS wal_archive_files_ready_count,
        CAST(COALESCE(SUM(CAST(archive_file ~ $r$\.done$$r$ AS INT)),
        0) AS INT)  AS wal_archive_files_done_count 
    FROM
        pg_catalog.pg_ls_dir('pg_wal/archive_status') AS archive_files (archive_file);
`
}

func queryCatalogRelations() string {
	// kind of same as
	// https://github.com/netdata/netdata/blob/750810e1798e09cc6210e83594eb9ed4905f8f12/collectors/python.d.plugin/postgres/postgres.chart.py#L336-L354
	// TODO: do we need that? It is optional and disabled by default in py version.
	return `
SELECT relkind,
       COUNT(1),
       SUM(relpages) * current_setting('block_size')::NUMERIC AS size
FROM pg_class
GROUP BY relkind;
`
}

func queryAutovacuumWorkers() string {
	// https://github.com/postgres/postgres/blob/9e4f914b5eba3f49ab99bdecdc4f96fac099571f/src/backend/postmaster/autovacuum.c#L3168-L3183
	return `
SELECT count(*) FILTER (
    WHERE
            query LIKE 'autovacuum: ANALYZE%%'
        AND query NOT LIKE '%%to prevent wraparound%%'
    )        AS autovacuum_analyze,
       count(*) FILTER (
           WHERE
                   query LIKE 'autovacuum: VACUUM ANALYZE%%'
               AND query NOT LIKE '%%to prevent wraparound%%'
           ) AS autovacuum_vacuum_analyze,
       count(*) FILTER (
           WHERE
                   query LIKE 'autovacuum: VACUUM %.%%'
               AND query NOT LIKE '%%to prevent wraparound%%'
           ) AS autovacuum_vacuum,
       count(*) FILTER (
           WHERE
           query LIKE '%%to prevent wraparound%%'
           ) AS autovacuum_vacuum_freeze,
       count(*) FILTER (
           WHERE
           query LIKE 'autovacuum: BRIN summarize%%'
           ) AS autovacuum_brin_summarize
FROM pg_stat_activity
WHERE query NOT LIKE '%%pg_stat_activity%%';
`
}

func queryActiveQueriesRunTime() string {
	return `
SELECT datname,
       EXTRACT(epoch from now() - query_start) as active_query_running_time
FROM pg_stat_activity
WHERE datname IS NOT NULL
  AND state = 'active'
  AND backend_type = 'client backend';
`
}

func queryReplicationStandbyAppDelta(version int) string {
	if version < pgVersion10 {
		return `
SELECT application_name,
       pg_xlog_location_diff(
               CASE pg_is_in_recovery()
                   WHEN true THEN pg_last_xlog_receive_location()
                   ELSE pg_current_xlog_location()
                   END,
               sent_location)   AS sent_delta,
       pg_xlog_location_diff(
               CASE pg_is_in_recovery()
                   WHEN true THEN pg_last_xlog_receive_location()
                   ELSE pg_current_xlog_location()
                   END,
               write_location)  AS write_delta,
       pg_xlog_location_diff(
               CASE pg_is_in_recovery()
                   WHEN true THEN pg_last_xlog_receive_location()
                   ELSE pg_current_xlog_location()
                   END,
               flush_location)  AS flush_delta,
       pg_xlog_location_diff(
               CASE pg_is_in_recovery()
                   WHEN true THEN pg_last_xlog_receive_location()
                   ELSE pg_current_xlog_location()
                   END,
               replay_location) AS replay_delta
FROM pg_stat_replication psr
WHERE application_name IS NOT NULL;
`
	}
	return `
SELECT application_name,
       pg_wal_lsn_diff(
               CASE pg_is_in_recovery()
                   WHEN true THEN pg_last_wal_receive_lsn()
                   ELSE pg_current_wal_lsn()
                   END,
               sent_lsn)   AS sent_delta,
       pg_wal_lsn_diff(
               CASE pg_is_in_recovery()
                   WHEN true THEN pg_last_wal_receive_lsn()
                   ELSE pg_current_wal_lsn()
                   END,
               write_lsn)  AS write_delta,
       pg_wal_lsn_diff(
               CASE pg_is_in_recovery()
                   WHEN true THEN pg_last_wal_receive_lsn()
                   ELSE pg_current_wal_lsn()
                   END,
               flush_lsn)  AS flush_delta,
       pg_wal_lsn_diff(
               CASE pg_is_in_recovery()
                   WHEN true THEN pg_last_wal_receive_lsn()
                   ELSE pg_current_wal_lsn()
                   END,
               replay_lsn) AS replay_delta
FROM pg_stat_replication
WHERE application_name IS NOT NULL;
`
}

func queryReplicationStandbyAppLag() string {
	return `
SELECT application_name,
       COALESCE(EXTRACT(EPOCH FROM write_lag)::bigint, 0)  AS write_lag,
       COALESCE(EXTRACT(EPOCH FROM flush_lag)::bigint, 0)  AS flush_lag,
       COALESCE(EXTRACT(EPOCH FROM replay_lag)::bigint, 0) AS replay_lag
FROM pg_stat_replication psr
WHERE application_name IS NOT NULL;
`
}

func queryReplicationSlotFiles(version int) string {
	if version < pgVersion11 {
		return `
WITH wal_size AS (
  SELECT
    current_setting('wal_block_size')::INT * setting::INT AS val
  FROM pg_settings
  WHERE name = 'wal_segment_size'
  )
SELECT
    slot_name,
    slot_type,
    replslot_wal_keep,
    count(slot_file) AS replslot_files
FROM
    (SELECT
        slot.slot_name,
        CASE
            WHEN slot_file <> 'state' THEN 1
        END AS slot_file ,
        slot_type,
        COALESCE (
          floor(
            CASE WHEN pg_is_in_recovery()
            THEN (
              pg_wal_lsn_diff(pg_last_wal_receive_lsn(), slot.restart_lsn)
              -- this is needed to account for whole WAL retention and
              -- not only size retention
              + (pg_wal_lsn_diff(restart_lsn, '0/0') % s.val)
            ) / s.val
            ELSE (
              pg_wal_lsn_diff(pg_current_wal_lsn(), slot.restart_lsn)
              -- this is needed to account for whole WAL retention and
              -- not only size retention
              + (pg_walfile_name_offset(restart_lsn)).file_offset
            ) / s.val
            END
          ),0) AS replslot_wal_keep
    FROM pg_replication_slots slot
    LEFT JOIN (
        SELECT
            slot2.slot_name,
            pg_ls_dir('pg_replslot/' || slot2.slot_name) AS slot_file
        FROM pg_replication_slots slot2
        ) files (slot_name, slot_file)
        ON slot.slot_name = files.slot_name
    CROSS JOIN wal_size s
    ) AS d
GROUP BY
    slot_name,
    slot_type,
    replslot_wal_keep;
`
	}

	return `
WITH wal_size AS (
  SELECT
    setting::int AS val
  FROM pg_settings
  WHERE name = 'wal_segment_size'
  )
SELECT
    slot_name,
    slot_type,
    replslot_wal_keep,
    count(slot_file) AS replslot_files
FROM
    (SELECT
        slot.slot_name,
        CASE
            WHEN slot_file <> 'state' THEN 1
        END AS slot_file ,
        slot_type,
        COALESCE (
          floor(
            CASE WHEN pg_is_in_recovery()
            THEN (
              pg_wal_lsn_diff(pg_last_wal_receive_lsn(), slot.restart_lsn)
              -- this is needed to account for whole WAL retention and
              -- not only size retention
              + (pg_wal_lsn_diff(restart_lsn, '0/0') % s.val)
            ) / s.val
            ELSE (
              pg_wal_lsn_diff(pg_current_wal_lsn(), slot.restart_lsn)
              -- this is needed to account for whole WAL retention and
              -- not only size retention
              + (pg_walfile_name_offset(restart_lsn)).file_offset
            ) / s.val
            END
          ),0) AS replslot_wal_keep
    FROM pg_replication_slots slot
    LEFT JOIN (
        SELECT
            slot2.slot_name,
            pg_ls_dir('pg_replslot/' || slot2.slot_name) AS slot_file
        FROM pg_replication_slots slot2
        ) files (slot_name, slot_file)
        ON slot.slot_name = files.slot_name
    CROSS JOIN wal_size s
    ) AS d
GROUP BY
    slot_name,
    slot_type,
    replslot_wal_keep;
`
}

func queryQueryableDatabaseList() string {
	return `
SELECT datname
FROM pg_database
WHERE datallowconn = true
  AND datistemplate = false
  AND datname != current_database()
  AND has_database_privilege((SELECT CURRENT_USER), datname, 'connect');
`
}

func queryDatabaseStats() string {
	// definition by version: https://pgpedia.info/p/pg_stat_database.html
	// docs: https://www.postgresql.org/docs/current/monitoring-stats.html#MONITORING-PG-STAT-DATABASE-VIEW
	// code: https://github.com/postgres/postgres/blob/366283961ac0ed6d89014444c6090f3fd02fce0a/src/backend/catalog/system_views.sql#L1018

	return `
SELECT stat.datname,
       numbackends,
       pg_database.datconnlimit,
       xact_commit,
       xact_rollback,
       blks_read,
       blks_hit,
       tup_returned,
       tup_fetched,
       tup_inserted,
       tup_updated,
       tup_deleted,
       conflicts,
       pg_database_size(stat.datname) AS size,
       temp_files,
       temp_bytes,
       deadlocks
FROM pg_stat_database stat
         INNER JOIN
     pg_database
     ON pg_database.datname = stat.datname
WHERE pg_database.datistemplate = false;
`
}

func queryDatabaseConflicts() string {
	// definition by version: https://pgpedia.info/p/pg_stat_database_conflicts.html
	// docs: https://www.postgresql.org/docs/current/monitoring-stats.html#MONITORING-PG-STAT-DATABASE-CONFLICTS-VIEW
	// code: https://github.com/postgres/postgres/blob/366283961ac0ed6d89014444c6090f3fd02fce0a/src/backend/catalog/system_views.sql#L1058

	return `
SELECT stat.datname,
       confl_tablespace,
       confl_lock,
       confl_snapshot,
       confl_bufferpin,
       confl_deadlock
FROM pg_stat_database_conflicts stat
         INNER JOIN
     pg_database
     ON pg_database.datname = stat.datname
WHERE pg_database.datistemplate = false;
`
}

func queryDatabaseLocks() string {
	// definition by version: https://pgpedia.info/p/pg_locks.html
	// docs: https://www.postgresql.org/docs/current/view-pg-locks.html

	return `
SELECT pg_database.datname,
       mode,
       granted,
       count(mode) AS locks_count
FROM pg_locks
         INNER JOIN
     pg_database
     ON pg_database.oid = pg_locks.database
WHERE pg_database.datistemplate = false
GROUP BY datname,
         mode,
         granted
ORDER BY datname,
         mode;
`
}

func queryUserTableStats() string {
	return `
SELECT current_database()                                   as datname,
       schemaname,
       relname,
       seq_scan,
       seq_tup_read,
       idx_scan,
       idx_tup_fetch,
       n_tup_ins,
       n_tup_upd,
       n_tup_del,
       n_tup_hot_upd,
       n_live_tup,
       n_dead_tup,
       EXTRACT(epoch from now() - last_vacuum)              as last_vacuum,
       EXTRACT(epoch from now() - last_autovacuum)          as last_autovacuum,
       EXTRACT(epoch from now() - last_analyze)             as last_analyze,
       EXTRACT(epoch from now() - last_autoanalyze)         as last_autoanalyze,
       vacuum_count,
       autovacuum_count,
       analyze_count,
       autoanalyze_count,
       pg_total_relation_size(schemaname || '.' || relname) as total_relation_size
FROM pg_stat_user_tables;
`
}
