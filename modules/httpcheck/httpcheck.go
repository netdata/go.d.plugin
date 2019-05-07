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
	"github.com/netdata/go.d.plugin/pkg/web"

	"github.com/netdata/go-orchestrator/module"
)

func init() {
	creator := module.Creator{
		Defaults: module.Defaults{
			Disabled:    true,
			UpdateEvery: 5,
		},
		Create: func() module.Module { return New() },
	}

	module.Register("httpcheck", creator)
}

var (
	defaultHTTPTimeout = time.Second
)

// New creates HTTPCheck with default values
func New() *HTTPCheck {
	config := Config{
		HTTP: web.HTTP{
			Client: web.Client{Timeout: web.Duration{Duration: defaultHTTPTimeout}},
		},
	}

	return &HTTPCheck{
		Config:   config,
		statuses: map[int]bool{200: true},
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

// Config is the HTTPCheck module configuration.
type Config struct {
	web.HTTP `yaml:",inline"`

	StatusAccepted []int  `yaml:"status_accepted"`
	ResponseMatch  string `yaml:"response_match"`
}

// HTTPCheck HTTPCheck module
type HTTPCheck struct {
	module.Base
	Config `yaml:",inline"`

	match    *regexp.Regexp
	statuses map[int]bool
	client   *http.Client
	metrics  metrics
}

// Cleanup makes cleanup
func (HTTPCheck) Cleanup() {}

// Init makes initialization
func (hc *HTTPCheck) Init() bool {
	var err error

	if err = hc.ParseUserURL(); err != nil {
		hc.Errorf("error on parsing url '%s' : %v", hc.UserURL, err)
		return false
	}

	if hc.URL.Host == "" {
		hc.Error("URL is not set")
		return false
	}

	// create HTTP client
	if hc.client, err = web.NewHTTPClient(hc.Client); err != nil {
		hc.Error(err)
		return false
	}

	// create response match
	if hc.match, err = regexp.Compile(hc.ResponseMatch); err != nil {
		hc.Errorf("error on creating regexp %s : %s", hc.ResponseMatch, err)
		return false
	}

	// post Init debug info
	hc.Debugf("using URL %s", hc.URL)
	hc.Debugf("using HTTP timeout %s", hc.Timeout.Duration)

	// populate accepted statuses
	if len(hc.StatusAccepted) != 0 {
		delete(hc.statuses, 200)

		for _, s := range hc.StatusAccepted {
			hc.statuses[s] = true
		}
	}

	hc.Debugf("using accepted HTTP statuses %s", hc.statuses)
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
		_ = c.Remove("response_check_content")
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

	req, err := web.NewHTTPRequest(hc.Request)
	if err != nil {
		return nil, err
	}

	r, err := hc.client.Do(req)
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
