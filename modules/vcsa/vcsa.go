package vcsa

import (
	"crypto/tls"
	"golang.org/x/net/proxy"
	"net/http"
	"time"

	"github.com/netdata/go.d.plugin/modules/vcsa/client"
	"github.com/netdata/go.d.plugin/pkg/web"

	"github.com/netdata/go-orchestrator/module"
)

func httpClientWithSocks5Proxy(proxyAddr string) (*http.Client, error) {
	dialer, err := proxy.SOCKS5("tcp", proxyAddr, nil, proxy.Direct)
	if err != nil {
		return nil, err
	}
	_ = dialer

	httpTransport := &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}
	httpTransport.Dial = dialer.Dial
	httpClient := &http.Client{Transport: httpTransport, Timeout: time.Second * 5}
	return httpClient, nil
}

func init() {
	creator := module.Creator{
		Create: func() module.Module { return New() },
	}

	module.Register("vcsa", creator)
}

// New creates VCenter with default values.
func New() *VCenter {
	config := Config{
		HTTP: web.HTTP{
			Request: web.Request{UserURL: "https://192.168.0.154", Username: "administrator@vsphere.local", Password: "123qwe!@#QWE"},
			Client:  web.Client{Timeout: web.Duration{Duration: time.Second * 2}},
		},
	}
	return &VCenter{
		Config: config,
	}
}

// Config is the VCenter module configuration.
type Config struct {
	web.HTTP `yaml:",inline"`
}

type healthClient interface {
	Login() error
	Logout() error
	Ping() error
	ApplMgmt() (string, error)
	DatabaseStorage() (string, error)
	Load() (string, error)
	Mem() (string, error)
	SoftwarePackages() (string, error)
	Storage() (string, error)
	Swap() (string, error)
	System() (string, error)
}

// VCenter VCenter module.
type VCenter struct {
	module.Base
	Config `yaml:",inline"`

	client healthClient
}

// Cleanup makes cleanup.
func (vc VCenter) Cleanup() {
	if vc.client == nil {
		return
	}

	err := vc.client.Logout()
	if err != nil {
		vc.Errorf("error on logout : %vc", err)
	}
}

func (vc *VCenter) createHealthClient() error {
	httpClient, err := web.NewHTTPClient(vc.Client)
	if err != nil {
		return err
	}

	httpClient, err = httpClientWithSocks5Proxy("127.0.0.1:8888")
	if err != nil {
		return err
	}

	vc.client = client.New(httpClient, vc.UserURL, vc.Username, vc.Password)
	return nil
}

// Init makes initialization.
func (vc *VCenter) Init() bool {
	err := vc.createHealthClient()
	if err != nil {
		vc.Errorf("error on creating health client : %vc", err)
		return false
	}

	vc.Debugf("using URL %s", vc.UserURL)
	vc.Debugf("using timeout: %s", vc.Timeout.Duration)
	return true
}

// Check makes check.
func (vc *VCenter) Check() bool {
	err := vc.client.Login()
	if err != nil {
		vc.Error(err)
		return false
	}
	return len(vc.Collect()) > 0
}

// Charts returns Charts.
func (vc VCenter) Charts() *module.Charts {
	return charts.Copy()
}

// Collect collects metrics.
func (vc *VCenter) Collect() map[string]int64 {
	mx, err := vc.collect()
	if err != nil {
		vc.Error(err)
	}

	if len(mx) == 0 {
		return nil
	}
	return mx
}
