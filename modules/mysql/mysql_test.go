// SPDX-License-Identifier: GPL-3.0-or-later

package mysql

import (
	"bufio"
	"bytes"
	"database/sql/driver"
	"errors"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/netdata/go.d.plugin/agent/module"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/blang/semver/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	dataMySQLV8030Version, _                = os.ReadFile("testdata/mysql/v8.0.30/version.txt")
	dataMySQLV8030GlobalStatus, _           = os.ReadFile("testdata/mysql/v8.0.30/global_status.txt")
	dataMySQLV8030GlobalVariables, _        = os.ReadFile("testdata/mysql/v8.0.30/global_variables.txt")
	dataMySQLV8030SlaveStatusMultiSource, _ = os.ReadFile("testdata/mysql/v8.0.30/slave_status_multi_source.txt")
	dataMySQLV8030ProcessList, _            = os.ReadFile("testdata/mysql/v8.0.30/process_list.txt")

	dataMariaV5564Version, _         = os.ReadFile("testdata/mariadb/v5.5.64/version.txt")
	dataMariaV5564GlobalStatus, _    = os.ReadFile("testdata/mariadb/v5.5.64/global_status.txt")
	dataMariaV5564GlobalVariables, _ = os.ReadFile("testdata/mariadb/v5.5.64/global_variables.txt")
	dataMariaV5564ProcessList, _     = os.ReadFile("testdata/mariadb/v5.5.64/process_list.txt")

	dataMariaV1084Version, _                     = os.ReadFile("testdata/mariadb/v10.8.4/version.txt")
	dataMariaV1084GlobalStatus, _                = os.ReadFile("testdata/mariadb/v10.8.4/global_status.txt")
	dataMariaV1084GlobalVariables, _             = os.ReadFile("testdata/mariadb/v10.8.4/global_variables.txt")
	dataMariaV1084AllSlavesStatusSingleSource, _ = os.ReadFile("testdata/mariadb/v10.8.4/all_slaves_status_single_source.txt")
	dataMariaV1084AllSlavesStatusMultiSource, _  = os.ReadFile("testdata/mariadb/v10.8.4/all_slaves_status_multi_source.txt")
	dataMariaV1084UserStatistics, _              = os.ReadFile("testdata/mariadb/v10.8.4/user_statistics.txt")
	dataMariaV1084ProcessList, _                 = os.ReadFile("testdata/mariadb/v10.8.4/process_list.txt")

	dataMariaGaleraClusterV1084Version, _         = os.ReadFile("testdata/mariadb/v10.8.4-galera-cluster/version.txt")
	dataMariaGaleraClusterV1084GlobalStatus, _    = os.ReadFile("testdata/mariadb/v10.8.4-galera-cluster/global_status.txt")
	dataMariaGaleraClusterV1084GlobalVariables, _ = os.ReadFile("testdata/mariadb/v10.8.4-galera-cluster/global_variables.txt")
	dataMariaGaleraClusterV1084UserStatistics, _  = os.ReadFile("testdata/mariadb/v10.8.4-galera-cluster/user_statistics.txt")
	dataMariaGaleraClusterV1084ProcessList, _     = os.ReadFile("testdata/mariadb/v10.8.4-galera-cluster/process_list.txt")
)

func Test_testDataIsValid(t *testing.T) {
	for name, data := range map[string][]byte{
		"dataMySQLV8030Version":                dataMySQLV8030Version,
		"dataMySQLV8030GlobalStatus":           dataMySQLV8030GlobalStatus,
		"dataMySQLV8030GlobalVariables":        dataMySQLV8030GlobalVariables,
		"dataMySQLV8030SlaveStatusMultiSource": dataMySQLV8030SlaveStatusMultiSource,
		"dataMySQLV8030ProcessList":            dataMySQLV8030ProcessList,

		"dataMariaV5564Version":         dataMariaV5564Version,
		"dataMariaV5564GlobalStatus":    dataMariaV5564GlobalStatus,
		"dataMariaV5564GlobalVariables": dataMariaV5564GlobalVariables,
		"dataMariaV5564ProcessList":     dataMariaV5564ProcessList,

		"dataMariaV1084Version":                     dataMariaV1084Version,
		"dataMariaV1084GlobalStatus":                dataMariaV1084GlobalStatus,
		"dataMariaV1084GlobalVariables":             dataMariaV1084GlobalVariables,
		"dataMariaV1084AllSlavesStatusSingleSource": dataMariaV1084AllSlavesStatusSingleSource,
		"dataMariaV1084AllSlavesStatusMultiSource":  dataMariaV1084AllSlavesStatusMultiSource,
		"dataMariaV1084UserStatistics":              dataMariaV1084UserStatistics,
		"dataMariaV1084ProcessList":                 dataMariaV1084ProcessList,

		"dataMariaGaleraClusterV1084Version":         dataMariaGaleraClusterV1084Version,
		"dataMariaGaleraClusterV1084GlobalStatus":    dataMariaGaleraClusterV1084GlobalStatus,
		"dataMariaGaleraClusterV1084GlobalVariables": dataMariaGaleraClusterV1084GlobalVariables,
		"dataMariaGaleraClusterV1084UserStatistics":  dataMariaGaleraClusterV1084UserStatistics,
		"dataMariaGaleraClusterV1084ProcessList":     dataMariaGaleraClusterV1084ProcessList,
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

func TestMySQL_Collect(t *testing.T) {
	type testCaseStep struct {
		prepareMock func(t *testing.T, m sqlmock.Sqlmock)
		check       func(t *testing.T, my *MySQL)
	}
	tests := map[string][]testCaseStep{
		"MariaDB-Standalone[v5.5.46]: success on all queries": {
			{
				prepareMock: func(t *testing.T, m sqlmock.Sqlmock) {
					mockExpect(t, m, queryShowVersion, dataMariaV5564Version)
					mockExpect(t, m, queryShowGlobalStatus, dataMariaV5564GlobalStatus)
					mockExpect(t, m, queryShowGlobalVariables, dataMariaV5564GlobalVariables)
					mockExpect(t, m, queryShowSlaveStatus, nil)
					mockExpect(t, m, queryShowProcessList, dataMariaV5564ProcessList)
				},
				check: func(t *testing.T, my *MySQL) {
					mx := my.Collect()

					expected := map[string]int64{
						"aborted_connects":                      0,
						"binlog_cache_disk_use":                 0,
						"binlog_cache_use":                      0,
						"binlog_stmt_cache_disk_use":            0,
						"binlog_stmt_cache_use":                 0,
						"bytes_received":                        639,
						"bytes_sent":                            41620,
						"com_delete":                            0,
						"com_insert":                            0,
						"com_replace":                           0,
						"com_select":                            4,
						"com_update":                            0,
						"connections":                           4,
						"created_tmp_disk_tables":               0,
						"created_tmp_files":                     6,
						"created_tmp_tables":                    5,
						"handler_commit":                        0,
						"handler_delete":                        0,
						"handler_prepare":                       0,
						"handler_read_first":                    0,
						"handler_read_key":                      0,
						"handler_read_next":                     0,
						"handler_read_prev":                     0,
						"handler_read_rnd":                      0,
						"handler_read_rnd_next":                 1264,
						"handler_rollback":                      0,
						"handler_savepoint":                     0,
						"handler_savepoint_rollback":            0,
						"handler_update":                        0,
						"handler_write":                         0,
						"innodb_buffer_pool_bytes_data":         2342912,
						"innodb_buffer_pool_bytes_dirty":        0,
						"innodb_buffer_pool_pages_data":         143,
						"innodb_buffer_pool_pages_dirty":        0,
						"innodb_buffer_pool_pages_flushed":      0,
						"innodb_buffer_pool_pages_free":         16240,
						"innodb_buffer_pool_pages_misc":         0,
						"innodb_buffer_pool_pages_total":        16383,
						"innodb_buffer_pool_read_ahead":         0,
						"innodb_buffer_pool_read_ahead_evicted": 0,
						"innodb_buffer_pool_read_ahead_rnd":     0,
						"innodb_buffer_pool_read_requests":      459,
						"innodb_buffer_pool_reads":              144,
						"innodb_buffer_pool_wait_free":          0,
						"innodb_buffer_pool_write_requests":     0,
						"innodb_data_fsyncs":                    3,
						"innodb_data_pending_fsyncs":            0,
						"innodb_data_pending_reads":             0,
						"innodb_data_pending_writes":            0,
						"innodb_data_read":                      4542976,
						"innodb_data_reads":                     155,
						"innodb_data_writes":                    3,
						"innodb_data_written":                   1536,
						"innodb_deadlocks":                      0,
						"innodb_log_waits":                      0,
						"innodb_log_write_requests":             0,
						"innodb_log_writes":                     1,
						"innodb_os_log_fsyncs":                  3,
						"innodb_os_log_pending_fsyncs":          0,
						"innodb_os_log_pending_writes":          0,
						"innodb_os_log_written":                 512,
						"innodb_row_lock_current_waits":         0,
						"innodb_rows_deleted":                   0,
						"innodb_rows_inserted":                  0,
						"innodb_rows_read":                      0,
						"innodb_rows_updated":                   0,
						"key_blocks_not_flushed":                0,
						"key_blocks_unused":                     107171,
						"key_blocks_used":                       0,
						"key_read_requests":                     0,
						"key_reads":                             0,
						"key_write_requests":                    0,
						"key_writes":                            0,
						"max_connections":                       100,
						"max_used_connections":                  1,
						"open_files":                            21,
						"open_tables":                           26,
						"opened_files":                          84,
						"opened_tables":                         0,
						"process_list_fetch_query_duration":     0,
						"process_list_longest_query_duration":   9,
						"process_list_queries_count_system":     0,
						"process_list_queries_count_user":       2,
						"qcache_free_blocks":                    1,
						"qcache_free_memory":                    67091120,
						"qcache_hits":                           0,
						"qcache_inserts":                        0,
						"qcache_lowmem_prunes":                  0,
						"qcache_not_cached":                     4,
						"qcache_queries_in_cache":               0,
						"qcache_total_blocks":                   1,
						"queries":                               12,
						"questions":                             11,
						"select_full_join":                      0,
						"select_full_range_join":                0,
						"select_range":                          0,
						"select_range_check":                    0,
						"select_scan":                           5,
						"slow_queries":                          0,
						"sort_merge_passes":                     0,
						"sort_range":                            0,
						"sort_scan":                             0,
						"table_locks_immediate":                 36,
						"table_locks_waited":                    0,
						"table_open_cache":                      400,
						"thread_cache_misses":                   2500,
						"threads_cached":                        0,
						"threads_connected":                     1,
						"threads_created":                       1,
						"threads_running":                       1,
					}

					copyProcessListQueryDuration(mx, expected)
					assert.Equal(t, expected, mx)
				},
			},
		},
		"MariaDB-Standalone[v10.8.4]: success on all queries": {
			{
				prepareMock: func(t *testing.T, m sqlmock.Sqlmock) {
					mockExpect(t, m, queryShowVersion, dataMariaV1084Version)
					mockExpect(t, m, queryShowGlobalStatus, dataMariaV1084GlobalStatus)
					mockExpect(t, m, queryShowGlobalVariables, dataMariaV1084GlobalVariables)
					mockExpect(t, m, queryShowAllSlavesStatus, nil)
					mockExpect(t, m, queryShowUserStatistics, dataMariaV1084UserStatistics)
					mockExpect(t, m, queryShowProcessList, dataMariaV1084ProcessList)
				},
				check: func(t *testing.T, my *MySQL) {
					mx := my.Collect()

					expected := map[string]int64{
						"aborted_connects":                      2,
						"binlog_cache_disk_use":                 0,
						"binlog_cache_use":                      0,
						"binlog_stmt_cache_disk_use":            0,
						"binlog_stmt_cache_use":                 0,
						"bytes_received":                        81392,
						"bytes_sent":                            56794,
						"com_delete":                            0,
						"com_insert":                            0,
						"com_replace":                           0,
						"com_select":                            6,
						"com_update":                            0,
						"connection_errors_accept":              0,
						"connection_errors_internal":            0,
						"connection_errors_max_connections":     0,
						"connection_errors_peer_address":        0,
						"connection_errors_select":              0,
						"connection_errors_tcpwrap":             0,
						"connections":                           12,
						"created_tmp_disk_tables":               0,
						"created_tmp_files":                     5,
						"created_tmp_tables":                    2,
						"handler_commit":                        30,
						"handler_delete":                        0,
						"handler_prepare":                       0,
						"handler_read_first":                    7,
						"handler_read_key":                      7,
						"handler_read_next":                     3,
						"handler_read_prev":                     0,
						"handler_read_rnd":                      0,
						"handler_read_rnd_next":                 626,
						"handler_rollback":                      0,
						"handler_savepoint":                     0,
						"handler_savepoint_rollback":            0,
						"handler_update":                        3,
						"handler_write":                         13,
						"innodb_buffer_pool_bytes_data":         5062656,
						"innodb_buffer_pool_bytes_dirty":        475136,
						"innodb_buffer_pool_pages_data":         309,
						"innodb_buffer_pool_pages_dirty":        29,
						"innodb_buffer_pool_pages_flushed":      0,
						"innodb_buffer_pool_pages_free":         7755,
						"innodb_buffer_pool_pages_misc":         0,
						"innodb_buffer_pool_pages_total":        8064,
						"innodb_buffer_pool_read_ahead":         0,
						"innodb_buffer_pool_read_ahead_evicted": 0,
						"innodb_buffer_pool_read_ahead_rnd":     0,
						"innodb_buffer_pool_read_requests":      1911,
						"innodb_buffer_pool_reads":              171,
						"innodb_buffer_pool_wait_free":          0,
						"innodb_buffer_pool_write_requests":     148,
						"innodb_data_fsyncs":                    17,
						"innodb_data_pending_fsyncs":            0,
						"innodb_data_pending_reads":             0,
						"innodb_data_pending_writes":            0,
						"innodb_data_read":                      2801664,
						"innodb_data_reads":                     185,
						"innodb_data_writes":                    16,
						"innodb_data_written":                   0,
						"innodb_deadlocks":                      0,
						"innodb_log_waits":                      0,
						"innodb_log_write_requests":             109,
						"innodb_log_writes":                     15,
						"innodb_os_log_written":                 6097,
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
						"open_files":                            29,
						"open_tables":                           10,
						"opened_files":                          100,
						"opened_tables":                         16,
						"process_list_fetch_query_duration":     0,
						"process_list_longest_query_duration":   9,
						"process_list_queries_count_system":     0,
						"process_list_queries_count_user":       2,
						"qcache_free_blocks":                    1,
						"qcache_free_memory":                    1031272,
						"qcache_hits":                           0,
						"qcache_inserts":                        0,
						"qcache_lowmem_prunes":                  0,
						"qcache_not_cached":                     0,
						"qcache_queries_in_cache":               0,
						"qcache_total_blocks":                   1,
						"queries":                               33,
						"questions":                             24,
						"select_full_join":                      0,
						"select_full_range_join":                0,
						"select_range":                          0,
						"select_range_check":                    0,
						"select_scan":                           2,
						"slow_queries":                          0,
						"sort_merge_passes":                     0,
						"sort_range":                            0,
						"sort_scan":                             0,
						"table_locks_immediate":                 60,
						"table_locks_waited":                    0,
						"table_open_cache":                      2000,
						"thread_cache_misses":                   1666,
						"threads_cached":                        0,
						"threads_connected":                     1,
						"threads_created":                       2,
						"threads_running":                       3,
						"userstats_netdata_cpu_time":            77,
						"userstats_netdata_other_commands":      0,
						"userstats_netdata_rows_deleted":        0,
						"userstats_netdata_rows_inserted":       0,
						"userstats_netdata_rows_read":           0,
						"userstats_netdata_rows_sent":           99,
						"userstats_netdata_rows_updated":        0,
						"userstats_netdata_select_commands":     33,
						"userstats_netdata_update_commands":     0,
						"userstats_root_cpu_time":               0,
						"userstats_root_other_commands":         0,
						"userstats_root_rows_deleted":           0,
						"userstats_root_rows_inserted":          0,
						"userstats_root_rows_read":              0,
						"userstats_root_rows_sent":              2,
						"userstats_root_rows_updated":           0,
						"userstats_root_select_commands":        0,
						"userstats_root_update_commands":        0,
						"wsrep_cluster_size":                    0,
						"wsrep_cluster_status":                  2,
						"wsrep_connected":                       0,
						"wsrep_local_bf_aborts":                 0,
						"wsrep_ready":                           0,
						"wsrep_thread_count":                    0,
					}

					copyProcessListQueryDuration(mx, expected)
					assert.Equal(t, expected, mx)
				},
			},
		},
		"MariaDB-SingleSourceReplication[v10.8.4]: success on all queries": {
			{
				prepareMock: func(t *testing.T, m sqlmock.Sqlmock) {
					mockExpect(t, m, queryShowVersion, dataMariaV1084Version)
					mockExpect(t, m, queryShowGlobalStatus, dataMariaV1084GlobalStatus)
					mockExpect(t, m, queryShowGlobalVariables, dataMariaV1084GlobalVariables)
					mockExpect(t, m, queryShowAllSlavesStatus, dataMariaV1084AllSlavesStatusSingleSource)
					mockExpect(t, m, queryShowUserStatistics, dataMariaV1084UserStatistics)
					mockExpect(t, m, queryShowProcessList, dataMariaV1084ProcessList)
				},
				check: func(t *testing.T, my *MySQL) {
					mx := my.Collect()

					expected := map[string]int64{
						"aborted_connects":                      2,
						"binlog_cache_disk_use":                 0,
						"binlog_cache_use":                      0,
						"binlog_stmt_cache_disk_use":            0,
						"binlog_stmt_cache_use":                 0,
						"bytes_received":                        81392,
						"bytes_sent":                            56794,
						"com_delete":                            0,
						"com_insert":                            0,
						"com_replace":                           0,
						"com_select":                            6,
						"com_update":                            0,
						"connection_errors_accept":              0,
						"connection_errors_internal":            0,
						"connection_errors_max_connections":     0,
						"connection_errors_peer_address":        0,
						"connection_errors_select":              0,
						"connection_errors_tcpwrap":             0,
						"connections":                           12,
						"created_tmp_disk_tables":               0,
						"created_tmp_files":                     5,
						"created_tmp_tables":                    2,
						"handler_commit":                        30,
						"handler_delete":                        0,
						"handler_prepare":                       0,
						"handler_read_first":                    7,
						"handler_read_key":                      7,
						"handler_read_next":                     3,
						"handler_read_prev":                     0,
						"handler_read_rnd":                      0,
						"handler_read_rnd_next":                 626,
						"handler_rollback":                      0,
						"handler_savepoint":                     0,
						"handler_savepoint_rollback":            0,
						"handler_update":                        3,
						"handler_write":                         13,
						"innodb_buffer_pool_bytes_data":         5062656,
						"innodb_buffer_pool_bytes_dirty":        475136,
						"innodb_buffer_pool_pages_data":         309,
						"innodb_buffer_pool_pages_dirty":        29,
						"innodb_buffer_pool_pages_flushed":      0,
						"innodb_buffer_pool_pages_free":         7755,
						"innodb_buffer_pool_pages_misc":         0,
						"innodb_buffer_pool_pages_total":        8064,
						"innodb_buffer_pool_read_ahead":         0,
						"innodb_buffer_pool_read_ahead_evicted": 0,
						"innodb_buffer_pool_read_ahead_rnd":     0,
						"innodb_buffer_pool_read_requests":      1911,
						"innodb_buffer_pool_reads":              171,
						"innodb_buffer_pool_wait_free":          0,
						"innodb_buffer_pool_write_requests":     148,
						"innodb_data_fsyncs":                    17,
						"innodb_data_pending_fsyncs":            0,
						"innodb_data_pending_reads":             0,
						"innodb_data_pending_writes":            0,
						"innodb_data_read":                      2801664,
						"innodb_data_reads":                     185,
						"innodb_data_writes":                    16,
						"innodb_data_written":                   0,
						"innodb_deadlocks":                      0,
						"innodb_log_waits":                      0,
						"innodb_log_write_requests":             109,
						"innodb_log_writes":                     15,
						"innodb_os_log_written":                 6097,
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
						"open_files":                            29,
						"open_tables":                           10,
						"opened_files":                          100,
						"opened_tables":                         16,
						"process_list_fetch_query_duration":     0,
						"process_list_longest_query_duration":   9,
						"process_list_queries_count_system":     0,
						"process_list_queries_count_user":       2,
						"qcache_free_blocks":                    1,
						"qcache_free_memory":                    1031272,
						"qcache_hits":                           0,
						"qcache_inserts":                        0,
						"qcache_lowmem_prunes":                  0,
						"qcache_not_cached":                     0,
						"qcache_queries_in_cache":               0,
						"qcache_total_blocks":                   1,
						"queries":                               33,
						"questions":                             24,
						"seconds_behind_master":                 0,
						"select_full_join":                      0,
						"select_full_range_join":                0,
						"select_range":                          0,
						"select_range_check":                    0,
						"select_scan":                           2,
						"slave_io_running":                      1,
						"slave_sql_running":                     1,
						"slow_queries":                          0,
						"sort_merge_passes":                     0,
						"sort_range":                            0,
						"sort_scan":                             0,
						"table_locks_immediate":                 60,
						"table_locks_waited":                    0,
						"table_open_cache":                      2000,
						"thread_cache_misses":                   1666,
						"threads_cached":                        0,
						"threads_connected":                     1,
						"threads_created":                       2,
						"threads_running":                       3,
						"userstats_netdata_cpu_time":            77,
						"userstats_netdata_other_commands":      0,
						"userstats_netdata_rows_deleted":        0,
						"userstats_netdata_rows_inserted":       0,
						"userstats_netdata_rows_read":           0,
						"userstats_netdata_rows_sent":           99,
						"userstats_netdata_rows_updated":        0,
						"userstats_netdata_select_commands":     33,
						"userstats_netdata_update_commands":     0,
						"userstats_root_cpu_time":               0,
						"userstats_root_other_commands":         0,
						"userstats_root_rows_deleted":           0,
						"userstats_root_rows_inserted":          0,
						"userstats_root_rows_read":              0,
						"userstats_root_rows_sent":              2,
						"userstats_root_rows_updated":           0,
						"userstats_root_select_commands":        0,
						"userstats_root_update_commands":        0,
						"wsrep_cluster_size":                    0,
						"wsrep_cluster_status":                  2,
						"wsrep_connected":                       0,
						"wsrep_local_bf_aborts":                 0,
						"wsrep_ready":                           0,
						"wsrep_thread_count":                    0,
					}

					copyProcessListQueryDuration(mx, expected)
					assert.Equal(t, expected, mx)
				},
			},
		},
		"MariaDB-MultiSourceReplication[v10.8.4]: success on all queries": {
			{
				prepareMock: func(t *testing.T, m sqlmock.Sqlmock) {
					mockExpect(t, m, queryShowVersion, dataMariaV1084Version)
					mockExpect(t, m, queryShowGlobalStatus, dataMariaV1084GlobalStatus)
					mockExpect(t, m, queryShowGlobalVariables, dataMariaV1084GlobalVariables)
					mockExpect(t, m, queryShowAllSlavesStatus, dataMariaV1084AllSlavesStatusMultiSource)
					mockExpect(t, m, queryShowUserStatistics, dataMariaV1084UserStatistics)
					mockExpect(t, m, queryShowProcessList, dataMariaV1084ProcessList)
				},
				check: func(t *testing.T, my *MySQL) {
					mx := my.Collect()

					expected := map[string]int64{
						"aborted_connects":                      2,
						"binlog_cache_disk_use":                 0,
						"binlog_cache_use":                      0,
						"binlog_stmt_cache_disk_use":            0,
						"binlog_stmt_cache_use":                 0,
						"bytes_received":                        81392,
						"bytes_sent":                            56794,
						"com_delete":                            0,
						"com_insert":                            0,
						"com_replace":                           0,
						"com_select":                            6,
						"com_update":                            0,
						"connection_errors_accept":              0,
						"connection_errors_internal":            0,
						"connection_errors_max_connections":     0,
						"connection_errors_peer_address":        0,
						"connection_errors_select":              0,
						"connection_errors_tcpwrap":             0,
						"connections":                           12,
						"created_tmp_disk_tables":               0,
						"created_tmp_files":                     5,
						"created_tmp_tables":                    2,
						"handler_commit":                        30,
						"handler_delete":                        0,
						"handler_prepare":                       0,
						"handler_read_first":                    7,
						"handler_read_key":                      7,
						"handler_read_next":                     3,
						"handler_read_prev":                     0,
						"handler_read_rnd":                      0,
						"handler_read_rnd_next":                 626,
						"handler_rollback":                      0,
						"handler_savepoint":                     0,
						"handler_savepoint_rollback":            0,
						"handler_update":                        3,
						"handler_write":                         13,
						"innodb_buffer_pool_bytes_data":         5062656,
						"innodb_buffer_pool_bytes_dirty":        475136,
						"innodb_buffer_pool_pages_data":         309,
						"innodb_buffer_pool_pages_dirty":        29,
						"innodb_buffer_pool_pages_flushed":      0,
						"innodb_buffer_pool_pages_free":         7755,
						"innodb_buffer_pool_pages_misc":         0,
						"innodb_buffer_pool_pages_total":        8064,
						"innodb_buffer_pool_read_ahead":         0,
						"innodb_buffer_pool_read_ahead_evicted": 0,
						"innodb_buffer_pool_read_ahead_rnd":     0,
						"innodb_buffer_pool_read_requests":      1911,
						"innodb_buffer_pool_reads":              171,
						"innodb_buffer_pool_wait_free":          0,
						"innodb_buffer_pool_write_requests":     148,
						"innodb_data_fsyncs":                    17,
						"innodb_data_pending_fsyncs":            0,
						"innodb_data_pending_reads":             0,
						"innodb_data_pending_writes":            0,
						"innodb_data_read":                      2801664,
						"innodb_data_reads":                     185,
						"innodb_data_writes":                    16,
						"innodb_data_written":                   0,
						"innodb_deadlocks":                      0,
						"innodb_log_waits":                      0,
						"innodb_log_write_requests":             109,
						"innodb_log_writes":                     15,
						"innodb_os_log_written":                 6097,
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
						"open_files":                            29,
						"open_tables":                           10,
						"opened_files":                          100,
						"opened_tables":                         16,
						"process_list_fetch_query_duration":     0,
						"process_list_longest_query_duration":   9,
						"process_list_queries_count_system":     0,
						"process_list_queries_count_user":       2,
						"qcache_free_blocks":                    1,
						"qcache_free_memory":                    1031272,
						"qcache_hits":                           0,
						"qcache_inserts":                        0,
						"qcache_lowmem_prunes":                  0,
						"qcache_not_cached":                     0,
						"qcache_queries_in_cache":               0,
						"qcache_total_blocks":                   1,
						"queries":                               33,
						"questions":                             24,
						"seconds_behind_master_master1":         0,
						"seconds_behind_master_master2":         0,
						"select_full_join":                      0,
						"select_full_range_join":                0,
						"select_range":                          0,
						"select_range_check":                    0,
						"select_scan":                           2,
						"slave_io_running_master1":              1,
						"slave_io_running_master2":              1,
						"slave_sql_running_master1":             1,
						"slave_sql_running_master2":             1,
						"slow_queries":                          0,
						"sort_merge_passes":                     0,
						"sort_range":                            0,
						"sort_scan":                             0,
						"table_locks_immediate":                 60,
						"table_locks_waited":                    0,
						"table_open_cache":                      2000,
						"thread_cache_misses":                   1666,
						"threads_cached":                        0,
						"threads_connected":                     1,
						"threads_created":                       2,
						"threads_running":                       3,
						"userstats_netdata_cpu_time":            77,
						"userstats_netdata_other_commands":      0,
						"userstats_netdata_rows_deleted":        0,
						"userstats_netdata_rows_inserted":       0,
						"userstats_netdata_rows_read":           0,
						"userstats_netdata_rows_sent":           99,
						"userstats_netdata_rows_updated":        0,
						"userstats_netdata_select_commands":     33,
						"userstats_netdata_update_commands":     0,
						"userstats_root_cpu_time":               0,
						"userstats_root_other_commands":         0,
						"userstats_root_rows_deleted":           0,
						"userstats_root_rows_inserted":          0,
						"userstats_root_rows_read":              0,
						"userstats_root_rows_sent":              2,
						"userstats_root_rows_updated":           0,
						"userstats_root_select_commands":        0,
						"userstats_root_update_commands":        0,
						"wsrep_cluster_size":                    0,
						"wsrep_cluster_status":                  2,
						"wsrep_connected":                       0,
						"wsrep_local_bf_aborts":                 0,
						"wsrep_ready":                           0,
						"wsrep_thread_count":                    0,
					}

					copyProcessListQueryDuration(mx, expected)
					assert.Equal(t, expected, mx)
				},
			},
		},
		"MariaDB-GaleraCluster[v10.8.4]: success on all queries": {
			{
				prepareMock: func(t *testing.T, m sqlmock.Sqlmock) {
					mockExpect(t, m, queryShowVersion, dataMariaGaleraClusterV1084Version)
					mockExpect(t, m, queryShowGlobalStatus, dataMariaGaleraClusterV1084GlobalStatus)
					mockExpect(t, m, queryShowGlobalVariables, dataMariaGaleraClusterV1084GlobalVariables)
					mockExpect(t, m, queryShowAllSlavesStatus, nil)
					mockExpect(t, m, queryShowUserStatistics, dataMariaGaleraClusterV1084UserStatistics)
					mockExpect(t, m, queryShowProcessList, dataMariaGaleraClusterV1084ProcessList)
				},
				check: func(t *testing.T, my *MySQL) {
					mx := my.Collect()

					expected := map[string]int64{
						"aborted_connects":                      0,
						"binlog_cache_disk_use":                 0,
						"binlog_cache_use":                      0,
						"binlog_stmt_cache_disk_use":            0,
						"binlog_stmt_cache_use":                 0,
						"bytes_received":                        3009,
						"bytes_sent":                            228856,
						"com_delete":                            6,
						"com_insert":                            0,
						"com_replace":                           0,
						"com_select":                            12,
						"com_update":                            0,
						"connection_errors_accept":              0,
						"connection_errors_internal":            0,
						"connection_errors_max_connections":     0,
						"connection_errors_peer_address":        0,
						"connection_errors_select":              0,
						"connection_errors_tcpwrap":             0,
						"connections":                           15,
						"created_tmp_disk_tables":               4,
						"created_tmp_files":                     5,
						"created_tmp_tables":                    17,
						"handler_commit":                        37,
						"handler_delete":                        7,
						"handler_prepare":                       0,
						"handler_read_first":                    3,
						"handler_read_key":                      9,
						"handler_read_next":                     1,
						"handler_read_prev":                     0,
						"handler_read_rnd":                      0,
						"handler_read_rnd_next":                 6222,
						"handler_rollback":                      0,
						"handler_savepoint":                     0,
						"handler_savepoint_rollback":            0,
						"handler_update":                        0,
						"handler_write":                         9,
						"innodb_buffer_pool_bytes_data":         5193728,
						"innodb_buffer_pool_bytes_dirty":        2260992,
						"innodb_buffer_pool_pages_data":         317,
						"innodb_buffer_pool_pages_dirty":        138,
						"innodb_buffer_pool_pages_flushed":      0,
						"innodb_buffer_pool_pages_free":         7747,
						"innodb_buffer_pool_pages_misc":         0,
						"innodb_buffer_pool_pages_total":        8064,
						"innodb_buffer_pool_read_ahead":         0,
						"innodb_buffer_pool_read_ahead_evicted": 0,
						"innodb_buffer_pool_read_ahead_rnd":     0,
						"innodb_buffer_pool_read_requests":      2298,
						"innodb_buffer_pool_reads":              184,
						"innodb_buffer_pool_wait_free":          0,
						"innodb_buffer_pool_write_requests":     203,
						"innodb_data_fsyncs":                    15,
						"innodb_data_pending_fsyncs":            0,
						"innodb_data_pending_reads":             0,
						"innodb_data_pending_writes":            0,
						"innodb_data_read":                      3014656,
						"innodb_data_reads":                     201,
						"innodb_data_writes":                    14,
						"innodb_data_written":                   0,
						"innodb_deadlocks":                      0,
						"innodb_log_waits":                      0,
						"innodb_log_write_requests":             65,
						"innodb_log_writes":                     13,
						"innodb_os_log_written":                 4785,
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
						"open_files":                            7,
						"open_tables":                           0,
						"opened_files":                          125,
						"opened_tables":                         24,
						"process_list_fetch_query_duration":     0,
						"process_list_longest_query_duration":   9,
						"process_list_queries_count_system":     0,
						"process_list_queries_count_user":       2,
						"qcache_free_blocks":                    1,
						"qcache_free_memory":                    1031272,
						"qcache_hits":                           0,
						"qcache_inserts":                        0,
						"qcache_lowmem_prunes":                  0,
						"qcache_not_cached":                     0,
						"qcache_queries_in_cache":               0,
						"qcache_total_blocks":                   1,
						"queries":                               75,
						"questions":                             62,
						"select_full_join":                      0,
						"select_full_range_join":                0,
						"select_range":                          0,
						"select_range_check":                    0,
						"select_scan":                           17,
						"slow_queries":                          0,
						"sort_merge_passes":                     0,
						"sort_range":                            0,
						"sort_scan":                             0,
						"table_locks_immediate":                 17,
						"table_locks_waited":                    0,
						"table_open_cache":                      2000,
						"thread_cache_misses":                   4000,
						"threads_cached":                        0,
						"threads_connected":                     1,
						"threads_created":                       6,
						"threads_running":                       1,
						"userstats_netdata_cpu_time":            77,
						"userstats_netdata_other_commands":      0,
						"userstats_netdata_rows_deleted":        0,
						"userstats_netdata_rows_inserted":       0,
						"userstats_netdata_rows_read":           0,
						"userstats_netdata_rows_sent":           99,
						"userstats_netdata_rows_updated":        0,
						"userstats_netdata_select_commands":     33,
						"userstats_netdata_update_commands":     0,
						"userstats_root_cpu_time":               0,
						"userstats_root_other_commands":         0,
						"userstats_root_rows_deleted":           0,
						"userstats_root_rows_inserted":          0,
						"userstats_root_rows_read":              0,
						"userstats_root_rows_sent":              2,
						"userstats_root_rows_updated":           0,
						"userstats_root_select_commands":        0,
						"userstats_root_update_commands":        0,
						"wsrep_cluster_size":                    3,
						"wsrep_cluster_status":                  0,
						"wsrep_cluster_weight":                  3,
						"wsrep_connected":                       1,
						"wsrep_flow_control_paused_ns":          0,
						"wsrep_local_bf_aborts":                 0,
						"wsrep_local_cert_failures":             0,
						"wsrep_local_recv_queue":                0,
						"wsrep_local_send_queue":                0,
						"wsrep_local_state":                     4,
						"wsrep_open_transactions":               0,
						"wsrep_ready":                           1,
						"wsrep_received":                        11,
						"wsrep_received_bytes":                  1410,
						"wsrep_replicated":                      0,
						"wsrep_replicated_bytes":                0,
						"wsrep_thread_count":                    5,
					}

					copyProcessListQueryDuration(mx, expected)
					assert.Equal(t, expected, mx)
				},
			},
		},
		"MySQL-MultiSourceReplication[v8.0.30]: success on all queries": {
			{
				prepareMock: func(t *testing.T, m sqlmock.Sqlmock) {
					mockExpect(t, m, queryShowVersion, dataMySQLV8030Version)
					mockExpect(t, m, queryShowGlobalStatus, dataMySQLV8030GlobalStatus)
					mockExpect(t, m, queryShowGlobalVariables, dataMySQLV8030GlobalVariables)
					mockExpect(t, m, queryShowSlaveStatus, dataMySQLV8030SlaveStatusMultiSource)
					mockExpect(t, m, queryShowProcessList, dataMySQLV8030ProcessList)
				},
				check: func(t *testing.T, my *MySQL) {
					mx := my.Collect()

					expected := map[string]int64{
						"aborted_connects":                      0,
						"binlog_cache_disk_use":                 0,
						"binlog_cache_use":                      6,
						"binlog_stmt_cache_disk_use":            0,
						"binlog_stmt_cache_use":                 0,
						"bytes_received":                        5584,
						"bytes_sent":                            70700,
						"com_delete":                            0,
						"com_insert":                            0,
						"com_replace":                           0,
						"com_select":                            2,
						"com_update":                            0,
						"connection_errors_accept":              0,
						"connection_errors_internal":            0,
						"connection_errors_max_connections":     0,
						"connection_errors_peer_address":        0,
						"connection_errors_select":              0,
						"connection_errors_tcpwrap":             0,
						"connections":                           25,
						"created_tmp_disk_tables":               0,
						"created_tmp_files":                     5,
						"created_tmp_tables":                    6,
						"handler_commit":                        720,
						"handler_delete":                        8,
						"handler_prepare":                       24,
						"handler_read_first":                    50,
						"handler_read_key":                      1914,
						"handler_read_next":                     4303,
						"handler_read_prev":                     0,
						"handler_read_rnd":                      0,
						"handler_read_rnd_next":                 4723,
						"handler_rollback":                      1,
						"handler_savepoint":                     0,
						"handler_savepoint_rollback":            0,
						"handler_update":                        373,
						"handler_write":                         1966,
						"innodb_buffer_pool_bytes_data":         17121280,
						"innodb_buffer_pool_bytes_dirty":        0,
						"innodb_buffer_pool_pages_data":         1045,
						"innodb_buffer_pool_pages_dirty":        0,
						"innodb_buffer_pool_pages_flushed":      361,
						"innodb_buffer_pool_pages_free":         7143,
						"innodb_buffer_pool_pages_misc":         4,
						"innodb_buffer_pool_pages_total":        8192,
						"innodb_buffer_pool_read_ahead":         0,
						"innodb_buffer_pool_read_ahead_evicted": 0,
						"innodb_buffer_pool_read_ahead_rnd":     0,
						"innodb_buffer_pool_read_requests":      16723,
						"innodb_buffer_pool_reads":              878,
						"innodb_buffer_pool_wait_free":          0,
						"innodb_buffer_pool_write_requests":     2377,
						"innodb_data_fsyncs":                    255,
						"innodb_data_pending_fsyncs":            0,
						"innodb_data_pending_reads":             0,
						"innodb_data_pending_writes":            0,
						"innodb_data_read":                      14453760,
						"innodb_data_reads":                     899,
						"innodb_data_writes":                    561,
						"innodb_data_written":                   6128128,
						"innodb_log_waits":                      0,
						"innodb_log_write_requests":             1062,
						"innodb_log_writes":                     116,
						"innodb_os_log_fsyncs":                  69,
						"innodb_os_log_pending_fsyncs":          0,
						"innodb_os_log_pending_writes":          0,
						"innodb_os_log_written":                 147968,
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
						"max_used_connections":                  2,
						"open_files":                            8,
						"open_tables":                           127,
						"opened_files":                          8,
						"opened_tables":                         208,
						"process_list_fetch_query_duration":     0,
						"process_list_longest_query_duration":   9,
						"process_list_queries_count_system":     0,
						"process_list_queries_count_user":       2,
						"queries":                               27,
						"questions":                             15,
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
						"table_locks_immediate":                 6,
						"table_locks_waited":                    0,
						"table_open_cache":                      4000,
						"thread_cache_misses":                   800,
						"threads_cached":                        1,
						"threads_connected":                     1,
						"threads_created":                       2,
						"threads_running":                       2,
					}

					copyProcessListQueryDuration(mx, expected)
					assert.Equal(t, expected, mx)
				},
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			db, mock, err := sqlmock.New(
				sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual),
			)
			require.NoError(t, err)
			my := New()
			my.db = db
			defer func() { _ = db.Close() }()

			require.True(t, my.Init())

			for i, step := range test {
				t.Run(fmt.Sprintf("step[%d]", i), func(t *testing.T) {
					step.prepareMock(t, mock)
					step.check(t, my)
				})
			}
			assert.NoError(t, mock.ExpectationsWereMet())
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

func copyProcessListQueryDuration(dst, src map[string]int64) {
	if _, ok := dst["process_list_fetch_query_duration"]; !ok {
		return
	}
	if _, ok := src["process_list_fetch_query_duration"]; !ok {
		return
	}
	dst["process_list_fetch_query_duration"] = src["process_list_fetch_query_duration"]
}

func mustMockRows(t *testing.T, data []byte) *sqlmock.Rows {
	rows, err := prepareMockRows(data)
	require.NoError(t, err)
	return rows
}

func mockExpect(t *testing.T, mock sqlmock.Sqlmock, query string, rows []byte) {
	mock.ExpectQuery(query).WillReturnRows(mustMockRows(t, rows)).RowsWillBeClosed()
}

func mockExpectErr(mock sqlmock.Sqlmock, query string) {
	mock.ExpectQuery(query).WillReturnError(fmt.Errorf("mock error (%s)", query))
}

func prepareMockRows(data []byte) (*sqlmock.Rows, error) {
	if len(data) == 0 {
		return sqlmock.NewRows(nil), nil
	}

	r := bytes.NewReader(data)
	sc := bufio.NewScanner(r)

	var numColumns int
	var rows *sqlmock.Rows

	for sc.Scan() {
		s := strings.TrimSpace(strings.Trim(sc.Text(), "|"))
		switch {
		case s == "",
			strings.HasPrefix(s, "+"),
			strings.HasPrefix(s, "ft_boolean_syntax"):
			continue
		}

		parts := strings.Split(s, "|")
		for i, v := range parts {
			parts[i] = strings.TrimSpace(v)
		}

		if rows == nil {
			numColumns = len(parts)
			rows = sqlmock.NewRows(parts)
			continue
		}

		if len(parts) != numColumns {
			return nil, fmt.Errorf("prepareMockRows(): columns != values (%d/%d)", numColumns, len(parts))
		}

		values := make([]driver.Value, len(parts))
		for i, v := range parts {
			values[i] = v
		}
		rows.AddRow(values...)
	}

	if rows == nil {
		return nil, errors.New("prepareMockRows(): nil rows result")
	}

	return rows, sc.Err()
}
