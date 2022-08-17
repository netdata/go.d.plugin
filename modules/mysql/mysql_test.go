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
	mariaV5546Version, _         = os.ReadFile("testdata/mariadb/v5.5.46/version.txt")
	mariaV5546GlobalStatus, _    = os.ReadFile("testdata/mariadb/v5.5.46/global_status.txt")
	mariaV5546GlobalVariables, _ = os.ReadFile("testdata/mariadb/v5.5.46/global_variables.txt")
	mariaV5546SlaveStatus, _     = os.ReadFile("testdata/mariadb/v5.5.46/slave_status.txt")
	mariaV5546ProcessList, _     = os.ReadFile("testdata/mariadb/v5.5.46/process_list.txt")

	mariaV1054Version, _         = os.ReadFile("testdata/mariadb/v10.5.4/version.txt")
	mariaV1054GlobalStatus, _    = os.ReadFile("testdata/mariadb/v10.5.4/global_status.txt")
	mariaV1054GlobalVariables, _ = os.ReadFile("testdata/mariadb/v10.5.4/global_variables.txt")
	mariaV1054UserStatistics, _  = os.ReadFile("testdata/mariadb/v10.5.4/user_statistics.txt")
	mariaV1054AllSlavesStatus, _ = os.ReadFile("testdata/mariadb/v10.5.4/all_slaves_status.txt")
	mariaV1054ProcessList, _     = os.ReadFile("testdata/mariadb/v10.5.4/process_list.txt")

	dataMariaV1083Version, _         = os.ReadFile("testdata/mariadb/standalone/v10.8.3/version.txt")
	dataMariaV1083GlobalStatus, _    = os.ReadFile("testdata/mariadb/standalone/v10.8.3/global_status.txt")
	dataMariaV1083GlobalVariables, _ = os.ReadFile("testdata/mariadb/standalone/v10.8.3/global_variables.txt")
	dataMariaV1083UserStatistics, _  = os.ReadFile("testdata/mariadb/standalone/v10.8.3/user_statistics.txt")
	dataMariaV1083ProcessList, _     = os.ReadFile("testdata/mariadb/standalone/v10.8.3/process_list.txt")

	mysqlV8021Version, _         = os.ReadFile("testdata/mysql/v8.0.21/version.txt")
	mysqlV8021GlobalStatus, _    = os.ReadFile("testdata/mysql/v8.0.21/global_status.txt")
	mysqlV8021GlobalVariables, _ = os.ReadFile("testdata/mysql/v8.0.21/global_variables.txt")
	mysqlV8021SlaveStatus, _     = os.ReadFile("testdata/mysql/v8.0.21/slave_status.txt")
	mysqlV8021ProcessList, _     = os.ReadFile("testdata/mysql/v8.0.21/process_list.txt")
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

		"dataMariaV1083Version":         dataMariaV1083Version,
		"dataMariaV1083GlobalStatus":    dataMariaV1083GlobalStatus,
		"dataMariaV1083GlobalVariables": dataMariaV1083GlobalVariables,
		"dataMariaV1083UserStatistics":  dataMariaV1083UserStatistics,
		"dataMariaV1083ProcessList":     dataMariaV1083ProcessList,

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

func TestMySQL_Collect(t *testing.T) {
	type testCaseStep struct {
		prepareMock func(t *testing.T, m sqlmock.Sqlmock)
		check       func(t *testing.T, my *MySQL)
	}
	tests := map[string][]testCaseStep{
		"MariaV10.8.3[Standalone]: success on all queries": {
			{
				prepareMock: func(t *testing.T, m sqlmock.Sqlmock) {
					mockExpect(t, m, queryShowVersion, dataMariaV1083Version)
					mockExpect(t, m, queryShowGlobalStatus, dataMariaV1083GlobalStatus)
					mockExpect(t, m, queryShowGlobalVariables, dataMariaV1083GlobalVariables)
					mockExpect(t, m, queryShowAllSlavesStatus, nil)
					mockExpect(t, m, queryShowUserStatistics, dataMariaV1083UserStatistics)
					mockExpect(t, m, queryShowProcessList, dataMariaV1083ProcessList)
				},
				check: func(t *testing.T, my *MySQL) {
					mx := my.Collect()

					expected := map[string]int64{
						"aborted_connects":                      0,
						"binlog_cache_disk_use":                 0,
						"binlog_cache_use":                      0,
						"binlog_stmt_cache_disk_use":            0,
						"binlog_stmt_cache_use":                 0,
						"bytes_received":                        892626,
						"bytes_sent":                            42783889,
						"com_delete":                            0,
						"com_insert":                            0,
						"com_replace":                           0,
						"com_select":                            2191,
						"com_update":                            0,
						"connection_errors_accept":              0,
						"connection_errors_internal":            0,
						"connection_errors_max_connections":     0,
						"connection_errors_peer_address":        0,
						"connection_errors_select":              0,
						"connection_errors_tcpwrap":             0,
						"connections":                           9,
						"created_tmp_disk_tables":               2185,
						"created_tmp_files":                     4,
						"created_tmp_tables":                    8743,
						"handler_commit":                        9,
						"handler_delete":                        0,
						"handler_prepare":                       0,
						"handler_read_first":                    3,
						"handler_read_key":                      0,
						"handler_read_next":                     1,
						"handler_read_prev":                     0,
						"handler_read_rnd":                      0,
						"handler_read_rnd_next":                 1253151,
						"handler_rollback":                      0,
						"handler_savepoint":                     0,
						"handler_savepoint_rollback":            0,
						"handler_update":                        0,
						"handler_write":                         0,
						"innodb_buffer_pool_bytes_data":         4653056,
						"innodb_buffer_pool_bytes_dirty":        98304,
						"innodb_buffer_pool_pages_data":         284,
						"innodb_buffer_pool_pages_dirty":        6,
						"innodb_buffer_pool_pages_flushed":      0,
						"innodb_buffer_pool_pages_free":         7780,
						"innodb_buffer_pool_pages_misc":         0,
						"innodb_buffer_pool_pages_total":        8064,
						"innodb_buffer_pool_read_ahead":         0,
						"innodb_buffer_pool_read_ahead_evicted": 0,
						"innodb_buffer_pool_read_ahead_rnd":     0,
						"innodb_buffer_pool_read_requests":      1594,
						"innodb_buffer_pool_reads":              153,
						"innodb_buffer_pool_wait_free":          0,
						"innodb_buffer_pool_write_requests":     521,
						"innodb_data_fsyncs":                    2,
						"innodb_data_pending_fsyncs":            0,
						"innodb_data_pending_reads":             0,
						"innodb_data_pending_writes":            0,
						"innodb_data_read":                      2506752,
						"innodb_data_reads":                     166,
						"innodb_data_writes":                    1,
						"innodb_data_written":                   0,
						"innodb_deadlocks":                      0,
						"innodb_log_waits":                      0,
						"innodb_log_write_requests":             6,
						"innodb_log_writes":                     1,
						"innodb_os_log_written":                 168,
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
						"max_used_connections":                  2,
						"open_files":                            23,
						"open_tables":                           10,
						"opened_files":                          8813,
						"opened_tables":                         17,
						"process_list_fetch_query_duration":     0,
						"process_list_longest_query_duration":   0,
						"process_list_queries_count_system":     0,
						"process_list_queries_count_user":       0,
						"qcache_free_blocks":                    1,
						"qcache_free_memory":                    1031272,
						"qcache_hits":                           0,
						"qcache_inserts":                        0,
						"qcache_lowmem_prunes":                  0,
						"qcache_not_cached":                     0,
						"qcache_queries_in_cache":               0,
						"qcache_total_blocks":                   1,
						"queries":                               8757,
						"questions":                             8757,
						"select_full_join":                      0,
						"select_full_range_join":                0,
						"select_range":                          0,
						"select_range_check":                    0,
						"select_scan":                           8743,
						"slow_queries":                          0,
						"sort_merge_passes":                     0,
						"sort_range":                            0,
						"sort_scan":                             2185,
						"table_locks_immediate":                 17,
						"table_locks_waited":                    0,
						"table_open_cache":                      2000,
						"thread_cache_misses":                   2222,
						"threads_cached":                        0,
						"threads_connected":                     2,
						"threads_created":                       2,
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
						"wsrep_cluster_size":                    0,
						"wsrep_cluster_status":                  2,
						"wsrep_connected":                       0,
						"wsrep_local_bf_aborts":                 0,
						"wsrep_ready":                           0,
						"wsrep_thread_count":                    0,
					}

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
