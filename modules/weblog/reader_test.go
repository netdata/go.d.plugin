package weblog

import (
	"fmt"
	"io"
	"os"
	"testing"

	"github.com/fsnotify/fsnotify"

	"github.com/stretchr/testify/assert"

	"github.com/stretchr/testify/require"
)

func TestReader_Read(t *testing.T) {
	const filename = "./tmp/test1"
	os.Mkdir("./tmp", 0755)
	defer os.RemoveAll("./tmp")

	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_APPEND|os.O_SYNC, 0644)
	require.NoError(t, err)

	watcher, err := fsnotify.NewWatcher()
	require.NoError(t, err)

	watcher.Add(filename)
	defer watcher.Close()

	go func() {
		for e := range watcher.Events {
			fmt.Printf("%v\n", e)
			if e.Op&fsnotify.Rename != 0 {
			}
		}
		fmt.Println("end watch event")
	}()
	go func() {
		for e := range watcher.Errors {
			fmt.Printf("%v\n", e)
		}
		fmt.Println("end watch error")
	}()

	r := NewReader(filename)
	buf := make([]byte, 80)
	n, err := r.Read(buf)
	assert.Equal(t, 0, n)
	assert.Equal(t, io.EOF, err)

	file.WriteString("hello\n")
	file.Close()

	n, err = r.Read(buf)
	assert.Equal(t, 6, n)
	assert.NoError(t, err)

	n, err = r.Read(buf)
	assert.Equal(t, 0, n)
	assert.Equal(t, io.EOF, err)

	os.Rename(filename, filename+".old")
	file, err = os.OpenFile(filename+".old", os.O_WRONLY|os.O_CREATE|os.O_APPEND|os.O_SYNC, 0644)
	require.NoError(t, err)
	file.WriteString("hello\n")
	file.Close()

	n, err = r.Read(buf)
	assert.Equal(t, 6, n)
	assert.NoError(t, err)

	n, err = r.Read(buf)
	assert.Equal(t, 0, n)
	assert.Equal(t, io.EOF, err)

	os.Remove(filename + ".old")

	n, err = r.Read(buf)
	assert.Equal(t, 0, n)
	assert.Equal(t, io.EOF, err)

	//n, err = r.Read(buf)
	//assert.Equal(t, 6, n)
	//assert.NoError(t, err)
}
