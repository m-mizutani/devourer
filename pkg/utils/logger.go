package utils

import (
	"log/slog"
	"os"
	"sync"
)

var (
	logger      = slog.Default()
	loggerMutex sync.Mutex
)

func init() {
	logger = slog.New(slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{
		AddSource: true,
		Level:     slog.LevelInfo,
	}))
}

func Logger() *slog.Logger {
	return logger
}

func ReconfigureLogger(handler slog.Handler) {
	loggerMutex.Lock()
	logger = slog.New(handler)
	loggerMutex.Unlock()
}

func ErrLog(err error) slog.Attr { return slog.Any("error", err) }
