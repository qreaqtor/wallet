package logger

import (
	"context"
	"log/slog"
	"os"
)

type logger struct {
	log *slog.Logger
}

func New(logLevel Level, pretty bool) Log {
	opts := &slog.HandlerOptions{
		Level: toSlogLogLevel(logLevel),
	}

	var logHandler slog.Handler = slog.NewJSONHandler(os.Stdout, opts)

	if pretty {
		logHandler = newPrettyHandler(os.Stdout, opts)
	}

	return &logger{
		log: slog.New(logHandler),
	}
}

func (l logger) Log(ctx context.Context, level Level, msg string, args ...any) {
	l.log.Log(ctx, toSlogLogLevel(level), msg)
}

func (l logger) Debug(ctx context.Context, msg string, args ...any) {
	l.log.DebugContext(ctx, msg, args...)
}

func (l logger) Error(ctx context.Context, msg string, args ...any) {
	l.log.ErrorContext(ctx, msg, args...)
}

func (l logger) Info(ctx context.Context, msg string, args ...any) {
	l.log.InfoContext(ctx, msg, args...)
}

func (l logger) Warn(ctx context.Context, msg string, args ...any) {
	l.log.WarnContext(ctx, msg, args...)
}

func (l logger) WithFields(args ...any) Log {
	return logger{
		log: l.log.With(args...),
	}
}

func toSlogLogLevel(logLevel Level) slog.Level {
	switch logLevel {
	case LevelError:
		return slog.LevelError
	case LevelWarn:
		return slog.LevelWarn
	case LevelInfo:
		return slog.LevelInfo
	default:
		return slog.LevelDebug
	}
}
