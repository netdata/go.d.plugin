// SPDX-License-Identifier: GPL-3.0-or-later

package dyncfg

import (
	"errors"
	"sync"
	"testing"

	"github.com/netdata/go.d.plugin/agent/confgroup"
	"github.com/stretchr/testify/assert"
)

func TestNewDiscovery(t *testing.T) {

}

func TestDiscovery_Register(t *testing.T) {
	tests := map[string]struct {
		regConfigs   []confgroup.Config
		wantApiStats *mockNetdataDyncfgAPI
		wantConfigs  int
	}{
		"register jobs created by Dyncfg and other providers": {
			regConfigs: []confgroup.Config{
				prepareConfig(
					"__provider__", dynCfg,
					"module", "test",
					"name", "first",
				),
				prepareConfig(
					"__provider__", "test",
					"module", "test",
					"name", "second",
				),
			},
			wantConfigs: 2,
			wantApiStats: &mockNetdataDyncfgAPI{
				callsDynCfgRegisterJob: 1,
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			var mock mockNetdataDyncfgAPI
			d := &Discovery{
				API:     &mock,
				mux:     &sync.Mutex{},
				configs: make(map[string]confgroup.Config),
			}

			for _, v := range test.regConfigs {
				d.Register(v)
			}

			assert.Equal(t, test.wantApiStats, &mock)
			assert.Equal(t, test.wantConfigs, len(d.configs))
		})
	}
}

func TestDiscovery_Unregister(t *testing.T) {
	tests := map[string]struct {
		regConfigs   []confgroup.Config
		unregConfigs []confgroup.Config
		wantApiStats *mockNetdataDyncfgAPI
		wantConfigs  int
	}{
		"register/unregister jobs created by Dyncfg and other providers": {
			wantConfigs: 0,
			wantApiStats: &mockNetdataDyncfgAPI{
				callsDynCfgRegisterJob: 1,
			},
			regConfigs: []confgroup.Config{
				prepareConfig(
					"__provider__", dynCfg,
					"module", "test",
					"name", "first",
				),
				prepareConfig(
					"__provider__", "test",
					"module", "test",
					"name", "second",
				),
			},
			unregConfigs: []confgroup.Config{
				prepareConfig(
					"__provider__", dynCfg,
					"module", "test",
					"name", "first",
				),
				prepareConfig(
					"__provider__", "test",
					"module", "test",
					"name", "second",
				),
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			var mock mockNetdataDyncfgAPI
			d := &Discovery{
				API:     &mock,
				mux:     &sync.Mutex{},
				configs: make(map[string]confgroup.Config),
			}

			for _, v := range test.regConfigs {
				d.Register(v)
			}
			for _, v := range test.unregConfigs {
				d.Unregister(v)
			}

			assert.Equal(t, test.wantApiStats, &mock)
			assert.Equal(t, test.wantConfigs, len(d.configs))
		})
	}
}

func TestDiscovery_UpdateStatus(t *testing.T) {

}

func TestDiscovery_Run(t *testing.T) {

}

type mockNetdataDyncfgAPI struct {
	errOnDynCfgEnable          bool
	errOnDyncCfgRegisterModule bool
	errOnDynCfgRegisterJob     bool
	errOnDynCfgReportJobStatus bool
	errOnFunctionResultSuccess bool
	errOnFunctionResultReject  bool

	callsDynCfgEnable          int
	callsDyncCfgRegisterModule int
	callsDynCfgRegisterJob     int
	callsDynCfgReportJobStatus int
	callsFunctionResultSuccess int
	callsFunctionResultReject  int
}

func (m *mockNetdataDyncfgAPI) DynCfgEnable(string) error {
	m.callsDynCfgEnable++
	if m.errOnDynCfgEnable {
		return errors.New("mock error on DynCfgEnable()")
	}
	return nil
}

func (m *mockNetdataDyncfgAPI) DyncCfgRegisterModule(string) error {
	m.callsDyncCfgRegisterModule++
	if m.errOnDyncCfgRegisterModule {
		return errors.New("mock error on DyncCfgRegisterModule()")
	}
	return nil
}

func (m *mockNetdataDyncfgAPI) DynCfgRegisterJob(_, _, _ string) error {
	m.callsDynCfgRegisterJob++
	if m.errOnDynCfgRegisterJob {
		return errors.New("mock error on DynCfgRegisterJob()")
	}
	return nil
}

func (m *mockNetdataDyncfgAPI) DynCfgReportJobStatus(_, _, _, _ string) error {
	m.callsDynCfgReportJobStatus++
	if m.errOnDynCfgReportJobStatus {
		return errors.New("mock error on DynCfgReportJobStatus()")
	}
	return nil
}

func (m *mockNetdataDyncfgAPI) FunctionResultSuccess(_, _, _ string) error {
	m.callsFunctionResultSuccess++
	if m.errOnFunctionResultSuccess {
		return errors.New("mock error on FunctionResultSuccess()")
	}
	return nil
}

func (m *mockNetdataDyncfgAPI) FunctionResultReject(_, _, _ string) error {
	m.callsFunctionResultReject++
	if m.errOnFunctionResultReject {
		return errors.New("mock error on FunctionResultReject()")
	}
	return nil
}

func prepareConfig(values ...string) confgroup.Config {
	cfg := confgroup.Config{}
	for i := 1; i < len(values); i += 2 {
		cfg[values[i-1]] = values[i]
	}
	return cfg
}
