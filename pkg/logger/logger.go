package logger

import (
	"log/slog"
	"os"
)

type Options struct {
	Level  slog.Level
	Format string
}

// Returns a configured *slog.Logger and sets it as the default.
func New(opts Options) *slog.Logger {
	var handler slog.Handler

	handlerOpts := &slog.HandlerOptions{
		Level:     opts.Level,
		AddSource: opts.Level == slog.LevelDebug,
	}

	if opts.Format == "json" {
		handler = slog.NewJSONHandler(os.Stdout, handlerOpts)
	} else {
		handler = slog.NewTextHandler(os.Stdout, handlerOpts)
	}

	log := slog.New(handler)
	slog.SetDefault(log)
	return log
}

// Converts a string like "debug", "info", "warn", "error".
func LevelFromString(s string) slog.Level {
	switch s {
	case "debug":
		return slog.LevelDebug
	case "warn":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}
