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
	l := New("", "")
	b := new(bytes.Buffer)
	r := bufio.NewReader(b)
	l.log.SetOutput(b)
	SetLevel(DEBUG)

	test := func(sev Severity) {
		s, err := r.ReadString('\n')
		assert.Nil(t, err)
		assert.NotEmpty(t, s)
		assert.True(t, strings.Contains(s, sev.String()))
	}

	l.Error("")
	test(ERROR)

	l.Warning("")
	test(WARNING)

	l.Info("")
	test(INFO)

	l.Debug("")
	test(DEBUG)

	l.Errorf("")
	test(ERROR)

	l.Warningf("")
	test(WARNING)

	l.Infof("")
	test(INFO)

	l.Debugf("")
	test(DEBUG)
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
