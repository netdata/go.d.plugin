package logger

import (
	"bytes"
	"strings"
	"testing"

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
	delete(GlobalMsgCountWatcher.items, logger.id)
}

func TestLogger_Critical(t *testing.T) {
	buf := bytes.Buffer{}
	logger := New("", "")
	logger.log.SetOutput(&buf)

	logger.Critical()
	assert.True(t, strings.Contains(buf.String(), CRITICAL.String()))
}

func TestLogger_Criticalf(t *testing.T) {
	buf := bytes.Buffer{}
	logger := New("", "")
	logger.log.SetOutput(&buf)

	logger.Criticalf("")
	assert.True(t, strings.Contains(buf.String(), CRITICAL.String()))
}

func TestLogger_Error(t *testing.T) {
	buf := bytes.Buffer{}
	logger := New("", "")
	logger.log.SetOutput(&buf)

	logger.Error()
	assert.True(t, strings.Contains(buf.String(), ERROR.String()))
}

func TestLogger_Errorf(t *testing.T) {
	buf := bytes.Buffer{}
	logger := New("", "")
	logger.log.SetOutput(&buf)

	logger.Errorf("")
	assert.True(t, strings.Contains(buf.String(), ERROR.String()))
}

func TestLogger_Warning(t *testing.T) {
	buf := bytes.Buffer{}
	logger := New("", "")
	logger.log.SetOutput(&buf)

	logger.Warning()
	assert.True(t, strings.Contains(buf.String(), WARNING.String()))
}

func TestLogger_Warningf(t *testing.T) {
	buf := bytes.Buffer{}
	logger := New("", "")
	logger.log.SetOutput(&buf)

	logger.Warningf("")
	assert.True(t, strings.Contains(buf.String(), WARNING.String()))
}

func TestLogger_Info(t *testing.T) {
	buf := bytes.Buffer{}
	logger := New("", "")
	logger.log.SetOutput(&buf)

	logger.Info()
	assert.True(t, strings.Contains(buf.String(), INFO.String()))
}

func TestLogger_Infof(t *testing.T) {
	buf := bytes.Buffer{}
	logger := New("", "")
	logger.log.SetOutput(&buf)

	logger.Infof("")
	assert.True(t, strings.Contains(buf.String(), INFO.String()))
}

func TestLogger_Debug(t *testing.T) {
	buf := bytes.Buffer{}
	logger := New("", "")
	logger.log.SetOutput(&buf)

	logger.Debug()
	assert.True(t, strings.Contains(buf.String(), DEBUG.String()))
}

func TestLogger_Debugf(t *testing.T) {
	buf := bytes.Buffer{}
	logger := New("", "")
	logger.log.SetOutput(&buf)

	logger.Debugf("")
	assert.True(t, strings.Contains(buf.String(), DEBUG.String()))
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
	logger.log.SetOutput(&wr)

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
	logger.log.SetOutput(&wr)

	num := 1000

	for i := 0; i < num; i++ {
		logger.Info()
	}

	require.Equal(t, msgPerSecondLimit, int(wr))
}

type countWriter int

func (c *countWriter) Write(b []byte) (n int, err error) {
	*c++
	return len(b), nil
}
