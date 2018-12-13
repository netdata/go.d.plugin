package logger

import (
	"fmt"
	"os"
	"sync/atomic"

	"github.com/mattn/go-isatty"
)

const (
	msgPerSecondLimit = 60
)

var (
	base      = New("base", "base")
	initialID = int64(1)
)

var defaultFormatter = newFormatter(os.Stderr, isatty.IsTerminal(os.Stderr.Fd()))

// Logger represents a logger object
type Logger struct {
	formatter *formatter

	id      int64
	modName string
	jobName string

	limited  bool
	msgCount int64
}

// New creates a new logger
func New(modName, jobName string) *Logger {
	return &Logger{
		formatter: defaultFormatter,
		modName:   modName,
		jobName:   jobName,
		id:        createUniqueID(),
	}
}

// NewLimited creates a new limited logger
func NewLimited(modName, jobName string) *Logger {
	logger := New(modName, jobName)
	logger.limited = true
	GlobalMsgCountWatcher.Register(logger)

	return logger
}

// Panic logs a message with the Critical severity then panic
func (l *Logger) Panic(a ...interface{}) {
	s := fmt.Sprint(a...)
	l.print(CRITICAL, s)
	panic(s)
}

// Critical logs a message with the Critical severity
func (l *Logger) Critical(a ...interface{}) {
	l.print(CRITICAL, a...)
}

// Error logs a message with the Error severity
func (l *Logger) Error(a ...interface{}) {
	l.print(ERROR, a...)
}

// Warning logs a message with the Warning severity
func (l *Logger) Warning(a ...interface{}) {
	l.print(WARNING, a...)
}

// Info logs a message with the Info severity
func (l *Logger) Info(a ...interface{}) {
	l.print(INFO, a...)
}

// Debug logs a message with the Debug severity
func (l *Logger) Debug(a ...interface{}) {
	l.print(DEBUG, a...)
}

// Panicf logs a message with the Critical severity using the same syntax and options as fmt.Printf then panic
func (l *Logger) Panicf(format string, a ...interface{}) {
	s := fmt.Sprintf(format, a...)
	l.print(CRITICAL, s)
	panic(s)
}

// Criticalf logs a message with the Critical severity using the same syntax and options as fmt.Printf
func (l *Logger) Criticalf(format string, a ...interface{}) {
	l.Critical(fmt.Sprintf(format, a...))
}

// Errorf logs a message with the Error severity using the same syntax and options as fmt.Printf
func (l *Logger) Errorf(format string, a ...interface{}) {
	l.Error(fmt.Sprintf(format, a...))
}

// Warningf logs a message with the Warning severity using the same syntax and options as fmt.Printf
func (l *Logger) Warningf(format string, a ...interface{}) {
	l.Warning(fmt.Sprintf(format, a...))
}

// Infof logs a message with the Info severity using the same syntax and options as fmt.Printf
func (l *Logger) Infof(format string, a ...interface{}) {
	l.Info(fmt.Sprintf(format, a...))
}

// Debugf logs a message with the Debug severity using the same syntax and options as fmt.Printf
func (l *Logger) Debugf(format string, a ...interface{}) {
	l.Debug(fmt.Sprintf(format, a...))
}

func (l *Logger) print(severity Severity, a ...interface{}) {
	if severity > globalSeverity {
		return
	}

	if l == nil || l.formatter == nil {
		base.formatter.Output(severity, base.modName, base.jobName, 3, fmt.Sprint(a...))
		return
	}

	if l.limited && globalSeverity < DEBUG && atomic.AddInt64(&l.msgCount, 1) > msgPerSecondLimit {
		return
	}
	l.formatter.Output(severity, l.modName, l.jobName, 3, fmt.Sprint(a...))
}

// SetSeverity sets global severity level
func SetSeverity(severity Severity) {
	globalSeverity = severity
}

func createUniqueID() int64 {
	return atomic.AddInt64(&initialID, 1)
}
