// SPDX-License-Identifier: GPL-3.0-or-later

package logger

import (
	"context"
	"log/slog"
	"os"
	"runtime"
	"strings"

	"github.com/netdata/go.d.plugin/agent/executable"

	"github.com/lmittmann/tint"
)

var pluginAttr = slog.String("plugin", executable.Name)

var base = newBaseLogger()

func newBaseLogger() *Logger {
	return New()
}

func newTerminalLogger() *Logger {
	// skip Callers, this function, 2 slog pkg calls, 2 this pkg calls
	h := withCallDepth(6, tint.NewHandler(os.Stderr, &tint.Options{
		AddSource: true,
		Level:     Level,
	}))

	return &Logger{sl: slog.New(h)}
}

func newLogger() *Logger {
	h := slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		Level: Level,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == slog.LevelKey {
				if v, ok := a.Value.Any().(slog.Level); ok {
					a.Value = slog.StringValue(strings.ToLower(v.String()))
				}
			}
			return a
		},
	})

	return &Logger{sl: slog.New(h)}
}

func withCallDepth(depth int, sh slog.Handler) slog.Handler {
	return &callDepthHandler{depth: depth, sh: sh}
}

type callDepthHandler struct {
	depth int
	sh    slog.Handler
}

func (h *callDepthHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return h.sh.Enabled(ctx, level)
}

func (h *callDepthHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &callDepthHandler{depth: h.depth, sh: h.sh.WithAttrs(attrs)}
}

func (h *callDepthHandler) WithGroup(name string) slog.Handler {
	return &callDepthHandler{depth: h.depth, sh: h.sh.WithGroup(name)}
}

func (h *callDepthHandler) Handle(ctx context.Context, r slog.Record) error {
	// https://pkg.go.dev/log/slog#example-package-Wrapping
	var pcs [1]uintptr
	runtime.Callers(h.depth, pcs[:])
	r.PC = pcs[0]

	return h.sh.Handle(ctx, r)
}
