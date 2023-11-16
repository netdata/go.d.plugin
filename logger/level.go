// SPDX-License-Identifier: GPL-3.0-or-later

package logger

import (
	"log/slog"
	"strings"
)

var Level = &level{lvl: &slog.LevelVar{}}

type level struct {
	lvl *slog.LevelVar
}

func (l *level) Enabled(level slog.Level) bool {
	return level >= l.lvl.Level()
}

func (l *level) Set(level slog.Level) {
	l.lvl.Set(level)
}

func (l *level) SetByName(level string) {
	switch strings.ToLower(level) {
	case "err", "error":
		Level.Set(slog.LevelError)
	case "warn", "warning":
		Level.Set(slog.LevelWarn)
	case "info":
		Level.Set(slog.LevelInfo)
	case "debug":
		Level.Set(slog.LevelDebug)
	}
}
