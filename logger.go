package bayaan

import (
	"fmt"
	"os"
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

	LoggerLevelCount
)

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

var level LoggerLevel = LoggerLevelInfo

// Helper function to log messages
func logMessage(lvl LoggerLevel, levelStr, f string, a ...any) {
	if level <= lvl {
		format := colors[lvl] + fmt.Sprintf("[%s] %s: ", timestamp(), levelStr) + Reset + f + "\n"
		fmt.Printf(format, a...)
	}
}

// Get the current timestamp
func timestamp() string {
	return time.Now().Format("2006-01-02 15:04:05")
}

func Trace(f string, a ...any) {
	logMessage(LoggerLevelTrace, "Trace", f, a...)
}

func Debug(f string, a ...any) {
	logMessage(LoggerLevelDebug, "Debug", f, a...)
}

func Info(f string, a ...any) {
	logMessage(LoggerLevelInfo, "Info", f, a...)
}

func Warn(f string, a ...any) {
	logMessage(LoggerLevelWarn, "Warn", f, a...)
}

func Error(f string, a ...any) (err error) {
	logMessage(LoggerLevelError, "Error", f, a...)
	return fmt.Errorf(f, a...)
}

func Fatal(f string, a ...any) {
	logMessage(LoggerLevelFatal, "Fatal", f, a...)
	os.Exit(1) // Exit after logging fatal error
}

func Panic(f string, a ...any) {
	logMessage(LoggerLevelPanic, "Panic", f, a...)
	panic(fmt.Sprintf(f, a...)) // Panic after logging
}

func SetLevel(newLevel LoggerLevel) {
	level = newLevel
}

func GetLevel() LoggerLevel {
	return level
}
