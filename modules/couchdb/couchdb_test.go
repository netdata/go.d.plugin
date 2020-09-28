package couchdb

import (
	"io/ioutil"
	"testing"

	"github.com/netdata/go.d.plugin/agent/module"
	"github.com/netdata/go.d.plugin/pkg/tlscfg"
	"github.com/netdata/go.d.plugin/pkg/web"
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

func TestElasticsearch_Init(t *testing.T) {
	tests := map[string]struct {
		config              Config
		expectedNumOfCharts int
		shouldFail          bool
	}{
		"default": {
			expectedNumOfCharts: numOfCharts(
				dbActivityCharts,
				httpTrafficBreakdownCharts,
				serverOperationsCharts,
				erlangStatisticsCharts,
			),
			config: New().Config,
		},
		"URL not set": {
			shouldFail: true,
			config: Config{
				HTTP: web.HTTP{
					Request: web.Request{URL: ""},
				}},
		},
		"invalid TLSCA": {
			shouldFail: true,
			config: Config{
				HTTP: web.HTTP{
					Client: web.Client{
						TLSConfig: tlscfg.TLSConfig{TLSCA: "testdata/tls"},
					},
				}},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			es := New()
			es.Config = test.config

			if test.shouldFail {
				assert.False(t, es.Init())
			} else {
				assert.True(t, es.Init())
				assert.Equal(t, test.expectedNumOfCharts, len(*es.Charts()))
			}
		})
	}
}

func numOfCharts(charts ...Charts) (num int) {
	for _, v := range charts {
		num += len(v)
	}
	return num
}
