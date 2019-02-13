package weblog

import (
	"encoding/csv"
	"io"
)

type (
	LogParser struct {
		parser *csv.Reader
	}
)

func NewLogParser() *LogParser {
	return &LogParser{}
}

func (p *LogParser) SetInput(r io.Reader) {
	p.parser = csv.NewReader(r)
	p.parser.Comma = ' '
	p.parser.ReuseRecord = true
	p.parser.TrimLeadingSpace = true
	p.parser.FieldsPerRecord = -1
}

func (p *LogParser) Read() ([]string, error) {
	return p.parser.Read()
}
