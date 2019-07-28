package vsphere

import (
	"crypto/tls"
	"net/url"
	"testing"
	"time"

	"github.com/netdata/go.d.plugin/modules/vsphere/client"
	"github.com/netdata/go.d.plugin/pkg/web"

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

func TestNew(t *testing.T) {

}

func TestVSphere_Init(t *testing.T) {

}

func TestVSphere_Check(t *testing.T) {

}

func TestVSphere_Charts(t *testing.T) {

}

func TestVSphere_Cleanup(t *testing.T) {

}

func TestVSphere_Collect(t *testing.T) {

}
