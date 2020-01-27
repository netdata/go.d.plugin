package cockroachdb

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"sort"
	"testing"
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

	m := cockroachDB.Collect()
	l := make([]string, 0)
	for k := range m {
		l = append(l, k)
	}
	sort.Strings(l)
	for _, value := range l {
		fmt.Println(fmt.Sprintf("\"%s\": %d,", value, m[value]))
	}

	expected := map[string]int64{
		"capacity":                           64202351837184,
		"capacity_available":                 40402062147584,
		"capacity_reserved":                  0,
		"capacity_unusable":                  23800157791684,
		"capacity_usable":                    40402194045500,
		"capacity_usable_used_percentage":    0,
		"capacity_used":                      131897916,
		"capacity_used_percentage":           37070,
		"clock_offset_meannanos":             -14326,
		"keybytes":                           6730852,
		"liveness_livenodes":                 3,
		"range_adds":                         0,
		"range_merges":                       0,
		"range_removes":                      0,
		"range_snapshots_generated":          0,
		"range_snapshots_learner_applied":    0,
		"range_snapshots_normal_applied":     0,
		"range_snapshots_preemptive_applied": 0,
		"range_splits":                       0,
		"ranges":                             34,
		"ranges_overreplicated":              0,
		"ranges_unavailable":                 0,
		"ranges_underreplicated":             0,
		"rebalancing_queriespersecond":       801,
		"rebalancing_writespersecond":        213023,
		"replicas":                           34,
		"replicas_leaders":                   7,
		"replicas_leaders_not_leaseholders":  0,
		"replicas_leaseholders":              7,
		"replicas_quiescent":                 34,
		"replicas_reserved":                  0,
		"requests_slow_latch":                0,
		"requests_slow_lease":                0,
		"requests_slow_raft":                 0,
		"rocksdb_block_cache_hit_rate":       92104,
		"rocksdb_block_cache_hits":           94825,
		"rocksdb_block_cache_misses":         8129,
		"rocksdb_block_cache_usage":          39397184,
		"rocksdb_compactions":                7,
		"rocksdb_flushes":                    13,
		"rocksdb_num_sstables":               8,
		"rocksdb_read_amplification":         1,
		"sql_bytesin":                        0,
		"sql_bytesout":                       0,
		"sql_conns":                          0,
		"sql_ddl_count":                      0,
		"sql_delete_count":                   0,
		"sql_distsql_flows_active":           0,
		"sql_distsql_flows_total":            1042,
		"sql_distsql_queries_active":         0,
		"sql_distsql_queries_total":          2660,
		"sql_failure_count":                  0,
		"sql_insert_count":                   0,
		"sql_select_count":                   0,
		"sql_txn_abort_count":                0,
		"sql_txn_begin_count":                0,
		"sql_txn_commit_count":               0,
		"sql_txn_rollback_count":             0,
		"sql_update_count":                   0,
		"sys_cgo_allocbytes":                 63363512,
		"sys_cgo_totalbytes":                 81698816,
		"sys_cgocalls":                       577778,
		"sys_cpu_sys_ns":                     154420000000,
		"sys_cpu_user_ns":                    227620000000,
		"sys_fd_open":                        47,
		"sys_fd_softlimit":                   1048576,
		"sys_gc_count":                       279,
		"sys_gc_pause_ns":                    60700450,
		"sys_go_allocbytes":                  106576224,
		"sys_go_totalbytes":                  197562616,
		"sys_goroutines":                     235,
		"sys_host_disk_iopsinprogress":       0,
		"sys_host_disk_read_bytes":           43319296,
		"sys_host_disk_read_count":           1176,
		"sys_host_disk_write_bytes":          942080,
		"sys_host_disk_write_count":          106,
		"sys_host_net_recv_bytes":            234392325,
		"sys_host_net_recv_packets":          593876,
		"sys_host_net_send_bytes":            461746036,
		"sys_host_net_send_packets":          644128,
		"sys_rss":                            314691584,
		"sys_uptime":                         12224,
		"timeseries_write_bytes":             82810041,
		"timeseries_write_errors":            0,
		"timeseries_write_samples":           845784,
		"valbytes":                           75527718,
	}

	collected := cockroachDB.Collect()

	assert.Equal(t, expected, collected)
	//testCharts(t, cockroachDB, collected)
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
