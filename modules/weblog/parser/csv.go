package parser

import (
	"bytes"
	"encoding/csv"
	"io"
	"strings"

	"golang.org/x/xerrors"
)

type (
	csvParser struct {
		timeMultiplier float64
		delimiter      rune
		reader         *csv.Reader
		format         *csvFormat
	}

	csvFormat struct {
		Raw          string
		maxIndex     int
		fieldIndexes map[string]int
	}
)

func newCSVParser(config Config, in io.Reader) (*csvParser, error) {
	format := newCSVFormat(config.CSV.Format)
	if !format.hasField(fieldClient) || !format.hasField(fieldStatus) || !format.hasField(fieldRespSize) {
		return nil, xerrors.New("missing some mandatory fields")
	}
	if !format.hasField(fieldRequest) && (!format.hasField(fieldMethod) || !format.hasField(fieldURI)) {
		return nil, xerrors.New("missing some mandatory fields")
	}
	return &csvParser{
		timeMultiplier: config.TimeMultiplier,
		delimiter:      config.CSV.Delimiter,
		reader:         newCSVReader(in, config.CSV.Delimiter),
		format:         format,
	}, nil
}

func newCSVReader(in io.Reader, delimiter rune) *csv.Reader {
	r := csv.NewReader(in)
	r.Comma = delimiter
	r.ReuseRecord = true
	r.FieldsPerRecord = -1
	return r
}

func (p *csvParser) ReadLine() (LogLine, error) {
	records, err := p.reader.Read()
	if err != nil {
		return LogLine{}, err
	}
	return p.format.parse(records, p.timeMultiplier)
}

func (p *csvParser) Parse(line []byte) (LogLine, error) {
	r := newCSVReader(bytes.NewBuffer(line), p.delimiter)
	records, err := r.Read()
	if err != nil {
		return LogLine{}, err
	}
	return p.format.parse(records, p.timeMultiplier)
}

var (
	errUnmatchedLine = xerrors.New("unmatched line")
)

func newCSVFormat(logFormat string) *csvFormat {
	format := &csvFormat{
		Raw:          logFormat,
		fieldIndexes: make(map[string]int, 13),
	}
	fields := strings.Fields(logFormat)
	offset := 0
	for i, field := range fields {
		field = strings.Trim(field, `'"[]`)
		switch field {
		case "$remote_addr":
			format.fieldIndexes[fieldClient] = i + offset
		case "$request":
			format.fieldIndexes[fieldRequest] = i + offset
		case "$request_method":
			format.fieldIndexes[fieldMethod] = i + offset
		case "$request_uri":
			format.fieldIndexes[fieldURI] = i + offset
		case "server_protocol":
			format.fieldIndexes[fieldProtocol] = i + offset
		case "$status":
			format.fieldIndexes[fieldStatus] = i + offset
		case "$body_bytes_sent", "$bytes_sent":
			format.fieldIndexes[fieldRespSize] = i + offset
		case "$request_length":
			format.fieldIndexes[fieldReqSize] = i + offset
		case "$request_time":
			format.fieldIndexes[fieldRespTime] = i + offset
		case "$upstream_response_time":
			format.fieldIndexes[fieldUpstreamRespTime] = i + offset
		case "$server_name", "$http_host", "$host", "$hostname":
			format.fieldIndexes[fieldVhost] = i + offset
		case "<custom>":
			format.fieldIndexes[fieldCustom] = i + offset
		case "$time_local":
			offset++
		}
	}
	format.maxIndex = len(fields) + offset

	return format
}

func (f *csvFormat) hasField(field string) bool {
	_, ok := f.fieldIndexes[field]
	return ok
}

func (f *csvFormat) parse(record []string, timeMultiplier float64) (log LogLine, err error) {
	log = emptyLogLine

	if len(record) < f.maxIndex {
		return log, errUnmatchedLine
	}

	for field, idx := range f.fieldIndexes {
		log.assign(field, record[idx], timeMultiplier)
	}
	return log, nil
}
