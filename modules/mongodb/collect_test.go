package mongo

import (
	"fmt"
	"testing"

	"github.com/netdata/go.d.plugin/agent/module"
	"github.com/netdata/go.d.plugin/modules/mongodb/testdata/v5.0.0"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
)

func TestMongo_serverStatusCollect(t *testing.T) {
	var status map[string]interface{}
	err := bson.UnmarshalExtJSON([]byte(v5_0_0.ServerStatus), true, &status)
	require.NoError(t, err)

	ms := map[string]int64{}
	m := New()
	m.charts, err = m.initCharts()
	assert.NoError(t, err)
	m.addOptionalCharts(status)
	iterateServerStatus(ms, status)
	for _, chart := range serverStatusCharts {
		for _, dim := range chart.Dims {
			_, ok := ms[toID(dim.ID)]
			if contains(v5_0_0.V500ServerStatusMissingValues, toID(dim.ID)) {
				assert.False(t, ok, fmt.Sprintf("value for dim.ID:%s not found", toID(dim.ID)))
			} else {
				assert.True(t, ok, fmt.Sprintf("value for dim.ID:%s not found", toID(dim.ID)))
			}
		}
	}
}

func TestMongo_serverStatusCollectOptionalCharts(t *testing.T) {
	var status map[string]interface{}
	err := bson.UnmarshalExtJSON([]byte(v5_0_0.ServerStatus), true, &status)
	require.NoError(t, err)

	m := New()
	m.charts = &module.Charts{}
	m.addOptionalCharts(status)
	assert.True(t, m.charts.Has(chartTransactionsCurrent.ID))
	assert.True(t, m.charts.Has(chartGlobalLockActiveClients.ID))
	assert.True(t, m.charts.Has(chartCollections.ID))
	assert.True(t, m.charts.Has(chartTcmallocGeneric.ID))
	assert.True(t, m.charts.Has(chartTcmalloc.ID))
	assert.True(t, m.charts.Has(chartGlobalLockCurrentQueue.ID))
	assert.True(t, m.charts.Has(chartMetricsCommands.ID))
	assert.True(t, m.charts.Has(chartGlobalLocks.ID))
	assert.True(t, m.charts.Has(chartFlowControl.ID))
	assert.True(t, m.charts.Has(chartWiredTigerBlockManager.ID))
	assert.True(t, m.charts.Has(chartWiredTigerCache.ID))
	assert.True(t, m.charts.Has(chartWiredTigerCapacity.ID))
	assert.True(t, m.charts.Has(chartWiredTigerConnection.ID))
	assert.True(t, m.charts.Has(chartWiredTigerCursor.ID))
	assert.True(t, m.charts.Has(chartWiredTigerLock.ID))
	assert.True(t, m.charts.Has(chartWiredTigerLockDuration.ID))
	assert.True(t, m.charts.Has(chartWiredTigerLogOps.ID))
	assert.True(t, m.charts.Has(chartWiredTigerLogBytes.ID))
	assert.True(t, m.charts.Has(chartWiredTigerTransactions.ID))
	assert.False(t, m.charts.Has("some random id"))

}

func TestMongo_metricNotExists(t *testing.T) {
	var status map[string]interface{}
	err := bson.UnmarshalExtJSON([]byte(v5_0_0.V500ServerStatusInvalidData), true, &status)
	require.NoError(t, err)

	m := New()
	m.metricExists(status, "invalid key", &chartAsserts)
	assert.False(t, m.charts.Has(chartAsserts.ID))
}

func TestMongo_wrongBsonData(t *testing.T) {
	var status map[string]interface{}
	err := bson.UnmarshalExtJSON([]byte(v5_0_0.V500ServerStatusInvalidBsonData), true, &status)
	require.NoError(t, err)

	m := New()
	m.metricExists(status, "asserts", &chartAsserts)
	assert.False(t, m.charts.Has(chartAsserts.ID))
}

func TestMongo_addChartTwice(t *testing.T) {
	var status map[string]interface{}
	err := bson.UnmarshalExtJSON([]byte(v5_0_0.ServerStatus), true, &status)
	require.NoError(t, err)

	m := New()
	m.metricExists(status, "asserts", &chartAsserts)
	assert.True(t, m.charts.Has(chartAsserts.ID))
	m.metricExists(status, "asserts", &chartAsserts)
}

func TestMongo_addIfExistsWrongData(t *testing.T) {
	var status map[string]interface{}
	err := bson.UnmarshalExtJSON([]byte(v5_0_0.V500ServerStatusInvalidBsonData), true, &status)
	require.NoError(t, err)

	ms := map[string]int64{}
	addIfExists(status, "asserts.warning", ms)
	assert.Len(t, ms, 0)
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
