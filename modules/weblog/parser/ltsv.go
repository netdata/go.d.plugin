package parser

//
//import (
//	"bufio"
//	"io"
//
//	"github.com/Wing924/ltsv"
//)
//
//type (
//	ltsvParser struct {
//		timeMultiplier float64
//		scanner        *bufio.Scanner
//		parser         ltsv.Parser
//		mapping        map[string]string
//	}
//)
//
//func newLTSVParser(config Config, in io.Reader) (*ltsvParser, error) {
//	return &ltsvParser{
//		timeMultiplier: config.TimeMultiplier,
//		scanner:        bufio.NewScanner(in),
//		parser: ltsv.Parser{
//			FieldDelimiter: config.LTSV.FieldDelimiter,
//			ValueDelimiter: config.LTSV.ValueDelimiter,
//			StrictMode:     false,
//		},
//		mapping: config.LTSV.Mapping,
//	}, nil
//}
//
//func (p *ltsvParser) ReadLine() (LogLine, error) {
//	log := emptyLogLine
//	if p.scanner.Scan() {
//		return p.Parse(p.scanner.Bytes())
//	}
//	err := p.scanner.Err()
//	if err == nil {
//		err = io.EOF
//	}
//	return log, err
//}
//
//func (p *ltsvParser) Parse(line []byte) (LogLine, error) {
//	log := emptyLogLine
//	err := p.parser.ParseLine(line, func(label []byte, value []byte) error {
//		labelString := string(label)
//		if mappedLabel, ok := p.mapping[labelString]; ok {
//			labelString = mappedLabel
//		}
//		return log.assign(labelString, string(value), p.timeMultiplier)
//	})
//	return log, err
//}
