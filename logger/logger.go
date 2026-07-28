package logger

import (
	"log/slog"
	"os"
	"strings"
)

type Config struct {
	Level  string
	Format string
}

var (
	defaultLogger *slog.Logger
)

func Init(cfg Config) error {
	var level slog.Level
	switch strings.ToLower(cfg.Level) {
	case "debug":
		level = slog.LevelDebug
	case "info":
		level = slog.LevelInfo
	case "warn":
		level = slog.LevelWarn
	case "error":
		level = slog.LevelError
	default:
		level = slog.LevelInfo
	}

	var handler slog.Handler
	switch strings.ToLower(cfg.Format) {
	case "json":
		handler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level:     level,
			AddSource: level == slog.LevelDebug,
		})
	default:
		handler = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level:     level,
			AddSource: level == slog.LevelDebug,
		})
	}
	defaultLogger = slog.New(handler)
	slog.SetDefault(defaultLogger)

	return nil
}

func Info(msg string, v ...any) {
	if defaultLogger != nil {
		defaultLogger.Info(msg, v...)
	}
	slog.Info(msg, v...)
}

func Warn(msg string, v ...any) {
	if defaultLogger != nil {
		defaultLogger.Warn(msg, v...)
	}
	slog.Warn(msg, v...)
}

func Error(msg string, v ...any) {
	if defaultLogger != nil {
		defaultLogger.Error(msg, v...)
	}
	slog.Error(msg, v...)
}

func Debug(msg string, v ...any) {
	if defaultLogger != nil {
		defaultLogger.Debug(msg, v...)
	}
	slog.Debug(msg, v...)
}

func Fatal(msg string, v ...any) {
	if defaultLogger != nil {
		defaultLogger.Error(msg, v...)
	} else {
		slog.Error(msg, v...)
	}
	os.Exit(1)
}

func With(v ...any) *slog.Logger {
	if defaultLogger != nil {
		return defaultLogger.With(v...)
	}
	return slog.With(v...)
}

func WithGroup(group string) *slog.Logger {
	if defaultLogger != nil {
		return defaultLogger.WithGroup(group)
	}
	return slog.Default().WithGroup(group)
}

func Get() *slog.Logger {
	if defaultLogger != nil {
		return defaultLogger
	}
	return slog.Default()
}
