package squidlog

import (
	"io"
	"runtime"

	"github.com/netdata/go.d.plugin/pkg/logs"
	"github.com/netdata/go.d.plugin/pkg/stm"
)

func (s SquidLog) logPanicStackIfAny() {
	err := recover()
	if err == nil {
		return
	}
	s.Errorf("[ERROR] %s\n", err)
	for depth := 0; ; depth++ {
		_, file, line, ok := runtime.Caller(depth)
		if !ok {
			break
		}
		s.Errorf("======> %d: %v:%d", depth, file, line)
	}
	panic(err)
}

func (s *SquidLog) collect() (map[string]int64, error) {
	defer s.logPanicStackIfAny()
	s.mx.reset()

	var mx map[string]int64

	n, err := s.collectLogLines()

	if n > 0 || err == nil {
		mx = stm.ToMap(s.mx)
	}
	return mx, err
}

func (s *SquidLog) collectLogLines() (int, error) {
	var n int
	for {
		s.line.reset()
		err := s.parser.ReadLine(s.line)
		if err != nil {
			if err == io.EOF {
				return n, nil
			}
			if !logs.IsParseError(err) {
				return n, err
			}
			n++
			s.collectUnmatched()
			continue
		}
		n++
		if s.line.empty() {
			s.collectUnmatched()
		} else {
			s.collectLogLine()
		}
	}
}

func (s *SquidLog) collectLogLine() {}

func (s *SquidLog) collectUnmatched() {

}
