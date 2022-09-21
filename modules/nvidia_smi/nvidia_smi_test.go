package nvidia_smi

import (
	"errors"
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	dataRTX2080Win, _ = os.ReadFile("testdata/rtx-2080-win.xml")
	dataRTX3060, _    = os.ReadFile("testdata/rtx-3060.xml")
	dataTeslaP100, _  = os.ReadFile("testdata/tesla-p100.xml")
)

func Test_testDataIsValid(t *testing.T) {
	for name, data := range map[string][]byte{
		"dataRTX2080Win": dataRTX2080Win,
		"dataRTX3060":    dataRTX3060,
		"dataTeslaP100":  dataTeslaP100,
	} {
		require.NotNilf(t, data, name)
	}
}

func TestNvidiaSMI_Init(t *testing.T) {
	tests := map[string]struct {
		prepare  func(nv *NvidiaSMI)
		wantFail bool
	}{
		"fails if can't local nvidia-smi": {
			wantFail: true,
			prepare: func(nv *NvidiaSMI) {
				nv.binName += "!!!"
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			nv := New()

			test.prepare(nv)

			if test.wantFail {
				assert.False(t, nv.Init())
			} else {
				assert.True(t, nv.Init())
			}
		})
	}
}

func TestNvidiaSMI_Charts(t *testing.T) {
	assert.NotNil(t, New().Charts())
}

func TestNvidiaSMI_Check(t *testing.T) {
	tests := map[string]struct {
		prepare  func(nv *NvidiaSMI)
		wantFail bool
	}{
		"success RTX 3060": {
			wantFail: false,
			prepare:  prepareCaseRTX3060,
		},
		"success Tesla P100": {
			wantFail: false,
			prepare:  prepareCaseTeslaP100,
		},
		"success RTX 2080 Win": {
			wantFail: false,
			prepare:  prepareCaseRTX2080Win,
		},
		"fail on queryXML error": {
			wantFail: true,
			prepare:  prepareCaseErrOnQueryXML,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			nv := New()

			test.prepare(nv)

			if test.wantFail {
				assert.False(t, nv.Check())
			} else {
				assert.True(t, nv.Check())
			}
		})
	}
}

func TestNvidiaSMI_Collect(t *testing.T) {
	type testCaseStep struct {
		prepare func(nv *NvidiaSMI)
		check   func(t *testing.T, nv *NvidiaSMI)
	}
	tests := map[string][]testCaseStep{
		"success RTX 3060": {
			{
				prepare: prepareCaseRTX3060,
				check: func(t *testing.T, nv *NvidiaSMI) {
					mx := nv.Collect()

					expected := map[string]int64{
						"gpu_GPU-473d8d0f-d462-185c-6b36-6fc23e23e571_bar1_memory_usage_free":             8586788864,
						"gpu_GPU-473d8d0f-d462-185c-6b36-6fc23e23e571_bar1_memory_usage_used":             3145728,
						"gpu_GPU-473d8d0f-d462-185c-6b36-6fc23e23e571_decoder_utilization":                0,
						"gpu_GPU-473d8d0f-d462-185c-6b36-6fc23e23e571_encoder_utilization":                0,
						"gpu_GPU-473d8d0f-d462-185c-6b36-6fc23e23e571_frame_buffer_memory_usage_free":     6228541440,
						"gpu_GPU-473d8d0f-d462-185c-6b36-6fc23e23e571_frame_buffer_memory_usage_reserved": 206569472,
						"gpu_GPU-473d8d0f-d462-185c-6b36-6fc23e23e571_frame_buffer_memory_usage_used":     5242880,
						"gpu_GPU-473d8d0f-d462-185c-6b36-6fc23e23e571_gpu_utilization":                    0,
						"gpu_GPU-473d8d0f-d462-185c-6b36-6fc23e23e571_graphics_clock":                     210,
						"gpu_GPU-473d8d0f-d462-185c-6b36-6fc23e23e571_mem_clock":                          405,
						"gpu_GPU-473d8d0f-d462-185c-6b36-6fc23e23e571_mem_utilization":                    0,
						"gpu_GPU-473d8d0f-d462-185c-6b36-6fc23e23e571_pcie_bandwidth_usage_rx":            0,
						"gpu_GPU-473d8d0f-d462-185c-6b36-6fc23e23e571_pcie_bandwidth_usage_tx":            0,
						"gpu_GPU-473d8d0f-d462-185c-6b36-6fc23e23e571_performance_state_P0":               0,
						"gpu_GPU-473d8d0f-d462-185c-6b36-6fc23e23e571_performance_state_P1":               0,
						"gpu_GPU-473d8d0f-d462-185c-6b36-6fc23e23e571_performance_state_P10":              0,
						"gpu_GPU-473d8d0f-d462-185c-6b36-6fc23e23e571_performance_state_P11":              0,
						"gpu_GPU-473d8d0f-d462-185c-6b36-6fc23e23e571_performance_state_P12":              0,
						"gpu_GPU-473d8d0f-d462-185c-6b36-6fc23e23e571_performance_state_P13":              0,
						"gpu_GPU-473d8d0f-d462-185c-6b36-6fc23e23e571_performance_state_P14":              0,
						"gpu_GPU-473d8d0f-d462-185c-6b36-6fc23e23e571_performance_state_P15":              0,
						"gpu_GPU-473d8d0f-d462-185c-6b36-6fc23e23e571_performance_state_P2":               0,
						"gpu_GPU-473d8d0f-d462-185c-6b36-6fc23e23e571_performance_state_P3":               0,
						"gpu_GPU-473d8d0f-d462-185c-6b36-6fc23e23e571_performance_state_P4":               0,
						"gpu_GPU-473d8d0f-d462-185c-6b36-6fc23e23e571_performance_state_P5":               0,
						"gpu_GPU-473d8d0f-d462-185c-6b36-6fc23e23e571_performance_state_P6":               0,
						"gpu_GPU-473d8d0f-d462-185c-6b36-6fc23e23e571_performance_state_P7":               0,
						"gpu_GPU-473d8d0f-d462-185c-6b36-6fc23e23e571_performance_state_P8":               1,
						"gpu_GPU-473d8d0f-d462-185c-6b36-6fc23e23e571_performance_state_P9":               0,
						"gpu_GPU-473d8d0f-d462-185c-6b36-6fc23e23e571_power_draw":                         8,
						"gpu_GPU-473d8d0f-d462-185c-6b36-6fc23e23e571_sm_clock":                           210,
						"gpu_GPU-473d8d0f-d462-185c-6b36-6fc23e23e571_temperature":                        45,
						"gpu_GPU-473d8d0f-d462-185c-6b36-6fc23e23e571_video_clock":                        555,
					}

					assert.Equal(t, expected, mx)
				},
			},
		},
		"success Tesla P100": {
			{
				prepare: prepareCaseTeslaP100,
				check: func(t *testing.T, nv *NvidiaSMI) {
					mx := nv.Collect()

					expected := map[string]int64{
						"gpu_GPU-d3da8716-eaab-75db-efc1-60e88e1cd55e_bar1_memory_usage_free":             17177772032,
						"gpu_GPU-d3da8716-eaab-75db-efc1-60e88e1cd55e_bar1_memory_usage_used":             2097152,
						"gpu_GPU-d3da8716-eaab-75db-efc1-60e88e1cd55e_decoder_utilization":                0,
						"gpu_GPU-d3da8716-eaab-75db-efc1-60e88e1cd55e_encoder_utilization":                0,
						"gpu_GPU-d3da8716-eaab-75db-efc1-60e88e1cd55e_frame_buffer_memory_usage_free":     17070817280,
						"gpu_GPU-d3da8716-eaab-75db-efc1-60e88e1cd55e_frame_buffer_memory_usage_reserved": 108003328,
						"gpu_GPU-d3da8716-eaab-75db-efc1-60e88e1cd55e_frame_buffer_memory_usage_used":     0,
						"gpu_GPU-d3da8716-eaab-75db-efc1-60e88e1cd55e_gpu_utilization":                    0,
						"gpu_GPU-d3da8716-eaab-75db-efc1-60e88e1cd55e_graphics_clock":                     405,
						"gpu_GPU-d3da8716-eaab-75db-efc1-60e88e1cd55e_mem_clock":                          715,
						"gpu_GPU-d3da8716-eaab-75db-efc1-60e88e1cd55e_mem_utilization":                    0,
						"gpu_GPU-d3da8716-eaab-75db-efc1-60e88e1cd55e_pcie_bandwidth_usage_rx":            0,
						"gpu_GPU-d3da8716-eaab-75db-efc1-60e88e1cd55e_pcie_bandwidth_usage_tx":            0,
						"gpu_GPU-d3da8716-eaab-75db-efc1-60e88e1cd55e_performance_state_P0":               1,
						"gpu_GPU-d3da8716-eaab-75db-efc1-60e88e1cd55e_performance_state_P1":               0,
						"gpu_GPU-d3da8716-eaab-75db-efc1-60e88e1cd55e_performance_state_P10":              0,
						"gpu_GPU-d3da8716-eaab-75db-efc1-60e88e1cd55e_performance_state_P11":              0,
						"gpu_GPU-d3da8716-eaab-75db-efc1-60e88e1cd55e_performance_state_P12":              0,
						"gpu_GPU-d3da8716-eaab-75db-efc1-60e88e1cd55e_performance_state_P13":              0,
						"gpu_GPU-d3da8716-eaab-75db-efc1-60e88e1cd55e_performance_state_P14":              0,
						"gpu_GPU-d3da8716-eaab-75db-efc1-60e88e1cd55e_performance_state_P15":              0,
						"gpu_GPU-d3da8716-eaab-75db-efc1-60e88e1cd55e_performance_state_P2":               0,
						"gpu_GPU-d3da8716-eaab-75db-efc1-60e88e1cd55e_performance_state_P3":               0,
						"gpu_GPU-d3da8716-eaab-75db-efc1-60e88e1cd55e_performance_state_P4":               0,
						"gpu_GPU-d3da8716-eaab-75db-efc1-60e88e1cd55e_performance_state_P5":               0,
						"gpu_GPU-d3da8716-eaab-75db-efc1-60e88e1cd55e_performance_state_P6":               0,
						"gpu_GPU-d3da8716-eaab-75db-efc1-60e88e1cd55e_performance_state_P7":               0,
						"gpu_GPU-d3da8716-eaab-75db-efc1-60e88e1cd55e_performance_state_P8":               0,
						"gpu_GPU-d3da8716-eaab-75db-efc1-60e88e1cd55e_performance_state_P9":               0,
						"gpu_GPU-d3da8716-eaab-75db-efc1-60e88e1cd55e_power_draw":                         26,
						"gpu_GPU-d3da8716-eaab-75db-efc1-60e88e1cd55e_sm_clock":                           405,
						"gpu_GPU-d3da8716-eaab-75db-efc1-60e88e1cd55e_temperature":                        38,
						"gpu_GPU-d3da8716-eaab-75db-efc1-60e88e1cd55e_video_clock":                        835,
					}

					assert.Equal(t, expected, mx)
				},
			},
		},
		"success RTX 2080 Win": {
			{
				prepare: prepareCaseRTX2080Win,
				check: func(t *testing.T, nv *NvidiaSMI) {
					mx := nv.Collect()

					expected := map[string]int64{
						"gpu_GPU-fbd55ed4-1eec-4423-0a47-ad594b4333e3_bar1_memory_usage_free":             266338304,
						"gpu_GPU-fbd55ed4-1eec-4423-0a47-ad594b4333e3_bar1_memory_usage_used":             2097152,
						"gpu_GPU-fbd55ed4-1eec-4423-0a47-ad594b4333e3_decoder_utilization":                0,
						"gpu_GPU-fbd55ed4-1eec-4423-0a47-ad594b4333e3_encoder_utilization":                0,
						"gpu_GPU-fbd55ed4-1eec-4423-0a47-ad594b4333e3_fan_speed_perc":                     37,
						"gpu_GPU-fbd55ed4-1eec-4423-0a47-ad594b4333e3_frame_buffer_memory_usage_free":     7494172672,
						"gpu_GPU-fbd55ed4-1eec-4423-0a47-ad594b4333e3_frame_buffer_memory_usage_reserved": 190840832,
						"gpu_GPU-fbd55ed4-1eec-4423-0a47-ad594b4333e3_frame_buffer_memory_usage_used":     903872512,
						"gpu_GPU-fbd55ed4-1eec-4423-0a47-ad594b4333e3_gpu_utilization":                    2,
						"gpu_GPU-fbd55ed4-1eec-4423-0a47-ad594b4333e3_graphics_clock":                     193,
						"gpu_GPU-fbd55ed4-1eec-4423-0a47-ad594b4333e3_mem_clock":                          403,
						"gpu_GPU-fbd55ed4-1eec-4423-0a47-ad594b4333e3_mem_utilization":                    7,
						"gpu_GPU-fbd55ed4-1eec-4423-0a47-ad594b4333e3_pcie_bandwidth_usage_rx":            93184000,
						"gpu_GPU-fbd55ed4-1eec-4423-0a47-ad594b4333e3_pcie_bandwidth_usage_tx":            13312000,
						"gpu_GPU-fbd55ed4-1eec-4423-0a47-ad594b4333e3_performance_state_P0":               0,
						"gpu_GPU-fbd55ed4-1eec-4423-0a47-ad594b4333e3_performance_state_P1":               0,
						"gpu_GPU-fbd55ed4-1eec-4423-0a47-ad594b4333e3_performance_state_P10":              0,
						"gpu_GPU-fbd55ed4-1eec-4423-0a47-ad594b4333e3_performance_state_P11":              0,
						"gpu_GPU-fbd55ed4-1eec-4423-0a47-ad594b4333e3_performance_state_P12":              0,
						"gpu_GPU-fbd55ed4-1eec-4423-0a47-ad594b4333e3_performance_state_P13":              0,
						"gpu_GPU-fbd55ed4-1eec-4423-0a47-ad594b4333e3_performance_state_P14":              0,
						"gpu_GPU-fbd55ed4-1eec-4423-0a47-ad594b4333e3_performance_state_P15":              0,
						"gpu_GPU-fbd55ed4-1eec-4423-0a47-ad594b4333e3_performance_state_P2":               0,
						"gpu_GPU-fbd55ed4-1eec-4423-0a47-ad594b4333e3_performance_state_P3":               0,
						"gpu_GPU-fbd55ed4-1eec-4423-0a47-ad594b4333e3_performance_state_P4":               0,
						"gpu_GPU-fbd55ed4-1eec-4423-0a47-ad594b4333e3_performance_state_P5":               0,
						"gpu_GPU-fbd55ed4-1eec-4423-0a47-ad594b4333e3_performance_state_P6":               0,
						"gpu_GPU-fbd55ed4-1eec-4423-0a47-ad594b4333e3_performance_state_P7":               0,
						"gpu_GPU-fbd55ed4-1eec-4423-0a47-ad594b4333e3_performance_state_P8":               1,
						"gpu_GPU-fbd55ed4-1eec-4423-0a47-ad594b4333e3_performance_state_P9":               0,
						"gpu_GPU-fbd55ed4-1eec-4423-0a47-ad594b4333e3_power_draw":                         14,
						"gpu_GPU-fbd55ed4-1eec-4423-0a47-ad594b4333e3_sm_clock":                           193,
						"gpu_GPU-fbd55ed4-1eec-4423-0a47-ad594b4333e3_temperature":                        29,
						"gpu_GPU-fbd55ed4-1eec-4423-0a47-ad594b4333e3_video_clock":                        539,
					}

					assert.Equal(t, expected, mx)
				},
			},
		},
		"fail on queryXML error": {
			{
				prepare: prepareCaseErrOnQueryXML,
				check: func(t *testing.T, nv *NvidiaSMI) {
					mx := nv.Collect()

					assert.Equal(t, map[string]int64(nil), mx)
				},
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			nv := New()

			for i, step := range test {
				t.Run(fmt.Sprintf("step[%d]", i), func(t *testing.T) {
					step.prepare(nv)
					step.check(t, nv)
				})
			}
		})
	}
}

type mockNvidiaSMI struct {
	respXML       []byte
	errOnQueryXML bool
}

func (m *mockNvidiaSMI) queryXML() ([]byte, error) {
	if m.errOnQueryXML {
		return nil, errors.New("error on mock.queryXML()")
	}
	return m.respXML, nil
}

func prepareCaseRTX3060(nv *NvidiaSMI) {
	nv.exec = &mockNvidiaSMI{respXML: dataRTX3060}
}

func prepareCaseTeslaP100(nv *NvidiaSMI) {
	nv.exec = &mockNvidiaSMI{respXML: dataTeslaP100}
}

func prepareCaseRTX2080Win(nv *NvidiaSMI) {
	nv.exec = &mockNvidiaSMI{respXML: dataRTX2080Win}
}

func prepareCaseErrOnQueryXML(nv *NvidiaSMI) {
	nv.exec = &mockNvidiaSMI{errOnQueryXML: true}
}
