package weblog

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type LogLine struct {
	RemoteAddr       string
	Request          string
	Method           string
	URI              string
	Version          string
	Status           int
	BytesSent        int
	ReqLength        int
	Host             string
	RespTime         float64
	UpstreamRespTime []float64
	Custom           string
}

type Format struct {
	Raw       string
	TimeScale float64
	maxIndex  int

	RemoteAddr       int
	Request          int
	Status           int
	BytesSent        int
	Host             int
	RespTime         int
	UpstreamRespTime int
	ReqLength        int
	Custom           int
}

var (
	errUnmatchedLine       = errors.New("unmatched line")
	errInvalidRequestField = errors.New("invalid request field")
)

var (
	common        = NewFormat(time.Microsecond.Seconds(), `$remote_addr - $remote_user [$time_local] "$request" $status $body_bytes_sent`)
	combined      = NewFormat(time.Microsecond.Seconds(), `$remote_addr - $remote_user [$time_local] "$request" $status $body_bytes_sent "$http_referer" "$http_user_agent"`)
	custom1       = NewFormat(time.Microsecond.Seconds(), `$remote_addr - $remote_user [$time_local] "$request" $status $body_bytes_sent "$http_referer" "$http_user_agent" $uid_got $request_time'`)
	custom2       = NewFormat(time.Microsecond.Seconds(), `$remote_addr - $remote_user [$time_local] "$request" $status $body_bytes_sent $request_length $request_time'`)
	custom3       = NewFormat(time.Microsecond.Seconds(), `$remote_addr - $remote_user [$time_local] "$request" $status $body_bytes_sent $request_length $request_time "$upstream_response_time"'`)
	vhostCommon   = NewFormat(time.Microsecond.Seconds(), `$http_host $remote_addr - $remote_user [$time_local] "$request" $status $body_bytes_sent`)
	vhostCombined = NewFormat(time.Microsecond.Seconds(), `$http_host $remote_addr - $remote_user [$time_local] "$request" $status $body_bytes_sent "$http_referer" "$http_user_agent"`)
	vhostCustom1  = NewFormat(time.Microsecond.Seconds(), `$http_host $remote_addr - $remote_user [$time_local] "$request" $status $body_bytes_sent "$http_referer" "$http_user_agent" $uid_got $request_time`)
	vhostCustom2  = NewFormat(time.Microsecond.Seconds(), `$http_host $remote_addr - $remote_user [$time_local] "$request" $status $body_bytes_sent $request_length $request_time`)
	vhostCustom3  = NewFormat(time.Microsecond.Seconds(), `$http_host $remote_addr - $remote_user [$time_local] "$request" $status $body_bytes_sent $request_length $request_time "$upstream_response_time"`)
)

func NewFormat(timeScale float64, logFormat string) Format {
	format := Format{
		Raw:              logFormat,
		TimeScale:        timeScale,
		RemoteAddr:       -1,
		Request:          -1,
		Status:           -1,
		BytesSent:        -1,
		Host:             -1,
		RespTime:         -1,
		UpstreamRespTime: -1,
		ReqLength:        -1,
		Custom:           -1,
	}
	fields := strings.Fields(logFormat)
	offset := 0
	for i, field := range fields {
		field = strings.Trim(field, `'"[]`)
		switch field {
		case "$remote_addr":
			format.RemoteAddr = i + offset
		case "$request":
			format.Request = i + offset
		case "$status":
			format.Status = i + offset
		case "$body_bytes_sent", "$bytes_sent":
			format.BytesSent = i + offset
		case "$request_length":
			format.ReqLength = i + offset
		case "$request_time":
			format.RespTime = i + offset
		case "$upstream_response_time":
			format.UpstreamRespTime = i + offset
		case "$server_name", "$http_host", "$host", "$hostname":
			format.Host = i + offset
		case "<custom>":
			format.Custom = i + offset
		case "$time_local":
			offset++
		}
	}
	format.maxIndex = len(fields) + offset

	return format
}

func (f Format) Parse(record []string) (LogLine, error) {
	line := LogLine{
		Status:    -1,
		BytesSent: -1,
		RespTime:  -1,
		ReqLength: -1,
	}

	if len(record) < f.maxIndex {
		return line, errUnmatchedLine
	}

	if f.RemoteAddr >= 0 {
		line.RemoteAddr = record[f.RemoteAddr]
	}
	if f.Request >= 0 {
		line.Request = record[f.Request]
		var err error
		line.Method, line.URI, line.Version, err = parseRequest(line.Request)
		if err != nil {
			return line, err
		}
	}
	if f.Status >= 0 {
		val, err := strconv.Atoi(record[f.Status])
		if err != nil {
			return line, err
		}
		line.Status = val
	}
	if f.BytesSent >= 0 && record[f.BytesSent] != "-" {
		val, err := strconv.Atoi(record[f.BytesSent])
		if err != nil {
			return line, err
		}
		line.BytesSent = val
	}
	if f.Host >= 0 {
		line.Host = record[f.Host]
	}
	if f.RespTime >= 0 && record[f.RespTime] != "-" {
		val, err := strconv.ParseFloat(record[f.RespTime], 64)
		if err != nil {
			return line, err
		}
		line.RespTime = val * f.TimeScale
	}
	if f.UpstreamRespTime >= 0 && record[f.UpstreamRespTime] != "-" {
		times := strings.Split(record[f.UpstreamRespTime], ", ")
		line.UpstreamRespTime = make([]float64, len(times))
		for i, t := range times {
			val, err := strconv.ParseFloat(t, 64)
			if err != nil {
				return line, err
			}
			line.UpstreamRespTime[i] = val * f.TimeScale
		}
	}
	if f.ReqLength >= 0 && record[f.ReqLength] != "-" {
		val, err := strconv.Atoi(record[f.ReqLength])
		if err != nil {
			return line, err
		}
		line.ReqLength = val
	}
	if f.Custom >= 0 {
		line.Custom = record[f.Custom]
	}

	return line, nil
}

var (
	reRemoteAddr = regexp.MustCompile(`[\da-f.:]+|localhost`)
	reRequest    = regexp.MustCompile(`[A-Z]+ [^\s]+ HTTP/[0-9.]+`)
	reHost       = regexp.MustCompile(`[\w.-]+`)
)

func (f *Format) Match(record []string) error {
	line, err := f.Parse(record)
	if err != nil {
		return err
	}
	if f.RemoteAddr >= 0 && !reRemoteAddr.MatchString(line.RemoteAddr) {
		return fmt.Errorf("remoteAddr field bad syntax: '%s'", line.RemoteAddr)
	}
	if f.Request >= 0 && !reRequest.MatchString(line.Request) {
		return fmt.Errorf("request field bad syntax: '%s'", line.Request)
	}
	if f.Status >= 0 && (line.Status < 100 || line.Status >= 600) {
		return fmt.Errorf("status field bad syntax: %d", line.Status)
	}
	if f.BytesSent >= 0 && line.BytesSent < 0 {
		return fmt.Errorf("bytesSent field bad syntax: %d", line.BytesSent)
	}
	if f.Host >= 0 && !reHost.MatchString(line.Host) {
		return fmt.Errorf("host field bad syntax: '%s'", line.Host)
	}
	if f.RespTime >= 0 && line.RespTime < 0 {
		return fmt.Errorf("respTime field bad syntax: %f", line.RespTime)
	}
	if f.UpstreamRespTime >= 0 {
		for _, t := range line.UpstreamRespTime {
			if t < 0 {
				return fmt.Errorf("upstreamRespTime field bad syntax: %v", line.UpstreamRespTime)
			}
		}
	}
	if f.ReqLength >= 0 && line.ReqLength < 0 {
		return fmt.Errorf("reqLength field bad syntax: %d", line.ReqLength)
	}
	return nil
}

func parseRequest(request string) (method string, uri string, version string, err error) {
	fields := strings.Fields(request)
	if len(fields) != 3 {
		err = errInvalidRequestField
		return
	}
	return fields[0], fields[1], fields[2], nil
}
