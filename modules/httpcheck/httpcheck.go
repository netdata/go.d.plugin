package httpcheck

import (
	"net/http"
	"regexp"
	"time"

	"github.com/netdata/go.d.plugin/pkg/web"

	"github.com/netdata/go-orchestrator/module"
)

func init() {
	creator := module.Creator{
		Defaults: module.Defaults{
			UpdateEvery: 5,
		},
		Create: func() module.Module { return New() },
	}

	module.Register("httpcheck", creator)
}

var (
	defaultHTTPTimeout      = time.Second
	defaultAcceptedStatuses = []int{200}
)

// New creates HTTPCheck with default values.
func New() *HTTPCheck {
	config := Config{
		HTTP: web.HTTP{
			Client: web.Client{Timeout: web.Duration{Duration: defaultHTTPTimeout}},
		},
		AcceptedStatuses: defaultAcceptedStatuses,
	}
	return &HTTPCheck{
		Config:           config,
		acceptedStatuses: make(map[int]bool),
	}
}

// Config is the HTTPCheck module configuration.
type Config struct {
	web.HTTP         `yaml:",inline"`
	AcceptedStatuses []int  `yaml:"status_accepted"`
	ResponseMatch    string `yaml:"response_match"`
}

// HTTPCheck HTTPCheck module.
type HTTPCheck struct {
	module.Base
	Config           `yaml:",inline"`
	acceptedStatuses map[int]bool
	reResponse       *regexp.Regexp
	client           *http.Client
}

// Cleanup makes cleanup.
func (HTTPCheck) Cleanup() {}

// Init makes initialization
func (hc *HTTPCheck) Init() bool {
	if err := hc.ParseUserURL(); err != nil {
		hc.Errorf("error on parsing url '%s' : %v", hc.UserURL, err)
		return false
	}

	if hc.URL.Host == "" {
		hc.Error("URL is not set")
		return false
	}

	c, err := web.NewHTTPClient(hc.Client)
	if err != nil {
		hc.Error("error on creating HTTP client : %v", err)
		return false
	}
	hc.client = c

	r, err := regexp.Compile(hc.ResponseMatch)
	if err != nil {
		hc.Errorf("error on creating regexp %s : %s", hc.ResponseMatch, err)
		return false
	}
	hc.reResponse = r

	for _, v := range hc.AcceptedStatuses {
		hc.acceptedStatuses[v] = true
	}

	// post Init debug info
	hc.Debugf("using URL %s", hc.URL)
	hc.Debugf("using HTTP timeout %s", hc.Timeout.Duration)

	hc.Debugf("using accepted HTTP statuses %s", hc.AcceptedStatuses)
	if hc.reResponse != nil {
		hc.Debugf("using response reResponse regexp %s", hc.reResponse)
	}

	return true
}

// Check makes check.
func (HTTPCheck) Check() bool { return true }

// Charts creates Charts
func (hc HTTPCheck) Charts() *Charts {
	cs := charts.Copy()
	if hc.reResponse != nil {
		_ = cs.Add(respCheckContentChart.Copy())
	}
	return cs
}

// Collect collects metrics
func (hc *HTTPCheck) Collect() map[string]int64 {
	mx, err := hc.collect()

	if err != nil {
		hc.Error(err)
	}

	if len(mx) == 0 {
		return nil
	}

	return mx
}
