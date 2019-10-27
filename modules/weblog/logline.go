package weblog

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// nginx: http://nginx.org/en/docs/varindex.html
// apache: http://httpd.apache.org/docs/current/mod/mod_log_config.html#logformat

// TODO: do we really we need "custom" :thinking:
/*
| name               | nginx                   | apache    |
|--------------------|-------------------------|-----------|
| vhost              | $host ($http_host)      | %v        | name of the server which accepted a request
| port               | $server_port            | %p        | port of the server which accepted a request
| req_scheme         | $scheme                 | -         | request scheme, “http” or “https”
| req_client         | $remote_addr            | %a (%h)   | apache %h: logs the IP address if HostnameLookups is Off
| request            | $request                | %r        | req_method + req_uri + req_protocol
| req_method         | $request_method         | %m        |
| req_uri            | $request_uri            | %U        | nginx: w/ queries, apache: w/o
| req_proto          | $server_protocol        | %H        | request protocol, usually “HTTP/1.0”, “HTTP/1.1”, or “HTTP/2.0”
| resp_status        | $status                 | %s (%>s)  | response respStatus
| req_size           | $request_length         | $I        | request length (including request line, header, and request body), apache: need mod_logio
| resp_size          | $bytes_sent             | %O        | number of bytes sent to a client, including headers
| resp_size          | $body_bytes_sent        | %B        | number of bytes sent to a client, not including headers
| req_time           | $request_time           | %D        | the time taken to serve the request. Apache: in microseconds, nginx: in seconds with a milliseconds resolution
| ups_resp_time      | $upstream_response_time | -         | keeps time spent on receiving the response from the upstream server; the time is kept in seconds with millisecond resolution. Times of several responses are separated by commas and colons
| custom             | -                       | -         |
*/

const (
	fieldVhost         = "vhost"
	fieldPort          = "port"
	fieldVhostWithPort = "vhost_port"
	fieldReqScheme     = "req_scheme"
	fieldReqClient     = "req_client"
	fieldRequest       = "request"
	fieldReqMethod     = "req_method"
	fieldReqURI        = "req_uri"
	fieldReqProto      = "req_proto"
	fieldReqSize       = "req_size"
	fieldRespStatus    = "resp_status"
	fieldRespSize      = "resp_size"
	fieldRespTime      = "resp_time"
	fieldUpsRespTime   = "ups_resp_time"
	fieldCustom        = "custom"
)

var fieldsMapping = map[string]string{
	"host":                   fieldVhost,
	"http_host":              fieldVhost,
	"server_port":            fieldPort,
	"host:$server_port":      fieldVhostWithPort,
	"scheme":                 fieldReqScheme,
	"remote_addr":            fieldReqClient,
	"request":                fieldRequest,
	"request_method":         fieldReqMethod,
	"request_uri":            fieldReqURI,
	"server_protocol":        fieldReqProto,
	"status":                 fieldRespStatus,
	"request_length":         fieldReqSize,
	"bytes_sent":             fieldRespSize,
	"body_bytes_sent":        fieldRespSize,
	"request_time":           fieldRespTime,
	"upstream_response_time": fieldUpsRespTime,
	"custom":                 fieldCustom,
	"v":                      fieldVhost,
	"p":                      fieldPort,
	"v:%p":                   fieldVhostWithPort,
	"a":                      fieldReqClient,
	"h":                      fieldReqClient,
	"r":                      fieldRequest,
	"m":                      fieldReqMethod,
	"U":                      fieldReqURI,
	"H":                      fieldReqProto,
	"s":                      fieldRespStatus,
	">s":                     fieldRespStatus,
	"I":                      fieldReqSize,
	"O":                      fieldRespSize,
	"B":                      fieldRespSize,
	"D":                      fieldRespTime,
}

var (
	errMandatoryField      = errors.New("missing mandatory field")
	errUnknownField        = errors.New("unknown field")
	errBadVhost            = errors.New("bad vhost")
	errBadVhostPort        = errors.New("bad vhost with port")
	errBadPort             = errors.New("bad port")
	errBadReqScheme        = errors.New("bad req scheme")
	errBadReqClient        = errors.New("bad req client")
	errBadRequest          = errors.New("bad request")
	errBadReqMethod        = errors.New("bad req method")
	errBadReqURI           = errors.New("bad req uri")
	errBadReqProto         = errors.New("bad req protocol")
	errBadReqSize          = errors.New("bad req size")
	errBadRespStatus       = errors.New("bad resp status")
	errBadRespSize         = errors.New("bad resp size")
	errBadRespTime         = errors.New("bad resp time")
	errBadUpstreamRespTime = errors.New("bad upstream resp time")
)

func newEmptyLogLine() *logLine {
	var l logLine
	l.reset()
	return &l
}

type logLine struct {
	vhost string
	port  string // Apache has no $scheme, this is workaround to collect per scheme requests. Lame.

	reqScheme string
	reqClient string
	reqMethod string
	reqURI    string
	reqProto  string
	reqSize   int

	respStatus int
	respSize   int
	respTime   float64

	upsRespTime float64

	custom string
}

func (l *logLine) Assign(variable string, value string) (err error) {
	if value == "" {
		return
	}
	field, ok := fieldsMapping[variable]
	if !ok {
		return
	}

	switch field {
	default:
		err = fmt.Errorf("assign '%s': %w", field, errUnknownField)
	case fieldVhost:
		l.vhost = value
	case fieldPort:
		l.port = value
	case fieldVhostWithPort:
		err = l.assignVhostWithPort(value)
	case fieldReqScheme:
		l.reqScheme = value
	case fieldReqClient:
		l.reqClient = value
	case fieldRequest:
		err = l.assignRequest(value)
	case fieldReqMethod:
		l.reqMethod = value
	case fieldReqURI:
		l.reqURI = value
	case fieldReqProto:
		err = l.assignReqProto(value)
	case fieldRespStatus:
		err = l.assignRespStatus(value)
	case fieldRespSize:
		err = l.assignRespSize(value)
	case fieldReqSize:
		err = l.assignReqSize(value)
	case fieldRespTime:
		err = l.assignRespTime(value)
	case fieldUpsRespTime:
		err = l.assignUpstreamRespTime(value)
	case fieldCustom:
		l.custom = value
	}
	return err
}

func (l *logLine) assignVhostWithPort(vhostPort string) error {
	idx := strings.LastIndexByte(vhostPort, ':')
	if idx == -1 {
		return fmt.Errorf("assign '%s' : %w", vhostPort, errBadVhostPort)
	}
	l.vhost = vhostPort[0:idx]
	l.port = vhostPort[idx+1:]
	return nil
}

func (l *logLine) assignRequest(request string) error {
	if request == "-" {
		return nil
	}
	idx := strings.IndexByte(request, ' ')
	if idx < 0 {
		return fmt.Errorf("assign '%s': %w", request, errBadRequest)
	}
	// TODO: fail or continue?
	method := request[0:idx]
	if !isValidReqMethod(method) {
		return fmt.Errorf("assign '%s': %w", request, errBadRequest)
	}

	rest := request[idx+1:]
	idx = strings.IndexByte(rest, ' ')
	if idx < 0 {
		return fmt.Errorf("assign '%s': %w", request, errBadRequest)
	}
	uri := rest[0:idx]
	rest = rest[idx+1:]

	if err := l.assignReqProto(rest); err != nil {
		return fmt.Errorf("assign '%s': %w", request, errBadRequest)
	}
	l.reqMethod = method
	l.reqURI = uri
	return nil
}

func (l *logLine) assignReqProto(proto string) error {
	if len(proto) <= 5 || !strings.HasPrefix(proto, "HTTP/") {
		return fmt.Errorf("assign '%s': %w", proto, errBadReqProto)
	}
	l.reqProto = proto[5:]
	return nil
}

func (l *logLine) assignRespStatus(status string) error {
	if status == "-" {
		// TODO: hm?
		return nil
	}
	v, err := strconv.Atoi(status)
	if err != nil {
		return fmt.Errorf("assign '%s': %w", status, errBadRespStatus)
	}
	l.respStatus = v
	return nil
}

func (l *logLine) assignReqSize(size string) error {
	if size == "-" {
		l.reqSize = 0
		return nil
	}
	v, err := strconv.Atoi(size)
	if err != nil {
		return fmt.Errorf("assign '%s': %w", size, errBadReqSize)
	}
	l.reqSize = v
	return nil
}

func (l *logLine) assignRespSize(size string) error {
	if size == "-" {
		l.respSize = 0
		return nil
	}
	v, err := strconv.Atoi(size)
	if err != nil {
		return fmt.Errorf("assign '%s': %w", size, errBadRespSize)
	}
	l.respSize = v
	return nil
}

func (l *logLine) assignRespTime(time string) error {
	if time == "-" {
		return nil
	}
	val, err := strconv.ParseFloat(time, 64)
	if err != nil {
		return fmt.Errorf("assign '%s': %w", time, errBadRespTime)
	}
	l.respTime = val * respTimeMultiplier(time)
	return nil
}

func (l *logLine) assignUpstreamRespTime(time string) error {
	if time == "-" {
		return nil
	}
	if idx := strings.IndexByte(time, ','); idx >= 0 {
		time = time[0:idx]
	}
	val, err := strconv.ParseFloat(time, 64)
	if err != nil {
		return fmt.Errorf("assign '%s': %w", time, errBadUpstreamRespTime)
	}
	l.upsRespTime = val * respTimeMultiplier(time)
	return nil
}

func (l logLine) verify() error {
	if !l.hasRespStatus() {
		return fmt.Errorf("%s: %w", fieldRespStatus, errMandatoryField)
	}
	if !l.validRespStatus() {
		return fmt.Errorf("verify '%d': %w", l.respStatus, errBadRespStatus)
	}

	// optional checks
	if l.hasVhost() && !l.validVhost() {
		return fmt.Errorf("verify '%s': %w", l.vhost, errBadVhost)
	}
	if l.hasPort() && !l.validPort() {
		return fmt.Errorf("verify '%s': %w", l.port, errBadPort)
	}
	if l.hasReqScheme() && !l.validReqScheme() {
		return fmt.Errorf("verify '%s': %w", l.reqScheme, errBadReqScheme)
	}
	if l.hasReqClient() && !l.validReqClient() {
		return fmt.Errorf("verify '%s': %w", l.reqClient, errBadReqClient)
	}
	if l.hasReqMethod() && !l.validReqMethod() {
		return fmt.Errorf("verify '%s': %w", l.reqMethod, errBadReqMethod)
	}
	if l.hasReqURI() && !l.validReqURI() {
		return fmt.Errorf("verify '%s': %w", l.reqURI, errBadReqURI)
	}
	if l.hasReqProto() && !l.validReqProto() {
		return fmt.Errorf("verify '%s': %w", l.reqProto, errBadReqProto)
	}
	if l.hasReqSize() && !l.validReqSize() {
		return fmt.Errorf("verify '%d': %w", l.reqSize, errBadReqSize)
	}
	if l.hasRespSize() && !l.validRespSize() {
		return fmt.Errorf("verify '%d': %w", l.respSize, errBadRespSize)
	}
	if l.hasRespTime() && !l.validRespTime() {
		return fmt.Errorf("verify '%f': %w", l.respTime, errBadRespTime)
	}
	if l.hasUpstreamRespTime() && !l.validUpstreamRespTime() {
		return fmt.Errorf("verify '%f': %w", l.upsRespTime, errBadUpstreamRespTime)
	}
	return nil
}

func (l logLine) hasVhost() bool { return !isEmptyString(l.vhost) }

func (l logLine) hasPort() bool { return !isEmptyString(l.port) }

func (l logLine) hasReqScheme() bool { return !isEmptyString(l.reqScheme) }

func (l logLine) hasReqClient() bool { return !isEmptyString(l.reqClient) }

func (l logLine) hasReqMethod() bool { return !isEmptyString(l.reqMethod) }

func (l logLine) hasReqURI() bool { return !isEmptyString(l.reqURI) }

func (l logLine) hasReqProto() bool { return !isEmptyString(l.reqProto) }

func (l logLine) hasRespStatus() bool { return !isEmptyNumber(l.respStatus) }

func (l logLine) hasReqSize() bool { return !isEmptyNumber(l.reqSize) }

func (l logLine) hasRespSize() bool { return !isEmptyNumber(l.respSize) }

func (l logLine) hasRespTime() bool { return !isEmptyNumber(int(l.respTime)) }

func (l logLine) hasUpstreamRespTime() bool { return !isEmptyNumber(int(l.upsRespTime)) }

func (l logLine) hasCustom() bool { return !isEmptyString(l.custom) }

var (
	// TODO: reClient doesnt work with %h when HostnameLookups is On.
	reVhost  = regexp.MustCompile(`^[a-zA-Z0-9-:.]+$`)
	reClient = regexp.MustCompile(`^([\da-f:.]+|localhost)$`)
	reURI    = regexp.MustCompile(`^/[^\s]*$`)
)

func (l logLine) validVhost() bool { return reVhost.MatchString(l.vhost) }

func (l logLine) validPort() bool { v, err := strconv.Atoi(l.port); return err == nil && v >= 80 }

func (l logLine) validReqScheme() bool { return l.reqScheme == "http" || l.reqScheme == "https" }

func (l logLine) validReqClient() bool { return reClient.MatchString(l.reqClient) }

func (l logLine) validReqMethod() bool { return isValidReqMethod(l.reqMethod) }

func (l logLine) validReqURI() bool { return reURI.MatchString(l.reqURI) }

func (l logLine) validReqProto() bool { return isValidReqProto(l.reqProto) }

func (l logLine) validReqSize() bool { return l.reqSize >= 0 }

func (l logLine) validRespSize() bool { return l.respSize >= 0 }

func (l logLine) validRespTime() bool { return l.respTime >= 0 }

func (l logLine) validUpstreamRespTime() bool { return l.upsRespTime >= 0 }

func (l logLine) validRespStatus() bool { return l.respStatus >= 100 && l.respStatus <= 600 }

func isEmptyString(s string) bool {
	return s == emptyString || s == ""
}

func isEmptyNumber(n int) bool {
	return n == emptyNumber
}

func isValidReqMethod(method string) bool {
	if method == "GET" {
		return true
	}
	switch method {
	case "HEAD", "POST", "PUT", "PATCH", "DELETE", "CONNECT", "OPTIONS", "TRACE":
		return true
	}
	return false
}

func isValidReqProto(version string) bool {
	if version == "1.1" {
		return true
	}
	switch version {
	case "1", "1.0", "2", "2.0":
		return true
	}
	return false
}

const (
	// Apache time is in microseconds, Nginx time is in seconds with a milliseconds resolution.
	apacheTimeMul = 1
	nginxTimeMul  = 1000000
)

func respTimeMultiplier(time string) float64 {
	if strings.IndexByte(time, '.') > 0 {
		return nginxTimeMul
	}
	return apacheTimeMul
}

const (
	emptyString = "__empty_string__"
	emptyNumber = -9999
)

func (l *logLine) reset() {
	l.vhost = emptyString
	l.port = emptyString
	l.reqScheme = emptyString
	l.reqClient = emptyString
	l.reqMethod = emptyString
	l.reqURI = emptyString
	l.reqProto = emptyString
	l.reqSize = emptyNumber
	l.respStatus = emptyNumber
	l.respSize = emptyNumber
	l.respTime = emptyNumber
	l.upsRespTime = emptyNumber
	l.custom = emptyString
}
