package wmi

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/netdata/go.d.plugin/agent/module"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	fullData, _    = ioutil.ReadFile("testdata/full.txt")
	partialData, _ = ioutil.ReadFile("testdata/partial.txt")
)

func Test_readTestData(t *testing.T) {
	assert.NotNil(t, fullData)
	assert.NotNil(t, partialData)
}

func TestNew(t *testing.T) {
	assert.Implements(t, (*module.Module)(nil), New())
}

func TestWMI_Init(t *testing.T) {
	wmi := New()
	wmi.URL = "http://127.0.0.1:38001/metrics"

	assert.True(t, wmi.Init())
}

func TestWMI_Init_ErrorOnValidatingConfigURLIsNotSet(t *testing.T) {
	wmi := New()

	assert.False(t, wmi.Init())
}

func TestWMI_Init_ErrorOnCreatingClientWrongTLSCA(t *testing.T) {
	wmi := New()
	wmi.URL = "http://127.0.0.1:38001/metrics"
	wmi.Client.TLSConfig.TLSCA = "testdata/tls"

	assert.False(t, wmi.Init())
}

func TestWMI_Check(t *testing.T) {
	wmi, ts := prepareClientServerFullData(t)
	defer ts.Close()

	assert.True(t, wmi.Check())
}

func TestWMI_Check_ErrorOnCollectConnectionRefused(t *testing.T) {
	wmi := New()
	wmi.URL = "http://127.0.0.1:38001/metrics"
	require.True(t, wmi.Init())

	assert.False(t, wmi.Check())
}

func TestWMI_Charts(t *testing.T) {
	assert.NotNil(t, New().Charts())
}

func TestWMI_Cleanup(t *testing.T) {
	assert.NotPanics(t, New().Cleanup)
}

func TestWMI_Collect(t *testing.T) {
	wmi, ts := prepareClientServerFullData(t)
	defer ts.Close()

	expected := map[string]int64{
		"cpu_collection_duration":                                     0,
		"cpu_collection_success":                                      1,
		"cpu_core_0,0_c1":                                             28905096,
		"cpu_core_0,0_c2":                                             0,
		"cpu_core_0,0_c3":                                             0,
		"cpu_core_0,0_dpc":                                            3828,
		"cpu_core_0,0_dpcs":                                           305034000,
		"cpu_core_0,0_idle":                                           28535546,
		"cpu_core_0,0_interrupt":                                      22734,
		"cpu_core_0,0_interrupts":                                     59163056000,
		"cpu_core_0,0_privileged":                                     2113093,
		"cpu_core_0,0_user":                                           860296,
		"cpu_core_0,1_c1":                                             29243332,
		"cpu_core_0,1_c2":                                             0,
		"cpu_core_0,1_c3":                                             0,
		"cpu_core_0,1_dpc":                                            1562,
		"cpu_core_0,1_dpcs":                                           170434000,
		"cpu_core_0,1_idle":                                           29413609,
		"cpu_core_0,1_interrupt":                                      12171,
		"cpu_core_0,1_interrupts":                                     10966489000,
		"cpu_core_0,1_privileged":                                     943984,
		"cpu_core_0,1_user":                                           1151109,
		"cpu_core_0,2_c1":                                             29914283,
		"cpu_core_0,2_c2":                                             0,
		"cpu_core_0,2_c3":                                             0,
		"cpu_core_0,2_dpc":                                            3328,
		"cpu_core_0,2_dpcs":                                           162396000,
		"cpu_core_0,2_idle":                                           29943328,
		"cpu_core_0,2_interrupt":                                      16843,
		"cpu_core_0,2_interrupts":                                     11142989000,
		"cpu_core_0,2_privileged":                                     634875,
		"cpu_core_0,2_user":                                           930500,
		"cpu_core_0,3_c1":                                             29465377,
		"cpu_core_0,3_c2":                                             0,
		"cpu_core_0,3_c3":                                             0,
		"cpu_core_0,3_dpc":                                            49390,
		"cpu_core_0,3_dpcs":                                           213509000,
		"cpu_core_0,3_idle":                                           29285500,
		"cpu_core_0,3_interrupt":                                      31078,
		"cpu_core_0,3_interrupts":                                     10967361000,
		"cpu_core_0,3_privileged":                                     1082906,
		"cpu_core_0,3_user":                                           1140296,
		"cpu_dpc":                                                     58109,
		"cpu_idle":                                                    117177984,
		"cpu_interrupt":                                               82828,
		"cpu_privileged":                                              4774859,
		"cpu_user":                                                    4082203,
		"logical_disk_C:_free_space":                                  31390171136000,
		"logical_disk_C:_idle_seconds_total":                          0,
		"logical_disk_C:_read_bytes_total":                            1531812864000,
		"logical_disk_C:_read_latency":                                37023,
		"logical_disk_C:_read_seconds_total":                          0,
		"logical_disk_C:_reads_total":                                 38137,
		"logical_disk_C:_requests_queued":                             0,
		"logical_disk_C:_split_ios_total":                             0,
		"logical_disk_C:_total_space":                                 53076819968000,
		"logical_disk_C:_used_space":                                  21686648832000,
		"logical_disk_C:_write_bytes_total":                           961305600000,
		"logical_disk_C:_write_latency":                               18309,
		"logical_disk_C:_write_seconds_total":                         0,
		"logical_disk_C:_writes_total":                                38039,
		"logical_disk_E:_free_space":                                  10694426624000,
		"logical_disk_E:_idle_seconds_total":                          0,
		"logical_disk_E:_read_bytes_total":                            904704000,
		"logical_disk_E:_read_latency":                                40,
		"logical_disk_E:_read_seconds_total":                          0,
		"logical_disk_E:_reads_total":                                 176,
		"logical_disk_E:_requests_queued":                             0,
		"logical_disk_E:_split_ios_total":                             0,
		"logical_disk_E:_total_space":                                 10733223936000,
		"logical_disk_E:_used_space":                                  38797312000,
		"logical_disk_E:_write_bytes_total":                           32892416000,
		"logical_disk_E:_write_latency":                               84,
		"logical_disk_E:_write_seconds_total":                         0,
		"logical_disk_E:_writes_total":                                294,
		"logical_disk_collection_duration":                            0,
		"logical_disk_collection_success":                             1,
		"logon_collection_duration":                                   115,
		"logon_collection_success":                                    1,
		"logon_type_batch":                                            0,
		"logon_type_cached_interactive":                               0,
		"logon_type_cached_remote_interactive":                        0,
		"logon_type_cached_unlock":                                    0,
		"logon_type_interactive":                                      2,
		"logon_type_network":                                          0,
		"logon_type_network_clear_text":                               0,
		"logon_type_new_credentials":                                  0,
		"logon_type_proxy":                                            0,
		"logon_type_remote_interactive":                               0,
		"logon_type_service":                                          0,
		"logon_type_system":                                           0,
		"logon_type_unlock":                                           0,
		"memory_available_bytes":                                      2621665280000,
		"memory_cache_bytes":                                          55283712000,
		"memory_cache_bytes_peak":                                     81985536000,
		"memory_cache_faults_total":                                   291802000,
		"memory_cache_total":                                          1866829824000,
		"memory_collection_duration":                                  0,
		"memory_collection_success":                                   1,
		"memory_commit_limit":                                         5770891264000,
		"memory_committed_bytes":                                      1653608448000,
		"memory_demand_zero_faults_total":                             26234351000,
		"memory_free_and_zero_page_list_bytes":                        816140288000,
		"memory_free_system_page_table_entries":                       12558385000,
		"memory_modified_page_list_bytes":                             61304832000,
		"memory_not_committed_bytes":                                  4117282816000,
		"memory_page_faults_total":                                    38420407000,
		"memory_pool_nonpaged_allocs_total":                           0,
		"memory_pool_nonpaged_bytes_total":                            74821632000,
		"memory_pool_paged_allocs_total":                              0,
		"memory_pool_paged_bytes":                                     118706176000,
		"memory_pool_paged_resident_bytes":                            99954688000,
		"memory_standby_cache_core_bytes":                             158302208000,
		"memory_standby_cache_normal_priority_bytes":                  310276096000,
		"memory_standby_cache_reserve_bytes":                          1336946688000,
		"memory_standby_cache_total":                                  1805524992000,
		"memory_swap_page_operations_total":                           524900000,
		"memory_swap_page_reads_total":                                73613000,
		"memory_swap_page_writes_total":                               61000,
		"memory_swap_pages_read_total":                                521808000,
		"memory_swap_pages_written_total":                             3092000,
		"memory_system_cache_resident_bytes":                          55283712000,
		"memory_system_code_resident_bytes":                           0,
		"memory_system_code_total_bytes":                              0,
		"memory_system_driver_resident_bytes":                         9621504000,
		"memory_system_driver_total_bytes":                            16408576000,
		"memory_transition_faults_total":                              12759543000,
		"memory_transition_pages_repurposed_total":                    0,
		"memory_used_bytes":                                           1672830976000,
		"memory_write_copies_total":                                   271024000,
		"net_Intel_R_PRO_1000_MT_Desktop_Adapter_2_bytes_received":    10661733000,
		"net_Intel_R_PRO_1000_MT_Desktop_Adapter_2_bytes_sent":        186479806000,
		"net_Intel_R_PRO_1000_MT_Desktop_Adapter_2_bytes_total":       197141539000,
		"net_Intel_R_PRO_1000_MT_Desktop_Adapter_2_current_bandwidth": 1000000000,
		"net_Intel_R_PRO_1000_MT_Desktop_Adapter_2_packets_outbound_discarded": 0,
		"net_Intel_R_PRO_1000_MT_Desktop_Adapter_2_packets_outbound_errors":    0,
		"net_Intel_R_PRO_1000_MT_Desktop_Adapter_2_packets_received_discarded": 0,
		"net_Intel_R_PRO_1000_MT_Desktop_Adapter_2_packets_received_errors":    0,
		"net_Intel_R_PRO_1000_MT_Desktop_Adapter_2_packets_received_total":     93134000,
		"net_Intel_R_PRO_1000_MT_Desktop_Adapter_2_packets_received_unknown":   0,
		"net_Intel_R_PRO_1000_MT_Desktop_Adapter_2_packets_sent_total":         94947000,
		"net_Intel_R_PRO_1000_MT_Desktop_Adapter_2_packets_total":              188081000,
		"net_Intel_R_PRO_1000_MT_Desktop_Adapter_bytes_received":               4541800000,
		"net_Intel_R_PRO_1000_MT_Desktop_Adapter_bytes_sent":                   1780848000,
		"net_Intel_R_PRO_1000_MT_Desktop_Adapter_bytes_total":                  6322648000,
		"net_Intel_R_PRO_1000_MT_Desktop_Adapter_current_bandwidth":            1000000000,
		"net_Intel_R_PRO_1000_MT_Desktop_Adapter_packets_outbound_discarded":   0,
		"net_Intel_R_PRO_1000_MT_Desktop_Adapter_packets_outbound_errors":      0,
		"net_Intel_R_PRO_1000_MT_Desktop_Adapter_packets_received_discarded":   0,
		"net_Intel_R_PRO_1000_MT_Desktop_Adapter_packets_received_errors":      0,
		"net_Intel_R_PRO_1000_MT_Desktop_Adapter_packets_received_total":       10575000,
		"net_Intel_R_PRO_1000_MT_Desktop_Adapter_packets_received_unknown":     0,
		"net_Intel_R_PRO_1000_MT_Desktop_Adapter_packets_sent_total":           9605000,
		"net_Intel_R_PRO_1000_MT_Desktop_Adapter_packets_total":                20180000,
		"net_collection_duration":           0,
		"net_collection_success":            1,
		"os_collection_duration":            74,
		"os_collection_success":             1,
		"os_paging_free_bytes":              1468342272000,
		"os_paging_limit_bytes":             1476395008000,
		"os_physical_memory_free_bytes":     2621657088000,
		"os_process_memory_limit_bytes":     0,
		"os_processes":                      121,
		"os_processes_limit":                4294967295,
		"os_time":                           1577997682,
		"os_users":                          2,
		"os_virtual_memory_bytes":           5770891264000,
		"os_virtual_memory_free_bytes":      4116611072000,
		"os_visible_memory_bytes":           4294496256000,
		"system_boot_time":                  1577966173,
		"system_calls_total":                430573960000,
		"system_collection_duration":        0,
		"system_collection_success":         1,
		"system_context_switches_total":     49284345000,
		"system_exception_dispatches_total": 59337000,
		"system_processor_queue_length":     0,
		"system_threads":                    994,
		"system_up_time":                    77504,
	}

	collected := wmi.Collect()
	collected["system_up_time"] = expected["system_up_time"]

	assert.Equal(t, expected, collected)
	testCharts(t, wmi, collected)
}

func TestWMI_Collect_Partial(t *testing.T) {
	wmi, ts := prepareClientServerPartialData(t)
	defer ts.Close()

	expected := map[string]int64{
		"cpu_collection_duration":                    0,
		"cpu_collection_success":                     1,
		"cpu_core_0,0_c1":                            9996822,
		"cpu_core_0,0_c2":                            0,
		"cpu_core_0,0_c3":                            0,
		"cpu_core_0,0_dpc":                           2375,
		"cpu_core_0,0_dpcs":                          119284000,
		"cpu_core_0,0_idle":                          10181359,
		"cpu_core_0,0_interrupt":                     1843,
		"cpu_core_0,0_interrupts":                    17693361000,
		"cpu_core_0,0_privileged":                    185484,
		"cpu_core_0,0_user":                          269031,
		"cpu_core_0,1_c1":                            8774705,
		"cpu_core_0,1_c2":                            0,
		"cpu_core_0,1_c3":                            0,
		"cpu_core_0,1_dpc":                           20062,
		"cpu_core_0,1_dpcs":                          306334000,
		"cpu_core_0,1_idle":                          8831375,
		"cpu_core_0,1_interrupt":                     6109,
		"cpu_core_0,1_interrupts":                    2346399000,
		"cpu_core_0,1_privileged":                    861703,
		"cpu_core_0,1_user":                          942796,
		"cpu_core_0,2_c1":                            9941251,
		"cpu_core_0,2_c2":                            0,
		"cpu_core_0,2_c3":                            0,
		"cpu_core_0,2_dpc":                           12625,
		"cpu_core_0,2_dpcs":                          185430000,
		"cpu_core_0,2_idle":                          9998625,
		"cpu_core_0,2_interrupt":                     11062,
		"cpu_core_0,2_interrupts":                    2212620000,
		"cpu_core_0,2_privileged":                    297359,
		"cpu_core_0,2_user":                          339890,
		"cpu_core_0,3_c1":                            9453105,
		"cpu_core_0,3_c2":                            0,
		"cpu_core_0,3_c3":                            0,
		"cpu_core_0,3_dpc":                           13781,
		"cpu_core_0,3_dpcs":                          211020000,
		"cpu_core_0,3_idle":                          9463031,
		"cpu_core_0,3_interrupt":                     105625,
		"cpu_core_0,3_interrupts":                    2611467000,
		"cpu_core_0,3_privileged":                    411406,
		"cpu_core_0,3_user":                          761437,
		"cpu_dpc":                                    48843,
		"cpu_idle":                                   38474390,
		"cpu_interrupt":                              124640,
		"cpu_privileged":                             1755953,
		"cpu_user":                                   2313156,
		"memory_available_bytes":                     2337222656000,
		"memory_cache_bytes":                         128589824000,
		"memory_cache_bytes_peak":                    195198976000,
		"memory_cache_faults_total":                  7675068000,
		"memory_cache_total":                         2052911104000,
		"memory_collection_duration":                 0,
		"memory_collection_success":                  1,
		"memory_commit_limit":                        5770891264000,
		"memory_committed_bytes":                     2006388736000,
		"memory_demand_zero_faults_total":            6882552000,
		"memory_free_and_zero_page_list_bytes":       304807936000,
		"memory_free_system_page_table_entries":      12558411000,
		"memory_modified_page_list_bytes":            20496384000,
		"memory_not_committed_bytes":                 3764502528000,
		"memory_page_faults_total":                   18061429000,
		"memory_pool_nonpaged_allocs_total":          0,
		"memory_pool_nonpaged_bytes_total":           164827136000,
		"memory_pool_paged_allocs_total":             0,
		"memory_pool_paged_bytes":                    372215808000,
		"memory_pool_paged_resident_bytes":           359211008000,
		"memory_standby_cache_core_bytes":            165449728000,
		"memory_standby_cache_normal_priority_bytes": 600199168000,
		"memory_standby_cache_reserve_bytes":         1266765824000,
		"memory_standby_cache_total":                 2032414720000,
		"memory_swap_page_operations_total":          5396970000,
		"memory_swap_page_reads_total":               676801000,
		"memory_swap_page_writes_total":              804000,
		"memory_swap_pages_read_total":               5368093000,
		"memory_swap_pages_written_total":            28877000,
		"memory_system_cache_resident_bytes":         128589824000,
		"memory_system_code_resident_bytes":          0,
		"memory_system_code_total_bytes":             0,
		"memory_system_driver_resident_bytes":        9486336000,
		"memory_system_driver_total_bytes":           16224256000,
		"memory_transition_faults_total":             5279307000,
		"memory_transition_pages_repurposed_total":   2163369000,
		"memory_used_bytes":                          1957273600000,
		"memory_write_copies_total":                  166632000,
		"os_collection_duration":                     44,
		"os_collection_success":                      1,
		"os_paging_free_bytes":                       1445355520000,
		"os_paging_limit_bytes":                      1476395008000,
		"os_physical_memory_free_bytes":              2335821824000,
		"os_process_memory_limit_bytes":              0,
		"os_processes":                               124,
		"os_processes_limit":                         4294967295,
		"os_time":                                    1578049740,
		"os_users":                                   2,
		"os_virtual_memory_bytes":                    5770891264000,
		"os_virtual_memory_free_bytes":               3764899840000,
		"os_visible_memory_bytes":                    4294496256000,
	}

	collected := wmi.Collect()
	assert.Equal(t, expected, collected)
	testCharts(t, wmi, collected)
}

func TestWMI_CollectNoResponse(t *testing.T) {
	wmi := New()
	wmi.URL = "http://127.0.0.1:38001/jmx"
	require.True(t, wmi.Init())

	assert.Nil(t, wmi.Collect())
}

func TestWMI_Collect_ReceiveInvalidResponse(t *testing.T) {
	wmi, ts := prepareClientServerInvalidResponse(t)
	defer ts.Close()

	assert.Nil(t, wmi.Collect())
}

func TestWMI_Collect_Receive404(t *testing.T) {
	wmi, ts := prepareClientServerResponse404(t)
	defer ts.Close()

	assert.Nil(t, wmi.Collect())
}

func testCharts(t *testing.T, wmi *WMI, collected map[string]int64) {
	ensureChartsCreated(t, wmi)
	ensureChartsDynamicDimsCreated(t, wmi)
	ensureCollectedHasAllChartsDimsVarsIDs(t, wmi, collected)
}

func ensureChartsCreated(t *testing.T, w *WMI) {
	for _, chart := range cpuCharts() {
		if w.cache.collectors[collectorCPU] {
			assert.Truef(t, w.Charts().Has(chart.ID), "chart '%s' not created", chart.ID)
		} else {
			assert.Falsef(t, w.Charts().Has(chart.ID), "chart '%s' created", chart.ID)
		}
	}
	for _, chart := range memCharts() {
		if w.cache.collectors[collectorMemory] {
			assert.Truef(t, w.Charts().Has(chart.ID), "chart '%s' not created", chart.ID)
		} else {
			assert.Falsef(t, w.Charts().Has(chart.ID), "chart '%s' created", chart.ID)
		}
	}
	for _, chart := range osCharts() {
		if w.cache.collectors[collectorOS] {
			assert.Truef(t, w.Charts().Has(chart.ID), "chart '%s' not created", chart.ID)
		} else {
			assert.Falsef(t, w.Charts().Has(chart.ID), "chart '%s' created", chart.ID)
		}
	}
	for _, chart := range systemCharts() {
		if w.cache.collectors[collectorSystem] {
			assert.Truef(t, w.Charts().Has(chart.ID), "chart '%s' not created", chart.ID)
		} else {
			assert.Falsef(t, w.Charts().Has(chart.ID), "chart '%s' created", chart.ID)
		}
	}
	for _, chart := range logonCharts() {
		if w.cache.collectors[collectorLogon] {
			assert.Truef(t, w.Charts().Has(chart.ID), "chart '%s' not created", chart.ID)
		} else {
			assert.Falsef(t, w.Charts().Has(chart.ID), "chart '%s' created", chart.ID)
		}
	}
	for _, chart := range *collectionCharts() {
		assert.Truef(t, w.Charts().Has(chart.ID), "chart '%s' not created", chart.ID)
	}

	for coreID := range w.cache.cores {
		for _, chart := range cpuCoreCharts() {
			id := fmt.Sprintf(chart.ID, coreID)
			assert.Truef(t, w.Charts().Has(id), "charts has no '%s' chart for '%s' core", id, coreID)
		}
	}

	for nicID := range w.cache.nics {
		for _, chart := range nicCharts() {
			id := fmt.Sprintf(chart.ID, nicID)
			assert.Truef(t, w.Charts().Has(id), "charts has no '%s' chart for '%s' nic", id, nicID)
		}
	}

	for diskID := range w.cache.volumes {
		for _, chart := range diskCharts() {
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

func prepareClientServerFullData(t *testing.T) (*WMI, *httptest.Server) {
	t.Helper()
	ts := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			_, _ = w.Write(fullData)
		}))

	wmi := New()
	wmi.URL = ts.URL
	require.True(t, wmi.Init())
	return wmi, ts
}

func prepareClientServerPartialData(t *testing.T) (*WMI, *httptest.Server) {
	t.Helper()
	ts := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			_, _ = w.Write(partialData)
		}))

	wmi := New()
	wmi.URL = ts.URL
	require.True(t, wmi.Init())
	return wmi, ts
}

func prepareClientServerInvalidResponse(t *testing.T) (*WMI, *httptest.Server) {
	t.Helper()
	ts := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			_, _ = w.Write([]byte("hello and\n goodbye"))
		}))

	wmi := New()
	wmi.URL = ts.URL
	require.True(t, wmi.Init())
	return wmi, ts
}

func prepareClientServerResponse404(t *testing.T) (*WMI, *httptest.Server) {
	t.Helper()
	ts := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotFound)
		}))

	wmi := New()
	wmi.URL = ts.URL
	require.True(t, wmi.Init())
	return wmi, ts
}
