package logs

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
)

type JSONConfig struct {
	Mapping map[string]string `yaml:"mapping"`
}

type JSONParser struct {
	reader *bufio.Reader
	mapping map[string]string
}

func NewJSONParser(config JSONConfig, in io.Reader) (*JSONParser, error) {
	fieldMap := config.Mapping
	if fieldMap == nil {
		fieldMap = make(map[string]string)
	}

	parser := &JSONParser{
		reader:		bufio.NewReader(in),
		mapping:	fieldMap,
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
	var parsedLine map[string]interface{}

	if err := json.Unmarshal(row, &parsedLine); err != nil {
		return err
	}

	/* Now map the fields */
	for logField,logValue := range parsedLine {
		/* Convert logValue to string */
		stringValue := fmt.Sprintf("%v", logValue)

		if err := line.Assign( p.mapField(logField), stringValue ); err != nil {
			return err
		}
	}

	return nil
}

func (p *JSONParser) mapField(field string) string {
	if newLogLineField,exist := p.mapping[field]; exist {
		return newLogLineField
	}
	return field
}

func (p *JSONParser) Info() string {
	return fmt.Sprintf("json: %q", p.mapping)
}