package dnsmasq_dhcp

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	testLeasesPath = "testdata/dnsmasq.leases"
	testConfPath   = "testdata/dnsmasq.conf"
	testConfDir    = "testdata/dnsmasq.d"
)

func TestNew(t *testing.T) {
	job := New()

	assert.IsType(t, (*DnsmasqDHCP)(nil), job)
	assert.Equal(t, defaultLeasesPath, job.LeasesPath)
	assert.Equal(t, defaultConfPath, job.ConfPath)
	assert.Equal(t, defaultConfDir, job.ConfDir)
}

func TestDnsmasqDHCP_Init(t *testing.T) {
	job := New()
	job.LeasesPath = testLeasesPath
	job.ConfPath = testConfPath
	job.ConfDir = testConfDir

	assert.True(t, job.Init())
}

func TestDnsmasqDHCP_InitNG(t *testing.T) {
	job := New()
	job.LeasesPath += "_"
	job.ConfPath += "_"
	job.ConfDir += "_"

	assert.False(t, job.Init())
}

func TestDnsmasqDHCP_Check(t *testing.T) {
	job := New()
	job.LeasesPath = testLeasesPath
	job.ConfPath = testConfPath
	job.ConfDir = testConfDir

	require.True(t, job.Init())
	assert.True(t, job.Check())
}

func TestDnsmasqDHCP_CheckNG(t *testing.T) {
	job := New()
	job.LeasesPath = testLeasesPath
	job.ConfPath = testConfPath
	job.ConfDir = testConfDir

	require.True(t, job.Init())

	job.LeasesPath += "_"
	job.ConfPath += "_"
	job.ConfDir += "_"
	assert.False(t, job.Check())
}

func TestDnsmasqDHCP_Charts(t *testing.T) {
	job := New()
	job.LeasesPath = testLeasesPath
	job.ConfPath = testConfPath
	job.ConfDir = testConfDir

	require.True(t, job.Init())

	assert.NotNil(t, job.Charts())
}

func TestDnsmasqDHCP_Cleanup(t *testing.T) { assert.NotPanics(t, New().Cleanup) }

func TestDnsmasqDHCP_Collect(t *testing.T) {
	job := New()
	job.LeasesPath = testLeasesPath
	job.ConfPath = testConfPath
	job.ConfDir = testConfDir

	require.True(t, job.Init())
	require.True(t, job.Check())

	expected := map[string]int64{
		"1230::-1230::9":                     6,
		"1230::-1230::9_percentage":          60,
		"1231::-1231::9":                     0,
		"1231::-1231::9_percentage":          0,
		"172.168.0.0-172.168.0.9":            0,
		"172.168.0.0-172.168.0.9_percentage": 0,
		"192.168.0.0-192.168.0.9":            5,
		"192.168.0.0-192.168.0.9_percentage": 50,
		"192.168.1.0-192.168.1.9":            4,
		"192.168.1.0-192.168.1.9_percentage": 40,
		"192.168.2.0-192.168.2.9":            3,
		"192.168.2.0-192.168.2.9_percentage": 30,
	}

	assert.Equal(t, expected, job.Collect())
}

func TestConfigDir_findConfigs(t *testing.T) {
	testcases := []struct {
		config   configDir
		expected []string
	}{
		0: {
			configDir{
				path:          "testdata",
				includeSuffix: nil,
				excludeSuffix: nil,
			},
			[]string{
				"testdata/dnsmasq.conf",
				"testdata/dnsmasq.leases",
				"testdata/dnsmasq.more.conf",
			},
		},

		1: {
			configDir{
				path:          "testdata",
				includeSuffix: []string{".leases"},
				excludeSuffix: nil,
			},
			[]string{
				"testdata/dnsmasq.leases",
			},
		},

		2: {
			configDir{
				path:          "testdata",
				includeSuffix: nil,
				excludeSuffix: []string{".leases"},
			},
			[]string{
				"testdata/dnsmasq.conf",
				"testdata/dnsmasq.more.conf",
			},
		},

		3: {
			// weird one, but possible
			configDir{
				path:          "testdata",
				includeSuffix: []string{".conf", ".leases"},
				excludeSuffix: []string{".conf"},
			},
			[]string{
				"testdata/dnsmasq.leases",
			},
		},
	}

	for i, testcase := range testcases {
		actual, err := testcase.config.findConfigs()
		assert.NoError(t, err, "testcase: %d", i)

		assert.Equal(t, testcase.expected, actual, "testcase: %d", i)
	}
}
