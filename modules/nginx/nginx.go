package nginx

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/netdata/go.d.plugin/modules"
	"github.com/netdata/go.d.plugin/pkg/stm"
	"github.com/netdata/go.d.plugin/pkg/web"
)

func init() {
	creator := modules.Creator{
		Create: func() modules.Module { return New() },
	}

	modules.Register("nginx", creator)
}

// New creates Nginx with default values
func New() *Nginx {
	var (
		defURL         = "http://localhost/stub_status"
		defHTTPTimeout = time.Second
	)
	return &Nginx{
		HTTP: web.HTTP{
			Request: web.Request{URL: defURL},
			Client:  web.Client{Timeout: web.Duration{Duration: defHTTPTimeout}},
		},
	}
}

type metrics struct {
	Active   int `stm:"active"`
	Requests int `stm:"requests"`
	Reading  int `stm:"reading"`
	Writing  int `stm:"writing"`
	Waiting  int `stm:"waiting"`
	Accepts  int `stm:"accepts"`
	Handled  int `stm:"handled"`
}

// Nginx nginx module
type Nginx struct {
	modules.Base // should be embedded by every module

	web.HTTP `yaml:",inline"`

	request *http.Request
	client  web.HTTPClient

	metrics metrics
}

// Cleanup makes cleanup
func (Nginx) Cleanup() {}

// Init makes initialization
func (n *Nginx) Init() bool {
	var err error

	// create HTTP request
	if n.request, err = web.NewHTTPRequest(n.Request); err != nil {
		n.Errorf("error on creating request to %s : %s", n.URL, err)
		return false
	}

	// create HTTP client
	n.client = web.NewHTTPClient(n.Client)

	// post Init debug info
	n.Debugf("using URL %s", n.request.URL)
	n.Debugf("using timeout: %s", n.Timeout.Duration)

	return true
}

// Check makes check
func (n *Nginx) Check() bool {
	return len(n.Collect()) > 0
}

// Charts creates Charts
func (Nginx) Charts() *Charts {
	return charts.Copy()
}

// Collect collects metrics
func (n *Nginx) Collect() map[string]int64 {
	resp, err := n.doRequest()

	if err != nil {
		n.Error(err)
		return nil
	}

	defer func() {
		_, _ = io.Copy(ioutil.Discard, resp.Body)
		_ = resp.Body.Close()
	}()

	if resp.StatusCode != http.StatusOK {
		n.Errorf("%s returned HTTP status %d", n.request.URL.String(), resp.StatusCode)
		return nil
	}

	if err := n.parseResponse(resp); err != nil {
		n.Errorf("error on parse response : %s", err)
		return nil
	}

	return stm.ToMap(n.metrics)
}

func (n *Nginx) doRequest() (*http.Response, error) {
	return n.client.Do(n.request)
}

func (n *Nginx) parseResponse(resp *http.Response) error {
	// Active connections: 2
	//server accepts handled requests
	// 2 2 3
	//Reading: 0 Writing: 1 Waiting: 1

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	lines := strings.Split(string(b), "\n")

	if len(lines) < 4 {
		return fmt.Errorf("unparsable data, expected 4 rows, got %d", len(lines))
	}

	if err := n.parseActiveConnections(lines[0]); err != nil {
		return err
	}

	if err := n.parseAcceptsHandledRequests(lines[2]); err != nil {
		return err
	}

	if err := n.parseReadWriteWait(lines[3]); err != nil {
		return err
	}

	return nil
}

func (n *Nginx) parseActiveConnections(line string) error {
	slice := strings.Fields(line)

	if len(slice) != 3 {
		return fmt.Errorf("not enough fields in %s", line)
	}

	v, err := strconv.Atoi(slice[2])
	if err != nil {
		return err
	}
	n.metrics.Active = v

	return nil
}

func (n *Nginx) parseAcceptsHandledRequests(line string) error {
	slice := strings.Fields(line)

	if len(slice) != 3 {
		return fmt.Errorf("not enough fields in %s", line)
	}

	v, err := strconv.Atoi(slice[0])
	if err != nil {
		return err
	}
	n.metrics.Accepts = v

	v, err = strconv.Atoi(slice[1])
	if err != nil {
		return err
	}
	n.metrics.Handled = v

	v, err = strconv.Atoi(slice[2])
	if err != nil {
		return err
	}
	n.metrics.Requests = v

	return nil
}

func (n *Nginx) parseReadWriteWait(line string) error {
	slice := strings.Fields(line)

	if len(slice) != 6 {
		return fmt.Errorf("not enough fields in %s", line)
	}

	v, err := strconv.Atoi(slice[1])
	if err != nil {
		return err
	}
	n.metrics.Reading = v

	v, err = strconv.Atoi(slice[3])
	if err != nil {
		return err
	}
	n.metrics.Writing = v

	v, err = strconv.Atoi(slice[5])
	if err != nil {
		return err
	}
	n.metrics.Waiting = v

	return nil
}
