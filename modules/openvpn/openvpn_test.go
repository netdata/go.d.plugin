package openvpn

import (
	"fmt"
	"io/ioutil"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	loadStatsData, _    = ioutil.ReadFile("testdata/load-stats.txt")
	versionData, _      = ioutil.ReadFile("testdata/version.txt")
	status3Data, _      = ioutil.ReadFile("testdata/status3.txt")
	status3EmptyData, _ = ioutil.ReadFile("testdata/status3-empty.txt")
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

func TestOpenVPN_Check(t *testing.T) {
	//mod := New()
	//
	//assert.True(t, mod.Check())
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
		return strings.Split(string(loadStatsData), "\n"), nil
	case commandVersion:
		return strings.Split(string(versionData), "\n"), nil
	case commandStatus:
		return strings.Split(string(status3Data), "\n"), nil
	}
	return nil, fmt.Errorf("unknown command : %s", m.lastCommand)
}
