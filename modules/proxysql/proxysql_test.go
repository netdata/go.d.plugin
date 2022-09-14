package proxysql

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
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	dataProxySQLV2010Version, _                    = os.ReadFile("testdata/v2.0.10/version.txt")
	dataProxySQLV2010StatsMySQLGlobal, _           = os.ReadFile("testdata/v2.0.10/stats_mysql_global.txt")
	dataProxySQLV2010StatsMemoryMetrics, _         = os.ReadFile("testdata/v2.0.10/stats_memory_metrics.txt")
	dataProxySQLV2010StatsMySQLCommandsCounters, _ = os.ReadFile("testdata/v2.0.10/stats_mysql_commands_counters.txt")
	dataProxySQLV2010StatsMySQLUsers, _            = os.ReadFile("testdata/v2.0.10/stats_mysql_users.txt")
)

func Test_testDataIsValid(t *testing.T) {
	for name, data := range map[string][]byte{
		"dataProxySQLV2010Version":                    dataProxySQLV2010Version,
		"dataProxySQLV2010StatsMySQLGlobal":           dataProxySQLV2010StatsMySQLGlobal,
		"dataProxySQLV2010StatsMemoryMetrics":         dataProxySQLV2010StatsMemoryMetrics,
		"dataProxySQLV2010StatsMySQLCommandsCounters": dataProxySQLV2010StatsMySQLCommandsCounters,
		"dataProxySQLV2010StatsMySQLUsers":            dataProxySQLV2010StatsMySQLUsers,
	} {
		require.NotNilf(t, data, name)
		_, err := prepareMockRows(data)
		require.NoErrorf(t, err, name)
	}
}

func TestNew(t *testing.T) {
	assert.Implements(t, (*module.Module)(nil), New())
}

func TestProxySQL_Init(t *testing.T) {
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
			proxySQL := New()
			proxySQL.Config = test.config

			if test.wantFail {
				assert.False(t, proxySQL.Init())
			} else {
				assert.True(t, proxySQL.Init())
			}
		})
	}
}

func TestProxySQL_Cleanup(t *testing.T) {
	tests := map[string]func(t *testing.T) (proxySQL *ProxySQL, cleanup func()){
		"db connection not initialized": func(t *testing.T) (proxySQL *ProxySQL, cleanup func()) {
			return New(), func() {}
		},
		"db connection initialized": func(t *testing.T) (proxySQL *ProxySQL, cleanup func()) {
			db, mock, err := sqlmock.New()
			require.NoError(t, err)

			mock.ExpectClose()
			proxySQL = New()
			proxySQL.db = db
			cleanup = func() { _ = db.Close() }

			return proxySQL, cleanup
		},
	}

	for name, prepare := range tests {
		t.Run(name, func(t *testing.T) {
			proxySQL, cleanup := prepare(t)
			defer cleanup()

			assert.NotPanics(t, proxySQL.Cleanup)
			assert.Nil(t, proxySQL.db)
		})
	}
}

func TestProxySQL_Charts(t *testing.T) {
	assert.NotNil(t, New().Charts())
}

func TestProxySQL_Check(t *testing.T) {
	tests := map[string]struct {
		prepareMock func(t *testing.T, m sqlmock.Sqlmock)
		wantFail    bool
	}{
		"success on all queries": {
			wantFail: false,
			prepareMock: func(t *testing.T, m sqlmock.Sqlmock) {
				mockExpect(t, m, queryStatsMySQLMemoryMetrics, dataProxySQLV2010StatsMemoryMetrics)
				mockExpect(t, m, queryStatsMySQLCommandsCounters, dataProxySQLV2010StatsMySQLCommandsCounters)
				mockExpect(t, m, queryStatsMySQLGlobal, dataProxySQLV2010StatsMySQLGlobal)
				mockExpect(t, m, queryStatsMySQLUsers, dataProxySQLV2010StatsMySQLUsers)
			},
		},
		"fails when error on querying global memory metrics": {
			wantFail: true,
			prepareMock: func(t *testing.T, m sqlmock.Sqlmock) {
				mockExpectErr(m, queryStatsMySQLMemoryMetrics)
			},
		},
		"fails when error on querying command counters": {
			wantFail: true,
			prepareMock: func(t *testing.T, m sqlmock.Sqlmock) {
				mockExpect(t, m, queryStatsMySQLMemoryMetrics, dataProxySQLV2010StatsMemoryMetrics)
				mockExpectErr(m, queryStatsMySQLCommandsCounters)
			},
		},
		"fails when error on querying mysql global status": {
			wantFail: true,
			prepareMock: func(t *testing.T, m sqlmock.Sqlmock) {
				mockExpect(t, m, queryStatsMySQLMemoryMetrics, dataProxySQLV2010StatsMemoryMetrics)
				mockExpect(t, m, queryStatsMySQLCommandsCounters, dataProxySQLV2010StatsMySQLCommandsCounters)
				mockExpectErr(m, queryStatsMySQLGlobal)
			},
		},
		"fails when error on querying command counter statistics": {
			wantFail: true,
			prepareMock: func(t *testing.T, m sqlmock.Sqlmock) {
				mockExpect(t, m, queryStatsMySQLMemoryMetrics, dataProxySQLV2010StatsMemoryMetrics)
				mockExpect(t, m, queryStatsMySQLCommandsCounters, dataProxySQLV2010StatsMySQLCommandsCounters)
				mockExpect(t, m, queryStatsMySQLGlobal, dataProxySQLV2010StatsMySQLGlobal)
				mockExpectErr(m, queryStatsMySQLUsers)
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			db, mock, err := sqlmock.New(
				sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual),
			)
			require.NoError(t, err)
			proxySQL := New()
			proxySQL.db = db
			defer func() { _ = db.Close() }()

			require.True(t, proxySQL.Init())

			test.prepareMock(t, mock)

			if test.wantFail {
				assert.False(t, proxySQL.Check())
			} else {
				assert.True(t, proxySQL.Check())
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestProxySQL_Collect(t *testing.T) {
	type testCaseStep struct {
		prepareMock func(t *testing.T, m sqlmock.Sqlmock)
		check       func(t *testing.T, my *ProxySQL)
	}
	tests := map[string][]testCaseStep{

		"ProxySQL[v2.0.10]: success on all queries": {
			{
				prepareMock: func(t *testing.T, m sqlmock.Sqlmock) {
					mockExpect(t, m, queryStatsMySQLMemoryMetrics, dataProxySQLV2010StatsMemoryMetrics)
					mockExpect(t, m, queryStatsMySQLCommandsCounters, dataProxySQLV2010StatsMySQLCommandsCounters)
					mockExpect(t, m, queryStatsMySQLGlobal, dataProxySQLV2010StatsMySQLGlobal)
					mockExpect(t, m, queryStatsMySQLUsers, dataProxySQLV2010StatsMySQLUsers)
				},
				check: func(t *testing.T, my *ProxySQL) {
					mx := my.Collect()

					expected := map[string]int64{
						"access_denied_max_connections":                   0,
						"access_denied_max_user_connections":              0,
						"access_denied_wrong_password":                    2,
						"active_transactions":                             0,
						"auth_memory":                                     1044,
						"automatic_detected_sql_injection":                0,
						"aws_aurora_replicas_skipped_during_query":        0,
						"backend_lagging_during_query":                    8880,
						"backend_offline_during_query":                    8,
						"backend_query_time_nsec":                         0,
						"client_connections_aborted":                      2,
						"client_connections_connected":                    3,
						"client_connections_created":                      5458991,
						"client_connections_hostgroup_locked":             0,
						"client_connections_non_idle":                     3,
						"com_autocommit":                                  0,
						"com_autocommit_filtered":                         0,
						"com_backend_change_user":                         188694,
						"com_backend_init_db":                             0,
						"com_backend_set_names":                           1517893,
						"com_backend_stmt_close":                          0,
						"com_backend_stmt_execute":                        36303146,
						"com_backend_stmt_prepare":                        16858208,
						"com_commit":                                      0,
						"com_commit_filtered":                             0,
						"com_frontend_init_db":                            2,
						"com_frontend_set_names":                          0,
						"com_frontend_stmt_close":                         32137933,
						"com_frontend_stmt_execute":                       36314138,
						"com_frontend_stmt_prepare":                       32185987,
						"com_frontend_use_db":                             0,
						"com_rollback":                                    0,
						"com_rollback_filtered":                           0,
						"connpool_get_conn_failure":                       212943,
						"connpool_get_conn_immediate":                     13361,
						"connpool_get_conn_latency_awareness":             0,
						"connpool_get_conn_success":                       36319474,
						"connpool_memory_bytes":                           932248,
						"generated_error_packets":                         231,
						"hostgroup_locked_queries":                        0,
						"hostgroup_locked_set_cmds":                       0,
						"jemalloc_active":                                 385101824,
						"jemalloc_allocated":                              379402432,
						"jemalloc_mapped":                                 430993408,
						"jemalloc_metadata":                               17418872,
						"jemalloc_resident":                               403759104,
						"jemalloc_retained":                               260542464,
						"max_connect_timeouts":                            227,
						"myhgm_myconnpoll_destroy":                        15150,
						"myhgm_myconnpoll_get":                            36519056,
						"myhgm_myconnpoll_get_ok":                         36306113,
						"myhgm_myconnpoll_push":                           37358734,
						"myhgm_myconnpoll_reset":                          2,
						"mysql_backend_buffers_bytes":                     0,
						"mysql_command_alter_table_cnt_100ms":             0,
						"mysql_command_alter_table_cnt_100us":             0,
						"mysql_command_alter_table_cnt_10ms":              0,
						"mysql_command_alter_table_cnt_10s":               0,
						"mysql_command_alter_table_cnt_1ms":               0,
						"mysql_command_alter_table_cnt_1s":                0,
						"mysql_command_alter_table_cnt_500ms":             0,
						"mysql_command_alter_table_cnt_500us":             0,
						"mysql_command_alter_table_cnt_5ms":               0,
						"mysql_command_alter_table_cnt_5s":                0,
						"mysql_command_alter_table_cnt_infs":              0,
						"mysql_command_alter_table_total_cnt":             0,
						"mysql_command_alter_table_total_time_us":         0,
						"mysql_command_alter_view_cnt_100ms":              0,
						"mysql_command_alter_view_cnt_100us":              0,
						"mysql_command_alter_view_cnt_10ms":               0,
						"mysql_command_alter_view_cnt_10s":                0,
						"mysql_command_alter_view_cnt_1ms":                0,
						"mysql_command_alter_view_cnt_1s":                 0,
						"mysql_command_alter_view_cnt_500ms":              0,
						"mysql_command_alter_view_cnt_500us":              0,
						"mysql_command_alter_view_cnt_5ms":                0,
						"mysql_command_alter_view_cnt_5s":                 0,
						"mysql_command_alter_view_cnt_infs":               0,
						"mysql_command_alter_view_total_cnt":              0,
						"mysql_command_alter_view_total_time_us":          0,
						"mysql_command_analyze_table_cnt_100ms":           0,
						"mysql_command_analyze_table_cnt_100us":           0,
						"mysql_command_analyze_table_cnt_10ms":            0,
						"mysql_command_analyze_table_cnt_10s":             0,
						"mysql_command_analyze_table_cnt_1ms":             0,
						"mysql_command_analyze_table_cnt_1s":              0,
						"mysql_command_analyze_table_cnt_500ms":           0,
						"mysql_command_analyze_table_cnt_500us":           0,
						"mysql_command_analyze_table_cnt_5ms":             0,
						"mysql_command_analyze_table_cnt_5s":              0,
						"mysql_command_analyze_table_cnt_infs":            0,
						"mysql_command_analyze_table_total_cnt":           0,
						"mysql_command_analyze_table_total_time_us":       0,
						"mysql_command_begin_cnt_100ms":                   0,
						"mysql_command_begin_cnt_100us":                   0,
						"mysql_command_begin_cnt_10ms":                    0,
						"mysql_command_begin_cnt_10s":                     0,
						"mysql_command_begin_cnt_1ms":                     0,
						"mysql_command_begin_cnt_1s":                      0,
						"mysql_command_begin_cnt_500ms":                   0,
						"mysql_command_begin_cnt_500us":                   0,
						"mysql_command_begin_cnt_5ms":                     0,
						"mysql_command_begin_cnt_5s":                      0,
						"mysql_command_begin_cnt_infs":                    0,
						"mysql_command_begin_total_cnt":                   0,
						"mysql_command_begin_total_time_us":               0,
						"mysql_command_call_cnt_100ms":                    0,
						"mysql_command_call_cnt_100us":                    0,
						"mysql_command_call_cnt_10ms":                     0,
						"mysql_command_call_cnt_10s":                      0,
						"mysql_command_call_cnt_1ms":                      0,
						"mysql_command_call_cnt_1s":                       0,
						"mysql_command_call_cnt_500ms":                    0,
						"mysql_command_call_cnt_500us":                    0,
						"mysql_command_call_cnt_5ms":                      0,
						"mysql_command_call_cnt_5s":                       0,
						"mysql_command_call_cnt_infs":                     0,
						"mysql_command_call_total_cnt":                    0,
						"mysql_command_call_total_time_us":                0,
						"mysql_command_change_master_cnt_100ms":           0,
						"mysql_command_change_master_cnt_100us":           0,
						"mysql_command_change_master_cnt_10ms":            0,
						"mysql_command_change_master_cnt_10s":             0,
						"mysql_command_change_master_cnt_1ms":             0,
						"mysql_command_change_master_cnt_1s":              0,
						"mysql_command_change_master_cnt_500ms":           0,
						"mysql_command_change_master_cnt_500us":           0,
						"mysql_command_change_master_cnt_5ms":             0,
						"mysql_command_change_master_cnt_5s":              0,
						"mysql_command_change_master_cnt_infs":            0,
						"mysql_command_change_master_total_cnt":           0,
						"mysql_command_change_master_total_time_us":       0,
						"mysql_command_commit_cnt_100ms":                  0,
						"mysql_command_commit_cnt_100us":                  0,
						"mysql_command_commit_cnt_10ms":                   0,
						"mysql_command_commit_cnt_10s":                    0,
						"mysql_command_commit_cnt_1ms":                    0,
						"mysql_command_commit_cnt_1s":                     0,
						"mysql_command_commit_cnt_500ms":                  0,
						"mysql_command_commit_cnt_500us":                  0,
						"mysql_command_commit_cnt_5ms":                    0,
						"mysql_command_commit_cnt_5s":                     0,
						"mysql_command_commit_cnt_infs":                   0,
						"mysql_command_commit_total_cnt":                  0,
						"mysql_command_commit_total_time_us":              0,
						"mysql_command_create_database_cnt_100ms":         0,
						"mysql_command_create_database_cnt_100us":         0,
						"mysql_command_create_database_cnt_10ms":          0,
						"mysql_command_create_database_cnt_10s":           0,
						"mysql_command_create_database_cnt_1ms":           0,
						"mysql_command_create_database_cnt_1s":            0,
						"mysql_command_create_database_cnt_500ms":         0,
						"mysql_command_create_database_cnt_500us":         0,
						"mysql_command_create_database_cnt_5ms":           0,
						"mysql_command_create_database_cnt_5s":            0,
						"mysql_command_create_database_cnt_infs":          0,
						"mysql_command_create_database_total_cnt":         0,
						"mysql_command_create_database_total_time_us":     0,
						"mysql_command_create_index_cnt_100ms":            0,
						"mysql_command_create_index_cnt_100us":            0,
						"mysql_command_create_index_cnt_10ms":             0,
						"mysql_command_create_index_cnt_10s":              0,
						"mysql_command_create_index_cnt_1ms":              0,
						"mysql_command_create_index_cnt_1s":               0,
						"mysql_command_create_index_cnt_500ms":            0,
						"mysql_command_create_index_cnt_500us":            0,
						"mysql_command_create_index_cnt_5ms":              0,
						"mysql_command_create_index_cnt_5s":               0,
						"mysql_command_create_index_cnt_infs":             0,
						"mysql_command_create_index_total_cnt":            0,
						"mysql_command_create_index_total_time_us":        0,
						"mysql_command_create_table_cnt_100ms":            0,
						"mysql_command_create_table_cnt_100us":            0,
						"mysql_command_create_table_cnt_10ms":             0,
						"mysql_command_create_table_cnt_10s":              0,
						"mysql_command_create_table_cnt_1ms":              0,
						"mysql_command_create_table_cnt_1s":               0,
						"mysql_command_create_table_cnt_500ms":            0,
						"mysql_command_create_table_cnt_500us":            0,
						"mysql_command_create_table_cnt_5ms":              0,
						"mysql_command_create_table_cnt_5s":               0,
						"mysql_command_create_table_cnt_infs":             0,
						"mysql_command_create_table_total_cnt":            0,
						"mysql_command_create_table_total_time_us":        0,
						"mysql_command_create_temporary_cnt_100ms":        0,
						"mysql_command_create_temporary_cnt_100us":        0,
						"mysql_command_create_temporary_cnt_10ms":         0,
						"mysql_command_create_temporary_cnt_10s":          0,
						"mysql_command_create_temporary_cnt_1ms":          0,
						"mysql_command_create_temporary_cnt_1s":           0,
						"mysql_command_create_temporary_cnt_500ms":        0,
						"mysql_command_create_temporary_cnt_500us":        0,
						"mysql_command_create_temporary_cnt_5ms":          0,
						"mysql_command_create_temporary_cnt_5s":           0,
						"mysql_command_create_temporary_cnt_infs":         0,
						"mysql_command_create_temporary_total_cnt":        0,
						"mysql_command_create_temporary_total_time_us":    0,
						"mysql_command_create_trigger_cnt_100ms":          0,
						"mysql_command_create_trigger_cnt_100us":          0,
						"mysql_command_create_trigger_cnt_10ms":           0,
						"mysql_command_create_trigger_cnt_10s":            0,
						"mysql_command_create_trigger_cnt_1ms":            0,
						"mysql_command_create_trigger_cnt_1s":             0,
						"mysql_command_create_trigger_cnt_500ms":          0,
						"mysql_command_create_trigger_cnt_500us":          0,
						"mysql_command_create_trigger_cnt_5ms":            0,
						"mysql_command_create_trigger_cnt_5s":             0,
						"mysql_command_create_trigger_cnt_infs":           0,
						"mysql_command_create_trigger_total_cnt":          0,
						"mysql_command_create_trigger_total_time_us":      0,
						"mysql_command_create_user_cnt_100ms":             0,
						"mysql_command_create_user_cnt_100us":             0,
						"mysql_command_create_user_cnt_10ms":              0,
						"mysql_command_create_user_cnt_10s":               0,
						"mysql_command_create_user_cnt_1ms":               0,
						"mysql_command_create_user_cnt_1s":                0,
						"mysql_command_create_user_cnt_500ms":             0,
						"mysql_command_create_user_cnt_500us":             0,
						"mysql_command_create_user_cnt_5ms":               0,
						"mysql_command_create_user_cnt_5s":                0,
						"mysql_command_create_user_cnt_infs":              0,
						"mysql_command_create_user_total_cnt":             0,
						"mysql_command_create_user_total_time_us":         0,
						"mysql_command_create_view_cnt_100ms":             0,
						"mysql_command_create_view_cnt_100us":             0,
						"mysql_command_create_view_cnt_10ms":              0,
						"mysql_command_create_view_cnt_10s":               0,
						"mysql_command_create_view_cnt_1ms":               0,
						"mysql_command_create_view_cnt_1s":                0,
						"mysql_command_create_view_cnt_500ms":             0,
						"mysql_command_create_view_cnt_500us":             0,
						"mysql_command_create_view_cnt_5ms":               0,
						"mysql_command_create_view_cnt_5s":                0,
						"mysql_command_create_view_cnt_infs":              0,
						"mysql_command_create_view_total_cnt":             0,
						"mysql_command_create_view_total_time_us":         0,
						"mysql_command_deallocate_cnt_100ms":              0,
						"mysql_command_deallocate_cnt_100us":              0,
						"mysql_command_deallocate_cnt_10ms":               0,
						"mysql_command_deallocate_cnt_10s":                0,
						"mysql_command_deallocate_cnt_1ms":                0,
						"mysql_command_deallocate_cnt_1s":                 0,
						"mysql_command_deallocate_cnt_500ms":              0,
						"mysql_command_deallocate_cnt_500us":              0,
						"mysql_command_deallocate_cnt_5ms":                0,
						"mysql_command_deallocate_cnt_5s":                 0,
						"mysql_command_deallocate_cnt_infs":               0,
						"mysql_command_deallocate_total_cnt":              0,
						"mysql_command_deallocate_total_time_us":          0,
						"mysql_command_delete_cnt_100ms":                  0,
						"mysql_command_delete_cnt_100us":                  0,
						"mysql_command_delete_cnt_10ms":                   0,
						"mysql_command_delete_cnt_10s":                    0,
						"mysql_command_delete_cnt_1ms":                    0,
						"mysql_command_delete_cnt_1s":                     0,
						"mysql_command_delete_cnt_500ms":                  0,
						"mysql_command_delete_cnt_500us":                  0,
						"mysql_command_delete_cnt_5ms":                    0,
						"mysql_command_delete_cnt_5s":                     0,
						"mysql_command_delete_cnt_infs":                   0,
						"mysql_command_delete_total_cnt":                  0,
						"mysql_command_delete_total_time_us":              0,
						"mysql_command_describe_cnt_100ms":                0,
						"mysql_command_describe_cnt_100us":                0,
						"mysql_command_describe_cnt_10ms":                 0,
						"mysql_command_describe_cnt_10s":                  0,
						"mysql_command_describe_cnt_1ms":                  0,
						"mysql_command_describe_cnt_1s":                   0,
						"mysql_command_describe_cnt_500ms":                0,
						"mysql_command_describe_cnt_500us":                0,
						"mysql_command_describe_cnt_5ms":                  0,
						"mysql_command_describe_cnt_5s":                   0,
						"mysql_command_describe_cnt_infs":                 0,
						"mysql_command_describe_total_cnt":                0,
						"mysql_command_describe_total_time_us":            0,
						"mysql_command_drop_database_cnt_100ms":           0,
						"mysql_command_drop_database_cnt_100us":           0,
						"mysql_command_drop_database_cnt_10ms":            0,
						"mysql_command_drop_database_cnt_10s":             0,
						"mysql_command_drop_database_cnt_1ms":             0,
						"mysql_command_drop_database_cnt_1s":              0,
						"mysql_command_drop_database_cnt_500ms":           0,
						"mysql_command_drop_database_cnt_500us":           0,
						"mysql_command_drop_database_cnt_5ms":             0,
						"mysql_command_drop_database_cnt_5s":              0,
						"mysql_command_drop_database_cnt_infs":            0,
						"mysql_command_drop_database_total_cnt":           0,
						"mysql_command_drop_database_total_time_us":       0,
						"mysql_command_drop_index_cnt_100ms":              0,
						"mysql_command_drop_index_cnt_100us":              0,
						"mysql_command_drop_index_cnt_10ms":               0,
						"mysql_command_drop_index_cnt_10s":                0,
						"mysql_command_drop_index_cnt_1ms":                0,
						"mysql_command_drop_index_cnt_1s":                 0,
						"mysql_command_drop_index_cnt_500ms":              0,
						"mysql_command_drop_index_cnt_500us":              0,
						"mysql_command_drop_index_cnt_5ms":                0,
						"mysql_command_drop_index_cnt_5s":                 0,
						"mysql_command_drop_index_cnt_infs":               0,
						"mysql_command_drop_index_total_cnt":              0,
						"mysql_command_drop_index_total_time_us":          0,
						"mysql_command_drop_table_cnt_100ms":              0,
						"mysql_command_drop_table_cnt_100us":              0,
						"mysql_command_drop_table_cnt_10ms":               0,
						"mysql_command_drop_table_cnt_10s":                0,
						"mysql_command_drop_table_cnt_1ms":                0,
						"mysql_command_drop_table_cnt_1s":                 0,
						"mysql_command_drop_table_cnt_500ms":              0,
						"mysql_command_drop_table_cnt_500us":              0,
						"mysql_command_drop_table_cnt_5ms":                0,
						"mysql_command_drop_table_cnt_5s":                 0,
						"mysql_command_drop_table_cnt_infs":               0,
						"mysql_command_drop_table_total_cnt":              0,
						"mysql_command_drop_table_total_time_us":          0,
						"mysql_command_drop_trigger_cnt_100ms":            0,
						"mysql_command_drop_trigger_cnt_100us":            0,
						"mysql_command_drop_trigger_cnt_10ms":             0,
						"mysql_command_drop_trigger_cnt_10s":              0,
						"mysql_command_drop_trigger_cnt_1ms":              0,
						"mysql_command_drop_trigger_cnt_1s":               0,
						"mysql_command_drop_trigger_cnt_500ms":            0,
						"mysql_command_drop_trigger_cnt_500us":            0,
						"mysql_command_drop_trigger_cnt_5ms":              0,
						"mysql_command_drop_trigger_cnt_5s":               0,
						"mysql_command_drop_trigger_cnt_infs":             0,
						"mysql_command_drop_trigger_total_cnt":            0,
						"mysql_command_drop_trigger_total_time_us":        0,
						"mysql_command_drop_user_cnt_100ms":               0,
						"mysql_command_drop_user_cnt_100us":               0,
						"mysql_command_drop_user_cnt_10ms":                0,
						"mysql_command_drop_user_cnt_10s":                 0,
						"mysql_command_drop_user_cnt_1ms":                 0,
						"mysql_command_drop_user_cnt_1s":                  0,
						"mysql_command_drop_user_cnt_500ms":               0,
						"mysql_command_drop_user_cnt_500us":               0,
						"mysql_command_drop_user_cnt_5ms":                 0,
						"mysql_command_drop_user_cnt_5s":                  0,
						"mysql_command_drop_user_cnt_infs":                0,
						"mysql_command_drop_user_total_cnt":               0,
						"mysql_command_drop_user_total_time_us":           0,
						"mysql_command_drop_view_cnt_100ms":               0,
						"mysql_command_drop_view_cnt_100us":               0,
						"mysql_command_drop_view_cnt_10ms":                0,
						"mysql_command_drop_view_cnt_10s":                 0,
						"mysql_command_drop_view_cnt_1ms":                 0,
						"mysql_command_drop_view_cnt_1s":                  0,
						"mysql_command_drop_view_cnt_500ms":               0,
						"mysql_command_drop_view_cnt_500us":               0,
						"mysql_command_drop_view_cnt_5ms":                 0,
						"mysql_command_drop_view_cnt_5s":                  0,
						"mysql_command_drop_view_cnt_infs":                0,
						"mysql_command_drop_view_total_cnt":               0,
						"mysql_command_drop_view_total_time_us":           0,
						"mysql_command_execute_cnt_100ms":                 0,
						"mysql_command_execute_cnt_100us":                 0,
						"mysql_command_execute_cnt_10ms":                  0,
						"mysql_command_execute_cnt_10s":                   0,
						"mysql_command_execute_cnt_1ms":                   0,
						"mysql_command_execute_cnt_1s":                    0,
						"mysql_command_execute_cnt_500ms":                 0,
						"mysql_command_execute_cnt_500us":                 0,
						"mysql_command_execute_cnt_5ms":                   0,
						"mysql_command_execute_cnt_5s":                    0,
						"mysql_command_execute_cnt_infs":                  0,
						"mysql_command_execute_total_cnt":                 0,
						"mysql_command_execute_total_time_us":             0,
						"mysql_command_explain_cnt_100ms":                 0,
						"mysql_command_explain_cnt_100us":                 0,
						"mysql_command_explain_cnt_10ms":                  0,
						"mysql_command_explain_cnt_10s":                   0,
						"mysql_command_explain_cnt_1ms":                   0,
						"mysql_command_explain_cnt_1s":                    0,
						"mysql_command_explain_cnt_500ms":                 0,
						"mysql_command_explain_cnt_500us":                 0,
						"mysql_command_explain_cnt_5ms":                   0,
						"mysql_command_explain_cnt_5s":                    0,
						"mysql_command_explain_cnt_infs":                  0,
						"mysql_command_explain_total_cnt":                 0,
						"mysql_command_explain_total_time_us":             0,
						"mysql_command_flush_cnt_100ms":                   0,
						"mysql_command_flush_cnt_100us":                   0,
						"mysql_command_flush_cnt_10ms":                    0,
						"mysql_command_flush_cnt_10s":                     0,
						"mysql_command_flush_cnt_1ms":                     0,
						"mysql_command_flush_cnt_1s":                      0,
						"mysql_command_flush_cnt_500ms":                   0,
						"mysql_command_flush_cnt_500us":                   0,
						"mysql_command_flush_cnt_5ms":                     0,
						"mysql_command_flush_cnt_5s":                      0,
						"mysql_command_flush_cnt_infs":                    0,
						"mysql_command_flush_total_cnt":                   0,
						"mysql_command_flush_total_time_us":               0,
						"mysql_command_grant_cnt_100ms":                   0,
						"mysql_command_grant_cnt_100us":                   0,
						"mysql_command_grant_cnt_10ms":                    0,
						"mysql_command_grant_cnt_10s":                     0,
						"mysql_command_grant_cnt_1ms":                     0,
						"mysql_command_grant_cnt_1s":                      0,
						"mysql_command_grant_cnt_500ms":                   0,
						"mysql_command_grant_cnt_500us":                   0,
						"mysql_command_grant_cnt_5ms":                     0,
						"mysql_command_grant_cnt_5s":                      0,
						"mysql_command_grant_cnt_infs":                    0,
						"mysql_command_grant_total_cnt":                   0,
						"mysql_command_grant_total_time_us":               0,
						"mysql_command_insert_cnt_100ms":                  0,
						"mysql_command_insert_cnt_100us":                  0,
						"mysql_command_insert_cnt_10ms":                   0,
						"mysql_command_insert_cnt_10s":                    0,
						"mysql_command_insert_cnt_1ms":                    0,
						"mysql_command_insert_cnt_1s":                     0,
						"mysql_command_insert_cnt_500ms":                  0,
						"mysql_command_insert_cnt_500us":                  0,
						"mysql_command_insert_cnt_5ms":                    0,
						"mysql_command_insert_cnt_5s":                     0,
						"mysql_command_insert_cnt_infs":                   0,
						"mysql_command_insert_total_cnt":                  0,
						"mysql_command_insert_total_time_us":              0,
						"mysql_command_kill_cnt_100ms":                    0,
						"mysql_command_kill_cnt_100us":                    0,
						"mysql_command_kill_cnt_10ms":                     0,
						"mysql_command_kill_cnt_10s":                      0,
						"mysql_command_kill_cnt_1ms":                      0,
						"mysql_command_kill_cnt_1s":                       0,
						"mysql_command_kill_cnt_500ms":                    0,
						"mysql_command_kill_cnt_500us":                    0,
						"mysql_command_kill_cnt_5ms":                      0,
						"mysql_command_kill_cnt_5s":                       0,
						"mysql_command_kill_cnt_infs":                     0,
						"mysql_command_kill_total_cnt":                    0,
						"mysql_command_kill_total_time_us":                0,
						"mysql_command_load_cnt_100ms":                    0,
						"mysql_command_load_cnt_100us":                    0,
						"mysql_command_load_cnt_10ms":                     0,
						"mysql_command_load_cnt_10s":                      0,
						"mysql_command_load_cnt_1ms":                      0,
						"mysql_command_load_cnt_1s":                       0,
						"mysql_command_load_cnt_500ms":                    0,
						"mysql_command_load_cnt_500us":                    0,
						"mysql_command_load_cnt_5ms":                      0,
						"mysql_command_load_cnt_5s":                       0,
						"mysql_command_load_cnt_infs":                     0,
						"mysql_command_load_total_cnt":                    0,
						"mysql_command_load_total_time_us":                0,
						"mysql_command_lock_table_cnt_100ms":              0,
						"mysql_command_lock_table_cnt_100us":              0,
						"mysql_command_lock_table_cnt_10ms":               0,
						"mysql_command_lock_table_cnt_10s":                0,
						"mysql_command_lock_table_cnt_1ms":                0,
						"mysql_command_lock_table_cnt_1s":                 0,
						"mysql_command_lock_table_cnt_500ms":              0,
						"mysql_command_lock_table_cnt_500us":              0,
						"mysql_command_lock_table_cnt_5ms":                0,
						"mysql_command_lock_table_cnt_5s":                 0,
						"mysql_command_lock_table_cnt_infs":               0,
						"mysql_command_lock_table_total_cnt":              0,
						"mysql_command_lock_table_total_time_us":          0,
						"mysql_command_optimize_cnt_100ms":                0,
						"mysql_command_optimize_cnt_100us":                0,
						"mysql_command_optimize_cnt_10ms":                 0,
						"mysql_command_optimize_cnt_10s":                  0,
						"mysql_command_optimize_cnt_1ms":                  0,
						"mysql_command_optimize_cnt_1s":                   0,
						"mysql_command_optimize_cnt_500ms":                0,
						"mysql_command_optimize_cnt_500us":                0,
						"mysql_command_optimize_cnt_5ms":                  0,
						"mysql_command_optimize_cnt_5s":                   0,
						"mysql_command_optimize_cnt_infs":                 0,
						"mysql_command_optimize_total_cnt":                0,
						"mysql_command_optimize_total_time_us":            0,
						"mysql_command_prepare_cnt_100ms":                 0,
						"mysql_command_prepare_cnt_100us":                 0,
						"mysql_command_prepare_cnt_10ms":                  0,
						"mysql_command_prepare_cnt_10s":                   0,
						"mysql_command_prepare_cnt_1ms":                   0,
						"mysql_command_prepare_cnt_1s":                    0,
						"mysql_command_prepare_cnt_500ms":                 0,
						"mysql_command_prepare_cnt_500us":                 0,
						"mysql_command_prepare_cnt_5ms":                   0,
						"mysql_command_prepare_cnt_5s":                    0,
						"mysql_command_prepare_cnt_infs":                  0,
						"mysql_command_prepare_total_cnt":                 0,
						"mysql_command_prepare_total_time_us":             0,
						"mysql_command_purge_cnt_100ms":                   0,
						"mysql_command_purge_cnt_100us":                   0,
						"mysql_command_purge_cnt_10ms":                    0,
						"mysql_command_purge_cnt_10s":                     0,
						"mysql_command_purge_cnt_1ms":                     0,
						"mysql_command_purge_cnt_1s":                      0,
						"mysql_command_purge_cnt_500ms":                   0,
						"mysql_command_purge_cnt_500us":                   0,
						"mysql_command_purge_cnt_5ms":                     0,
						"mysql_command_purge_cnt_5s":                      0,
						"mysql_command_purge_cnt_infs":                    0,
						"mysql_command_purge_total_cnt":                   0,
						"mysql_command_purge_total_time_us":               0,
						"mysql_command_rename_table_cnt_100ms":            0,
						"mysql_command_rename_table_cnt_100us":            0,
						"mysql_command_rename_table_cnt_10ms":             0,
						"mysql_command_rename_table_cnt_10s":              0,
						"mysql_command_rename_table_cnt_1ms":              0,
						"mysql_command_rename_table_cnt_1s":               0,
						"mysql_command_rename_table_cnt_500ms":            0,
						"mysql_command_rename_table_cnt_500us":            0,
						"mysql_command_rename_table_cnt_5ms":              0,
						"mysql_command_rename_table_cnt_5s":               0,
						"mysql_command_rename_table_cnt_infs":             0,
						"mysql_command_rename_table_total_cnt":            0,
						"mysql_command_rename_table_total_time_us":        0,
						"mysql_command_replace_cnt_100ms":                 0,
						"mysql_command_replace_cnt_100us":                 0,
						"mysql_command_replace_cnt_10ms":                  0,
						"mysql_command_replace_cnt_10s":                   0,
						"mysql_command_replace_cnt_1ms":                   0,
						"mysql_command_replace_cnt_1s":                    0,
						"mysql_command_replace_cnt_500ms":                 0,
						"mysql_command_replace_cnt_500us":                 0,
						"mysql_command_replace_cnt_5ms":                   0,
						"mysql_command_replace_cnt_5s":                    0,
						"mysql_command_replace_cnt_infs":                  0,
						"mysql_command_replace_total_cnt":                 0,
						"mysql_command_replace_total_time_us":             0,
						"mysql_command_reset_master_cnt_100ms":            0,
						"mysql_command_reset_master_cnt_100us":            0,
						"mysql_command_reset_master_cnt_10ms":             0,
						"mysql_command_reset_master_cnt_10s":              0,
						"mysql_command_reset_master_cnt_1ms":              0,
						"mysql_command_reset_master_cnt_1s":               0,
						"mysql_command_reset_master_cnt_500ms":            0,
						"mysql_command_reset_master_cnt_500us":            0,
						"mysql_command_reset_master_cnt_5ms":              0,
						"mysql_command_reset_master_cnt_5s":               0,
						"mysql_command_reset_master_cnt_infs":             0,
						"mysql_command_reset_master_total_cnt":            0,
						"mysql_command_reset_master_total_time_us":        0,
						"mysql_command_reset_slave_cnt_100ms":             0,
						"mysql_command_reset_slave_cnt_100us":             0,
						"mysql_command_reset_slave_cnt_10ms":              0,
						"mysql_command_reset_slave_cnt_10s":               0,
						"mysql_command_reset_slave_cnt_1ms":               0,
						"mysql_command_reset_slave_cnt_1s":                0,
						"mysql_command_reset_slave_cnt_500ms":             0,
						"mysql_command_reset_slave_cnt_500us":             0,
						"mysql_command_reset_slave_cnt_5ms":               0,
						"mysql_command_reset_slave_cnt_5s":                0,
						"mysql_command_reset_slave_cnt_infs":              0,
						"mysql_command_reset_slave_total_cnt":             0,
						"mysql_command_reset_slave_total_time_us":         0,
						"mysql_command_revoke_cnt_100ms":                  0,
						"mysql_command_revoke_cnt_100us":                  0,
						"mysql_command_revoke_cnt_10ms":                   0,
						"mysql_command_revoke_cnt_10s":                    0,
						"mysql_command_revoke_cnt_1ms":                    0,
						"mysql_command_revoke_cnt_1s":                     0,
						"mysql_command_revoke_cnt_500ms":                  0,
						"mysql_command_revoke_cnt_500us":                  0,
						"mysql_command_revoke_cnt_5ms":                    0,
						"mysql_command_revoke_cnt_5s":                     0,
						"mysql_command_revoke_cnt_infs":                   0,
						"mysql_command_revoke_total_cnt":                  0,
						"mysql_command_revoke_total_time_us":              0,
						"mysql_command_rollback_cnt_100ms":                0,
						"mysql_command_rollback_cnt_100us":                0,
						"mysql_command_rollback_cnt_10ms":                 0,
						"mysql_command_rollback_cnt_10s":                  0,
						"mysql_command_rollback_cnt_1ms":                  0,
						"mysql_command_rollback_cnt_1s":                   0,
						"mysql_command_rollback_cnt_500ms":                0,
						"mysql_command_rollback_cnt_500us":                0,
						"mysql_command_rollback_cnt_5ms":                  0,
						"mysql_command_rollback_cnt_5s":                   0,
						"mysql_command_rollback_cnt_infs":                 0,
						"mysql_command_rollback_total_cnt":                0,
						"mysql_command_rollback_total_time_us":            0,
						"mysql_command_savepoint_cnt_100ms":               0,
						"mysql_command_savepoint_cnt_100us":               0,
						"mysql_command_savepoint_cnt_10ms":                0,
						"mysql_command_savepoint_cnt_10s":                 0,
						"mysql_command_savepoint_cnt_1ms":                 0,
						"mysql_command_savepoint_cnt_1s":                  0,
						"mysql_command_savepoint_cnt_500ms":               0,
						"mysql_command_savepoint_cnt_500us":               0,
						"mysql_command_savepoint_cnt_5ms":                 0,
						"mysql_command_savepoint_cnt_5s":                  0,
						"mysql_command_savepoint_cnt_infs":                0,
						"mysql_command_savepoint_total_cnt":               0,
						"mysql_command_savepoint_total_time_us":           0,
						"mysql_command_select_cnt_100ms":                  4909816,
						"mysql_command_select_cnt_100us":                  32185976,
						"mysql_command_select_cnt_10ms":                   2955830,
						"mysql_command_select_cnt_10s":                    497,
						"mysql_command_select_cnt_1ms":                    481335,
						"mysql_command_select_cnt_1s":                     1321917,
						"mysql_command_select_cnt_500ms":                  11123900,
						"mysql_command_select_cnt_500us":                  36650,
						"mysql_command_select_cnt_5ms":                    4600948,
						"mysql_command_select_cnt_5s":                     403451,
						"mysql_command_select_cnt_infs":                   1870,
						"mysql_command_select_for_update_cnt_100ms":       0,
						"mysql_command_select_for_update_cnt_100us":       0,
						"mysql_command_select_for_update_cnt_10ms":        0,
						"mysql_command_select_for_update_cnt_10s":         0,
						"mysql_command_select_for_update_cnt_1ms":         0,
						"mysql_command_select_for_update_cnt_1s":          0,
						"mysql_command_select_for_update_cnt_500ms":       0,
						"mysql_command_select_for_update_cnt_500us":       0,
						"mysql_command_select_for_update_cnt_5ms":         0,
						"mysql_command_select_for_update_cnt_5s":          0,
						"mysql_command_select_for_update_cnt_infs":        0,
						"mysql_command_select_for_update_total_cnt":       0,
						"mysql_command_select_for_update_total_time_us":   0,
						"mysql_command_select_total_cnt":                  68490650,
						"mysql_command_select_total_time_us":              4673958076637,
						"mysql_command_set_cnt_100ms":                     0,
						"mysql_command_set_cnt_100us":                     0,
						"mysql_command_set_cnt_10ms":                      0,
						"mysql_command_set_cnt_10s":                       0,
						"mysql_command_set_cnt_1ms":                       0,
						"mysql_command_set_cnt_1s":                        0,
						"mysql_command_set_cnt_500ms":                     0,
						"mysql_command_set_cnt_500us":                     0,
						"mysql_command_set_cnt_5ms":                       0,
						"mysql_command_set_cnt_5s":                        0,
						"mysql_command_set_cnt_infs":                      0,
						"mysql_command_set_total_cnt":                     0,
						"mysql_command_set_total_time_us":                 0,
						"mysql_command_show_cnt_100ms":                    0,
						"mysql_command_show_cnt_100us":                    0,
						"mysql_command_show_cnt_10ms":                     0,
						"mysql_command_show_cnt_10s":                      0,
						"mysql_command_show_cnt_1ms":                      0,
						"mysql_command_show_cnt_1s":                       0,
						"mysql_command_show_cnt_500ms":                    0,
						"mysql_command_show_cnt_500us":                    0,
						"mysql_command_show_cnt_5ms":                      1,
						"mysql_command_show_cnt_5s":                       0,
						"mysql_command_show_cnt_infs":                     0,
						"mysql_command_show_table_status_cnt_100ms":       0,
						"mysql_command_show_table_status_cnt_100us":       0,
						"mysql_command_show_table_status_cnt_10ms":        0,
						"mysql_command_show_table_status_cnt_10s":         0,
						"mysql_command_show_table_status_cnt_1ms":         0,
						"mysql_command_show_table_status_cnt_1s":          0,
						"mysql_command_show_table_status_cnt_500ms":       0,
						"mysql_command_show_table_status_cnt_500us":       0,
						"mysql_command_show_table_status_cnt_5ms":         0,
						"mysql_command_show_table_status_cnt_5s":          0,
						"mysql_command_show_table_status_cnt_infs":        0,
						"mysql_command_show_table_status_total_cnt":       0,
						"mysql_command_show_table_status_total_time_us":   0,
						"mysql_command_show_total_cnt":                    1,
						"mysql_command_show_total_time_us":                2158,
						"mysql_command_start_transaction_cnt_100ms":       0,
						"mysql_command_start_transaction_cnt_100us":       0,
						"mysql_command_start_transaction_cnt_10ms":        0,
						"mysql_command_start_transaction_cnt_10s":         0,
						"mysql_command_start_transaction_cnt_1ms":         0,
						"mysql_command_start_transaction_cnt_1s":          0,
						"mysql_command_start_transaction_cnt_500ms":       0,
						"mysql_command_start_transaction_cnt_500us":       0,
						"mysql_command_start_transaction_cnt_5ms":         0,
						"mysql_command_start_transaction_cnt_5s":          0,
						"mysql_command_start_transaction_cnt_infs":        0,
						"mysql_command_start_transaction_total_cnt":       0,
						"mysql_command_start_transaction_total_time_us":   0,
						"mysql_command_truncate_table_cnt_100ms":          0,
						"mysql_command_truncate_table_cnt_100us":          0,
						"mysql_command_truncate_table_cnt_10ms":           0,
						"mysql_command_truncate_table_cnt_10s":            0,
						"mysql_command_truncate_table_cnt_1ms":            0,
						"mysql_command_truncate_table_cnt_1s":             0,
						"mysql_command_truncate_table_cnt_500ms":          0,
						"mysql_command_truncate_table_cnt_500us":          0,
						"mysql_command_truncate_table_cnt_5ms":            0,
						"mysql_command_truncate_table_cnt_5s":             0,
						"mysql_command_truncate_table_cnt_infs":           0,
						"mysql_command_truncate_table_total_cnt":          0,
						"mysql_command_truncate_table_total_time_us":      0,
						"mysql_command_unknown_cnt_100ms":                 0,
						"mysql_command_unknown_cnt_100us":                 0,
						"mysql_command_unknown_cnt_10ms":                  0,
						"mysql_command_unknown_cnt_10s":                   0,
						"mysql_command_unknown_cnt_1ms":                   0,
						"mysql_command_unknown_cnt_1s":                    0,
						"mysql_command_unknown_cnt_500ms":                 0,
						"mysql_command_unknown_cnt_500us":                 0,
						"mysql_command_unknown_cnt_5ms":                   0,
						"mysql_command_unknown_cnt_5s":                    0,
						"mysql_command_unknown_cnt_infs":                  0,
						"mysql_command_unknown_total_cnt":                 0,
						"mysql_command_unknown_total_time_us":             0,
						"mysql_command_unlock_tables_cnt_100ms":           0,
						"mysql_command_unlock_tables_cnt_100us":           0,
						"mysql_command_unlock_tables_cnt_10ms":            0,
						"mysql_command_unlock_tables_cnt_10s":             0,
						"mysql_command_unlock_tables_cnt_1ms":             0,
						"mysql_command_unlock_tables_cnt_1s":              0,
						"mysql_command_unlock_tables_cnt_500ms":           0,
						"mysql_command_unlock_tables_cnt_500us":           0,
						"mysql_command_unlock_tables_cnt_5ms":             0,
						"mysql_command_unlock_tables_cnt_5s":              0,
						"mysql_command_unlock_tables_cnt_infs":            0,
						"mysql_command_unlock_tables_total_cnt":           0,
						"mysql_command_unlock_tables_total_time_us":       0,
						"mysql_command_update_cnt_100ms":                  0,
						"mysql_command_update_cnt_100us":                  0,
						"mysql_command_update_cnt_10ms":                   0,
						"mysql_command_update_cnt_10s":                    0,
						"mysql_command_update_cnt_1ms":                    0,
						"mysql_command_update_cnt_1s":                     0,
						"mysql_command_update_cnt_500ms":                  0,
						"mysql_command_update_cnt_500us":                  0,
						"mysql_command_update_cnt_5ms":                    0,
						"mysql_command_update_cnt_5s":                     0,
						"mysql_command_update_cnt_infs":                   0,
						"mysql_command_update_total_cnt":                  0,
						"mysql_command_update_total_time_us":              0,
						"mysql_command_use_cnt_100ms":                     0,
						"mysql_command_use_cnt_100us":                     0,
						"mysql_command_use_cnt_10ms":                      0,
						"mysql_command_use_cnt_10s":                       0,
						"mysql_command_use_cnt_1ms":                       0,
						"mysql_command_use_cnt_1s":                        0,
						"mysql_command_use_cnt_500ms":                     0,
						"mysql_command_use_cnt_500us":                     0,
						"mysql_command_use_cnt_5ms":                       0,
						"mysql_command_use_cnt_5s":                        0,
						"mysql_command_use_cnt_infs":                      0,
						"mysql_command_use_total_cnt":                     0,
						"mysql_command_use_total_time_us":                 0,
						"mysql_firewall_rules_config":                     329,
						"mysql_firewall_rules_table":                      0,
						"mysql_firewall_users_config":                     0,
						"mysql_firewall_users_table":                      0,
						"mysql_frontend_buffers_bytes":                    196608,
						"mysql_killed_backend_connections":                0,
						"mysql_killed_backend_queries":                    0,
						"mysql_max_allowed_packet":                        4194304,
						"mysql_monitor_connect_check_err":                 130,
						"mysql_monitor_connect_check_ok":                  3548306,
						"mysql_monitor_ping_check_err":                    108271,
						"mysql_monitor_ping_check_ok":                     21289849,
						"mysql_monitor_read_only_check_err":               19610,
						"mysql_monitor_read_only_check_ok":                106246409,
						"mysql_monitor_replication_lag_check_err":         482,
						"mysql_monitor_replication_lag_check_ok":          28702388,
						"mysql_monitor_workers":                           10,
						"mysql_monitor_workers_aux":                       0,
						"mysql_monitor_workers_started":                   10,
						"mysql_query_rules_memory":                        22825,
						"mysql_session_internal_bytes":                    20232,
						"mysql_thread_workers":                            4,
						"mysql_unexpected_frontend_com_quit":              0,
						"mysql_unexpected_frontend_packets":               0,
						"mysql_user_first_user_frontend_connections":      0,
						"mysql_user_first_user_frontend_max_connections":  200,
						"mysql_user_second_user_frontend_connections":     3,
						"mysql_user_second_user_frontend_max_connections": 15,
						"proxysql_uptime":                                 26748286,
						"queries_backends_bytes_recv":                     5896210168,
						"queries_backends_bytes_sent":                     4329581500,
						"queries_frontends_bytes_recv":                    7434816962,
						"queries_frontends_bytes_sent":                    11643634097,
						"queries_with_max_lag_ms":                         0,
						"queries_with_max_lag_ms__delayed":                0,
						"queries_with_max_lag_ms__total_wait_time_us":     0,
						"query_cache_bytes_in":                            0,
						"query_cache_bytes_out":                           0,
						"query_cache_count_get":                           0,
						"query_cache_count_get_ok":                        0,
						"query_cache_count_set":                           0,
						"query_cache_entries":                             0,
						"query_cache_memory_bytes":                        0,
						"query_cache_purged":                              0,
						"query_digest_memory":                             13688,
						"query_processor_time_nsec":                       0,
						"questions":                                       100638067,
						"selects_for_update__autocommit0":                 0,
						"server_connections_aborted":                      9979,
						"server_connections_connected":                    13,
						"server_connections_created":                      2122254,
						"server_connections_delayed":                      0,
						"servers_table_version":                           37,
						"slow_queries":                                    405818,
						"sqlite3_memory_bytes":                            6021248,
						"stack_memory_admin_threads":                      16777216,
						"stack_memory_cluster_threads":                    0,
						"stack_memory_mysql_threads":                      33554432,
						"stmt_cached":                                     65,
						"stmt_client_active_total":                        18,
						"stmt_client_active_unique":                       18,
						"stmt_max_stmt_id":                                66,
						"stmt_server_active_total":                        101,
						"stmt_server_active_unique":                       39,
						"whitelisted_sqli_fingerprint":                    0,
					}

					require.Equal(t, expected, mx)
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
