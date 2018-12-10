package httpcheck

import (
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/netdata/go.d.plugin/pkg/stm"

	"github.com/netdata/go.d.plugin/modules"
	"github.com/netdata/go.d.plugin/pkg/web"
)

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

// HTTPCheck httpcheck module
type HTTPCheck struct {
	modules.Base

	StatusAccepted []int  `yaml:"status_accepted"`
	ResponseMatch  string `yaml:"response_match"`
	web.HTTP       `yaml:",inline"`

	match    *regexp.Regexp
	statuses map[int]bool

	request *http.Request
	client  web.Client

	metrics metrics
}

// Cleanup makes cleanup
func (HTTPCheck) Cleanup() {}

// Init makes initialization
func (hc *HTTPCheck) Init() bool {
	if hc.Timeout.Duration == 0 {
		hc.Timeout.Duration = time.Second
	}

	hc.Debugf("using timeout: %s", hc.Timeout.Duration)

	if len(hc.StatusAccepted) != 0 {
		delete(hc.statuses, 200)
		for _, s := range hc.StatusAccepted {
			hc.statuses[s] = true
		}
	}

	req, err := hc.CreateHTTPRequest()

	if err != nil {
		hc.Error(err)
		return false
	}

	hc.request = req

	hc.client = hc.CreateHTTPClient()

	re, err := regexp.Compile(hc.ResponseMatch)

	if err != nil {
		hc.Error(err)
		return false
	}

	hc.match = re

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

// GatherMetrics gathers metrics
func (hc *HTTPCheck) GatherMetrics() map[string]int64 {
	hc.metrics.reset()

	resp, err := hc.doRequest()

	if err != nil {
		hc.processErrResponse(err)
	} else {
		hc.processOKResponse(resp)
	}

	return hc.metrics.toMap()
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

func init() {
	creator := modules.Creator{
		UpdateEvery: 5,
		Create:      func() modules.Module { return New() },
	}

	modules.Register("httpcheck", creator)
}
