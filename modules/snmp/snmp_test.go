package snmp

import (
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/gosnmp/gosnmp"
	snmpmock "github.com/gosnmp/gosnmp/mocks"
	"github.com/netdata/go.d.plugin/agent/module"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Define some configs
var (
	community   = "public"
	cType       = "area"
	cFamily     = "lan"
	cAlgorithm  = module.Incremental
	cMultiplier = 8
	cDivisor    = 1024
)

func TestNew(t *testing.T) {
	// We want to ensure that module is a reference type, nothing more.

	assert.IsType(t, (*SNMP)(nil), New())
}

func TestSNMP_Init(t *testing.T) {
	// 'Init() bool' initializes the module with an appropriate config, so to test it we need:
	// - provide the config.
	// - set module.Config field with the config.
	// - call Init() and compare its return value with the expected value.

	// 'test' map contains different test cases.
	mockSNMP, cleanup := mockInit(t)
	defer cleanup()

	defaultMockExpects(mockSNMP)

	tests := map[string]struct {
		config   Config
		wantFail bool
	}{
		"fail without 'charts' set": {
			config:   prepareConfigWithoutChart(),
			wantFail: true,
		},
		"success with 'charts' and 'dimensions' set": {
			config: prepareConfigWithDimensions(),
		},
		"fail with 'charts' but no 'dimensions' set": {
			config:   prepareConfigWithoutDimensions(),
			wantFail: true,
		},
		"success with 'charts' but invalid 'multiply_range' set": {
			config: prepareConfigWithMultiplyRange(),
		},
		"success with 'community' set for 'options.version=2' set": {
			config: prepareConfigWithCommunity(),
		},
		"fail when 'user' unset for 'options.version=3' set": {
			config:   prepareConfigWithoutUser(),
			wantFail: true,
		},
		"fail when 'community' unset for 'options.version=2' set": {
			config:   prepareConfigWithoutCommunity(),
			wantFail: true,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			SNMP := New()
			SNMP.Config = test.config
			newSNMPClient = func() gosnmp.Handler {
				return mockSNMP
			}
			if test.wantFail {
				assert.False(t, SNMP.Init())
			} else {
				assert.True(t, SNMP.Init())
			}
		})
	}
}

func TestSNMP_Check(t *testing.T) {
	// 'Check() bool' reports whether the module is able to collect any data, so to test it we need:
	// - provide the module with a specific config.
	// - initialize the module (call Init()).
	// - call Check() and compare its return value with the expected value.

	// 'test' map contains different test cases.

	returnSNMPpacket := gosnmp.SnmpPacket{
		Variables: []gosnmp.SnmpPDU{
			{Value: 10}, // Our configs/defaults have 2 OIDs
			{Value: 20},
		},
	}

	tests := map[string]struct {
		prepare  func(m *snmpmock.MockHandler) (s *SNMP)
		wantFail bool
	}{
		"success when 'dimensions' set": {
			prepare: func(m *snmpmock.MockHandler) *SNMP {
				snmp := New()
				snmp.Config = prepareConfigWithDimensions()
				newSNMPClient = func() gosnmp.Handler {
					return m
				}
				m.EXPECT().Get(gomock.Any()).Return(&returnSNMPpacket, nil).Times(1)
				return snmp
			},
		},
		"success when 'max_request_size' set": {
			prepare: func(m *snmpmock.MockHandler) *SNMP {
				snmp := New()
				snmp.Config = prepareConfigWithDimensions()
				newSNMPClient = func() gosnmp.Handler {
					return m
				}
				snmp.Options.MaxOIDs = 1

				//Get() must be called twice if MaxOIDs = 1
				m.EXPECT().Get(gomock.Any()).Return(&returnSNMPpacket, nil).Times(2)
				return snmp
			},
		},
		"success when 'multiply_range' set": {
			prepare: func(m *snmpmock.MockHandler) *SNMP {
				snmp := New()
				snmpPacket := gosnmp.SnmpPacket{
					Variables: []gosnmp.SnmpPDU{
						{Value: 10},
						{Value: 20},
						{Value: 30},
						{Value: 40},
						{Value: 50},
					},
				}

				snmp.Config = prepareConfigWithMultiplyRange()
				newSNMPClient = func() gosnmp.Handler {
					return m
				}
				m.EXPECT().Get(gomock.Any()).Return(&snmpPacket, nil).Times(1)
				return snmp
			},
		},
		"fail when chart collection fails": {
			prepare: func(m *snmpmock.MockHandler) *SNMP {
				snmp := New()
				snmp.Config = prepareConfigWithDimensions()
				newSNMPClient = func() gosnmp.Handler {
					return m
				}

				m.EXPECT().Get(gomock.Any()).Return(nil,
					fmt.Errorf("error from mock function")).Times(1)
				return snmp
			},
			wantFail: true,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			mS, cleanup := mockInit(t)
			defer cleanup()

			SNMP := test.prepare(mS)
			defaultMockExpects(mS)
			require.True(t, SNMP.Init())

			collectCount := 0
			for _, c := range *SNMP.charts {
				collectCount += len(c.Dims)
			}

			if test.wantFail {
				assert.False(t, SNMP.Check())
			} else {
				assert.Equal(t, collectCount, len(SNMP.Collect()))
			}
		})
	}
}

func TestSNMP_Cleanup(t *testing.T) {
	snmpC := New()
	snmpC.snmpClient = nil
	assert.NotPanics(t, snmpC.Cleanup)
}

func mockInit(t *testing.T) (*snmpmock.MockHandler, func()) {
	mockCtl := gomock.NewController(t)
	cleanup := func() { mockCtl.Finish() }
	mockSNMP := snmpmock.NewMockHandler(mockCtl)

	return mockSNMP, cleanup
}

func defaultMockExpects(m *snmpmock.MockHandler) {
	m.EXPECT().SetTarget(gomock.Any()).AnyTimes()
	m.EXPECT().SetPort(gomock.Any()).AnyTimes()
	m.EXPECT().SetMaxOids(gomock.Any()).AnyTimes()
	m.EXPECT().SetLogger(gomock.Any()).AnyTimes()
	m.EXPECT().SetTimeout(gomock.Any()).AnyTimes()
	m.EXPECT().SetCommunity(gomock.Any()).AnyTimes()
	m.EXPECT().SetVersion(gomock.Any()).AnyTimes()
	m.EXPECT().SetSecurityModel(gomock.Any()).AnyTimes()
	m.EXPECT().SetMsgFlags(gomock.Any()).AnyTimes()
	m.EXPECT().SetSecurityParameters(gomock.Any()).AnyTimes()
	m.EXPECT().Connect().Return(nil).AnyTimes()
}

func createCharts() []ChartsConfig {
	return []ChartsConfig{
		{
			ID:       "test_chart",
			Title:    "Test chart",
			Priority: 1,
			Type:     &cType,
			Family:   &cFamily,
			Dimensions: []Dimension{
				{
					Name:       "in",
					OID:        "1.3.6.1.2.1.2.2.1.10.2",
					Algorithm:  (*string)(&cAlgorithm),
					Multiplier: &cMultiplier,
					Divisor:    &cDivisor,
				},
				{
					Name:       "out",
					OID:        "1.3.6.1.2.1.2.2.1.16.2",
					Algorithm:  (*string)(&cAlgorithm),
					Multiplier: &cMultiplier,
					Divisor:    &cDivisor,
				},
			},
		},
	}
}

func prepareConfigWithoutUser() Config {
	return Config{
		Hostname:    "test",
		UpdateEvery: 2,

		Options: &Options{
			Port:    161,
			Retries: 1,
			Timeout: 2,
			Version: 3,
			MaxOIDs: 4,
		},
		ChartInput: createCharts(),
	}
}

func prepareConfigWithCommunity() Config {
	return Config{
		Hostname:    "test",
		UpdateEvery: 2,
		Options: &Options{
			Port:    161,
			Retries: 1,
			Timeout: 2,
			Version: 2, //Version 2
			MaxOIDs: 4,
		},
		Community:  &community,
		ChartInput: createCharts(),
	}
}

func prepareConfigWithoutCommunity() Config {
	return Config{
		Hostname:    "test",
		UpdateEvery: 2,
		Options: &Options{
			Port:    161,
			Retries: 1,
			Timeout: 2,
			Version: 2, //Version 2
			MaxOIDs: 4,
		},
		ChartInput: createCharts(),
	}
}

func prepareConfigWithoutChart() Config {
	return Config{
		Hostname:    "test",
		UpdateEvery: 2,
		Options: &Options{
			Port:    161,
			Retries: 1,
			Timeout: 2,
			Version: 3,
			MaxOIDs: 4,
		},
		User: &User{
			Name:          "test",
			SecurityLevel: "3",
			AuthProto:     "2",
			AuthKey:       "test_auth_key",
			PrivProto:     "2",
			PrivKey:       "test_priv_key",
		},
	}
}

func prepareConfigWithDimensions() Config {
	return Config{
		Hostname:    "test",
		UpdateEvery: 2,
		Options: &Options{
			Port:    161,
			Retries: 1,
			Timeout: 2,
			Version: 3,
			MaxOIDs: 4,
		},
		User: &User{
			Name:          "test",
			SecurityLevel: "3",
			AuthProto:     "2",
			AuthKey:       "test_auth_key",
			PrivProto:     "2",
			PrivKey:       "test_priv_key",
		},
		ChartInput: createCharts(),
	}
}

func prepareConfigWithoutDimensions() Config {
	return Config{
		Hostname:    "test",
		UpdateEvery: 2,
		Options: &Options{
			Port:    161,
			Retries: 1,
			Timeout: 2,
			Version: 3,
			MaxOIDs: 4,
		},
		User: &User{
			Name:          "test",
			SecurityLevel: "3",
			AuthProto:     "2",
			AuthKey:       "test_auth_key",
			PrivProto:     "2",
			PrivKey:       "test_priv_key",
		},
		ChartInput: []ChartsConfig{
			{
				ID:       "test_chart",
				Title:    "Test chart",
				Priority: 1,
			},
		},
	}
}

func prepareConfigWithMultiplyRange() Config {
	return Config{
		Hostname:    "test",
		UpdateEvery: 2,
		Options: &Options{
			Port:    161,
			Retries: 1,
			Timeout: 2,
			Version: 3,
			MaxOIDs: 5,
		},
		User: &User{
			Name:          "test",
			SecurityLevel: "3",
			AuthProto:     "2",
			AuthKey:       "test_auth_key",
			PrivProto:     "2",
			PrivKey:       "test_priv_key",
		},
		ChartInput: []ChartsConfig{
			{
				ID:            "test_chart",
				Title:         "Test chart",
				Priority:      1,
				Type:          &cType,
				Family:        &cFamily,
				MultiplyRange: []int{1, 5},
				Dimensions: []Dimension{
					{
						Name:       "in",
						OID:        "1.3.6.1.2.1.2.2.1.10",
						Algorithm:  (*string)(&cAlgorithm),
						Multiplier: &cMultiplier,
						Divisor:    &cDivisor,
					},
				},
			},
		},
	}
}
