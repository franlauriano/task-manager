package logger

import (
	"log/slog"
	"os"
	"strings"
)

type Configuration struct {
	Level string `toml:"level"`
}

func (config Configuration) GetLevel() slog.Level {
	levelLog := slog.LevelInfo
	switch strings.ToLower(config.Level) {
	case "error":
		levelLog = slog.LevelError
	case "warn", "warning":
		levelLog = slog.LevelWarn
	case "info":
		levelLog = slog.LevelInfo
	case "debug", "trace":
		// trace is mapped to debug in slog
		levelLog = slog.LevelDebug
	}

	return levelLog
}

// Initialize configures and sets the default slog logger with JSON format
func (config Configuration) Initialize() {
	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: config.GetLevel(),
	})

	logger := slog.New(handler)
	slog.SetDefault(logger)
}
