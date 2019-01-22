package mysql

import (
	"database/sql"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/netdata/go.d.plugin/modules"

	_ "github.com/go-sql-driver/mysql"
)

func init() {
	creator := modules.Creator{
		Create: func() modules.Module { return New() },
	}

	modules.Register("mysql", creator)
}

// MySQL is the mysql database module.
type MySQL struct {
	modules.Base
	db *sql.DB
	// i.e user:password@/dbname
	DSN  string `yaml:"dsn"`
	user string
}

// New creates and returns a new empty MySQL module.
func New() *MySQL {
	// return &MySQL{SlowQueriesMaxCount: 0, SlowQueriesInPercentile: 95, SlowQueriesSeconds: 4}
	return &MySQL{}
}

// // CompatibleMinimumVersion is the minimum required version of the mysql server.
// const CompatibleMinimumVersion = 5.1
//
// func (m *MySQL) getMySQLVersion() float64 {
// 	var versionStr string
// 	var versionNum float64
// 	if err := m.db.QueryRow("SELECT @@version").Scan(&versionStr); err == nil {
// 		versionNum, _ = strconv.ParseFloat(regexp.MustCompile(`^\d+\.\d+`).FindString(versionStr), 64)
// 	}
//
// 	return versionNum
// }

// Cleanup performs cleanup.
func (m *MySQL) Cleanup() {
	err := m.db.Close()
	if err != nil {
		m.Errorf("cleanup: error on closing the mysql database [%s]: %v", m.DSN, err)
	}
}

// Init makes initialization of the MySQL mod.
func (m *MySQL) Init() bool {
	if m.DSN == "" {
		return false
	}

	// test the connectivity here.
	if err := m.openConnection(); err != nil {
		return false
	}

	m.user = parseUser(m.DSN)

	// if min, got := CompatibleMinimumVersion, m.getMySQLVersion(); min > 0 && got < min {
	// 	m.Warningf("running with uncompatible mysql version [%v<%v]", got, min)
	// }

	// post Init debug info.
	m.Debugf("connected using DSN [%s]", m.DSN)
	return true
}

func parseUser(dsn string) string {
	if userIdx := strings.IndexRune(dsn, ':'); userIdx != -1 {
		return dsn[:userIdx]
	} else if userIdx = strings.IndexRune(dsn, '@'); userIdx != -1 {
		return dsn[:userIdx]
	}

	return ""
}

func (m *MySQL) openConnection() error {
	if m.db != nil {
		if err := m.db.Ping(); err != nil {
			m.db.Close()
			m.db = nil

			return m.openConnection()
		}

		return nil
	}

	db, err := sql.Open("mysql", m.DSN)
	if err != nil {
		m.Errorf("error on opening a connection with the mysql database [%s]: %v", m.DSN, err)
		return err
	}
	db.SetConnMaxLifetime(1 * time.Minute)

	if err = db.Ping(); err != nil {
		db.Close()
		m.Errorf("error on pinging the mysql database [%s]: %v", m.DSN, err)
		return err
	}

	m.db = db
	return nil
}

// Check makes check.
func (m *MySQL) Check() bool {
	return len(m.Collect()) > 0
}

// Charts creates Charts.
func (m *MySQL) Charts() *Charts {
	return charts.Copy()
}

// Collect collects health checks and metrics for MySQL.
func (m *MySQL) Collect() map[string]int64 {
	if err := m.openConnection(); err != nil {
		return nil
	}

	metrics := make(map[string]int64)

	if err := m.collectGlobalStats(metrics); err != nil {
		m.Errorf("error on collecting global stats: %v", err)
		return nil
	}

	if err := m.collectSlaveStatus(metrics); err != nil {
		m.Errorf("error on collecting slave status: %v", err)
		return nil
	}

	if err := m.collectMaxConnections(metrics); err != nil {
		m.Errorf("error on determinating max connections: %v", err)
		return nil
	}

	return metrics
}

var globalStats = map[string]bool{
	"Bytes_received":                        true,
	"Bytes_sent":                            true,
	"Queries":                               true,
	"Questions":                             true,
	"Slow_queries":                          true, // can be configured by the mysql user.
	"Handler_commit":                        true,
	"Handler_delete":                        true,
	"Handler_prepare":                       true,
	"Handler_read_first":                    true,
	"Handler_read_key":                      true,
	"Handler_read_next":                     true,
	"Handler_read_prev":                     true,
	"Handler_read_rnd":                      true,
	"Handler_read_rnd_next":                 true,
	"Handler_rollback":                      true,
	"Handler_savepoint":                     true,
	"Handler_savepoint_rollback":            true,
	"Handler_update":                        true,
	"Handler_write":                         true,
	"Table_locks_immediate":                 true,
	"Table_locks_waited":                    true,
	"Select_full_join":                      true,
	"Select_full_range_join":                true,
	"Select_range":                          true,
	"Select_range_check":                    true,
	"Select_scan":                           true,
	"Sort_merge_passes":                     true,
	"Sort_range":                            true,
	"Sort_scan":                             true,
	"Created_tmp_disk_tables":               true,
	"Created_tmp_files":                     true,
	"Created_tmp_tables":                    true,
	"Connections":                           true,
	"Aborted_connects":                      true,
	"Max_used_connections":                  true,
	"Binlog_cache_disk_use":                 true,
	"Binlog_cache_use":                      true,
	"Threads_connected":                     true,
	"Threads_created":                       true,
	"Threads_cached":                        true,
	"Threads_running":                       true,
	"Innodb_data_read":                      true,
	"Innodb_data_written":                   true,
	"Innodb_data_reads":                     true,
	"Innodb_data_writes":                    true,
	"Innodb_data_fsyncs":                    true,
	"Innodb_data_pending_reads":             true,
	"Innodb_data_pending_writes":            true,
	"Innodb_data_pending_fsyncs":            true,
	"Innodb_log_waits":                      true,
	"Innodb_log_write_requests":             true,
	"Innodb_log_writes":                     true,
	"Innodb_os_log_fsyncs":                  true,
	"Innodb_os_log_pending_fsyncs":          true,
	"Innodb_os_log_pending_writes":          true,
	"Innodb_os_log_written":                 true,
	"Innodb_row_lock_current_waits":         true,
	"Innodb_rows_inserted":                  true,
	"Innodb_rows_read":                      true,
	"Innodb_rows_updated":                   true,
	"Innodb_rows_deleted":                   true,
	"Innodb_buffer_pool_pages_data":         true,
	"Innodb_buffer_pool_pages_dirty":        true,
	"Innodb_buffer_pool_pages_free":         true,
	"Innodb_buffer_pool_pages_flushed":      true,
	"Innodb_buffer_pool_pages_misc":         true,
	"Innodb_buffer_pool_pages_total":        true,
	"Innodb_buffer_pool_bytes_data":         true,
	"Innodb_buffer_pool_bytes_dirty":        true,
	"Innodb_buffer_pool_read_ahead":         true,
	"Innodb_buffer_pool_read_ahead_evicted": true,
	"Innodb_buffer_pool_read_ahead_rnd":     true,
	"Innodb_buffer_pool_read_requests":      true,
	"Innodb_buffer_pool_write_requests":     true,
	"Innodb_buffer_pool_reads":              true,
	"Innodb_buffer_pool_wait_free":          true,
	"Qcache_hits":                           true,
	"Qcache_lowmem_prunes":                  true,
	"Qcache_inserts":                        true,
	"Qcache_not_cached":                     true,
	"Qcache_queries_in_cache":               true,
	"Qcache_free_memory":                    true,
	"Qcache_free_blocks":                    true,
	"Qcache_total_blocks":                   true,
	"Key_blocks_unused":                     true,
	"Key_blocks_used":                       true,
	"Key_blocks_not_flushed":                true,
	"Key_read_requests":                     true,
	"Key_write_requests":                    true,
	"Key_reads":                             true,
	"Key_writes":                            true,
	"Open_files":                            true,
	"Opened_files":                          true,
	"Binlog_stmt_cache_disk_use":            true,
	"Binlog_stmt_cache_use":                 true,
	"Connection_errors_accept":              true,
	"Connection_errors_internal":            true,
	"Connection_errors_max_connections":     true,
	"Connection_errors_peer_address":        true,
	"Connection_errors_select":              true,
	"Connection_errors_tcpwrap":             true,
	"wsrep_local_recv_queue":                true,
	"wsrep_local_send_queue":                true,
	"wsrep_received":                        true,
	"wsrep_replicated":                      true,
	"wsrep_received_bytes":                  true,
	"wsrep_replicated_bytes":                true,
	"wsrep_local_bf_aborts":                 true,
	"wsrep_local_cert_failures":             true,
	"wsrep_flow_control_paused_ns":          true,
	"Com_delete":                            true,
	"Com_insert":                            true,
	"Com_select":                            true,
	"Com_update":                            true,
	"Com_replace":                           true,
}

func extractStatKey(s string) (string, bool) {
	_, ok := globalStats[s]
	return toKey(s), ok
}

func toKey(s string) string {
	return strings.ToLower(s)
}

func (m *MySQL) collectGlobalStats(metrics map[string]int64) error {
	rows, err := m.db.Query("SHOW GLOBAL STATUS")
	if err != nil {
		return err
	}

	defer rows.Close()

	for rows.Next() {
		var (
			varName string
			// it does not always returns int64, i.e public key []uint8(bytes) we can't catch it before Scan, so interface{}.
			value interface{}
		)

		// Variable_name
		// Value
		err = rows.Scan(&varName, &value)
		if err != nil {
			return err
		}

		if key, ok := extractStatKey(varName); ok {
			metrics[key], err = strconv.ParseInt(fmt.Sprintf("%s", value), 10, 64)
			if err != nil {
				return err
			}
		}
	}

	v1, ok1 := metrics["threads_created"]
	v2, ok2 := metrics["connections"]

	// NOTE: not sure this check is needed
	if ok1 && ok2 {
		metrics["thread_cache_misses"] = int64(float64(v1) / float64(v2) * 10000)
	}

	return nil

	/*
	   [sort_merge_passes] = [0]
	   [com_delete] = [0]
	   [innodb_buffer_pool_read_ahead_rnd] = [0]
	   [innodb_data_pending_fsyncs] = [0]
	   [key_blocks_unused] = [6698]
	   [key_read_requests] = [0]
	   [key_reads] = [0]
	   [select_full_range_join] = [0]
	   [bytes_received] = [4014]
	   [connections] = [22]
	   [handler_read_prev] = [0]
	   [innodb_buffer_pool_wait_free] = [0]
	   [innodb_data_writes] = [362]
	   [innodb_data_written] = [5153792]
	   [table_locks_waited] = [0]
	   [created_tmp_files] = [5]
	   [innodb_data_read] = [15864832]
	   [open_files] = [7]
	   [com_select] = [18]
	   [innodb_log_waits] = [0]
	   [select_range] = [0]
	   [handler_read_next] = [6081]
	   [innodb_buffer_pool_bytes_data] = [4210688]
	   [innodb_buffer_pool_pages_misc] = [6]
	   [sort_scan] = [4]
	   [handler_savepoint] = [0]
	   [handler_update] = [346]
	   [innodb_buffer_pool_pages_total] = [512]
	   [innodb_os_log_fsyncs] = [71]
	   [queries] = [61]
	   [threads_cached] = [0]
	   [innodb_buffer_pool_read_ahead] = [14]
	   [bytes_sent] = [156380]
	   [connection_errors_peer_address] = [0]
	   [created_tmp_disk_tables] = [0]
	   [handler_rollback] = [0]
	   [handler_savepoint_rollback] = [0]
	   [handler_write] = [5222]
	   [innodb_buffer_pool_pages_free] = [249]
	   [innodb_buffer_pool_write_requests] = [1945]
	   [sort_range] = [0]
	   [threads_connected] = [3]
	   [created_tmp_tables] = [20]
	   [handler_read_first] = [39]
	   [innodb_buffer_pool_read_requests] = [21518]
	   [key_blocks_used] = [0]
	   [max_used_connections] = [3]
	   [select_full_join] = [2]
	   [innodb_buffer_pool_pages_flushed] = [215]
	   [binlog_cache_use] = [0]
	   [com_insert] = [0]
	   [com_replace] = [0]
	   [connection_errors_internal] = [0]
	   [connection_errors_tcpwrap] = [0]
	   [innodb_buffer_pool_pages_data] = [257]
	   [innodb_buffer_pool_pages_dirty] = [0]
	   [innodb_os_log_written] = [89600]
	   [threads_created] = [3]
	   [innodb_data_fsyncs] = [105]
	   [innodb_data_reads] = [987]
	   [innodb_rows_read] = [7208]
	   [key_write_requests] = [0]
	   [key_writes] = [0]
	   [innodb_data_pending_writes] = [0]
	   [innodb_os_log_pending_fsyncs] = [0]
	   [aborted_connects] = [1]
	   [handler_prepare] = [0]
	   [innodb_log_writes] = [98]
	   [innodb_rows_deleted] = [0]
	   [innodb_rows_updated] = [346]
	   [key_blocks_not_flushed] = [0]
	   [binlog_cache_disk_use] = [0]
	   [connection_errors_accept] = [0]
	   [connection_errors_max_connections] = [0]
	   [connection_errors_select] = [0]
	   [innodb_buffer_pool_read_ahead_evicted] = [0]
	   [innodb_os_log_pending_writes] = [0]
	   [binlog_stmt_cache_use] = [0]
	   [com_update] = [0]
	   [handler_delete] = [0]
	   [handler_read_key] = [2363]
	   [handler_read_rnd] = [360]
	   [handler_read_rnd_next] = [10166]
	   [innodb_data_pending_reads] = [0]
	   [innodb_rows_inserted] = [0]
	   [handler_commit] = [678]
	   [innodb_log_write_requests] = [765]
	   [opened_files] = [28]
	   [questions] = [37]
	   [select_range_check] = [0]
	   [binlog_stmt_cache_disk_use] = [0]
	   [innodb_buffer_pool_bytes_dirty] = [0]
	   [innodb_buffer_pool_reads] = [950]
	   [innodb_row_lock_current_waits] = [0]
	   [select_scan] = [35]
	   [table_locks_immediate] = [12]
	   [threads_running] = [2]
	*/
}

func slaveSeconds(value interface{}) int64 {
	v, err := strconv.ParseInt(fmt.Sprintf("%s", value), 10, 64)
	if err != nil {
		return -1
	}
	return v
}

func slaveRunning(value interface{}) int64 {
	if v, ok := value.(string); ok {
		if v == "Yes" {
			return 1
		}

		return -1
	}

	if v := fmt.Sprintf("%s", value); v == "Yes" {
		return 1
	}

	return -1
}

var slaveStats = map[string]func(value interface{}) int64{
	"Seconds_Behind_Master": slaveSeconds,
	"Slave_SQL_Running":     slaveRunning,
	"Slave_IO_Running":      slaveRunning,
}

// hasReplPriv returns if the current user has privileges for the "SHOW SLAVE STATUS" query.
func (m *MySQL) hasReplPriv() bool {
	if m.user == "" {
		return false
	}

	rows, err := m.db.Query("select User,Host,Super_priv,Repl_client_priv from mysql.user")
	if err != nil {
		return false
	}

	defer rows.Close()

	for rows.Next() {
		var (
			user           string
			superPriv      string // Y or N
			replClientPriv string // Y or N
		)

		err = rows.Scan(&user, &superPriv, &replClientPriv)
		if err != nil {
			return false
		}

		// either  SUPER or REPLICATION CLIENT.
		if user == m.user && (superPriv == "Y" || replClientPriv == "Y") {
			return true
		}
	}

	return false

}

// https://dev.mysql.com/doc/refman/8.0/en/show-slave-status.html
func (m *MySQL) collectSlaveStatus(metrics map[string]int64) error {
	if !m.hasReplPriv() {
		return nil
	}

	rows, err := m.db.Query("SHOW SLAVE STATUS")
	if err != nil {
		return err
	}

	defer rows.Close()

	for rows.Next() {
		var (
			varName string
			value   interface{}
		)

		err = rows.Scan(&varName, &value)
		if err != nil {
			return err
		}

		// fmt.Printf("var name: %s | value: %s\n", varName, value)

		if valFn, ok := slaveStats[varName]; ok { // hmmm...
			value := valFn(value)
			metrics[toKey(varName)] = value
		}
	}

	return nil
}

// https://dev.mysql.com/doc/refman/8.0/en/show-variables.html
func (m *MySQL) collectMaxConnections(metrics map[string]int64) error {
	rows, err := m.db.Query("SHOW GLOBAL VARIABLES LIKE 'max_connections'") // only one result, i.e "max_conections" = 151
	if err != nil {
		return err
	}

	defer rows.Close()

	if !rows.Next() {
		return nil
	}

	var (
		varName string
		value   int64
	)

	err = rows.Scan(&varName, &value)
	if err != nil {
		return err
	}

	metrics["max_connections"] = value

	return nil

	/*
		[max_connections] = [151]
	*/
}
