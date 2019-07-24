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
		Client           string
		Method           string
		URI              string
		Version          string
		Status           int
		ReqSize          int
		RespSize         int
		RespTime         float64
		UpstreamRespTime float64
		Custom           string
	}
)

const (
	fieldClient           = "client"
	fieldRequest          = "request"
	fieldMethod           = "method"
	fieldURI              = "uri"
	fieldProtocol         = "protocol"
	fieldVersion          = "version"
	fieldStatus           = "status"
	fieldReqSize          = "req_size"
	fieldRespSize         = "resp_size"
	fieldRespTime         = "resp_time"
	fieldUpstreamRespTime = "upstream_resp_time"
	fieldVhost            = "vhost"
	fieldCustom           = "custom"
)

const (
	EmptyString = "__empty_string__"
	EmptyNumber = -9999
)

var (
	reClient  = regexp.MustCompile(`^([\da-f.:]+|localhost)$`)
	reMethod  = regexp.MustCompile(`^[A-Z]+$`)
	reURI     = regexp.MustCompile(`^/[^\s]*$`)
	reVersion = regexp.MustCompile(`^\d+(\.\d+)?$`)
	reVhost   = regexp.MustCompile(`^[a-zA-Z0-9.-:]+$`)

	emptyLogLine = LogLine{
		Vhost:            EmptyString,
		Client:           EmptyString,
		Method:           EmptyString,
		URI:              EmptyString,
		Version:          EmptyString,
		Custom:           EmptyString,
		Status:           EmptyNumber,
		ReqSize:          EmptyNumber,
		RespSize:         EmptyNumber,
		RespTime:         EmptyNumber,
		UpstreamRespTime: EmptyNumber,
	}
)

func (l *LogLine) Verify() error {
	if l.Client == EmptyString || l.Status == EmptyNumber || l.RespSize == EmptyNumber || l.Method == EmptyString || l.URI == EmptyString {
		return xerrors.New("missing some mandatory fields: client, status,resp_size, method, uri")
	}
	if l.Vhost != EmptyString && !reVhost.MatchString(l.Vhost) {
		return xerrors.Errorf("invalid vhost field: '%s'", l.Vhost)
	}
	if !reClient.MatchString(l.Client) {
		return xerrors.Errorf("invalid client field: %q", l.Client)
	}
	if !reMethod.MatchString(l.Method) {
		return xerrors.Errorf("invalid method field: '%s'", l.Method)
	}
	if !reURI.MatchString(l.URI) {
		return xerrors.Errorf("invalid URI field: '%s'", l.URI)
	}
	if l.Version != EmptyString && !reVersion.MatchString(l.Version) {
		return xerrors.Errorf("invalid protocol field: '%s'", l.Version)
	}
	if l.Status < 100 || l.Status >= 600 {
		return xerrors.Errorf("invalid status field: %d", l.Status)
	}
	if l.RespSize < 0 {
		return xerrors.Errorf("invalid response size field: %d", l.RespSize)
	}
	if l.ReqSize != EmptyNumber && l.ReqSize < 0 {
		return xerrors.Errorf("invalid request size field: %d", l.ReqSize)
	}
	if l.RespTime != EmptyNumber && l.RespTime < 0 {
		return xerrors.Errorf("invalid response time field: %f", l.RespTime)
	}
	if l.UpstreamRespTime != EmptyNumber && l.UpstreamRespTime < 0 {
		return xerrors.Errorf("invalid upstream response time field: %f", l.UpstreamRespTime)
	}
	return nil
}

func (l *LogLine) assign(field string, value string, timeMultiplier float64) error {
	switch field {
	case fieldClient:
		l.Client = value
	case fieldRequest:
		if err := l.assignRequest(value); err != nil {
			return err
		}
	case fieldMethod:
		l.Method = value
	case fieldURI:
		l.URI = value
	case fieldProtocol:
		if err := l.assignProtocol(value); err != nil {
			return err
		}
	case fieldVersion:
		l.Version = value
	case fieldStatus:
		if err := l.assignStatus(value); err != nil {
			return err
		}
	case fieldRespSize:
		if err := l.assignRespSize(value); err != nil {
			return err
		}
	case fieldReqSize:
		if err := l.assignReqSize(value); err != nil {
			return err
		}
	case fieldRespTime:
		if err := l.assignRespTime(value, timeMultiplier); err != nil {
			return err
		}
	case fieldUpstreamRespTime:
		if err := l.assignUpstreamRespTime(value, timeMultiplier); err != nil {
			return err
		}
	case fieldVhost:
		l.Vhost = value
	case fieldCustom:
		l.Custom = value
	}
	return nil
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
	l.Method = req[0:idx]
	req = req[idx+1:]

	idx = strings.IndexByte(req, ' ')
	if idx < 0 {
		return xerrors.Errorf("invalid request: %q", request)
	}
	l.URI = req[0:idx]
	req = req[idx+1:]

	return l.assignProtocol(req)
}

func (l *LogLine) assignProtocol(proto string) error {
	if len(proto) <= 5 || !strings.HasPrefix(proto, "HTTP/") {
		return xerrors.Errorf("invalid protocol: %q", proto)
	}
	l.Version = proto[5:]
	return nil
}

func (l *LogLine) assignStatus(status string) error {
	if status == "-" {
		return nil
	}
	var err error
	l.Status, err = strconv.Atoi(status)
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
