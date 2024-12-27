package bayaan

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"time"
)

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

type logEntry struct {
	level  LoggerLevel
	msg    string
	fields Fields
}

type output struct {
	writer   io.Writer
	useColor bool
}

type Logger struct {
	level      LoggerLevel
	outputs    []output
	timeFormat string
	mu         sync.RWMutex
	fields     Fields
	logChan    chan logEntry
	done       chan struct{}
}

type Fields map[string]interface{}

type LoggerOption func(*Logger)

func NewLogger(options ...LoggerOption) *Logger {
	l := &Logger{
		level:      LoggerLevelInfo,
		outputs:    []output{{writer: os.Stdout, useColor: true}},
		timeFormat: "2006-01-02 15:04:05",
		fields:     make(Fields),
		logChan:    make(chan logEntry, 1000), // Buffered channel to prevent blocking
		done:       make(chan struct{}),
	}

	for _, option := range options {
		option(l)
	}

	go l.processLogs()

	return l
}

func WithLevel(level LoggerLevel) LoggerOption {
	return func(l *Logger) {
		l.mu.Lock()
		l.level = level
		l.mu.Unlock()
	}
}

func WithOutput(writer io.Writer, additive bool, useColor bool) LoggerOption {
	return func(l *Logger) {
		l.mu.Lock()
		if additive {
			l.outputs = append(l.outputs, output{writer: writer, useColor: useColor})
		} else {
			l.outputs = []output{{writer: writer, useColor: useColor}}
		}
		l.mu.Unlock()
	}
}

func WithTimeFormat(format string) LoggerOption {
	return func(l *Logger) {
		l.mu.Lock()
		l.timeFormat = format
		l.mu.Unlock()
	}
}

func WithFields(fields Fields) LoggerOption {
	return func(l *Logger) {
		l.mu.Lock()
		for k, v := range fields {
			l.fields[k] = v
		}
		l.mu.Unlock()
	}
}

func (l *Logger) processLogs() {
	for {
		select {
		case entry := <-l.logChan:
			l.writeLog(entry)
		case <-l.done:
			return
		}
	}
}

func (l *Logger) writeLog(entry logEntry) {
	if entry.level < l.level {
		return
	}

	l.mu.RLock()
	timeFormat := l.timeFormat
	defaultFields := make(Fields, len(l.fields))
	for k, v := range l.fields {
		defaultFields[k] = v
	}
	outputs := make([]output, len(l.outputs))
	copy(outputs, l.outputs)
	l.mu.RUnlock()

	mergedFields := make(Fields)
	for k, v := range defaultFields {
		mergedFields[k] = v
	}
	for k, v := range entry.fields {
		mergedFields[k] = v
	}

	mergedFields["timestamp"] = time.Now().Format(timeFormat)
	mergedFields["level"] = entry.level.String()
	mergedFields["message"] = entry.msg

	if _, file, line, ok := runtime.Caller(3); ok {
		mergedFields["caller"] = fmt.Sprintf("%s:%d", filepath.Base(file), line)
	}

	output, err := json.Marshal(mergedFields)
	if err != nil {
		output = []byte(fmt.Sprintf("error marshaling log entry: %v", err))
	}

	for _, out := range outputs {
		logLine := string(output) + "\n"
		if out.useColor {
			logLine = colors[entry.level] + logLine + Reset
		}
		_, _ = fmt.Fprint(out.writer, logLine)
	}
}

func (l *Logger) Close() {
	close(l.done)
}

func (l *Logger) log(level LoggerLevel, msg string, fields Fields) {
	select {
	case l.logChan <- logEntry{level: level, msg: msg, fields: fields}:
	default:
		// Channel is full, log a warning and drop the message
		fmt.Fprintf(os.Stderr, "Warning: Logger channel full, dropping message: %s\n", msg)
	}
}

func (l *Logger) With(fields Fields) *Logger {
	l.mu.RLock()
	newLogger := &Logger{
		level:      l.level,
		outputs:    make([]output, len(l.outputs)),
		timeFormat: l.timeFormat,
		fields:     make(Fields),
		logChan:    l.logChan,
		done:       l.done,
	}
	copy(newLogger.outputs, l.outputs)

	for k, v := range l.fields {
		newLogger.fields[k] = v
	}
	l.mu.RUnlock()

	for k, v := range fields {
		newLogger.fields[k] = v
	}

	return newLogger
}

func (l *Logger) Trace(msg string, fields Fields) {
	l.log(LoggerLevelTrace, msg, fields)
}

func (l *Logger) Debug(msg string, fields Fields) {
	l.log(LoggerLevelDebug, msg, fields)
}

func (l *Logger) Info(msg string, fields Fields) {
	l.log(LoggerLevelInfo, msg, fields)
}

func (l *Logger) Warn(msg string, fields Fields) {
	l.log(LoggerLevelWarn, msg, fields)
}

func (l *Logger) Error(msg string, fields Fields) error {
	l.log(LoggerLevelError, msg, fields)
	return errors.New(msg)
}

func (l *Logger) Fatal(msg string, fields Fields) {
	l.log(LoggerLevelFatal, msg, fields)
	os.Exit(1)
}

func (l *Logger) Panic(msg string, fields Fields) {
	l.log(LoggerLevelPanic, msg, fields)
	panic(msg)
}

var defaultLogger *Logger

// Setup initializes the default logger with the provided options.
// If no options are provided, it uses sensible defaults.
// This should be called early in your application's lifecycle.
func Setup(options ...LoggerOption) {
	if len(options) == 0 {
		// Set sensible defaults
		options = []LoggerOption{
			WithLevel(LoggerLevelInfo),
			WithTimeFormat("2006-01-02 15:04:05"),
			WithOutput(os.Stdout, false, true), // Set stdout as default output with color enabled
			WithFields(Fields{
				"app": os.Getenv("APP_NAME"),
				"env": os.Getenv("APP_ENV"),
			}),
		}

		// Add file output if LOG_FILE is set
		if logFile := os.Getenv("LOG_FILE"); logFile != "" {
			if f, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644); err == nil {
				options = append(options, WithOutput(f, true, false)) // Append file output with color disabled
			}
		}
	}

	// Close existing logger if it exists
	if defaultLogger != nil {
		defaultLogger.Close()
	}

	defaultLogger = NewLogger(options...)
}

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

func SetLevel(level LoggerLevel) {
	defaultLogger = NewLogger(WithLevel(level))
}

func GetLevel() LoggerLevel {
	var level LoggerLevel
	defaultLogger.mu.RLock()
	level = defaultLogger.level
	defaultLogger.mu.RUnlock()
	return level
}
