package couchdb

import (
	"io/ioutil"
	"testing"

	"github.com/netdata/go.d.plugin/agent/module"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	responseRoot, _        = ioutil.ReadFile("testdata/root.json")
	responseNodeStats, _   = ioutil.ReadFile("testdata/node_stats.json")
	responseActiveTasks, _ = ioutil.ReadFile("testdata/active_tasks.json")
	responseNodeSystem, _  = ioutil.ReadFile("testdata/node_system.json")
)

func Test_testDataIsCorrectlyReadAndValid(t *testing.T) {
	for name, data := range map[string][]byte{
		"responseRoot":        responseRoot,
		"responseNodeStats":   responseNodeStats,
		"responseActiveTasks": responseActiveTasks,
		"responseNodeSystem":  responseNodeSystem,
	} {
		require.NotNilf(t, data, name)
	}
}

func TestNew(t *testing.T) {
	assert.Implements(t, (*module.Module)(nil), New())
}
