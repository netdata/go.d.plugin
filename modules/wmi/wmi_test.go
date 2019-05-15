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
		"cpu_core_0_c1":         40718296,
		"cpu_core_0_c2":         0,
		"cpu_core_0_c3":         0,
		"cpu_core_0_dpc":        2640,
		"cpu_core_0_idle":       40718296,
		"cpu_core_0_interrupt":  4156,
		"cpu_core_0_interrupts": 87289519000,
		"cpu_core_0_privileged": 362437,
		"cpu_core_0_user":       258953,
		"cpu_core_1_c1":         40765354,
		"cpu_core_1_c2":         0,
		"cpu_core_1_c3":         0,
		"cpu_core_1_dpc":        7843,
		"cpu_core_1_idle":       40765354,
		"cpu_core_1_interrupt":  22031,
		"cpu_core_1_interrupts": 7798273000,
		"cpu_core_1_privileged": 270187,
		"cpu_core_1_user":       319828,
		"cpu_dpc":               10484,
		"cpu_idle":              81483650,
		"cpu_interrupt":         26187,
		"cpu_privileged":        632625,
		"cpu_user":              578781,
		"net_Intel_R_PRO_1000_MT_Desktop_Adapter_2_bytes_received":             3075019000,
		"net_Intel_R_PRO_1000_MT_Desktop_Adapter_2_bytes_sent":                 24108436000,
		"net_Intel_R_PRO_1000_MT_Desktop_Adapter_2_bytes_total":                27183455000,
		"net_Intel_R_PRO_1000_MT_Desktop_Adapter_2_current_bandwidth":          1000000000,
		"net_Intel_R_PRO_1000_MT_Desktop_Adapter_2_packets_outbound_discarded": 0,
		"net_Intel_R_PRO_1000_MT_Desktop_Adapter_2_packets_outbound_errors":    0,
		"net_Intel_R_PRO_1000_MT_Desktop_Adapter_2_packets_received_discarded": 0,
		"net_Intel_R_PRO_1000_MT_Desktop_Adapter_2_packets_received_errors":    0,
		"net_Intel_R_PRO_1000_MT_Desktop_Adapter_2_packets_received_total":     22450000,
		"net_Intel_R_PRO_1000_MT_Desktop_Adapter_2_packets_received_unknown":   0,
		"net_Intel_R_PRO_1000_MT_Desktop_Adapter_2_packets_sent_total":         21808000,
		"net_Intel_R_PRO_1000_MT_Desktop_Adapter_2_packets_total":              44258000,
		"net_Intel_R_PRO_1000_MT_Desktop_Adapter_bytes_received":               5342635000,
		"net_Intel_R_PRO_1000_MT_Desktop_Adapter_bytes_sent":                   1218739000,
		"net_Intel_R_PRO_1000_MT_Desktop_Adapter_bytes_total":                  6561374000,
		"net_Intel_R_PRO_1000_MT_Desktop_Adapter_current_bandwidth":            1000000000,
		"net_Intel_R_PRO_1000_MT_Desktop_Adapter_packets_outbound_discarded":   0,
		"net_Intel_R_PRO_1000_MT_Desktop_Adapter_packets_outbound_errors":      0,
		"net_Intel_R_PRO_1000_MT_Desktop_Adapter_packets_received_discarded":   0,
		"net_Intel_R_PRO_1000_MT_Desktop_Adapter_packets_received_errors":      0,
		"net_Intel_R_PRO_1000_MT_Desktop_Adapter_packets_received_total":       9900000,
		"net_Intel_R_PRO_1000_MT_Desktop_Adapter_packets_received_unknown":     0,
		"net_Intel_R_PRO_1000_MT_Desktop_Adapter_packets_sent_total":           6893000,
		"net_Intel_R_PRO_1000_MT_Desktop_Adapter_packets_total":                16793000,
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
