package vcsa

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"sync"

	"github.com/netdata/go.d.plugin/pkg/web"
)

// https://vmware.github.io/vsphere-automation-sdk-rest/vsphere/index.html
const (
	pathCISSession             = "/rest/com/vmware/cis/session"
	pathHealthSystem           = "/rest/appliance/health/system"
	pathHealthSwap             = "/rest/appliance/health/swap"
	pathHealthStorage          = "/rest/appliance/health/storage"
	pathHealthSoftwarePackager = "/rest/appliance/health/software-packages"
	pathHealthMem              = "/rest/appliance/health/mem"
	pathHealthLoad             = "/rest/appliance/health/load"
	pathHealthDatabaseStorage  = "/rest/appliance/health/database-storage"
	pathHealthApplMgmt         = "/rest/appliance/health/applmgmt"
)

func newClient(httpClient *http.Client, url, username, password string) *client {
	return &client{
		httpClient: httpClient,
		url:        url,
		username:   username,
		password:   password,
		lock:       new(sync.RWMutex),
	}
}

type client struct {
	httpClient *http.Client

	url      string
	username string
	password string

	lock   *sync.RWMutex
	sessID string
}

func (c *client) setSessID(v string) {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.sessID = v
}

func (c *client) getSessID() string {
	c.lock.RLock()
	defer c.lock.RUnlock()
	return c.sessID
}

func (c *client) doOK(req web.Request) (*http.Response, error) {
	httpReq, err := web.NewHTTPRequest(req)
	if err != nil {
		return nil, fmt.Errorf("error on creating http request to %s : %v", req.UserURL, err)
	}

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return resp, fmt.Errorf("%s returned %d", httpReq.URL, resp.StatusCode)
	}
	return resp, nil
}

func (c *client) doOKWithDecode(req web.Request, dst interface{}) error {
	resp, err := c.doOK(req)
	if err != nil {
		return err
	}
	defer closeBody(resp)

	err = json.NewDecoder(resp.Body).Decode(dst)
	if err != nil {
		return fmt.Errorf("error on decoding response from %s : %v", req.UserURL, err)
	}
	return nil
}

func (c *client) Login() error {
	req := web.Request{
		UserURL:  fmt.Sprintf("%s%s", c.url, pathCISSession),
		Username: c.username,
		Password: c.password,
		Method:   http.MethodPost,
	}

	s := struct{ Value string }{}
	err := c.doOKWithDecode(req, &s)
	if err == nil {
		c.setSessID(s.Value)
	}
	return err
}

func (c *client) Logout() error {
	req := web.Request{
		UserURL: fmt.Sprintf("%s%s", c.url, pathCISSession),
		Method:  http.MethodDelete,
		Headers: map[string]string{"vmware-api-session-id": c.getSessID()},
	}

	resp, err := c.doOK(req)
	closeBody(resp)
	c.setSessID("")
	return err
}

func (c *client) Ping() error {
	req := web.Request{
		UserURL: fmt.Sprintf("%s%s?~action=get", c.url, pathCISSession),
		Method:  http.MethodPost,
		Headers: map[string]string{"vmware-api-session-id": c.getSessID()},
	}

	resp, err := c.doOK(req)
	defer closeBody(resp)
	if resp != nil && resp.StatusCode == http.StatusUnauthorized {
		return c.Login()
	}
	return err
}

func (c *client) health(urlPath string) (string, error) {
	req := web.Request{
		UserURL: fmt.Sprintf("%s%s", c.url, urlPath),
		Headers: map[string]string{"vmware-api-session-id": c.getSessID()},
	}

	s := struct{ Value string }{}
	err := c.doOKWithDecode(req, &s)
	return s.Value, err
}

func (c *client) ApplMgmt() (string, error) {
	return c.health(pathHealthApplMgmt)
}

func (c *client) DatabaseStorage() (string, error) {
	return c.health(pathHealthDatabaseStorage)
}

func (c *client) Load() (string, error) {
	return c.health(pathHealthLoad)
}

func (c *client) Mem() (string, error) {
	return c.health(pathHealthMem)
}

func (c *client) SoftwarePackages() (string, error) {
	return c.health(pathHealthSoftwarePackager)
}

func (c *client) Storage() (string, error) {
	return c.health(pathHealthStorage)
}

func (c *client) Swap() (string, error) {
	return c.health(pathHealthSwap)
}

func (c *client) System() (string, error) {
	return c.health(pathHealthSystem)
}

func closeBody(resp *http.Response) {
	if resp != nil && resp.Body != nil {
		_, _ = io.Copy(ioutil.Discard, resp.Body)
		_ = resp.Body.Close()
	}
}
