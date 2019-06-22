package portcheck

import (
	"errors"
	"net"
	"testing"
	"time"

	"github.com/netdata/go-orchestrator/module"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	job := New()

	assert.Implements(t, (*module.Module)(nil), job)
	assert.Equal(t, defaultConnectTimeout, job.Timeout.Duration)
}

func TestPortCheck_Init(t *testing.T) {
	job := New()

	job.Host = "127.0.0.1"
	job.Ports = []int{38001, 38002}
	assert.True(t, job.Init())
	assert.Len(t, job.ports, 2)
}
func TestPortCheck_InitNG(t *testing.T) {
	job := New()

	assert.False(t, job.Init())
	job.Host = "127.0.0.1"
	assert.False(t, job.Init())
	job.Ports = []int{38001, 38002}
	assert.True(t, job.Init())
}

func TestPortCheck_Check(t *testing.T) { assert.True(t, New().Check()) }

func TestPortCheck_Cleanup(t *testing.T) { New().Cleanup() }

func TestPortCheck_Charts(t *testing.T) {
	job := New()
	job.Ports = []int{1, 2}
	assert.Len(t, *job.Charts(), len(portCharts)*len(job.Ports))
}

func TestPortCheck_Collect(t *testing.T) {
	job := New()

	job.Host = "127.0.0.1"
	job.Ports = []int{38001, 38002}
	job.UpdateEvery = 5
	job.dial = testDial(nil)
	require.True(t, job.Init())
	require.True(t, job.Check())

	expected := map[string]int64{
		"port_38001_current_state_duration": int64(job.UpdateEvery),
		"port_38001_failed":                 0,
		"port_38001_latency":                0,
		"port_38001_success":                1,
		"port_38001_timeout":                0,
		"port_38002_current_state_duration": int64(job.UpdateEvery),
		"port_38002_failed":                 0,
		"port_38002_latency":                0,
		"port_38002_success":                1,
		"port_38002_timeout":                0,
	}
	assert.Equal(t, expected, job.Collect())

	expected = map[string]int64{
		"port_38001_current_state_duration": int64(job.UpdateEvery) * 2,
		"port_38001_failed":                 0,
		"port_38001_latency":                0,
		"port_38001_success":                1,
		"port_38001_timeout":                0,
		"port_38002_current_state_duration": int64(job.UpdateEvery) * 2,
		"port_38002_failed":                 0,
		"port_38002_latency":                0,
		"port_38002_success":                1,
		"port_38002_timeout":                0,
	}
	assert.Equal(t, expected, job.Collect())

	job.dial = testDial(errors.New("failed"))

	expected = map[string]int64{
		"port_38001_current_state_duration": int64(job.UpdateEvery),
		"port_38001_failed":                 1,
		"port_38001_latency":                0,
		"port_38001_success":                0,
		"port_38001_timeout":                0,
		"port_38002_current_state_duration": int64(job.UpdateEvery),
		"port_38002_failed":                 1,
		"port_38002_latency":                0,
		"port_38002_success":                0,
		"port_38002_timeout":                0,
	}
	assert.Equal(t, expected, job.Collect())

	job.dial = testDial(timeoutError{})

	expected = map[string]int64{
		"port_38001_current_state_duration": int64(job.UpdateEvery),
		"port_38001_failed":                 0,
		"port_38001_latency":                0,
		"port_38001_success":                0,
		"port_38001_timeout":                1,
		"port_38002_current_state_duration": int64(job.UpdateEvery),
		"port_38002_failed":                 0,
		"port_38002_latency":                0,
		"port_38002_success":                0,
		"port_38002_timeout":                1,
	}
	assert.Equal(t, expected, job.Collect())
}

func testDial(err error) dialFunc {
	return func(_, _ string, _ time.Duration) (net.Conn, error) { return &net.TCPConn{}, err }
}

type timeoutError struct{}

func (timeoutError) Error() string { return "timeout" }

func (timeoutError) Timeout() bool { return true }
