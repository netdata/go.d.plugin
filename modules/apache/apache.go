package apache

import (
	"io"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/netdata/go.d.plugin/modules"
	"github.com/netdata/go.d.plugin/pkg/web"
)

// -- Extended On --
// Total Accesses: 7
// Total kBytes: 5
// Uptime: 6
// ReqPerSec: 1.16667
// BytesPerSec: 853.333
// BytesPerReq: 731.429
// BusyWorkers: 1
// IdleWorkers: 49
// ConnsTotal: 1
// ConnsAsyncWriting: 0
// ConnsAsyncKeepAlive: 1
// ConnsAsyncClosing: 0

// -- Extended Off --
// BusyWorkers: 1
// IdleWorkers: 49
// ConnsTotal: 1
// ConnsAsyncWriting: 0
// ConnsAsyncKeepAlive: 1
// ConnsAsyncClosing: 0

func init() {
	creator := modules.Creator{
		Create: func() modules.Module { return New() },
	}

	modules.Register("apache", creator)
}

// New creates Apache with default values
func New() *Apache {
	return &Apache{}
}

// Apache apache module
type Apache struct {
	modules.Base // should be embedded by every module

	web.HTTP `yaml:",inline"`

	request *http.Request
	client  web.Client
}

func (Apache) Cleanup() {

}

func (a *Apache) Init() bool {
	req, err := a.CreateHTTPRequest()

	if err != nil {
		a.Error(err)
		return false
	}

	if a.Timeout.Duration == 0 {
		a.Timeout.Duration = time.Second
	}

	a.request = req
	a.client = a.CreateHTTPClient()

	return true
}

func (Apache) Check() bool {
	return false
}

func (Apache) Charts() *modules.Charts {
	return nil
}

func (a *Apache) GatherMetrics() map[string]int64 {
	resp, err := a.doRequest()

	if err != nil {
		a.Error(err)
		return nil
	}

	defer func() {
		_, _ = io.Copy(ioutil.Discard, resp.Body)
		_ = resp.Body.Close()
	}()

	return nil
}

func (a *Apache) doRequest() (*http.Response, error) {
	return a.client.Do(a.request)
}
