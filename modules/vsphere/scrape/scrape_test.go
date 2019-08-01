package scrape

import (
	"crypto/tls"
	"net/url"
	"testing"
	"time"

	"github.com/netdata/go.d.plugin/modules/vsphere/client"
	"github.com/netdata/go.d.plugin/modules/vsphere/discover"
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

func TestNewVSphereMetricScraper(t *testing.T) {

}

func TestVSphereMetricScraper_ScrapeHostsMetrics(t *testing.T) {
	model, srv, err := createSim()
	require.NoError(t, err)
	defer model.Remove()
	defer srv.Close()

	c, err := newTestClient(srv.URL)
	require.NoError(t, err)

	d := discover.NewVSphereDiscoverer(c)
	res, err := d.Discover()
	require.NoError(t, err)

	mc := NewVSphereMetricScraper(c)
	metrics := mc.ScrapeHostsMetrics(res.Hosts)
	assert.Len(t, metrics, len(res.Hosts))
}

func TestVSphereMetricScraper_ScrapeVMsMetrics(t *testing.T) {
	model, srv, err := createSim()
	require.NoError(t, err)
	defer model.Remove()
	defer srv.Close()

	c, err := newTestClient(srv.URL)
	require.NoError(t, err)

	d := discover.NewVSphereDiscoverer(c)
	res, err := d.Discover()
	require.NoError(t, err)

	mc := NewVSphereMetricScraper(c)
	metrics := mc.ScrapeVMsMetrics(res.VMs)
	assert.Len(t, metrics, len(res.VMs))
}
