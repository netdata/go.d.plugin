package vcsa

import (
	"crypto/tls"
	"errors"
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
		Defaults: module.Defaults{
			// it seems health checks freq is 5 seconds, at least this is true for Overall Health according
			// Last checked info on the dashboard (:5480)
			UpdateEvery: 5,
		},
		Create: func() module.Module { return New() },
	}

	module.Register("vcsa", creator)
}

// New creates VCSA with default values.
func New() *VCSA {
	config := Config{
		HTTP: web.HTTP{
			Request: web.Request{UserURL: "https://192.168.0.154", Username: "administrator@vsphere.local", Password: "123qwe!@#QWE"},
			Client:  web.Client{Timeout: web.Duration{Duration: time.Second * 5}},
		},
	}
	return &VCSA{
		Config: config,
	}
}

// Config is the VCSA module configuration.
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

// VCSA VCSA module.
type VCSA struct {
	module.Base
	Config `yaml:",inline"`

	client healthClient
}

// Cleanup makes cleanup.
func (vc VCSA) Cleanup() {
	if vc.client == nil {
		return
	}
	err := vc.client.Logout()
	if err != nil {
		vc.Errorf("error on logout : %v", err)
	}
}

func (vc VCSA) validateInitParameters() error {
	if vc.UserURL == "" {
		return errors.New("URL not set")
	}
	if vc.Username == "" || vc.Password == "" {
		return errors.New("username or password not set")
	}
	return nil
}

func (vc *VCSA) createHealthClient() error {
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
func (vc *VCSA) Init() bool {
	err := vc.validateInitParameters()
	if err != nil {
		vc.Error(err)
		return false
	}

	err = vc.createHealthClient()
	if err != nil {
		vc.Errorf("error on creating health client : %vc", err)
		return false
	}

	vc.Debugf("using URL %s", vc.UserURL)
	vc.Debugf("using timeout: %s", vc.Timeout.Duration)
	return true
}

// Check makes check.
func (vc *VCSA) Check() bool {
	err := vc.client.Login()
	if err != nil {
		vc.Error(err)
		return false
	}
	return len(vc.Collect()) > 0
}

// Charts returns Charts.
func (vc VCSA) Charts() *module.Charts {
	return charts.Copy()
}

// Collect collects metrics.
func (vc *VCSA) Collect() map[string]int64 {
	mx, err := vc.collect()
	if err != nil {
		vc.Error(err)
	}

	if len(mx) == 0 {
		return nil
	}
	return mx
}
