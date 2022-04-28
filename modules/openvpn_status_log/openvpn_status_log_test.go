package openvpn_status_log

import (
	// "fmt"
	// "io/ioutil"
	// "os"
	"testing"

	"github.com/netdata/go.d.plugin/agent/module"
	"github.com/netdata/go.d.plugin/pkg/matcher"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var expectedWithClients = map[string]int64{
	"bytes_in":              22017,
	"bytes_out":             265176,
	"clients":               2,
	"gofle_bytes_in":        19265,
	"gofle_bytes_out":       261631,
	"client_bsd2_bytes_in":  2752,
	"client_bsd2_bytes_out": 3545,
}

var expectedWithoutClients = map[string]int64{
	"bytes_in":  0,
	"bytes_out": 0,
	"clients":   0,
}

func TestNew(t *testing.T) {
	job := New()

	assert.Implements(t, (*module.Module)(nil), job)
	assert.Equal(t, defaultFilePath, job.StatusPath)
	assert.NotNil(t, job.charts)
	assert.NotNil(t, job.collectedUsers)
}

func TestOpenVPN_Status_Init(t *testing.T) {
	assert.True(t, New().Init())
}

func TestOpenVPN_Status_WithClients_Collect(t *testing.T) {
	logFiles := []string{
		"testdata/status_version_1.txt",
		"testdata/status_version_2.txt",
		"testdata/status_version_3.txt",
	}
	for _, log := range logFiles {
		job := New()
		job.StatusPath = log

		require.True(t, job.Init())
		job.perUserMatcher = matcher.TRUE()
		require.True(t, job.Check())

		mx := job.Collect()
		require.NotNil(t, mx)
		delete(mx, "gofle_connection_time")
		delete(mx, "client_bsd2_connection_time")
		assert.Equal(t, expectedWithClients, mx)
	}
}

func TestOpenVPN_Status_WithoutClients_Collect(t *testing.T) {
	logFiles := []string{
		"testdata/status_version_1_wo_clients.txt",
		"testdata/status_version_2_wo_clients.txt",
		"testdata/status_version_3_wo_clients.txt",
	}
	for _, log := range logFiles {
		job := New()
		job.StatusPath = log

		require.True(t, job.Init())
		job.perUserMatcher = matcher.TRUE()
		require.True(t, job.Check())

		mx := job.Collect()
		require.NotNil(t, mx)
		assert.Equal(t, expectedWithoutClients, mx)
	}
}
