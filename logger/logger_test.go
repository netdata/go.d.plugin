// SPDX-License-Identifier: GPL-3.0-or-later

package logger

import (
	"bytes"
	"io"
	"log"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSetSeverity(t *testing.T) {
	require.Equal(t, globalSeverity, INFO)
	SetSeverity(DEBUG)

	assert.Equal(t, globalSeverity, DEBUG)
}

func TestNew(t *testing.T) {
	assert.IsType(
		t,
		(*Logger)(nil),
		New("", ""),
	)
}

func TestNewLimited(t *testing.T) {
	logger := NewLimited("", "")
	assert.True(t, logger.limited)

	_, ok := GlobalMsgCountWatcher.items[logger.id]
	require.True(t, ok)
	GlobalMsgCountWatcher.Unregister(logger)
}

func TestLogger_Critical(t *testing.T) {
	buf := bytes.Buffer{}
	logger := New("", "")
	logger.formatter.SetOutput(&buf)
	logger.formatter.flag = log.Lshortfile
	logger.Critical()
	assert.Contains(t, buf.String(), CRITICAL.ShortString())
	assert.Contains(t, buf.String(), " logger_test.go")
}

func TestLogger_Criticalf(t *testing.T) {
	buf := bytes.Buffer{}
	logger := New("", "")
	logger.formatter.SetOutput(&buf)
	logger.formatter.flag = log.Lshortfile
	logger.Criticalf("")
	assert.Contains(t, buf.String(), CRITICAL.ShortString())
	assert.Contains(t, buf.String(), " logger_test.go")
}

func TestLogger_Error(t *testing.T) {
	buf := bytes.Buffer{}
	logger := New("", "")
	logger.formatter.SetOutput(&buf)

	logger.Error()
	assert.Contains(t, buf.String(), ERROR.ShortString())
}

func TestLogger_Errorf(t *testing.T) {
	buf := bytes.Buffer{}
	logger := New("", "")
	logger.formatter.SetOutput(&buf)

	logger.Errorf("")
	assert.Contains(t, buf.String(), ERROR.ShortString())
}

func TestLogger_Warning(t *testing.T) {
	buf := bytes.Buffer{}
	logger := New("", "")
	logger.formatter.SetOutput(&buf)

	logger.Warning()
	assert.Contains(t, buf.String(), WARNING.ShortString())
}

func TestLogger_Warningf(t *testing.T) {
	buf := bytes.Buffer{}
	logger := New("", "")
	logger.formatter.SetOutput(&buf)

	logger.Warningf("")
	assert.Contains(t, buf.String(), WARNING.ShortString())
}

func TestLogger_Info(t *testing.T) {
	buf := bytes.Buffer{}
	logger := New("", "")
	logger.formatter.SetOutput(&buf)

	logger.Info()
	assert.Contains(t, buf.String(), INFO.ShortString())
}

func TestLogger_Infof(t *testing.T) {
	buf := bytes.Buffer{}
	logger := New("", "")
	logger.formatter.SetOutput(&buf)

	logger.Infof("")
	assert.Contains(t, buf.String(), INFO.ShortString())
}

func TestLogger_Debug(t *testing.T) {
	buf := bytes.Buffer{}
	logger := New("", "")
	logger.formatter.SetOutput(&buf)

	logger.Debug()
	assert.Contains(t, buf.String(), DEBUG.ShortString())
}

func TestLogger_Debugf(t *testing.T) {
	buf := bytes.Buffer{}
	logger := New("", "")
	logger.formatter.SetOutput(&buf)

	logger.Debugf("")
	assert.Contains(t, buf.String(), DEBUG.ShortString())
}

func TestLogger_NotInitialized(t *testing.T) {
	var logger Logger
	f := func() {
		logger.Info()
	}
	assert.NotPanics(t, f)
}

func TestLogger_NotInitializedPtr(t *testing.T) {
	var logger *Logger
	f := func() {
		logger.Info()
	}
	assert.NotPanics(t, f)
}

func TestLogger_Unlimited(t *testing.T) {
	logger := New("", "")

	wr := countWriter(0)
	logger.formatter.SetOutput(&wr)

	num := 1000

	for i := 0; i < num; i++ {
		logger.Info()
	}

	require.Equal(t, num, int(wr))
}

func TestLogger_Limited(t *testing.T) {
	SetSeverity(INFO)

	logger := New("", "")
	logger.limited = true

	wr := countWriter(0)
	logger.formatter.SetOutput(&wr)

	num := 1000

	for i := 0; i < num; i++ {
		logger.Info()
	}

	require.Equal(t, msgPerSecondLimit, int(wr))
}

func TestLogger_Info_race(t *testing.T) {
	logger := New("", "")
	logger.formatter.SetOutput(io.Discard)
	for i := 0; i < 10; i++ {
		go func() {
			for j := 0; j < 10; j++ {
				logger.Info("hello ", "world")
			}
		}()
	}
	time.Sleep(time.Second)
}

type countWriter int

func (c *countWriter) Write(b []byte) (n int, err error) {
	*c++
	return len(b), nil
}

func BenchmarkLogger_Infof(b *testing.B) {
	l := New("test", "test")
	l.formatter.SetOutput(io.Discard)
	for i := 0; i < b.N; i++ {
		l.Infof("hello %s", "world")
	}
}

func BenchmarkLog_Printf(b *testing.B) {
	logger := log.New(io.Discard, "", log.Lshortfile)
	for i := 0; i < b.N; i++ {
		logger.Printf("hello %s", "world")
	}
}
