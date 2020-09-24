package logger

import (
	"io/ioutil"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMsgCountWatcher_Register(t *testing.T) {
	cw := newMsgCountWatcher(time.Second)
	defer cw.stop()

	require.Len(t, cw.items, 0)

	logger := New("", "")
	cw.Register(logger)

	require.Len(t, cw.items, 1)
	require.Equal(t, logger, cw.items[logger.id])

}

func TestMsgCountWatcher_Unregister(t *testing.T) {
	cw := newMsgCountWatcher(time.Second)
	defer cw.stop()

	require.Len(t, cw.items, 0)

	logger := New("", "")
	cw.items[logger.id] = logger
	cw.Unregister(logger)

	require.Len(t, cw.items, 0)
}

func TestMsgCountWatcher(t *testing.T) {
	reset := time.Millisecond * 500
	cw := newMsgCountWatcher(reset)
	defer cw.stop()

	logger := New("", "")
	logger.limited = true
	logger.formatter.SetOutput(ioutil.Discard)
	cw.Register(logger)

	for i := 0; i < 3; i++ {
		for m := 0; m < 100; m++ {
			logger.Info()
		}
		time.Sleep(reset * 2)
		assert.Equal(t, int64(0), atomic.LoadInt64(&logger.msgCount))
	}
}
