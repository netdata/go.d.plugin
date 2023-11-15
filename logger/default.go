// SPDX-License-Identifier: GPL-3.0-or-later

package logger

func newDefaultLogger() *Logger {
	return New()
}

var defaultLogger = newDefaultLogger()

func Error(a ...any)                   { defaultLogger.Error(a...) }
func Warning(a ...any)                 { defaultLogger.Warning(a...) }
func Info(a ...any)                    { defaultLogger.Info(a...) }
func Debug(a ...any)                   { defaultLogger.Debug(a...) }
func Errorf(format string, a ...any)   { defaultLogger.Errorf(format, a...) }
func Warningf(format string, a ...any) { defaultLogger.Warningf(format, a...) }
func Infof(format string, a ...any)    { defaultLogger.Infof(format, a...) }
func Debugf(format string, a ...any)   { defaultLogger.Debugf(format, a...) }
func With(args ...any) *Logger         { return defaultLogger.With(args...) }
