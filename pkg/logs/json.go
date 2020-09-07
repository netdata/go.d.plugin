package logs

import (
	"bufio"
	"fmt"
	"io"

	"github.com/valyala/fastjson"
)

type JSONConfig struct {
	Mapping map[string]string `yaml:"mapping"`
}

type JSONParser struct {
	reader  *bufio.Reader
	parser  fastjson.Parser
	buf     []byte
	mapping map[string]string
}

func NewJSONParser(config JSONConfig, in io.Reader) (*JSONParser, error) {
	parser := &JSONParser{
		reader:  bufio.NewReader(in),
		mapping: config.Mapping,
		buf:     make([]byte, 0, 100),
	}
	return parser, nil
}

func (p *JSONParser) ReadLine(line LogLine) error {
	row, err := p.reader.ReadSlice('\n')
	if err != nil && len(row) == 0 {
		return err
	}
	return p.Parse(row, line)
}

func (p *JSONParser) Parse(row []byte, line LogLine) error {
	val, err := p.parser.ParseBytes(row)
	if err != nil {
		return err
	}
	obj, err := val.Object()
	if err != nil {
		return err
	}

	obj.Visit(func(key []byte, v *fastjson.Value) {
		if err != nil {
			return
		}
		switch v.Type() {
		case fastjson.TypeString, fastjson.TypeNumber:
		default:
			return
		}

		name := string(key)
		if mapped, ok := p.mapping[name]; ok {
			name = mapped
		}

		p.buf = p.buf[:0]
		if p.buf = v.MarshalTo(p.buf); len(p.buf) == 0 {
			return
		}

		switch v.Type() {
		case fastjson.TypeString:
			// trim "
			err = line.Assign(name, string(p.buf[1:len(p.buf)-1]))
		default:
			err = line.Assign(name, string(p.buf))
		}
	})
	if err != nil {
		return &ParseError{msg: fmt.Sprintf("json parse: %v", err), err: err}
	}
	return nil
}

func (p *JSONParser) Info() string {
	return fmt.Sprintf("json: %q", p.mapping)
}
