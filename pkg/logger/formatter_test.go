package logger

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFormatter_Output_cli(t *testing.T) {
	out := &bytes.Buffer{}
	fmtter := newFormatter(out, true, "test")

	fmtter.Output(INFO, "mod1", "job1", 1, "hello")
	assert.NotRegexp(t, `\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2}: `, out.String())
	assert.Contains(t, out.String(), "INFO")
	assert.Contains(t, out.String(), "mod1[job1]")
	assert.Contains(t, out.String(), "formatter_test.go:")
	assert.Contains(t, out.String(), "hello")
}

func TestFormatter_Output_file(t *testing.T) {
	out := &bytes.Buffer{}
	fmtter := newFormatter(out, false, "test")

	fmtter.Output(INFO, "mod1", "job1", 1, "hello")
	assert.Regexp(t, `\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2}: `, out.String())
	assert.Contains(t, out.String(), "INFO")
	assert.Contains(t, out.String(), "mod1[job1]")
	assert.NotContains(t, out.String(), "formatter_test.go:")
	assert.Contains(t, out.String(), "hello")
}
