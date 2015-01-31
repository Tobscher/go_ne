/*
Package logging provides basic logging features (modules, log levels).
*/
package logging

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/mgutz/ansi"
)

// Level defines the log level. Can be one of the following:
// * OFF
// * FATAL
// * ERROR
// * WARN
// * INFO
// * DEBUG
// * TRACE
type Level int

const (
	// OFF disables logging
	OFF Level = iota
	// FATAL logs fatals
	FATAL
	// ERROR logs at least errors
	ERROR
	// WARN logs at least warnings
	WARN
	// INFO logs at least infos
	INFO
	// DEBUG logs at least debug output
	DEBUG
	// TRACE logs everything (verbose!)
	TRACE
)

var levelNames = []string{
	"OFF",
	"FATAL",
	"ERROR",
	"WARNING",
	"INFO",
	"DEBUG",
	"TRACE",
}

var levelColors = []string{
	"white",
	"red+h",
	"red",
	"yellow+h",
	"green+h",
	"black+h",
	"black+h",
}

// DefaultLogLevel is set to INFO.
var DefaultLogLevel = INFO

func (ll Level) String() string {
	return levelNames[ll]
}

// Logger is used to extend log.Logger.
type Logger struct {
	*log.Logger
	Level  Level
	Module string
}

// GetLogger creates a new logger object with the given prefix.
func GetLogger(module string) *Logger {
	logger := &Logger{
		Logger: log.New(os.Stdout, "", 0),
		Level:  DefaultLogLevel,
		Module: module,
	}

	return logger
}

// SetLevel sets the current log leve.
// Messages with a lower level than the given level
// will be omitted.
func (l *Logger) SetLevel(level Level) {
	l.Level = level
}

// Trace logs trace level messages.
func (l *Logger) Trace(message string) {
	l.logLevel(TRACE, message)
}

// Tracef logs formatted trace level messages.
func (l *Logger) Tracef(message string, a ...interface{}) {
	l.logLevelf(TRACE, message, a...)
}

// Debug logs debug level messages.
func (l *Logger) Debug(message string) {
	l.logLevel(DEBUG, message)
}

// Debugf logs formatted debug level messages.
func (l *Logger) Debugf(message string, a ...interface{}) {
	l.logLevelf(DEBUG, message, a...)
}

// Info logs info level messages.
func (l *Logger) Info(message string) {
	l.logLevel(INFO, message)
}

// Infof logs formatted info level messages.
func (l *Logger) Infof(message string, a ...interface{}) {
	l.logLevelf(INFO, message, a...)
}

// Warn logs warn level messages.
func (l *Logger) Warn(message string) {
	l.logLevel(WARN, message)
}

// Warnf logs formatted warn level messages.
func (l *Logger) Warnf(message string, a ...interface{}) {
	l.logLevelf(WARN, message, a...)
}

// Error logs error level messages.
func (l *Logger) Error(message string) {
	l.logLevel(ERROR, message)
}

// Errorf logs formatted error level messages.
func (l *Logger) Errorf(message string, a ...interface{}) {
	l.logLevelf(ERROR, message, a...)
}

// Fatal logs fatal level messages.
func (l *Logger) Fatal(message string) {
	l.logLevel(FATAL, message)
}

// Fatalf logs formatted fatal level messages.
func (l *Logger) Fatalf(message string, a ...interface{}) {
	l.logLevelf(FATAL, message, a...)
}

func (l *Logger) logLevelWithLineEnding(level Level, message string, newLine string) {
	if l.Level < level {
		return
	}

	time := time.Now()
	formattedTime := time.Format("2006-01-02 15:04:05")

	formatted := ansi.Color(fmt.Sprintf("%v [%v] - %v - %v", formattedTime, level.String()[0:4], l.Module, message), levelColors[level])
	formatted = fmt.Sprintf("%v%v", formatted, newLine)
	l.Print(formatted)
}

func (l *Logger) logLevel(level Level, message string) {
	l.logLevelWithLineEnding(level, message, "\n")
}

func (l *Logger) logLevelf(level Level, message string, a ...interface{}) {
	l.logLevelWithLineEnding(level, fmt.Sprintf(message, a...), "")
}
