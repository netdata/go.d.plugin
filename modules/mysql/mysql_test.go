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

	"github.com/netdata/go.d.plugin/agent/module"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/blang/semver/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	mariaV5546Version, _         = ioutil.ReadFile("testdata/mariadb/v5.5.46/version.txt")
	mariaV5546GlobalStatus, _    = ioutil.ReadFile("testdata/mariadb/v5.5.46/global_status.txt")
	mariaV5546GlobalVariables, _ = ioutil.ReadFile("testdata/mariadb/v5.5.46/global_variables.txt")
	mariaV5546SlaveStatus, _     = ioutil.ReadFile("testdata/mariadb/v5.5.46/slave_status.txt")

	mariaV1054Version, _         = ioutil.ReadFile("testdata/mariadb/v10.5.4/version.txt")
	mariaV1054GlobalStatus, _    = ioutil.ReadFile("testdata/mariadb/v10.5.4/global_status.txt")
	mariaV1054GlobalVariables, _ = ioutil.ReadFile("testdata/mariadb/v10.5.4/global_variables.txt")
	mariaV1054UserStatistics, _  = ioutil.ReadFile("testdata/mariadb/v10.5.4/user_statistics.txt")
	mariaV1054AllSlavesStatus, _ = ioutil.ReadFile("testdata/mariadb/v10.5.4/all_slaves_status.txt")

	mysqlV8021Version, _         = ioutil.ReadFile("testdata/mysql/v8.0.21/version.txt")
	mysqlV8021GlobalStatus, _    = ioutil.ReadFile("testdata/mysql/v8.0.21/global_status.txt")
	mysqlV8021GlobalVariables, _ = ioutil.ReadFile("testdata/mysql/v8.0.21/global_variables.txt")
	mysqlV8021SlaveStatus, _     = ioutil.ReadFile("testdata/mysql/v8.0.21/slave_status.txt")
)

var (
	errSQLSyntax = errors.New("you have an error in your SQL syntax")
)

func Test_testDataIsValid(t *testing.T) {
	for name, data := range map[string][]byte{
		"mariaV5546Version":         mariaV5546Version,
		"mariaV5546GlobalStatus":    mariaV5546GlobalStatus,
		"mariaV5546GlobalVariables": mariaV5546GlobalVariables,
		"mariaV5546SlaveStatus":     mariaV5546SlaveStatus,

		"mariaV1054Version":         mariaV1054Version,
		"mariaV1054GlobalStatus":    mariaV1054GlobalStatus,
		"mariaV1054GlobalVariables": mariaV1054GlobalVariables,
		"mariaV1054UserStatistics":  mariaV1054UserStatistics,
		"mariaV1054AllSlavesStatus": mariaV1054AllSlavesStatus,

		"mysqlV8021Version":         mysqlV8021Version,
		"mysqlV8021GlobalStatus":    mysqlV8021GlobalStatus,
		"mysqlV8021GlobalVariables": mysqlV8021GlobalVariables,
		"mysqlV8021SlaveStatus":     mysqlV8021SlaveStatus,
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
		prepare   func(t *testing.T) (mySQL *MySQL, mock sqlmock.Sqlmock, cleanup func())
		wantFalse bool
	}{
		"all queries success (MariaDB)": {
			prepare: func(t *testing.T) (mySQL *MySQL, mock sqlmock.Sqlmock, cleanup func()) {
				db, mock, err := sqlmock.New()
				require.NoError(t, err)
				mySQL = New()
				mySQL.db = db
				cleanup = func() { _ = db.Close() }

				mock.ExpectQuery(queryVersion).
					WillReturnRows(mustMockRows(t, mariaV1054Version))
				mock.ExpectQuery(queryGlobalStatus).
					WillReturnRows(mustMockRows(t, mariaV1054GlobalStatus))
				mock.ExpectQuery(queryGlobalVariables).
					WillReturnRows(mustMockRows(t, mariaV1054GlobalVariables))
				mock.ExpectQuery(queryAllSlavesStatus).
					WillReturnRows(mustMockRows(t, mariaV1054AllSlavesStatus))
				mock.ExpectQuery(queryUserStatistics).
					WillReturnRows(mustMockRows(t, mariaV1054UserStatistics))

				return mySQL, mock, cleanup
			},
		},
		"'SELECT VERSION()' fails": {
			wantFalse: true,
			prepare: func(t *testing.T) (mySQL *MySQL, mock sqlmock.Sqlmock, cleanup func()) {
				db, mock, err := sqlmock.New()
				require.NoError(t, err)
				mySQL = New()
				mySQL.db = db
				cleanup = func() { _ = db.Close() }

				mock.ExpectQuery(queryVersion).
					WillReturnError(errSQLSyntax)

				return mySQL, mock, cleanup
			},
		},
		"'SHOW GLOBAL STATUS' fails": {
			wantFalse: true,
			prepare: func(t *testing.T) (mySQL *MySQL, mock sqlmock.Sqlmock, cleanup func()) {
				db, mock, err := sqlmock.New()
				require.NoError(t, err)
				mySQL = New()
				mySQL.db = db
				cleanup = func() { _ = db.Close() }

				mock.ExpectQuery(queryVersion).
					WillReturnRows(mustMockRows(t, mariaV1054Version))
				mock.ExpectQuery(queryGlobalStatus).
					WillReturnError(errSQLSyntax)

				return mySQL, mock, cleanup
			},
		},
		"'SHOW GLOBAL VARIABLES' fails": {
			wantFalse: true,
			prepare: func(t *testing.T) (mySQL *MySQL, mock sqlmock.Sqlmock, cleanup func()) {
				db, mock, err := sqlmock.New()
				require.NoError(t, err)
				mySQL = New()
				mySQL.db = db
				cleanup = func() { _ = db.Close() }

				mock.ExpectQuery(queryVersion).
					WillReturnRows(mustMockRows(t, mariaV1054Version))
				mock.ExpectQuery(queryGlobalStatus).
					WillReturnRows(mustMockRows(t, mariaV1054GlobalStatus))
				mock.ExpectQuery(queryGlobalVariables).
					WillReturnError(errSQLSyntax)

				return mySQL, mock, cleanup
			},
		},
		"'SHOW ALL SLAVES STATUS' fails (MariaDB)": {
			prepare: func(t *testing.T) (mySQL *MySQL, mock sqlmock.Sqlmock, cleanup func()) {
				db, mock, err := sqlmock.New()
				require.NoError(t, err)
				mySQL = New()
				mySQL.db = db
				cleanup = func() { _ = db.Close() }

				mock.ExpectQuery(queryVersion).
					WillReturnRows(mustMockRows(t, mariaV1054Version))
				mock.ExpectQuery(queryGlobalStatus).
					WillReturnRows(mustMockRows(t, mariaV1054GlobalStatus))
				mock.ExpectQuery(queryGlobalVariables).
					WillReturnRows(mustMockRows(t, mariaV1054GlobalVariables))
				mock.ExpectQuery(queryAllSlavesStatus).
					WillReturnError(errSQLSyntax)
				mock.ExpectQuery(queryUserStatistics).
					WillReturnRows(mustMockRows(t, mariaV1054UserStatistics))

				return mySQL, mock, cleanup
			},
		},
		"'SHOW USER_STATISTICS' fails (MariaDB)": {
			prepare: func(t *testing.T) (mySQL *MySQL, mock sqlmock.Sqlmock, cleanup func()) {
				db, mock, err := sqlmock.New()
				require.NoError(t, err)
				mySQL = New()
				mySQL.db = db
				cleanup = func() { _ = db.Close() }

				mock.ExpectQuery(queryVersion).
					WillReturnRows(mustMockRows(t, mariaV1054Version))
				mock.ExpectQuery(queryGlobalStatus).
					WillReturnRows(mustMockRows(t, mariaV1054GlobalStatus))
				mock.ExpectQuery(queryGlobalVariables).
					WillReturnRows(mustMockRows(t, mariaV1054GlobalVariables))
				mock.ExpectQuery(queryAllSlavesStatus).
					WillReturnRows(mustMockRows(t, mariaV1054AllSlavesStatus))
				mock.ExpectQuery(queryUserStatistics).
					WillReturnError(errSQLSyntax)

				return mySQL, mock, cleanup
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			mySQL, mock, cleanup := test.prepare(t)
			defer cleanup()

			if test.wantFalse {
				assert.False(t, mySQL.Check())
			} else {
				assert.True(t, mySQL.Check())
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestMySQL_Collect(t *testing.T) {
	tests := map[string]struct {
		prepare  func(t *testing.T) (mySQL *MySQL, mock sqlmock.Sqlmock, cleanup func())
		expected map[string]int64
	}{
		"MariaDBv5.5.46: all queries": {
			prepare: func(t *testing.T) (mySQL *MySQL, mock sqlmock.Sqlmock, cleanup func()) {
				db, mock, err := sqlmock.New()
				require.NoError(t, err)
				mySQL = New()
				mySQL.db = db
				cleanup = func() { _ = db.Close() }

				mock.ExpectQuery(queryVersion).
					WillReturnRows(mustMockRows(t, mariaV5546Version))
				mock.ExpectQuery(queryGlobalStatus).
					WillReturnRows(mustMockRows(t, mariaV5546GlobalStatus))
				mock.ExpectQuery(queryGlobalVariables).
					WillReturnRows(mustMockRows(t, mariaV5546GlobalVariables))
				mock.ExpectQuery(querySlaveStatus).
					WillReturnRows(mustMockRows(t, mariaV5546SlaveStatus))

				return mySQL, mock, cleanup
			},
			expected: map[string]int64{
				"aborted_connects":                      2,
				"binlog_cache_disk_use":                 0,
				"binlog_cache_use":                      0,
				"binlog_stmt_cache_disk_use":            0,
				"binlog_stmt_cache_use":                 1,
				"bytes_received":                        438,
				"bytes_sent":                            31090,
				"com_delete":                            0,
				"com_insert":                            0,
				"com_replace":                           0,
				"com_select":                            3,
				"com_update":                            0,
				"connections":                           5,
				"created_tmp_disk_tables":               0,
				"created_tmp_files":                     6,
				"created_tmp_tables":                    3,
				"handler_commit":                        0,
				"handler_delete":                        0,
				"handler_prepare":                       0,
				"handler_read_first":                    0,
				"handler_read_key":                      0,
				"handler_read_next":                     0,
				"handler_read_prev":                     0,
				"handler_read_rnd":                      0,
				"handler_read_rnd_next":                 834,
				"handler_rollback":                      0,
				"handler_savepoint":                     0,
				"handler_savepoint_rollback":            0,
				"handler_update":                        0,
				"handler_write":                         0,
				"innodb_buffer_pool_bytes_data":         5029888,
				"innodb_buffer_pool_bytes_dirty":        0,
				"innodb_buffer_pool_pages_data":         307,
				"innodb_buffer_pool_pages_dirty":        0,
				"innodb_buffer_pool_pages_flushed":      317,
				"innodb_buffer_pool_pages_free":         7883,
				"innodb_buffer_pool_pages_misc":         1,
				"innodb_buffer_pool_pages_total":        8191,
				"innodb_buffer_pool_read_ahead":         0,
				"innodb_buffer_pool_read_ahead_evicted": 0,
				"innodb_buffer_pool_read_ahead_rnd":     0,
				"innodb_buffer_pool_read_requests":      2576,
				"innodb_buffer_pool_reads":              0,
				"innodb_buffer_pool_wait_free":          0,
				"innodb_buffer_pool_write_requests":     2397,
				"innodb_data_fsyncs":                    19,
				"innodb_data_pending_fsyncs":            0,
				"innodb_data_pending_reads":             0,
				"innodb_data_pending_writes":            0,
				"innodb_data_read":                      0,
				"innodb_data_reads":                     0,
				"innodb_data_writes":                    359,
				"innodb_data_written":                   9131520,
				"innodb_deadlocks":                      0,
				"innodb_log_waits":                      0,
				"innodb_log_write_requests":             3177,
				"innodb_log_writes":                     4,
				"innodb_os_log_fsyncs":                  10,
				"innodb_os_log_pending_fsyncs":          0,
				"innodb_os_log_pending_writes":          0,
				"innodb_os_log_written":                 1591296,
				"innodb_row_lock_current_waits":         0,
				"innodb_rows_deleted":                   0,
				"innodb_rows_inserted":                  0,
				"innodb_rows_read":                      0,
				"innodb_rows_updated":                   0,
				"key_blocks_not_flushed":                0,
				"key_blocks_unused":                     107167,
				"key_blocks_used":                       4,
				"key_read_requests":                     14,
				"key_reads":                             4,
				"key_write_requests":                    16,
				"key_writes":                            5,
				"max_connections":                       151,
				"max_used_connections":                  1,
				"open_files":                            29,
				"open_tables":                           27,
				"opened_files":                          98,
				"opened_tables":                         0,
				"qcache_free_blocks":                    0,
				"qcache_free_memory":                    0,
				"qcache_hits":                           0,
				"qcache_inserts":                        0,
				"qcache_lowmem_prunes":                  0,
				"qcache_not_cached":                     0,
				"qcache_queries_in_cache":               0,
				"qcache_total_blocks":                   0,
				"queries":                               17,
				"questions":                             9,
				"seconds_behind_master":                 0,
				"select_full_join":                      0,
				"select_full_range_join":                0,
				"select_range":                          0,
				"select_range_check":                    0,
				"select_scan":                           3,
				"slave_io_running":                      0,
				"slave_sql_running":                     0,
				"slow_queries":                          0,
				"sort_merge_passes":                     0,
				"sort_range":                            0,
				"sort_scan":                             0,
				"table_locks_immediate":                 58,
				"table_locks_waited":                    0,
				"table_open_cache":                      400,
				"thread_cache_misses":                   6000,
				"threads_cached":                        0,
				"threads_connected":                     1,
				"threads_created":                       3,
				"threads_running":                       1,
			},
		},
		"MariaDBv10.5.4: all queries (single source replication)": {
			prepare: func(t *testing.T) (mySQL *MySQL, mock sqlmock.Sqlmock, cleanup func()) {
				db, mock, err := sqlmock.New()
				require.NoError(t, err)
				mySQL = New()
				mySQL.db = db
				cleanup = func() { _ = db.Close() }

				mock.ExpectQuery(queryVersion).
					WillReturnRows(mustMockRows(t, mariaV1054Version))
				mock.ExpectQuery(queryGlobalStatus).
					WillReturnRows(mustMockRows(t, mariaV1054GlobalStatus))
				mock.ExpectQuery(queryGlobalVariables).
					WillReturnRows(mustMockRows(t, mariaV1054GlobalVariables))
				mock.ExpectQuery(queryAllSlavesStatus).
					WillReturnRows(mustMockRows(t, mariaV1054AllSlavesStatus))
				mock.ExpectQuery(queryUserStatistics).
					WillReturnRows(mustMockRows(t, mariaV1054UserStatistics))

				return mySQL, mock, cleanup
			},
			expected: map[string]int64{
				"aborted_connects":                      0,
				"binlog_cache_disk_use":                 0,
				"binlog_cache_use":                      0,
				"binlog_stmt_cache_disk_use":            0,
				"binlog_stmt_cache_use":                 0,
				"bytes_received":                        1001,
				"bytes_sent":                            195182,
				"com_delete":                            0,
				"com_insert":                            0,
				"com_replace":                           0,
				"com_select":                            3,
				"com_update":                            0,
				"connection_errors_accept":              0,
				"connection_errors_internal":            0,
				"connection_errors_max_connections":     0,
				"connection_errors_peer_address":        0,
				"connection_errors_select":              0,
				"connection_errors_tcpwrap":             0,
				"connections":                           13,
				"created_tmp_disk_tables":               0,
				"created_tmp_files":                     5,
				"created_tmp_tables":                    12,
				"handler_commit":                        26,
				"handler_delete":                        0,
				"handler_prepare":                       0,
				"handler_read_first":                    7,
				"handler_read_key":                      7,
				"handler_read_next":                     3,
				"handler_read_prev":                     0,
				"handler_read_rnd":                      0,
				"handler_read_rnd_next":                 5201,
				"handler_rollback":                      1,
				"handler_savepoint":                     0,
				"handler_savepoint_rollback":            0,
				"handler_update":                        6,
				"handler_write":                         1,
				"innodb_buffer_pool_bytes_data":         5357568,
				"innodb_buffer_pool_bytes_dirty":        0,
				"innodb_buffer_pool_pages_data":         327,
				"innodb_buffer_pool_pages_dirty":        0,
				"innodb_buffer_pool_pages_flushed":      134,
				"innodb_buffer_pool_pages_free":         7727,
				"innodb_buffer_pool_pages_misc":         0,
				"innodb_buffer_pool_pages_total":        8054,
				"innodb_buffer_pool_read_ahead":         0,
				"innodb_buffer_pool_read_ahead_evicted": 0,
				"innodb_buffer_pool_read_ahead_rnd":     0,
				"innodb_buffer_pool_read_requests":      2369,
				"innodb_buffer_pool_reads":              196,
				"innodb_buffer_pool_wait_free":          0,
				"innodb_buffer_pool_write_requests":     853,
				"innodb_data_fsyncs":                    25,
				"innodb_data_pending_fsyncs":            0,
				"innodb_data_pending_reads":             0,
				"innodb_data_pending_writes":            0,
				"innodb_data_read":                      3211264,
				"innodb_data_reads":                     215,
				"innodb_data_writes":                    157,
				"innodb_data_written":                   2244608,
				"innodb_deadlocks":                      0,
				"innodb_log_waits":                      0,
				"innodb_log_write_requests":             0,
				"innodb_log_writes":                     20,
				"innodb_os_log_fsyncs":                  20,
				"innodb_os_log_pending_fsyncs":          0,
				"innodb_os_log_pending_writes":          0,
				"innodb_os_log_written":                 10240,
				"innodb_row_lock_current_waits":         0,
				"innodb_rows_deleted":                   0,
				"innodb_rows_inserted":                  0,
				"innodb_rows_read":                      0,
				"innodb_rows_updated":                   0,
				"key_blocks_not_flushed":                0,
				"key_blocks_unused":                     107163,
				"key_blocks_used":                       0,
				"key_read_requests":                     0,
				"key_reads":                             0,
				"key_write_requests":                    0,
				"key_writes":                            0,
				"max_connections":                       151,
				"max_used_connections":                  1,
				"open_files":                            24,
				"open_tables":                           13,
				"opened_files":                          80,
				"opened_tables":                         19,
				"qcache_free_blocks":                    1,
				"qcache_free_memory":                    1031304,
				"qcache_hits":                           0,
				"qcache_inserts":                        0,
				"qcache_lowmem_prunes":                  0,
				"qcache_not_cached":                     0,
				"qcache_queries_in_cache":               0,
				"qcache_total_blocks":                   1,
				"queries":                               32,
				"questions":                             24,
				"seconds_behind_master_master1":         0,
				"seconds_behind_master_master2":         0,
				"select_full_join":                      0,
				"select_full_range_join":                0,
				"select_range":                          0,
				"select_range_check":                    0,
				"select_scan":                           12,
				"slave_io_running_master1":              1,
				"slave_io_running_master2":              1,
				"slave_sql_running_master1":             1,
				"slave_sql_running_master2":             1,
				"slow_queries":                          0,
				"sort_merge_passes":                     0,
				"sort_range":                            0,
				"sort_scan":                             0,
				"table_locks_immediate":                 59,
				"table_locks_waited":                    0,
				"table_open_cache":                      2000,
				"thread_cache_misses":                   4615,
				"threads_cached":                        0,
				"threads_connected":                     1,
				"threads_created":                       6,
				"threads_running":                       1,
				"userstats_netdata_cpu_time":            2,
				"userstats_netdata_other_commands":      0,
				"userstats_netdata_rows_deleted":        0,
				"userstats_netdata_rows_inserted":       0,
				"userstats_netdata_rows_read":           0,
				"userstats_netdata_rows_sent":           15,
				"userstats_netdata_rows_updated":        0,
				"userstats_netdata_select_commands":     1,
				"userstats_netdata_update_commands":     0,
				"userstats_root_cpu_time":               40,
				"userstats_root_other_commands":         1,
				"userstats_root_rows_deleted":           0,
				"userstats_root_rows_inserted":          1,
				"userstats_root_rows_read":              17,
				"userstats_root_rows_sent":              4541,
				"userstats_root_rows_updated":           3,
				"userstats_root_select_commands":        2,
				"userstats_root_update_commands":        4,
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
		"MariaDBv10.5.4: minimal: global status and variables": {
			prepare: func(t *testing.T) (mySQL *MySQL, mock sqlmock.Sqlmock, cleanup func()) {
				db, mock, err := sqlmock.New()
				require.NoError(t, err)
				mySQL = New()
				mySQL.db = db
				cleanup = func() { _ = db.Close() }

				mock.ExpectQuery(queryVersion).
					WillReturnRows(mustMockRows(t, mariaV1054Version))
				mock.ExpectQuery(queryGlobalStatus).
					WillReturnRows(mustMockRows(t, mariaV1054GlobalStatus))
				mock.ExpectQuery(queryGlobalVariables).
					WillReturnRows(mustMockRows(t, mariaV1054GlobalVariables))
				mock.ExpectQuery(queryAllSlavesStatus).
					WillReturnError(errSQLSyntax)
				mock.ExpectQuery(queryUserStatistics).
					WillReturnError(errSQLSyntax)

				return mySQL, mock, cleanup
			},
			expected: map[string]int64{
				"aborted_connects":                      0,
				"binlog_cache_disk_use":                 0,
				"binlog_cache_use":                      0,
				"binlog_stmt_cache_disk_use":            0,
				"binlog_stmt_cache_use":                 0,
				"bytes_received":                        1001,
				"bytes_sent":                            195182,
				"com_delete":                            0,
				"com_insert":                            0,
				"com_replace":                           0,
				"com_select":                            3,
				"com_update":                            0,
				"connection_errors_accept":              0,
				"connection_errors_internal":            0,
				"connection_errors_max_connections":     0,
				"connection_errors_peer_address":        0,
				"connection_errors_select":              0,
				"connection_errors_tcpwrap":             0,
				"connections":                           13,
				"created_tmp_disk_tables":               0,
				"created_tmp_files":                     5,
				"created_tmp_tables":                    12,
				"handler_commit":                        26,
				"handler_delete":                        0,
				"handler_prepare":                       0,
				"handler_read_first":                    7,
				"handler_read_key":                      7,
				"handler_read_next":                     3,
				"handler_read_prev":                     0,
				"handler_read_rnd":                      0,
				"handler_read_rnd_next":                 5201,
				"handler_rollback":                      1,
				"handler_savepoint":                     0,
				"handler_savepoint_rollback":            0,
				"handler_update":                        6,
				"handler_write":                         1,
				"innodb_buffer_pool_bytes_data":         5357568,
				"innodb_buffer_pool_bytes_dirty":        0,
				"innodb_buffer_pool_pages_data":         327,
				"innodb_buffer_pool_pages_dirty":        0,
				"innodb_buffer_pool_pages_flushed":      134,
				"innodb_buffer_pool_pages_free":         7727,
				"innodb_buffer_pool_pages_misc":         0,
				"innodb_buffer_pool_pages_total":        8054,
				"innodb_buffer_pool_read_ahead":         0,
				"innodb_buffer_pool_read_ahead_evicted": 0,
				"innodb_buffer_pool_read_ahead_rnd":     0,
				"innodb_buffer_pool_read_requests":      2369,
				"innodb_buffer_pool_reads":              196,
				"innodb_buffer_pool_wait_free":          0,
				"innodb_buffer_pool_write_requests":     853,
				"innodb_data_fsyncs":                    25,
				"innodb_data_pending_fsyncs":            0,
				"innodb_data_pending_reads":             0,
				"innodb_data_pending_writes":            0,
				"innodb_data_read":                      3211264,
				"innodb_data_reads":                     215,
				"innodb_data_writes":                    157,
				"innodb_data_written":                   2244608,
				"innodb_deadlocks":                      0,
				"innodb_log_waits":                      0,
				"innodb_log_write_requests":             0,
				"innodb_log_writes":                     20,
				"innodb_os_log_fsyncs":                  20,
				"innodb_os_log_pending_fsyncs":          0,
				"innodb_os_log_pending_writes":          0,
				"innodb_os_log_written":                 10240,
				"innodb_row_lock_current_waits":         0,
				"innodb_rows_deleted":                   0,
				"innodb_rows_inserted":                  0,
				"innodb_rows_read":                      0,
				"innodb_rows_updated":                   0,
				"key_blocks_not_flushed":                0,
				"key_blocks_unused":                     107163,
				"key_blocks_used":                       0,
				"key_read_requests":                     0,
				"key_reads":                             0,
				"key_write_requests":                    0,
				"key_writes":                            0,
				"max_connections":                       151,
				"max_used_connections":                  1,
				"open_files":                            24,
				"opened_files":                          80,
				"open_tables":                           13,
				"opened_tables":                         19,
				"qcache_free_blocks":                    1,
				"qcache_free_memory":                    1031304,
				"qcache_hits":                           0,
				"qcache_inserts":                        0,
				"qcache_lowmem_prunes":                  0,
				"qcache_not_cached":                     0,
				"qcache_queries_in_cache":               0,
				"qcache_total_blocks":                   1,
				"queries":                               32,
				"questions":                             24,
				"select_full_join":                      0,
				"select_full_range_join":                0,
				"select_range":                          0,
				"select_range_check":                    0,
				"select_scan":                           12,
				"slow_queries":                          0,
				"sort_merge_passes":                     0,
				"sort_range":                            0,
				"sort_scan":                             0,
				"table_locks_immediate":                 59,
				"table_locks_waited":                    0,
				"table_open_cache":                      2000,
				"thread_cache_misses":                   4615,
				"threads_cached":                        0,
				"threads_connected":                     1,
				"threads_created":                       6,
				"threads_running":                       1,
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

				mock.ExpectQuery(queryVersion).
					WillReturnRows(mustMockRows(t, mysqlV8021Version))
				mock.ExpectQuery(queryGlobalStatus).
					WillReturnRows(mustMockRows(t, mysqlV8021GlobalStatus))
				mock.ExpectQuery(queryGlobalVariables).
					WillReturnRows(mustMockRows(t, mysqlV8021GlobalVariables))
				mock.ExpectQuery(querySlaveStatus).
					WillReturnRows(mustMockRows(t, mysqlV8021SlaveStatus))

				return mySQL, mock, cleanup
			},
			expected: map[string]int64{
				"aborted_connects":                      0,
				"binlog_cache_disk_use":                 0,
				"binlog_cache_use":                      2,
				"binlog_stmt_cache_disk_use":            0,
				"binlog_stmt_cache_use":                 0,
				"bytes_received":                        13552,
				"bytes_sent":                            21281,
				"com_delete":                            0,
				"com_insert":                            0,
				"com_replace":                           0,
				"com_select":                            3,
				"com_update":                            0,
				"connection_errors_accept":              0,
				"connection_errors_internal":            0,
				"connection_errors_max_connections":     0,
				"connection_errors_peer_address":        0,
				"connection_errors_select":              0,
				"connection_errors_tcpwrap":             0,
				"connections":                           67,
				"created_tmp_disk_tables":               0,
				"created_tmp_files":                     5,
				"created_tmp_tables":                    2,
				"handler_commit":                        552,
				"handler_delete":                        0,
				"handler_prepare":                       8,
				"handler_read_first":                    34,
				"handler_read_key":                      1635,
				"handler_read_next":                     3891,
				"handler_read_prev":                     0,
				"handler_read_rnd":                      0,
				"handler_read_rnd_next":                 1011,
				"handler_rollback":                      0,
				"handler_savepoint":                     0,
				"handler_savepoint_rollback":            0,
				"handler_update":                        316,
				"handler_write":                         467,
				"innodb_buffer_pool_bytes_data":         15761408,
				"innodb_buffer_pool_bytes_dirty":        0,
				"innodb_buffer_pool_pages_data":         962,
				"innodb_buffer_pool_pages_dirty":        0,
				"innodb_buffer_pool_pages_flushed":      170,
				"innodb_buffer_pool_pages_free":         7226,
				"innodb_buffer_pool_pages_misc":         4,
				"innodb_buffer_pool_pages_total":        8192,
				"innodb_buffer_pool_read_ahead":         0,
				"innodb_buffer_pool_read_ahead_evicted": 0,
				"innodb_buffer_pool_read_ahead_rnd":     0,
				"innodb_buffer_pool_read_requests":      14452,
				"innodb_buffer_pool_reads":              818,
				"innodb_buffer_pool_wait_free":          0,
				"innodb_buffer_pool_write_requests":     1696,
				"innodb_data_fsyncs":                    76,
				"innodb_data_pending_fsyncs":            0,
				"innodb_data_pending_reads":             0,
				"innodb_data_pending_writes":            0,
				"innodb_data_read":                      13472768,
				"innodb_data_reads":                     840,
				"innodb_data_writes":                    252,
				"innodb_data_written":                   3002368,
				"innodb_log_waits":                      0,
				"innodb_log_write_requests":             664,
				"innodb_log_writes":                     26,
				"innodb_os_log_fsyncs":                  25,
				"innodb_os_log_pending_fsyncs":          0,
				"innodb_os_log_pending_writes":          0,
				"innodb_os_log_written":                 38912,
				"innodb_row_lock_current_waits":         0,
				"innodb_rows_deleted":                   0,
				"innodb_rows_inserted":                  0,
				"innodb_rows_read":                      0,
				"innodb_rows_updated":                   0,
				"key_blocks_not_flushed":                0,
				"key_blocks_unused":                     6698,
				"key_blocks_used":                       0,
				"key_read_requests":                     0,
				"key_reads":                             0,
				"key_write_requests":                    0,
				"key_writes":                            0,
				"max_connections":                       151,
				"max_used_connections":                  1,
				"open_files":                            2,
				"opened_files":                          2,
				"open_tables":                           64,
				"opened_tables":                         143,
				"queries":                               125,
				"questions":                             67,
				"seconds_behind_master_master1":         0,
				"seconds_behind_master_master2":         0,
				"select_full_join":                      0,
				"select_full_range_join":                0,
				"select_range":                          0,
				"select_range_check":                    0,
				"select_scan":                           4,
				"slave_io_running_master1":              0,
				"slave_io_running_master2":              0,
				"slave_sql_running_master1":             1,
				"slave_sql_running_master2":             1,
				"slow_queries":                          0,
				"sort_merge_passes":                     0,
				"sort_range":                            0,
				"sort_scan":                             0,
				"table_locks_immediate":                 2,
				"table_locks_waited":                    0,
				"table_open_cache":                      4000,
				"thread_cache_misses":                   149,
				"threads_cached":                        0,
				"threads_connected":                     1,
				"threads_created":                       1,
				"threads_running":                       2,
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
			ensureCollectedHasAllChartsDimsVarsIDs(t, mySQL, collected)
		})
	}
}

func ensureCollectedHasAllChartsDimsVarsIDs(t *testing.T, mySQL *MySQL, collected map[string]int64) {
	for _, chart := range *mySQL.Charts() {
		// https://mariadb.com/kb/en/server-status-variables/#connection_errors_accept
		minVer := semver.Version{Major: 10, Minor: 0, Patch: 4}
		if mySQL.isMariaDB && mySQL.version.LT(minVer) && chart.ID == "connection_errors" {
			continue
		}

		for _, dim := range chart.Dims {
			_, ok := collected[dim.ID]
			assert.Truef(t, ok, "collected metrics has no data for dim '%s' chart '%s'", dim.ID, chart.ID)
		}
		for _, v := range chart.Vars {
			_, ok := collected[v.ID]
			assert.Truef(t, ok, "collected metrics has no data for var '%s' chart '%s'", v.ID, chart.ID)
		}
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
