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
}

func TestDnsmasqDHCP_Init(t *testing.T) {
	job := New()
	job.LeasesPath = testLeasesPath
	job.ConfPath = testConfPath
	job.ConfDir = testConfDir

	assert.True(t, job.Init())
}

func TestDnsmasqDHCP_InitEmptyLeasesPath(t *testing.T) {
	job := New()
	job.LeasesPath = ""

	assert.False(t, job.Init())
}

func TestDnsmasqDHCP_InitInvalidLeasesPath(t *testing.T) {
	job := New()
	job.LeasesPath = testLeasesPath
	job.LeasesPath += "!"

	assert.False(t, job.Init())
}

func TestDnsmasqDHCP_InitZeroDHCPRanges(t *testing.T) {
	job := New()
	job.LeasesPath = testLeasesPath
	job.ConfPath = "testdata/dnsmasq3.conf"
	job.ConfDir = ""

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

func TestDnsmasqDHCP_Charts(t *testing.T) {
	job := New()
	job.LeasesPath = testLeasesPath
	job.ConfPath = testConfPath
	job.ConfDir = testConfDir

	require.True(t, job.Init())

	assert.NotNil(t, job.Charts())
}

func TestDnsmasqDHCP_Cleanup(t *testing.T) {
	assert.NotPanics(t, New().Cleanup)
}

func TestDnsmasqDHCP_Collect(t *testing.T) {
	job := New()
	job.LeasesPath = testLeasesPath
	job.ConfPath = testConfPath
	job.ConfDir = testConfDir

	require.True(t, job.Init())
	require.True(t, job.Check())

	expected := map[string]int64{
		"1230::1-1230::64":                         7,
		"1230::1-1230::64_percentage":              7,
		"1231::1-1231::64":                         1,
		"1231::1-1231::64_percentage":              1,
		"1232::1-1232::64":                         1,
		"1232::1-1232::64_percentage":              1,
		"1233::1-1233::64":                         1,
		"1233::1-1233::64_percentage":              1,
		"1234::1-1234::64":                         1,
		"1234::1-1234::64_percentage":              1,
		"192.168.0.1-192.168.0.100":                6,
		"192.168.0.1-192.168.0.100_percentage":     6,
		"192.168.1.1-192.168.1.100":                5,
		"192.168.1.1-192.168.1.100_percentage":     5,
		"192.168.2.1-192.168.2.100":                4,
		"192.168.2.1-192.168.2.100_percentage":     4,
		"192.168.3.1-192.168.3.100":                1,
		"192.168.3.1-192.168.3.100_percentage":     1,
		"192.168.4.1-192.168.4.100":                1,
		"192.168.4.1-192.168.4.100_percentage":     1,
		"192.168.200.1-192.168.200.100":            1,
		"192.168.200.1-192.168.200.100_percentage": 1,
	}

	assert.Equal(t, expected, job.Collect())
}

func TestDnsmasqDHCP_CollectFailedToOpenLeasesPath(t *testing.T) {
	job := New()
	job.LeasesPath = testLeasesPath
	job.ConfPath = testConfPath
	job.ConfDir = testConfDir

	require.True(t, job.Init())
	require.True(t, job.Check())

	job.LeasesPath = ""
	assert.Nil(t, job.Collect())
}
