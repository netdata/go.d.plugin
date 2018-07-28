package logger

import (
	"fmt"
	"log"
	"os"
)

type namer interface {
	ModuleName() string
	JobName() string
}

type n struct{}

func (n) ModuleName() string {
	return "module"
}

func (n) JobName() string {
	return "job"
}

func New(n namer) *Logger {
	v := CacheGet(n)
	if v != nil {
		return v
	}
	v = &Logger{
		log:   log.New(colored{}, "", log.Ldate|log.Ltime),
		namer: n,
	}
	add(v)
	return v
}

func NewTest() *Logger {
	return &Logger{
		log:   log.New(&colored{}, "", log.Ldate|log.Ltime),
		namer: n{},
	}
}

type Logger struct {
	log *log.Logger
	namer
}

func (l *Logger) Critical(a ...interface{}) {
	l.print(CRITICAL, a...)
	os.Exit(1)
}

func (l *Logger) Error(a ...interface{}) {
	l.print(ERROR, a...)
}

func (l *Logger) Warning(a ...interface{}) {
	l.print(WARNING, a...)
}

func (l *Logger) Info(a ...interface{}) {
	l.print(INFO, a...)
}

func (l *Logger) Debug(a ...interface{}) {
	l.print(DEBUG, a...)
}

func (l *Logger) Criticalf(format string, a ...interface{}) {
	l.Critical(fmt.Sprintf(format, a...))
}

func (l *Logger) Errorf(format string, a ...interface{}) {
	l.Error(fmt.Sprintf(format, a...))
}

func (l *Logger) Warningf(format string, a ...interface{}) {
	l.Warning(fmt.Sprintf(format, a...))
}

func (l *Logger) Infof(format string, a ...interface{}) {
	l.Info(fmt.Sprintf(format, a...))
}

func (l *Logger) Debugf(format string, a ...interface{}) {
	l.Debug(fmt.Sprintf(format, a...))
}

func (l *Logger) print(level Severity, a ...interface{}) {
	if level > sevLevel {
		return
	}
	l.log.Printf(
		"go.d: %s: %s: %s: %s",
		level,
		l.ModuleName(),
		l.JobName(),
		fmt.Sprintln(a...))
}

func (l *Logger) SetLevel(lev Severity) {
	sevLevel = lev
}
