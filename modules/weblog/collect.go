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
		}

		n++
		w.collectLogLine()
	}
}

func (w *WebLog) collectLogLine() {
	// TODO: chart filtered?
	if w.line.hasReqURL() && !w.filter.MatchString(w.line.reqURL) {
		w.mx.ReqFiltered.Inc()
		return
	}

	w.mx.Requests.Inc()
	w.collectVhost()
	w.collectPort()
	w.collectReqScheme()
	w.collectReqClient()
	w.collectReqMethod()
	w.collectReqURL()
	w.collectReqProto()
	w.collectRespStatus()
	w.collectReqSize()
	w.collectRespSize()
	w.collectRespTime()
	w.collectUpstreamRespTime()
	w.collectCustom()
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

	c, _ := w.mx.ReqVhost.GetP(w.line.vhost)
	c.Inc()
}

func (w *WebLog) collectPort() {
	if !w.line.hasPort() {
		return
	}
	w.col.port = true

	c, _ := w.mx.ReqPort.GetP(w.line.port)
	c.Inc()
}

func (w *WebLog) collectReqClient() {
	if !w.line.hasReqClient() {
		return
	}
	w.col.client = true

	// TODO: not always IP address
	if strings.ContainsRune(w.line.reqClient, ':') {
		w.mx.ReqIpv6.Inc()
		w.mx.UniqueIPv6.Insert(w.line.reqClient)
		return
	}

	w.mx.ReqIpv4.Inc()
	w.mx.UniqueIPv4.Insert(w.line.reqClient)
}

func (w *WebLog) collectReqScheme() {
	if !w.line.hasReqScheme() {
		return
	}
	w.col.scheme = true

	if w.line.reqScheme == "https" {
		w.mx.ReqHTTPSScheme.Inc()
	} else {
		w.mx.ReqHTTPScheme.Inc()
	}
}

func (w *WebLog) collectReqMethod() {
	if !w.line.hasReqMethod() {
		return
	}
	w.col.method = true

	c, _ := w.mx.ReqMethod.GetP(w.line.reqMethod)
	c.Inc()
}

func (w *WebLog) collectReqURL() {
	if !w.line.hasReqURL() || len(w.urlCats) == 0 {
		return
	}
	w.col.uri = true

	for _, cat := range w.urlCats {
		if !cat.MatchString(w.line.reqURL) {
			continue
		}

		c, _ := w.mx.ReqURI.GetP(cat.name)
		c.Inc()

		w.collectStatsPerURL(cat.name)
		return
	}
}

func (w *WebLog) collectReqProto() {
	if !w.line.hasReqProto() {
		return
	}
	w.col.version = true

	c, _ := w.mx.ReqVersion.GetP(w.line.reqProto)
	c.Inc()
}

func (w *WebLog) collectRespStatus() {
	if !w.line.hasRespStatus() {
		return
	}
	w.col.status = true
	status := w.line.respStatus

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

	w.mx.BytesSent.Add(float64(w.line.reqSize))
}

func (w *WebLog) collectRespSize() {
	if !w.line.hasRespSize() {
		return
	}
	w.col.respSize = true

	w.mx.BytesReceived.Add(float64(w.line.respSize))
}

func (w *WebLog) collectRespTime() {
	if !w.line.hasRespTime() {
		return
	}
	w.col.respTime = true

	w.mx.RespTime.Observe(w.line.respTime)
	if w.mx.RespTimeHist == nil {
		return
	}
	w.mx.RespTimeHist.Observe(w.line.respTime)
}

func (w *WebLog) collectUpstreamRespTime() {
	if !w.line.hasUpstreamRespTime() {
		return
	}
	w.col.upRespTime = true

	w.mx.RespTimeUpstream.Observe(w.line.upsRespTime)
	if w.mx.RespTimeUpstreamHist == nil {
		return
	}
	w.mx.RespTimeUpstreamHist.Observe(w.line.upsRespTime)
}

func (w *WebLog) collectCustom() {
	if !w.line.hasCustom() || len(w.userCats) == 0 {
		return
	}
	w.col.custom = true

	for _, cat := range w.userCats {
		if !cat.MatchString(w.line.custom) {
			continue
		}
		c, _ := w.mx.ReqCustom.GetP(cat.name)
		c.Inc()
		return
	}
}

func (w *WebLog) collectStatsPerURL(uriCat string) {
	v, ok := w.mx.CategorizedStats[uriCat]
	if !ok {
		return
	}

	if w.line.hasRespStatus() {
		status := strconv.Itoa(w.line.respStatus)
		c, _ := v.RespCode.GetP(status)
		c.Inc()
	}

	if w.line.hasReqSize() {
		v.BytesSent.Add(float64(w.line.reqSize))
	}

	if w.line.hasRespSize() {
		v.BytesReceived.Add(float64(w.line.respSize))
	}

	if w.line.hasRespTime() {
		v.RespTime.Observe(w.line.respTime)
	}
}
