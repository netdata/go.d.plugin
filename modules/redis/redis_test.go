package redis

import (
	"io/ioutil"
	"testing"

	"github.com/netdata/go-orchestrator/module"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	infoData, _ = ioutil.ReadFile("testdata/info.txt")
)

func TestNew(t *testing.T) {
	job := New()

	assert.Implements(t, (*module.Module)(nil), job)
	assert.Equal(t, defaultHost, job.Host)
	assert.Equal(t, defaultPort, job.Port)
}

func TestRedis_Charts(t *testing.T) {
	assert.NotNil(t, New().Charts())
}

func TestRedis_Collect(t *testing.T) {
	var bgSaveTime int64 = 2
	metrics := make(map[string]int64)
	err := parseMetrics(string(infoData), charts.Copy(), metrics, &bgSaveTime, 3)

	require.Nil(t, err)

	expectedMetrics := map[string]int64{
		"total_commands_processed":    9590,
		"instantaneous_ops_per_sec":   3,
		"hit_rate":                    98,
		"used_memory":                 840400,
		"used_memory_lua":             37888,
		"total_net_input_bytes":       134274,
		"total_net_output_bytes":      26591079,
		"db0":                         1,
		"db1":                         3,
		"db2":                         2,
		"evicted_keys":                123,
		"total_connections_received":  8,
		"rejected_connections":        7,
		"connected_clients":           6,
		"blocked_clients":             8,
		"connected_slaves":            3,
		"rdb_changes_since_last_save": 7,
		"rdb_bgsave_in_progress":      5,
		"rdb_last_bgsave_status":      1,
		"uptime_in_seconds":           10277,
	}

	assert.Equal(t, expectedMetrics, metrics)
	assert.EqualValues(t, 5, bgSaveTime)
}
