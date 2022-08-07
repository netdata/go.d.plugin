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
	v0160Metrics, _ = os.ReadFile("testdata/v0.16.0/metrics.txt")
)

func Test_TestData(t *testing.T) {
	for name, data := range map[string][]byte{
		"v0150Metrics": v0150Metrics,
		"v0160Metrics": v0160Metrics,
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
		"success on valid response v0.16.0": {
			prepare: prepareWMIv0160,
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
		"success on valid response v0.16.0": {
			prepare: prepareWMIv0160,
			wantCollected: map[string]int64{
				"cpu_collection_duration":                                     1000,
				"cpu_collection_success":                                      1,
				"cpu_core_0,0_c1":                                             684265,
				"cpu_core_0,0_c2":                                             1000,
				"cpu_core_0,0_c3":                                             1000,
				"cpu_core_0,0_dpc":                                            6890,
				"cpu_core_0,0_dpcs":                                           367230000,
				"cpu_core_0,0_idle":                                           715406,
				"cpu_core_0,0_interrupt":                                      21828,
				"cpu_core_0,0_interrupts":                                     3928344000,
				"cpu_core_0,0_privileged":                                     530593,
				"cpu_core_0,0_user":                                           424578,
				"cpu_dpc":                                                     6890,
				"cpu_idle":                                                    715406,
				"cpu_interrupt":                                               21828,
				"cpu_privileged":                                              530593,
				"cpu_user":                                                    424578,
				"logical_disk_C:_free_space":                                  8392802304000,
				"logical_disk_C:_idle_seconds_total":                          0,
				"logical_disk_C:_read_bytes_total":                            8469489664000,
				"logical_disk_C:_read_latency":                                144007,
				"logical_disk_C:_read_seconds_total":                          0,
				"logical_disk_C:_reads_total":                                 101190,
				"logical_disk_C:_requests_queued":                             0,
				"logical_disk_C:_split_ios_total":                             0,
				"logical_disk_C:_total_space":                                 21371027456000,
				"logical_disk_C:_used_space":                                  12978225152000,
				"logical_disk_C:_write_bytes_total":                           7485627392000,
				"logical_disk_C:_write_latency":                               39909,
				"logical_disk_C:_write_seconds_total":                         0,
				"logical_disk_C:_writes_total":                                56687,
				"logical_disk_collection_duration":                            1000,
				"logical_disk_collection_success":                             1,
				"logon_collection_duration":                                   1044,
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
				"memory_available_bytes":                                      809349120000,
				"memory_cache_bytes":                                          68198400000,
				"memory_cache_bytes_peak":                                     102326272000,
				"memory_cache_faults_total":                                   915945000,
				"memory_cache_total":                                          838234112000,
				"memory_collection_duration":                                  1000,
				"memory_collection_success":                                   1,
				"memory_commit_limit":                                         3547709440000,
				"memory_committed_bytes":                                      2642300928000,
				"memory_demand_zero_faults_total":                             6452835000,
				"memory_free_and_zero_page_list_bytes":                        3993600000,
				"memory_free_system_page_table_entries":                       12530059000,
				"memory_modified_page_list_bytes":                             32878592000,
				"memory_not_committed_bytes":                                  905408512000,
				"memory_page_faults_total":                                    19064737000,
				"memory_pool_nonpaged_allocs_total":                           1000,
				"memory_pool_nonpaged_bytes_total":                            97280000000,
				"memory_pool_paged_allocs_total":                              1000,
				"memory_pool_paged_bytes":                                     172818432000,
				"memory_pool_paged_resident_bytes":                            153276416000,
				"memory_standby_cache_core_bytes":                             124502016000,
				"memory_standby_cache_normal_priority_bytes":                  485228544000,
				"memory_standby_cache_reserve_bytes":                          195624960000,
				"memory_standby_cache_total":                                  805355520000,
				"memory_swap_page_operations_total":                           1472952000,
				"memory_swap_page_reads_total":                                128201000,
				"memory_swap_page_writes_total":                               3752000,
				"memory_swap_pages_read_total":                                1242858000,
				"memory_swap_pages_written_total":                             230094000,
				"memory_system_cache_resident_bytes":                          68198400000,
				"memory_system_code_resident_bytes":                           4321280000,
				"memory_system_code_total_bytes":                              4636672000,
				"memory_system_driver_resident_bytes":                         3309568000,
				"memory_system_driver_total_bytes":                            17526784000,
				"memory_transition_faults_total":                              11958820000,
				"memory_transition_pages_repurposed_total":                    1381522000,
				"memory_used_bytes":                                           1337663488000,
				"memory_write_copies_total":                                   106130000,
				"net_Intel_R_PRO_1000_MT_Desktop_Adapter_2_bytes_received":    83866000,
				"net_Intel_R_PRO_1000_MT_Desktop_Adapter_2_bytes_sent":        110493000,
				"net_Intel_R_PRO_1000_MT_Desktop_Adapter_2_bytes_total":       194359000,
				"net_Intel_R_PRO_1000_MT_Desktop_Adapter_2_current_bandwidth": 1000000000,
				"net_Intel_R_PRO_1000_MT_Desktop_Adapter_2_packets_outbound_discarded": 1000,
				"net_Intel_R_PRO_1000_MT_Desktop_Adapter_2_packets_outbound_errors":    1000,
				"net_Intel_R_PRO_1000_MT_Desktop_Adapter_2_packets_received_discarded": 1000,
				"net_Intel_R_PRO_1000_MT_Desktop_Adapter_2_packets_received_errors":    1000,
				"net_Intel_R_PRO_1000_MT_Desktop_Adapter_2_packets_received_total":     736000,
				"net_Intel_R_PRO_1000_MT_Desktop_Adapter_2_packets_received_unknown":   1000,
				"net_Intel_R_PRO_1000_MT_Desktop_Adapter_2_packets_sent_total":         730000,
				"net_Intel_R_PRO_1000_MT_Desktop_Adapter_2_packets_total":              1466000,
				"net_Intel_R_PRO_1000_MT_Desktop_Adapter_bytes_received":               424915241000,
				"net_Intel_R_PRO_1000_MT_Desktop_Adapter_bytes_sent":                   7656073000,
				"net_Intel_R_PRO_1000_MT_Desktop_Adapter_bytes_total":                  432571314000,
				"net_Intel_R_PRO_1000_MT_Desktop_Adapter_current_bandwidth":            1000000000,
				"net_Intel_R_PRO_1000_MT_Desktop_Adapter_packets_outbound_discarded":   1000,
				"net_Intel_R_PRO_1000_MT_Desktop_Adapter_packets_outbound_errors":      1000,
				"net_Intel_R_PRO_1000_MT_Desktop_Adapter_packets_received_discarded":   1000,
				"net_Intel_R_PRO_1000_MT_Desktop_Adapter_packets_received_errors":      1000,
				"net_Intel_R_PRO_1000_MT_Desktop_Adapter_packets_received_total":       290889000,
				"net_Intel_R_PRO_1000_MT_Desktop_Adapter_packets_received_unknown":     1000,
				"net_Intel_R_PRO_1000_MT_Desktop_Adapter_packets_sent_total":           96184000,
				"net_Intel_R_PRO_1000_MT_Desktop_Adapter_packets_total":                387073000,
				"net_collection_duration":           1000,
				"net_collection_success":            1,
				"os_collection_duration":            1046,
				"os_collection_success":             1,
				"os_paging_free_bytes":              1042104320000,
				"os_paging_limit_bytes":             1400696832000,
				"os_paging_used_bytes":              358592512000,
				"os_physical_memory_free_bytes":     810295296000,
				"os_process_memory_limit_bytes":     0,
				"os_processes":                      78,
				"os_processes_limit":                4294967295,
				"os_time":                           1615228173,
				"os_users":                          2,
				"os_virtual_memory_bytes":           3547709440000,
				"os_virtual_memory_free_bytes":      904642560000,
				"os_visible_memory_bytes":           2147012608000,
				"os_visible_memory_used_bytes":      1336717312000,
				"system_boot_time":                  1615226502,
				"system_calls_total":                44265421000,
				"system_collection_duration":        1000,
				"system_collection_success":         1,
				"system_context_switches_total":     4685379000,
				"system_exception_dispatches_total": 3704000,
				"system_processor_queue_length":     1,
				"system_threads":                    943,
				"system_up_time":                    10424271,
				"thermalzone_THRM_temperature":      30050,
				"thermalzone_TZ00_temperature":      27850,
				"thermalzone_TZ01_temperature":      29850,
				"thermalzone_collection_duration":   1000,
				"thermalzone_collection_success":    1,
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

func prepareWMIv0160() (wmi *WMI, cleanup func()) {
	ts := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			_, _ = w.Write(v0160Metrics)
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
