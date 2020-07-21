package mysql

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"strings"
	"testing"

	"database/sql/driver"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"
)

var (
	globalStatusBitnamiMariaDBv1054, _    = ioutil.ReadFile("testdata/MariaDB-10.5.4-galera-global_status.txt")
	globalVariablesBitnamiMariaDBv1054, _ = ioutil.ReadFile("testdata/MariaDB-10.5.4-galera-global_variables.txt")
	userStatisticsBitnamiMariaDBv1054, _  = ioutil.ReadFile("testdata/MariaDB-10.5.4-galera-user_statistics.txt")
)

func Test_readFile(t *testing.T) {
	require.NotNil(t, globalStatusBitnamiMariaDBv1054)
	require.NotNil(t, globalVariablesBitnamiMariaDBv1054)
	require.NotNil(t, userStatisticsBitnamiMariaDBv1054)
}

func TestNew(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	mock.ExpectQuery("SHOW GLOBAL STATUS").
		WillReturnRows(prepareMockRows(t, globalStatusBitnamiMariaDBv1054))
	mock.ExpectQuery("SHOW GLOBAL VARIABLES").
		WillReturnRows(prepareMockRows(t, globalVariablesBitnamiMariaDBv1054))
	mock.ExpectQuery("SHOW USER_STATISTICS").
		WillReturnRows(prepareMockRows(t, userStatisticsBitnamiMariaDBv1054))

	mySQL := New()
	mySQL.db = db

	mx := mySQL.Collect()
	for key, value := range mx {
		fmt.Println(key, value)
	}
}

func prepareMockRows(t *testing.T, data []byte) *sqlmock.Rows {
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
		require.NotEqual(t, -1, idx)

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
		require.Equal(t, len(columns), len(values))
		rows.AddRow(values...)
	}
	return rows
}
