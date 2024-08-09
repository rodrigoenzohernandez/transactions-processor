package logger

import (
	"log/slog"
	"os"
	"strings"
)

var baseLogger *slog.Logger

func init() {

	logLevel := os.Getenv("LOG_LEVEL")
	if logLevel == "" {
		logLevel = "DEBUG"
	}

	var level slog.Level
	switch strings.ToUpper(logLevel) {
	case "DEBUG":
		level = slog.LevelDebug
	case "INFO":
		level = slog.LevelInfo
	case "WARN":
		level = slog.LevelWarn
	case "ERROR":
		level = slog.LevelError
	}

	baseLogger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: level}))
}

func GetLogger(prefix string) *slog.Logger {
	return baseLogger.With(slog.String("module", prefix))
}
