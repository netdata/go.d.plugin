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

func New(n namer) *logger {
	l := &logger{
		log:   log.New(&coloredWriter{}, "", log.Ldate|log.Ltime),
		namer: n}
	if cache.get(n) == nil {
		cache.add(l)
	}
	return l
}

type logger struct {
	log *log.Logger
	namer
}

func (m *logger) Critical(a ...interface{}) {
	m.print(CRITICAL, a...)
	os.Exit(1)
}

func (m *logger) Error(a ...interface{}) {
	m.print(ERROR, a...)
}

func (m *logger) Warning(a ...interface{}) {
	m.print(WARNING, a...)
}

func (m *logger) Info(a ...interface{}) {
	m.print(INFO, a...)
}

func (m *logger) Debug(a ...interface{}) {
	m.print(DEBUG, a...)
}

func (m *logger) Criticalf(format string, a ...interface{}) {
	m.Critical(fmt.Sprintf(format, a...))
}

func (m *logger) Errorf(format string, a ...interface{}) {
	m.Error(fmt.Sprintf(format, a...))
}

func (m *logger) Warningf(format string, a ...interface{}) {
	m.Warning(fmt.Sprintf(format, a...))
}

func (m *logger) Infof(format string, a ...interface{}) {
	m.Info(fmt.Sprintf(format, a...))
}

func (m *logger) Debugf(format string, a ...interface{}) {
	m.Debug(fmt.Sprintf(format, a...))
}

func (m *logger) print(level Severity, a ...interface{}) {
	if level > sevLevel {
		return
	}
	m.log.Printf(
		"go.d: %s: %s: %s: %s",
		level,
		m.ModuleName(),
		m.JobName(),
		fmt.Sprintln(a...))
}

func (m *logger) SetLevel(l Severity) {
	sevLevel = l
}

func (m *logger) Level() Severity {
	return sevLevel
}
