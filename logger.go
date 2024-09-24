package bayaan

import "fmt"

type LoggerLevel int

const (
	LoggerLevelDebug LoggerLevel = iota
	LoggerLevelInfo
	LoggerLevelWarn
	LoggerLevelError
	LoggerLevelFatal
	LoggerLevelPanic

	LoggerLevelCount
)

var colors = map[LoggerLevel]string{
	LoggerLevelDebug: "\033[32m",
	LoggerLevelInfo:  "\033[34m",
	LoggerLevelWarn:  "\033[36m",
	LoggerLevelError: "\033[33m",
	LoggerLevelFatal: "\033[35m",
	LoggerLevelPanic: "\033[31m",
}

var (
	Reset = "\033[0m"
	Gray  = "\033[37m"
	White = "\033[97m"

	Green   = "\033[32m"
	Blue    = "\033[34m"
	Cyan    = "\033[36m"
	Yellow  = "\033[33m"
	Magenta = "\033[35m"
	Red     = "\033[31m"
)

var level LoggerLevel = LoggerLevelInfo

func Debug(f string, a ...any) {
	format := colors[LoggerLevelDebug] + "Debug: " + Reset + f
	if level <= LoggerLevelDebug {
		fmt.Printf(format, a...)
	}
}

func Info(f string, a ...any) {
	format := colors[LoggerLevelDebug] + "Info: " + Reset + f
	if level <= LoggerLevelInfo {
		fmt.Printf(format, a...)
	}
}

func Warn(f string, a ...any) {
	format := colors[LoggerLevelDebug] + "Warn: " + Reset + f
	if level <= LoggerLevelWarn {
		fmt.Printf(format, a...)
	}
}

func Error(f string, a ...any) {
	format := colors[LoggerLevelDebug] + "Error: " + Reset + f
	if level <= LoggerLevelError {
		fmt.Printf(format, a...)
	}
}

func Fatal(f string, a ...any) {
	format := colors[LoggerLevelDebug] + "Fatal: " + Reset + f
	if level <= LoggerLevelFatal {
		fmt.Printf(format, a...)
	}
}

func Panic(f string, a ...any) {
	format := colors[LoggerLevelDebug] + "Panic: " + Reset + f
	if level <= LoggerLevelPanic {
		fmt.Printf(format, a...)
	}
}

func SetLevel(newLevel LoggerLevel) {
	level = newLevel
}

func GetLevel() LoggerLevel {
	return level
}
