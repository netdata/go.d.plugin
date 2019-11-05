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

func (w *WebLog) collect() (map[string]int64, error) {
	defer w.logPanicStackIfAny()
	w.mx.reset()

	var mx map[string]int64

	n, err := w.collectLogLines()

	if n > 0 || err == nil {
		mx = stm.ToMap(w.mx)
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
			n++
			w.collectUnmatched()
			continue
		}
		n++
		if !w.line.hasRespStatusCode() {
			w.collectUnmatched()
			continue
		}
		w.collectLogLine()
	}
}

func (w *WebLog) collectLogLine() {
	w.mx.Requests.Inc()
	if w.line.hasReqURL() && !w.filter.MatchString(w.line.reqURL) {
		w.mx.ReqFiltered.Inc()
		return
	}
	w.collectVhost()
	w.collectPort()
	w.collectReqScheme()
	w.collectReqClient()
	w.collectReqMethod()
	w.collectReqURL()
	w.collectReqProto()
	w.collectRespStatusCode()
	w.collectReqSize()
	w.collectRespSize()
	w.collectReqProcTime()
	w.collectUpstreamRespTime()
	w.collectSSLProto()
	w.collectSSLCipherSuite()
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
	c, ok := w.mx.ReqVhost.GetP(w.line.vhost)
	if !ok {
		w.addDimToVhostChart(w.line.vhost)
	}
	c.Inc()
}

func (w *WebLog) collectPort() {
	if !w.line.hasPort() {
		return
	}
	c, ok := w.mx.ReqPort.GetP(w.line.port)
	if !ok {
		w.addDimToPortChart(w.line.port)
	}
	c.Inc()
}

func (w *WebLog) collectReqClient() {
	if !w.line.hasReqClient() {
		return
	}
	if strings.ContainsRune(w.line.reqClient, ':') {
		w.mx.ReqIPv6.Inc()
		w.mx.UniqueIPv6.Insert(w.line.reqClient)
		return
	}
	// NOTE: count hostname as IPv4 address
	w.mx.ReqIPv4.Inc()
	w.mx.UniqueIPv4.Insert(w.line.reqClient)
}

func (w *WebLog) collectReqScheme() {
	if !w.line.hasReqScheme() {
		return
	}
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
	c, ok := w.mx.ReqMethod.GetP(w.line.reqMethod)
	if !ok {
		w.addDimToReqMethodChart(w.line.reqMethod)
	}
	c.Inc()
}

func (w *WebLog) collectReqURL() {
	if !w.line.hasReqURL() || len(w.patURL) == 0 {
		return
	}
	for _, p := range w.patURL {
		if !p.MatchString(w.line.reqURL) {
			continue
		}
		c, _ := w.mx.ReqURLPattern.GetP(p.name)
		c.Inc()

		w.collectURLPatternStats(p.name)
		return
	}
}

func (w *WebLog) collectReqProto() {
	if !w.line.hasReqProto() {
		return
	}
	c, ok := w.mx.ReqVersion.GetP(w.line.reqProto)
	if !ok {
		w.addDimToReqVersionChart(w.line.reqProto)
	}
	c.Inc()
}

func (w *WebLog) collectRespStatusCode() {
	if !w.line.hasRespStatusCode() {
		return
	}
	code := w.line.respStatusCode
	//  1xx (Informational): The request was received, continuing process.
	//  2xx (Successful): The request was successfully received, understood, and accepted.
	//  3xx (Redirection): Further action needs to be taken in order to complete the request.
	//  4xx (Client Error): The request contains bad syntax or cannot be fulfilled.
	//  5xx (Server Error): The server failed to fulfill an apparently valid request.
	// TODO: this grouping is confusing since it uses terms from rfc7231
	switch {
	case code >= 100 && code < 300, code == 304:
		w.mx.RespSuccessful.Inc()
	case code >= 300 && code < 400:
		w.mx.RespRedirect.Inc()
	case code >= 400 && code < 500:
		w.mx.RespClientError.Inc()
	case code >= 500 && code < 600:
		w.mx.RespServerError.Inc()
	}

	switch code / 100 {
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

	codeStr := strconv.Itoa(code)
	c, ok := w.mx.RespStatusCode.GetP(codeStr)
	if !ok {
		w.addDimToRespStatusCodeChart(codeStr)
	}
	c.Inc()
}

func (w *WebLog) collectReqSize() {
	if !w.line.hasReqSize() {
		return
	}
	w.mx.BytesSent.Add(float64(w.line.reqSize))
}

func (w *WebLog) collectRespSize() {
	if !w.line.hasRespSize() {
		return
	}
	w.mx.BytesReceived.Add(float64(w.line.respSize))
}

func (w *WebLog) collectReqProcTime() {
	if !w.line.hasReqProcTime() {
		return
	}
	w.mx.ReqProcTime.Observe(w.line.reqProcTime)
	if w.mx.ReqProcTimeHist == nil {
		return
	}
	w.mx.ReqProcTimeHist.Observe(w.line.reqProcTime)
}

func (w *WebLog) collectUpstreamRespTime() {
	if !w.line.hasUpstreamRespTime() {
		return
	}
	w.mx.UpsRespTime.Observe(w.line.upsRespTime)
	if w.mx.UpsRespTimeHist == nil {
		return
	}
	w.mx.UpsRespTimeHist.Observe(w.line.upsRespTime)
}

func (w *WebLog) collectSSLProto() {
	if !w.line.hasSSLProto() {
		return
	}
	c, ok := w.mx.ReqSSLProto.GetP(w.line.sslProto)
	if !ok {
		w.addDimToSSLProtoChart(w.line.sslProto)
	}
	c.Inc()
}

func (w *WebLog) collectSSLCipherSuite() {
	if !w.line.hasSSLCipherSuite() {
		return
	}
	c, ok := w.mx.ReqSSLCipherSuite.GetP(w.line.sslCipherSuite)
	if !ok {
		w.addDimToSSLCipherSuiteChart(w.line.sslCipherSuite)
	}
	c.Inc()
}

func (w *WebLog) collectCustom() {
	if !w.line.hasCustom() || len(w.patCustom) == 0 {
		return
	}
	for _, p := range w.patCustom {
		if !p.MatchString(w.line.custom) {
			continue
		}
		c, _ := w.mx.ReqCustomPattern.GetP(p.name)
		c.Inc()
		return
	}
}

func (w *WebLog) collectURLPatternStats(name string) {
	v, ok := w.mx.URLPatternStats[name]
	if !ok {
		return
	}
	if w.line.hasRespStatusCode() {
		status := strconv.Itoa(w.line.respStatusCode)
		c, ok := v.RespStatusCode.GetP(status)
		if !ok {
			w.addDimToURLPatternRespStatusCodeChart(name, status)
		}
		c.Inc()
	}

	if w.line.hasReqSize() {
		v.BytesSent.Add(float64(w.line.reqSize))
	}

	if w.line.hasRespSize() {
		v.BytesReceived.Add(float64(w.line.respSize))
	}

	if w.line.hasReqProcTime() {
		v.ReqProcTime.Observe(w.line.reqProcTime)
	}
}
