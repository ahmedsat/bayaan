package bayaan

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"time"
)

// LoggerLevel represents the severity level of a log message
type LoggerLevel int

const (
	LoggerLevelTrace LoggerLevel = iota
	LoggerLevelDebug
	LoggerLevelInfo
	LoggerLevelWarn
	LoggerLevelError
	LoggerLevelFatal
	LoggerLevelPanic
)

// String returns the string representation of the log level
func (l LoggerLevel) String() string {
	switch l {
	case LoggerLevelTrace:
		return "TRACE"
	case LoggerLevelDebug:
		return "DEBUG"
	case LoggerLevelInfo:
		return "INFO"
	case LoggerLevelWarn:
		return "WARN"
	case LoggerLevelError:
		return "ERROR"
	case LoggerLevelFatal:
		return "FATAL"
	case LoggerLevelPanic:
		return "PANIC"
	default:
		return "UNKNOWN"
	}
}

var colors = map[LoggerLevel]string{
	LoggerLevelTrace: "\033[36m", // Cyan
	LoggerLevelDebug: "\033[34m", // Blue
	LoggerLevelInfo:  "\033[32m", // Green
	LoggerLevelWarn:  "\033[33m", // Yellow
	LoggerLevelError: "\033[31m", // Red
	LoggerLevelFatal: "\033[35m", // Magenta
	LoggerLevelPanic: "\033[95m", // Bright Magenta
}

const Reset = "\033[0m"

// Logger represents a logger instance with its configuration
type Logger struct {
	level      LoggerLevel
	outputs    []io.Writer
	timeFormat string
	useColor   bool
	mu         sync.Mutex
	fields     Fields
}

// Fields represents structured logging fields
type Fields map[string]interface{}

// LoggerOption is a function that configures a Logger
type LoggerOption func(*Logger)

// NewLogger creates a new logger with the given options
func NewLogger(options ...LoggerOption) *Logger {
	l := &Logger{
		level:      LoggerLevelInfo,
		outputs:    []io.Writer{os.Stdout},
		timeFormat: "2006-01-02 15:04:05",
		useColor:   true,
		fields:     make(Fields),
	}

	for _, option := range options {
		option(l)
	}

	return l
}

// WithLevel sets the logger level
func WithLevel(level LoggerLevel) LoggerOption {
	return func(l *Logger) {
		l.level = level
	}
}

// WithOutput adds an output writer
func WithOutput(output io.Writer) LoggerOption {
	return func(l *Logger) {
		l.outputs = append(l.outputs, output)
	}
}

// WithTimeFormat sets the time format
func WithTimeFormat(format string) LoggerOption {
	return func(l *Logger) {
		l.timeFormat = format
	}
}

// WithColor enables or disables colored output
func WithColor(useColor bool) LoggerOption {
	return func(l *Logger) {
		l.useColor = useColor
	}
}

// WithFields adds default fields to the logger
func WithFields(fields Fields) LoggerOption {
	return func(l *Logger) {
		for k, v := range fields {
			l.fields[k] = v
		}
	}
}

// log formats and writes a log message
func (l *Logger) log(level LoggerLevel, msg string, fields Fields) {
	if level < l.level {
		return
	}

	l.mu.Lock()
	defer l.mu.Unlock()

	// Merge default fields with message fields
	mergedFields := make(Fields)
	for k, v := range l.fields {
		mergedFields[k] = v
	}
	for k, v := range fields {
		mergedFields[k] = v
	}

	// Add basic log information
	mergedFields["timestamp"] = time.Now().Format(l.timeFormat)
	mergedFields["level"] = level.String()
	mergedFields["message"] = msg

	// Add caller information
	if _, file, line, ok := runtime.Caller(2); ok {
		mergedFields["caller"] = fmt.Sprintf("%s:%d", filepath.Base(file), line)
	}

	var output []byte
	var err error
	output, err = json.Marshal(mergedFields)
	if err != nil {
		output = []byte(fmt.Sprintf("error marshaling log entry: %v", err))
	}

	logLine := string(output) + "\n"
	if l.useColor {
		logLine = colors[level] + logLine + Reset
	}

	for _, w := range l.outputs {
		_, _ = fmt.Fprint(w, logLine)
	}
}

// With creates a new logger with additional fields
func (l *Logger) With(fields Fields) *Logger {
	newLogger := &Logger{
		level:      l.level,
		outputs:    l.outputs,
		timeFormat: l.timeFormat,
		useColor:   l.useColor,
		fields:     make(Fields),
	}

	// Copy existing fields
	for k, v := range l.fields {
		newLogger.fields[k] = v
	}
	// Add new fields
	for k, v := range fields {
		newLogger.fields[k] = v
	}

	return newLogger
}

// Trace logs a message at trace level
func (l *Logger) Trace(msg string, fields Fields) {
	l.log(LoggerLevelTrace, msg, fields)
}

// Debug logs a message at debug level
func (l *Logger) Debug(msg string, fields Fields) {
	l.log(LoggerLevelDebug, msg, fields)
}

// Info logs a message at info level
func (l *Logger) Info(msg string, fields Fields) {
	l.log(LoggerLevelInfo, msg, fields)
}

// Warn logs a message at warn level
func (l *Logger) Warn(msg string, fields Fields) {
	l.log(LoggerLevelWarn, msg, fields)
}

// Error logs a message at error level
func (l *Logger) Error(msg string, fields Fields) error {
	l.log(LoggerLevelError, msg, fields)
	return fmt.Errorf(msg)
}

// Fatal logs a message at fatal level and exits
func (l *Logger) Fatal(msg string, fields Fields) {
	l.log(LoggerLevelFatal, msg, fields)
	os.Exit(1)
}

// Panic logs a message at panic level and panics
func (l *Logger) Panic(msg string, fields Fields) {
	l.log(LoggerLevelPanic, msg, fields)
	panic(msg)
}

// Default logger instance
var defaultLogger = NewLogger()

// Global functions that use the default logger

func Trace(msg string, fields Fields) {
	defaultLogger.Trace(msg, fields)
}

func Debug(msg string, fields Fields) {
	defaultLogger.Debug(msg, fields)
}

func Info(msg string, fields Fields) {
	defaultLogger.Info(msg, fields)
}

func Warn(msg string, fields Fields) {
	defaultLogger.Warn(msg, fields)
}

func Error(msg string, fields Fields) error {
	return defaultLogger.Error(msg, fields)
}

func Fatal(msg string, fields Fields) {
	defaultLogger.Fatal(msg, fields)
}

func Panic(msg string, fields Fields) {
	defaultLogger.Panic(msg, fields)
}

// SetLevel sets the level of the default logger
func SetLevel(level LoggerLevel) {
	defaultLogger = NewLogger(WithLevel(level))
}

// GetLevel returns the level of the default logger
func GetLevel() LoggerLevel {
	return defaultLogger.level
}
