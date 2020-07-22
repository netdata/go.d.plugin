package mysql

import (
	"bufio"
	"bytes"
	"database/sql/driver"
	"errors"
	"fmt"
	"io/ioutil"
	"strings"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/netdata/go-orchestrator/module"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	globalStatusGaleraMariaDBv1054, _    = ioutil.ReadFile("testdata/MariaDBv10.5.4-galera-[global_status].txt")
	globalVariablesGaleraMariaDBv1054, _ = ioutil.ReadFile("testdata/MariaDBv10.5.4-galera-[global_variables].txt")
	userStatisticsGaleraMariaDBv1054, _  = ioutil.ReadFile("testdata/MariaDBv10.5.4-galera-[user_statistics].txt")
	slaveStatusSingleSrcMariaDBv1054, _  = ioutil.ReadFile("testdata/MariaDBv10.5.4-single-source-[slave_status].txt")

	globalStatusMySQLv8021, _         = ioutil.ReadFile("testdata/MySQLv8.0.21-[global_status].txt")
	globalVariablesMySQLv8021, _      = ioutil.ReadFile("testdata/MySQLv8.0.21-[global_variables].txt")
	slaveStatusSingleSrcMySQLv8021, _ = ioutil.ReadFile("testdata/MySQLv8.0.21-single-source-[slave_status].txt")
	slaveStatusMultiSrcMySQLv8021, _  = ioutil.ReadFile("testdata/MySQLv8.0.21-multi-source-[slave_status].txt")
)

var (
	errSQLSyntax = errors.New("you have an error in your SQL syntax")
)

func Test_testDataIsCorrectlyReadAndValid(t *testing.T) {
	for name, data := range map[string][]byte{
		"globalStatusGaleraMariaDBv1054":    globalStatusGaleraMariaDBv1054,
		"globalVariablesGaleraMariaDBv1054": globalVariablesGaleraMariaDBv1054,
		"userStatisticsGaleraMariaDBv1054":  userStatisticsGaleraMariaDBv1054,
		"slaveStatusSingleSrcMariaDBv1054":  slaveStatusSingleSrcMariaDBv1054,
		"globalStatusMySQLv8021":            globalStatusMySQLv8021,
		"globalVariablesMySQLv8021":         globalVariablesMySQLv8021,
		"slaveStatusSingleSrcMySQLv8021":    slaveStatusSingleSrcMySQLv8021,
		"slaveStatusMultiSrcMySQLv8021":     slaveStatusMultiSrcMySQLv8021,
	} {
		require.NotNilf(t, data, fmt.Sprintf("read data: %s", name))
		_, err := prepareMockRows(data)
		require.NoErrorf(t, err, fmt.Sprintf("prepare mock rows: %s", name))
	}
}

func TestNew(t *testing.T) {
	assert.Implements(t, (*module.Module)(nil), New())
}

func TestMySQL_Init(t *testing.T) {
	tests := map[string]struct {
		config   Config
		wantFail bool
	}{
		"empty DSN": {
			config:   Config{DSN: ""},
			wantFail: true,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			mySQL := New()
			mySQL.Config = test.config

			if test.wantFail {
				assert.False(t, mySQL.Init())
			} else {
				assert.True(t, mySQL.Init())
			}
		})
	}
}

func TestMySQL_Cleanup(t *testing.T) {
	tests := map[string]func(t *testing.T) (mySQL *MySQL, cleanup func()){
		"db connection not initialized": func(t *testing.T) (mySQL *MySQL, cleanup func()) {
			return New(), func() {}
		},
		"db connection initialized": func(t *testing.T) (mySQL *MySQL, cleanup func()) {
			db, mock, err := sqlmock.New()
			require.NoError(t, err)

			mock.ExpectClose()
			mySQL = New()
			mySQL.db = db
			cleanup = func() { _ = db.Close() }

			return mySQL, cleanup
		},
	}

	for name, prepare := range tests {
		t.Run(name, func(t *testing.T) {
			mySQL, cleanup := prepare(t)
			defer cleanup()

			assert.NotPanics(t, mySQL.Cleanup)
			assert.Nil(t, mySQL.db)
		})
	}
}

func TestMySQL_Charts(t *testing.T) {
	assert.NotNil(t, New().Charts())
}

func TestMySQL_Check(t *testing.T) {
	tests := map[string]struct {
		prepare   func(t *testing.T) (mySQL *MySQL, cleanup func())
		wantFalse bool
	}{
		"all queries success": {
			prepare: func(t *testing.T) (mySQL *MySQL, cleanup func()) {
				db, mock, err := sqlmock.New()
				require.NoError(t, err)
				mySQL = New()
				mySQL.db = db
				cleanup = func() { _ = db.Close() }

				mock.ExpectQuery(queryGlobalStatus).
					WillReturnRows(mustMockRows(t, globalStatusGaleraMariaDBv1054))
				mock.ExpectQuery(queryGlobalVariables).
					WillReturnRows(mustMockRows(t, globalVariablesGaleraMariaDBv1054))
				mock.ExpectQuery(querySlaveStatus).
					WillReturnRows(mustMockRows(t, slaveStatusSingleSrcMariaDBv1054))
				mock.ExpectQuery(queryUserStatistics).
					WillReturnRows(mustMockRows(t, userStatisticsGaleraMariaDBv1054))

				return mySQL, cleanup
			},
		},
		"'SHOW SLAVE STATUS' fails": {
			prepare: func(t *testing.T) (mySQL *MySQL, cleanup func()) {
				db, mock, err := sqlmock.New()
				require.NoError(t, err)
				mySQL = New()
				mySQL.db = db
				cleanup = func() { _ = db.Close() }

				mock.ExpectQuery(queryGlobalStatus).
					WillReturnRows(mustMockRows(t, globalStatusGaleraMariaDBv1054))
				mock.ExpectQuery(queryGlobalVariables).
					WillReturnRows(mustMockRows(t, globalVariablesGaleraMariaDBv1054))
				mock.ExpectQuery(querySlaveStatus).
					WillReturnError(errSQLSyntax)
				mock.ExpectQuery(queryUserStatistics).
					WillReturnRows(mustMockRows(t, userStatisticsGaleraMariaDBv1054))

				return mySQL, cleanup
			},
		},
		"'SHOW USER_STATISTICS' fails": {
			prepare: func(t *testing.T) (mySQL *MySQL, cleanup func()) {
				db, mock, err := sqlmock.New()
				require.NoError(t, err)
				mySQL = New()
				mySQL.db = db
				cleanup = func() { _ = db.Close() }

				mock.ExpectQuery(queryGlobalStatus).
					WillReturnRows(mustMockRows(t, globalStatusGaleraMariaDBv1054))
				mock.ExpectQuery(queryGlobalVariables).
					WillReturnRows(mustMockRows(t, globalVariablesGaleraMariaDBv1054))
				mock.ExpectQuery(querySlaveStatus).
					WillReturnRows(mustMockRows(t, slaveStatusSingleSrcMariaDBv1054))
				mock.ExpectQuery(queryUserStatistics).
					WillReturnError(errSQLSyntax)

				return mySQL, cleanup
			},
		},
		"'SHOW GLOBAL STATUS' fails": {
			prepare: func(t *testing.T) (mySQL *MySQL, cleanup func()) {
				db, mock, err := sqlmock.New()
				require.NoError(t, err)
				mySQL = New()
				mySQL.db = db
				cleanup = func() { _ = db.Close() }

				mock.ExpectQuery(queryGlobalStatus).
					WillReturnError(errSQLSyntax)

				return mySQL, cleanup
			},
			wantFalse: true,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			mySQL, cleanup := test.prepare(t)
			defer cleanup()

			if test.wantFalse {
				assert.False(t, mySQL.Check())
			} else {
				assert.True(t, mySQL.Check())
			}
		})
	}
}

func TestMySQL_Collect(t *testing.T) {
	tests := map[string]struct {
		prepare  func(t *testing.T) (mySQL *MySQL, mock sqlmock.Sqlmock, cleanup func())
		expected map[string]int64
	}{
		"MariaDBv10.5.4-galera: all queries (single source replication)": {
			prepare: func(t *testing.T) (mySQL *MySQL, mock sqlmock.Sqlmock, cleanup func()) {
				db, mock, err := sqlmock.New()
				require.NoError(t, err)
				mySQL = New()
				mySQL.db = db
				cleanup = func() { _ = db.Close() }

				mock.ExpectQuery(queryGlobalStatus).
					WillReturnRows(mustMockRows(t, globalStatusGaleraMariaDBv1054))
				mock.ExpectQuery(queryGlobalVariables).
					WillReturnRows(mustMockRows(t, globalVariablesGaleraMariaDBv1054))
				mock.ExpectQuery(querySlaveStatus).
					WillReturnRows(mustMockRows(t, slaveStatusSingleSrcMariaDBv1054))
				mock.ExpectQuery(queryUserStatistics).
					WillReturnRows(mustMockRows(t, userStatisticsGaleraMariaDBv1054))

				return mySQL, mock, cleanup
			},
			expected: map[string]int64{
				"Aborted_connects":                      0,
				"Binlog_cache_disk_use":                 0,
				"Binlog_cache_use":                      0,
				"Binlog_stmt_cache_disk_use":            0,
				"Binlog_stmt_cache_use":                 0,
				"Bytes_received":                        1001,
				"Bytes_sent":                            195182,
				"Com_delete":                            0,
				"Com_insert":                            0,
				"Com_replace":                           0,
				"Com_select":                            3,
				"Com_update":                            0,
				"Connection_errors_accept":              0,
				"Connection_errors_internal":            0,
				"Connection_errors_max_connections":     0,
				"Connection_errors_peer_address":        0,
				"Connection_errors_select":              0,
				"Connection_errors_tcpwrap":             0,
				"Connections":                           13,
				"Created_tmp_disk_tables":               0,
				"Created_tmp_files":                     5,
				"Created_tmp_tables":                    12,
				"Handler_commit":                        26,
				"Handler_delete":                        0,
				"Handler_prepare":                       0,
				"Handler_read_first":                    7,
				"Handler_read_key":                      7,
				"Handler_read_next":                     3,
				"Handler_read_prev":                     0,
				"Handler_read_rnd":                      0,
				"Handler_read_rnd_next":                 5201,
				"Handler_rollback":                      1,
				"Handler_savepoint":                     0,
				"Handler_savepoint_rollback":            0,
				"Handler_update":                        6,
				"Handler_write":                         1,
				"Innodb_buffer_pool_bytes_data":         5357568,
				"Innodb_buffer_pool_bytes_dirty":        0,
				"Innodb_buffer_pool_pages_data":         327,
				"Innodb_buffer_pool_pages_dirty":        0,
				"Innodb_buffer_pool_pages_flushed":      134,
				"Innodb_buffer_pool_pages_free":         7727,
				"Innodb_buffer_pool_pages_misc":         0,
				"Innodb_buffer_pool_pages_total":        8054,
				"Innodb_buffer_pool_read_ahead":         0,
				"Innodb_buffer_pool_read_ahead_evicted": 0,
				"Innodb_buffer_pool_read_ahead_rnd":     0,
				"Innodb_buffer_pool_read_requests":      2369,
				"Innodb_buffer_pool_reads":              196,
				"Innodb_buffer_pool_wait_free":          0,
				"Innodb_buffer_pool_write_requests":     853,
				"Innodb_data_fsyncs":                    25,
				"Innodb_data_pending_fsyncs":            0,
				"Innodb_data_pending_reads":             0,
				"Innodb_data_pending_writes":            0,
				"Innodb_data_read":                      3211264,
				"Innodb_data_reads":                     215,
				"Innodb_data_writes":                    157,
				"Innodb_data_written":                   2244608,
				"Innodb_deadlocks":                      0,
				"Innodb_log_waits":                      0,
				"Innodb_log_write_requests":             0,
				"Innodb_log_writes":                     20,
				"Innodb_os_log_fsyncs":                  20,
				"Innodb_os_log_pending_fsyncs":          0,
				"Innodb_os_log_pending_writes":          0,
				"Innodb_os_log_written":                 10240,
				"Innodb_row_lock_current_waits":         0,
				"Innodb_rows_deleted":                   0,
				"Innodb_rows_inserted":                  0,
				"Innodb_rows_read":                      0,
				"Innodb_rows_updated":                   0,
				"Key_blocks_not_flushed":                0,
				"Key_blocks_unused":                     107163,
				"Key_blocks_used":                       0,
				"Key_read_requests":                     0,
				"Key_reads":                             0,
				"Key_write_requests":                    0,
				"Key_writes":                            0,
				"Max_used_connections":                  1,
				"Open_files":                            24,
				"Opened_files":                          80,
				"Qcache_free_blocks":                    1,
				"Qcache_free_memory":                    1031304,
				"Qcache_hits":                           0,
				"Qcache_inserts":                        0,
				"Qcache_lowmem_prunes":                  0,
				"Qcache_not_cached":                     0,
				"Qcache_queries_in_cache":               0,
				"Qcache_total_blocks":                   1,
				"Queries":                               32,
				"Questions":                             24,
				"Seconds_Behind_Master":                 0,
				"Select_full_join":                      0,
				"Select_full_range_join":                0,
				"Select_range":                          0,
				"Select_range_check":                    0,
				"Select_scan":                           12,
				"Slave_IO_Running":                      1,
				"Slave_SQL_Running":                     1,
				"Slow_queries":                          0,
				"Sort_merge_passes":                     0,
				"Sort_range":                            0,
				"Sort_scan":                             0,
				"Table_locks_immediate":                 59,
				"Table_locks_waited":                    0,
				"Threads_cached":                        0,
				"Threads_connected":                     1,
				"Threads_created":                       6,
				"Threads_running":                       1,
				"max_connections":                       151,
				"table_open_cache":                      2000,
				"userstats_netdata_Other_commands":      0,
				"userstats_netdata_Rows_deleted":        0,
				"userstats_netdata_Rows_inserted":       0,
				"userstats_netdata_Rows_read":           0,
				"userstats_netdata_Rows_sent":           15,
				"userstats_netdata_Rows_updated":        0,
				"userstats_netdata_Select_commands":     1,
				"userstats_netdata_Update_commands":     0,
				"userstats_root_Other_commands":         1,
				"userstats_root_Rows_deleted":           0,
				"userstats_root_Rows_inserted":          1,
				"userstats_root_Rows_read":              17,
				"userstats_root_Rows_sent":              4541,
				"userstats_root_Rows_updated":           3,
				"userstats_root_Select_commands":        2,
				"userstats_root_Update_commands":        4,
				"wsrep_cluster_size":                    2,
				"wsrep_cluster_status":                  0,
				"wsrep_cluster_weight":                  2,
				"wsrep_connected":                       1,
				"wsrep_flow_control_paused_ns":          0,
				"wsrep_local_bf_aborts":                 0,
				"wsrep_local_cert_failures":             0,
				"wsrep_local_recv_queue":                0,
				"wsrep_local_send_queue":                0,
				"wsrep_local_state":                     4,
				"wsrep_open_transactions":               0,
				"wsrep_ready":                           1,
				"wsrep_received":                        8,
				"wsrep_received_bytes":                  2608,
				"wsrep_replicated":                      5,
				"wsrep_replicated_bytes":                2392,
				"wsrep_thread_count":                    5,
			},
		},
		"MariaDBv10.5.4-galera: minimal: global status and variables": {
			prepare: func(t *testing.T) (mySQL *MySQL, mock sqlmock.Sqlmock, cleanup func()) {
				db, mock, err := sqlmock.New()
				require.NoError(t, err)
				mySQL = New()
				mySQL.db = db
				cleanup = func() { _ = db.Close() }

				mock.ExpectQuery(queryGlobalStatus).
					WillReturnRows(mustMockRows(t, globalStatusGaleraMariaDBv1054))
				mock.ExpectQuery(queryGlobalVariables).
					WillReturnRows(mustMockRows(t, globalVariablesGaleraMariaDBv1054))
				mock.ExpectQuery(querySlaveStatus).
					WillReturnError(errSQLSyntax)
				mock.ExpectQuery(queryUserStatistics).
					WillReturnError(errSQLSyntax)

				return mySQL, mock, cleanup
			},
			expected: map[string]int64{
				"Aborted_connects":                      0,
				"Binlog_cache_disk_use":                 0,
				"Binlog_cache_use":                      0,
				"Binlog_stmt_cache_disk_use":            0,
				"Binlog_stmt_cache_use":                 0,
				"Bytes_received":                        1001,
				"Bytes_sent":                            195182,
				"Com_delete":                            0,
				"Com_insert":                            0,
				"Com_replace":                           0,
				"Com_select":                            3,
				"Com_update":                            0,
				"Connection_errors_accept":              0,
				"Connection_errors_internal":            0,
				"Connection_errors_max_connections":     0,
				"Connection_errors_peer_address":        0,
				"Connection_errors_select":              0,
				"Connection_errors_tcpwrap":             0,
				"Connections":                           13,
				"Created_tmp_disk_tables":               0,
				"Created_tmp_files":                     5,
				"Created_tmp_tables":                    12,
				"Handler_commit":                        26,
				"Handler_delete":                        0,
				"Handler_prepare":                       0,
				"Handler_read_first":                    7,
				"Handler_read_key":                      7,
				"Handler_read_next":                     3,
				"Handler_read_prev":                     0,
				"Handler_read_rnd":                      0,
				"Handler_read_rnd_next":                 5201,
				"Handler_rollback":                      1,
				"Handler_savepoint":                     0,
				"Handler_savepoint_rollback":            0,
				"Handler_update":                        6,
				"Handler_write":                         1,
				"Innodb_buffer_pool_bytes_data":         5357568,
				"Innodb_buffer_pool_bytes_dirty":        0,
				"Innodb_buffer_pool_pages_data":         327,
				"Innodb_buffer_pool_pages_dirty":        0,
				"Innodb_buffer_pool_pages_flushed":      134,
				"Innodb_buffer_pool_pages_free":         7727,
				"Innodb_buffer_pool_pages_misc":         0,
				"Innodb_buffer_pool_pages_total":        8054,
				"Innodb_buffer_pool_read_ahead":         0,
				"Innodb_buffer_pool_read_ahead_evicted": 0,
				"Innodb_buffer_pool_read_ahead_rnd":     0,
				"Innodb_buffer_pool_read_requests":      2369,
				"Innodb_buffer_pool_reads":              196,
				"Innodb_buffer_pool_wait_free":          0,
				"Innodb_buffer_pool_write_requests":     853,
				"Innodb_data_fsyncs":                    25,
				"Innodb_data_pending_fsyncs":            0,
				"Innodb_data_pending_reads":             0,
				"Innodb_data_pending_writes":            0,
				"Innodb_data_read":                      3211264,
				"Innodb_data_reads":                     215,
				"Innodb_data_writes":                    157,
				"Innodb_data_written":                   2244608,
				"Innodb_deadlocks":                      0,
				"Innodb_log_waits":                      0,
				"Innodb_log_write_requests":             0,
				"Innodb_log_writes":                     20,
				"Innodb_os_log_fsyncs":                  20,
				"Innodb_os_log_pending_fsyncs":          0,
				"Innodb_os_log_pending_writes":          0,
				"Innodb_os_log_written":                 10240,
				"Innodb_row_lock_current_waits":         0,
				"Innodb_rows_deleted":                   0,
				"Innodb_rows_inserted":                  0,
				"Innodb_rows_read":                      0,
				"Innodb_rows_updated":                   0,
				"Key_blocks_not_flushed":                0,
				"Key_blocks_unused":                     107163,
				"Key_blocks_used":                       0,
				"Key_read_requests":                     0,
				"Key_reads":                             0,
				"Key_write_requests":                    0,
				"Key_writes":                            0,
				"Max_used_connections":                  1,
				"Open_files":                            24,
				"Opened_files":                          80,
				"Qcache_free_blocks":                    1,
				"Qcache_free_memory":                    1031304,
				"Qcache_hits":                           0,
				"Qcache_inserts":                        0,
				"Qcache_lowmem_prunes":                  0,
				"Qcache_not_cached":                     0,
				"Qcache_queries_in_cache":               0,
				"Qcache_total_blocks":                   1,
				"Queries":                               32,
				"Questions":                             24,
				"Select_full_join":                      0,
				"Select_full_range_join":                0,
				"Select_range":                          0,
				"Select_range_check":                    0,
				"Select_scan":                           12,
				"Slow_queries":                          0,
				"Sort_merge_passes":                     0,
				"Sort_range":                            0,
				"Sort_scan":                             0,
				"Table_locks_immediate":                 59,
				"Table_locks_waited":                    0,
				"Threads_cached":                        0,
				"Threads_connected":                     1,
				"Threads_created":                       6,
				"Threads_running":                       1,
				"max_connections":                       151,
				"table_open_cache":                      2000,
				"wsrep_cluster_size":                    2,
				"wsrep_cluster_status":                  0,
				"wsrep_cluster_weight":                  2,
				"wsrep_connected":                       1,
				"wsrep_flow_control_paused_ns":          0,
				"wsrep_local_bf_aborts":                 0,
				"wsrep_local_cert_failures":             0,
				"wsrep_local_recv_queue":                0,
				"wsrep_local_send_queue":                0,
				"wsrep_local_state":                     4,
				"wsrep_open_transactions":               0,
				"wsrep_ready":                           1,
				"wsrep_received":                        8,
				"wsrep_received_bytes":                  2608,
				"wsrep_replicated":                      5,
				"wsrep_replicated_bytes":                2392,
				"wsrep_thread_count":                    5,
			},
		},
		"MySQLv8.0.21: all queries (multi source replication)": {
			prepare: func(t *testing.T) (mySQL *MySQL, mock sqlmock.Sqlmock, cleanup func()) {
				db, mock, err := sqlmock.New()
				require.NoError(t, err)
				mySQL = New()
				mySQL.db = db
				cleanup = func() { _ = db.Close() }

				mock.ExpectQuery(queryGlobalStatus).
					WillReturnRows(mustMockRows(t, globalStatusMySQLv8021))
				mock.ExpectQuery(queryGlobalVariables).
					WillReturnRows(mustMockRows(t, globalVariablesMySQLv8021))
				mock.ExpectQuery(querySlaveStatus).
					WillReturnRows(mustMockRows(t, slaveStatusMultiSrcMySQLv8021))
				mock.ExpectQuery(queryUserStatistics).
					WillReturnError(errSQLSyntax)

				return mySQL, mock, cleanup
			},
			expected: map[string]int64{
				"Aborted_connects":                      0,
				"Binlog_cache_disk_use":                 0,
				"Binlog_cache_use":                      2,
				"Binlog_stmt_cache_disk_use":            0,
				"Binlog_stmt_cache_use":                 0,
				"Bytes_received":                        13552,
				"Bytes_sent":                            21281,
				"Com_delete":                            0,
				"Com_insert":                            0,
				"Com_replace":                           0,
				"Com_select":                            3,
				"Com_update":                            0,
				"Connection_errors_accept":              0,
				"Connection_errors_internal":            0,
				"Connection_errors_max_connections":     0,
				"Connection_errors_peer_address":        0,
				"Connection_errors_select":              0,
				"Connection_errors_tcpwrap":             0,
				"Connections":                           67,
				"Created_tmp_disk_tables":               0,
				"Created_tmp_files":                     5,
				"Created_tmp_tables":                    2,
				"Handler_commit":                        552,
				"Handler_delete":                        0,
				"Handler_prepare":                       8,
				"Handler_read_first":                    34,
				"Handler_read_key":                      1635,
				"Handler_read_next":                     3891,
				"Handler_read_prev":                     0,
				"Handler_read_rnd":                      0,
				"Handler_read_rnd_next":                 1011,
				"Handler_rollback":                      0,
				"Handler_savepoint":                     0,
				"Handler_savepoint_rollback":            0,
				"Handler_update":                        316,
				"Handler_write":                         467,
				"Innodb_buffer_pool_bytes_data":         15761408,
				"Innodb_buffer_pool_bytes_dirty":        0,
				"Innodb_buffer_pool_pages_data":         962,
				"Innodb_buffer_pool_pages_dirty":        0,
				"Innodb_buffer_pool_pages_flushed":      170,
				"Innodb_buffer_pool_pages_free":         7226,
				"Innodb_buffer_pool_pages_misc":         4,
				"Innodb_buffer_pool_pages_total":        8192,
				"Innodb_buffer_pool_read_ahead":         0,
				"Innodb_buffer_pool_read_ahead_evicted": 0,
				"Innodb_buffer_pool_read_ahead_rnd":     0,
				"Innodb_buffer_pool_read_requests":      14452,
				"Innodb_buffer_pool_reads":              818,
				"Innodb_buffer_pool_wait_free":          0,
				"Innodb_buffer_pool_write_requests":     1696,
				"Innodb_data_fsyncs":                    76,
				"Innodb_data_pending_fsyncs":            0,
				"Innodb_data_pending_reads":             0,
				"Innodb_data_pending_writes":            0,
				"Innodb_data_read":                      13472768,
				"Innodb_data_reads":                     840,
				"Innodb_data_writes":                    252,
				"Innodb_data_written":                   3002368,
				"Innodb_log_waits":                      0,
				"Innodb_log_write_requests":             664,
				"Innodb_log_writes":                     26,
				"Innodb_os_log_fsyncs":                  25,
				"Innodb_os_log_pending_fsyncs":          0,
				"Innodb_os_log_pending_writes":          0,
				"Innodb_os_log_written":                 38912,
				"Innodb_row_lock_current_waits":         0,
				"Innodb_rows_deleted":                   0,
				"Innodb_rows_inserted":                  0,
				"Innodb_rows_read":                      0,
				"Innodb_rows_updated":                   0,
				"Key_blocks_not_flushed":                0,
				"Key_blocks_unused":                     6698,
				"Key_blocks_used":                       0,
				"Key_read_requests":                     0,
				"Key_reads":                             0,
				"Key_write_requests":                    0,
				"Key_writes":                            0,
				"Max_used_connections":                  1,
				"Open_files":                            2,
				"Opened_files":                          2,
				"Queries":                               125,
				"Questions":                             67,
				"Seconds_Behind_Master_master1":         0,
				"Seconds_Behind_Master_master2":         0,
				"Select_full_join":                      0,
				"Select_full_range_join":                0,
				"Select_range":                          0,
				"Select_range_check":                    0,
				"Select_scan":                           4,
				"Slave_IO_Running_master1":              2,
				"Slave_IO_Running_master2":              2,
				"Slave_SQL_Running_master1":             1,
				"Slave_SQL_Running_master2":             1,
				"Slow_queries":                          0,
				"Sort_merge_passes":                     0,
				"Sort_range":                            0,
				"Sort_scan":                             0,
				"Table_locks_immediate":                 2,
				"Table_locks_waited":                    0,
				"Threads_cached":                        0,
				"Threads_connected":                     1,
				"Threads_created":                       1,
				"Threads_running":                       2,
				"max_connections":                       151,
				"table_open_cache":                      4000,
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			mySQL, mock, cleanup := test.prepare(t)
			defer cleanup()

			collected := mySQL.Collect()
			assert.Equal(t, test.expected, collected)
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func mustMockRows(t *testing.T, data []byte) *sqlmock.Rows {
	rows, err := prepareMockRows(data)
	require.NoError(t, err)
	return rows
}

func prepareMockRows(data []byte) (*sqlmock.Rows, error) {
	r := bytes.NewReader(data)
	sc := bufio.NewScanner(r)

	set := make(map[string]bool)
	var columns []string
	var lines [][]driver.Value
	var values []driver.Value

	for sc.Scan() {
		text := strings.TrimSpace(sc.Text())
		if text == "" {
			continue
		}
		if isNewRow := text[0] == '*'; isNewRow {
			if len(values) != 0 {
				lines = append(lines, values)
				values = []driver.Value{}
			}
			continue
		}

		idx := strings.IndexByte(text, ':')
		// not interested in multi line values
		if idx == -1 {
			continue
		}

		name := strings.TrimSpace(text[:idx])
		value := strings.TrimSpace(text[idx+1:])
		if !set[name] {
			set[name] = true
			columns = append(columns, name)
		}
		values = append(values, value)
	}
	if len(values) != 0 {
		lines = append(lines, values)
	}

	rows := sqlmock.NewRows(columns)
	for _, values := range lines {
		if len(columns) != len(values) {
			return nil, fmt.Errorf("columns != values (%d/%d)", len(columns), len(values))
		}
		rows.AddRow(values...)
	}
	return rows, nil
}
