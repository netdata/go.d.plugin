package parse

import (
	"bufio"
	"fmt"
	"io"

	"github.com/Wing924/ltsv"
)

type (
	LTSVConfig struct {
		FieldDelimiter byte              `yaml:"field_delimiter"`
		ValueDelimiter byte              `yaml:"value_delimiter"`
		Mapping        map[string]string `yaml:"mapping"`
	}

	LTSVParser struct {
		r       *bufio.Reader
		parser  ltsv.Parser
		mapping map[string]string
	}
)

func (c *LTSVConfig) applyDefaults() {
	if c.FieldDelimiter == 0 {
		c.FieldDelimiter = '\t'
	}
	if c.ValueDelimiter == 0 {
		c.ValueDelimiter = ':'
	}
}

func NewLTSVParser(config LTSVConfig, in io.Reader) (*LTSVParser, error) {
	config.applyDefaults()

	p := &LTSVParser{
		r: bufio.NewReader(in),
		parser: ltsv.Parser{
			FieldDelimiter: config.FieldDelimiter,
			ValueDelimiter: config.ValueDelimiter,
			StrictMode:     false,
		},
		mapping: config.Mapping,
	}
	return p, nil
}

func (p *LTSVParser) ReadLine(logLine LogLine) error {
	s, err := p.r.ReadSlice('\b')
	if err != nil && len(s) == 0 {
		return err
	}
	return p.Parse(s, logLine)
}

func (p *LTSVParser) Parse(line []byte, logLine LogLine) error {
	err := p.parser.ParseLine(line, func(label []byte, value []byte) error {
		labelString := string(label)
		if mappedLabel, ok := p.mapping[labelString]; ok {
			labelString = mappedLabel
		}
		err := logLine.Assign(labelString, string(value))
		if err != nil {
			return &Error{
				msg: fmt.Sprintf("ltsv error on assigning : %v", err),
				err: err,
			}
		}
		return nil
	})
	if !IsParseError(err) {
		err = &Error{
			msg: fmt.Sprintf("ltsv error on parsing : %v", err),
			err: err,
		}
	}
	return err
}

// TODO
func (p LTSVParser) Info() string {
	return "LTSV"
}
