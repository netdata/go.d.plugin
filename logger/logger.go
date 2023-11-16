// SPDX-License-Identifier: GPL-3.0-or-later

package logger

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/netdata/go.d.plugin/agent/executable"

	"github.com/mattn/go-isatty"
)

var pluginAttr = slog.String("plugin", executable.Name)

func New() *Logger {
	if isatty.IsTerminal(os.Stderr.Fd()) {
		// skip 2 slog pkg calls, 2 this pkg calls
		return &Logger{sl: slog.New(withCallDepth(4, newTerminalHandler()))}
	}
	return &Logger{sl: slog.New(newTextHandler()).With(pluginAttr)}
}

type Logger struct {
	sl *slog.Logger
}

func (l *Logger) Error(a ...any)                   { l.log(slog.LevelError, fmt.Sprint(a...)) }
func (l *Logger) Warning(a ...any)                 { l.log(slog.LevelWarn, fmt.Sprint(a...)) }
func (l *Logger) Info(a ...any)                    { l.log(slog.LevelInfo, fmt.Sprint(a...)) }
func (l *Logger) Debug(a ...any)                   { l.log(slog.LevelDebug, fmt.Sprint(a...)) }
func (l *Logger) Errorf(format string, a ...any)   { l.log(slog.LevelError, fmt.Sprintf(format, a...)) }
func (l *Logger) Warningf(format string, a ...any) { l.log(slog.LevelWarn, fmt.Sprintf(format, a...)) }
func (l *Logger) Infof(format string, a ...any)    { l.log(slog.LevelInfo, fmt.Sprintf(format, a...)) }
func (l *Logger) Debugf(format string, a ...any)   { l.log(slog.LevelDebug, fmt.Sprintf(format, a...)) }

func (l *Logger) With(args ...any) *Logger {
	if l.isNil() {
		return &Logger{sl: New().sl.With(args...)}
	}
	return &Logger{sl: l.sl.With(args...)}
}

func (l *Logger) log(level slog.Level, msg string) {
	if l.isNil() {
		nilLogger.sl.Log(context.Background(), level, msg)
	} else {
		l.sl.Log(context.Background(), level, msg)
	}
}

func (l *Logger) isNil() bool {
	return l == nil || l.sl == nil
}

var nilLogger = New()
