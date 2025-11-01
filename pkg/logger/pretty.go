package logger

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	stdLog "log"
	"log/slog"
	"strings"

	"github.com/mitchellh/colorstring"
)

type prettyHandler struct {
	slog.Handler
	l     *stdLog.Logger
	attrs []slog.Attr
}

func newPrettyHandler(w io.Writer, opts *slog.HandlerOptions) *prettyHandler {
	return &prettyHandler{
		Handler: slog.NewTextHandler(w, opts),
		l:       stdLog.New(w, "", 0),
	}
}

func (h *prettyHandler) Handle(_ context.Context, r slog.Record) error {
	level := r.Level.String() + ":"

	switch r.Level {
	case slog.LevelDebug:
		level = colorstring.Color("[green]" + level)
	case slog.LevelInfo:
		level = colorstring.Color("[light_blue]" + level)
	case slog.LevelWarn:
		level = colorstring.Color("[yellow]" + level)
	case slog.LevelError:
		level = colorstring.Color("[red]" + level)
	}

	fields := make(map[string]any, r.NumAttrs())

	r.Attrs(func(a slog.Attr) bool {
		fields[a.Key] = a.Value.Any()
		return true
	})

	for _, a := range h.attrs {
		fields[a.Key] = a.Value.Any()
	}

	args := []any{
		colorstring.Color("[dark_gray]" + r.Time.Format("[15:04:05.000]")),
		level,
		colorstring.Color("[light_magenta]" + r.Message),
	}

	if len(fields) != 0 {
		var buf bytes.Buffer
		enc := json.NewEncoder(&buf)
		enc.SetEscapeHTML(false)
		enc.SetIndent("", "    ")
		enc.Encode(fields)

		args = append(args, strings.TrimSpace(buf.String()))
	}

	h.l.Println(args...)

	return nil
}

func (h *prettyHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &prettyHandler{
		Handler: h.Handler,
		l:       h.l,
		attrs:   attrs,
	}
}

func (h *prettyHandler) WithGroup(name string) slog.Handler {
	return &prettyHandler{
		Handler: h.Handler.WithGroup(name),
		l:       h.l,
	}
}
