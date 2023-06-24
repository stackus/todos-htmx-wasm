package log

import (
	"time"

	"github.com/rs/zerolog"
)

// skip the first 2 callers when logging
const skipFrameCount = 2

// DefaultLogger is the default logger used by the package-level functions.
var DefaultLogger = zerolog.New(
	zerolog.NewConsoleWriter(
		func(w *zerolog.ConsoleWriter) {
			w.TimeFormat = time.RFC822
		},
	),
).
	Level(zerolog.DebugLevel).
	With().
	CallerWithSkipFrameCount(skipFrameCount).
	Timestamp().
	Logger()

// Trace returns a new event with the given message and level set to zerolog.TraceLevel.
func Trace() *zerolog.Event {
	return DefaultLogger.Trace()
}

// Debug returns a new event with the given message and level set to zerolog.DebugLevel.
func Debug() *zerolog.Event {
	return DefaultLogger.Debug()
}

// Info returns a new event with the given message and level set to zerolog.InfoLevel.
func Info() *zerolog.Event {
	return DefaultLogger.Info()
}

// Warn returns a new event with the given message and level set to zerolog.WarnLevel.
func Warn() *zerolog.Event {
	return DefaultLogger.Warn()
}

// Error returns a new event with the given message and level set to zerolog.ErrorLevel.
func Error() *zerolog.Event {
	return DefaultLogger.Error()
}

// Fatal returns a new event with the given message and level set to zerolog.PanicLevel.
func Fatal() *zerolog.Event {
	return DefaultLogger.Fatal()
}

// Err returns a new event with the given error and level set to zerolog.ErrorLevel.
func Err(err error) *zerolog.Event {
	return DefaultLogger.Err(err)
}

// With creates a child logger.
func With() zerolog.Context {
	return DefaultLogger.With()
}
