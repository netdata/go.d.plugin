package cockroachdb

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	testMetricsData, _ = ioutil.ReadFile("testdata/metrics.txt")
)

func Test_readTestData(t *testing.T) {
	assert.NotNil(t, testMetricsData)
}

func TestNew(t *testing.T) {

}

func TestCockroachDB_Init(t *testing.T) {

}

func TestCockroachDB_Check(t *testing.T) {

}

func TestCockroachDB_Charts(t *testing.T) {

}

func TestCockroachDB_Cleanup(t *testing.T) {

}

func TestCockroachDB_Collect(t *testing.T) {
	cockroachDB, srv := prepareClientServer(t)
	defer srv.Close()

	expected := map[string]int64{
		"storage_capacity_available":              40402062147584,
		"storage_capacity_unusable":               23800157791684,
		"storage_capacity_usable":                 40402194045500,
		"storage_capacity_total_used_percentage":  37070,
		"storage_capacity_usable_used_percentage": 0,
		"storage_capacity_reserved":               0,
		"storage_capacity_total":                  64202351837184,
		"storage_capacity_used":                   131897916,
		"storage_file_descriptors_open":           47,
		"storage_file_descriptors_soft_limit":     1048576,
		"storage_live_bytes":                      81979227,
		"storage_rocksdb_block_cache_hit_rate":    92104,
		"storage_rocksdb_block_cache_hits":        94825,
		"storage_rocksdb_block_cache_bytes":       39397184,
		"storage_rocksdb_block_cache_misses":      8129,
		"storage_rocksdb_compactions":             7,
		"storage_rocksdb_flushes":                 13,
		"storage_rocksdb_num_sstables":            8,
		"storage_rocksdb_read_amplification":      1,
		"storage_sys_bytes":                       13327,
		"storage_timeseries_write_bytes":          82810041,
		"storage_timeseries_write_errors":         0,
		"storage_timeseries_write_samples":        845784,
	}

	collected := cockroachDB.Collect()

	assert.Equal(t, expected, collected)
	testCharts(t, cockroachDB, collected)
}

func testCharts(t *testing.T, cockroachDB *CockroachDB, collected map[string]int64) {
	ensureCollectedHasAllChartsDimsVarsIDs(t, cockroachDB, collected)
}

func ensureCollectedHasAllChartsDimsVarsIDs(t *testing.T, c *CockroachDB, collected map[string]int64) {
	for _, chart := range *c.Charts() {
		for _, dim := range chart.Dims {
			_, ok := collected[dim.ID]
			assert.Truef(t, ok, "collected metrics has no data for dim '%s' chart '%s'", dim.ID, chart.ID)
		}
		for _, v := range chart.Vars {
			_, ok := collected[v.ID]
			assert.Truef(t, ok, "collected metrics has no data for var '%s' chart '%s'", v.ID, chart.ID)
		}
	}
}

func prepareClientServer(t *testing.T) (*CockroachDB, *httptest.Server) {
	t.Helper()
	ts := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			_, _ = w.Write(testMetricsData)
		}))

	cockroachDB := New()
	cockroachDB.UserURL = ts.URL
	require.True(t, cockroachDB.Init())

	return cockroachDB, ts
}
