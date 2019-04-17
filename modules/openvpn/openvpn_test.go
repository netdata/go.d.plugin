package openvpn

import (
	"fmt"
	"io/ioutil"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	testLoadStatsData, _    = ioutil.ReadFile("testdata/load-stats.txt")
	testVersionData, _      = ioutil.ReadFile("testdata/version.txt")
	testStatus3Data, _      = ioutil.ReadFile("testdata/status3.txt")
	testStatus3EmptyData, _ = ioutil.ReadFile("testdata/status3-empty.txt")

	testCommandStatus3Empty = "status 3 empty/n"
)

func TestNew(t *testing.T) {
	job := New()
	assert.IsType(t, (*OpenVPN)(nil), job)
	assert.Equal(t, defaultAddress, job.Address)
	assert.Equal(t, defaultConnectTimeout, job.ConnectTimeout.Duration)
	assert.Equal(t, defaultReadTimeout, job.ReadTimeout.Duration)
}

func TestOpenVPN_Init(t *testing.T) {
	job := New()
	assert.True(t, job.Init())
	assert.NotNil(t, job.apiClient)
}

func TestOpenVPN_InitNG(t *testing.T) {
	job := New()
	job.Address = ""
	assert.False(t, job.Init())
	assert.Nil(t, job.apiClient)
}

func TestOpenVPN_Check(t *testing.T) {
	job := New()
	assert.True(t, job.Init())
	job.apiClient = &mockApiClient{}
	assert.True(t, job.Check())
}

func TestOpenVPN_CheckNG(t *testing.T) {
	job := New()
	job.Address = "127.0.0.1:38001"
	assert.True(t, job.Init())
	assert.False(t, job.Check())
}

func TestOpenVPN_Charts(t *testing.T) { assert.NotNil(t, New().Charts()) }

func TestOpenVPN_Cleanup(t *testing.T) { assert.NotPanics(t, New().Cleanup) }

func TestExample_Collect(t *testing.T) {
	//mod := New()
	//
	//assert.NotNil(t, mod.Collect())
}

type mockApiClient struct{ lastCommand string }

func (mockApiClient) connect() error { return nil }

func (mockApiClient) reconnect() error { return nil }

func (mockApiClient) disconnect() error { return nil }

func (mockApiClient) isConnected() bool { return true }

func (m *mockApiClient) send(command string) error {
	m.lastCommand = command
	return nil
}

func (m *mockApiClient) read(stop func(string) bool) ([]string, error) {
	switch m.lastCommand {
	case commandLoadStats:
		return strings.Split(string(testLoadStatsData), "\n"), nil
	case commandVersion:
		return strings.Split(string(testVersionData), "\n"), nil
	case commandStatus:
		return strings.Split(string(testStatus3Data), "\n"), nil
	}
	return nil, fmt.Errorf("unknown command : %s", m.lastCommand)
}
