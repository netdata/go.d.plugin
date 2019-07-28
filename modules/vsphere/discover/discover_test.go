package discover

import (
	"crypto/tls"
	"fmt"
	"net/url"
	"testing"
	"time"

	"github.com/netdata/go.d.plugin/modules/vsphere/client"
	"github.com/netdata/go.d.plugin/modules/vsphere/resources"
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

	err := model.Create()
	if err != nil {
		return nil, nil, err
	}

	model.Service.TLS = new(tls.Config)

	s := model.Service.NewServer()
	return model, s, nil
}

func numOfFolders(m *simulator.Model) int {
	if m.Datacenter >= m.Folder {
		return (m.Datacenter-m.Folder)*4 + m.Folder*9
	}
	return m.Datacenter * 9
}

func numOfClusters(m *simulator.Model) int {
	return m.Cluster*m.Datacenter + m.Host
}

func numOfHosts(m *simulator.Model) int {
	return m.Host*m.Datacenter + m.ClusterHost*m.Cluster
}

func numOfVMs(m *simulator.Model) int {
	return m.Host*m.Machine + m.Cluster*m.Machine + m.Cluster*m.App*m.Machine
}

func TestVSphereDiscoverer_Discover(t *testing.T) {
	model, srv, err := createSim()
	require.NoError(t, err)

	defer model.Remove()
	defer srv.Close()

	c, err := newTestClient(srv.URL)
	require.NoError(t, err)

	d := NewVSphereDiscoverer(c)

	raw, err := d.discoverRawResources()
	require.NoError(t, err)

	assert.Equalf(t, model.Datacenter, len(raw.dcs), "discover raw resources, datacenters")
	assert.Equalf(t, numOfFolders(model), len(raw.folders), "discover raw resources, folders")
	assert.Equalf(t, numOfClusters(model), len(raw.clusters), "discover raw resources, clusters")
	assert.Equalf(t, numOfHosts(model), len(raw.hosts), "discover raw resources, hosts")
	assert.Equalf(t, numOfVMs(model), len(raw.vms), "discover raw resources, vms")

	res := d.buildResources(raw)
	assert.Equalf(t, len(raw.dcs), len(res.Dcs), "build resources, datacenters")
	assert.Equalf(t, len(raw.folders), len(res.Folders), "build resources, folders")
	assert.Equalf(t, len(raw.clusters), len(res.Clusters), "build resources, clusters")
	assert.Equalf(t, len(raw.hosts), len(res.Hosts), "build resources, hosts")
	assert.Equalf(t, len(raw.vms), len(res.VMs), "build resources, vms")

	assert.Zero(t, d.removeUnmatched(res))

	assert.NoErrorf(t, d.setHierarchy(res), "set hierarchy")
	assert.NoErrorf(t, checkHierarchy(res), "set hierarchy")

	assert.NoErrorf(t, d.collectMetricLists(res), "collect metric lists")
	assert.NoErrorf(t, checkMetricLists(res), "collect metric lists")
}

func checkHierarchy(res *resources.Resources) error {
	for _, c := range res.Clusters {
		if c.Hier.IsSet() {
			continue
		}
		return fmt.Errorf("hierarchy not set for cluster %s/%s", c.ID, c.Name)
	}
	for _, h := range res.Hosts {
		if h.Hier.IsSet() {
			continue
		}
		return fmt.Errorf("hierarchy not set for host %s/%s", h.ID, h.Name)
	}
	for _, v := range res.VMs {
		if v.Hier.IsSet() {
			continue
		}
		return fmt.Errorf("hierarchy not set for vm %s/%s", v.ID, v.Name)
	}
	return nil
}

func checkMetricLists(res *resources.Resources) error {
	for _, h := range res.Hosts {
		if h.MetricList != nil {
			continue
		}
		return fmt.Errorf("metric list not set for host %s/%s", h.ID, h.Name)
	}
	for _, v := range res.VMs {
		if v.MetricList != nil {
			continue
		}
		return fmt.Errorf("metric list not set for vm %s/%s", v.ID, v.Name)
	}
	return nil
}
