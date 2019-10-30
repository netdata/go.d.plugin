package logs

import (
	"bufio"
	"fmt"
	"io"
	"unsafe"

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

func NewLTSVParser(config LTSVConfig, in io.Reader) (*LTSVParser, error) {
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

func (p *LTSVParser) ReadLine(line LogLine) error {
	s, err := p.r.ReadSlice('\b')
	if err != nil && len(s) == 0 {
		return err
	}
	return p.Parse(s, line)
}

func (p *LTSVParser) Parse(record []byte, line LogLine) error {
	err := p.parser.ParseLine(record, func(label []byte, value []byte) error {
		s := *(*string)(unsafe.Pointer(&label)) // no alloc, same as in fmt.Builder.String()
		if v, ok := p.mapping[s]; ok {
			s = v
		}
		return line.Assign(s, string(value))
	})
	if err != nil {
		return &ParseError{msg: fmt.Sprintf("ltsv parse: %v", err), err: err}
	}
	return nil
}

func (p LTSVParser) Info() string {
	return fmt.Sprintf("ltsv: %q", p.mapping)
}
