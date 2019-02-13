package weblog

import (
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/stretchr/testify/assert"
)

var logPatternCases = []struct {
	name    string
	pattern *LogPattern
}{
	{"common", logFmtCommon},
	{"combined", logFmtCommon},
	{"custom1", logFmtCustom1},
	{"vhost_common", logFmtHostCommon},
	{"vhost_combined", logFmtHostCommon},
	{"vhost_custom1", logFmtHostCustom1},
}

func TestLogPattern_Match(t *testing.T) {
	for _, c := range logPatternCases {
		t.Run(c.name, func(t *testing.T) {
			file, err := os.Open("tests/" + c.name + ".log")
			require.NoError(t, err)
			parser := NewLogParser()
			parser.SetInput(file)
			for {
				record, err := parser.Read()
				if err == io.EOF {
					break
				}
				require.NoError(t, err)
				err = c.pattern.Match(record)
				assert.NoError(t, err)
			}
		})
	}
}

func TestLogPattern_guess(t *testing.T) {
	for _, c := range logPatternCases {
		t.Run(c.name, func(t *testing.T) {
			file, err := os.Open("tests/" + c.name + ".log")
			require.NoError(t, err)
			parser := NewLogParser()
			parser.SetInput(file)
			for {
				record, err := parser.Read()
				if err == io.EOF {
					break
				}
				require.NoError(t, err)
				pattern := guessPattern(record)
				assert.Equal(t, c.pattern, pattern)
			}
		})
	}
}
