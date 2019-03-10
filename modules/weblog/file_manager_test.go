package weblog

import (
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/stretchr/testify/require"
)

func TestFile_Read(t *testing.T) {
	tmpFileName1 := filepath.Join(os.TempDir(), "test_read.1.log")
	tmpFileName2 := filepath.Join(os.TempDir(), "test_read.2.log")
	tmpFileName3 := filepath.Join(os.TempDir(), "test_read.3.log")

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		defer assert.NoError(t, os.Remove(tmpFileName1))
		defer assert.NoError(t, os.Remove(tmpFileName2))
		defer assert.NoError(t, os.Remove(tmpFileName3))

		writeLog(t, tmpFileName1, time.Millisecond*10)
		writeLog(t, tmpFileName2, time.Millisecond*10)
		writeLog(t, tmpFileName3, time.Millisecond*10)
	}()

	time.Sleep(time.Millisecond * 10)

	go func() {
		defer wg.Done()
		mgr, err := NewFileManager(filepath.Join(os.TempDir(), "test_read.*.log"), "")
		require.NoError(t, err)
		file, err := mgr.OpenFile()
		require.NoError(t, err)
		r := csv.NewReader(file)
		r.Comma = ' '
		r.ReuseRecord = true
		r.FieldsPerRecord = -1
		for i := 0; i < 31; i++ {
			record, err := r.Read()
			if err == nil {
				fmt.Printf("[%d] line:  %v\n", i, record)
			} else {
				fmt.Printf("[%d] error: %v\n", i, err)
			}
			time.Sleep(time.Millisecond * 20)
		}
	}()

	wg.Wait()
}

func TestFile_Read2(t *testing.T) {
	tmpFileName1 := filepath.Join(os.TempDir(), "test_read.1.log")

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		defer assert.NoError(t, os.Remove(tmpFileName1))

		writeLog(t, tmpFileName1, time.Millisecond*10)
		//os.Truncate(tmpFileName1, 0)
		//writeLog(t, tmpFileName1, time.Millisecond*10)
		//os.Truncate(tmpFileName1, 0)
		//writeLog(t, tmpFileName1, time.Millisecond*10)
		//os.Truncate(tmpFileName1, 0)
	}()

	time.Sleep(time.Millisecond * 10)

	go func() {
		defer wg.Done()
		mgr, err := NewFileManager(filepath.Join(os.TempDir(), "test_read.*.log"), "")
		require.NoError(t, err)
		file, err := mgr.OpenFile()
		require.NoError(t, err)
		r := csv.NewReader(file)
		r.Comma = ' '
		r.ReuseRecord = true
		r.FieldsPerRecord = -1
		for i := 0; i < 31; i++ {
			record, err := r.Read()
			if err == nil {
				fmt.Printf("[%d] line:  %v\n", i, record)
			} else {
				fmt.Printf("[%d] error: %v\n", i, err)
			}
			time.Sleep(time.Millisecond * 20)
		}
	}()

	wg.Wait()
}

func writeLog(t *testing.T, filename string, interval time.Duration) {
	t.Helper()
	base := filepath.Base(filename)
	file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, 0600)
	require.NoError(t, err)
	require.NotNil(t, file)
	defer file.Close()

	for i := 0; i < 10; i++ {
		fmt.Fprintln(file, "line", i, "filename", base)
		time.Sleep(interval)
	}
}
