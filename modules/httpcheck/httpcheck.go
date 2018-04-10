package httpcheck

import (
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/l2isbad/go.d.plugin/charts/raw"
	"github.com/l2isbad/go.d.plugin/modules"
	"github.com/l2isbad/go.d.plugin/shared/funcs"
	"github.com/l2isbad/go.d.plugin/shared/url_helper"
)

type (
	Charts      = raw.Charts
	Order       = raw.Order
	Definitions = raw.Definitions
	Chart       = raw.Chart
	Options     = raw.Options
	Dimensions  = raw.Dimensions
	Dimension   = raw.Dimension
	Variables   = raw.Variables
	Variable    = raw.Variable
)

var uCharts = Charts{
	Order: Order{
		"response_time", "response_length", "response_status", "response_check_status", "response_check_content"},
	Definitions: Definitions{
		Chart{
			ID:      "response_time",
			Options: Options{"HTTP Response Time", "ms", "response", "httpcheck.response_time"},
			Dimensions: Dimensions{
				Dimension{"response_time", "time", "", 1, 1e6},
			},
		},
		Chart{
			ID:      "response_length",
			Options: Options{"HTTP Response Body Length", "characters", "response", "httpcheck.response_length"},
			Dimensions: Dimensions{
				Dimension{"response_length", "length"},
			},
		},
		Chart{
			ID:      "response_status",
			Options: Options{"HTTP Response Status", "boolean", "status", "httpcheck.status"},
			Dimensions: Dimensions{
				Dimension{"success"},
				Dimension{"failed"},
				Dimension{"timeout"},
			},
		},
		Chart{
			ID:      "response_check_status",
			Options: Options{"HTTP Response Check Status", "boolean", "status", "httpcheck.check_status"},
			Dimensions: Dimensions{
				Dimension{"bad_status", "bad status"},
			},
		},
		Chart{
			ID:      "response_check_content",
			Options: Options{"HTTP Response Check Content", "boolean", "status", "httpcheck.check_content"},
			Dimensions: Dimensions{
				Dimension{"bad_content", "bad content"},
			},
		},
	},
}

type data struct {
	Success        int `stm:"success"`
	Failed         int `stm:"failed"`
	Timeout        int `stm:"timeout"`
	BadContent     int `stm:"bad_content"`
	BadStatus      int `stm:"bad_status"`
	ResponseTime   int `stm:"response_time"`
	ResponseLength int `stm:"response_length"`
}

type HttpCheck struct {
	modules.Charts
	modules.Logger
	modules.BaseConfHook

	StatusAccepted []int  `toml:"status_accepted"`
	ResponseMatch  string `toml:"response_match"`
	url_helper.Request
	url_helper.Client

	responseMatch  *regexp.Regexp
	statusAccepted map[int]bool
	client         *http.Client
	request        *http.Request

	data data
}

func (h *HttpCheck) Check() bool {
	rawCharts := uCharts.Copy()

	if len(h.ResponseMatch) == 0 {
		rawCharts.DeleteChartByID("response_check_content")
	} else {
		if re, err := regexp.Compile(h.ResponseMatch); err != nil {
			h.Errorf("regex compile failed: %s", err)
			return false
		} else {
			h.responseMatch = re
		}
	}

	if len(h.StatusAccepted) != 0 {
		delete(h.statusAccepted, 200)
		for _, s := range h.StatusAccepted {
			h.statusAccepted[s] = true
		}
	}

	if h.Timeout.Duration == 0 {
		h.Timeout.Duration = time.Duration(h.UpdateEvery()) * time.Second
		h.Warningf("timeout not specified. Setting to %s", h.Timeout.Duration)
	}

	req, err := url_helper.CreateRequest(&h.Request)
	if err != nil {
		h.Error(err)
		return false
	}
	h.request = req
	h.client = url_helper.CreateHttpClient(&h.Client)

	h.AddMany(rawCharts)
	return true
}

func (h *HttpCheck) GetData() *map[string]int64 {
	h.data = data{}

	start := time.Now()
	resp, err := h.client.Do(h.request)
	h.data.ResponseTime = int(time.Since(start))

	if err != nil {
		h.Debug(err)
		v, ok := err.(net.Error)
		switch {
		case ok && v.Timeout():
			h.data.Timeout = 1
		case ok && strings.Contains(v.Error(), "connection refused"):
			h.data.Failed = 1
		default:
			h.Error(err)
			return nil
		}
		return funcs.ToMap(&h.data)
	}

	defer func() {
		io.Copy(ioutil.Discard, resp.Body)
		resp.Body.Close()
	}()

	h.data.Success = 1
	// TODO error check ?
	bodyBytes, _ := ioutil.ReadAll(resp.Body)
	h.data.ResponseLength = len(bodyBytes)

	if !h.statusAccepted[resp.StatusCode] {
		h.data.BadStatus = 1
	}
	if h.responseMatch != nil && !h.responseMatch.Match(bodyBytes) {
		h.data.BadContent = 1
	}

	return funcs.ToMap(&h.data)
}

func init() {
	modules.SetDefault(modules.UpdateEvery).Set(5)

	f := func() modules.Module {
		return &HttpCheck{
			statusAccepted: map[int]bool{200: true}}
	}
	modules.Add(f)
}
