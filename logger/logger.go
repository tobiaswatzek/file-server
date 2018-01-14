package logger

import (
	"bytes"
	"log"
	"net/http"
	"os"
	"time"

	uuid "github.com/satori/go.uuid"
)

// LogLevel is the log level type
type LogLevel int

const (
	// Undefined is returned when no log level is set.
	Undefined LogLevel = iota
	// Debug is the debug log level
	Debug
	// Info is the info log level
	Info
	// Warning is the warning log level
	Warning
	// Error is the error log level
	Error
)

// Logger is a wrapper for the default log.
type Logger interface {
	Middleware(next http.Handler) http.Handler
	SetLogLevel(level LogLevel)
	GetLogLevel() LogLevel
	Debug(format string, p ...interface{})
	DebugWithContext(ctx string, format string, p ...interface{})
	Info(format string, p ...interface{})
	InfoWithContext(ctx string, format string, p ...interface{})
	Warning(format string, p ...interface{})
	WarningWithContext(ctx string, format string, p ...interface{})
	Error(format string, p ...interface{})
	ErrorWithContext(ctx string, format string, p ...interface{})
}

type logger struct {
	logLevel LogLevel
	logger   *log.Logger
}

// NewLoggerWithLevel creates a new logger with the given level.
func NewLoggerWithLevel(level LogLevel) Logger {
	l := &logger{}
	l.logger = log.New(os.Stderr, "", log.LstdFlags)
	l.SetLogLevel(level)
	return l
}

// NewLogger creates a new logger with the Info log level.
func NewLogger() Logger {
	return NewLoggerWithLevel(Info)
}

func (l *logger) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := uuid.NewV4().String()
		l.InfoWithContext(ctx, "Received request for '%s'.", r.RequestURI)
		begin := time.Now()
		next.ServeHTTP(w, r)
		end := time.Now()
		processedIn := end.Sub(begin)
		l.InfoWithContext(ctx, "Processed request for '%s' in %s.", r.RequestURI, processedIn)
	})
}

func (l *logger) SetLogLevel(level LogLevel) {
	if l == nil {
		return
	}
	l.logLevel = level
}

func (l *logger) GetLogLevel() LogLevel {
	if l == nil {
		return Undefined
	}
	return l.logLevel
}

func (l *logger) Debug(format string, p ...interface{}) {
	l.logWithLevel(Debug, format, p)
}

func (l *logger) DebugWithContext(ctx string, format string, p ...interface{}) {
	l.logWithContext(Debug, ctx, format, p)
}

func (l *logger) Info(format string, p ...interface{}) {
	l.logWithLevel(Info, format, p)
}

func (l *logger) InfoWithContext(ctx string, format string, p ...interface{}) {
	l.logWithContext(Info, ctx, format, p)
}

func (l *logger) Warning(format string, p ...interface{}) {
	l.logWithLevel(Warning, format, p)
}

func (l *logger) WarningWithContext(ctx string, format string, p ...interface{}) {
	l.logWithContext(Warning, ctx, format, p)
}

func (l *logger) Error(format string, p ...interface{}) {
	l.logWithLevel(Error, format, p)
}

func (l *logger) ErrorWithContext(ctx string, format string, p ...interface{}) {
	l.logWithContext(Error, ctx, format, p)
}

func (l *logger) shouldLog(level LogLevel) bool {
	return l != nil && l.GetLogLevel() <= level
}

func logLevelToString(level LogLevel) string {
	switch level {
	case Debug:
		return "Debug"
	case Info:
		return "Info"
	case Warning:
		return "Warning"
	case Error:
		return "Error"
	default:
		return "Undefined"
	}
}

func prependLogLevel(level LogLevel, s string) string {
	var buffer bytes.Buffer
	buffer.WriteRune('[')
	buffer.WriteString(logLevelToString(level))
	buffer.WriteString("] ")
	buffer.WriteString(s)
	return buffer.String()
}

func (l *logger) logWithContext(level LogLevel, ctx, format string, p []interface{}) {
	if !l.shouldLog(level) {
		return
	}
	var buffer bytes.Buffer
	buffer.WriteRune('{')
	buffer.WriteString(ctx)
	buffer.WriteString("} ")
	buffer.WriteString(format)
	l.logWithLevel(level, buffer.String(), p)
}

func (l *logger) logWithLevel(level LogLevel, format string, p []interface{}) {
	if !l.shouldLog(level) {
		return
	}
	l.log(prependLogLevel(level, format), p)
}

func (l *logger) log(format string, p []interface{}) {
	l.logger.Printf(format, p...)
}
