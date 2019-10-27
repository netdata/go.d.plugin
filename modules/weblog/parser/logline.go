package parser

//
//import (
//	"fmt"
//	"regexp"
//	"strconv"
//	"strings"
//)
//
//type (
//	LogLine struct {
//		vhost            string
//		client       string
//		reqHTTPMethod    string
//		reqURI           string
//		reqHTTPVersion   string
//		respCode   int
//		reqSize          int
//		respSize         int
//		respTime         float64
//		upstreamRespTime float64
//		custom           string
//	}
//)
//
//// nginx: http://nginx.org/en/docs/varindex.html
//// apache: http://httpd.apache.org/docs/current/mod/mod_log_config.html#logformat
//
//// TODO: do we really we need "custom" :thinking:
///*
//| name               | nginx                   | apache           |
//|--------------------|-------------------------|------------------|
//| vhost              | $server_name            | %v               | name of the server which accepted a request
//| client_addr        | $remote_addr            | %a (%h)          | apache %h: logs the IP address if HostnameLookups is Off
//| request            | $request                | %r               | req_method + req_uri + req_protocol
//| req_method         | $request_method         | %m               |
//| req_uri            | $request_uri            | %U               | nginx: w/ queries, apache: w/o
//| req_protocol       | $server_protocol        | %H               | request protocol, usually “HTTP/1.0”, “HTTP/1.1”, or “HTTP/2.0”
//| resp_status        | $status                 | %s (%>s)         | response status
//| req_size           | $request_length         | $I               | request length (including request line, header, and request body), apache: need mod_logio
//| resp_size          | $bytes_sent             | %O               | number of bytes sent to a client, including headers
//| resp_size          | $body_bytes_sent        | %B               | number of bytes sent to a client, not including headers
//| resp_time          | $request_time           | %D               | the time taken to serve the request. Apache: in microseconds, nginx: in seconds with a milliseconds resolution
//| upstream_resp_time | $upstream_response_time | -                | keeps time spent on receiving the response from the upstream server; the time is kept in seconds with millisecond resolution. Times of several responses are separated by commas and colons
//| custom             | -                       | -                |
//*/
//
//const (
//	fieldVhost            = "vhost"
//	fieldClientAddr       = "client_addr"
//	fieldRequest          = "request"
//	fieldReqMethod        = "req_method"
//	fieldReqURI           = "req_uri"
//	fieldReqProtocol      = "req_protocol"
//	fieldRespStatus       = "resp_status"
//	fieldReqSize          = "req_size"
//	fieldRespSize         = "resp_size"
//	fieldRespTime         = "resp_time"
//	fieldUpstreamRespTime = "upstream_resp_time"
//	fieldCustom           = "custom"
//)
//
//const (
//	EmptyString = "__empty_string__"
//	EmptyNumber = -9999
//)
//
//var (
//	// TODO: reClientAddr doesnt work with %h and HostnameLookups is On.
//	reVhost          = regexp.MustCompile(`^[a-zA-Z0-9.-:]+$`)
//	reClientAddr     = regexp.MustCompile(`^([\da-f.:]+|localhost)$`)
//	reReqHTTPMethod  = regexp.MustCompile(`^[A-Z]+$`)
//	reURI            = regexp.MustCompile(`^/[^\s]*$`)
//	reReqHTTPVersion = regexp.MustCompile(`^\d+(\.\d+)?$`)
//
//	emptyLogLine = LogLine{
//		vhost:            EmptyString,
//		client:       EmptyString,
//		reqHTTPMethod:    EmptyString,
//		reqURI:           EmptyString,
//		reqHTTPVersion:   EmptyString,
//		custom:           EmptyString,
//		respCode:   EmptyNumber,
//		reqSize:          EmptyNumber,
//		respSize:         EmptyNumber,
//		respTime:         EmptyNumber,
//		upstreamRespTime: EmptyNumber,
//	}
//)
//
//func (l LogLine) hasVhost() bool { return l.vhost != EmptyString }
//
//func (l LogLine) hasClientAddr() bool { return l.client != EmptyString }
//
//func (l LogLine) hasReqHTTPMethod() bool { return l.reqHTTPMethod != EmptyString }
//
//func (l LogLine) hasReqURI() bool { return l.reqURI != EmptyString }
//
//func (l LogLine) hasReqHTTPVersion() bool { return l.reqHTTPVersion != EmptyString }
//
//func (l LogLine) hasRespCodeStatus() bool { return l.respCode != EmptyNumber }
//
//func (l LogLine) hasReqSize() bool { return l.reqSize != EmptyNumber }
//
//func (l LogLine) hasRespSize() bool { return l.respSize != EmptyNumber }
//
//func (l LogLine) hasRespTime() bool { return l.respTime != EmptyNumber }
//
//func (l LogLine) hasUpstreamRespTime() bool { return l.upstreamRespTime != EmptyNumber }
//
//func (l LogLine) hasCustom() bool { return l.custom != EmptyString }
//
//func (l LogLine) verify() error {
//	err := l.verifyMandatoryFields()
//	if err != nil {
//		return err
//	}
//	return l.verifyOptionalFields()
//}
//
//func (l LogLine) verifyMandatoryFields() error {
//	if !l.hasRespCodeStatus() {
//		return fmt.Errorf("missing mandatory field: %s", fieldRespStatus)
//	}
//	if l.respCode < 100 || l.respCode >= 600 {
//		return fmt.Errorf("invalid '%s' field: %d", fieldRespStatus, l.respCode)
//	}
//	return nil
//}
//
//func (l LogLine) verifyOptionalFields() error {
//	if l.hasVhost() && !reVhost.MatchString(l.vhost) {
//		return fmt.Errorf("invalid '%s' field: %s", fieldVhost, l.vhost)
//	}
//	if l.hasClientAddr() && !reClientAddr.MatchString(l.client) {
//		return fmt.Errorf("invalid  '%s' field: %s", fieldClientAddr, l.client)
//	}
//	if l.hasReqHTTPMethod() && !reReqHTTPMethod.MatchString(l.reqHTTPMethod) {
//		return fmt.Errorf("invalid '%s' field: %s", fieldReqMethod, l.reqHTTPMethod)
//	}
//	if l.hasReqURI() && !reURI.MatchString(l.reqURI) {
//		return fmt.Errorf("invalid '%s' field: %s", fieldReqURI, l.reqURI)
//	}
//	if l.hasReqHTTPVersion() && !reReqHTTPVersion.MatchString(l.reqHTTPVersion) {
//		return fmt.Errorf("invalid '%s' field: %s", fieldReqProtocol, l.reqHTTPVersion)
//	}
//	if l.hasReqSize() && l.reqSize < 0 {
//		return fmt.Errorf("invalid '%s' field: %d", fieldReqSize, l.reqSize)
//	}
//	if l.hasRespSize() && l.respSize < 0 {
//		return fmt.Errorf("invalid '%s' field: %d", fieldRespSize, l.respSize)
//	}
//	if l.hasRespTime() && l.respTime < 0 {
//		return fmt.Errorf("invalid '%s' field: %f", fieldRespTime, l.respTime)
//	}
//	if l.hasUpstreamRespTime() && l.upstreamRespTime < 0 {
//		return fmt.Errorf("invalid '%s' field: %f", fieldUpstreamRespTime, l.upstreamRespTime)
//	}
//	return nil
//}
//
//func (l *LogLine) assign(field string, value string, timeMultiplier float64) (err error) {
//	switch field {
//	case fieldVhost:
//		l.vhost = value
//	case fieldClientAddr:
//		l.client = value
//	case fieldRequest:
//		err = l.assignRequest(value)
//	case fieldReqMethod:
//		l.reqHTTPMethod = value
//	case fieldReqURI:
//		l.reqURI = value
//	case fieldReqProtocol:
//		err = l.assignReqHTTPVersion(value)
//	case fieldRespStatus:
//		err = l.assignReqCodeStatus(value)
//	case fieldRespSize:
//		err = l.assignRespSize(value)
//	case fieldReqSize:
//		err = l.assignReqSize(value)
//	case fieldRespTime:
//		err = l.assignRespTime(value, timeMultiplier)
//	case fieldUpstreamRespTime:
//		err = l.assignUpstreamRespTime(value, timeMultiplier)
//	case fieldCustom:
//		l.custom = value
//	}
//	return err
//}
//
//func (l *LogLine) assignRequest(request string) error {
//	if request == "-" {
//		return nil
//	}
//	req := request
//	idx := strings.IndexByte(req, ' ')
//	if idx < 0 {
//		return fmt.Errorf("invalid request: %q", request)
//	}
//	l.reqHTTPMethod = req[0:idx]
//	req = req[idx+1:]
//
//	idx = strings.IndexByte(req, ' ')
//	if idx < 0 {
//		return fmt.Errorf("invalid request: %q", request)
//	}
//	l.reqURI = req[0:idx]
//	req = req[idx+1:]
//
//	return l.assignReqHTTPVersion(req)
//}
//
//func (l *LogLine) assignReqHTTPVersion(proto string) error {
//	if len(proto) <= 5 || !strings.HasPrefix(proto, "HTTP/") {
//		return fmt.Errorf("invalid protocol: %q", proto)
//	}
//	l.reqHTTPVersion = proto[5:]
//	return nil
//}
//
//func (l *LogLine) assignReqCodeStatus(status string) error {
//	if status == "-" {
//		return nil
//	}
//	var err error
//	l.respCode, err = strconv.Atoi(status)
//	if err != nil {
//		return fmt.Errorf("invalid status: %q: %w", status, err)
//	}
//	return nil
//}
//
//func (l *LogLine) assignReqSize(size string) error {
//	if size == "-" {
//		l.reqSize = 0
//		return nil
//	}
//	var err error
//	l.reqSize, err = strconv.Atoi(size)
//	if err != nil {
//		return fmt.Errorf("invalid request size: %q: %w", size, err)
//	}
//	return nil
//}
//
//func (l *LogLine) assignRespSize(size string) error {
//	if size == "-" {
//		l.respSize = 0
//		return nil
//	}
//	var err error
//	l.respSize, err = strconv.Atoi(size)
//	if err != nil {
//		return fmt.Errorf("invalid response size: %q: %w", size, err)
//	}
//	return nil
//}
//
//func (l *LogLine) assignRespTime(time string, timeScale float64) error {
//	if time == "-" {
//		return nil
//	}
//	val, err := strconv.ParseFloat(time, 64)
//	if err != nil {
//		return fmt.Errorf("invalid response time: %q: %w", time, err)
//	}
//	l.respTime = val * timeScale
//	return nil
//}
//
//func (l *LogLine) assignUpstreamRespTime(time string, timeScale float64) error {
//	if time == "-" {
//		return nil
//	}
//	if idx := strings.IndexByte(time, ','); idx >= 0 {
//		time = time[0:idx]
//	}
//	val, err := strconv.ParseFloat(time, 64)
//	if err != nil {
//		return fmt.Errorf("invalid upstream response time: %q: %w", time, err)
//	}
//	l.upstreamRespTime = val * timeScale
//	return nil
//}
