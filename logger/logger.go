package logger

import (
	"fmt"
	"log"
	"sync/atomic"
)

const (
	msgPerSecondLimit = 60
)

var (
	base      = New("base", "base")
	initialID = int64(1)
)

func createUniqueID() int64 {
	return atomic.AddInt64(&initialID, 1)
}

// New creates a new logger
func New(modName, jobName string) *Logger {
	return &Logger{
		log:     log.New(colored{}, "", log.Ldate|log.Ltime),
		modName: modName,
		jobName: jobName,
		id:      createUniqueID(),
	}
}

// NewLimited creates a new limited logger
func NewLimited(modName, jobName string) *Logger {
	logger := New(modName, jobName)
	logger.limited = true
	GlobalMsgCountWatcher.Register(logger)

	return logger
}

// Logger represents a logger object
type Logger struct {
	log *log.Logger

	id      int64
	modName string
	jobName string

	limited  bool
	msgCount int64
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

	if l == nil || l.log == nil {
		base.log.Printf(
			"go.d: %s: %s: %s: %s",
			severity,
			base.modName,
			base.jobName,
			fmt.Sprintln(a...),
		)
		return
	}

	if l.limited && globalSeverity < DEBUG && atomic.AddInt64(&l.msgCount, 1) > msgPerSecondLimit {
		return
	}

	l.log.Printf(
		"go.d: %s: %s: %s: %s",
		severity,
		l.modName,
		l.jobName,
		fmt.Sprintln(a...),
	)
}

// SetGlobalSeverity sets global severity level
func SetGlobalSeverity(severity Severity) {
	globalSeverity = severity
}
