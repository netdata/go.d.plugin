package wmi

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	testData, _ = ioutil.ReadFile("testdata/metrics.txt")
)

func TestNew(t *testing.T) {
	job := New()

	assert.IsType(t, (*WMI)(nil), job)
	assert.Equal(t, defaultURL, job.UserURL)
	assert.Equal(t, defaultHTTPTimeout, job.Timeout.Duration)
}

func TestWMI_Init(t *testing.T) {
	job := New()

	require.True(t, job.Init())
	assert.NotNil(t, job.prom)
	assert.NotNil(t, job.charts)
	require.NotNil(t, job.collected)
	assert.NotNil(t, job.collected.collectors)
	assert.NotNil(t, job.collected.cores)
	assert.NotNil(t, job.collected.nics)
}

func TestWMI_InitNG(t *testing.T) {
	job := New()
	job.UserURL = ""

	require.False(t, job.Init())
}

func TestWMI_Check(t *testing.T) {
	ts := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				_, _ = w.Write(testData)
			}))
	defer ts.Close()

	job := New()
	job.UserURL = ts.URL
	require.True(t, job.Init())
	assert.True(t, job.Check())
}

func TestWMI_CheckNG(t *testing.T) {
	job := New()
	job.UserURL = "http://127.0.0.1:38001/metrics"
	require.True(t, job.Init())
	assert.False(t, job.Check())
}

func TestWMI_Charts(t *testing.T) { assert.NotNil(t, New().Charts()) }

func TestWMI_Cleanup(t *testing.T) { assert.NotPanics(t, New().Cleanup) }

func TestWMI_Collect(t *testing.T) {
	ts := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				_, _ = w.Write(testData)
			}))
	defer ts.Close()

	job := New()
	job.UserURL = ts.URL
	require.True(t, job.Init())
	require.True(t, job.Check())

	//m := job.Collect()
	//l := make([]string, 0)
	//for k := range m {
	//	l = append(l, k)
	//}
	//sort.Strings(l)
	//for _, v := range l {
	//	fmt.Println(fmt.Sprintf("\"%s\": %d,", v, m[v]))
	//}

	expected := map[string]int64{
		"cpu":                                        627,
		"cpu_core_0_c1":                              1370542,
		"cpu_core_0_c2":                              0,
		"cpu_core_0_c3":                              0,
		"cpu_core_0_dpc":                             3656,
		"cpu_core_0_dpcs":                            174098000,
		"cpu_core_0_idle":                            1370542,
		"cpu_core_0_interrupt":                       1968,
		"cpu_core_0_interrupts":                      4385862000,
		"cpu_core_0_privileged":                      167906,
		"cpu_core_0_user":                            209250,
		"cpu_core_1_c1":                              1499317,
		"cpu_core_1_c2":                              0,
		"cpu_core_1_c3":                              0,
		"cpu_core_1_dpc":                             3125,
		"cpu_core_1_dpcs":                            134032000,
		"cpu_core_1_idle":                            1499317,
		"cpu_core_1_interrupt":                       28546,
		"cpu_core_1_interrupts":                      1143867000,
		"cpu_core_1_privileged":                      130906,
		"cpu_core_1_user":                            102218,
		"cpu_dpc":                                    6781,
		"cpu_idle":                                   2869859,
		"cpu_interrupt":                              30515,
		"cpu_privileged":                             298812,
		"cpu_user":                                   311468,
		"cs":                                         490,
		"cs_logical_processors":                      2,
		"cs_physical_memory_bytes":                   4294496256000,
		"memory":                                     595,
		"memory_available_bytes":                     2483294208000,
		"memory_cache_bytes":                         164429824000,
		"memory_cache_bytes_peak":                    275238912000,
		"memory_cache_faults_total":                  1054890000,
		"memory_commit_limit":                        5770891264000,
		"memory_committed_bytes":                     1607307264000,
		"memory_demand_zero_faults_total":            4485151000,
		"memory_free_and_zero_page_list_bytes":       1172709376000,
		"memory_free_system_page_table_entries":      12301246000,
		"memory_modified_page_list_bytes":            25849856000,
		"memory_page_faults_total":                   6374808000,
		"memory_pool_nonpaged_allocs_total":          162025000,
		"memory_pool_nonpaged_bytes_total":           124121088000,
		"memory_pool_paged_allocs_total":             282312000,
		"memory_pool_paged_bytes":                    271446016000,
		"memory_pool_paged_resident_bytes":           264552448000,
		"memory_standby_cache_core_bytes":            167759872000,
		"memory_standby_cache_normal_priority_bytes": 773877760000,
		"memory_standby_cache_reserve_bytes":         368947200000,
		"memory_swap_page_operations_total":          646571000,
		"memory_swap_page_reads_total":               172307000,
		"memory_swap_page_writes_total":              236000,
		"memory_swap_pages_read_total":               637983000,
		"memory_swap_pages_written_total":            8588000,
		"memory_system_cache_resident_bytes":         164429824000,
		"memory_system_code_resident_bytes":          0,
		"memory_system_code_total_bytes":             0,
		"memory_system_driver_resident_bytes":        10760192000,
		"memory_system_driver_total_bytes":           16642048000,
		"memory_transition_faults_total":             1205411000,
		"memory_transition_pages_repurposed_total":   238279000,
		"memory_write_copies_total":                  90377000,
		"net":                                        415,
		"net_Intel_R_PRO_1000_MT_Desktop_Adapter_2_bytes_received":             21238000,
		"net_Intel_R_PRO_1000_MT_Desktop_Adapter_2_bytes_sent":                 241197000,
		"net_Intel_R_PRO_1000_MT_Desktop_Adapter_2_bytes_total":                262435000,
		"net_Intel_R_PRO_1000_MT_Desktop_Adapter_2_current_bandwidth":          1000000000,
		"net_Intel_R_PRO_1000_MT_Desktop_Adapter_2_packets_outbound_discarded": 0,
		"net_Intel_R_PRO_1000_MT_Desktop_Adapter_2_packets_outbound_errors":    0,
		"net_Intel_R_PRO_1000_MT_Desktop_Adapter_2_packets_received_discarded": 0,
		"net_Intel_R_PRO_1000_MT_Desktop_Adapter_2_packets_received_errors":    0,
		"net_Intel_R_PRO_1000_MT_Desktop_Adapter_2_packets_received_total":     158000,
		"net_Intel_R_PRO_1000_MT_Desktop_Adapter_2_packets_received_unknown":   0,
		"net_Intel_R_PRO_1000_MT_Desktop_Adapter_2_packets_sent_total":         159000,
		"net_Intel_R_PRO_1000_MT_Desktop_Adapter_2_packets_total":              317000,
		"net_Intel_R_PRO_1000_MT_Desktop_Adapter_bytes_received":               39181285000,
		"net_Intel_R_PRO_1000_MT_Desktop_Adapter_bytes_sent":                   1568020000,
		"net_Intel_R_PRO_1000_MT_Desktop_Adapter_bytes_total":                  40749305000,
		"net_Intel_R_PRO_1000_MT_Desktop_Adapter_current_bandwidth":            1000000000,
		"net_Intel_R_PRO_1000_MT_Desktop_Adapter_packets_outbound_discarded":   0,
		"net_Intel_R_PRO_1000_MT_Desktop_Adapter_packets_outbound_errors":      0,
		"net_Intel_R_PRO_1000_MT_Desktop_Adapter_packets_received_discarded":   0,
		"net_Intel_R_PRO_1000_MT_Desktop_Adapter_packets_received_errors":      0,
		"net_Intel_R_PRO_1000_MT_Desktop_Adapter_packets_received_total":       48369000,
		"net_Intel_R_PRO_1000_MT_Desktop_Adapter_packets_received_unknown":     0,
		"net_Intel_R_PRO_1000_MT_Desktop_Adapter_packets_sent_total":           20685000,
		"net_Intel_R_PRO_1000_MT_Desktop_Adapter_packets_total":                69054000,
		"os":                                567,
		"os_paging_free_bytes":              1451696128000,
		"os_paging_limit_bytes":             1476395008000,
		"os_physical_memory_free_bytes":     2483429376000,
		"os_process_memory_limit_bytes":     0,
		"os_processes":                      116,
		"os_processes_limit":                4294967295,
		"os_time":                           1558031630,
		"os_users":                          2,
		"os_virtual_memory_bytes":           5770891264000,
		"os_virtual_memory_free_bytes":      4163584000000,
		"os_visible_memory_bytes":           4294496256000,
		"system":                            511,
		"system_context_switches_total":     4549390000,
		"system_exception_dispatches_total": 4441000,
		"system_processor_queue_length":     0,
		"system_system_calls_total":         36309080000,
		"system_system_threads":             1320,
		"system_system_up_time":             1558029847,
	}

	assert.Equal(t, expected, job.Collect())
}

func TestWMI_InvalidData(t *testing.T) {
	ts := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				_, _ = w.Write([]byte("hello and goodbye"))
			}))
	defer ts.Close()

	job := New()
	job.UserURL = ts.URL
	require.True(t, job.Init())
	assert.False(t, job.Check())
}

func TestWMI_404(t *testing.T) {
	ts := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusNotFound)
			}))
	defer ts.Close()

	job := New()
	job.UserURL = ts.URL
	require.True(t, job.Init())
	assert.False(t, job.Check())
}
