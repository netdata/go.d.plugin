package weblog

import (
	"encoding/csv"
	"io"
)

type (
	LogParser interface {
		ReadLine() (LogLine, error)
	}

	LogFormatConfig struct {
		LogFormatType string `yaml:"log_format_type"`
	}
)

type ParserConfig struct {
	Type      string            `yaml:"type"` // 'csv' or 'ltsv'
	Delimiter rune              `yaml:"delimiter"`
	Mapping   map[string]string `yaml:"mapping"`
}

var DefaultParserConfig = ParserConfig{
	Type:      "csv",
	Delimiter: 0,
	Mapping:   map[string]string{},
}

func NewCSVParser(config ParserConfig, r io.Reader) *csv.Reader {
	if config.Delimiter == 0 {
		config.Delimiter = ' '
	}
	parser := csv.NewReader(r)
	parser.Comma = config.Delimiter
	parser.ReuseRecord = true
	parser.TrimLeadingSpace = true
	parser.FieldsPerRecord = -1
	return parser
}

type LtsvReader struct {
	r io.Reader
}

func NewLtsvParser(config ParserConfig, r io.Reader) {

}
