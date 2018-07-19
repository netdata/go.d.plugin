package httpcheck

import (
	"io/ioutil"
	"net"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/l2isbad/go.d.plugin/internal/modules"
	"github.com/l2isbad/go.d.plugin/internal/pkg/helpers/web"
	"github.com/l2isbad/go.d.plugin/internal/pkg/utils"
	"io"
)

const (
	timeout = iota
	failed
	unknown
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

type HttpCheck struct {
	modules.Charts
	modules.Logger

	StatusAccepted []int  `yaml:"status_accepted"`
	ResponseMatch  string `yaml:"response_match"`
	web.Request    `yaml:",inline"`
	web.Client     `yaml:",inline"`

	match    *regexp.Regexp
	statuses map[int]bool
	client   *http.Client
	request  *http.Request

	data data
}

func (hc *HttpCheck) Check() bool {
	// Set Timeout
	if hc.Timeout.Duration == 0 {
		hc.Timeout.Duration = time.Second
	}
	hc.Debugf("Using timeout: %s", hc.Timeout.Duration)

	// Get Request and Client
	req, err := web.CreateRequest(&hc.Request)

	if err != nil {
		hc.Error(err)
		return false
	}

	hc.request = req
	hc.client = web.CreateHttpClient(&hc.Client)

	// Get Response Match Regex
	re, err := regexp.Compile(hc.ResponseMatch)

	if err != nil {
		hc.Error(err)
		return false
	}

	hc.match = re

	// Get Response Statuses
	for _, s := range hc.StatusAccepted {
		hc.statuses[s] = true
	}

	if len(hc.statuses) == 0 {
		hc.statuses[200] = true
	}

	// Get Charts
	c := charts.Copy()
	if len(hc.ResponseMatch) == 0 {
		c.DeleteChartByID("response_check_content")
	}
	hc.AddMany(c)

	return true
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
		return utils.StrToMap(&hc.data)
	}

	defer func() {
		io.Copy(ioutil.Discard, resp.Body)
		resp.Body.Close()
	}()

	hc.data.Success = 1
	// TODO error check ?
	bodyBytes, _ := ioutil.ReadAll(resp.Body)
	hc.data.ResponseLength = len(bodyBytes)

	if !hc.statuses[resp.StatusCode] {
		hc.data.BadStatus = 1
	}

	if hc.match != nil && !hc.match.Match(bodyBytes) {
		hc.data.BadContent = 1
	}

	return utils.StrToMap(&hc.data)
}

func (hc *HttpCheck) doRequest() (*http.Response, error) {
	t := time.Now()
	r, err := hc.client.Do(hc.request)
	hc.data.ResponseTime = int(time.Since(t))
	return r, err
}

func errCheck(err error) int {
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
	modules.SetDefault().SetUpdateEvery(5)

	f := func() modules.Module {
		return &HttpCheck{
			statuses: make(map[int]bool),
			data:     data{},
		}
	}
	modules.Add(f)
}
