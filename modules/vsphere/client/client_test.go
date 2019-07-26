package client

import (
	"crypto/tls"
	"net/url"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/vmware/govmomi/simulator"
	"github.com/vmware/govmomi/vim25/mo"
	"github.com/vmware/govmomi/vim25/types"
)

func newTestClient(vCenterURL *url.URL) (*Client, error) {
	return New(Config{
		URL:      vCenterURL.String(),
		Timeout:  time.Second,
		User:     "admin",
		Password: "password",
	})
}

func createSim() (*simulator.Model, *simulator.Server, error) {
	model := simulator.VPX()

	err := model.Create()
	if err != nil {
		return nil, nil, err
	}

	model.Service.TLS = new(tls.Config)

	s := model.Service.NewServer()
	return model, s, nil
}

func TestNew(t *testing.T) {
	model, srv, err := createSim()
	require.NoError(t, err)

	defer model.Remove()
	defer srv.Close()

	c, err := New(Config{URL: srv.URL.String(), Timeout: time.Second, User: "admin", Password: "password"})
	assert.NoError(t, err)
	assert.NotNil(t, c.Client)
	assert.NotNil(t, c.Root)
	assert.NotNil(t, c.Perf)
	assert.NotNil(t, c.Lock)
}

func TestClient_Version(t *testing.T) {
	model, srv, err := createSim()
	require.NoError(t, err)

	defer model.Remove()
	defer srv.Close()

	c, err := newTestClient(srv.URL)
	require.NoError(t, err)

	assert.NotEmpty(t, c.Version())
}

func TestClient_CounterInfoByName(t *testing.T) {
	model, srv, err := createSim()
	require.NoError(t, err)

	defer model.Remove()
	defer srv.Close()

	c, err := newTestClient(srv.URL)
	require.NoError(t, err)

	v, err := c.CounterInfoByName()
	assert.NoError(t, err)
	assert.IsType(t, map[string]*types.PerfCounterInfo{}, v)
	assert.NotEmpty(t, v)
}

func TestClient_IsSessionActive(t *testing.T) {
	model, srv, err := createSim()
	require.NoError(t, err)

	defer model.Remove()
	defer srv.Close()

	c, err := newTestClient(srv.URL)
	require.NoError(t, err)

	v, err := c.IsSessionActive()
	assert.NoError(t, err)
	assert.True(t, v)
}

func TestClient_Logout(t *testing.T) {
	model, srv, err := createSim()
	require.NoError(t, err)

	defer model.Remove()
	defer srv.Close()

	c, err := newTestClient(srv.URL)
	require.NoError(t, err)

	v, err := c.IsSessionActive()
	assert.NoError(t, err)
	assert.True(t, v)

	err = c.Logout()
	assert.NoError(t, err)

	v, err = c.IsSessionActive()
	assert.NoError(t, err)
	assert.False(t, v)
}

func TestClient_Login(t *testing.T) {
	model, srv, err := createSim()
	require.NoError(t, err)

	defer model.Remove()
	defer srv.Close()

	c, err := newTestClient(srv.URL)
	require.NoError(t, err)

	err = c.Logout()
	assert.NoError(t, err)

	v, err := c.IsSessionActive()
	assert.NoError(t, err)
	assert.False(t, v)

	err = c.Login()
	assert.NoError(t, err)

	v, err = c.IsSessionActive()
	assert.NoError(t, err)
	assert.True(t, v)
}

func TestClient_Datacenters(t *testing.T) {
	model, srv, err := createSim()
	require.NoError(t, err)

	defer model.Remove()
	defer srv.Close()

	c, err := newTestClient(srv.URL)
	require.NoError(t, err)

	dcs, err := c.Datacenters([]string{})
	assert.NoError(t, err)
	assert.IsType(t, []mo.Datacenter{}, dcs)
	assert.NotEmpty(t, dcs)
}

func TestClient_Folders(t *testing.T) {
	model, srv, err := createSim()
	require.NoError(t, err)

	defer model.Remove()
	defer srv.Close()

	c, err := newTestClient(srv.URL)
	require.NoError(t, err)

	folders, err := c.Folders([]string{})
	assert.NoError(t, err)
	assert.IsType(t, []mo.Folder{}, folders)
	assert.NotEmpty(t, folders)
}

func TestClient_ComputeResources(t *testing.T) {
	model, srv, err := createSim()
	require.NoError(t, err)

	defer model.Remove()
	defer srv.Close()

	c, err := newTestClient(srv.URL)
	require.NoError(t, err)

	computes, err := c.ComputeResources([]string{})
	assert.NoError(t, err)
	assert.IsType(t, []mo.ComputeResource{}, computes)
	assert.NotEmpty(t, computes)
}

func TestClient_Hosts(t *testing.T) {
	model, srv, err := createSim()
	require.NoError(t, err)

	defer model.Remove()
	defer srv.Close()

	c, err := newTestClient(srv.URL)
	require.NoError(t, err)

	hosts, err := c.Hosts([]string{})
	assert.NoError(t, err)
	assert.IsType(t, []mo.HostSystem{}, hosts)
	assert.NotEmpty(t, hosts)
}

func TestClient_VirtualMachines(t *testing.T) {
	model, srv, err := createSim()
	require.NoError(t, err)

	defer model.Remove()
	defer srv.Close()

	c, err := newTestClient(srv.URL)
	require.NoError(t, err)

	vms, err := c.VirtualMachines([]string{})
	assert.NoError(t, err)
	assert.IsType(t, []mo.VirtualMachine{}, vms)
	assert.NotEmpty(t, vms)
}

func TestClient_Reconnect(t *testing.T) {
	model, srv, err := createSim()
	require.NoError(t, err)

	defer model.Remove()
	defer srv.Close()

	c, err := newTestClient(srv.URL)
	require.NoError(t, err)

	cl := c.Client
	root := c.Root
	perf := c.Perf

	err = c.Reconnect()
	assert.NoError(t, err)

	assert.False(t, cl == c.Client)
	assert.False(t, root == c.Root)
	assert.False(t, perf == c.Perf)
}
