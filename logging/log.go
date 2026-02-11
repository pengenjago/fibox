package logging

import (
	"strings"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var Logger zerolog.Logger = log.Logger

// Info logs an info message
func Info(msg string) {
	Logger.Info().Msg(msg)
}

// InfoWithFields logs an info message with additional fields
func InfoWithFields(msg string, fields map[string]interface{}) {
	event := Logger.Info()
	for k, v := range fields {
		event = event.Interface(k, v)
	}
	event.Msg(msg)
}

// Error logs an error message
func Error(msg string, err error) {
	if err != nil {
		Logger.Error().Err(err).Msg(msg)
	} else {
		Logger.Error().Msg(msg)
	}
}

// ErrorWithFields logs an error message with additional fields
func ErrorWithFields(msg string, err error, fields map[string]interface{}) {
	event := Logger.Error()
	if err != nil {
		event = event.Err(err)
	}
	for k, v := range fields {
		event = event.Interface(k, v)
	}
	event.Msg(msg)
}

// Debug logs a debug message
func Debug(msg string) {
	Logger.Debug().Msg(msg)
}

// DebugWithFields logs a debug message with additional fields
func DebugWithFields(msg string, fields map[string]interface{}) {
	event := Logger.Debug()
	for k, v := range fields {
		event = event.Interface(k, v)
	}
	event.Msg(msg)
}

// Warn logs a warning message
func Warn(msg string) {
	Logger.Warn().Msg(msg)
}

// WarnWithFields logs a warning message with additional fields
func WarnWithFields(msg string, fields map[string]interface{}) {
	event := Logger.Warn()
	for k, v := range fields {
		event = event.Interface(k, v)
	}
	event.Msg(msg)
}

// SetLogLevel sets the global log level
func SetLogLevel(level string) {
	zerolog.SetGlobalLevel(parseLogLevel(level))
}

func parseLogLevel(level string) zerolog.Level {
	switch strings.ToLower(level) {
	case "debug":
		return zerolog.DebugLevel
	case "info":
		return zerolog.InfoLevel
	case "warn", "warning":
		return zerolog.WarnLevel
	case "error":
		return zerolog.ErrorLevel
	case "fatal":
		return zerolog.FatalLevel
	case "panic":
		return zerolog.PanicLevel
	case "disabled", "none":
		return zerolog.Disabled
	default:
		return zerolog.InfoLevel
	}
}
