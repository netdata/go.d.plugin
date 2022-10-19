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
	v0200Metrics, _ = os.ReadFile("testdata/v0.20.0/metrics.txt")
)

func Test_TestData(t *testing.T) {
	for name, data := range map[string][]byte{
		"v0150Metrics": v0150Metrics,
		"v0160Metrics": v0160Metrics,
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
		"success on valid response v0.16.0": {
			prepare: prepareWMIv0160,
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
		"success on valid response v0.20.0": {
			prepare: prepareWMIv0200,
			wantCollected: map[string]int64{
				"cpu_collection_duration":                                  2,
				"cpu_collection_success":                                   1,
				"cpu_core_0,0_c1":                                          83364,
				"cpu_core_0,0_c2":                                          54090,
				"cpu_core_0,0_c3":                                          364706,
				"cpu_core_0,0_dpc":                                         9406,
				"cpu_core_0,0_dpcs":                                        383958000,
				"cpu_core_0,0_idle":                                        687843,
				"cpu_core_0,0_interrupt":                                   7390,
				"cpu_core_0,0_interrupts":                                  1789909000,
				"cpu_core_0,0_privileged":                                  92015,
				"cpu_core_0,0_user":                                        91406,
				"cpu_core_0,1_c1":                                          67504,
				"cpu_core_0,1_c2":                                          49498,
				"cpu_core_0,1_c3":                                          415924,
				"cpu_core_0,1_dpc":                                         1375,
				"cpu_core_0,1_dpcs":                                        121705000,
				"cpu_core_0,1_idle":                                        704140,
				"cpu_core_0,1_interrupt":                                   2125,
				"cpu_core_0,1_interrupts":                                  1254593000,
				"cpu_core_0,1_privileged":                                  67828,
				"cpu_core_0,1_user":                                        98578,
				"cpu_core_0,2_c1":                                          82466,
				"cpu_core_0,2_c2":                                          47310,
				"cpu_core_0,2_c3":                                          391016,
				"cpu_core_0,2_dpc":                                         781,
				"cpu_core_0,2_dpcs":                                        101032000,
				"cpu_core_0,2_idle":                                        729937,
				"cpu_core_0,2_interrupt":                                   1343,
				"cpu_core_0,2_interrupts":                                  1159698000,
				"cpu_core_0,2_privileged":                                  54609,
				"cpu_core_0,2_user":                                        86000,
				"cpu_core_0,3_c1":                                          65836,
				"cpu_core_0,3_c2":                                          47179,
				"cpu_core_0,3_c3":                                          431330,
				"cpu_core_0,3_dpc":                                         468,
				"cpu_core_0,3_dpcs":                                        84391000,
				"cpu_core_0,3_idle":                                        731906,
				"cpu_core_0,3_interrupt":                                   984,
				"cpu_core_0,3_interrupts":                                  886968000,
				"cpu_core_0,3_privileged":                                  51250,
				"cpu_core_0,3_user":                                        87390,
				"cpu_core_0,4_c1":                                          64077,
				"cpu_core_0,4_c2":                                          30752,
				"cpu_core_0,4_c3":                                          491402,
				"cpu_core_0,4_dpc":                                         781,
				"cpu_core_0,4_dpcs":                                        113844000,
				"cpu_core_0,4_idle":                                        763265,
				"cpu_core_0,4_interrupt":                                   1562,
				"cpu_core_0,4_interrupts":                                  1254399000,
				"cpu_core_0,4_privileged":                                  51218,
				"cpu_core_0,4_user":                                        56062,
				"cpu_core_0,5_c1":                                          59385,
				"cpu_core_0,5_c2":                                          30267,
				"cpu_core_0,5_c3":                                          521170,
				"cpu_core_0,5_dpc":                                         593,
				"cpu_core_0,5_dpcs":                                        82929000,
				"cpu_core_0,5_idle":                                        769093,
				"cpu_core_0,5_interrupt":                                   1281,
				"cpu_core_0,5_interrupts":                                  997639000,
				"cpu_core_0,5_privileged":                                  46796,
				"cpu_core_0,5_user":                                        54656,
				"cpu_core_0,6_c1":                                          44442,
				"cpu_core_0,6_c2":                                          15357,
				"cpu_core_0,6_c3":                                          579625,
				"cpu_core_0,6_dpc":                                         1125,
				"cpu_core_0,6_dpcs":                                        98283000,
				"cpu_core_0,6_idle":                                        774671,
				"cpu_core_0,6_interrupt":                                   1187,
				"cpu_core_0,6_interrupts":                                  829347000,
				"cpu_core_0,6_privileged":                                  41312,
				"cpu_core_0,6_user":                                        54562,
				"cpu_core_0,7_c1":                                          41024,
				"cpu_core_0,7_c2":                                          14045,
				"cpu_core_0,7_c3":                                          597029,
				"cpu_core_0,7_dpc":                                         468,
				"cpu_core_0,7_dpcs":                                        53923000,
				"cpu_core_0,7_idle":                                        772703,
				"cpu_core_0,7_interrupt":                                   781,
				"cpu_core_0,7_interrupts":                                  670161000,
				"cpu_core_0,7_privileged":                                  38703,
				"cpu_core_0,7_user":                                        59125,
				"cpu_dpc":                                                  15000,
				"cpu_idle":                                                 5933562,
				"cpu_interrupt":                                            16656,
				"cpu_privileged":                                           443734,
				"cpu_user":                                                 587781,
				"iis_collection_duration":                                  2,
				"iis_collection_success":                                   1,
				"logical_disk_C:_free_space":                               11994660864000,
				"logical_disk_C:_idle_seconds_total":                       0,
				"logical_disk_C:_read_bytes_total":                         8356269056000,
				"logical_disk_C:_read_latency":                             176005,
				"logical_disk_C:_read_seconds_total":                       0,
				"logical_disk_C:_reads_total":                              153413,
				"logical_disk_C:_requests_queued":                          0,
				"logical_disk_C:_split_ios_total":                          0,
				"logical_disk_C:_total_space":                              238910701568000,
				"logical_disk_C:_used_space":                               226916040704000,
				"logical_disk_C:_write_bytes_total":                        2213975040000,
				"logical_disk_C:_write_latency":                            35557,
				"logical_disk_C:_write_seconds_total":                      0,
				"logical_disk_C:_writes_total":                             44137,
				"logical_disk_collection_duration":                         0,
				"logical_disk_collection_success":                          1,
				"logon_collection_duration":                                200,
				"logon_collection_success":                                 1,
				"logon_type_batch":                                         0,
				"logon_type_cached_interactive":                            0,
				"logon_type_cached_remote_interactive":                     0,
				"logon_type_cached_unlock":                                 0,
				"logon_type_interactive":                                   6,
				"logon_type_network":                                       0,
				"logon_type_network_clear_text":                            0,
				"logon_type_new_credentials":                               0,
				"logon_type_proxy":                                         0,
				"logon_type_remote_interactive":                            0,
				"logon_type_service":                                       4,
				"logon_type_system":                                        1,
				"logon_type_unlock":                                        0,
				"memory_available_bytes":                                   1021739008000,
				"memory_cache_bytes":                                       123760640000,
				"memory_cache_bytes_peak":                                  255131648000,
				"memory_cache_faults_total":                                2451301000,
				"memory_cache_total":                                       1053057024000,
				"memory_collection_duration":                               0,
				"memory_collection_success":                                1,
				"memory_commit_limit":                                      16599695360000,
				"memory_committed_bytes":                                   10447548416000,
				"memory_demand_zero_faults_total":                          11553396000,
				"memory_free_and_zero_page_list_bytes":                     2314240000,
				"memory_free_system_page_table_entries":                    16431393000,
				"memory_modified_page_list_bytes":                          33632256000,
				"memory_not_committed_bytes":                               6152146944000,
				"memory_page_faults_total":                                 19728764000,
				"memory_pool_nonpaged_allocs_total":                        0,
				"memory_pool_nonpaged_bytes_total":                         0,
				"memory_pool_paged_allocs_total":                           0,
				"memory_pool_paged_bytes":                                  310181888000,
				"memory_pool_paged_resident_bytes":                         236302336000,
				"memory_standby_cache_core_bytes":                          141012992000,
				"memory_standby_cache_normal_priority_bytes":               877178880000,
				"memory_standby_cache_reserve_bytes":                       1232896000,
				"memory_standby_cache_total":                               1019424768000,
				"memory_swap_page_operations_total":                        2769288000,
				"memory_swap_page_reads_total":                             250944000,
				"memory_swap_page_writes_total":                            1213000,
				"memory_swap_pages_read_total":                             2475232000,
				"memory_swap_pages_written_total":                          294056000,
				"memory_system_cache_resident_bytes":                       123760640000,
				"memory_system_code_resident_bytes":                        17092608000,
				"memory_system_code_total_bytes":                           8192000,
				"memory_system_driver_resident_bytes":                      34639872000,
				"memory_system_driver_total_bytes":                         18665472000,
				"memory_transition_faults_total":                           5924756000,
				"memory_transition_pages_repurposed_total":                 2139419000,
				"memory_used_bytes":                                        7278108672000,
				"memory_write_copies_total":                                151687000,
				"net_Intel_R_Wireless_AC_9462_bytes_received":              129484170000,
				"net_Intel_R_Wireless_AC_9462_bytes_sent":                  4647982000,
				"net_Intel_R_Wireless_AC_9462_bytes_total":                 134132152000,
				"net_Intel_R_Wireless_AC_9462_current_bandwidth":           0,
				"net_Intel_R_Wireless_AC_9462_packets_outbound_discarded":  0,
				"net_Intel_R_Wireless_AC_9462_packets_outbound_errors":     0,
				"net_Intel_R_Wireless_AC_9462_packets_received_discarded":  0,
				"net_Intel_R_Wireless_AC_9462_packets_received_errors":     0,
				"net_Intel_R_Wireless_AC_9462_packets_received_total":      98374000,
				"net_Intel_R_Wireless_AC_9462_packets_received_unknown":    0,
				"net_Intel_R_Wireless_AC_9462_packets_sent_total":          17817000,
				"net_Intel_R_Wireless_AC_9462_packets_total":               116191000,
				"net_Realtek_PCIe_GbE_Family_Controller_bytes_received":    0,
				"net_Realtek_PCIe_GbE_Family_Controller_bytes_sent":        0,
				"net_Realtek_PCIe_GbE_Family_Controller_bytes_total":       0,
				"net_Realtek_PCIe_GbE_Family_Controller_current_bandwidth": 0,
				"net_Realtek_PCIe_GbE_Family_Controller_packets_outbound_discarded": 0,
				"net_Realtek_PCIe_GbE_Family_Controller_packets_outbound_errors":    0,
				"net_Realtek_PCIe_GbE_Family_Controller_packets_received_discarded": 0,
				"net_Realtek_PCIe_GbE_Family_Controller_packets_received_errors":    0,
				"net_Realtek_PCIe_GbE_Family_Controller_packets_received_total":     0,
				"net_Realtek_PCIe_GbE_Family_Controller_packets_received_unknown":   0,
				"net_Realtek_PCIe_GbE_Family_Controller_packets_sent_total":         0,
				"net_Realtek_PCIe_GbE_Family_Controller_packets_total":              0,
				"net_collection_duration":           0,
				"net_collection_success":            1,
				"os_collection_duration":            20,
				"os_collection_success":             1,
				"os_paging_free_bytes":              7307608064000,
				"os_paging_limit_bytes":             8299847680000,
				"os_paging_used_bytes":              992239616000,
				"os_physical_memory_free_bytes":     1022656512000,
				"os_process_memory_limit_bytes":     140737488224256000,
				"os_processes":                      257,
				"os_processes_limit":                4294967295,
				"os_time":                           1666132477,
				"os_users":                          2,
				"os_virtual_memory_bytes":           16599695360000,
				"os_virtual_memory_free_bytes":      6151704576000,
				"os_visible_memory_bytes":           8299847680000,
				"os_visible_memory_used_bytes":      7277191168000,
				"system_boot_time":                  1666131606,
				"system_calls_total":                73154314000,
				"system_collection_duration":        0,
				"system_collection_success":         1,
				"system_context_switches_total":     17853582000,
				"system_exception_dispatches_total": 19751000,
				"system_processor_queue_length":     0,
				"system_threads":                    3128,
				"system_up_time":                    0,
				"thermalzone_TZ00_temperature":      44050,
				"thermalzone_collection_duration":   306,
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
