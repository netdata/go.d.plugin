// SPDX-License-Identifier: GPL-3.0-or-later

package slogger

import (
	"log/slog"
	"os"
	"strings"

	"github.com/netdata/go.d.plugin/agent/executable"

	"github.com/lmittmann/tint"
)

var base = newBaseLogger()

var pluginAttr = slog.String("plugin", executable.Name)

func newBaseLogger() *Logger {
	return New()
}

func newTerminalLogger() *Logger {
	return &Logger{
		sl: slog.New(tint.NewHandler(os.Stderr, &tint.Options{
			AddSource: Level.Level() == slog.LevelDebug,
			Level:     Level,
		})),
	}
}

func newLogger() *Logger {
	return &Logger{
		sl: slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
			Level: Level,
			ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
				if a.Key == slog.LevelKey {
					if v, ok := a.Value.Any().(slog.Level); ok {
						a.Value = slog.StringValue(strings.ToLower(v.String()))
					}
				}
				return a
			},
		})),
	}
}
