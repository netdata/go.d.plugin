// SPDX-License-Identifier: GPL-3.0-or-later

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
	isCLI     = func() bool {
		switch os.Getenv("NETDATA_FORCE_COLOR") {
		case "1", "true":
			return true
		case "0", "false":
			return true
		default:
			return isatty.IsTerminal(os.Stderr.Fd())
		}
	}()
	Prefix = "goplugin"
)

// Logger represents a logger object
type Logger struct {
	formatter *formatter

	id      int64
	modName string
	jobName string

	limited  bool
	msgCount int64
}

// New creates a new logger.
func New(modName, jobName string) *Logger {
	return &Logger{
		formatter: newFormatter(os.Stderr, isCLI, Prefix),
		modName:   modName,
		jobName:   jobName,
		id:        uniqueID(),
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
	l.output(CRITICAL, 1, s)
	panic(s)
}

// Critical logs a message with the Critical severity
func (l *Logger) Critical(a ...interface{}) {
	l.output(CRITICAL, 1, fmt.Sprint(a...))
}

// Error logs a message with the Error severity
func (l *Logger) Error(a ...interface{}) {
	l.output(ERROR, 1, fmt.Sprint(a...))
}

// Warning logs a message with the Warning severity
func (l *Logger) Warning(a ...interface{}) {
	l.output(WARNING, 1, fmt.Sprint(a...))
}

// Info logs a message with the Info severity
func (l *Logger) Info(a ...interface{}) {
	l.output(INFO, 1, fmt.Sprint(a...))
}

// Print logs a message with the Info severity (same as Info)
func (l *Logger) Print(a ...interface{}) {
	l.output(INFO, 1, fmt.Sprint(a...))
}

// Debug logs a message with the Debug severity
func (l *Logger) Debug(a ...interface{}) {
	l.output(DEBUG, 1, fmt.Sprint(a...))
}

// Panicln logs a message with the Critical severity then panic
func (l *Logger) Panicln(a ...interface{}) {
	s := fmt.Sprintln(a...)
	l.output(CRITICAL, 1, s)
	panic(s)
}

// Criticalln logs a message with the Critical severity
func (l *Logger) Criticalln(a ...interface{}) {
	l.output(CRITICAL, 1, fmt.Sprintln(a...))
}

// Errorln logs a message with the Error severity
func (l *Logger) Errorln(a ...interface{}) {
	l.output(ERROR, 1, fmt.Sprintln(a...))
}

// Warningln logs a message with the Warning severity
func (l *Logger) Warningln(a ...interface{}) {
	l.output(WARNING, 1, fmt.Sprintln(a...))
}

// Infoln logs a message with the Info severity
func (l *Logger) Infoln(a ...interface{}) {
	l.output(INFO, 1, fmt.Sprintln(a...))
}

// Println logs a message with the Info severity (same as Infoln)
func (l *Logger) Println(a ...interface{}) {
	l.output(INFO, 1, fmt.Sprintln(a...))
}

// Debugln logs a message with the Debug severity
func (l *Logger) Debugln(a ...interface{}) {
	l.output(DEBUG, 1, fmt.Sprintln(a...))
}

// Panicf logs a message with the Critical severity using the same syntax and options as fmt.Printf then panic
func (l *Logger) Panicf(format string, a ...interface{}) {
	s := fmt.Sprintf(format, a...)
	l.output(CRITICAL, 1, s)
	panic(s)
}

// Criticalf logs a message with the Critical severity using the same syntax and options as fmt.Printf
func (l *Logger) Criticalf(format string, a ...interface{}) {
	l.output(CRITICAL, 1, fmt.Sprintf(format, a...))
}

// Errorf logs a message with the Error severity using the same syntax and options as fmt.Printf
func (l *Logger) Errorf(format string, a ...interface{}) {
	l.output(ERROR, 1, fmt.Sprintf(format, a...))
}

// Warningf logs a message with the Warning severity using the same syntax and options as fmt.Printf
func (l *Logger) Warningf(format string, a ...interface{}) {
	l.output(WARNING, 1, fmt.Sprintf(format, a...))
}

// Infof logs a message with the Info severity using the same syntax and options as fmt.Printf
func (l *Logger) Infof(format string, a ...interface{}) {
	l.output(INFO, 1, fmt.Sprintf(format, a...))
}

// Printf logs a message with the Info severity using the same syntax and options as fmt.Printf
func (l *Logger) Printf(format string, a ...interface{}) {
	l.output(INFO, 1, fmt.Sprintf(format, a...))
}

// Debugf logs a message with the Debug severity using the same syntax and options as fmt.Printf
func (l *Logger) Debugf(format string, a ...interface{}) {
	l.output(DEBUG, 1, fmt.Sprintf(format, a...))
}

func (l *Logger) output(severity Severity, callDepth int, msg string) {
	if severity > globalSeverity {
		return
	}

	if l == nil || l.formatter == nil {
		base.formatter.Output(severity, base.modName, base.jobName, callDepth+2, msg)
		return
	}

	if l.limited && globalSeverity < DEBUG && atomic.AddInt64(&l.msgCount, 1) > msgPerSecondLimit {
		return
	}
	l.formatter.Output(severity, l.modName, l.jobName, callDepth+2, msg)
}

func uniqueID() int64 {
	return atomic.AddInt64(&initialID, 1)
}
