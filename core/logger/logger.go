package logger

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"github.com/lmittmann/tint"
)

type Log struct {
	Handler       slog.Handler
	customHandler *slog.Handler
}

// WithColor provides an optional configuration to colorize the output logs.
func WithColor() func(*Log) {
	return func(l *Log) {
		ch := tint.NewHandler(os.Stdout, &tint.Options{
			Level:      slog.LevelInfo,
			TimeFormat: time.Kitchen,
			AddSource:  true,
		})

		l.customHandler = &ch
	}
}

func New(w io.Writer, serviceName string, opts ...func(*Log)) *Log {
	log := &Log{}

	for _, o := range opts {
		o(log)
	}

	// Attributes to add to every log
	attrs := []slog.Attr{
		{Key: "service", Value: slog.StringValue(serviceName)},
	}

	if log.customHandler != nil {
		ch := *log.customHandler
		ch.WithAttrs(attrs)

		return &Log{
			Handler: ch.WithAttrs(attrs),
		}
	}

	// convert filename to just the name.ext
	f := func(groups []string, a slog.Attr) slog.Attr {
		if a.Key == slog.SourceKey {
			if source, ok := a.Value.Any().(*slog.Source); ok {
				v := fmt.Sprintf("%s:%d", filepath.Base(source.File), source.Line)
				return slog.Attr{Key: "file", Value: slog.StringValue(v)}
			}
		}

		return a
	}

	handler := slog.Handler(slog.NewJSONHandler(w, &slog.HandlerOptions{
		AddSource:   true,
		Level:       slog.LevelInfo,
		ReplaceAttr: f,
	}))
	handler = handler.WithAttrs(attrs)

	return &Log{
		Handler: handler,
	}
}

func (log *Log) write(ctx context.Context, level slog.Level, caller int, msg string, args ...any) {
	var pcs [1]uintptr
	runtime.Callers(caller, pcs[:])

	r := slog.NewRecord(time.Now(), level, msg, pcs[0])
	r.Add(args...)

	err := log.Handler.Handle(ctx, r)
	if err != nil {
		panic(fmt.Errorf("unable to write logs: %w", err))
	}
}

func (log *Log) Info(ctx context.Context, msg string, args ...any) {
	log.write(ctx, slog.LevelInfo, 3, msg, args...)
}

func (log *Log) Debug(ctx context.Context, msg string, args ...any) {
	log.write(ctx, slog.LevelDebug, 3, msg, args...)
}

func (log *Log) Warn(ctx context.Context, msg string, args ...any) {
	log.write(ctx, slog.LevelWarn, 3, msg, args...)
}

func (log *Log) Error(ctx context.Context, msg string, args ...any) {
	log.write(ctx, slog.LevelError, 3, msg, args...)
}

// ErrorWithCaller lets you define your own custom call stack.
func (log *Log) ErrorWithCaller(ctx context.Context, caller int, msg string, args ...any) {
	log.write(ctx, slog.LevelError, caller, msg, args...)
}
