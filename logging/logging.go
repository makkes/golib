// Package logging provides a dead-simple, level-aware logging facility.
package logging

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"runtime"
	"strings"
)

// Level defines all available log levels.
type Level uint8

// The various log levels.
const (
	DEBUG Level = iota
	INFO
	WARN
	ERROR
)

var levelNames = map[Level]string{
	DEBUG: "DEBUG",
	INFO:  "INFO",
	WARN:  "WARN",
	ERROR: "ERROR",
}

// LevelFromString parses a string and returns the approriate Level with the
// same name. If no such Level exists, a non-nil error is returned.
func LevelFromString(s string) (*Level, error) {
	s = strings.ToUpper(s)
	for k, v := range levelNames {
		if v == s {
			return &k, nil
		}
	}
	return nil, errors.New("no such level name exists")
}

func (l Level) String() string {
	return levelNames[l]
}

// A Logger is used to write a string to some output using a certain log level
// that is part of the written string.
type Logger interface {
	Debug(format string, v ...interface{})
	Info(format string, v ...interface{})
	Warn(format string, v ...interface{})
	Error(format string, v ...interface{})
	Fatal(format string, v ...interface{})
}

type stdoutLogger struct {
	level  Level
	prefix string
}

func (l *stdoutLogger) Debug(format string, v ...interface{}) {
	l.log(DEBUG, format, v...)
}

func (l *stdoutLogger) Info(format string, v ...interface{}) {
	l.log(INFO, format, v...)
}

func (l *stdoutLogger) Warn(format string, v ...interface{}) {
	l.log(WARN, format, v...)
}

func (l *stdoutLogger) Error(format string, v ...interface{}) {
	l.log(ERROR, format, v...)
}

func (l *stdoutLogger) Fatal(format string, v ...interface{}) {
	log.Fatalf("[%s] %s", ERROR.String(), fmt.Sprintf(format, v...))
}

func (l *stdoutLogger) log(level Level, format string, v ...interface{}) {
	if l.level > level {
		return
	}
	var prolog bytes.Buffer
	if len(l.prefix) > 0 {
		fmt.Fprintf(&prolog, "[%s] [%s] ", level.String(), l.prefix)
	} else {
		fmt.Fprintf(&prolog, "[%s] ", level.String())
	}
	_, file, line, _ := runtime.Caller(2)
	fileParts := strings.Split(file, "/")
	fmt.Fprintf(&prolog, "[%s:%d] ", fileParts[len(fileParts)-1], line)
	fmt.Fprintf(&prolog, format, v...)
	err := log.Output(4, prolog.String())
	if err != nil {
		Fatal("Could not output log message: %s", err)
	}
}

func (l *stdoutLogger) SetLevel(level Level) {
	l.level = level
}

var defaultLogger = &stdoutLogger{DEBUG, ""}

// NewLogger creates a new logging object that only logs messages above the
// given level and swallows all others. The optional prefix is prepended to
// each log output.
func NewLogger(level Level, prefix string) Logger {
	return &stdoutLogger{level, prefix}
}

// NewDefaultLevelLogger creates a new logging object. The log level is
// inherited from the default logger which is set via SetLevel. The optional
// prefix is prepended to each log output.
func NewDefaultLevelLogger(prefix string) Logger {
	return &stdoutLogger{defaultLogger.level, prefix}
}

// SetLevel sets the logging level of the default logger.
func SetLevel(level Level) {
	defaultLogger.SetLevel(level)
}

// Debug logs the given message with DEBUG level.
func Debug(format string, v ...interface{}) {
	defaultLogger.log(DEBUG, format, v...)
}

// Info logs the given message with INFO level.
func Info(format string, v ...interface{}) {
	defaultLogger.log(INFO, format, v...)
}

// Warn logs the given message with WARN level.
func Warn(format string, v ...interface{}) {
	defaultLogger.log(WARN, format, v...)
}

// Error logs the given message with ERROR level.
func Error(format string, v ...interface{}) {
	defaultLogger.log(ERROR, format, v...)
}

// Fatal logs the given message with ERROR level and exits (same as
// log.Fatal).
func Fatal(format string, v ...interface{}) {
	defaultLogger.Fatal(format, v...)
}
