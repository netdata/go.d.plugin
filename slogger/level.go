// SPDX-License-Identifier: GPL-3.0-or-later

package slogger

import (
	"log/slog"
	"strings"
)

var Level = &slog.LevelVar{}

func SetLevelByName(level string) {
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
