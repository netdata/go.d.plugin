package scaleio

import (
	"encoding/json"
	"io/ioutil"
	"net/http/httptest"
	"testing"

	"github.com/netdata/go.d.plugin/modules/scaleio/client"

	"github.com/netdata/go-orchestrator/module"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	selectedStatisticsData, _ = ioutil.ReadFile("testdata/selected_statistics.json")
	instancesData, _          = ioutil.ReadFile("testdata/instances.json")
)

func Test_readTestData(t *testing.T) {
	assert.NotNil(t, selectedStatisticsData)
	assert.NotNil(t, instancesData)
}

func TestNew(t *testing.T) {
	assert.Implements(t, (*module.Module)(nil), New())
}

func TestScaleIO_Init(t *testing.T) {
	scaleIO := New()
	scaleIO.Username = "username"
	scaleIO.Password = "password"

	assert.True(t, scaleIO.Init())
}
func TestScaleIO_Init_UsernameOrPasswordNotSet(t *testing.T) {
	assert.False(t, New().Init())
}

func TestScaleIO_Init_ErrorOnCreatingClientWrongTLSCA(t *testing.T) {
	job := New()
	job.Username = "username"
	job.Password = "password"
	job.ClientTLSConfig.TLSCA = "testdata/tls"

	assert.False(t, job.Init())
}

func TestScaleIO_Check(t *testing.T) {
	srv, _, scaleIO := prepareSrvMockScaleIO(t)
	defer srv.Close()
	require.True(t, scaleIO.Init())

	assert.True(t, scaleIO.Check())
}

func TestScaleIO_Check_ErrorOnLogin(t *testing.T) {
	srv, mock, scaleIO := prepareSrvMockScaleIO(t)
	defer srv.Close()
	require.True(t, scaleIO.Init())
	mock.Password = "new password"

	assert.False(t, scaleIO.Check())
}

func TestScaleIO_Charts(t *testing.T) {
	assert.NotNil(t, New().Charts())
}

func TestScaleIO_Cleanup(t *testing.T) {
	srv, _, scaleIO := prepareSrvMockScaleIO(t)
	defer srv.Close()

	require.True(t, scaleIO.Init())
	require.True(t, scaleIO.Check())

	_ = scaleIO.client.Logout()
	assert.False(t, scaleIO.client.LoggedIn())
}

func TestScaleIO_Collect(t *testing.T) {
	srv, _, scaleIO := prepareSrvMockScaleIO(t)
	defer srv.Close()

	require.True(t, scaleIO.Init())
	require.True(t, scaleIO.Check())

	expected := map[string]int64{
		"sdc_6076fd0f00000000_bandwidth_read":                0,
		"sdc_6076fd0f00000000_bandwidth_read_write":          0,
		"sdc_6076fd0f00000000_bandwidth_write":               0,
		"sdc_6076fd0f00000000_io_size_read":                  0,
		"sdc_6076fd0f00000000_io_size_read_write":            0,
		"sdc_6076fd0f00000000_io_size_write":                 0,
		"sdc_6076fd0f00000000_iops_read":                     0,
		"sdc_6076fd0f00000000_iops_read_write":               0,
		"sdc_6076fd0f00000000_iops_write":                    0,
		"sdc_6076fd0f00000000_mdm_connection_state":          1,
		"sdc_6076fd0f00000000_num_of_mapped_volumes":         1,
		"sdc_6076fd1000000001_bandwidth_read":                2000,
		"sdc_6076fd1000000001_bandwidth_read_write":          127840000,
		"sdc_6076fd1000000001_bandwidth_write":               127838000,
		"sdc_6076fd1000000001_io_size_read":                  2000,
		"sdc_6076fd1000000001_io_size_read_write":            791123,
		"sdc_6076fd1000000001_io_size_write":                 789123,
		"sdc_6076fd1000000001_iops_read":                     1000,
		"sdc_6076fd1000000001_iops_read_write":               163000,
		"sdc_6076fd1000000001_iops_write":                    162000,
		"sdc_6076fd1000000001_mdm_connection_state":          0,
		"sdc_6076fd1000000001_num_of_mapped_volumes":         1,
		"sdc_6076fd1100000002_bandwidth_read":                0,
		"sdc_6076fd1100000002_bandwidth_read_write":          129580000,
		"sdc_6076fd1100000002_bandwidth_write":               129580000,
		"sdc_6076fd1100000002_io_size_read":                  0,
		"sdc_6076fd1100000002_io_size_read_write":            1004496,
		"sdc_6076fd1100000002_io_size_write":                 1004496,
		"sdc_6076fd1100000002_iops_read":                     0,
		"sdc_6076fd1100000002_iops_read_write":               129000,
		"sdc_6076fd1100000002_iops_write":                    129000,
		"sdc_6076fd1100000002_mdm_connection_state":          0,
		"sdc_6076fd1100000002_num_of_mapped_volumes":         1,
		"storage_pool_40395b7b00000000_capacity_utilization": 1510,
		"storage_pool_40395b7b00000000_num_of_devices":       3,
		"storage_pool_40395b7b00000000_num_of_snapshots":     1,
		"storage_pool_40395b7b00000000_num_of_volumes":       3,
		"storage_pool_40395b7b00000000_num_of_vtrees":        2,
		"storage_pool_4039828b00000001_capacity_utilization": 0,
		"storage_pool_4039828b00000001_num_of_devices":       3,
		"storage_pool_4039828b00000001_num_of_snapshots":     0,
		"storage_pool_4039828b00000001_num_of_volumes":       0,
		"storage_pool_4039828b00000001_num_of_vtrees":        0,
		"system_backend_primary_bandwidth_read":              800,
		"system_backend_primary_bandwidth_read_write":        249035199,
		"system_backend_primary_bandwidth_write":             249034400,
		"system_backend_primary_io_size_read":                4000,
		"system_backend_primary_io_size_read_write":          892139,
		"system_backend_primary_io_size_write":               888139,
		"system_backend_primary_iops_read":                   200,
		"system_backend_primary_iops_read_write":             280599,
		"system_backend_primary_iops_write":                  280400,
		"system_backend_secondary_bandwidth_read":            0,
		"system_backend_secondary_bandwidth_read_write":      250278400,
		"system_backend_secondary_bandwidth_write":           250278400,
		"system_backend_secondary_io_size_read":              0,
		"system_backend_secondary_io_size_read_write":        890670,
		"system_backend_secondary_io_size_write":             890670,
		"system_backend_secondary_iops_read":                 0,
		"system_backend_secondary_iops_read_write":           281000,
		"system_backend_secondary_iops_write":                281000,
		"system_backend_total_bandwidth_read":                800,
		"system_backend_total_bandwidth_read_write":          499313600,
		"system_backend_total_bandwidth_write":               499312800,
		"system_backend_total_io_size_read":                  4000,
		"system_backend_total_io_size_read_write":            1782810,
		"system_backend_total_io_size_write":                 1778810,
		"system_backend_total_iops_read":                     200,
		"system_backend_total_iops_read_write":               561600,
		"system_backend_total_iops_write":                    561400,
		"system_capacity_available_for_volume_allocation":    243269632,
		"system_capacity_decreased":                          0,
		"system_capacity_degraded":                           0,
		"system_capacity_failed":                             0,
		"system_capacity_in_maintenance":                     0,
		"system_capacity_in_use":                             42342400,
		"system_capacity_max_capacity":                       643819520,
		"system_capacity_protected":                          42342400,
		"system_capacity_snapshot":                           743424,
		"system_capacity_spare":                              64380928,
		"system_capacity_thick_in_use":                       0,
		"system_capacity_thin_in_use":                        41598976,
		"system_capacity_unreachable_unused":                 0,
		"system_capacity_unused":                             536352768,
		"system_frontend_user_data_bandwidth_read":           3000,
		"system_frontend_user_data_bandwidth_read_write":     257519000,
		"system_frontend_user_data_bandwidth_write":          257516000,
		"system_frontend_user_data_io_size_read":             1500,
		"system_frontend_user_data_io_size_read_write":       859886,
		"system_frontend_user_data_io_size_write":            858386,
		"system_frontend_user_data_iops_read":                2000,
		"system_frontend_user_data_iops_read_write":          302000,
		"system_frontend_user_data_iops_write":               300000,
		"system_num_of_devices":                              6,
		"system_num_of_fault_sets":                           0,
		"system_num_of_mapped_to_all_volumes":                0,
		"system_num_of_mapped_volumes":                       3,
		"system_num_of_protection_domains":                   1,
		"system_num_of_rfcache_devices":                      0,
		"system_num_of_sdc":                                  3,
		"system_num_of_sds":                                  3,
		"system_num_of_snapshots":                            1,
		"system_num_of_storage_pools":                        2,
		"system_num_of_thick_base_volumes":                   0,
		"system_num_of_thin_base_volumes":                    2,
		"system_num_of_unmapped_volumes":                     0,
		"system_num_of_volumes":                              3,
		"system_num_of_vtrees":                               2,
		"system_rebalance_bandwidth_read":                    0,
		"system_rebalance_bandwidth_read_write":              0,
		"system_rebalance_bandwidth_write":                   0,
		"system_rebalance_io_size_read":                      0,
		"system_rebalance_io_size_read_write":                0,
		"system_rebalance_io_size_write":                     0,
		"system_rebalance_iops_read":                         0,
		"system_rebalance_iops_read_write":                   0,
		"system_rebalance_iops_write":                        0,
		"system_rebalance_pending_capacity_in_Kb":            0,
		"system_rebalance_time_until_finish":                 0,
		"system_rebuild_backward_bandwidth_read":             0,
		"system_rebuild_backward_bandwidth_read_write":       0,
		"system_rebuild_backward_bandwidth_write":            0,
		"system_rebuild_backward_io_size_read":               0,
		"system_rebuild_backward_io_size_read_write":         0,
		"system_rebuild_backward_io_size_write":              0,
		"system_rebuild_backward_iops_read":                  0,
		"system_rebuild_backward_iops_read_write":            0,
		"system_rebuild_backward_iops_write":                 0,
		"system_rebuild_backward_pending_capacity_in_Kb":     0,
		"system_rebuild_forward_bandwidth_read":              0,
		"system_rebuild_forward_bandwidth_read_write":        0,
		"system_rebuild_forward_bandwidth_write":             0,
		"system_rebuild_forward_io_size_read":                0,
		"system_rebuild_forward_io_size_read_write":          0,
		"system_rebuild_forward_io_size_write":               0,
		"system_rebuild_forward_iops_read":                   0,
		"system_rebuild_forward_iops_read_write":             0,
		"system_rebuild_forward_iops_write":                  0,
		"system_rebuild_forward_pending_capacity_in_Kb":      0,
		"system_rebuild_normal_bandwidth_read":               0,
		"system_rebuild_normal_bandwidth_read_write":         0,
		"system_rebuild_normal_bandwidth_write":              0,
		"system_rebuild_normal_io_size_read":                 0,
		"system_rebuild_normal_io_size_read_write":           0,
		"system_rebuild_normal_io_size_write":                0,
		"system_rebuild_normal_iops_read":                    0,
		"system_rebuild_normal_iops_read_write":              0,
		"system_rebuild_normal_iops_write":                   0,
		"system_rebuild_normal_pending_capacity_in_Kb":       0,
		"system_rebuild_total_bandwidth_read":                0,
		"system_rebuild_total_bandwidth_read_write":          0,
		"system_rebuild_total_bandwidth_write":               0,
		"system_rebuild_total_io_size_read":                  0,
		"system_rebuild_total_io_size_read_write":            0,
		"system_rebuild_total_io_size_write":                 0,
		"system_rebuild_total_iops_read":                     0,
		"system_rebuild_total_iops_read_write":               0,
		"system_rebuild_total_iops_write":                    0,
		"system_rebuild_total_pending_capacity_in_Kb":        0,
		"system_total_bandwidth_read":                        800,
		"system_total_bandwidth_read_write":                  499313600,
		"system_total_bandwidth_write":                       499312800,
		"system_total_io_size_read":                          4000,
		"system_total_io_size_read_write":                    893406,
		"system_total_io_size_write":                         889406,
		"system_total_iops_read":                             200,
		"system_total_iops_read_write":                       561600,
		"system_total_iops_write":                            561400,
	}

	assert.Equal(t, expected, scaleIO.Collect())
}

func TestScaleIO_Collect_ConnectionRefused(t *testing.T) {
	srv, _, scaleIO := prepareSrvMockScaleIO(t)
	defer srv.Close()

	require.True(t, scaleIO.Init())
	require.True(t, scaleIO.Check())
	scaleIO.client.Request.URL.Host = "127.0.0.1:38001"

	assert.Nil(t, scaleIO.Collect())
}

func prepareSrvMockScaleIO(t *testing.T) (*httptest.Server, *client.MockScaleIOAPIServer, *ScaleIO) {
	t.Helper()
	const (
		user     = "user"
		password = "password"
		version  = "2.5"
		token    = "token"
	)
	var stats client.SelectedStatistics
	err := json.Unmarshal(selectedStatisticsData, &stats)
	require.NoError(t, err)

	var ins client.Instances
	err = json.Unmarshal(instancesData, &ins)
	require.NoError(t, err)

	mock := client.MockScaleIOAPIServer{
		User:       user,
		Password:   password,
		Version:    version,
		Token:      token,
		Instances:  ins,
		Statistics: stats,
	}
	srv := httptest.NewServer(&mock)
	require.NoError(t, err)

	scaleIO := New()
	scaleIO.UserURL = srv.URL
	scaleIO.Username = user
	scaleIO.Password = password
	return srv, &mock, scaleIO
}
