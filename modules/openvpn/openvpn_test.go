package openvpn

import (
	"testing"

	"github.com/netdata/go.d.plugin/modules/openvpn/client"
	"github.com/netdata/go.d.plugin/pkg/matcher"

	"github.com/netdata/go-orchestrator/module"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	testVersion   = client.Version{Major: 1, Minor: 1, Patch: 1, Management: 1}
	testLoadStats = client.LoadStats{NumOfClients: 1, BytesIn: 1, BytesOut: 1}
	testUsers     = client.Users{{
		CommonName:     "common_name",
		RealAddress:    "1.2.3.4:4321",
		VirtualAddress: "1.2.3.4",
		BytesReceived:  1,
		BytesSent:      2,
		ConnectedSince: 3,
		Username:       "name",
	}}
)

func TestNew(t *testing.T) {
	job := New()

	assert.Implements(t, (*module.Module)(nil), job)
	assert.Equal(t, defaultAddress, job.Address)
	assert.Equal(t, defaultConnectTimeout, job.ConnectTimeout.Duration)
	assert.Equal(t, defaultReadTimeout, job.ReadTimeout.Duration)
	assert.Equal(t, defaultWriteTimeout, job.WriteTimeout.Duration)
	assert.NotNil(t, job.charts)
	assert.NotNil(t, job.collectedUsers)
}

func TestOpenVPN_Init(t *testing.T) { assert.True(t, New().Init()) }

func TestOpenVPN_Check(t *testing.T) {
	job := New()

	require.True(t, job.Init())
	job.apiClient = &mockOKOpenVPNClient{}
	require.True(t, job.Check())
}

func TestOpenVPN_Charts(t *testing.T) { assert.NotNil(t, New().Charts()) }

func TestOpenVPN_Cleanup(t *testing.T) {
	job := New()

	assert.NotPanics(t, job.Cleanup)
	require.True(t, job.Init())
	job.apiClient = &mockOKOpenVPNClient{}
	require.True(t, job.Check())
	job.Cleanup()
	assert.False(t, job.apiClient.IsConnected())
}

func TestOpenVPN_Collect(t *testing.T) {
	job := New()

	require.True(t, job.Init())
	job.perUserMatcher = matcher.TRUE()
	job.apiClient = &mockOKOpenVPNClient{}
	require.True(t, job.Check())

	expected := map[string]int64{
		"bytes_in":            1,
		"bytes_out":           1,
		"clients":             1,
		"name_bytes_received": 1,
		"name_bytes_sent":     2,
	}

	mx := job.Collect()
	require.NotNil(t, mx)
	delete(mx, "name_connection_time")
	assert.Equal(t, expected, mx)
}

type mockOKOpenVPNClient struct {
	isConnected bool
}

func (m *mockOKOpenVPNClient) Connect() error {
	m.isConnected = true
	return nil
}

func (m *mockOKOpenVPNClient) Disconnect() error {
	m.isConnected = false
	return nil
}

func (m mockOKOpenVPNClient) IsConnected() bool { return m.isConnected }

func (mockOKOpenVPNClient) GetVersion() (*client.Version, error) { return &testVersion, nil }

func (mockOKOpenVPNClient) GetLoadStats() (*client.LoadStats, error) { return &testLoadStats, nil }

func (mockOKOpenVPNClient) GetUsers() (client.Users, error) { return testUsers, nil }
