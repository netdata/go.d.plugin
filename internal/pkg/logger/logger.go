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

func New(n namer) *logger {
	v := CacheGet(n)
	if v != nil {
		return v
	}
	v = &logger{
		log:   log.New(&colored{}, "", log.Ldate|log.Ltime),
		namer: n,
	}
	add(v)
	return v
}

func NewTest() *logger {
	return &logger{
		log:   log.New(&colored{}, "", log.Ldate|log.Ltime),
		namer: n{},
	}
}

type logger struct {
	log *log.Logger
	namer
}

func (l *logger) Critical(a ...interface{}) {
	l.print(CRITICAL, a...)
	os.Exit(1)
}

func (l *logger) Error(a ...interface{}) {
	l.print(ERROR, a...)
}

func (l *logger) Warning(a ...interface{}) {
	l.print(WARNING, a...)
}

func (l *logger) Info(a ...interface{}) {
	l.print(INFO, a...)
}

func (l *logger) Debug(a ...interface{}) {
	l.print(DEBUG, a...)
}

func (l *logger) Criticalf(format string, a ...interface{}) {
	l.Critical(fmt.Sprintf(format, a...))
}

func (l *logger) Errorf(format string, a ...interface{}) {
	l.Error(fmt.Sprintf(format, a...))
}

func (l *logger) Warningf(format string, a ...interface{}) {
	l.Warning(fmt.Sprintf(format, a...))
}

func (l *logger) Infof(format string, a ...interface{}) {
	l.Info(fmt.Sprintf(format, a...))
}

func (l *logger) Debugf(format string, a ...interface{}) {
	l.Debug(fmt.Sprintf(format, a...))
}

func (l *logger) print(level Severity, a ...interface{}) {
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

func (l *logger) SetLevel(lev Severity) {
	sevLevel = lev
}

func (l *logger) SetNamer(n namer) {
	l.namer = n
}

func (l logger) Level() Severity {
	return sevLevel
}
