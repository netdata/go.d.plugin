package parser

import (
	"regexp"
	"strconv"
	"strings"

	"golang.org/x/xerrors"
)

type (
	LogLine struct {
		Vhost            string
		ClientAddr       string
		ReqMethod        string
		ReqURI           string
		ReqHTTPVersion   string
		RespCodeStatus   int
		ReqSize          int
		RespSize         int
		RespTime         float64
		UpstreamRespTime float64
		Custom           string
	}
)

// nginx: http://nginx.org/en/docs/varindex.html
// apache: http://httpd.apache.org/docs/current/mod/mod_log_config.html#logformat

// TODO: do we really we need "custom" :thinking:
/*
| name               | nginx                   | apache           |
|--------------------|-------------------------|------------------|
| vhost              | $server_name            | %v               | name of the server which accepted a request
| client_addr        | $remote_addr            | %a (%h)          | apache %h: logs the IP address if HostnameLookups is Off
| request            | $request                | %r               | req_method + req_uri + req_protocol
| req_method         | $request_method         | %m               |
| req_uri            | $request_uri            | %U               | nginx: w/ queries, apache: w/o
| req_protocol       | $server_protocol        | %H               | request protocol, usually “HTTP/1.0”, “HTTP/1.1”, or “HTTP/2.0”
| resp_status        | $status                 | %s (%>s)         | response status
| req_size           | $request_length         | $I               | request length (including request line, header, and request body), apache: need mod_logio
| resp_size          | $bytes_sent             | %O               | number of bytes sent to a client, including headers
| resp_size          | $body_bytes_sent        | %B               | number of bytes sent to a client, not including headers
| req_time           | $request_time           | %D               | the time taken to serve the request. Apache: in microseconds, nginx: in seconds with a milliseconds resolution
| upstream_resp_time | $upstream_response_time | -                | keeps time spent on receiving the response from the upstream server; the time is kept in seconds with millisecond resolution. Times of several responses are separated by commas and colons
| custom             | -                       | -                |
*/

const (
	fieldVhost            = "vhost"
	fieldClientAddr       = "client_addr"
	fieldRequest          = "request"
	fieldReqMethod        = "req_method"
	fieldReqURI           = "req_uri"
	fieldReqProtocol      = "req_protocol"
	fieldRespStatus       = "resp_status"
	fieldReqSize          = "req_size"
	fieldRespSize         = "resp_size"
	fieldReqTime          = "req_time"
	fieldUpstreamRespTime = "upstream_resp_time"
	fieldCustom           = "custom"
)

const (
	EmptyString = "__empty_string__"
	EmptyNumber = -9999
)

var (
	reClient     = regexp.MustCompile(`^([\da-f.:]+|localhost)$`)
	reReqMethod  = regexp.MustCompile(`^[A-Z]+$`)
	reURI        = regexp.MustCompile(`^/[^\s]*$`)
	reReqVersion = regexp.MustCompile(`^\d+(\.\d+)?$`)
	reServerName = regexp.MustCompile(`^[a-zA-Z0-9.-:]+$`)

	emptyLogLine = LogLine{
		Vhost:            EmptyString,
		ClientAddr:       EmptyString,
		ReqMethod:        EmptyString,
		ReqURI:           EmptyString,
		ReqHTTPVersion:   EmptyString,
		Custom:           EmptyString,
		RespCodeStatus:   EmptyNumber,
		ReqSize:          EmptyNumber,
		RespSize:         EmptyNumber,
		RespTime:         EmptyNumber,
		UpstreamRespTime: EmptyNumber,
	}
)

func (l LogLine) HasVhost() bool { return l.Vhost != EmptyString }

func (l LogLine) HasClientAddr() bool { return l.ClientAddr != EmptyString }

func (l LogLine) HasReqMethod() bool { return l.ReqMethod != EmptyString }

func (l LogLine) HasReqURI() bool { return l.ReqURI != EmptyString }

func (l LogLine) HasReqHTTPVersion() bool { return l.ReqHTTPVersion != EmptyString }

func (l LogLine) HasRespCodeStatus() bool { return l.RespCodeStatus != EmptyNumber }

func (l LogLine) HasReqSize() bool { return l.ReqSize != EmptyNumber }

func (l LogLine) HasRespSize() bool { return l.RespSize != EmptyNumber }

func (l LogLine) HasRespTime() bool { return l.RespTime != EmptyNumber }

func (l LogLine) HasUpstreamRespTime() bool { return l.UpstreamRespTime != EmptyNumber }

func (l LogLine) HasCustom() bool { return l.Custom != EmptyString }

func (l LogLine) Verify() error {
	err := l.verifyMandatoryFields()
	if err != nil {
		return err
	}
	return l.verifyOptionalFields()
}

func (l LogLine) verifyMandatoryFields() error {
	if !l.HasRespCodeStatus() {
		return xerrors.New("missing mandatory field: status")
	}
	if l.RespCodeStatus < 100 || l.RespCodeStatus >= 600 {
		return xerrors.Errorf("invalid status field: %d", l.RespCodeStatus)
	}
	return nil
}

func (l LogLine) verifyOptionalFields() error {
	if l.HasVhost() && !reServerName.MatchString(l.Vhost) {
		return xerrors.Errorf("invalid vhost field: '%s'", l.Vhost)
	}
	if l.HasClientAddr() && !reClient.MatchString(l.ClientAddr) {
		return xerrors.Errorf("invalid client field: %q", l.ClientAddr)
	}
	if l.HasReqMethod() && !reReqMethod.MatchString(l.ReqMethod) {
		return xerrors.Errorf("invalid method field: '%s'", l.ReqMethod)
	}
	if l.HasReqURI() && !reURI.MatchString(l.ReqURI) {
		return xerrors.Errorf("invalid ReqURI field: '%s'", l.ReqURI)
	}
	if l.HasReqHTTPVersion() && !reReqVersion.MatchString(l.ReqHTTPVersion) {
		return xerrors.Errorf("invalid protocol field: '%s'", l.ReqHTTPVersion)
	}
	if l.HasReqSize() && l.ReqSize < 0 {
		return xerrors.Errorf("invalid request size field: %d", l.ReqSize)
	}
	if l.HasRespSize() && l.RespSize < 0 {
		return xerrors.Errorf("invalid response size field: %d", l.RespSize)
	}
	if l.HasRespTime() && l.RespTime < 0 {
		return xerrors.Errorf("invalid response time field: %f", l.RespTime)
	}
	if l.HasUpstreamRespTime() && l.UpstreamRespTime < 0 {
		return xerrors.Errorf("invalid upstream response time field: %f", l.UpstreamRespTime)
	}
	return nil
}

func (l *LogLine) assign(field string, value string, timeMultiplier float64) (err error) {
	switch field {
	case fieldVhost:
		l.Vhost = value
	case fieldClientAddr:
		l.ClientAddr = value
	case fieldRequest:
		err = l.assignRequest(value)
	case fieldReqMethod:
		l.ReqMethod = value
	case fieldReqURI:
		l.ReqURI = value
	case fieldReqProtocol:
		err = l.assignProtocol(value)
	case fieldRespStatus:
		err = l.assignStatus(value)
	case fieldRespSize:
		err = l.assignRespSize(value)
	case fieldReqSize:
		err = l.assignReqSize(value)
	case fieldReqTime:
		err = l.assignRespTime(value, timeMultiplier)
	case fieldUpstreamRespTime:
		err = l.assignUpstreamRespTime(value, timeMultiplier)
	case fieldCustom:
		l.Custom = value
	}
	return err
}

func (l *LogLine) assignRequest(request string) error {
	if request == "-" {
		return nil
	}
	req := request
	idx := strings.IndexByte(req, ' ')
	if idx < 0 {
		return xerrors.Errorf("invalid request: %q", request)
	}
	l.ReqMethod = req[0:idx]
	req = req[idx+1:]

	idx = strings.IndexByte(req, ' ')
	if idx < 0 {
		return xerrors.Errorf("invalid request: %q", request)
	}
	l.ReqURI = req[0:idx]
	req = req[idx+1:]

	return l.assignProtocol(req)
}

func (l *LogLine) assignProtocol(proto string) error {
	if len(proto) <= 5 || !strings.HasPrefix(proto, "HTTP/") {
		return xerrors.Errorf("invalid protocol: %q", proto)
	}
	l.ReqHTTPVersion = proto[5:]
	return nil
}

func (l *LogLine) assignStatus(status string) error {
	if status == "-" {
		return nil
	}
	var err error
	l.RespCodeStatus, err = strconv.Atoi(status)
	if err != nil {
		return xerrors.Errorf("invalid status: %q: %w", status, err)
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
		return xerrors.Errorf("invalid request size: %q: %w", size, err)
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
		return xerrors.Errorf("invalid response size: %q: %w", size, err)
	}
	return nil
}

func (l *LogLine) assignRespTime(time string, timeScale float64) error {
	if time == "-" {
		return nil
	}
	val, err := strconv.ParseFloat(time, 64)
	if err != nil {
		return xerrors.Errorf("invalid response time: %q: %w", time, err)
	}
	l.RespTime = val * timeScale
	return nil
}

func (l *LogLine) assignUpstreamRespTime(time string, timeScale float64) error {
	if time == "-" {
		return nil
	}
	if idx := strings.IndexByte(time, ','); idx >= 0 {
		time = time[0:idx]
	}
	val, err := strconv.ParseFloat(time, 64)
	if err != nil {
		return xerrors.Errorf("invalid upstream response time: %q: %w", time, err)
	}
	l.UpstreamRespTime = val * timeScale
	return nil
}
