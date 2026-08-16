package logger

import (
	"io"
	"os"
	"strings"

	"github.com/rs/zerolog"
)

func New(level string, appEnv string) *zerolog.Logger {
	var l zerolog.Level

	switch strings.ToLower(level) {
	case "error":
		l = zerolog.ErrorLevel
	case "warn":
		l = zerolog.WarnLevel
	case "info":
		l = zerolog.InfoLevel
	case "debug":
		l = zerolog.DebugLevel
	default:
		l = zerolog.InfoLevel
	}

	zerolog.SetGlobalLevel(l)

	var skipFrameCount int
	var w io.Writer

	if appEnv == "production" {
		w = os.Stdout
		skipFrameCount = 3
	} else {
		w = zerolog.ConsoleWriter{Out: os.Stderr}
		skipFrameCount = 0
	}

	loggerContext := zerolog.New(w).With().Timestamp()

	if appEnv == "production" {
		loggerContext = loggerContext.CallerWithSkipFrameCount(zerolog.CallerSkipFrameCount + skipFrameCount)
	}

	logger := loggerContext.Logger()

	return &logger
}
