package logger

import (
	"fmt"
	"log"
	"os"
)

var dummy = &Logger{
	log: log.New(colored{}, "", log.Ldate|log.Ltime),
}

func New(modName, jobName string) *Logger {
	return &Logger{
		log:     log.New(colored{}, "", log.Ldate|log.Ltime),
		modName: modName,
		jobName: jobName,
	}
}

type Logger struct {
	log     *log.Logger
	modName string
	jobName string
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

	if l != nil && l.log != nil {
		l.log.Printf("go.d: %s: %s: %s: %s", level, l.modName, l.jobName, fmt.Sprintln(a...))
	} else {
		dummy.log.Printf("go.d: %s: dummy: dummy: %s", level, fmt.Sprintln(a...))
	}
}

func SetLevel(lev Severity) {
	sevLevel = lev
}

func SetModName(l *Logger, modName string) {
	l.modName = modName
}
