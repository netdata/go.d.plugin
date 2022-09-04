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
	dataProxySQLV2010Version, _              = os.ReadFile("testdata/v2.0.10/version.txt")
	dataProxySQLV2010GlobalVariables, _      = os.ReadFile("testdata/v2.0.10/global_variables.txt")
	dataProxySQLV2010StatMemorMetrics, _     = os.ReadFile("testdata/v2.0.10/stats_memory_metrics.txt")
	dataProxySQLV2010StatCommandsCounters, _ = os.ReadFile("testdata/v2.0.10/stats_mysql_commands_counters.txt")
	dataProxySQLV2010StatMySQLGlobal, _      = os.ReadFile("testdata/v2.0.10/stats_mysql_global.txt")
	dataProxySQLV2010StatMySQLUsers, _       = os.ReadFile("testdata/v2.0.10/stats_mysql_users.txt")
)

func Test_testDataIsValid(t *testing.T) {
	for name, data := range map[string][]byte{
		"dataProxySQLV2010Version": dataProxySQLV2010Version,
	} {
		require.NotNilf(t, data, fmt.Sprintf("read data: %s", name))
		_, err := prepareMockRows(data)
		require.NoErrorf(t, err, fmt.Sprintf("prepare mock rows: %s", name))
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

func TestMySQL_Charts(t *testing.T) {
	assert.NotNil(t, New().Charts())
}

func TestProxySQL_Check(t *testing.T) {
	tests := map[string]struct {
		prepareMock func(t *testing.T, m sqlmock.Sqlmock)
		wantFail    bool
	}{
		// "success on all queries": {
		// 	wantFail: false,
		// 	prepareMock: func(t *testing.T, m sqlmock.Sqlmock) {
		// 		mockExpect(t, m, queryGlobalVars, dataProxySQLV2010GlobalVariables)
		// 		mockExpect(t, m, queryMemoryMetrics, dataProxySQLV2010StatMemorMetrics)
		// 		mockExpect(t, m, queryMysqlCommandCounters, dataProxySQLV2010StatCommandsCounters)
		// 		mockExpect(t, m, queryMysqlGlobalStatus, dataProxySQLV2010StatMySQLGlobal)
		// 		mockExpect(t, m, queryMysqlusers, dataProxySQLV2010StatMySQLUsers)
		// 	},
		// },
		"fails when error on querying global status": {
			wantFail: true,
			prepareMock: func(t *testing.T, m sqlmock.Sqlmock) {
				mockExpectErr(m, queryGlobalVars)
			},
		},
		"fails when error on querying global variables": {
			wantFail: true,
			prepareMock: func(t *testing.T, m sqlmock.Sqlmock) {
				mockExpect(t, m, queryGlobalVars, dataProxySQLV2010GlobalVariables)
				mockExpectErr(m, queryMemoryMetrics)
			},
		},
		// "success when error on querying memory status": {
		// 	wantFail: false,
		// 	prepareMock: func(t *testing.T, m sqlmock.Sqlmock) {
		// 		mockExpectErr(m, queryGlobalVars)
		// 		mockExpect(t, m, queryMemoryMetrics, dataProxySQLV2010GlobalVariables)
		// 		// mockExpectErr(m, queryMemoryMetrics)
		// 	},
		// },
		// "success when error on querying command counter statistics": {
		// 	wantFail: false,
		// 	prepareMock: func(t *testing.T, m sqlmock.Sqlmock) {
		// 		mockExpect(t, m, queryMysqlCommandCounters, dataProxySQLV2010StatCommandsCounters)
		// 		// mockExpectErr(m, queryMysqlCommandCounters)
		// 	},
		// },
		// "success when error on querying process list": {
		// 	wantFail: false,
		// 	prepareMock: func(t *testing.T, m sqlmock.Sqlmock) {
		// 		mockExpect(t, m, queryMysqlusers, dataProxySQLV2010StatMySQLUsers)
		// 		// mockExpectErr(m, queryMysqlusers)
		// 	},
		// },
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
