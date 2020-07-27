package mysql

import (
	"strconv"
	"strings"
)

const queryGlobalStatus = "SHOW GLOBAL STATUS"

var globalStatusMetrics = []string{
	"Bytes_received",
	"Bytes_sent",
	"Queries",
	"Questions",
	"Slow_queries",
	"Handler_commit",
	"Handler_delete",
	"Handler_prepare",
	"Handler_read_first",
	"Handler_read_key",
	"Handler_read_next",
	"Handler_read_prev",
	"Handler_read_rnd",
	"Handler_read_rnd_next",
	"Handler_rollback",
	"Handler_savepoint",
	"Handler_savepoint_rollback",
	"Handler_update",
	"Handler_write",
	"Table_locks_immediate",
	"Table_locks_waited",
	"Select_full_join",
	"Select_full_range_join",
	"Select_range",
	"Select_range_check",
	"Select_scan",
	"Sort_merge_passes",
	"Sort_range",
	"Sort_scan",
	"Created_tmp_disk_tables",
	"Created_tmp_files",
	"Created_tmp_tables",
	"Connections",
	"Aborted_connects",
	"Max_used_connections",
	"Binlog_cache_disk_use",
	"Binlog_cache_use",
	"Threads_connected",
	"Threads_created",
	"Threads_cached",
	"Threads_running",
	"Thread_cache_misses",
	"Innodb_data_read",
	"Innodb_data_written",
	"Innodb_data_reads",
	"Innodb_data_writes",
	"Innodb_data_fsyncs",
	"Innodb_data_pending_reads",
	"Innodb_data_pending_writes",
	"Innodb_data_pending_fsyncs",
	"Innodb_log_waits",
	"Innodb_log_write_requests",
	"Innodb_log_writes",
	"Innodb_os_log_fsyncs",
	"Innodb_os_log_pending_fsyncs",
	"Innodb_os_log_pending_writes",
	"Innodb_os_log_written",
	"Innodb_row_lock_current_waits",
	"Innodb_rows_inserted",
	"Innodb_rows_read",
	"Innodb_rows_updated",
	"Innodb_rows_deleted",
	"Innodb_buffer_pool_pages_data",
	"Innodb_buffer_pool_pages_dirty",
	"Innodb_buffer_pool_pages_free",
	"Innodb_buffer_pool_pages_flushed",
	"Innodb_buffer_pool_pages_misc",
	"Innodb_buffer_pool_pages_total",
	"Innodb_buffer_pool_bytes_data",
	"Innodb_buffer_pool_bytes_dirty",
	"Innodb_buffer_pool_read_ahead",
	"Innodb_buffer_pool_read_ahead_evicted",
	"Innodb_buffer_pool_read_ahead_rnd",
	"Innodb_buffer_pool_read_requests",
	"Innodb_buffer_pool_write_requests",
	"Innodb_buffer_pool_reads",
	"Innodb_buffer_pool_wait_free",
	"Innodb_deadlocks",
	"Qcache_hits",
	"Qcache_lowmem_prunes",
	"Qcache_inserts",
	"Qcache_not_cached",
	"Qcache_queries_in_cache",
	"Qcache_free_memory",
	"Qcache_free_blocks",
	"Qcache_total_blocks",
	"Key_blocks_unused",
	"Key_blocks_used",
	"Key_blocks_not_flushed",
	"Key_read_requests",
	"Key_write_requests",
	"Key_reads",
	"Key_writes",
	"Open_files",
	"Opened_files",
	"Binlog_stmt_cache_disk_use",
	"Binlog_stmt_cache_use",
	"Connection_errors_accept",
	"Connection_errors_internal",
	"Connection_errors_max_connections",
	"Connection_errors_peer_address",
	"Connection_errors_select",
	"Connection_errors_tcpwrap",
	"Com_delete",
	"Com_insert",
	"Com_select",
	"Com_update",
	"Com_replace",
	"Opened_tables",
	"Open_tables",

	"wsrep_local_recv_queue",
	"wsrep_local_send_queue",
	"wsrep_received",
	"wsrep_replicated",
	"wsrep_received_bytes",
	"wsrep_replicated_bytes",
	"wsrep_local_bf_aborts",
	"wsrep_local_cert_failures",
	"wsrep_flow_control_paused_ns",
	"wsrep_cluster_weight",
	"wsrep_cluster_size",
	"wsrep_cluster_status",
	"wsrep_local_state",
	"wsrep_open_transactions",
	"wsrep_connected",
	"wsrep_ready",
	"wsrep_thread_count",
}

func (m *MySQL) collectGlobalStatus(collected map[string]int64) error {
	// MariaDB: https://mariadb.com/kb/en/server-status-variables/
	// MySQL: https://dev.mysql.com/doc/refman/8.0/en/server-status-variable-reference.html
	m.Debugf("executing query: '%s'", queryGlobalStatus)

	rows, err := m.db.Query(queryGlobalStatus)
	if err != nil {
		return err
	}
	defer rows.Close()

	set, err := rowsAsMap(rows)
	if err != nil {
		return err
	}

	for _, name := range globalStatusMetrics {
		strValue, ok := set[name]
		if !ok {
			continue
		}
		value, err := parseGlobalStatusValue(name, strValue)
		if err != nil {
			continue
		}
		collected[strings.ToLower(name)] = value
	}
	return nil
}

func parseGlobalStatusValue(name, value string) (int64, error) {
	if strings.HasPrefix(name, "wsrep_") {
		value = convertWsrepValue(name, value)
	}
	return strconv.ParseInt(value, 10, 64)
}

func convertWsrepValue(name, val string) string {
	switch name {
	case "wsrep_connected":
		return convertWsrepConnected(val)
	case "wsrep_ready":
		return convertWsrepReady(val)
	case "wsrep_cluster_status":
		return convertWsrepClusterStatus(val)
	default:
		return val
	}
}

func convertWsrepConnected(val string) string {
	// https://www.percona.com/doc/percona-xtradb-cluster/LATEST/wsrep-status-index.html#wsrep_connected
	switch val {
	case "OFF":
		return "0"
	case "ON":
		return "1"
	default:
		return "-1"
	}
}

func convertWsrepReady(val string) string {
	// https://www.percona.com/doc/percona-xtradb-cluster/LATEST/wsrep-status-index.html#wsrep_ready
	switch val {
	case "OFF":
		return "0"
	case "ON":
		return "1"
	default:
		return "-1"
	}
}

func convertWsrepClusterStatus(val string) string {
	// https://www.percona.com/doc/percona-xtradb-cluster/LATEST/wsrep-status-index.html#wsrep_cluster_status
	// https://github.com/codership/wsrep-API/blob/eab2d5d5a31672c0b7d116ef1629ff18392fd7d0/wsrep_api.h
	// typedef enum wsrep_view_status {
	//   WSREP_VIEW_PRIMARY,      //!< primary group configuration (quorum present)
	//   WSREP_VIEW_NON_PRIMARY,  //!< non-primary group configuration (quorum lost)
	//   WSREP_VIEW_DISCONNECTED, //!< not connected to group, retrying.
	//   WSREP_VIEW_MAX
	// } wsrep_view_status_t;
	switch strings.ToUpper(val) {
	case "PRIMARY":
		return "0"
	case "NON-PRIMARY":
		return "1"
	case "DISCONNECTED":
		return "2"
	default:
		return "-1"
	}
}
