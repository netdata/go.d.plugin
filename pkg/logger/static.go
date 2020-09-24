package logger

import "fmt"

// Panic logs a message with the Critical severity then panic
func Panic(a ...interface{}) {
	s := fmt.Sprint(a...)
	base.output(CRITICAL, 1, s)
	panic(s)
}

// Critical logs a message with the Critical severity
func Critical(a ...interface{}) {
	base.output(CRITICAL, 1, fmt.Sprint(a...))
}

// Error logs a message with the Error severity
func Error(a ...interface{}) {
	base.output(ERROR, 1, fmt.Sprint(a...))
}

// Warning logs a message with the Warning severity
func Warning(a ...interface{}) {
	base.output(WARNING, 1, fmt.Sprint(a...))
}

// Info logs a message with the Info severity
func Info(a ...interface{}) {
	base.output(INFO, 1, fmt.Sprint(a...))
}

// Debug logs a message with the Debug severity
func Debug(a ...interface{}) {
	base.output(DEBUG, 1, fmt.Sprint(a...))
}

// Panicln logs a message with the Critical severity then panic
func Panicln(a ...interface{}) {
	s := fmt.Sprintln(a...)
	base.output(CRITICAL, 1, s)
	panic(s)
}

// Criticalln logs a message with the Critical severity
func Criticalln(a ...interface{}) {
	base.output(CRITICAL, 1, fmt.Sprintln(a...))
}

// Errorln logs a message with the Error severity
func Errorln(a ...interface{}) {
	base.output(ERROR, 1, fmt.Sprintln(a...))
}

// Warningln logs a message with the Warning severity
func Warningln(a ...interface{}) {
	base.output(WARNING, 1, fmt.Sprintln(a...))
}

// Infoln logs a message with the Info severity
func Infoln(a ...interface{}) {
	base.output(INFO, 1, fmt.Sprintln(a...))
}

// Debugln logs a message with the Debug severity
func Debugln(a ...interface{}) {
	base.output(DEBUG, 1, fmt.Sprintln(a...))
}

// Panicf logs a message with the Critical severity using the same syntax and options as fmt.Printf then panic
func Panicf(format string, a ...interface{}) {
	s := fmt.Sprintf(format, a...)
	base.output(CRITICAL, 1, s)
	panic(s)
}

// Criticalf logs a message with the Critical severity using the same syntax and options as fmt.Printf
func Criticalf(format string, a ...interface{}) {
	base.output(CRITICAL, 1, fmt.Sprintf(format, a...))
}

// Errorf logs a message with the Error severity using the same syntax and options as fmt.Printf
func Errorf(format string, a ...interface{}) {
	base.output(ERROR, 1, fmt.Sprintf(format, a...))
}

// Warningf logs a message with the Warning severity using the same syntax and options as fmt.Printf
func Warningf(format string, a ...interface{}) {
	base.output(WARNING, 1, fmt.Sprintf(format, a...))
}

// Infof logs a message with the Info severity using the same syntax and options as fmt.Printf
func Infof(format string, a ...interface{}) {
	base.output(INFO, 1, fmt.Sprintf(format, a...))
}

// Debugf logs a message with the Debug severity using the same syntax and options as fmt.Printf
func Debugf(format string, a ...interface{}) {
	base.output(DEBUG, 1, fmt.Sprintf(format, a...))
}
