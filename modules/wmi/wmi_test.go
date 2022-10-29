// SPDX-License-Identifier: GPL-3.0-or-later

package wmi

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/netdata/go.d.plugin/pkg/web"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	v0150Metrics, _ = os.ReadFile("testdata/v0.15.0/metrics.txt")
	v0200Metrics, _ = os.ReadFile("testdata/v0.20.0/metrics.txt")
)

func Test_TestData(t *testing.T) {
	for name, data := range map[string][]byte{
		"v0150Metrics": v0150Metrics,
		"v0200Metrics": v0200Metrics,
	} {
		assert.NotNilf(t, data, name)
	}
}

func TestNew(t *testing.T) {
	assert.IsType(t, (*WMI)(nil), New())
}

func TestWMI_Init(t *testing.T) {
	tests := map[string]struct {
		config   Config
		wantFail bool
	}{
		"success if 'url' is set": {
			config: Config{
				HTTP: web.HTTP{Request: web.Request{URL: "http://127.0.0.1:9182/metrics"}}},
		},
		"fails on default config": {
			wantFail: true,
			config:   New().Config,
		},
		"fails if 'url' is unset": {
			wantFail: true,
			config:   Config{HTTP: web.HTTP{Request: web.Request{URL: ""}}},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			wmi := New()
			wmi.Config = test.config

			if test.wantFail {
				assert.False(t, wmi.Init())
			} else {
				assert.True(t, wmi.Init())
			}
		})
	}
}

func TestWMI_Check(t *testing.T) {
	tests := map[string]struct {
		prepare  func() (wmi *WMI, cleanup func())
		wantFail bool
	}{
		"success on valid response v0.15.0": {
			prepare: prepareWMIv0150,
		},
		"success on valid response v0.20.0": {
			prepare: prepareWMIv0200,
		},
		"fails if endpoint returns invalid data": {
			wantFail: true,
			prepare:  prepareWMIReturnsInvalidData,
		},
		"fails on connection refused": {
			wantFail: true,
			prepare:  prepareWMIConnectionRefused,
		},
		"fails on 404 response": {
			wantFail: true,
			prepare:  prepareWMIResponse404,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			wmi, cleanup := test.prepare()
			defer cleanup()

			require.True(t, wmi.Init())

			if test.wantFail {
				assert.False(t, wmi.Check())
			} else {
				assert.True(t, wmi.Check())
			}
		})
	}
}

func TestWMI_Charts(t *testing.T) {
	assert.NotNil(t, New().Charts())
}

func TestWMI_Cleanup(t *testing.T) {
	assert.NotPanics(t, New().Cleanup)
}

func TestWMI_Collect(t *testing.T) {
	tests := map[string]struct {
		prepare       func() (wmi *WMI, cleanup func())
		wantCollected map[string]int64
	}{
		"success on valid response v0.15.0": {
			prepare: prepareWMIv0150,
			wantCollected: map[string]int64{
				"cpu_collection_duration":                                     1000,
				"cpu_collection_success":                                      1,
				"cpu_core_0,0_c1":                                             666516,
				"cpu_core_0,0_c2":                                             1000,
				"cpu_core_0,0_c3":                                             1000,
				"cpu_core_0,0_dpc":                                            6234,
				"cpu_core_0,0_dpcs":                                           352862000,
				"cpu_core_0,0_idle":                                           696218,
				"cpu_core_0,0_interrupt":                                      21359,
				"cpu_core_0,0_interrupts":                                     3799540000,
				"cpu_core_0,0_privileged":                                     517031,
				"cpu_core_0,0_user":                                           402703,
				"cpu_dpc":                                                     6234,
				"cpu_idle":                                                    696218,
				"cpu_interrupt":                                               21359,
				"cpu_privileged":                                              517031,
				"cpu_user":                                                    402703,
				"logical_disk_C:_free_space":                                  8434745344000,
				"logical_disk_C:_idle_seconds_total":                          0,
				"logical_disk_C:_read_bytes_total":                            8458891776000,
				"logical_disk_C:_read_latency":                                143835,
				"logical_disk_C:_read_seconds_total":                          0,
				"logical_disk_C:_reads_total":                                 101079,
				"logical_disk_C:_requests_queued":                             0,
				"logical_disk_C:_split_ios_total":                             0,
				"logical_disk_C:_total_space":                                 21371027456000,
				"logical_disk_C:_used_space":                                  12936282112000,
				"logical_disk_C:_write_bytes_total":                           7427673600000,
				"logical_disk_C:_write_latency":                               39666,
				"logical_disk_C:_write_seconds_total":                         0,
				"logical_disk_C:_writes_total":                                56260,
				"logical_disk_collection_duration":                            1000,
				"logical_disk_collection_success":                             1,
				"logon_collection_duration":                                   1256,
				"logon_collection_success":                                    1,
				"logon_type_batch":                                            1,
				"logon_type_cached_interactive":                               1,
				"logon_type_cached_remote_interactive":                        1,
				"logon_type_cached_unlock":                                    1,
				"logon_type_interactive":                                      2,
				"logon_type_network":                                          1,
				"logon_type_network_clear_text":                               1,
				"logon_type_new_credentials":                                  1,
				"logon_type_proxy":                                            1,
				"logon_type_remote_interactive":                               1,
				"logon_type_service":                                          1,
				"logon_type_system":                                           1,
				"logon_type_unlock":                                           1,
				"memory_available_bytes":                                      788783104000,
				"memory_cache_bytes":                                          68575232000,
				"memory_cache_bytes_peak":                                     102326272000,
				"memory_cache_faults_total":                                   915557000,
				"memory_cache_total":                                          842539008000,
				"memory_collection_duration":                                  1000,
				"memory_collection_success":                                   1,
				"memory_commit_limit":                                         3547709440000,
				"memory_committed_bytes":                                      2657218560000,
				"memory_demand_zero_faults_total":                             6242530000,
				"memory_free_and_zero_page_list_bytes":                        2531328000,
				"memory_free_system_page_table_entries":                       12529874000,
				"memory_modified_page_list_bytes":                             56287232000,
				"memory_not_committed_bytes":                                  890490880000,
				"memory_page_faults_total":                                    17047959000,
				"memory_pool_nonpaged_allocs_total":                           1000,
				"memory_pool_nonpaged_bytes_total":                            97243136000,
				"memory_pool_paged_allocs_total":                              1000,
				"memory_pool_paged_bytes":                                     172675072000,
				"memory_pool_paged_resident_bytes":                            153165824000,
				"memory_standby_cache_core_bytes":                             124506112000,
				"memory_standby_cache_normal_priority_bytes":                  441131008000,
				"memory_standby_cache_reserve_bytes":                          220614656000,
				"memory_standby_cache_total":                                  786251776000,
				"memory_swap_page_operations_total":                           1466380000,
				"memory_swap_page_reads_total":                                127979000,
				"memory_swap_page_writes_total":                               3618000,
				"memory_swap_pages_read_total":                                1240157000,
				"memory_swap_pages_written_total":                             226223000,
				"memory_system_cache_resident_bytes":                          68575232000,
				"memory_system_code_resident_bytes":                           4321280000,
				"memory_system_code_total_bytes":                              4636672000,
				"memory_system_driver_resident_bytes":                         3244032000,
				"memory_system_driver_total_bytes":                            17526784000,
				"memory_transition_faults_total":                              10153909000,
				"memory_transition_pages_repurposed_total":                    1375981000,
				"memory_used_bytes":                                           1358229504000,
				"memory_write_copies_total":                                   105886000,
				"net_Intel_R_PRO_1000_MT_Desktop_Adapter_2_bytes_received":    76499000,
				"net_Intel_R_PRO_1000_MT_Desktop_Adapter_2_bytes_sent":        88865000,
				"net_Intel_R_PRO_1000_MT_Desktop_Adapter_2_bytes_total":       165364000,
				"net_Intel_R_PRO_1000_MT_Desktop_Adapter_2_current_bandwidth": 1000000000,
				"net_Intel_R_PRO_1000_MT_Desktop_Adapter_2_packets_outbound_discarded": 1000,
				"net_Intel_R_PRO_1000_MT_Desktop_Adapter_2_packets_outbound_errors":    1000,
				"net_Intel_R_PRO_1000_MT_Desktop_Adapter_2_packets_received_discarded": 1000,
				"net_Intel_R_PRO_1000_MT_Desktop_Adapter_2_packets_received_errors":    1000,
				"net_Intel_R_PRO_1000_MT_Desktop_Adapter_2_packets_received_total":     676000,
				"net_Intel_R_PRO_1000_MT_Desktop_Adapter_2_packets_received_unknown":   1000,
				"net_Intel_R_PRO_1000_MT_Desktop_Adapter_2_packets_sent_total":         686000,
				"net_Intel_R_PRO_1000_MT_Desktop_Adapter_2_packets_total":              1362000,
				"net_Intel_R_PRO_1000_MT_Desktop_Adapter_bytes_received":               383489027000,
				"net_Intel_R_PRO_1000_MT_Desktop_Adapter_bytes_sent":                   6755954000,
				"net_Intel_R_PRO_1000_MT_Desktop_Adapter_bytes_total":                  390244981000,
				"net_Intel_R_PRO_1000_MT_Desktop_Adapter_current_bandwidth":            1000000000,
				"net_Intel_R_PRO_1000_MT_Desktop_Adapter_packets_outbound_discarded":   1000,
				"net_Intel_R_PRO_1000_MT_Desktop_Adapter_packets_outbound_errors":      1000,
				"net_Intel_R_PRO_1000_MT_Desktop_Adapter_packets_received_discarded":   1000,
				"net_Intel_R_PRO_1000_MT_Desktop_Adapter_packets_received_errors":      1000,
				"net_Intel_R_PRO_1000_MT_Desktop_Adapter_packets_received_total":       262638000,
				"net_Intel_R_PRO_1000_MT_Desktop_Adapter_packets_received_unknown":     1000,
				"net_Intel_R_PRO_1000_MT_Desktop_Adapter_packets_sent_total":           84041000,
				"net_Intel_R_PRO_1000_MT_Desktop_Adapter_packets_total":                346679000,
				"net_collection_duration":           1000,
				"net_collection_success":            1,
				"os_collection_duration":            1127,
				"os_collection_success":             1,
				"os_paging_free_bytes":              1056043008000,
				"os_paging_limit_bytes":             1400696832000,
				"os_paging_used_bytes":              344653824000,
				"os_physical_memory_free_bytes":     798547968000,
				"os_process_memory_limit_bytes":     0,
				"os_processes":                      79,
				"os_processes_limit":                4294967295,
				"os_time":                           1615228118,
				"os_users":                          2,
				"os_virtual_memory_bytes":           3547709440000,
				"os_virtual_memory_free_bytes":      905457664000,
				"os_visible_memory_bytes":           2147012608000,
				"os_visible_memory_used_bytes":      1348464640000,
				"system_boot_time":                  1615226502,
				"system_calls_total":                41784660000,
				"system_collection_duration":        1000,
				"system_collection_success":         1,
				"system_context_switches_total":     4616757000,
				"system_exception_dispatches_total": 3695000,
				"system_processor_queue_length":     10,
				"system_threads":                    967,
				"system_up_time":                    10424243,
			},
		},
		"success on valid response v0.20.0": {
			prepare: prepareWMIv0200,
			wantCollected: map[string]int64{
				"cpu_collection_duration":                                       0,
				"cpu_collection_success":                                        1,
				"cpu_core_0,0_c1":                                               1841769627,
				"cpu_core_0,0_c2":                                               0,
				"cpu_core_0,0_c3":                                               0,
				"cpu_core_0,0_dpc":                                              89046,
				"cpu_core_0,0_dpcs":                                             9205934000,
				"cpu_core_0,0_idle":                                             1844928953,
				"cpu_core_0,0_interrupt":                                        79750,
				"cpu_core_0,0_interrupts":                                       164487760000,
				"cpu_core_0,0_privileged":                                       3507703,
				"cpu_core_0,0_user":                                             3949109,
				"cpu_core_0,1_c1":                                               1843147617,
				"cpu_core_0,1_c2":                                               0,
				"cpu_core_0,1_c3":                                               0,
				"cpu_core_0,1_dpc":                                              2140,
				"cpu_core_0,1_dpcs":                                             1467612000,
				"cpu_core_0,1_idle":                                             1846598625,
				"cpu_core_0,1_interrupt":                                        130796,
				"cpu_core_0,1_interrupts":                                       173675048000,
				"cpu_core_0,1_privileged":                                       2806625,
				"cpu_core_0,1_user":                                             2979500,
				"cpu_dpc":                                                       91187,
				"cpu_idle":                                                      3691527578,
				"cpu_interrupt":                                                 210546,
				"cpu_privileged":                                                6314328,
				"cpu_user":                                                      6928609,
				"logical_disk_C:_free_space":                                    31603032064000,
				"logical_disk_C:_idle_seconds_total":                            0,
				"logical_disk_C:_read_bytes_total":                              121832960000,
				"logical_disk_C:_read_latency":                                  3586,
				"logical_disk_C:_read_seconds_total":                            0,
				"logical_disk_C:_reads_total":                                   5316,
				"logical_disk_C:_requests_queued":                               0,
				"logical_disk_C:_split_ios_total":                               0,
				"logical_disk_C:_total_space":                                   53683945472000,
				"logical_disk_C:_used_space":                                    22080913408000,
				"logical_disk_C:_write_bytes_total":                             195316224000,
				"logical_disk_C:_write_latency":                                 8492,
				"logical_disk_C:_write_seconds_total":                           0,
				"logical_disk_C:_writes_total":                                  9374,
				"logical_disk_collection_duration":                              0,
				"logical_disk_collection_success":                               1,
				"logon_collection_duration":                                     72,
				"logon_collection_success":                                      1,
				"logon_type_batch":                                              0,
				"logon_type_cached_interactive":                                 0,
				"logon_type_cached_remote_interactive":                          1,
				"logon_type_cached_unlock":                                      0,
				"logon_type_interactive":                                        7,
				"logon_type_network":                                            1,
				"logon_type_network_clear_text":                                 0,
				"logon_type_new_credentials":                                    0,
				"logon_type_proxy":                                              0,
				"logon_type_remote_interactive":                                 1,
				"logon_type_service":                                            5,
				"logon_type_system":                                             1,
				"logon_type_unlock":                                             0,
				"memory_available_bytes":                                        1938132992000,
				"memory_cache_bytes":                                            145616896000,
				"memory_cache_bytes_peak":                                       195719168000,
				"memory_cache_faults_total":                                     14334571000,
				"memory_cache_total":                                            1570410496000,
				"memory_collection_duration":                                    0,
				"memory_collection_success":                                     1,
				"memory_commit_limit":                                           5720535040000,
				"memory_committed_bytes":                                        2370998272000,
				"memory_demand_zero_faults_total":                               61122648000,
				"memory_free_and_zero_page_list_bytes":                          375119872000,
				"memory_free_system_page_table_entries":                         12524106000,
				"memory_modified_page_list_bytes":                               7401472000,
				"memory_not_committed_bytes":                                    3349536768000,
				"memory_page_faults_total":                                      224955133000,
				"memory_pool_nonpaged_allocs_total":                             0,
				"memory_pool_nonpaged_bytes_total":                              0,
				"memory_pool_paged_allocs_total":                                0,
				"memory_pool_paged_bytes":                                       275705856000,
				"memory_pool_paged_resident_bytes":                              246333440000,
				"memory_standby_cache_core_bytes":                               0,
				"memory_standby_cache_normal_priority_bytes":                    1093529600000,
				"memory_standby_cache_reserve_bytes":                            469479424000,
				"memory_standby_cache_total":                                    1563009024000,
				"memory_swap_page_operations_total":                             7595386000,
				"memory_swap_page_reads_total":                                  608527000,
				"memory_swap_page_writes_total":                                 2934000,
				"memory_swap_pages_read_total":                                  6847572000,
				"memory_swap_pages_written_total":                               747814000,
				"memory_system_cache_resident_bytes":                            145616896000,
				"memory_system_code_resident_bytes":                             17100800000,
				"memory_system_code_total_bytes":                                8192000,
				"memory_system_driver_resident_bytes":                           38674432000,
				"memory_system_driver_total_bytes":                              28409856000,
				"memory_transition_faults_total":                                143833525000,
				"memory_transition_pages_repurposed_total":                      4769338000,
				"memory_used_bytes":                                             2306007040000,
				"memory_write_copies_total":                                     6732296000,
				"net_Amazon_Elastic_Network_Adapter_bytes_received":             1002322186000,
				"net_Amazon_Elastic_Network_Adapter_bytes_sent":                 601090821000,
				"net_Amazon_Elastic_Network_Adapter_bytes_total":                1603413007000,
				"net_Amazon_Elastic_Network_Adapter_current_bandwidth":          0,
				"net_Amazon_Elastic_Network_Adapter_packets_outbound_discarded": 0,
				"net_Amazon_Elastic_Network_Adapter_packets_outbound_errors":    0,
				"net_Amazon_Elastic_Network_Adapter_packets_received_discarded": 0,
				"net_Amazon_Elastic_Network_Adapter_packets_received_errors":    0,
				"net_Amazon_Elastic_Network_Adapter_packets_received_total":     2740032000,
				"net_Amazon_Elastic_Network_Adapter_packets_received_unknown":   0,
				"net_Amazon_Elastic_Network_Adapter_packets_sent_total":         2083285000,
				"net_Amazon_Elastic_Network_Adapter_packets_total":              4823317000,
				"net_collection_duration":                                       0,
				"net_collection_success":                                        1,
				"os_collection_duration":                                        4,
				"os_collection_success":                                         1,
				"os_paging_free_bytes":                                          1275879424000,
				"os_paging_limit_bytes":                                         1476395008000,
				"os_paging_used_bytes":                                          200515584000,
				"os_physical_memory_free_bytes":                                 1937801216000,
				"os_process_memory_limit_bytes":                                 140737488224256000,
				"os_processes":                                                  119,
				"os_processes_limit":                                            4294967295,
				"os_time":                                                       1666394276,
				"os_users":                                                      1,
				"os_virtual_memory_bytes":                                       5720535040000,
				"os_virtual_memory_free_bytes":                                  3349536768000,
				"os_visible_memory_bytes":                                       4244140032000,
				"os_visible_memory_used_bytes":                                  2306338816000,
				"process_Idle_cpu_time":                                         4533459734,
				"process_Idle_handles":                                          0,
				"process_Idle_io_bytes":                                         0,
				"process_Idle_io_operations":                                    0,
				"process_Idle_page_faults":                                      9,
				"process_Idle_page_file_bytes":                                  61440,
				"process_Idle_pool_bytes":                                       272,
				"process_Idle_threads":                                          4,
				"process_LogonUI_cpu_time":                                      5859,
				"process_LogonUI_handles":                                       447,
				"process_LogonUI_io_bytes":                                      19334,
				"process_LogonUI_io_operations":                                 3636319,
				"process_LogonUI_page_faults":                                   3650284,
				"process_LogonUI_page_file_bytes":                               11075584,
				"process_LogonUI_pool_bytes":                                    465048,
				"process_LogonUI_threads":                                       9,
				"process_MsMpEng_cpu_time":                                      76984,
				"process_MsMpEng_handles":                                       748,
				"process_MsMpEng_io_bytes":                                      3243817169,
				"process_MsMpEng_io_operations":                                 759480,
				"process_MsMpEng_page_faults":                                   1303224,
				"process_MsMpEng_page_file_bytes":                               218316800,
				"process_MsMpEng_pool_bytes":                                    668568,
				"process_MsMpEng_threads":                                       33,
				"process_NisSrv_cpu_time":                                       62,
				"process_NisSrv_handles":                                        211,
				"process_NisSrv_io_bytes":                                       96137,
				"process_NisSrv_io_operations":                                  811,
				"process_NisSrv_page_faults":                                    3700,
				"process_NisSrv_page_file_bytes":                                3629056,
				"process_NisSrv_pool_bytes":                                     132232,
				"process_NisSrv_threads":                                        4,
				"process_Registry_cpu_time":                                     4906,
				"process_Registry_handles":                                      0,
				"process_Registry_io_bytes":                                     412273567,
				"process_Registry_io_operations":                                57804,
				"process_Registry_page_faults":                                  517444,
				"process_Registry_page_file_bytes":                              1585152,
				"process_Registry_pool_bytes":                                   204912,
				"process_Registry_threads":                                      4,
				"process_RuntimeBroker_cpu_time":                                3468,
				"process_RuntimeBroker_handles":                                 934,
				"process_RuntimeBroker_io_bytes":                                1266140,
				"process_RuntimeBroker_io_operations":                           18580,
				"process_RuntimeBroker_page_faults":                             73689,
				"process_RuntimeBroker_page_file_bytes":                         17141760,
				"process_RuntimeBroker_pool_bytes":                              711800,
				"process_RuntimeBroker_threads":                                 14,
				"process_SearchApp_cpu_time":                                    10625,
				"process_SearchApp_handles":                                     1117,
				"process_SearchApp_io_bytes":                                    49046573,
				"process_SearchApp_io_operations":                               16800,
				"process_SearchApp_page_faults":                                 191428,
				"process_SearchApp_page_file_bytes":                             91574272,
				"process_SearchApp_pool_bytes":                                  942176,
				"process_SearchApp_threads":                                     35,
				"process_SecurityHealthService_cpu_time":                        281,
				"process_SecurityHealthService_handles":                         198,
				"process_SecurityHealthService_io_bytes":                        10454,
				"process_SecurityHealthService_io_operations":                   507,
				"process_SecurityHealthService_page_faults":                     9001,
				"process_SecurityHealthService_page_file_bytes":                 2121728,
				"process_SecurityHealthService_pool_bytes":                      102992,
				"process_SecurityHealthService_threads":                         3,
				"process_StartMenuExperienceHost_cpu_time":                      5562,
				"process_StartMenuExperienceHost_handles":                       577,
				"process_StartMenuExperienceHost_io_bytes":                      47667,
				"process_StartMenuExperienceHost_io_operations":                 966,
				"process_StartMenuExperienceHost_page_faults":                   52930,
				"process_StartMenuExperienceHost_page_file_bytes":               20430848,
				"process_StartMenuExperienceHost_pool_bytes":                    712240,
				"process_StartMenuExperienceHost_threads":                       9,
				"process_System_cpu_time":                                       233546,
				"process_System_handles":                                        2633,
				"process_System_io_bytes":                                       968210777,
				"process_System_io_operations":                                  65017,
				"process_System_page_faults":                                    3207,
				"process_System_page_file_bytes":                                40960,
				"process_System_pool_bytes":                                     272,
				"process_System_threads":                                        124,
				"process_TextInputHost_cpu_time":                                1343,
				"process_TextInputHost_handles":                                 541,
				"process_TextInputHost_io_bytes":                                10120,
				"process_TextInputHost_io_operations":                           787,
				"process_TextInputHost_page_faults":                             47606,
				"process_TextInputHost_page_file_bytes":                         10588160,
				"process_TextInputHost_pool_bytes":                              504952,
				"process_TextInputHost_threads":                                 12,
				"process_WMIC_cpu_time":                                         0,
				"process_WMIC_handles":                                          97,
				"process_WMIC_io_bytes":                                         114,
				"process_WMIC_io_operations":                                    29,
				"process_WMIC_page_faults":                                      1374,
				"process_WMIC_page_file_bytes":                                  1081344,
				"process_WMIC_pool_bytes":                                       74584,
				"process_WMIC_threads":                                          3,
				"process_WUDFHost_cpu_time":                                     46,
				"process_WUDFHost_handles":                                      331,
				"process_WUDFHost_io_bytes":                                     4890,
				"process_WUDFHost_io_operations":                                336,
				"process_WUDFHost_page_faults":                                  5128,
				"process_WUDFHost_page_file_bytes":                              3493888,
				"process_WUDFHost_pool_bytes":                                   172744,
				"process_WUDFHost_threads":                                      12,
				"process_WmiPrvSE_cpu_time":                                     5250,
				"process_WmiPrvSE_handles":                                      818,
				"process_WmiPrvSE_io_bytes":                                     14282749,
				"process_WmiPrvSE_io_operations":                                17186,
				"process_WmiPrvSE_page_faults":                                  166367,
				"process_WmiPrvSE_page_file_bytes":                              22876160,
				"process_WmiPrvSE_pool_bytes":                                   361528,
				"process_WmiPrvSE_threads":                                      27,
				"process_amazon-ssm-agent_cpu_time":                             510562,
				"process_amazon-ssm-agent_handles":                              147,
				"process_amazon-ssm-agent_io_bytes":                             43298219,
				"process_amazon-ssm-agent_io_operations":                        2520787,
				"process_amazon-ssm-agent_page_faults":                          8700668,
				"process_amazon-ssm-agent_page_file_bytes":                      17702912,
				"process_amazon-ssm-agent_pool_bytes":                           71832,
				"process_amazon-ssm-agent_threads":                              11,
				"process_collection_duration":                                   1056,
				"process_collection_success":                                    1,
				"process_conhost_cpu_time":                                      22687,
				"process_conhost_handles":                                       656,
				"process_conhost_io_bytes":                                      34146517,
				"process_conhost_io_operations":                                 716385,
				"process_conhost_page_faults":                                   16641,
				"process_conhost_page_file_bytes":                               21712896,
				"process_conhost_pool_bytes":                                    611384,
				"process_conhost_threads":                                       14,
				"process_csrss_cpu_time":                                        11468,
				"process_csrss_handles":                                         1163,
				"process_csrss_io_bytes":                                        13846412,
				"process_csrss_io_operations":                                   242189,
				"process_csrss_page_faults":                                     99474,
				"process_csrss_page_file_bytes":                                 6041600,
				"process_csrss_pool_bytes":                                      590384,
				"process_csrss_threads":                                         32,
				"process_ctfmon_cpu_time":                                       1234,
				"process_ctfmon_handles":                                        404,
				"process_ctfmon_io_bytes":                                       60266,
				"process_ctfmon_io_operations":                                  2152,
				"process_ctfmon_page_faults":                                    5769,
				"process_ctfmon_page_file_bytes":                                3870720,
				"process_ctfmon_pool_bytes":                                     206560,
				"process_ctfmon_threads":                                        9,
				"process_dllhost_cpu_time":                                      328,
				"process_dllhost_handles":                                       242,
				"process_dllhost_io_bytes":                                      20566486,
				"process_dllhost_io_operations":                                 4898,
				"process_dllhost_page_faults":                                   7240,
				"process_dllhost_page_file_bytes":                               4747264,
				"process_dllhost_pool_bytes":                                    176032,
				"process_dllhost_threads":                                       7,
				"process_dwm_cpu_time":                                          23500,
				"process_dwm_handles":                                           1404,
				"process_dwm_io_bytes":                                          1060550,
				"process_dwm_io_operations":                                     35004,
				"process_dwm_page_faults":                                       1492556,
				"process_dwm_page_file_bytes":                                   43937792,
				"process_dwm_pool_bytes":                                        900528,
				"process_dwm_threads":                                           33,
				"process_explorer_cpu_time":                                     43656,
				"process_explorer_handles":                                      2109,
				"process_explorer_io_bytes":                                     13917564,
				"process_explorer_io_operations":                                89512,
				"process_explorer_page_faults":                                  383795,
				"process_explorer_page_file_bytes":                              43700224,
				"process_explorer_pool_bytes":                                   1074560,
				"process_explorer_threads":                                      67,
				"process_fontdrvhost_cpu_time":                                  343,
				"process_fontdrvhost_handles":                                   117,
				"process_fontdrvhost_io_bytes":                                  2264,
				"process_fontdrvhost_io_operations":                             84,
				"process_fontdrvhost_page_faults":                               3961,
				"process_fontdrvhost_page_file_bytes":                           4407296,
				"process_fontdrvhost_pool_bytes":                                137952,
				"process_fontdrvhost_threads":                                   15,
				"process_lsass_cpu_time":                                        133484,
				"process_lsass_handles":                                         1307,
				"process_lsass_io_bytes":                                        150024212,
				"process_lsass_io_operations":                                   690357,
				"process_lsass_page_faults":                                     403572,
				"process_lsass_page_file_bytes":                                 7086080,
				"process_lsass_pool_bytes":                                      175536,
				"process_lsass_threads":                                         8,
				"process_msdtc_cpu_time":                                        125,
				"process_msdtc_handles":                                         234,
				"process_msdtc_io_bytes":                                        87454,
				"process_msdtc_io_operations":                                   219,
				"process_msdtc_page_faults":                                     3177,
				"process_msdtc_page_file_bytes":                                 3158016,
				"process_msdtc_pool_bytes":                                      107704,
				"process_msdtc_threads":                                         9,
				"process_powershell_cpu_time":                                   1906,
				"process_powershell_handles":                                    1385,
				"process_powershell_io_bytes":                                   1618302,
				"process_powershell_io_operations":                              9197,
				"process_powershell_page_faults":                                54849,
				"process_powershell_page_file_bytes":                            131256320,
				"process_powershell_pool_bytes":                                 1041088,
				"process_powershell_threads":                                    38,
				"process_rdpclip_cpu_time":                                      2203,
				"process_rdpclip_handles":                                       360,
				"process_rdpclip_io_bytes":                                      1731036,
				"process_rdpclip_io_operations":                                 821,
				"process_rdpclip_page_faults":                                   6065,
				"process_rdpclip_page_file_bytes":                               2818048,
				"process_rdpclip_pool_bytes":                                    198840,
				"process_rdpclip_threads":                                       9,
				"process_services_cpu_time":                                     27406,
				"process_services_handles":                                      606,
				"process_services_io_bytes":                                     1584188,
				"process_services_io_operations":                                193393,
				"process_services_page_faults":                                  289962,
				"process_services_page_file_bytes":                              5115904,
				"process_services_pool_bytes":                                   157208,
				"process_services_threads":                                      5,
				"process_sihost_cpu_time":                                       1796,
				"process_sihost_handles":                                        515,
				"process_sihost_io_bytes":                                       674100,
				"process_sihost_io_operations":                                  4753,
				"process_sihost_page_faults":                                    36976,
				"process_sihost_page_file_bytes":                                5259264,
				"process_sihost_pool_bytes":                                     287312,
				"process_sihost_threads":                                        12,
				"process_smartscreen_cpu_time":                                  140,
				"process_smartscreen_handles":                                   452,
				"process_smartscreen_io_bytes":                                  23200,
				"process_smartscreen_io_operations":                             675,
				"process_smartscreen_page_faults":                               6770,
				"process_smartscreen_page_file_bytes":                           8376320,
				"process_smartscreen_pool_bytes":                                259200,
				"process_smartscreen_threads":                                   10,
				"process_smss_cpu_time":                                         187,
				"process_smss_handles":                                          60,
				"process_smss_io_bytes":                                         28818,
				"process_smss_io_operations":                                    514,
				"process_smss_page_faults":                                      1004,
				"process_smss_page_file_bytes":                                  1118208,
				"process_smss_pool_bytes":                                       16992,
				"process_smss_threads":                                          2,
				"process_spoolsv_cpu_time":                                      921,
				"process_spoolsv_handles":                                       455,
				"process_spoolsv_io_bytes":                                      39008,
				"process_spoolsv_io_operations":                                 975,
				"process_spoolsv_page_faults":                                   21133,
				"process_spoolsv_page_file_bytes":                               5890048,
				"process_spoolsv_pool_bytes":                                    186024,
				"process_spoolsv_threads":                                       9,
				"process_sqlceip_cpu_time":                                      22093,
				"process_sqlceip_handles":                                       793,
				"process_sqlceip_io_bytes":                                      263412813,
				"process_sqlceip_io_operations":                                 273865,
				"process_sqlceip_page_faults":                                   1208534,
				"process_sqlceip_page_file_bytes":                               42561536,
				"process_sqlceip_pool_bytes":                                    426752,
				"process_sqlceip_threads":                                       12,
				"process_sqlservr_cpu_time":                                     618203,
				"process_sqlservr_handles":                                      732,
				"process_sqlservr_io_bytes":                                     393531102,
				"process_sqlservr_io_operations":                                248517,
				"process_sqlservr_page_faults":                                  805403,
				"process_sqlservr_page_file_bytes":                              399376384,
				"process_sqlservr_pool_bytes":                                   589888,
				"process_sqlservr_threads":                                      71,
				"process_sqlwriter_cpu_time":                                    15,
				"process_sqlwriter_handles":                                     139,
				"process_sqlwriter_io_bytes":                                    21888,
				"process_sqlwriter_io_operations":                               462,
				"process_sqlwriter_page_faults":                                 2242,
				"process_sqlwriter_page_file_bytes":                             1814528,
				"process_sqlwriter_pool_bytes":                                  94680,
				"process_sqlwriter_threads":                                     2,
				"process_ssm-agent-worker_cpu_time":                             3840500,
				"process_ssm-agent-worker_handles":                              381,
				"process_ssm-agent-worker_io_bytes":                             525476070,
				"process_ssm-agent-worker_io_operations":                        13226879,
				"process_ssm-agent-worker_page_faults":                          991729,
				"process_ssm-agent-worker_page_file_bytes":                      24788992,
				"process_ssm-agent-worker_pool_bytes":                           148640,
				"process_ssm-agent-worker_threads":                              19,
				"process_svchost_cpu_time":                                      2828546,
				"process_svchost_handles":                                       18080,
				"process_svchost_io_bytes":                                      9910037367,
				"process_svchost_io_operations":                                 14895436,
				"process_svchost_page_faults":                                   22339174,
				"process_svchost_page_file_bytes":                               387096576,
				"process_svchost_pool_bytes":                                    8922104,
				"process_svchost_threads":                                       419,
				"process_taskhostw_cpu_time":                                    62,
				"process_taskhostw_handles":                                     215,
				"process_taskhostw_io_bytes":                                    8369888,
				"process_taskhostw_io_operations":                               4613,
				"process_taskhostw_page_faults":                                 7169,
				"process_taskhostw_page_file_bytes":                             4538368,
				"process_taskhostw_pool_bytes":                                  191784,
				"process_taskhostw_threads":                                     4,
				"process_windows_exporter-0.20.0-amd64_cpu_time":                1828,
				"process_windows_exporter-0.20.0-amd64_handles":                 304,
				"process_windows_exporter-0.20.0-amd64_io_bytes":                1205058,
				"process_windows_exporter-0.20.0-amd64_io_operations":           5442,
				"process_windows_exporter-0.20.0-amd64_page_faults":             145453,
				"process_windows_exporter-0.20.0-amd64_page_file_bytes":         32546816,
				"process_windows_exporter-0.20.0-amd64_pool_bytes":              201160,
				"process_windows_exporter-0.20.0-amd64_threads":                 15,
				"process_wininit_cpu_time":                                      1578,
				"process_wininit_handles":                                       153,
				"process_wininit_io_bytes":                                      39300,
				"process_wininit_io_operations":                                 344,
				"process_wininit_page_faults":                                   47683,
				"process_wininit_page_file_bytes":                               1490944,
				"process_wininit_pool_bytes":                                    88456,
				"process_wininit_threads":                                       2,
				"process_winlogon_cpu_time":                                     390,
				"process_winlogon_handles":                                      470,
				"process_winlogon_io_bytes":                                     545790,
				"process_winlogon_io_operations":                                974,
				"process_winlogon_page_faults":                                  13396,
				"process_winlogon_page_file_bytes":                              4476928,
				"process_winlogon_pool_bytes":                                   299336,
				"process_winlogon_threads":                                      8,
				"service_collection_duration":                                   0,
				"service_collection_success":                                    1,
				"system_boot_time":                                              0,
				"system_calls_total":                                            0,
				"system_collection_duration":                                    0,
				"system_collection_success":                                     1,
				"system_context_switches_total":                                 0,
				"system_exception_dispatches_total":                             0,
				"system_processor_queue_length":                                 0,
				"system_threads":                                                0,
				"system_up_time":                                                1667076471,
				"tcp_collection_duration":                                       0,
				"tcp_collection_success":                                        1,
				"tcp_conns_active_ipv4":                                         96883,
				"tcp_conns_active_ipv6":                                         67,
				"tcp_conns_established_ipv4":                                    3,
				"tcp_conns_established_ipv6":                                    2,
				"tcp_conns_failures_ipv4":                                       7587,
				"tcp_conns_failures_ipv6":                                       6,
				"tcp_conns_passive_ipv4":                                        20108,
				"tcp_conns_passive_ipv6":                                        61,
				"tcp_conns_resets_ipv4":                                         14793,
				"tcp_conns_resets_ipv6":                                         2,
				"tcp_segments_received_ipv4":                                    2322665,
				"tcp_segments_received_ipv6":                                    6350,
				"tcp_segments_retransmitted_ipv4":                               29907,
				"tcp_segments_retransmitted_ipv6":                               24,
				"tcp_segments_sent_ipv4":                                        1846067,
				"tcp_segments_sent_ipv6":                                        6076,
				"thermalzone_collection_duration":                               1714,
				"thermalzone_collection_success":                                0,
			},
		},
		"fails if endpoint returns invalid data": {
			prepare: prepareWMIReturnsInvalidData,
		},
		"fails on connection refused": {
			prepare: prepareWMIConnectionRefused,
		},
		"fails on 404 response": {
			prepare: prepareWMIResponse404,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			wmi, cleanup := test.prepare()
			defer cleanup()

			require.True(t, wmi.Init())

			collected := wmi.Collect()

			if collected != nil && test.wantCollected != nil {
				collected["system_up_time"] = test.wantCollected["system_up_time"]
			}

			assert.Equal(t, test.wantCollected, collected)
			if len(test.wantCollected) > 0 {
				testCharts(t, wmi, collected)
			}
		})
	}
}
func testCharts(t *testing.T, wmi *WMI, collected map[string]int64) {
	ensureChartsCreated(t, wmi)
	ensureChartsDynamicDimsCreated(t, wmi)
	ensureCollectedHasAllChartsDimsVarsIDs(t, wmi, collected)
}

func ensureChartsCreated(t *testing.T, w *WMI) {
	for _, chart := range newCPUCharts() {
		if w.cache.collectors[collectorCPU] {
			assert.Truef(t, w.Charts().Has(chart.ID), "chart '%s' not created", chart.ID)
		} else {
			assert.Falsef(t, w.Charts().Has(chart.ID), "chart '%s' created", chart.ID)
		}
	}
	for _, chart := range newMemCharts() {
		if w.cache.collectors[collectorMemory] {
			assert.Truef(t, w.Charts().Has(chart.ID), "chart '%s' not created", chart.ID)
		} else {
			assert.Falsef(t, w.Charts().Has(chart.ID), "chart '%s' created", chart.ID)
		}
	}
	for _, chart := range newOSCharts() {
		if w.cache.collectors[collectorOS] {
			assert.Truef(t, w.Charts().Has(chart.ID), "chart '%s' not created", chart.ID)
		} else {
			assert.Falsef(t, w.Charts().Has(chart.ID), "chart '%s' created", chart.ID)
		}
	}
	for _, chart := range newSystemCharts() {
		if w.cache.collectors[collectorSystem] {
			assert.Truef(t, w.Charts().Has(chart.ID), "chart '%s' not created", chart.ID)
		} else {
			assert.Falsef(t, w.Charts().Has(chart.ID), "chart '%s' created", chart.ID)
		}
	}
	for _, chart := range newLogonCharts() {
		if w.cache.collectors[collectorLogon] {
			assert.Truef(t, w.Charts().Has(chart.ID), "chart '%s' not created", chart.ID)
		} else {
			assert.Falsef(t, w.Charts().Has(chart.ID), "chart '%s' created", chart.ID)
		}
	}
	for _, chart := range newThermalzoneCharts() {
		if w.cache.collectors[collectorThermalzone] {
			assert.Truef(t, w.Charts().Has(chart.ID), "chart '%s' not created", chart.ID)
		} else {
			assert.Falsef(t, w.Charts().Has(chart.ID), "chart '%s' created", chart.ID)
		}
	}
	for _, chart := range *newCollectionCharts() {
		assert.Truef(t, w.Charts().Has(chart.ID), "chart '%s' not created", chart.ID)
	}

	for coreID := range w.cache.cores {
		for _, chart := range newCPUCoreCharts() {
			id := fmt.Sprintf(chart.ID, coreID)
			assert.Truef(t, w.Charts().Has(id), "charts has no '%s' chart for '%s' core", id, coreID)
		}
	}

	for nicID := range w.cache.nics {
		for _, chart := range newNICCharts() {
			id := fmt.Sprintf(chart.ID, nicID)
			assert.Truef(t, w.Charts().Has(id), "charts has no '%s' chart for '%s' nic", id, nicID)
		}
	}

	for diskID := range w.cache.volumes {
		for _, chart := range newDiskCharts() {
			id := fmt.Sprintf(chart.ID, diskID)
			assert.Truef(t, w.Charts().Has(id), "charts has no '%s' chart for '%s' disk", id, diskID)
		}
	}
}

func ensureChartsDynamicDimsCreated(t *testing.T, w *WMI) {
	for coreID := range w.cache.cores {
		chart := w.Charts().Get(cpuDPCsChart.ID)
		if chart != nil {
			dimID := fmt.Sprintf("cpu_core_%s_dpc", coreID)
			assert.Truef(t, chart.HasDim(dimID), "chart '%s' has not dim '%s' for core '%s'", chart.ID, dimID, coreID)
		}

		chart = w.Charts().Get(cpuInterruptsChart.ID)
		if chart != nil {
			dimID := fmt.Sprintf("cpu_core_%s_interrupts", coreID)
			assert.Truef(t, chart.HasDim(dimID), "chart '%s' has not dim '%s' for core '%s'", chart.ID, dimID, coreID)
		}
	}

	for zone := range w.cache.thermalZones {
		chart := w.Charts().Get(thermalzoneTemperatureChart.ID)
		if chart != nil {
			dimID := fmt.Sprintf("thermalzone_%s_temperature", zone)
			assert.Truef(t, chart.HasDim(dimID), "chart '%s' has not dim '%s' for core '%s'", chart.ID, dimID, zone)
		}
	}

	for colID := range w.cache.collectors {
		chart := w.Charts().Get(collectionDurationChart.ID)
		if chart != nil {
			dimID := colID + "_collection_duration"
			assert.Truef(t, chart.HasDim(dimID), "chart '%s' has not dim '%s' for collector '%s'", chart.ID, dimID, colID)
		}

		chart = w.Charts().Get(collectionsStatusChart.ID)
		if chart != nil {
			dimID := colID + "_collection_success"
			assert.Truef(t, chart.HasDim(dimID), "chart '%s' has not dim '%s' for collector '%s'", chart.ID, dimID, colID)
		}
	}
}

func ensureCollectedHasAllChartsDimsVarsIDs(t *testing.T, w *WMI, collected map[string]int64) {
	for _, chart := range *w.Charts() {
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

func prepareWMIv0150() (wmi *WMI, cleanup func()) {
	ts := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			_, _ = w.Write(v0150Metrics)
		}))

	wmi = New()
	wmi.URL = ts.URL
	return wmi, ts.Close
}

func prepareWMIv0200() (wmi *WMI, cleanup func()) {
	ts := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			_, _ = w.Write(v0200Metrics)
		}))

	wmi = New()
	wmi.URL = ts.URL
	return wmi, ts.Close
}

func prepareWMIReturnsInvalidData() (wmi *WMI, cleanup func()) {
	ts := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			_, _ = w.Write([]byte("hello and\n goodbye"))
		}))

	wmi = New()
	wmi.URL = ts.URL
	return wmi, ts.Close
}

func prepareWMIConnectionRefused() (wmi *WMI, cleanup func()) {
	wmi = New()
	wmi.URL = "http://127.0.0.1:38001"
	return wmi, func() {}
}

func prepareWMIResponse404() (wmi *WMI, cleanup func()) {
	ts := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotFound)
		}))

	wmi = New()
	wmi.URL = ts.URL
	return wmi, ts.Close
}
