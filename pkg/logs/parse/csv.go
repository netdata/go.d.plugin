package parse

import (
	"bytes"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"strings"
)

type (
	CSVConfig struct {
		Delimiter        rune   `yaml:"delimiter"`
		TrimLeadingSpace bool   `yaml:"trim_leading_space"`
		Format           string `yaml:"format"`
	}

	CSVParser struct {
		config CSVConfig
		reader *csv.Reader
		format *csvFormat
	}

	csvFormat struct {
		Raw          string
		maxIndex     int
		fieldIndexes map[string]int
	}
)

func (c *CSVConfig) applyDefaults() {
	if c.Delimiter == 0 {
		c.Delimiter = ' '
	}
}

func newCSVReader(in io.Reader, config CSVConfig) *csv.Reader {
	r := csv.NewReader(in)
	r.Comma = config.Delimiter
	r.TrimLeadingSpace = config.TrimLeadingSpace
	r.ReuseRecord = true
	r.FieldsPerRecord = -1
	return r
}

func NewCSVParser(config CSVConfig, in io.Reader) (*CSVParser, error) {
	if config.Format == "" {
		return nil, errors.New("empty csv format")
	}

	config.applyDefaults()

	format, err := newCSVFormat(config)
	if err != nil {
		return nil, fmt.Errorf("error on creating csv format : %v", err)
	}

	p := &CSVParser{
		config: config,
		reader: newCSVReader(in, config),
		format: format,
	}
	return p, nil
}

func (p *CSVParser) ReadLine(logLine LogLine) error {
	records, err := p.reader.Read()
	if err != nil {
		return handleCSVReadError(err)
	}
	return p.format.parse(records, logLine)
}

func (p *CSVParser) Parse(line []byte, logLine LogLine) error {
	r := newCSVReader(bytes.NewBuffer(line), p.config)
	records, err := r.Read()
	if err != nil {
		return handleCSVReadError(err)
	}
	return p.format.parse(records, logLine)
}

func (p CSVParser) Info() string {
	return p.format.Raw
}

func newCSVFormat(config CSVConfig) (*csvFormat, error) {
	r := csv.NewReader(strings.NewReader(config.Format))
	r.Comma = config.Delimiter
	r.TrimLeadingSpace = config.TrimLeadingSpace

	fields, err := r.Read()
	if err != nil {
		return nil, err
	}

	format := &csvFormat{
		Raw:          config.Format,
		fieldIndexes: make(map[string]int),
	}

	var max int
	for i, field := range fields {
		field = strings.Trim(field, `"`) // TODO: fix?
		if !strings.HasPrefix(field, "$") {
			continue
		}
		format.fieldIndexes[field[1:]] = i
		if max < i {
			max = i
		}
	}

	format.maxIndex = max
	return format, nil
}

func (f *csvFormat) parse(records []string, logLine LogLine) error {
	if len(records) <= f.maxIndex {
		return &Error{
			msg: fmt.Sprintf("csv unmatched line, expect at least %d fields, got %d", f.maxIndex+1, len(records)),
		}
	}

	for field, idx := range f.fieldIndexes {
		err := logLine.Assign(field, records[idx])
		if err != nil {
			return &Error{
				msg: fmt.Sprintf("csv error on assigning : %v", err),
				err: err,
			}
		}
	}
	return nil
}

func handleCSVReadError(err error) error {
	if !isCSVParseError(err) {
		return err
	}
	return &Error{
		msg: fmt.Sprintf("csv error on parsing : %v", err),
		err: err,
	}
}

func isCSVParseError(err error) bool {
	return errors.Is(err, csv.ErrBareQuote) || errors.Is(err, csv.ErrFieldCount) || errors.Is(err, csv.ErrQuote)
}
