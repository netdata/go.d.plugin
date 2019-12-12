package squidlog

import (
	"io"
	"runtime"
	"strconv"
	"strings"

	"github.com/netdata/go.d.plugin/pkg/logs"
	"github.com/netdata/go.d.plugin/pkg/stm"

	"github.com/netdata/go-orchestrator/module"
)

func (s SquidLog) logPanicStackIfAny() {
	err := recover()
	if err == nil {
		return
	}
	s.Errorf("[ERROR] %s\n", err)
	for depth := 0; ; depth++ {
		_, file, line, ok := runtime.Caller(depth)
		if !ok {
			break
		}
		s.Errorf("======> %d: %v:%d", depth, file, line)
	}
	panic(err)
}

func (s *SquidLog) collect() (map[string]int64, error) {
	defer s.logPanicStackIfAny()
	s.mx.reset()

	var mx map[string]int64

	n, err := s.collectLogLines()

	if n > 0 || err == nil {
		mx = stm.ToMap(s.mx)
	}
	return mx, err
}

func (s *SquidLog) collectLogLines() (int, error) {
	var n int
	for {
		s.line.reset()
		err := s.parser.ReadLine(s.line)
		if err != nil {
			if err == io.EOF {
				return n, nil
			}
			if !logs.IsParseError(err) {
				return n, err
			}
			n++
			s.collectUnmatched()
			continue
		}
		n++
		if s.line.empty() {
			s.collectUnmatched()
		} else {
			s.collectLogLine()
		}
	}
}

func (s *SquidLog) collectLogLine() {
	s.mx.Requests.Inc()
	s.collectRespSize()
	s.collectClientAddress()
	s.collectCacheCode()
	s.collectHTTPCode()
	s.collectRespSize()
	s.collectReqMethod()
	s.collectHierCode()
	s.collectServerAddress()
	s.collectMimeType()
}

func (s *SquidLog) collectUnmatched() {
	s.mx.Requests.Inc()
	s.mx.ReqUnmatched.Inc()
}

func (s *SquidLog) collectRespTime() {
	if !s.line.hasRespTime() {
		return
	}
	s.mx.RespTime.Observe(float64(s.line.respTime))
}

func (s *SquidLog) collectClientAddress() {
	if !s.line.hasClientAddress() {
		return
	}
	s.mx.UniqueClients.Insert(s.line.clientAddr)
}

func (s *SquidLog) collectCacheCode() {
	if !s.line.hasCacheCode() {
		return
	}

	c, ok := s.mx.CacheCode.GetP(s.line.cacheCode)
	if !ok {
		s.addDimToCacheCodeChart(s.line.cacheCode)
	}
	c.Inc()

	parts := strings.Split(s.line.cacheCode, "_")
	for _, part := range parts {
		s.collectCacheCodePart(part)
	}
}

func (s *SquidLog) collectHTTPCode() {
	if !s.line.hasHTTPCode() {
		return
	}

	code := s.line.httpCode
	switch {
	case code >= 100 && code < 300, code == 304, code == 401, code == 0:
		s.mx.ReqSuccess.Inc()
	case code >= 300 && code < 400:
		s.mx.ReqRedirect.Inc()
	case code >= 400 && code < 500:
		s.mx.ReqBad.Inc()
	case code >= 500 && code < 603:
		s.mx.ReqError.Inc()
	}

	switch code / 100 {
	case 0:
		s.mx.HTTP0xx.Inc()
	case 1:
		s.mx.HTTP1xx.Inc()
	case 2:
		s.mx.HTTP2xx.Inc()
	case 3:
		s.mx.HTTP3xx.Inc()
	case 4:
		s.mx.HTTP4xx.Inc()
	case 5:
		s.mx.HTTP5xx.Inc()
	case 6:
		s.mx.HTTP6xx.Inc()
	}

	codeStr := strconv.Itoa(code)
	c, ok := s.mx.HTTPCode.GetP(codeStr)
	if !ok {
		s.addDimToRespCodesChart(codeStr)
	}
	c.Inc()
}

func (s *SquidLog) collectRespSize() {
	if !s.line.hasRespSize() {
		return
	}
	s.mx.BytesSent.Add(float64(s.line.respSize))
}

func (s *SquidLog) collectReqMethod() {
	if !s.line.hasReqMethod() {
		return
	}
	c, ok := s.mx.ReqMethod.GetP(s.line.reqMethod)
	if !ok {
		s.addDimToReqMethodChart(s.line.reqMethod)
	}
	c.Inc()
}

func (s *SquidLog) collectHierCode() {
	if !s.line.hasHierCode() {
		return
	}
	c, ok := s.mx.HierCode.GetP(s.line.hierCode)
	if !ok {
		s.addDimToHierCodeChart(s.line.hierCode)
	}
	c.Inc()
}

func (s *SquidLog) collectServerAddress() {
	if !s.line.hasServerAddress() {
		return
	}
	c, ok := s.mx.Server.GetP(s.line.serverAddr)
	if !ok {
		s.addDimToServerAddressChart(s.line.serverAddr)
	}
	c.Inc()
}

func (s *SquidLog) collectMimeType() {
	if !s.line.hasMimeType() {
		return
	}
	c, ok := s.mx.MimeType.GetP(s.line.mimeType)
	if !ok {
		s.addDimToMimeTypeChart(s.line.mimeType)
	}
	c.Inc()
}

func (s *SquidLog) collectCacheCodePart(codePart string) {
	// https://wiki.squid-cache.org/SquidFaq/SquidLogs#Squid_result_codes
	switch codePart {
	default:
	case "TCP", "UDP", "NONE":
		c, ok := s.mx.CacheCodeTransport.GetP(codePart)
		if !ok {
			s.addDimToCacheCodeTransportChart(codePart)
		}
		c.Inc()
	case "CF", "CLIENT", "IMS", "ASYNC", "SWAPFAIL", "REFRESH", "SHARED", "REPLY":
		c, ok := s.mx.CacheCodeHandling.GetP(codePart)
		if !ok {
			s.addDimToCacheCodeHandlingChart(codePart)
		}
		c.Inc()
	case "NEGATIVE", "STALE", "OFFLINE", "INVALID", "FAIL", "MODIFIED", "UNMODIFIED", "REDIRECT":
		c, ok := s.mx.CacheCodeObject.GetP(codePart)
		if !ok {
			s.addDimToCacheCodeObjectChart(codePart)
		}
		c.Inc()
	case "HIT", "MEM", "MISS", "DENIED", "NOFETCH", "TUNNEL":
		c, ok := s.mx.CacheCodeLoadSource.GetP(codePart)
		if !ok {
			s.addDimToCacheCodeLoadSourceChart(codePart)
		}
		c.Inc()
	case "ABORTED", "TIMEOUT", "IGNORED":
		c, ok := s.mx.CacheCodeError.GetP(codePart)
		if !ok {
			s.addDimToCacheCodeErrorChart(codePart)
		}
		c.Inc()
	}
}

func (s *SquidLog) addDimToCacheCodeTransportChart(codePart string) {
	s.addDimToChart(cacheCodeTransport.ID, "cache_code_transport_"+codePart, codePart)
}

func (s *SquidLog) addDimToCacheCodeHandlingChart(codePart string) {
	s.addDimToChart(cacheCodeHandling.ID, "cache_code_handling_"+codePart, codePart)
}

func (s *SquidLog) addDimToCacheCodeObjectChart(codePart string) {
	s.addDimToChart(cacheCodeObject.ID, "cache_code_object_"+codePart, codePart)
}

func (s *SquidLog) addDimToCacheCodeLoadSourceChart(codePart string) {
	s.addDimToChart(cacheCodeLoadSource.ID, "cache_code_load_source_"+codePart, codePart)
}

func (s *SquidLog) addDimToCacheCodeErrorChart(codePart string) {
	s.addDimToChart(cacheCodeError.ID, "cache_code_error_"+codePart, codePart)
}

func (s *SquidLog) addDimToMimeTypeChart(mime string) {
	s.addDimToChart(reqByMimeType.ID, "mime_type_"+mime, mime)
}

func (s *SquidLog) addDimToServerAddressChart(address string) {
	s.addDimToChart(reqByServer.ID, "server_address_"+address, address)
}

func (s *SquidLog) addDimToHierCodeChart(code string) {
	s.addDimToChart(reqByHierCode.ID, "hier_code_"+code, code)
}

func (s *SquidLog) addDimToReqMethodChart(method string) {
	s.addDimToChart(reqByMethod.ID, "req_method_"+method, method)
}

func (s *SquidLog) addDimToRespCodesChart(code string) {
	s.addDimToChart(respCodes.ID, "http_code_"+code, code)
}

func (s *SquidLog) addDimToCacheCodeChart(code string) {
	s.addDimToChart(cacheCode.ID, "cache_code_"+code, code)
}

func (s *SquidLog) addDimToChart(chartID, dimID, dimName string) {
	chart := s.Charts().Get(chartID)
	if chart == nil {
		s.Warningf("add '%s' dim: couldn't find '%s' chart", dimID, chartID)
		return
	}
	dim := &Dim{ID: dimID, Name: dimName, Algo: module.Incremental}
	if err := chart.AddDim(dim); err != nil {
		s.Warningf("add '%s' dim: %v", dimID, err)
		return
	}
	chart.MarkNotCreated()
}
