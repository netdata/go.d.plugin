package parser

//
//import (
//	"bufio"
//	"io"
//	"regexp"
//
//	"golang.org/x/xerrors"
//)
//
//type (
//	regexpParser struct {
//		timeMultiplier float64
//		scanner        *bufio.Scanner
//		pattern        *regexp.Regexp
//	}
//)
//
//func newRegExpParser(config Config, in io.Reader) (*regexpParser, error) {
//	if config.RegExp.Pattern == "" {
//		return nil, xerrors.New("empty RegExp pattern")
//	}
//	pattern, err := regexp.Compile(config.RegExp.Pattern)
//	if err != nil {
//		return nil, err
//	}
//	return &regexpParser{
//		timeMultiplier: config.TimeMultiplier,
//		scanner:        bufio.NewScanner(in),
//		pattern:        pattern,
//	}, nil
//}
//
//func (p *regexpParser) ReadLine() (LogLine, error) {
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
//func (p *regexpParser) Parse(line []byte) (LogLine, error) {
//	log := emptyLogLine
//	match := p.pattern.FindSubmatch(line)
//	for i, name := range p.pattern.SubexpNames() {
//		if name == "" || match[i] == nil {
//			continue
//		}
//		if err := log.assign(name, string(match[i]), p.timeMultiplier); err != nil {
//			return log, err
//		}
//
//	}
//	return log, nil
//}
