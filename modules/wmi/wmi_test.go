// SPDX-License-Identifier: GPL-3.0-or-later

package wmi

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/netdata/go.d.plugin/pkg/web"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	v0200Metrics, _ = os.ReadFile("testdata/v0.20.0/metrics.txt")
)

func Test_TestData(t *testing.T) {
	for name, data := range map[string][]byte{
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
		"success on valid response v0.20.0": {
			prepare: prepareWMIv0200,
			wantCollected: map[string]int64{
				"collector_cpu_duration":                                        0,
				"collector_cpu_status_fail":                                     0,
				"collector_cpu_status_success":                                  1,
				"collector_iis_duration":                                        0,
				"collector_iis_status_fail":                                     0,
				"collector_iis_status_success":                                  1,
				"collector_logical_disk_duration":                               0,
				"collector_logical_disk_status_fail":                            0,
				"collector_logical_disk_status_success":                         1,
				"collector_logon_duration":                                      113,
				"collector_logon_status_fail":                                   0,
				"collector_logon_status_success":                                1,
				"collector_memory_duration":                                     0,
				"collector_memory_status_fail":                                  0,
				"collector_memory_status_success":                               1,
				"collector_net_duration":                                        0,
				"collector_net_status_fail":                                     0,
				"collector_net_status_success":                                  1,
				"collector_os_duration":                                         2,
				"collector_os_status_fail":                                      0,
				"collector_os_status_success":                                   1,
				"collector_process_duration":                                    115,
				"collector_process_status_fail":                                 0,
				"collector_process_status_success":                              1,
				"collector_service_duration":                                    101,
				"collector_service_status_fail":                                 0,
				"collector_service_status_success":                              1,
				"collector_system_duration":                                     0,
				"collector_system_status_fail":                                  0,
				"collector_system_status_success":                               1,
				"collector_tcp_duration":                                        0,
				"collector_tcp_status_fail":                                     0,
				"collector_tcp_status_success":                                  1,
				"cpu_core_0,0_cstate_c1":                                        160233427,
				"cpu_core_0,0_cstate_c2":                                        0,
				"cpu_core_0,0_cstate_c3":                                        0,
				"cpu_core_0,0_dpc_time":                                         67109,
				"cpu_core_0,0_dpcs":                                             4871900,
				"cpu_core_0,0_idle_time":                                        162455593,
				"cpu_core_0,0_interrupt_time":                                   77281,
				"cpu_core_0,0_interrupts":                                       155194331,
				"cpu_core_0,0_privileged_time":                                  1182109,
				"cpu_core_0,0_user_time":                                        1073671,
				"cpu_core_0,1_cstate_c1":                                        159528054,
				"cpu_core_0,1_cstate_c2":                                        0,
				"cpu_core_0,1_cstate_c3":                                        0,
				"cpu_core_0,1_dpc_time":                                         11093,
				"cpu_core_0,1_dpcs":                                             1650552,
				"cpu_core_0,1_idle_time":                                        159478125,
				"cpu_core_0,1_interrupt_time":                                   58093,
				"cpu_core_0,1_interrupts":                                       79325847,
				"cpu_core_0,1_privileged_time":                                  1801234,
				"cpu_core_0,1_user_time":                                        3432000,
				"cpu_core_0,2_cstate_c1":                                        159891723,
				"cpu_core_0,2_cstate_c2":                                        0,
				"cpu_core_0,2_cstate_c3":                                        0,
				"cpu_core_0,2_dpc_time":                                         16062,
				"cpu_core_0,2_dpcs":                                             2236469,
				"cpu_core_0,2_idle_time":                                        159848437,
				"cpu_core_0,2_interrupt_time":                                   53515,
				"cpu_core_0,2_interrupts":                                       67305419,
				"cpu_core_0,2_privileged_time":                                  1812546,
				"cpu_core_0,2_user_time":                                        3050250,
				"cpu_core_0,3_cstate_c1":                                        159544117,
				"cpu_core_0,3_cstate_c2":                                        0,
				"cpu_core_0,3_cstate_c3":                                        0,
				"cpu_core_0,3_dpc_time":                                         8140,
				"cpu_core_0,3_dpcs":                                             1185046,
				"cpu_core_0,3_idle_time":                                        159527546,
				"cpu_core_0,3_interrupt_time":                                   44484,
				"cpu_core_0,3_interrupts":                                       60766938,
				"cpu_core_0,3_privileged_time":                                  1760828,
				"cpu_core_0,3_user_time":                                        3422875,
				"cpu_dpc_time":                                                  102404,
				"cpu_idle_time":                                                 641309701,
				"cpu_interrupt_time":                                            233373,
				"cpu_privileged_time":                                           6556717,
				"cpu_user_time":                                                 10978796,
				"iis_Default_Web_Site_iis_active_conn":                          0,
				"iis_Default_Web_Site_iis_active_request_anon":                  0,
				"iis_Default_Web_Site_iis_active_request_non_anon":              0,
				"iis_Default_Web_Site_iis_bandwidth_recv":                       10289,
				"iis_Default_Web_Site_iis_bandwidth_sent":                       105882,
				"iis_Default_Web_Site_iis_conns_atemp":                          1,
				"iis_Default_Web_Site_iis_error_404":                            1,
				"iis_Default_Web_Site_iis_error_423":                            0,
				"iis_Default_Web_Site_iis_extension_req":                        0,
				"iis_Default_Web_Site_iis_file_transfer_recv":                   0,
				"iis_Default_Web_Site_iis_file_transfer_sent":                   2,
				"iis_Default_Web_Site_iis_isapi_ext":                            0,
				"iis_Default_Web_Site_iis_logon":                                4,
				"iis_Default_Web_Site_iis_req_total":                            3,
				"iis_Default_Web_Site_iis_total_req":                            3,
				"iis_Default_Web_Site_iis_total_req_anon":                       3,
				"iis_Default_Web_Site_iis_uptime":                               258633,
				"logical_disk_C:_free_space":                                    43636490240,
				"logical_disk_C:_read_bytes_total":                              17676328448,
				"logical_disk_C:_read_latency":                                  97420,
				"logical_disk_C:_reads_total":                                   350593,
				"logical_disk_C:_total_space":                                   67938287616,
				"logical_disk_C:_used_space":                                    24301797376,
				"logical_disk_C:_write_bytes_total":                             9135282688,
				"logical_disk_C:_write_latency":                                 123912,
				"logical_disk_C:_writes_total":                                  450705,
				"logon_type_batch_sessions":                                     0,
				"logon_type_cached_interactive_sessions":                        0,
				"logon_type_cached_remote_interactive_sessions":                 0,
				"logon_type_cached_unlock_sessions":                             0,
				"logon_type_interactive_sessions":                               2,
				"logon_type_network_clear_text_sessions":                        0,
				"logon_type_network_sessions":                                   0,
				"logon_type_new_credentials_sessions":                           0,
				"logon_type_proxy_sessions":                                     0,
				"logon_type_remote_interactive_sessions":                        0,
				"logon_type_service_sessions":                                   0,
				"logon_type_system_sessions":                                    0,
				"logon_type_unlock_sessions":                                    0,
				"memory_available_bytes":                                        1379942400,
				"memory_cache_faults_total":                                     8009603,
				"memory_cache_total":                                            1392185344,
				"memory_commit_limit":                                           5733113856,
				"memory_committed_bytes":                                        3447439360,
				"memory_modified_page_list_bytes":                               32653312,
				"memory_not_committed_bytes":                                    2285674496,
				"memory_page_faults_total":                                      119093924,
				"memory_pool_nonpaged_bytes_total":                              126865408,
				"memory_pool_paged_bytes":                                       303906816,
				"memory_standby_cache_core_bytes":                               107376640,
				"memory_standby_cache_normal_priority_bytes":                    1019121664,
				"memory_standby_cache_reserve_bytes":                            233033728,
				"memory_standby_cache_total":                                    1359532032,
				"memory_swap_page_reads_total":                                  402087,
				"memory_swap_page_writes_total":                                 7012,
				"memory_swap_pages_read_total":                                  4643279,
				"memory_swap_pages_written_total":                               312896,
				"memory_used_bytes":                                             2876776448,
				"net_nic_Intel_R_PRO_1000_MT_Network_Connection_bytes_received": 38290755856,
				"net_nic_Intel_R_PRO_1000_MT_Network_Connection_bytes_sent":     8211165504,
				"net_nic_Intel_R_PRO_1000_MT_Network_Connection_packets_outbound_discarded": 0,
				"net_nic_Intel_R_PRO_1000_MT_Network_Connection_packets_outbound_errors":    0,
				"net_nic_Intel_R_PRO_1000_MT_Network_Connection_packets_received_discarded": 0,
				"net_nic_Intel_R_PRO_1000_MT_Network_Connection_packets_received_errors":    0,
				"net_nic_Intel_R_PRO_1000_MT_Network_Connection_packets_received_total":     4120869,
				"net_nic_Intel_R_PRO_1000_MT_Network_Connection_packets_sent_total":         1332466,
				"os_paging_free_bytes":                     1414107136,
				"os_paging_limit_bytes":                    1476395008,
				"os_paging_used_bytes":                     62287872,
				"os_physical_memory_free_bytes":            1379946496,
				"os_processes":                             152,
				"os_processes_limit":                       4294967295,
				"os_users":                                 2,
				"os_visible_memory_bytes":                  4256718848,
				"os_visible_memory_used_bytes":             2876772352,
				"process_msedge_cpu_time":                  1919893,
				"process_msedge_handles":                   5779,
				"process_msedge_io_bytes":                  3978227378,
				"process_msedge_io_operations":             16738642,
				"process_msedge_page_faults":               5355941,
				"process_msedge_page_file_bytes":           681603072,
				"process_msedge_threads":                   213,
				"process_msedge_working_set_private_bytes": 461344768,
				"service_dhcp_state_continue_pending":      0,
				"service_dhcp_state_pause_pending":         0,
				"service_dhcp_state_paused":                0,
				"service_dhcp_state_running":               1,
				"service_dhcp_state_start_pending":         0,
				"service_dhcp_state_stop_pending":          0,
				"service_dhcp_state_stopped":               0,
				"service_dhcp_state_unknown":               0,
				"service_dhcp_status_degraded":             0,
				"service_dhcp_status_error":                0,
				"service_dhcp_status_lost_comm":            0,
				"service_dhcp_status_no_contact":           0,
				"service_dhcp_status_nonrecover":           0,
				"service_dhcp_status_ok":                   1,
				"service_dhcp_status_pred_fail":            0,
				"service_dhcp_status_service":              0,
				"service_dhcp_status_starting":             0,
				"service_dhcp_status_stopping":             0,
				"service_dhcp_status_stressed":             0,
				"service_dhcp_status_unknown":              0,
				"system_threads":                           1559,
				"system_up_time":                           221423,
				"tcp_ipv4_conns_active":                    4301,
				"tcp_ipv4_conns_established":               7,
				"tcp_ipv4_conns_failures":                  137,
				"tcp_ipv4_conns_passive":                   501,
				"tcp_ipv4_conns_resets":                    1282,
				"tcp_ipv4_segments_received":               676388,
				"tcp_ipv4_segments_retransmitted":          2120,
				"tcp_ipv4_segments_sent":                   871379,
				"tcp_ipv6_conns_active":                    214,
				"tcp_ipv6_conns_established":               0,
				"tcp_ipv6_conns_failures":                  214,
				"tcp_ipv6_conns_passive":                   0,
				"tcp_ipv6_conns_resets":                    0,
				"tcp_ipv6_segments_received":               1284,
				"tcp_ipv6_segments_retransmitted":          428,
				"tcp_ipv6_segments_sent":                   856,
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

			mx := wmi.Collect()

			if mx != nil && test.wantCollected != nil {
				mx["system_up_time"] = test.wantCollected["system_up_time"]
			}

			assert.Equal(t, test.wantCollected, mx)
			if len(test.wantCollected) > 0 {
				testCharts(t, wmi, mx)
			}
		})
	}
}

func testCharts(t *testing.T, wmi *WMI, mx map[string]int64) {
	ensureChartsDimsCreated(t, wmi)
	ensureCollectedHasAllChartsDimsVarsIDs(t, wmi, mx)
}

func ensureChartsDimsCreated(t *testing.T, w *WMI) {
	for _, chart := range cpuCharts {
		if w.cache.collection[collectorCPU] {
			assert.Truef(t, w.Charts().Has(chart.ID), "chart '%s' not created", chart.ID)
		} else {
			assert.Falsef(t, w.Charts().Has(chart.ID), "chart '%s' created", chart.ID)
		}
	}
	for _, chart := range memCharts {
		if w.cache.collection[collectorMemory] {
			assert.Truef(t, w.Charts().Has(chart.ID), "chart '%s' not created", chart.ID)
		} else {
			assert.Falsef(t, w.Charts().Has(chart.ID), "chart '%s' created", chart.ID)
		}
	}
	for _, chart := range tcpCharts {
		if w.cache.collection[collectorTCP] {
			assert.Truef(t, w.Charts().Has(chart.ID), "chart '%s' not created", chart.ID)
		} else {
			assert.Falsef(t, w.Charts().Has(chart.ID), "chart '%s' created", chart.ID)
		}
	}
	for _, chart := range osCharts {
		if w.cache.collection[collectorOS] {
			assert.Truef(t, w.Charts().Has(chart.ID), "chart '%s' not created", chart.ID)
		} else {
			assert.Falsef(t, w.Charts().Has(chart.ID), "chart '%s' created", chart.ID)
		}
	}
	for _, chart := range systemCharts {
		if w.cache.collection[collectorSystem] {
			assert.Truef(t, w.Charts().Has(chart.ID), "chart '%s' not created", chart.ID)
		} else {
			assert.Falsef(t, w.Charts().Has(chart.ID), "chart '%s' created", chart.ID)
		}
	}
	for _, chart := range logonCharts {
		if w.cache.collection[collectorLogon] {
			assert.Truef(t, w.Charts().Has(chart.ID), "chart '%s' not created", chart.ID)
		} else {
			assert.Falsef(t, w.Charts().Has(chart.ID), "chart '%s' created", chart.ID)
		}
	}
	for _, chart := range processesCharts {
		if w.cache.collection[collectorProcess] {
			assert.Truef(t, w.Charts().Has(chart.ID), "chart '%s' not created", chart.ID)
		} else {
			assert.Falsef(t, w.Charts().Has(chart.ID), "chart '%s' created", chart.ID)
		}
	}

	for core := range w.cache.cores {
		for _, chart := range cpuCoreChartsTmpl {
			id := fmt.Sprintf(chart.ID, core)
			assert.Truef(t, w.Charts().Has(id), "charts has no '%s' chart for '%s' core", id, core)
		}
	}
	for disk := range w.cache.volumes {
		for _, chart := range diskChartsTmpl {
			id := fmt.Sprintf(chart.ID, disk)
			assert.Truef(t, w.Charts().Has(id), "charts has no '%s' chart for '%s' disk", id, disk)
		}
	}
	for nic := range w.cache.nics {
		for _, chart := range nicChartsTmpl {
			id := fmt.Sprintf(chart.ID, nic)
			assert.Truef(t, w.Charts().Has(id), "charts has no '%s' chart for '%s' nic", id, nic)
		}
	}
	for zone := range w.cache.thermalZones {
		for _, chart := range thermalzoneChartsTmpl {
			id := fmt.Sprintf(chart.ID, zone)
			assert.Truef(t, w.Charts().Has(id), "charts has no '%s' chart for '%s' thermalzone", id, zone)
		}
	}
	for svc := range w.cache.services {
		for _, chart := range serviceChartsTmpl {
			id := fmt.Sprintf(chart.ID, svc)
			assert.Truef(t, w.Charts().Has(id), "charts has no '%s' chart for '%s' service", id, svc)
		}
	}
	for name := range w.cache.collectors {
		for _, chart := range collectorChartsTmpl {
			id := fmt.Sprintf(chart.ID, name)
			assert.Truef(t, w.Charts().Has(id), "charts has no '%s' chart for '%s' collector", id, name)
		}
	}

	for _, chart := range processesCharts {
		if chart = w.Charts().Get(chart.ID); chart == nil {
			continue
		}
		for proc := range w.cache.processes {
			var found bool
			for _, dim := range chart.Dims {
				if found = strings.HasPrefix(dim.ID, "process_"+proc); found {
					break
				}
			}
			assert.Truef(t, found, "chart '%s' has not dim for '%s' process", chart.ID, proc)
		}
	}
}

func ensureCollectedHasAllChartsDimsVarsIDs(t *testing.T, w *WMI, mx map[string]int64) {
	for _, chart := range *w.Charts() {
		for _, dim := range chart.Dims {
			_, ok := mx[dim.ID]
			assert.Truef(t, ok, "collected metrics has no data for dim '%s' chart '%s'", dim.ID, chart.ID)
		}
		for _, v := range chart.Vars {
			_, ok := mx[v.ID]
			assert.Truef(t, ok, "collected metrics has no data for var '%s' chart '%s'", v.ID, chart.ID)
		}
	}
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
