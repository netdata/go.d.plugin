// SPDX-License-Identifier: GPL-3.0-or-later

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
	dataXMLRTX2080Win, _ = os.ReadFile("testdata/rtx-2080-win.xml")
	dataXMLRTX3060, _    = os.ReadFile("testdata/rtx-3060.xml")
	dataXMLTeslaP100, _  = os.ReadFile("testdata/tesla-p100.xml")

	dataHelpQueryGPU, _ = os.ReadFile("testdata/help-query-gpu.txt")
	dataCSVTeslaP100, _ = os.ReadFile("testdata/tesla-p100.csv")
)

func Test_testDataIsValid(t *testing.T) {
	for name, data := range map[string][]byte{
		"dataXMLRTX2080Win": dataXMLRTX2080Win,
		"dataXMLRTX3060":    dataXMLRTX3060,
		"dataXMLTeslaP100":  dataXMLTeslaP100,
		"dataHelpQueryGPU":  dataHelpQueryGPU,
		"dataCSVTeslaP100":  dataCSVTeslaP100,
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
		"success RTX 3060 [XML]": {
			wantFail: false,
			prepare:  prepareCaseRTX3060formatXML,
		},
		"success Tesla P100 [XML]": {
			wantFail: false,
			prepare:  prepareCaseTeslaP100formatXML,
		},
		"success Tesla P100 [CSV]": {
			wantFail: false,
			prepare:  prepareCaseTeslaP100formatCSV,
		},
		"success RTX 2080 Win [XML]": {
			wantFail: false,
			prepare:  prepareCaseRTX2080WinFormatXML,
		},
		"fail on queryGPUInfoXML error": {
			wantFail: true,
			prepare:  prepareCaseErrOnQueryGPUInfoXML,
		},
		"fail on queryGPUInfoCSV error": {
			wantFail: true,
			prepare:  prepareCaseErrOnQueryGPUInfoCSV,
		},
		"fail on queryHelpQueryGPU error": {
			wantFail: true,
			prepare:  prepareCaseErrOnQueryHelpQueryGPU,
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
		"success RTX 3060 [XML]": {
			{
				prepare: prepareCaseRTX3060formatXML,
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
		"success Tesla P100 [XML]": {
			{
				prepare: prepareCaseTeslaP100formatXML,
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
		"success Tesla P100 [CSV]": {
			{
				prepare: prepareCaseTeslaP100formatCSV,
				check: func(t *testing.T, nv *NvidiaSMI) {
					mx := nv.Collect()

					expected := map[string]int64{
						"gpu_GPU-ef1b2c9b-38d8-2090-2bd1-f567a3eb42a6_frame_buffer_memory_usage_free":     17070817280,
						"gpu_GPU-ef1b2c9b-38d8-2090-2bd1-f567a3eb42a6_frame_buffer_memory_usage_reserved": 108003328,
						"gpu_GPU-ef1b2c9b-38d8-2090-2bd1-f567a3eb42a6_frame_buffer_memory_usage_used":     0,
						"gpu_GPU-ef1b2c9b-38d8-2090-2bd1-f567a3eb42a6_gpu_utilization":                    0,
						"gpu_GPU-ef1b2c9b-38d8-2090-2bd1-f567a3eb42a6_graphics_clock":                     405,
						"gpu_GPU-ef1b2c9b-38d8-2090-2bd1-f567a3eb42a6_mem_clock":                          715,
						"gpu_GPU-ef1b2c9b-38d8-2090-2bd1-f567a3eb42a6_mem_utilization":                    0,
						"gpu_GPU-ef1b2c9b-38d8-2090-2bd1-f567a3eb42a6_performance_state_P0":               1,
						"gpu_GPU-ef1b2c9b-38d8-2090-2bd1-f567a3eb42a6_performance_state_P1":               0,
						"gpu_GPU-ef1b2c9b-38d8-2090-2bd1-f567a3eb42a6_performance_state_P10":              0,
						"gpu_GPU-ef1b2c9b-38d8-2090-2bd1-f567a3eb42a6_performance_state_P11":              0,
						"gpu_GPU-ef1b2c9b-38d8-2090-2bd1-f567a3eb42a6_performance_state_P12":              0,
						"gpu_GPU-ef1b2c9b-38d8-2090-2bd1-f567a3eb42a6_performance_state_P13":              0,
						"gpu_GPU-ef1b2c9b-38d8-2090-2bd1-f567a3eb42a6_performance_state_P14":              0,
						"gpu_GPU-ef1b2c9b-38d8-2090-2bd1-f567a3eb42a6_performance_state_P15":              0,
						"gpu_GPU-ef1b2c9b-38d8-2090-2bd1-f567a3eb42a6_performance_state_P2":               0,
						"gpu_GPU-ef1b2c9b-38d8-2090-2bd1-f567a3eb42a6_performance_state_P3":               0,
						"gpu_GPU-ef1b2c9b-38d8-2090-2bd1-f567a3eb42a6_performance_state_P4":               0,
						"gpu_GPU-ef1b2c9b-38d8-2090-2bd1-f567a3eb42a6_performance_state_P5":               0,
						"gpu_GPU-ef1b2c9b-38d8-2090-2bd1-f567a3eb42a6_performance_state_P6":               0,
						"gpu_GPU-ef1b2c9b-38d8-2090-2bd1-f567a3eb42a6_performance_state_P7":               0,
						"gpu_GPU-ef1b2c9b-38d8-2090-2bd1-f567a3eb42a6_performance_state_P8":               0,
						"gpu_GPU-ef1b2c9b-38d8-2090-2bd1-f567a3eb42a6_performance_state_P9":               0,
						"gpu_GPU-ef1b2c9b-38d8-2090-2bd1-f567a3eb42a6_power_draw":                         28,
						"gpu_GPU-ef1b2c9b-38d8-2090-2bd1-f567a3eb42a6_sm_clock":                           405,
						"gpu_GPU-ef1b2c9b-38d8-2090-2bd1-f567a3eb42a6_temperature":                        37,
						"gpu_GPU-ef1b2c9b-38d8-2090-2bd1-f567a3eb42a6_video_clock":                        835,
					}

					assert.Equal(t, expected, mx)
				},
			},
		},
		"success RTX 2080 Win [XML]": {
			{
				prepare: prepareCaseRTX2080WinFormatXML,
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
		"fail on queryGPUInfoXML error [XML]": {
			{
				prepare: prepareCaseErrOnQueryGPUInfoXML,
				check: func(t *testing.T, nv *NvidiaSMI) {
					mx := nv.Collect()

					assert.Equal(t, map[string]int64(nil), mx)
				},
			},
		},
		"fail on queryGPUInfoCSV error [CSV]": {
			{
				prepare: prepareCaseErrOnQueryGPUInfoCSV,
				check: func(t *testing.T, nv *NvidiaSMI) {
					mx := nv.Collect()

					assert.Equal(t, map[string]int64(nil), mx)
				},
			},
		},
		"fail on queryHelpQueryGPU error": {
			{
				prepare: prepareCaseErrOnQueryHelpQueryGPU,
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
	gpuInfoXML           []byte
	errOnQueryGPUInfoXML bool

	gpuInfoCSV           []byte
	errOnQueryGPUInfoCSV bool

	helpQueryGPU           []byte
	errOnQueryHelpQueryGPU bool
}

func (m *mockNvidiaSMI) queryGPUInfoXML() ([]byte, error) {
	if m.errOnQueryGPUInfoXML {
		return nil, errors.New("error on mock.queryGPUInfoXML()")
	}
	return m.gpuInfoXML, nil
}

func (m *mockNvidiaSMI) queryGPUInfoCSV(_ []string) ([]byte, error) {
	if m.errOnQueryGPUInfoCSV {
		return nil, errors.New("error on mock.queryGPUInfoCSV()")
	}
	return m.gpuInfoCSV, nil
}

func (m *mockNvidiaSMI) queryHelpQueryGPU() ([]byte, error) {
	if m.errOnQueryHelpQueryGPU {
		return nil, errors.New("error on mock.queryHelpQueryGPU()")
	}
	return m.helpQueryGPU, nil
}

func prepareCaseRTX3060formatXML(nv *NvidiaSMI) {
	nv.UseCSVFormat = false
	nv.exec = &mockNvidiaSMI{gpuInfoXML: dataXMLRTX3060}
}

func prepareCaseTeslaP100formatXML(nv *NvidiaSMI) {
	nv.UseCSVFormat = false
	nv.exec = &mockNvidiaSMI{gpuInfoXML: dataXMLTeslaP100}
}

func prepareCaseRTX2080WinFormatXML(nv *NvidiaSMI) {
	nv.UseCSVFormat = false
	nv.exec = &mockNvidiaSMI{gpuInfoXML: dataXMLRTX2080Win}
}

func prepareCaseErrOnQueryGPUInfoXML(nv *NvidiaSMI) {
	nv.UseCSVFormat = false
	nv.exec = &mockNvidiaSMI{errOnQueryGPUInfoXML: true}
}

func prepareCaseTeslaP100formatCSV(nv *NvidiaSMI) {
	nv.UseCSVFormat = true
	nv.exec = &mockNvidiaSMI{helpQueryGPU: dataHelpQueryGPU, gpuInfoCSV: dataCSVTeslaP100}
}

func prepareCaseErrOnQueryHelpQueryGPU(nv *NvidiaSMI) {
	nv.UseCSVFormat = true
	nv.exec = &mockNvidiaSMI{errOnQueryHelpQueryGPU: true}
}

func prepareCaseErrOnQueryGPUInfoCSV(nv *NvidiaSMI) {
	nv.UseCSVFormat = true
	nv.exec = &mockNvidiaSMI{helpQueryGPU: dataHelpQueryGPU, errOnQueryGPUInfoCSV: true}
}
