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
	assert.Implements(t, (*module.Module)(nil), job)
	assert.Equal(t, defaultHost, job.Host)
	assert.Equal(t, defaultPort, job.Port)
}

func TestRedis_Charts(t *testing.T) {
	assert.NotNil(t, New().Charts())
}

func TestRedis_Collect(t *testing.T) {
	metrics := make(map[string]int64)
	err := parseMetrics(string(infoData), metrics)

	require.Nil(t, err)

	expectedMetrics := map[string]int64{
		"total_commands_processed":  9590,
		"instantaneous_ops_per_sec": 3,
		"hit_rate":                  98,
		"used_memory":               840400,
		"used_memory_lua":           37888,
		"total_net_input_bytes":     134274,
		"total_net_output_bytes":    26591079,
	}

	assert.Equal(t, expectedMetrics, metrics)
}
