package weblog

import (
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
| scheme             | $scheme                 | -         | request scheme, “http” or “https”
| client             | $remote_addr            | %a (%h)   | apache %h: logs the IP address if HostnameLookups is Off
| request            | $request                | %r        | req_method + req_uri + req_protocol
| req_method         | $request_method         | %m        |
| req_uri            | $request_uri            | %U        | nginx: w/ queries, apache: w/o
| req_proto          | $server_protocol        | %H        | request protocol, usually “HTTP/1.0”, “HTTP/1.1”, or “HTTP/2.0”
| resp_status        | $status                 | %s (%>s)  | response status
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
	fieldVhostWithPort = "vhost:port"
	fieldScheme        = "scheme"
	fieldClient        = "client"
	fieldRequest       = "request"
	fieldReqMethod     = "req_method"
	fieldReqURI        = "req_uri"
	fieldReqProto      = "req_proto"
	fieldStatus        = "status"
	fieldReqSize       = "req_size"
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
	"scheme":                 fieldScheme,
	"remote_addr":            fieldClient,
	"request_method":         fieldReqMethod,
	"request_uri":            fieldReqURI,
	"server_protocol":        fieldReqProto,
	"status":                 fieldStatus,
	"request_length":         fieldReqSize,
	"bytes_sent":             fieldRespSize,
	"body_bytes_sent":        fieldRespSize,
	"request_time":           fieldRespTime,
	"upstream_response_time": fieldUpsRespTime,
	"custom":                 fieldCustom,
	"v":                      fieldVhost,
	"p":                      fieldPort,
	"v:%p":                   fieldVhostWithPort,
	"a":                      fieldClient,
	"h":                      fieldClient,
	"m":                      fieldReqMethod,
	"U":                      fieldReqURI,
	"H":                      fieldReqProto,
	"s":                      fieldStatus,
	">s":                     fieldStatus,
	"I":                      fieldReqSize,
	"O":                      fieldRespSize,
	"B":                      fieldRespSize,
	"D":                      fieldRespTime,
}

func newEmptyLogLine() *LogLine {
	var l LogLine
	l.reset()
	return &l
}

type LogLine struct {
	Vhost            string
	Scheme           string
	Port             string // Apache has no $scheme, this is workaround to collect per scheme requests. Lame.
	ClientAddr       string
	ReqHTTPMethod    string
	ReqURI           string
	ReqHTTPVersion   string
	RespCode         int
	ReqSize          int
	RespSize         int
	RespTime         float64
	UpstreamRespTime float64
	Custom           string

	timeScale float64
}

func (l *LogLine) Assign(variable string, value string) (err error) {
	field, ok := fieldsMapping[variable]
	if !ok {
		return
	}

	switch field {
	default:
		err = fmt.Errorf("unknown field : %s", field)
	case fieldVhost:
		l.Vhost = value
	case fieldPort:
		l.Port = value
	case fieldVhostWithPort:
		err = l.assignVhostWithPort(value)
	case fieldScheme:
		l.Scheme = value
	case fieldClient:
		l.ClientAddr = value
	case fieldRequest:
		err = l.assignRequest(value)
	case fieldReqMethod:
		l.ReqHTTPMethod = value
	case fieldReqURI:
		l.ReqURI = value
	case fieldReqProto:
		err = l.assignReqHTTPVersion(value)
	case fieldStatus:
		err = l.assignReqCodeStatus(value)
	case fieldRespSize:
		err = l.assignRespSize(value)
	case fieldReqSize:
		err = l.assignReqSize(value)
	case fieldRespTime:
		err = l.assignRespTime(value)
	case fieldUpsRespTime:
		err = l.assignUpstreamRespTime(value)
	case fieldCustom:
		l.Custom = value
	}
	return err
}

func (l *LogLine) assignVhostWithPort(vhostPort string) error {
	idx := strings.LastIndexByte(vhostPort, ':')
	if idx == -1 {
		return fmt.Errorf("invalid vhost with port: %q", vhostPort)
	}
	l.Vhost = vhostPort[0:idx]
	l.Port = vhostPort[idx+1:]
	return nil
}

func (l *LogLine) assignRequest(request string) error {
	if request == "-" {
		return nil
	}
	req := request
	idx := strings.IndexByte(req, ' ')
	if idx < 0 {
		return fmt.Errorf("invalid request: %q", request)
	}
	l.ReqHTTPMethod = req[0:idx]
	req = req[idx+1:]

	idx = strings.IndexByte(req, ' ')
	if idx < 0 {
		return fmt.Errorf("invalid request: %q", request)
	}
	l.ReqURI = req[0:idx]
	req = req[idx+1:]

	return l.assignReqHTTPVersion(req)
}

func (l *LogLine) assignReqHTTPVersion(proto string) error {
	if len(proto) <= 5 || !strings.HasPrefix(proto, "HTTP/") {
		return fmt.Errorf("invalid protocol: %q", proto)
	}
	l.ReqHTTPVersion = proto[5:]
	return nil
}

func (l *LogLine) assignReqCodeStatus(status string) error {
	if status == "-" {
		return nil
	}
	var err error
	l.RespCode, err = strconv.Atoi(status)
	if err != nil {
		return fmt.Errorf("invalid status: %q: %w", status, err)
	}
	return nil
}

func (l *LogLine) assignReqSize(size string) error {
	if size == "-" {
		l.ReqSize = 0
		return nil
	}
	var err error
	l.ReqSize, err = strconv.Atoi(size)
	if err != nil {
		return fmt.Errorf("invalid request size: %q: %w", size, err)
	}
	return nil
}

func (l *LogLine) assignRespSize(size string) error {
	if size == "-" {
		l.RespSize = 0
		return nil
	}
	var err error
	l.RespSize, err = strconv.Atoi(size)
	if err != nil {
		return fmt.Errorf("invalid response size: %q: %w", size, err)
	}
	return nil
}

func (l *LogLine) assignRespTime(time string) error {
	if time == "-" {
		return nil
	}
	val, err := strconv.ParseFloat(time, 64)
	if err != nil {
		return fmt.Errorf("invalid response time: %q: %w", time, err)
	}
	l.RespTime = val * l.timeScale
	return nil
}

func (l *LogLine) assignUpstreamRespTime(time string) error {
	if time == "-" {
		return nil
	}
	if idx := strings.IndexByte(time, ','); idx >= 0 {
		time = time[0:idx]
	}
	val, err := strconv.ParseFloat(time, 64)
	if err != nil {
		return fmt.Errorf("invalid upstream response time: %q: %w", time, err)
	}
	l.UpstreamRespTime = val * l.timeScale
	return nil
}

func (l LogLine) Verify() error {
	if !l.hasRespCode() {
		return fmt.Errorf("missing mandatory field: %s", fieldStatus)
	}
	if !l.validRespCode() {
		return fmt.Errorf("invalid '%s' field: %d", fieldStatus, l.RespCode)
	}

	// optional checks
	if l.hasVhost() && !l.validVhost() {
		return fmt.Errorf("invalid '%s' field: %s", fieldVhost, l.Vhost)
	}
	if l.hasPort() && !l.validPort() {
		return fmt.Errorf("invalid '%s' field: %s", fieldPort, l.Port)
	}
	if l.hasScheme() && !l.validScheme() {
		return fmt.Errorf("invalid '%s' field: %s", fieldScheme, l.Scheme)
	}
	if l.hasClientAddr() && !l.validClientAddr() {
		return fmt.Errorf("invalid  '%s' field: %s", fieldClient, l.ClientAddr)
	}
	if l.hasReqHTTPMethod() && !l.validHTTPMethod() {
		return fmt.Errorf("invalid '%s' field: %s", fieldReqMethod, l.ReqHTTPMethod)
	}
	if l.hasReqURI() && !l.validReqURI() {
		return fmt.Errorf("invalid '%s' field: %s", fieldReqURI, l.ReqURI)
	}
	if l.hasReqHTTPVersion() && !l.validHTTPVersion() {
		return fmt.Errorf("invalid '%s' field: %s", fieldReqProto, l.ReqHTTPVersion)
	}
	if l.hasReqSize() && !l.validReqSize() {
		return fmt.Errorf("invalid '%s' field: %d", fieldReqSize, l.ReqSize)
	}
	if l.hasRespSize() && !l.validRespSize() {
		return fmt.Errorf("invalid '%s' field: %d", fieldRespSize, l.RespSize)
	}
	if l.hasRespTime() && !l.validRespTime() {
		return fmt.Errorf("invalid '%s' field: %f", fieldRespTime, l.RespTime)
	}
	if l.hasUpstreamRespTime() && !l.validUpstreamRespTime() {
		return fmt.Errorf("invalid '%s' field: %f", fieldUpsRespTime, l.UpstreamRespTime)
	}
	return nil
}

func (l LogLine) hasVhost() bool { return !isEmptyString(l.Vhost) }

func (l LogLine) hasPort() bool { return !isEmptyString(l.Port) }

func (l LogLine) hasScheme() bool { return !isEmptyString(l.Scheme) }

func (l LogLine) hasClientAddr() bool { return !isEmptyString(l.ClientAddr) }

func (l LogLine) hasReqHTTPMethod() bool { return !isEmptyString(l.ReqHTTPMethod) }

func (l LogLine) hasReqURI() bool { return !isEmptyString(l.ReqURI) }

func (l LogLine) hasReqHTTPVersion() bool { return !isEmptyString(l.ReqHTTPVersion) }

func (l LogLine) hasRespCode() bool { return !isEmptyNumber(l.RespCode) }

func (l LogLine) hasReqSize() bool { return !isEmptyNumber(l.ReqSize) }

func (l LogLine) hasRespSize() bool { return !isEmptyNumber(l.RespSize) }

func (l LogLine) hasRespTime() bool { return !isEmptyNumber(int(l.RespTime)) }

func (l LogLine) hasUpstreamRespTime() bool { return !isEmptyNumber(int(l.UpstreamRespTime)) }

func (l LogLine) hasCustom() bool { return !isEmptyString(l.Custom) }

func isEmptyString(s string) bool { return s == emptyString || s == "" }

func isEmptyNumber(n int) bool { return n == emptyNumber }

var (
	// TODO: reClientAddr doesnt work with %h when HostnameLookups is On.
	reVhost          = regexp.MustCompile(`^[a-zA-Z0-9.-:]+$`)
	reClientAddr     = regexp.MustCompile(`^([\da-f.:]+|localhost)$`)
	reReqHTTPMethod  = regexp.MustCompile(`^GET|HEAD|POST|PUT|PATCH|DELETE|CONNECT|OPTIONS|TRACE$`)
	reURI            = regexp.MustCompile(`^/[^\s]*$`)
	reReqHTTPVersion = regexp.MustCompile(`^\d+(\.\d+)?$`)
)

func (l LogLine) validVhost() bool { return reVhost.MatchString(l.Vhost) }

func (l LogLine) validScheme() bool { return l.Scheme == "http" || l.Scheme == "https" }

func (l LogLine) validClientAddr() bool { return reClientAddr.MatchString(l.ClientAddr) }

func (l LogLine) validHTTPMethod() bool { return reReqHTTPMethod.MatchString(l.ReqHTTPMethod) }

func (l LogLine) validReqURI() bool { return reURI.MatchString(l.ReqURI) }

func (l LogLine) validHTTPVersion() bool { return reReqHTTPVersion.MatchString(l.ReqHTTPVersion) }

func (l LogLine) validReqSize() bool { return l.ReqSize >= 0 }

func (l LogLine) validRespSize() bool { return l.RespSize >= 0 }

func (l LogLine) validRespTime() bool { return l.RespTime >= 0 }

func (l LogLine) validUpstreamRespTime() bool { return l.UpstreamRespTime > 0 }

func (l LogLine) validRespCode() bool { return l.RespCode >= 100 || l.RespCode <= 600 }

func (l LogLine) validPort() bool { v, err := strconv.Atoi(l.Port); return err != nil && v > 80 }

const (
	emptyString = "__empty_string__"
	emptyNumber = -9999
)

func (l *LogLine) reset() {
	l.Vhost = emptyString
	l.Port = emptyString
	l.Scheme = emptyString
	l.ClientAddr = emptyString
	l.ReqHTTPMethod = emptyString
	l.ReqURI = emptyString
	l.ReqHTTPVersion = emptyString
	l.RespCode = emptyNumber
	l.ReqSize = emptyNumber
	l.RespSize = emptyNumber
	l.RespTime = emptyNumber
	l.UpstreamRespTime = emptyNumber
	l.Custom = emptyString
}
