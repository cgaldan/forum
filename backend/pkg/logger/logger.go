package logger

import (
	"fmt"
	"io"
	"log"
	"os"
	"runtime"
	"strings"
	"time"
)

// Level represents log level
type Level int

const (
	// DebugLevel for debug messages
	DebugLevel Level = iota
	// InfoLevel for informational messages
	InfoLevel
	// WarnLevel for warning messages
	WarnLevel
	// ErrorLevel for error messages
	ErrorLevel
	// FatalLevel for fatal messages
	FatalLevel
)

// String returns string representation of log level
func (l Level) String() string {
	switch l {
	case DebugLevel:
		return "DEBUG"
	case InfoLevel:
		return "INFO"
	case WarnLevel:
		return "WARN"
	case ErrorLevel:
		return "ERROR"
	case FatalLevel:
		return "FATAL"
	default:
		return "UNKNOWN"
	}
}

// Logger represents a logger
type Logger struct {
	logger *log.Logger
	level  Level
}

// NewLogger creates a new logger
func NewLogger(output io.Writer, level Level) *Logger {
	return &Logger{
		logger: log.New(output, "", 0),
		level:  level,
	}
}

// SetLevel sets the logging level
func (l *Logger) SetLevel(level Level) {
	l.level = level
}

func (l *Logger) log(level Level, message string, args ...interface{}) {
	if level < l.level {
		return
	}

	// Get caller information
	_, file, line, ok := runtime.Caller(2)
	caller := "???"
	if ok {
		// Extract just the filename
		parts := strings.Split(file, "/")
		caller = fmt.Sprintf("%s:%d", parts[len(parts)-1], line)
	}

	timestamp := time.Now().Format("2006-01-02 15:04:05")
	prefix := fmt.Sprintf("[%s] %s %s", timestamp, level.String(), caller)

	// Format additional arguments
	var argsStr string
	if len(args) > 0 {
		var pairs []string
		for i := 0; i < len(args); i += 2 {
			if i+1 < len(args) {
				pairs = append(pairs, fmt.Sprintf("%v=%v", args[i], args[i+1]))
			}
		}
		if len(pairs) > 0 {
			argsStr = " | " + strings.Join(pairs, ", ")
		}
	}

	l.logger.Printf("%s - %s%s\n", prefix, message, argsStr)
}

// Debug logs a debug message
func (l *Logger) Debug(message string, args ...interface{}) {
	l.log(DebugLevel, message, args...)
}

// Info logs an informational message
func (l *Logger) Info(message string, args ...interface{}) {
	l.log(InfoLevel, message, args...)
}

// Warn logs a warning message
func (l *Logger) Warn(message string, args ...interface{}) {
	l.log(WarnLevel, message, args...)
}

// Error logs an error message
func (l *Logger) Error(message string, args ...interface{}) {
	l.log(ErrorLevel, message, args...)
}

// Fatal logs a fatal message and exits
func (l *Logger) Fatal(message string, args ...interface{}) {
	l.log(FatalLevel, message, args...)
	os.Exit(1)
}

