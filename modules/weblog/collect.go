package weblog

import (
	"io"
	"runtime"
	"strconv"
	"strings"

	"github.com/netdata/go.d.plugin/pkg/logs"
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
	w.mx.Reset()
	var n int

	n, err = w.collectLogLines()

	if n > 0 || err == nil {
		mx = stm.ToMap(w.mx)
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
			if !logs.IsParseError(err) {
				return n, err
			}
			w.collectUnmatched()
			continue
		}

		n++
		w.collectLogLine()
	}
}

func (w *WebLog) collectLogLine() {
	if w.line.hasReqURI() && !w.filter.MatchString(w.line.ReqURI) {
		return
	}

	w.mx.Requests.Inc()
	//w.collectVhost()
	//w.collectPort()
	//w.collectScheme()
	//w.collectClientAddr()
	//w.collectReqHTTPMethod()
	w.collectReqURI()
	//w.collectReqHTTPVersion()
	//w.collectRespStatusCode()
	//w.collectReqSize()
	//w.collectRespSize()
	//w.collectRespTime()
	//w.collectUpstreamRespTime()
	//w.collectCustom()
}

func (w *WebLog) collectUnmatched() {
	w.mx.Requests.Inc()
	w.mx.ReqUnmatched.Inc()
}

func (w *WebLog) collectVhost() {
	if !w.line.hasVhost() {
		return
	}
	w.col.vhost = true

	c, _ := w.mx.ReqVhost.GetP(w.line.Vhost)
	c.Inc()
}

func (w *WebLog) collectPort() {
	if !w.line.hasPort() {
		return
	}
	w.col.port = true

	c, _ := w.mx.ReqPort.GetP(w.line.Port)
	c.Inc()
}

func (w *WebLog) collectClientAddr() {
	if !w.line.hasClientAddr() {
		return
	}
	w.col.client = true

	// TODO: not always IP address
	if strings.ContainsRune(w.line.ClientAddr, ':') {
		w.mx.ReqIpv6.Inc()
		w.mx.UniqueIPv6.Insert(w.line.ClientAddr)
		return
	}

	w.mx.ReqIpv4.Inc()
	w.mx.UniqueIPv4.Insert(w.line.ClientAddr)
}

func (w *WebLog) collectScheme() {
	if !w.line.hasScheme() {
		return
	}
	w.col.scheme = true

	if w.line.Scheme == "https" {
		w.mx.ReqHTTPSScheme.Inc()
	} else {
		w.mx.ReqHTTPScheme.Inc()
	}
}

func (w *WebLog) collectReqHTTPMethod() {
	if !w.line.hasReqHTTPMethod() {
		return
	}
	w.col.method = true

	c, _ := w.mx.ReqMethod.GetP(w.line.ReqHTTPMethod)
	c.Inc()
}

func (w *WebLog) collectReqURI() {
	if !w.line.hasReqURI() || len(w.urlCats) == 0 {
		return
	}
	w.col.uri = true

	for _, cat := range w.urlCats {
		if !cat.MatchString(w.line.ReqURI) {
			continue
		}

		c, _ := w.mx.ReqURI.GetP(cat.name)
		c.Inc()

		w.collectStatsPerURI(cat.name)
	}
}

func (w *WebLog) collectReqHTTPVersion() {
	if !w.line.hasReqHTTPVersion() {
		return
	}
	w.col.version = true

	c, _ := w.mx.ReqVersion.GetP(w.line.ReqHTTPVersion)
	c.Inc()
}

func (w *WebLog) collectRespStatusCode() {
	if !w.line.hasRespCode() {
		return
	}
	w.col.status = true
	status := w.line.RespCode

	switch {
	case status >= 100 && status < 300, status == 304:
		w.mx.RespSuccessful.Inc()
	case status >= 300 && status < 400:
		w.mx.RespRedirect.Inc()
	case status >= 400 && status < 500:
		w.mx.RespClientError.Inc()
	case status >= 500 && status < 600:
		w.mx.RespServerError.Inc()
	}

	switch status / 100 {
	case 1:
		w.mx.Resp1xx.Inc()
	case 2:
		w.mx.Resp2xx.Inc()
	case 3:
		w.mx.Resp3xx.Inc()
	case 4:
		w.mx.Resp4xx.Inc()
	case 5:
		w.mx.Resp5xx.Inc()
	}

	statusStr := strconv.Itoa(status)
	c, _ := w.mx.RespCode.GetP(statusStr)
	c.Inc()
}

func (w *WebLog) collectReqSize() {
	if !w.line.hasReqSize() {
		return
	}
	w.col.reqSize = true

	w.mx.BytesSent.Add(float64(w.line.ReqSize))
}

func (w *WebLog) collectRespSize() {
	if !w.line.hasRespSize() {
		return
	}
	w.col.respSize = true

	w.mx.BytesReceived.Add(float64(w.line.RespSize))
}

func (w *WebLog) collectRespTime() {
	if !w.line.hasRespTime() {
		return
	}
	w.col.respTime = true

	w.mx.RespTime.Observe(w.line.RespTime)
	if w.mx.RespTimeHist == nil {
		return
	}
	w.mx.RespTimeHist.Observe(w.line.RespTime)
}

func (w *WebLog) collectUpstreamRespTime() {
	if !w.line.hasUpstreamRespTime() {
		return
	}
	w.col.upRespTime = true

	w.mx.RespTimeUpstream.Observe(w.line.UpstreamRespTime)
	if w.mx.RespTimeUpstreamHist == nil {
		return
	}
	w.mx.RespTimeUpstreamHist.Observe(w.line.UpstreamRespTime)
}

func (w *WebLog) collectCustom() {
	if !w.line.hasCustom() || len(w.userCats) == 0 {
		return
	}
	w.col.custom = true

	for _, cat := range w.userCats {
		if !cat.MatchString(w.line.Custom) {
			continue
		}
		c, _ := w.mx.ReqCustom.GetP(cat.name)
		c.Inc()
		return
	}
}

func (w *WebLog) collectStatsPerURI(uriCat string) {
	v, ok := w.mx.CategorizedStats[uriCat]
	if !ok {
		return
	}

	if w.line.hasRespCode() {
		status := strconv.Itoa(w.line.RespCode)
		c, _ := v.RespCode.GetP(status)
		c.Inc()
	}

	if w.line.hasReqSize() {
		v.BytesSent.Add(float64(w.line.ReqSize))
	}

	if w.line.hasRespSize() {
		v.BytesReceived.Add(float64(w.line.RespSize))
	}

	if w.line.hasRespTime() {
		v.RespTime.Observe(w.line.RespTime)
	}
}
