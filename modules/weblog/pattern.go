package weblog

//
//import (
//	"errors"
//	"fmt"
//	"regexp"
//	"strconv"
//	"strings"
//)
//
//type (
//	LogPattern struct {
//		Name     string
//		maxIndex int
//
//		RemoteAddr       int
//		Request          int
//		Status           int
//		BytesSent        int
//		Host             int
//		ReqTime         int
//		RespTimeUpstream int
//		ReqLength        int
//		UserDefined      int
//	}
//
//	LogLine struct {
//		RemoteAddr       string
//		Request          string
//		Method           string
//		URI              string
//		Version          string
//		Status           int
//		BytesSent        int
//		Host             string
//		ReqTime         float64
//		RespTimeUpstream []float64
//		ReqLength        int
//		UserDefined      string
//	}
//)
//
//const (
//	keyRemoteAddr       = "remote_addr"
//	keyRequest          = "request"
//	keyStatus           = "status"
//	keyBytesSent        = "bytes_sent"
//	keyHost             = "host"
//	keyRespTime         = "respTime"
//	keyRespTimeUpstream = "resp_time_upstream"
//	keyReqLength        = "request_length"
//	keyUserDefined      = "user_defined"
//)
//
//var (
//	reRemoteAddr = regexp.MustCompile(`[\da-f.:]+|localhost`)
//	reRequest    = regexp.MustCompile(`[A-Z]+ [^\s]+ HTTP/[0-9.]+`)
//	reHost       = regexp.MustCompile(`[\da-z.:-]+`) // TODO: not sure about this
//)
//
//var (
//	errUnmatchedLine       = errors.New("unmatched line")
//	errInvalidRequestField = errors.New("invalid request field")
//)
//
//func NewLogPattern(name string, mapping map[string]int) (*LogPattern, error) {
//	pattern := LogPattern{
//		Name:     name,
//		maxIndex: -1,
//
//		RemoteAddr:       -1,
//		Request:          -1,
//		Status:           -1,
//		BytesSent:        -1,
//		Host:             -1,
//		ReqTime:         -1,
//		RespTimeUpstream: -1,
//		ReqLength:        -1,
//		UserDefined:      -1,
//	}
//
//	set := map[int]bool{}
//
//	for key, idx := range mapping {
//		if key != keyUserDefined {
//			if set[idx] {
//				return nil, fmt.Errorf("duplicate index in log pattern: %d", idx)
//			}
//			set[idx] = true
//		}
//		switch key {
//		case keyRemoteAddr:
//			pattern.RemoteAddr = idx
//		case keyRequest:
//			pattern.Request = idx
//		case keyStatus:
//			pattern.Status = idx
//		case keyBytesSent:
//			pattern.BytesSent = idx
//		case keyHost:
//			pattern.Host = idx
//		case keyRespTime:
//			pattern.ReqTime = idx
//		case keyRespTimeUpstream:
//			pattern.RespTimeUpstream = idx
//		case keyReqLength:
//			pattern.ReqLength = idx
//		case keyUserDefined:
//			pattern.UserDefined = idx
//		default:
//			return nil, fmt.Errorf("unknown field in log pattern: %s", key)
//		}
//		if pattern.maxIndex < idx {
//			pattern.maxIndex = idx
//		}
//	}
//	return &pattern, nil
//}
//
//func (p LogPattern) MaxIndex() int {
//	return p.maxIndex
//}
//
//func (p *LogPattern) Match(records []string) error {
//	line, err := p.Parse(records)
//	if err != nil {
//		return err
//	}
//	if p.RemoteAddr >= 0 && !reRemoteAddr.MatchString(line.RemoteAddr) {
//		return fmt.Errorf("'%s' field bad syntax: '%s'", keyRemoteAddr, line.RemoteAddr)
//	}
//	if p.Request >= 0 && !reRequest.MatchString(line.Request) {
//		return fmt.Errorf("'%s' field bad syntax: '%s'", keyRequest, line.Request)
//	}
//	if p.Status >= 0 && (line.Status < 100 || line.Status >= 600) {
//		return fmt.Errorf("'%s' field bad syntax: %d", keyStatus, line.Status)
//	}
//	if p.BytesSent >= 0 && line.BytesSent < 0 {
//		return fmt.Errorf("'%s' field bad syntax: %d", keyBytesSent, line.BytesSent)
//	}
//	if p.Host >= 0 && !reHost.MatchString(line.Host) {
//		return fmt.Errorf("'%s' field bad syntax: '%s'", keyHost, line.Host)
//	}
//	if p.ReqTime >= 0 && line.ReqTime < 0 {
//		return fmt.Errorf("'%s' field bad syntax: %f", keyRespTime, line.ReqTime)
//	}
//	if p.RespTimeUpstream >= 0 {
//		for _, time := range line.RespTimeUpstream {
//			if time < 0 {
//				return fmt.Errorf("'%s' field bad syntax: %v", keyRespTimeUpstream, line.RespTimeUpstream)
//			}
//		}
//	}
//	if p.ReqLength >= 0 && line.ReqLength < 0 {
//		return fmt.Errorf("'%s' field bad syntax: %d", keyReqLength, line.ReqLength)
//	}
//	return nil
//}
//
//func (p *LogPattern) Parse(records []string) (LogLine, error) {
//	line := LogLine{
//		Status:    -1,
//		BytesSent: -1,
//		ReqTime:  -1,
//		ReqLength: -1,
//	}
//
//	if len(records) <= p.MaxIndex() {
//		return line, errUnmatchedLine
//	}
//
//	if p.RemoteAddr >= 0 {
//		line.RemoteAddr = records[p.RemoteAddr]
//	}
//	if p.Request >= 0 {
//		line.Request = records[p.Request]
//		var err error
//		line.Method, line.URI, line.Version, err = parseRequest(line.Request)
//		if err != nil {
//			return line, err
//		}
//	}
//	if p.Status >= 0 {
//		val, err := strconv.Atoi(records[p.Status])
//		if err != nil {
//			return line, err
//		}
//		line.Status = val
//	}
//	if p.BytesSent >= 0 {
//		val, err := strconv.Atoi(records[p.BytesSent])
//		if err != nil {
//			return line, err
//		}
//		line.BytesSent = val
//	}
//	if p.Host >= 0 {
//		line.Host = records[p.Host]
//	}
//	if p.ReqTime >= 0 {
//		val, err := strconv.ParseFloat(records[p.ReqTime], 64)
//		if err != nil {
//			return line, err
//		}
//		line.ReqTime = val
//	}
//	if p.RespTimeUpstream >= 0 && records[p.RespTimeUpstream] != "-" {
//		times := strings.Split(records[p.RespTimeUpstream], ", ")
//		line.RespTimeUpstream = make([]float64, len(times))
//		for i, time := range times {
//			val, err := strconv.ParseFloat(time, 64)
//			if err != nil {
//				return line, err
//			}
//			line.RespTimeUpstream[i] = val
//		}
//	}
//	if p.ReqLength >= 0 {
//		val, err := strconv.Atoi(records[p.ReqLength])
//		if err != nil {
//			return line, err
//		}
//		line.ReqLength = val
//	}
//	if p.UserDefined >= 0 {
//		line.UserDefined = records[p.UserDefined]
//	}
//
//	return line, nil
//}
//
//func parseRequest(request string) (method string, uri string, version string, err error) {
//	fields := strings.Fields(request)
//	if len(fields) != 3 {
//		err = errInvalidRequestField
//		return
//	}
//	return fields[0], fields[1], fields[2], nil
//}
//
//func guessPattern(records []string) *LogPattern {
//	for _, pattern := range defaultLogFmtPatterns {
//		err := pattern.Match(records)
//		if err == nil {
//			return pattern
//		}
//	}
//	return nil
//}
//
//var (
//	// <combined> cookie request_time
//	logFmtCustom1, _ = NewLogPattern(
//		`$remote_addr - $remote_user [$time_local] "$request" $status $body_bytes_sent "$http_referer" "$http_user_agent" $uid_got $request_time`,
//		map[string]int{
//			keyRemoteAddr: 0,
//			keyRequest:    5,
//			keyStatus:     6,
//			keyBytesSent:  7,
//			keyRespTime:   11,
//		},
//	)
//	// host <combined> cookie request_time
//	logFmtHostCustom1, _ = NewLogPattern(
//		`$http_host $remote_addr - $remote_user [$time_local] "$request" $status $body_bytes_sent "$http_referer" "$http_user_agent" $uid_got $request_time`,
//		map[string]int{
//			keyHost:       0,
//			keyRemoteAddr: 1,
//			keyRequest:    6,
//			keyStatus:     7,
//			keyBytesSent:  8,
//			keyRespTime:   12,
//		},
//	)
//
//	// <common> request_length request_time
//	logFmtCustom2, _ = NewLogPattern(
//		`$remote_addr - $remote_user [$time_local] "$request" $status $body_bytes_sent $request_length $request_time`,
//		map[string]int{
//			keyRemoteAddr: 0,
//			keyRequest:    5,
//			keyStatus:     6,
//			keyBytesSent:  7,
//			keyReqLength:  8,
//			keyRespTime:   9,
//		},
//	)
//	// host <common> request_length request_time
//	logFmtHostCustom2, _ = NewLogPattern(
//		`$http_host $remote_addr - $remote_user [$time_local] "$request" $status $body_bytes_sent $request_length $request_time`,
//		map[string]int{
//			keyHost:       0,
//			keyRemoteAddr: 1,
//			keyRequest:    6,
//			keyStatus:     7,
//			keyBytesSent:  8,
//			keyReqLength:  9,
//			keyRespTime:   10,
//		},
//	)
//
//	// <common> request_length request_time upstream_response_time
//	logFmtCustom3, _ = NewLogPattern(
//		`$remote_addr - $remote_user [$time_local] "$request" $status $body_bytes_sent $request_length $request_time "$upstream_response_time"`,
//		map[string]int{
//			keyRemoteAddr:       0,
//			keyRequest:          5,
//			keyStatus:           6,
//			keyBytesSent:        7,
//			keyReqLength:        8,
//			keyRespTime:         9,
//			keyRespTimeUpstream: 10,
//		},
//	)
//
//	// host <common> request_length request_time upstream_response_time
//	logFmtHostCustom3, _ = NewLogPattern(
//		`$http_host $remote_addr - $remote_user [$time_local] "$request" $status $body_bytes_sent $request_length $request_time "$upstream_response_time"`,
//		map[string]int{
//			keyHost:             0,
//			keyRemoteAddr:       1,
//			keyRequest:          6,
//			keyStatus:           7,
//			keyBytesSent:        8,
//			keyReqLength:        9,
//			keyRespTime:         10,
//			keyRespTimeUpstream: 11,
//		},
//	)
//
//	logFmtCommon, _ = NewLogPattern(
//		`$remote_addr - $remote_user [$time_local] "$request" $status $body_bytes_sent`,
//		map[string]int{
//			keyRemoteAddr: 0,
//			keyRequest:    5,
//			keyStatus:     6,
//			keyBytesSent:  7,
//		},
//	)
//	logFmtHostCommon, _ = NewLogPattern(
//		`$http_host $remote_addr - $remote_user [$time_local] "$request" $status $body_bytes_sent`,
//		map[string]int{
//			keyHost:       0,
//			keyRemoteAddr: 1,
//			keyRequest:    6,
//			keyStatus:     7,
//			keyBytesSent:  8,
//		},
//	)
//
//	defaultLogFmtPatterns = []*LogPattern{
//		logFmtHostCustom1,
//		logFmtCustom1,
//		logFmtHostCustom2,
//		logFmtCustom2,
//		logFmtHostCustom3,
//		logFmtCustom3,
//		logFmtHostCommon,
//		logFmtCommon,
//	}
//)
