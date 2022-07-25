package postgres

import (
	"bufio"
	"bytes"
	"database/sql/driver"
	"errors"
	"fmt"
	"io/ioutil"
	"strings"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	dataV140004ServerVersionNum, _   = ioutil.ReadFile("testdata/v14.4/server_version_num.txt")
	dataV140004IsSuperUserFalse, _   = ioutil.ReadFile("testdata/v14.4/is_super_user-false.txt")
	dataV140004IsSuperUserTrue, _    = ioutil.ReadFile("testdata/v14.4/is_super_user-true.txt")
	dataV140004DatabasesList1DB, _   = ioutil.ReadFile("testdata/v14.4/databases_list-1db.txt")
	dataV140004DatabasesList2DB, _   = ioutil.ReadFile("testdata/v14.4/databases_list-2db.txt")
	dataV140004DatabasesList3DB, _   = ioutil.ReadFile("testdata/v14.4/databases_list-3db.txt")
	dataV140004DatabasesStats, _     = ioutil.ReadFile("testdata/v14.4/databases_stats.txt")
	dataV140004DatabasesConflicts, _ = ioutil.ReadFile("testdata/v14.4/databases_conflicts.txt")
	dataV140004Checkpoints, _        = ioutil.ReadFile("testdata/v14.4/checkpoints.txt")
)

func Test_testDataIsValid(t *testing.T) {
	for name, data := range map[string][]byte{
		"dataV140004ServerVersionNum":   dataV140004ServerVersionNum,
		"dataV140004IsSuperUserFalse":   dataV140004IsSuperUserFalse,
		"dataV140004IsSuperUserTrue":    dataV140004IsSuperUserTrue,
		"dataV140004DatabasesList1DB":   dataV140004DatabasesList1DB,
		"dataV140004DatabasesList2DB":   dataV140004DatabasesList2DB,
		"dataV140004DatabasesList3DB":   dataV140004DatabasesList3DB,
		"dataV140004DatabasesStats":     dataV140004DatabasesStats,
		"dataV140004DatabasesConflicts": dataV140004DatabasesConflicts,
		"dataV140004Checkpoints":        dataV140004Checkpoints,
	} {
		require.NotNilf(t, data, name)
	}
}

func TestPostgres_Init(t *testing.T) {
	tests := map[string]struct {
		wantFail bool
		config   Config
	}{
		"Success with default": {
			wantFail: false,
			config:   New().Config,
		},
		"Fail when DSN not set": {
			wantFail: true,
			config:   Config{DSN: ""},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			apache := New()
			apache.Config = test.config

			if test.wantFail {
				assert.False(t, apache.Init())
			} else {
				assert.True(t, apache.Init())
			}
		})
	}
}

func TestPostgres_Cleanup(t *testing.T) {

}

func TestPostgres_Charts(t *testing.T) {
	assert.NotNil(t, New().Charts())
}

func TestPostgres_Check(t *testing.T) {
	dbs := []string{"postgres", "production"}
	tests := map[string]struct {
		prepareMock func(t *testing.T, mock sqlmock.Sqlmock)
		wantFail    bool
	}{
		"Success when all queries are successful (v14.4)": {
			wantFail: false,
			prepareMock: func(t *testing.T, mock sqlmock.Sqlmock) {
				mock.ExpectQuery(queryServerVersion()).
					WillReturnRows(mustMockRows(t, dataV140004ServerVersionNum)).RowsWillBeClosed()
				mock.ExpectQuery(queryDatabasesList()).
					WillReturnRows(mustMockRows(t, dataV140004DatabasesList2DB)).RowsWillBeClosed()
				mock.ExpectQuery(queryDatabasesStats(dbs)).
					WillReturnRows(mustMockRows(t, dataV140004DatabasesStats)).RowsWillBeClosed()
				mock.ExpectQuery(queryDatabasesConflicts(dbs)).
					WillReturnRows(mustMockRows(t, dataV140004DatabasesConflicts)).RowsWillBeClosed()
				mock.ExpectQuery(queryCheckpoints()).
					WillReturnRows(mustMockRows(t, dataV140004Checkpoints)).RowsWillBeClosed()
			},
		},
		"Success when the first query is successful (v14.4)": {
			wantFail: false,
			prepareMock: func(t *testing.T, mock sqlmock.Sqlmock) {
				mock.ExpectQuery(queryServerVersion()).
					WillReturnRows(mustMockRows(t, dataV140004ServerVersionNum)).RowsWillBeClosed()
				mock.ExpectQuery(queryDatabasesList()).
					WillReturnRows(mustMockRows(t, dataV140004DatabasesList2DB)).RowsWillBeClosed()
				mock.ExpectQuery(queryDatabasesStats(dbs)).
					WillReturnRows(mustMockRows(t, dataV140004DatabasesStats)).RowsWillBeClosed()
				mock.ExpectQuery(queryDatabasesConflicts(dbs)).
					WillReturnError(errors.New("mock queryDatabasesConflicts() error"))
			},
		},
		"Fail when querying the database version returns an error": {
			wantFail: true,
			prepareMock: func(t *testing.T, mock sqlmock.Sqlmock) {
				mock.ExpectQuery(queryServerVersion()).
					WillReturnError(errors.New("mock queryServerVersion() error"))
			},
		},
		"Fail when querying the databases list returns an error": {
			wantFail: true,
			prepareMock: func(t *testing.T, mock sqlmock.Sqlmock) {
				mock.ExpectQuery(queryServerVersion()).
					WillReturnRows(mustMockRows(t, dataV140004ServerVersionNum)).RowsWillBeClosed()
				mock.ExpectQuery(queryDatabasesList()).
					WillReturnError(errors.New("mock queryDatabaseList() error"))
			},
		},
		"Fail when querying the databases stats returns an error": {
			wantFail: true,
			prepareMock: func(t *testing.T, mock sqlmock.Sqlmock) {
				mock.ExpectQuery(queryServerVersion()).
					WillReturnRows(mustMockRows(t, dataV140004ServerVersionNum)).RowsWillBeClosed()
				mock.ExpectQuery(queryDatabasesList()).
					WillReturnRows(mustMockRows(t, dataV140004DatabasesList2DB)).RowsWillBeClosed()
				mock.ExpectQuery(queryDatabasesStats(dbs)).
					WillReturnError(errors.New("mock queryDatabasesStats() error"))
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			db, mock, err := sqlmock.New(
				sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual),
			)
			require.NoError(t, err)
			pg := New()
			pg.db = db
			defer func() { _ = db.Close() }()

			require.True(t, pg.Init())

			test.prepareMock(t, mock)

			if test.wantFail {
				assert.False(t, pg.Check())
			} else {
				assert.True(t, pg.Check())
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestPostgres_Collect(t *testing.T) {
	dbs1 := []string{"postgres"}
	dbs2 := []string{"postgres", "production"}
	dbs3 := []string{"postgres", "production", "staging"}
	_ = dbs3
	type testCaseStep struct {
		prepareMock func(t *testing.T, mock sqlmock.Sqlmock)
		check       func(t *testing.T, pg *Postgres)
	}
	tests := map[string][]testCaseStep{
		"Success on all queries (v14.4)": {
			{
				prepareMock: func(t *testing.T, mock sqlmock.Sqlmock) {
					mock.ExpectQuery(queryServerVersion()).
						WillReturnRows(mustMockRows(t, dataV140004ServerVersionNum)).RowsWillBeClosed()
					mock.ExpectQuery(queryDatabasesList()).
						WillReturnRows(mustMockRows(t, dataV140004DatabasesList2DB)).RowsWillBeClosed()
					mock.ExpectQuery(queryDatabasesStats(dbs2)).
						WillReturnRows(mustMockRows(t, dataV140004DatabasesStats)).RowsWillBeClosed()
					mock.ExpectQuery(queryDatabasesConflicts(dbs2)).
						WillReturnRows(mustMockRows(t, dataV140004DatabasesConflicts)).RowsWillBeClosed()
					mock.ExpectQuery(queryCheckpoints()).
						WillReturnRows(mustMockRows(t, dataV140004Checkpoints)).RowsWillBeClosed()
				},
				check: func(t *testing.T, pg *Postgres) {
					mx := pg.Collect()

					expected := map[string]int64{
						"buffers_alloc":                  27295744,
						"buffers_backend":                0,
						"buffers_backend_fsync":          0,
						"buffers_checkpoint":             32768,
						"buffers_clean":                  0,
						"checkpoint_sync_time":           47,
						"checkpoint_write_time":          167,
						"checkpoints_req":                16,
						"checkpoints_timed":              1814,
						"db_postgres_blks_hit":           9245,
						"db_postgres_blks_read":          246,
						"db_postgres_confl_bufferpin":    0,
						"db_postgres_confl_deadlock":     0,
						"db_postgres_confl_lock":         0,
						"db_postgres_confl_snapshot":     0,
						"db_postgres_confl_tablespace":   0,
						"db_postgres_conflicts":          0,
						"db_postgres_deadlocks":          0,
						"db_postgres_numbackends":        2,
						"db_postgres_size":               8758051,
						"db_postgres_temp_bytes":         0,
						"db_postgres_temp_files":         0,
						"db_postgres_tup_deleted":        0,
						"db_postgres_tup_fetched":        3577,
						"db_postgres_tup_inserted":       0,
						"db_postgres_tup_returned":       65095,
						"db_postgres_tup_updated":        0,
						"db_postgres_xact_commit":        1636,
						"db_postgres_xact_rollback":      2,
						"db_production_blks_hit":         0,
						"db_production_blks_read":        0,
						"db_production_confl_bufferpin":  0,
						"db_production_confl_deadlock":   0,
						"db_production_confl_lock":       0,
						"db_production_confl_snapshot":   0,
						"db_production_confl_tablespace": 0,
						"db_production_conflicts":        0,
						"db_production_deadlocks":        0,
						"db_production_numbackends":      0,
						"db_production_size":             8602115,
						"db_production_temp_bytes":       0,
						"db_production_temp_files":       0,
						"db_production_tup_deleted":      0,
						"db_production_tup_fetched":      0,
						"db_production_tup_inserted":     0,
						"db_production_tup_returned":     0,
						"db_production_tup_updated":      0,
						"db_production_xact_commit":      0,
						"db_production_xact_rollback":    0,
						"maxwritten_clean":               0,
					}
					assert.Equal(t, expected, mx)
				},
			},
		},
		"DB removed/added on relisting databases": {
			{
				prepareMock: func(t *testing.T, mock sqlmock.Sqlmock) {
					mock.ExpectQuery(queryServerVersion()).
						WillReturnRows(mustMockRows(t, dataV140004ServerVersionNum)).RowsWillBeClosed()
					mock.ExpectQuery(queryDatabasesList()).
						WillReturnRows(mustMockRows(t, dataV140004DatabasesList2DB)).RowsWillBeClosed()
					mock.ExpectQuery(queryDatabasesStats(dbs2)).
						WillReturnRows(mustMockRows(t, dataV140004DatabasesStats)).RowsWillBeClosed()
					mock.ExpectQuery(queryDatabasesConflicts(dbs2)).
						WillReturnRows(mustMockRows(t, dataV140004DatabasesConflicts)).RowsWillBeClosed()
					mock.ExpectQuery(queryCheckpoints()).
						WillReturnRows(mustMockRows(t, dataV140004Checkpoints)).RowsWillBeClosed()
				},
				check: func(t *testing.T, pg *Postgres) { _ = pg.Collect() },
			},
			{
				prepareMock: func(t *testing.T, mock sqlmock.Sqlmock) {
					mock.ExpectQuery(queryDatabasesList()).
						WillReturnRows(mustMockRows(t, dataV140004DatabasesList1DB)).RowsWillBeClosed()
					mock.ExpectQuery(queryDatabasesStats(dbs1)).
						WillReturnRows(mustMockRows(t, dataV140004DatabasesStats)).RowsWillBeClosed()
					mock.ExpectQuery(queryDatabasesConflicts(dbs1)).
						WillReturnRows(mustMockRows(t, dataV140004DatabasesConflicts)).RowsWillBeClosed()
					mock.ExpectQuery(queryCheckpoints()).
						WillReturnRows(mustMockRows(t, dataV140004Checkpoints)).RowsWillBeClosed()
				},
				check: func(t *testing.T, pg *Postgres) {
					pg.relistDatabasesEvery = time.Second
					time.Sleep(time.Second)
					_ = pg.Collect()
					assert.Len(t, pg.databases, 1)
				},
			},
			{
				prepareMock: func(t *testing.T, mock sqlmock.Sqlmock) {
					mock.ExpectQuery(queryDatabasesList()).
						WillReturnRows(mustMockRows(t, dataV140004DatabasesList3DB)).RowsWillBeClosed()
					mock.ExpectQuery(queryDatabasesStats(dbs3)).
						WillReturnRows(mustMockRows(t, dataV140004DatabasesStats)).RowsWillBeClosed()
					mock.ExpectQuery(queryDatabasesConflicts(dbs3)).
						WillReturnRows(mustMockRows(t, dataV140004DatabasesConflicts)).RowsWillBeClosed()
					mock.ExpectQuery(queryCheckpoints()).
						WillReturnRows(mustMockRows(t, dataV140004Checkpoints)).RowsWillBeClosed()
				},
				check: func(t *testing.T, pg *Postgres) {
					pg.relistDatabasesEvery = time.Second
					time.Sleep(time.Second)
					_ = pg.Collect()
					assert.Len(t, pg.databases, 3)
				},
			},
		},
		"Fail when querying the database version returns an error": {
			{
				prepareMock: func(t *testing.T, mock sqlmock.Sqlmock) {
					mock.ExpectQuery(queryServerVersion()).
						WillReturnError(errors.New("mock queryServerVersion() error"))
				},
				check: func(t *testing.T, pg *Postgres) {
					mx := pg.Collect()
					var excepted map[string]int64
					assert.Equal(t, excepted, mx)
				},
			},
		},
		"Fail when querying the databases list returns an error": {
			{
				prepareMock: func(t *testing.T, mock sqlmock.Sqlmock) {
					mock.ExpectQuery(queryServerVersion()).
						WillReturnRows(mustMockRows(t, dataV140004ServerVersionNum)).RowsWillBeClosed()
					mock.ExpectQuery(queryDatabasesList()).
						WillReturnError(errors.New("mock queryDatabaseList() error"))
				},
				check: func(t *testing.T, pg *Postgres) {
					mx := pg.Collect()
					var excepted map[string]int64
					assert.Equal(t, excepted, mx)
				},
			},
		},
		"Fail when querying the databases stats returns an error": {
			{
				prepareMock: func(t *testing.T, mock sqlmock.Sqlmock) {
					mock.ExpectQuery(queryServerVersion()).
						WillReturnRows(mustMockRows(t, dataV140004ServerVersionNum)).RowsWillBeClosed()
					mock.ExpectQuery(queryDatabasesList()).
						WillReturnRows(mustMockRows(t, dataV140004DatabasesList2DB)).RowsWillBeClosed()
					mock.ExpectQuery(queryDatabasesStats(dbs2)).
						WillReturnError(errors.New("mock queryDatabasesStats() error"))
				},
				check: func(t *testing.T, pg *Postgres) {
					mx := pg.Collect()
					var excepted map[string]int64
					assert.Equal(t, excepted, mx)
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
			pg := New()
			pg.db = db
			defer func() { _ = db.Close() }()

			require.True(t, pg.Init())

			for i, step := range test {
				t.Run(fmt.Sprintf("step[%d]", i), func(t *testing.T) {
					step.prepareMock(t, mock)
					step.check(t, pg)
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

func prepareMockRows(data []byte) (*sqlmock.Rows, error) {
	r := bytes.NewReader(data)
	sc := bufio.NewScanner(r)

	var numColumns int
	var rows *sqlmock.Rows

	for sc.Scan() {
		s := strings.TrimSpace(sc.Text())
		if s == "" || strings.HasPrefix(s, "---") {
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

	return rows, nil
}
