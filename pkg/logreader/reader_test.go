package logreader

import (
	"encoding/csv"
	"fmt"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"
	"testing"
	"time"

	"github.com/netdata/go-orchestrator/logger"

	"github.com/stretchr/testify/require"
)

func TestFile_Read(t *testing.T) {
	tmpFileName1 := filepath.Join(os.TempDir(), "test_read.1.log")
	tmpFileName2 := filepath.Join(os.TempDir(), "test_read.2.log")
	tmpFileName3 := filepath.Join(os.TempDir(), "test_read.3.log")
	defer os.Remove(tmpFileName1)
	defer os.Remove(tmpFileName2)
	defer os.Remove(tmpFileName3)

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		writeLog(t, tmpFileName1, time.Millisecond*10)
		writeLog(t, tmpFileName2, time.Millisecond*10)
		writeLog(t, tmpFileName3, time.Millisecond*10)
	}()

	time.Sleep(time.Millisecond * 10)

	go func() {
		defer wg.Done()
		readLog(t)
	}()

	wg.Wait()
}

func TestFile_Read2(t *testing.T) {
	tmpFileName1 := filepath.Join(os.TempDir(), "test_read.1.log")
	tmpFileName2 := filepath.Join(os.TempDir(), "test_read.1.log")
	tmpFileName3 := filepath.Join(os.TempDir(), "test_read.1.log")
	defer os.Remove(tmpFileName1)
	defer os.Remove(tmpFileName2)
	defer os.Remove(tmpFileName3)

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		writeLog(t, tmpFileName1, time.Millisecond*10)
		writeLog(t, tmpFileName2, time.Millisecond*10)
		writeLog(t, tmpFileName3, time.Millisecond*10)
	}()

	time.Sleep(time.Millisecond * 10)

	go func() {
		defer wg.Done()
		readLog(t)
	}()

	wg.Wait()
}

func readLog(t *testing.T) {
	t.Helper()
	file, err := Open(filepath.Join(os.TempDir(), "test_read.*.log"), "", logger.New("go.d", "web_log", "test"))
	require.NoError(t, err)
	defer file.Close()
	r := csv.NewReader(file)
	r.Comma = ' '
	r.ReuseRecord = true
	r.FieldsPerRecord = -1
	for i := 0; i < 50; i++ {
		record, err := r.Read()
		if err == nil {
			fmt.Printf("[%d] line:  %v\n", i, record)
		} else {
			fmt.Printf("[%d] error: %v\n", i, err)
		}
		time.Sleep(time.Millisecond * 15)
	}
}

func writeLog(t *testing.T, filename string, interval time.Duration) {
	t.Helper()
	base := filepath.Base(filename)
	file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	require.NoError(t, err)
	require.NotNil(t, file)
	defer file.Close()

	for i := 0; i < 15; i++ {
		fmt.Fprintln(file, "line", i, "filename", base)
		time.Sleep(interval)
	}
}

func TestReadLastLine(t *testing.T) {
	tests := []struct {
		name     string
		content  string
		expected string
		err      error
	}{
		{"empty", "", "", nil},
		{"empty-ln", "\n", "\n", nil},
		{"one-line", "hello", "hello", nil},
		{"one-line-ln", "hello\n", "hello\n", nil},
		{"multi-line", "hello\nworld", "world", nil},
		{"multi-line-ln", "hello\nworld\n", "world\n", nil},
		{"long-line", "hello hello hello", "", ErrTooLongLine},
		{"long-line-ln", "hello hello hello\n", "", ErrTooLongLine},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			filename := prepareFile(t, test.content)
			defer os.Remove(filename)
			line, err := ReadLastLine(filename, 10)
			if test.err != nil {
				assert.Contains(t, err.Error(), test.err.Error())
			} else {
				assert.Equal(t, test.expected, string(line))
			}
		})
	}
}

func prepareFile(t *testing.T, content string) string {
	t.Helper()
	file, err := ioutil.TempFile("", "go-test")
	require.NoError(t, err)
	defer file.Close()

	_, _ = file.WriteString(content)
	return file.Name()
}
