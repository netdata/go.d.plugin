package logger

import (
	"bufio"
	"bytes"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSetGlobLevel(t *testing.T) {
	SetGlobalSeverity(WARNING)
	assert.Equal(t, globalSeverity, WARNING)
}

func TestNew(t *testing.T) {
	assert.IsType(t, (*Logger)(nil), New("", ""))
}

func TestLogger_Log(t *testing.T) {
	logger := New("", "")
	buf := new(bytes.Buffer)
	reader := bufio.NewReader(buf)
	logger.log.SetOutput(buf)
	SetGlobalSeverity(DEBUG)

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
	logger := NewLimited("", "")
	buf := new(bytes.Buffer)
	scan := bufio.NewScanner(buf)

	logger.log.SetOutput(buf)
	SetGlobalSeverity(DEBUG)

	num := 500

	for i := 0; i < num; i++ {
		logger.Error("")
	}

	var c int
	for scan.Scan() {
		c++
	}

	assert.True(t, num == c)

	SetGlobalSeverity(INFO)

	for i := 0; i < num; i++ {
		logger.Error("")
	}

	scan = bufio.NewScanner(buf)
	for scan.Scan() {
		c++
	}
	assert.True(t, num+msgPerSecondLimit == c)
}
