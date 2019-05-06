package scaleio

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"testing"

	"github.com/netdata/go-orchestrator/module"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	testSelectedStatsData, _ = ioutil.ReadFile("testdata/selected_statistics.json")
)

func TestNew(t *testing.T) {
	job := New()

	assert.Implements(t, (*module.Module)(nil), job)
	//assert.Equal(t, defaultURL, job.UserURL)
	//assert.Equal(t, defaultHTTPTimeout, job.Timeout.Duration)
}

func TestScaleIO_Init(t *testing.T) {
	job := New()
	job.UserURL = "http://127.0.0.1:38001"

	require.True(t, job.Init())
	assert.NotNil(t, job.apiClient)
}

func TestScaleIO_InitNG(t *testing.T) {
	job := New()

	require.True(t, job.Init())
	assert.NotNil(t, job.apiClient)
}

func TestScaleIO_Check(t *testing.T) {
	job := New()
	job.UserURL = "http://127.0.0.1:38001"

	require.True(t, job.Init())
	job.apiClient = &okAPIClient{}
	require.True(t, job.Check())
}

func TestScaleIO_CheckNG(t *testing.T) {
	job := New()
	job.UserURL = "http://127.0.0.1:38001"

	require.True(t, job.Init())
	assert.False(t, job.Check())
}

func TestScaleIO_Charts(t *testing.T) { assert.NotNil(t, New().Charts()) }

func TestScaleIO_Cleanup(t *testing.T) {
	job := New()
	job.UserURL = "http://127.0.0.1:38001"

	require.True(t, job.Init())
	job.apiClient = &okAPIClient{}
	require.True(t, job.Check())

	assert.True(t, job.apiClient.IsLoggedIn())
	job.Cleanup()
	assert.False(t, job.apiClient.IsLoggedIn())
}

func TestScaleIO_Collect(t *testing.T) {
	job := New()
	job.UserURL = "http://127.0.0.1:38001"

	require.True(t, job.Init())
	job.apiClient = &okAPIClient{}
	require.True(t, job.Check())

	expected := map[string]int64{
		"system_backend_primary_bandwidth_read":           0,
		"system_backend_primary_bandwidth_read_write":     82616000,
		"system_backend_primary_bandwidth_write":          82616000,
		"system_backend_primary_io_size_read":             0,
		"system_backend_primary_io_size_read_write":       320963,
		"system_backend_primary_io_size_write":            320963,
		"system_backend_primary_iops_read":                0,
		"system_backend_primary_iops_read_write":          257399,
		"system_backend_primary_iops_write":               257399,
		"system_backend_secondary_bandwidth_read":         0,
		"system_backend_secondary_bandwidth_read_write":   82432000,
		"system_backend_secondary_bandwidth_write":        82432000,
		"system_backend_secondary_io_size_read":           0,
		"system_backend_secondary_io_size_read_write":     396689,
		"system_backend_secondary_io_size_write":          396689,
		"system_backend_secondary_iops_read":              0,
		"system_backend_secondary_iops_read_write":        207800,
		"system_backend_secondary_iops_write":             207800,
		"system_backend_total_bandwidth_read":             0,
		"system_backend_total_bandwidth_read_write":       165048000,
		"system_backend_total_bandwidth_write":            165048000,
		"system_backend_total_io_size_read":               0,
		"system_backend_total_io_size_read_write":         717652,
		"system_backend_total_io_size_write":              717652,
		"system_backend_total_iops_read":                  0,
		"system_backend_total_iops_read_write":            465200,
		"system_backend_total_iops_write":                 465200,
		"system_capacity_available_for_volume_allocation": 67108864,
		"system_capacity_decreased":                       0,
		"system_capacity_degraded":                        0,
		"system_capacity_failed":                          0,
		"system_capacity_in_maintenance":                  0,
		"system_capacity_in_use":                          65386496,
		"system_capacity_limit":                           337638400,
		"system_capacity_max_capacity":                    337638400,
		"system_capacity_protected":                       65386496,
		"system_capacity_semi_protected":                  0,
		"system_capacity_snap_in_use":                     17451008,
		"system_capacity_snap_in_use_occupied":            0,
		"system_capacity_spare":                           118172672,
		"system_capacity_thick_in_use":                    16777216,
		"system_capacity_thin_allocated":                  67108864,
		"system_capacity_thin_free":                       18499584,
		"system_capacity_thin_in_use":                     48609280,
		"system_capacity_unreachable_unused":              0,
		"system_capacity_unused":                          154079232,
		"system_frontend_user_data_bandwidth_read":        0,
		"system_frontend_user_data_bandwidth_read_write":  87404000,
		"system_frontend_user_data_bandwidth_write":       87404000,
		"system_frontend_user_data_io_size_read":          0,
		"system_frontend_user_data_io_size_read_write":    346841,
		"system_frontend_user_data_io_size_write":         346841,
		"system_frontend_user_data_iops_read":             0,
		"system_frontend_user_data_iops_read_write":       252000,
		"system_frontend_user_data_iops_write":            252000,
		"system_num_of_devices":                           3,
		"system_num_of_fault_sets":                        2,
		"system_num_of_mapped_to_all_volumes":             0,
		"system_num_of_mapped_volumes":                    2,
		"system_num_of_protection_domains":                2,
		"system_num_of_rfcache_devices":                   0,
		"system_num_of_scsi_initiators":                   0,
		"system_num_of_sdc":                               3,
		"system_num_of_sds":                               3,
		"system_num_of_snapshots":                         0,
		"system_num_of_storage_pools":                     2,
		"system_num_of_thick_base_volumes":                1,
		"system_num_of_thin_base_volumes":                 3,
		"system_num_of_unmapped_volumes":                  2,
		"system_num_of_volumes":                           4,
		"system_num_of_volumes_in_deletion":               0,
		"system_num_of_vtrees":                            4,
		"system_rebalance_bandwidth_read":                 0,
		"system_rebalance_bandwidth_read_write":           0,
		"system_rebalance_bandwidth_write":                0,
		"system_rebalance_io_size_read":                   0,
		"system_rebalance_io_size_read_write":             0,
		"system_rebalance_io_size_write":                  0,
		"system_rebalance_iops_read":                      0,
		"system_rebalance_iops_read_write":                0,
		"system_rebalance_iops_write":                     0,
		"system_rebalance_pending_capacity_in_Kb":         0,
		"system_rebuild_backward_bandwidth_read":          0,
		"system_rebuild_backward_bandwidth_read_write":    0,
		"system_rebuild_backward_bandwidth_write":         0,
		"system_rebuild_backward_io_size_read":            0,
		"system_rebuild_backward_io_size_read_write":      0,
		"system_rebuild_backward_io_size_write":           0,
		"system_rebuild_backward_iops_read":               0,
		"system_rebuild_backward_iops_read_write":         0,
		"system_rebuild_backward_iops_write":              0,
		"system_rebuild_backward_pending_capacity_in_Kb":  0,
		"system_rebuild_forward_bandwidth_read":           0,
		"system_rebuild_forward_bandwidth_read_write":     0,
		"system_rebuild_forward_bandwidth_write":          0,
		"system_rebuild_forward_io_size_read":             0,
		"system_rebuild_forward_io_size_read_write":       0,
		"system_rebuild_forward_io_size_write":            0,
		"system_rebuild_forward_iops_read":                0,
		"system_rebuild_forward_iops_read_write":          0,
		"system_rebuild_forward_iops_write":               0,
		"system_rebuild_forward_pending_capacity_in_Kb":   0,
		"system_rebuild_normal_bandwidth_read":            0,
		"system_rebuild_normal_bandwidth_read_write":      0,
		"system_rebuild_normal_bandwidth_write":           0,
		"system_rebuild_normal_io_size_read":              0,
		"system_rebuild_normal_io_size_read_write":        0,
		"system_rebuild_normal_io_size_write":             0,
		"system_rebuild_normal_iops_read":                 0,
		"system_rebuild_normal_iops_read_write":           0,
		"system_rebuild_normal_iops_write":                0,
		"system_rebuild_normal_pending_capacity_in_Kb":    0,
		"system_rebuild_total_bandwidth_read":             0,
		"system_rebuild_total_bandwidth_read_write":       0,
		"system_rebuild_total_bandwidth_write":            0,
		"system_rebuild_total_io_size_read":               0,
		"system_rebuild_total_io_size_read_write":         0,
		"system_rebuild_total_io_size_write":              0,
		"system_rebuild_total_iops_read":                  0,
		"system_rebuild_total_iops_read_write":            0,
		"system_rebuild_total_iops_write":                 0,
		"system_rebuild_total_pending_capacity_in_Kb":     0,
		"system_total_bandwidth_read":                     0,
		"system_total_bandwidth_read_write":               165048000,
		"system_total_bandwidth_write":                    165048000,
		"system_total_io_size_read":                       0,
		"system_total_io_size_read_write":                 354789,
		"system_total_io_size_write":                      354789,
		"system_total_iops_read":                          0,
		"system_total_iops_read_write":                    465200,
		"system_total_iops_write":                         465200,
	}

	assert.Equal(t, expected, job.Collect())
}

type okAPIClient struct {
	loggedIn bool
}

func (o *okAPIClient) Login() error {
	o.loggedIn = true
	return nil
}

func (o *okAPIClient) Logout() error {
	o.loggedIn = false
	return nil
}

func (o okAPIClient) IsLoggedIn() bool { return o.loggedIn }

func (okAPIClient) GetSelectedStatistics(dst interface{}, query string) error {
	r := bytes.NewBuffer(testSelectedStatsData)
	return json.NewDecoder(r).Decode(dst)
}
