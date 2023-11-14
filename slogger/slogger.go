// SPDX-License-Identifier: GPL-3.0-or-later

package slogger

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/mattn/go-isatty"
)

func New() *Logger {
	var l *Logger

	if isatty.IsTerminal(os.Stderr.Fd()) {
		l = newTerminalLogger()
	} else {
		l = newLogger()
	}

	l = l.With(pluginAttr)

	return l
}

type Logger struct {
	sl *slog.Logger
}

func (l *Logger) Errorf(format string, a ...any) {
	l.logf(slog.LevelError, format, a...)
}

func (l *Logger) Warningf(format string, a ...any) {
	l.logf(slog.LevelWarn, format, a...)
}

func (l *Logger) Infof(format string, a ...any) {
	l.logf(slog.LevelInfo, format, a...)
}

func (l *Logger) Debugf(format string, a ...any) {
	l.logf(slog.LevelDebug, format, a...)
}

func (l *Logger) logf(level slog.Level, format string, a ...any) {
	if level < Level.Level() {
		return
	}

	msg := fmt.Sprintf(format, a...)

	if l.isNil() {
		base.sl.Log(context.Background(), level, msg)
	} else {
		l.sl.Log(context.Background(), level, msg)
	}
}

func (l *Logger) With(args ...any) *Logger {
	if l.isNil() {
		return &Logger{sl: newBaseLogger().sl.With(args...)}
	}
	return &Logger{sl: l.sl.With(args...)}
}

func (l *Logger) isNil() bool {
	return l == nil || l.sl == nil
}
