package client

import (
	"crypto/tls"
	"net/url"
	"testing"
	"time"

	"github.com/netdata/go.d.plugin/pkg/web"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/vmware/govmomi/simulator"
	"github.com/vmware/govmomi/vim25/mo"
	"github.com/vmware/govmomi/vim25/types"
)

func newTestClient(vCenterURL *url.URL) (*Client, error) {
	return New(Config{
		URL:             vCenterURL.String(),
		User:            "admin",
		Password:        "password",
		Timeout:         time.Second,
		ClientTLSConfig: web.ClientTLSConfig{InsecureSkipVerify: true},
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

	c, err := newTestClient(srv.URL)
	require.NoError(t, err)

	assert.NotNil(t, c.client)
	assert.NotNil(t, c.root)
	assert.NotNil(t, c.perf)
	v, err := c.IsSessionActive()
	assert.NoError(t, err)
	assert.True(t, v)
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

	err = c.Login(url.UserPassword("admin", "password"))
	assert.NoError(t, err)

	v, err = c.IsSessionActive()
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

func TestClient_Datacenters(t *testing.T) {
	model, srv, err := createSim()
	require.NoError(t, err)

	defer model.Remove()
	defer srv.Close()

	c, err := newTestClient(srv.URL)
	require.NoError(t, err)

	dcs, err := c.Datacenters()
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

	folders, err := c.Folders()
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

	computes, err := c.ComputeResources()
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

	hosts, err := c.Hosts()
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

	vms, err := c.VirtualMachines()
	assert.NoError(t, err)
	assert.IsType(t, []mo.VirtualMachine{}, vms)
	assert.NotEmpty(t, vms)
}
