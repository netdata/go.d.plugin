package discover

import (
	"crypto/tls"
	"net/url"
	"testing"
	"time"

	"github.com/netdata/go.d.plugin/modules/vsphere/client"
	rs "github.com/netdata/go.d.plugin/modules/vsphere/resources"
	"github.com/netdata/go.d.plugin/pkg/web"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/vmware/govmomi/simulator"
)

func newTestClient(vCenterURL *url.URL) (*client.Client, error) {
	return client.New(client.Config{
		URL:             vCenterURL.String(),
		User:            "admin",
		Password:        "password",
		Timeout:         time.Second * 3,
		ClientTLSConfig: web.ClientTLSConfig{InsecureSkipVerify: true},
	})
}

func createSim() (*simulator.Model, *simulator.Server, error) {
	model := simulator.VPX()
	model.Datacenter = 2
	model.Folder = 3

	err := model.Create()
	if err != nil {
		return nil, nil, err
	}

	model.Service.TLS = new(tls.Config)

	s := model.Service.NewServer()
	return model, s, nil
}

func TestVSphereDiscoverer_Discover(t *testing.T) {
	model, srv, err := createSim()
	require.NoError(t, err)

	defer model.Remove()
	defer srv.Close()

	c, err := newTestClient(srv.URL)
	require.NoError(t, err)

	d := NewVSphereDiscoverer(c)
	res, err := d.Discover()
	require.NoError(t, err)

	assert.True(t, len(res.Dcs) > 0)
	assert.True(t, len(res.Folders) > 0)
	assert.True(t, len(res.Clusters) > 0)
	assert.True(t, len(res.Hosts) > 0)
	assert.True(t, len(res.VMs) > 0)
	assert.True(t, isHierarchySet(res))
	assert.True(t, isMetricListsCollected(res))
}

func TestVSphereDiscoverer_discover(t *testing.T) {
	model, srv, err := createSim()
	require.NoError(t, err)
	defer model.Remove()
	defer srv.Close()

	c, err := newTestClient(srv.URL)
	require.NoError(t, err)

	d := NewVSphereDiscoverer(c)
	raw, err := d.discover()
	require.NoError(t, err)

	count := model.Count()
	assert.Lenf(t, raw.dcs, count.Datacenter, "datacenters")
	assert.Lenf(t, raw.folders, count.Folder-1, "folders") // minus root folder
	dummyClusters := model.Host * count.Datacenter
	assert.Lenf(t, raw.clusters, count.Cluster+dummyClusters, "clusters")
	assert.Lenf(t, raw.hosts, count.Host, "hosts")
	assert.Lenf(t, raw.vms, count.Machine, "hosts")
}

func TestVSphereDiscoverer_build(t *testing.T) {
	model, srv, err := createSim()
	require.NoError(t, err)
	defer model.Remove()
	defer srv.Close()

	c, err := newTestClient(srv.URL)
	require.NoError(t, err)

	d := NewVSphereDiscoverer(c)
	raw, err := d.discover()
	require.NoError(t, err)

	res := d.build(raw)
	assert.Lenf(t, res.Dcs, len(raw.dcs), "datacenters")
	assert.Lenf(t, res.Folders, len(raw.folders), "folders")
	assert.Lenf(t, res.Clusters, len(raw.clusters), "clusters")
	assert.Lenf(t, res.Hosts, len(raw.hosts), "hosts")
	assert.Lenf(t, res.VMs, len(raw.vms), "hosts")
}

func TestVSphereDiscoverer_setHierarchy(t *testing.T) {
	model, srv, err := createSim()
	require.NoError(t, err)
	defer model.Remove()
	defer srv.Close()

	c, err := newTestClient(srv.URL)
	require.NoError(t, err)

	d := NewVSphereDiscoverer(c)
	raw, err := d.discover()
	require.NoError(t, err)
	res := d.build(raw)

	err = d.setHierarchy(res)
	require.NoError(t, err)
	assert.True(t, isHierarchySet(res))
}

func TestVSphereDiscoverer_removeUnmatched(t *testing.T) {
	model, srv, err := createSim()
	require.NoError(t, err)
	defer model.Remove()
	defer srv.Close()

	c, err := newTestClient(srv.URL)
	require.NoError(t, err)

	d := NewVSphereDiscoverer(c)
	d.HostMatcher = testHostMatcher{}
	d.VMMatcher = testVMMatcher{}
	raw, err := d.discover()
	require.NoError(t, err)
	res := d.build(raw)

	numVMs, numHosts := len(res.VMs), len(res.Hosts)
	assert.Equal(t, numVMs+numHosts, d.removeUnmatched(res))
	assert.Lenf(t, res.Hosts, 0, "hosts")
	assert.Lenf(t, res.VMs, 0, "vms")
}

func TestVSphereDiscoverer_collectMetricLists(t *testing.T) {
	model, srv, err := createSim()
	require.NoError(t, err)
	defer model.Remove()
	defer srv.Close()

	c, err := newTestClient(srv.URL)
	require.NoError(t, err)

	d := NewVSphereDiscoverer(c)
	raw, err := d.discover()
	require.NoError(t, err)

	res := d.build(raw)
	err = d.collectMetricLists(res)
	require.NoError(t, err)
	assert.True(t, isMetricListsCollected(res))
}

func isHierarchySet(res *rs.Resources) bool {
	for _, c := range res.Clusters {
		if !c.Hier.IsSet() {
			return false
		}
	}
	for _, h := range res.Hosts {
		if !h.Hier.IsSet() {
			return false
		}
	}
	for _, v := range res.VMs {
		if !v.Hier.IsSet() {
			return false
		}
	}
	return true
}

func isMetricListsCollected(res *rs.Resources) bool {
	for _, h := range res.Hosts {
		if h.MetricList == nil {
			return false
		}
	}
	for _, v := range res.VMs {
		if v.MetricList == nil {
			return false
		}
	}
	return true
}

type testHostMatcher struct{}

func (testHostMatcher) Match(host rs.Host) bool {
	return false
}

type testVMMatcher struct{}

func (testVMMatcher) Match(vm rs.VM) bool {
	return false
}
