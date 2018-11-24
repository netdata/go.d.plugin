package httpcheck

import (
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/netdata/go.d.plugin/modules"
	"github.com/netdata/go.d.plugin/pkg/utils"
	"github.com/netdata/go.d.plugin/pkg/web"
)

type state string

var (
	// success state = "success"
	timeout state = "timeout"
	failed  state = "failed"
	unknown state = "unknown"
)

type data struct {
	Success        int `stm:"success"`
	Failed         int `stm:"failed"`
	Timeout        int `stm:"timeout"`
	BadContent     int `stm:"bad_content"`
	BadStatus      int `stm:"bad_status"`
	ResponseTime   int `stm:"response_time"`
	ResponseLength int `stm:"response_length"`
}

func (d *data) reset() {
	d.Success = 0
	d.Failed = 0
	d.Timeout = 0
	d.BadContent = 0
	d.BadStatus = 0
	d.ResponseTime = 0
	d.ResponseLength = 0
}

func (d data) toMap() map[string]int64 {
	return utils.ToMap(d)
}

func New() *HttpCheck {
	return &HttpCheck{
		statuses: make(map[int]bool),
		data:     data{},
	}
}

type HttpCheck struct {
	modules.Base

	StatusAccepted []int  `yaml:"status_accepted"`
	ResponseMatch  string `yaml:"response_match"`
	web.HTTP       `yaml:",inline"`

	match    *regexp.Regexp
	statuses map[int]bool

	request *http.Request
	client  web.Client

	data data
}

func (hc *HttpCheck) Init() bool {
	if hc.Timeout.Duration == 0 {
		hc.Timeout.Duration = time.Second
	}

	hc.Debugf("Using timeout: %s", hc.Timeout.Duration)

	for _, s := range hc.StatusAccepted {
		hc.statuses[s] = true
	}

	if len(hc.statuses) == 0 {
		hc.statuses[200] = true
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

func (hc HttpCheck) Check() bool {
	return true
}

func (hc HttpCheck) GetCharts() *modules.Charts {
	c := charts.Copy()

	if len(hc.ResponseMatch) == 0 {
		c.Remove("response_check_content")
	}

	return c

}

func (hc *HttpCheck) GetData() map[string]int64 {
	hc.data.reset()
	resp, err := hc.doRequest()

	if err != nil {
		switch errCheck(err) {
		case timeout:
			hc.data.Timeout = 1
		case failed:
			hc.data.Failed = 1
		case unknown:
			hc.Error(err)
			return nil
		}
		return hc.data.toMap()
	}

	defer func() {
		io.Copy(ioutil.Discard, resp.Body)
		resp.Body.Close()
	}()

	hc.data.Success = 1
	bodyBytes, _ := ioutil.ReadAll(resp.Body)
	hc.data.ResponseLength = len(bodyBytes)

	if !hc.statuses[resp.StatusCode] {
		hc.data.BadStatus = 1
	}

	if hc.match != nil && !hc.match.Match(bodyBytes) {
		hc.data.BadContent = 1
	}

	return hc.data.toMap()
}

func (hc *HttpCheck) doRequest() (*http.Response, error) {
	t := time.Now()
	r, err := hc.client.Do(hc.request)
	hc.data.ResponseTime = int(time.Since(t))

	return r, err
}

func errCheck(err error) state {
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
