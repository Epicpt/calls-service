package logger

import (
	"os"
	"strings"
	"time"

	"github.com/rs/zerolog"
)

func New(level string) zerolog.Logger {
	logLevel, err := zerolog.ParseLevel(strings.ToLower(level))
	if err != nil {
		logLevel = zerolog.InfoLevel
	}

	logger := zerolog.New(os.Stdout).
		With().
		Timestamp().
		Logger().
		Level(logLevel)

	if logLevel <= zerolog.DebugLevel {
		logger = logger.Output(zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339})
	}

	return logger
}
