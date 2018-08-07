package logger

import (
	"bufio"
	"bytes"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSetLevel(t *testing.T) {
	SetLevel(WARNING)
	assert.Equal(t, sevLevel, WARNING)
}

func TestSetModName(t *testing.T) {
	l := New("", "")
	SetModName(l, "name")
	assert.Equal(t, l.modName, "name")
}

func TestSetLimit(t *testing.T) {
	l := New("", "")
	SetLimit(l)
	assert.Len(t, globalTicker.loggers, 1)

}

func TestNew(t *testing.T) {
	assert.IsType(t, (*Logger)(nil), New("", ""))
}

func TestLogger_Log(t *testing.T) {
	logger := New("", "")
	buf := new(bytes.Buffer)
	reader := bufio.NewReader(buf)
	logger.log.SetOutput(buf)
	SetLevel(DEBUG)

	check := func(sev Severity) {
		s, err := reader.ReadString('\n')
		assert.Nil(t, err)
		assert.NotEmpty(t, s)
		assert.True(t, strings.Contains(s, sev.String()))
	}

	logger.Error("")
	check(ERROR)

	logger.Warning("")
	check(WARNING)

	logger.Info("")
	check(INFO)

	logger.Debug("")
	check(DEBUG)

	logger.Errorf("")
	check(ERROR)

	logger.Warningf("")
	check(WARNING)

	logger.Infof("")
	check(INFO)

	logger.Debugf("")
	check(DEBUG)
}

func TestLogger_Limit(t *testing.T) {
	logger := New("", "")
	buf := new(bytes.Buffer)
	scan := bufio.NewScanner(buf)

	logger.log.SetOutput(buf)
	SetLimit(logger)
	SetLevel(DEBUG)

	num := 500

	for i := 0; i < num; i++ {
		logger.Error("")
	}

	var c int
	for scan.Scan() {
		c++
	}

	assert.True(t, num == c)

	SetLevel(INFO)

	for i := 0; i < num; i++ {
		logger.Error("")
	}

	scan = bufio.NewScanner(buf)
	for scan.Scan() {
		c++
	}
	assert.True(t, num+int(msgPerSecond) == c)
}
