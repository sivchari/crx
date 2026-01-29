package logger

import (
	"os"

	"github.com/charmbracelet/log"
)

// Init initializes the global logger with the specified verbosity.
func Init(verbose bool) {
	level := log.InfoLevel
	if verbose {
		level = log.DebugLevel
	}

	log.SetDefault(log.NewWithOptions(os.Stderr, log.Options{
		Level:           level,
		ReportTimestamp: false,
	}))
}

// Debug logs a debug message.
func Debug(msg string, keyvals ...any) {
	log.Debug(msg, keyvals...)
}

// Info logs an info message.
func Info(msg string, keyvals ...any) {
	log.Info(msg, keyvals...)
}

// Warn logs a warning message.
func Warn(msg string, keyvals ...any) {
	log.Warn(msg, keyvals...)
}

// Error logs an error message.
func Error(msg string, keyvals ...any) {
	log.Error(msg, keyvals...)
}

// Fatal logs a fatal message and exits.
func Fatal(msg string, keyvals ...any) {
	log.Fatal(msg, keyvals...)
}
