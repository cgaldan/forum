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

type Level int

const (
	DebugLevel Level = iota
	InfoLevel
	WarnLevel
	ErrorLevel
	FatalLevel
)

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

type Logger struct {
	logger *log.Logger
	level  Level
}

func NewLogger(output io.Writer, level Level) *Logger {
	return &Logger{
		logger: log.New(output, "", 0),
		level:  level,
	}
}

func (l *Logger) SetLevel(level Level) {
	l.level = level
}

func (l *Logger) log(level Level, message string, args ...any) {
	if level < l.level {
		return
	}

	_, file, line, ok := runtime.Caller(2)
	caller := "???"
	if ok {
		parts := strings.Split(file, "/")
		caller = fmt.Sprintf("%s:%d", parts[len(parts)-1], line)
	}

	timestamp := time.Now().Format("2006-01-02 15:04:05")
	prefix := fmt.Sprintf("[%s] %s %s ", timestamp, level.String(), caller)

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

	l.logger.Printf("%s%s%s", prefix, message, argsStr)
}

func (l *Logger) Debug(message string, args ...any) {
	l.log(DebugLevel, message, args...)
}

func (l *Logger) Info(message string, args ...any) {
	l.log(InfoLevel, message, args...)
}

func (l *Logger) Warn(message string, args ...any) {
	l.log(WarnLevel, message, args...)
}

func (l *Logger) Error(message string, args ...any) {
	l.log(ErrorLevel, message, args...)
}

func (l *Logger) Fatal(message string, args ...any) {
	l.log(FatalLevel, message, args...)
	os.Exit(1)
}
