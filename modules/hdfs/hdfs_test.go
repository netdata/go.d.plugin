package hdfs

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/netdata/go-orchestrator/module"
)

var (
	testUnknownNodeData, _ = ioutil.ReadFile("testdata/unknownnode.json")
	testDataNodeData, _    = ioutil.ReadFile("testdata/datanode.json")
	testNameNodeData, _    = ioutil.ReadFile("testdata/namenode.json")
)

func Test_readTestData(t *testing.T) {
	assert.NotNil(t, testUnknownNodeData)
	assert.NotNil(t, testDataNodeData)
	assert.NotNil(t, testNameNodeData)
}

func TestNew(t *testing.T) {
	assert.Implements(t, (*module.Module)(nil), New())
}

func TestHDFS_Init(t *testing.T) {
	job := New()

	assert.True(t, job.Init())
}

func TestHDFS_InitErrorOnCreatingClientWrongTLSCA(t *testing.T) {
	job := New()
	job.ClientTLSConfig.TLSCA = "testdata/tls"

	assert.False(t, job.Init())
}

func TestHDFS_Check(t *testing.T) {
	ts := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				_, _ = w.Write(testNameNodeData)
			}))
	defer ts.Close()

	job := New()
	job.UserURL = ts.URL
	require.True(t, job.Init())

	assert.True(t, job.Check())
	assert.NotZero(t, job.nodeType)
}

func TestHDFS_CheckUnknownNodeType(t *testing.T) {
	ts := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				_, _ = w.Write(testUnknownNodeData)
			}))
	defer ts.Close()

	job := New()
	job.UserURL = ts.URL
	require.True(t, job.Init())

	assert.True(t, job.Check())
	assert.Equal(t, unknownNodeType, job.nodeType)
}

func TestHDFS_CheckDataNodeType(t *testing.T) {
	ts := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				_, _ = w.Write(testDataNodeData)
			}))
	defer ts.Close()

	job := New()
	job.UserURL = ts.URL
	require.True(t, job.Init())

	assert.True(t, job.Check())
	assert.Equal(t, dataNodeType, job.nodeType)
}

func TestHDFS_CheckNameNodeType(t *testing.T) {
	ts := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				_, _ = w.Write(testNameNodeData)
			}))
	defer ts.Close()

	job := New()
	job.UserURL = ts.URL
	require.True(t, job.Init())

	assert.True(t, job.Check())
	assert.Equal(t, nameNodeType, job.nodeType)
}

func TestHDFS_CheckErrorOnNodeTypeDetermination(t *testing.T) {
	ts := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				_, _ = w.Write([]byte("{}"))
			}))
	defer ts.Close()

	job := New()
	job.UserURL = ts.URL
	require.True(t, job.Init())

	assert.False(t, job.Check())
}

func TestHDFS_CheckNoResponse(t *testing.T) {
	job := New()
	job.UserURL = "http://127.0.0.1:38001/jmx"
	require.True(t, job.Init())

	assert.False(t, job.Check())
}

func TestHDFS_Charts(t *testing.T) {
	assert.NotNil(t, New().Charts())
}

func TestHDFS_ChartsUnknownNode(t *testing.T) {
	job := New()
	job.nodeType = unknownNodeType

	assert.Equal(t, unknownNodeCharts(), job.Charts())
}

func TestHDFS_ChartsDataNode(t *testing.T) {
	job := New()
	job.nodeType = dataNodeType

	assert.Equal(t, dataNodeCharts(), job.Charts())
}

func TestHDFS_ChartsNameNode(t *testing.T) {
	job := New()
	job.nodeType = nameNodeType

	assert.Equal(t, nameNodeCharts(), job.Charts())
}

func TestHDFS_Cleanup(t *testing.T) {
	New().Cleanup()
}

func TestHDFS_CollectDataNode(t *testing.T) {
	ts := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				_, _ = w.Write(testDataNodeData)
			}))
	defer ts.Close()

	job := New()
	job.UserURL = ts.URL
	require.True(t, job.Init())
	require.True(t, job.Check())

	expected := map[string]int64{
		"jvm_gc_count":                       155,
		"jvm_gc_num_info_threshold_exceeded": 0,
		"jvm_gc_num_warn_threshold_exceeded": 0,
		"jvm_gc_time_millis":                 672,
		"jvm_gc_total_extra_sleep_time":      8783,
		"jvm_log_error":                      1,
		"jvm_log_fatal":                      0,
		"jvm_log_info":                       257,
		"jvm_log_warn":                       2,
		"jvm_mem_heap_committed":             60,
		"jvm_mem_heap_max":                   843,
		"jvm_mem_heap_used":                  18,
		"jvm_mem_max":                        843,
		"jvm_mem_non_heap_committed":         54,
		"jvm_mem_non_heap_max":               -1,
		"jvm_mem_non_heap_used":              53,
		"jvm_threads_blocked":                0,
		"jvm_threads_new":                    0,
		"jvm_threads_runnable":               11,
		"jvm_threads_terminated":             0,
		"jvm_threads_timed_waiting":          25,
		"jvm_threads_waiting":                11,
	}

	assert.Equal(t, expected, job.Collect())
}

func TestHDFS_CollectNameNode(t *testing.T) {
	ts := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				_, _ = w.Write(testNameNodeData)
			}))
	defer ts.Close()

	job := New()
	job.UserURL = ts.URL
	require.True(t, job.Init())
	require.True(t, job.Check())

	expected := map[string]int64{
		"fsn_block_capacity":                                    2097152,
		"fsn_blocks_total":                                      15,
		"fsn_bytes_in_future_ec_block_groups":                   0,
		"fsn_bytes_in_future_replicated_blocks":                 0,
		"fsn_capacity_remaining":                                65861697536,
		"fsn_capacity_remaining_gb":                             61,
		"fsn_capacity_total":                                    107351072768,
		"fsn_capacity_total_gb":                                 100,
		"fsn_capacity_used":                                     2372116480,
		"fsn_capacity_used_gb":                                  2,
		"fsn_capacity_used_non_dfs":                             39117258752,
		"fsn_corrupt_blocks":                                    0,
		"fsn_corrupt_ec_block_groups":                           0,
		"fsn_corrupt_replicated_blocks":                         0,
		"fsn_estimated_capacity_lost_total":                     0,
		"fsn_excess_blocks":                                     0,
		"fsn_expired_heartbeats":                                0,
		"fsn_files_total":                                       12,
		"fsn_highest_priority_low_redundancy_ec_blocks":         0,
		"fsn_highest_priority_low_redundancy_replicated_blocks": 0,
		"fsn_last_checkpoint_time":                              1566814983890,
		"fsn_last_written_transaction_id":                       624,
		"fsn_lock_queue_length":                                 0,
		"fsn_low_redundancy_blocks":                             0,
		"fsn_low_redundancy_ec_block_groups":                    0,
		"fsn_low_redundancy_replicated_blocks":                  0,
		"fsn_millis_since_last_loaded_edits":                    0,
		"fsn_missing_blocks":                                    0,
		"fsn_missing_ec_block_groups":                           0,
		"fsn_missing_repl_one_blocks":                           0,
		"fsn_missing_replicated_blocks":                         0,
		"fsn_missing_replication_one_blocks":                    0,
		"fsn_num_active_clients":                                0,
		"fsn_num_dead_data_nodes":                               0,
		"fsn_num_decom_dead_data_nodes":                         0,
		"fsn_num_decom_live_data_nodes":                         0,
		"fsn_num_decommissioning_data_nodes":                    0,
		"fsn_num_encryption_zones":                              0,
		"fsn_num_entering_maintenance_data_nodes":               0,
		"fsn_num_files_under_construction":                      0,
		"fsn_num_in_maintenance_dead_data_nodes":                0,
		"fsn_num_in_maintenance_live_data_nodes":                0,
		"fsn_num_live_data_nodes":                               2,
		"fsn_num_stale_storages":                                0,
		"fsn_num_timed_out_pending_reconstructions":             0,
		"fsn_pending_data_node_message_count":                   0,
		"fsn_pending_deletion_blocks":                           0,
		"fsn_pending_deletion_ec_blocks":                        0,
		"fsn_pending_deletion_replicated_blocks":                0,
		"fsn_pending_reconstruction_blocks":                     0,
		"fsn_pending_replication_blocks":                        0,
		"fsn_postponed_misreplicated_blocks":                    0,
		"fsn_provided_capacity_total":                           0,
		"fsn_scheduled_replication_blocks":                      0,
		"fsn_snapshots":                                         0,
		"fsn_snapshottable_directories":                         0,
		"fsn_stale_data_nodes":                                  0,
		"fsn_total_ec_block_groups":                             0,
		"fsn_total_load":                                        2,
		"fsn_total_replicated_blocks":                           15,
		"fsn_total_sync_count":                                  2,
		"fsn_transactions_since_last_checkpoint":                1,
		"fsn_transactions_since_last_log_roll":                  1,
		"fsn_under_replicated_blocks":                           0,
		"fsn_volume_failures_total":                             0,
		"jvm_gc_count":                                          1699,
		"jvm_gc_num_info_threshold_exceeded":                    0,
		"jvm_gc_num_warn_threshold_exceeded":                    0,
		"jvm_gc_time_millis":                                    3483,
		"jvm_gc_total_extra_sleep_time":                         1944,
		"jvm_log_error":                                         0,
		"jvm_log_fatal":                                         0,
		"jvm_log_info":                                          3382077,
		"jvm_log_warn":                                          3378983,
		"jvm_mem_heap_committed":                                67,
		"jvm_mem_heap_max":                                      843,
		"jvm_mem_heap_used":                                     26,
		"jvm_mem_max":                                           843,
		"jvm_mem_non_heap_committed":                            67,
		"jvm_mem_non_heap_max":                                  -1,
		"jvm_mem_non_heap_used":                                 66,
		"jvm_threads_blocked":                                   0,
		"jvm_threads_new":                                       0,
		"jvm_threads_runnable":                                  7,
		"jvm_threads_terminated":                                0,
		"jvm_threads_timed_waiting":                             34,
		"jvm_threads_waiting":                                   6,
	}

	assert.Equal(t, expected, job.Collect())
}

func TestHDFS_CollectNoResponse(t *testing.T) {
	job := New()
	job.UserURL = "http://127.0.0.1:38001/jmx"
	require.True(t, job.Init())

	assert.Nil(t, job.Collect())
}

func TestHDFS_CollectReceiveInvalidResponse(t *testing.T) {
	ts := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				_, _ = w.Write([]byte("hello and\ngoodbye!\n"))
			}))
	defer ts.Close()

	job := New()
	job.UserURL = ts.URL
	require.True(t, job.Init())

	assert.Nil(t, job.Collect())
}

func TestHDFS_CollectReceive404(t *testing.T) {
	ts := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusNotFound)
			}))
	defer ts.Close()

	job := New()
	job.UserURL = ts.URL
	require.True(t, job.Init())

	assert.Nil(t, job.Collect())
}
