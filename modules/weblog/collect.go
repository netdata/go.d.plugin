package weblog

import (
	"io"
	"runtime"
	"strconv"
	"strings"

	"github.com/netdata/go.d.plugin/pkg/logs/parse"
	"github.com/netdata/go.d.plugin/pkg/stm"
)

func (w WebLog) logPanicStackIfAny() {
	err := recover()
	if err == nil {
		return
	}
	w.Errorf("[ERROR] %s\n", err)
	for depth := 0; ; depth++ {
		_, file, line, ok := runtime.Caller(depth)
		if !ok {
			break
		}
		w.Errorf("======> %d: %v:%d", depth, file, line)
	}
	panic(err)
}

func (w *WebLog) collect() (mx map[string]int64, err error) {
	defer w.logPanicStackIfAny()
	w.metrics.Reset()
	var n int

	n, err = w.collectLogLines()

	if n > 0 || n == 0 && err == nil {
		mx = stm.ToMap(w.metrics)
	}
	if n > 0 {
		w.updateCharts()
	}
	return mx, err
}

func (w *WebLog) collectLogLines() (int, error) {
	var n int
	for {
		w.line.reset()
		err := w.parser.ReadLine(w.line)
		if err != nil {
			if err == io.EOF {
				return n, nil
			}
			if !parse.IsParseError(err) {
				return n, err
			}
			w.collectUnmatched()
		}
		n++
		w.collectLogLine()
	}
}

func (w *WebLog) collectLogLine() {
	if w.line.hasReqURI() && !w.filter.MatchString(w.line.ReqURI) {
		return
	}

	w.metrics.Requests.Inc()
	w.collectVhost()
	w.collectClientAddr()
	w.collectReqHTTPMethod()
	w.collectReqURI()
	w.collectReqHTTPVersion()
	w.collectRespStatusCode()
	w.collectReqSize()
	w.collectRespSize()
	w.collectRespTime()
	w.collectUpstreamRespTime()
	w.collectCustom()
}

func (w *WebLog) collectUnmatched() {
	w.metrics.Requests.Inc()
	w.metrics.ReqUnmatched.Inc()
}

func (w *WebLog) collectVhost() {
	if !w.line.hasVhost() {
		return
	}
	w.collected.vhost = true

	c, _ := w.metrics.ReqVhost.GetP(w.line.Vhost)
	c.Inc()
}

func (w *WebLog) collectClientAddr() {
	if !w.line.hasClientAddr() {
		return
	}

	w.collected.client = true

	// TODO: not always IP address
	if strings.ContainsRune(w.line.ClientAddr, ':') {
		w.metrics.ReqIpv6.Inc()
		w.metrics.UniqueIPv6.Insert(w.line.ClientAddr)
		return
	}

	w.metrics.ReqIpv4.Inc()
	w.metrics.UniqueIPv4.Insert(w.line.ClientAddr)
}

func (w *WebLog) collectReqHTTPMethod() {
	if !w.line.hasReqHTTPMethod() {
		return
	}
	w.collected.method = true

	c, _ := w.metrics.ReqMethod.GetP(w.line.ReqHTTPMethod)
	c.Inc()
}

func (w *WebLog) collectReqURI() {
	if !w.line.hasReqURI() || len(w.urlCategories) == 0 {
		return
	}
	w.collected.uri = true

	for _, cat := range w.urlCategories {
		if !cat.Matcher.MatchString(w.line.ReqURI) {
			continue
		}
		c, _ := w.metrics.ReqURI.GetP(cat.name)
		c.Inc()
		return
	}
}

func (w *WebLog) collectReqHTTPVersion() {
	if !w.line.hasReqHTTPVersion() {
		return
	}
	w.collected.version = true

	c, _ := w.metrics.ReqVersion.GetP(w.line.ReqHTTPVersion)
	c.Inc()
}

func (w *WebLog) collectRespStatusCode() {
	if !w.line.hasRespCodeStatus() {
		return
	}
	w.collected.status = true
	status := w.line.RespCodeStatus

	switch {
	case status >= 100 && status < 300, status == 304:
		w.metrics.RespSuccessful.Inc()
	case status >= 300 && status < 400:
		w.metrics.RespRedirect.Inc()
	case status >= 400 && status < 500:
		w.metrics.RespClientError.Inc()
	case status >= 500 && status < 600:
		w.metrics.RespServerError.Inc()
	}

	switch status / 100 {
	case 1:
		w.metrics.Resp1xx.Inc()
	case 2:
		w.metrics.Resp2xx.Inc()
	case 3:
		w.metrics.Resp3xx.Inc()
	case 4:
		w.metrics.Resp4xx.Inc()
	case 5:
		w.metrics.Resp5xx.Inc()
	}

	statusStr := strconv.Itoa(status)
	c, _ := w.metrics.RespCode.GetP(statusStr)
	c.Inc()
}

func (w *WebLog) collectReqSize() {
	if !w.line.hasReqSize() {
		return
	}
	w.collected.reqSize = true

	w.metrics.BytesSent.Add(float64(w.line.ReqSize))
}

func (w *WebLog) collectRespSize() {
	if !w.line.hasRespSize() {
		return
	}
	w.collected.respSize = true

	w.metrics.BytesReceived.Add(float64(w.line.RespSize))
}

func (w *WebLog) collectRespTime() {
	if !w.line.hasRespTime() {
		return
	}
	w.collected.respTime = true

	w.metrics.RespTime.Observe(w.line.RespTime)
	if w.metrics.RespTimeHist == nil {
		return
	}
	w.metrics.RespTimeHist.Observe(w.line.RespTime)
}

func (w *WebLog) collectUpstreamRespTime() {
	if !w.line.hasUpstreamRespTime() {
		return
	}
	w.collected.upRespTime = true

	w.metrics.RespTimeUpstream.Observe(w.line.UpstreamRespTime)
	if w.metrics.RespTimeUpstreamHist == nil {
		return
	}
	w.metrics.RespTimeUpstreamHist.Observe(w.line.UpstreamRespTime)
}

func (w *WebLog) collectCustom() {
	if !w.line.hasCustom() || len(w.userCategories) == 0 {
		return
	}
	w.collected.custom = true

	for _, cat := range w.userCategories {
		if !cat.Matcher.MatchString(w.line.Custom) {
			continue
		}
		c, _ := w.metrics.ReqCustom.GetP(cat.name)
		c.Inc()
		return
	}
}
