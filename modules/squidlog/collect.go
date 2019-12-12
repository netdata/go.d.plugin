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

func (s *SquidLog) collectLogLine() {
	s.mx.Requests.Inc()
	s.collectRespSize()
	s.collectClientAddress()
	s.collectCacheCode()
	s.collectHTPCode()
	s.collectRespSize()
	s.collectReqMethod()
	s.collectHierCode()
	s.collectServerAddress()
	s.collectMimeType()
}

func (s *SquidLog) collectUnmatched() {
	s.mx.Requests.Inc()
	s.mx.ReqUnmatched.Inc()
}

func (s *SquidLog) collectRespTime() {
	if !s.line.hasRespTime() {
		return
	}
}

func (s *SquidLog) collectClientAddress() {
	if !s.line.hasClientAddress() {
		return
	}
}

func (s *SquidLog) collectCacheCode() {
	if !s.line.hasCacheCode() {
		return
	}
}

func (s *SquidLog) collectHTPCode() {
	if !s.line.hasHTTPCode() {
		return
	}
}

func (s *SquidLog) collectRespSize() {
	if !s.line.hasRespSize() {
		return
	}
}

func (s *SquidLog) collectReqMethod() {
	if !s.line.hasReqMethod() {
		return
	}
}

func (s *SquidLog) collectHierCode() {
	if !s.line.hasHierCode() {
		return
	}
}

func (s *SquidLog) collectServerAddress() {
	if !s.line.hasServerAddress() {
		return
	}
}

func (s *SquidLog) collectMimeType() {
	if !s.line.hasMimeType() {
		return
	}
}
