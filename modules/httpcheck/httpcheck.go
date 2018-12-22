package httpcheck

import (
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/netdata/go.d.plugin/pkg/stm"

	"github.com/netdata/go.d.plugin/modules"
	"github.com/netdata/go.d.plugin/pkg/stm"
	"github.com/netdata/go.d.plugin/pkg/web"
)

func init() {
	creator := modules.Creator{
		UpdateEvery: 5,
		Create:      func() modules.Module { return New() },
	}

	modules.Register("httpcheck", creator)
}

// New creates HTTPCheck with default values
func New() *HTTPCheck {
	var (
		defHTTPTimeout    = time.Second
		defStatusAccepted = map[int]bool{200: true}
	)

	return &HTTPCheck{
		HTTP: web.HTTP{
			Client: web.Client{Timeout: web.Duration{Duration: defHTTPTimeout}},
		},
		statuses: defStatusAccepted,
		metrics:  metrics{},
	}
}

type state string

var (
	timeout state = "timeout"
	failed  state = "failed"
	unknown state = "unknown"
)

type metrics struct {
	Success        int `stm:"success"`
	Failed         int `stm:"failed"`
	Timeout        int `stm:"timeout"`
	BadContent     int `stm:"bad_content"`
	BadStatus      int `stm:"bad_status"`
	ResponseTime   int `stm:"response_time"`
	ResponseLength int `stm:"response_length"`
}

func (d *metrics) reset() {
	d.Success = 0
	d.Failed = 0
	d.Timeout = 0
	d.BadContent = 0
	d.BadStatus = 0
	d.ResponseTime = 0
	d.ResponseLength = 0
}

<<<<<<< HEAD
func (d metrics) toMap() map[string]int64 {
	return stm.ToMap(d)
}

// New creates HTTPCheck with default values
func New() *HTTPCheck {
	return &HTTPCheck{
		statuses: map[int]bool{
			200: true,
		},
		metrics: metrics{},
	}
}

=======
>>>>>>> master
// HTTPCheck httpcheck module
type HTTPCheck struct {
	modules.Base

	web.HTTP `yaml:",inline"`

	StatusAccepted []int  `yaml:"status_accepted"`
	ResponseMatch  string `yaml:"response_match"`

	match    *regexp.Regexp
	statuses map[int]bool

	request *http.Request
	client  web.HTTPClient

	metrics metrics
}

// Cleanup makes cleanup
func (HTTPCheck) Cleanup() {}

// Init makes initialization
func (hc *HTTPCheck) Init() bool {
	// populate accepted statuses
	if len(hc.StatusAccepted) != 0 {
		delete(hc.statuses, 200)

		for _, s := range hc.StatusAccepted {
			hc.statuses[s] = true
		}
	}

	var err error

	// create HTTP request
	if hc.request, err = web.NewHTTPRequest(hc.Request); err != nil {
		hc.Errorf("error on creating request to %s : %s", hc.URL, err)
		return false
	}

	// create HTTP client
	hc.client = web.NewHTTPClient(hc.Client)

	// create response match
	if hc.match, err = regexp.Compile(hc.ResponseMatch); err != nil {
		hc.Errorf("error on creating regexp %s : %s", hc.ResponseMatch, err)
		return false
	}

	// post Init debug info
	hc.Debugf("using URL %s", hc.request.URL)
	hc.Debugf("using HTTP timeout %s", hc.Timeout.Duration)
	var statuses []int
	for status := range hc.statuses {
		statuses = append(statuses, status)
	}
	sort.Ints(statuses)
	hc.Debugf("using accepted HTTP statuses %s", statuses)
	if hc.match != nil {
		hc.Debugf("using response match regexp %s", hc.match)
	}

	return true
}

// Check makes check
func (hc HTTPCheck) Check() bool {
	return true
}

// Charts creates Charts
func (hc HTTPCheck) Charts() *Charts {
	c := charts.Copy()

	if len(hc.ResponseMatch) == 0 {
		c.Remove("response_check_content")
	}

	return c

}

// Collect collects metrics
func (hc *HTTPCheck) Collect() map[string]int64 {
	hc.metrics.reset()

	resp, err := hc.doRequest()

	if err != nil {
		hc.processErrResponse(err)
	} else {
		hc.processOKResponse(resp)
	}

	return stm.ToMap(hc.metrics)
}

func (hc *HTTPCheck) doRequest() (*http.Response, error) {
	t := time.Now()
	r, err := hc.client.Do(hc.request)
	hc.metrics.ResponseTime = int(time.Since(t))

	return r, err
}

func (hc *HTTPCheck) processErrResponse(err error) {
	switch parseErr(err) {
	case timeout:
		hc.metrics.Timeout = 1
	case failed:
		hc.metrics.Failed = 1
	case unknown:
		hc.Error(err)
		panic("unknown state")
	}
}

func (hc *HTTPCheck) processOKResponse(resp *http.Response) {
	defer func() {
		_, _ = io.Copy(ioutil.Discard, resp.Body)
		_ = resp.Body.Close()
	}()

	hc.metrics.Success = 1
	bodyBytes, _ := ioutil.ReadAll(resp.Body)
	hc.metrics.ResponseLength = len(bodyBytes)

	if !hc.statuses[resp.StatusCode] {
		hc.metrics.BadStatus = 1
	}

	if hc.match != nil && !hc.match.Match(bodyBytes) {
		hc.metrics.BadContent = 1
	}
}

func parseErr(err error) state {
	v, ok := err.(net.Error)

	if ok && v.Timeout() {
		return timeout
	}

	if ok && strings.Contains(v.Error(), "connection refused") {
		return failed
	}

	return unknown
}
