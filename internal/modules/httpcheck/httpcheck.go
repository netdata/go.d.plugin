package httpcheck

import (
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/l2isbad/go.d.plugin/internal/modules"
	"github.com/l2isbad/go.d.plugin/internal/pkg/helpers/web"
	"github.com/l2isbad/go.d.plugin/internal/pkg/utils"
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

type HttpCheck struct {
	modules.Charts
	modules.Logger
	modules.BaseConfHook

	StatusAccepted []int  `yaml:"status_accepted"`
	ResponseMatch  string `yaml:"response_match"`
	web.Request
	web.Client

	responseMatch  *regexp.Regexp
	statusAccepted map[int]bool
	client         *http.Client
	request        *http.Request

	data data
}

func (h *HttpCheck) Check() bool {
	rawCharts := charts.Copy()

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
		h.Timeout.Duration = time.Duration(h.GetUpdateEvery()) * time.Second
		h.Warningf("timeout not specified. Setting to %s", h.Timeout.Duration)
	}

	req, err := web.CreateRequest(&h.Request)
	if err != nil {
		h.Error(err)
		return false
	}
	h.request = req
	h.client = web.CreateHttpClient(&h.Client)

	h.AddMany(rawCharts)
	return true
}

func (h *HttpCheck) GetData() map[string]int64 {
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
		return utils.StrToMap(&h.data)
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

	return utils.StrToMap(&h.data)
}

func init() {
	modules.SetDefault().SetUpdateEvery(5)

	f := func() modules.Module {
		return &HttpCheck{
			statusAccepted: map[int]bool{200: true}}
	}
	modules.Add(f)
}
