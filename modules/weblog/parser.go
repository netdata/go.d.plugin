package weblog

import (
	"encoding/csv"
	"errors"
	"io"
	"strconv"
	"strings"
)

type (
	logParser struct {
		parser  *csv.Reader
		pattern logPattern
	}
	logLine struct {
		RemoteAddr       string
		Request          string
		Method           string
		URI              string
		Version          string
		Status           int
		BytesSent        int
		Host             string
		RespTime         float64
		RespTimeUpstream string
		ReqLength        int
		UserDefined      string
	}
)

var (
	errUnmatchedLine       = errors.New("unmatched line")
	errInvalidRequestField = errors.New("invalid request field")
)

func newLogParser(pattern logPattern) *logParser {
	return &logParser{
		pattern: pattern,
	}
}

func (p *logParser) SetDataSource(r io.Reader) {
	p.parser = csv.NewReader(r)
	p.parser.Comma = ' '
	p.parser.ReuseRecord = true
	p.parser.TrimLeadingSpace = true
	p.parser.FieldsPerRecord = -1
}

func (p *logParser) Read() (logLine, error) {
	log := logLine{
		Status:    -1,
		BytesSent: -1,
		RespTime:  -1,
		ReqLength: -1,
	}
	records, err := p.parser.Read()
	if err != nil {
		return log, err
	}

	if len(records) <= p.pattern.MaxIndex() {
		return log, errUnmatchedLine
	}

	for field, idx := range p.pattern.Mapping {
		if idx < 0 {
			continue
		}
		switch fieldID(field) {
		case fieldRemoteAddr:
			log.RemoteAddr = records[idx]
		case fieldRequest:
			log.Request = records[idx]
			var err error
			log.Method, log.URI, log.Version, err = parseRequest(log.Request)
			if err != nil {
				return log, err
			}
		case fieldStatus:
			val, err := strconv.Atoi(records[idx])
			if err != nil {
				return log, err
			}
			log.Status = val
		case fieldBytesSent:
			val, err := strconv.Atoi(records[idx])
			if err != nil {
				return log, err
			}
			log.BytesSent = val
		case fieldHost:
			log.Host = records[idx]
		case fieldRespTime:
			val, err := strconv.ParseFloat(records[idx], 64)
			if err != nil {
				return log, err
			}
			log.RespTime = val
		case fieldRespTimeUpstream:
			log.RespTimeUpstream = records[idx]
		case fieldReqLength:
			val, err := strconv.Atoi(records[idx])
			if err != nil {
				return log, err
			}
			log.ReqLength = val
		case fieldUserDefined:
			log.UserDefined = records[idx]
		}
	}

	return log, nil
}

func parseRequest(request string) (method string, uri string, version string, err error) {
	fields := strings.Fields(request)
	if len(fields) != 3 {
		err = errInvalidRequestField
		return
	}
	return fields[0], fields[1], fields[2], nil
}
