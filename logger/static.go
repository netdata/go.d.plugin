package logger

import "fmt"

// Panic logs a message with the Critical severity then panic
func Panic(a ...interface{}) {
	base.Panic(a...)
	s := fmt.Sprint(a...)
	base.print(CRITICAL, s)
	panic(s)
}

// Critical logs a message with the Critical severity
func Critical(a ...interface{}) {
	base.print(CRITICAL, a...)
}

// Error logs a message with the Error severity
func Error(a ...interface{}) {
	base.print(ERROR, a...)
}

// Warning logs a message with the Warning severity
func Warning(a ...interface{}) {
	base.print(WARNING, a...)
}

// Info logs a message with the Info severity
func Info(a ...interface{}) {
	base.print(INFO, a...)
}

// Debug logs a message with the Debug severity
func Debug(a ...interface{}) {
	base.print(DEBUG, a...)
}

// Panicf logs a message with the Critical severity using the same syntax and options as fmt.Printf then panic
func Panicf(format string, a ...interface{}) {
	s := fmt.Sprintf(format, a...)
	base.print(CRITICAL, s)
	panic(s)
}

// Criticalf logs a message with the Critical severity using the same syntax and options as fmt.Printf
func Criticalf(format string, a ...interface{}) {
	base.Critical(fmt.Sprintf(format, a...))
}

// Errorf logs a message with the Error severity using the same syntax and options as fmt.Printf
func Errorf(format string, a ...interface{}) {
	base.Error(fmt.Sprintf(format, a...))
}

// Warningf logs a message with the Warning severity using the same syntax and options as fmt.Printf
func Warningf(format string, a ...interface{}) {
	base.Warning(fmt.Sprintf(format, a...))
}

// Infof logs a message with the Info severity using the same syntax and options as fmt.Printf
func Infof(format string, a ...interface{}) {
	base.Info(fmt.Sprintf(format, a...))
}

// Debugf logs a message with the Debug severity using the same syntax and options as fmt.Printf
func Debugf(format string, a ...interface{}) {
	base.Debug(fmt.Sprintf(format, a...))
}
