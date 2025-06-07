package logger

import (
	"os"
	"time"

	"github.com/rs/zerolog"
)

var Logger zerolog.Logger

func init() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.SetGlobalLevel(zerolog.InfoLevel)

	// Set up the logger to write to standard output
	Logger = zerolog.New(zerolog.ConsoleWriter{
		Out:        os.Stderr,
		NoColor:    false,
		TimeFormat: time.RFC3339,
	}).With().Timestamp().Logger()
}

func Warn() *zerolog.Event {
	return Logger.Warn()
}

func Info() *zerolog.Event {
	return Logger.Info()
}

func Panic() *zerolog.Event {
	return Logger.Panic()
}

func Debug() *zerolog.Event {
	return Logger.Debug()
}

func Error() *zerolog.Event {
	return Logger.Error()
}

func Fatal() *zerolog.Event {
	return Logger.Fatal()
}

func Trace() *zerolog.Event {
	return Logger.Trace()
}
