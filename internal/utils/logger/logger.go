package logger

import (
	"log/slog"
	"os"
)

var baseLogger *slog.Logger

func init() {
	baseLogger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{}))
}

func GetLogger(prefix string) *slog.Logger {
	return baseLogger.With(slog.String("module", prefix))
}
