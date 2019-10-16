package weblog

import (
	"io"
	"runtime"
	"strconv"
	"strings"

	"github.com/netdata/go.d.plugin/modules/weblog/parser"
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

func (w *WebLog) collect() (map[string]int64, error) {
	defer w.logPanicStackIfAny()
	w.metrics.Reset()

	n, err := w.collectLogLines()
	if err != nil {
		return nil, err
	}

	if n > 0 {
		w.updateCharts()
	}

	result := stm.ToMap(w.metrics)
	return result, nil
}

func (w *WebLog) collectLogLines() (int, error) {
	var n int
	for {
		line, err := w.parser.ReadLine()
		if err != nil {
			if err == io.EOF {
				break
			}
			// TODO: collect unmatched
			// w.collectUnmatched()
			return n, err
		}
		n++
		w.collectLogLine(line)
	}
	return n, nil
}

func (w *WebLog) collectLogLine(line parser.LogLine) {
	if !isEmptyString(line.ReqURI) && !w.filter.MatchString(line.ReqURI) {
		return
	}

	w.metrics.Requests.Inc()
	w.collectVhost(line.Vhost)
	w.collectClientAddr(line.ClientAddr)
	w.collectReqHTTPMethod(line.ReqHTTPMethod)
	w.collectReqURI(line.ReqURI)
	w.collectReqHTTPVersion(line.ReqHTTPVersion)
	w.collectRespStatusCode(line.RespCodeStatus)
	w.collectReqSize(line.ReqSize)
	w.collectRespSize(line.RespSize)
	w.collectRespTime(line.RespTime)
	w.collectUpstreamRespTime(line.UpstreamRespTime)
	w.collectCustom(line.Custom)
}

func (w *WebLog) collectUnmatched() {
	w.metrics.Requests.Inc()
	w.metrics.ReqUnmatched.Inc()
}

func (w *WebLog) collectVhost(vhost string) {
	if isEmptyString(vhost) {
		return
	}
	w.collected.vhost = true

	c, _ := w.metrics.ReqVhost.GetP(vhost)
	c.Inc()
}

func (w *WebLog) collectClientAddr(client string) {
	if isEmptyString(client) {
		return
	}
	w.collected.client = true

	// TODO: not always IP address
	if strings.ContainsRune(client, ':') {
		w.metrics.ReqIpv6.Inc()
		w.metrics.UniqueIPv6.Insert(client)
		return
	}

	w.metrics.ReqIpv4.Inc()
	w.metrics.UniqueIPv4.Insert(client)
}

func (w *WebLog) collectReqHTTPMethod(method string) {
	if isEmptyString(method) {
		return
	}
	w.collected.method = true

	c, _ := w.metrics.ReqMethod.GetP(method)
	c.Inc()
}

func (w *WebLog) collectReqURI(uri string) {
	if isEmptyString(uri) || len(w.urlCategories) == 0 {
		return
	}
	w.collected.uri = true

	for _, cat := range w.urlCategories {
		if !cat.Matcher.MatchString(uri) {
			continue
		}
		c, _ := w.metrics.ReqURI.GetP(cat.name)
		c.Inc()
		return
	}
}

func (w *WebLog) collectReqHTTPVersion(version string) {
	if isEmptyString(version) {
		return
	}
	w.collected.version = true

	c, _ := w.metrics.ReqVersion.GetP(version)
	c.Inc()
}

func (w *WebLog) collectRespStatusCode(status int) {
	if isEmptyNumber(status) {
		return
	}
	w.collected.status = true

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

func (w *WebLog) collectReqSize(size int) {
	if isEmptyNumber(size) {
		return
	}
	w.collected.reqSize = true

	w.metrics.BytesSent.Add(float64(size))
}

func (w *WebLog) collectRespSize(size int) {
	if isEmptyNumber(size) {
		return
	}
	w.collected.respSize = true

	w.metrics.BytesReceived.Add(float64(size))
}

func (w *WebLog) collectRespTime(respTime float64) {
	if isEmptyNumber(int(respTime)) {
		return
	}
	w.collected.respTime = true

	w.metrics.RespTime.Observe(respTime)
	if w.metrics.RespTimeHist == nil {
		return
	}
	w.metrics.RespTimeHist.Observe(respTime)
}

func (w *WebLog) collectUpstreamRespTime(respTime float64) {
	if isEmptyNumber(int(respTime)) {
		return
	}
	w.collected.upRespTime = true

	w.metrics.RespTimeUpstream.Observe(respTime)
	if w.metrics.RespTimeUpstreamHist == nil {
		return
	}
	w.metrics.RespTimeUpstreamHist.Observe(respTime)
}

func (w *WebLog) collectCustom(custom string) {
	if isEmptyString(custom) || len(w.userCategories) == 0 {
		return
	}
	w.collected.custom = true

	for _, cat := range w.userCategories {
		if !cat.Matcher.MatchString(custom) {
			continue
		}
		c, _ := w.metrics.ReqCustom.GetP(cat.name)
		c.Inc()
		return
	}
}

func isEmptyString(s string) bool { return s != parser.EmptyString }

func isEmptyNumber(n int) bool { return n != parser.EmptyNumber }
