package logger

import (
	"io"
	"log"
	"os"
	"runtime"
	"sync"
	"time"

	"github.com/mattn/go-isatty"
)

type (
	formatter struct {
		colored bool
		prefix  string
		out     io.Writer // destination for output
		flag    int       // properties

		mu  sync.Mutex // ensures atomic writes; protects the following fields
		buf []byte     // for accumulating text to write
	}
)

func newFormatter(out io.Writer) *formatter {
	if isatty.IsTerminal(os.Stderr.Fd()) {
		return &formatter{
			out:     out,
			colored: true,
			flag:    log.Lshortfile,
			buf:     make([]byte, 0, 120),
		}
	}
	return &formatter{
		out:     out,
		colored: false,
		prefix:  "go.d ",
		flag:    log.Ldate | log.Ltime,
		buf:     make([]byte, 0, 120),
	}
}

func (l *formatter) SetOutput(out io.Writer) {
	l.out = out
}

func (l *formatter) Output(severity Severity, module, job string, callDepth int, s string) {
	now := time.Now() // get this early.
	var file string
	var line int
	if l.flag&(log.Lshortfile|log.Llongfile) != 0 {
		var ok bool
		_, file, line, ok = runtime.Caller(callDepth)
		if !ok {
			file = "???"
			line = 0
		}
	}

	l.mu.Lock()
	defer l.mu.Unlock()

	l.formatTimestamp(now)
	l.buf = append(l.buf, l.prefix...)
	l.formatSeverity(severity)
	l.formatModuleJob(module, job)
	l.formatFile(file, line)
	l.buf = append(l.buf, s...)
	if s == "" || s[len(s)-1] != '\n' {
		l.buf = append(l.buf, '\n')
	}
	_, _ = l.out.Write(l.buf)
	l.buf = l.buf[:0]
}

func (l *formatter) formatModuleJob(module string, job string) {
	l.buf = append(l.buf, module...)
	l.buf = append(l.buf, '[')
	l.buf = append(l.buf, job...)
	l.buf = append(l.buf, "] "...)
}

// formatTimestamp writes timestamp to buf
func (l *formatter) formatTimestamp(t time.Time) {
	if l.flag&(log.Ldate|log.Ltime|log.Lmicroseconds) != 0 {
		if l.flag&log.LUTC != 0 {
			t = t.UTC()
		}
		if l.flag&log.Ldate != 0 {
			year, month, day := t.Date()
			itoa(&l.buf, year, 4)
			l.buf = append(l.buf, '-')
			itoa(&l.buf, int(month), 2)
			l.buf = append(l.buf, '-')
			itoa(&l.buf, day, 2)
			l.buf = append(l.buf, ' ')
		}
		if l.flag&(log.Ltime|log.Lmicroseconds) != 0 {
			hour, min, sec := t.Clock()
			itoa(&l.buf, hour, 2)
			l.buf = append(l.buf, ':')
			itoa(&l.buf, min, 2)
			l.buf = append(l.buf, ':')
			itoa(&l.buf, sec, 2)
			if l.flag&log.Lmicroseconds != 0 {
				l.buf = append(l.buf, '.')
				itoa(&l.buf, t.Nanosecond()/1e3, 6)
			}
			l.buf = append(l.buf, ' ')
		}
		l.buf[len(l.buf)-1] = ':'
		l.buf = append(l.buf, ' ')
	}
}

func (l *formatter) formatSeverity(severity Severity) {
	if l.colored {
		switch severity {
		case DEBUG:
			l.buf = append(l.buf, "\x1b[0;36m[ "...)
		case INFO:
			l.buf = append(l.buf, "\x1b[0;32m[ "...)
		case WARNING:
			l.buf = append(l.buf, "\x1b[0;33m[ "...)
		case ERROR:
			l.buf = append(l.buf, "\x1b[0;31m[ "...)
		case CRITICAL:
			l.buf = append(l.buf, "\x1b[0;37;41m[ "...)
		}
		putString(&l.buf, severity.ShortString(), 5)
		l.buf = append(l.buf, " ]\x1b[0m "...)
	} else {
		l.buf = append(l.buf, severity.String()...)
		l.buf = append(l.buf, ": "...)
	}
}

// formatFile writes file info to buf
func (l *formatter) formatFile(file string, line int) {
	if l.flag&(log.Lshortfile|log.Llongfile) == 0 {
		return
	}
	if l.flag&log.Lshortfile != 0 {
		short := file
		for i := len(file) - 1; i > 0; i-- {
			if file[i] == '/' {
				short = file[i+1:]
				break
			}
		}
		file = short
	}

	if l.colored {
		l.buf = append(l.buf, "\x1b[1;30m"...)
	}
	l.buf = append(l.buf, file...)
	l.buf = append(l.buf, ':')
	itoa(&l.buf, line, -1)
	if l.colored {
		l.buf = append(l.buf, "\x1b[0m "...)
	} else {
		l.buf = append(l.buf, ' ')
	}
}

// Cheap integer to fixed-width decimal ASCII. Give a negative width to avoid zero-padding.
func itoa(buf *[]byte, i int, wid int) {
	// Assemble decimal in reverse order.
	var b [20]byte
	bp := len(b) - 1
	for i >= 10 || wid > 1 {
		wid--
		q := i / 10
		b[bp] = byte('0' + i - q*10)
		bp--
		i = q
	}
	// i < 10
	b[bp] = byte('0' + i)
	*buf = append(*buf, b[bp:]...)
}

func putString(buf *[]byte, s string, wid int) {
	*buf = append(*buf, s...)
	space := wid - len(s)
	if space > 0 {
		for i := 0; i < space; i++ {
			*buf = append(*buf, ' ')
		}
	}
}
